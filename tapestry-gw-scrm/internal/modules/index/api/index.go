package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/index"
	"gw-scrm/internal/pkg/sdk/cdp"
	"gw-scrm/internal/pkg/sdk/microauthority"
	"gw-scrm/internal/pkg/sdk/microboss"
	"gw-scrm/internal/pkg/sdk/microconsole"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/sdk/microshop"
)

type IndexApi struct {
}

// 商城变现
func (m IndexApi) ShopFee(c ctx.Context, req index.ShopFeeParam, body index.ShopFeeBody) (interface{}, error) {
	res, err := cdp.CompositeIndicator(c, req, body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 公众号粉丝
func (m IndexApi) PublicSubFans(c ctx.Context, req index.ShopFeeParam, cdpSign index.FansBody) (interface{}, error) {
	res, err := cdp.FansTargetQuery(c, req, cdpSign)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 导购排行榜
func (m IndexApi) ShopperList(c ctx.Context, req index.ShopperListParam, cdpSign index.ShopperListBody) (interface{}, error) {
	res, err := cdp.ShopperTargetQuery(c, req, cdpSign)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 公告
func (m IndexApi) Notice(c ctx.Context, req map[string]interface{}) (*microboss.AnnouncementRes, error) {
	res, err := microboss.AnnouncementList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 门店列表
func (m IndexApi) StoreList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microshop.StoreList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 基本配置校验
func (m IndexApi) BaseConfigCheck(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microauthority.BaseInfoCheck(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 企微好友数
func (m IndexApi) QwFriend(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.WorkRelationCalc(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 群成员数
func (m IndexApi) GroupMemberNum(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microscrm.GroupRelationCalc(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 帮助中心
func (m IndexApi) HelpList(c ctx.Context, req interface{}) (*index.HelpListRes, error) {
	res, err := microconsole.HelpList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m IndexApi) BuyNotice(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microboss.BuyNotice(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
