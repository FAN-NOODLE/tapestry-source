package membershttp

import (
	"net/http"
	"strings"

	"micro-scrm/internal/modules/department"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

var service *members.Service
var departmentService *department.Service

// memberDetail 成员详情
func memberDetail(c *gin.Context) {
	req := &members.MemberDetailRequest{}
	err := c.BindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if req.ID == 0 {
		ginplus.ResError(c, errs.MemberIDIsEmpty, errs.ErrorStatus)
		return
	}
	entity, _ := service.Detail(ctx.Wrap(c), req.ID)
	ginplus.ResSuccess(c, entity)
}

// memberDetailByUserid 成员详情
func memberDetailByUserid(c *gin.Context) {
	req := members.Member{}
	err := c.BindJSON(&req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	entity := service.FindMemberByUserid(ctx.Wrap(c), req)
	ginplus.ResSuccess(c, entity)
}

// memberList 成员列表
func memberList(c *gin.Context) {
	req := &members.MemberListRequest{}
	err := c.BindJSON(req)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	req.IDCount = len(req.IDs)
	req.UserIDCount = len(req.UserIDs)
	if req.Name != "" {
		var build strings.Builder
		build.WriteString("%")
		build.WriteString(req.Name)
		build.WriteString("%")
		req.Name = build.String()
	}
	count, entity := service.Page(ctx.Wrap(c), *req)
	res := members.MemberListResponse{}

	if entity != nil {
		res.MemberList = *entity
	} else {
		res.MemberList = make([]members.Member, 0)
	}

	res.Total = count
	res.Page = req.Page
	res.Size = req.Size
	ginplus.ResSuccess(c, res)
}

// memberAdd 成员添加
func memberAdd(c *gin.Context) {
	entity := &members.Member{}
	err := c.BindJSON(entity)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	service.Add(ctx.Wrap(c), entity)
	ginplus.ResSuccess(c, entity)
}

// syncContactAndDept 同步通讯录
func syncContactAndDept(c *gin.Context) {
	// 获取请求参数
	syncContactReq := &members.SyncContact{}
	reqErr := c.ShouldBind(syncContactReq)
	if reqErr != nil {
		log.Info(ctx.Wrap(c)).Msgf("syncContactAndDept请求参数错误：%+v, 错误信息：%s", syncContactReq, reqErr.Error())
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	cc := ctx.Wrap(c)
	ok := service.GetSyncContact(cc, syncContactReq.EnterpriseID)
	if !ok {
		ginplus.ResError(c, errs.SyncContactRunning)
		return
	}
	if err := service.SyncContact(cc, syncContactReq.EnterpriseID, syncContactReq.AdminID); err != nil {
		response.Notice(c, errs.ActionErr, err)
		return
	}
	response.Success(c, nil)
}

// momentList 朋友圈选择员工
func momentList(c *gin.Context) {
	request := &members.MomentListRequest{}
	jsonError := c.BindJSON(request)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	result, logicError := service.MomentList(ctx.Wrap(c), request)
	if len(result.List) > 0 {
		for i, info := range result.List {
			if info.Mobile != "" {
				result.List[i].Mobile = info.Mobile[:len(info.Mobile)-8] + "****" + info.Mobile[len(info.Mobile)-4:]
			}

		}
	}

	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
	}
	ginplus.ResSuccess(c, result)
}

// simpleList 简单列表查询
func simpleList(c *gin.Context) {
	memberList := &members.MemberListRequest{}
	err := c.BindJSON(memberList)
	if err != nil {
		response.Notice(c, errs.ParamsErr, err)
		return
	}

	// 验证企业ID必传
	if memberList.EnterpriseID == 0 {
		response.Notice(c, errs.ParamsErr, err)
		return
	}

	// 初始化分页参数
	if memberList.Page == 0 {
		memberList.Page = 1
	}
	if memberList.Size == 0 {
		memberList.Size = 10
	}

	// 设置默认值
	if memberList.ID != 0 {
		memberList.IDs = append(memberList.IDs, memberList.ID)
	}
	if memberList.Userid != "" {
		memberList.UserIDs = append(memberList.UserIDs, memberList.Userid)
	}

	memberList.IDCount = len(memberList.IDs)
	memberList.UserIDCount = len(memberList.UserIDs)
	total, list, err := service.SimpleMemberList(ctx.Wrap(c), memberList)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errs.ParamsErr, err)
		return
	}

	response.SuccessList(c, total, memberList.Page, memberList.Size, list)
}

// authInfo 通过 code 获取用户信息
func authInfo(c *gin.Context) {
	request := members.AuthInfoRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	auth, err := service.AuthInfoByCode(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, auth)
}

// momentList 朋友圈选择BA员工
func membersBaList(c *gin.Context) {
	request := &members.MomentListRequest{}
	jsonError := c.BindJSON(request)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	result, logicError := service.MemberList(ctx.Wrap(c), request)

	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, result)
}

// getMembersOfDepartments 获取部门下所有员工列表
func getMembersOfDepartments(c *gin.Context) {
	var req struct {
		Departments  []int64 `json:"departments" desc:"所属部门"`
		EnterpriseID int64   `json:"enterprise_id"`
		MemType      string  `json:"mem_type" desc:"员工类型"`
	}
	if err := c.BindJSON(&req); err != nil {
		ginplus.ResError(c, err)
		return
	}

	var totalDepartments []int64
	dpMap, _ := departmentService.GetChildrenForDepartment(ctx.Wrap(c), req.EnterpriseID)
	for _, dpId := range req.Departments {
		if children, ok := dpMap[dpId]; ok {
			totalDepartments = append(totalDepartments, children...)
		}
	}
	dpMembers, err := service.GetMembersInfoByDepartmentIds(ctx.Wrap(c), req.EnterpriseID, totalDepartments, req.MemType)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, dpMembers)
	return
}

// code2Session 通过 code 获取用户信息
func code2Session(c *gin.Context) {
	request := members.Code2SessionRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	auth, err := service.Code2Session(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, auth)
}
