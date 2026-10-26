package drainage

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/redis"
	"strconv"
)

func (s *StatisticsAllReq) GetCache(ctx ctx.Context, key string) (int64, bool) {
	result := redis.Client().Get(ctx, key).Val()
	if result == "" {
		return 0, false
	}
	data, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, false
	}

	return data, true
}

func (s *StatisticsAllReq) SetCache(ctx ctx.Context, key string, data int64) {
	redis.Client().Set(ctx, key, data, 0)
}

func (s *StatisticsAllReq) CacheKey(key string) string {
	id := strconv.FormatInt(s.EnterpriseID, 10)
	start := strconv.FormatInt(s.StartTime, 10)
	end := strconv.FormatInt(s.EndTime, 10)
	return "statistics:" + id + ":" + start + end + key
}

func (s *StatisticsChannelReq) CacheKey(key string) string {
	id := strconv.FormatInt(s.EnterpriseID, 10)
	start := strconv.FormatInt(s.StartTime, 10)
	end := strconv.FormatInt(s.EndTime, 10)
	channels := ""
	if s.ChannelsState != 0 {
		channels = strconv.FormatInt(s.ChannelsState, 10) + ":"
	}
	return "statistics:" + id + ":" + channels + start + end + key
}
