package openSvc

import (
	"context"
	"errors"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

	openReq "micro-scrm/internal/app/http/request/open"
	pkgKMS "micro-scrm/internal/pkg/kms"
)

type decryptFakeKMS struct {
	dataKey      []byte
	encryptedKey string
	decryptCalls int
}

func (f *decryptFakeKMS) GenerateDataKey(context.Context) ([]byte, string, error) {
	return append([]byte(nil), f.dataKey...), f.encryptedKey, nil
}

func (f *decryptFakeKMS) DecryptDataKey(context.Context, string) ([]byte, error) {
	f.decryptCalls++
	return append([]byte(nil), f.dataKey...), nil
}

func TestDecryptRejectsInvalidItemCount(t *testing.T) {
	service := NewDecryptServiceWithClient(newDecryptFakeKMS())
	assertOpenAPIParamsError(t, decryptError(service, openReq.DecryptReq{}))

	items := make([]openReq.DecryptItem, 101)
	for index := range items {
		items[index] = openReq.DecryptItem{RequestKey: "key", Algorithm: pkgKMS.AES256GCMAlgorithm}
	}
	assertOpenAPIParamsError(t, decryptError(service, openReq.DecryptReq{Items: items}))
}

func TestDecryptRejectsInvalidItemStructure(t *testing.T) {
	service := NewDecryptServiceWithClient(newDecryptFakeKMS())
	tests := []openReq.DecryptItem{
		{Algorithm: pkgKMS.AES256GCMAlgorithm},
		{RequestKey: "key", Algorithm: "AES-256-CBC"},
		{RequestKey: "key", Ciphertext: "ciphertext", Algorithm: pkgKMS.AES256GCMAlgorithm},
	}
	for index, item := range tests {
		err := decryptError(service, openReq.DecryptReq{Items: []openReq.DecryptItem{item}})
		if err == nil {
			t.Fatalf("第%d个非法请求应返回参数错误", index)
		}
		assertOpenAPIParamsError(t, err)
	}
}

func decryptError(service *DecryptService, req openReq.DecryptReq) error {
	_, err := service.Decrypt(context.Background(), req)
	return err
}

func assertOpenAPIParamsError(t *testing.T, err error) {
	t.Helper()
	var businessErr *errorx.Error
	if !errors.As(err, &businessErr) {
		t.Fatalf("应返回业务参数错误，实际=%v", err)
	}
	if businessErr.Code != 400 {
		t.Fatalf("参数错误码应为400，实际=%d", businessErr.Code)
	}
}

func TestDecryptReusesDataKeyWithinRequest(t *testing.T) {
	client := newDecryptFakeKMS()
	first := encryptForDecryptTest(t, client, "763514")
	second := encryptForDecryptTestWithPage(t, client, "Susu Wang", &first)
	client.decryptCalls = 0

	service := NewDecryptServiceWithClient(client)
	response, err := service.Decrypt(context.Background(), openReq.DecryptReq{Items: []openReq.DecryptItem{
		decryptRequestItem("employee_code", first),
		decryptRequestItem("name", second),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if client.decryptCalls != 1 {
		t.Fatalf("相同encrypted_data_key应只调用一次KMS，实际%d次", client.decryptCalls)
	}
	if !response.Items[0].Success || response.Items[0].Plaintext != "763514" {
		t.Fatalf("第一项结果=%+v", response.Items[0])
	}
	if !response.Items[1].Success || response.Items[1].Plaintext != "Susu Wang" {
		t.Fatalf("第二项结果=%+v", response.Items[1])
	}
}

func TestDecryptContinuesAfterSingleItemFailure(t *testing.T) {
	client := newDecryptFakeKMS()
	valid := encryptForDecryptTest(t, client, "763514")
	invalid := valid
	invalid.Ciphertext = "不是Base64"
	client.decryptCalls = 0

	service := NewDecryptServiceWithClient(client)
	response, err := service.Decrypt(context.Background(), openReq.DecryptReq{Items: []openReq.DecryptItem{
		decryptRequestItem("invalid", invalid),
		decryptRequestItem("valid", valid),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if response.Items[0].Success || response.Items[0].Plaintext != "" || response.Items[0].ErrorMessage == "" {
		t.Fatalf("失败项结果=%+v", response.Items[0])
	}
	if !response.Items[1].Success || response.Items[1].Plaintext != "763514" {
		t.Fatalf("成功项结果=%+v", response.Items[1])
	}
}

func TestDecryptEmptyEncryptedValueDoesNotCallKMS(t *testing.T) {
	client := newDecryptFakeKMS()
	service := NewDecryptServiceWithClient(client)
	response, err := service.Decrypt(context.Background(), openReq.DecryptReq{Items: []openReq.DecryptItem{
		decryptRequestItem("empty", pkgKMS.EmptyEncryptedValue()),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !response.Items[0].Success || response.Items[0].Plaintext != "" {
		t.Fatalf("空对象结果=%+v", response.Items[0])
	}
	if client.decryptCalls != 0 {
		t.Fatalf("空对象不应调用KMS，实际%d次", client.decryptCalls)
	}
}

func newDecryptFakeKMS() *decryptFakeKMS {
	return &decryptFakeKMS{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "shared-encrypted-data-key",
	}
}

func encryptForDecryptTest(t *testing.T, client *decryptFakeKMS, plaintext string) pkgKMS.EncryptedValue {
	t.Helper()
	page, err := pkgKMS.NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer page.Close()
	value, err := page.EncryptString(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func encryptForDecryptTestWithPage(t *testing.T, client *decryptFakeKMS, plaintext string, first *pkgKMS.EncryptedValue) pkgKMS.EncryptedValue {
	t.Helper()
	page, err := pkgKMS.NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer page.Close()
	value, err := page.EncryptString(plaintext)
	if err != nil {
		t.Fatal(err)
	}
	value.EncryptedDataKey = first.EncryptedDataKey
	return value
}

func decryptRequestItem(requestKey string, value pkgKMS.EncryptedValue) openReq.DecryptItem {
	return openReq.DecryptItem{
		RequestKey:       requestKey,
		Ciphertext:       value.Ciphertext,
		EncryptedDataKey: value.EncryptedDataKey,
		IV:               value.IV,
		AuthTag:          value.AuthTag,
		Algorithm:        value.Algorithm,
	}
}
