package chatJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/chat"
	chatRepository "micro-scrm/internal/modules/chat/repository/mysql"
	"micro-scrm/internal/pkg/sdk"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

// GetMsgJob 拉取会话存档消息
func GetMsgJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	chatService := chat.NewMsg(chatRepository.NewChatMsg())

	receiveUserService := *chat.NewReceive(chatRepository.NewChatReceiveUser())
	objectService := *chat.NewObject(chatRepository.NewChatObject())
	memberRelationService := *chat.NewMember(chatRepository.NewChatMember())
	groupChatService := *chat.NewGroupChat(chatRepository.NewWorkGroupChat())
	groupChatMemberService := *chat.NewGroupChatMember(chatRepository.NewWorkGroupChatMember())
	chatService.GetMsg(c, receiveUserService, objectService, memberRelationService, groupChatService, groupChatMemberService)

	return nil
}

// GetMsgMediaJob 获取会话存档媒体文件
func GetMsgMediaJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 设置请求的 HTTP 客户端超时时间
	sdk.ClientObject = sdk.NewSdkClient(clients.Timeout(time.Second * 20))

	chatService := chat.NewMsg(chatRepository.NewChatMsg())

	msgMediaService := *chat.NewMsgMedia(chatRepository.NewChatMsgMedia())
	chatService.GetMsgMedia(c, msgMediaService)

	return nil
}
