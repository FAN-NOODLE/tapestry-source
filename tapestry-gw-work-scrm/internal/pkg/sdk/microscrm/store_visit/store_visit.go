package storeVisitSdk

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type StoreVisit struct {
	microScrm.BaseSDK
}

func NewStoreVisit() *StoreVisit {
	return &StoreVisit{}
}

func (a *StoreVisit) SavePlan(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-plan"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) List(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) GetStoreManager(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/get-store-manager"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) GetStoreMM(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/get-store-mm"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) DeptMemberList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/dept-member-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
func (a *StoreVisit) EmailContact(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/email-contact"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveEmail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-email"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
func (a *StoreVisit) DelEmail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/del-email"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SearchMember(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/search-member"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
func (a *StoreVisit) AddCopyMember(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/add-copy-member"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) LoginUserInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/login-user-info"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) StartVisit(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/start-visit"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveShowRecordContent(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-show-record-content"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveShowRecordMedia(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-show-record-media"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) DelShowRecord(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/del-show-record"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) DelShowRecordMedia(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/del-show-record-media"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) ShowRecordList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/show-record-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) StoreVisitInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/info"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) CompleteVisit(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/complete-visit"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SendCompleteReport(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/send-complete-report"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveVisitInfo(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-visit-info"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveReportPermission(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-report-permission"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveComment(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-comment"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) CommentList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/comment-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) DelComment(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/del-comment"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) DashboardTree(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/dashboard-tree"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) ExportDashboardList(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/export-dashboard-list"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SendPdfToEmail(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/send-pdf-to-email"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) GetConfig(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/get-config"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

func (a *StoreVisit) SaveLanguageConfig(c ctx.Context, request interface{}) (interface{}, error) {
	res, err := a.Post(c, a.GetUrl("%s/micro-scrm/v1/store-visit/save-language-config"), request)

	if err != nil {
		return nil, err
	}

	result := &microScrm.Response{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != microScrm.StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
