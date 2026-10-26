package callback

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"io"
	"micro-scrm/internal/app/http/handle"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	callbackSvc "micro-scrm/internal/app/service/callback"
	"micro-scrm/internal/modules/callback"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
)

type Contacts struct {
	handle.BaseHandle
	contactsUserSvc *callbackSvc.User
	contactsDepSvc  *callbackSvc.Department
	contactsTagSvc  *callbackSvc.Tag
}

func NewContactsHandle() Contacts {
	return Contacts{
		contactsUserSvc: callbackSvc.NewUser(),
		contactsDepSvc:  callbackSvc.NewDepartment(),
		contactsTagSvc:  callbackSvc.NewTag(),
	}
}

func (h Contacts) AutoHandle(c *gin.Context) {
	cc := h.WarpContext(c)

	// 取出xml里面的appid
	reqBody, _ := io.ReadAll(c.Request.Body)
	h.Log().Info(cc).Msgf("通讯录回调,回调信息：%s", string(reqBody))

	var callbackMsg callback.QwExternalCallback
	if err := xml.Unmarshal(reqBody, &callbackMsg); err != nil {
		h.Log().Error(cc).Msgf("通讯录回调,企业微信获取到的回调参数：%+v, 错误信息：%+v", callbackMsg, err)
		h.ResError(c, errs.CallbackXmlErr)
		return
	}
	h.Log().Info(cc).Msgf("通讯录回调,企业微信获取到的回调参数：%v", callbackMsg)

	// 取企业微信的回调参数，签名.时间戳
	if err := c.BindQuery(&callbackMsg); err != nil {
		h.Log().Error(cc).Msgf("通讯录回调,企业微信Query参数回调参数，错误信息：%+v", err)
		h.ResError(c, errs.CallbackQueryErr)
		return
	}
	h.Log().Info(cc).Msgf("通讯录回调,企业微信Query参数回调参数：%v", callbackMsg)

	// 解密回调参数
	decryptReq := work.QyDecryptRequest{
		MsgSignature:  callbackMsg.MsgSignature,
		Nonce:         callbackMsg.Nonce,
		Timestamp:     callbackMsg.Timestamp,
		CorpID:        callbackMsg.ToUserName,
		EncryptString: string(reqBody),
	}

	workImpl := new(work.QyWorkImpl)
	callbackTextData, err := workImpl.QyContactDecrypt(cc, decryptReq)
	if err != nil {
		h.Log().Error(cc).Msgf("通讯录回调,解密企业微信回调信息报错:%+v", err)
		h.ResError(c, err)
		return
	}

	// 解析基础字段
	var baseEventResponse struct {
		Data work.ContactCallback `json:"data"`
	}
	if err = json.Unmarshal([]byte(callbackTextData), &baseEventResponse); err != nil {
		h.Log().Error(cc).Msgf("通讯录回调,解析事件,错误信息：%+v", err)
		h.ResError(c, errs.ParseResultErr)
		return
	}
	h.Log().Info(cc).Msgf("通讯录回调,解析事件-数据：%v", baseEventResponse)

	// 查询企业信息
	bossRes, err := orm.BossService.Enterprise.GetInfo(boss.EnterpriseModel{
		CorpID: callbackMsg.ToUserName,
	})

	switch baseEventResponse.Data.ChangeType {
	// 员工事件
	case callbackSvc.CreateUser, callbackSvc.UpdateUser, callbackSvc.DeleteUser:
		var userEveRes struct {
			Data work.ContactCallbackUser `json:"data"`
		}
		if err = json.Unmarshal([]byte(callbackTextData), &userEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析员工事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsUserSvc.Handel(cc, userEveRes.Data, microBossModel.EnterpriseModel(bossRes))
		break
	// 部门事件
	case callbackSvc.CreateParty, callbackSvc.UpdateParty, callbackSvc.DeleteParty:
		var depEveRes struct {
			Data work.ContactCallbackDepartment `json:"data"`
		}
		if err = json.Unmarshal([]byte(callbackTextData), &depEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析部门事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsDepSvc.Handel(cc, depEveRes.Data, bossRes)
		break
	// 标签事件
	case callbackSvc.UpdateTag:
		var tagEveRes struct {
			Data work.ContactCallbackTag `json:"data"`
		}
		if err = json.Unmarshal([]byte(callbackTextData), &tagEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析标签事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsTagSvc.Handel(cc, tagEveRes.Data, bossRes)
		break
	}

	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
	return
}

func (h Contacts) ManualHandle(c *gin.Context) {
	cc := h.WarpContext(c)
	callbackTextData, _ := io.ReadAll(c.Request.Body)

	// 解析基础字段
	var baseEventResponse struct {
		Data work.ContactCallback `json:"data"`
	}

	if err := json.Unmarshal(callbackTextData, &baseEventResponse); err != nil {
		h.ResError(c, errs.ParseResultErr)
		return
	}

	h.Log().Info(cc).Msgf("通讯录回调,解析事件-数据：%v", baseEventResponse)

	// 查询企业信息
	bossRes, err := orm.BossService.Enterprise.GetInfo(boss.EnterpriseModel{
		CorpID: baseEventResponse.Data.ToUserName,
	})

	switch baseEventResponse.Data.ChangeType {
	// 员工事件
	case callbackSvc.CreateUser, callbackSvc.UpdateUser, callbackSvc.DeleteUser:
		var userEveRes struct {
			Data work.ContactCallbackUser `json:"data"`
		}
		if err = json.Unmarshal(callbackTextData, &userEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析员工事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsUserSvc.Handel(cc, userEveRes.Data, microBossModel.EnterpriseModel(bossRes))
		break
	// 部门事件
	case callbackSvc.CreateParty, callbackSvc.UpdateParty, callbackSvc.DeleteParty:
		var depEveRes struct {
			Data work.ContactCallbackDepartment `json:"data"`
		}
		if err = json.Unmarshal(callbackTextData, &depEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析部门事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsDepSvc.Handel(cc, depEveRes.Data, bossRes)
		break
	// 标签事件
	case callbackSvc.UpdateTag:
		var tagEveRes struct {
			Data work.ContactCallbackTag `json:"data"`
		}
		if err = json.Unmarshal(callbackTextData, &tagEveRes); err != nil {
			h.Log().Error(cc).Msgf("通讯录回调,解析标签事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}
		err = h.contactsTagSvc.Handel(cc, tagEveRes.Data, bossRes)
		break
	}

	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
	return
}
