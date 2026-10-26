package marketingPromotion

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/golang/freetype"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/imagex/poster/core"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/imagex/poster/handler"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRespository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/microauthority"
	"micro-scrm/internal/pkg/sdk/uc"
	"micro-scrm/internal/pkg/sdk/work"
)

const (
	EnterpriseExpiration = 24 * time.Hour
	rootpath             = "source/"
	fontPath             = rootpath + "font/msyh.ttc"
)

var MarketingPath = map[int]string{
	MarketingGoods:         "goods/goods-detail/goods-detail",
	MarketingCoupon:        "other/receive-coupons/receive-coupons",
	MarketingHelperPost:    "goods/goods-detail/goods-detail",
	MarketingHelperBuyGive: "goods/goods-detail/goods-detail",
	MarketingSecKill:       "goods/goods-detail/goods-detail",
	MarketingLottery:       "goods/goods-detail/goods-detail",
	MarketingGroup:         "goods/goods-detail/goods-detail",
	MarketingBargain:       "goods/goods-detail/goods-detail",
}

// Service 应用服务
type Service struct {
	Repo Repository
}

// 获取太阳码携带参数
func getCodeParams(c ctx.Context, domain string, req GetPosterRequest) (CodeParam, error) {
	members, err := uc.MemberByUserId(c, uc.MemberByUserIdRequest{
		MemUserID: req.MemUserID,
	}, "crs-manu-customer", domain)
	if err != nil {
		return nil, err
	}
	memUid := ""
	if members != nil {
		memUid = members.MemUid
	}
	switch req.Type {
	case MarketingGoods:
		goods := &GoodsParams{}
		json.Unmarshal([]byte(req.Params), goods)
		params := GoodsCodeParam{
			SpuID:  goods.SpuID,
			SkuID:  goods.SkuID,
			MemUID: memUid,
			//MemUID:  "5DD6A3397F000001655FB786D1E60C04",
			StoreID:  0,
			ShareUid: req.MemUserID,
			//ShareUid: "zx",
			Source: "scrm",
		}
		return &params, nil
	case MarketingCoupon:
		coupon := &CouponParams{}
		json.Unmarshal([]byte(req.Params), coupon)
		params := CouponCodeParam{
			CouponID: coupon.CouponID,
			MemUID:   memUid,
			StoreID:  0,
			ShareUid: req.MemUserID,
			Source:   "scrm",
			ShareImg: "",
		}
		return &params, nil
	case MarketingGroup, MarketingHelperPost, MarketingHelperBuyGive, MarketingSecKill, MarketingLottery, MarketingBargain:
		activity := &ActivityParams{}
		json.Unmarshal([]byte(req.Params), activity)
		params := GroupCodeParam{
			SpuID:    activity.SpuID,
			SkuID:    activity.SkuID,
			MemUID:   memUid,
			StoreID:  0,
			Type:     req.Type,
			RuleID:   activity.ShopActivityID,
			ShareUid: req.MemUserID,
			Source:   "scrm",
		}
		return &params, nil
	}
	return nil, nil
}

// GetPoster 生成海报
func (s Service) GetPoster(c ctx.Context, domain string, req GetPosterRequest) (*GetPosterResponse, error) {
	startTime := time.Now().Unix()
	costTime := int64(0)
	var getPosterResponse = &GetPosterResponse{}
	// 获取缓存中的太阳码
	codeResp, err := s.getSunCode(c, domain, req)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("生成太阳码错误：%+v", err)
		return getPosterResponse, err
	}
	sunTime := time.Now().Unix()
	costTime = sunTime - startTime
	if req.Type == MarketingCoupon {
		req.ShareImg = codeResp.Data.ShareImg
	}
	req.SunCode = codeResp.Data.Url
	// 生成海报
	switch req.Type {
	case MarketingGoods:
		getPosterResponse = s.GoodsPoster(c, req)
	case MarketingCoupon:
		getPosterResponse = s.CouponPoster(c, req)
	default:
		getPosterResponse = s.ActivityPoster(c, req)
	}
	posterTime := time.Now().Unix()
	costTime = posterTime - sunTime
	log.Info(ctx.Wrap(c)).Msgf("生成海报: startTime:%d, NowTime:%d, costTime:%d", sunTime, time.Now().Unix(), costTime)
	return getPosterResponse, nil
}

// 生成活动海报
func (s Service) ActivityPoster(c ctx.Context, req GetPosterRequest) *GetPosterResponse {
	activityParams := &ActivityParams{}
	err := json.Unmarshal([]byte(req.Params), activityParams)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("Unmarshal ERR：%+v", err)

	}
	log.Info(ctx.Wrap(c)).Msgf("Unmarshal结果：%+v", activityParams)
	// 获取缓存
	posterKey := generateKey(req)
	posterCache := redis.Client().Get(c, posterKey).Val()
	poster := &GetPosterResponse{}
	if posterCache != "" {
		jErr := json.Unmarshal([]byte(posterCache), poster)
		if jErr != nil {
			return &GetPosterResponse{}
		}
		if poster.Object != "" {
			return &GetPosterResponse{
				Object: poster.Object,
			}
		}
	}

	// 获取小程序的logo
	appInfo, err := getAppInfo(c, req.EnterpriseID)
	if err != nil {
		return &GetPosterResponse{}
	}
	nullHandler := &handler.NullHandler{}
	multiple := 2
	des := &handler.Context{
		//图片都绘在这个PNG载体上
		PngCarrier: core.NewPNG(0, 0, 375*multiple, 660*multiple),
	}
	// 00.绘制header
	//rectangHandler := &handler.ImageRectangleHandler{
	//	X:      0,
	//	Y:      0,
	//	Width:  375 * multiple,
	//	Height: 70 * multiple,
	//	R:      210,
	//	G:      0,
	//	B:      0,
	//	A:      80,
	//}
	// 0.绘制背景图
	backgroundHandler := &handler.BackgroundHandler{
		X:    0,                                          // 图片x坐标
		Y:    0,                                          // 图片y坐标
		Path: rootpath + "sourceImages/activityback.png", //图片路径
	}
	activityParams.AppLogo = appInfo.Data.Logo
	// 1.小程序logo
	appLogoHandler := &handler.ImageCircleRemoteHandler{
		X:      15 * multiple, // 图片x坐标
		Y:      20 * multiple, // 图片y坐标
		Width:  30 * multiple,
		Height: 30 * multiple,
		R:      255,
		B:      0,
		G:      0,
		A:      0,
		URL:    activityParams.AppLogo,
	}
	activityParams.AppName = appInfo.Data.Name
	// 2.小程序名称
	appNameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        55 * multiple,          // 文字x坐标
		Y:        43 * multiple,          // 文字y坐标
		Size:     float64(20 * multiple), // 文字大小
		R:        255,                    // 文字颜色RGB值
		G:        255,
		B:        255,
		Text:     activityParams.AppName, // 文字内容
		FontPath: fontPath,               // 字体文件
	}
	icon := GetActivityIcon(activityParams.Type)
	// 3. 砍价iocn
	topIconHandler := &handler.ImageLocalHandler{
		X:    icon.X * multiple, // 二维码x坐标
		Y:    icon.Y * multiple, // 二维码y坐标
		Path: icon.TopIcon,
	}
	// 4.员工头像
	mavatarHandler := &handler.ImageCircleRemoteHandler{
		X:      15 * multiple, // 图片x坐标
		Y:      85 * multiple, // 图片y坐标
		Width:  45 * multiple,
		Height: 45 * multiple,
		URL:    req.MemberAvatar,
	}

	// 2.员工姓名
	mnameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        75 * multiple,          // 文字x坐标
		Y:        103 * multiple,         // 文字y坐标
		Size:     float64(15 * multiple), // 文字大小
		R:        127,                    // 文字颜色RGB值
		G:        127,
		B:        127,
		Text:     req.MemberName, // 文字内容
		FontPath: fontPath,       // 字体文件
	}
	// 计算员工姓名字节长度
	font := freetype.NewContext()
	MemberNameCount := utf8.RuneCountInString(req.MemberName)
	// 3.为你推荐
	recommendHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        75*multiple + int(font.PointToFixed(float64(15*multiple*MemberNameCount))>>6), // todo 计算坐标
		Y:        103 * multiple,
		Size:     float64(15 * multiple),
		R:        50,
		G:        50,
		B:        50,
		Text:     "为你推荐",
		FontPath: fontPath,
	}
	// 双引号
	//marksHandler := &handler.ImageLocalHandler{
	//	X:    300 * multiple, // 二维码x坐标
	//	Y:    91 * multiple,  // 二维码y坐标
	//	Path: rootpath + "sourceImages/marks@2x.png",
	//}
	//// 4.认证ICON
	//authIconHandler := &handler.ImageLocalHandler{
	//	X:    75 * multiple,  // 二维码x坐标
	//	Y:    111 * multiple, // 二维码y坐标
	//	Path: rootpath + "sourceImages/renzheng@2x.png",
	//}
	// 5.商品名称
	//计算商品名称长度
	skuNameCount := utf8.RuneCountInString(activityParams.SkuName)
	//计算总共需要几行
	base := (375 - 15*2) * multiple
	lineCount := math.Ceil(float64(skuNameCount*multiple*18) / float64(base))
	// 计算每行可放多少字
	each := base / (multiple * 18)
	rs := []rune(activityParams.SkuName)
	one := 0
	if len(rs) > each {
		one = each
	} else {
		one = len(rs)
	}
	goodsNameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        15 * multiple,          // 文字x坐标
		Y:        167 * multiple,         // 文字y坐标
		Size:     18 * float64(multiple), // 文字大小
		R:        50,                     // 文字颜色RGB值
		G:        50,
		B:        50,
		Text:     string(rs[0:one]), // 文字内容
		FontPath: fontPath,          // 字体文件
	}
	sgoodsNameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}

	if lineCount > 1 {
		sgoodsNameHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        15 * multiple,          // 文字x坐标
			Y:        192 * multiple,         // 文字y坐标
			Size:     18 * float64(multiple), // 文字大小
			R:        50,                     // 文字颜色RGB值
			G:        50,
			B:        50,
			Text:     string(rs[each:len(rs)]), // 文字内容
			FontPath: fontPath,                 // 字体文件
		}
	}
	// 6.商品图片
	skuImgHandler := &handler.ImageRemoteHandler{
		X:      15 * multiple,  // 图片x坐标
		Y:      213 * multiple, // 图片y坐标
		Width:  345 * multiple,
		Height: 345 * multiple,
		URL:    activityParams.Image,
	}
	// 7.¥
	priceMarkHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        15 * multiple,          // 文字x坐标
		Y:        605 * multiple,         // 文字y坐标
		Size:     float64(21 * multiple), // 文字大小
		R:        250,                    // 文字颜色RGB值
		G:        64,
		B:        70,
		Text:     "¥",      // 文字内容
		FontPath: fontPath, // 字体文件
	}
	// 8.售卖价
	sellPrice := fmt.Sprintf("%.2f", float64(activityParams.ActivityPrice)/float64(100))
	sell := []string{"", ""}
	sell = strings.Split(sellPrice, ".")
	sellPriceHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}

	sellHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}
	sellFirstCount := utf8.RuneCountInString(sell[0])
	sellSecondCount := utf8.RuneCountInString(sell[1])
	priceFirst := 0
	if sell[0] != "" {
		sellPriceHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        27 * multiple,          // 文字x坐标
			Y:        605 * multiple,         // 文字y坐标
			Size:     float64(30 * multiple), // 文字大小
			R:        250,                    // 文字颜色RGB值
			G:        64,
			B:        70,
			Text:     sell[0],  // 文字内容
			FontPath: fontPath, // 字体文件
		}
		priceFirst = int(font.PointToFixed(float64(18*multiple*sellFirstCount)) >> 6)
	}
	priceSecond := 0
	if sell[1] != "" {
		sellHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        27*multiple + priceFirst, // 文字x坐标
			Y:        605 * multiple,           // 文字y坐标
			Size:     float64(24 * multiple),   // 文字大小
			R:        250,                      // 文字颜色RGB值
			G:        64,
			B:        70,
			Text:     "." + sell[1], // 文字内容
			FontPath: fontPath,      // 字体文件
		}
		priceSecond = int(font.PointToFixed(float64(24*multiple*sellSecondCount)) >> 6)
	}
	// 9.吊牌价
	SignPrice := fmt.Sprintf("%.2f", float64(activityParams.SellPrice)/float64(100))
	SignPriceCount := utf8.RuneCountInString(SignPrice)
	markerPriceHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        15 * multiple,          // 文字x坐标
		Y:        630 * multiple,         // 文字y坐标
		Size:     float64(15 * multiple), // 文字大小
		R:        166,                    // 文字颜色RGB值
		G:        166,
		B:        166,
		Text:     "¥" + SignPrice, // 文字内容
		FontPath: fontPath,        // 字体文件
	}
	// 10. 删除线
	lineHandler := &handler.TextLineHandler{
		Next:   handler.Next{},
		X:      15 * multiple,
		Y:      625 * multiple,
		R:      166,
		G:      166,
		B:      166,
		A:      166,
		Width:  (SignPriceCount + 1) * 15,
		Height: 2,
	}
	// 11. 太阳码
	sunCodeLocalHandler := &handler.ImageCircleRemoteHandler{
		X:      280 * multiple, // 图片x坐标
		Y:      568 * multiple, // 图片y坐标
		Width:  75 * multiple,
		Height: 75 * multiple,
		URL:    req.SunCode,
	}
	// 12. 下方icon
	bottomIconHandler := &handler.ImageLocalHandler{
		X:    24*multiple + priceFirst + priceSecond, // 二维码x坐标
		Y:    580 * multiple,                         // 二维码y坐标
		Path: icon.BottomIcon,
	}
	//结束绘制，把前面的内容合并成一张图片,输出到build目录
	fileName := "/poster_" + core.RandString(10) + ".png"
	outputUrl := rootpath + "cacheimg" + fileName
	endHandler := &handler.EndHandler{
		Output: outputUrl,
	}
	// 链式调用绘制过程
	nullHandler.
		//SetNext(rectangHandler).
		SetNext(backgroundHandler).
		SetNext(appNameHandler).
		SetNext(appLogoHandler).
		SetNext(topIconHandler).
		SetNext(mnameHandler).
		SetNext(mavatarHandler).
		SetNext(recommendHandler).
		//SetNext(marksHandler).
		//SetNext(authIconHandler).
		SetNext(goodsNameHandler).
		SetNext(sgoodsNameHandler).
		SetNext(skuImgHandler).
		SetNext(sellHandler).
		SetNext(sellPriceHandler).
		SetNext(priceMarkHandler).
		SetNext(bottomIconHandler).
		SetNext(markerPriceHandler).
		SetNext(lineHandler).
		SetNext(sunCodeLocalHandler).
		SetNext(endHandler)

	// 开始执行业务
	wg := sync.WaitGroup{}
	if err := nullHandler.Run(des, &wg); err != nil {
		// 异常
		log.Info(ctx.Wrap(c)).Msgf("海报生成错误：%+v", err)
		return &GetPosterResponse{}
	}
	req.FileTmpUrl = outputUrl
	req.FileName = fileName
	url, err := changeUrlImage(c, req)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("changeUrlImage Error: %+v", err)
	}
	removeTmpImgFile(c, outputUrl)
	//成功
	posterByte, _ := json.Marshal(&GetPosterResponse{
		Object: url,
	})
	err = redis.Client().Set(c, posterKey, posterByte, EnterpriseExpiration).Err()
	if err != nil {
		return &GetPosterResponse{}
	}
	//成功
	return &GetPosterResponse{
		Object: url,
	}
}
func GetActivityIcon(types int) *ActivityIcon {
	icon := &ActivityIcon{}
	switch types {
	case MarketingLottery:
		icon.TopIcon = rootpath + "sourceImages/zhuliTop@2x.png"
		icon.BottomIcon = rootpath + "sourceImages/zhuliBottom@2x.png"
		icon.X = 258
		icon.Y = 24
	case MarketingSecKill:
		icon.TopIcon = rootpath + "sourceImages/xianshiTop@2x.png"
		icon.BottomIcon = rootpath + "sourceImages/xianshiBottom@2x.png"
		icon.X = 272
		icon.Y = 24
	case MarketingBargain:
		icon.TopIcon = rootpath + "sourceImages/kanjiaTop@2x.png"
		icon.BottomIcon = rootpath + "sourceImages/kanjiaBottom@2x.png"
		icon.X = 272
		icon.Y = 24
	case MarketingGroup:
		icon.TopIcon = rootpath + "sourceImages/pintuanTop@2x.png"
		icon.BottomIcon = rootpath + "sourceImages/pintuanBottom@2x.png"
		icon.X = 272
		icon.Y = 24
	default:
		icon.TopIcon = rootpath + "sourceImages/kanjiaTop@2x.png"
		icon.BottomIcon = rootpath + "sourceImages/kanjiaBottom@2x.png"
		icon.X = 272
		icon.Y = 24
	}
	return icon
}

// 生成普通商品二维码
func (s Service) GoodsPoster(c ctx.Context, req GetPosterRequest) *GetPosterResponse {
	goodsParams := &GoodsParams{}
	json.Unmarshal([]byte(req.Params), goodsParams)
	log.Info(ctx.Wrap(c)).Msgf("请求数据:%+v,GoodsPoster数据：%+v", req, goodsParams)
	// 获取缓存
	posterKey := generateKey(req)
	posterCache := redis.Client().Get(c, posterKey).Val()
	log.Info(ctx.Wrap(c)).Msgf("posterCache：%+v", posterCache)
	poster := &GetPosterResponse{}
	if posterCache != "" {
		jErr := json.Unmarshal([]byte(posterCache), poster)
		if jErr != nil {
			log.Info(ctx.Wrap(c)).Msgf("GoodsPosterCacheErr：%+v", jErr)
			return &GetPosterResponse{}
		}
		if poster.Object != "" {
			return &GetPosterResponse{
				Object: poster.Object,
			}
		}
	}

	nullHandler := &handler.NullHandler{}
	multiple := 2
	des := &handler.Context{
		//图片都绘在这个PNG载体上
		PngCarrier: core.NewPNG(0, 0, 375*multiple, 632*multiple),
	}
	// 0.绘制背景图
	backgroundHandler := &handler.BackgroundHandler{
		X:    0,                                       // 图片x坐标
		Y:    0,                                       // 图片y坐标
		Path: rootpath + "sourceImages/goodsBack.png", //图片路径
	}
	// 1.员工头像
	imageCircleLocalHandler := &handler.ImageCircleRemoteHandler{
		X:      15 * multiple, // 图片x坐标
		Y:      15 * multiple, // 图片y坐标
		Width:  60 * multiple,
		Height: 60 * multiple,
		URL:    req.MemberAvatar,
	}
	// 2.员工姓名
	nameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        90 * multiple,          // 文字x坐标
		Y:        40 * multiple,          // 文字y坐标
		Size:     float64(15 * multiple), // 文字大小
		R:        127,                    // 文字颜色RGB值
		G:        127,
		B:        127,
		Text:     req.MemberName, // 文字内容
		FontPath: fontPath,       // 字体文件
	}
	// 计算员工姓名字节长度
	font := freetype.NewContext()
	MemberNameCount := utf8.RuneCountInString(req.MemberName)

	// 3.为你推荐
	recommendHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        90*multiple + int(font.PointToFixed(float64(15*multiple*MemberNameCount))>>6), // todo 计算坐标
		Y:        40 * multiple,
		Size:     float64(15 * multiple),
		R:        50,
		G:        50,
		B:        50,
		Text:     "为你推荐",
		FontPath: fontPath,
	}
	// 双引号
	//marksHandler := &handler.ImageLocalHandler{
	//	X:    300 * multiple, // 二维码x坐标
	//	Y:    25 * multiple,  // 二维码y坐标
	//	Path: rootpath + "sourceImages/marks@2x.png",
	//}
	//// 4.认证ICON
	//iconHandler := &handler.ImageLocalHandler{
	//	X:    90 * multiple, // 二维码x坐标
	//	Y:    49 * multiple, // 二维码y坐标
	//	Path: rootpath + "sourceImages/renzheng@2x.png",
	//}

	// 5.商品名称
	//计算商品名称长度
	skuNameCount := utf8.RuneCountInString(goodsParams.SkuName)
	//计算总共需要几行
	base := (375 - 15*2) * multiple
	lineCount := math.Ceil(float64(skuNameCount*multiple*18) / float64(base))
	// 计算每行可放多少字
	each := base / (multiple * 18)
	rs := []rune(goodsParams.SkuName)
	one := 0
	if len(rs) > each {
		one = each
	} else {
		one = len(rs)
	}

	log.Info(ctx.New()).Msgf("总共%+v行,每行可放:%d字,商品名称长度:%+v", lineCount, each, len(rs))
	goodsNameHandler := &handler.TextHandler{
		Next: handler.Next{},
		X:    15 * multiple,          // 文字x坐标
		Y:    110 * multiple,         // 文字y坐标
		Size: 18 * float64(multiple), // 文字大小
		R:    127,                    // 文字颜色RGB值
		G:    127,
		B:    127,
		//Text:     goodsParams.SkuName, // 文字内容
		Text:     string(rs[0:one]), // 文字内容
		FontPath: fontPath,          // 字体文件
	}
	sgoodsNameHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}
	if lineCount > 1 {
		sgoodsNameHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        15 * multiple,          // 文字x坐标
			Y:        135 * multiple,         // 文字y坐标
			Size:     18 * float64(multiple), // 文字大小
			R:        127,                    // 文字颜色RGB值
			G:        127,
			B:        127,
			Text:     string(rs[each:len(rs)]), // 文字内容
			FontPath: fontPath,                 // 字体文件
		}
	}
	// 6.商品图片
	skuImgHandler := &handler.ImageRemoteHandler{
		X:      15 * multiple,  // 图片x坐标
		Y:      155 * multiple, // 图片y坐标
		Width:  345 * multiple,
		Height: 345 * multiple,
		URL:    goodsParams.SkuPicture,
	}
	// 7.¥
	priceMarkHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        15 * multiple,          // 文字x坐标
		Y:        540 * multiple,         // 文字y坐标
		Size:     float64(21 * multiple), // 文字大小
		R:        250,                    // 文字颜色RGB值
		G:        64,
		B:        70,
		Text:     "¥",      // 文字内容
		FontPath: fontPath, // 字体文件
	}
	// 8.售卖价
	sellPrice := fmt.Sprintf("%.2f", float64(goodsParams.SellPrice)/float64(100))
	sell := []string{"", ""}
	sell = strings.Split(sellPrice, ".")
	sellPriceHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}
	sellHandler := &handler.TextHandler{
		Next:     handler.Next{},
		FontPath: fontPath, // 字体文件
	}
	sellFirstCount := utf8.RuneCountInString(sell[0])
	sellSecondCount := utf8.RuneCountInString(sell[1])
	priceFirst := 0
	if sell[0] != "" {
		sellPriceHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        27 * multiple,          // 文字x坐标
			Y:        541 * multiple,         // 文字y坐标
			Size:     float64(30 * multiple), // 文字大小
			R:        250,                    // 文字颜色RGB值
			G:        64,
			B:        70,
			Text:     sell[0],  // 文字内容
			FontPath: fontPath, // 字体文件
		}
		priceFirst = int(font.PointToFixed(float64(18*multiple*sellFirstCount)) >> 6)
	}
	priceSecond := 0
	if sell[1] != "" {
		sellHandler = &handler.TextHandler{
			Next:     handler.Next{},
			X:        27*multiple + priceFirst, // 文字x坐标
			Y:        541 * multiple,           // 文字y坐标
			Size:     float64(24 * multiple),   // 文字大小
			R:        250,                      // 文字颜色RGB值
			G:        64,
			B:        70,
			Text:     "." + sell[1], // 文字内容
			FontPath: fontPath,      // 字体文件
		}
		priceSecond = int(font.PointToFixed(float64(24*multiple*sellSecondCount)) >> 6)
	}
	// 9.吊牌价
	SignPrice := fmt.Sprintf("%.2f", float64(goodsParams.SignPrice)/float64(100))
	SignPriceCount := utf8.RuneCountInString(SignPrice)
	markerPriceHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        35*multiple + priceFirst + priceSecond, // 文字x坐标
		Y:        541 * multiple,                         // 文字y坐标
		Size:     float64(15 * multiple),                 // 文字大小
		R:        166,                                    // 文字颜色RGB值
		G:        166,
		B:        166,
		Text:     "¥" + SignPrice, // 文字内容
		FontPath: fontPath,        // 字体文件
	}

	// 10. 删除线
	lineHandler := &handler.TextLineHandler{
		Next:   handler.Next{},
		X:      35*multiple + priceFirst + priceSecond,
		Y:      536 * multiple,
		R:      166,
		G:      166,
		B:      166,
		A:      166,
		Width:  (SignPriceCount + 1) * 15,
		Height: 2,
	}
	// 11. 太阳码
	sunCodeLocalHandler := &handler.ImageCircleRemoteHandler{
		X:      279 * multiple, // 图片x坐标
		Y:      508 * multiple, // 图片y坐标
		Width:  75 * multiple,
		Height: 75 * multiple,
		URL:    req.SunCode,
	}

	// 12.长按识别抢购
	//purchaseHandler := &handler.TextHandler{
	//	Next:     handler.Next{},
	//	X:        150 * multiple,         // 文字x坐标
	//	Y:        615 * multiple,         // 文字y坐标
	//	Size:     float64(18 * multiple), // 文字大小
	//	R:        127,                    // 文字颜色RGB值
	//	G:        127,
	//	B:        127,
	//	Text:     "时间不等人 长按识别抢购", // 文字内容
	//	FontPath: fontPath,       // 字体文件
	//}
	//结束绘制，把前面的内容合并成一张图片,输出到build目录
	fileName := "/poster_" + core.RandString(10) + ".png"
	outputUrl := rootpath + "cacheimg" + fileName
	endHandler := &handler.EndHandler{
		Output: outputUrl,
	}
	// 链式调用绘制过程
	nullHandler.
		SetNext(backgroundHandler).
		SetNext(imageCircleLocalHandler).
		SetNext(sunCodeLocalHandler).
		SetNext(skuImgHandler).
		//SetNext(marksHandler).
		SetNext(nameHandler).
		SetNext(recommendHandler).
		//SetNext(iconHandler).
		SetNext(goodsNameHandler).
		SetNext(sgoodsNameHandler).
		SetNext(priceMarkHandler).
		SetNext(sellPriceHandler).
		SetNext(sellHandler).
		SetNext(markerPriceHandler).
		SetNext(lineHandler).
		//SetNext(purchaseHandler).
		SetNext(endHandler)

	// 开始执行业务
	wg := sync.WaitGroup{}
	if err := nullHandler.Run(des, &wg); err != nil {
		// 异常
		log.Info(ctx.Wrap(c)).Msgf("海报生成错误：%+v", err)
		return poster
	}
	req.FileTmpUrl = outputUrl
	req.FileName = fileName
	url, err := changeUrlImage(c, req)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("changeUrlImage Error: %+v", err)
		return poster
	}
	removeTmpImgFile(c, outputUrl)
	//成功
	posterByte, _ := json.Marshal(&GetPosterResponse{
		Object: url,
	})
	err = redis.Client().Set(c, posterKey, posterByte, EnterpriseExpiration).Err()
	if err != nil {
		return &GetPosterResponse{}
	}
	return &GetPosterResponse{
		Object: url,
	}
}

// 生成优惠券海报
func (s Service) CouponPoster(c ctx.Context, req GetPosterRequest) *GetPosterResponse {
	couponParams := &CouponParams{}
	json.Unmarshal([]byte(req.Params), couponParams)
	log.Info(ctx.Wrap(c)).Msgf("请求数据:%+v,CouponPoster数据：%+v", req, couponParams)
	// 获取缓存
	posterKey := generateKey(req)
	posterCache := redis.Client().Get(c, posterKey).Val()
	poster := &GetPosterResponse{}
	if posterCache != "" {
		jErr := json.Unmarshal([]byte(posterCache), poster)
		if jErr != nil {
			return poster
		}
		if poster.Object != "" {
			return &GetPosterResponse{
				Object: poster.Object,
			}
		}
	}

	// 获取小程序的logo
	appInfo, err := getAppInfo(c, req.EnterpriseID)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("获取小程序的logo错误：%+v", err)
		return poster
	}
	couponParams.LogoUrl = appInfo.Data.Logo

	nullHandler := &handler.NullHandler{}
	multiple := 2
	des := &handler.Context{
		//图片都绘在这个PNG载体上
		PngCarrier: core.NewPNG(0, 0, 375*multiple, 575*multiple),
	}
	// 0.绘制背景图
	backgroundHandler := &handler.BackgroundHandler{
		X:    0,                                       // 图片x坐标
		Y:    0,                                       // 图片y坐标
		Path: rootpath + "sourceImages/coupon@2x.png", //图片路径
	}
	// 1.小程序logo
	imageCircleHandler := &handler.ImageCircleRemoteHandler{
		X:      33 * multiple,  // 图片x坐标
		Y:      469 * multiple, // 图片y坐标
		Width:  55 * multiple,
		Height: 55 * multiple,
		URL:    couponParams.LogoUrl,
	}

	// 2.太阳码
	sunCodeHandler := &handler.ImageCircleRemoteHandler{
		X:      265 * multiple, // 图片x坐标
		Y:      455 * multiple, // 图片y坐标
		Width:  84 * multiple,
		Height: 84 * multiple,
		URL:    req.SunCode,
	}
	// 3. 优惠券面值
	valueCount := fmt.Sprintf("%.2f", float64(couponParams.Value)/float64(100))
	valueLength := utf8.RuneCountInString(valueCount)
	ValueHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        ((375 * multiple) - (valueLength * 30 * multiple)) / 2,
		Y:        270 * multiple,
		Size:     float64(48 * multiple),
		R:        253,
		G:        246,
		B:        234,
		Text:     "¥" + valueCount,
		FontPath: fontPath,
	}
	// 4.使用期限
	timeTemplate3 := "2006.01.02"
	startTime := time.Unix(couponParams.StartTime, 0).Format(timeTemplate3)
	endTime := time.Unix(couponParams.EndTime, 0).Format(timeTemplate3)
	timeStr := startTime + "-" + endTime
	//timeCount := utf8.RuneCountInString(timeStr)
	TimeHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        57 * multiple,
		Y:        400 * multiple,
		Size:     float64(24 * multiple),
		R:        240,
		G:        213,
		B:        160,
		Text:     timeStr,
		FontPath: fontPath,
	}
	// 条件
	CouponMini := fmt.Sprintf("%.2f", float64(couponParams.CouponMini)/float64(100))
	couponParams.Condition = "满" + CouponMini + "减" + valueCount
	condCount := utf8.RuneCountInString(couponParams.Condition)
	CondHandler := &handler.TextHandler{
		Next:     handler.Next{},
		X:        ((375 * multiple) - (condCount * 16 * multiple)) / 2,
		Y:        355 * multiple,
		Size:     float64(28 * multiple),
		R:        240,
		G:        213,
		B:        160,
		Text:     couponParams.Condition,
		FontPath: fontPath,
	}
	//结束绘制，把前面的内容合并成一张图片,输出到build目录
	fileName := "/poster_" + core.RandString(10) + ".png"
	outputUrl := rootpath + "cacheimg" + fileName
	endHandler := &handler.EndHandler{
		Output: outputUrl,
	}
	// 链式调用绘制过程

	nullHandler.
		SetNext(backgroundHandler).
		SetNext(imageCircleHandler).
		SetNext(sunCodeHandler).
		SetNext(ValueHandler).
		SetNext(TimeHandler).
		SetNext(CondHandler).
		SetNext(endHandler)

	// 开始执行业务
	wg := sync.WaitGroup{}
	if runErr := nullHandler.Run(des, &wg); runErr != nil {
		// 异常
		log.Info(ctx.Wrap(c)).Msgf("海报生成错误：%+v", runErr)
		return &GetPosterResponse{}
	}
	log.Info(ctx.Wrap(c)).Msgf("海报生成成功")

	// 获取海报对象
	req.FileTmpUrl = outputUrl
	req.FileName = fileName
	url, changeErr := changeUrlImage(c, req)
	if changeErr != nil {
		log.Info(ctx.Wrap(c)).Msgf("changeUrlImage Error: %+v", changeErr)
		return poster
	}
	// 上传腾讯云
	cloudUrl, uploadErr := uploadCloud(c, cos.COS, outputUrl, req.ShareImg)
	if uploadErr != nil {
		// 异常
		return &GetPosterResponse{}
	}
	removeTmpImgFile(c, outputUrl)
	posterByte, _ := json.Marshal(&GetPosterResponse{
		Object: url,
	})
	err = redis.Client().Set(c, posterKey, posterByte, EnterpriseExpiration).Err()
	if req.IsTest == 1 {
		url = cloudUrl
	}
	log.Info(ctx.Wrap(c)).Msgf("ShareImg:%+v,cloudUrl: %+v", req.ShareImg, cloudUrl)
	//成功
	return &GetPosterResponse{
		Object: url,
	}
}

// getWxQrcode 获取太阳码
func getWxQrcode(c ctx.Context, request WxQrcodeRequest) (*microauthority.WxQrcodeResponse, error) {
	response, err := microauthority.WxQrcode(c, &request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("getWxQrcode Error: %+v", err)
		return &microauthority.WxQrcodeResponse{}, err
	}
	return response, nil
}

func generateKey(redisReq interface{}) string {
	keyData, _ := json.Marshal(redisReq)
	m := md5.New()
	m.Write([]byte(keyData))
	redisKey := fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
	return redisKey
}

// changeUrlImage 换取资源或mediaid
func changeUrlImage(c ctx.Context, request GetPosterRequest) (string, error) {
	if request.OperatorType == OperatorTypeSend {
		enterpriseService := enterprise.New(enterpriseRespository.New())
		cacheMediaRequest := work.CacheMediaWorkRequest{
			FileName:     request.FileName,
			FileType:     ImgFileType,
			FileOrg:      request.FileTmpUrl,
			EnterpriseID: request.EnterpriseID,
		}
		corpRequest := enterprise.Corp{
			EnterpriseID: request.EnterpriseID,
		}
		corp, err := enterpriseService.Detail(c, corpRequest)
		if err != nil {
			log.Info(ctx.Wrap(c)).Msgf("企业查询失败,Error: %+v", err)
			return "", err
		}
		workImpl := new(work.QyWorkImpl).SetDomain(corp.Domain)
		mediaInfo, err := workImpl.CacheMedia(c, cacheMediaRequest)
		if err != nil {
			log.Info(ctx.Wrap(c)).Msgf("Error: %+v", err)
			return "", err
		}
		return mediaInfo.Data.MediaId, nil
	}
	return imagesToBase64(request.FileTmpUrl), nil
}

// imagesToBase64 图片转base64
func imagesToBase64(strImages string) string {
	//读原图片
	ff, _ := os.Open(strImages)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return sourcestring
}

// 上传云存储 eg: file.UploadCloud(*c, storage.COS, tmpFilePath, "/uploads")
func uploadCloud(c ctx.Context, pkg string, tmpFilePath string, toFilePath string) (string, error) {
	cs := storage.Init(pkg)
	path, err := cs.Upload(c, tmpFilePath, toFilePath)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("uploadCloud Error: %+v", err)
		return "", err
	}
	return cs.GetURL(path), err
}

// removeTmpImgFile 删除临时文件
func removeTmpImgFile(c ctx.Context, file string) {
	log.Info(c).Msgf("【===删除临时文件===】")
	os.Remove(file)
}

// 获取小程序授权信息
func getAppInfo(c ctx.Context, enterpriseID int64) (*boss.AuthDetailResponse, error) {
	// 获取小程序的logo+name
	appReq := &boss.AuthDetailRequest{
		EnterpriseID: enterpriseID,
		AuthType:     1,
	}
	appInfoKey := generateKey(appReq)
	appInfoCache := redis.Client().Get(c, appInfoKey).Val()
	if appInfoCache != "" {
		appInfo := &boss.AuthDetailResponse{}
		jErr := json.Unmarshal([]byte(appInfoCache), appInfo)
		if jErr != nil {
			return nil, jErr
		}
		return appInfo, nil
	} else {
		appInfo, err := boss.AuthDetail(c, appReq)
		if err != nil || appInfo.Data.Logo == "" || appInfo.Data.Name == "" {
			return nil, err
		}
		appInfoByte, _ := json.Marshal(appInfo)
		redisErr := redis.Client().Set(c, appInfoKey, appInfoByte, EnterpriseExpiration).Err()
		if redisErr != nil {
			return nil, redisErr
		}
		return appInfo, nil
	}
}

// 获取太阳码
func (s Service) getSunCode(c ctx.Context, domain string, req GetPosterRequest) (*microauthority.WxQrcodeResponse, error) {
	// 获取太阳码携带参数
	qrcodeParams, err := getCodeParams(ctx.Wrap(c), domain, req)
	if err != nil {
		return &microauthority.WxQrcodeResponse{}, err
	}
	qrcodeReq := WxQrcodeRequest{
		EnterpriseID: req.EnterpriseID,
		Path:         MarketingPath[req.Type],
		Params:       qrcodeParams,
	}
	if req.Type == MarketingCoupon {
		qrcodeParams.SetShareImg("/uploads/" + common.GetDate("/") + "/poster_" + core.RandString(10) + ".png")
		req.ShareImg = qrcodeParams.GetShareImg()
	}
	// 获取缓存中的太阳码
	sunCodeKey := generateKey(qrcodeReq.Params)
	sunCodeCache := redis.Client().Get(c, sunCodeKey).Val()
	if sunCodeCache != "" {
		sunCodeResp := &microauthority.WxQrcodeResponse{}
		jErr := json.Unmarshal([]byte(sunCodeCache), sunCodeResp)
		if jErr != nil {
			return sunCodeResp, jErr
		}
		return sunCodeResp, nil
	}
	// 重新生成太阳码
	sunCodeResp, codeErr := getWxQrcode(ctx.Wrap(c), qrcodeReq)

	sunCodeResp.Data.ShareImg = req.ShareImg
	if codeErr != nil || sunCodeResp.Data.Url == "" {
		log.Info(ctx.New()).Msgf("重新生成太阳码错误:%+v,请求参数:%+v, 携带参数:%+v", codeErr, qrcodeReq, qrcodeReq.Params)
		return sunCodeResp, codeErr
	}
	// 写入缓存K

	sunCodeByte, _ := json.Marshal(sunCodeResp)
	err = redis.Client().Set(c, sunCodeKey, sunCodeByte, EnterpriseExpiration).Err()
	if err != nil {
		log.Info(ctx.New()).Msgf("写入缓存失败:%+v", err)
		return sunCodeResp, err
	}
	return sunCodeResp, err
}
