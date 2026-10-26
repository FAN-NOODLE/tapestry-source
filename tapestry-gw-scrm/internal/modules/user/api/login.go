package api

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gw-scrm/internal/modules/user"
	"gw-scrm/internal/modules/user/entity"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microauthority"
)

type WorkWechatImpl struct{}
type UcPostReq struct{}
type AuthCodeStruct struct {
	WxAuthCode string `json:"wxAuthCode"`
}

// boss企业信息
type BossEnterpriseDetail struct {
	EnterpriseId string `json:"enterprise_id"`
	Hash         string `json:"hash"`
	CorpId       string `json:"corp_id"`
	Domain       string `json:"domain"`
}

// 调用UC的接口查询企业信息
func AuthCodeInfo(authCode string) (interface{}, error) {
	clientMode := &common.ClientMode{}
	ucClientHost := clientMode.GetUcUrl()
	url := fmt.Sprintf("%s/s/qy/login/wx-qrcode-login", ucClientHost)
	sendParams := &AuthCodeStruct{WxAuthCode: authCode}
	res, err := client.Post(ctx.New(), url, clients.JSONBody(sendParams))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		ucResponse := &user.UcResponse{}
		json.Unmarshal([]byte(rt), ucResponse)
		err := errorx.New(ucResponse.Msg)
		return rt, err
	})
}

// 根据epId查询企业信息
func EnSettingQuery(epId string) (interface{}, error) {
	clientMode := &common.ClientMode{}
	ucClientHost := clientMode.GetUcUrl()

	ucPostReq := &UcPostReq{}
	url := fmt.Sprintf("%s/b/%s/ensetting/query", ucClientHost, epId)
	res, err := client.Post(ctx.New(), url, clients.JSONBody(ucPostReq))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		ucResponse := &user.UcResponse{}
		json.Unmarshal([]byte(rt), ucResponse)
		err := errorx.New(ucResponse.Msg)
		return rt, err
	})
}

// 从boss查询企业信息
func EnterpriseDetail(corpId string) (interface{}, error) {
	clientMode := &common.ClientMode{}
	sdkHost := clientMode.GetClientUrl(common.GetClientUrl{"boss"})

	bossEnterpriseDetail := &BossEnterpriseDetail{}
	bossEnterpriseDetail.CorpId = corpId
	url := fmt.Sprintf("%s/service/v1/enterprise/detail", sdkHost)
	res, err := client.Post(ctx.New(), url, clients.JSONBody(bossEnterpriseDetail))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		ucResponse := &user.UcResponse{}
		json.Unmarshal([]byte(rt), ucResponse)
		err := errorx.New(ucResponse.Msg)
		return rt, err
	})
}

func AdminDetail(domain string, eaId string) (interface{}, error) {
	clientMode := &common.ClientMode{}
	ucClientHost := clientMode.GetUcUrl()
	url := fmt.Sprintf("%s/b/%s/adminer/member/detail", ucClientHost, domain)

	ucPostReq := &user.UcAdminDetail{EaId: eaId}
	res, err := client.Post(ctx.New(), url, clients.JSONBody(ucPostReq))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		ucResponse := &user.UcResponse{}
		json.Unmarshal([]byte(rt), ucResponse)
		err := errorx.New(ucResponse.Msg)
		return rt, err
	})
}

// 查询管理员信息
func BossAdminDetail(req *entity.BossEnterpriseAdminDetailRequest) (interface{}, error) {
	clientMode := &common.ClientMode{}
	sdkHost := clientMode.GetClientUrl(common.GetClientUrl{"boss"})

	url := fmt.Sprintf("%s/service/v1/enterprise_admin/detail", sdkHost)
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

func AuthorityCheckSmsCode(req *entity.AuthorityCheckSmsCodeRequest) (interface{}, error) {
	clientMode := &common.ClientMode{}
	sdkHost := clientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := fmt.Sprintf("%s/micro-authority/v1/common/check-sms-code", sdkHost)

	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 保存管理员信息
func BossAdminSave(req *entity.BossEnterpriseAdminSaveRequest) (interface{}, error) {
	clientMode := &common.ClientMode{}
	sdkHost := clientMode.GetClientUrl(common.GetClientUrl{"boss"})
	url := fmt.Sprintf("%s/service/v1/enterprise-admin/save", sdkHost)

	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

func AccountBaseInfo(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microauthority.CheckAccount(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EnterpriseList(c ctx.Context, req interface{}) (interface{}, error) {
	res, err := microauthority.EnterpriseList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
