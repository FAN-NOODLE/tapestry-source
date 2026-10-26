package delivery

import (
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/groupmsg"
	groupMsgApi "micro-scrm/internal/modules/groupmsg/api"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/ws_trace"
	wsTraceMysql "micro-scrm/internal/modules/ws_trace/respository/mysql"
)

// 初始化
func InitServiceDefault() {
	groupmsg.MemberService = members.New(membersMysql.New())
	groupmsg.WorkApi = &groupMsgApi.WorkApi{}
	groupmsg.EnterpriseService = enterprise.New(enterpriseMysql.New())
	groupmsg.WsTraceService = ws_trace.NewService(wsTraceMysql.NewWsTrace())
}
