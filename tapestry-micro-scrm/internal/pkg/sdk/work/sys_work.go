package work

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
)

// AddCallbackUrl 配置外部联系人回调
func (w *QyWorkImpl) AddCallbackUrl(cc ctx.Context, req SCallbackAdd) error {
	url := fmt.Sprintf(CallbackAdd, config.CfgData.Service.Work.Domain)
	return Post(cc, url, req, nil)
}

// CallbackDel 外部联系人回调删除
func (w *QyWorkImpl) CallbackDel(cc ctx.Context, req SCallbackDel) error {
	url := fmt.Sprintf(CallbackDel, config.CfgData.Service.Work.Domain)
	return Post(cc, url, req, nil)
}

// CallbackList 外部联系人回调地址列表
func (w *QyWorkImpl) CallbackList(cc ctx.Context, req SCallbackList) (*[]CallbackItem, error) {
	url := fmt.Sprintf(CallbackList, config.CfgData.Service.Work.Domain)
	type response struct {
		SBaseResponse
		Data *[]CallbackItem `json:"data"`
	}
	list := &response{}
	err := Post(cc, url, req, list)
	if err != nil {
		return nil, err
	}
	return list.Data, nil
}
