package workDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type GroupChat struct {
	dao       dao.Dao
	tableName string
}

func NewGroupChat() *GroupChat {
	return &GroupChat{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: workModels.GroupChatModel{}.TableName(),
	}
}

func (d *GroupChat) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]workModels.GroupChatModel, error) {
	var result []workModels.GroupChatModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *GroupChat) GetList(c app.Ctx, page, size int, cond interface{}) ([]workModels.GroupChatModel, int64, error) {
	var result []workModels.GroupChatModel
	var total int64
	offset := (page - 1) * size
	query := d.dao.WithContext(c).Table(d.tableName).Where(cond)
	err := query.Count(&total).Offset(offset).Limit(size).
		Order("created_at desc").
		Find(&result).Error

	return result, total, err
}

func (d *GroupChat) GetListByChatIds(c app.Ctx, enterpriseId int64, chatIds []string) ([]workModels.GroupChatModel, error) {
	var result []workModels.GroupChatModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ? and chat_id In(?)", enterpriseId, chatIds).
		Find(&result).Error
	return result, err
}

func (d *GroupChat) GetGroupChatAllotList(c app.Ctx, enterpriseID int64, groups []string) (data []customerReq.GroupChatInfo, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"chat_id":       groups,
		"type":          1,
	}).Select("id,chat_id,name").Find(&data).Error
	return
}
