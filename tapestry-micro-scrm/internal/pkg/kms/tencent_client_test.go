package kms

import (
	"encoding/base64"
	"testing"

	"micro-scrm/internal/config"
)

func TestNewTencentClientRejectsMissingRequiredConfig(t *testing.T) {
	valid := config.KMSConfig{
		Region:    "ap-shanghai",
		KeyID:     "9fb110c3-9aab-11f1-a81d-5254001f6749",
		SecretID:  "REDACTED_FOR_SOURCE_SHARING",
		SecretKey: "REDACTED_FOR_SOURCE_SHARING",
	}

	tests := []struct {
		name   string
		mutate func(*config.KMSConfig)
	}{
		{name: "缺少Region", mutate: func(cfg *config.KMSConfig) { cfg.Region = "" }},
		{name: "缺少KeyId", mutate: func(cfg *config.KMSConfig) { cfg.KeyID = "" }},
		{name: "缺少SecretId", mutate: func(cfg *config.KMSConfig) { cfg.SecretID = "" }},
		{name: "缺少SecretKey", mutate: func(cfg *config.KMSConfig) { cfg.SecretKey = "" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := valid
			test.mutate(&cfg)
			if _, err := NewTencentClient(cfg); err == nil {
				t.Fatal("缺少必填配置时应返回错误")
			}
		})
	}
}

func TestDecodePlaintextDataKeyRequires32Bytes(t *testing.T) {
	invalid := base64.StdEncoding.EncodeToString(make([]byte, 31))
	if _, err := decodePlaintextDataKey(invalid); err == nil {
		t.Fatal("KMS明文数据密钥不是32字节时应返回错误")
	}

	valid := base64.StdEncoding.EncodeToString(make([]byte, 32))
	key, err := decodePlaintextDataKey(valid)
	if err != nil {
		t.Fatalf("解析32字节数据密钥失败: %v", err)
	}
	if len(key) != 32 {
		t.Fatalf("数据密钥长度=%d，期望32", len(key))
	}
}
