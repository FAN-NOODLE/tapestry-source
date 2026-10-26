package microScrmModel

import (
	"fmt"
	"gorm.io/gorm"
)

type WxUserCorpExternalRepeatCustomerModel struct {
	WorkExternalUserID int64 `json:"work_external_user_id" gorm:"primaryKey;column:work_external_user_id;type:bigint(20);comment:外部联系人ID"`
	MemberNum          int   `json:"member_num" gorm:"column:member_num;type:int(11);not null;comment:重复员工数量"`
	CreatedAt          int   `json:"created_at" gorm:"column:created_at;type:int(10);not null"`
	UpdateAt           int   `json:"update_at" gorm:"column:update_at;type:int(10);not null"`
}

// TableName returns the table name of the WxUserCorpExternalRepeatCustomer model
func (WxUserCorpExternalRepeatCustomerModel) TableName() string {
	return "wx_user_corp_external_repeat_customer"
}

// BranchTable 切表
func (m WxUserCorpExternalRepeatCustomerModel) BranchTable(enterpriseID int64) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("micro_scrm_enterprise_%d.%s", enterpriseID, m.TableName())
		return tx.Table(tableName)
	}
}
