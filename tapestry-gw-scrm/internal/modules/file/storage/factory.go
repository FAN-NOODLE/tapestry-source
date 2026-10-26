package storage

import (
	"context"
	"gw-scrm/internal/modules/file/storage/blob"
	"gw-scrm/internal/modules/file/storage/cos"
)

type Storage interface {
	Upload(ctx context.Context, filePath string, fileName string, toPath string) (string, error)
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
