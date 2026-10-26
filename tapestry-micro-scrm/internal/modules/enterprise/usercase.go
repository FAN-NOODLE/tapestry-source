package enterprise

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// Service 应用服务
type Service struct {
	repo Repository
}

// Add 新增
func (s *Service) Add(ctx ctx.Context, req interface{}) int64 {
	id, err := s.repo.Add(ctx, req)
	if err != nil {
		return 0
	}
	return id
}

// Update 更新
func (s Service) Update(ctx ctx.Context, req interface{}) int64 {
	id, err := s.repo.Update(ctx, req)
	if err != nil {
		return 0
	}
	return id
}

// UpdateEnterprise 更新企业信息
func (s Service) UpdateEnterprise(ctx ctx.Context, req Corp) error {
	return s.repo.UpdateEnterprise(ctx, req)
}

// FindEnterpriseFromDb 从数据库里查询企业信息
func (s Service) FindEnterpriseFromDb(ctx ctx.Context, req Corp) (*Corp, error) {
	return s.repo.FindEnterpriseDetail(ctx, req)
}

// Detail 查询企业详情信息
func (s Service) Detail(ctx ctx.Context, c Corp) (*Corp, error) {
	redisKey := ""
	if c.EnterpriseID > 0 {
		redisKey = strconv.FormatInt(c.EnterpriseID, 10)
	}

	if c.CorpID != "" {
		redisKey = c.CorpID
	}
	if redisKey != "" {
		enterpriseStr := redis.Client().Get(ctx, redisKey).Val()
		if enterpriseStr != "" {
			enterprise := &Corp{}
			jErr := json.Unmarshal([]byte(enterpriseStr), enterprise)

			if jErr != nil {
				return nil, jErr
			}
			return enterprise, nil
		}
	}
	detail, err := s.repo.FindEnterpriseDetail(ctx, c)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if detail == nil {
		req := &boss.EnterpriseDetailRequest{
			EnterpriseID: c.EnterpriseID,
			Hash:         c.Hash,
			CorpID:       c.CorpID,
			Domain:       c.Domain,
		}
		enterpriseData, err := boss.EnterpriseDetail(ctx, req)

		if err != nil {
			return nil, err
		}

		c.CorpID = enterpriseData.CorpID
		c.Domain = enterpriseData.Domain
		c.Hash = enterpriseData.Hash
		c.EnterpriseID = enterpriseData.ID
		s.repo.Add(ctx, c)

		detail = &c
	}
	enterpriseByte, _ := json.Marshal(detail)
	// 入缓存
	redis.Client().Set(ctx, detail.CorpID, string(enterpriseByte), 2*time.Hour)
	redis.Client().Set(ctx, strconv.FormatInt(detail.EnterpriseID, 10), string(enterpriseByte), 2*time.Hour)
	return detail, nil
}

// GetEnterprises 返回企业信息列表
func (s Service) GetEnterprises(ctx ctx.Context, req CorpListRequest) (int64, []Corp, error) {
	var list []Corp
	return 0, list, nil
}

// GetDomain 通过企业id获取domain
func (s *Service) GetInfo(c ctx.Context, enterpriseID int64) (string, string, string, error) {
	var cond = Corp{
		EnterpriseID: enterpriseID,
	}
	corp, err := s.Detail(c, cond)
	if err != nil {
		return "", "", "", err
	}
	return corp.Domain, corp.CorpID, corp.Hash, nil
}

// GetCorpID 获取企业的 CORP ID
func (s *Service) GetCorpID(c ctx.Context, enterpriseID int64) (string, error) {
	_, corpID, _, err := s.GetInfo(c, enterpriseID)
	if err != nil {
		return "", err
	}

	return corpID, err
}

// EnterpriseJsTicket 获取企业 ticket
func (s *Service) EnterpriseJsTicket(c ctx.Context, enterpriseID int64) (string, error) {
	// 从缓存中取缓存的ticket
	key := fmt.Sprintf("ticket:enterprise_%s", strconv.FormatInt(enterpriseID, 10))
	jsTicket := redis.Client().Get(c, key).Val()
	if jsTicket != "" {
		return jsTicket, nil
	}

	// 从企业微信获取ticket
	domain, _, _, err := s.GetInfo(c, enterpriseID)
	if err != nil {
		return "", err
	}

	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	ticket, exp, err := workImpl.JsApiTicket(c)
	if err != nil {
		return "", err
	}

	if exp < 60 {
		return "", errors.New("[enterprise]enterprise ticket有效期小于60S")
	}

	ttl := time.Duration(exp-60) * time.Second
	redis.Client().Set(c, key, ticket, ttl)
	return ticket, nil
}

func (s *Service) AgentJsTicket(c ctx.Context, enterpriseID int64) (string, error) {
	// 从缓存中取缓存的ticket
	key := fmt.Sprintf("ticket:enterprise_agent_%s", strconv.FormatInt(enterpriseID, 10))
	jsTicket := redis.Client().Get(c, key).Val()
	if jsTicket != "" {
		return jsTicket, nil
	}

	// 从企业微信获取ticket
	domain, _, _, err := s.GetInfo(c, enterpriseID)
	if err != nil {
		return "", err
	}

	workImpl := new(work.QyWorkImpl).SetDomain(domain)
	ticket, exp, err := workImpl.JsAgentTicket(c)
	if err != nil {
		return "", err
	}

	if exp < 60 {
		return "", errors.New("[enterprise]agent ticket有效期小于60S")
	}

	ttl := time.Duration(exp-60) * time.Second
	redis.Client().Set(c, key, ticket, ttl)
	return ticket, nil
}

// ListEnterprise 企业列表
func (s *Service) ListEnterprise(cc ctx.Context, enterpriseIds []int64, page int, size int) (int64, *[]*Corp, error) {
	return s.repo.ListEnterprise(cc, enterpriseIds, page, size)
}

func (s *Service) GetOfficialInfo(c ctx.Context, enterpriseID int64) (appID string, secret string, err error) {
	officialInfo, err := s.OfficialInfo(c, enterpriseID)
	if err != nil {
		return
	}
	appID = officialInfo.AppID
	secret = officialInfo.Secret
	return
}

// OfficialInfo 获取公众号信息
func (s *Service) OfficialInfo(ctx ctx.Context, enterpriseID int64) (*OfficialInfo, error) {
	return s.repo.OfficialInfo(ctx, enterpriseID)
}

// CreateOfficialInfo 创建公众号信息
func (s Service) CreateOfficialInfo(ctx ctx.Context, data *EpConfig) error {
	return s.repo.CreateOfficialInfo(ctx, data)
}

// OfficialUser 查询公众号用户
func (s *Service) OfficialUser(ctx ctx.Context, req *OfficialUserReq) (*OfficialUser, error) {
	return s.repo.OfficialUser(ctx, req)
}

// CreateOfficialUser 创建公众号用户
func (s *Service) CreateOfficialUser(ctx ctx.Context, user *OfficialUser) error {
	return s.repo.CreateOfficialUser(ctx, user)
}

// GetEnterpriseID 获取企业ID
func (s Service) GetEnterpriseID(c ctx.Context, corpID string) (int64, error) {
	corp, err := s.Detail(c, Corp{
		CorpID: corpID,
	})

	if err != nil {
		return 0, err
	}

	if corp == nil {
		return 0, errors.New("未查询到企业信息")
	}

	return corp.EnterpriseID, nil
}

// AddWxAppConfig 添加公众号信息
func (s *Service) AddWxAppConfig(c ctx.Context, entity *WxAppConfigEntity) error {
	entity.ID = uid.ID(c)
	entity.CreatedAt = time.Now().Unix()
	entity.UpdatedAt = entity.CreatedAt
	_, err := s.repo.AddWxAppConfig(c, entity)
	if err != nil {
		return err
	}
	return nil
}
