package work

// 企业微信的接口列表
var (
	// SToken 获取token
	SToken = "%s/s/token"
	// QyGw 调用企业微信接口
	QyGw = "%s/b/%s/qy/gw"
	// QyGwPic 调用企业微信接口
	QyGwPic = "%s/a/%s/%s/qy/media/upload"
	// QyDecryptUrl 解密接口
	QyDecryptUrl = "%s/s/qy/decrypt"
	// CallbackAdd 添加外部联系人回调
	CallbackAdd = "%s/s/callback/add"
	// CallbackDel 外部联系人回调删除
	CallbackDel = "%s/s/callback/del"
	// CallbackList 外部联系回调地址人列表
	CallbackList = "%s/s/callback/list"
	/*----------------------企业微信接口地址------------------------*/
	// DepartmentList 组织列表
	DepartmentList = "cgi-bin/department/list"
	// GetAgent 获取指定的应用详情
	GetAgent = "cgi-bin/agent/get"
	// MemberList 成员列表
	MemberList = "cgi-bin/user/list"
	// AddTagUsers 给标签添加成员
	AddTagUsers = "cgi-bin/tag/addtagusers"
	// SaTagList 获取标签列表
	SaTagList = "cgi-bin/tag/list"
	// MemberGet 成员详情
	MemberGet = "cgi-bin/user/get"
	// 发送应用消息
	SendMessage = "cgi-bin/message/send"
	// 外部联系人详情
	ExternalContactGet      = "cgi-bin/externalcontact/get"
	BatchExternalContactGet = "cgi-bin/externalcontact/batch/get_by_user"
	// 获取客户群详情
	ExternalGroupChatGet = "cgi-bin/externalcontact/groupchat/get"
	// 修改外部联系人备注信息
	EditExternalRemarks = "cgi-bin/externalcontact/remark"
	// 获取会话同意情况
	CheckSingleAgree = "cgi-bin/msgaudit/check_single_agree"
	// 添加联系我二维码
	AddExternalContact = "cgi-bin/externalcontact/add_contact_way"
	// 编辑联系我二维码
	UpdateExternalContact = "cgi-bin/externalcontact/update_contact_way"
	// 删除联系我二维码
	DeleteExternalContact = "cgi-bin/externalcontact/del_contact_way"
	// 发送欢迎语
	SendWelcome = "cgi-bin/externalcontact/send_welcome_msg"
	// 用户标签处理
	UserTagHandlerUrl = "cgi-bin/externalcontact/mark_tag"
	// 获取机器人信息
	GetRobotInfoUrl = "cgi-bin/msgaudit/get_robot_info"
	// 获取会话内容存档内部群信息
	GetGroupChatInfo = "cgi-bin/msgaudit/groupchat/get"
	// 获取企业群发消息发送结果
	GetGroupMsgResult = "cgi-bin/externalcontact/get_group_msg_result"
	// 添加企业群发消息任务
	QwMassHair = "cgi-bin/externalcontact/add_msg_template"
	// 外部联系人列表
	ExternalContactGetUrl = "cgi-bin/externalcontact/list"
	// 获取客户群列表
	ExternalGroupChatList = "cgi-bin/externalcontact/groupchat/list"
	// 获取客户群详情
	ExternalGroupChatListURL = "cgi-bin/externalcontact/groupchat/list"
	// 获取用户身份
	AuthUserInfoUrl = "cgi-bin/user/getuserinfo"
	// 获取企业的jsapi_ticket
	JsTicket = "cgi-bin/get_jsapi_ticket"
	// 获取应用的jsapi_ticket
	AgentTicket = "cgi-bin/ticket/get"
	// 获取企业群发成员执行结果
	GroupMsgSendResult = "cgi-bin/externalcontact/get_groupmsg_send_result"
	// 上传图片
	MediaUploadURL = "cgi-bin/media/upload"
	// 获取「联系客户统计」数据
	GetUserBehaviorData = "cgi-bin/externalcontact/get_user_behavior_data"
	// 获取用户身份
	Jscode2session = "cgi-bin/miniprogram/jscode2session"
	/******************标签*********************/
	// 获取企业标签信息
	GetCorpTagListUrl = "cgi-bin/externalcontact/get_corp_tag_list"
	// 删除企业标签
	DelCorpTagUrl = "cgi-bin/externalcontact/del_corp_tag"
	// 新增企业标签
	AddCorpTagUrl = "cgi-bin/externalcontact/add_corp_tag"
	// 更新企业标签
	UpdateCorpTagUrl = "cgi-bin/externalcontact/edit_corp_tag"
	// TagCreate 创建标签
	TagCreate = "cgi-bin/tag/create"
	// TagUpdate 更新标签名字
	TagUpdate = "cgi-bin/tag/update"
	// TagDelete 更新标签名字
	TagDelete = "cgi-bin/tag/delete"
	// TagUsers 获取标签成员
	TagUsers = "cgi-bin/tag/get"
	// TagUsersAdd 增加标签成员
	TagUsersAdd = "cgi-bin/tag/addtagusers"
	// TagUsersDel 删除标签成员
	TagUsersDel = "cgi-bin/tag/deltagusers"
	// Taglist 获取标签列表
	Taglist = "cgi-bin/tag/list"
	// AddMomentTask 企业发表内容到客户的朋友圈
	AddMomentTask = "cgi-bin/externalcontact/add_moment_task"
	// GetMomentTaskResult 获取朋友圈任务创建结果
	GetMomentTaskResult = "cgi-bin/externalcontact/get_moment_task_result"
	// GetMomentTask 获取客户朋友圈企业发表的列表
	GetMomentTask = "cgi-bin/externalcontact/get_moment_task"
	// GetMomentComments  获取客户朋友圈企业发表的评论
	GetMomentComments = "cgi-bin/externalcontact/get_moment_comments"
	// GetMomentSendResult 获取客户朋友圈发表后的可见客户列表
	GetMomentSendResult = "cgi-bin/externalcontact/get_moment_send_result"
	// CancelMomentTask  停止发表企业朋友圈
	CancelMomentTask = "cgi-bin/externalcontact/cancel_moment_task"
	// CancelGroupMsgSend 停止企业群发
	CancelGroupMsgSend = "cgi-bin/externalcontact/cancel_groupmsg_send"
)
