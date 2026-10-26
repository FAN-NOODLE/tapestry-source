package member

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/handle"
	memberReq "micro-scrm/internal/app/http/request/member"
	departmentsModels "micro-scrm/internal/app/models/departments"
	membersModel "micro-scrm/internal/app/models/members"
	cacheSvc "micro-scrm/internal/app/service/cache"
	membersSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"strconv"
	"strings"
)

type MembersHandle struct {
	handle.BaseHandle
	membersService *membersSvc.MembersService
}

func NewMembersHandle() MembersHandle {
	return MembersHandle{
		membersService: membersSvc.NewMembersService(),
	}
}

// GetDetail 导购部门详情
func (s MembersHandle) GetDetail(c *gin.Context) {

	var request memberReq.MemberDetailReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, err := s.membersService.GetDetail(s.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, result)
	return
}

// Info 导购信息
func (s MembersHandle) Info(c *gin.Context) {

	var request memberReq.MemberDetailReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, err := s.membersService.GetInfoByMobile(s.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err != nil {
		s.ResError(c, err)
		return
	}

	var res membersModel.MemberInfo
	res.Id = result.ID
	res.EnterpriseId = result.EnterpriseID
	res.Name = result.Name
	res.Userid = result.Userid
	res.Mobile = result.Mobile
	res.Type = result.Type
	res.Avatar = result.Avatar
	res.QrPosterCode = result.QrPosterCode
	res.Status = result.Status
	if len(result.MemberDepartments) > 0 {
		res.DepartmentName = result.MemberDepartments[0].Department.Name
		res.StoreId = result.MemberDepartments[0].Department.ThirdID
	}
	s.ResSuccess(c, res)
	return
}

func (s MembersHandle) StoreLists(c *gin.Context) {
	var request memberReq.MemberStoreReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, err := s.membersService.GetStoreLists(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, result)
	return
}

// StaffLists coach 店长(SM)获取导购(SA)列表
func (s MembersHandle) StaffLists(c *gin.Context) {
	var request memberReq.MemberStaffReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	lists, err := s.membersService.GetStaffLists(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, lists)
}

func (s MembersHandle) MemberListByThirdIds(c *gin.Context) {
	var request memberReq.MemberListByThirdIdsReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	dpList, err := departmentsDao.NewDepartments().GetDepartmentsInfoThirdIds(s.WarpContext(c), request.EnterpriseId, request.ThirdIds)
	if err != nil {
		return
	}
	dpMap := make(map[int64]departmentsModels.DepartmentsModel)
	var dpIds []int64
	for _, dpInfo := range dpList {
		dpMap[dpInfo.ID] = dpInfo
		dpIds = append(dpIds, dpInfo.ID)
	}
	request.DpIds = dpIds
	list, err := membersDao.NewMemberDepartments().GetMemberListByDpIds(s.WarpContext(c), request)

	var resList = make([]membersModel.MemberInfo, 0)

	for _, row := range list {
		if row.DeletedAt > 0 || row.Member.DeletedAt > 0 || row.Member.ID == 0 {
			continue
		}
		dpName := ""
		thirdId := 0
		if _, ok := dpMap[row.DepartmentID]; ok {
			dpName = dpMap[row.DepartmentID].Name
			thirdId = dpMap[row.DepartmentID].ThirdID
		}
		tmp := membersModel.MemberInfo{
			Id:             row.MemberID,
			EnterpriseId:   row.EnterpriseID,
			Name:           row.Member.Name,
			Userid:         row.Member.Userid,
			Mobile:         row.Member.Mobile,
			Type:           row.Member.Type,
			Avatar:         row.Member.Avatar,
			QrPosterCode:   row.Member.QrPosterCode,
			Status:         row.Member.Status,
			DepartmentName: dpName,
			ThirdId:        thirdId,
		}
		resList = append(resList, tmp)
	}
	s.ResSuccess(c, resList)
}

func (s MembersHandle) ReportData(c *gin.Context) {
	var request memberReq.ReportDataReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	cc := s.WarpContext(c)

	threading.GoSafe(func() error {
		queryStr := fmt.Sprintf("enterprise_id = %d and userid = %s", request.EnterpriseId, request.UserId)
		member, err := membersDao.NewMembers().GetOne(cc, queryStr)
		if err != nil {
			return err
		}

		shenCeClient := microwork.NewShenCe()
		memberDepartments, _ := departmentsDao.NewMemberDepartments().GetThirdsByMemberIds(cc, request.EnterpriseId, member.ID)
		if len(memberDepartments) > 0 {
			// 获取当前员工的门店编码
			storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(cc, map[string]any{
				"enterprise_id":        request.EnterpriseId,
				"wechat_department_id": memberDepartments[0],
			})

			typeStr := "sa"
			if member.Type == 4 {
				typeStr = "sm"
			}
			_, _ = shenCeClient.ProfileSet(cc, microwork.ProfileSetReq{
				EnterpriseId: request.EnterpriseId,
				DistinctId:   member.Userid,
				IsLoginId:    true,
				Profiles: map[string]any{
					"userid":     member.Userid,
					"mobile":     member.Mobile,
					"position":   member.Position,
					"type":       typeStr,
					"store_code": storeDepartmentInfo.ThirdStoreCode,
					"gender":     member.Gender,
					"avatar":     member.Avatar,
				},
			})
		}
		return nil
	})

	s.ResOk(c)

	return
}

// DownloadWhiteList 下载中心白名单列表
func (s MembersHandle) DownloadWhiteList(c *gin.Context) {
	var request struct {
		EnterpriseId int64 `json:"enterprise_id"`
		AdminId      int64 `json:"admin_id"`
	}
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(s.WarpContext(c), request.EnterpriseId, "download_white_list")
	if err != nil {
		s.ResError(c, err)
		return
	}
	var isWhite bool
	adminIds := strings.Split(config.Name, ",")
	if len(adminIds) == 0 {
		s.ResSuccess(c, isWhite)
		return
	}
	adminIdsStr := strconv.FormatInt(request.AdminId, 10)
	if helper.InArray(adminIdsStr, adminIds) {
		isWhite = true
	}
	s.ResSuccess(c, isWhite)
	return
}

// FirstOpen 第一次进入某页面给的前端引导提示
func (s MembersHandle) FirstOpen(c *gin.Context) {
	var request memberReq.FirstOpenReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	isFirstOpen, err := membersSvc.NewMembersService().GetFirstOpen(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, isFirstOpen)
	return
}

func (s MembersHandle) GetMemberListByStoreCodes(c *gin.Context) {
	var request memberReq.StoreCodesReq
	if err := c.ShouldBindJSON(&request); err != nil {
		s.ResError(c, err)
		return
	}
	list, err := membersSvc.NewMembersService().GetMemberListByStoreCodes(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, list)
}
