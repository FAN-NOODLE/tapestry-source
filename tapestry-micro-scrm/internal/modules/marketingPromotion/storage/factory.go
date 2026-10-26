package storage

import (
	"context"
	"micro-scrm/internal/modules/marketingPromotion/storage/blob"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
)

type Storage interface {
	Upload(ctx context.Context, filePath string, toFilePath string) (string, error)
	Download(ctx context.Context, fromFilePath string, toFilePath string) error
	GetDomain() string
	GetURL(objectKey string) string
}

func Init(pkg string) Storage {
	switch pkg {
	case cos.COS:
		return new(cos.Storage)
	case blob.BLOB:
		return new(blob.Storage)
	}
	return nil
}
