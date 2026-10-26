package notificationSvc

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/notification"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

type CrowdSvc struct {
	service.Base
	crowdDao     *notificationDao.CrowdDao
	crowdDataDao *notificationDao.CrowdDataDao
}

func NewCrowdSvc() *CrowdSvc {
	return &CrowdSvc{
		crowdDao:     notificationDao.NewCrowdDao(),
		crowdDataDao: notificationDao.NewCrowdDataDao(),
	}
}

func (s *CrowdSvc) Save(c ctx.Context, params notificationReq.CrowdReq) error {
	memberInfo, err := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseId, params.Mobile)
	if err != nil {
		return err
	}
	nameCheck, _ := s.crowdDao.GetOne(c, params.EnterpriseId, "name=?", params.Name)
	if params.Id > 0 {
		info, err := s.crowdDao.GetOne(c, params.EnterpriseId, "id=?", params.Id)
		if info.ID == 0 {
			return errors.New("找不到该数据")
		}
		if nameCheck.ID > 0 && nameCheck.ID != info.ID {
			return errors.New("该白名单名称已存在")
		}
		err = s.crowdDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"id": params.Id,
		}, map[string]interface{}{
			"name":       params.Name,
			"file_url":   params.FileUrl,
			"type":       params.Type,
			"member_num": len(params.MemberIds),
			//"admin_id":   params.AdminId,
			//"admin_name": memberInfo.Name,
			"updated_at": time.Now().Unix(),
		})
		if err != nil {
			return err
		}
		err = s.crowdDataDao.UnscopedDeleteData(c, params.EnterpriseId, map[string]interface{}{
			"crowd_id": params.Id,
		})
		if err != nil {
			return err
		}
		err = s.saveCrowdData(c, params.EnterpriseId, params.Id, params.MemberIds)
		if err != nil {
			return err
		}
	} else {
		if nameCheck.ID > 0 {
			return errors.New("该白名单名称已存在")
		}
		crowdId := app.UID()
		err = s.crowdDao.CreateData(c, params.EnterpriseId, microScrmModel.NotificationCrowdModel{
			BaseModel: models.CreateBase(crowdId),
			Name:      params.Name,
			FileUrl:   params.FileUrl,
			Type:      params.Type,
			MemberNum: len(params.MemberIds),
			AdminID:   params.AdminId,
			AdminName: memberInfo.Name,
		})
		if err != nil {
			return err
		}
		err = s.saveCrowdData(c, params.EnterpriseId, crowdId, params.MemberIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CrowdSvc) List(c ctx.Context, params notificationReq.CrowdListReq) ([]microScrmModel.NotificationCrowdModel, int64, error) {
	list, total, err := s.crowdDao.GetPagination(c, params)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (s *CrowdSvc) Detail(c ctx.Context, params notificationReq.CrowdDetailReq) (notificationResp.CrowdDetailResp, error) {
	info, err := s.crowdDao.GetOneData(c, params.EnterpriseId, "id=?", params.Id)
	if err != nil {
		return notificationResp.CrowdDetailResp{}, err
	}
	resp := notificationResp.CrowdDetailResp{
		Id:   info.ID,
		Name: info.Name,
		Type: info.Type,
	}
	if len(info.DataList) > 0 {
		var memberIds []int64
		for _, row := range info.DataList {
			memberIds = append(memberIds, row.MemberId)
		}
		list, err := membersDao.NewMembers().GetAllById(c, params.EnterpriseId, memberIds)
		if err != nil {
			return notificationResp.CrowdDetailResp{}, err
		}
		for _, data := range list {
			var dpId int64
			if len(data.MemberDepartments) > 0 {
				dpId = data.MemberDepartments[0].DepartmentID
			}
			resp.MemberList = append(resp.MemberList, notificationResp.CrowdDataResp{
				MemberId:     data.ID,
				MemberName:   data.Name,
				DepartmentId: dpId,
			})
		}
	}
	return resp, nil
}

func (s *CrowdSvc) saveCrowdData(c ctx.Context, enterpriseId int64, crowdId int64, memberIds []int64) error {
	var dataArr []microScrmModel.NotificationCrowdDataModel
	for _, memberId := range memberIds {
		dataArr = append(dataArr, microScrmModel.NotificationCrowdDataModel{
			BaseModel: models.CreateBase(app.UID()),
			CrowdId:   crowdId,
			MemberId:  memberId,
		})
	}
	err := s.crowdDataDao.CreateInBatchesData(c, enterpriseId, dataArr)
	if err != nil {
		return err
	}
	return nil
}

func (s *CrowdSvc) DeleteCrowd(c ctx.Context, params notificationReq.CrowdDetailReq) error {
	list, err := notificationDao.NewInfoMergeDao().GetMany(c, params.EnterpriseId, "status = ? AND crowd_type = ?", microScrmModel.NotificationMergeNotSend, microScrmModel.NotificationCrowdWhite)
	for _, model := range list {
		tmpCrowdIds := strings.Split(model.CrowdIds, ",")
		for _, id := range tmpCrowdIds {
			if int(params.Id) == helper.StringToInt(id) {
				return errors.New("该白名单已关联尚未发布的整合传讯")
			}
		}
	}
	err = s.crowdDao.DeleteData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.Id,
	})
	if err != nil {
		return err
	}
	err = s.crowdDataDao.DeleteData(c, params.EnterpriseId, map[string]interface{}{
		"crowd_id": params.Id,
	})
	if err != nil {
		return err
	}
	return nil
}
