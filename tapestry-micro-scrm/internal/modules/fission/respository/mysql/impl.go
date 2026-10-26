package mysql

import (
	"database/sql"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/uid"
	"time"
)

// FissionRepo 裂变活动仓库实现
type FissionRepo struct {
	*repos.SQLRepo
}

var currentTime = func() int64 {
	return time.Now().Unix()
}

// ActivityList 裂变活动列表
func (f *FissionRepo) ActivityList(c ctx.Context, request *fission.ActivityListRequest) (int64, []*fission.ActivityEntity, error) {
	request.Search.NowTime = time.Now().Unix()
	var resp = make([]*fission.ActivityEntity, 0)
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{
		Offset: int(offset),
		Limit:  int(request.Size),
	}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission_activity/select-activity-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &resp
		},
	}
	total, list, err := f.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, resp, err
	}
	resList := *list.(*[]*fission.ActivityEntity)

	return total, resList, nil
}

// MutualActivityList 选择互斥裂变活动
func (f *FissionRepo) MutualActivityList(c ctx.Context, request *fission.MutualActivityListRequest) (int64, []*fission.MutualActivityEntity, error) {
	cmp := make([]*fission.MutualActivityEntity, 0)
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{
		Offset: int(offset),
		Limit:  int(request.Size),
	}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission_activity/select-activity-mutual",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &cmp
		},
	}
	total, list, err := f.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, cmp, err
	}
	resList := *list.(*[]*fission.MutualActivityEntity)

	return total, resList, nil
}

// fissionMembersSet 设置裂变活动员工
func (f *FissionRepo) fissionMembersSet(c ctx.Context, tx *sqlplus.Tx, activityID int64, enterpriseID int64, members *fission.MembersSet) error {
	// 批量删除员工
	if len(members.DelMembers) > 0 {
		err := f.DeleteFissionMembers(c, tx, &fission.BatchDelete{
			DeleteIDs: members.DelMembers,
			DeletedAt: currentTime(),
		})
		if err != nil {
			return err
		}
	}
	for _, memberID := range members.AddMembers {
		stmt, err := tx.Template("t_fission_members/insert")
		if err != nil {
			continue
		}
		_, err = stmt.NamedExec(c, fission.MembersEntity{
			ID:           uid.ID(c),
			ActivityID:   activityID,
			EnterpriseID: enterpriseID,
			MemberID:     memberID,
			CreatedAt:    currentTime(),
			UpdatedAt:    currentTime(),
		})
		if err != nil {
			log.Info(c).Msgf("新增裂变活动员工失败：%+v，memberID：%d", err, memberID)
			continue
		}
	}
	return nil
}

// fissionPrizesSet 设置裂变活动奖品
func (f *FissionRepo) fissionPrizesSet(c ctx.Context, tx *sqlplus.Tx, activityID int64, enterpriseID int64, prizes *fission.PrizeSet) error {
	if len(prizes.DelPrizes) > 0 {
		err := f.DeleteFissionPrize(c, tx, &fission.BatchDelete{
			DeleteIDs: prizes.DelPrizes,
			DeletedAt: currentTime(),
		})
		if err != nil {
			return err
		}
	}
	for _, prize := range prizes.AddPrizes {
		stmt, err := tx.Template("t_fission_prize/insert")
		if err != nil {
			continue
		}
		_, err = stmt.NamedExec(c, fission.PrizeEntity{
			ID:           uid.ID(c),
			ActivityID:   activityID,
			EnterpriseID: enterpriseID,
			Name:         prize.Name,
			Pic:          prize.Pic,
			Stocks:       prize.Stocks,
			NeedNum:      prize.NeedNum,
			LimitNum:     prize.LimitNum,
			Sort:         prize.Sort,
			CreatedAt:    currentTime(),
			UpdatedAt:    currentTime(),
		})
		if err != nil {
			continue
		}
	}
	for _, prize := range prizes.UpdatePrizes {
		stmt, err := tx.Template("t_fission_prize/update", f.Cond2Args(prize))
		if err != nil {
			continue
		}
		prize.UpdatedAt = currentTime()
		_, err = stmt.NamedExec(c, prize)
		if err != nil {
			log.Info(c).Msgf("修改裂变活动奖品失败：%+v，prize：%+v", err, prize)
			continue
		}
	}
	return nil
}

// fissionGroupSet 设置裂变活动社群
func (f *FissionRepo) fissionGroupSet(c ctx.Context, tx *sqlplus.Tx, activityID int64, enterpriseID int64, groups *fission.ChatGroupSet) error {
	// 批量删除社群
	if len(groups.DelGroups) > 0 {
		err := f.DeleteChatGroup(c, tx, &fission.BatchDelete{
			DeletedAt: currentTime(),
			DeleteIDs: groups.DelGroups,
		})
		if err != nil {
			return err
		}
	}
	for _, group := range groups.AddGroups {
		stmt, err := tx.Template("t_fission_chat_group/insert")
		if err != nil {
			continue
		}
		_, err = stmt.NamedExec(c, fission.ChatGroupEntity{
			ID:           uid.ID(c),
			ActivityID:   activityID,
			EnterpriseID: enterpriseID,
			ChatID:       group.ChatID,
			QrCode:       group.QrCode,
			Sort:         group.Sort,
			CreatedAt:    currentTime(),
			UpdatedAt:    currentTime(),
		})
		if err != nil {
			continue
		}
	}
	for _, group := range groups.UpdateGroups {
		stmt, err := tx.Template("t_fission_chat_group/update")
		if err != nil {
			continue
		}
		_, err = stmt.NamedExec(c, fission.ChatGroupEntity{
			ID:        group.ID,
			ChatID:    group.ChatID,
			QrCode:    group.QrCode,
			Sort:      group.Sort,
			UpdatedAt: currentTime(),
		})
		if err != nil {
			log.Info(c).Msgf("新增添加裂变活动群失败：%+v，group：%+v", err, group)
			continue
		}
	}
	return nil
}

// mutualActivitySet 设置互斥活动
func (f *FissionRepo) mutualActivitySet(c ctx.Context, tx *sqlplus.Tx, activityID int64, enterpriseID int64, mutual *fission.MutualSet) error {
	// 批量删除互斥活动
	if len(mutual.DelMutual) > 0 {
		err := f.DeleteMutualActivity(c, tx, &fission.BatchDelete{
			DeleteIDs: mutual.DelMutual,
			DeletedAt: currentTime(),
		})
		if err != nil {
			return err
		}
	}
	for _, mutualID := range mutual.AddMutual {
		stmt, err := tx.Template("t_fission_mutual_relation/insert")
		if err != nil {
			continue
		}
		_, err = stmt.NamedExec(c, fission.MutualRelationEntity{
			ID:               uid.ID(c),
			ActivityID:       activityID,
			EnterpriseID:     enterpriseID,
			MutualActivityID: mutualID,
			CreatedAt:        currentTime(),
			UpdatedAt:        currentTime(),
		})
		if err != nil {
			continue
		}
	}
	return nil
}

// welcomeSet 设置活动欢迎语
func (f *FissionRepo) welcomeSet(c ctx.Context, tx *sqlplus.Tx, enterpriseID int64, drainageWelcomesID int64, channelID int64, contentContent string, action int) error {
	if action == fission.ActionUpdate {
		stmt, err := tx.Template("t_drainage_welcomes/welcome-msg-update")
		if err != nil {
			return err
		}
		_, err = stmt.NamedExec(c, drainage.Welcome{
			ID:           drainageWelcomesID,
			Content:      contentContent,
			UpdatedAt:    currentTime(),
			EnterpriseId: enterpriseID,
			ChannelId:    channelID,
		})
		if err != nil {
			return err
		}
		return nil
	}
	stmt, err := tx.Template("t_drainage_welcomes/insert")
	_, err = stmt.NamedExec(c, drainage.Welcome{
		ID:           drainageWelcomesID,
		Content:      contentContent,
		EnterpriseId: enterpriseID,
		IsDefault:    drainage.NotDefaultWelcome,
		Type:         drainage.ActivityWelcome,
		ChannelId:    channelID,
		CreatedAt:    currentTime(),
		UpdatedAt:    currentTime(),
		DeletedAt:    0,
	})
	if err != nil {
		return err
	}
	return nil
}

// activityAdd 添加裂变活动
func (f *FissionRepo) activityAdd(c ctx.Context, tx *sqlplus.Tx, request *fission.ActivityEntity) error {
	request.CreatedAt = currentTime()
	request.UpdatedAt = currentTime()

	stmt, err := tx.Template("t_fission_activity/insert")
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// ActivitySet 裂变活动设置
func (f *FissionRepo) ActivitySet(c ctx.Context, enterpriseID int64, action int, activity *fission.ActivityEntity, members *fission.MembersSet, chatGroups *fission.ChatGroupSet, welcomeContent string, prize *fission.PrizeSet, mutual *fission.MutualSet) (int64, error) {
	activityID := activity.ID
	drainageWelcomesID := uid.ID(c)
	txErr := f.DB.Tx(c, func(tx *sqlplus.Tx) error {
		runSafeErr, crash := fx.RunSafe(func() error {
			// 活动主表
			if action == fission.ActionAdd {
				activity.DrainageWelcomesID = drainageWelcomesID
				activity.EnterpriseID = enterpriseID
				err := f.activityAdd(c, tx, activity)
				if err != nil {
					return err
				}
			} else {
				drainageWelcomesID = activity.DrainageWelcomesID
				action = fission.ActionUpdate
				err := f.activityUpdate(c, tx, activity)
				if err != nil {
					return err
				}
			}
			// 裂变活动员工
			err := f.fissionMembersSet(c, tx, activityID, enterpriseID, members)
			if err != nil {
				return err
			}
			// 裂变活动奖品
			err = f.fissionPrizesSet(c, tx, activityID, enterpriseID, prize)
			if err != nil {
				return err
			}
			//添加欢迎语
			err = f.welcomeSet(c, tx, enterpriseID, drainageWelcomesID, activity.ChannelID, welcomeContent, action)
			if err != nil {
				return err
			}
			// 添加裂变活动群
			if activity.FissionType == fission.GroupFissionType {
				err = f.fissionGroupSet(c, tx, activityID, enterpriseID, chatGroups)
				if err != nil {
					return err
				}
			}
			// 添加互斥裂变活动
			if activity.MutualType == fission.IsMutualType {
				err = f.mutualActivitySet(c, tx, activityID, enterpriseID, mutual)
				if err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
			return nil
		})
		if crash != nil {
			return crash
		}
		if runSafeErr != nil {
			return runSafeErr
		}
		return nil
	})
	if txErr != nil {
		return 0, txErr
	}
	return activityID, nil
}

// ActivityDelete 活动删除
func (f *FissionRepo) ActivityDelete(c ctx.Context, activityID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_fission_activity/delete", fission.ActivityEntity{
		ID:        activityID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateStatus 更新活动状态
func (f *FissionRepo) UpdateStatus(c ctx.Context, activityID int64, status int) error {
	_, err := db.Update(c, f.SQLRepo, "t_fission_activity/update-activity-status", fission.ActivityEntity{
		ID:        activityID,
		Status:    status,
		UpdatedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// ChatGroupList 社群列表
func (f *FissionRepo) ChatGroupList(c ctx.Context, request *fission.ChatGroupListRequest) (int64, []*fission.WorkChatGroupEntity, error) {
	resp := make([]*fission.WorkChatGroupEntity, 0)
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{
		Offset: int(offset),
		Limit:  int(request.Size),
	}
	option := &repos.FindPageOption{
		TplPrefix: "t_work_group_chat/chat-group-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &resp
		},
	}
	total, list, err := f.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, resp, err
	}
	resList := *list.(*[]*fission.WorkChatGroupEntity)

	return total, resList, nil
}

// ActivityDetail 查询活动详情
func (f *FissionRepo) ActivityDetail(c ctx.Context, activityID int64, enterpriseID int64) (*fission.ActivityEntity, error) {
	var activity = fission.ActivityEntity{}
	err := db.Get(c, f.SQLRepo, "t_fission_activity/activity-detail", &activity, activityID, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return &activity, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &activity, nil
}

// ByContactIdActivityDetail 通过联系我二维码ID查询裂变活动
func (f *FissionRepo) ByContactIdActivityDetail(c ctx.Context, contactID int64, enterpriseID int64) (*fission.ActivityEntity, error) {
	var activity = fission.ActivityEntity{}
	err := db.Get(c, f.SQLRepo, "t_fission_activity/by-contact-id-activity-detail", &activity, contactID, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &activity, nil
}

// activityUpdate 修改裂变活动
func (f *FissionRepo) activityUpdate(c ctx.Context, tx *sqlplus.Tx, request *fission.ActivityEntity) error {
	request.UpdatedAt = currentTime()
	stmt, err := tx.Template("t_fission_activity/update-activity")
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// FissionMemberIDs 裂变活动员工ID
func (f *FissionRepo) FissionMemberIDs(c ctx.Context, activityID int64) ([]*fission.Member, error) {
	members := make([]*fission.Member, 0)
	err := db.Query(c, f.SQLRepo, "t_fission_members/select-fission-member-ids", &members, activityID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return members, nil
}

// MemberList 裂变活动员工
func (f *FissionRepo) MemberList(c ctx.Context, memberIDs []int64, enterpriseID int64) ([]*fission.Member, error) {
	mapping := make(map[string]interface{})
	mapping["member_ids"] = memberIDs
	mapping["enterprise_id"] = enterpriseID
	members := make([]*fission.Member, 0)
	err := db.NamedSelect(c, f.SQLRepo, "t_members/select-member-name-by-ids", &members, mapping, &mapping)
	if err != nil && err != sql.ErrNoRows {
		return members, err
	}
	return members, nil
}

// QwActivityList 裂变活动列表
func (f *FissionRepo) QwActivityList(c ctx.Context, request *fission.QwActivityListRequest) (int64, *[]fission.QwActivityList, error) {
	cmp := &[]fission.QwActivityList{}
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{
		Offset: int(offset),
		Limit:  int(request.Size),
	}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission_activity/qw-activity-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, list, err := f.FindPage(c, pageVo, &request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, cmp, err
	}
	resList := list.(*[]fission.QwActivityList)

	return total, resList, nil
}

// MutualActivityIDs 查询互斥活动IDs
func (f *FissionRepo) MutualActivityIDs(c ctx.Context, activityID int64, enterpriseID int64) ([]*fission.MutualRelationEntity, error) {
	var mutualActivity = make([]*fission.MutualRelationEntity, 0)
	var mapping = make(map[string]interface{})
	mapping["activityID"] = activityID
	mapping["enterpriseID"] = enterpriseID
	err := db.NamedSelect(c, f.SQLRepo, "t_fission_mutual_relation/select-mutual-activity-list", &mutualActivity, mapping, &mapping)
	if err != nil && err != sql.ErrNoRows {
		return mutualActivity, err
	}
	return mutualActivity, nil
}

// StatisticsInfo 统计活动数据
func (f *FissionRepo) StatisticsInfo(c ctx.Context, request *fission.QwActivityListRequest) (*fission.Statistics, error) {
	stmt, stmtErr := f.DB.Template("t_fission_activity/activity-statistics-info", f.Cond2Args(request))
	if stmtErr != nil {
		return nil, stmtErr
	}
	ResponseInfo := &fission.Statistics{}
	getErr := stmt.NamedGet(c, ResponseInfo, request)
	if getErr != nil {
		return nil, getErr
	}
	return ResponseInfo, nil
}

// FissionGroupIDs 获取裂变群IDs
func (f *FissionRepo) FissionGroupIDs(c ctx.Context, activityID int64, enterpriseID int64, withDel bool) ([]*fission.ChatGroupEntity, error) {
	var mapping = make(map[string]interface{})
	mapping["activityID"] = activityID
	mapping["enterpriseID"] = enterpriseID
	mapping["withDel"] = withDel
	chatGroup := make([]*fission.ChatGroupEntity, 0)
	err := db.NamedSelect(c, f.SQLRepo, "t_fission_chat_group/select-fission-chat-group", &chatGroup, mapping, &mapping)
	if err != nil && err != sql.ErrNoRows {
		return chatGroup, err
	}
	return chatGroup, nil
}

// WelComeInfo 获取活动欢迎语
func (f *FissionRepo) WelComeInfo(c ctx.Context, enterpriseID int64, drainageWelcomesID int64) (string, error) {
	type WelComeContent struct {
		Content string `json:"content" db:"content"`
	}
	content := &WelComeContent{}
	err := db.Get(c, f.SQLRepo, "t_drainage_welcomes/select-fission-activity-welcome", content, enterpriseID, drainageWelcomesID, 0)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return content.Content, nil
}

//AssistantUserList 助力人信息
func (f *FissionRepo) AssistantUserList(c ctx.Context, ids []int64, enterpriseID int64) ([]*fission.AssistantRecordInfo, error) {
	var (
		result = make([]*fission.AssistantRecordInfo, 0)
	)
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	mapping["enterpriseID"] = enterpriseID
	stmt, err := f.DB.Template("t_fission/assistant-user-list", &mapping)
	if err != nil {
		return nil, err
	}
	err = stmt.NamedSelect(c, &result, mapping)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	return result, nil
}

//FissionMemberList 获取活动下的裂变成员
func (f *FissionRepo) FissionMemberList(c ctx.Context, enterpriseID int64, acvitityID int64) ([]*fission.Member, error) {

	var (
		members = make([]*fission.Member, 0)
	)
	stmt, err := f.DB.Template("t_fission/fission-member-list")
	if nil != err {
		return nil, err
	}
	err = stmt.Select(c, &members, enterpriseID, acvitityID)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	return members, nil
}

//StoreExchangeWarnConfig 配置兑换提醒设置
func (f *FissionRepo) StoreExchangeWarnConfig(c ctx.Context, req fission.ExchangeWarnRequest) error {

	var (
		entry = make(map[string]interface{})
		bytes []byte
		err   error
		value string
	)
	entry["status"] = req.Status
	entry["member_ids"] = req.MemberIds
	bytes, err = json.Marshal(entry)
	if nil != err {
		return err
	}
	value = string(bytes)
	stmt, err := f.DB.Template("t_fission/store-exchange-warn-config")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, req.EnterpriseID, fission.ExchangeWarn, value)
	if nil != err {
		return err
	}
	return nil
}

//ExchangeTotal 奖品兑换记录总数
func (f *FissionRepo) ExchangeTotal(c ctx.Context, enterpriseID int64) (int64, error) {
	var (
		total int64
	)
	stmt, err := f.DB.Template("t_fission/exchange-total")

	if nil != err {
		return 0, err
	}
	err = stmt.Get(c, &total, enterpriseID)
	if nil != err {
		return 0, err
	}
	return total, nil
}

//ActivityListSimple 活动列表 全部有效活动
func (f *FissionRepo) ActivityListSimple(c ctx.Context, enterpriseID int64) ([]*fission.ActivityEntity, error) {
	var (
		list = make([]*fission.ActivityEntity, 0)
	)
	stmt, err := f.DB.Template("t_fission_activity/activity-list-simple")
	if nil != err {
		return list, err
	}
	err = stmt.Select(c, &list, enterpriseID)
	if nil != errorIgnore.Ignore(err) {
		return nil, err
	}
	return list, nil
}

//ExchangeRecord 发起人兑换记录查询
func (f *FissionRepo) ExchangeRecord(c ctx.Context, req fission.ExchangeDetailRequest) (int64, *[]fission.ExchangeRecordInfo, error) {
	var (
		list   *[]fission.ExchangeRecordInfo
		err    error
		total  int64
		cmp    = &[]fission.ExchangeRecordInfo{}
		result interface{}
	)
	req.DeliverName, _ = common.AesEncrypt(req.DeliverName)
	req.DeliverMobile, _ = common.AesEncrypt(req.DeliverMobile)
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(req.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission/exchange-record-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, result, err = f.FindPage(c, pageVo, &req, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, list, err
	}
	list = result.(*[]fission.ExchangeRecordInfo)
	for k, item := range *list {
		(*list)[k].DeliverName, _ = common.AesDecrypt(item.DeliverName)
		(*list)[k].DeliverAddress, _ = common.AesDecrypt(item.DeliverAddress)
		(*list)[k].DeliverArea, _ = common.AesDecrypt(item.DeliverArea)
		(*list)[k].DeliverAreaCode, _ = common.AesDecrypt(item.DeliverAreaCode)
		mobileText, _ := common.AesDecrypt(item.DeliverMobile)
		(*list)[k].DeliverMobile = mobileText
	}
	return total, list, nil
}

//ActivityDetailList 活动详情 列表
func (f *FissionRepo) ActivityDetailList(c ctx.Context, req fission.ActivityDetailRequest) (int64, *[]fission.ActivityUserEntity, error) {

	var (
		list   *[]fission.ActivityUserEntity
		err    error
		total  int64
		cmp    = &[]fission.ActivityUserEntity{}
		result interface{}
	)
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(req.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission/activity-detail-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, result, err = f.FindPage(c, pageVo, &req, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	list = result.(*[]fission.ActivityUserEntity)
	return total, list, nil
}

//InviteDetail发起人邀请明细
func (f *FissionRepo) InviteDetail(c ctx.Context, req fission.InviteDetailRequest) (int64, []*fission.AssistantRecordInfo, error) {
	var (
		cmp           = &[]*fission.AssistantRecordInfo{}
		result        interface{}
		assistantIds  []int64
		assistantInfo []*fission.AssistantRecordInfo
		err           error
		total         int64
		list          []*fission.AssistantRecordInfo
	)
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(req.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_fission_assistant_records/invite-detail-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return cmp
		},
	}
	total, result, err = f.FindPage(c, pageVo, &req, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	list = *result.(*[]*fission.AssistantRecordInfo)
	for _, item := range list {
		assistantIds = append(assistantIds, item.AssistantUserID)
	}
	//查找对应的头像和昵称
	if ok := len(assistantIds) > 0; ok {
		assistantInfo, err = f.AssistantUserList(c, assistantIds, req.EnterpriseID)
		if err != nil {
			return 0, nil, err
		}
		for i := 0; i < len(assistantInfo); i++ {
			for j := 0; j < len(list); j++ {
				if list[j].AssistantUserID == assistantInfo[i].AssistantUserID {
					list[j].ImgUrl = assistantInfo[i].ImgUrl
					list[j].Nickname = assistantInfo[i].Nickname
				}
			}
		}
	}
	return total, list, nil
}

// PrizeList 获取奖品列表
func (f *FissionRepo) PrizeList(c ctx.Context, enterpriseID int64, activityID int64, page int64, size int64) (int64, *[]*fission.PrizeEntity, error) {
	var resList = &[]*fission.PrizeEntity{}
	mapping := make(map[string]interface{})
	mapping["activity_id"] = activityID
	mapping["enterprise_id"] = enterpriseID
	total, list, err := db.FindPage(c, f.SQLRepo, "t_fission_prize/list", page, size, mapping, resList)
	if err != nil {
		return total, resList, err
	}
	resList = list.(*[]*fission.PrizeEntity)
	return total, resList, nil
}

// PrizeInfo 查询奖品信息
func (f *FissionRepo) PrizeInfo(c ctx.Context, ID int64, withDeleted bool) (*fission.PrizeEntity, error) {
	var prizeInfo fission.PrizeEntity
	var tpl = "t_fission_prize/info"
	if withDeleted {
		tpl = "t_fission_prize/info-with-deleted"
	}
	err := db.Get(c, f.SQLRepo, tpl, &prizeInfo, ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &prizeInfo, nil
}

func (f *FissionRepo) GetExchangeWarnConfig(c ctx.Context, enterpriseId int64, key string) (string, error) {

	var (
		config string
		err    error
	)
	err = db.Get(c, f.SQLRepo, "t_fission/exchange-warn-config", &config, enterpriseId, fission.ExchangeWarn)
	return config, err
}

// GetActivityUser 获取活动用户信息
func (f *FissionRepo) GetActivityUser(c ctx.Context, enterpriseID int64, activityID int64, userID int64) (*fission.ActivityUserEntity, error) {
	var userInfo fission.ActivityUserEntity
	err := db.Get(c, f.SQLRepo, "t_fission_user/info", &userInfo, enterpriseID, activityID, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &userInfo, nil
}

// GetActivityUser 获取活动用户信息
func (f *FissionRepo) GetActivityUserByContactID(c ctx.Context, ContactID int64) (*fission.ActivityUserEntity, error) {
	userInfo := new(fission.ActivityUserEntity)
	err := db.Get(c, f.SQLRepo, "t_fission_user/info-by-contact", userInfo, ContactID)

	switch true {
	case err == sql.ErrNoRows || userInfo == nil:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return userInfo, nil
	}

}

// GetMpUserInfo 查询用户信息
func (f *FissionRepo) GetMpUserInfo(c ctx.Context, ID int64) (*fission.MpUser, error) {
	var userInfo fission.MpUser
	err := db.Get(c, f.SQLRepo, "t_fission_mp_user/detail", &userInfo, ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &userInfo, nil
}

// RecoveryFissionRecord 查询是否存在删除的记录
func (f *FissionRepo) RecoveryFissionRecord(ctx ctx.Context, customerID, memberID int64) (*fission.AssistantRecordsEntity, error) {
	stmt, err := f.DB.Template("t_fission_assistant_records/select-recovery-record")
	if err != nil {
		return nil, err
	}

	data := new(fission.AssistantRecordsEntity)
	err = stmt.Get(ctx, data, customerID, memberID)

	switch true {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return data, nil
	}
}

// RecoveryFissionRecord 查询是否存在记录
func (f *FissionRepo) GetRecordByRelation(ctx ctx.Context, req *fission.ExistRelationReq) (*fission.AssistantRecordsEntity, error) {
	stmt, err := f.DB.Template("t_fission_assistant_records/select-exist-relation", f.Cond2Args(req))
	if err != nil {
		return nil, err
	}
	data := new(fission.AssistantRecordsEntity)
	err = stmt.Get(ctx, data)

	switch true {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return data, nil
	}
}

// UpdateAssistantRecordWithTx 使用事务更新 AssistantRecord 记录
func (f *FissionRepo) UpdateAssistantRecordWithTx(c ctx.Context, tx *sqlplus.Tx, record *fission.UpdateAssistantRecord) (int64, error) {
	stmt, err := tx.Template("t_fission_assistant_records/update-record-by-id", f.Cond2Args(record))
	if err != nil {
		return 0, err
	}
	result, err := stmt.NamedExec(c, record)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()

}

// CreateAssistantRecordWithTx 更新记录
func (f *FissionRepo) CreateAssistantRecordWithTx(c ctx.Context, tx *sqlplus.Tx, record *fission.AssistantRecordsEntity) (int64, error) {
	stmt, err := tx.Template("t_fission_assistant_records/insert-record", f.Cond2Args(record))
	if err != nil {
		return 0, err
	}
	result, err := stmt.NamedExec(c, record)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// DeleteFissionMembers 批量删除
func (f *FissionRepo) DeleteFissionMembers(c ctx.Context, tx *sqlplus.Tx, request *fission.BatchDelete) error {
	stmt, err := tx.Template("t_fission_members/delete", f.Cond2Args(request))
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFissionPrize 删除裂变奖品
func (f *FissionRepo) DeleteFissionPrize(c ctx.Context, tx *sqlplus.Tx, request *fission.BatchDelete) error {
	stmt, err := tx.Template("t_fission_prize/delete", f.Cond2Args(request))
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteChatGroup 删除裂变群
func (f *FissionRepo) DeleteChatGroup(c ctx.Context, tx *sqlplus.Tx, request *fission.BatchDelete) error {
	stmt, err := tx.Template("t_fission_chat_group/delete", f.Cond2Args(request))
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMutualActivity 删除互斥活动
func (f *FissionRepo) DeleteMutualActivity(c ctx.Context, tx *sqlplus.Tx, request *fission.BatchDelete) error {
	stmt, err := tx.Template("t_fission_mutual_relation/delete", f.Cond2Args(request))
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

// ByActivityContactQrCode 通过活动Id获取群二维码
func (f *FissionRepo) ByActivityContactQrCode(c ctx.Context, activityID int64, enterpriseID int64) (*fission.ChatGroupEntity, error) {
	chatGroupEntity := new(fission.ChatGroupEntity)
	mapping := make(map[string]interface{})
	mapping["activity_id"] = activityID
	mapping["enterprise_id"] = enterpriseID
	stmt, stmtErr := f.DB.Template("t_fission_chat_group/get-chat-group-qr-code", &mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	err := stmt.NamedGet(c, chatGroupEntity, mapping)

	switch true {
	case err != nil:
		return nil, err
	default:
		return chatGroupEntity, nil
	}
}

// GetActivityUserList 获取用户参与指定活动列表信息
func (f *FissionRepo) GetActivityUserList(c ctx.Context, userID int64, activityIDList []int64, enterpriseID int64) ([]*fission.ActivityUserEntity, error) {
	var list = make([]*fission.ActivityUserEntity, 0)
	var mapping = make(map[string]interface{})
	mapping["user_id"] = userID
	mapping["activity_id_list"] = activityIDList
	mapping["enterprise_id"] = enterpriseID
	err := db.NamedSelect(c, f.SQLRepo, "t_fission_user/list-by-activity-ids", &list, mapping, &mapping)
	if err != nil {
		return list, err
	}
	return list, nil
}

// AddFissionUser 添加用户活动数据
func (f *FissionRepo) AddFissionUser(c ctx.Context, entity *fission.ActivityUserEntity) (int64, error) {
	return db.Add(c, f.SQLRepo, "t_fission_user/insert", entity)
}

// GetExchangeList 获取用户兑换列表
func (f *FissionRepo) GetExchangeList(c ctx.Context, enterpriseID int64, activityID int64, userID int64, page int64, size int64) (int64, *[]*fission.ExchangePrizeRecordsEntity, error) {
	var resList = &[]*fission.ExchangePrizeRecordsEntity{}
	var mapping = make(map[string]interface{})
	mapping["activity_id"] = activityID
	mapping["user_id"] = userID
	mapping["enterprise_id"] = enterpriseID
	total, list, err := db.FindPage(c, f.SQLRepo, "t_fission_exchange_prize_records/list", page, size, mapping, resList)
	if err != nil {
		return 0, resList, err
	}
	resList = list.(*[]*fission.ExchangePrizeRecordsEntity)
	return total, resList, nil
}

// PrizeListByIDs 根据id获取奖品列表
func (f *FissionRepo) PrizeListByIDs(c ctx.Context, IDList []int64) (*[]*fission.PrizeEntity, error) {
	var list = make([]*fission.PrizeEntity, 0)
	if len(IDList) == 0 {
		return &list, nil
	}

	var mapping = make(map[string]interface{})
	mapping["id_list"] = IDList
	err := db.NamedSelect(c, f.SQLRepo, "t_fission_prize/list-by-ids", &list, mapping, &mapping)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// GetExchangeInfo 获取兑换详情
func (f *FissionRepo) GetExchangeInfo(c ctx.Context, ID int64) (*fission.ExchangePrizeRecordsEntity, error) {
	var data fission.ExchangePrizeRecordsEntity
	err := db.Get(c, f.SQLRepo, "t_fission_exchange_prize_records/detail", &data, ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &data, nil
}

// GetAssistantList 获取助力列表
func (f *FissionRepo) GetAssistantList(c ctx.Context, activityID int64, userID int64, page int64, size int64) (int64, *[]*fission.AssistantRecordsEntity, error) {
	var resList = &[]*fission.AssistantRecordsEntity{}
	var mapping = make(map[string]interface{})
	mapping["activity_id"] = activityID
	mapping["user_id"] = userID
	total, list, err := db.FindPage(c, f.SQLRepo, "t_fission_assistant_records/list", page, size, mapping, resList)
	if err != nil {
		return 0, nil, err
	}
	resList = list.(*[]*fission.AssistantRecordsEntity)
	return total, resList, nil
}

// UserExchange 用户兑换奖品
func (f *FissionRepo) UserExchange(c ctx.Context, userActivity *fission.ActivityUserEntity, prize *fission.PrizeEntity, record *fission.ExchangePrizeRecordsEntity) (bool, error) {
	txErr := f.DB.Tx(c, func(tx *sqlplus.Tx) error {
		runSafeErr, crash := fx.RunSafe(func() error {
			// 更新用户信息
			_, err := f.UpdateFissionUser(c, userActivity)
			if err != nil {
				return err
			}
			// 更新奖品信息
			_, err = f.UpdateFissionPrize(c, prize)
			if err != nil {
				return err
			}
			_, err = f.AddExchangePrizeRecords(c, record)
			if err != nil {
				return err
			}
			return nil
		})
		if crash != nil {
			return crash
		}
		if runSafeErr != nil {
			return runSafeErr
		}
		return nil
	})
	if txErr != nil {
		return false, txErr
	}
	return true, nil
}

// UpdateFissionUser 更新用户活动信息
func (f *FissionRepo) UpdateFissionUser(c ctx.Context, entity *fission.ActivityUserEntity) (int64, error) {
	return db.Update(c, f.SQLRepo, "t_fission_user/update", entity)
}

// UpdateFissionPrize 更新奖品信息
func (f *FissionRepo) UpdateFissionPrize(c ctx.Context, entity *fission.PrizeEntity) (int64, error) {
	return db.Update(c, f.SQLRepo, "t_fission_prize/update-all", entity)
}

// AddExchangePrizeRecords 添加奖品兑换记录
func (f *FissionRepo) AddExchangePrizeRecords(c ctx.Context, entity *fission.ExchangePrizeRecordsEntity) (int64, error) {
	return db.Add(c, f.SQLRepo, "t_fission_exchange_prize_records/insert", entity)
}

// GetUserExchangeNum 获取用户兑换奖品数量
func (f *FissionRepo) GetUserExchangeNum(c ctx.Context, userID int64, prizeID int64) (int64, error) {
	var count struct {
		Count sql.NullInt64 `db:"count"`
	}
	err := db.Get(c, f.SQLRepo, "t_fission_exchange_prize_records/query-user-exchange-count", &count, prizeID, userID)
	if err != nil {
		return 0, err
	}
	return count.Count.Int64, nil
}

// QueryUserLatestExchangeRecord 获取用户最近最近兑换的记录
func (f *FissionRepo) QueryUserLatestExchangeRecord(c ctx.Context, userID int64) (*fission.ExchangePrizeRecordsEntity, error) {
	var entity fission.ExchangePrizeRecordsEntity
	err := db.Get(c, f.SQLRepo, "t_fission_exchange_prize_records/query-latest", &entity, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &entity, nil
}

// UpdateActivityStatisticsWithTx 更新活动统计信息
func (f *FissionRepo) UpdateActivityStatisticsWithTx(c ctx.Context, tx *sqlplus.Tx, ID int64, statistics *fission.Statistics) error {
	stmt, err := tx.Template("t_fission_activity/update-activity-statistics", f.Cond2Args(statistics))
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateFissionUserStatisticsWithTx 通过事务去更新 FissionUser 记录
func (f *FissionRepo) UpdateFissionUserStatisticsWithTx(c ctx.Context, tx *sqlplus.Tx, ID int64, statistics *fission.UpdateStatisticsData) error {
	stmt, err := tx.Template("t_fission_user/update-activity-statistics", f.Cond2Args(statistics))
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, ID)
	if err != nil {
		return err
	}
	return nil
}

// GetActivityByChatID 通过聊天ID查询关联的的活动ID
func (f *FissionRepo) GetActivityByChatID(c ctx.Context, req *fission.ActivityByChatIDReq) ([]int64, error) {
	stmt, err := f.DB.Template("t_fission_activity/get-activity-by-chat-id", f.Cond2Args(req))
	if err != nil {
		return nil, err
	}
	list := make([]*fission.ChatGroupEntity, 0)
	res, err := stmt.NamedQuery(c, req)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	sqlErr := sqlplus.ScanAll(res, &list)
	if sqlErr != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

// FindGroupRecordsByChatID 通过聊天ID查询需要入群的记录
func (f *FissionRepo) FindGroupRecordsByChatID(c ctx.Context, activityID []int64, chatID, enterpriseID int64, joinOrLeave bool) ([]*fission.AssistantRecordsWithTime, error) {
	record := make(map[string]interface{}, 0)
	record["ActivityID"] = activityID
	record["ChatID"] = chatID
	record["JoinOrLeave"] = joinOrLeave
	record["EnterpriseID"] = enterpriseID
	stmt, err := f.DB.Template("t_fission_assistant_records/find-group-records-by-chat", f.Cond2Args(record))
	if err != nil {
		return nil, err
	}
	list := make([]*fission.AssistantRecordsWithTime, 0)
	res, err := stmt.NamedQuery(c, record)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, &list)
	if sqlplusErr != nil {
		return nil, err
	}
	return list, nil
}

// GetActivityDetailWithDeleted 获取活动详情，查询删除的数据
func (f *FissionRepo) GetActivityDetailWithDeleted(c ctx.Context, activityID int64, enterpriseID int64) (*fission.ActivityEntity, error) {
	var entity fission.ActivityEntity
	err := db.Get(c, f.SQLRepo, "t_fission_activity/detail-with-deleted", &entity, enterpriseID, activityID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &entity, nil
}

// RecoveryRecordAndStatistics 编辑好友执行操作
func (f *FissionRepo) RecoveryRecordAndStatistics(c ctx.Context, assistantRecord *fission.AssistantRecordsEntity, fissionType int, updateData *fission.UpdateAssistantRecord, fissionUserStatistics *fission.UpdateStatisticsData) error {
	// 使用事务
	err := f.DB.Tx(c, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {

			// 1. 更新记录表的状态
			rows, err := f.UpdateAssistantRecordWithTx(c, tx, updateData)
			if err != nil {
				return err
			}
			if rows <= 0 {
				return errors.New("未更新数据")
			}

			// 2. 更新裂变人统计数据
			err = f.UpdateFissionUserStatisticsWithTx(c, tx, assistantRecord.FissionUserID, fissionUserStatistics)
			if err != nil {
				return err
			}

			// 3. 更新活动统计数据
			s := fissionUserStatistics.Statistics
			err = f.UpdateActivityStatisticsWithTx(c, tx, assistantRecord.ActivityID, &s)
			if err != nil {
				return err
			}
			return nil
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return err
}

// CreateAssistantRecordAndStatistics 添加好友操作
func (f *FissionRepo) CreateAssistantRecordAndStatistics(c ctx.Context, record *fission.AssistantRecordsEntity, fissionUserStatistics *fission.UpdateStatisticsData) error {
	txErr := f.DB.Tx(c, func(tx *sqlplus.Tx) error {
		runSafeErr, crash := fx.RunSafe(func() error {
			// 1. 创建记录表的状态
			rows, err := f.CreateAssistantRecordWithTx(c, tx, record)
			if err != nil {
				return err
			}
			if rows <= 0 {
				return errors.New("未插入成功")
			}

			// 2. 更新裂变人统计数据数据
			err = f.UpdateFissionUserStatisticsWithTx(c, tx, record.FissionUserID, fissionUserStatistics)
			if err != nil {
				return err
			}

			// 3. 更新活动统计数据
			s := fissionUserStatistics.Statistics
			err = f.UpdateActivityStatisticsWithTx(c, tx, record.ActivityID, &s)
			if err != nil {
				return err
			}
			return nil
		})
		if crash != nil {
			return crash
		}
		if runSafeErr != nil {
			return runSafeErr
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

// HasWorkExternalRelation 查询用户是否有外部联系人关系数据
func (f *FissionRepo) HasWorkExternalRelation(c ctx.Context, enterpriseID int64, workUserID int64, memberIDs []int64) (bool, error) {
	type ID struct {
		ID int64 `db:"id"`
	}
	var relation = make([]*ID, 0)
	var mapping = make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	mapping["work_user_id"] = workUserID
	if len(memberIDs) > 0 {
		mapping["member_ids"] = memberIDs
	}

	err := db.NamedSelect(c, f.SQLRepo, "t_work_external_user_relation/has-relation-with-members", &relation, mapping, &mapping)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	if len(relation) == 0 {
		return false, nil
	}
	return true, nil
}

// GetAssistantUserList 获取活动最新的助力数据
func (f *FissionRepo) GetAssistantUserList(c ctx.Context, enterpriseID, activityID, index, limit int64, needGroup bool) ([]*fission.AssistantRecordsEntity, error) {
	var list = make([]*fission.AssistantRecordsEntity, 0)
	var tpl = "t_fission_assistant_records/query-activity-latest-list"
	if needGroup {
		tpl = "t_fission_assistant_records/query-activity-latest-list-group"
	}
	err := db.Query(c, f.SQLRepo, tpl, &list, enterpriseID, activityID, index, limit)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetTotalAssistantUser 获取活动总助力数
func (f *FissionRepo) GetTotalAssistantUser(c ctx.Context, enterpriseID int64, activityID int64, needGroup bool) (int64, error) {
	var count struct {
		Count sql.NullInt64 `db:"count"`
	}
	var tpl = "t_fission_assistant_records/query-activity-total"
	if needGroup {
		tpl = "t_fission_assistant_records/query-activity-total-group"
	}
	err := db.Get(c, f.SQLRepo, tpl, &count, enterpriseID, activityID)
	if err != nil {
		return 0, err
	}
	return count.Count.Int64, nil
}

// ActivityUserNumAdd 添加活动发起人数
func (f *FissionRepo) ActivityUserNumAdd(c ctx.Context, activityID int64, num int) (int64, error) {
	return db.Exec(c, f.SQLRepo, "t_fission_activity/add-user-num", num, activityID)
}

// MutualDeleteByActivityID 删除该活动互斥活动
func (f *FissionRepo) MutualDeleteByActivityID(c ctx.Context, activityID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_fission_mutual_relation/delete-by-activity-id", fission.ActivityEntity{
		ID:        activityID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// ChatDeleteByActivityID 删除该活动社群
func (f *FissionRepo) ChatDeleteByActivityID(c ctx.Context, activityID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_fission_chat_group/delete-by-activity-id", fission.ActivityEntity{
		ID:        activityID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// PrizeDeleteByActivityID 删除该活动的奖品
func (f *FissionRepo) PrizeDeleteByActivityID(c ctx.Context, activityID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_fission_prize/delete-by-activity-id", fission.ActivityEntity{
		ID:        activityID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// WelcomeDeleteByID 删除该活动的欢迎语
func (f *FissionRepo) WelcomeDeleteByID(c ctx.Context, ID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_drainage_welcomes/delete-one-welcome", fission.ActivityEntity{
		ID:        ID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}

// ContactDeleteByID 删除该活动的默认联系我二维码
func (f *FissionRepo) ContactDeleteByID(c ctx.Context, ID int64) error {
	_, err := db.Update(c, f.SQLRepo, "t_drainage_contact/contacts-delete-by-id", fission.ActivityEntity{
		ID:        ID,
		DeletedAt: time.Now().Unix(),
	})
	if err != nil {
		return err
	}
	return nil
}
