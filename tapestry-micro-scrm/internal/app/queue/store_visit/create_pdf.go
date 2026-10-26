package store_visit

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	storeVisitReq "micro-scrm/internal/app/http/request/store_visit"
	storeVisitSvc "micro-scrm/internal/app/service/store_visit"
)

type CreateStoreVisitPDFProduce struct {
	groupId string
	topics  []string
}

func NewCreateStoreVisitPDFProduce() *CreateStoreVisitPDFProduce {
	return &CreateStoreVisitPDFProduce{
		groupId: "store_visit",
		topics:  []string{"create_store_visit_pdf"},
	}
}

func (q *CreateStoreVisitPDFProduce) Topics() []string {
	return q.topics
}

func (q *CreateStoreVisitPDFProduce) GroupID() string {
	return q.groupId
}

func (q *CreateStoreVisitPDFProduce) Run(c ctx.Context, key string, value *string) error {
	// 处理调整人群包数据
	var params storeVisitReq.CreatePdfReq
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		app.Logger().Error(c).Msgf("%s计划人群人群包导入的数据解析失败，错误原因为:%v", key, err)
		return err
	}
	err = storeVisitSvc.NewCreatePDFQueueSvc().CreatePDFFile(c, params)
	if err != nil {
		return err
	}
	return err
}
