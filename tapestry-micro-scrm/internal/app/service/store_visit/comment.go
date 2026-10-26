package storeVisitSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	storeVisitDao "micro-scrm/internal/app/dao/store_visit"
	"micro-scrm/internal/app/http/request"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitResp "micro-scrm/internal/app/http/response/store_visit"
	"micro-scrm/internal/app/models"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
	"micro-scrm/internal/app/service"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
)

type CommentSvc struct {
	service.Base
	visitCommentDao *storeVisitDao.CommentDao
	visitMediaDao   *storeVisitDao.MediaDao
}

func NewCommentSvc() *CommentSvc {
	return &CommentSvc{
		visitCommentDao: storeVisitDao.NewCommentDao(),
		visitMediaDao:   storeVisitDao.NewMediaDao(),
	}
}

func (s *CommentSvc) SaveComment(c ctx.Context, req storeVisitReq.SaveCommentReq) (commentId int64, err error) {
	var commentInfo storeVisitModel.StoreVisitCommentModel
	var msgMemberIds []int64
	if req.ReplyCommentId > 0 {
		commentInfo, err = s.visitCommentDao.GetOne(c, req.EnterpriseId, "id = ?", req.ReplyCommentId)
		if err != nil {
			return 0, err
		}
		msgMemberIds = append(msgMemberIds, commentInfo.MemberID)
	}
	commentId = app.UID()
	err = s.visitCommentDao.CreateData(c, req.EnterpriseId, storeVisitModel.StoreVisitCommentModel{
		BaseModel:     models.CreateBase(commentId),
		VisitID:       req.ID,
		MemberID:      req.MemberId,
		Content:       req.Content,
		ParentID:      helper.RuleThree(commentInfo.ParentID > 0, commentInfo.ParentID, req.ReplyCommentId),
		ReplyMemberID: commentInfo.MemberID,
	})
	if err != nil {
		return 0, err
	}
	if len(req.MediaList) > 0 {
		var addMedia = make([]storeVisitModel.StoreVisitMediaModel, 0)
		for _, media := range req.MediaList {
			addMedia = append(addMedia, storeVisitModel.StoreVisitMediaModel{
				BaseModel:     models.CreateBase(app.UID()),
				VisitID:       req.ID,
				MemberID:      req.MemberId,
				RelID:         commentId,
				RelType:       storeVisitModel.Comments,
				MediaType:     media.MediaType,
				URL:           media.Url,
				VideoImageUrl: media.VideoImageUrl,
				VideoSecond:   media.VideoSeconds,
			})
		}
		err = s.visitMediaDao.CreateInBatchesData(c, req.EnterpriseId, addMedia)
		if err != nil {
			return 0, err
		}
	}
	go func() {
		visitInfo, _ := storeVisitDao.NewMainDao().GetInfoPreloadCreator(c, req.VisitInfoReq)
		if len(visitInfo.MemberList) > 0 {
			msgMemberIds = append(msgMemberIds, visitInfo.MemberList[0].MemberID)
		}
		msgMemberIds = helper.ArrayDiff(msgMemberIds, []int64{req.MemberId})
		if len(msgMemberIds) > 0 {
			memberList, _ := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND id in ?", req.EnterpriseId, helper.ArrayUniq(msgMemberIds))
			if len(memberList) > 0 {
				storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
					"enterprise_id":        req.EnterpriseId,
					"wechat_department_id": visitInfo.ThirdId,
				})
				for _, member := range memberList {
					_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
						EnterpriseBaseReq: request.EnterpriseBaseReq{
							EnterpriseId: req.EnterpriseId,
						},
						Title: helper.RuleThree(member.ID == commentInfo.MemberID,
							helper.RuleThree(req.Language == "en", "New Reply Notification", "留言回复提醒"),
							helper.RuleThree(req.Language == "en", "New Comment Notification", "巡店留言提醒")),
						Uri:     fmt.Sprintf("store-visit-report?id=%d&comment_id=%d", req.ID, commentId),
						Party:   nil,
						UserIds: []string{member.Userid},
						Content: []map[string]string{
							{
								"keyname": helper.RuleThree(req.Language == "en", "Content", "内容"),
								"value": helper.RuleThree(req.Language == "en",
									"You've received a new "+helper.RuleThree(member.ID == commentInfo.MemberID, "reply", "comment")+".  Click to view details.",
									"你收到了一条新的留言"+helper.RuleThree(member.ID == commentInfo.MemberID, "回复", "")+"，点击查看详情。",
								),
							},
							{
								"keyname": helper.RuleThree(req.Language == "en", "Store", "巡店门店"),
								"value":   helper.RuleThree(req.Language == "en" && storeDepartmentInfo.InventoryName != "", storeDepartmentInfo.InventoryName, storeDepartmentInfo.StoreName),
							},
							{
								"keyname": helper.RuleThree(req.Language == "en", "Time", "巡店时间"),
								"value":   fmt.Sprintf("%s %s", visitInfo.Date, helper.FormatStoreVisitTimeSpan(req.Language, visitInfo.TimeSpan)),
							},
						},
						AppSign:  qyWork.StoreVisit,
						Language: req.Language,
					})
					if err != nil {
						app.Logger().Error(c).Msgf("巡店留言通知消息发送失败：%+v", err)
					}
				}
			}
		}
	}()
	return commentId, err
}

func (s *CommentSvc) CommentList(c ctx.Context, req storeVisitReq.CommentListReq) ([]storeVisitResp.VisitCommentAllInfo, int64, error) {
	var resList = make([]storeVisitResp.VisitCommentAllInfo, 0)
	list, total, err := s.visitCommentDao.GetListPreloadMedia(c, req)
	if err != nil {
		return nil, 0, err
	}
	var memberIds []int64
	for _, row := range list {
		memberIds = append(memberIds, row.MemberID)
		if len(row.Children) > 0 {
			for _, reply := range row.Children {
				memberIds = append(memberIds, reply.MemberID)
			}
			memberIds = append(memberIds)
		}
	}
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, req.EnterpriseId, helper.ArrayUniq(memberIds))
	if err != nil {
		return nil, 0, err
	}
	// 整理导购id和导购信息的key=>value
	var memberMap = make(map[int64]storeVisitResp.MemberInfo)
	for _, info := range memberDepartmentList {
		dpName := ""
		if len(info.MemberDepartments) > 0 && info.MemberDepartments[0].Department.ID > 0 {
			dpName = info.MemberDepartments[0].Department.Name
		}
		memberMap[info.ID] = storeVisitResp.MemberInfo{
			MemberId: info.ID,
			Name:     info.Name,
			DeptName: dpName,
			Avatar:   info.Avatar,
		}
	}
	for _, row := range list {
		tmp := storeVisitResp.VisitCommentAllInfo{
			VisitCommentInfo: storeVisitResp.VisitCommentInfo{
				ID:         row.ID,
				VisitID:    row.VisitID,
				Content:    row.Content,
				ParentID:   row.ParentID,
				IsDel:      row.IsDel,
				MemberInfo: storeVisitResp.MemberInfo{},
				CreatedAt:  row.CreatedAt,
			},
			Children: nil,
		}
		if member, exist := memberMap[row.MemberID]; exist {
			tmp.MemberInfo = member
		}
		if member, exist := memberMap[row.ReplyMemberID]; exist {
			tmp.ReplyMemberInfo = member
		}
		if len(row.MediaList) > 0 {
			for _, image := range row.MediaList {
				tmp.MediaList = append(tmp.MediaList, storeVisitResp.MediaInfo{
					MediaType:     image.MediaType,
					Url:           image.URL,
					VideoImageUrl: image.VideoImageUrl,
					VideoSeconds:  image.VideoSecond,
				})
			}
		}
		if len(row.Children) > 0 {
			for _, childrenRow := range row.Children {
				tmpChildren := storeVisitResp.VisitCommentInfo{
					ID:        childrenRow.ID,
					VisitID:   childrenRow.VisitID,
					Content:   childrenRow.Content,
					ParentID:  childrenRow.ParentID,
					CreatedAt: childrenRow.CreatedAt,
				}
				if member, exist := memberMap[childrenRow.MemberID]; exist {
					tmpChildren.MemberInfo = member
				}
				if member, exist := memberMap[childrenRow.ReplyMemberID]; exist {
					tmpChildren.ReplyMemberInfo = member
				}
				if len(childrenRow.MediaList) > 0 {
					for _, image := range childrenRow.MediaList {
						tmpChildren.MediaList = append(tmpChildren.MediaList, storeVisitResp.MediaInfo{
							MediaType:     image.MediaType,
							Url:           image.URL,
							VideoImageUrl: image.VideoImageUrl,
							VideoSeconds:  image.VideoSecond,
						})
					}
				}
				tmp.Children = append(tmp.Children, tmpChildren)
			}
		}
		resList = append(resList, tmp)
	}
	return resList, total, nil
}

func (s *CommentSvc) DelComment(c ctx.Context, req storeVisitReq.DelCommentReq) error {
	one, err := s.visitCommentDao.GetOne(c, req.EnterpriseId, "id = ?", req.CommentID)
	if err != nil {
		return err
	}
	if one.MemberID != req.MemberId {
		return errors.New("您不能删除其他人的评论")
	}
	_ = s.visitMediaDao.DeleteData(c, req.EnterpriseId, "rel_id = ? AND rel_type = ?", req.CommentID, storeVisitModel.Comments)
	return s.visitCommentDao.DeleteData(c, req.EnterpriseId, "id = ?", req.CommentID)

	//return s.visitCommentDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{"id": req.CommentID}, map[string]interface{}{
	//	"content":    "该评论已删除",
	//	"is_del":     storeVisitModel.CommentIsDel,
	//	"updated_at": time.Now().Unix(),
	//})
}
