package kms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	tencentCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentKMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"micro-scrm/internal/config"
)

const dataKeySize = 32

type tencentClient struct {
	client *tencentKMS.Client
	keyID  string
}

// NewTencentClient 创建腾讯云 KMS 客户端。
func NewTencentClient(cfg config.KMSConfig) (DataKeyClient, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	credential := tencentCommon.NewCredential(cfg.SecretID, cfg.SecretKey)
	clientProfile := profile.NewClientProfile()
	if cfg.Endpoint != "" {
		clientProfile.HttpProfile.Endpoint = cfg.Endpoint
	}
	client, err := tencentKMS.NewClient(credential, cfg.Region, clientProfile)
	if err != nil {
		return nil, fmt.Errorf("创建腾讯云KMS客户端失败: %w", err)
	}

	return &tencentClient{client: client, keyID: cfg.KeyID}, nil
}

func validateConfig(cfg config.KMSConfig) error {
	required := []struct {
		name  string
		value string
	}{
		{name: "region", value: cfg.Region},
		{name: "key_id", value: cfg.KeyID},
		{name: "secret_id", value: cfg.SecretID},
		{name: "secret_key", value: cfg.SecretKey},
	}
	for _, item := range required {
		if strings.TrimSpace(item.value) == "" {
			return fmt.Errorf("KMS配置缺少%s", item.name)
		}
	}
	return nil
}

func (c *tencentClient) GenerateDataKey(ctx context.Context) ([]byte, string, error) {
	request := tencentKMS.NewGenerateDataKeyRequest()
	request.KeyId = tencentCommon.StringPtr(c.keyID)
	request.NumberOfBytes = tencentCommon.Uint64Ptr(dataKeySize)

	response, err := c.client.GenerateDataKeyWithContext(ctx, request)
	if err != nil {
		return nil, "", fmt.Errorf("腾讯云KMS生成数据密钥失败: %w", err)
	}
	if response == nil || response.Response == nil || response.Response.Plaintext == nil || response.Response.CiphertextBlob == nil {
		return nil, "", errors.New("腾讯云KMS生成数据密钥响应不完整")
	}

	plaintext, err := decodePlaintextDataKey(*response.Response.Plaintext)
	if err != nil {
		return nil, "", err
	}
	if *response.Response.CiphertextBlob == "" {
		return nil, "", errors.New("腾讯云KMS返回的密文数据密钥为空")
	}
	return plaintext, *response.Response.CiphertextBlob, nil
}

func (c *tencentClient) DecryptDataKey(ctx context.Context, encryptedDataKey string) ([]byte, error) {
	if strings.TrimSpace(encryptedDataKey) == "" {
		return nil, errors.New("密文数据密钥不能为空")
	}

	request := tencentKMS.NewDecryptRequest()
	request.CiphertextBlob = tencentCommon.StringPtr(encryptedDataKey)
	response, err := c.client.DecryptWithContext(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("腾讯云KMS解密数据密钥失败: %w", err)
	}
	if response == nil || response.Response == nil || response.Response.Plaintext == nil {
		return nil, errors.New("腾讯云KMS解密数据密钥响应不完整")
	}
	return decodePlaintextDataKey(*response.Response.Plaintext)
}

func decodePlaintextDataKey(encoded string) ([]byte, error) {
	plaintext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("解析KMS明文数据密钥失败: %w", err)
	}
	if len(plaintext) != dataKeySize {
		return nil, fmt.Errorf("KMS明文数据密钥长度为%d，期望%d", len(plaintext), dataKeySize)
	}
	return plaintext, nil
}
