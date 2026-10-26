package appJob

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
)

var errCosPrivateLegacyJobDisabled = errors.New("legacy COS private migration job is disabled; use cos_private_db_migrate")

var cosPrivatePublicPrefixes = []string{
	"local-upload/",
	"static/",
	"uploads/",
}

type CosPrivateCopy struct{}

type CosPrivateDelete struct{}

type CosPrivateDbURLRefresh struct{}

func NewCosPrivateCopy() *CosPrivateCopy {
	return &CosPrivateCopy{}
}

func NewCosPrivateDelete() *CosPrivateDelete {
	return &CosPrivateDelete{}
}

func NewCosPrivateDbURLRefresh() *CosPrivateDbURLRefresh {
	return &CosPrivateDbURLRefresh{}
}

func (job *CosPrivateCopy) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	return errCosPrivateLegacyJobDisabled
}

func (job *CosPrivateDelete) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	return errCosPrivateLegacyJobDisabled
}

func (job *CosPrivateDbURLRefresh) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	return errCosPrivateLegacyJobDisabled
}

func newCosClient(bucket string) *cos.Client {
	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosOriginHost(bucket))
	b := &cos.BaseURL{BucketURL: u}
	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})
}

func cosPrivateBucket() string {
	if config.CfgData.Storage.COS.PrivateBucket != "" {
		return config.CfgData.Storage.COS.PrivateBucket
	}
	return helper.RuleThree(config.CfgData.Env != "prod", "tapestry-uat-private-1308608494", "tapestry-prod-private-1308608494")
}

func cosOriginHost(bucket string) string {
	return bucket + ".cos." + config.CfgData.Storage.COS.Origin + ".myqcloud.com"
}

func cosPrivateIsPublicKey(key string) bool {
	key = strings.TrimLeft(key, "/")
	for _, prefix := range cosPrivatePublicPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
