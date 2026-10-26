package gomicrouser

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	ConfigGuideListUrl   string `desc:"配置引导-列表"`
	LevelDirectListUrl   string `desc:"等级直达-列表"`
	LevelDirectAddUrl    string `desc:"等级直达-添加"`
	LevelDirectUpdateUrl string `desc:"等级直达-更新"`
	LevelDirectDeleteUrl string `desc:"等级直达-删除"`
	MemberLevelListUrl   string `desc:"会员等级-列表"`
	MemberLevelDetailUrl string `desc:"会员等级-详情"`
	MemberLevelSetUrl    string `desc:"会员等级-设置"`
	UpdateStatusUrl      string `desc:"会员等级-状态修改"`
	ContentSortUrl       string `desc:"会员等级-内容排序"`
	LevelUpListURL       string `desc:"等级跃迁-列表"`
	AddLevelUpURL        string `desc:"等级跃迁-新增"`
	EditLevelUpURL       string `desc:"等级跃迁-修改"`
	DeleteLevelUpURL     string `desc:"等级跃迁-删除"`
	LevelUpInfoURL       string `desc:"等级跃迁-详情"`
	SaveCouponRuleUrl    string `desc:"添加编辑优惠券规则"`
	WxCouponListUrl      string `desc:"优惠券列表"`
	InfoCouponRuleUrl    string `desc:"优惠券权益规则详情"`
	CouponPackageListUrl string `desc:"券包列表"`
	RuleListUrl          string `desc:"券包列表"`
	ScoreRuleSave        string `desc:"积分权益规则新增保存"`
	ScoreDetail          string `desc:"积分权益规则详情"`
	SetUse               string `desc:"积分规则启用禁用"`
	RuleDelete           string `desc:"积分规则删除"`
	RightSave            string `desc:"权益保存"`
	RightList            string `desc:"权益列表"`
	SelectLevelListUrl   string `desc:"会员等级-下拉框等级列表"`
	UpdateRightStatusUrl string `desc:"更新权益状态"`
	LevelDirectDetailUrl string `desc:"等级直达详情"`
}{
	// 配置引导-获取
	ConfigGuideListUrl:   "%s/micro-user/v1/config-guide/list",
	LevelDirectListUrl:   "%s/micro-user/v1/level-direct/list",
	LevelDirectAddUrl:    "%s/micro-user/v1/level-direct/add",
	LevelDirectUpdateUrl: "%s/micro-user/v1/level-direct/update",
	LevelDirectDeleteUrl: "%s/micro-user/v1/level-direct/delete",
	MemberLevelListUrl:   "%s/micro-user/v1/member-level/list",
	MemberLevelSetUrl:    "%s/micro-user/v1/member-level/set",
	MemberLevelDetailUrl: "%s/micro-user/v1/member-level/detail",
	UpdateStatusUrl:      "%s/micro-user/v1/member-level/update-status",
	ContentSortUrl:       "%s/micro-user/v1/member-level/content-sort",
	LevelUpListURL:       "%s/micro-user/v1/member/level-up-list",
	AddLevelUpURL:        "%s/micro-user/v1/member/level-up-add",
	EditLevelUpURL:       "%s/micro-user/v1/member/level-up-edit",
	DeleteLevelUpURL:     "%s/micro-user/v1/member/level-up-delete",
	LevelUpInfoURL:       "%s/micro-user/v1/member/level-up-info",
	SaveCouponRuleUrl:    "%s/micro-user/v1/right-rule/save-coupon-rule",
	WxCouponListUrl:      "%s/micro-user/v1/right-rule/wx-coupon-list",
	InfoCouponRuleUrl:    "%s/micro-user/v1/right-rule/info-coupon-rule",
	CouponPackageListUrl: "%s/micro-user/v1/right-rule/coupon-package-list",
	RuleListUrl:          "%s/micro-user/v1/right-rule/list",
	ScoreRuleSave:        "%s/micro-user/v1/right-rule/score-rule-save",
	ScoreDetail:          "%s/micro-user/v1/right-rule/score-detail",
	SetUse:               "%s/micro-user/v1/right-rule/set-use",
	RuleDelete:           "%s/micro-user/v1/right-rule/delete",
	RightSave:            "%s/micro-user/v1/legal-right/save",
	RightList:            "%s/micro-user/v1/legal-right/list",
	SelectLevelListUrl:   "%s/micro-user/v1/member-level/level-list",
	UpdateRightStatusUrl: "%s/micro-user/v1/legal-right/update-status",
	LevelDirectDetailUrl: "%s/micro-user/v1/level-direct/detail",
}

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.GoMicroUser.Domain)
}

var httpClient = sdk.NewHttpClient(sdk.ParseJson{})

type GoUserResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func (r *GoUserResponse) CheckCode() error {
	if r.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(r.Code, r.Msg)

}

// CommonResponse 公共默认返回值
type CommonResponse struct {
	GoUserResponse
	Data interface{} `json:"data"`
}

var scrmClient = sdk.NewHttpClient(sdk.ParseJson{})
