package microScrmModel

import (
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
)

const (
	ManageIsSA = iota + 1
	ManageNotSA
)
const RelIsKeySa = 1
const (
	// WorkUserRelationStatusNormal 好友关系状态正常
	WorkUserRelationStatusNormal = 0
	// WorkUserRelationStatusDeleteEachOther 外部联系人关系互相删除
	WorkUserRelationStatusDeleteEachOther = 3
)

type WorkExternalUserRelationModel struct {
	models.BaseModel
	EnterpriseID        int64  `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;index:enterprise_id_index,good_friend_index,index_0;default:0;comment:企业id"`
	WorkExternalUserID  int64  `json:"work_external_user_id" gorm:"column:work_external_user_id;type:bigint(20);not null;index:index_0,work_external_user_id_index;default:0;comment:企业外部联系人信息表id"`
	MemberID            int64  `json:"member_id" gorm:"column:member_id;type:bigint(11);not null;index:good_friend_index,index_0,member_id_index;default:0;comment:成员表id"`
	Remark              string `json:"remark" gorm:"column:remark;type:varchar(255);not null;comment:成员设置的外部联系人备注"`
	Description         string `json:"description" gorm:"column:description;type:varchar(255);not null;comment:成员设置的外部联系人描述"`
	Createtime          int    `json:"createtime" gorm:"column:createtime;type:int(11);not null;index:createtime_index;default:0;comment:关系建立时间"`
	RemarkCorpName      string `json:"remark_corp_name" gorm:"column:remark_corp_name;type:varchar(255);not null;comment:成员备注的外部联系人企业名称"`
	RemarkMobiles       string `json:"remark_mobiles" gorm:"column:remark_mobiles;type:varchar(255);not null;comment:成员备注的外部联系人手机号"`
	AddWay              int32  `json:"add_way" gorm:"column:add_way;type:int(10) unsigned;not null;default:0;comment:该成员添加此外部联系人的来源"`
	State               string `json:"state" gorm:"column:state;type:varchar(255);not null;comment:企业自定义的state参数"`
	OperUserid          string `json:"oper_userid" gorm:"column:oper_userid;type:varchar(255);not null;default:0;comment:发起添加的userid，如果成员主动添加，为成员的userid；如果是客户主动添加，则为客户的外部联系人userid；如果是内部成员共享/管理员分配，则为对应的成员/管理员userid"`
	RelationStatus      int    `json:"relation_status" gorm:"column:relation_status;type:int(11);not null;index:relation_status_index;default:0;comment:外部联系人关系状态：0-正常，1-员工删除联系人，2-联系人删除员工，3-双方互相删除，4-关系转移中"`
	MemberDeletedTime   int    `json:"member_deleted_time" gorm:"column:member_deleted_time;type:int(11);not null;index:member_deleted_time_index;default:0;comment:成员主动删除关系时间"`
	ExternalDeletedTime int    `json:"external_deleted_time" gorm:"column:external_deleted_time;type:int(11);not null;index:external_deleted_time_index;default:0;comment:外部联系人删除关系时间"`
	UserID              string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;index:user_id_index;comment:成员表user_id"`
	Unionid             string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null;index:unionid_index;comment:外部联系人的unionid"`
	IsKeySa             int    `json:"is_key_sa" gorm:"column:is_key_sa;type:int(1);not null;comment:是否是专属导购：0否，1是"`

	Members membersModel.MemberModel `json:"member" gorm:"foreignKey:MemberID;references:ID"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (WorkExternalUserRelationModel) TableName() string {
	return "work_external_user_relation"
}

type RepeatCustomerRelation struct {
	WorkExternalUserId int64 `json:"work_external_user_id"`
	Num                int64 `json:"num"`
}

// UserSaInfo 外部联系人SA信息
type UserSaInfo struct {
	Id                 int64  `json:"id"`
	WorkExternalUserId int64  `json:"work_external_user_id"`
	UnionId            string `json:"unionid"`
	MemberId           int64  `json:"member_id"`
	Description        string `json:"description"`
	DepartmentId       int64  `json:"department_id"`
	DepartmentName     string `json:"department_name"`
	IsKeySa            int    `json:"is_key_sa" db:"is_key_sa"`   // 是否专属导购：0否，1是
	SaUserId           string `json:"sa_user_id" db:"sa_user_id"` // 是否专属导购：0否，1是
	SaUserName         string `json:"sa_user_name" db:"sa_user_name"`
	Remark             string `json:"remark" db:"remark"`
}

type ExternalUserList struct {
	Id                    int64   `json:"id"`
	EnterpriseID          int64   `json:"enterprise_id"`
	WorkExternalUserID    int64   `json:"work_external_user_id"`
	MemberID              int64   `json:"member_id"`
	Remark                string  `json:"remark"`
	Unionid               string  `json:"unionid"`
	Openid                string  `json:"openid"`
	Euid                  int64   `json:"euid"`
	Nick                  string  `json:"nick"`
	Name                  string  `json:"name"`
	Avatar                string  `json:"avatar"`
	Birthday              string  `json:"birthday"`
	UnionCode             string  `json:"union_code"` // 会员卡号
	Mobile                string  `json:"mobile"`
	CreatedAt             int     `json:"created_at"`
	MemberName            string  `json:"member_name"`
	IsKeySa               int     `json:"is_key_sa"` // 是否专属导购：0否，1是
	SaUserId              string  `json:"sa_user_id"`
	KeySaMemberID         int64   `json:"key_sa_member_id"`
	KeySaMemberName       string  `json:"key_sa_member_name"`
	IsRelation            int     `json:"is_relation"`
	CardLevel             string  `json:"card_level"`
	CustomerCard          string  `json:"customer_card"`
	CltType               string  `json:"clt_type"`
	ClientType            string  `json:"client_type"`
	AmountTotal           int64   `json:"amount_total"`
	Last12mPurchaseAmount float32 `json:"last_12m_purchase_amount"`
	Last36mPurchaseAmount float32 `json:"last_36m_purchase_amount"`
	CustomerCode          string  `json:"customer_code"`
	LastPurchaseDate      int     `json:"last_purchase_date"`
	LastPurchaseDateStr   string  `json:"last_purchase_date_str"`
	LastPurchaseStoreCode string  `json:"last_purchase_store_code"`
	LastPurchaseStoreName string  `json:"last_purchase_store_name"`
	UpgradeNeedAmount     int     `json:"upgrade_need_amount"`
	CustomerType          string  `json:"customer_type"`
	VipType               int     `json:"vip_type"`
	ValidReserveAt        int     `json:"valid_reserve_at"`
	SuccessReserveAt      int     `json:"success_reserve_at"`
	TaskFinishAt          int     `json:"task_finish_at"`
	CouponStatus          int     `json:"coupon_status"`
	CouponBy30daysStatus  int     `json:"coupon_by_30_days_status"`
}

type UserSaData struct {
	MemberId   int64  `json:"member_id"`
	Unionid    string `json:"unionid"`
	Userid     string `json:"userid"`
	QrCode     string `json:"qr_code"`
	IsKeySa    int    `json:"is_key_sa" db:"is_key_sa"` // 是否专属导购：0否，1是
	SaUserId   string `json:"sa_user_id" db:"sa_user_id"`
	SaUserName string `json:"sa_user_name" db:"sa_user_name"`
}

type UserSaListData struct {
	MemberId       int64  `json:"member_id"`
	Unionid        string `json:"unionid"`
	Userid         string `json:"userid"`
	ThumbAvatar    string `json:"thumb_avatar"`
	QrCode         string `json:"qr_code"`
	IsKeySa        int    `json:"is_key_sa" db:"is_key_sa"` // 是否专属导购：0否，1是
	SaUserId       string `json:"sa_user_id" db:"sa_user_id"`
	SaUserName     string `json:"sa_user_name" db:"sa_user_name"`
	ThirdId        int    `json:"third_id"`
	DepartmentName string `json:"department_name"`
	CreatedAt      int    `json:"created_at"`
}
