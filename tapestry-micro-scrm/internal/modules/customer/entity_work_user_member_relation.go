package customer

import (
	"strings"

	"micro-scrm/internal/pkg/sdk/work"
)

// WorkUserRelationStatusNormal 好友关系状态正常
const WorkUserRelationStatusNormal = 0

// WorkUserRelationStatusDeletedByMember 成员主动删除好友关系
const WorkUserRelationStatusDeletedByMember = 1

// WorkUserRelationStatusDeletedByExternal 外部联系人主动删除成员
const WorkUserRelationStatusDeletedByExternal = 2

// WorkUserRelationStatusDeleteEachOther 外部联系人关系互相删除
const WorkUserRelationStatusDeleteEachOther = 3

// WorkUserRelationStatusTransferWaitingForConfirmation 外部联系人关系-新关系确认中
const WorkUserRelationStatusTransferWaitingForConfirmation = 4

// WorkUserMemberRelation 外部联系人
type WorkUserMemberRelation struct {
	ID                  int64  `json:"id" db:"id" desc:"外部联系人关系主键"`
	EnterpriseID        int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	WorkExternalUserID  int64  `json:"work_external_user_id" db:"work_external_user_id" desc:"外部联系人主键"`
	MemberID            int64  `json:"member_id" db:"member_id" desc:"成员ID"`
	Remark              string `json:"remark" db:"remark" desc:"备注"`
	Description         string `json:"description" db:"description" desc:"成员设置的外部联系人描述"`
	CreateTime          int64  `json:"createtime" db:"createtime" desc:"关系建立时间"`
	RemarkCorpName      string `json:"remark_corp_name" db:"remark_corp_name" desc:"成员备注的外部联系人企业名称"`
	RemarkMobiles       string `json:"remark_mobiles" db:"remark_mobiles" desc:"成员备注的外部联系人手机号"`
	AddWay              int    `json:"add_way" db:"add_way" desc:"该成员添加此外部联系人的来源"`
	State               string `json:"state" db:"state" desc:"企业自定义的state参数"`
	OperUserID          string `json:"oper_userid" db:"oper_userid"`
	RelationStatus      int8   `json:"relation_status" db:"relation_status" desc:"外部联系人关系状态"`
	ExternalDeletedTime int64  `json:"external_deleted_time" db:"external_deleted_time" desc:"外部联系人删除成员时间"`
	MemberDeletedTime   int64  `json:"member_deleted_time" db:"member_deleted_time" desc:"成员删除外部联系人时间"`
	CreatedAt           int64  `json:"created_at" db:"created_at"`
	UpdatedAt           int64  `json:"updated_at" db:"updated_at"`
	DeletedAt           int64  `json:"deleted_at" db:"deleted_at"`
	// 非数据库
	ExternalUserID string
	UserID         string
	UserId         string `json:"user_id" db:"user_id" desc:"成员表user_id"`
	UnionId        string `json:"unionid" db:"unionid" desc:"外部联系人的unionid"`
	IsKeySa        uint   `json:"is_key_sa" db:"is_key_sa" desc:"是否专属导购：0否，1是"`
}

// WorkUserMemberRelationData 外部联系人信息
type WorkUserMemberRelationData struct {
	ID                  int64  `json:"id" db:"id" desc:"外部联系人关系主键"`
	EnterpriseID        int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	WorkExternalUserID  string `json:"work_external_user_id" db:"work_external_user_id" desc:"外部联系人主键"`
	MemberID            string `json:"member_id" db:"member_id" desc:"成员ID"`
	Remark              string `json:"remark" db:"remark" desc:"备注"`
	Description         string `json:"description" db:"description" desc:"成员设置的外部联系人描述"`
	CreateTime          int64  `json:"createtime" db:"createtime" desc:"关系建立时间"`
	RemarkCorpName      string `json:"remark_corp_name" db:"remark_corp_name" desc:"成员备注的外部联系人企业名称"`
	RemarkMobiles       string `json:"remark_mobiles" db:"remark_mobiles" desc:"成员备注的外部联系人手机号"`
	AddWay              int    `json:"add_way" db:"add_way" desc:"该成员添加此外部联系人的来源"`
	State               string `json:"state" db:"state" desc:"企业自定义的state参数"`
	OperUserID          string `json:"oper_userid" db:"oper_userid"`
	RelationStatus      int8   `json:"relation_status" db:"relation_status" desc:"外部联系人关系状态"`
	ExternalDeletedTime int64  `json:"external_deleted_time" db:"external_deleted_time" desc:"外部联系人删除成员时间"`
	MemberDeletedTime   int64  `json:"member_deleted_time" db:"member_deleted_time" desc:"成员删除外部联系人时间"`
	CreatedAt           int64  `json:"created_at" db:"created_at"`
	UpdatedAt           int64  `json:"updated_at" db:"updated_at"`
	DeletedAt           int64  `json:"deleted_at" db:"deleted_at"`
	// 非数据库
	ExternalUserID string `json:"external_userid" db:"openid" desc:"外部联系人userid"`
	UserID         string `json:"userid" db:"userid" desc:"成员userid"`
}

// WorkCustomerRelationResponse 外部联系人信息
type WorkCustomerRelationResponse struct {
	WorkUserMemberRelation
	MemberName   string `json:"member_name" db:"member_name"`
	MemberMobile string `json:"member_mobile" db:"member_mobile"`
	QyUserID     string `json:"qy_userid" db:"qy_userid"`
}

// BatchWorkUserMemberRelationRequest 外部联系人列表请求
type BatchWorkUserMemberRelationRequest struct {
	EnterpriseID       int64  `json:"enterprise_id" db:"enterprise_id"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id"`
	UserID             string `json:"userid"`
	UserIDCnt          int
	UserIds            []string `json:"userids"`
}

// BatchWorkUserMemberRelationResponse 外部联系人列表返回结果
type BatchWorkUserMemberRelationResponse struct {
	RelationID          int64  `json:"relation_id" db:"relation_id"`
	MemberUserID        string `json:"member_userid" db:"member_userid"`
	State               string `json:"state" db:"state"`
	RelationStatus      int8   `json:"relation_status" db:"relation_status"`
	MemberDeletedTime   int64  `json:"member_deleted_time" db:"member_deleted_time"`
	ExternalDeletedTime int64  `json:"external_deleted_time" db:"external_deleted_time"`
}

// ExternalRemark 联系人备注
type ExternalRemark struct {
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id"`
	MemberID           int64  `json:"member_id" db:"member_id"`
	Remark             string `json:"remark" db:"remark"`
}

// TransferExternalRelation 外部联系人数据转
func TransferExternalRelation(vItem work.ExternalFollowUser, response *WorkUserMemberRelation) {
	response.AddWay = vItem.AddWay
	response.Description = vItem.Description
	response.OperUserID = vItem.OperUserID
	response.CreateTime = vItem.CreateTime
	response.Remark = vItem.Remark
	response.RemarkCorpName = vItem.RemarkCorpName
	response.RemarkMobiles = strings.Join(vItem.RemarkMobiles, ",")
	response.State = vItem.State
}
