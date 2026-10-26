package storeVisitSvc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type PdfInfoSvc struct {
	service.Base
	visitMainDao  *storeVisitDao.MainDao
	showRecordDao *storeVisitDao.ShowRecordDao
	deptMemberDao *storeVisitDao.DeptMemberDao
}

func NewPdfInfoSvc() *PdfInfoSvc {
	return &PdfInfoSvc{
		visitMainDao:  storeVisitDao.NewMainDao(),
		showRecordDao: storeVisitDao.NewShowRecordDao(),
		deptMemberDao: storeVisitDao.NewDeptMemberDao(),
	}
}
func (s *PdfInfoSvc) getPdfInfo(c ctx.Context, req storeVisitReq.VisitInfoReq) (info storeVisitResp.PdfReportInfo, err error) {
	visitInfo, _ := s.visitMainDao.GetInfoPreloadCreator(c, req)
	helper.CoverStruct(&info, &visitInfo)
	info.UpdatedAt = visitInfo.UpdatedAt
	storeDepartmentInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": visitInfo.ThirdId,
	})
	info.StoreName = storeDepartmentInfo.StoreName
	if req.Language == "en" && storeDepartmentInfo.InventoryName != "" {
		info.StoreName = storeDepartmentInfo.InventoryName
	}
	if len(visitInfo.MemberList) > 0 {
		creator, _ := s.deptMemberDao.GetOne(c, req.EnterpriseId, "member_id = ?", visitInfo.MemberList[0].MemberID)
		if creator.ID > 0 {
			info.CreatorName = creator.Name
			info.CreatorDeptName = creator.DeptName
		} else {
			memberDepartmentList, _ := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, req.EnterpriseId, []int64{visitInfo.MemberList[0].MemberID})
			if len(memberDepartmentList) > 0 {
				memberInfo := memberDepartmentList[0]
				info.CreatorName = memberInfo.Name
				if len(memberInfo.MemberDepartments) > 0 {
					info.CreatorDeptName = memberInfo.MemberDepartments[0].Department.Name
				}
			}
		}
	}
	recordList, err := s.showRecordDao.GetListPreloadMediaOrderByType(c, req.EnterpriseId, "visit_id = ?", req.ID)
	for _, record := range recordList {
		tmpRecord := storeVisitResp.PdfShowRecord{
			TypeName: record.TypeName,
			Content:  record.Content,
		}
		if len(record.MediaList) > 0 {
			for _, media := range record.MediaList {
				if media.URL == "" {
					continue
				}
				thumbUrl := fmt.Sprintf("%s?imageMogr2/thumbnail/800x800/pad/1/color/I0Y0RjdGYg==", media.URL)
				if media.MediaType == storeVisitModel.Video {
					if media.VideoImageUrl == "" {
						continue
					}
					thumbUrl = fmt.Sprintf("%s?imageMogr2/thumbnail/800x800/pad/1/color/I0Y0RjdGYg==", media.VideoImageUrl)
				}
				tmpRecord.MediaList = append(tmpRecord.MediaList, storeVisitResp.PdfMedia{
					MediaType:     media.MediaType,
					URL:           media.URL,
					ThumbUrl:      thumbUrl,
					VideoImageUrl: media.VideoImageUrl,
				})

			}
		}
		info.ShowRecordList = append(info.ShowRecordList, tmpRecord)
	}
	formatPdfInfo(req, &info)
	return info, nil

}
func formatPdfInfo(req storeVisitReq.VisitInfoReq, info *storeVisitResp.PdfReportInfo) {

	if len(info.ShowRecordList) > 0 {
		if info.Overview == "" {
			info.Overview += helper.RuleThree(req.Language == "en", "This store visit has been completed,primarily focusing on ", "本次巡店已完成，主要聚焦在")
			for i, record := range info.ShowRecordList {
				info.Overview += record.TypeName
				if i < len(info.ShowRecordList)-1 {
					info.Overview += helper.RuleThree(req.Language == "en", ",", "、")
				}
			}
			info.Overview += helper.RuleThree(req.Language == "en", " Below is a summary of the report contents.", "方面的巡访报告内容概览")
		}
	}
}
