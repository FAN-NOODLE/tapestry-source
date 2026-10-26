package externalSearchModel

import (
	"micro-scrm/internal/app/models"
)

type ExternalRelationshipModel struct {
	models.BaseModel
	WorkExternalUserID  int64   `gorm:"column:work_external_user_id;type:bigint(20);not null;default:0;comment:企业外部联系人信息表id"`
	Unionid             string  `gorm:"column:unionid;type:varchar(255);not null;comment:外部联系人的unionid"`
	MemberID            int64   `gorm:"column:member_id;type:bigint(11);not null;default:0;comment:成员表id"`
	UserID              string  `gorm:"column:user_id;type:varchar(255);not null;comment:成员表user_id"`
	Remark              string  `gorm:"column:remark;type:varchar(255);not null;comment:成员设置的外部联系人备注"`
	Description         string  `gorm:"column:description;type:varchar(255);not null;comment:成员设置的外部联系人描述"`
	Createtime          int64   `gorm:"column:createtime;type:int(11);not null;default:0;comment:关系建立时间"`
	RemarkMobiles       string  `gorm:"column:remark_mobiles;type:varchar(255);not null;comment:成员备注的外部联系人手机号"`
	AddWay              int64   `gorm:"column:add_way;type:int(10) unsigned;not null;default:0;comment:该成员添加此外部联系人的来源"`
	State               string  `gorm:"column:state;type:varchar(255);not null;comment:企业自定义的state参数"`
	RelationStatus      int64   `gorm:"column:relation_status;type:int(11);not null;default:0;comment:外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中"`
	MemberDeletedTime   int64   `gorm:"column:member_deleted_time;type:int(11);not null;default:0;comment:成员主动删除关系时间"`
	ExternalDeletedTime int64   `gorm:"column:external_deleted_time;type:int(11);not null;default:0;comment:外部联系人删除关系时间"`
	IsKeySa             int64   `gorm:"column:is_key_sa;type:int(11);not null;default:0;comment:是否是专属导购关系，0否 1是"`
	InteractTime        int64   `gorm:"column:interact_time;type:int(11);not null;default:0;comment:最近一次联系时间"`
	ActiveTime          int64   `gorm:"column:active_time;type:int(11);not null;default:0;comment:最近一次活跃沟通时间"`
	ExternalTagIds      []int64 `json:"external_tag_ids" gorm:"column:external_tag_ids;type:Array(Int64);"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (ExternalRelationshipModel) TableName() string {
	return "external_relationship"
}

const (
	RelationNormal     = iota //好友关系正常
	RelationDelContact        // 导购删除联系人
	RelationDelFollow         // 联系人删除导购
	RelationDelAll            // 双方互删
)
