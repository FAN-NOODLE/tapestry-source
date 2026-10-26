package dl

import (
	"encoding/json"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

// 和Data Lake对接

// TIMEOUT_TIMESTAMP http请求常量
const (
	TIMEOUT_TIMESTAMP = 15 * time.Second
)

// DL 会员信息URL
const (
	MEMBER_INFO_QUERY_URL  = "/service/v1/data-lake/query-member-info"
	MEMBER_GRADE_QUERY_URL = "/service/v1/data-lake/query-member-grade"
	BA_INFO_URL            = "/service/v1/data-compass/customer-statistic/rank-for-console"
)

type BaseResponse struct {
	ID string `json:"id"`
	//Code         int    `json:"code"`
	Status       int    `json:"status"`
	Message      string `json:"msg"`
	Error        string `json:"error"`
	ConStartTime int64  `json:"con_start_time"`
	ConEndTime   int64  `json:"con_end_time"`
}

type BaseMemberQuery struct {
	QueryType string `json:"query_type"`
	Value     string `json:"value"`
}

type DataLake struct {
	prefixURL  string
	httpClient *sdk.DefaultClient
}

func NewDL() *DataLake {
	return &DataLake{
		prefixURL:  config.CfgData.Service.MicroWork.Domain,
		httpClient: sdk.NewSdkClient(clients.Timeout(TIMEOUT_TIMESTAMP)),
	}
}

func (DataLake) makeResponseData(resp *clients.ResponseProxy, data interface{}) error {
	resStr, rErr := resp.Read(clients.Text)
	if rErr != nil {
		log.Error(ctx.New()).Msgf("makeResponseData data %s,err %+v", resStr.(string), rErr)
		return rErr
	}

	return json.Unmarshal([]byte(resStr.(string)), data)
}
