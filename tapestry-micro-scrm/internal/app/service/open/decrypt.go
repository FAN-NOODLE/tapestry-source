package openSvc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errs"
	pkgKMS "micro-scrm/internal/pkg/kms"
)

const maxDecryptItems = 100

type DecryptService struct {
	client  pkgKMS.DataKeyClient
	initErr error
}

type cachedDataKey struct {
	plaintext []byte
	err       error
}

func NewDecryptService() *DecryptService {
	client, err := pkgKMS.NewTencentClient(config.CfgData.KMS)
	return &DecryptService{client: client, initErr: err}
}

// NewDecryptServiceWithClient 注入数据密钥客户端，供测试和内部组合使用。
func NewDecryptServiceWithClient(client pkgKMS.DataKeyClient) *DecryptService {
	if client == nil {
		return &DecryptService{initErr: errors.New("KMS数据密钥客户端不能为空")}
	}
	return &DecryptService{client: client}
}

func (s *DecryptService) Decrypt(ctx context.Context, req openReq.DecryptReq) (openResp.DecryptResp, error) {
	if s.initErr != nil {
		return openResp.DecryptResp{}, fmt.Errorf("初始化KMS客户端失败: %w", s.initErr)
	}
	if err := validateDecryptRequest(req); err != nil {
		return openResp.DecryptResp{}, err
	}

	cache := make(map[string]cachedDataKey)
	defer clearDataKeyCache(cache)
	response := openResp.DecryptResp{Items: make([]openResp.DecryptResult, 0, len(req.Items))}
	for _, item := range req.Items {
		result := openResp.DecryptResult{RequestKey: item.RequestKey}
		value := encryptedValueFromRequest(item)
		if isEmptyDecryptItem(item) {
			result.Success = true
			response.Items = append(response.Items, result)
			continue
		}

		cached, exists := cache[item.EncryptedDataKey]
		if !exists {
			cached.plaintext, cached.err = s.client.DecryptDataKey(ctx, item.EncryptedDataKey)
			cache[item.EncryptedDataKey] = cached
		}
		if cached.err != nil {
			result.ErrorMessage = "解密失败"
			response.Items = append(response.Items, result)
			continue
		}

		plaintext, err := pkgKMS.DecryptValueWithDataKey(cached.plaintext, value)
		if err != nil {
			result.ErrorMessage = "解密失败"
			response.Items = append(response.Items, result)
			continue
		}
		result.Plaintext = plaintext
		result.Success = true
		response.Items = append(response.Items, result)
	}
	return response, nil
}

func validateDecryptRequest(req openReq.DecryptReq) error {
	if len(req.Items) < 1 || len(req.Items) > maxDecryptItems {
		return errs.OpenAPIParamsErr
	}
	for _, item := range req.Items {
		if strings.TrimSpace(item.RequestKey) == "" || item.Algorithm != pkgKMS.AES256GCMAlgorithm {
			return errs.OpenAPIParamsErr
		}
		emptyFields := 0
		for _, value := range []string{item.Ciphertext, item.EncryptedDataKey, item.IV, item.AuthTag} {
			if value == "" {
				emptyFields++
			}
		}
		if emptyFields != 0 && emptyFields != 4 {
			return errs.OpenAPIParamsErr
		}
	}
	return nil
}

func encryptedValueFromRequest(item openReq.DecryptItem) pkgKMS.EncryptedValue {
	return pkgKMS.EncryptedValue{
		Ciphertext:       item.Ciphertext,
		EncryptedDataKey: item.EncryptedDataKey,
		IV:               item.IV,
		AuthTag:          item.AuthTag,
		Algorithm:        item.Algorithm,
	}
}

func isEmptyDecryptItem(item openReq.DecryptItem) bool {
	return item.Ciphertext == "" && item.EncryptedDataKey == "" && item.IV == "" && item.AuthTag == ""
}

func clearDataKeyCache(cache map[string]cachedDataKey) {
	for encryptedDataKey, cached := range cache {
		for index := range cached.plaintext {
			cached.plaintext[index] = 0
		}
		delete(cache, encryptedDataKey)
	}
}
