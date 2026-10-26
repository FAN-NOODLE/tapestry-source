package http

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"gw-scrm/internal/config"
)

func TestSignAutoDownloadPathReturnsPrivateCOSPresignedURL(t *testing.T) {
	originalConfig := config.CfgData
	t.Cleanup(func() {
		config.CfgData = originalConfig
	})

	config.CfgData = &config.Config{
		Storage: config.StorageCfg{
			COS: config.CloudStorage{
				Schema:        "https",
				SecretID:      "REDACTED_FOR_SOURCE_SHARING",
				SecretKey:     "REDACTED_FOR_SOURCE_SHARING",
				Bucket:        "tapestry-uat-1308608494",
				PrivateBucket: "tapestry-uat-private-1308608494",
				Origin:        "ap-shanghai",
				CdnHost:       "cdn.example.com",
			},
		},
	}

	got, err := signAutoDownloadPath(
		context.Background(),
		"https://tapestry-uat-private-1308608494.cos.ap-shanghai.myqcloud.com/download/report.xlsx",
	)
	if err != nil {
		t.Fatalf("signAutoDownloadPath() error = %v", err)
	}

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	if parsed.Host != "tapestry-uat-private-1308608494.cos.ap-shanghai.myqcloud.com" {
		t.Fatalf("host = %q", parsed.Host)
	}
	if parsed.Path != "/download/report.xlsx" {
		t.Fatalf("path = %q", parsed.Path)
	}
	if !strings.Contains(parsed.RawQuery, "q-sign-algorithm=") {
		t.Fatalf("signed query missing: %q", parsed.RawQuery)
	}
}
