package activityQRCodeReq

import (
	"micro-scrm/internal/app/http/request"
)

type BaseReq struct {
	request.EnterpriseBaseReq
	AdminId int64 `json:"admin_id"`
}

type ActivityCategoryListReq struct {
	BaseReq
	AddWelcome bool `json:"add_welcome"`
}

type ActivityListReq struct {
	BaseReq
	Name       string        `json:"name"`
	CategoryID int64         `json:"category_id"`
	EffectTime TimeSearchReq `json:"effect_time"`
	CreateTime TimeSearchReq `json:"create_time"`
	request.PageReq
}

type TimeSearchReq struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type InfoReq struct {
	BaseReq
	ID int64 `json:"id"`
}

type ActivityInfoReq struct {
	InfoReq
	State string `json:"state"`
}

type ActivitySaveReq struct {
	InfoReq
	Name           string         `json:"name"`
	AliasName      string         `json:"alias_name"`
	BannerUrl      string         `json:"banner_url"`
	TimeType       int            `json:"time_type"`
	StartTime      int            `json:"start_time"`
	EndTime        int            `json:"end_time"`
	SkipVerify     int            `json:"skip_verify"`
	AddNameTag     int            `json:"add_name_tag"`
	AddCategoryTag int            `json:"add_category_tag"`
	TagIds         []int64        `json:"tag_ids"`
	ScopeMemberIds []int64        `json:"scope_member_ids"`
	ScopeAreaIds   []int          `json:"scope_area_ids"`
	CategoryList   []CategoryInfo `json:"category_list"`
	WelcomeList    []WelcomeInfo  `json:"welcome_list"`
}

type CategoryInfo struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Points []string `json:"points"`
}
type WelcomeInfo struct {
	ID      int64                    `json:"id" db:"id" desc:"id"`
	Type    int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content []map[string]interface{} `json:"content" binding:"required" desc:"内容"`
}

type WelcomeListReq struct {
	BaseReq
	Content    string `json:"content"`
	CategoryID int64  `json:"category_id"`
	request.PageReq
}

type WelcomeInfoReq struct {
	BaseReq
	ID         int64                    `json:"id"`
	CategoryID int64                    `json:"category_id"  binding:"required"`
	Content    []map[string]interface{} `json:"content" binding:"required" desc:"内容"`
}

type StatisticsReq struct {
	BaseReq
	NameStr      string `json:"name_str"`
	PointNameStr string `json:"point_name_str"`
	CategoryID   int64  `json:"category_id"`
	TimeSearchReq
	request.PageReq
}
