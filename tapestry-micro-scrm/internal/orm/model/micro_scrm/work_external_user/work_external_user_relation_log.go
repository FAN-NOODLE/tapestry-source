package work_external_user

import (
	"fmt"
	"micro-scrm/internal/orm/model"
)

type WorkExternalUserRelationLogModel struct {
	model.MicroScrmBaseModel

	WorkExternalUserID int64  `gorm:"column:work_external_user_id;type:bigint(20);not null;index:index_0,work_external_user_id_index;default:0;comment:企业外部联系人信息表id"`
	MemberID           int64  `gorm:"column:member_id;type:bigint(11);not null;index:index_0,member_id_index;default:0;comment:成员表id"`
	State              string `gorm:"column:state;type:varchar(255);not null;comment:企业自定义的state参数"`
	RelationStatus     int    `gorm:"column:relation_status;type:int(11);not null;index:relation_status_index;default:0;comment:外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中"`
}

// TableName returns the table name of the struct_name model
func (w WorkExternalUserRelationLogModel) TableName() string {
	return "work_external_user_relation_log"
}

func (w WorkExternalUserRelationLogModel) DBPrefix() string {
	return "micro_scrm"
}

func (w WorkExternalUserRelationLogModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", w.DBPrefix(), enterpriseID, w.TableName())
}
