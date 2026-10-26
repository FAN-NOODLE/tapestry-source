package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app/dao"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type MemberLearnRank struct {
	dao       dao.Dao
	tableName string
}

func NewMemberLearnRank() *MemberLearnRank {
	return &MemberLearnRank{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.MemberLearnRankModel{}.TableName(),
	}
}

func (d *MemberLearnRank) BatchCreate(c ctx.Context, enterpriseId int64, list []microScrmModel.MemberLearnRankModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&list).Error
}

func (d *MemberLearnRank) GetTopList(c ctx.Context, enterpriseId int64, fileType string, pMonth int, inviteTop int) (list []microScrmModel.MemberLearnRankModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("file_type = ? AND p_month = ? AND member_id > 0 AND deleted_at = 0", fileType, pMonth).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if fileType == excitationModels.MedalNewCoursesList {
				return db.Where("value = ?", "100.00%")
			}
			return db
		}).Order("rank asc").
		Limit(inviteTop).
		Find(&list).Error
	return
}

func (d *MemberLearnRank) ConsumeData(c ctx.Context, enterpriseId int64, fileType string, pMonth int) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("file_type = ? AND p_month = ?", fileType, pMonth).
		UpdateColumns(map[string]interface{}{
			"deleted_at": time.Now().Unix(),
		}).Error
}
