package index_http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/index"
	"gw-scrm/internal/modules/index/api"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/cdp"
	"sort"
	"time"
)

// 私域客户 - 企微客户 群客户 粉丝客户
func indexData(c *gin.Context) {
	shopFeeParam := &index.ShopFeeParam{}
	err := c.BindJSON(shopFeeParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	shopFeeParam.EnterpriseId = tokenInfo.EnterpriseId
	shopFeeParam.Type = 1

	fansCdpSign := index.FansBody{}
	fansCdpSign.SignTime = time.Now().Unix()
	IndexApi := api.IndexApi{}
	fansNewRes := index.FansCDPRes{}

	friendRes, err := IndexApi.QwFriend(ctx.Wrap(c), *shopFeeParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	FriendRes := index.Friend{}
	newFriendRes, _ := json.Marshal(friendRes.(map[string]interface{}))
	_ = json.Unmarshal(newFriendRes, &FriendRes)

	groupRes, err := IndexApi.GroupMemberNum(ctx.Wrap(c), *shopFeeParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	GroupRes := index.Group{}
	newGroupRes, _ := json.Marshal(groupRes.(map[string]interface{}))
	_ = json.Unmarshal(newGroupRes, &GroupRes)

	fans := make([]index.Fans, len(fansNewRes.Data))
	privateDomain := index.PrivateDomain{}
	if len(fansNewRes.Data) > 0 {
		for k, v := range fansNewRes.Data {
			fans[k].Name = v.Name
			fans[k].Appid = v.Appid
			fans[k].V.Now.Uv = v.V.Now[0]
			fans[k].V.Now.UvThan = v.V.Now[1]
			fans[k].V.Now.Direction = index.ThreeDirection(int64(v.V.Now[0]), int64(v.V.Now[1]))
			fans[k].V.Now.Percent = index.Percent(int64(v.V.Now[0]), int64(v.V.Now[1]))

			fans[k].V.New.Uv = v.V.New[0]
			fans[k].V.New.UvThan = v.V.New[1]
			fans[k].V.New.Direction = index.ThreeDirection(int64(v.V.New[0]), int64(v.V.New[1]))
			fans[k].V.New.Percent = index.Percent(int64(v.V.New[0]), int64(v.V.New[1]))

			fans[k].V.Del.Uv = v.V.Del[0]
			fans[k].V.Del.UvThan = v.V.Del[1]
			fans[k].V.Del.Direction = index.ThreeDirection(int64(v.V.Del[0]), int64(v.V.Del[1]))
			fans[k].V.Del.Percent = index.Percent(int64(v.V.Del[0]), int64(v.V.Del[1]))
		}
	}
	privateDomain.Fans = fans

	// 处理好友数
	privateDomain.Friend.AddUv = FriendRes.AddUv
	privateDomain.Friend.AddUvThan = FriendRes.AddUvThan
	privateDomain.Friend.AddUvDirection = index.ThreeDirection(FriendRes.AddUv, FriendRes.AddUvThan)
	privateDomain.Friend.AddUvPercent = index.Percent(FriendRes.AddUv, FriendRes.AddUvThan)
	privateDomain.Friend.ReUv = FriendRes.ReUv
	privateDomain.Friend.ReUvThan = FriendRes.ReUvThan
	privateDomain.Friend.ReUvDirection = index.ThreeDirection(FriendRes.ReUv, FriendRes.ReUvThan)
	privateDomain.Friend.ReUvPercent = index.Percent(FriendRes.ReUv, FriendRes.ReUvThan)
	privateDomain.Friend.LossUv = FriendRes.LossUv
	privateDomain.Friend.LossUvThan = FriendRes.LossUvThan
	privateDomain.Friend.LossUvDirection = index.ThreeDirection(FriendRes.LossUv, FriendRes.LossUvThan)
	privateDomain.Friend.LossUvPercent = index.Percent(FriendRes.LossUv, FriendRes.LossUvThan)
	// 处理群成员
	privateDomain.Group.AddGroupUv = GroupRes.AddGroupUv
	privateDomain.Group.AddGroupUvThan = GroupRes.AddGroupUvThan
	privateDomain.Group.AddGroupUvDirection = index.ThreeDirection(GroupRes.AddGroupUv, GroupRes.AddGroupUvThan)
	privateDomain.Group.AddGroupUvPercent = index.Percent(GroupRes.AddGroupUv, GroupRes.AddGroupUvThan)
	privateDomain.Group.ReGroupUv = GroupRes.ReGroupUv
	privateDomain.Group.ReGroupUvThan = GroupRes.ReGroupUvThan
	privateDomain.Group.ReGroupUvDirection = index.ThreeDirection(GroupRes.ReGroupUv, GroupRes.ReGroupUvThan)
	privateDomain.Group.ReGroupUvPercent = index.Percent(GroupRes.ReGroupUv, GroupRes.ReGroupUvThan)
	privateDomain.Group.LossGroupUv = GroupRes.LossGroupUv
	privateDomain.Group.LossGroupUvThan = GroupRes.LossGroupUvThan
	privateDomain.Group.LossGroupUvDirection = index.ThreeDirection(GroupRes.LossGroupUv, GroupRes.LossGroupUvThan)
	privateDomain.Group.LossGroupUvPercent = index.Percent(GroupRes.LossGroupUv, GroupRes.LossGroupUvThan)
	privateDomain.Group.GroupNum = GroupRes.GroupNum
	privateDomain.Group.GroupNumThan = GroupRes.GroupNumThan
	privateDomain.Group.GroupNumUvDirection = index.ThreeDirection(GroupRes.GroupNum, GroupRes.GroupNumThan)
	privateDomain.Group.GroupNumUvPercent = index.Percent(GroupRes.GroupNum, GroupRes.GroupNumThan)

	ginplus.ResSuccess(c, privateDomain)
	return
}

// 商城变现 - 完成
func ShopFee(c *gin.Context) {
	shopFeeParam := &index.ShopFeeParam{}
	err := c.BindJSON(shopFeeParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	shopFeeParam.EnterpriseId = tokenInfo.EnterpriseId
	shopFeeParam.Type = 0

	cdpSign := index.ShopFeeBody{}
	cdpSign.SignTime = time.Now().Unix()
	cdpSign.Indices = []string{"ss_fee", "re_fee", "uv", "order_uv", "pay_uv", "pay_num"}
	sign, err := cdp.GetSign(map[string]interface{}{
		"sign_time": cdpSign.SignTime,
		"indices":   cdpSign.Indices,
		"store_id":  cdpSign.StoreId,
	})
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	cdpSign.Sign = sign

	IndexApi := api.IndexApi{}
	res, err := IndexApi.ShopFee(ctx.Wrap(c), *shopFeeParam, cdpSign)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	newRes := index.ShopFeeRes{}

	var m map[string]interface{}
	_ = json.Unmarshal([]byte(res.(string)), &m)
	if ok := m["code"]; ok != nil {
		ginplus.ResSuccess(c, newRes.Data)
	} else {
		_ = json.Unmarshal([]byte(res.(string)), &newRes)

		ShopFeeRealRes := index.ShopFeeRealRes{}
		ShopFeeRealRes.Uv = newRes.Data[0].Uv
		ShopFeeRealRes.SsFee = newRes.Data[0].SsFee
		ShopFeeRealRes.ReFee = newRes.Data[0].ReFee
		ShopFeeRealRes.OrderUv = newRes.Data[0].OrderUv
		ShopFeeRealRes.PayUv = newRes.Data[0].PayUv
		ShopFeeRealRes.PayNum = newRes.Data[0].PayNum
		ShopFeeRealRes.UnitPrice = index.UnitPrice(newRes.Data[0].SsFee, float64(newRes.Data[0].PayUv))
		ShopFeeRealRes.UvOrderThan = index.ConversionPercent(newRes.Data[0].OrderUv, newRes.Data[0].Uv)
		ShopFeeRealRes.OrderPayThan = index.ConversionPercent(newRes.Data[0].PayUv, newRes.Data[0].OrderUv)
		ShopFeeRealRes.UvPayThan = index.ConversionPercent(newRes.Data[0].PayUv, newRes.Data[0].Uv)

		ginplus.ResSuccess(c, ShopFeeRealRes)
	}

	return
}

// 总客户转化
func TotalUv(c *gin.Context) {
	shopFeeParam := &index.ShopFeeParam{}
	err := c.BindJSON(shopFeeParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	shopFeeParam.EnterpriseId = tokenInfo.EnterpriseId
	shopFeeParam.Type = 0

	cdpSign := index.ShopFeeBody{}
	cdpSign.SignTime = time.Now().Unix()
	cdpSign.Indices = []string{"all_uv", "pays_uv", "fg_uv", "fg_time", "zh_time", "min_time"}
	cdpSign.StoreId = shopFeeParam.StoreId
	sign, err := cdp.GetSign(map[string]interface{}{
		"sign_time": cdpSign.SignTime,
		"indices":   cdpSign.Indices,
		"store_id":  cdpSign.StoreId,
	})
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	cdpSign.Sign = sign

	IndexApi := api.IndexApi{}
	res, err := IndexApi.ShopFee(ctx.Wrap(c), *shopFeeParam, cdpSign)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	newRes := index.ShopConversion{}
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(res.(string)), &m)

	if ok := m["code"]; ok != nil {
		ginplus.ResSuccess(c, newRes.Data)
	} else {
		_ = json.Unmarshal([]byte(res.(string)), &newRes)

		ShopConversionRealRes := index.ShopConversionRealRes{}
		ShopConversionRealRes.AllUv = newRes.Data[0].AllUv
		ShopConversionRealRes.PaysUv = newRes.Data[0].PaysUv
		ShopConversionRealRes.FgUv = newRes.Data[0].FgUv
		ShopConversionRealRes.FgTime = newRes.Data[0].FgTime
		ShopConversionRealRes.ZhTime = newRes.Data[0].ZhTime
		ShopConversionRealRes.MinTime = newRes.Data[0].MinTime
		ShopConversionRealRes.UvPayThan = index.ConversionPercent(newRes.Data[0].PaysUv, newRes.Data[0].AllUv)
		ShopConversionRealRes.RepeatPayThan = index.ConversionPercent(newRes.Data[0].FgUv, newRes.Data[0].PaysUv)

		ginplus.ResSuccess(c, ShopConversionRealRes)
	}

	return
}

// 导购排行 - 完成
func ShopperList(c *gin.Context) {
	shopperParam := &index.ShopperListParam{}
	err := c.BindJSON(shopperParam)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	shopperParam.EnterpriseId = tokenInfo.EnterpriseId

	cdpSign := index.ShopperListBody{}
	cdpSign.SignTime = time.Now().Unix()
	cdpSign.Indices = []string{"yj_fee", "uv2", "qw_uv", "coupon_fee", "order_fee"}
	sign, err := cdp.GetSign(map[string]interface{}{
		"sign_time": cdpSign.SignTime,
		"indices":   cdpSign.Indices,
	})
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	cdpSign.Sign = sign

	IndexApi := api.IndexApi{}
	res, err := IndexApi.ShopperList(ctx.Wrap(c), *shopperParam, cdpSign)
	if err != nil {
		ginplus.ResError(c, err)
		c.Abort()
		return
	}
	newRes := index.ShopperListRes{}
	newRealRes := index.ShopperListRealRes{}
	_ = json.Unmarshal([]byte(res.(string)), &newRes)

	var m map[string]interface{}
	_ = json.Unmarshal([]byte(res.(string)), &m)

	if ok := m["code"]; ok != nil {
		ginplus.ResSuccess(c, newRealRes)
	} else {
		newRes.Total.Discount = index.ConversionPercentStr(newRes.Total.CouponFee, newRes.Total.OrderFee)
		newRes.Total.Contribution = "100%"
		newRes.Total.UV = newRes.Total.UV2
		newRes.Total.QwUv = newRes.Total.QwLayUv
		newRealRes.Total = newRes.Total

		for _, v := range newRes.Data[0] {
			v.UV = v.UV2
			v.QwUv = v.QwLayUv
			v.Discount = index.ConversionPercentFloat(v.CouponFee, v.OrderFee)
			v.Discount = index.ConversionPercentFloat(v.CouponFee, v.OrderFee)
			v.Contribution = index.ConversionPercentFloat(v.YjFee, newRes.Total.YjFee)
			newRealRes.List = append(newRealRes.List, v)
		}

		sort.SliceStable(newRealRes.List, func(i, j int) bool {
			// 1为降序
			if shopperParam.YjSort == 1 {
				return newRealRes.List[i].YjFee < newRealRes.List[j].YjFee
			} else {
				return newRealRes.List[i].YjFee > newRealRes.List[j].YjFee
			}
		})
		// 1为降序
		var indexKey int64 = 0
		var indexLen = int64(len(newRealRes.List)) + 1

		if shopperParam.YjSort == 1 {
			for k, _ := range newRealRes.List {
				indexLen = indexLen - 1
				newRealRes.List[k].Id = indexLen
			}
		} else {
			for k, _ := range newRealRes.List {
				indexKey = indexKey + 1
				newRealRes.List[k].Id = indexKey
			}
		}
		ginplus.ResSuccess(c, newRealRes)
	}

	return
}

// 首页公告 - 完成
func NoticeHelp(c *gin.Context) {
	// 公告
	var param = map[string]interface{}{
		"not_page": 1,
		"search": map[string]interface{}{
			"time_type": "released_at",
			"start_at":  0,
			"end_at":    time.Now().Unix(),
		},
	}
	indexApi := api.IndexApi{}
	res, err := indexApi.Notice(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	helpRes, err := indexApi.HelpList(ctx.Wrap(c), map[int]string{})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	indexNotice, err := json.Marshal(res.List)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	for k, v := range helpRes.Data.List {
		if v.Key == "myapp_home" {
			helpRes.Data.List[k].Content = string(indexNotice)
		}
	}

	ginplus.ResSuccess(c, helpRes.Data)
	return
}

// 门店列表
func StoreList(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	param := map[string]int64{
		"enterprise_id": tokenInfo.EnterpriseId,
	}

	indexApi := api.IndexApi{}
	res, err := indexApi.StoreList(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

func BaseConfigCheck(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	param := map[string]int64{
		"enterprise_id": tokenInfo.EnterpriseId,
	}

	indexApi := api.IndexApi{}
	res, err := indexApi.BaseConfigCheck(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)

	return
}

func BuyNotice(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	param := map[string]int64{
		"enterprise_id": tokenInfo.EnterpriseId,
	}

	indexApi := api.IndexApi{}
	res, err := indexApi.BuyNotice(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}
