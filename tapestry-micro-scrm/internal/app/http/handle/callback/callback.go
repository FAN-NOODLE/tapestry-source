package callback

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"io"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/pkg/errs"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"net/http"
)

type Callback struct {
	handle.BaseHandle
}

func NewCallbackHandle() Callback {
	return Callback{}
}

func (h Callback) Handle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		h.verify(c)
		return
	}

	if c.Request.Method == http.MethodPost {
		cc := h.WarpContext(c)

		var request qyWork.Request

		// 取企业微信的回调参数，签名.时间戳,随机数
		if err := c.ShouldBindQuery(&request); err != nil {
			h.Log().Error(cc).Msgf("企业微信Query参数回调参数，错误信息：%+v", err)
			h.ResError(c, errs.CallbackQueryErr)
			return
		}
		h.Log().Debug(cc).Msgf("企业微信Query参数回调参数：%v", request)

		// 取出xml里面的企业id
		reqBody, _ := io.ReadAll(c.Request.Body)
		h.Log().Debug(cc).Msgf("企微回调,请求参数%s", string(reqBody))

		if err := xml.Unmarshal(reqBody, &request); err != nil {
			h.Log().Error(cc).Msgf("企业微信获取到的回调参数：%+v, 错误信息：%+v", request, err)
			h.ResError(c, errs.CallbackXmlErr)
			return
		}
		h.Log().Debug(cc).Msgf("企微回调,回调参数：%v", request)

		token, encodingAesKey, err := getCallbackVerify(cc, request.ToUserName)
		if err != nil {
			h.ResError(c, err)
			return
		}

		var wxError *qyWork.CryptError
		// 解密企微加密数据
		wxCpt := qyWork.NewWXBizMsgCrypt(token, encodingAesKey, request.ToUserName, qyWork.XmlType)
		resultByte, wxError := wxCpt.DecryptMsg(request.MsgSignature, request.Timestamp, request.Nonce, reqBody)
		if wxError != nil {
			h.Log().Error(cc).Msgf("企微回调,解密异常：%v", wxError)
			h.ResError(c, errorx.NewErrorMsg(wxError.ErrCode, wxError.ErrMsg))
			return
		}

		// 解析基础字段
		var baseEventResponse qyWork.DecryptBaseData
		if err = json.Unmarshal(resultByte, &baseEventResponse); err != nil {
			h.Log().Error(cc).Msgf("企微回调,解析事件,错误信息：%+v", err)
			h.ResError(c, errs.ParseResultErr)
			return
		}

		h.Log().Debug(cc).Msgf("企微回调,解析事件-数据：%v", baseEventResponse)

		// 针对不同类型时间分发处理
		switch baseEventResponse.Event {
		case qyWork.ChangeContactEvent:
		case qyWork.ChangeExternalContactEvent:
		case qyWork.ChangeExternalChatEvent:
		case qyWork.ChangeExternalTagEvent:
		}
	}

	return
}

// 回调验证
func (h Callback) verify(c *gin.Context) {
	cc := h.WarpContext(c)

	h.Log().Debug(cc).Msgf("微信回调验证来啦....")

	var callbackMsg qyWork.Verify

	// 取企业微信的回调参数，签名.时间戳
	if err := c.ShouldBindQuery(&callbackMsg); err != nil {
		h.Log().Error(cc).Msgf("微信回调验证,Query参数，错误信息：%+v", err)
		h.ResError(c, err)
		return
	}

	h.Log().Debug(cc).Msgf("微信回调验证,Query参数: %+v", callbackMsg)

	token, encodingAesKey, err := getCallbackVerify(cc, callbackMsg.CorpID)
	if err != nil {
		h.ResError(c, err)
		return
	}

	wxCpt := qyWork.NewWXBizMsgCrypt(token, encodingAesKey, callbackMsg.CorpID, qyWork.XmlType)

	echoStrByte, cryptErr := wxCpt.VerifyURL(callbackMsg.MsgSignature, callbackMsg.Timestamp, callbackMsg.Nonce, callbackMsg.EchoStr)
	if nil != cryptErr {
		h.Log().Error(cc).Msgf("微信回调验证失败: %+v", cryptErr)
		c.String(http.StatusOK, "微信回调验证失败: %+v", cryptErr)
		return
	}

	h.Log().Debug(cc).Msgf("微信回调验证,解密值: %+v", string(echoStrByte))

	c.String(http.StatusOK, string(echoStrByte))

	return
}

func getCallbackVerify(cc ctx.Context, corpId string) (string, string, error) {
	// 从缓存中取出 验证参数
	cacheKey := fmt.Sprintf("callback:verify:%s", corpId)
	cacheResult, err := app.Redis().Get(cc, cacheKey).Result()
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("微信[%s]回调: 验证参数不存在", corpId))
	}

	var result struct {
		Token          string `json:"token"`
		EncodingAESKey string `json:"encoding_aes_key"`
	}
	err = json.Unmarshal([]byte(cacheResult), &result)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("微信[%s]回调,解密缓存消息失败:%v", corpId, err))
	}

	return result.Token, result.EncodingAESKey, nil
}
