package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	WelcomeDetailError       = errorx.NewErrorMsg(160001, "欢迎语获取详情失败")
	WelcomeEditError         = errorx.NewErrorMsg(160001, "编辑欢迎语获取失败")
	WelcomeDeleteError       = errorx.NewErrorMsg(160001, "删除欢迎语获取失败")
	WelcomeListError         = errorx.NewErrorMsg(160001, "欢迎语列表查询失败")
	WelcomeSendError         = errorx.NewErrorMsg(160001, "欢迎语发送失败")
	WelcomeValidateError     = errorx.NewErrorMsg(160001, "获取渠道列表失败")
	WelcomeCreateError       = errorx.NewErrorMsg(140001, "欢迎语创建失败")
	WelcomeUpdateError       = errorx.NewErrorMsg(140002, "欢迎语更新失败")
	WelcomeMemberCreateError = errorx.NewErrorMsg(140001, "欢迎语员工创建失败")
	WelcomeMemberDeleteError = errorx.NewErrorMsg(140001, "欢迎语员工删除失败")
	WeChatImgUrlCacheError   = errorx.NewErrorMsg(140001, "微信图片素材库失败")

	LabelNotNllError           = errorx.NewErrorMsg(160100, "请添加标签")
	ChannelNameNotNllError     = errorx.NewErrorMsg(160101, "请输入渠道名称")
	ChannelCreateError         = errorx.NewErrorMsg(160102, "渠道创建失败")
	ChannelListError           = errorx.NewErrorMsg(160103, "渠道查询失败")
	ChannelDetailSearchError   = errorx.NewErrorMsg(160104, "渠道详情查询失败")
	ChannelDetailExists        = errorx.NewErrorMsg(160105, "渠道名称已经存在")
	WelcomeChannelDetailExists = errorx.NewErrorMsg(160106, "此渠道已经配置了欢迎语")
	ChannelNotExists           = errorx.NewErrorMsg(160107, "渠道信息不存在")
)
