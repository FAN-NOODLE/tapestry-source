package appJob

import (
	"fmt"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	bossModel "micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

const conversationArchiveTagRepairWindow = 24 * time.Hour

type ConversationArchiveTagRepair struct {
	name string
}

type conversationArchiveTagRepairMember struct {
	UserID         string `gorm:"column:userid"`
	Name           string `gorm:"column:name"`
	Position       string `gorm:"column:position"`
	DepartmentName string `gorm:"column:department_name"`
	DepartmentID   int64  `gorm:"column:department_third_id"`
	CreatedAt      int64  `gorm:"column:created_at"`
}

func NewConversationArchiveTagRepair() *ConversationArchiveTagRepair {
	return &ConversationArchiveTagRepair{name: "会话存档标签补打"}
}

func (job *ConversationArchiveTagRepair) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIDs, err := GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	beginAt := time.Now().Add(-conversationArchiveTagRepairWindow).Unix()
	for _, enterpriseID := range enterpriseIDs {
		if err = job.repair(c, enterpriseID, beginAt); err != nil {
			return err
		}
	}
	return nil
}

func (job *ConversationArchiveTagRepair) repair(c ctx.Context, enterpriseID, beginAt int64) error {
	tagID := config.CfgData.CoachSaTag.TagId
	if tagID == 0 {
		log.Info(c).Msgf("%s跳过, 未配置coach_sa_tag.tag_id, enterprise_id:%d", job.name, enterpriseID)
		return nil
	}

	enterpriseInfo, err := orm.BossService.Enterprise.GetInfo(bossModel.EnterpriseModel{ID: enterpriseID})
	if err != nil {
		return fmt.Errorf("查询企业信息失败, enterprise_id:%d: %w", enterpriseID, err)
	}

	storeRootDeptID, err := job.getStoreRootDeptID(c, enterpriseID)
	if err != nil {
		return err
	}

	shouldTagMembers, err := job.getShouldTagMembers(c, enterpriseID, beginAt, storeRootDeptID)
	if err != nil {
		return err
	}

	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	tagUsers, err := job.getTagUsers(c, enterpriseInfo.Domain, workImpl, tagID)
	if err != nil {
		return fmt.Errorf("查询企微会话存档标签成员失败, enterprise_id:%d: %w", enterpriseID, err)
	}

	taggedUserIDs := make(map[string]struct{}, len(tagUsers.Userlist))
	for _, user := range tagUsers.Userlist {
		taggedUserIDs[user.Userid] = struct{}{}
	}

	seenShouldTagUserIDs := make(map[string]struct{}, len(shouldTagMembers))
	missingMembers := make([]conversationArchiveTagRepairMember, 0)
	missingUserIDs := make([]string, 0)
	for _, member := range shouldTagMembers {
		if _, exists := seenShouldTagUserIDs[member.UserID]; exists {
			continue
		}
		seenShouldTagUserIDs[member.UserID] = struct{}{}
		if _, exists := taggedUserIDs[member.UserID]; exists {
			continue
		}
		missingMembers = append(missingMembers, member)
		missingUserIDs = append(missingUserIDs, member.UserID)
	}

	log.Info(c).Msgf(
		"%s审计, enterprise_id:%d, enterprise:%s, tag_id:%d, begin_at:%s, 应打:%d, 直接已打:%d, 待补打:%d",
		job.name,
		enterpriseID,
		enterpriseInfo.Name,
		tagID,
		time.Unix(beginAt, 0).Format("2006-01-02 15:04:05"),
		len(seenShouldTagUserIDs),
		len(taggedUserIDs),
		len(missingMembers),
	)
	for _, member := range missingMembers {
		log.Info(c).Msgf(
			"%s待补打成员, enterprise_id:%d, userid:%s, name:%s, position:%s, created_at:%s, department_third_id:%d, department:%s",
			job.name,
			enterpriseID,
			member.UserID,
			member.Name,
			member.Position,
			time.Unix(member.CreatedAt, 0).Format("2006-01-02 15:04:05"),
			member.DepartmentID,
			member.DepartmentName,
		)
	}

	if len(missingUserIDs) == 0 {
		return nil
	}

	if err = workImpl.AddTagUsers(c, work.AddTagUsersRequest{
		QyGwRequest: work.QyGwRequest{ReqPL: work.QWCommunalContact},
		TagID:       tagID,
		UserList:    missingUserIDs,
	}); err != nil {
		return fmt.Errorf("批量补打会话存档标签失败, enterprise_id:%d, tag_id:%d, user_count:%d: %w", enterpriseID, tagID, len(missingUserIDs), err)
	}

	log.Info(c).Msgf(
		"%s成功, enterprise_id:%d, tag_id:%d, user_count:%d, userids:%v",
		job.name,
		enterpriseID,
		tagID,
		len(missingUserIDs),
		missingUserIDs,
	)
	return nil
}

func (job *ConversationArchiveTagRepair) getStoreRootDeptID(c ctx.Context, enterpriseID int64) (int, error) {
	configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"key":           "store_root_dept_id",
	})
	if err != nil {
		return 0, fmt.Errorf("获取门店根部门配置失败, enterprise_id:%d: %w", enterpriseID, err)
	}
	if len(configData.Data) == 0 {
		return 0, fmt.Errorf("未配置门店根部门, enterprise_id:%d", enterpriseID)
	}
	return int(configData.Data[0].Value), nil
}

func (job *ConversationArchiveTagRepair) getShouldTagMembers(c ctx.Context, enterpriseID, beginAt int64, storeRootDeptID int) ([]conversationArchiveTagRepairMember, error) {
	var members []conversationArchiveTagRepairMember
	err := gormx.GetInstance().WithContext(c).
		Table("members AS m").
		Select("m.userid, m.name, m.position, m.created_at, d.name AS department_name, d.third_id AS department_third_id").
		Joins("INNER JOIN member_departments AS md ON md.member_id = m.id AND md.enterprise_id = m.enterprise_id AND md.deleted_at = 0").
		Joins("INNER JOIN departments AS d ON d.id = md.department_id AND d.enterprise_id = m.enterprise_id AND d.deleted_at = 0").
		Where("m.enterprise_id = ?", enterpriseID).
		Where("m.deleted_at = 0").
		Where("m.created_at >= ?", beginAt).
		Where("m.position IN ?", membersDao.ConversationArchiveTag).
		Where("d.third_parent_id = ?", storeRootDeptID).
		Order("m.created_at ASC, m.userid ASC").
		Scan(&members).Error
	if err != nil {
		return nil, fmt.Errorf("查询应打会话存档标签成员失败, enterprise_id:%d: %w", enterpriseID, err)
	}
	return members, nil
}

func (job *ConversationArchiveTagRepair) getTagUsers(c ctx.Context, domain string, workImpl *work.QyWorkImpl, tagID int) (work.TagUsersResponse, error) {
	var response work.TagUsersResponse
	token, err := workImpl.GetToken(c)
	if err != nil {
		return response, err
	}

	req := struct {
		work.QyGwRequest
		TagID int `json:"tagid"`
	}{
		QyGwRequest: work.QyGwRequest{
			ReqPL:     work.QWCommunalContact,
			ReqMethod: work.GET,
			ReqURL:    work.TagUsers,
			Token:     token,
		},
		TagID: tagID,
	}
	url := fmt.Sprintf(work.QyGw, config.CfgData.Service.Work.Domain, domain)
	err = work.WorkRequest(c, url, req, &response)
	return response, err
}
