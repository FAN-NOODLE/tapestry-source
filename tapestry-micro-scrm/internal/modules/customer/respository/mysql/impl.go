package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/uid"

	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/sdk/work"
)

type CustomerRepo struct {
	*repos.SQLRepo
}

var workImpl *work.QyWorkImpl

// QueryEnterpriseUser 获取聚合用户
func (c *CustomerRepo) QueryEnterpriseUser(cc ctx.Context, enterpriseID int64, source *customer.SourceUserBase) (*customer.EnterpriseUser, error) {
	mapping := make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	if source.EuID != 0 {
		mapping["euid"] = source.EuID
		mapping["app_table"] = "wx_user_corp_external"
	} else if source.UnionID != "" {
		if source.CustomerType == customer.CustomerTypeExternal {
			mapping["app_table"] = "wx_user_corp_external"
			mapping["unionid"] = source.UnionID
		} else {
			mapping["app_table"] = "wx_user_union"
			mapping["unionid"] = source.UnionID
		}
	} else {
		mapping["appid"] = source.AppID
		mapping["openid"] = source.OpenID
		if source.ID > 0 {
			mapping["id"] = source.ID
		} else if source.OpenID == "" {
			return nil, nil
		}
		switch source.CustomerType {
		case customer.CustomerTypeShop:
			mapping["app_table"] = "wx_user_shop"
		case customer.CustomerTypeExternal:
			mapping["app_table"] = "wx_user_corp_external"
		case customer.CustomerTypeOpen:
			mapping["app_table"] = "wx_user_open"
		default:
			return nil, nil
		}
	}
	stmt, err := c.DB.Template("t_customer/enterprise_user", &mapping)
	if err != nil {
		return nil, err
	}

	var enterpriseUser struct {
		customer.EnterpriseUserBindDB
		customer.EnterpriseUserBase
		customer.EnterpriseUserLocation
		customer.EnterpriseUserAttr
		Remark   sql.NullString `db:"remark"`
		DBTags   []byte         `db:"tags"`
		DBGroups []byte         `db:"groups"`
	}

	if err = stmt.Get(cc, &enterpriseUser); err != nil {
		return nil, err
	}

	if enterpriseUser.BindTo != 0 && enterpriseUser.BindTo != enterpriseUser.ID {
		return c.QueryEnterpriseUser(cc, enterpriseID, &customer.SourceUserBase{EuID: enterpriseUser.BindTo})
	}

	user := &customer.EnterpriseUser{
		EnterpriseUserBind:     *enterpriseUser.ToUser(),
		EnterpriseUserBase:     enterpriseUser.EnterpriseUserBase,
		EnterpriseUserLocation: enterpriseUser.EnterpriseUserLocation,
		EnterpriseUserAttr:     enterpriseUser.EnterpriseUserAttr,
		Remark:                 enterpriseUser.Remark.String,
	}

	json.Unmarshal(enterpriseUser.DBTags, &user.Tags)
	json.Unmarshal(enterpriseUser.DBGroups, &user.Groups)

	return user, nil
}

// UpdateEnterpriseUserAttr 更新用户属性
func (c *CustomerRepo) UpdateEnterpriseUserAttr(cc ctx.Context, enterpriseID int64, id int64, attr *customer.EnterpriseUserAttr) error {
	var sets []string
	var args []interface{}
	if attr.MarketCount > 0 {
		sets = append(sets, "market_count = market_count + ?", "latest_market_time=?")
		args = append(args, attr.MarketCount, time.Now().Unix())
	}
	if attr.PayOrder > 0 {
		sets = append(sets, "pay_order=?", "total_fee=?")
		args = append(args, attr.PayOrder, attr.TotalFee)
	}
	if attr.RefundOrder > 0 {
		sets = append(sets, "refund_order=?", "refund_fee=?")
		args = append(args, attr.RefundOrder, attr.RefundFee)
	}
	if attr.Tags != nil {
		sets = append(sets, "tags=?")
		b, _ := json.Marshal(attr.Tags)
		args = append(args, string(b))
	}
	if attr.Groups != nil {
		sets = append(sets, "groups=?")
		b, _ := json.Marshal(attr.Groups)
		args = append(args, string(b))
	}

	if attr.LevelID != 0 {
		sets = append(sets, "level_id=?")
		args = append(args, attr.LevelID)
	}
	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("update micro_scrm_enterprise_%d.enterprise_user set %s where id=?", enterpriseID, strings.Join(sets, ","))
	_, err := c.DB.DB.Exec(cc, query, args...)
	return err
}

func (c *CustomerRepo) SelectEUidBYDepartment(cc ctx.Context, enterpriseID int64, departmentIDs []int64, page repos.PageVO) (int64, *[]*customer.EUidList, error) {
	cond := struct {
		EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id"`
		DepartmentIDs []int64 `json:"department_ids"`
	}{}
	cond.EnterpriseID = enterpriseID
	cond.DepartmentIDs = departmentIDs
	findOption := &repos.FindPageOption{
		TplPrefix: "t_customer/select-euid-by-department",
		Fn: func(ctx.Context, int64) interface{} {
			return &[]*customer.EUidList{}
		},
	}

	total, list, err := c.FindPage(cc, page, cond, findOption)

	if err != nil {
		return 0, nil, err
	}

	return total, list.(*[]*customer.EUidList), err

}

// UpdateEnterpriseUserBase 更新用户基本属性
func (c *CustomerRepo) UpdateEnterpriseUserBase(cc ctx.Context, enterpriseID int64, id int64, base *customer.EnterpriseUserBase) error {
	var sets []string
	var args []interface{}
	if base.Nick != "" {
		sets = append(sets, "nick=?")
		args = append(args, base.Nick)
	}
	if base.Name != "" {
		sets = append(sets, "name=?")
		args = append(args, base.Name)
	}
	if base.Avatar != "" {
		sets = append(sets, "avatar=?")
		args = append(args, base.Avatar)
	}
	if base.Birthday != "" {
		sets = append(sets, "birthday=?")
		args = append(args, base.Birthday)
	}
	if base.Sex != 0 {
		sets = append(sets, "sex=?")
		args = append(args, base.Sex)
	}
	if base.Mobile != "" {
		sets = append(sets, "mobile=?")
		args = append(args, base.Mobile)
	}
	if base.DeliveryMobile != "" {
		sets = append(sets, "delivery_mobile=?")
		args = append(args, base.DeliveryMobile)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("update micro_scrm_enterprise_%d.enterprise_user set %s where id=?", enterpriseID, strings.Join(sets, ","))
	_, err := c.DB.DB.Exec(cc, query, args...)
	return err
}

// BelongToMembers 客户所属员工列表
func (c *CustomerRepo) BelongToMembers(cc ctx.Context, externalUserIDs []int64, page repos.PageVO) (int64, *customer.ResponseMemberList, error) {

	cond := struct {
		WorkExternalUserID []int64
	}{}
	for _, id := range externalUserIDs {
		cond.WorkExternalUserID = append(cond.WorkExternalUserID, id)
	}
	if len(cond.WorkExternalUserID) == 0 {
		return 0, nil, nil
	}
	findOption := &repos.FindPageOption{
		TplPrefix: "t_customer/belong-to-members",
		Fn: func(ctx.Context, int64) interface{} {
			return &customer.ResponseMemberList{}
		},
	}

	total, list, err := c.FindPage(cc, page, cond, findOption)

	if err != nil {
		return 0, nil, err
	}

	return total, list.(*customer.ResponseMemberList), err
}

// BelongToChats 客户加入所属员工群列表
func (c *CustomerRepo) BelongToChats(cc ctx.Context, enterpriseID int64, externalUserIDs []string, memberIDs []int64, status int, pageVo repos.PageVO) (int64, *customer.BelongToGroupChatList, error) {
	cond := &struct {
		ExternalUserIDs []string
		MemberIDs       []int64
		Status          int
		EnterpriseID    int64
	}{
		MemberIDs:    memberIDs,
		Status:       status,
		EnterpriseID: enterpriseID,
	}

	if len(externalUserIDs) == 0 {
		return 0, &customer.BelongToGroupChatList{}, nil
	}

	for _, id := range externalUserIDs {
		cond.ExternalUserIDs = append(cond.ExternalUserIDs, id)
	}

	findPageOption := &repos.FindPageOption{
		TplPrefix: "t_customer/belong-to-chats",
		Fn: func(ctx.Context, int64) interface{} {
			return &customer.BelongToGroupChatList{}
		},
	}

	total, list, err := c.FindPage(cc, pageVo, cond, findPageOption)
	if err != nil {
		return 0, &customer.BelongToGroupChatList{}, err
	}

	return total, list.(*customer.BelongToGroupChatList), err
}

// BelongToChatsCount 查询客户群个数
func (c *CustomerRepo) BelongToChatsCount(cc ctx.Context, enterpriseID int64, externalUserIDs map[string]*customer.ThirdUserID, memberIDs []int64, Status int) (int64, error) {
	mapping := make(map[string]interface{})
	cond := &struct {
		EnterpriseID    int64
		ExternalUserIDs []string
		MemberIDs       []int64
		Status          int
	}{
		EnterpriseID: enterpriseID,
		MemberIDs:    memberIDs,
		Status:       Status,
	}

	if len(externalUserIDs) == 0 {
		return 0, nil
	}

	for _, v := range externalUserIDs {
		cond.ExternalUserIDs = append(cond.ExternalUserIDs, v.ThirdID)
	}

	mapping["cond"] = cond
	stmt, err := c.DB.Template("t_customer/belong-to-chats-count", &mapping)
	if err != nil {
		return 0, err
	}
	type selectFromDb struct {
		Count int64 `json:"count" db:"count"`
	}
	dest := &selectFromDb{}
	sErr := stmt.Get(cc, dest)
	if sErr != nil {
		return 0, sErr
	}
	return dest.Count, nil
}

// FindExternal 查询外部联系人信息
func (c *CustomerRepo) FindExternal(ctx ctx.Context, user customer.WorkExternalUser) (*customer.WorkExternalUser, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = user
	stmt, err := c.DB.Template("t_work_external_user/find-work-external-user-first", &mapping)
	if err != nil {
		return nil, err
	}
	workUser := &customer.WorkExternalUser{}
	errStmt := stmt.NamedGet(ctx, workUser, user)
	if errStmt != nil {
		return nil, errStmt
	}
	workUser.EnterpriseID = user.EnterpriseID
	return workUser, nil
}

// FindExternalRelation 查询外部联系人关系
func (c *CustomerRepo) FindExternalRelation(ctx ctx.Context, r customer.WorkUserMemberRelation) (*customer.WorkUserMemberRelation, error) {
	stmt, err := c.DB.Template("t_work_external_user_relation/find-work-external-user-relation")
	if err != nil {
		return nil, err
	}
	relation := &customer.WorkUserMemberRelation{}
	errStmt := stmt.NamedGet(ctx, relation, r)
	if errStmt != nil {
		return nil, errStmt
	}
	return relation, nil
}

// FindCustomerRelation 查询外部联系人关系数据
func (c *CustomerRepo) FindCustomerRelation(ctx ctx.Context, memberID int64, workExternalUserID int64, enterpriseID int64) (*customer.WorkCustomerRelationResponse, error) {
	stmt, err := c.DB.Template("t_work_external_user_relation/find-one")
	if err != nil {
		return nil, err
	}
	relation := &customer.WorkCustomerRelationResponse{}
	errStmt := stmt.Get(ctx, relation, enterpriseID, workExternalUserID, memberID)
	if errStmt != nil {
		if errStmt == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errStmt
	}
	return relation, nil
}

// AddExternalRelation 添加外部联系人关系
func (c *CustomerRepo) AddExternalRelation(ctx ctx.Context, r []customer.WorkUserMemberRelation) (int64, error) {
	log.Info(ctx).Msgf("AddExternalRelationBengin %+v", r)
	var workExternalUserRelations []work_external_user.WorkExternalUserRelationModel
	for _, relation := range r {
		var userRelation work_external_user.WorkExternalUserRelationModel
		userRelation.ID = relation.ID
		userRelation.EnterpriseID = relation.EnterpriseID
		userRelation.WorkExternalUserID = relation.WorkExternalUserID
		userRelation.MemberID = relation.MemberID
		userRelation.Remark = relation.Remark
		userRelation.Description = relation.Description
		userRelation.RemarkCorpName = relation.RemarkCorpName
		userRelation.RemarkMobiles = relation.RemarkMobiles
		userRelation.AddWay = uint(relation.AddWay)
		userRelation.State = relation.State
		userRelation.Createtime = int(relation.CreateTime)
		userRelation.OperUserid = relation.OperUserID
		userRelation.RelationStatus = int(relation.RelationStatus)
		userRelation.CreatedAt = relation.CreatedAt
		userRelation.UpdatedAt = relation.UpdatedAt
		userRelation.UserID = relation.UserID
		userRelation.Unionid = relation.UnionId
		workExternalUserRelations = append(workExternalUserRelations, userRelation)
	}
	tx := gormx.GetInstance().Model(work_external_user.WorkExternalUserRelationModel{}).CreateInBatches(workExternalUserRelations, len(workExternalUserRelations))
	if tx.Error != nil {
		return 0, tx.Error
	}
	log.Info(ctx).Msgf("AddExternalRelationRepeatCustomerCount %+v", r)
	if len(r) > 0 {
		go func() {
			fx.RunSafe(func() error {
				for _, val := range r {
					err := c.RepeatCustomerCount(ctx, val.EnterpriseID, val.ID)
					if err != nil {
						log.Error(ctx).Msgf("AddExternalRelation-RepeatCustomerCount err %+v", err)
					}
				}
				return nil
			})
		}()
	}
	return tx.RowsAffected, nil
}

// AddExternalRelationLog 添加外部联系人关系
func (c *CustomerRepo) AddExternalRelationLog(ctx ctx.Context, r []customer.WorkUserMemberRelation) (int64, error) {
	stmt, err := c.DB.Template("t_work_external_user_relation_log/insert")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, r)
	if errStmt != nil {
		return 0, errStmt
	}

	return res.RowsAffected()
}

// UpdateExternalRelation 更新外部联系人关系
func (c *CustomerRepo) UpdateExternalRelation(ctx ctx.Context, r customer.WorkUserMemberRelation) (int64, error) {
	log.Info(ctx).Msgf("UpdateExternalRelationBengin %+v", r)
	stmt, err := c.DB.Template("t_work_external_user_relation/update")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, r)
	if errStmt != nil {
		return 0, errStmt
	}
	log.Info(ctx).Msgf("UpdateExternalRelationRepeatCustomerCount %+v", r)
	// 更新外部联系人信息不需要更新重复数 2023.2.21
	// go func() {
	//	fx.RunSafe(func() error {
	//		err = c.RepeatCustomerCount(ctx, r.EnterpriseID, r.ID)
	//		if err != nil {
	//			log.Error(ctx).Msgf("UpdateExternalRelation-RepeatCustomerCount err %+v", err)
	//		}
	//		return nil
	//	})
	// }()
	return res.RowsAffected()
}

func (c *CustomerRepo) BatchDeleteUpdateExternalReleation(ctx ctx.Context, enterpriseID int64, lastUpdateAt int64) (int64, error) {
	result, err := c.DB.DB.Exec(ctx, "update work_external_user_relation set relation_status = ?, updated_at = ? where enterprise_id = ? and updated_at < ?",
		customer.WorkUserRelationStatusDeleteEachOther, time.Now().Unix(), enterpriseID, lastUpdateAt)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// BatchTagGroupByName 根据标签组名称批量查询标签组
func (c *CustomerRepo) BatchTagGroupByName(ctx ctx.Context, g []string, enterpriseID int64) ([]customer.WorkExternalTagGroup, error) {
	mapping := make(map[string]interface{})
	mapping["groupNames"] = g
	stmt, err := c.DB.Template("t_work_external_tag_group/find-group-by-groupnames", mapping)

	if err != nil {
		return nil, err
	}
	res, errStmt := stmt.Query(ctx, enterpriseID, 0)
	if errStmt != nil {
		return nil, errStmt
	}
	defer res.Close()
	var list []customer.WorkExternalTagGroup
	sErr := sqlplus.ScanAll(res, list)
	if sErr != nil {
		return nil, sErr
	}
	return list, nil
}

// AddTagGroup 添加标签组
func (c *CustomerRepo) AddTagGroup(ctx ctx.Context, t customer.WorkExternalTagGroup) (int64, error) {
	stmt, err := c.DB.Template("t_work_external_user_relation/insert")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, t)
	if errStmt != nil {
		return 0, errStmt
	}

	return res.LastInsertId()
}

// FindWorkGroupChatDetail 群聊详情
func (c *CustomerRepo) FindWorkGroupChatDetail(ctx ctx.Context, req customer.WorkGroupChat) (*customer.WorkGroupChat, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = req
	stmt, err := c.DB.Template("t_work_group_chat/detail", &mapping)
	if err != nil {
		return nil, err
	}
	chat := &customer.WorkGroupChat{}
	sErr := stmt.NamedGet(ctx, chat, req)

	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, sErr
	}

	return chat, nil
}

// AddWorkGroupChat 新建群聊
func (c *CustomerRepo) AddWorkGroupChat(ctx ctx.Context, req customer.WorkGroupChat) (int64, error) {
	stmt, err := c.DB.Template("t_work_group_chat/insert")
	if err != nil {
		return 0, err
	}
	_, sErr := stmt.NamedExec(ctx, req)
	if sErr != nil {
		return 0, sErr
	}
	return req.ID, nil
}

// UpdateWorkGroupChat 更新群聊
func (c *CustomerRepo) UpdateWorkGroupChat(ctx ctx.Context, req customer.WorkGroupChat) error {
	stmt, err := c.DB.Template("t_work_group_chat/update")
	if err != nil {
		return err
	}
	_, sErr := stmt.NamedExec(ctx, req)
	if sErr != nil {
		return sErr
	}
	return nil
}

// FindWorkGroupChatMember 查询群成员
func (c *CustomerRepo) FindWorkGroupChatMember(ctx ctx.Context, req customer.WorkGroupChatMember) (*customer.WorkGroupChatMember, error) {
	stmt, err := c.DB.Template("t_work_group_chat_member/find-work-group-chat-member")
	if err != nil {
		return nil, err
	}
	chatMember := &customer.WorkGroupChatMember{}
	sErr := stmt.NamedGet(ctx, chatMember, req)
	if sErr != nil {
		return nil, sErr
	}
	return chatMember, nil
}

// AddWorkGroupChatMember 新增群成员
func (c *CustomerRepo) AddWorkGroupChatMember(ctx ctx.Context, req []customer.WorkGroupChatMember) (int64, error) {
	stmt, err := c.DB.Template("t_work_group_chat_member/insert")
	if err != nil {
		return 0, err
	}
	var result int64
	for _, val := range req {
		_, sErr := stmt.NamedExec(ctx, val)
		if sErr != nil {
			result++
		}
	}
	return result, nil
}

func (c *CustomerRepo) BatchFindWorGroupChatMember(ctxContext ctx.Context, req customer.BatchWorkGroupMemberRequest) (int64, []customer.WorkGroupChatMember, error) {
	offset := (req.Page - 1) * req.Size
	pageVo := repos.PageVO{Offset: offset, Limit: req.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_work_group_chat_member/select",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]customer.WorkGroupChatMember{}
		},
	}
	total, list, err := c.FindPage(ctxContext, pageVo, req, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]customer.WorkGroupChatMember)
	return total, *resList, nil
}

// BatchDeleteWorkGroupChatMembers 批量删除
func (c *CustomerRepo) BatchDeleteWorkGroupChatMembers(ctx ctx.Context, ids []int64) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := c.DB.Template("t_work_group_chat_member/batch-deleted-by-ids", &mapping)
	if err != nil {
		return err
	}

	_, stmtErr := stmt.NamedExec(ctx, nil)
	if stmtErr != nil {
		return stmtErr
	}
	return nil
}

// MemberMultiDepartments 根据memberIDs查询所在部门信息
func (c *CustomerRepo) MemberMultiDepartments(cc ctx.Context, memberIDs []int64) (map[int64]customer.MemberDepartsList, error) {
	memberIDsCond := map[string]interface{}{
		"memberIDs": memberIDs,
	}

	memberDepartsList := &customer.MemberDepartsList{}
	memberOption := &repos.ListOption{
		TplPrefix: "t_departments/select-dept-by-memberids",
		Fn: func(context ctx.Context, i int64) interface{} {
			return memberDepartsList
		},
	}

	departmentListById := make(map[int64]customer.MemberDepartsList)
	_, err := c.List(cc, memberIDsCond, memberOption)
	if err != nil {
		return departmentListById, err
	}

	for _, department := range *memberDepartsList {
		departmentListById[department.MemberID] = append(departmentListById[department.MemberID], department)
	}

	return departmentListById, err
}

// FindTagGroupFirst 查询标签组信息
func (c *CustomerRepo) FindTagGroupFirst(ctx ctx.Context, req customer.WorkExternalTagGroup) (*customer.WorkExternalTagGroup, error) {
	stmt, err := c.DB.Template("t_work_external_tag_group/find-tag-group-first")
	if err != nil {
		return nil, err
	}
	dest := &customer.WorkExternalTagGroup{}
	sErr := stmt.NamedGet(ctx, dest, req)
	if sErr != nil {
		return nil, sErr
	}
	return dest, nil
}

// UpdateTagGroup 更新标签组
func (c *CustomerRepo) UpdateTagGroup(ctx ctx.Context, req customer.WorkExternalTagGroup) error {
	stmt, err := c.DB.Template("t_work_external_tag_group/update")
	if err != nil {
		return err
	}
	res, sErr := stmt.NamedExec(ctx, req)
	if sErr != nil {
		return sErr
	}
	_, rErr := res.RowsAffected()
	return rErr
}

// BatchUpdateWorkGroupChatMembers 批量更新群聊表
func (c *CustomerRepo) BatchUpdateWorkGroupChatMembers(ctx ctx.Context, where map[string]interface{}, update customer.WorkGroupChatMember) (int64, error) {
	stmt, err := c.DB.Template("t_work_group_chat_member/batch-update", &where)
	if err != nil {
		return 0, err
	}
	res, sErr := stmt.NamedExec(ctx, update)
	if sErr != nil {
		return 0, sErr
	}
	return res.RowsAffected()
}

// AddCustomerCallbackRecord 新增回调记录
func (c *CustomerRepo) AddCustomerCallbackRecord(ctx ctx.Context, req *work.MemberExternalCallback) (int64, error) {
	stmt, err := c.DB.Template("t_relation_callback_records/insert")
	if err != nil {
		return 0, err
	}
	res, sErr := stmt.NamedExec(ctx, req)
	if sErr != nil {
		return 0, sErr
	}
	_, rErr := res.RowsAffected()
	if rErr != nil {
		return 0, rErr
	}

	return req.ID, nil
}

// UpdateCustomerCallbackRecord 更新回调记录
func (c *CustomerRepo) UpdateCustomerCallbackRecord(ctx ctx.Context, req work.MemberExternalCallback) (int64, error) {
	stmt, err := c.DB.Template("t_relation_callback_records/update", c.Cond2Args(req))
	if err != nil {
		return 0, err
	}
	res, sErr := stmt.NamedExec(ctx, req)
	if sErr != nil {
		return 0, sErr
	}
	return res.RowsAffected()
}

// QueryCustomerCallbackRecord 获取回调记录
func (c *CustomerRepo) QueryCustomerCallbackRecord(cc ctx.Context, page int, size int, cond work.MemberExternalCallback) ([]*work.MemberExternalCallback, error) {
	offset := (page - 1) * size
	limit := size
	option := &repos.FindPageOption{
		TplPrefix: "t_relation_callback_records/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*work.MemberExternalCallback{}
		},
	}
	_, list, err := c.FindPage(cc, repos.PageVO{Offset: offset, Limit: limit}, cond, option)
	if err != nil {
		return nil, err
	}
	data := list.(*[]*work.MemberExternalCallback)
	return *data, nil
}

// CallbackRecordDetail 获取回调记录
func (c *CustomerRepo) CallbackRecordDetail(cc ctx.Context, cond *work.MemberExternalCallback) (*work.MemberExternalCallback, error) {
	stmt, err := c.DB.Template("t_relation_callback_records/detail-by-relation", c.Cond2Args(cond))
	if err != nil {
		return nil, err
	}

	detail := new(work.MemberExternalCallback)
	sErr := stmt.NamedGet(cc, detail, cond)
	// 未查询到数据
	if sErr == sql.ErrNoRows {
		return nil, nil
	}

	if sErr != nil {
		return nil, sErr
	}
	return detail, nil
}

// DelCallbackRecord 查询删除的数据
func (c *CustomerRepo) DelCallbackRecord(cc ctx.Context, page int, size int, corpID string) ([]*work.MemberExternalCallback, error) {
	offset := (page - 1) * size
	limit := size
	cond := make(map[string]interface{})
	cond["corpID"] = corpID
	cond["offset"] = offset
	cond["limit"] = limit
	stmt, err := c.DB.Template("t_relation_callback_records/del-list", c.Cond2Args(cond))

	if err != nil {
		return nil, err
	}
	list := make([]*work.MemberExternalCallback, 0, size)
	res, err := stmt.NamedQuery(cc, cond)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	sqlErr := sqlplus.ScanAll(res, &list)
	if sqlErr != nil {
		return nil, err
	}
	return list, nil
}

// CallbackRecordDetail 获取回调记录
func (c *CustomerRepo) LastRecordDetail(cc ctx.Context, cond *work.MemberExternalCallback) (*work.MemberExternalCallback, error) {
	stmt, err := c.DB.Template("t_relation_callback_records/last-by-relation", c.Cond2Args(cond))
	if err != nil {
		return nil, err
	}

	detail := new(work.MemberExternalCallback)
	sErr := stmt.NamedGet(cc, detail, cond)
	// 未查询到数据
	if sErr == sql.ErrNoRows {
		return nil, nil
	}

	if sErr != nil {
		return nil, sErr
	}
	return detail, nil
}

// BatchExternalRelations 批量查询外部联系关系列表
func (c *CustomerRepo) BatchExternalRelations(cc ctx.Context, req customer.BatchWorkUserMemberRelationRequest) ([]customer.BatchWorkUserMemberRelationResponse, error) {
	var relations []customer.BatchWorkUserMemberRelationResponse
	stmt, err := c.DB.Template("t_work_external_user_relation/select", &map[string]interface{}{"cond": req})
	if err != nil {
		return nil, err
	}

	if err = stmt.NamedSelect(cc, &relations, &req); err != nil {
		return nil, err
	}

	return relations, nil
}

// GetExternalTagsOfficial 获取标签
func (c *CustomerRepo) GetExternalTagsOfficial(ctx ctx.Context, enterpriseID int64, tagIDList []string) ([]*customer.WorkExternalTag, error) {
	mapping := make(map[string]interface{})
	mapping["tag_id_list"] = tagIDList
	mapping["enterprise_id"] = enterpriseID
	mapping["deleted_at"] = 0
	stmt, err := c.DB.Template("t_work_external_tag/query-by-tag-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.WorkExternalTag, 0)
	err = stmt.NamedSelect(ctx, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetExternalUserTagsByTagId 获取指定关系和标签数据
func (c *CustomerRepo) GetExternalUserTagsByTagId(ctx ctx.Context, enterpriseId, relationId int64, tagIds []int64) ([]int64, error) {
	// selectSql := "SELECT external_tag_id FROM work_external_user_tag WHERE enterprise_id = ? AND relation_id = ? AND external_tag_id IN(?) AND deleted_at = 0;"
	var result []int64
	err := gormx.GetInstance().Table("work_external_user_tag").
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"relation_id":   relationId,
			"deleted_at":    0,
		}).Where("external_tag_id IN(?)", tagIds).Pluck("external_tag_id", &result).Error
	return result, err
}

// AddExternalTag 添加外部联系人标签
func (c *CustomerRepo) AddExternalTag(ctx ctx.Context, data []*customer.WorkExternalTag) (int64, error) {
	stmt, err := c.DB.Template("t_work_external_tag/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetExternalTagGroups 查询标签组数据
func (c *CustomerRepo) GetExternalTagGroups(ctx ctx.Context, enterpriseID int64, tagGroupIDList []string) ([]*customer.WorkExternalTagGroup, error) {
	mapping := make(map[string]interface{})
	mapping["group_id_list"] = tagGroupIDList
	mapping["enterprise_id"] = enterpriseID
	mapping["deleted_at"] = 0
	stmt, err := c.DB.Template("t_work_external_tag_group/query-by-group-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.WorkExternalTagGroup, 0)
	err = stmt.NamedSelect(ctx, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AddExternalTagGroup 添加标签组
func (c *CustomerRepo) AddExternalTagGroup(ctx ctx.Context, data []*customer.WorkExternalTagGroup) (int64, error) {
	stmt, err := c.DB.Template("t_work_external_tag_group/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (c *CustomerRepo) AddExternalUserTags(ctx ctx.Context, data []*customer.WorkExternalUserTag) (int64, error) {
	// 查询已存在管理
	stmt, err := c.DB.Template("t_work_external_user_tag/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryExternalUserTags 查询外部联系人标签列表
func (c *CustomerRepo) QueryExternalUserTags(ctx ctx.Context, relationID int64) ([]*customer.WorkExternalUserTag, error) {
	var list []*customer.WorkExternalUserTag
	stmt, err := c.DB.Template("t_work_external_user_tag/query-tags")
	if err != nil {
		return nil, err
	}
	res, errStmt := stmt.Query(ctx, relationID, 0)
	if errStmt != nil {
		log.Info(ctx).Msgf("查询标签数据失败:%d, err:%v", relationID, errStmt)
		return nil, errStmt
	}
	err = sqlx.StructScan(res, &list)
	if err != nil {
		log.Info(ctx).Msgf("读取标签数据失败: %d, err:%s", relationID, err.Error())
		return nil, err
	}
	return list, nil
}

// GetExternalUserRelation 获取外部联系人关系
func (c *CustomerRepo) GetExternalUserRelation(ctx ctx.Context, enterpriseID int64, workExternalUserID int64, memberUserID string) (*customer.WorkUserMemberRelation, error) {
	var relation customer.WorkUserMemberRelation
	stmt, err := c.DB.Template("t_work_external_user_relation/query-by-member-user-id")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Get(ctx, &relation, enterpriseID, workExternalUserID, memberUserID)
	if errStmt != nil {
		log.Info(ctx).Msgf("查询外部联系人关系失败:%v, err:%v", workExternalUserID, errStmt)
		return nil, errStmt
	}
	return &relation, nil
}

// GetExternalTagsByID 根据id获取标签数据
func (c *CustomerRepo) GetExternalTagsByID(ctx ctx.Context, idList []int64) ([]*customer.WorkExternalTag, error) {
	if len(idList) == 0 {
		return nil, nil
	}
	mapping := make(map[string]interface{})
	mapping["id_list"] = idList
	mapping["deleted_at"] = 0
	stmt, err := c.DB.Template("t_work_external_tag/query-by-ids", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.WorkExternalTag, 0)
	err = stmt.NamedSelect(ctx, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteExternalUserTags 删除外部联系人标签
func (c *CustomerRepo) DeleteExternalUserTags(ctx ctx.Context, idList []int64) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["id_list"] = idList
	mapping["now"] = time.Now().Unix()
	stmt, err := c.DB.Template("t_work_external_user_tag/delete", &mapping)
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, mapping)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// QueryCustomerTags 查询外部联系人自定义标签
func (c *CustomerRepo) QueryCustomerTags(ctx ctx.Context, enterpriseID int64, memberID int64, tags []string) ([]*customer.WorkExternalTag, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	mapping := make(map[string]interface{})
	mapping["tag_list"] = tags
	mapping["enterprise_id"] = enterpriseID
	mapping["member_id"] = memberID
	mapping["deleted_at"] = 0
	mapping["type"] = customer.ExternalTagTypeCustomer
	stmt, err := c.DB.Template("t_work_external_tag/query-customer-tags", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.WorkExternalTag, 0)
	err = stmt.NamedSelect(ctx, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetCustomerOfficialAccounts 查询用户关注的公众号列表
func (c *CustomerRepo) GetCustomerOfficialAccounts(ctx ctx.Context, enterpriseID int64, wxOpen map[string]*customer.ThirdUserID) ([]*customer.DetailOfficialAccount, error) {
	var userids []int64
	for _, u := range wxOpen {
		userids = append(userids, u.ID)
	}
	mapping := make(map[string]interface{})
	mapping["user_ids"] = userids
	mapping["enterprise_id"] = enterpriseID
	stmt, err := c.DB.Template("t_customer/select-user-open", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.DetailOfficialAccount, 0)
	err = stmt.NamedSelect(ctx, &list, mapping)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// ExternalListRemark 查询联系人列表备注
func (c *CustomerRepo) ExternalListRemark(cc ctx.Context, enterpriseID int64, memberID int64, externalIDList []int64) ([]customer.ExternalRemark, error) {
	if len(externalIDList) == 0 {

		return []customer.ExternalRemark{}, nil
	}
	option := &repos.ListOption{
		TplPrefix: "t_work_external_user_relation/external-list-remark",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]customer.ExternalRemark{}
		},
	}
	cond := struct {
		EnterpriseID   int64
		MemberID       int64
		ExternalIDList []int64
	}{
		EnterpriseID:   enterpriseID,
		MemberID:       memberID,
		ExternalIDList: externalIDList,
	}
	res, err := c.List(cc, cond, option)

	return *res.(*[]customer.ExternalRemark), err
}

// UpdateWorkChatMemberInfo 更新企微群聊表
func (c *CustomerRepo) UpdateWorkChatMemberInfo(cc ctx.Context, entity *customer.WorkGroupChatMember) (int64, error) {
	stmt, err := c.DB.Template("t_work_group_chat_member/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(cc, entity)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetUserRemark 客户备注信息查询
func (c *CustomerRepo) GetUserRemark(cc ctx.Context, data customer.ReturnUID) (info *customer.RemarkInfo, err error) {

	var entity customer.RemarkInfo
	stmt, err := c.DB.Template("t_customer_remark_info/customer-remark-info")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Get(cc, &entity, data.UserID)
	if errStmt != nil {
		return nil, errStmt
	}

	return &entity, nil
}

// RemarkSave 客户备注信息保存
func (c *CustomerRepo) RemarkSave(cc ctx.Context, data customer.RemarkInfo) (id int64, err error) {

	cond := c.Cond2Args(data)

	stmt, err := c.DB.Template("t_customer_remark_info/customer-remark-save", cond)
	if err != nil {
		return 0, err
	}
	res, ok := stmt.NamedExec(cc, data)
	if ok != nil || res == nil {
		return 0, ok
	}

	return data.UserID, nil
}

// SetExternalTags 设置外部联系人企业标签并同步到企微
func (c *CustomerRepo) SetExternalTags(cc ctx.Context, params customer.SetExternalTags) (int64, error) {

	_ = c.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {

			// 获取关系数据
			relationInfo, error := c.ExternalRelation(cc, params.EnterpriseID, params.RelationID)
			fmt.Println(relationInfo, error)
			if error != nil && relationInfo == nil {
				return error
			}

			// 企微请求数据
			req := work.UserTagHandlerRequest{
				Userid:         relationInfo.UserID,
				ExternalUserid: relationInfo.ExternalUserID,
				AddTag:         params.AddTag,
				RemoveTag:      params.RemoveTag,
			}

			// 企微数据同步
			if len(req.AddTag) > 0 || len(req.RemoveTag) > 0 {
				// 查询企业信息
				var workImpl work.QyWorkImpl
				corpReq := enterprise.Corp{
					EnterpriseID: params.EnterpriseID,
				}
				enterpriseService := enterprise.New(enterpriseMysql.New())
				enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(cc), corpReq)
				if enterpriseErr != nil {
					log.Error(ctx.Wrap(cc)).Msgf("查询企业信息失败，错误信息：%+v", enterpriseErr)
					return enterpriseErr
				}
				workImpl.SetDomain(enterpriseData.Domain)

				isOK, err := workImpl.UserTagHandler(cc, req)
				if err != nil || isOK == false {
					return err
				}
			}

			// 新增的外部联系人标签
			if params.AddTag != nil {
				// 根据企微标签获取标签表ID
				addTagList, tagErr := c.GetExternalTagsOfficial(cc, params.EnterpriseID, params.AddTag)
				if addTagList != nil && tagErr == nil {
					var addTagIds []int64
					for _, tag := range addTagList {
						addTagIds = append(addTagIds, tag.ID)
					}
					externalUserTagIds, err := c.GetExternalUserTagsByTagId(cc, params.EnterpriseID, params.RelationID, addTagIds)
					if err != nil {
						return err
					}
					var addList []*customer.WorkExternalUserTag
					for _, tag := range addTagList {
						if !helper.InArray[int64](tag.ID, externalUserTagIds) {
							tagData := &customer.WorkExternalUserTag{
								ID:            uid.ID(cc),
								EnterpriseID:  params.EnterpriseID,
								RelationID:    params.RelationID,
								ExternalTagId: tag.ID,
								CreatedAt:     assist.Timestamp(),
								UpdatedAt:     assist.Timestamp(),
							}
							addList = append(addList, tagData)
						}
					}
					if len(addList) > 0 {
						_, err = c.AddExternalUserTags(cc, addList)
						if err != nil {
							return err
						}
					}
				}
			}

			// 删除的外部联系人标签
			if params.RemoveTag != nil {
				removeTagList, tagErr := c.GetExternalTagsOfficial(cc, params.EnterpriseID, params.RemoveTag)
				if removeTagList != nil && tagErr == nil {
					var delTagIDs []int64
					for _, tag := range removeTagList {
						delTagIDs = append(delTagIDs, tag.ID)
					}

					_, error := c.RemoveTags(cc, params.RelationID, delTagIDs)

					if error != nil {
						return error
					}
				}
			}

			return nil
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return params.RelationID, nil
}

// formatTags 格式化标签
func (c *CustomerRepo) formatTags(params []*customer.WorkExternalTag) []string {

	var tagData []string
	for _, tag := range params {
		tagData = append(tagData, tag.TagID)
	}

	return tagData
}

// ExternalRelation 查询企业外部联系人关系
func (c *CustomerRepo) ExternalRelation(cc ctx.Context, enterpriseID int64, relationID int64) (*customer.WorkUserMemberRelationData, error) {

	var entity customer.WorkUserMemberRelationData
	stmt, err := c.DB.Template("t_work_external_user_relation/detail", &map[string]interface{}{"enterprise_id": enterpriseID})
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Get(cc, &entity, relationID)
	if errStmt != nil {
		return nil, errStmt
	}

	return &entity, nil
}

// RemoveTags 移除企业外部联系人标签
func (c *CustomerRepo) RemoveTags(cc ctx.Context, relationID int64, externalTagID []int64) (int64, error) {

	mapping := make(map[string]interface{})
	mapping["now"] = assist.Timestamp()
	mapping["external_id"] = externalTagID
	stmt, err := c.DB.Template("t_work_external_user_tag/remove-tags", &mapping)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(cc, relationID)
	if err != nil && res != nil {
		return 0, err
	}

	return relationID, nil
}

// FindExternalBindMembers 根据外部联系人信息查询绑定的成员信息
func (c *CustomerRepo) FindExternalBindMembers(cc ctx.Context, params *customer.ExternalRelationByUserIdRequest) (int64, *customer.WorkUserMemberRelationMemberIds, error) {
	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_customer/find-work-external-user-relation-by-userid",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &customer.WorkUserMemberRelationMemberIds{}
		},
	}
	total, list, err := c.FindPage(cc, pageVo, params, option)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	resList := list.(*customer.WorkUserMemberRelationMemberIds)

	return total, resList, nil
}

// HasNormalRelation 是否存在正常的关系
func (c *CustomerRepo) HasNormalRelation(ctx ctx.Context, workExternalUserID int64, enterpriseId, channelId int64) (bool, error) {
	req := make(map[string]interface{})
	req["work_user_id"] = workExternalUserID
	req["enterprise_id"] = enterpriseId
	req["channel_id"] = channelId
	stmt, err := c.DB.Template("t_customer/has-normal-relation", &req)
	if err != nil {
		return false, err
	}

	data := new(customer.WorkUserMemberRelationData)
	err = stmt.NamedGet(ctx, data, req)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if data.ID == 0 {
		return false, nil
	}

	return true, nil
}

// ExistsRelationByStartTime 查询在开始时间之前，客户是否添加过这些员工
func (c *CustomerRepo) ExistsRelationByStartTime(ctx ctx.Context, req *customer.ExistsRelationByTimeReq) (bool, error) {
	stmt, err := c.DB.Template("t_customer/select-exists-relation-by-time", c.Cond2Args(req))
	if err != nil {
		return false, err
	}

	data := new(customer.WorkUserMemberRelationData)
	err = stmt.NamedGet(ctx, data, req)

	switch true {
	case err == sql.ErrNoRows || data == nil:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

// FindFirstJoinTimeInGroups 获取用户和群最早添加的关系
func (c *CustomerRepo) FindFirstJoinTimeInGroups(ctx ctx.Context, externalUserID string, groupIDs []int64, enterpriseID int64) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["ExternalUserID"] = externalUserID
	mapping["GroupIDs"] = groupIDs
	mapping["EnterpriseID"] = enterpriseID
	stmt, err := c.DB.Template("t_work_group_chat_member/find-first-join-time", c.Cond2Args(mapping))
	if err != nil {
		return 0, err
	}

	data := &struct {
		JoinTime int64 `db:"join_time"`
	}{}
	err = stmt.Get(ctx, data)

	switch true {
	case err == sql.ErrNoRows || data == nil:
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return data.JoinTime, nil
	}
}

// AllEterpriseUserForEsList 将初始数据存到es中
func (c *CustomerRepo) AllEterpriseUserForEsList(cc ctx.Context, enterpriseId int64, beginTime, endTime string, pageSize int, elist chan<- []customer.EnterpriseUserEntity) {
	defer close(elist)

	for {
		var euserList = []customer.EnterpriseUserEntity{}
		sql := fmt.Sprintf("SELECT `id`,`wx_union`,`wx_shop`,`wx_external`,`wx_open`,`tags`,`groups`,`name`,`nick`,`avatar`,`mobile`, `delivery_mobile`, `birthday`,"+
			"`sex`,`language`,`country`,`province`,`city`,`remark`,`version`,`pay_order`,"+
			"`total_fee`,`refund_order`,`refund_fee`,`per_customer_transaction`,`latest_market_time`,`market_count`,`bind_to`,`level_id`,`created_at`,"+
			"`updated_at` FROM micro_scrm_enterprise_%d.`enterprise_user` WHERE `es_update` IS NULL || ( `dt_update` >= \"%s\" AND `dt_update` <= \"%s\"  AND `es_update` != `dt_update`) ORDER BY `id` ASC LIMIT  %d ", enterpriseId, beginTime, endTime, pageSize)
		fmt.Println("sql =", sql, beginTime, endTime)
		err := c.DB.Select(cc, &euserList, sql)
		if err != nil {
			log.Error(cc).Msgf("AllEterpriseUserForEsList err %v", err)
			return
		}
		if len(euserList) > 0 {
			elist <- euserList
		}
		if len(euserList) < pageSize {
			log.Info(cc).Msgf("AllEterpriseUserForEsList over %d", len(euserList))
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type WorkRepo struct {
	*repos.SQLRepo
}

func (w *WorkRepo) QueryExternalRelationInfo(c ctx.Context, params customer.QueryRelationParams) (*customer.ExternalUserRelationEntity, error) {
	// 构建sql查询mapping
	queryMap := make(map[string]interface{})
	queryMap["enterprise_id"] = params.EnterpriseID
	queryMap["external_userid"] = params.ExternalUserID
	queryMap["member_id"] = params.MemberID

	// 构建查询模版
	stmt, err := w.DB.Template("t_work_external_user_relation/select-relation-info", &queryMap)
	if err != nil {
		return nil, err
	}
	// 查询
	var data = &customer.ExternalUserRelationEntity{}
	if err = stmt.NamedGet(c, data, queryMap); err != nil {
		return nil, err
	}

	return data, nil
}

// QueryExternalRelationRemark 查询外部联系人关系及描述信息
func (w *WorkRepo) QueryExternalRelationRemark(c ctx.Context, relationID int64, enterpriseID int64) (*customer.WorkExternalUserRemark, error) {
	mapping := make(map[string]interface{})
	mapping["id"] = relationID
	mapping["enterprise_id"] = enterpriseID
	stmt, err := w.DB.Template("t_work_external_user_relation/query-relation-for-edit", &mapping)
	if err != nil {
		return nil, err
	}

	user := &customer.WorkExternalUserRemark{}

	if err = stmt.NamedGet(c, user, mapping); err != nil {
		return nil, err
	}

	return user, nil
}

// QueryExternalRelationRemarks 查询外部联系人关系及描述信息
func (w *WorkRepo) QueryExternalRelationRemarks(c ctx.Context, enterpriseID int64, workExternalUserIDs []int64) ([]*customer.WorkExternalUserRemark, error) {
	mapping := make(map[string]interface{})
	mapping["user_id_list"] = workExternalUserIDs
	mapping["enterprise_id"] = enterpriseID
	stmt, err := w.DB.Template("t_work_external_user_relation/query_relation_remarks", &mapping)
	if err != nil {
		return nil, err
	}
	var list = make([]*customer.WorkExternalUserRemark, 0)
	err = stmt.NamedSelect(c, &list, mapping)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// UpdateExternalRelationRemark 更新外部联系人关系备注
func (w *WorkRepo) UpdateExternalRelationRemark(c ctx.Context, id int64, remark string) error {
	stmt, err := w.DB.Template("t_work_external_user_relation/update-desc")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c, remark, time.Now().Unix(), id)
	if err != nil {
		return err
	}
	return nil
}

func (w *WorkRepo) QueryWorkUserByUserIDs(c ctx.Context, enterpriseID int64, userIDs []int64) (*[]*customer.WorkExternalUser, error) {
	var resList = make([]*customer.WorkExternalUser, 0)
	var mapping = make(map[string]interface{})
	mapping["enterprise_id"] = enterpriseID
	mapping["user_ids"] = userIDs
	err := db.NamedSelect(c, w.SQLRepo, "t_work_external_user/query-by-userids", &resList, mapping, &mapping)
	if err != nil {
		return nil, err
	}
	return &resList, nil
}

func (c *CustomerRepo) UpdateEsResult(cc ctx.Context, esIndex, dtUpdate string, ids []string) {
	if len(ids) > 0 {
		params := strings.Join(ids, "','")
		var updateSQL = "UPDATE %s SET `es_update` = '%s',`dt_update` = '%s' WHERE `id` IN ('%s'); "
		updateSQL = fmt.Sprintf(updateSQL, esIndex, dtUpdate, dtUpdate, params)
		_, err := db.DB().DB.Exec(ctx.New(), updateSQL)
		log.Info(cc).Msgf("UpdateEsResult -> %s,err = %v", updateSQL, err)
	}
}

// GetExternalUserID 获取指定userid的外部联系人user_id
func (c *CustomerRepo) GetExternalUserID(ctx ctx.Context, userId int64, enterpriseId int64) (int64, string, error) {
	data := &struct {
		Id             int64  `db:"id"`
		ExternalUserID string `db:"external_userid"`
	}{}

	mapping := make(map[string]interface{})

	mapping["user_id"] = userId
	mapping["enterprise_id"] = enterpriseId
	stmt, err := c.DB.Template("t_work_external_user/select-user-by-user-id", &mapping)
	if err != nil {
		return 0, "", err
	}

	stmtErr := stmt.NamedGet(ctx, data, mapping)
	if stmtErr != nil {
		return 0, "", stmtErr
	}
	return data.Id, data.ExternalUserID, nil
}

func (c *CustomerRepo) GetExternalUserUsedTag(ctx ctx.Context, enterpriseId int64, tagID int64) ([]customer.WorkExternalRelationTag, error) {
	list := make([]customer.WorkExternalRelationTag, 0)
	sqlStr := "SELECT a.enterprise_id,a.relation_id,a.external_tag_id,b.member_id,b.work_external_user_id FROM micro_scrm_ck.work_external_user_tag a,micro_scrm_ck.work_external_user_relation b WHERE a.enterprise_id = ? AND a.external_tag_id =? AND a.relation_id = b.id AND a.deleted_at = 0 AND b.deleted_at = 0"
	err := gormx.GetCkInstance().WithContext(ctx).Raw(sqlStr, enterpriseId, tagID).Error
	// err := c.DB.Select(ctx, &list, sqlStr, enterpriseId, tagID)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (c *CustomerRepo) DeleteUserTagsByUserID(ctx ctx.Context, enterpriseID int64, externalUserID int64) {
	sqlStr := fmt.Sprintf("delete from work_external_user_tag where relation_id in (select DISTINCT(id) from work_external_user_relation where enterprise_id = ? AND work_external_user_id=?)")
	c.DB.DB.Exec(ctx, sqlStr, enterpriseID, externalUserID)
	sqlStr2 := fmt.Sprintf("delete from work_external_user_relation where enterprise_id = ? AND work_external_user_id=?)")
	c.DB.DB.Exec(ctx, sqlStr2, enterpriseID, externalUserID)
}

// RepeatCustomerCountJob 定时更新重复客户数量
func (c *CustomerRepo) RepeatCustomerCountJob(ctx ctx.Context, enterpriseID int64) error {
	var (
		timeNow = assist.Timestamp()
	)

	sqlStr := `SELECT tmp.work_external_user_id,member_num
FROM (
    SELECT work_external_user_id,count(work_external_user_id) AS member_num
    FROM work_external_user_relation as user_relation
    LEFT JOIN members on members.id = user_relation.member_id
    where  user_relation.enterprise_id = %d
    AND  relation_status = 0
    AND members.deleted_at = 0
    AND user_relation.deleted_at = 0
    GROUP BY work_external_user_id
    HAVING member_num > 1
) as tmp
LEFT JOIN micro_scrm_enterprise_%d.wx_user_corp_external wu on tmp.work_external_user_id = wu.id
WHERE wu.deleted_at = 0`
	sqlStr = fmt.Sprintf(sqlStr, enterpriseID, enterpriseID)
	var rcResult []customer.WorkExternalRepeatCustomerResult
	err := c.DB.Select(ctx, &rcResult, sqlStr)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	// 更新最新的重复数据
	addSql := `INSERT INTO micro_scrm_enterprise_%d.wx_user_corp_external_repeat_customer (work_external_user_id, member_num, created_at, update_at)
VALUES (?,?,%d,%d)  ON DUPLICATE KEY UPDATE member_num= ?,update_at=%d;`
	addSql = fmt.Sprintf(addSql, enterpriseID, timeNow, timeNow, timeNow)
	if len(rcResult) > 0 {
		for _, val := range rcResult {
			_, err = c.DB.DB.Exec(ctx, addSql, val.WorkExternalUserId, val.MemberNum, val.MemberNum)
			if err != nil {
				return err
			}
		}
	}
	// 未更新的数据，就是要被清理的
	delSql := `DELETE FROM micro_scrm_enterprise_%d.wx_user_corp_external_repeat_customer WHERE update_at < %d; `
	delSql = fmt.Sprintf(delSql, enterpriseID, timeNow)
	_, err = c.DB.DB.Exec(ctx, delSql)
	return nil
}

func (c *CustomerRepo) RepeatCustomerCount(ctx ctx.Context, enterpriseID int64, workExternalUserRelationId int64) error {
	log.Info(ctx).Msgf("RepeatCustomerCount %d-%d", enterpriseID, workExternalUserRelationId)
	var (
		timeNow = assist.Timestamp()
	)

	var userExternalId int64
	sqlExternal := fmt.Sprintf("SELECT `work_external_user_id` FROM `work_external_user_relation` WHERE `id` = %d limit 1", workExternalUserRelationId)
	err := c.DB.DB.Get(ctx, &userExternalId, sqlExternal)
	if err != nil {
		log.Error(ctx).Msgf("RepeatCustomerCount-sqlExternal %s,err:%+v", sqlExternal, err)
		return err
	}

	// 先清理指定外部联系人ID的重复数据，后新增
	delSql := `DELETE FROM micro_scrm_enterprise_%d.wx_user_corp_external_repeat_customer WHERE work_external_user_id = %d; `
	delSql = fmt.Sprintf(delSql, enterpriseID, userExternalId)
	_, err = c.DB.DB.Exec(ctx, delSql)
	if err != nil {
		return err
	}
	// 查询指定外部联系ID的重复数量
	numSql := `SELECT work_external_user_id,COUNT(work_external_user_id) AS member_num
    FROM work_external_user_relation AS user_relation
    LEFT JOIN members ON members.id = user_relation.member_id
    WHERE  user_relation.enterprise_id = ?
    AND work_external_user_id = ?
    AND  user_relation.relation_status = 0
    AND members.deleted_at = 0
    AND user_relation.deleted_at = 0 
	GROUP BY work_external_user_id 
	HAVING member_num > 1`
	var rcResult []customer.WorkExternalRepeatCustomerResult
	err = c.DB.Select(ctx, &rcResult, numSql, enterpriseID, userExternalId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	addSql := `INSERT INTO micro_scrm_enterprise_%d.wx_user_corp_external_repeat_customer (work_external_user_id, member_num, created_at, update_at)
VALUES (?,?,%d,%d)  ON DUPLICATE KEY UPDATE member_num= ?,update_at=%d;`
	addSql = fmt.Sprintf(addSql, enterpriseID, timeNow, timeNow, timeNow)
	if len(rcResult) > 0 {
		for _, val := range rcResult {
			_, err = c.DB.DB.Exec(ctx, addSql, val.WorkExternalUserId, val.MemberNum, val.MemberNum)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LossCustomerAddRelationTimeJob 最近一次建立关系的时间：先查relation_callback_records表，再查work_external_user_relation，没有则删除记录无效，不展示.
func (c *CustomerRepo) LossCustomerAddRelationTimeJob(ctx ctx.Context, enterpriseId int64, corpID string) error {
	pageLimit := 5000
	beginID := int64(0)

	// 更新删除事件的关系建立时间
	updateSql := "UPDATE `relation_callback_records` SET `add_relation_time` = ? WHERE `id` = ? "

	// relation_callback_records最近一次建立关系的时间
	nearestCallbackSql := "SELECT `create_time` FROM `relation_callback_records` WHERE `id` < ?  AND `to_user_name` = '%s' AND `change_type` = 'add_external_contact' " +
		"AND `user_id`=? AND `external_user_id`=?  ORDER BY `id` DESC LIMIT 1;"
	nearestCallbackSql = fmt.Sprintf(nearestCallbackSql, corpID)
	// work_external_user_relation最近一次建立关系的时间,如果relation_callback_records没有查到，则查该表
	nearestRelationSql := "SELECT ur.`createtime` FROM `work_external_user_relation` ur LEFT JOIN micro_scrm_enterprise_%d.`wx_user_corp_external` ce  " +
		"ON ur.`work_external_user_id` = ce.`id` LEFT JOIN `members` m ON ur.`member_id` = m.id WHERE  m.`userid` = ? AND ce.`openid` = ?  " +
		"AND ur.`createtime` < ?  ORDER BY ur.`createtime` DESC LIMIT 1"
	nearestRelationSql = fmt.Sprintf(nearestRelationSql, enterpriseId)

	baseSql := "SELECT id,create_time,`user_id`,`external_user_id` FROM `relation_callback_records` WHERE `id`> ? AND `to_user_name` = '%s' " +
		"AND (`change_type` = 'del_external_contact' OR `change_type` = 'del_follow_user' )  AND `add_relation_time` = 0 ORDER BY `id` ASC LIMIT %d;"
	baseSql = fmt.Sprintf(baseSql, corpID, pageLimit)
	for {
		var rcResult []customer.LossCustomerAddRelationTime
		err := c.DB.Select(ctx, &rcResult, baseSql, beginID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if len(rcResult) > 0 {
			for _, row := range rcResult {
				var createTime int64
				err = c.DB.DB.Get(ctx, &createTime, nearestCallbackSql, row.Id, row.UserId, row.ExternalUserId)
				if err != nil {
					log.Error(ctx).Msgf("LossCustomerAddRelationTimeJob-nearestSql err %+v", err)
				}
				if createTime == 0 {
					err = c.DB.DB.Get(ctx, &createTime, nearestRelationSql, row.UserId, row.ExternalUserId, row.CreateTime)
					if err != nil {
						log.Error(ctx).Msgf("LossCustomerAddRelationTimeJob-nearestRelationSql err %+v", err)
					}
				}
				if createTime > 0 {
					log.Info(ctx).Msgf("LossCustomerAddRelationTimeJob-updateSql-ok %d,%+v", createTime, row)
					_, err = c.DB.DB.Exec(ctx, updateSql, createTime, row.Id)
					if err != nil {
						log.Error(ctx).Msgf("LossCustomerAddRelationTimeJob-updateSql err %+V,%d,%+v", err, createTime, row)
					}
				} else {
					log.Info(ctx).Msgf("LossCustomerAddRelationTimeJob-updateSql-没有找到添加时间 %d,%+v", createTime, row)
				}
			}
			beginID = rcResult[len(rcResult)-1].Id
		}

		if len(rcResult) < pageLimit {
			break
		}
	}
	return nil
}

func (c *CustomerRepo) LossCustomerAddRelationTime(ctx ctx.Context, req *work.MemberExternalCallback) work.LossCustomerAddRelationTimeChannelID {
	// 查企业ID
	var resData work.LossCustomerAddRelationTimeChannelID
	var enterpriseId int64
	enterpriseSql := "SELECT `enterprise_id` FROM `enterprise_corp` WHERE `corp_id` = ? LIMIT 1;"
	err := c.DB.DB.Get(ctx, &enterpriseId, enterpriseSql, req.ToUserName)
	if err != nil {
		log.Error(ctx).Msgf("LossCustomerAddRelationTime-enterpriseSql err %+v", err)
		return resData
	}
	nearestCallbackSql := "SELECT `channel_id`,`create_time` as add_relation_time FROM `relation_callback_records` WHERE `id` < ?  AND `to_user_name` = '%s' AND `change_type` = 'add_external_contact' " +
		"AND `user_id`=? AND `external_user_id`=?  ORDER BY `id` DESC LIMIT 1;"
	nearestCallbackSql = fmt.Sprintf(nearestCallbackSql, req.ToUserName)

	err = c.DB.DB.Get(ctx, &resData, nearestCallbackSql, req.ID, req.UserID, req.ExternalUserID)
	if err != nil {
		log.Error(ctx).Msgf("LossCustomerAddRelationTime-nearestSql err %+v", err)
	}
	return resData
}

func (c *CustomerRepo) FindMemberByUserid(ctx ctx.Context, enterpriseId int64, userId string) (customer.LossMember, error) {
	var res customer.LossMember
	sql := "SELECT id,name,avatar FROM `members` WHERE  `enterprise_id` = ? AND `userid` = ? ORDER BY `deleted_at` ASC LIMIT 1;"
	err := c.DB.DB.Get(ctx, &res, sql, enterpriseId, userId)
	return res, err
}
