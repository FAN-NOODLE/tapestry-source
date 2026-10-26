package fission

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

//go:generate mockgen -source=dependency.go -destination=mock_dependency.go -package=fission
type ContactInterface interface {
	GetContactQrCode(c ctx.Context, ID int64, enterpriseID int64) (string, string, error)
	CreateContact(c ctx.Context, enterpriseID int64, ID int64, name string, channelID int64, qrCodeUrl string, configID string, cType int) error
}

type EnterpriseInterface interface {
	GetInfo(c ctx.Context, enterpriseID int64) (domain string, corpID string, hash string, err error)
	GetOfficialInfo(c ctx.Context, enterpriseID int64) (appID string, secret string, err error)
}

type UserInterface interface {
	GetWorkUserInfoByUserIDs(c ctx.Context, enterpriseID int64, userIDs []int64) (*[]*WorkUserBasicInfo, error)
	GetRobotUserList(c ctx.Context, num int) ([]*AssistantUserInfo, error)
}

type BossApiInterface interface {
	GetShopAuthInfo(c ctx.Context, enterpriseID int64, authType int) (*ShopAuthInfo, error)
}
