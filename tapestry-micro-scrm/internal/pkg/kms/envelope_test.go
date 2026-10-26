package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
)

type fakeDataKeyClient struct {
	dataKey       []byte
	encryptedKey  string
	generateCalls int
	decryptCalls  int
	generateErr   error
	decryptErr    error
}

func (f *fakeDataKeyClient) GenerateDataKey(context.Context) ([]byte, string, error) {
	f.generateCalls++
	if f.generateErr != nil {
		return nil, "", f.generateErr
	}
	return append([]byte(nil), f.dataKey...), f.encryptedKey, nil
}

func (f *fakeDataKeyClient) DecryptDataKey(context.Context, string) ([]byte, error) {
	f.decryptCalls++
	if f.decryptErr != nil {
		return nil, f.decryptErr
	}
	return append([]byte(nil), f.dataKey...), nil
}

func TestPageCipherSharesDataKeyAndUsesIndependentIV(t *testing.T) {
	client := newFakeDataKeyClient()
	cipher, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer cipher.Close()

	first, err := cipher.EncryptString("Jean Yuan")
	if err != nil {
		t.Fatal(err)
	}
	second, err := cipher.EncryptString("Sandy Tang")
	if err != nil {
		t.Fatal(err)
	}

	if client.generateCalls != 1 {
		t.Fatalf("GenerateDataKey调用次数=%d，期望1", client.generateCalls)
	}
	if first.EncryptedDataKey != second.EncryptedDataKey || first.EncryptedDataKey != client.encryptedKey {
		t.Fatal("同一页面的加密字段应共用encrypted_data_key")
	}
	if first.IV == second.IV {
		t.Fatal("不同字段必须使用不同IV")
	}
	for _, value := range []EncryptedValue{first, second} {
		if value.Algorithm != AES256GCMAlgorithm {
			t.Fatalf("algorithm=%q", value.Algorithm)
		}
		assertBase64(t, value.Ciphertext)
		assertBase64(t, value.IV)
		assertBase64(t, value.AuthTag)
	}
}

func TestPageCipherEmptyValueDoesNotCallKMS(t *testing.T) {
	client := newFakeDataKeyClient()
	cipher, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer cipher.Close()

	value, err := cipher.EncryptString("")
	if err != nil {
		t.Fatal(err)
	}
	if client.generateCalls != 0 {
		t.Fatalf("空值不应调用KMS，实际调用%d次", client.generateCalls)
	}
	if value != EmptyEncryptedValue() {
		t.Fatalf("空值加密对象=%+v", value)
	}
}

func TestEncryptAndDecryptValueRoundTrip(t *testing.T) {
	client := newFakeDataKeyClient()
	cipher, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	value, err := cipher.EncryptString("000931251")
	cipher.Close()
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := DecryptValue(context.Background(), client, value)
	if err != nil {
		t.Fatal(err)
	}
	if plaintext != "000931251" {
		t.Fatalf("解密结果=%q", plaintext)
	}
	if client.decryptCalls != 1 {
		t.Fatalf("DecryptDataKey调用次数=%d，期望1", client.decryptCalls)
	}
}

func TestDecryptValueRejectsTampering(t *testing.T) {
	client := newFakeDataKeyClient()
	cipher, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	original, err := cipher.EncryptString("Susu Wang")
	cipher.Close()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*EncryptedValue)
	}{
		{name: "密文", mutate: func(value *EncryptedValue) { value.Ciphertext = flipBase64Byte(t, value.Ciphertext) }},
		{name: "IV", mutate: func(value *EncryptedValue) { value.IV = flipBase64Byte(t, value.IV) }},
		{name: "认证标签", mutate: func(value *EncryptedValue) { value.AuthTag = flipBase64Byte(t, value.AuthTag) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := original
			test.mutate(&value)
			if _, err := DecryptValue(context.Background(), client, value); err == nil {
				t.Fatal("被篡改的加密对象应解密失败")
			}
		})
	}
}

func TestPageCipherReturnsKMSFailure(t *testing.T) {
	client := newFakeDataKeyClient()
	client.generateErr = errors.New("kms unavailable")
	cipher, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	defer cipher.Close()

	if _, err := cipher.EncryptString("Jean Yuan"); err == nil {
		t.Fatal("KMS失败时加密应返回错误")
	}
}

func TestPageCipherCloseClearsPlaintextDataKey(t *testing.T) {
	client := newFakeDataKeyClient()
	page, err := NewPageCipher(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	cipher := page.(*pageCipher)
	if _, err := cipher.EncryptString("Jean Yuan"); err != nil {
		t.Fatal(err)
	}
	keyReference := cipher.plaintextDataKey

	cipher.Close()

	for index, value := range keyReference {
		if value != 0 {
			t.Fatalf("Close后明文数据密钥第%d字节未清零", index)
		}
	}
}

func newFakeDataKeyClient() *fakeDataKeyClient {
	return &fakeDataKeyClient{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "kms-encrypted-data-key",
	}
}

func assertBase64(t *testing.T, value string) {
	t.Helper()
	if _, err := base64.StdEncoding.DecodeString(value); err != nil {
		t.Fatalf("%q不是有效Base64: %v", value, err)
	}
}

func flipBase64Byte(t *testing.T, value string) string {
	t.Helper()
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	decoded[0] ^= 0xff
	return base64.StdEncoding.EncodeToString(decoded)
}
