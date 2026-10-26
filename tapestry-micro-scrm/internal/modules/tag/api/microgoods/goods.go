package microgoods

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/sdk/microgoods"
)

// GetSpuListBySpuIds 根据spuId获取spu信息
func GetSpuListBySpuIds(c ctx.Context, spuIds []int64, enterpriseID int64) (map[int64]*entity.TagGoodInfo, error) {
	var request = microgoods.SpuListRequest{
		EnterpriseID: enterpriseID,
		SpuIds:       spuIds,
		Page:         1,
		Size:         len(spuIds),
	}
	response, err := microgoods.SpuList(c, &request)
	if err != nil {
		return nil, err
	}
	var ret = make(map[int64]*entity.TagGoodInfo)
	for _, v := range response.Data.List {
		ret[v.ID] = &entity.TagGoodInfo{
			ImgURL: v.SpuSource.URL,
			Name:   v.SpuName,
		}
	}
	return ret, nil
}
