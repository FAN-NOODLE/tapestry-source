package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"golang.org/x/time/rate"
	"log"
	"sync"
)

func NewAccountRateMiddleware(limit *AccountRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.GetInt64(UserId)
		log.Println("------OpenAccountID-----------", accountID)
		limiter := limit.GetLimiter(accountID)
		if !limiter.Allow() {
			ginplus.ResError(c, errors.New("访问超出限制"))
			c.Abort()
			return
		} else {
			log.Println("---------------------------------success---------------------------")
			c.Next()
		}
	}
}

type AccountRateLimiter struct {
	accounts map[int64]*rate.Limiter
	mu       *sync.RWMutex
	r        rate.Limit
	b        int
}

// SetupAccountRateLimiter 创建一个按账号限流的RateLimiter
func SetupAccountRateLimiter(r, b int) *AccountRateLimiter {
	return &AccountRateLimiter{
		accounts: make(map[int64]*rate.Limiter),
		mu:       &sync.RWMutex{},
		r:        rate.Limit(r),
		b:        b,
	}
}

// AddAccount 添加一个账号到map
func (i *AccountRateLimiter) AddAccount(accountID int64) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.accounts[accountID] = limiter
	return limiter
}

// GetLimiter 通过账号得到limiter
func (i *AccountRateLimiter) GetLimiter(accountID int64) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.accounts[accountID]
	if !exists {
		i.mu.Unlock()
		return i.AddAccount(accountID)
	}
	i.mu.Unlock()
	return limiter
}
