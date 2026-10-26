package customerSvc

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"micro-scrm/internal/pkg/helper"
)

type ExternalCrowd struct {
	externalCrowdDao *externalSearchDao.ExternalCrowd
}

func NewExternalCrowd() *ExternalCrowd {
	return &ExternalCrowd{
		externalCrowdDao: externalSearchDao.NewExternalCrowd(),
	}
}
func (s *ExternalCrowd) Add(c ctx.Context, req customerSearchReq.AddCrowdReq, externalIds []int64) error {
	crowdId, err := s.externalCrowdDao.AddCrowd(c, req)
	if err != nil {
		return err
	}

	err = s.externalCrowdDao.AddCrowdData(c, req.EnterpriseId, crowdId, externalIds)

	return err
}

func (s *ExternalCrowd) List(c ctx.Context, req customerSearchReq.CrowdListReq) ([]customerSearchReq.AddCrowdReq, error) {
	result := make([]customerSearchReq.AddCrowdReq, 0)
	list, err := s.externalCrowdDao.List(c, req)
	if err != nil {
		return result, err
	}
	for _, item := range list {
		var info customerSearchReq.AddCrowdReq
		helper.Unmarshal(item.SearchConnJson, &info.Search)
		info.Sort = item.Sort
		info.MemberId = item.MemberID
		info.Name = item.Name
		info.CrowdId = item.ID
		result = append(result, info)
	}
	return result, nil
}

func (s *ExternalCrowd) ListByConn(c ctx.Context, req customerSearchReq.CrowdInfoReq) ([]externalSearchModel.ExternalCrowdDataModel, error) {
	return s.externalCrowdDao.ListByConn(c, req)
}

func (s *ExternalCrowd) PageListByExternalConn(c ctx.Context, req customerSearchReq.CrowdInfoReq, externalSearchQuery customerSearchResp.SearchQueryResp) ([]externalSearchModel.ExternalCrowdDataModel, int64, error) {
	return s.externalCrowdDao.PageListByExternalConn(c, req, externalSearchQuery)
}

func (s *ExternalCrowd) CheckAuth(c ctx.Context, req customerSearchReq.CrowdCheckReq) error {
	info, err := s.externalCrowdDao.GetOne(c, req.EnterpriseId, "id = ?", req.Id)
	if err != nil {
		return err
	}
	if info.MemberID != req.MemberId {
		return errors.New("系统错误")
	}
	return nil
}

func (s *ExternalCrowd) Edit(c ctx.Context, req customerSearchReq.CrowdEditReq) error {
	exist, err := s.externalCrowdDao.GetOne(c, req.EnterpriseId, "member_id=? AND name=?", req.MemberId, req.Name)
	if err != nil {
		return err
	}
	if exist.ID > 0 && exist.ID != req.CrowdId {
		return errors.New("人群包名称有重复")
	}
	return s.externalCrowdDao.Edit(c, req)
}

func (s *ExternalCrowd) Del(c ctx.Context, req customerSearchReq.CrowdEditReq) error {
	return s.externalCrowdDao.Del(c, req)
}

func (s *ExternalCrowd) UpdateSort(c ctx.Context, req customerSearchReq.CrowdUpdateSortReq) error {
	return s.externalCrowdDao.UpdateSort(c, req)
}

func (s *ExternalCrowd) UpdateContact(c ctx.Context, req customerSearchReq.CrowdUpdateContactReq) error {
	return s.externalCrowdDao.UpdateContact(c, req)
}

func (s *ExternalCrowd) BatchUpdateContact(c ctx.Context, req customerSearchReq.BatchCrowdUpdateContactReq) error {
	return s.externalCrowdDao.BatchUpdateContact(c, req)
}

func (s *ExternalCrowd) PageListByExternalConnSw(c ctx.Context, req customerSearchReq.CrowdInfoSwReq, externalSearchQuery customerSearchResp.SearchQuerySwResp) ([]externalSearchModel.ExternalCrowdDataModel, int64, error) {
	return s.externalCrowdDao.PageListByExternalConnSw(c, req, externalSearchQuery)
}

func (s *ExternalCrowd) ListByConnSw(c ctx.Context, req customerSearchReq.CrowdInfoSwReq) ([]externalSearchModel.ExternalCrowdDataModel, error) {
	return s.externalCrowdDao.ListByConnSw(c, req)
}
