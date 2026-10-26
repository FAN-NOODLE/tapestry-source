package display

import (
	"gw-work-scrm/internal/app/http/request"
)

type TreeReq struct {
	request.EnterpriseBaseReq
	Language         string `json:"language"`
	ParentId         int    `json:"parent_id"`
	HasChildren      bool   `json:"has_children"`
	HasVirtual       bool   `json:"has_virtual"`
	MemberId         int64  `json:"member_id"`
	IgnoreCloseStore bool   `json:"ignore_close_store"`
}

type SearchStore struct {
	request.EnterpriseBaseReq
	Language         string `json:"language"`
	MemberId         int64  `json:"member_id"`
	StoreStr         string `json:"store_str"`
	IgnoreCloseStore bool   `json:"ignore_close_store"`
}

type MemberAuthDetail struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}
