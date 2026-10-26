package workwechat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	workSDK "micro-scrm/internal/pkg/sdk/work"
	"sort"
)

// ManagerCustomerTag 管理外部联系人的标签，删除、新增
func (w *WorkWeChat) ManagerCustomerTag(request workSDK.UserTagHandlerRequest) (isSuccess bool, err error) {

	isSuccess, err = w.obj.UserTagHandler(w.ctx, request)
	if err != nil {
		log.Info(w.ctx).Msgf("更新外部联系人企微标签失败:" + err.Error())
		return
	}

	return
}

// 获取企业微信中企业设置的标签
func (w *WorkWeChat) TagList() ([]*workSDK.CorpTagGroupEntity, error) {
	data, err := w.obj.GetCorpTagList(w.ctx, []string{})
	if err != nil {
		log.Info(w.ctx).Msgf("获取企业微信标签列表错误:" + err.Error())
		return nil, err
	}
	sort.SliceStable(data.TagGroup, func(i, j int) bool {
		return data.TagGroup[i].CreateTime > data.TagGroup[j].CreateTime
	})
	return data.TagGroup, nil
}

// 根据企业微信的标签ID，删除标签
func (w *WorkWeChat) DeleteTag(params workSDK.DelCorpTagParam) error {
	err := w.obj.DelCorpTagInfo(w.ctx, params)
	if err != nil {
		log.Info(w.ctx).Msgf("删除企业微信标签错误:" + err.Error())
		return err
	}

	return nil
}

// 根据企业微信的标签ID，更新企业微信的名称
func (w *WorkWeChat) UpdateTagName(request workSDK.UpdateTagRequest) error {
	err := w.obj.UpdateCorpTagInfo(w.ctx, request)
	if err != nil {
		log.Info(w.ctx).Msgf("更新企业微信标签错误:" + err.Error())
		return err
	}

	return nil
}

// 新增企业微信标签/组
func (w *WorkWeChat) AddTag(request workSDK.AddCorpTagRequest) (*workSDK.AddCorpTagResponse, error) {
	response, err := w.obj.AddCorpTag(w.ctx, request)
	if err != nil {
		log.Info(w.ctx).Msgf("新增企业微信标签错误:" + err.Error())
		return nil, err
	}

	return response, nil
}
