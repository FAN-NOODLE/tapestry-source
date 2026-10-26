package work

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"micro-scrm/internal/orm"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk"
	"micro-scrm/internal/pkg/sdk/boss"

	bossModel "micro-scrm/internal/orm/model/boss"
)

// token key
var sTokenKey = "REDACTED_FOR_SOURCE_SHARING"

// SetDomain 设置domain
func (w *QyWorkImpl) SetDomain(domain string) *QyWorkImpl {
	w.domain = domain
	return w
}

func (w *QyWorkImpl) SetEnterprise(enterpriseID int64) *QyWorkImpl {
	enterpriseInfo, _ := orm.BossService.Enterprise.GetInfo(bossModel.EnterpriseModel{ID: enterpriseID})
	w.domain = enterpriseInfo.Domain
	return w
}

func NewQyWorkImpl() *QyWorkImpl {
	return &QyWorkImpl{}
}

// SetDomainValue 实现Api接口，设置domain
func (w *QyWorkImpl) SetDomainValue(domain string) {
	w.SetDomain(domain)
}

// GetToken 获取token
func (w *QyWorkImpl) GetToken(ctx ctx.Context) (string, error) {
	tokenCache := redis.Client().Get(ctx, sTokenKey).Val()
	if tokenCache != "" {
		return tokenCache, nil
	}
	sTokenRes, err := w.GetTokenFromQy(ctx)
	if err != nil {
		return "", err
	}
	ttl := sTokenRes.Expire/1000 - time.Now().Unix()
	redis.Client().Set(ctx, sTokenKey, sTokenRes.AccessToken, time.Duration(ttl)*time.Second)
	return sTokenRes.AccessToken, nil
}

// GetTokenFromQy 从企业微信获取AccessToken和过期值
func (w *QyWorkImpl) GetTokenFromQy(ctx ctx.Context) (*STokenResponseData, error) {
	sTokenRes := &STokenResponse{}
	url := fmt.Sprintf(SToken, config.CfgData.Service.Work.Domain)
	req := STokenRequest{
		ID:     config.CfgData.Service.Work.ID,
		Secret: config.CfgData.Service.Work.Secret,
	}
	rErr := workRequest(ctx, url, req, sTokenRes)
	if rErr != nil {
		return nil, rErr
	}
	log.Info(ctx).Msgf("获取到的参数信息：%+v", sTokenRes)
	return &sTokenRes.Data, nil
}

// QyDecrypt 解密回调数据
func (w *QyWorkImpl) QyDecrypt(ctx ctx.Context, req QyDecryptRequest) (*MemberExternalCallback, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.PlIdentifier = QWExternalIdentifier
	req.Token = token
	// 返回结果
	type data struct {
		Data MemberExternalCallback `json:"data"`
	}
	response := &data{}
	url := fmt.Sprintf(QyDecryptUrl, config.CfgData.Service.Work.Domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}

func (w *QyWorkImpl) QyContactDecrypt(ctx ctx.Context, req QyDecryptRequest) (string, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return "", tErr
	}
	req.PlIdentifier = QWExternalIdentifier
	req.Token = token

	url := fmt.Sprintf(QyDecryptUrl, config.CfgData.Service.Work.Domain)
	return Request(ctx, url, req)
}

// WorkDepartmentList 组织列表
func (w *QyWorkImpl) WorkDepartmentList(ctx ctx.Context, req WorkDepartmentListRequest) ([]WorkDepartment, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = DepartmentList
	req.ReqMethod = "GET"
	// 返回结果
	type depList struct {
		Department []WorkDepartment `json:"department"`
	}
	response := &depList{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response.Department, nil
}

// GetAgent 获取自建应用范围
func (w *QyWorkImpl) GetAgent(ctx ctx.Context, req GetAgentRequest) (*GetAgentResponse, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = GetAgent
	req.ReqMethod = "GET"
	response := &GetAgentResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AddTagUsers 添加标签
func (w *QyWorkImpl) AddTagUsers(ctx ctx.Context, req AddTagUsersRequest) error {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return tErr
	}
	req.Token = token
	req.ReqURL = AddTagUsers
	req.ReqMethod = "POST"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	var response QyResponse
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, &response)
	if err != nil {
		return err
	}
	return nil
}

// TagList 标签列表
func (w *QyWorkImpl) SaTagList(ctx ctx.Context, req QyGwRequest) (*SaTagListResp, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = SaTagList
	req.ReqMethod = "GET"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	response := &SaTagListResp{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if tErr != nil {
		return nil, err
	}
	return response, nil
}

// WorkMembers 成员列表
func (w *QyWorkImpl) WorkMembers(ctx ctx.Context, req WorkMemberRequest) ([]MemberDetail, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = MemberList
	req.ReqMethod = "GET"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	type member struct {
		UserList []MemberDetail `json:"userlist"`
	}
	response := &member{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response.UserList, nil
}

// WorkMemberDetail 成员详情
func (w *QyWorkImpl) WorkMemberDetail(ctx ctx.Context, req MemberDetailRequest) (*MemberDetail, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = MemberGet
	req.ReqMethod = "GET"
	if req.ReqPL == "" {
		req.ReqPL = QWExternalIdentifier
	}
	response := &MemberDetail{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response, nil
}

// SendMessage发送消息
func (w *QyWorkImpl) SendMessage(ctx ctx.Context, req SendMessageRequest) (*SendMessageResponse, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = SendMessage
	req.ReqMethod = "POST"
	response := &SendMessageResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response, nil
}

// ExternalDetail 外部联系人详情
func (w *QyWorkImpl) ExternalDetail(ctx ctx.Context, req ExternalContactRequest) (*ExternalContactResponse, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ExternalContactGet
	req.ReqMethod = "GET"
	req.ReqPL = QWExternalIdentifier
	response := &ExternalContactResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response, nil
}

func (w *QyWorkImpl) BatchExternalContactGet(ctx ctx.Context, req BatchExternalContactRequest) (*BatchExternalContactResponse, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = BatchExternalContactGet
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	response := &BatchExternalContactResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response, nil
}

// ExternalGroupChatGet 查询群详情
func (w *QyWorkImpl) ExternalGroupChatGet(ctx ctx.Context, req ExternalGroupRequest) (*ExternalGroupResponse, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ExternalGroupChatGet
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	type groupChat struct {
		GroupChatData ExternalGroupResponse `json:"group_chat"`
	}
	response := &groupChat{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return &response.GroupChatData, nil
}

// UpdateExternalRelationDescription 更新外部联系人描述
func (w *QyWorkImpl) UpdateExternalRelationDescription(c ctx.Context, req EditExternalDescriptionRequest) error {
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return err
	}
	req.Token = token
	req.ReqURL = EditExternalRemarks
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier

	response := struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	}{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return err
	}
	if response.Code == 0 {
		return nil
	}
	return errorx.NewErrorMsg(errs.ExternalEditRemarksErrCode, fmt.Sprintf("[%d]%s", response.Code, response.Msg))
}

// GetWxEexternalContact 获取联系我二维码
func (w *QyWorkImpl) GetWxEexternalContact(ctx ctx.Context, request ContactCreateRequest) (*GetWxCodeResult, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	res := &GetWxCodeResult{}
	req := GetWxCodeRequest{
		Type:  request.Type,
		Scene: 2,
		User:  request.UserIds,
		Party: request.ThirdIds,
		State: strconv.FormatInt(request.ID, 10),
	}
	if request.State != "" {
		req.State = request.State
	}
	if request.AutoFriend == 0 {
		req.SkipVerify = false
	}
	if request.AutoFriend == 1 {
		req.SkipVerify = true
	}
	req.Token = token
	req.ReqURL = AddExternalContact
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	log.Info(ctx).Msgf("获取联系人企业微信请求参数：%+v", req)
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	requestErr := workRequest(ctx, url, req, &res)
	if requestErr != nil {
		return nil, requestErr
	}
	return res, nil
}

// 获取会话同意情况
func (w *QyWorkImpl) CheckSingleAgree(ctx ctx.Context, req CheckSingleAgreeRequest) ([]CheckSingleAgreeResponse, error) {
	// ctx.Secret = "REDACTED_FOR_SOURCE_SHARING"
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = CheckSingleAgree
	req.ReqMethod = "POST"
	req.ReqPL = QWArchiveIdentifier
	type agreeInfo struct {
		Data []CheckSingleAgreeResponse `json:"agreeinfo"`
	}
	response := &agreeInfo{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response.Data, nil
}

// GetCorpTagList 获取企业标签
func (w *QyWorkImpl) GetCorpTagList(c ctx.Context, idList []string) (*GetCorpTagListResponse, error) {
	var req GetCorpTagListRequest
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = GetCorpTagListUrl
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	req.TagID = idList
	var response GetCorpTagListResponse
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, &response)
	if err != nil {
		return nil, err
	}
	if response.Code == 0 {
		return &response, nil
	}
	return nil, errorx.NewErrorMsg(errs.GetCorpTagListErr, fmt.Sprintf("[%d]%s", response.Code, response.Msg))
}

// AddCorpTag 新增企业微信标签
func (w *QyWorkImpl) AddCorpTag(c ctx.Context, req AddCorpTagRequest) (*AddCorpTagResponse, error) {
	// 获取数据请求的token
	token, tokenErr := w.GetToken(c)
	if tokenErr != nil {
		return nil, tokenErr
	}

	req.Token = token
	req.ReqURL = AddCorpTagUrl
	req.ReqMethod = POST
	req.ReqPL = QWExternalIdentifier

	var response AddCorpTagResponse
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	reqErr := workRequest(c, url, req, &response)
	if reqErr != nil {
		return nil, reqErr
	}

	if response.ErrCode == 0 {
		return &response, nil
	}
	return nil, errorx.NewErrorMsg(errs.DelCorpTagErrCode, fmt.Sprintf("[%d]%s", response.ErrCode, response.ErrMsg))
}

// DelCorpTagInfo 删除企业微信标签
func (w *QyWorkImpl) DelCorpTagInfo(c ctx.Context, param DelCorpTagParam) error {
	// 获取数据请求的token
	token, tokenErr := w.GetToken(c)
	if tokenErr != nil {
		return tokenErr
	}

	var req DelCorpTagRequest
	req.Token = token
	req.ReqURL = DelCorpTagUrl
	req.ReqMethod = POST
	req.ReqPL = QWExternalIdentifier

	req.TagID = param.TagID
	req.GroupID = param.GroupID
	var response QyResponse
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	reqErr := workRequest(c, url, req, &response)
	if reqErr != nil {
		return nil
	}

	if response.ErrCode == 0 {
		return nil
	}
	return errorx.NewErrorMsg(errs.DelCorpTagErrCode, fmt.Sprintf("[%d]%s", response.ErrCode, response.ErrMsg))
}

// 更新企业微信标签
func (w *QyWorkImpl) UpdateCorpTagInfo(c ctx.Context, params UpdateTagRequest) error {
	// 获取数据请求的token
	token, tokenErr := w.GetToken(c)
	if tokenErr != nil {
		return tokenErr
	}

	// 包装求请的参数
	params.Token = token
	params.ReqURL = UpdateCorpTagUrl
	params.ReqMethod = POST
	params.ReqPL = QWExternalIdentifier

	var response QyResponse
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	reqErr := workRequest(c, url, params, &response)
	if reqErr != nil {
		return nil
	}

	if response.ErrCode == 0 {
		return nil
	}
	return errorx.NewErrorMsg(errs.DelCorpTagErrCode, fmt.Sprintf("[%d]%s", response.ErrCode, response.ErrMsg))
}

// 修改联系我 请求企业微信
func (w *QyWorkImpl) UpdateWxEexternalContact(ctx ctx.Context, request ContactCreateRequest, configId string) (*GetWxCodeResult, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	res := &GetWxCodeResult{}
	req := UpdateWxCodeRequest{
		User:     request.UserIds,
		Party:    request.ThirdIds,
		ConfigId: configId,
		// State:    strconv.FormatInt(request.ChannelId, 10),
		State: strconv.FormatInt(request.ID, 10),
	}
	if request.AutoFriend == 0 {
		req.SkipVerify = false
	}
	if request.AutoFriend == 1 {
		req.SkipVerify = true
	}
	req.Token = token
	req.ReqURL = UpdateExternalContact
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	requestErr := workRequest(ctx, url, req, &res)
	if requestErr != nil {
		log.Info(ctx).Msgf("编辑企业微信联系我失败%+v", requestErr)
		return nil, requestErr
	}
	return res, nil
}

// GetRobotInfo 获取机器人信息详情
func (w *QyWorkImpl) GetRobotInfo(ctx ctx.Context, req RobotInfoRequest) (*RobotInfo, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = GetRobotInfoUrl
	req.ReqMethod = "GET"
	req.ReqPL = QWArchiveIdentifier
	response := &RobotInfoResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return &response.Data, nil
}

// 删除联系我 请求企业微信
func (w *QyWorkImpl) DeleteWxEexternalContact(ctx ctx.Context, configId string) (*GetWxCodeResult, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	res := &GetWxCodeResult{}
	req := UpdateWxCodeRequest{
		ConfigId: configId,
	}
	req.Token = token
	req.ReqURL = DeleteExternalContact
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	requestErr := workRequest(ctx, url, req, &res)
	if requestErr != nil {
		log.Info(ctx).Msgf("【===删除企业微信联系我失败===】%+v", requestErr)
		return res, requestErr
		// helper.Must(requestErr)
	}

	return res, nil
}

// 获取内部群信息
func (w *QyWorkImpl) InsideGroupChatGet(ctx ctx.Context, req InsideGroupRequest) (*InsideGroupRequestResponse, error) {
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = GetGroupChatInfo
	req.ReqMethod = "POST"
	req.ReqPL = QWArchiveIdentifier
	response := &InsideGroupRequestResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response, nil
}

// CacheMedia 缓存微信图片素材
func (w *QyWorkImpl) CacheMedia(ctx ctx.Context, req CacheMediaWorkRequest) (*CacheMediaWorkResponse, error) {
	// 去BOSS查询应用列表
	reqBoss := boss.WorkPlugin{
		EnterpriseID: req.EnterpriseID,
	}
	plugins, err := boss.WorkPluginList(ctx, reqBoss)

	if err != nil || len(plugins) == 0 {
		log.Info(ctx).Msgf("企业%d从boss下未找到应用列表", req.EnterpriseID)
		return nil, err
	}

	response := CacheMediaWorkResponse{}

	// 模拟POST请求
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	file, err := os.Open(req.FileOrg)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	tmpFile := filepath.Base(req.FileOrg)
	part, err := writer.CreateFormFile("file", tmpFile)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.WriteField("fileType", req.FileType)
	_ = writer.WriteField("filename", req.FileName)
	err = writer.Close()
	if err != nil {
		log.Info(ctx).Msgf("写文件报错, err:%s", err.Error())
		return nil, err
	}
	contentType := writer.FormDataContentType()
	for _, plugin := range plugins {
		url := fmt.Sprintf(QyGwPic, config.CfgData.Service.Work.Domain, w.domain, plugin.Identifier)
		b := &bytes.Buffer{}
		b.Write(body.Bytes())
		request, _ := http.NewRequest("POST", url, b)
		request.Header.Set("Content-Type", contentType)
		request.Header.Set(sdk.HeaderKeyEnum.XRequestID, ctx.ID())
		httpResponse, err := client.Do(request)
		b.Reset()
		if err == nil {
			returnStr, err := ioutil.ReadAll(httpResponse.Body)
			httpResponse.Body.Close()
			if err != nil {
				log.Info(ctx).Msgf("response解析失败, err:%s", err)
				// 继续用下一个应用上传
				continue
			}
			log.Info(ctx).Msgf("上传临时素材, requestUrl:%s,file:%+v, response:%s", url, req, returnStr)
			contentJson := CacheMediaWorkResponse{}
			err = json.Unmarshal([]byte(returnStr), &contentJson)
			if err != nil {
				log.Info(ctx).Msgf("【===请求上传素材库JSON解析失败%s===】", err.Error())
				continue
			}
			// 代表上传失败
			if contentJson.Code != "SUCCESS" {
				continue
			}
			response.Code = contentJson.Code
			response.Data.MediaId = contentJson.Data.MediaId
			break
		} else {
			log.Info(ctx).Msgf("上传素材失败, err:%s, url:%s", err.Error(), url)
		}
	}
	body.Reset()
	if response.Code != "SUCCESS" {
		response.Code = "FAIL"
	}

	return &response, nil
}

// SendWelcome 发送欢迎语
func (w *QyWorkImpl) SendWelcome(c ctx.Context, req SendWelcomeRequest) error {
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return err
	}
	qyGwRequest := QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    SendWelcome,
	}
	req.QyGwRequest = qyGwRequest

	reqj, _ := json.Marshal(req)
	log.Info(c).Msgf("workSendWelcomeData %s", string(reqj))
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// UserTagHandler 用户标签处理
func (w *QyWorkImpl) UserTagHandler(c ctx.Context, request UserTagHandlerRequest) (bool, error) {
	// 获取Token
	token, err := w.GetToken(c)
	if err != nil {
		return false, err
	}
	request.Token = token
	request.ReqURL = UserTagHandlerUrl
	request.ReqMethod = "POST"
	request.ReqPL = QWExternalIdentifier
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	response := struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	}{}
	err = workRequest(c, url, request, &response)
	if err != nil || response.Code != 0 {
		return false, err
	}

	return true, nil
}

// GetGroupMsgResult 群发发送结果
func (w *QyWorkImpl) GetGroupMsgResult(ctx ctx.Context, req MsgResultInfoRequest) ([]MsgResultInfoResponse, error) {
	// 获取Token
	token, tErr := w.GetToken(ctx)
	if tErr != nil {
		return nil, tErr
	}
	req.Token = token
	req.ReqURL = GetGroupMsgResult
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	// 返回结果
	type detailList struct {
		DetailList []MsgResultInfoResponse `json:"detail_list"`
	}
	response := &detailList{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err := workRequest(ctx, url, req, response)
	if err != nil {
		return nil, err
	}
	return response.DetailList, nil
}

// QwMassHair 添加企业群发消息任务
func (w *QyWorkImpl) QwMassHair(c ctx.Context, req QyMassHairRequest) (*QyMassHairResponse, error) {
	// 获取Token
	response := &QyMassHairResponse{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	qyGwRequest := QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    QwMassHair,
	}
	req.Link = LimitLinkForQy(req.Link)
	var actReq interface{}
	if req.Image.MediaId != "" || req.Image.PicUrl != "" {
		actReq = SendMassHairTextImage{
			QyGwRequest:    qyGwRequest,
			Text:           req.Text,
			Image:          req.Image,
			ChatType:       req.ChatType,
			ExternalUserid: req.ExternalUserid,
		}
	}
	if req.Link.Title != "" {
		actReq = SendMassHairTextLink{
			QyGwRequest:    qyGwRequest,
			Text:           req.Text,
			Link:           req.Link,
			ChatType:       req.ChatType,
			ExternalUserid: req.ExternalUserid,
		}
	}
	if req.MiniProgram.Appid != "" {
		actReq = SendMassHairMiniProgram{
			QyGwRequest:    qyGwRequest,
			Text:           req.Text,
			MiniProgram:    req.MiniProgram,
			ChatType:       req.ChatType,
			ExternalUserid: req.ExternalUserid,
		}
	}
	if actReq == nil {
		actReq = SendMassHairText{
			QyGwRequest:    qyGwRequest,
			Text:           req.Text,
			ChatType:       req.ChatType,
			ExternalUserid: req.ExternalUserid,
		}
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, actReq, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// ExternalContactList 获取外部联系人列表
func (w *QyWorkImpl) ExternalContactList(c ctx.Context, request ExternalContactListRequest) ([]string, error) {
	token, err := w.GetToken(c)
	if err != nil {
		return nil, err
	}
	request.Token = token
	request.ReqURL = ExternalContactGetUrl
	request.ReqMethod = "GET"
	request.ReqPL = QWExternalIdentifier
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	var response ExternalContactListResponse
	err = workRequest(c, url, request, &response)
	if err != nil {
		log.Info(c).Msgf("请求企业微信接口失败, err:%s, request:%v", err.Error(), request)
		return nil, err
	}
	if response.ErrCode != 0 {
		log.Info(c).Msgf("企微接口返回错误, response:%+v", response)
		return nil, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
	}
	return response.ExternalUserIDList, nil
}

// ExternalGroupChatList 查询群列表
func (w *QyWorkImpl) ExternalGroupChatList(ctx ctx.Context, req ExternalGroupChatListRequest) ([]*ExternalGroupChatItem, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ExternalGroupChatListURL
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	response := &ExternalGroupChatListResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response.GroupChatList, nil
}

func (w *QyWorkImpl) ExternalGroupChatListUpdate(ctx ctx.Context, req ExternalGroupChatListRequestV2) ([]*ExternalGroupChatItem, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	req.Token = token
	req.ReqURL = ExternalGroupChatListURL
	req.ReqMethod = "POST"
	req.ReqPL = QWExternalIdentifier
	response := &ExternalGroupChatListResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}
	return response.GroupChatList, nil
}

// AuthUserInfo 获取用户身份
func (w *QyWorkImpl) AuthUserInfo(ctx ctx.Context, req UserInfoRequest) (*AuthAllUserInfo, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Token = token
	req.ReqURL = AuthUserInfoUrl
	req.ReqMethod = "GET"
	response := &AuthInfoResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}

	return response.AuthAllUserInfo, nil
}

// JsApiTicket 获取企业的jsapi_ticket
func (w *QyWorkImpl) JsApiTicket(ctx ctx.Context) (string, int, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return "", 0, err
	}

	req := &QyGwRequest{
		Token:     token,
		ReqURL:    JsTicket,
		ReqMethod: "GET",
		ReqPL:     QWExternalIdentifier,
	}

	response := &TicketResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return "", 0, wErr
	}

	return response.Ticket, response.Expires, nil
}

// JsAgentTicket 获取应用的jsapi_ticket
func (w *QyWorkImpl) JsAgentTicket(ctx ctx.Context) (string, int, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return "", 0, err
	}

	req := AgentTicketRequest{
		QyGwRequest{
			Token:     token,
			ReqURL:    AgentTicket,
			ReqMethod: "GET",
			ReqPL:     QWGadgetIdentifier,
		},
		"agent_config",
	}

	response := &TicketResponse{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return "", 0, wErr
	}

	return response.Ticket, response.Expires, nil
}

// QwMassHair 添加企业群发消息任务
func (w *QyWorkImpl) SendMsgTemplate(c ctx.Context, req SendMsgTemplateReq) (*QyMassHairResponse, error) {
	// 获取Token
	response := &QyMassHairResponse{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    QwMassHair,
	}
	for i := range req.Attachments {
		if req.Attachments[i].Link != nil {
			limitedLink := LimitLinkForQy(*req.Attachments[i].Link)
			req.Attachments[i].Link = &limitedLink
		}
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (w *QyWorkImpl) GetGroupMsgSendResult(c ctx.Context, req GroupMsgSendResultRequest) (*GroupMsgSendResultResponse, error) {
	// 获取Token
	response := &GroupMsgSendResultResponse{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    GroupMsgSendResult,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// SetMediaID 设置mediaID
func (w *QyWorkImpl) SetMediaID(c ctx.Context, content string, enterpriseID int64) (string, error) {
	if content == "" {
		return "", nil
	}

	var mediaData MediaContentData
	err := json.Unmarshal([]byte(content), &mediaData)
	if err != nil {
		log.Error(c).Msgf("解析素材文本错误, err:%s", err.Error())
		return content, err
	}

	sendMsgTemplateReq := SendMsgTemplateReq{
		Text: mediaData.Text,
	}
	for _, media := range mediaData.Content {
		if media.Link != nil {
			sendMsgTemplateReq.Attachments = append(sendMsgTemplateReq.Attachments, SendMsgTemplateAttachments{
				MsgType: "link",
				Link:    media.Link,
			})
		}

		if media.Image != nil {
			if media.Image.PicUrl != "" && media.Image.MediaId == "" {
				mediaID, err := w.getMediaIDOfURL(c, media.Image.PicUrl, "1", enterpriseID)
				if err != nil {
					return content, err
				}
				media.Image.MediaId = mediaID
				sendMsgTemplateReq.Attachments = append(sendMsgTemplateReq.Attachments, SendMsgTemplateAttachments{
					MsgType: "image",
					Image:   media.Image,
				})
			}
		}

		if media.MiniProgram != nil {
			if media.MiniProgram.CoverUrl != "" && media.MiniProgram.PicMediaId == "" {
				mediaID, err := w.getMediaIDOfURL(c, media.MiniProgram.CoverUrl, "1", enterpriseID)
				if err != nil {
					return content, err
				}
				media.MiniProgram.PicMediaId = mediaID
				sendMsgTemplateReq.Attachments = append(sendMsgTemplateReq.Attachments, SendMsgTemplateAttachments{
					MsgType:     "miniprogram",
					MiniProgram: media.MiniProgram,
				})
			}
		}

		if media.Video != nil {
			if media.Video.Url != "" && media.Video.MediaId == "" {
				mediaID, err := w.getMediaIDOfURL(c, media.Video.Url, "3", enterpriseID)
				if err != nil {
					return content, err
				}
				media.Video.MediaId = mediaID
				sendMsgTemplateReq.Attachments = append(sendMsgTemplateReq.Attachments, SendMsgTemplateAttachments{
					MsgType: "video",
					Video:   media.Video,
				})
			}
		}

		if media.Files != nil {
			if media.Files.PicUrl != "" && media.Files.MediaId == "" {
				mediaID, err := w.getMediaIDOfURL(c, media.Files.PicUrl, "99", enterpriseID)
				if err != nil {
					return content, err
				}
				media.Files.MediaId = mediaID
				sendMsgTemplateReq.Attachments = append(sendMsgTemplateReq.Attachments, SendMsgTemplateAttachments{
					MsgType: "file",
					File:    media.Files,
				})
			}
		}
	}

	newContent, _ := json.Marshal(sendMsgTemplateReq)

	return string(newContent), nil
}

// GetTmpFileOfURL 获取url指定的内容，生成临时文件
func makeMediaOfURL(imgURL, fileName, fileType string) (*bytes.Buffer, string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(imgURL)
	if err != nil {
		return nil, "", err
	}

	defer res.Body.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, res.Body)
	if err != nil {
		return nil, "", err
	}
	_ = writer.WriteField("fileType", fileType)
	_ = writer.WriteField("filename", fileName)
	contentType := writer.FormDataContentType()

	return body, contentType, writer.Close()
}

// getMediaIDOfURL url文件上传素材库生成mediaID
func (w *QyWorkImpl) getMediaIDOfURL(c ctx.Context, URL string, fileType string, enterpriseID int64) (string, error) {
	ext := path.Ext(URL)
	// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
	if u, err := url.Parse(URL); err == nil && u.Path != "" {
		ext = path.Ext(u.Path)
	}
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext

	body, contentType, err := makeMediaOfURL(URL, fileName, fileType)
	if err != nil {
		return "", err
	}

	response, err := w.CacheMediaStream(c, enterpriseID, URL, body, contentType)
	if err != nil {
		log.Info(c).Msgf("生成mediaID失败, err:%s", err.Error())
		return "", err
	}
	return response.Data.MediaId, nil
}

// CacheMediaStream 缓存微信图片素材
func (w *QyWorkImpl) CacheMediaStream(ctx ctx.Context, enterpriseID int64, cacheKey string, body *bytes.Buffer, contentType string) (*CacheMediaWorkResponse, error) {
	// 去BOSS查询应用列表
	plugins, err := boss.WorkPluginList(ctx, boss.WorkPlugin{
		EnterpriseID: enterpriseID,
	})
	if err != nil || len(plugins) == 0 {
		log.Info(ctx).Msgf("企业%d从boss下未找到应用列表", enterpriseID)
		return nil, err
	}

	response := CacheMediaWorkResponse{}

	// 模拟POST请求
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	for _, plugin := range plugins {
		url := fmt.Sprintf(QyGwPic, config.CfgData.Service.Work.Domain, w.domain, plugin.Identifier)
		b := &bytes.Buffer{}
		b.Write(body.Bytes())
		request, _ := http.NewRequest("POST", url, b)
		request.Header.Set("Content-Type", contentType)
		request.Header.Set(sdk.HeaderKeyEnum.XRequestID, ctx.ID())
		httpResponse, err := client.Do(request)
		if err == nil {
			returnStr, err := ioutil.ReadAll(httpResponse.Body)
			httpResponse.Body.Close()
			if err != nil {
				log.Info(ctx).Msgf("response解析失败, err:%s", err)
				// 继续用下一个应用上传
				continue
			}
			log.Info(ctx).Msgf("上传临时素材, requestUrl:%s, response:%s", url, returnStr)
			contentJson := CacheMediaWorkResponse{}
			err = json.Unmarshal([]byte(returnStr), &contentJson)
			if err != nil {
				log.Info(ctx).Msgf("【===请求上传素材库JSON解析失败%s===】", err.Error())
				continue
			}

			response.Code = contentJson.Code
			response.Data.MediaId = contentJson.Data.MediaId
			response.Msg = contentJson.Msg
			// 代表上传成功
			if contentJson.Code == "SUCCESS" {
				break
			}
		} else {
			log.Info(ctx).Msgf("上传素材失败, err:%s, url:%s", err.Error(), url)
		}
	}
	if response.Code != "SUCCESS" {
		return nil, errors.New(fmt.Sprintf("上传临时素材失败 code:%s,msg:%s", response.Code, response.Msg))
	} else if cacheKey != "" {
		redis.Client().Set(ctx, cacheKey, response.Data.MediaId, 2*24*60*60*time.Second)
	}

	return &response, nil
}

// GetUserBehaviorData o获取「联系客户统计」数据
func (w *QyWorkImpl) GetUserBehaviorData(c ctx.Context, req UserBehaviorDataReq) (*UserBehaviorDataResp, error) {
	// 获取Token
	response := &UserBehaviorDataResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    GetUserBehaviorData,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// Code2Session 获取导购身份信息
func (w *QyWorkImpl) Code2Session(ctx ctx.Context, req Code2SessionRequest) (*Code2SessionReq, error) {
	// 获取Token
	token, err := w.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Token = token
	req.ReqURL = Jscode2session
	req.ReqMethod = "GET"
	response := &Code2SessionReq{}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	wErr := workRequest(ctx, url, req, response)
	if wErr != nil {
		return nil, wErr
	}

	return response, nil
}

// AddMomentTask 企业发表内容到客户的朋友圈
func (w *QyWorkImpl) AddMomentTask(c ctx.Context, req AddMomentTaskReq) (*AddMomentTaskResp, error) {
	// 获取Token
	response := &AddMomentTaskResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    AddMomentTask,
	}
	for i := range req.Attachments {
		if req.Attachments[i].Link != nil {
			limitedLink := LimitMomentLinkForQy(*req.Attachments[i].Link)
			req.Attachments[i].Link = &limitedLink
		}
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetMomentTaskResult 获取朋友圈任务的创建结果
func (w *QyWorkImpl) GetMomentTaskResult(c ctx.Context, req MomentTaskResultRequest) (*MomentTaskResultResp, error) {
	// 获取Token
	response := &MomentTaskResultResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "GET",
		ReqURL:    GetMomentTaskResult,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetMomentTask 获取客户朋友圈企业发表的列表
func (w *QyWorkImpl) GetMomentTask(c ctx.Context, req GetMomentTaskRequest) (*GetMomentTaskResp, error) {
	// 获取Token
	response := &GetMomentTaskResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    GetMomentTask,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetMomentSendResult 获取客户朋友圈发表后的可见客户列表
func (w *QyWorkImpl) GetMomentSendResult(c ctx.Context, req GetMomentSendTaskRequest) (*GetMomentSendTaskResp, error) {
	// 获取Token
	response := &GetMomentSendTaskResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    GetMomentSendResult,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// CancelMomentTask 停止发表企业朋友圈
func (w *QyWorkImpl) CancelMomentTask(c ctx.Context, req CancelMomentTaskRequest) (*QyResponse, error) {
	// 获取Token
	response := &QyResponse{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    CancelMomentTask,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// CancelGroupMsgSend 停止企业群发
func (w *QyWorkImpl) CancelGroupMsgSend(c ctx.Context, req CancelGroupMsgSendRequest) (*QyResponse, error) {
	// 获取Token
	response := &QyResponse{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}
	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    CancelGroupMsgSend,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetMediaId 根据cos地址获取企微的临时素材id
func (w *QyWorkImpl) GetMediaId(c ctx.Context, enterpriseId int64, url string, mediaType string) (string, error) {
	fileType := "1"
	if mediaType == "video" {
		fileType = "3"
	}
	return w.getMediaIDOfURL(c, url, fileType, enterpriseId)
}

// GetMomentComments  获取朋友圈的评论列表
func (w *QyWorkImpl) GetMomentComments(c ctx.Context, req GetMomentCommentsRequest) (*GetMomentCommentResp, error) {
	// 获取Token
	response := &GetMomentCommentResp{}
	token, err := w.GetToken(c)
	if err != nil {
		return response, err
	}

	req.QyGwRequest = QyGwRequest{
		Token:     token,
		ReqPL:     QWExternalIdentifier,
		ReqMethod: "POST",
		ReqURL:    GetMomentComments,
	}
	url := fmt.Sprintf(QyGw, config.CfgData.Service.Work.Domain, w.domain)
	err = workRequest(c, url, req, response)
	if err != nil {
		return response, err
	}
	return response, nil
}
