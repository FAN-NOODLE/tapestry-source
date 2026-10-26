package sms_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/sms-template/list", "短信模板列表",
			endpoint.Handler(templateList),
			endpoint.Tags("短信模板列表"),
			endpoint.Response(http.StatusOK, "", "templateList", "接口调用成功")),

		endpoint.Post("/sms-template/create-sign", "添加短信签名",
			endpoint.Handler(createSign),
			endpoint.Tags("短信模板列表"),
			endpoint.Response(http.StatusOK, "", "createSign", "接口调用成功")),

		endpoint.Post("/sms-template/create", "添加短信模板",
			endpoint.Handler(createTemplate),
			endpoint.Tags("添加短信模板"),
			endpoint.Response(http.StatusOK, "", "createTemplate", "接口调用成功")),

		endpoint.Post("/sms-template/arraignment", "短信模板提审",
			endpoint.Handler(arraignment),
			endpoint.Tags("短信模板提审"),
			endpoint.Response(http.StatusOK, "", "arraignment", "接口调用成功")),

		endpoint.Post("/sms-template/params-list", "短信模板提审",
			endpoint.Handler(paramsList),
			endpoint.Tags("短信模板提审"),
			endpoint.Response(http.StatusOK, "", "paramsList", "接口调用成功")),

		endpoint.Post("/sms-template/params-list", "短信参数列表",
			endpoint.Handler(paramsList),
			endpoint.Tags("短信参数列表"),
			endpoint.Response(http.StatusOK, "", "paramsList", "接口调用成功")),

		endpoint.Post("/sms-template/sign-list", "短信签名列表",
			endpoint.Handler(signList),
			endpoint.Tags("短信签名列表"),
			endpoint.Response(http.StatusOK, "", "signList", "接口调用成功")),

		endpoint.Post("/sms-template/delete", "删除短信模板",
			endpoint.Handler(templateDelete),
			endpoint.Tags("删除短信模板"),
			endpoint.Response(http.StatusOK, "", "templateDelete", "接口调用成功")),

		endpoint.Post("/sms-template/edit", "编辑短信模板",
			endpoint.Handler(templateEdit),
			endpoint.Tags("编辑短信模板"),
			endpoint.Response(http.StatusOK, "", "templateEdit", "接口调用成功")),

		endpoint.Post("/sms-template/info", "短信模板详情",
			endpoint.Handler(templateInfo),
			endpoint.Tags("短信模板详情"),
			endpoint.Response(http.StatusOK, "", "templateInfo", "接口调用成功")),

		endpoint.Post("/sms-template/sign-delete", "删除短信签名",
			endpoint.Handler(templateSignDelete),
			endpoint.Tags("删除短信签名"),
			endpoint.Response(http.StatusOK, "", "templateSignDelete", "接口调用成功")),

		endpoint.Post("/sms-template/send", "发送短信",
			endpoint.Handler(send),
			endpoint.Tags("删除短信签名"),
			endpoint.Response(http.StatusOK, "", "send", "接口调用成功")),

		endpoint.Post("/sms-template/total", "企业短信统计",
			endpoint.Handler(total),
			endpoint.Tags("企业短信统计"),
			endpoint.Response(http.StatusOK, "", "total", "接口调用成功")),
	}
}
