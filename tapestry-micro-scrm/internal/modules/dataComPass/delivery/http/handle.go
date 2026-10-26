package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	cache2 "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"micro-scrm/internal/logic/smart_shopping_guide"
	"micro-scrm/internal/modules/dataComPass"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	response "micro-scrm/internal/pkg/resoponse"
)

type GoodFriendsHandler struct {
	Name string
}

func NewGoodFriendsHandler() *GoodFriendsHandler {
	return &GoodFriendsHandler{
		Name: "好友数据",
	}
}

func (g GoodFriendsHandler) getDataAll(c *gin.Context) {
	var request dataComPass.DataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
}

func (g GoodFriendsHandler) getGroupNumber(c *gin.Context) {
	var request dataComPass.DataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}
	var keyString string
	newTime := dataComPass.EndTimeTOCorde(request.EndTime)
	var sString string
	for _, s := range request.UserIds {
		sString += s
	}
	searchString := helper.Md5V(sString)
	keyString = fmt.Sprintf("%s:%d:%d:%d:%s", "dataCompass", request.EnterpriseID, request.StartTime, newTime, searchString)
	//keyString = fmt.Sprintf("%s:%d:%d:%d", "dataCompass", request.EnterpriseID, request.StartTime, newTime)
	//} else {
	//	var sString string
	//	for _, s := range request.ShoppingGuide {
	//		sString += s
	//	}
	//	searchString := helper.Md5V(sString)
	//	keyString = fmt.Sprintf("%s:%d:%d:%d:%s", "dataCompass", request.EnterpriseID, request.StartTime, newTime, searchString)
	//}

	//先查ck如果ck异常了在查mysql
	newSearchCk := &dataComPass.CompassCacheCk{
		Cache: cache2.NewRedisCache(redis.Client()),
		Key:   keyString,
	}
	resultCk, errCk := newSearchCk.DataInfoResultCk(ctx.Wrap(c), request)
	if errCk != nil {

		//Mysql查询到的业务
		newSearch := &dataComPass.CompassCache{
			Cache: cache2.NewRedisCache(redis.Client()),
			Key:   keyString,
		}
		result, err := newSearch.DataInfoResult(ctx.Wrap(c), request)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		ginplus.ResSuccess(c, &result)

	} else {
		ginplus.ResSuccess(c, &resultCk)
	}
}

func (g GoodFriendsHandler) datacompassCommunityData(c *gin.Context) {
	params := new(dataComPass.DataRequest)
	if err := c.BindJSON(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.StartTime == 0 || params.EndTime == 0 {
		ginplus.ResError(c, errors.New("请求的时间有误"))
		return
	}
	if params.StartTime == params.EndTime {
		params.StartTime, params.EndTime = helper.StartAndEndByDayTime(params.StartTime)
	}

	var keyString string
	var countString string
	newTime := dataComPass.EndTimeTOCorde(params.EndTime)
	if len(params.ShoppingGuide) == 0 {
		keyString = fmt.Sprintf("%s:%d:%d:%d:%d:%d", "dataCompassCommunityDetail", params.EnterpriseID, params.StartTime, newTime, params.PageSize, params.PageNum)
		countString = fmt.Sprintf("%s:%d:%d:%d:%s", "dataCompassCommunityDetail", params.EnterpriseID, params.StartTime, newTime, "count")
	} else {
		var sString string
		for _, s := range params.ShoppingGuide {
			sString += s
		}
		keyString = fmt.Sprintf("%s:%d:%d:%d:%d:%d:%s", "dataCompassCommunityDetail", params.EnterpriseID, params.StartTime, newTime, params.PageSize, params.PageNum, sString)
		countString = fmt.Sprintf("%s:%d:%d:%d:%s:%s", "dataCompassCommunityDetail", params.EnterpriseID, params.StartTime, newTime, sString, "count")
	}
	newResult := &dataComPass.CompassCommunityCache{
		Cache: cache2.NewRedisCache(redis.Client()),
		Key:   keyString,
		Count: countString,
	}
	data, count, err := newResult.GetDataCompassCommunityData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	//result, total := params.DatacompassCommunityData(ctx.Wrap(c))
	response.SuccessList(c, count, params.Page.PageNum, params.Page.PageSize, data)
}

func (g GoodFriendsHandler) datacompassFriendData(c *gin.Context) {
	params := new(dataComPass.DataRequest)
	if err := c.ShouldBindJSON(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.StartTime == 0 || params.EndTime == 0 {
		ginplus.ResError(c, errors.New("请求的时间有误"))
		return
	}
	if params.StartTime == params.EndTime {
		params.StartTime, params.EndTime = helper.StartAndEndByDayTime(params.StartTime)
	}
	var keyString string
	var countString string
	newTime := dataComPass.EndTimeTOCorde(params.EndTime)

	keyString = fmt.Sprintf("%s:%d:%d:%d:%d:%d", "dataCompassFriendDetail", params.EnterpriseID, params.StartTime, newTime, params.PageSize, params.PageNum)
	countString = fmt.Sprintf("%s:%d:%d:%d:%s", "dataCompassFriendDetail", params.EnterpriseID, params.StartTime, newTime, "count")

	newResult := &dataComPass.CompassCommunityCache{
		Cache: cache2.NewRedisCache(redis.Client()),
		Key:   keyString,
		Count: countString,
	}
	data, count, err := newResult.GetDataCompassFriendData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, count, params.Page.PageNum, params.Page.PageSize, data)
}

// 根据智慧导购的导购的mem_uid，获取导购的一级部门
func (g GoodFriendsHandler) getBaFirstDepartment(c *gin.Context) {
	param := dataComPass.GetBaFirstDepartmentParam{}
	if paramErr := c.ShouldBindJSON(&param); paramErr != nil {
		ginplus.ResError(c, paramErr)
		return
	}

	userID, userErr := smart_shopping_guide.NewSmartShoppingGuideLogic(param.EnterpriseID).GetBaUserID(ctx.Wrap(c), param.MemID)
	if userErr != nil {
		ginplus.ResError(c, userErr)
		return
	}

	userInfo, infoErr := dataComPass.GetBaFirstDepartment(param.EnterpriseID, userID)
	if infoErr != nil {
		ginplus.ResError(c, userErr)
		return
	}

	ginplus.ResSuccess(c, struct {
		MemID          string `json:"mem_id"`
		UserID         string `json:"user_id"`
		DepartmentName string `json:"department_name"`
		Name           string `json:"name"`
	}{
		MemID:          param.MemID,
		UserID:         userID,
		DepartmentName: userInfo.Department.Name,
		Name:           userInfo.Member.Name,
	})
}

func (g GoodFriendsHandler) rankData(c *gin.Context) {

	params := new(dataComPass.DataRequest)
	if err := c.ShouldBindJSON(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if params.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID不能为空"))
		return
	}
	if params.EndTime == 0 {
		ginplus.ResError(c, errors.New("请求的时间有误"))
		return
	}
	if params.Type != 1 && params.Type != 2 {
		ginplus.ResError(c, errors.New("请求的类型有误"))
		return
	}
	if params.StartTime == params.EndTime {
		params.StartTime, params.EndTime = helper.StartAndEndByDayTime(params.StartTime)
	}

	var keyString string
	var countString string
	newTime := dataComPass.EndTimeTOCorde(params.EndTime)
	if len(params.ShoppingGuide) == 0 {
		keyString = fmt.Sprintf("%s:%d:%d:%d:%d:%d:%d", "dataCompassRank", params.EnterpriseID, params.StartTime, newTime, params.PageSize, params.PageNum, params.Type)
		countString = fmt.Sprintf("%s:%d:%d:%d:%s:%d", "dataCompassRank", params.EnterpriseID, params.StartTime, newTime, "count", params.Type)
	} else {
		var sString string
		for _, s := range params.ShoppingGuide {
			sString += s
		}
		keyString = fmt.Sprintf("%s:%d:%d:%d:%d:%d:%d", "dataCompassRank", params.EnterpriseID, params.StartTime, newTime, params.PageSize, params.PageNum, params.Type)
		countString = fmt.Sprintf("%s:%d:%d:%d:%s:%d", "dataCompassRank", params.EnterpriseID, params.StartTime, newTime, "count", params.Type)
	}
	newResult := &dataComPass.CompassCommunityCache{
		Cache: cache2.NewRedisCache(redis.Client()),
		Key:   keyString,
		Count: countString,
	}

	data, count, err := newResult.GetRankData(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, count, 1, int(count), data)
}

// getGroupNumberCk  数据罗盘统计数据Ck
func (g GoodFriendsHandler) getGroupNumberCk(c *gin.Context) {
	var request dataComPass.DataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if request.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业ID为空！"))
		return
	}

	var keyString string
	newTime := dataComPass.EndTimeTOCorde(request.EndTime)
	var sString string
	for _, s := range request.UserIds {
		sString += s
	}
	searchString := helper.Md5V(sString)
	keyString = fmt.Sprintf("%s:%d:%d:%d:%s", "dataCompass", request.EnterpriseID, request.StartTime, newTime, searchString)
	newSearchCk := &dataComPass.CompassCacheCk{
		Cache: cache2.NewRedisCache(redis.Client()),
		Key:   keyString,
	}
	result, err := newSearchCk.DataInfoResultCk(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, &result)

}
