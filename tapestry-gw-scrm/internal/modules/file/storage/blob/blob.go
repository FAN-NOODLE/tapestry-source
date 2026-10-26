package blob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"gw-scrm/internal/config"
	"net/url"
	"os"
)

const BLOB = "blob"

type Storage struct{}

// 上传微软云
func (cs *Storage) Upload(ctx context.Context, filePath string, fileName string, toPath string) (string, error) {
	toPath += "/"
	accountName := config.CfgData.Storage.Blob.AccountName
	accountKey := config.CfgData.Storage.Blob.AccountKey
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	bucket := config.CfgData.Storage.Blob.Bucket
	rawUrl := fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s%s", accountName, bucket, toPath)
	URL, _ := url.Parse(rawUrl)
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	// 创建bucket
	//_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessBlob)
	//if err != nil {
	//	fmt.Println("11111111111111111111111111",err.Error())
	//	return "", err
	//}

	blobURL := containerURL.NewBlockBlobURL(fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	return toPath + fileName, nil
}

// 获取存储域
func (cs *Storage) GetDomain() string {
	return config.CfgData.Storage.Blob.Domain
}

func (cs *Storage) GetURL(objectKey string) string {
	return cs.GetDomain() + objectKey
}
