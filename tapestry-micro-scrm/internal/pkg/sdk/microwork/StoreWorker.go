package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type StoreWorker struct {
}

func NewStoreWorker() *StoreWorker {
	return &StoreWorker{}
}

// AddGuide 添加导购(企业导购,社群导购)
func (s StoreWorker) AddGuide(c ctx.Context, request AddGuideRequest) (AddGuideResponse, error) {
	var response AddGuideResponse
	err := httpClient.Post(c, GetURL("%s/service/v1/store-worker/add-guide"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
