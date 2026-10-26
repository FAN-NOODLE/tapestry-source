package microAuthority

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gw-work-scrm/internal/pkg/sdk"
	"time"
)

func GetAccountDataRangeList(ctx ctx.Context, req AccountDataRangeReq) (res interface{}, err error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL("%s/micro-authority/v1/account/data-range-list"), clients.JSONBody(req))
	if err != nil {
		return res, err
	}
	var response = AccountDataRangeRes{}
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return res, rErr
	}
	
	if response.Code != 200 {
		return res, errorx.NewErrorMsg(response.Code, response.Message)
	}

	return response.Data, nil
}

func GetAccountDataRangeThirdIds(ctx ctx.Context, req AccountDataRangeReq) ([]int, error) {
	dataRangeListRes, err := GetAccountDataRangeList(ctx, req)

	if err != nil {
		return nil, err
	}

	var thirdIDs []int
	dataRangeList, err := json.Marshal(dataRangeListRes)
	var dataRangeLists []AccountDataRangeInfo
	err = json.Unmarshal(dataRangeList, &dataRangeLists)
	if err != nil {
		return nil, err
	}

	for _, v := range dataRangeLists {
		thirdIDs = append(thirdIDs, v.ObjectId)
	}

	return thirdIDs, nil
}
