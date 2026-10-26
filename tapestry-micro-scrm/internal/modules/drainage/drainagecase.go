package drainage

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"micro-scrm/internal/pkg/redis"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"micro-scrm/internal/pkg/common"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

const DefaultWelcome = 1
const NotDefaultWelcome = 0
const FilePath = "H:\\"
const FileExt = ".jpg"
const ImgFileType = "1"
const SystemType = 1          // 系统默认欢迎语
const FissionActivityType = 2 // 裂变活动欢迎语

type DrainageService struct {
	repo Repository
}

func (s *DrainageService) FixDrainageWelcomeContent(c ctx.Context, enterpriseId int64) error {
	err := s.repo.FixDrainageWelcomeContentRep(ctx.Wrap(c), enterpriseId)
	return err
}

// WelcomeCreate 欢迎语新建
func (s *DrainageService) WelcomeCreate(c ctx.Context, request WelcomeCreateRequest) (*WelcomeCreateResponse, error) {

	response := &WelcomeCreateResponse{
		Result: false,
	}

	currentTime := time.Now().Unix()

	contentJson, jsonErr := json.Marshal(request.Content)
	if jsonErr != nil {
		return response, errs.JsonParseErr
	}

	welcome := Welcome{
		ID:               uid.ID(c),
		Content:          string(contentJson),
		EnterpriseId:     request.EnterpriseId,
		IsDefault:        request.IsDefault,
		ChannelId:        request.ChannelId,
		Type:             DefaultWelcome,
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
		DeletedAt:        0,
		ConditionContent: "",
	}

	_, err := s.repo.Add(c, welcome)
	if err != nil {
		return response, err
	}

	if welcome.IsDefault == DefaultWelcome {
		response.Result = true
		return response, nil
	}

	response.Result = true
	return response, nil
}

// WelcomeDefault 欢迎语默认详情
func (s *DrainageService) WelcomeDefault(c ctx.Context, request WelcomeDetailRequest) (*Welcome, error) {
	welcome, err := s.repo.SelectWelcomeDefault(c, request)
	if err != nil {
		// 没有找到的话新增一条
		currentNow := time.Now().Unix()
		welcomeDetailResponse := &WelcomeDetailResponse{}
		welcomeDetailResponse.ID = uid.ID(c)
		welcomeDetailResponse.Content = make([]map[string]interface{}, 0)
		tmpMap := make(map[string]interface{})
		tmpMap["text"] = map[string]interface{}{
			"content": "您好，很高兴认识您",
		}
		welcomeDetailResponse.Content = append(welcomeDetailResponse.Content, tmpMap)
		welcomeDetailResponse.EnterpriseId = request.EnterpriseId
		welcomeDetailResponse.IsDefault = DefaultWelcome
		welcomeDetailResponse.CreatedAt = currentNow
		welcomeDetailResponse.UpdatedAt = currentNow
		welcomeDetailResponse.DeletedAt = 0
		welcomeDetailResponse.Page = 1
		welcomeDetailResponse.Total = 1
		welcomeDetailResponse.Size = 10
		welcomeDetailResponse.ChannelId = 0
		contentJson, jsonErr := json.Marshal(welcomeDetailResponse.Content)
		if jsonErr != nil {
			return nil, errs.JsonParseErr
		}
		welcome = &Welcome{
			ID:           uid.ID(c),
			Content:      string(contentJson),
			EnterpriseId: request.EnterpriseId,
			IsDefault:    DefaultWelcome,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			DeletedAt:    0,
		}
		s.repo.Add(c, welcome)
	}

	return welcome, nil
}

// WelcomeDetail 欢迎语详情
func (s *DrainageService) WelcomeDetail(c ctx.Context, request WelcomeDetailRequest) (*WelcomeDetailResponse, error) {
	log.Info(ctx.Wrap(c)).Msgf("WelcomeDetail-request %+v", request)
	welcome, err := s.repo.SelectWelcomeOne(c, request)
	if err != nil {
		return nil, err
	}
	log.Info(ctx.Wrap(c)).Msgf("WelcomeDetail-SelectWelcomeOne %+v", welcome)
	response := &WelcomeDetailResponse{}
	contentJson := make([]map[string]interface{}, 0)
	jsonError := json.Unmarshal([]byte(welcome.Content), &contentJson)
	if jsonError != nil {
		log.Error(ctx.Wrap(c)).Msgf("WelcomeDetail-jsonError %+v", jsonError)
		return nil, errs.JsonParseErr
	}

	response.ID = welcome.ID
	response.Content = contentJson
	response.EnterpriseId = welcome.EnterpriseId
	response.IsDefault = welcome.IsDefault
	response.CreatedAt = welcome.CreatedAt
	response.UpdatedAt = welcome.UpdatedAt
	response.DeletedAt = welcome.DeletedAt
	response.ChannelId = welcome.ChannelId
	response.Channel = welcome.Channel

	return response, nil
}

// WelcomeEdit 编辑欢迎语
func (s *DrainageService) WelcomeEdit(c ctx.Context, request WelcomeEditRequest) (*WelcomeEditResponse, error) {

	response := &WelcomeEditResponse{
		Result: false,
	}

	currentTime := time.Now().Unix()
	contentJson, jsonErr := json.Marshal(request.Content)
	if jsonErr != nil {
		return nil, errs.JsonParseErr
	}
	welcomeCreateRequest := &Welcome{
		ID:           request.OldID,
		Content:      string(contentJson),
		EnterpriseId: request.EnterpriseId,
		IsDefault:    request.IsDefault,
		ChannelId:    request.ChannelId,
		UpdatedAt:    currentTime,
	}
	_, uErr := s.UpdateWelcomeMsg(ctx.Wrap(c), welcomeCreateRequest)
	if uErr != nil {
		return response, errs.WelcomeUpdateError
	}
	response.Result = true
	return response, nil
}

// WelcomeDelete 删除欢迎语
func (s *DrainageService) WelcomeDelete(c ctx.Context, request WelcomeDeleteRequest) (*WelcomeDeleteResponse, error) {
	response := &WelcomeDeleteResponse{
		Result: false,
	}

	request.DeletedAt = time.Now().Unix()

	_, err := s.repo.DeleteWelcomeOne(c, request)
	if err != nil {
		return response, errs.WelcomeDeleteError
	}

	response.Result = true
	return response, nil
}

// WelcomeList 欢迎语列表
func (s *DrainageService) WelcomeList(c ctx.Context, request WelcomeListRequest) (*WelcomeListResponse, error) {
	response := &WelcomeListResponse{
		Page:  request.Page,
		Size:  request.Size,
		Total: 0,
		List:  make([]WelcomeListStruct, 0),
	}

	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(request.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_drainage_welcomes/select-list-welcomes",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]WelcomeListStruct{}
		},
	}
	request.Type = SystemType
	total, list, err := s.repo.FindPage(c, pageVo, request, option)

	if err != nil || total == 0 {
		return response, nil
	}

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	resList := list.(*[]WelcomeListStruct)
	returnList := make([]WelcomeListStruct, 0)

	for _, res := range *resList {
		contentJson := make([]map[string]interface{}, 0)
		jsonError := json.Unmarshal([]byte(res.ContentStr), &contentJson)
		if jsonError != nil {
			return nil, errs.JsonParseErr
		}
		res.Content = contentJson
		res.ContentStr = ""
		returnList = append(returnList, res)
	}

	response.Total = total
	response.List = returnList

	return response, nil
}

// WelcomeSend 回调发送欢迎语
func (s *DrainageService) WelcomeSend(c ctx.Context, request WelcomeSendRequest) (*WelcomeSendResponse, error) {
	fmt.Println(request)
	response := &WelcomeSendResponse{}
	return response, nil
}

// WelcomeChannelList 查询渠道列表
func (s *DrainageService) WelcomeChannelList(c ctx.Context, request WelcomeValidateRequest) (*[]WelcomeChannelsList, error) {
	// 过滤欢迎语渠道，过滤裂变活动渠道
	var filter, activity, welcome []int64
	// 查询关联活动的渠道
	activityChannels, err := s.repo.GetActivityChannelIDs(c, request.EnterpriseId)
	for _, res := range activityChannels {
		activity = append(activity, res.ChannelID)
	}

	// 查询关联欢迎语的渠道
	welcomeChannels, err := s.repo.GetWelcomeChannelIDs(c, request.EnterpriseId)
	for _, res := range welcomeChannels {
		welcome = append(welcome, res.ChannelID)
	}
	welcome = append(welcome, activity...)
	filter = append(filter, welcome...)
	// 查询渠道数据表
	request.Filter = filter
	channelsList, err := s.repo.WelcomeChannelList(c, request)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	if channelsList == nil {
		channelsList = &[]WelcomeChannelsList{}
	}

	return channelsList, nil
}

// ChannelsData 查询渠道数据统计
func (s *DrainageService) ChannelsData(c ctx.Context, request ChannelsDataRequest) (*ChannelsDataResponse, error) {
	// 通过渠道查询当前渠道进入的用户信息
	WorkExternalUserRelation := WorkExternalUserRelation{
		EnterpriseID: request.EnterpriseID,
		State:        strconv.FormatInt(request.ChannelsId, 10),
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		ChannelId:    request.ChannelsId,
	}
	list, err := s.repo.GetWorkExternalUserRelationList(c, WorkExternalUserRelation)

	if err != nil {
		return nil, err
	}
	var response ChannelsDataResponse
	// 按日期返回渠道中统计数据
	if request.Type == 1 {
		response, err = s.DateChannelsData(c, request, list)
		if err != nil {
			return nil, err
		}
	}
	if request.Type == 2 {
		response, err = s.ChannelChannelsData(c, request, list)
		if err != nil {
			return nil, err
		}
	}
	// 分页处理
	total := len(response.ChannelsDataListResponse)
	size := request.Size
	page := request.Page
	min := size * (page - 1)
	max := page * size
	if max > int64(total) {
		max = int64(total)
	}
	if total == 1 {
		min = 0
	}
	// totalPage := math.Ceil(float64(total) / float64(size))
	// fmt.Println(totalPage)
	response.ChannelsDataListResponse = response.ChannelsDataListResponse[min:max]
	return &response, nil
}

// 按渠道处理统计数据
func (s *DrainageService) ChannelChannelsData(c ctx.Context, request ChannelsDataRequest, list *[]WorkExternalUserRelation) (ChannelsDataResponse, error) {
	var contactChannelsData = ContactChannelsData{
		EnterpriseID: request.EnterpriseID,
		ID:           request.ChannelsId,
	}
	var response ChannelsDataResponse
	// 查询渠道列表
	channelsList, err := s.repo.GetDrainageContactChannelsList(c, contactChannelsData)
	if err != nil {
		return response, err
	}

	for _, channelValue := range *channelsList {
		var mapping ChannelsDataListResponse
		var totalCustomer int64 = 0
		var lossCustomerCount int64 = 0
		var customerDelMember int64 = 0
		var memberDelCustomer int64 = 0

		for _, value := range *list {
			if value.ChannelsName == channelValue.Name {
				// 累计客户总数 净增客户总数
				totalCustomer++

				// 客户删除员工关系时间
				if value.ExternalDeletedTimeStr != "1970-01-01" && value.RelationStatus != 0 {
					customerDelMember++
				}
				// 员工删除客户关系时间
				if value.MemberDeletedTimeStr != "1970-01-01" && value.RelationStatus != 0 {
					memberDelCustomer++
				}

				// 累计流失客户数
				if (value.RelationStatus != 0 && value.ExternalDeletedTimeStr != "1970-01-01") || value.RelationStatus != 0 && value.MemberDeletedTimeStr != "1970-01-01" {
					lossCustomerCount++
				}
			}
		}

		mapping.TimeOrName = channelValue.Name
		mapping.TotalCustomer = totalCustomer
		mapping.NewlyAddedCustomerCount = totalCustomer
		mapping.LossCustomerCount = lossCustomerCount
		mapping.IncreaseCustomerCount = totalCustomer
		mapping.CustomerDelMember = customerDelMember
		mapping.MemberDelCustomer = memberDelCustomer
		response.ChannelsDataListResponse = append(response.ChannelsDataListResponse, &mapping)

		response.TotalCustomer += totalCustomer
		response.NewlyAddedCustomerCount += totalCustomer
		response.LossCustomerCount += lossCustomerCount
		response.IncreaseCustomerCount += totalCustomer
		response.CustomerDelMember += customerDelMember
		response.MemberDelCustomer += memberDelCustomer
	}
	response.Total = len(response.ChannelsDataListResponse)
	response.Size = request.Size
	response.Page = request.Page

	return response, nil
}

// 按时间处理统计数据
func (s *DrainageService) DateChannelsData(c ctx.Context, request ChannelsDataRequest, list *[]WorkExternalUserRelation) (ChannelsDataResponse, error) {
	var timeIndex = request.StartTime
	var endTime = request.EndTime
	var response ChannelsDataResponse

	timeLayout := "2006-01-02" // 转化所需模板
	for {
		if timeIndex > endTime {
			break
		}
		var mapping ChannelsDataListResponse
		var totalCustomer int64 = 0
		var lossCustomerCount int64 = 0
		var customerDelMember int64 = 0
		var memberDelCustomer int64 = 0
		mapping.TimeOrName = time.Unix(timeIndex, 0).Format(timeLayout)
		mapping.TotalCustomer = 1
		for _, value := range *list {
			// 客户删除员工关系时间
			if value.ExternalDeletedTimeStr == mapping.TimeOrName && value.RelationStatus != 0 {
				customerDelMember++
			}
			// 员工删除客户关系时间
			if value.MemberDeletedTimeStr == mapping.TimeOrName && value.RelationStatus != 0 {
				memberDelCustomer++
			}

			if value.Time == mapping.TimeOrName {
				// 累计客户总数 净增客户总数
				totalCustomer++

				// 累计流失客户数
				if value.RelationStatus != 0 && value.ExternalDeletedTimeStr == mapping.TimeOrName || value.MemberDeletedTimeStr == mapping.TimeOrName {
					lossCustomerCount++
				}
			}
		}

		mapping.TotalCustomer = totalCustomer
		mapping.NewlyAddedCustomerCount = totalCustomer
		mapping.LossCustomerCount = lossCustomerCount
		mapping.IncreaseCustomerCount = totalCustomer
		mapping.CustomerDelMember = customerDelMember
		mapping.MemberDelCustomer = memberDelCustomer
		response.ChannelsDataListResponse = append(response.ChannelsDataListResponse, &mapping)
		timeIndex += 86400

		response.TotalCustomer += totalCustomer
		response.NewlyAddedCustomerCount += totalCustomer
		response.LossCustomerCount += lossCustomerCount
		response.IncreaseCustomerCount += totalCustomer
		response.CustomerDelMember += customerDelMember
		response.MemberDelCustomer += memberDelCustomer
	}
	response.Total = len(response.ChannelsDataListResponse)
	response.Size = request.Size
	response.Page = request.Page

	return response, nil
}

// MakeMediaBySourceUrl 生成临时素材ID
func (s *DrainageService) MakeMediaBySourceUrl(ctx ctx.Context, enterpriseID int64, sourceUrl string, sourceType string, sourceName string, domain string) (string, error) {
	req := CacheMediaRequest{
		EnterpriseId: enterpriseID,
		ImgUrl:       sourceUrl,
		FileName:     sourceName,
	}
	if sourceType == "image" {
		req.FileType = "1"
	} else if sourceType == "video" {
		req.FileType = "3"
	} else if sourceType == "file" {
		req.FileType = "99"
	} else {
		return "", errs.SourceTypeErr
	}

	res, err := s.CacheMedia(ctx, req, domain)
	if err != nil {
		return "", err
	}
	return res.MediaId, nil
}

// CacheMedia 缓存微信图片素材
func (s *DrainageService) CacheMedia(c ctx.Context, request CacheMediaRequest, domain string) (*CacheMediaResponse, error) {
	log.Info(c).Msgf("【===开始缓存微信图片素材===】请求参数：%+v", request)
	response := &CacheMediaResponse{
		Result: false,
	}
	cacheKey := fmt.Sprintf("CacheMedia:%s", request.ImgUrl)
	imgMediaId := redis.Client().Get(c, cacheKey).Val()
	if imgMediaId != "" {
		response.Result = true
		response.MediaId = imgMediaId
		return response, nil
	}

	file, err := s.getTmpImgFile(c, request)
	if err != nil {
		log.Error(c).Msgf("将素材存储到本地失败，错误信息：%+v", err)
		return nil, err
	}
	// 获取图片名称
	if request.FileName == "" {
		ext := path.Ext(request.ImgUrl)
		// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
		if u, err := url.Parse(request.ImgUrl); err == nil && u.Path != "" {
			ext = path.Ext(u.Path)
		}
		request.FileName = strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	}
	cacheMediaRequest := work.CacheMediaWorkRequest{
		FileName:     request.FileName,
		FileType:     request.FileType,
		FileOrg:      file,
		EnterpriseID: request.EnterpriseId,
		ImgUrl:       request.ImgUrl,
	}
	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	mediaInfo, err := workImpl.CacheMedia(c, cacheMediaRequest)
	defer s.removeTmpImgFile(c, file)
	if err != nil {
		return nil, err
	}

	if mediaInfo.Data.MediaId != "" {
		response.Result = true
		response.MediaId = mediaInfo.Data.MediaId
		redis.Client().Set(c, cacheKey, mediaInfo.Data.MediaId, 2*12*time.Hour)
	} else {
		response.Result = false
	}

	return response, nil
}

// getTmpImgFile 获取临时文件
func (s *DrainageService) getTmpImgFile(c ctx.Context, request CacheMediaRequest) (name string, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(request.ImgUrl)
	if err != nil {
		return "", err
	}

	urlBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10)

	fileContent := []byte(urlBody)
	FilePath := config.CfgData.UploadPath

	fileType := ""
	file := ""
	if "" == fileType {
		ext := path.Ext(request.ImgUrl)
		// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
		if u, err := url.Parse(request.ImgUrl); err == nil && u.Path != "" {
			ext = path.Ext(u.Path)
		}
		file = FilePath + fileName + ext
	} else {
		file = fmt.Sprintf("%s%s.%s", FilePath, fileName, common.GetFileType(fileContent))
	}

	err = ioutil.WriteFile(file, fileContent, os.ModePerm)
	size, err := FileSize(file)
	log.Info(c).Msgf("【===文件大小为%d,err:%s===】", size, err)
	if err != nil {
		return "", err
	}

	return file, nil
}

// removeTmpImgFile 删除临时文件
func (s *DrainageService) removeTmpImgFile(c ctx.Context, file string) {
	log.Info(c).Msgf("【===删除临时文件===】")
	os.Remove(file)
}

// TagList 标签列表
func (s *DrainageService) TagList(c ctx.Context, request *ContactChannelLabelsListRequest) ([]*ContactChannelLabels, error) {
	return s.repo.SelectLabelsByChannelID(c, request)
}

// ContactSelect 二维码详情
func (s *DrainageService) ContactSelect(c ctx.Context, request ContactDetailsRequest) (*Contacts, error) {
	return s.repo.ContactSelect(c, request)
}

// 累计添加人数修改
func (s *DrainageService) SumAddMemberNumber(c ctx.Context, contactsID int64) error {
	err := s.repo.SumAddMemberNumber(c, contactsID)
	return err
}

func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}

func (s *DrainageService) GetWorkExternalUserRelation(c ctx.Context, relation customer.ExternalUserRelation) (*customer.ExternalUserRelation, error) {
	return s.repo.GetWorkExternalUserRelation(c, relation)
}

type Channel struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Channel      string `json:"channel" db:"channel"`
	SourceType   int64  `json:"source_type" db:"source_type"`
	AdminID      int64  `json:"admin_id" db:"admin_id"`
	AdminName    string `json:"admin_name" db:"admin_name"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// ContactQrCode 联系我二维码表
type ContactQrCode struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Name         string `json:"name" db:"name"`
	ChannelID    int64  `json:"channel_id" db:"channel_id"`
	Channel      string `json:"channel" db:"channel"`
	ContactType  int64  `json:"contact_type" db:"contact_type"`
	SourceType   int64  `json:"source_type" db:"source_type"`
	IsVerify     int    `json:"is_verify" db:"is_verify"`
	OfficialURL  string `json:"official_url" db:"official_url"`
	TransformURL string `json:"transform_url" db:"transform_url"`
	UcID         string `json:"uc_id" db:"uc_id"`
	ConfigID     string `json:"config_id" db:"config_id"`
	Object       string `json:"object" db:"object"`
	ObjectExtra  string `json:"object_extra" db:"object_extra"`
	ParamsType   string `json:"params_type" db:"params_type"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// UcMember uc用户表
type UcMember struct {
	WorkMemUid   string `json:"mem_uid" db:"mem_uid"`
	MemberId     int64  `json:"member_id" db:"member_id"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
}

// SynContact 同步联系我二维码
func (s *DrainageService) SynContact(c ctx.Context, enterpriseId int64) error {
	// 同步渠道数据
	_ = s.SynChannel(c, enterpriseId)

	// 同步联系我二维码
	_ = s.SynContactQrCode(c, enterpriseId)

	return nil
}

// SynContactQrCode 同步联系我二维码数据
func (s *DrainageService) SynContactQrCode(c ctx.Context, enterpriseId int64) error {
	var channelRequest = ContactQrCode{
		EnterpriseID: enterpriseId,
	}

	// 查询联系我列表
	ContactList, err := s.repo.GetSynContactList(c, channelRequest)
	if err != nil {
		log.Info(c).Msgf("=========同步联系我二维码业务--查询联系我列表===========：%+v", err)
	}
	for _, value := range *ContactList {
		var UcMember UcMember
		UcMember.EnterpriseId = enterpriseId
		if value.ParamsType == "1" {
			// 查询员工关联的企业微信UserId并返回平台需要的memberId
			UcMember.WorkMemUid = value.Object
		}
		if value.ParamsType == "2" {
			// 查询员工关联的企业微信UserId并返回平台需要的memberId
			UcMember.WorkMemUid = value.ObjectExtra
		}

		// 查询联系相对应的员工Id
		MemberData, err := s.repo.GetSynMemberId(c, UcMember)
		if err != nil {
			log.Info(c).Msgf("=========同步联系我二维码业务--查询联系我列表===========：%+v", err)
			continue
		}
		contacts := Contacts{
			ID:           value.ID,
			Name:         value.Name,
			ChannelId:    value.ChannelID,
			AutoFriend:   value.IsVerify,
			CreatedAt:    value.CreatedAt,
			UpdatedAt:    value.UpdatedAt,
			EnterpriseId: value.EnterpriseID,
			QrCodeUrl:    value.OfficialURL,
			MemberCount:  0,
			DeletedAt:    0,
			ConfigId:     value.ConfigID,
		}
		// 查询联系我二维码是否已同步
		ContactsInfo, _ := s.repo.GetSynContacts(c, contacts)
		if ContactsInfo != nil {
			continue
		}

		addErr := s.repo.ContactAdd(c, contacts)
		if addErr != nil {
			log.Info(c).Msgf("=========同步联系我二维码业务--添加联系我二维码数据===========：%+v", err)
		}

		for _, item := range *MemberData {
			// 添加联系人
			contactMember := ContactMembers{
				ID:           uid.ID(c),
				ContactId:    contacts.ID,
				EnterpriseId: contacts.EnterpriseId,
				MemberId:     item.MemberId,
				CreatedAt:    time.Now().Unix(),
				UpdatedAt:    time.Now().Unix(),
				DeletedAt:    0,
			}
			contactMembersAddError := s.repo.ContactMembersAdd(c, []*ContactMembers{&contactMember})
			if contactMembersAddError != nil {
				log.Info(c).Msgf("=========同步联系我二维码业务--联系我二维码添加员工===========：%+v", err)
			}
		}
	}

	return nil
}

// SynChannel 同步渠道信息
func (s *DrainageService) SynChannel(c ctx.Context, enterpriseId int64) error {
	var channelRequest = Channel{
		EnterpriseID: enterpriseId,
	}
	// 查询渠道数据表
	channelsList, err := s.repo.GetSynChannelsList(c, channelRequest)
	if err != nil {
		log.Info(c).Msgf("=========同步联系我二维码业务--同步渠道数据报错===========：%+v", err)
	}
	for _, value := range *channelsList {
		// 添加渠道数据
		ContactChannelsData := ContactChannels{
			ID:           value.ID,
			EnterpriseID: value.EnterpriseID,
			Name:         value.Channel,
			CreatedAt:    value.CreatedAt,
			UpdatedAt:    value.UpdatedAt,
		}
		// 判断渠道是否已存在
		channelsInfo, _ := s.repo.GetSynChannelsInfo(c, ContactChannelsData)
		if channelsInfo != nil {
			continue
		}

		_, channelsCreateErr := s.repo.ChannelsCreate(c, ContactChannelsData)
		if channelsCreateErr != nil {
			log.Info(c).Msgf("=========同步联系我二维码业务--同步渠道数据报错===========：%+v", channelsCreateErr)
		}
	}
	return nil
}

func (s *DrainageService) UpdateWelcomeMsg(ctx ctx.Context, req *Welcome) (int64, error) {
	return s.repo.UpdateWelcomeMsg(ctx, req)
}

// CountWelcomeByChannelID 根据渠道ID查询欢迎语数
func (s *DrainageService) CountWelcomeByChannelID(ctx ctx.Context, req *WelcomeListRequest) (int64, error) {
	return s.repo.CountWelcomeByChannelID(ctx, req)
}

// ExistRelationCustomer 留存好友数
// 从系统上线都筛选时间范围结束时间，查询外部联系人存在正常关系的的外部联系人个人
func (s *DrainageService) ExistRelationCustomer(c ctx.Context, req *StatisticsAllReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("exist")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.ExistRelationCount(c, req)
	//
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// ExistRelationCustomerByChannel 渠道留存好友数
// 从系统上线都筛选时间范围结束时间，查询外部联系人存在该渠道正常关系的的外部联系人个人
func (s *DrainageService) ExistRelationCustomerByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("exist")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.ExistRelationCountByChannel(c, req)

	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// IncreaseCustomerCount 新增好友总数
// 1. 在筛选时间范围内，该外部联系人要触发添加好友操作（不管在筛选时间内该外部联系人与该员工最后关系是否正常）
// 2. 该外部联系人在筛选时间之前与任何员工不存在正常的好友关系
func (s *DrainageService) IncreaseCustomerCount(c ctx.Context, req *StatisticsAllReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("increase")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.IncreaseCustomerCount(c, req)

	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// IncreaseCustomerCountByChannel 渠道新增好友总数
// 1. 在筛选时间范围内，该渠道外部联系人要触发添加好友操作（不管在筛选时间内该外部联系人与该员工最后关系是否正常）
// 2. 该外部联系人在筛选时间之前与该渠道任何员工不存在正常的好友关系
func (s *DrainageService) IncreaseCustomerCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("increase")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.IncreaseCustomerCountByChannel(c, req)

	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// LossCustomerCount 流失好友数量
// 1.在筛选时间范围内，该外部联系人需触发删除（任何一方删除都算）操作
// 2.从系统开始到筛选时间的结束日期,该外部联系人跟所有员工不是正常好友关系(不管员工删除客户还是客户删除员工)
func (s *DrainageService) LossCustomerCount(c ctx.Context, req *StatisticsAllReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("loss")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.LossCustomerCount(c, req)
	//
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// LossCustomerCountByChannel 渠道流失好友数量
// 1.在筛选时间范围内，该渠道外部联系人需触发删除（任何一方删除都算）操作
// 2.从系统开始到筛选时间的结束日期,该外部联系人跟渠道所有员工不是正常好友关系(不管员工删除客户还是客户删除员工)
func (s *DrainageService) LossCustomerCountByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("loss")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.LossCustomerCountByChannel(c, req)
	//
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// CustomerDelMember 好友删员工
// 1.在筛选时间范围内，该外部联系人需触发主动删员工操作。
// 2.在筛选时间范围内，删除该关系最后的状态为主动删员工状态。（当删除多个关系的时候，只要满足一条最后的状态为主动删员工状态即可）
func (s *DrainageService) CustomerDelMember(c ctx.Context, req *StatisticsAllReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("del_m")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.CustomerDelMemberCount(c, req)
	//
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// CustomerDelMemberByChannel 渠道好友删员工
// 1.在筛选时间范围内，该渠道外部联系人需触发主动删员工操作。
// 2.在筛选时间范围内，删除该渠道关系最后的状态为主动删员工状态。（当删除多个关系的时候，只要满足一条最后的状态为主动删员工状态即可）
func (s *DrainageService) CustomerDelMemberByChannel(c ctx.Context, req *StatisticsChannelReq) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("del_m")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.CustomerDelMemberCountByChannel(c, req)
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

// MemberDelCustomerOrChannel 员工删除客户总数
// 在筛选时间范围内，该外部联系人需触发被员工删除操作
func (s *DrainageService) MemberDelCustomerOrChannel(c ctx.Context, req *StatisticsChannelReq, corpId string) (int64, error) {
	// isToday := isToday(req.EndTime)
	// key := req.CacheKey("del_c")
	// if !isToday {
	//	data, ok := req.GetCache(c, key)
	//	if ok {
	//		return data, nil
	//	}
	// }

	result, err := s.repo.MemberDelCustomerCountOrChannel(c, req, corpId)
	// if err == nil && !isToday {
	//	req.SetCache(c, key, result)
	// }

	return result, err
}

func isToday(endTime int64) bool {
	a := time.Now().Format("2006-01-02")
	b := time.Unix(endTime, 0).Format("2006-01-02")
	return strings.EqualFold(a, b)
}
