package work_external_user

import (
	"fmt"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
)

type WorkExternalUserRelationModel struct {
	model.MicroScrmBaseModel

	WorkExternalUserID  int64  `gorm:"column:work_external_user_id;type:bigint(20);not null;index:index_0,work_external_user_id_index;default:0;comment:企业外部联系人信息表id"`
	MemberID            int64  `gorm:"column:member_id;type:bigint(11);not null;index:index_0,member_id_index;default:0;comment:成员表id"`
	Remark              string `gorm:"column:remark;type:varchar(255);not null;comment:成员设置的外部联系人备注"`
	Description         string `gorm:"column:description;type:varchar(255);not null;comment:成员设置的外部联系人描述"`
	Createtime          int    `gorm:"column:createtime;type:int(11);not null;index:createtime_index;default:0;comment:关系建立时间"`
	RemarkCorpName      string `gorm:"column:remark_corp_name;type:varchar(255);not null;comment:成员备注的外部联系人企业名称"`
	RemarkMobiles       string `gorm:"column:remark_mobiles;type:varchar(255);not null;comment:成员备注的外部联系人手机号"`
	AddWay              uint   `gorm:"column:add_way;type:int(10) unsigned;not null;default:0;comment:该成员添加此外部联系人的来源"`
	State               string `gorm:"column:state;type:varchar(255);not null;comment:企业自定义的state参数"`
	OperUserid          string `gorm:"column:oper_userid;type:varchar(255);not null;default:0;comment:发起添加的userid，如果成员主动添加，为成员的userid；如果是客户主动添加，则为客户的外部联系人userid；如果是内部成员共享/管理员分配，则为对应的成员/管理员userid"`
	RelationStatus      int    `gorm:"column:relation_status;type:int(11);not null;index:relation_status_index;default:0;comment:外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中"`
	MemberDeletedTime   int    `gorm:"column:member_deleted_time;type:int(11);not null;index:member_deleted_time_index;default:0;comment:成员主动删除关系时间"`
	ExternalDeletedTime int    `gorm:"column:external_deleted_time;type:int(11);not null;index:external_deleted_time_index;default:0;comment:外部联系人删除关系时间"`

	UserID  string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:成员表user_id"`
	Unionid string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null;comment:外部联系人的unionid"`

	MemberInfo micro_scrm.MemberModel `gorm:"foreignKey:MemberID;references:ID"`
}

// TableName returns the table name of the struct_name model
func (w WorkExternalUserRelationModel) TableName() string {
	return "work_external_user_relation"
}

func (w WorkExternalUserRelationModel) DBPrefix() string {
	return "micro_scrm"
}

func (w WorkExternalUserRelationModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", w.DBPrefix(), enterpriseID, w.TableName())
}

type Bastatistical struct {
	MemberId int64 `json:"member_id"`
	Number   int64 `json:"number"`
}
