package guide

import (
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/uid"
)

// Service  应用服务
type Service struct {
	repo         Repository
	Member       MemberInterface
	customerRepo customer.Repository
}

type MemberInterface interface {
	GetMemberId(c ctx.Context, enterpriseID int64, curUserID string) (int64, error)
}

// FindGuideByUserIDForCustomer 查询导购信息
func (s Service) FindGuideByUserIDForCustomer(ctx ctx.Context, enterpriseID int64, userID int64) (*customer.PortraitBindGuide, error) {
	rest, err := s.FindGuideRelation(ctx, enterpriseID, userID)
	if err != nil {
		return nil, err
	}
	return &customer.PortraitBindGuide{
		BindMemberID:   rest.MembersID,
		BindMemberName: rest.MemberUsername,
		UserID:         userID,
	}, nil
}

// FindGuideRelation 查询外部联系人关系
func (s Service) FindGuideRelation(cc ctx.Context, enterpriseID int64, userID int64) (*CustomerGuideRelation, error) {
	return s.repo.FindGuideRelationByUserID(cc, enterpriseID, userID)
}

// SyncUserBingGuide 同步导购数据
func (s *Service) SyncUserBingGuide(c ctx.Context, request *[]BasicUserBingGuideInfoEntity, event common.DdcEvent) (bool, error) {
	log.Debug(c).Msgf("请求同步导购数据参数，request:%+v,事件:%s", request, event)

	nowTime := time.Now().Unix()
	// 3.1 删除事件
	if event == common.DdcEventEnum.UserBindingGuideDelete {

		data := &BatchUserBingGuideInfoEntity{
			DeletedAt: nowTime,
		}

		for _, value := range *request {
			data.ThirdUserIds = append(data.ThirdUserIds, value.ThirdUserID)
			data.EnterpriseID = value.EnterpriseID
		}

		err := s.repo.DeleteBatchUserBingGuide(c, data)
		if err != nil {
			return true, err
		}
		return true, nil
	}

	// 3.2 新增事件
	if event == common.DdcEventEnum.UserBindingGuideAdd {
		for _, value := range *request {
			// 查询导购用户ID
			memberId, err := s.Member.GetMemberId(c, value.EnterpriseID, value.CurUserID)
			if err != nil {
				memberId = 0
			}
			value.MembersID = memberId
			if value.UnionID != "" {
				user, err := s.customerRepo.QueryEnterpriseUser(c, value.EnterpriseID, &customer.SourceUserBase{UnionID: value.UnionID})
				if err != nil {
					log.Err(c, err).Msgf("用户不存在: %d %s", value.EnterpriseID, value.UnionID)
				} else {
					value.UserID = user.ID
				}
			}

			value.ID = uid.ID(c)
			value.CreatedAt = nowTime
			value.UpdatedAt = nowTime
			_, err = s.repo.AddUserBingGuide(c, &value)
			if err != nil {
				return true, err
			}
		}
	}

	// 3.3 更新事件走
	if event == common.DdcEventEnum.UserBindingGuideUpdate {

		where := &BatchUserBingGuideInfoEntity{}

		for _, value := range *request {
			where.ThirdUserIds = append(where.ThirdUserIds, value.ThirdUserID)
			where.EnterpriseID = value.EnterpriseID
		}

		list, _ := s.repo.SelectUserBingGuide(c, where)

		var thirdUserIds []int64
		if list != nil {
			for _, val := range *list {
				thirdUserIds = append(thirdUserIds, val.ThirdUserID)
			}
		}

		for _, value := range *request {
			// 查询导购用户ID
			memberId, _ := s.Member.GetMemberId(c, value.EnterpriseID, value.CurUserID)
			value.MembersID = memberId
			if value.UnionID != "" {
				user, err := s.customerRepo.QueryEnterpriseUser(c, value.EnterpriseID, &customer.SourceUserBase{UnionID: value.UnionID})
				if err != nil {
					log.Err(c, err).Msgf("用户不存在: %d %s", value.EnterpriseID, value.UnionID)
				} else {
					value.UserID = user.ID
				}
			}

			value.UpdatedAt = nowTime
			if IsContain(thirdUserIds, value.ThirdUserID) {
				_, err := s.repo.UpdateUserBingGuide(c, &value)
				if err != nil {
					return true, err
				}
			} else {
				value.ID = uid.ID(c)
				value.CreatedAt = time.Now().Unix()
				_, err := s.repo.AddUserBingGuide(c, &value)
				if err != nil {
					return true, err
				}
			}
		}
	}

	return true, nil
}

func IsContain(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
