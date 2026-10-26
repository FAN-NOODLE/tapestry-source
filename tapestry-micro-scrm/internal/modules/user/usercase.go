package user

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// Service 用户服务
type Service struct {
	repo           Repository
	enterpriseRepo enterprise.Repository
	esRepo         customer.EsRepository
	eventRepo      EventRepository
	workRepo       customer.WorkUserRepository
}

// MergeUser 合并用户
func (s *Service) MergeUser(c ctx.Context, source *SourceUser) (bool, error) {
	return s.repo.MergeUser(c, source)
}

// 同步公众号用户
func (s *Service) SynWxOpenUser(c ctx.Context, requestUserInfo *OfficialUserEntity) (bool, error) {
	//调用合并接口，处理合并用户逻辑
	source := &SourceUser{
		EnterpriseID: requestUserInfo.EnterpriseID,
		SourceUserBase: customer.SourceUserBase{
			AppID:        requestUserInfo.AuthorizerAppid,
			OpenID:       requestUserInfo.OpenId,
			UnionID:      requestUserInfo.UnionID,
			CustomerType: customer.CustomerTypeOpen,
		},
		EnterpriseUserBase: customer.EnterpriseUserBase{
			Name:     requestUserInfo.Name,
			Nick:     requestUserInfo.Nickname,
			Avatar:   requestUserInfo.Avatar,
			Birthday: requestUserInfo.Birthday,
			Sex:      int8(requestUserInfo.Sex),
			Mobile:   requestUserInfo.Mobile,
		},
		EnterpriseUserLocation: customer.EnterpriseUserLocation{
			Country:  requestUserInfo.Country,
			City:     requestUserInfo.City,
			Province: requestUserInfo.Province,
		},
		EnterpriseUserSubscribe: customer.EnterpriseUserSubscribe{
			SubscribeTime:   requestUserInfo.SubscribeTime,
			SubscribeStatus: requestUserInfo.SubscribeStatus,
			SubscribeScene:  requestUserInfo.SubscribeScene,
		},
		EnterpriseUserOther: customer.EnterpriseUserOther{
			QRScene:    requestUserInfo.QRScene,
			QRSceneStr: requestUserInfo.QRSceneStr,
			Remark:     requestUserInfo.Remark,
		},
		CreateAt: int64(requestUserInfo.SubscribeTime),
	}
	if source.AppID == "" {
		source.AppID, _ = s.enterpriseRepo.QueryOneApp(c, source.EnterpriseID, customer.CustomerTypeOpen)
	}
	_, err := s.repo.MergeUser(c, source)
	if err != nil {
		log.Info(c).Msgf("SynWxOpenUser合并用户数据MergeUser失败,err:%s,data:%+v", err.Error(), requestUserInfo)
		return false, err
	}
	return true, nil
}

// 同步商城用户
func (s *Service) SynWxShopUser(c ctx.Context, request *ShopUserEntity) (bool, error) {
	//调用合并接口，处理合并用户逻辑
	source := &SourceUser{
		EnterpriseID: request.EnterpriseID,
		SourceUserBase: customer.SourceUserBase{
			AppID:        request.AuthorizerAppid,
			OpenID:       request.OpenId,
			UnionID:      request.UnionID,
			CustomerType: customer.CustomerTypeShop,
		},
		ThirdUserId: request.ThirdUserID,
		EnterpriseUserBase: customer.EnterpriseUserBase{
			Name:     request.Name,
			Nick:     request.Nickname,
			Avatar:   request.Avatar,
			Birthday: request.Birthday,
			Sex:      int8(request.Sex),
			Mobile:   request.Mobile,
		},
		EnterpriseUserLocation: customer.EnterpriseUserLocation{
			Country:  request.Country,
			City:     request.City,
			Province: request.Province,
		},
		EnterpriseUserOther: customer.EnterpriseUserOther{
			QRScene:    request.QRScene,
			QRSceneStr: request.QRSceneStr,
			Remark:     request.Remark,
		},
		CreateAt: int64(request.CreatedAt),
	}

	//对不传app_id的地方做兼容
	if source.AppID == "" {
		source.AppID, _ = s.enterpriseRepo.QueryOneApp(c, source.EnterpriseID, customer.CustomerTypeShop)
	}

	bAddUser, err := s.repo.MergeUser(c, source)
	if err != nil {
		log.Info(c).Msgf("SynWxShopUser合并用户数据MergeUser失败,err:%s,data:%+v", err.Error(), request)
		return false, err
	}

	if bAddUser {
		shopName, _ := s.getShopName(c, request.EnterpriseID)
		if shopName != "" {
			_ = s.AddUserEvent(c, request.EnterpriseID, source.EuID, common.EventEnum.FirstVisitShop.String(), request.ShopName, strconv.FormatInt(request.EnterpriseID, 10), "", int64(request.CreatedAt))
		}
	}

	return true, nil
}

// AddUserEvent 添加用户事件
func (s *Service) AddUserEvent(c ctx.Context, enterpriseID int64, userID int64, event string, eventDesc string, object string, extra interface{}, eventTime int64) error {

	var extraStr string
	switch extraV := extra.(type) {
	case string:
		extraStr = extraV
	default:
		bytes, _ := json.Marshal(extra)
		extraStr = string(bytes)
	}

	var eventEntity = EventEntity{
		ID:           uid.ID(c),
		EnterpriseID: enterpriseID,
		UserID:       userID,
		Event:        event,
		Object:       object,
		EventDesc:    eventDesc,
		Extra:        extraStr,
		CreatedAt:    int(eventTime),
		UpdatedAt:    int(time.Now().Unix()),
	}
	_, err := s.eventRepo.Add(c, eventEntity)
	return err
}

// QueryUserCount 计算标签用户数
func (s *Service) QueryUserCount(c ctx.Context, cond *customer.ListCond) (*customer.TagUserCount, error) {
	return s.esRepo.QueryUserCount(c, cond)
}

func (s *Service) QueryUserCountSql(c ctx.Context, cond *customer.ListCond) (int, error) {
	return s.repo.QueryUserCount(c, cond)
}

func (s *Service) EditRemarks(c ctx.Context, request EditRemarksRequest, enterpriseService *enterprise.Service) (*EditRemarksResponse, error) {
	switch request.EditType {
	case EditTypeOne:
		return s.EditRemarksOne(c, request, enterpriseService)
	case EditTypeMultiple:
		if len(request.UserIDs) > 0 {
			userList, err := s.repo.EnterpriseUsers(c, request.EnterpriseID, request.UserIDs)
			if err != nil {
				return nil, err
			}
			return s.EditRemarkMultiple(c, request.EnterpriseID, userList, request.Remarks, enterpriseService)
		}
	case EditTypeAll:
		return s.EditRemarksAll(c, request, enterpriseService)
	default:
		return nil, errs.ParamsErr
	}
	return nil, errs.ParamsErr
}

func (s *Service) EditRemarksAll(c ctx.Context, request EditRemarksRequest, enterpriseService *enterprise.Service) (*EditRemarksResponse, error) {
	cond := customer.ListRequest{
		ListCond: request.ListCond,
		Page:     1,
		Size:     999,
	}
	cond.EnterpriseID = common.EnterpriseID(c)
	_, userList, err := s.esRepo.QueryUserList(c, &cond)

	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	return s.EditRemarkMultiple(c, request.EnterpriseID, userList, request.Remarks, enterpriseService)
}

func (s *Service) EditRemarkMultiple(c ctx.Context, enterpriseID int64, userList []*customer.EnterpriseUser, remark string, enterpriseService *enterprise.Service) (*EditRemarksResponse, error) {
	wg := &sync.WaitGroup{}
	for _, value := range userList {
		go func(user *customer.EnterpriseUser) {
			wg.Add(1)
			defer wg.Done()
			var desc string = remark
			if user.Remark != "" {
				desc = remark + "," + user.Remark
			}

			if utf8.RuneCountInString(desc) > ExternalDescriptionLength {
				log.Info(c).Msgf("用户备注长度超过150，更新失败, desc:%s, value:%+v", desc, user)
				return
			}

			if err := s.repo.UpdateDescription(c, enterpriseID, user.ID, desc); err != nil {
				log.Info(c).Msgf("更新用户描述失败, err:%s, user:%v", err.Error(), user)
				return
			}

			var workExternalUserIDs []int64
			for _, u := range user.WxExternal {
				workExternalUserIDs = append(workExternalUserIDs, u.ID)
			}

			if len(workExternalUserIDs) > 0 {
				relations, err := s.workRepo.QueryExternalRelationRemarks(c, enterpriseID, workExternalUserIDs)
				if err != nil {
					log.Info(c).Msgf("获取外部联系人关系数据失败, %s, %v", err.Error(), user)
					return
				}
				for _, relation := range relations {
					var desc string = remark
					if relation.Remark != "" {
						desc = remark + "," + user.Remark
					}
					if err := s.workRepo.UpdateExternalRelationRemark(c, relation.RelationID, desc); err != nil {
						log.Info(c).Msgf("更新外部联系人备注失败, %s, %v", err.Error(), relation)
					}

					err = s.UpdateExternalRelationDescription(c, enterpriseID, &relation.WorkExternalRelationForEdit, enterpriseService)
					if err != nil {
						log.Info(c).Msgf("同步企业微信备注失败：err:%v", err)
					}
				}
			}
		}(value)
	}
	wg.Wait()
	return nil, nil
}

func (s *Service) EditRemarksOne(c ctx.Context, request EditRemarksRequest, enterpriseService *enterprise.Service) (*EditRemarksResponse, error) {
	if utf8.RuneCountInString(request.Remarks) > ExternalDescriptionLength {
		return nil, errs.ExternalRemarksLengthErr
	}

	if request.RelationID != 0 {
		relation, err := s.workRepo.QueryExternalRelationRemark(c, request.RelationID, request.EnterpriseID)
		if err != nil {
			return nil, err
		}
		if err := s.workRepo.UpdateExternalRelationRemark(c, request.RelationID, request.Remarks); err != nil {
			return nil, err
		}

		err = s.UpdateExternalRelationDescription(c, request.EnterpriseID, &relation.WorkExternalRelationForEdit, enterpriseService)
		if err != nil {
			return nil, err
		}

		userExternal, err := s.repo.QueryByID(c, request.EnterpriseID, customer.CustomerTypeExternal, relation.WorkExternalUserID)
		if err != nil {
			return nil, err
		}
		if err := s.repo.UpdateDescription(c, request.EnterpriseID, userExternal.EuID, request.Remarks); err != nil {
			return nil, err
		}
	} else if request.UserID != 0 {
		if err := s.repo.UpdateDescription(c, request.EnterpriseID, request.UserID, request.Remarks); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *Service) UpdateExternalRelationDescription(c ctx.Context, enterpriseID int64, editInfo *customer.WorkExternalRelationForEdit, enterpriseService *enterprise.Service) error {
	corp := enterprise.Corp{
		EnterpriseID: enterpriseID,
	}

	enterpriseInfo, err := enterpriseService.Detail(c, corp)
	if err != nil {
		return err
	}

	qyWork := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	updateData := work.EditExternalDescriptionRequest{
		UserID:         editInfo.MemberUserID,
		ExternalUserid: editInfo.ExternalUserID,
		Description:    editInfo.Remark,
	}
	err = qyWork.UpdateExternalRelationDescription(c, updateData)
	return err
}

// ListConfig 列表配置
func (s *Service) ListConfig(c ctx.Context) (*ListConfigResponse, error) {
	enterpriseID := common.EnterpriseID(c)
	var response = ListConfigResponse{
		ShowOfficial: false,
		ShowShop:     false,
	}
	cntmap, err := s.enterpriseRepo.QueryAppTypeCount(c, enterpriseID)
	if err != nil {
		return nil, err
	}

	response.ShowOfficial = cntmap[customer.CustomerTypeOpen] > 0
	response.ShowShop = cntmap[customer.CustomerTypeShop] > 0
	return &response, nil
}

// GetWorkUserInfoByUserIDs fission模块，根据用户id查询外部联系人信息
func (s *Service) GetWorkUserInfoByUserIDs(c ctx.Context, enterpriseID int64, userIDs []int64) (*[]*fission.WorkUserBasicInfo, error) {
	userList, err := s.GetWorkUserByUserIDs(c, enterpriseID, userIDs)
	if err != nil {
		return nil, err
	}
	if userList == nil {
		return nil, nil
	}
	var resList []*fission.WorkUserBasicInfo
	for _, v := range *userList {
		resList = append(resList, &fission.WorkUserBasicInfo{
			Name:   v.Name,
			ImgURL: v.Avatar,
			UserID: v.ID,
		})
	}
	return &resList, nil
}

// GetWorkUserByUserIDs 根据work_external_user表id查询外部联系人数据
func (s *Service) GetWorkUserByUserIDs(c ctx.Context, enterpriseID int64, userIDs []int64) (*[]*customer.WorkExternalUser, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	list, err := s.workRepo.QueryWorkUserByUserIDs(c, enterpriseID, userIDs)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetRobotUserList 获取假用户列表
func (s *Service) getRobotList(c ctx.Context, num int) ([]*RobotUserEntity, error) {
	return s.repo.GetRobotUserList(c, num)
}

// GetRobotUserList 获取假用户列表
func (s *Service) GetRobotUserList(c ctx.Context, num int) ([]*fission.AssistantUserInfo, error) {
	robotList, err := s.getRobotList(c, num)
	if err != nil {
		return nil, err
	}
	var userList []*fission.AssistantUserInfo
	for _, v := range robotList {
		userList = append(userList, &fission.AssistantUserInfo{
			Name:   v.NickName,
			ImgUrl: v.ImgUrl,
		})
	}
	return userList, nil
}

// getShopName 获取商城小程序名称
func (s *Service) getShopName(c ctx.Context, enterpriseID int64) (string, error) {
	appInfo, err := s.getAppInfo(c, enterpriseID)
	if err != nil {
		log.Info(c).Msgf("获取商城小程序名称失败, err:%s, enterpriseID:%d", err.Error(), enterpriseID)
		return "", err
	}
	return appInfo.Data.Name, nil
}

// 获取小程序授权信息
func (s *Service) getAppInfo(c ctx.Context, enterpriseID int64) (*boss.AuthDetailResponse, error) {
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
		redisErr := redis.Client().Set(c, appInfoKey, appInfoByte, 30*24*time.Hour).Err()
		if redisErr != nil {
			return nil, redisErr
		}
		return appInfo, nil
	}
}

func generateKey(redisReq interface{}) string {
	keyData, _ := json.Marshal(redisReq)
	m := md5.New()
	m.Write([]byte(keyData))
	redisKey := fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
	return redisKey
}
