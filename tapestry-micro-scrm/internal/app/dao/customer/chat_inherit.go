package customerDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerModels "micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
)

type ChatInherit struct {
	dao       dao.Dao
	tableName string
}

func NewChatInherit() *ChatInherit {
	return &ChatInherit{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: customerModels.CustomerChatInheritModel{}.TableName(),
	}
}

func (d *ChatInherit) CreateInBatches(c app.Ctx, data []customerModels.CustomerChatInheritModel) error {
	return d.dao.WithContext(c).Table(d.tableName).CreateInBatches(data, len(data)).Error
}

func (d *ChatInherit) GetChatCountForUserid(c app.Ctx, enterpriseID int64, handoverUserid string, status int) (total int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id=? AND handover_userid = ? AND status = ?", enterpriseID, handoverUserid, status).Count(&total).Error
	return total, err
}
func (d *ChatInherit) GetTakeOverGroupCount(c app.Ctx, request customerReq.GroupChatList) (total int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("takeover_department_name,takeover_name,takeover_userid,chat_name as name,created_at as takeover_time,status,chat_id").
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"handover_userid": request.HandoverUserId,
		}).Count(&total).Error
	return
}

func (d *ChatInherit) GetTakeOverGroupList(c app.Ctx, request customerReq.GroupChatList) (data []customerReq.TakeOverGroupInfo, err error) {
	offset := (request.Page - 1) * request.Size
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("chat_id,chat_name as name,takeover_userid,takeover_name,takeover_department_name,created_at as takeover_time,status").
		Where(map[string]interface{}{
			"enterprise_id":   request.EnterpriseID,
			"handover_userid": request.HandoverUserId,
		}).Offset(offset).Limit(request.Size).Find(&data).Error
	return
}

func (d *ChatInherit) GetChatCountNumForUserid(c app.Ctx, enterpriseID int64, handoverUserids []string, status int) (res []customerReq.InsertChatNumList, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("count(*) as num,handover_userid").
		Where("enterprise_id = ? AND  status = ?", enterpriseID, status).
		Where("handover_userid in  ?", handoverUserids).Find(&res).Error
	return res, err
}
