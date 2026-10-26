package mysql

import (
	"database/sql"
	"fmt"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"github.com/jmoiron/sqlx"
	"micro-scrm/internal/pkg/uid"

	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/sdk/work"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/members"
)

//MemberRepo 定义成员存储
type MemberRepo struct {
	*repos.SQLRepo
}

// FindAll 查询所有成员
func (d *MemberRepo) FindAll(c ctx.Context, enterpriseID int64) ([]members.Member, error) {
	var list []members.Member
	stmt, err := d.DB.Template("t_members/select-all")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Select(c, &list, enterpriseID, 0)
	if errStmt != nil {
		return nil, err
	}
	return list, nil
}

// MemberDepartmentRelation 成员组织关系入库
func (d *MemberRepo) MemberDepartmentRelation(c ctx.Context, md []*members.MemberDepartment) (int64, error) {
	stmt, err := d.DB.Template("t_member_departments/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, md)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// FindMemberByUserid 根据user_id 查询成员信息
func (d *MemberRepo) FindMemberByUserid(c ctx.Context, m members.Member) (*members.Member, error) {
	stmt, err := d.DB.Template("t_members/select-member-by-userid")
	if err != nil {
		return nil, err
	}
	mem := &members.Member{}
	errStmt := stmt.Get(c, mem, m.EnterpriseID, m.Userid, m.DeletedAt)
	if errStmt != nil {
		return nil, err
	}
	return mem, nil
}

// GetMemberIDByUserid 获取客户详情
func (d *MemberRepo) GetMemberIDByUserid(c ctx.Context, userID string, enterpriseID int64) (*members.Member, error) {
	req := members.Member{
		Userid:       userID,
		EnterpriseID: enterpriseID,
	}

	user, err := d.FindMemberByUserid(c, req)
	switch true {
	case err == sql.ErrNoRows || user == nil:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

// FindMemberByUseridStatus 根据user_id 查询成员信息过滤状态
func (d *MemberRepo) FindMemberByUseridStatus(c ctx.Context, m members.Member) (members.Member, error) {
	var mem members.Member
	stmt, err := d.DB.Template("t_members/select-member-by-userid-status", d.Cond2Args(m))
	if err != nil {
		return mem, err
	}
	errStmt := stmt.NamedGet(c, &mem, m)

	return mem, errStmt
}

// SelectList 查询成员列表
func (d *MemberRepo) SelectList(c ctx.Context, request members.MemberListRequest) (*[]members.Member, error) {
	stmt, err := d.DB.Template("t_members/select-list", d.Cond2Args(request))
	if err != nil {
		return nil, err
	}
	dept := &[]members.Member{}
	res, errStmt := stmt.NamedQuery(c, request)
	defer res.Close()
	serr := sqlplus.ScanAll(res, dept)
	if serr != nil {
		return nil, serr
	}
	if errStmt != nil {
		return nil, err
	}
	return dept, nil
}

//DeleteAllMemberDepartmentRelation 删除成员组织关系
func (d *MemberRepo) DeleteAllMemberDepartmentRelation(c ctx.Context, request interface{}) error {
	stmt, err := d.DB.Template("t_member_departments/delete_relation")
	if err != nil {
		return err
	}
	_, errStmt := stmt.NamedExec(c, request)
	if errStmt != nil {
		return errStmt
	}
	return nil
}

// MemberDepartmentBatchInsert 批量插入组织成员关系数据
func (d *MemberRepo) MemberDepartmentBatchInsert(ctx ctx.Context, md []*members.MemberDepartment) error {
	stmt, err := d.DB.Template("t_member_departments/insert-batch")
	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(ctx, md)
	if errStmt != nil {
		return errStmt
	}
	return nil
}

// BatchDeleteMembersByIds 批量删除成员
func (d *MemberRepo) BatchDeleteMembersByIds(ctx ctx.Context, md *members.BatchDeleteMember) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = md.IDs
	stmt, err := d.DB.Template("t_members/delete-batch", &mapping)
	if err != nil {
		return err
	}
	_, sErr := stmt.NamedExec(ctx, md)

	if sErr != nil {
		return sErr
	}

	return nil
}

func (d *MemberRepo) GetMembersByDepartIds(c ctx.Context, bm *members.BatchGetMember) (*[]members.MemberDepartment, error) {
	mapping := make(map[string]interface{})
	mapping["depart_ids"] = bm.DepartIDs
	mapping["enterprise_id"] = bm.EnterpriseID
	stmt, stmtErr := d.DB.Template("t_member_departments/batch-get-departments-by-thirds", &mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	memberDeparts := &[]members.MemberDepartment{}
	res, sErr := stmt.NamedQuery(c, mapping)
	if sErr != nil {
		return nil, sErr
	}
	scanErr := sqlplus.ScanAll(res, memberDeparts)
	if scanErr != nil {
		return nil, scanErr
	}

	return memberDeparts, nil
}

// GetMembersInfoByDepartmentIds 多部门获取员工列表
func (d *MemberRepo) GetMembersInfoByDepartmentIds(c ctx.Context, enterpriseId int64, departmentIds []int64, baUserIds []string) ([]*members.SimpleMember, error) {
	var (
		mems    []*members.SimpleMember
		baWhere = " and 1=1"
	)
	if len(baUserIds) > 0 {
		baWhere = fmt.Sprintf(" and m.userid in ('%s') ", strings.Join(baUserIds, "','"))
	}
	preQuery := fmt.Sprintf(`select m.id,m.mobile,m.position,m.userid,m.name,m.status,md.department_id from members m inner join member_departments md on m.id = md.member_id where m.enterprise_id = ? 
		and md.department_id in (?) and m.deleted_at=0 and m.status in (0,1) and md.deleted_at=0 %s group by m.userid`, baWhere)
	query, args, err := sqlx.In(preQuery, enterpriseId, departmentIds)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Select(c, &mems, query, args...); err != nil {
		return nil, err
	}
	return mems, nil
}

// GetMemberByMemberIds 通过员工ID获取员工
func (d *MemberRepo) GetMemberByMemberIds(c ctx.Context, batchMemberIDs *members.BatchGetMemberByMemberId) (*[]members.Member, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = batchMemberIDs.IDs
	mapping["enterprise_id"] = batchMemberIDs.EnterpriseID
	stmt, stmtErr := d.DB.Template("t_members/select-member-by-ids", &mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	memberinfos := &[]members.Member{}
	res, sErr := stmt.NamedQuery(c, mapping)
	if sErr != nil {
		return nil, sErr
	}
	scanErr := sqlplus.ScanAll(res, memberinfos)
	if scanErr != nil {
		return nil, scanErr
	}

	return memberinfos, nil
}

// GetMemberByDpId 通过部门ID获取员工
func (d *MemberRepo) GetMemberByDpId(c ctx.Context, deptId int64, enterpriseId int64) (*[]members.MemberDepartment, error) {
	stmt, stmtErr := d.DB.Template("t_members/select-member-by-dpid")
	if stmtErr != nil {
		return nil, stmtErr
	}
	momentMember := &[]members.MemberDepartment{}
	res, getErr := stmt.Query(c, deptId, enterpriseId)
	if getErr != nil {
		return nil, getErr
	}
	scanErr := sqlplus.ScanAll(res, momentMember)
	if scanErr != nil {
		return nil, scanErr
	}
	return momentMember, nil
}

// GetNoDepartmentMember 获取无部门的员工
func (d *MemberRepo) GetNoDepartmentMember(c ctx.Context, momentListRequest *members.MomentListRequest) (*[]members.MomentList, error) {
	mapping := make(map[string]interface{})
	mapping["mem_username"] = momentListRequest.Search.MemUsername
	stmt, stmtErr := d.DB.Template("t_members/select-no-department-member")
	if stmtErr != nil {
		return nil, stmtErr
	}
	momentList := &[]members.MomentList{}
	res, getErr := stmt.Query(c, momentListRequest.EnterpriseId, momentListRequest.EnterpriseId)
	if getErr != nil {
		return nil, getErr
	}
	scanErr := sqlplus.ScanAll(res, momentList)
	if scanErr != nil {
		return nil, scanErr
	}
	return momentList, nil
}

// AddExternalMigration 添加外部联系人迁移记录
func (d *MemberRepo) AddExternalMigration(c ctx.Context, data []*members.WorkExternalUserMigrationEntity) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	stmt, err := d.DB.Template("t_work_external_user_migration/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetExternalMigration 获取外部联系人迁移数据
func (d *MemberRepo) GetExternalMigration(c ctx.Context, page int, size int) ([]*members.WorkExternalUserMigrationEntity, error) {
	offset := (page - 1) * size
	limit := size
	option := &repos.FindPageOption{
		TplPrefix: "t_work_external_user_migration/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*members.WorkExternalUserMigrationEntity{}
		},
	}
	_, list, err := d.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, sqlplus.Cond{}, option)
	if err != nil {
		return nil, err
	}
	data := list.(*[]*members.WorkExternalUserMigrationEntity)
	return *data, nil
}

// UpdateExternalMigration 更新迁移记录状态
func (d *MemberRepo) UpdateExternalMigration(c ctx.Context, data *members.WorkExternalUserMigrationEntity) (int64, error) {
	stmt, err := d.DB.Template("t_work_external_user_migration/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// BatchUpdateExternalMigration 批量更新迁移数据状态
func (d *MemberRepo) BatchUpdateExternalMigration(c ctx.Context, idList []int64, status int) (int64, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = idList
	mapping["status"] = status
	mapping["updated_at"] = time.Now().Unix()
	stmt, err := d.DB.Template("t_work_external_user_migration/batch-update", &mapping)
	if err != nil {
		return 0, err
	}
	res, sErr := stmt.NamedExec(c, mapping)

	if sErr != nil {
		return 0, sErr
	}
	return res.RowsAffected()
}

// GetDepartmentsById 获取企业成员所属部门
func (d *MemberRepo) GetDepartmentsById(c ctx.Context, memberId int64, enterpriseID int64) ([]*members.MemberDepartment, error) {
	stmt, err := d.DB.Template("t_member_departments/select-departments-by-user")
	if err != nil {
		return nil, err
	}
	var data []*members.MemberDepartment
	if err := stmt.NamedSelect(c, &data, members.MemberDepartment{
		MemberID:     memberId,
		EnterpriseID: enterpriseID,
	}); err != nil {
		return nil, err
	}
	return data, nil
}

// GetMembersDepartments 批量获取员工所属部门
func (d *MemberRepo) GetMembersDepartments(c ctx.Context, enterpriseID int64, memberIds []int64) (map[int64][]int64, error) {
	stmt, err := d.DB.Template("t_member_departments/select-departments-by-memberids")
	if err != nil {
		return nil, err
	}
	var (
		dps      []*members.MemberDepartment
		memberDp = make(map[int64][]int64)
	)
	query, args, err := sqlx.In(stmt.SQLRaw(), enterpriseID, memberIds)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Select(c, &dps, query, args...); err != nil {
		return nil, err
	}
	for _, v := range dps {
		memberDp[v.MemberID] = append(memberDp[v.MemberID], v.DepartmentID)
	}
	return memberDp, nil
}

// GetMembersByUserIds 根据userid批量获取员工信息
func (d *MemberRepo) GetMembersByUserIds(c ctx.Context, enterpriseID int64, userIds []string) ([]*members.Member, error) {
	stmt, err := d.DB.Template("t_members/select-members-info-by-userids")
	if err != nil {
		return nil, err
	}
	query, args, err := sqlx.In(stmt.SQLRaw(), enterpriseID, userIds)
	if err != nil {
		return nil, err
	}
	var data []*members.Member
	if err := d.DB.Select(c, &data, query, args...); err != nil {
		return nil, err
	}
	return data, nil
}

// GetMembersByMobiles 根据手机号批量获取员工信息
func (d *MemberRepo) GetMembersByMobiles(c ctx.Context, enterpriseID int64, mobiles []string) ([]*members.Member, error) {
	query, args, err := sqlx.In(`select id, userid, name, mobile, position, gender, email, avatar, thumb_avatar, telephone, alias, status, qr_code, address, open_userid, enterprise_id, created_at
		from members where deleted_at = 0 and enterprise_id = ? and mobile in (?) and mobile != ''`, enterpriseID, mobiles)
	if err != nil {
		return nil, err
	}
	var data []*members.Member
	if err := d.DB.Select(c, &data, query, args...); err != nil {
		return nil, err
	}
	return data, nil
}

// GetMembersByMarsIds 根据mars_ids批量获取员工信息
func (d *MemberRepo) GetMembersByMarsIds(c ctx.Context, enterpriseID int64, marsIds []string) ([]*members.Member, error) {
	var req = make(map[string]interface{})
	req["EnterpriseID"] = enterpriseID
	stmt, err := d.DB.Template("t_members/select-members-info-by-marsids", &req)
	if err != nil {
		return nil, err
	}
	query, args, err := sqlx.In(stmt.SQLRaw(), marsIds)
	if err != nil {
		log.Debug(c).Msgf("%+v", err)
		return nil, err
	}
	var data []*members.Member
	if err := d.DB.Select(c, &data, query, args...); err != nil {
		log.Debug(c).Msgf("%+v", err)
		return nil, err
	}
	return data, nil
}

func (d *MemberRepo) BatchRefreshDepartments(c ctx.Context, enterpriseID int64, departs map[int64]*department.Department) error {
	nowtm := time.Now().Unix()
	query := db.MakeUpsertSQL("departments", []db.UpsertColumn{
		{Name: "id", Insert: true},
		{Name: "third_id", Insert: true},
		{Name: "enterprise_id", Insert: true},
		{Name: "third_parent_id", Insert: true, Update: true},
		{Name: "name", Insert: true, Update: true},
		{Name: "name_en", Insert: true, Update: true},
		{Name: "order", Insert: true, Update: true},
		{Name: "created_at", Insert: true},
		{Name: "updated_at", Insert: true, Update: true},
		{Name: "deleted_at", Update: true},
	})
	// 先更新
	for _, depart := range departs {
		depart.UpdatedAt = nowtm
		if err := helper.WithRetry(3, time.Second, func() error {
			_, err := d.DB.DB.Exec(c, query, uid.ID(c), depart.ID, enterpriseID, depart.ParentID, depart.Name, depart.NameEn, depart.Order, depart.CreatedAt, depart.UpdatedAt,
				depart.ParentID, depart.Name, depart.NameEn, depart.Order, depart.UpdatedAt, 0)
			return err
		}); err != nil {
			return err
		}
	}
	query = "update departments set deleted_at=? where enterprise_id=? and updated_at<?"
	_, err := d.DB.DB.Exec(c, query, nowtm, enterpriseID, nowtm)
	return err
}

type tmpMember struct {
	ID     int64  `db:"id"`
	Userid string `db:"userid"`
}

// BatchInsertMemberDepartment 批量插入组织成员关系数据
func (d *MemberRepo) BatchInsertMemberDepartment(ctx ctx.Context, enterpriseID int64, md []work.MemberDetail, bRefresh bool, dbDeparts map[int64]department.Department) error {
	nowtm := time.Now().Unix()
	query := db.MakeUpsertSQL("members", []db.UpsertColumn{
		{Name: "id", Insert: true},
		{Name: "userid", Insert: true},
		{Name: "name", Insert: true, Update: true},
		{Name: "mobile", Insert: true, Update: true},
		{Name: "position", Insert: true, Update: true},
		{Name: "gender", Insert: true, Update: true},
		{Name: "email", Insert: true, Update: true},
		{Name: "avatar", Insert: true, Update: true},
		{Name: "thumb_avatar", Insert: true, Update: true},
		{Name: "telephone", Insert: true, Update: true},
		{Name: "alias", Insert: true, Update: true},
		{Name: "status", Insert: true, Update: true},
		{Name: "qr_code", Insert: true, Update: true},
		{Name: "address", Insert: true, Update: true},
		{Name: "open_userid", Insert: true, Update: true},
		{Name: "enterprise_id", Insert: true},
		{Name: "created_at", Insert: true},
		{Name: "updated_at", Insert: true, Update: true},
		{Name: "deleted_at", Update: true},
	})
	for _, detail := range md {
		if err := helper.WithRetry(3, time.Second, func() error {
			tmpM := tmpMember{}
			if err := d.DB.DB.Get(ctx, &tmpM, "select id,userid from members where userid=? AND enterprise_id=? ", detail.Userid, enterpriseID); err != nil {
				// 查询手机号是否存在
				if detail.Mobile != "" {
					d.DB.DB.Get(ctx, &tmpM, "select id,userid from members where enterprise_id=? AND mobile=?", enterpriseID, detail.Mobile)
					if tmpM.ID > 0 && detail.Userid != tmpM.Userid {
						d.DB.DB.Exec(ctx, "update members set userid=?,updated_at=? where enterprise_id=? AND id=?", detail.Userid, nowtm, enterpriseID, tmpM.ID)
					}
				}
			}
			mid := tmpM.ID
			if mid == 0 {
				mid = uid.ID(ctx)
			}

			_, err := d.DB.DB.Exec(ctx, query, mid, detail.Userid, detail.Name, detail.Mobile, detail.Position, detail.Gender, detail.Email, detail.Avatar, detail.ThumbAvatar, detail.Telephone, detail.Alias, detail.Status, detail.QrCode, detail.Address, detail.OpenUserid, enterpriseID, nowtm, nowtm,
				detail.Name, detail.Mobile, detail.Position, detail.Gender, detail.Email, detail.Avatar, detail.ThumbAvatar, detail.Telephone, detail.Alias, detail.Status, detail.QrCode, detail.Address, detail.OpenUserid, nowtm, 0)
			if err != nil {
				return err
			}
			mdquery := "insert into member_departments(id, member_id, enterprise_id, department_id, created_at, updated_at) values (?, ?, ?, ?, ?, ?) on duplicate key update updated_at=?, deleted_at=0"
			for _, dp := range detail.Department {
				var dpID int64
				// 在数据库中不存在的部门不更新
				if _, ok := dbDeparts[dp]; !ok {
					continue
				}
				dpID = dbDeparts[dp].ID
				if _, err := d.DB.DB.Exec(ctx, mdquery, uid.ID(ctx), mid, enterpriseID, dpID, nowtm, nowtm, nowtm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if bRefresh {
		query = "update members set deleted_at=? where enterprise_id=? and updated_at<?"
		if _, err := d.DB.DB.Exec(ctx, query, nowtm, enterpriseID, nowtm); err != nil {
			return err
		}
		query = "update member_departments set deleted_at=?, updated_at=? where enterprise_id=? and updated_at<?"
		if _, err := d.DB.DB.Exec(ctx, query, nowtm, nowtm, enterpriseID, nowtm); err != nil {
			return err
		}
	}

	return nil
}

// CommitSyncMembers 提交同步通讯录申请
func (d *MemberRepo) CommitSyncMembers(c ctx.Context, enterpriseID int64, adminID int64) error {
	query := "insert into work_member_sync(enterprise_id, commit_at, admin_id) values(?, ?, ?) on duplicate key update commit_at=?, admin_id=?"
	nowtm := time.Now().Unix()
	_, err := d.DB.DB.Exec(c, query, enterpriseID, nowtm, adminID, nowtm, adminID)
	return err
}

type syncHistory struct {
	SuccAt     int64  `db:"succ_at"`
	CommitAt   int64  `db:"commit_at"`
	SyncAt     int64  `db:"sync_at"`
	DepartHash string `db:"succ_depart_hash"`
	MemberHash string `db:"succ_member_hash"`
}

func (d *MemberRepo) GetSyncContact(c ctx.Context, enterpriseID int64) bool {
	var history syncHistory
	d.DB.DB.Get(c, &history, "select succ_at, commit_at, sync_at, succ_depart_hash, succ_member_hash from work_member_sync where enterprise_id=?", enterpriseID)
	if history.SuccAt < history.CommitAt {
		return false
	} else {
		return true
	}
}

func (d *MemberRepo) GetUserIdsByDepId(c ctx.Context, deptId int64, enterpriseID int64) ([]string, error) {
	var userIdList []string
	var sql = "SELECT m.`userid` FROM `members` m LEFT JOIN `member_departments` d ON m.`id` = d.`member_id` " +
		"WHERE d.`enterprise_id` = ? AND  d.`department_id`= ?;"
	err := d.DB.DB.Select(c, &userIdList, sql, enterpriseID, deptId)
	return userIdList, err
}

// GetExternalUserRelatedMembers 获取外部联系人关联的员工数据
func (d *MemberRepo) GetExternalUserRelatedMembers(c ctx.Context, enterpriseId int64, workExternalUserId []int64) (map[int64][]int64, error) {
	var relatedMembers []struct {
		MemberID           int64 `db:"member_id"`
		WorkExternalUserID int64 `db:"work_external_user_id"`
	}
	query, args, err := sqlx.In(`select member_id,work_external_user_id from work_external_user_relation where enterprise_id = ? and work_external_user_id in (?) and relation_status in (0,4)`,
		enterpriseId, workExternalUserId)
	if err != nil {
		return nil, err
	}
	if err := d.DB.Select(c, &relatedMembers, query, args...); err != nil {
		return nil, err
	}

	var data = make(map[int64][]int64)
	for _, v := range relatedMembers {
		data[v.WorkExternalUserID] = append(data[v.WorkExternalUserID], v.MemberID)
	}
	return data, nil
}
