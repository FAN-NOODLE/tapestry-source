package cosroute

import (
	"net/url"
	"strings"

	"gw-scrm/internal/config"
)

var publicPrefixes = []string{
	"local-upload/",
	"static/",
	"uploads/",
}

type ObjectInfo struct {
	Bucket string
	Key    string
}

func NormalizeKey(key string) string {
	return strings.TrimLeft(key, "/")
}

func IsPublicKey(key string) bool {
	key = NormalizeKey(key)
	for _, prefix := range publicPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

func PrivateBucket() string {
	if config.CfgData.Storage.COS.PrivateBucket != "" {
		return config.CfgData.Storage.COS.PrivateBucket
	}
	if config.CfgData.Env == "prod" {
		return "tapestry-prod-private-1308608494"
	}
	return "tapestry-uat-private-1308608494"
}

func BucketForKey(key string) string {
	if IsPublicKey(key) {
		return config.CfgData.Storage.COS.Bucket
	}
	return PrivateBucket()
}

func OriginHost(bucket string) string {
	return bucket + ".cos." + config.CfgData.Storage.COS.Origin + ".myqcloud.com"
}

func PublicOriginHost() string {
	return OriginHost(config.CfgData.Storage.COS.Bucket)
}

func PrivateOriginHost() string {
	return OriginHost(PrivateBucket())
}

func URLForKey(key string) string {
	key = NormalizeKey(key)
	if IsPublicKey(key) {
		return config.CfgData.Storage.COS.Schema + "://" + strings.TrimRight(config.CfgData.Storage.COS.CdnHost, "/") + "/" + key
	}
	return config.CfgData.Storage.COS.Schema + "://" + PrivateOriginHost() + "/" + key
}

func ParseObjectURL(fileURL string) (ObjectInfo, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return ObjectInfo{}, err
	}

	key := NormalizeKey(parsedURL.Path)
	host := strings.ToLower(parsedURL.Host)

	switch host {
	case strings.ToLower(config.CfgData.Storage.COS.CdnHost), strings.ToLower(PublicOriginHost()):
		return ObjectInfo{Bucket: config.CfgData.Storage.COS.Bucket, Key: key}, nil
	case strings.ToLower(PrivateOriginHost()):
		return ObjectInfo{Bucket: PrivateBucket(), Key: key}, nil
	default:
		return ObjectInfo{}, errUnsupportedHost
	}
}

var errUnsupportedHost = &unsupportedHostError{}

type unsupportedHostError struct{}

func (e *unsupportedHostError) Error() string {
	return "仅支持当前 COS 文件地址"
}
