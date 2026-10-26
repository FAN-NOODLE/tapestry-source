package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/excitation"
)

func Excitation(group *gin.RouterGroup) {
	// 奖章
	medalGroup := group.Group("medal")
	{
		medalHandle := excitation.NewMedalHandle()
		medalGroup.POST("list", medalHandle.List)
		medalGroup.POST("create", medalHandle.Create)
		medalGroup.POST("detail", medalHandle.Detail)
		medalGroup.POST("update", medalHandle.Update)
		medalGroup.POST("remove", medalHandle.Remove)
		medalGroup.POST("status", medalHandle.Status)

		medalRecordHandle := excitation.NewMedalRecordHandle()
		// 领取记录列表
		medalGroup.POST("record-list", medalRecordHandle.List)
		// 勋章补发
		medalGroup.POST("record-reissue", medalRecordHandle.Reissue)
		medalGroup.POST("record-remove", medalRecordHandle.Remove)
		// 奖章规则列表
		medalGroup.POST("rule-list", excitation.NewMedalRuleHandle().List)

		// 奖章墙
		medalGroup.POST("medal-wall", medalHandle.MedalWall)
		// 奖章墙详情
		medalGroup.POST("medal-wall-detail", medalHandle.MedalWallDetail)
		// 领取奖章
		medalGroup.POST("record-member-received", medalRecordHandle.MemberReceived)
		// 已领取奖章信息
		medalGroup.POST("record-already-got", medalRecordHandle.AlreadyGot)

	}

	// 成长激励后台维度配置
	incentiveGroup := group.Group("incentive")
	{
		incentiveHandle := excitation.NewIncentiveHandle()
		// 积分规则维度配置
		incentiveGroup.POST("/config", incentiveHandle.IncentiveConfig)
		// 关系能力图谱-新增或修改指标配置
		incentiveGroup.POST("/relation/save", incentiveHandle.RelationSave)
		// 关系能力图谱-关联图谱规则详情
		incentiveGroup.POST("/relation/detail", incentiveHandle.RelationDetail)
		// 关系能力图谱-关联图谱规则详情列表
		incentiveGroup.POST("/relation/detail-list", incentiveHandle.RelationDetailList)

	}

	// 导购等级配置
	levelRuleGroup := group.Group("level-rule")
	{
		levelRuleHandle := excitation.NewLevelRuleHandle()
		// 新增或修改等级规则
		levelRuleGroup.POST("/save", levelRuleHandle.Save)
		// 等级规则详情
		levelRuleGroup.POST("/detail-list", levelRuleHandle.DetailList)
		// 等级规则列表
		levelRuleGroup.POST("/list", levelRuleHandle.List)

	}

	// 导购等级数据
	memberDataGroup := group.Group("member-data")
	{
		memberDataHandle := excitation.NewMemberDataHandle()
		// 后台列表
		memberDataGroup.POST("/list", memberDataHandle.List)
		// 导购等级数据-新增或更新
		memberDataGroup.POST("/save", memberDataHandle.Save)
		// 导购等级数据-更新导购等级
		memberDataGroup.POST("/update-level", memberDataHandle.UpdateLevel)
		// 导购导购积分&等级-首页
		memberDataGroup.POST("data-index", memberDataHandle.DataIndex)
		// 导购等级积分规则信息
		memberDataGroup.POST("grade-points", memberDataHandle.GradePoints)
		// 积分获取规则
		memberDataGroup.POST("get-points-rules", memberDataHandle.GetPointsRules)
		// 导购等级更新记录
		memberDataGroup.POST("level-record", memberDataHandle.LevelRecord)

	}

	// 积分规则
	integralRuleGroup := group.Group("integral-rule")
	{
		integralRuleHandle := excitation.NewIntegralRuleHandle()
		// 积分规则列表
		integralRuleGroup.POST("list", integralRuleHandle.List)
		// 积分规则配置信息
		integralRuleGroup.POST("configuration", integralRuleHandle.ParameterConfiguration)
		// 积分规则添加
		integralRuleGroup.POST("create", integralRuleHandle.Create)
		// 积分规则详情
		integralRuleGroup.POST("detail", integralRuleHandle.Detail)
		// 积分规则删除
		integralRuleGroup.POST("destroy", integralRuleHandle.Destroy)
		// 积分规则编辑
		integralRuleGroup.POST("update", integralRuleHandle.Update)

		// 积分规则说明
		integralRuleDescriptionHandle := excitation.NewIntegralRuleDescriptionHandle()
		// 积分规则说明修改添加
		integralRuleGroup.POST("description-save", integralRuleDescriptionHandle.CreateOUpdate)
		// 积分规则说明信息
		integralRuleGroup.POST("description-detail", integralRuleDescriptionHandle.Detail)
	}

	// 导购积分相关
	pointsGroup := group.Group("points")
	{
		pointsHandle := excitation.NewPointsHandle()
		// 管理后台积分列表
		pointsGroup.POST("record-list", pointsHandle.RecordList)
		// 积分手动导入
		pointsGroup.POST("import-manual-points", pointsHandle.AddManualPoints)
		// 小程序端导购积分列表
		pointsGroup.POST("list", pointsHandle.GetPointsList)

	}

	// 导购金币相关
	coinsGroup := group.Group("coins")
	{
		coinsHandle := excitation.NewCoinsHandle()
		// 管理后台金币列表
		coinsGroup.POST("record-list", coinsHandle.RecordList)
		// 金币手动导入
		coinsGroup.POST("import-manual-coins", coinsHandle.AddManualCoins)
		// 小程序端导购金币列表
		coinsGroup.POST("mini-coins-list", coinsHandle.GetCoinsList)
		// 金币配置详情
		coinsGroup.POST("coins-settings-info", coinsHandle.GetCoinsSettingsInfo)
		// 金币配置保存
		coinsGroup.POST("save-coins-settings", coinsHandle.SaveCoinsSettings)
	}
}
