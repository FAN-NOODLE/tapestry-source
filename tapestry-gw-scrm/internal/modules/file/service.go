package file

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	qcos "github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/file/storage"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/cosroute"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type File struct {
	Path   string `json:"path"`
	Domain string `json:"domain"`
	Url    string `json:"url"`
}

type CredentialInfo struct {
	Bucket        string `json:"bucket"`
	PrivateBucket string `json:"private_bucket"`
	Region        string `json:"region"`
	SecretID      string `json:"secret_id"`
	SecretKey     string `json:"secret_key"`
	SessionToken  string `json:"session_token"`
	StartTime     int    `json:"start_time"`
	ExpiredTime   int    `json:"expired_time"`
}

func UploadCloud(ctx gin.Context, pkg string, tmpFilePath string, rootDir string) (File, error) {
	var f = File{}
	cs := storage.Init(pkg)

	filePathArr := strings.Split(tmpFilePath, "/")
	fileName := filePathArr[len(filePathArr)-1]

	path, err := cs.Upload(&ctx, tmpFilePath, fileName, rootDir+"/"+common.GetDate("/"))
	if err != nil {
		return f, err
	}

	_ = os.Remove(tmpFilePath)

	f.Path = path
	f.Domain = cs.GetDomain()
	f.Url = cs.GetURL(path)

	return f, err
}

// Credential 获取临时凭证
func Credential() (CredentialInfo, error) {

	bucket := config.CfgData.Storage.COS.Bucket
	privateBucket := cosroute.PrivateBucket()
	region := config.CfgData.Storage.COS.Origin
	c := sts.NewClient(
		config.CfgData.Storage.COS.SecretID,
		config.CfgData.Storage.COS.SecretKey,
		nil,
	)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"*",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + region + ":uid/" + config.CfgData.Storage.COS.AppID + ":" + bucket + "/*",
						"qcs::cos:" + region + ":uid/" + config.CfgData.Storage.COS.AppID + ":" + privateBucket + "/*",
					},
				},
			},
		},
	}

	res, err := c.GetCredential(opt)
	if err != nil {
		panic(err)
	}

	credential := &CredentialInfo{
		Region:        region,
		Bucket:        bucket,
		PrivateBucket: privateBucket,
		SessionToken:  res.Credentials.SessionToken,
		SecretID:      res.Credentials.TmpSecretID,
		SecretKey:     res.Credentials.TmpSecretKey,
		StartTime:     res.StartTime,
		ExpiredTime:   res.ExpiredTime,
	}

	return *credential, nil
}

func GetPrivateFileURL(ctx context.Context, fileURL string) (string, error) {
	objectInfo, err := getCOSObjectInfo(fileURL)
	if err != nil {
		return "", err
	}

	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosroute.OriginHost(objectInfo.Bucket))
	b := &qcos.BaseURL{BucketURL: u}
	client := qcos.NewClient(b, &http.Client{
		Transport: &qcos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})

	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, objectInfo.Key, config.CfgData.Storage.COS.SecretID, config.CfgData.Storage.COS.SecretKey, time.Hour, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func getCOSObjectInfo(fileURL string) (cosroute.ObjectInfo, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return cosroute.ObjectInfo{}, err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return cosroute.ObjectInfo{}, errors.New("文件地址格式错误")
	}

	objectInfo, err := cosroute.ParseObjectURL(fileURL)
	if err != nil {
		return cosroute.ObjectInfo{}, err
	}
	if objectInfo.Key == "" {
		return cosroute.ObjectInfo{}, errors.New("文件地址路径为空")
	}

	return objectInfo, nil
}

//文本加密
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var key = "aaabbbcccxaaabbbcccxaaabbbcccx88"

func AesEncrypt(tex string) (string, error) {
	//需要去加密的字符串
	plaintext := []byte(tex)
	// 创建加密算法aes
	c, err := newCipher(key)
	if err != nil {
		return "error", err
	}
	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return hex.EncodeToString(ciphertext), nil
}

func AesDecrypt(enText string) (string, error) {
	enByte, err := hex.DecodeString(enText)
	if err != nil {
		return "error", err
	}
	// 创建加密算法aes
	c, cipherErr := newCipher(key)
	if cipherErr != nil {
		return "error", cipherErr
	}
	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(enByte))
	cfbdec.XORKeyStream(plaintextCopy, enByte)
	return string(plaintextCopy), nil
}

//创建加密算法aes
func newCipher(key string) (c cipher.Block, err error) {
	// 创建加密算法aes
	c, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	return
}
