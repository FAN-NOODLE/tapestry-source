package memberReq

import (
	membersModel "micro-scrm/internal/app/models/members"
)

type MemberDetailReq struct {
	EnterpriseId int64   `json:"enterprise_id" binding:"required" comment:"企业ID"`
	MemberIds    []int64 `json:"member_ids" comment:"服务导购的id"`
	StoreIds     []int64 `json:"store_ids" comment:"服务的门店id"`
	Mobile       string  `json:"mobile" comment:"手机号"`
}

type MemberDetailRes struct {
	membersModel.MemberDepartment
	Address   string `json:"address"`
	StoreCode string `json:"store_code"`
	StoreName string `json:"store_name"`
}

type MemberStoreReq struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	ThirdIds     []int  `json:"third_ids" binding:"required" comment:"门店id"`
	Mobile       string `json:"mobile"`
	IsSm         int    `json:"is_sm" desc:"0为店员,即导购,1为店长，即企业管理"`
}

type MemberStoreRes struct {
	DepartmentId int64  `json:"department_id"`
	StoreId      int    `json:"store_id"`
	StoreName    string `json:"store_name"`
}

type MemberStaffReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	StoreIds     []int  `json:"store_ids"`
	Mobile       string `json:"mobile"`
	UnionId      string `json:"unionid"`
}

type MemberStaffRes struct {
	StoreId       int    `json:"store_id"`
	StoreName     string `json:"store_name"`
	MembersId     int64  `json:"members_id"`
	MembersAvatar string `json:"members_avatar"`
	MembersName   string `json:"members_name"`
	MembersUserid string `json:"members_userid"`
	Position      string `json:"position"`
	Mobile        string `json:"mobile"`
}

type MemberListByThirdIdsReq struct {
	EnterpriseId int64   `json:"enterprise_id"`
	ThirdIds     []int   `json:"third_ids"`
	DpIds        []int64 `json:"dp_ids"`
}

type ReportDataReq struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	UserId       string `json:"userid" comment:"导购的id"`
}

type FirstOpenReq struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Type         int    `json:"type"`
	Mobile       string `json:"mobile"`
}

type StoreCodesReq struct {
	EnterpriseId int64    `json:"enterprise_id" binding:"required" comment:"企业ID"`
	StoreCodes   []string `json:"store_codes" binding:"required" comment:"企业ID"`
}

// MemberSearch  基本企业成员信息
type MemberSearch struct {
	Userid string `json:"user_id" db:"userid" desc:"企业微信userid"`
	Name   string `json:"user_name" db:"name" desc:"员工姓名"`
	Mobile string `json:"mobile"  db:"mobile" desc:"员工手机号"`
}
