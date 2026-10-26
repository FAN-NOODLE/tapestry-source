package smshttp

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/marketing"
	marketingMysql "micro-scrm/internal/modules/marketing/respository/mysql"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/modules/sms/respository/mysql"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {

	smsService = sms.InsideNew(mysql.New(), marketing.New(marketingMysql.New(), userEs.NewEsRepo()))
	return []*swagger.Endpoint{
		endpoint.Post("/sms-template/create", "创建短信模板",
			endpoint.Handler(createTemplate),
			endpoint.Tags("创建短信模板"),
			endpoint.Response(http.StatusOK, "",
				"createTemplate", "接口调用成功")),

		endpoint.Post("/sms-template/create-sign", "创建短信模板",
			endpoint.Handler(createSign),
			endpoint.Tags("创建短信签名"),
			endpoint.Response(http.StatusOK, "",
				"createSign", "接口调用成功")),

		endpoint.Post("/sms-template/arraignment", "短信模板提审",
			endpoint.Handler(arraignment),
			endpoint.Tags("短信模板提审"),
			endpoint.Response(http.StatusOK, "",
				"arraignment", "接口调用成功")),

		endpoint.Post("/sms-template/params-list", "短信模板参数",
			endpoint.Handler(paramsList),
			endpoint.Tags("短信模板参数"),
			endpoint.Response(http.StatusOK, "",
				"paramsList", "接口调用成功")),

		endpoint.Post("/sms-template/sign-list", "签名列表接口",
			endpoint.Handler(templateSignList),
			endpoint.Tags("签名列表接口"),
			endpoint.Response(http.StatusOK, "",
				"templateSignList", "接口调用成功")),

		endpoint.Post("/sms-template/template-list", "短信模板列表",
			endpoint.Handler(templateList),
			endpoint.Tags("短信模板列表接口"),
			endpoint.Response(http.StatusOK, "",
				"templateList", "接口调用成功")),

		endpoint.Post("/sms-template/delete", "短信模板删除",
			endpoint.Handler(templateDel),
			endpoint.Tags("短信模板删除"),
			endpoint.Response(http.StatusOK, "",
				"templateList", "接口调用成功")),

		endpoint.Post("/sms-template/edit", "短信模板编辑",
			endpoint.Handler(templateEdit),
			endpoint.Tags("短信模板编辑"),
			endpoint.Response(http.StatusOK, "",
				"templateList", "接口调用成功")),

		endpoint.Post("/sms-template/info", "短信模板详情",
			endpoint.Handler(templateInfo),
			endpoint.Tags("短信模板编辑"),
			endpoint.Response(http.StatusOK, "",
				"templateList", "接口调用成功")),

		endpoint.Post("/sms-template/sign-delete", "标签模板删除",
			endpoint.Handler(templateSignDelete),
			endpoint.Tags("标签模板删除"),
			endpoint.Response(http.StatusOK, "",
				"templateList", "接口调用成功")),

		endpoint.Post("/sms-template/syn-sms-status", "同步短信模板状态",
			endpoint.Handler(synSmsStatus),
			endpoint.Tags("同步短信模板状态"),
			endpoint.Response(http.StatusOK, "",
				"synSmsStatus", "接口调用成功")),

		endpoint.Post("/sms-template/send", "发送短信",
			endpoint.Handler(send),
			endpoint.Tags("发送短信"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功")),

		endpoint.Post("/sms-template/total", "企业短信数据统计",
			endpoint.Handler(total),
			endpoint.Tags("企业短信数据统计"),
			endpoint.Response(http.StatusOK, "",
				"send", "接口调用成功")),
	}
}
