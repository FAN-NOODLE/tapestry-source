package redis

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"time"
)

var (
	EsDocKey    = "user_es:%s:%s-%s" //es文档
	canalPosKey = "user_es:pos:"
)

const (
	EsDocKeyTtl = 5
)

func getEsDocKey(ctx ctx.Context, index, id, md5v string) string {
	return fmt.Sprintf(EsDocKey, index, id, md5v)
}
func SetEsDoc(ctx ctx.Context, index, id, md5v, doc string) (bool, error) {
	key := getEsDocKey(ctx, index, id, md5v)
	b, err := redis.Client().SetNX(ctx, key, doc, EsDocKeyTtl*time.Second).Result()
	return b, err
}
func GetEsDoc(ctx ctx.Context, index, id, md5v string) string {
	key := getEsDocKey(ctx, index, id, md5v)
	return redis.Client().Get(ctx, key).String()
}
func DelEsDoc(ctx ctx.Context, index, id, md5v string) {
	key := getEsDocKey(ctx, index, id, md5v)
	redis.Client().Del(ctx, key)
}

func getcanalPosKey() string {
	masterdsn := helper.Md5V(config.CfgData.DB.MasterDSN)
	return canalPosKey + config.CfgData.Env + ":" + masterdsn
}

//获取最新的canal pos
func GetCanalPos(ctx ctx.Context) (string, error) {
	return redis.Client().Get(ctx, getcanalPosKey()).Result()
}

//设置最新的canal pos
func SetCanalPos(ctx ctx.Context, val string) (string, error) {
	return redis.Client().Set(ctx, getcanalPosKey(), val, 240*time.Hour).Result()
}

//删除canal pos
func DelCanPos(ctx ctx.Context) (int64, error) {
	return redis.Client().Del(ctx, getcanalPosKey()).Result()
}
