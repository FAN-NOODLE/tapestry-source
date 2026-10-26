package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/sdk"
	"time"
)

// 用户会员类型相关业务

type vipUserCommon struct {
	ctx        ctx.Context
	name       string
	httpClient *sdk.DefaultClient
	param      VipUserParam
}

const (
	TIMEOUT_TIMESTAMP = 15 * time.Second // 请求超时时间戳
	JSON              = "application/json"
)

type VipUserParam struct {
	EnterpriseID   int64
	UnionID        string
	ExternalUserID string
}

func NewVipUser(ctx ctx.Context) *vipUserCommon {
	return &vipUserCommon{
		ctx:        ctx,
		name:       "VIP会员Logic",
		httpClient: sdk.NewSdkClient(clients.Timeout(TIMEOUT_TIMESTAMP)),
	}
}

// 设置校验的用户
func (v *vipUserCommon) SetUser(param VipUserParam) *vipUserCommon {
	if param.UnionID == "" {
		info, _ := orm.Service.WorkExternalUser.GetWorkExternalUserInfoByWorkExternalUserID(param.EnterpriseID, param.ExternalUserID)
		log.Info(v.ctx).Msgf("会员逻辑::用户信息：%v", info)
		v.param.UnionID = info.Unionid
	}

	v.param = param
	return v
}

// 判断是否是否为VIP
func (v vipUserCommon) IsVip() (bool, error) {
	if v.param.UnionID == "" {
		return false, nil
	}

	info, err := orm.Service.WorkExternalUser.GetWxUserCorpExternalInfoByUnionID(v.param.EnterpriseID, v.param.UnionID)
	log.Info(v.ctx).Msgf("会员逻辑::本地会员信息：%v", info)
	if err != nil {
		return false, err
	}

	if info.UnionCode != "" {
		return true, nil
	}

	// 请求三方
	return v.IsVipThird()
}

func (v vipUserCommon) IsVipThird() (bool, error) {
	vipInfo, err := v.VipThirdInfo()
	log.Info(v.ctx).Msgf("会员逻辑::接口会员信息：%v", vipInfo)
	if err != nil {
		return false, err
	}

	if vipInfo.Status == 200 {
		return true, nil
	}

	return false, nil
}

type VipThirdResponse struct {
	ID      string      `json:"id"`
	Code    interface{} `json:"code"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    struct {
		Code           int    `json:"code"`
		Msg            string `json:"msg"`
		IsUpdateGrade  string `json:"isUpdateGrade"`
		Grade          string `json:"grade"`
		GradeEndDate   string `json:"gradeEndDate"`
		GradeCode      string `json:"gradeCode"`
		GradeStartDate string `json:"gradeStartDate"`
	} `json:"data"`
}

func (v vipUserCommon) VipThirdInfo() (data VipThirdResponse, err error) {
	vipInfoThirdURL := config.CfgData.Service.MicroWork.Domain + "/service/v1/data-lake/query-member-grade"
	params := gin.H{
		"enterprise_id": v.param.EnterpriseID,
		"type":          "unionId",
		"value":         v.param.UnionID,
	}

	resp, err := v.httpClient.Post(v.ctx, vipInfoThirdURL, clients.JSONBody(params))
	if err != nil {
		return
	}
	resStr, rErr := resp.Read(clients.Text)
	if rErr != nil {
		return
	}

	err = json.Unmarshal([]byte(resStr.(string)), &data)
	if err != nil {
		return
	}

	return
}
