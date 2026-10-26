package notificationSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/notification"
	"micro-scrm/internal/app/http/request"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"sort"
	"strings"
	"time"
)

type InfoMergeSvc struct {
	service.Base
	infoMergeDao        *notificationDao.InfoMergeDao
	infoMergeContentDao *notificationDao.InfoMergeContentDao
}

func NewInfoMergeSvc() *InfoMergeSvc {
	return &InfoMergeSvc{
		infoMergeDao:        notificationDao.NewInfoMergeDao(),
		infoMergeContentDao: notificationDao.NewInfoMergeContentDao(),
	}
}

func (s *InfoMergeSvc) Save(c ctx.Context, params notificationReq.InfoMergeReq) error {
	memberInfo, err := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseId, params.Mobile)
	if err != nil {
		return err
	}
	var crowdIdStr string
	for _, crowdId := range params.CrowdIds {
		crowdIdStr += helper.Int64ToString(crowdId) + ","
	}
	crowdIdStr = strings.TrimRight(crowdIdStr, ",")
	var mergeId int64
	if params.Id > 0 {
		info, err := s.infoMergeDao.GetOne(c, params.EnterpriseId, "id=?", params.Id)
		if err != nil || info.ID == 0 {
			return errors.New("找不到该数据")
		}
		err = s.infoMergeDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"id": params.Id,
		}, map[string]interface{}{
			"type":           params.Type,
			"week_num":       params.WeekNum,
			"title":          params.Title,
			"cover":          params.Cover,
			"crowd_type":     params.CrowdType,
			"crowd_ids":      crowdIdStr,
			"send_type":      params.SendType,
			"send_time":      params.SendTime,
			"admin_id":       params.AdminId,
			"member_id":      memberInfo.Id,
			"member_name":    memberInfo.Name,
			"send_ro_name":   memberInfo.Name,
			"updated_at":     time.Now().Unix(),
			"pdf_url":        "",
			"email_loop_url": params.EmailLoopUrl,
		})
		if err != nil {
			return err
		}
		err = s.infoMergeContentDao.UnscopedDeleteData(c, params.EnterpriseId, map[string]interface{}{
			"merge_id": params.Id,
		})
		if err != nil {
			return err
		}
		err = s.saveInfoMergeContent(c, params.EnterpriseId, params.Id, params.Content)
		if err != nil {
			return err
		}
		mergeId = params.Id
	} else {
		mergeId = app.UID()
		err := s.infoMergeDao.CreateData(c, params.EnterpriseId, microScrmModel.NotificationInfoMergeModel{
			BaseModel:    models.CreateBase(mergeId),
			WeekNum:      params.WeekNum,
			Type:         params.Type,
			Title:        params.Title,
			Cover:        params.Cover,
			CrowdType:    params.CrowdType,
			CrowdIds:     crowdIdStr,
			SendType:     params.SendType,
			SendTime:     params.SendTime,
			AdminID:      params.AdminId,
			MemberID:     memberInfo.Id,
			MemberName:   memberInfo.Name,
			SendRoName:   memberInfo.Name,
			EmailLoopUrl: params.EmailLoopUrl,
		})
		if err != nil {
			return err
		}
		err = s.saveInfoMergeContent(c, params.EnterpriseId, mergeId, params.Content)
		if err != nil {
			return err
		}
	}
	// 即刻发布
	if params.SendType == microScrmModel.NotificationMergeNowSend {
		go func() {
			_ = s.Release(c, notificationReq.InfoMergeDetailReq{
				EnterpriseBaseReq: params.EnterpriseBaseReq,
				Id:                mergeId,
				Mobile:            memberInfo.Mobile,
			})
		}()
	}
	return nil
}

func (s *InfoMergeSvc) List(c ctx.Context, params notificationReq.InfoMergeListReq) (resList []notificationResp.InfoMergeListResp, total int64, err error) {
	list, total, err := s.infoMergeDao.GetPagination(c, params)
	if err != nil {
		return nil, 0, err
	}
	var mergeIds []int64
	for _, mergeInfo := range list {
		if mergeInfo.Status == microScrmModel.NotificationMergeSend {
			mergeIds = append(mergeIds, mergeInfo.ID)
		}
	}
	var numMap = make(map[int64]notificationResp.InfoMergeMemberNumResp)
	if len(mergeIds) > 0 {
		numMap, err = notificationDao.NewInfoMergeMemberDao().GetListNum(c, params.EnterpriseId, mergeIds)
		if err != nil {
			return nil, 0, err
		}
	}

	for _, mergeInfo := range list {
		tmpNum := notificationResp.InfoMergeMemberNumResp{
			MergeId: mergeInfo.ID,
		}
		if val, ok := numMap[mergeInfo.ID]; ok {
			tmpNum = val
		}
		resList = append(resList, notificationResp.InfoMergeListResp{
			NotificationInfoMergeModel: mergeInfo,
			InfoMergeMemberNumResp:     tmpNum,
		})
	}
	return resList, total, nil
}

func (s *InfoMergeSvc) Detail(c ctx.Context, params notificationReq.InfoMergeDetailReq) (notificationResp.InfoMergeDetailResp, error) {
	infoMerge, err := s.infoMergeDao.GetOneData(c, params.EnterpriseId, "id=?", params.Id)
	if err != nil {
		return notificationResp.InfoMergeDetailResp{}, err
	}
	numMap, err := notificationDao.NewInfoMergeMemberDao().GetListNum(c, params.EnterpriseId, []int64{params.Id})
	if err != nil {
		return notificationResp.InfoMergeDetailResp{}, err
	}
	memberInfo, _ := membersDao.NewMembers().GetMemberInfoById(c, params.EnterpriseId, infoMerge.MemberID)

	resp := notificationResp.InfoMergeDetailResp{
		Id:                     infoMerge.ID,
		WeekNum:                infoMerge.WeekNum,
		Title:                  infoMerge.Title,
		Type:                   infoMerge.Type,
		Cover:                  infoMerge.Cover,
		SendType:               infoMerge.SendType,
		SendTime:               infoMerge.SendTime,
		Status:                 infoMerge.Status,
		CrowdType:              infoMerge.CrowdType,
		PdfUrl:                 infoMerge.PdfUrl,
		EmailLoopUrl:           infoMerge.EmailLoopUrl,
		MemberName:             memberInfo.Name,
		CrowdList:              nil,
		ContentList:            nil,
		InfoMergeMemberNumResp: numMap[infoMerge.ID],
	}
	if len(infoMerge.ContentList) > 0 {
		var infoIdSlice []int64
		for _, row := range infoMerge.ContentList {
			infoIdSlice = append(infoIdSlice, row.InfoId)
		}
		infoList, err := notificationDao.NewInfoDao().GetMany(c, params.EnterpriseId, "id in ?", infoIdSlice)
		if err != nil {
			return notificationResp.InfoMergeDetailResp{}, err
		}
		var infoMap = make(map[int64]microScrmModel.NotificationInfoModel)
		for _, info := range infoList {
			infoMap[info.ID] = info
		}
		var mergeGroupIdsContentMap = make(map[int64]notificationResp.MergeContentResp)
		for i, row := range infoMerge.ContentList {
			var tmpInfoList []notificationResp.InfoRes
			if info, ok := infoMap[row.InfoId]; ok {
				attachmentInfo := make([]notificationResp.AttachmentInfo, 0)
				contactInfo := make([]notificationResp.ContactList, 0)
				helper.Unmarshal(infoMap[row.InfoId].Attachment, &attachmentInfo)
				helper.Unmarshal(infoMap[row.InfoId].ContactInfo, &contactInfo)
				tmpInfoList = append(tmpInfoList, notificationResp.InfoRes{
					Id:             info.ID,
					ThemeTitle:     info.ThemeTitle,
					Title:          info.Title,
					UntilName:      info.UntilName,
					Content:        info.Content,
					Perform:        info.Perform,
					DateTime:       info.DateTime,
					ContactInfo:    contactInfo,
					Attachment:     attachmentInfo,
					Type:           info.Type,
					DepartmentName: info.DepartmentName,
					AdminID:        info.AdminID,
					Status:         info.Status,
					ReviewTime:     info.ReviewTime,
					ReviewName:     info.ReviewName,
					SubmitTime:     info.SubmitTime,
					MemberID:       info.MemberID,
					MemberName:     info.MemberName,
					Comment:        info.Comment,
					PageUrl:        info.PageURL,
					IsDraft:        info.IsDraft,
					TypeStr:        microScrmModel.NotificationInfoModel{}.NotificationTypeStr(info.Type),
					StatusStr:      microScrmModel.NotificationInfoModel{}.NotificationStatusStr(info.Status),
				})
			}
			if contentInfo, ok := mergeGroupIdsContentMap[row.GroupId]; ok {
				mergeGroupIdsContentMap[row.GroupId] = notificationResp.MergeContentResp{
					DepartmentName: row.DepartmentName,
					Sort:           i,
					InfoList:       append(contentInfo.InfoList, tmpInfoList...),
				}
			} else {
				mergeGroupIdsContentMap[row.GroupId] = notificationResp.MergeContentResp{
					DepartmentName: row.DepartmentName,
					Sort:           i,
					InfoList:       tmpInfoList,
				}
			}
		}
		for _, contentResp := range mergeGroupIdsContentMap {
			resp.ContentList = append(resp.ContentList, contentResp)
		}
		sort.Slice(resp.ContentList, func(i, j int) bool {
			return resp.ContentList[i].Sort < resp.ContentList[j].Sort
		})
	}

	if infoMerge.CrowdIds != "" {
		var crowdIdSlice []int64
		for _, crowdId := range strings.Split(infoMerge.CrowdIds, ",") {
			crowdIdSlice = append(crowdIdSlice, int64(helper.StringToInt(crowdId)))
		}
		if infoMerge.CrowdType == microScrmModel.NotificationCrowdWhite {
			crowdList, err := notificationDao.NewCrowdDao().GetMany(c, params.EnterpriseId, "id in ?", crowdIdSlice)
			if err != nil {
				return notificationResp.InfoMergeDetailResp{}, err
			}
			for _, info := range crowdList {
				resp.CrowdList = append(resp.CrowdList, notificationResp.MergeCrowdResp{
					Id:   info.ID,
					Name: info.Name,
				})
			}
		} else if infoMerge.CrowdType == microScrmModel.NotificationMember {
			memberList, err := membersDao.NewMembers().GetAllById(c, params.EnterpriseId, crowdIdSlice)
			if err != nil {
				return notificationResp.InfoMergeDetailResp{}, err
			}
			for _, info := range memberList {
				var dpId int64
				if len(info.MemberDepartments) > 0 {
					dpId = info.MemberDepartments[0].DepartmentID
				}
				resp.CrowdList = append(resp.CrowdList, notificationResp.MergeCrowdResp{
					Id:           info.ID,
					Name:         info.Name,
					DepartmentId: dpId,
				})
			}
		}
	}
	// 具体发给了谁
	if infoMerge.Status == microScrmModel.NotificationMergeSend {
		memberIds, _ := notificationDao.NewInfoMergeMemberDao().GetMemberIds(c, params.EnterpriseId, "merge_id = ?", params.Id)
		memberList, _ := membersDao.NewMembers().GetUnscopedMemberInfoByIds(c, params.EnterpriseId, memberIds)
		for _, info := range memberList {
			resp.SendNameList = append(resp.SendNameList, info.Name)
		}
	}
	// 处理查看次数
	_ = s.readMergeInfo(c, params)
	return resp, nil
}

func (s *InfoMergeSvc) saveInfoMergeContent(c ctx.Context, enterpriseId int64, mergeId int64, contentList []notificationReq.InfoMergeContentReq) error {
	var dataArr []microScrmModel.NotificationInfoMergeContentModel
	for _, content := range contentList {
		groupId := app.UID()
		for _, infoId := range content.InfoIds {
			dataArr = append(dataArr, microScrmModel.NotificationInfoMergeContentModel{
				BaseModel:      models.CreateBase(app.UID()),
				GroupId:        groupId,
				MergeId:        mergeId,
				DepartmentName: content.DepartmentName,
				InfoId:         infoId,
			})
		}
	}
	err := s.infoMergeContentDao.CreateInBatchesData(c, enterpriseId, dataArr)
	if err != nil {
		return err
	}
	return nil
}

// Release 发布传讯
func (s *InfoMergeSvc) Release(c ctx.Context, params notificationReq.InfoMergeDetailReq) error {
	nowUnix := time.Now().Unix()
	var releaseName string
	if params.Mobile != "" {
		memberInfo, _ := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseId, params.Mobile)
		releaseName = memberInfo.Name
	} else if params.MemberId > 0 {
		memberInfo, _ := membersDao.NewMembers().GetMemberInfoById(c, params.EnterpriseId, params.MemberId)
		releaseName = memberInfo.Name
	}
	infoMerge, err := s.infoMergeDao.GetOneData(c, params.EnterpriseId, "id=?", params.Id)
	if err != nil {
		return err
	}
	if infoMerge.Status != microScrmModel.NotificationMergeNotSend {
		return errors.New("传讯状态异常")
	}
	var memberIds []int64
	var crowdIdSlice []int64
	for _, crowdId := range strings.Split(infoMerge.CrowdIds, ",") {
		crowdIdSlice = append(crowdIdSlice, int64(helper.StringToInt(crowdId)))
	}
	if infoMerge.CrowdType == microScrmModel.NotificationMember {
		// 直接选导购，存的就是导购id集合
		memberIds = crowdIdSlice
	} else if infoMerge.CrowdType == microScrmModel.NotificationCrowdWhite {
		// 选的是白名单，存的就是白名单的id
		crowdList, err := notificationDao.NewCrowdDao().GetManyData(c, params.EnterpriseId, "id in ?", crowdIdSlice)
		if err != nil {
			return err
		}
		for _, info := range crowdList {
			for _, crowdMember := range info.DataList {
				memberIds = append(memberIds, crowdMember.MemberId)
			}
		}
	}
	memberIds = helper.ArrayUniq(memberIds)
	// 组装插入数据
	var addData []microScrmModel.NotificationInfoMergeMemberModel
	for _, memberId := range memberIds {
		addData = append(addData, microScrmModel.NotificationInfoMergeMemberModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: nowUnix,
				UpdatedAt: nowUnix,
			},
			MergeId:  params.Id,
			MemberId: memberId,
		})
	}
	err = notificationDao.NewInfoMergeMemberDao().CreateInBatchesData(c, params.EnterpriseId, addData)
	if err != nil {
		return err
	}
	err = s.infoMergeDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.Id,
	}, map[string]interface{}{
		"status":       microScrmModel.NotificationMergeSend,
		"send_ro_name": releaseName,
		"send_time":    nowUnix,
	})
	var infoIdSlice []int64
	for _, row := range infoMerge.ContentList {
		infoIdSlice = append(infoIdSlice, row.InfoId)
	}
	_ = notificationDao.NewInfoDao().UpdateDataOnly(c, params.EnterpriseId, map[string]interface{}{
		"id": infoIdSlice,
	}, map[string]interface{}{
		"merge_release": microScrmModel.NotificationMergeSend,
	})
	if err != nil {
		return err
	}
	// 获取要发送的导购userid
	userIds, err := membersDao.NewMembers().GetUserIdsByMemberIDs(c, params.EnterpriseId, memberIds)
	if err != nil {
		return err
	}
	content := make([]map[string]string, 0)
	content = append(content, map[string]string{
		"keyname": "传讯标题",
		"value":   infoMerge.Title,
	})
	content = append(content, map[string]string{
		"keyname": "传讯时间",
		"value":   helper.GetFormatDateByParams(infoMerge.SendTime),
	})

	// 定义每组的大小
	groupSize := 1000
	uri := fmt.Sprintf("corporate-detail?id=%d", params.Id)

	// 分组并执行函数
	for i := 0; i < len(userIds); i += groupSize {
		end := i + groupSize
		if end > len(userIds) {
			end = len(userIds)
		}
		tmpUserIds := userIds[i:end]
		// 发布通知
		_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: params.EnterpriseId,
			},
			Title:   "您有1条企业传讯~",
			Uri:     uri,
			UserIds: tmpUserIds,
			Content: content,
			AppSign: qyWork.Notification,
		})
	}

	if err != nil {
		log.Info(c).Msgf("传讯模版消息发送失败:%d,错误原因：%+v", params.Id, err)
	}
	return nil
}

func (s *InfoMergeSvc) ListShow(c ctx.Context, params notificationReq.InfoMergeListReq) (list []notificationResp.InfoMergeListByH5Resp, total int64, err error) {
	list, total, err = s.infoMergeDao.GetShowList(c, params)
	if err != nil {
		return nil, 0, err
	}
	var mergeIds []int64
	for i, mergeInfo := range list {
		if mergeInfo.ViewNum > 0 {
			list[i].IsRead = true
		}
		mergeIds = append(mergeIds, mergeInfo.ID)
	}
	if len(mergeIds) > 0 {
		numMap, err := notificationDao.NewInfoMergeMemberDao().GetListNum(c, params.EnterpriseId, mergeIds)
		if err != nil {
			return nil, 0, err
		}
		for i, mergeInfo := range list {
			if val, ok := numMap[mergeInfo.ID]; ok {
				list[i].ViewNum = val.ViewNum //这个是总的viewNum
				list[i].InfoMergeMemberNumResp = val
			}
		}
	}
	return list, total, nil
}

func (s *InfoMergeSvc) SavePdf(c ctx.Context, params notificationReq.SavePdfReq) (err error) {
	if params.MergeId > 0 {
		return s.infoMergeDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"id": params.MergeId,
		}, map[string]interface{}{
			"pdf_url": params.PdfUrl,
		})
	} else if params.InfoId > 0 {
		return notificationDao.NewInfoDao().UpdateDataColumns(c, params.EnterpriseId, map[string]interface{}{
			"id": params.InfoId,
		}, map[string]interface{}{
			"pdf_url": params.PdfUrl,
		})
	}
	return nil
}

func (s *InfoMergeSvc) readMergeInfo(c ctx.Context, params notificationReq.InfoMergeDetailReq) (err error) {
	if params.MemberId > 0 {
		return notificationDao.NewInfoMergeMemberDao().Increment(c, params.EnterpriseId, map[string]interface{}{
			"merge_id":  params.Id,
			"member_id": params.MemberId,
		}, "view_num", 1)
	}
	return nil
}

// TimedRelease 定时发布
func (s *InfoMergeSvc) TimedRelease(c ctx.Context, enterpriseId int64) error {
	list, err := s.infoMergeDao.GetMany(c, enterpriseId, "send_type = ? AND status = ? AND send_time < ? AND send_time > 0", microScrmModel.NotificationMergeTimeSend, microScrmModel.NotificationMergeNotSend, time.Now().Unix())
	if err != nil {
		return err
	}
	for _, model := range list {
		err = s.Release(c, notificationReq.InfoMergeDetailReq{
			EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: enterpriseId},
			Id:                model.ID,
			MemberId:          model.MemberID,
		})
		if err != nil {
			log.Info(c).Msgf("定时发布传讯失败:%d,错误原因：%+v", model.ID, err)
		}
	}
	return nil
}
