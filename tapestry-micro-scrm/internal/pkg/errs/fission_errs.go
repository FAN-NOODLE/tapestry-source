package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	FissionActivityNotExistErr     = errorx.NewErrorMsg(200001, "裂变活动不存在")
	FissionPrizeIDEmpty            = errorx.NewErrorMsg(200002, "裂变奖品id为空")
	FissionPrizeNotExist           = errorx.NewErrorMsg(200003, "裂变奖品不存在")
	FissionUpdateErr               = errorx.NewErrorMsg(200004, "更新状态与活动状态不匹配")
	FissionDeleteErr               = errorx.NewErrorMsg(200005, "进行中的活动不可删除")
	FissionActivityNotStart        = errorx.NewErrorMsg(200006, "活动未开始")
	FissionPrizeExchangeOut        = errorx.NewErrorMsg(200007, "抱歉,奖品已被抢光")
	FissionPrizeOff                = errorx.NewErrorMsg(200008, "抱歉,奖品已被下架")
	FissionPrizeExchangeEnd        = errorx.NewErrorMsg(200009, "抱歉，兑换活动已结束。谢谢你的关注")
	FissionJoinParamErr            = errorx.NewErrorMsg(200010, "参数错误")
	FissionActivityOff             = errorx.NewErrorMsg(200011, "活动已下架")
	FissionActivityEnd             = errorx.NewErrorMsg(200012, "活动已结束")
	FissionUserJoinRepeat          = errorx.NewErrorMsg(200013, "用户已参加当前活动")
	FissionUserJoinMutual          = errorx.NewErrorMsg(200014, "抱歉，你已参与过类似活动 无法参加本次活动")
	FissionJoinFailedNoRelation    = errorx.NewErrorMsg(200015, "参加活动失败：未添加企业好友")
	FissionQueryRelationParamErr   = errorx.NewErrorMsg(200016, "查询外部联系人关系参数错误")
	FissionActivityNotExist        = errorx.NewErrorMsg(200017, "裂变活动不存在")
	FissionExchangeAssistNotEnough = errorx.NewErrorMsg(200018, "助力人数不够，兑换失败")
	FissionActivityUserNotExist    = errorx.NewErrorMsg(200019, "未找到用户活动数据")
	FissionUserNotLogin            = errorx.NewErrorMsg(200020, "用户未登录")
	FissionUserExchangeLimit       = errorx.NewErrorMsg(FissionUserExchangeLimitCode, "已达到该奖品兑换上限")
	ExchangeEndTimeErr             = errorx.NewErrorMsg(200022, "奖品兑换结束时间不得小于活动结束时间")
	ExchangeStartTimeErr           = errorx.NewErrorMsg(200023, "活动开始时间不得小于当前时间")
	GetWxExternalContactErr        = errorx.NewErrorMsg(200024, "裂变员工有误，操作失败")
)

var (
	FissionExchangeAssistNotEnoughCode = 200018
	FissionExchangeAssistNotEnoughMsg  = "还差%d人助力就可以免费兑换这个奖品啦"
	FissionUserExchangeLimitCode       = 200021
	FissionUserExchangeLimitMsg        = "抱歉，该奖品只允许兑换%d次"
)
