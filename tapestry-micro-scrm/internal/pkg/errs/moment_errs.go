package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	MomentMemberIsEmpty   = errorx.NewErrorMsg(140001, "用户朋友圈形象为空")
	MomentMemberUpdateErr = errorx.NewErrorMsg(140001, "用户朋友圈形象编辑失败")
	GetDepartmentsErr     = errorx.NewErrorMsg(140001, "获取部门失败")
	GetMembersErr         = errorx.NewErrorMsg(140001, "部门下无员工")
	JsonParseErr          = errorx.NewErrorMsg(140001, "JSON转换失败")
	PublishTimeErr        = errorx.NewErrorMsg(140001, "发布时间不能小于当前时间")
	TopicInsertErr        = errorx.NewErrorMsg(140001, "moment_topic新增失败")
	TopicDisplayErr       = errorx.NewErrorMsg(140001, "该朋友圈不存在")
	TopicDeleteErr        = errorx.NewErrorMsg(140001, "朋友圈删除失败")
	MemberInsertErr       = errorx.NewErrorMsg(140001, "moment_member新增失败")
	TopicMemberInsertErr  = errorx.NewErrorMsg(140001, "moment_topic_member新增失败")
	MemberDeleteError     = errorx.NewErrorMsg(140001, "moment_member删除失败")
	SignLengthErr         = errorx.NewErrorMsg(140001, "朋友圈形象内容最多30个字")
	NameLengthErr         = errorx.NewErrorMsg(140001, "朋友圈形象名称最多15个字")
	TopicNotMember        = errorx.NewErrorMsg(140001, "该朋友圈下无员工")
	UpdateMemberTimeErr   = errorx.NewErrorMsg(140001, "员工更新时间失败")
	ActivityInfoErr       = errorx.NewErrorMsg(140001, "获取活动数据失败")
	WelcomeDetailErr      = errorx.NewErrorMsg(140001, "获取裂变活动欢迎语数据失败")
	FissionActivityErr    = errorx.NewErrorMsg(140001, "获取用户发起的裂变活动数据失败")
	ContactQrCodeErr      = errorx.NewErrorMsg(140001, "获取用户发起的裂变活动群二维码数据失败")
)
