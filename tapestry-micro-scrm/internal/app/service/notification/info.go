package notificationSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io/ioutil"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/dao/notification"
	"micro-scrm/internal/app/http/request"
	emailReq "micro-scrm/internal/app/http/request/email"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	notificationResp "micro-scrm/internal/app/http/response/notification"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"strings"
	"time"
)

type InfoSvc struct {
	service.Base
	infoDao *notificationDao.InfoDao
}

func NewInfoSvc() *InfoSvc {
	return &InfoSvc{
		infoDao: notificationDao.NewInfoDao(),
	}
}

func (s *InfoSvc) Save(c ctx.Context, req notificationReq.InfoReq, memberId int64, name string) error {
	var err error
	id := app.UID()
	nowTime := time.Now().Unix()
	var notificationTitle, departmentName string
	var notificationType int32
	if req.Id > 0 { // 编辑
		saveInfo := map[string]interface{}{
			"theme_title":     req.ThemeTitle,
			"title":           req.Title,
			"until_name":      req.UntilName,
			"content":         req.Content,
			"perform":         req.Perform,
			"date_time":       req.DateTime,
			"contact_info":    helper.MarshalToString(req.ContactInfo),
			"attachment":      helper.MarshalToString(req.Attachment),
			"type":            req.Type,
			"department_name": req.DepartmentName,
			"admin_id":        req.AdminID,
			"member_id":       memberId,
			"member_name":     name,
			"status":          req.Status,
			"is_draft":        req.IsDraft,
			"submit_time":     int(nowTime),
			"page_url":        req.PageUrl,
			"pdf_url":         req.PdfUrl,
			"read_member_ids": "0",
			"updated_at":      nowTime,
		}
		// info := microScrmModel.NotificationInfoModel{
		// 	ThemeTitle:     req.ThemeTitle,
		// 	Title:          req.Title,
		// 	UntilName:      req.UntilName,
		// 	Content:        req.Content,
		// 	Perform:        req.Perform,
		// 	DateTime:       req.DateTime,
		// 	ContactInfo:    helper.MarshalToString(req.ContactInfo),
		// 	Attachment:     helper.MarshalToString(req.Attachment),
		// 	Type:           req.Type,
		// 	DepartmentName: req.DepartmentName,
		// 	AdminID:        req.AdminID,
		// 	MemberID:       memberId,
		// 	MemberName:     name,
		// 	Status:         req.Status,
		// 	IsDraft:        req.IsDraft,
		// 	SubmitTime:     int(nowTime),
		// 	PageURL:        req.PageUrl,
		// 	ReadMemberIds:  "0", // 只要再次编辑就清空已读
		// }
		if req.IsDraft == microScrmModel.NotificationNoDraft { // 保存提交之后，传讯状态变更为待审核状态
			saveInfo["status"] = microScrmModel.NotificationStatusChecking
		}
		// info.BaseModel.UpdatedAt = nowTime
		detail, _ := s.infoDao.GetOne(c, req.EnterpriseId, "id=?", req.Id)
		if detail.AdminID != req.AdminID {
			return errors.New("不能修改非自己创建的部门传讯")
		}
		if !helper.InArray(detail.Status, []int32{microScrmModel.NotificationStatusWait, microScrmModel.NotificationStatusRecall, microScrmModel.NotificationStatusFail}) {
			return errors.New("传讯状态为:" + microScrmModel.NotificationInfoModel{}.NotificationStatusStr(saveInfo["status"].(int32)) + "时不可修改")
		}
		err = s.infoDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": req.Id,
		}, saveInfo)
		notificationTitle = saveInfo["title"].(string)
		notificationType = saveInfo["type"].(int32)
		departmentName = saveInfo["department_name"].(string)
		id = req.Id
	} else {
		info := microScrmModel.NotificationInfoModel{
			ThemeTitle:     req.ThemeTitle,
			Title:          req.Title,
			UntilName:      req.UntilName,
			Content:        req.Content,
			Perform:        req.Perform,
			DateTime:       req.DateTime,
			ContactInfo:    helper.MarshalToString(req.ContactInfo),
			Attachment:     helper.MarshalToString(req.Attachment),
			Type:           req.Type,
			DepartmentName: req.DepartmentName,
			AdminID:        req.AdminID,
			MemberID:       memberId,
			MemberName:     name,
			Status:         req.Status,
			IsDraft:        req.IsDraft,
			SubmitTime:     int(nowTime),
			PageURL:        req.PageUrl,
			PdfUrl:         req.PdfUrl,
			ReadMemberIds:  "0", // 只要再次编辑就清空已读
		}
		info.BaseModel = models.BaseModel{
			ID:        id,
			CreatedAt: nowTime,
			UpdatedAt: nowTime,
			DeletedAt: 0,
		}
		if req.IsDraft == microScrmModel.NotificationNoDraft { // 保存提交之后，传讯状态变更为待审核状态
			info.Status = microScrmModel.NotificationStatusChecking
		}
		err = s.infoDao.CreateData(c, req.EnterpriseId, info)
		notificationTitle = info.Title
		notificationType = info.Type
		departmentName = info.DepartmentName
	}
	if err != nil {
		return err
	}
	if req.IsDraft == microScrmModel.NotificationNoDraft { // 部门传讯提交了之后,需要发送模板消息
		title := "有部门提交新的传讯~"
		uri := fmt.Sprintf("communication-detail/ro?id=%d", id)
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "传讯标题",
			"value":   notificationTitle,
		})
		content = append(content, map[string]string{
			"keyname": "传讯类型",
			"value":   microScrmModel.NotificationInfoModel{}.NotificationTypeStr(notificationType),
		})
		content = append(content, map[string]string{
			"keyname": "提交部门",
			"value":   departmentName,
		})
		submitTime := helper.GetFormatDateByParams(nowTime)
		content = append(content, map[string]string{
			"keyname": "提交时间",
			"value":   submitTime,
		})
		ro, err := s.GetAllRo(c, req.EnterpriseId)
		if err != nil {
			return err
		}
		go func() {
			_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: req.EnterpriseId,
				},
				Title:   title,
				Uri:     uri,
				Party:   nil,
				UserIds: ro,
				Content: content,
				AppSign: qyWork.Notification,
			})
			if err != nil {
				app.Logger().Error(c).Msgf("部门传讯通知消息发送失败：%+v", err)
			}
		}()
	}
	return err
}

func (s *InfoSvc) List(c ctx.Context, req notificationReq.SearchListReq, infoManager notificationResp.ManagerInfo) (int64, []notificationResp.ListRes, error) {
	result := make([]notificationResp.ListRes, 0)
	if infoManager.IsRo != notificationReq.IsRo && infoManager.IsFT != notificationReq.IsFT {
		return 0, result, nil
	}
	total, list, err := s.infoDao.GetList(c, req.EnterpriseId, req, infoManager)
	if err != nil {
		return total, result, err
	}
	var infoIds []int64
	for _, model := range list {
		infoIds = append(infoIds, model.ID)
	}
	contentList, _ := notificationDao.NewInfoMergeContentDao().GetMany(c, req.EnterpriseId, "info_id in ?", infoIds)

	var mergeIds []int64
	var infoMergeIdMap = make(map[int64]int64)
	for _, content := range contentList {
		mergeIds = append(mergeIds, content.MergeId)
		infoMergeIdMap[content.InfoId] = content.MergeId
	}
	mergeList, _ := notificationDao.NewInfoMergeDao().GetMany(c, req.EnterpriseId, "id in ?", mergeIds)

	var mergeStatusMap = make(map[int64]int)
	for _, merge := range mergeList {
		mergeStatusMap[merge.ID] = merge.Status
	}
	for _, model := range list {
		isRead := true // 默认已读
		if req.MemberId > 0 && helper.InArray(int(model.Status), []int{microScrmModel.NotificationStatusChecking, microScrmModel.NotificationStatusFail}) {
			memberIds := strings.Split(model.ReadMemberIds, ",")
			if !helper.InArray(helper.Int64ToString(req.MemberId), memberIds) {
				if req.IsAll != notificationReq.IsAll && model.Status == microScrmModel.NotificationStatusFail && model.MemberID == req.MemberId {
					isRead = false // 未读
				} else if req.IsAll == notificationReq.IsAll && model.Status == microScrmModel.NotificationStatusChecking {
					isRead = false // 未读
				}
			}
		}
		relStatus := -1 // 未引用
		if mergeId, ok := infoMergeIdMap[model.ID]; ok {
			if _, ok := mergeStatusMap[mergeId]; ok {
				relStatus = int(model.MergeRelease)
			}
		}
		item := notificationResp.ListRes{
			ID:             model.ID,
			Title:          model.Title,
			Type:           model.Type,
			DepartmentName: model.DepartmentName,
			SubmitName:     model.MemberName,
			UpdateAt:       model.UpdatedAt,
			ReviewName:     model.ReviewName,
			ReviewTime:     model.ReviewTime,
			Status:         model.Status,
			IsDraft:        model.IsDraft,
			RelStatus:      relStatus,
			IsRead:         isRead,
		}
		result = append(result, item)
	}
	return total, result, err
}

func (s *InfoSvc) Detail(c ctx.Context, params notificationReq.SearchId) (notificationResp.InfoRes, error) {
	var res notificationResp.InfoRes
	info, err := s.infoDao.GetOne(c, params.EnterpriseId, "id=?", params.Id)
	if err != nil {
		return res, err
	}
	//获取查看人数和查看次数
	readInfo, err := notificationDao.NewInfoMemberDao().GetInfoReadNum(c, params.EnterpriseId, params.Id)
	if err != nil {
		return notificationResp.InfoRes{}, err
	}
	contactInfo := make([]notificationResp.ContactList, 0)
	attachmentInfo := make([]notificationResp.AttachmentInfo, 0)
	helper.Unmarshal(info.ContactInfo, &contactInfo)
	helper.Unmarshal(info.Attachment, &attachmentInfo)
	res = notificationResp.InfoRes{
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
		PdfUrl:         info.PdfUrl,
		TypeStr:        microScrmModel.NotificationInfoModel{}.NotificationTypeStr(info.Type),
		StatusStr:      microScrmModel.NotificationInfoModel{}.NotificationStatusStr(info.Status),
		InfoViewNumResp: notificationResp.InfoViewNumResp{
			ViewPeopleNum: readInfo.ViewPeopleNum,
			ViewNum:       readInfo.ViewNum,
		},
	}
	if params.MemberId > 0 {
		// 处理为已读
		if helper.InArray(int(info.Status), []int{microScrmModel.NotificationStatusChecking, microScrmModel.NotificationStatusFail}) {
			s.reviewRead(c, params.EnterpriseId, params.MemberId, info)
		}
		// 处理h5查看次数
		if info.MergeRelease == microScrmModel.NotificationMergeSend {
			s.readInfo(c, params)
		}
	}
	return res, err
}
func (s *InfoSvc) reviewRead(c ctx.Context, enterpriseId, memberId int64, detail microScrmModel.NotificationInfoModel) {
	memberIdStr := helper.Int64ToString(memberId)
	readMemberIds := strings.Split(detail.ReadMemberIds, ",")
	if !helper.InArray(memberIdStr, readMemberIds) {
		_ = s.infoDao.UpdateDataOnly(c, enterpriseId, map[string]interface{}{
			"id": detail.ID,
		}, map[string]interface{}{
			"read_member_ids": strings.Trim(strings.Join(append(readMemberIds, memberIdStr), ","), ","),
		})
	}
}

func (s *InfoSvc) GetDepartments(c ctx.Context, enterpriseId int64) ([]microScrmModel.NotificationDepartmentModel, error) {
	return s.infoDao.GetDepartments(c, enterpriseId)
}

func (s *InfoSvc) Copy(c ctx.Context, enterpriseId, id, adminId int64) (int64, error) {
	newId := app.UID()
	one, err := s.infoDao.GetOne(c, enterpriseId, "id=?", id)
	if err != nil {
		return newId, err
	}
	if one.AdminID != adminId {
		return newId, errors.New("不能复制非本人提交的部门传讯")
	}
	nowTime := time.Now().Unix()
	info := microScrmModel.NotificationInfoModel{
		ThemeTitle:     one.ThemeTitle,
		Title:          one.Title,
		UntilName:      one.UntilName,
		Content:        one.Content,
		Perform:        one.Perform,
		DateTime:       one.DateTime,
		ContactInfo:    one.ContactInfo,
		Attachment:     one.Attachment,
		Type:           one.Type,
		DepartmentName: one.DepartmentName,
		AdminID:        one.AdminID,
		MemberID:       one.MemberID,
		MemberName:     one.MemberName,
	}
	info.BaseModel = models.BaseModel{
		ID:        newId,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
		DeletedAt: 0,
	}
	err = s.infoDao.CreateData(c, enterpriseId, info)
	return newId, err
}

func (s *InfoSvc) Recall(c ctx.Context, params notificationReq.SearchId) error {
	one, err := s.infoDao.GetOne(c, params.EnterpriseId, "id=?", params.Id)
	if err != nil {
		return err
	}
	if (params.AdminID > 0 && one.AdminID != params.AdminID) || (params.MemberId > 0 && one.MemberID != params.MemberId) {
		return errors.New("非本人提交的部门传讯不能撤回")
	}
	if one.Status != microScrmModel.NotificationStatusChecking {
		return errors.New("只有待审核才能撤回")
	}
	return s.infoDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
		"id": params.Id,
	}, map[string]interface{}{
		"status": microScrmModel.NotificationStatusRecall,
	})
}

func (s *InfoSvc) Check(c ctx.Context, req notificationReq.InfoReq, managerInfo notificationResp.ManagerInfo) error {
	if !helper.InArray(req.Status, []int32{microScrmModel.NotificationStatusFail, microScrmModel.NotificationStatusSuccess}) { // 部门传讯审核了之后,需要发送模板消息
		return errors.New("参数错误")
	}
	info, err := s.infoDao.GetOne(c, req.EnterpriseId, "id=?", req.Id)
	if err != nil {
		return err
	}
	if info.Status != microScrmModel.NotificationStatusChecking {
		return errors.New("该传讯状态已变更，暂不可审核")
	}
	key := fmt.Sprintf("check_notification_:%d", req.Id)
	defer s.Redis().Del(c, key)
	if !s.Redis().SetNX(c, key, 1, 5*time.Second).Val() { // 提交审核加下锁，防止覆盖
		return errors.New("该传讯已被其他RO审核")
	}
	memberInfo, err := membersDao.NewMembers().GetMemberInfoById(c, req.EnterpriseId, managerInfo.MemberId)
	if err != nil {
		return err
	}
	nowTime := time.Now().Unix()
	if req.IsH5 {
		err = s.infoDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": req.Id,
		}, map[string]interface{}{
			"status":           req.Status,
			"review_name":      memberInfo.Name,
			"review_time":      nowTime,
			"review_member_id": memberInfo.ID,
			"review_admin_id":  managerInfo.AdminId,
			"comment":          req.Comment,
			"read_member_ids":  "",
			"updated_at":       nowTime,
		})
	} else {
		checkInfo := map[string]interface{}{
			"theme_title":      req.ThemeTitle,
			"title":            req.Title,
			"until_name":       req.UntilName,
			"content":          req.Content,
			"perform":          req.Perform,
			"date_time":        req.DateTime,
			"contact_info":     helper.MarshalToString(req.ContactInfo),
			"attachment":       helper.MarshalToString(req.Attachment),
			"type":             req.Type,
			"department_name":  req.DepartmentName,
			"review_name":      memberInfo.Name,
			"review_time":      int(nowTime),
			"review_admin_id":  managerInfo.AdminId,
			"review_member_id": memberInfo.ID,
			"status":           req.Status,
			"page_url":         req.PageUrl,
			"comment":          req.Comment,
			"read_member_ids":  "0",
			"pdf_url":          req.PdfUrl,
			"updated_at":       nowTime,
		}
		err = s.infoDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": req.Id,
		}, checkInfo)
	}
	if err != nil {
		return err
	}
	title := "您有1条传讯审核通过了~"
	reviewResult := "已通过"
	if req.Status == microScrmModel.NotificationStatusFail {
		title = "您有1条传讯审核未通过~"
		reviewResult = "未通过"
	}
	uri := fmt.Sprintf("communication-detail/ft?id=%d", req.Id)
	content := make([]map[string]string, 0)
	content = append(content, map[string]string{
		"keyname": "传讯标题",
		"value":   info.Title,
	})
	content = append(content, map[string]string{
		"keyname": "审核结果",
		"value":   reviewResult,
	})
	reviewTime := helper.GetFormatDateByParams(nowTime)
	content = append(content, map[string]string{
		"keyname": "审核时间",
		"value":   reviewTime,
	})
	// 获取创建人userid
	memberInfo, err = membersDao.NewMembers().GetMemberInfoById(c, req.EnterpriseId, info.MemberID)
	if err != nil {
		return err
	}
	go func() {
		_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: req.EnterpriseId,
			},
			Title:   title,
			Uri:     uri,
			Party:   nil,
			UserIds: []string{memberInfo.Userid},
			Content: content,
			AppSign: qyWork.Notification,
		})
		if err != nil {
			app.Logger().Error(c).Msgf("部门传讯通知消息发送失败：%+v", err)
		}
	}()
	return err
}

func (s *InfoSvc) GetManagerInfo(c ctx.Context, params notificationReq.AdminManager) (notificationResp.ManagerInfo, error) {
	var result notificationResp.ManagerInfo
	queryStr := "admin_id = ?"
	args := params.AdminID
	if params.MemberId > 0 {
		queryStr = "member_id = ?"
		args = params.MemberId
	}
	info, err := s.infoDao.GetManagerInfo(c, params.EnterpriseId, queryStr, args)
	if err != nil {
		return result, err
	}
	result.IsFT = info.IsFt
	result.IsRo = info.IsRo
	result.MemberId = info.MemberID
	result.AdminId = info.AdminID
	return result, err
}

func (s *InfoSvc) GetAllRo(c ctx.Context, enterpriseId int64) (userIds []string, err error) {

	list, err := s.infoDao.GetAllManager(c, enterpriseId)
	if err != nil {
		return nil, err
	}
	var memberIds []int64
	for _, model := range list {
		if model.IsRo == microScrmModel.NotificationIsRo {
			memberIds = append(memberIds, model.MemberID)
		}
	}
	if len(memberIds) > 0 {
		userIds, err = membersDao.NewMembers().GetUserIdsByMemberIDs(c, enterpriseId, memberIds)
		if err != nil {
			return nil, err
		}
	}
	return userIds, err
}

func (s *InfoSvc) Cancel(c ctx.Context, enterpriseId, id int64) error {
	return s.infoDao.UpdateData(c, enterpriseId, map[string]interface{}{
		"id": id,
	}, map[string]interface{}{
		"status": microScrmModel.NotificationStatusCancel,
	})
}

func (s *InfoSvc) IsRoRole(c ctx.Context, params notificationReq.AdminManager) (bool, error) {
	queryStr := "admin_id = ?"
	args := params.AdminID
	if params.MemberId > 0 {
		queryStr = "member_id = ?"
		args = params.MemberId
	}
	info, err := s.infoDao.GetManagerInfo(c, params.EnterpriseId, queryStr, args)
	if err != nil {
		return false, err
	}
	if info.IsRo == 1 {
		return true, err
	}
	return false, err

}

func (s *InfoSvc) ListReadNum(c ctx.Context, enterpriseId, MemberId int64, managerInfo notificationResp.ManagerInfo) (res notificationResp.ListReadNum, err error) {
	if managerInfo.IsFT == microScrmModel.NotificationIsFt || managerInfo.IsRo == microScrmModel.NotificationIsRo {
		infoList, err := s.infoDao.GetMany(c, enterpriseId, "status in ?", []int{microScrmModel.NotificationStatusChecking, microScrmModel.NotificationStatusFail})
		if err != nil {
			return notificationResp.ListReadNum{}, err
		}
		for _, model := range infoList {
			memberIds := strings.Split(model.ReadMemberIds, ",")
			if !helper.InArray(helper.Int64ToString(MemberId), memberIds) {
				if model.Status == microScrmModel.NotificationStatusChecking {
					res.InfoCheckingNum++ // 待审核未读数量
				}
				if model.MemberID == MemberId && model.Status == microScrmModel.NotificationStatusFail {
					res.InfoFailNum++ // 未通过未读数量
				}
			}
		}
	}
	// 整合传讯的未读数量
	res.MergeInfoNum, err = notificationDao.NewInfoMergeMemberDao().GetCount(c, enterpriseId, "member_id = ? AND view_num = 0", MemberId)

	return res, err
}

func (s *InfoSvc) GetTitle(c ctx.Context, req notificationReq.SendFileByEmailReq) (title string, err error) {
	if req.Id > 0 {
		detail, err := notificationDao.NewInfoDao().GetOne(c, req.EnterpriseId, "id=?", req.Id)
		if err != nil {
			return title, err
		}
		title = detail.Title
	} else if req.MergeId > 0 {
		detail, err := notificationDao.NewInfoMergeDao().GetOne(c, req.EnterpriseId, "id=?", req.MergeId)
		if err != nil {
			return title, err
		}
		title = detail.Title
	}
	return title, err
}

func (s *InfoSvc) GetEmailData(c ctx.Context, req notificationReq.SendFileByEmailReq) (emailData emailReq.SendMailReq, err error) {
	var infoList []microScrmModel.NotificationInfoModel
	var title, zipUrl string
	fileDir := `./storage/notification/`
	infoDao := notificationDao.NewInfoDao()
	infoMergeDao := notificationDao.NewInfoMergeDao()

	if req.Id > 0 {
		info, err := infoDao.GetOne(c, req.EnterpriseId, "id=?", req.Id)
		if err != nil {
			return emailReq.SendMailReq{}, err
		}
		if info.ZipUrl != "" {
			zipUrl = info.ZipUrl
		} else {
			infoList = []microScrmModel.NotificationInfoModel{info}
		}
		title = info.Title
	} else if req.MergeId > 0 {
		mergeInfo, err := infoMergeDao.GetOneData(c, req.EnterpriseId, "id=?", req.MergeId)
		if err != nil {
			return emailReq.SendMailReq{}, err
		}
		title = mergeInfo.Title
		if mergeInfo.Type == microScrmModel.NotificationTypeWeek {
			content, _ := ioutil.ReadFile(fileDir + `zip_email_tpl_week.html`)
			emailData.Content = fmt.Sprintf(string(content), title, mergeInfo.EmailLoopUrl, mergeInfo.EmailLoopUrl)
			return emailData, nil
		}
		if mergeInfo.ZipUrl != "" {
			zipUrl = mergeInfo.ZipUrl
		} else {
			var infoIdSlice []int64
			// 获取子传讯的pdf及附件
			for _, content := range mergeInfo.ContentList {
				infoIdSlice = append(infoIdSlice, content.InfoId)
			}
			infoList, err = infoDao.GetMany(c, req.EnterpriseId, "id in ?", infoIdSlice)
			if err != nil {
				return emailReq.SendMailReq{}, err
			}

		}
	}
	if zipUrl == "" && len(infoList) > 0 {
		zipUrl, err = s.formatZip(c, fileDir, title, infoList)
		if err != nil {
			return emailReq.SendMailReq{}, err
		}
		//保存zip的地址
		if req.Id > 0 {
			_ = infoDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.Id), map[string]interface{}{
				"zip_url": zipUrl,
			})

		} else if req.MergeId > 0 {
			_ = infoMergeDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", req.MergeId), map[string]interface{}{
				"zip_url": zipUrl,
			})
		}
	}
	//获取发送模版
	content, _ := ioutil.ReadFile(fileDir + `zip_email_tpl.html`)
	encrypt, err := common.AesEncrypt(zipUrl)
	if err != nil {
		return emailReq.SendMailReq{}, err
	}
	pageUrl := fmt.Sprintf("https://%s.coach.com.cn/web-admin/download?secret=%s",
		helper.RuleThree(app.Cfg().Env == "uat", "tapestry-web-uat", "we"),
		encrypt)

	emailData.Content = fmt.Sprintf(string(content), title, helper.Md5V(encrypt), pageUrl)
	return emailData, nil
}

// h5查看次数处理
func (s *InfoSvc) readInfo(c ctx.Context, params notificationReq.SearchId) {
	if params.MemberId > 0 {
		infoMemberDao := notificationDao.NewInfoMemberDao()
		one, _ := infoMemberDao.GetOne(c, params.EnterpriseId, "info_id=? AND member_id = ?", params.Id, params.MemberId)

		if one.ID > 0 {
			_ = notificationDao.NewInfoMemberDao().Increment(c, params.EnterpriseId, map[string]interface{}{
				"id": one.ID,
			}, "view_num", 1)
		} else {
			_ = notificationDao.NewInfoMemberDao().CreateData(c, params.EnterpriseId, microScrmModel.NotificationInfoMemberModel{
				BaseModel: models.CreateBase(app.UID()),
				InfoId:    params.Id,
				MemberId:  params.MemberId,
				ViewNum:   1,
			})
		}
	}
}

// 生成zip
func (s *InfoSvc) formatZip(c ctx.Context, infoDir, title string, infoList []microScrmModel.NotificationInfoModel) (string, error) {
	files := make(map[string]string)
	for _, info := range infoList {
		if info.PdfUrl == "" {
			continue
		}
		pdfLocalPath := infoDir + info.Title + `.pdf`
		err := helper.DownloadFileToLocal(info.PdfUrl, pdfLocalPath)
		if err != nil {
			return "", err
		}
		if len(infoList) == 1 {
			//只有一个传讯的时候不要目录
			files[pdfLocalPath] = info.Title + `.pdf`
		} else {
			files[pdfLocalPath] = info.Title + `/` + info.Title + `.pdf`
		}
		attachmentInfo := make([]notificationResp.AttachmentInfo, 0)
		helper.Unmarshal(info.Attachment, &attachmentInfo)
		for _, urlInfo := range attachmentInfo {
			if urlInfo.Url == "" {
				continue
			}
			LocalFilePath := infoDir + urlInfo.AttachmentName
			err = helper.DownloadFileToLocal(urlInfo.Url, LocalFilePath)
			if err != nil {
				return "", err
			}
			if len(infoList) == 1 {
				//只有一个传讯的时候不要目录
				files[LocalFilePath] = urlInfo.AttachmentName
			} else {
				files[LocalFilePath] = info.Title + `/` + urlInfo.AttachmentName
			}
		}

	}
	zipFileName := infoDir + `/` + title + `.zip`
	err := helper.CreateZip(zipFileName, files)
	if err != nil {
		return "", err
	}
	zipUrl, err := helper.UploadFileToCos(c, zipFileName, true)
	if err != nil {
		return "", err
	}

	return zipUrl, nil
}
