// 定义存储使用的接口和辅助方法
package user

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository
}
