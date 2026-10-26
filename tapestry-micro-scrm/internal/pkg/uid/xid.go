package uid

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"micro-scrm/internal/config"
	"time"
)

var XID *ToolBoxUID

type ToolBoxUID struct {
	xidUrl  string
	xIdsUrl string
}

func InitUID() {
	xidUrl := config.CfgData.Xid
	if xidUrl == "" {
		panic("xid访问地址没有配置")
	}

	xIdsUrl := config.CfgData.XIds
	if xIdsUrl == "" {
		panic("xids访问地址没有配置")
	}

	XID = &ToolBoxUID{
		xidUrl:  xidUrl,
		xIdsUrl: xIdsUrl,
	}
}

func (t ToolBoxUID) ID(c ctx.Context) int64 {
	client := clients.NewHTTPClient()
	res, err := client.Get(c, t.xidUrl)
	if err != nil {
		<-time.After(10 * time.Millisecond)
		res, err = client.Get(c, t.xidUrl)
	}
	helper.Must(err)
	id, err := res.Read(clients.Int64)
	helper.Must(err)
	return id.(int64)
}

func (t ToolBoxUID) Ids(c ctx.Context, num int64) []int64 {
	if num <= 0 {
		num = 1
	}
	t.xIdsUrl = t.xIdsUrl + fmt.Sprintf("/batch?num=%d", num)
	client := clients.NewHTTPClient()
	res, err := client.Get(c, t.xIdsUrl)
	if err != nil {
		<-time.After(10 * time.Millisecond)
		res, err = client.Get(c, t.xIdsUrl)
	}
	helper.Must(err)
	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Ids []int64 `json:"ids"`
		} `json:"data"`
	}
	_, err = res.Read(clients.JSON(&resp))
	helper.Must(err)
	return resp.Data.Ids
}

func ID(c ctx.Context) int64 {
	return XID.ID(c)
}

func Ids(c ctx.Context, num int64) []int64 {
	return XID.Ids(c, num)
}
