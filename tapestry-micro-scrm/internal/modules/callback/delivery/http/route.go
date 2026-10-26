package callback_http

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/callback"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/fission"
	fissionBossApi "micro-scrm/internal/modules/fission/api/bossapi"
	fissionMysql "micro-scrm/internal/modules/fission/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/scripts"
	scriptsMysql "micro-scrm/internal/modules/scripts/respository/mysql"
	userModule "micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	workExternal "micro-scrm/internal/modules/work-external"
	workExternalMysql "micro-scrm/internal/modules/work-external/respository/mysql"
)

// Endpoints 路由
func Endpoints() []*swagger.Endpoint {
	service = new(callback.Service)
	contactService = drainage.NewContact(mysql.NewDrainage())
	enterpriseService = enterprise.New(enterpriseMysql.New())
	memberService = members.New(membersMysql.New())
	customerService = customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	messageService = customer.NewDelMessage(customerMysql.NewDelMessage())
	drainageService = drainage.NewDrainage(mysql.NewDrainage())
	externalTagService = customer.NewExternalTag(customerMysql.NewCustomer())

	userService = userModule.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())

	externalTrace := workExternal.NewExternalTrace(workExternalMysql.NewRepo())
	service.ExternalService = externalTrace
	scriptService = scripts.New(scriptsMysql.New())
	fissionService = fission.NewService(fissionMysql.New(), contactService, userService, fissionBossApi.NewApi(), enterpriseService)
	return []*swagger.Endpoint{

		endpoint.Post("/callback/enterprise-add", "企业新增回调",
			endpoint.Handler(enterpriseAdd),
			endpoint.Tags("回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/tag-init", "初始化标签结果回调",
			endpoint.Handler(tagInit),
			endpoint.Tags("回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/member-relation", "外部联系人回调",
			endpoint.Handler(externalCallback),
			endpoint.Tags("企业微信外部联系人回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/member-relation-other", "多余的外部联系人回调",
			endpoint.Handler(func(c *gin.Context) {
				ginplus.ResOK(c)
			}),
			endpoint.Tags("多余的企业微信外部联系人回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback-other", "接口调用成功")),

		endpoint.Get("/callback/member-relation-other", "多余的外部联系人回调企微验证",
			endpoint.Handler(func(c *gin.Context) {
				ginplus.ResOK(c)
			}),
		),

		endpoint.Get("/callback/member-relation-manual", "手动处理外部联系人回调接口",
			endpoint.Handler(externalManual),
			endpoint.Tags("企业微信外部联系人回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/user-address", "用户地址数据同步",
			endpoint.Handler(userAddressSync),
			endpoint.Tags("ddc用户地址数据同步回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/user-event", "用户事件同步",
			endpoint.Handler(userEvent),
			endpoint.Tags("用户事件同步"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "callback", "接口调用成功")),

		endpoint.Post("/callback/enterprise-official-account", "企业绑定公众号通知",
			endpoint.Handler(enterpriseOfficialAccount),
			endpoint.Tags("回调"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: callback.Callback{},
			}, "testUserTagHandler", "接口调用成功")),
	}
}
