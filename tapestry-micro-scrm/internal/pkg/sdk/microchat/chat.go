package microchat

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/iox"
	"micro-scrm/internal/pkg/sdk"
)

type Service interface {
	GetChatMsg(c ctx.Context, request *GetChatMsgReq) ([]*PlainMsg, error)
	GetMediaData(c ctx.Context, request *GetMediaDataReq, iox iox.WriterCloser, size int) error
}

type service struct {
	client *sdk.Client
	host   string
}

func NewService(host string) Service {
	client := sdk.NewHTTPClient(sdk.ParseJSON{})
	s := service{
		client,
		host,
	}

	return s
}

func (s service) GetChatMsg(c ctx.Context, request *GetChatMsgReq) ([]*PlainMsg, error) {
	url := fmt.Sprintf("%s/micro-chat-scrm/v1/chat/msg", s.host)
	type response struct {
		sdk.ResponseBase
		Data []*PlainMsg
	}

	res := new(response)
	err := s.client.Post(c, url, request, res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (s service) GetMediaData(c ctx.Context, request *GetMediaDataReq, iox iox.WriterCloser, size int) error {
	url := fmt.Sprintf("%s/micro-chat-scrm/v1/chat/store-media", s.host)
	res := new(sdk.ResponseBase)
	err := s.client.Post(c, url, request, res)
	if err != nil {
		return err
	}

	if res.Code != 200 {
		return errorx.NewErrorMsg(res.Code, res.Msg)
	}
	return nil
}
