package workexternal

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"time"
)

type Service struct {
	repo Repository
}

func (ser Service) SelectReUv(c ctx.Context, param SelectKeep) (int64, error) {
	return ser.repo.SelectReUv(c, param)
}

func (ser Service) SelectLossUv(c ctx.Context, param SelectKeep) (int64, error) {
	return ser.repo.SelectLossUv(c, param)
}

func (ser Service) SelectAddUv(c ctx.Context, param SelectKeep) (int64, error) {
	return ser.repo.SelectAddUv(c, param)
}

// ExternalExists 查询外部联系人是否存在
func (ser Service) ExternalExists(userid string, enterpriseID int64) (bool, error) {
	return false, nil
}

// GetBetweenByTime 获取时间范围
func (ser Service) GetBetweenByTime(day int64) (start, end int64) {
	now := time.Now()
	// 今日
	if day <= 0 {
		end = now.Unix()
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
		return
	}
	//N天前的统计，结束时间为今天0点
	end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	// 获取n天前的日期，开始时间往前对n天
	start = end - day*86400
	return
}
