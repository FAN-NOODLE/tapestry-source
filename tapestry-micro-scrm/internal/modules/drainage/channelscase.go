package drainage

import (
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/uid"
)

// ChannelsCreate 创建渠道
func (s *DrainageService) ChannelsCreate(c ctx.Context, request *CreateContactChannelLabelsRequest) (*WelcomeCreateResponse, error) {
	if request.Name == "" {
		return nil, errs.ChannelNameNotNllError
	}
	// 查询渠道是否已经存在
	req := ContactChannels{
		Name:         request.Name,
		EnterpriseID: request.EnterpriseID,
	}
	detail, _ := s.repo.ChannelDetail(c, req)
	if detail != nil {
		return nil, errs.ChannelDetailExists
	}
	// 添加渠道数据
	ContactChannelsData := ContactChannels{
		ID:           uid.ID(ctx.New()),
		EnterpriseID: request.EnterpriseID,
		Name:         request.Name,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	_, channelsCreateErr := s.repo.ChannelsCreate(c, ContactChannelsData)

	if channelsCreateErr != nil {
		return nil, channelsCreateErr
	}
	if len(request.Label) > 0 {
		// 保存渠道标签
		for k := range request.Label {
			request.Label[k].ID = uid.ID(ctx.New())
			request.Label[k].CreatedAt = time.Now().Unix()
			request.Label[k].UpdatedAt = time.Now().Unix()
			request.Label[k].EnterpriseID = request.EnterpriseID
			request.Label[k].ChannelID = ContactChannelsData.ID

			var labelInfo []ContactChannelLabels
			labelInfo = append(labelInfo, request.Label[k])
			cErr := s.repo.ChannelsLabelCreate(c, labelInfo)
			if cErr != nil {
				return nil, cErr
			}
		}
	}

	response := &WelcomeCreateResponse{}
	response.Result = true
	return response, nil
}

// ChannelsList 渠道列表
func (s *DrainageService) ChannelsList(c ctx.Context, request ContactChannelListRequest, tagService *tag.Service) (int64, *[]ContactChannelsList, error) {
	total, list, err := s.repo.ChannelsList(c, request)
	if err != nil {
		return 0, nil, err
	}
	if total <= 0 {
		return 0, nil, nil
	}
	// 获取到渠道ID
	var channelIds []int64
	for _, res := range *list {
		channelIds = append(channelIds, res.ID)
	}
	// 查询渠道标签数据
	req := ContactChannelLabelsListRequest{
		EnterpriseID:    request.EnterpriseID,
		ChannelIds:      channelIds,
		ChannelIdsCount: len(channelIds),
	}
	_, labels, cErr := s.repo.ChannelsLabelsList(c, req)
	if cErr != nil {
		return 0, nil, cErr
	}

	if len(*labels) > 0 {
		labelMap := make(map[int64][]string)
		// 再组装渠道数据
		for _, vItem := range *labels {
			labelMap[vItem.ChannelID] = append(labelMap[vItem.ChannelID], vItem.LabelName)
		}
		for k, item := range *list {
			temp := make([]string, 0)
			if labelItem, ok := labelMap[item.ID]; ok {
				temp = labelItem
			}
			(*list)[k].ChannelLabels = temp
		}
	}
	return total, list, nil
}

// ChannelDetail 渠道详情
func (s *DrainageService) ChannelDetail(c ctx.Context, req *ContactChannels) (*ContactChannels, error) {
	detail, err := s.repo.ChannelDetail(c, *req)
	if err != nil {
		return nil, err
	}
	// 查询渠道绑定的标签
	reqLabel := ContactChannelLabelsListRequest{
		ChannelID:    detail.ID,
		EnterpriseID: req.EnterpriseID,
	}
	_, list, lErr := s.repo.ChannelsLabelsList(c, reqLabel)
	if lErr != nil {
		return nil, lErr
	}
	detail.Labels = *list
	return detail, nil
}

// ChannelEdit 渠道编辑
func (s *DrainageService) ChannelEdit(ctx ctx.Context, req CreateContactChannelLabelsRequest) error {
	return s.repo.ChannelEdit(ctx, req)
}

// SimpleChannelsList 简单渠道列表
func (s *DrainageService) SimpleChannelsList(c ctx.Context, request ContactChannelListRequest) (int64, *[]ContactChannelsList, error) {
	return s.repo.ChannelsList(c, request)
}

// SimpleChannelsDetail 简单渠道详情
func (s *DrainageService) SimpleChannelsDetail(c ctx.Context, request ContactChannels) (*ContactChannels, error) {
	return s.repo.ChannelDetail(c, request)
}
