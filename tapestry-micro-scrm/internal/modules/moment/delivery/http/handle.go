package momenthttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/common"
	"math/rand"
	"micro-scrm/internal/logic"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/moment"
	"micro-scrm/internal/pkg/errs"
	"net/http"
	"time"
	"unicode/utf8"
)

var topicService *moment.TopicService
var memberService *moment.MemberService
var topMenService *moment.TopicMemberService
var departmentService *department.Service
var originalMemberService *members.Service
var enterpriseService *enterprise.Service

//var accessTokenServer *core.DefaultAccessTokenServer
//var ticketServer *jssdk.DefaultTicketServer
//var client *core.Client
var accessTokenServer *common.AccessTokenCacheService

// appearanceGet 获取朋友圈形象
func appearanceGet(c *gin.Context) {
	request := &moment.MomentMemberGetRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	momentMember, queryError := memberService.MemberGetInfo(ctx.Wrap(c), *request, enterpriseService)
	if queryError != nil {
		ginplus.ResError(c, errs.MomentMemberIsEmpty, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, momentMember)
}

// appearanceEdit 编辑朋友圈形象
func appearanceEdit(c *gin.Context) {
	request := &moment.MomentMemberEditRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	if utf8.RuneCountInString(request.Sign) > 30 {
		ginplus.ResError(c, errs.SignLengthErr, http.StatusOK)
		return
	}
	if utf8.RuneCountInString(request.Name) > 15 {
		ginplus.ResError(c, errs.NameLengthErr, http.StatusOK)
		return
	}
	affectedNum, updateError := memberService.MemberEditInfo(ctx.Wrap(c), *request)
	if updateError != nil {
		ginplus.ResError(c, errs.MomentMemberUpdateErr, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, affectedNum)
}

// topicCreate 发布朋友圈
func topicCreate(c *gin.Context) {
	request := &moment.TopicCreateRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, logicError := memberService.TopicCreate(ctx.Wrap(c), *request, topicService, topMenService, memberService, originalMemberService)
	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// topicEdit 编辑朋友圈
func topicEdit(c *gin.Context) {
	request := &moment.TopicEditRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, logicError := memberService.TopicEdit(ctx.Wrap(c), *request, topicService, topMenService, memberService, originalMemberService)
	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// topicDisplay 一个朋友圈展示
func topicDisplay(c *gin.Context) {
	request := &moment.TopicDisplayRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, logicError := memberService.TopicDisplay(ctx.Wrap(c), *request)
	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// topicDisplay 一个朋友圈删除
func topicDelete(c *gin.Context) {
	request := &moment.TopicDeleteRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, displayError := memberService.TopicDelete(ctx.Wrap(c), *request)
	if displayError != nil {
		ginplus.ResError(c, displayError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

func jsConfig(c *gin.Context) {
	var request struct {
		Url          string `json:"url"`
		EnterpriseID int64  `json:"enterprise_id"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	officialAccount, officialAccountErr := logic.NewOfficialAccount(ctx.Wrap(c), request.EnterpriseID)
	if officialAccountErr != nil {
		ginplus.ResError(c, officialAccountErr)
		return
	}

	jsConfig, jsConfigErr := officialAccount.GetJsConfig(request.Url)
	if jsConfigErr != nil {
		ginplus.ResError(c, officialAccountErr)
		return
	}

	ginplus.ResSuccess(c, jsConfig)

	// 不用了
	//officialInfo, _ := enterpriseService.OfficialInfo(ctx.Wrap(c), request.EnterpriseID)
	//var tokenServer *common.AccessTokenCacheService
	//var appid, secret string
	//if officialInfo != nil {
	//	appid = officialInfo.AppID
	//	secret = officialInfo.Secret
	//} else {
	//	appid = config.CfgData.Mp.AppID
	//	secret = config.CfgData.Mp.Secret
	//}
	//tokenServer = common.NewAccessTokenCacheService(cache.NewRedisCache(redis.Client()), appid, secret, clients.NewHTTPClient())
	//accessToken, err := tokenServer.GetToken(ctx.Wrap(c))
	//if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	//}
	//log.Info(ctx.Wrap(c)).Msgf("accessToken:%s", accessToken)
	//ticketService := jssdk.NewTicketCacheService(cache.NewRedisCache(redis.Client()), appid, accessToken, clients.NewHTTPClient())
	//ticket, err := ticketService.GetTicket(ctx.Wrap(c))
	//log.Info(ctx.Wrap(c)).Msgf("ticket:%s", ticket)
	//if err != nil {
	//	log.Info(ctx.Wrap(c)).Msgf("获取jsticket失败, err:%s", err.Error())
	//	ginplus.ResError(c, err)
	//	return
	//}
	//nonceStr := getRandomString(32)
	//timeStamp := time.Now().Unix()
	//signature := jssdk.WXConfigSign(ticket, nonceStr, strconv.FormatInt(timeStamp, 10), request.Url)
	//ginplus.ResSuccess(c, struct {
	//	AppID     string `json:"appId"`
	//	Timestamp int64  `json:"timestamp"`
	//	NonceStr  string `json:"nonce_str"`
	//	Signature string `json:"signature"`
	//}{
	//	AppID:     appid,
	//	Timestamp: timeStamp,
	//	NonceStr:  nonceStr,
	//	Signature: signature,
	//})
}

func getRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// contentList 按内容搜索
func contentList(c *gin.Context) {
	request := &moment.ContentListRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, displayError := memberService.ContentList(ctx.Wrap(c), *request)
	if displayError != nil {
		ginplus.ResError(c, displayError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// memberList 按员工搜索
func memberList(c *gin.Context) {
	request := &moment.MemberListRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, displayError := memberService.MemberList(ctx.Wrap(c), *request, departmentService, originalMemberService)
	if displayError != nil {
		ginplus.ResError(c, displayError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// memberTopic 一个员工的朋友圈列表
func memberTopic(c *gin.Context) {
	request := &moment.MemberTopicRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, error := memberService.MemberTopic(ctx.Wrap(c), *request)
	if error != nil {
		ginplus.ResError(c, error, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// topicMemberDelete 一个员工的朋友圈删除
func topicMemberDelete(c *gin.Context) {
	request := &moment.TopicMemberDeleteRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, error := memberService.TopicMemberDelete(ctx.Wrap(c), *request)
	if error != nil {
		ginplus.ResError(c, error, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// topicMemberEdit 一个员工的朋友圈编辑
func topicMemberEdit(c *gin.Context) {
	request := &moment.TopicMemberEditRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	result, error := memberService.TopicMemberEdit(ctx.Wrap(c), *request, topicService, topMenService, memberService, originalMemberService)
	if error != nil {
		ginplus.ResError(c, error, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}
