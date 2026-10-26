package ws_trace

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	InsertWsTrace(c ctx.Context, enterpriseID int64, wst WsTrace) (wstr *WsTrace, err error)
	GetWsTraceByPrimary(c ctx.Context, enterpriseID int64, ukVoucher string) (*WsTrace, error)
	GetWsTraceByID(c ctx.Context, enterpriseID int64, wsTraceID int64) (*WsTrace, error)
	SelectWsTraces(c ctx.Context, enterpriseID int64, sourceType int64) ([]WsTrace, error)
}
