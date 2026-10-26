package user

import (
	"time"

	"micro-scrm/internal/modules/customer"
)

const (
	ExternalDescriptionLength = 150

	OfficialSubscribeStatusOn  = 1 // 公众号关注状态：关注中
	OfficialSubscribeStatusOff = 2 // 公众号关注状态：已取关

	EnterpriseOfficialBindTypeOn  = 1 // 企业公众号绑定状态
	EnterpriseOfficialBindTypeOff = 2 // 企业公众号绑定状态
)

type DDCRequestBase struct {
	Method   string `json:"method"`
	Platform string `json:"platform"`
}

// DDCRequestBody DDC请求数据结构
type DDCRequestBody struct {
	DDCRequestBase
	Data interface{} `json:"data"`
}

// SyncShopUserRequest 同步商城用户接口请求参数
type SyncShopUserRequest struct {
	DDCRequestBase
	UserInfo ShopUserEntity `json:"data"`
}

// SyncShopUserAddressRequest 同步用户地址接口请求参数
type SyncShopUserAddressRequest struct {
	DDCRequestBase
	UserAddress ShopUserAddress `json:"data"`
}

type SyncUserEventRequest struct {
	DDCRequestBase
	Event EventEntity `json:"data"`
}

type SyncOfficialUserRequest struct {
	DDCRequestBase
	UserInfo OfficialUserEntity `json:"data"`
}

// ShopUserEntity 商城用户数据结构
type ShopUserEntity struct {
	BasicUser
	SubscribeTime   int    `json:"subscribe_time"`
	Remark          string `json:"remark"`
	Groupid         int    `json:"groupid"`
	SubscribeScene  int    `json:"subscribe_scene"`
	QRScene         string `json:"qr_scene"`
	QRSceneStr      string `json:"qr_scene_str"`
	AuthorizerAppid string `json:"authorizer_appid"`
	RealName        string `json:"real_name"`
}

type BasicUser struct {
	ID            int64  `json:"id" db:"id"`
	ThirdUserID   int64  `json:"user_id"`
	UserID        int64  `db:"user_id"`
	UnionID       string `json:"union_id" db:"union_id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	ShopName      string `json:"shop_name"`
	Avatar        string `json:"avatar" db:"avatar"`
	OpenId        string `json:"open_id" db:"open_id"`
	Nickname      string `json:"nick_name" db:"nick_name"`
	Name          string `json:"name" db:"name"`
	Mobile        string `json:"mobile" db:"mobile"`
	Birthday      string `json:"birthday" db:"birthday"`
	Sex           int    `json:"gender" db:"sex"`
	Source        int    `db:"source"`
	ThirdMemberNo string `json:"thirdMemberNo" db:"third_member_no"`
	Country       string `json:"country" db:"country"`
	Province      string `json:"province" db:"province"`
	City          string `json:"city" db:"city"`
	CreatedAt     int    `json:"created_at" db:"created_at"`
	UpdatedAt     int    `db:"updated_at"`
	DeletedAt     int    `db:"deleted_at"`
}

// OfficialUserEntity 公众号用户结构体
type OfficialUserEntity struct {
	BasicUser
	SubscribeTime   int    `json:"subscribe_time" db:"subscribe_time"`
	Remark          string `json:"remark" db:"remark"`
	Groupid         int    `json:"groupid" db:"groupid"`
	SubscribeScene  string `json:"subscribe_scene" db:"subscribe_scene"`
	SubscribeStatus int    `json:"subscribe_status" db:"subscribe_status"`
	QRScene         string `json:"qr_scene" db:"qr_scene"`
	QRSceneStr      string `json:"qr_scene_str" db:"qr_scene_str"`
	AuthorizerAppid string `json:"authorizer_appid" db:"appid"`
	Language        string `json:"language" db:"language"`
	OfficialSex     int    `json:"sex"`
}

// SourceType 用户来源类型
type SourceType int

// SourceTypeEnum 用户来源枚举
var SourceTypeEnum = struct {
	Shop SourceType
	Work SourceType
}{
	Shop: 1,
	Work: 2,
}

// TranslateToUser 公众号用户表数据转主用户表数据
func (s *OfficialUserEntity) TranslateToUser() *UserEntity {
	return &UserEntity{
		EnterpriseID:    s.EnterpriseID,
		UnionID:         s.UnionID,
		Nickname:        s.Nickname,
		Name:            s.Name,
		Mobile:          s.Mobile,
		Avatar:          s.Avatar,
		Gender:          s.Sex,
		FirstSourceType: int(SourceTypeEnum.Shop),
		CreatedAt:       int(time.Now().Unix()),
		UpdatedAt:       int(time.Now().Unix()),
	}
}

// TranslateToUser 商城用户表数据转主用户表数据
func (b *BasicUser) TranslateToUser() *UserEntity {
	return &UserEntity{
		EnterpriseID:    b.EnterpriseID,
		UnionID:         b.UnionID,
		Nickname:        b.Nickname,
		Name:            b.Name,
		Mobile:          b.Mobile,
		Avatar:          b.Avatar,
		Gender:          b.Sex,
		FirstSourceType: int(SourceTypeEnum.Shop),
		CreatedAt:       int(time.Now().Unix()),
		UpdatedAt:       int(time.Now().Unix()),
	}
}

// UserEntity user表实体结构
type UserEntity struct {
	ID              int64  `json:"id" db:"id"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	UnionID         string `json:"union_id" db:"union_id"`
	Nickname        string `json:"nickname" db:"nickname"`
	Name            string `json:"name" db:"name"`
	Gender          int    `json:"gender" db:"gender"`
	Mobile          string `json:"mobile" db:"mobile"`
	Avatar          string `json:"avatar" db:"avatar"`
	Remark          string `json:"remark" db:"remark"`
	FirstSourceType int    `json:"first_source_type" db:"first_source_type"`
	CreatedAt       int    `db:"created_at"`
	UpdatedAt       int    `db:"updated_at"`
	DeletedAt       int    `db:"deleted_at"`
}

// Event 事件结构体
type EventEntity struct {
	ID           int64       `json:"id" db:"id"`
	EnterpriseID int64       `json:"enterprise_id" db:"enterprise_id"`
	UnionID      string      `json:"union_id"`
	UserID       int64       `db:"user_id"`
	Event        string      `json:"event" db:"event"`
	Object       string      `json:"object" db:"object"`
	EventDesc    string      `json:"event_desc" db:"event_desc"`
	Extra        interface{} `json:"extra" db:"extra"`
	CreatedAt    int         `json:"event_at" db:"created_at"`
	UpdatedAt    int         `db:"updated_at"`
	DeletedAt    int         `db:"deleted_at"`
}

// ShopUserAddress 商城用户地址实体
type ShopUserAddress struct {
	ID             int64  `db:"id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`
	ShopUserID     int64  `db:"shop_user_id"` // ShopUserID 商城用户id
	AppID          string `json:"app_id"`
	OpenID         string `json:"open_id"`
	UnionID        string `json:"union_id"`                             // UnionID 用户union_id
	DeliveryName   string `json:"name" db:"delivery_name"`              // DeliveryName 收货人姓名
	DeliveryMobile string `json:"mobile" db:"delivery_mobile"`          // DeliveryMobile 收货人手机号码
	Province       string `json:"province" db:"province"`               // Province 省
	City           string `json:"city" db:"city"`                       // City 市
	Area           string `json:"area" db:"area"`                       // Area 区
	Address        string `json:"address" db:"address"`                 // Address 详细地址
	ProvinceCode   string `json:"province_code" db:"province_code"`     // ProvinceCode 省编码
	CityCode       string `json:"city_code" db:"city_code"`             // CityCode 市编码
	AreaCode       string `json:"area_code" db:"area_code"`             // AreaCode 区编码
	ZipCode        string `json:"zip_code" db:"zip_code"`               // ZipCode 邮编
	Default        int    `json:"default" db:"default"`                 // Default 是否默认地址(1:是，2:否)
	IsAuthAddress  int    `json:"is_auth_address" db:"is_auth_address"` // IsAuthAddress 商品类型 1=容许，2=不容许
	Source         int    `json:"source" db:"source"`                   // Source 微信地址(1:是，2:否)
	AddressID      int64  `json:"address_id" db:"address_id"`
	CreatedAt      int    `db:"created_at"`
	UpdatedAt      int    `db:"updated_at"`
	DeletedAt      int    `db:"deleted_at"`
}

const (
	EditTypeOne      = 1
	EditTypeMultiple = 2
	EditTypeAll      = 3
)

// EditRemarksRequest 编辑备注请求结构体
type EditRemarksRequest struct {
	customer.ListCond
	EditType   int     `json:"edit_type"`
	RelationID int64   `json:"relation_id"`
	UserIDs    []int64 `json:"user_ids"`
	Remarks    string  `json:"remarks"`
	UserID     int64   `json:"user_id"`
}

// EditRemarksResponse 编辑备注返回结构体
type EditRemarksResponse struct {
	RecordID      int64 `json:"record_id"`
	SuccessNumber int   `json:"success_number"`
	FailedNumber  int   `json:"failed_number"`
}

// EditRemarksResponse 编辑备注返回结构体
type ListConfigResponse struct {
	ShowOfficial bool `json:"show_official"`
	ShowShop     bool `json:"show_shop"`
}

// ExternalRelationAndUserEntity 用户及外部联系人数据
type ExternalRelationAndUserEntity struct {
	UserEntity
	RelationID         int64  `json:"relation_id" db:"relation_id"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id" desc:"外部联系人主键"`
	MemberID           int64  `json:"member_id" db:"member_id" desc:"成员ID"`
	Description        string `json:"description" db:"description" desc:"成员设置的外部联系人描述"`
}

type SourceUserID struct {
	UnionUID int64 `db:"union_uid"`
	OpenUID  int64 `db:"open_uid"`
}

type SourceUser struct {
	EnterpriseID int64
	customer.SourceUserBase
	ThirdUserId int64 //第三方唯一标识
	customer.EnterpriseUserBase
	customer.EnterpriseUserLocation
	customer.EnterpriseUserSubscribe
	customer.EnterpriseUserExternal
	customer.EnterpriseUserOther
	CreateAt int64 `db:"create_at"`
	SourceUserID
}

// RobotUserEntity 机器人实体
type RobotUserEntity struct {
	ID        int64  `db:"id" json:"id"`
	NickName  string `db:"nick_name" json:"nick_name"`
	Sex       int    `db:"sex" json:"sex"`
	ImgUrl    string `db:"img_url" json:"img_url"`
	CreatedAt int    `db:"created_at" json:"created_at"`
	UpdatedAt int    `db:"updated_at" json:"updated_at"`
	DeletedAt int    `db:"deleted_at" json:"deleted_at"`
}
