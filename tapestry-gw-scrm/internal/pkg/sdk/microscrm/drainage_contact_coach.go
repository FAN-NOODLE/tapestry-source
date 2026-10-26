package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type DrainageContactCoachApi struct {
}

func NewDrainageContactCoachApi() DrainageContactCoachApi {
	return DrainageContactCoachApi{}
}

func (a DrainageContactCoachApi) Save(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/save"), request)
}

func (a DrainageContactCoachApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/detail"), request)
}

func (a DrainageContactCoachApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/list"), request)
}

func (a DrainageContactCoachApi) DetailList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/detail-list"), request)
}

func (a DrainageContactCoachApi) Delete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/delete"), request)
}

func (a DrainageContactCoachApi) ZipUrl(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/drainage-contact-coach/zip-url"), request)
}
