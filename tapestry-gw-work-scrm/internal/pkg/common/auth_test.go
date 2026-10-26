package common

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/redis"
	"reflect"
	"testing"
)

func TestGetEnterpriseInfo(t *testing.T) {
	config.CfgData.Redis.Addr = "127.0.0.1:6379"
	config.CfgData.Redis.Pwd = "123456"
	config.CfgData.Service.Boss.Domain = "https://boss-api-test-crs.vchangyi.com"
	helper.Must(redis.Init(config.CfgData.Redis))
	type args struct {
		c            ctx.Context
		enterpriseID int64
	}
	tests := []struct {
		name    string
		args    args
		want    *EnterpriseInfo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "redis数据获取",
			args: args{
				c:            ctx.New(),
				enterpriseID: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEnterpriseInfo(tt.args.c, tt.args.enterpriseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnterpriseInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnterpriseInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
