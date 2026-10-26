package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/modules/drainage"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk"
)

// ChannelsCreate 渠道新增
func ChannelsCreate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ChannelCreateUrl), request)
	return response, err
}

// ChannelsList 渠道列表
func ChannelsList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ChannelListUrl), request)
	return response, err
}

// ChannelsDetail 渠道详情
func ChannelsDetail(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ChannelDetailUrl), request)
	return response, err
}

// ChannelsEdit 渠道编辑
func ChannelsEdit(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ChannelEditUrl), request)
	return response, err
}

//联系我列表
func ContactList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ContactListUrl), request)
	return response, err
}

//联系我新增
func ContactCreate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ContactCreateUrl), request)
	return response, err
}

//联系我修改
func ContactUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ContactUpdateUrl), request)
	return response, err
}

//联系我删除
func ContactDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ContactDeleteUrl), request)
	return response, err
}

//联系我详情
func ContactDetails(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ContactDetailsUrl), request)
	return response, err
}

//下载联系我二维码
func DownloadQrCode(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.DownloadQrCodeUrl), request)
	return response, err
}

// welcomeCreate 欢迎语新建
func WelcomeCreate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeCreateUrl), request)
	return response, err
}

// WelcomeDetail 欢迎语详情
func WelcomeDetail(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeDetailUrl), request)
	return response, err
}

// WelcomeDefault 欢迎语默认详情
func WelcomeDefault(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeDefaultUrl), request)
	return response, err
}

// WelcomeEdit 欢迎语编辑
func WelcomeEdit(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeEditUrl), request)
	return response, err
}

// WelcomeDelete 欢迎语删除
func WelcomeDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeDeleteUrl), request)
	return response, err
}

// WelcomeList 欢迎语列表
func WelcomeList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeListUrl), request)
	return response, err
}

// WelcomeChannelList 用户欢迎语渠道列表
func WelcomeChannelList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.WelcomeValidateUrl), request)
	return response, err
}

//渠道统计导出
func ChannelsCountData(c ctx.Context, request interface{}) (*drainage.ChannelsCountData, error) {
	var response drainage.ChannelsCountData
	err := client.Post(c, GetUrl(urlEnum.ChannelsDataUrl), request, &response)
	return &response, err
}

// ChannelsData 渠道数据统计
func ChannelsData(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ChannelsDataUrl), request)
	return response, err
}

// CacheMedia 缓存图片素材
func CacheMedia(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CacheMediaUrl), request)
	return response, err
}

// StatisticsItemByDay 渠道统计总览
func StatisticsAllBySum(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.StatisticsAllBySumUrl), request)
	return response, err
}

// StatisticsChannelFacades 渠道统计门面
func StatisticsChannelFacades(ctx ctx.Context, params *drainage.ChannelsDataRequest) (interface{}, error) {
	var url string
	switch {
	case params.ChannelsId == 0 && params.Type == 1:
		// 渠道统计按天展示
		url = urlEnum.StatisticsAllByDayUrl
	case params.ChannelsId != 0 && params.Type == 1:
		// 渠道统计单一渠道按天展示
		url = urlEnum.StatisticsItemByDayUrl
	case params.Type == 2:
		// 渠道统计渠道列表总览
		url = urlEnum.StatisticsAllByChannelUrl
	default:
		err := errorx.NewErrorMsg(scrm_errs.DefaultErrorCode, scrm_errs.DefaultErrorMsg)
		return nil, err
	}

	var response drainage.ChannelsCountData
	err := client.Post(ctx, GetUrl(url), params, &response)
	return response.Data, err
}

func WelcomeNewEdit(ctx ctx.Context, params *drainage.WelcomeNewEditRequest) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(urlEnum.WelcomeNewEditConfigUrl), params)
	return response, err
}

func WelcomeDetailNew(ctx ctx.Context, params *drainage.WelcomeBase) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(urlEnum.WelcomeDetailNewConfigUrl), params)
	return response, err
}

func DelWelcome(ctx ctx.Context, params *drainage.WelcomeBase) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(urlEnum.DelWelcomeConfigUrl), params)
	return response, err
}

func GetAllWelcome(ctx ctx.Context, params *drainage.WelcomeAllRequest) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(urlEnum.GetAllWelcomeConfigUrl), params)
	return response, err
}

func GetAllKd(ctx ctx.Context, params drainage.KdTypaAll) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(urlEnum.GetAllKDConfigUrl), params)
	return response, err
}

func GetSimpleChannelList(ctx ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(urlEnum.DrainageSimpleChannelListUrl), req)
}

func ExternalWelcomeTypeList(ctx ctx.Context, req interface{}) (interface{}, error) {
	return sdk.SendPost(ctx, GetUrl(urlEnum.ExternalWelcomeTypeUrl), req)
}
