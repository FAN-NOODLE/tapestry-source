package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/db"
)

//New 初始化存储
func NewWsTrace() *WsTraceRepo {
	return &WsTraceRepo{repos.NewSQLRepo("", db.DB())}
}
