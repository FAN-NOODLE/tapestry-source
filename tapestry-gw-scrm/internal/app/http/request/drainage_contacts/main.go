package drainageContactsReq

import "gw-scrm/internal/app/http/request"

type BaseReq struct {
	request.EnterpriseBaseReq
	AdminId int64 `json:"admin_id"`
}

type GroupListReq struct {
	BaseReq
	Name       string        `json:"name"`
	Type       int           `json:"type"`
	ChannelId  int64         `json:"channel_id"`
	CreateTime TimeSearchReq `json:"create_time"`
	request.PageReq
}

type TimeSearchReq struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type DetailListReq struct {
	InfoReq
	request.PageReq
}

type InfoReq struct {
	BaseReq
	ID int64 `json:"id"`
}

type GroupInfoReq struct {
	InfoReq
	State string `json:"state"`
}

type GroupSaveReq struct {
	InfoReq
	Name        string        `json:"name" binding:"required" `
	ChannelId   int64         `json:"channel_id" binding:"required" desc:"渠道id"`
	AutoFriend  int           `json:"auto_friend"   db:"auto_friend" desc:"是否自动开启"`
	MemberIds   []int64       `json:"member_ids"`
	ThirdIds    []int         `json:"third_ids" desc:"部门Id"`
	Type        int           `json:"type"  db:"type" desc:"类型 0 多人共享码（原本的） 1 员工码+门店码"`
	WelcomeList []WelcomeInfo `json:"welcome_list"`
}

type WelcomeInfo struct {
	ID      int64                    `json:"id" db:"id" desc:"id"`
	Type    int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content []map[string]interface{} `json:"content" binding:"required" desc:"内容"`
}

type ZipUrlReq struct {
	InfoReq
	Width int `json:"width"`
}
