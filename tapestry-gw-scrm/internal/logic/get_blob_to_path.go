package logic

import "gw-scrm/internal/config"

const (
	PROP_BLOB_PATH = `/tapestry-prod`
	TEST_BLOB_PATH = `/tapestry-uat`
)

func GetBolbPrefix() string {
	if config.CfgData.Env == "prop" {
		return PROP_BLOB_PATH
	}
	return TEST_BLOB_PATH
}
