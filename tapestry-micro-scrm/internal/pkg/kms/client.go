package kms

import "context"

// DataKeyClient 提供信封加密所需的数据密钥生成和解密能力。
type DataKeyClient interface {
	GenerateDataKey(context.Context) ([]byte, string, error)
	DecryptDataKey(context.Context, string) ([]byte, error)
}
