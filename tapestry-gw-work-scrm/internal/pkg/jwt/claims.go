package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/config"
	"time"
)

// 关于 Aud 的一些常量
const (
	MemberUser   = "member_user"
	ExternalUser = "external_user"
	OpenUser     = "open_user"
)

type CustomClaims struct {
	EnterpriseID int64  `json:"enterprise_id"`
	OpenId       string `json:"open_id"`
	jwt.StandardClaims
}

func (c CustomClaims) MemberUserId() (string, error) {
	if c.Audience != MemberUser {
		return "", errors.New("[jwt]非企业人员Token")
	}

	return c.Subject, nil
}

func (c CustomClaims) ExternalOpenId() (string, error) {
	if c.Audience == MemberUser {
		return "", errors.New("[jwt]非外部人员Token")
	}

	return c.Subject, nil
}

func (c CustomClaims) ExternalUserId() (string, error) {
	return c.ExternalOpenId()
}

// NewCustomClaims 生成负载信息
func NewCustomClaims(c ctx.Context, userId, openId, extUserId, url string, enterpriseID int64) CustomClaims {
	now := time.Now().Unix()
	ttl := config.CfgData.JWT.Ttl
	exp := time.Now().Add(time.Duration(ttl) * time.Minute).Unix()
	var sub, aud string

	if userId != "" {
		sub = userId
		aud = MemberUser
	}

	if openId != "" {
		sub = extUserId
		if extUserId != "" {
			aud = ExternalUser
		} else {
			aud = OpenUser
		}
	}

	return CustomClaims{
		enterpriseID,
		openId,
		jwt.StandardClaims{
			// 表示该人员是外部联系人还是内部人员
			Audience:  aud,
			ExpiresAt: exp,
			Id:        c.ID(),
			IssuedAt:  now,
			Issuer:    url,
			NotBefore: now,
			// 保存用户唯一信息(外部联系人保存userid，员工保存账号)
			Subject: sub,
		},
	}
}
