package qyWork

import (
	"encoding/xml"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/pkg/sdk/work"
)

// ExternalDecrypt 外部联系人回调解密
func ExternalDecrypt(req work.QyDecryptRequest) (*ExternalContent, error) {
	wxCrypt := NewWXBizMsgCrypt("xksYN3", "3Qt7eyHzEPTSqYAGOoRspVU0maCf8Z295nDFkJi1NLw", req.CorpID, XmlType)

	msg, cryptErr := wxCrypt.VerifyURL(req.MsgSignature, req.Timestamp, req.Nonce, req.EncryptString)
	if nil != cryptErr {
		return nil, errorx.NewErrorMsg(cryptErr.ErrCode, cryptErr.ErrMsg)
	}

	var content ExternalContent
	err := xml.Unmarshal(msg, &content)
	if nil != err {
		return nil, errors.New("外部联系人回调解密异常")
	}

	return &content, nil
}

// ContactDecrypt 通讯录回调解密
func ContactDecrypt(ctx ctx.Context, req work.QyDecryptRequest) {

}
