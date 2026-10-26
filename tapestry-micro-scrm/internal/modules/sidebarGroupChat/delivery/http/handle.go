package sidebarGroupChathttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/sidebarGroupChat"
	"net/http"
)

var sidebarGroupChatService *sidebarGroupChat.SidebarGroupChatService
var dbClient *sqlplus.ClusterDB

// groupChatInfo 侧边栏社群详情
func groupChatInfo(c *gin.Context) {
	//-------------
	//ChatMember := sidebarGroupChat.WorkGroupCorpUser{}
	//sql := "SELECT unionid from micro_scrm_enterprise_%d.wx_user_corp_external  where openid='%s'"
	//sql = fmt.Sprintf(sql, 3306120225798016, "wm3XPqCgAAoApPViINhNW9KqN3Hz6hEw")
	//dbErr := dbClient.DB.Select(ctx.Wrap(c), ChatMember, sql)
	//if dbErr != nil {
	//	return
	//}
	//dbClient = db.DB()
	//unionid := []string{}
	//sql := "SELECT unionid from micro_scrm_enterprise_%d.wx_user_corp_external  where openid='%s'"
	//sql = fmt.Sprintf(sql, 3306120225798016, "wm3XPqCgAAoApPViINhNW9KqN3Hz6hEw")
	//err := dbClient.Get(ctx.Wrap(c), &unionid, sql)
	//if err != nil {
	//	return
	//}
	//if len(unionid) > 0 {
	//	//所有会员等级
	//	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	//	params := gin.H{
	//		"enterprise_id": 3306120225798016,
	//		"type":          "unionId",
	//		"value":         unionid[0],
	//	}
	//	// 查询用户信息
	//	url := config.CfgData.Service.MicroWork.Domain + "/service/v1/data-lake/query-member-grade"
	//	resp, err := client.Post(ctx.Wrap(c), url, clients.JSONBody(params))
	//	if err != nil {
	//		return
	//	}
	//	resStr, rErr := resp.Read(clients.Text)
	//	if rErr != nil {
	//		return
	//	}
	//	var response = gradeResultData{}
	//	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	//	if jErr != nil {
	//		return
	//	}
	//	if response.Status != 200 {
	//		return
	//	}
	//	updateSql := "update  micro_scrm_enterprise_%d.wx_user_corp_external set grade='%s' where unionid='%s'"
	//	updateSql = fmt.Sprintf(updateSql, 3306120225798016, response.Data.Grade, unionid[0])
	//	_, err = db.SqlxDB().Exec(updateSql)
	//	if err != nil {
	//		return
	//	}
	//}
	//---------------------------------------
	var groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest
	jsonError := c.BindJSON(&groupChatInfoRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := sidebarGroupChatService.GroupChatInfo(ctx.Wrap(c), groupChatInfoRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// groupChatMember 侧边栏社群成员列表
func groupChatMember(c *gin.Context) {
	var groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest
	jsonError := c.BindJSON(&groupChatInfoRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	if groupChatInfoRequest.Size == 0 {
		groupChatInfoRequest.Size = 10
	}
	if groupChatInfoRequest.Page == 0 {
		groupChatInfoRequest.Page = 1
	}
	res, err := sidebarGroupChatService.GroupChatUserList(ctx.Wrap(c), groupChatInfoRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// groupUsersAddEmployees 侧边栏社群成员列表
func groupUsersAddEmployees(c *gin.Context) {
	var groupUsersAddEmployeesRequest sidebarGroupChat.GroupUsersAddEmployeesRequest
	jsonError := c.BindJSON(&groupUsersAddEmployeesRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := sidebarGroupChatService.GroupUsersAddEmployees(ctx.Wrap(c), groupUsersAddEmployeesRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// groupUsersAddEmployees 侧边栏社群成员列表
func groupUsersJoinActivity(c *gin.Context) {
	var groupUsersJoinActivityRequest sidebarGroupChat.GroupUsersJoinActivityRequest
	jsonError := c.BindJSON(&groupUsersJoinActivityRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := sidebarGroupChatService.GroupUsersJoinActivity(ctx.Wrap(c), groupUsersJoinActivityRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

type gradeResultData struct {
	ID      string    `json:"id"`
	Code    string    `json:"code"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    gradeData `json:"data"`
}
type gradeData struct {
	IsUpdateGrade  string `json:"isUpdateGrade"`
	Grade          string `json:"grade"`
	GradeEndDate   string `json:"gradeEndDate"`
	GradeCode      string `json:"gradeCode"`
	GradeStartDate string `json:"gradeStartDate"`
}
