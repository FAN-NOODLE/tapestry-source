package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/source"
)

type LinkKOSDao struct {
	dao.Repo[source.LinkKOSModel]
}

func NewLinkKOSDao() *LinkKOSDao {
	return &LinkKOSDao{Repo: dao.Repo[source.LinkKOSModel]{DB: app.DB()}}
}

func (d *LinkKOSDao) WithTx(tx *gorm.DB) *LinkKOSDao {
	d.SetDb(tx)
	return d
}

func (d *LinkKOSDao) GetByPostIDUnscoped(c app.Ctx, enterpriseID int64, postID string) (source.LinkKOSModel, error) {
	var result source.LinkKOSModel
	err := d.DB.WithContext(c).Model(source.LinkKOSModel{}).Unscoped().
		Where("enterprise_id = ? AND post_id = ?", enterpriseID, postID).
		Find(&result).Error
	return result, err
}

func (d *LinkKOSDao) SaveData(c app.Ctx, data source.LinkKOSModel) error {
	return d.DB.WithContext(c).Table(data.TableName()).Unscoped().Save(&data).Error
}

func (d *LinkKOSDao) GetSourceURLMap(c app.Ctx, enterpriseID int64, sourceIDs []int64) (map[int64]string, error) {
	result := make(map[int64]string)
	if len(sourceIDs) == 0 {
		return result, nil
	}

	var list []source.LinkKOSModel
	err := d.DB.WithContext(c).Model(source.LinkKOSModel{}).
		Select("source_id, source_url").
		Where("enterprise_id = ? AND source_id IN ? AND deleted_at = 0", enterpriseID, sourceIDs).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		result[item.SourceID] = item.SourceURL
	}
	return result, nil
}
