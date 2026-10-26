package drainagehttp

import (
	"encoding/json"
	"micro-scrm/internal/pkg/helper"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
)

var drainageService *drainage.DrainageService
var tagService *tag.Service

// welcomeCreate 欢迎语新建
func welcomeCreate(c *gin.Context) {
	request := &drainage.WelcomeCreateRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	// 非默认欢迎语
	if request.IsDefault != 1 {
		// 校验此渠道是否已经添加过欢迎语
		req := &drainage.WelcomeListRequest{
			EnterpriseId: request.EnterpriseId,
			ChannelId:    request.ChannelId,
		}
		total, cErr := drainageService.CountWelcomeByChannelID(ctx.Wrap(c), req)
		if cErr != nil {
			ginplus.ResError(c, cErr)
			return
		}
		if total > 0 {
			ginplus.ResError(c, errs.WelcomeChannelDetailExists)
			return
		}
	}

	response, queryError := drainageService.WelcomeCreate(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeCreateError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeDetail 欢迎语详情
func welcomeDetail(c *gin.Context) {
	request := &drainage.WelcomeDetailRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	response, queryError := drainageService.WelcomeDetail(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeDetailError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeDefault 欢迎语默认详情
func welcomeDefault(c *gin.Context) {
	request := &drainage.WelcomeDetailRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	welcome, err := drainageService.WelcomeDefault(ctx.Wrap(c), *request)
	if err != nil {
		ginplus.ResError(c, errs.WelcomeDetailError, http.StatusOK)
		return
	}
	contentJson := make([]map[string]interface{}, 0)
	jsonError = json.Unmarshal([]byte(welcome.Content), &contentJson)
	if jsonError != nil {
		ginplus.ResError(c, errs.WelcomeDetailError, http.StatusOK)
		return
	}

	response := &drainage.WelcomeDetailResponse{
		ID:           welcome.ID,
		Content:      contentJson,
		EnterpriseId: welcome.EnterpriseId,
		IsDefault:    welcome.IsDefault,
		CreatedAt:    welcome.CreatedAt,
		UpdatedAt:    welcome.UpdatedAt,
		Channel:      "",
		ChannelId:    0,
		Page:         1,
		Size:         1,
		Total:        1,
	}

	ginplus.ResSuccess(c, response)
}

// welcomeEdit 欢迎语编辑
func welcomeEdit(c *gin.Context) {
	request := &drainage.WelcomeEditRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	if request.IsDefault != 1 {
		// 校验此渠道是否已经添加过欢迎语
		req := &drainage.WelcomeListRequest{
			EnterpriseId: request.EnterpriseId,
			ChannelId:    request.ChannelId,
			IdNot:        request.OldID,
		}
		total, cErr := drainageService.CountWelcomeByChannelID(ctx.Wrap(c), req)
		if cErr != nil {
			ginplus.ResError(c, cErr)
			return
		}
		if total > 0 {
			ginplus.ResError(c, errs.WelcomeChannelDetailExists)
			return
		}
	}
	response, queryError := drainageService.WelcomeEdit(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeEditError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
	return
}

// welcomeDelete 删除欢迎语
func welcomeDelete(c *gin.Context) {
	request := &drainage.WelcomeDeleteRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	response, queryError := drainageService.WelcomeDelete(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeDeleteError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeList 欢迎语列表
func welcomeList(c *gin.Context) {
	request := &drainage.WelcomeListRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	response, queryError := drainageService.WelcomeList(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeListError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeSend 回调发送欢迎语
func welcomeSend(c *gin.Context) {
	request := &drainage.WelcomeSendRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	response, queryError := drainageService.WelcomeSend(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeSendError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// welcomeChannelList 查询显示的渠道列表
func welcomeChannelList(c *gin.Context) {
	request := &drainage.WelcomeValidateRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	response, queryError := drainageService.WelcomeChannelList(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeValidateError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// channelsList 渠道列表
func channelsList(c *gin.Context) {
	request := &drainage.ContactChannelListRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	total, list, queryError := drainageService.ChannelsList(ctx.Wrap(c), *request, tagService)
	if queryError != nil {
		ginplus.ResError(c, errs.ChannelListError, http.StatusOK)
		return
	}
	response.SuccessList(c, total, request.Page, request.Size, list)
}

// channelSimpleList 获取渠道总列表数据
func channelSimpleList(c *gin.Context) {
	request := drainage.ContactChannelListRequest{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	request.Size = 0
	_, list, _ := drainageService.SimpleChannelsList(ctx.Wrap(c), request)

	var resp []*drainage.SimpleChannelList
	for _, v := range *list {
		resp = append(resp, &drainage.SimpleChannelList{ID: v.ID, Name: v.Name})
	}

	ginplus.ResSuccess(c, resp)
	return
}

// channelsCreate 渠道创建
func channelsCreate(c *gin.Context) {
	request := &drainage.CreateContactChannelLabelsRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	responseData, queryError := drainageService.ChannelsCreate(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, queryError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, responseData)
}

// channelsDetail 渠道详情
func channelsDetail(c *gin.Context) {
	request := &drainage.ContactChannels{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	detail, queryError := drainageService.ChannelDetail(ctx.Wrap(c), request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeCreateError, http.StatusOK)
		return
	}
	// 处理标签名称
	if len(detail.Labels) > 0 {
		var tagIds []int64
		for _, res := range detail.Labels {
			tagIds = append(tagIds, res.LabelID)
		}
		tags, err := tagService.QueryWorkExternalTags(ctx.Wrap(c), tagIds)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		tagMap := make(map[string]entity.WorkExternalTag)
		for _, tagItem := range tags {
			tagMap[tagItem.TagID] = tagItem
		}
		// 组装标签数据
		for k, item := range detail.Labels {
			if tagExists, ok := tagMap[strconv.FormatInt(item.LabelID, 10)]; ok {
				detail.Labels[k].LabelName = tagExists.TagName
			}
		}
	}

	ginplus.ResSuccess(c, detail)
}

// channelsEdit 渠道编辑
func channelsEdit(c *gin.Context) {
	request := &drainage.CreateContactChannelLabelsRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	if request.ID == 0 {
		ginplus.ResError(c, errs.ChannelNameNotNllError)
		return
	}
	queryError := drainageService.ChannelEdit(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, queryError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, nil)
}

func channelsData(c *gin.Context) {
	request := &drainage.ChannelsDataRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	response, queryError := drainageService.ChannelsData(ctx.Wrap(c), *request)
	if queryError != nil {
		ginplus.ResError(c, errs.WelcomeCreateError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// cacheMedia 缓存微信图片素材
func cacheMedia(c *gin.Context) {
	request := &drainage.CacheMediaRequest{}
	jsonError := c.BindJSON(request)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	domain, _, _, _ := enterpriseService.GetInfo(ctx.Wrap(c), request.EnterpriseId)
	response, queryError := drainageService.CacheMedia(ctx.Wrap(c), *request, domain)
	if queryError != nil {
		ginplus.ResError(c, queryError, http.StatusOK)
		return
	}
	ginplus.ResSuccess(c, response)
}

// Enterprise
type Enterprise struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" desc:"企业id"`
}

// synContact 同步渠道联系我二维码
func synContact(c *gin.Context) {
	var request Enterprise
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	err = drainageService.SynContact(ctx.Wrap(c), request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// statisticsAllByDay 时间范围+所有渠道,按天展示
func statisticsAllByDay(c *gin.Context) {
	request := &drainage.ChannelsDataRequest{}
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	// 获取所有的天数
	allDays := helper.GetBetweenDates(
		time.Unix(request.StartTime, 0).Format("2006-01-02"),
		time.Unix(request.EndTime, 0).Format("2006-01-02"),
		false,
	)
	allDaysLen := int64(len(allDays))

	cc := ctx.Wrap(c)
	var err error
	// 保存最后的结果集
	results := make([]*drainage.StatisticsItemResponse, 0, request.Size)

	pageStart := (request.Page - 1) * request.Size
	pageEnd := pageStart + request.Size
	if pageEnd > allDaysLen {
		pageEnd = allDaysLen
	}

	// 查询企业详情
	_, corpId, _, _ := enterpriseService.GetInfo(cc, request.EnterpriseID)
	log.Info(cc).Msgf("渠道统计-----statisticsAllByDay---request.StartTime----%v", request.StartTime)
	log.Info(cc).Msgf("渠道统计-----statisticsAllByDay---request.EndTime----%v", request.EndTime)
	log.Info(cc).Msgf("渠道统计-----statisticsAllByDay---allDays----%v", allDays)
	log.Info(cc).Msgf("渠道统计-----statisticsAllByDay---pageStart----%v", pageStart)
	log.Info(cc).Msgf("渠道统计-----statisticsAllByDay---pageEnd----%v", pageEnd)
	// 循环本页所有的天数
	for _, date := range allDays[pageStart:pageEnd] {
		data := new(drainage.StatisticsItemResponse)
		data.TimeOrName = date

		// 创建请求结构体
		start, end := helper.GetTimestamp(date)
		itemReq := &drainage.StatisticsAllReq{
			EnterpriseID: request.EnterpriseID,
			StartTime:    start,
			EndTime:      end - 1,
		}

		group := threading.NewRoutineGroup()

		// 留存好友数
		group.RunSafe(func() {
			data.ExistRelationCustomer, err = drainageService.ExistRelationCustomer(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("留存好友数查询失败：%v", err)
			}
		})

		// 新增好友总数
		group.RunSafe(func() {
			data.IncreaseCustomerCount, err = drainageService.IncreaseCustomerCount(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("新增好友总数查询失败：%v", err)
			}
		})

		// 流失好友数量
		group.RunSafe(func() {
			data.LossCustomerCount, err = drainageService.LossCustomerCount(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("流失好友数量查询失败：%v", err)
			}
		})

		// 好友删除员工数
		group.RunSafe(func() {
			data.CustomerDelMember, err = drainageService.CustomerDelMember(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("好友删除员工数查询失败：%v", err)
			}
		})

		// 员工删除客户总数
		group.RunSafe(func() {
			data.MemberDelCustomer, err = drainageService.MemberDelCustomerOrChannel(
				cc,
				&drainage.StatisticsChannelReq{
					StatisticsAllReq: *itemReq,
				},
				corpId,
			)

			if err != nil {
				log.Error(cc).Msgf("员工删除客户总数查询失败：%v", err)
			}
		})

		group.Wait()

		// 净增好友数
		data.NewlyAddedCustomerCount = data.IncreaseCustomerCount - data.LossCustomerCount
		results = append(results, data)
	}

	statisticsResponse := new(drainage.StatisticsResponse)
	statisticsResponse.Size = request.Size
	statisticsResponse.Page = request.Page
	statisticsResponse.List = results
	statisticsResponse.Total = allDaysLen

	ginplus.ResSuccess(c, statisticsResponse)
}

// statisticsAllByDay 时间范围+所有渠道,按渠道展示
func statisticsAllByChannel(c *gin.Context) {
	request := &drainage.ChannelsDataRequest{}
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	cc := ctx.Wrap(c)
	channelReq := drainage.ContactChannelListRequest{
		EnterpriseID: request.EnterpriseID,
		Page:         int(request.Page),
		Size:         int(request.Size),
		ChannelID:    request.ChannelsId,
	}
	total, list, err := drainageService.SimpleChannelsList(cc, channelReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 查询企业详情
	_, corpId, _, _ := enterpriseService.GetInfo(cc, request.EnterpriseID)

	// 查询所有渠道的数据
	results := make([]*drainage.StatisticsItemResponse, 0, request.Size)
	for _, channel := range *list {
		data := new(drainage.StatisticsItemResponse)
		// 创建请求结构体
		itemReq := new(drainage.StatisticsChannelReq)
		itemReq.EnterpriseID = request.EnterpriseID
		itemReq.StartTime = request.StartTime
		itemReq.EndTime = request.EndTime
		itemReq.ChannelsState = channel.ID

		data.TimeOrName = channel.Name
		group := threading.NewRoutineGroup()

		// 留存好友数
		group.RunSafe(func() {
			data.ExistRelationCustomer, err = drainageService.ExistRelationCustomerByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("留存好友数查询失败：%v", err)
			}
		})

		// 新增好友总数
		group.RunSafe(func() {
			data.IncreaseCustomerCount, err = drainageService.IncreaseCustomerCountByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("新增好友总数查询失败：%v", err)
			}
		})

		// 流失好友数量
		group.RunSafe(func() {
			data.LossCustomerCount, err = drainageService.LossCustomerCountByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("流失好友数量查询失败：%v", err)
			}
		})

		// 好友删除员工数
		group.RunSafe(func() {
			data.CustomerDelMember, err = drainageService.CustomerDelMemberByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("好友删除员工数查询失败：%v", err)
			}
		})

		// 员工删除客户总数
		group.RunSafe(func() {
			data.MemberDelCustomer, err = drainageService.MemberDelCustomerOrChannel(cc, itemReq, corpId)
			if err != nil {
				log.Error(cc).Msgf("员工删除客户总数查询失败：%v", err)
			}
		})

		group.Wait()
		// 净增好友数
		data.NewlyAddedCustomerCount = data.IncreaseCustomerCount - data.LossCustomerCount
		results = append(results, data)
	}

	statisticsResponse := drainage.StatisticsResponse{
		Page:  request.Page,
		Size:  request.Size,
		Total: total,
		List:  results,
	}
	ginplus.ResSuccess(c, statisticsResponse)
}

// statisticsAllByDay 时间范围+所有渠道,所有数据之和
func statisticsAllBySum(c *gin.Context) {
	request := &drainage.ChannelsDataRequest{}
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	cc := ctx.Wrap(c)
	var err error
	data := new(drainage.StatisticsItemResponse)
	// 创建请求结构体
	itemReq := &drainage.StatisticsAllReq{
		EnterpriseID: request.EnterpriseID,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
	}

	// 查询企业详情
	_, corpId, _, _ := enterpriseService.GetInfo(cc, request.EnterpriseID)

	group := threading.NewRoutineGroup()

	// 留存好友数
	group.RunSafe(func() {
		data.ExistRelationCustomer, err = drainageService.ExistRelationCustomer(cc, itemReq)
		if err != nil {
			log.Error(cc).Msgf("留存好友数错误：%v", err)
		}
	})

	// 新增好友总数
	group.RunSafe(func() {
		data.IncreaseCustomerCount, err = drainageService.IncreaseCustomerCount(cc, itemReq)
		if err != nil {
			log.Error(cc).Msgf("新增好友总数错误：%v", err)
		}
	})

	// 流失好友数量
	group.RunSafe(func() {
		data.LossCustomerCount, err = drainageService.LossCustomerCount(cc, itemReq)
		if err != nil {
			log.Error(cc).Msgf("流失好友数量错误：%v", err)
		}
	})

	// 好友删除员工数
	group.RunSafe(func() {
		data.CustomerDelMember, err = drainageService.CustomerDelMember(cc, itemReq)
		if err != nil {
			log.Error(cc).Msgf("好友删除员工数错误：%v", err)
		}
	})

	// 员工删除客户总数
	group.RunSafe(func() {
		data.MemberDelCustomer, err = drainageService.MemberDelCustomerOrChannel(
			cc,
			&drainage.StatisticsChannelReq{
				StatisticsAllReq: *itemReq,
			},
			corpId,
		)
		if err != nil {
			log.Error(cc).Msgf("员工删除客户总数错误：%v", err)
		}
	})
	group.Wait()

	// 净增好友数
	data.NewlyAddedCustomerCount = data.IncreaseCustomerCount - data.LossCustomerCount
	ginplus.ResSuccess(c, data)
}

// statisticsItemByDay 时间范围+单个渠道,按天展示
func statisticsItemByDay(c *gin.Context) {
	request := &drainage.ChannelsDataRequest{}
	if err := c.BindJSON(request); err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	// 获取所有的天数
	allDays := helper.GetBetweenDates(
		time.Unix(request.StartTime, 0).Format("2006-01-02"),
		time.Unix(request.EndTime, 0).Format("2006-01-02"),
		false,
	)
	allDaysLen := int64(len(allDays))

	cc := ctx.Wrap(c)
	var err error
	// 保存最后的结果集
	results := make([]*drainage.StatisticsItemResponse, 0, request.Size)

	pageStart := (request.Page - 1) * request.Size
	pageEnd := pageStart + request.Size
	if pageEnd > allDaysLen {
		pageEnd = allDaysLen
	}

	// 查询企业详情
	_, corpId, _, _ := enterpriseService.GetInfo(cc, request.EnterpriseID)

	// 循环本页所有的天数
	for _, date := range allDays[pageStart:pageEnd] {
		data := new(drainage.StatisticsItemResponse)
		data.TimeOrName = date

		// 创建请求结构体
		start, end := helper.GetTimestamp(date)
		itemReq := new(drainage.StatisticsChannelReq)
		itemReq.EnterpriseID = request.EnterpriseID
		itemReq.StartTime = start
		itemReq.EndTime = end - 1
		if request.ChannelsId != 0 {
			itemReq.ChannelsState = request.ChannelsId
		}

		group := threading.NewRoutineGroup()

		// 留存好友数
		group.RunSafe(func() {
			data.ExistRelationCustomer, err = drainageService.ExistRelationCustomerByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("留存好友数错误：%v", err)
			}
		})

		// 新增好友总数
		group.RunSafe(func() {
			data.IncreaseCustomerCount, err = drainageService.IncreaseCustomerCountByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("新增好友总数错误：%v", err)
			}
		})

		// 流失好友数量
		group.RunSafe(func() {
			data.LossCustomerCount, err = drainageService.LossCustomerCountByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("流失好友数量错误：%v", err)
			}
		})

		// 好友删除员工数
		group.RunSafe(func() {
			data.CustomerDelMember, err = drainageService.CustomerDelMemberByChannel(cc, itemReq)
			if err != nil {
				log.Error(cc).Msgf("好友删除员工数错误：%v", err)
			}
		})

		// 员工删除客户总数
		group.RunSafe(func() {
			data.MemberDelCustomer, err = drainageService.MemberDelCustomerOrChannel(cc, itemReq, corpId)
			if err != nil {
				log.Error(cc).Msgf("员工删除客户总数错误：%v", err)
			}
		})

		group.Wait()
		// 净增好友数
		data.NewlyAddedCustomerCount = data.IncreaseCustomerCount - data.LossCustomerCount
		results = append(results, data)
	}

	statisticsResponse := new(drainage.StatisticsResponse)
	statisticsResponse.Size = request.Size
	statisticsResponse.Page = request.Page
	statisticsResponse.List = results
	statisticsResponse.Total = allDaysLen

	ginplus.ResSuccess(c, statisticsResponse)
}
