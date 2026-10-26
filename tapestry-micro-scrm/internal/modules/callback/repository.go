// Package callback 对接Open 平台
package callback

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储定义
type Repository interface {
	repos.IRepository
}
