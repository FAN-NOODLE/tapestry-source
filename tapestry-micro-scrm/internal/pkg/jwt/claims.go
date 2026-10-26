package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
	"time"
)

type CustomClaims struct {
	EnterpriseID int64 `json:"enterprise_id"`
	ID           int64 `json:"id"`
	jwt.StandardClaims
}

// NewCustomClaims 生成负载信息
func NewCustomClaims(c ctx.Context, ID, enterpriseID int64) CustomClaims {
	now := time.Now()
	ttl := config.CfgData.JWT.Ttl

	return CustomClaims{enterpriseID, ID,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Duration(ttl) * time.Second).Unix(),
			Id:        c.ID(),
			IssuedAt:  now.Unix(),
			NotBefore: time.Now().Unix() - 1000,
			Subject:   config.CfgData.App,
		},
	}
}
