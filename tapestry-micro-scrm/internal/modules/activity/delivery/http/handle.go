package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/activity"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/modules/marketingPromotion"
	"micro-scrm/internal/pkg/sdk/boss"
	"strconv"
)

var activityService *activity.Service

func syncActivity(c *gin.Context) {
	var request activity.SyncRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = activityService.SyncActivityData(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// goodsList 活动商品列表
func goodsList(c *gin.Context) {
	var request activity.GoodsListRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	cx := ctx.Wrap(c)
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}

	// 查询当前生效的活动
	reqActivity := new(activity.ShopActivityEntity)
	reqActivity.EnterpriseID = request.EnterpriseID
	reqActivity.ActivityStatus = activity.NormalStatus
	reqActivity.ActivityType = request.ActivityType
	activities, err := activityService.AllActivity(cx, reqActivity)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if len(activities) == 0 {
		ginplus.ResSuccess(c, activity.GoodsListResponse{
			List:  []*activity.GoodsListItem{},
			Page:  request.Page,
			Size:  request.Size,
			Total: 0,
		})
		return
	}

	// 提取活动ID
	activityIds := make([]int64, 0, len(activities))
	activityInfo := make(map[int64]*activity.ShopActivityEntity)
	for _, activity := range activities {
		activityIds = append(activityIds, activity.ID)
		activityInfo[activity.ID] = activity
	}

	// 查询活动相关的SKU
	skuReq := activity.SkuListByActivityIdsReq{
		ActivityIds:  activityIds,
		EnterpriseID: request.EnterpriseID,
		GoodsName:    request.GoodsName,
		// 启用
		Status: activity.SKuActivityUpStatus,
		Page:   request.Page,
		Size:   request.Size,
	}
	total, list, err := activityService.SkuListByActivity(cx, skuReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if list == nil {
		ginplus.ResSuccess(c, activity.GoodsListResponse{
			List:  []*activity.GoodsListItem{},
			Page:  request.Page,
			Size:  request.Size,
			Total: 0,
		})
		return
	}
	// 获取小程序的logo+name
	appReq := &boss.AuthDetailRequest{
		EnterpriseID: request.EnterpriseID,
		AuthType:     1,
	}
	authDetailResponse, _ := boss.AuthDetail(ctx.Wrap(c), appReq)

	for key, item := range list {
		// 填充活动字段
		v, ok := activityInfo[item.ActivityID]
		if ok {
			// 关联活动类型
			item.ActivityType = v.ActivityType
			// 关联人数
			item.NeedNum = v.NeedNum
			// 关联商城活动ID
			item.ActivityID = v.ActivityID
		}

		// 填充商品属性字段
		tmp := new(goods.SkuItem)
		tmp.SkuAttribute = item.SkuAttribute
		tmp.SetAttributes()
		item.Attribute = tmp.Attribute
		(list[key]).ShareAppId = authDetailResponse.Data.AppID
		(list[key]).SharePath = getParams(ctx.Wrap(c), request, *item)
	}

	ginplus.ResSuccess(c, activity.GoodsListResponse{
		List:  list,
		Page:  request.Page,
		Size:  request.Size,
		Total: total,
	})

}

// 组装携带参数
func getParams(c ctx.Context, req activity.GoodsListRequest, value activity.GoodsListItem) string {
	params := "/" + marketingPromotion.MarketingPath[marketingPromotion.MarketingHelperPost] + ".html"
	params += "?spu_id=" + strconv.FormatInt(value.SpuID, 10) + "&id=" + strconv.FormatInt(value.SpuID, 10) + "&sku_id=" + strconv.FormatInt(value.SkuID, 10) + "&type=" + strconv.Itoa(value.ActivityType) + "&rule_id=" + strconv.FormatInt(value.ActivityID, 10) + "&store_id=0&mem_uid=" + req.MemberUserId + "&share_uid=" + req.MemberUserId + "&source=scrm"
	return params
}
