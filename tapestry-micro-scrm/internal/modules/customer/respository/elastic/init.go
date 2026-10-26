package elastic

import (
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/elastic"
)

func NewEsRepo() *EsRepo {
	return &EsRepo{cli: elastic.ShareEsClient, env: config.CfgData.Env}
}
