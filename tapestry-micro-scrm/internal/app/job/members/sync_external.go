package membersJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/xjob"
	"strings"
	"time"
)

type SyncExternal struct {
	name string
}

func NewSyncExternal() *SyncExternal {
	return &SyncExternal{
		name: "同步外部联系人关系",
	}
}

func (job SyncExternal) Run(c ctx.Context, jobHandler *xjob.JobHandler) error {
	log.Info(c).Msgf("%s,开始执行", job.name)
	enterpriseIds, err := appJob.GetEnterpriseIds(jobHandler)
	if err != nil {
		return err
	}

	for _, enterpriseId := range enterpriseIds {
		err = job.handle(c, enterpriseId)
		if err != nil {
			log.Err(c, err).Msgf("%s,处理企业%d异常", job.name, enterpriseId)
		}
	}

	return nil
}

func (job SyncExternal) handle(c ctx.Context, enterpriseId int64) error {
	// 获取当前企业的所有员工
	members, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{"enterprise_id": enterpriseId, "userid": "15249279779"})
	if err != nil {
		return err
	}
	enterpriseInfo, _ := microBossDao.NewEnterprise().GetOne(c, map[string]interface{}{"id": enterpriseId})
	qyWork := work.NewQyWorkImpl().SetDomain(enterpriseInfo.Domain)

	for _, member := range members {
		// 获取当前员工的所有外部联系人关系
		memberUserRelations, err := microScrmDao.NewUserRelation().GetMany(c, "enterprise_id = ? and member_id = ?", enterpriseId, member.ID)
		if err != nil {
			log.Err(c, err).Msgf("%s,获取当前员工的所有外部联系人关系", job.name)
			continue
		}

		// 已存在的外部联系人UnionId
		var exitUnionIds []string
		for _, relation := range memberUserRelations {
			exitUnionIds = append(exitUnionIds, relation.Unionid)
		}

		// 获取当前员工所有外部联系人信息
		memberExternalUsers := job.fetchMemberExternalUser(c, member.Userid, qyWork)

		// 企微外部联系人UnionId
		var externalUnionIds []string
		// 取出所有外部联系人openid
		for _, externalUser := range memberExternalUsers {
			externalUnionIds = append(externalUnionIds, externalUser.Base.UnionID)
		}

		// 取出要添加的id和删除的id
		delExternalUnionIds, _ := helper.ArrayDifference(exitUnionIds, externalUnionIds)

		for _, externalUser := range memberExternalUsers {
			sourceUserId, err := job.mergeUser(c, enterpriseInfo, externalUser)
			if err != nil {
				continue
			}
			err = job.syncExternalRelation(c, enterpriseId, externalUser, member, sourceUserId, delExternalUnionIds)
		}
	}

	return err
}

// 拉取员工的外部联系人
func (job SyncExternal) fetchMemberExternalUser(c ctx.Context, userId string, qyWork *work.QyWorkImpl) []work.WorkExternalContact {
	// 员工外部联系人集合
	var memberExternalUsers []work.WorkExternalContact
	// 页码
	var cursor string

	for {
		externalUsers, err := qyWork.BatchExternalContactGet(c, work.BatchExternalContactRequest{
			MemberUserID: userId,
			Cursor:       cursor,
			Limit:        100,
		})
		if err != nil {
			log.Err(c, err).Msgf("%s,批量获取客户详情异常", job.name)
			continue
		}

		memberExternalUsers = append(memberExternalUsers, externalUsers.ExternalContactList...)
		if externalUsers.NextCursor == "" {
			break
		}
		cursor = externalUsers.NextCursor
	}

	return memberExternalUsers
}

// 合并用户
func (job SyncExternal) mergeUser(c ctx.Context, enterpriseInfo microBossModel.EnterpriseModel, externalUser work.WorkExternalContact) (int64, error) {
	sourceUser := &user.SourceUser{
		EnterpriseID: enterpriseInfo.ID,
		SourceUserBase: customer.SourceUserBase{
			AppID:        enterpriseInfo.CorpID,
			OpenID:       externalUser.Base.ExternalUserID,
			UnionID:      externalUser.Base.UnionID,
			CustomerType: customer.CustomerTypeExternal,
		},
		EnterpriseUserBase: customer.EnterpriseUserBase{
			Name:   externalUser.Base.Name,
			Avatar: externalUser.Base.Avatar,
			Sex:    externalUser.Base.Gender,
		},
		EnterpriseUserExternal: customer.EnterpriseUserExternal{
			Position:        externalUser.Base.Position,
			CorpName:        externalUser.Base.CorpName,
			CorpFullName:    externalUser.Base.CorpFullName,
			Type:            externalUser.Base.Type,
			ExternalProfile: string(externalUser.Base.ExternalProfile),
		},
		CreateAt: time.Now().Unix(),
	}

	_, err := userMysql.NewUser().MergeUser(c, sourceUser)
	if err != nil {
		log.Err(c, err).Msgf("FixExternalContact-%+v外部联系人插入数据库失败: %+v", enterpriseInfo.ID, sourceUser)
		return 0, err
	}

	return sourceUser.ID, nil
}

// 同步外部联系人关系
func (job SyncExternal) syncExternalRelation(
	c ctx.Context,
	enterpriseId int64,
	externalUser work.WorkExternalContact,
	member membersModel.MemberModel,
	externalUserID int64,
	delExternalUnionIds []string,
) error {
	// 更新关系
	if len(delExternalUnionIds) > 0 {
		updateData := map[string]interface{}{
			"relation_status": customer.WorkUserRelationStatusDeleteEachOther,
			// "deleted_at":      time.Now().Unix(),
		}
		err := microScrmDao.NewUserRelation().UpdateData(c, map[string]interface{}{
			"enterprise_id": enterpriseId,
			"member_id":     member.ID,
			"unionid":       delExternalUnionIds,
		}, updateData)

		if err != nil {
			log.Err(c, err).Msgf("FixExternalContact-%+v更新外部联系人关系数据库失败: %+v", enterpriseId, updateData)
		}
		return err
	}

	// 获取当前的员工是否存在关系
	workRelation, err := microScrmDao.NewUserRelation().GetOneUnscoped(c, map[string]interface{}{
		"enterprise_id":         enterpriseId,
		"member_id":             member.ID,
		"work_external_user_id": externalUserID,
	})

	// 更新关系
	if workRelation.ID > 0 {
		if workRelation.RelationStatus != customer.WorkUserRelationStatusNormal {
			workRelation.RelationStatus = customer.WorkUserRelationStatusNormal
		}
		workRelation.UpdatedAt = time.Now().Unix()
		workRelation.DeletedAt = 0
		err = microScrmDao.NewUserRelation().SaveData(c, workRelation)
		if err != nil {
			log.Err(c, err).Msgf("syncExternalRelation-%d更新外部联系人关系数据库失败: %+v", enterpriseId, workRelation)
		}
		return err
	}

	// 添加外部联系人关系
	relation := workExternalModel.WorkExternalUserRelationModel{
		BaseModel: models.BaseModel{
			ID:        uid.ID(c),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		EnterpriseID:       enterpriseId,
		WorkExternalUserID: externalUserID,
		MemberID:           member.ID,
		Remark:             externalUser.FollowInfo.Remark,
		Description:        externalUser.FollowInfo.Description,
		Createtime:         int(externalUser.FollowInfo.CreateTime),
		RemarkCorpName:     externalUser.FollowInfo.RemarkCorpName,
		RemarkMobiles:      strings.Join(externalUser.FollowInfo.RemarkMobiles, ","),
		AddWay:             int32(externalUser.FollowInfo.AddWay),
		OperUserid:         externalUser.FollowInfo.OperUserID,
		RelationStatus:     customer.WorkUserRelationStatusNormal,
		UserID:             member.Userid,
		Unionid:            externalUser.Base.UnionID,
	}

	err = microScrmDao.NewUserRelation().CreateData(c, relation)
	if err != nil {
		log.Err(c, err).Msgf("syncExternalRelation-%d外部联系人关系插入数据库失败: %+v", enterpriseId, relation)
	}

	return err
}
