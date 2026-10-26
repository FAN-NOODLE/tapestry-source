package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"micro-scrm/internal/pkg/assist"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/uid"
)

type DrainageRepo struct {
	*repos.SQLRepo
}

func (d *DrainageRepo) ContactAdd(c ctx.Context, request drainage.Contacts) error {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-add")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}
func (d *DrainageRepo) ContactMembersAdd(c ctx.Context, ages []*drainage.ContactMembers) error {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-members-add")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, ages)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}
func (d *DrainageRepo) ContactUpdate(c ctx.Context, ages drainage.Contacts) error {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-update")
	if err != nil {
		return err
	}
	result, sTmtErr := stmt.NamedExec(c, ages)
	if sTmtErr != nil {
		return sTmtErr
	}
	_, _ = result.RowsAffected()
	return sTmtErr
}
func (d *DrainageRepo) ContactMembersDelete(c ctx.Context, ages drainage.Contacts) error {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-members-delete")
	if err != nil {
		return err
	}
	result, sTmtErr := stmt.NamedExec(c, ages)
	if sTmtErr != nil {
		return sTmtErr
	}
	_, _ = result.RowsAffected()
	return sTmtErr
}

// ChannelsLabelCreate 渠道标签创建
func (d *DrainageRepo) ChannelsLabelCreate(c ctx.Context, req []drainage.ContactChannelLabels) error {
	stmt, err := d.DB.Template("t_contact_channels/insert-channels-labels")

	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(c, req)

	if errStmt != nil {
		return errStmt
	}
	return nil
}

// ChannelsCreate 添加渠道
func (d *DrainageRepo) ChannelsCreate(ctx ctx.Context, t drainage.ContactChannels) (int64, error) {
	stmt, err := d.DB.Template("t_contact_channels/insert-channels")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return 0, errStmt
	}

	return res.LastInsertId()
}

// ChannelsList 渠道列表
func (d *DrainageRepo) ChannelsList(cc ctx.Context, req drainage.ContactChannelListRequest) (int64, *[]drainage.ContactChannelsList, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_contact_channels/select",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]drainage.ContactChannelsList{}
		},
	}
	offset := 0
	if req.Page > 0 && req.Size > 0 {
		offset = (req.Page - 1) * req.Size
	}
	total, list, err := d.FindPage(cc, repos.PageVO{
		Limit:  req.Size,
		Offset: offset,
	}, req, option)

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]drainage.ContactChannelsList)

	return total, resList, nil
}

// ChannelDetail 渠道详情
func (d *DrainageRepo) ChannelDetail(ctx ctx.Context, req drainage.ContactChannels) (*drainage.ContactChannels, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = req
	stmt, err := d.DB.Template("t_contact_channels/detail", &mapping)
	if err != nil {
		return nil, err
	}
	dest := &drainage.ContactChannels{}
	sErr := stmt.NamedGet(ctx, dest, req)
	if sErr != nil {
		return nil, sErr
	}
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	return dest, nil
}

// ChannelsLabelsList 渠道标签列表
func (d *DrainageRepo) ChannelsLabelsList(cc ctx.Context, req drainage.ContactChannelLabelsListRequest) (int64, *[]drainage.ContactChannelLabels, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_contact_channel_labels/select",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]drainage.ContactChannelLabels{}
		},
	}
	offset := 0
	if req.Page > 0 && req.Size > 0 {
		offset = (req.Page - 1) * req.Size
	}
	total, list, err := d.FindPage(cc, repos.PageVO{
		Limit:  req.Size,
		Offset: offset,
	}, req, option)

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]drainage.ContactChannelLabels)

	return total, resList, nil
}

// ContactSelect 联系我二维码检错
func (d *DrainageRepo) ContactSelect(ctx ctx.Context, request drainage.ContactDetailsRequest) (*drainage.Contacts, error) {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-select", d.Cond2Args(request))
	var contacts drainage.Contacts
	if err != nil {
		log.Info(ctx).Msgf("%+v", err)
	}
	err = stmt.Get(ctx, &contacts, request.ID)
	return &contacts, err
}
func (d *DrainageRepo) ContactMemberSelect(ctx ctx.Context, request drainage.ContactDetailsRequest) ([]drainage.MemberData, error) {
	var memberData []drainage.MemberData
	stmt, err := d.DB.Template("t_drainage_contact/contacts-member-select")
	if err != nil {
		return memberData, err
	}
	errStmt := stmt.Select(ctx, &memberData, request.ID, 0)
	if errStmt != nil {
		return memberData, err
	}
	return memberData, nil
}

func (d *DrainageRepo) ContactDepartmentSelect(ctx ctx.Context, request drainage.ContactDetailsRequest) ([]drainage.DepartmentData, error) {
	var departmentData []drainage.DepartmentData
	stmt, err := d.DB.Template("t_drainage_contact/contacts-department-select")
	if err != nil {
		return departmentData, err
	}
	errStmt := stmt.Select(ctx, &departmentData, request.EnterpriseId, request.ID, 0)
	if errStmt != nil {
		return departmentData, err
	}
	return departmentData, nil
}

func (d *DrainageRepo) ContactDelete(c ctx.Context, request drainage.ContactDeleteRequest) error {
	var contacts drainage.Contacts
	contacts.DeletedAt = time.Now().Unix()
	mapping := make(map[string]interface{})
	mapping["ids"] = request.ID
	stmt, err := d.DB.Template("t_drainage_contact/contacts-delete", &mapping)
	if err != nil {
		return err
	}
	result, sTmtErr := stmt.NamedExec(c, contacts)
	if sTmtErr != nil {
		return sTmtErr
	}
	_, _ = result.RowsAffected()
	return sTmtErr
}

// InsertWelcomeMember 新增drainage_welcome_members
func (d *DrainageRepo) InsertWelcomeMember(c ctx.Context, request []*drainage.WelcomeMember) (int64, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/insert-welcome-member")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	return result.LastInsertId()
}

// SelectWelcomeOne 获取一条SelectWelcomeOne
func (d *DrainageRepo) SelectWelcomeOne(c ctx.Context, request drainage.WelcomeDetailRequest) (*drainage.Welcome, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/select-one-welcome")
	if err != nil {
		return nil, err
	}
	welcome := &drainage.Welcome{}
	sErr := stmt.Get(c, welcome, request.EnterpriseId, request.ID, request.IsDefault)
	if sErr != nil {
		return nil, sErr
	}
	return welcome, nil
}

// SelectWelcomeDefault 获取一条SelectWelcomeDefault
func (d *DrainageRepo) SelectWelcomeDefault(c ctx.Context, request drainage.WelcomeDetailRequest) (*drainage.Welcome, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/select-default-welcome")
	if err != nil {
		return nil, err
	}
	welcome := &drainage.Welcome{}
	sErr := stmt.Get(c, welcome, request.EnterpriseId)
	if sErr != nil {
		return nil, sErr
	}
	return welcome, nil
}

// DeleteWelcomeOne 删除一条欢迎语
func (d *DrainageRepo) DeleteWelcomeOne(c ctx.Context, request drainage.WelcomeDeleteRequest) (int64, error) {
	stmt, stmtErr := d.DB.Template("t_drainage_welcomes/delete-one-welcome")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, request)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// DeleteWelcomeMemberList 删除欢迎语相关员工
func (d *DrainageRepo) DeleteWelcomeMemberList(c ctx.Context, request drainage.WelcomeDeleteRequest) (int64, error) {
	stmt, stmtErr := d.DB.Template("t_drainage_welcomes/delete-list-welcome-member")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, request)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// ChannelEdit 渠道编辑
func (d *DrainageRepo) ChannelEdit(cc ctx.Context, req drainage.CreateContactChannelLabelsRequest) error {
	return d.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			// 先更新主表渠道表
			updateChannelData := drainage.ContactChannels{
				ID:           req.ID,
				EnterpriseID: req.EnterpriseID,
				Name:         req.Name,
				UpdatedAt:    time.Now().Unix(),
			}
			cErr := d.UpdateContactChannel(cc, updateChannelData)
			if cErr != nil {
				return cErr
			}
			ids := []int64{req.ID}
			// 删除原来的标签
			dErr := d.DeleteChannelLabelByChannelID(cc, ids, req.EnterpriseID)
			if dErr != nil {
				return dErr
			}
			// 渠道有关联的标签
			if len(req.Label) > 0 {
				// 再新增现在关联的标签
				var insertLabel []drainage.ContactChannelLabels
				for _, item := range req.Label {
					item.EnterpriseID = req.EnterpriseID
					item.ChannelID = req.ID
					item.UpdatedAt = time.Now().Unix()
					item.CreatedAt = time.Now().Unix()
					item.ID = uid.ID(cc)
					insertLabel = append(insertLabel, item)
				}
				iErr := d.ChannelsLabelCreate(cc, insertLabel)
				return iErr
			}
			return nil
		})

		if crash != nil {
			return crash
		}

		return e
	})
}

// DeleteChannelLabelByChannelID 根据渠道ID 删除关联标签
func (d *DrainageRepo) DeleteChannelLabelByChannelID(ctx ctx.Context, ids []int64, enterpriseID int64) error {
	if len(ids) <= 0 {
		return errs.ChannelNameNotNllError
	}
	mapping := make(map[string]interface{})
	mapping["ids"] = ids

	stmt, stmtErr := d.DB.Template("t_contact_channel_labels/delete-label-by-channel-id", &mapping)
	if stmtErr != nil {
		return stmtErr
	}
	_, execError := stmt.Exec(ctx, time.Now().Unix(), enterpriseID)
	return execError
}

// UpdateContactChannel 更新渠道
func (d *DrainageRepo) UpdateContactChannel(ctx ctx.Context, channels drainage.ContactChannels) error {
	stmt, stmtErr := d.DB.Template("t_contact_channels/update")
	if stmtErr != nil {
		return stmtErr
	}
	_, err := stmt.NamedExec(ctx, channels)
	return err
}

// GetWorkExternalUserRelationList 通过渠道添加员工数据
func (d *DrainageRepo) GetWorkExternalUserRelationList(c ctx.Context, request drainage.WorkExternalUserRelation) (*[]drainage.WorkExternalUserRelation, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/work-external-user-relation", d.Cond2Args(request))
	if err != nil {
		return nil, err
	}
	List := &[]drainage.WorkExternalUserRelation{}
	res, errStmt := stmt.NamedQuery(c, request)
	err = sqlplus.ScanAll(res, List)
	if err != nil {
		return nil, err
	}
	if errStmt != nil {
		return nil, err
	}
	return List, nil
}

func (d *DrainageRepo) GetDrainageContactChannelsList(c ctx.Context, request drainage.ContactChannelsData) (*[]drainage.ContactChannelsData, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/drainage-contact-channels-list", d.Cond2Args(request))
	if err != nil {
		return nil, err
	}
	List := &[]drainage.ContactChannelsData{}
	res, errStmt := stmt.NamedQuery(c, request)
	err = sqlplus.ScanAll(res, List)
	if err != nil {
		return nil, err
	}
	if errStmt != nil {
		return nil, err
	}
	return List, nil
}

//SelectLabelsByChannelID 通过渠道ID获取标签列表
func (d *DrainageRepo) SelectLabelsByChannelID(c ctx.Context, request *drainage.ContactChannelLabelsListRequest) ([]*drainage.ContactChannelLabels, error) {
	stmt, err := d.DB.Template("t_contact_channel_labels/drainage-contact-channel-labels-list")
	if err != nil {
		return nil, err
	}
	var List []*drainage.ContactChannelLabels
	res, errStmt := stmt.NamedQuery(c, request)
	if errStmt != nil {
		return nil, err
	}
	err = sqlplus.ScanAll(res, &List)
	if err != nil {
		return nil, err
	}
	return List, nil
}

func (d *DrainageRepo) SumAddMemberNumber(c ctx.Context, contactsID int64) error {
	stmt, err := d.DB.Template("t_drainage_contact/contacts-update-sum-add-member-number")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.Exec(c, contactsID)
	return sTmtErr
}

func (d *DrainageRepo) GetWorkExternalUserRelation(c ctx.Context, relation customer.ExternalUserRelation) (*customer.ExternalUserRelation, error) {
	stmt, err := d.DB.Template("t_work_external_user_relation/select-work-external-user-relation-one")
	if err != nil {
		return nil, err
	}
	res := &customer.ExternalUserRelation{}
	sErr := stmt.Get(c, res, relation.WorkExternalUserId, relation.State, relation.EnterpriseId)
	if sErr != nil && sErr != sql.ErrNoRows {
		return nil, sErr
	}
	return res, nil
}

func (d *DrainageRepo) GetSynChannelsList(c ctx.Context, request drainage.Channel) (*[]drainage.Channel, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/syn-channels-list")
	if err != nil {
		return nil, err
	}
	List := &[]drainage.Channel{}
	res, errStmt := stmt.NamedQuery(c, request)
	err = sqlplus.ScanAll(res, List)
	if err != nil {
		return nil, err
	}
	if errStmt != nil {
		return nil, err
	}
	return List, nil
}

func (d *DrainageRepo) GetSynContactList(c ctx.Context, request drainage.ContactQrCode) (*[]drainage.ContactQrCode, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/syn-contact-qr-code-list")
	if err != nil {
		return nil, err
	}
	List := &[]drainage.ContactQrCode{}
	res, errStmt := stmt.NamedQuery(c, request)
	err = sqlplus.ScanAll(res, List)
	if err != nil {
		return nil, err
	}
	if errStmt != nil {
		return nil, err
	}
	return List, nil
}

func (d *DrainageRepo) GetSynMemberId(c ctx.Context, request drainage.UcMember) (*[]drainage.UcMember, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/syn-get-member-id")
	if err != nil {
		return nil, err
	}
	List := &[]drainage.UcMember{}
	res, errStmt := stmt.NamedQuery(c, request)
	err = sqlplus.ScanAll(res, List)
	if err != nil {
		return nil, err
	}
	if errStmt != nil {
		return nil, err
	}
	return List, nil
}

// GetSynContacts 联系我详情
func (d *DrainageRepo) GetSynContacts(ctx ctx.Context, request drainage.Contacts) (*drainage.Contacts, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/syn-contacts-detail")
	if err != nil {
		return nil, err
	}
	contacts := &drainage.Contacts{}
	sErr := stmt.Get(ctx, contacts, request.ID)
	if sErr != nil {
		return nil, sErr
	}
	return contacts, nil
}

// GetSynChannelsInfo 渠道详情
func (d *DrainageRepo) GetSynChannelsInfo(ctx ctx.Context, request drainage.ContactChannels) (*drainage.ContactChannels, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/syn-channels-detail")
	if err != nil {
		return nil, err
	}
	channels := &drainage.ContactChannels{}
	sErr := stmt.Get(ctx, channels, request.ID)
	if sErr != nil {
		return nil, sErr
	}
	return channels, nil
}

// WelcomeChannelList 查询渠道列表选择项
func (d *DrainageRepo) WelcomeChannelList(c ctx.Context, request drainage.WelcomeValidateRequest) (*[]drainage.WelcomeChannelsList, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = request.EnterpriseId
	mapping["add_channel"] = request.AddChannel
	mapping["filter"] = request.Filter
	mapping["filter_count"] = len(request.Filter)
	stmt, err := d.DB.Template("t_drainage_welcomes/select-channel-list-welcomes", &mapping)
	if err != nil {
		return nil, err
	}

	list := &[]drainage.WelcomeChannelsList{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// CountWelcomeByChannelID 查询渠道配置的欢迎语个数
func (d *DrainageRepo) CountWelcomeByChannelID(ctx ctx.Context, req *drainage.WelcomeListRequest) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = req
	countStmt, err := d.DB.Template("t_drainage_welcomes/select-list-welcomes-count", &mapping)
	if err != nil {
		return 0, err
	}
	type countChannel struct {
		Count int64 `db:"count"`
	}
	dest := &countChannel{}
	err = countStmt.NamedGet(ctx, dest, req)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return dest.Count, nil
}

// UpdateWelcomeMsg 更新欢迎语
func (d *DrainageRepo) UpdateWelcomeMsg(ctx ctx.Context, req *drainage.Welcome) (int64, error) {
	stmt, err := d.DB.Template("t_drainage_welcomes/welcome-msg-update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, req)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// ExistRelationCount 查询留存好友数
func (d *DrainageRepo) ExistRelationCount(ctx ctx.Context, req *drainage.StatisticsAllReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_trace/select_exist_external_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// ExistRelationCountByChannel 查询渠道留存好友数
func (d *DrainageRepo) ExistRelationCountByChannel(ctx ctx.Context, req *drainage.StatisticsChannelReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_channel_trace/select_exist_external_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// IncreaseCustomerCount 新增好友总数
func (d *DrainageRepo) IncreaseCustomerCount(ctx ctx.Context, req *drainage.StatisticsAllReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_trace/select_increase_customer_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// IncreaseCustomerCountByChannel 渠道新增好友总数
func (d *DrainageRepo) IncreaseCustomerCountByChannel(ctx ctx.Context, req *drainage.StatisticsChannelReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_channel_trace/select_increase_customer_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// LossCustomerCount 流失好友总数
func (d *DrainageRepo) LossCustomerCount(ctx ctx.Context, req *drainage.StatisticsAllReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_trace/select_loss_customer_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// LossCustomerCountByChannel 流失好友总数
func (d *DrainageRepo) LossCustomerCountByChannel(ctx ctx.Context, req *drainage.StatisticsChannelReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_channel_trace/select_loss_customer_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// MemberDelCustomerCountOrChannel 员工删除客户总数
func (d *DrainageRepo) MemberDelCustomerCountOrChannel(ctx ctx.Context, req *drainage.StatisticsChannelReq, corpId string) (int64, error) {
	// 定义匿名结构体来查询
	reqSqlParams := struct {
		*drainage.StatisticsChannelReq
		CorpId string `json:"corp_id" db:"corp_id"`
	}{
		req,
		corpId,
	}
	stmt, err := d.DB.Template("t_relation_callback_records/select_member_del_customer_count", d.Cond2Args(&reqSqlParams))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, reqSqlParams)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// CustomerDelMemberCount 好友删员工
func (d *DrainageRepo) CustomerDelMemberCount(ctx ctx.Context, req *drainage.StatisticsAllReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_del_relation_channel_trace/select_customer_del_member_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// CustomerDelMemberCountByChannel 好友删员工
func (d *DrainageRepo) CustomerDelMemberCountByChannel(ctx ctx.Context, req *drainage.StatisticsChannelReq) (int64, error) {
	stmt, err := d.DB.Template("t_work_del_relation_channel_trace/select_customer_del_member_by_channel_count", d.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(ctx, countResponseInfo, req)
	if getErr != nil {
		return 0, getErr
	}
	return countResponseInfo.Count, nil
}

// GetActivityChannelIDs 获取活动关联渠道
func (d *DrainageRepo) GetActivityChannelIDs(c ctx.Context, enterpriseID int64) ([]drainage.ChannelID, error) {
	channels := make([]drainage.ChannelID, 0)
	err := db.Query(c, d.SQLRepo, "t_fission_activity/select-activity-channel-ids", &channels, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return channels, err
	}
	return channels, nil
}

// GetWelcomeChannelIDs 获取欢迎语关联渠道
func (d *DrainageRepo) GetWelcomeChannelIDs(c ctx.Context, enterpriseID int64) ([]drainage.ChannelID, error) {
	channels := make([]drainage.ChannelID, 0)
	err := db.Query(c, d.SQLRepo, "t_drainage_welcomes/select-welcome-channel-ids", &channels, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return channels, err
	}
	return channels, nil
}

// FixDrainageWelcomeContentRep 获取欢迎语关联渠道
func (d *DrainageRepo) FixDrainageWelcomeContentRep(c ctx.Context, enterpriseID int64) error {
	updatedAt := assist.Timestamp()
	updateSql := fmt.Sprintf("UPDATE `drainage_welcomes` SET `content` = ?,`updated_at` = %d WHERE `id` = ? ", updatedAt)
	var data []drainage.WelcomeContent
	baseSql := fmt.Sprintf("SELECT `id`,`content` FROM `drainage_welcomes` WHERE `enterprise_id` = %d", enterpriseID)
	err := d.DB.DB.Select(c, &data, baseSql)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		for _, row := range data {
			if row.Content == "" {
				continue
			}
			var con = make(map[string]interface{})
			var newCon = make([]map[string]interface{}, 0)
			err = json.Unmarshal([]byte(row.Content), &con)
			if err != nil {
				log.Info(c).Msgf("FixDrainageWelcomeContentRepConError %+v,%+v", err, row)
				continue
			}
			newCon = append(newCon, con)
			newConJson, _ := json.Marshal(newCon)
			_, err = d.DB.DB.Exec(c, updateSql, string(newConJson), row.Id)
			if err != nil {
				log.Error(c).Msgf("FixDrainageWelcomeContentRepExecError %+v,%+v", err, row)
				return err
			}
		}
	}
	return nil
}
