package mysql

import (
	"database/sql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/guide"
)

type GuideRepo struct {
	*repos.SQLRepo
}

// 添加导购关系
func (u *GuideRepo) AddUserBingGuide(c ctx.Context, data *guide.BasicUserBingGuideInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_guide/insert-customer-relation-guide")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 更新导购关系
func (u *GuideRepo) UpdateUserBingGuide(c ctx.Context, data *guide.BasicUserBingGuideInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_guide/update-customer-relation-guide")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 统计关系是否存在
func (u *GuideRepo) CountUserBingGuide(c ctx.Context, request *guide.BasicUserBingGuideInfoEntity) int64 {
	stmt, err := u.DB.Template("t_guide/count-customer-relation-guide", u.Cond2Args(request))
	if err != nil {
		return 0
	}
	countResponseInfo := &guide.CountResponse{}
	getErr := stmt.NamedGet(c, countResponseInfo, request)
	if getErr != nil {
		return 0
	}
	return countResponseInfo.Count
}

// 删除关系
func (u *GuideRepo) DeleteUserBingGuide(c ctx.Context, data *guide.BasicUserBingGuideInfoEntity) (int64, error) {
	stmt, err := u.DB.Template("t_guide/delete-customer-relation-guide")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// DeleteBatchUserBingGuide 批量删除
func (u *GuideRepo) DeleteBatchUserBingGuide(c ctx.Context, deleted *guide.BatchUserBingGuideInfoEntity) error {
	mapping := make(map[string]interface{})
	mapping["thirdUserIds"] = deleted.ThirdUserIds
	mapping["enterprise_id"] = deleted.EnterpriseID
	mapping["deleted_at"] = deleted.DeletedAt
	stmt, err := u.DB.Template("t_guide/delete-batch-customer-relation-guide", &mapping)
	if err != nil {
		return err
	}
	_, sErr := stmt.NamedExec(c, deleted)
	if sErr != nil {
		return sErr
	}

	return nil
}

func (u *GuideRepo) SelectUserBingGuide(c ctx.Context, batch *guide.BatchUserBingGuideInfoEntity) (*[]guide.BasicUserBingGuideInfoEntity, error) {
	mapping := make(map[string]interface{})
	mapping["thirdUserIds"] = batch.ThirdUserIds
	mapping["enterprise_id"] = batch.EnterpriseID

	stmt, err := u.DB.Template("t_guide/select-customer-relation-guide", &mapping)
	if err != nil {
		return nil, err
	}

	list := &[]guide.BasicUserBingGuideInfoEntity{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// FindGuideRelationByUserID 查询帮定导购关系
func (u *GuideRepo) FindGuideRelationByUserID(cc ctx.Context, enterpriseID int64, userID int64) (*guide.CustomerGuideRelation, error) {
	stmt, err := u.DB.Template("t_guide/find-customer-relation-guide")
	if err != nil {
		return nil, err
	}
	dest := &guide.CustomerGuideRelation{}
	sErr := stmt.Get(cc, dest, enterpriseID, userID)
	if sErr != nil && sErr != sql.ErrNoRows {
		return nil, sErr
	}
	return dest, nil
}
