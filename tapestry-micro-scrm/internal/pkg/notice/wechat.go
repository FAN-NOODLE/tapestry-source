package notice

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
)

type WeChat struct {
}

func NewWechat() *WeChat {
	return &WeChat{}
}

func (w WeChat) SendMarkdownMessage(c ctx.Context, content string) {
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}
	endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + app.Cfg().QYNoticeKey
	_, _ = app.Client().Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
}

func (w WeChat) SendTextMessage(c ctx.Context, content string) {
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content":        content,
			"mentioned_list": []string{"@all"},
		},
	}
	endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + app.Cfg().QYNoticeKey
	_, _ = app.Client().Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
}

func (w WeChat) SendTaskMarkdownMessage(c ctx.Context, content string) {
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}
	endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + app.Cfg().NoticeKeys.Task
	_, _ = app.Client().Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
}

func (w WeChat) SendTaskTextMessage(c ctx.Context, content string) {
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content":        content,
			"mentioned_list": []string{"@all"},
		},
	}
	endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + app.Cfg().NoticeKeys.Task
	_, _ = app.Client().Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
}

func (w WeChat) SendMarkMarkdownMessage(c ctx.Context, content string) {
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}
	endPoint := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + app.Cfg().NoticeKeys.Mark
	_, _ = app.Client().Post(ctx.Wrap(c), endPoint, clients.JSONBody(message))
}
