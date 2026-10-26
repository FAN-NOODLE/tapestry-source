package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/moment"
)

type MomentRepo struct {
	*repos.SQLRepo
}

// MemberFindOne 获取一条 moment_member
func (d *MomentRepo) MemberFindOne(c ctx.Context, memberRequest moment.MomentMemberGetRequest) (*moment.MomentMember, error) {
	stmt, stmtErr := d.DB.Template("t_moments/member-find-one")
	if stmtErr != nil {
		return nil, stmtErr
	}
	momentMember := &moment.MomentMember{}
	getErr := stmt.Get(c, momentMember, memberRequest.Mark)
	if getErr != nil {
		return nil, getErr
	}
	return momentMember, nil
}

// MemberUpdateOne 更新 moment-member
func (d *MomentRepo) MemberUpdateOne(c ctx.Context, momentMember moment.MomentMember) (int64, error) {
	stmt, stmtErr := d.DB.Template("t_moments/member-update-one")

	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, momentMember)
	if execError != nil {
		return 0, execError
	}

	return result.RowsAffected()
}

// MomentTopicAddOne 新增topic
func (d *MomentRepo) MomentTopicAddOne(c ctx.Context, momentTopic moment.MomentTopic) (int64, error) {
	stmt, err := d.DB.Template("t_moments/topic-add-one")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(c, momentTopic)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	return result.LastInsertId()
}

// MomentMemberAddMany 新增moment-member
func (d *MomentRepo) MomentMemberAddMany(c ctx.Context, momentMembers []*moment.MomentMember) (int64, error) {
	stmt, err := d.DB.Template("t_moments/member-add-one")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(c, momentMembers)
	if sTmtErr != nil {
		return 0, sTmtErr
	}

	return result.LastInsertId()
}

// MomentTopicMemberAddMany 新增 moment_member
func (d *MomentRepo) MomentTopicMemberAddMany(c ctx.Context, momentTopicMembers []*moment.MomentTopicMember) (int64, error) {
	stmt, err := d.DB.Template("t_moments/topic-member-add-one")
	if err != nil {
		return 0, err
	}

	result, sTmtErr := stmt.NamedExec(c, momentTopicMembers)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	return result.LastInsertId()
	//var result int64
	//for _, val := range momentTopicMembers {
	//	_, sErr := stmt.NamedExec(c, val)
	//	if sErr != nil {
	//		result++
	//	}
	//}
	//return result,nil
}

// MomentTopicFindOne 查找一条topic
func (d *MomentRepo) MomentTopicFindOne(c ctx.Context, momentTopicRequest *moment.MomentTopic) (*moment.MomentTopic, error) {
	stmt, stmtErr := d.DB.Template("t_moments/select-topic-one")

	if stmtErr != nil {
		return nil, stmtErr
	}

	momentTopic := &moment.MomentTopic{}
	getErr := stmt.Get(c, momentTopic, momentTopicRequest.ID)

	if getErr != nil {
		return nil, getErr
	}
	return momentTopic, nil
}

// MomentTopicUpdateOne 更新一条topic
func (d *MomentRepo) MomentTopicUpdateOne(c ctx.Context, momentTopic *moment.MomentTopic) (int64, error) {
	stmt, stmtErr := d.DB.Template("t_moments/topic-delete-one")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, momentTopic)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// MomentTopicMemberBatchDelete 批量删除topic-member
func (d *MomentRepo) MomentTopicMemberBatchDelete(c ctx.Context, topicMemberBatchDelete *moment.TopicMemberBatchDelete) (int64, error) {
	stmt, stmtErr := d.DB.Template("t_moments/topic-member-batch-delete")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, topicMemberBatchDelete)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// MemberBatchGetByMemberIds 批量获取朋友圈用户信息
func (d *MomentRepo) MemberBatchGetByMemberIds(c ctx.Context, request moment.MemberBatchGetByMemberIds) (*[]moment.MomentMember, error) {
	mapping := make(map[string]interface{})
	mapping["memberIds"] = request.MemberIds
	mapping["enterprise_id"] = request.EnterpriseId

	stmt, stmtErr := d.DB.Template("t_moments/member-batch-get-by-memberids", &mapping)

	if stmtErr != nil {
		return nil, stmtErr
	}
	momentMembers := &[]moment.MomentMember{}
	res, sErr := stmt.NamedQuery(c, mapping)
	if sErr != nil {
		return nil, sErr
	}

	scanErr := sqlplus.ScanAll(res, momentMembers)

	if scanErr != nil {
		return nil, scanErr
	}
	return momentMembers, nil
}

// TopicMemberDeleteOne 删除moment_topic_members
func (d *MomentRepo) TopicMemberDeleteOne(c ctx.Context, request moment.TopicMemberDeleteRequest) (bool, error) {
	stmt, stmtErr := d.DB.Template("t_moments/topic-member-delete-one")
	if stmtErr != nil {
		return false, stmtErr
	}
	_, execError := stmt.Exec(c, request.DeletedAt, request.MomentTopicMemberId)
	if execError != nil {
		return false, execError
	}
	return true, nil
}

// TopicUpdateMemberCount 更改moment_topic的发布人数
func (d *MomentRepo) TopicUpdateMemberCount(c ctx.Context, request moment.MomentTopic) (bool, error) {

	stmt, stmtErr := d.DB.Template("t_moments/topic-update-count")
	if stmtErr != nil {
		return false, stmtErr
	}

	_, execError := stmt.Exec(c, request.MemberCount, request.ID)

	if execError != nil {
		return false, execError
	}

	return true, nil
}

// BatchUpdateMemberTime 批量修改moment_member的更新时间
func (d *MomentRepo) BatchUpdateMemberTime(c ctx.Context, request moment.BatchUpdateMemberTime) (bool, error) {
	mapping := make(map[string]interface{})
	mapping["memberids"] = request.MemberIDs
	mapping["enterprise_id"] = request.EnterpriseId
	mapping["updated_at"] = request.UpdatedAt

	stmt, err := d.DB.Template("t_moments/member-batch-update-time", &mapping)

	if err != nil {
		return false, err
	}
	_, sErr := stmt.NamedExec(c, mapping)

	if sErr != nil {
		return false, sErr
	}
	return true, nil
}
