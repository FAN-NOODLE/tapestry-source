package work

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/config"
)

// Taglist 获取标签列表
func (w *QyWorkImpl) Taglist(c ctx.Context) (TaglistResponse, error) {
	var response TaglistResponse

	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}
	req := QyGwRequest{
		ReqPL:     QWExternalIdentifier,
		ReqMethod: POST,
		ReqURL:    Taglist,
		Token:     token,
	}

	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return response, err
	}
	if response.ErrCode == 0 {
		return response, nil
	}

	return response, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
}

// Taglist 创建标签
func (w *QyWorkImpl) TagCreate(c ctx.Context, tagName string, tagId int) (int, error) {
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return 0, err
	}
	var req struct {
		QyGwRequest
		TagName string `json:"tagname"`
		TagId   int    `json:"tagid"`
	}
	req.Token = token
	req.ReqURL = TagCreate
	req.ReqMethod = POST
	req.ReqPL = QWExternalIdentifier
	req.TagName = tagName
	req.TagId = tagId

	var response struct {
		QyResponse
		TagId int `json:"tagid"`
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return 0, err
	}
	if response.ErrCode == 0 {
		return response.TagId, nil
	}

	return 0, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
}

// 获取标签成员
func (w *QyWorkImpl) TagUsers(c ctx.Context, tagId int) (TagUsersResponse, error) {
	var response TagUsersResponse
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	var req struct {
		QyGwRequest
		TagId int `json:"tagid"`
	}
	req.Token = token
	req.ReqURL = TagUsers
	req.ReqMethod = POST
	req.ReqPL = QWExternalIdentifier
	req.TagId = tagId

	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return response, err
	}
	if response.ErrCode == 0 {
		return response, nil
	}

	return response, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
}

// 增加标签成员
func (w *QyWorkImpl) TagUsersAdd(c ctx.Context, tagid int, userList []string) (QyResponse, error) {
	var response QyResponse
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}
	var req struct {
		QyGwRequest
		TagId    int      `json:"id"`
		UserList []string `json:"name"`
	}
	req.Token = token
	req.ReqURL = TagUsersAdd
	req.ReqMethod = POST
	req.ReqPL = QWExternalIdentifier
	req.TagId = tagid
	req.UserList = userList

	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return response, err
	}
	if response.ErrCode == 0 {
		return response, nil
	}

	return response, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
}
