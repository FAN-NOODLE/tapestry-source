package kms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const AES256GCMAlgorithm = "AES-256-GCM"

// EncryptedValue 是 OpenAPI 返回的 PI 字段加密对象。
type EncryptedValue struct {
	Ciphertext       string `json:"ciphertext"`
	EncryptedDataKey string `json:"encrypted_data_key"`
	IV               string `json:"iv"`
	AuthTag          string `json:"auth_tag"`
	Algorithm        string `json:"algorithm"`
}

// PageCipher 在一个列表响应页面内复用同一数据密钥。
type PageCipher interface {
	EncryptString(string) (EncryptedValue, error)
	Close()
}

type pageCipher struct {
	ctx              context.Context
	client           DataKeyClient
	plaintextDataKey []byte
	encryptedDataKey string
}

// NewPageCipher 创建按页面复用数据密钥的加密器，数据密钥在首次加密非空字段时生成。
func NewPageCipher(ctx context.Context, client DataKeyClient) (PageCipher, error) {
	if client == nil {
		return nil, errors.New("KMS数据密钥客户端不能为空")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &pageCipher{ctx: ctx, client: client}, nil
}

// EmptyEncryptedValue 返回文档约定的空 PI 加密对象。
func EmptyEncryptedValue() EncryptedValue {
	return EncryptedValue{Algorithm: AES256GCMAlgorithm}
}

func (c *pageCipher) EncryptString(plaintext string) (EncryptedValue, error) {
	if plaintext == "" {
		return EmptyEncryptedValue(), nil
	}
	if err := c.ensureDataKey(); err != nil {
		return EncryptedValue{}, err
	}

	block, err := aes.NewCipher(c.plaintextDataKey)
	if err != nil {
		return EncryptedValue{}, fmt.Errorf("创建AES-256加密器失败: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptedValue{}, fmt.Errorf("创建AES-GCM加密器失败: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptedValue{}, fmt.Errorf("生成AES-GCM随机IV失败: %w", err)
	}

	sealed := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	tagStart := len(sealed) - gcm.Overhead()
	return EncryptedValue{
		Ciphertext:       base64.StdEncoding.EncodeToString(sealed[:tagStart]),
		EncryptedDataKey: c.encryptedDataKey,
		IV:               base64.StdEncoding.EncodeToString(nonce),
		AuthTag:          base64.StdEncoding.EncodeToString(sealed[tagStart:]),
		Algorithm:        AES256GCMAlgorithm,
	}, nil
}

func (c *pageCipher) ensureDataKey() error {
	if len(c.plaintextDataKey) > 0 {
		return nil
	}
	plaintextDataKey, encryptedDataKey, err := c.client.GenerateDataKey(c.ctx)
	if err != nil {
		return fmt.Errorf("获取KMS数据密钥失败: %w", err)
	}
	if len(plaintextDataKey) != dataKeySize {
		zeroBytes(plaintextDataKey)
		return fmt.Errorf("KMS明文数据密钥长度为%d，期望%d", len(plaintextDataKey), dataKeySize)
	}
	if encryptedDataKey == "" {
		zeroBytes(plaintextDataKey)
		return errors.New("KMS密文数据密钥为空")
	}
	c.plaintextDataKey = plaintextDataKey
	c.encryptedDataKey = encryptedDataKey
	return nil
}

func (c *pageCipher) Close() {
	zeroBytes(c.plaintextDataKey)
	c.plaintextDataKey = nil
	c.encryptedDataKey = ""
}

// DecryptValue 先通过 KMS 解密数据密钥，再解密单个 PI 字段。
func DecryptValue(ctx context.Context, client DataKeyClient, value EncryptedValue) (string, error) {
	if isEmptyEncryptedValue(value) {
		return validateAndDecrypt(nil, value)
	}
	if client == nil {
		return "", errors.New("KMS数据密钥客户端不能为空")
	}
	plaintextDataKey, err := client.DecryptDataKey(ctx, value.EncryptedDataKey)
	if err != nil {
		return "", fmt.Errorf("解密KMS数据密钥失败: %w", err)
	}
	defer zeroBytes(plaintextDataKey)
	return DecryptValueWithDataKey(plaintextDataKey, value)
}

// DecryptValueWithDataKey 使用已解密的数据密钥解密字段，供批量接口复用数据密钥。
func DecryptValueWithDataKey(plaintextDataKey []byte, value EncryptedValue) (string, error) {
	return validateAndDecrypt(plaintextDataKey, value)
}

func validateAndDecrypt(plaintextDataKey []byte, value EncryptedValue) (string, error) {
	if value.Algorithm != AES256GCMAlgorithm {
		return "", fmt.Errorf("不支持的加密算法: %s", value.Algorithm)
	}
	if isEmptyEncryptedValue(value) {
		return "", nil
	}
	if value.Ciphertext == "" || value.EncryptedDataKey == "" || value.IV == "" || value.AuthTag == "" {
		return "", errors.New("加密对象字段不完整")
	}
	if len(plaintextDataKey) != dataKeySize {
		return "", fmt.Errorf("明文数据密钥长度为%d，期望%d", len(plaintextDataKey), dataKeySize)
	}

	ciphertext, err := decodeBase64Field("ciphertext", value.Ciphertext)
	if err != nil {
		return "", err
	}
	nonce, err := decodeBase64Field("iv", value.IV)
	if err != nil {
		return "", err
	}
	authTag, err := decodeBase64Field("auth_tag", value.AuthTag)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(plaintextDataKey)
	if err != nil {
		return "", fmt.Errorf("创建AES-256解密器失败: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建AES-GCM解密器失败: %w", err)
	}
	if len(nonce) != gcm.NonceSize() {
		return "", fmt.Errorf("AES-GCM IV长度为%d，期望%d", len(nonce), gcm.NonceSize())
	}
	if len(authTag) != gcm.Overhead() {
		return "", fmt.Errorf("AES-GCM认证标签长度为%d，期望%d", len(authTag), gcm.Overhead())
	}

	sealed := make([]byte, 0, len(ciphertext)+len(authTag))
	sealed = append(sealed, ciphertext...)
	sealed = append(sealed, authTag...)
	plaintext, err := gcm.Open(nil, nonce, sealed, nil)
	if err != nil {
		return "", errors.New("AES-GCM认证或解密失败")
	}
	return string(plaintext), nil
}

func isEmptyEncryptedValue(value EncryptedValue) bool {
	return value.Ciphertext == "" && value.EncryptedDataKey == "" && value.IV == "" && value.AuthTag == ""
}

func decodeBase64Field(name, value string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("%s不是有效Base64: %w", name, err)
	}
	return decoded, nil
}

func zeroBytes(value []byte) {
	for index := range value {
		value[index] = 0
	}
}
