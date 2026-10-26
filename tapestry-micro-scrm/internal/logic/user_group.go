package logic

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/logic/entity"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"sync"
)

type userGroupLogic struct {
	Name           string
	enterpriseID   int64
	externalUserID int64
	relationID     int64
	ctx            ctx.Context
}

func NewUserGroupLogic(ctx ctx.Context, enterpriseID int64) *userGroupLogic {
	return &userGroupLogic{
		Name:         "客户分组逻辑",
		enterpriseID: enterpriseID,
		ctx:          ctx,
	}
}

func (u *userGroupLogic) SetExternalUserID(externalUserID int64) *userGroupLogic {
	u.externalUserID = externalUserID
	return u
}

func (u *userGroupLogic) SetRelationID(relationID int64) *userGroupLogic {
	u.relationID = relationID

	// 根据关系ID，获取外部联系人ID
	u.externalUserID, _ = orm.Service.WorkExternalUser.GetExternalUserIDByRelationID(u.enterpriseID, u.relationID)

	return u
}

// syncEnterpriseUserGroups 同步用户的客户分组到enterprise_user.groups
func (u userGroupLogic) SyncEnterpriseUserGroups() error {
	userCustomerGroup, err := u.GetUserCustomerGroup()
	if err != nil {
		return err
	}

	if len(userCustomerGroup) == 0 {
		return nil
	}

	customerGroupIDList := []int64{}
	for _, customerGroup := range userCustomerGroup {
		customerGroupIDList = append(customerGroupIDList, customerGroup.ID)
	}

	err = orm.Service.WorkExternalUser.ReplaceExternalUserGroups(u.enterpriseID, u.externalUserID, customerGroupIDList)
	if err != nil {
		log.Info(u.ctx).Msgf("%s::同步客户分组失败,err: %s", u.Name, err)
		return err
	}

	return nil
}

// GetUserCustomerGroup 获取用户的客户分组
func (u userGroupLogic) GetUserCustomerGroup() (userCustomerGroup []entity.CustomerGroup, err error) {

	if 0 == u.externalUserID && 0 == u.relationID {
		err = errors.New(fmt.Sprintf("%s::无效用户ID"))
		return
	}

	// 获取用户的标签集合
	tagIDList, tiErr := orm.Service.Tag.GetTagIDListByExternalUserID(u.enterpriseID, u.externalUserID)
	if tiErr != nil {
		err = tiErr
		return
	}

	userCustomerGroup, err = u.GetUserCustomerGroupByTagID(tagIDList)
	return
}

// GetUserCustomerGroupByWxTagID 根据外部联系人的企业微信标签，获取用户的客户分组
func (u userGroupLogic) GetUserCustomerGroupByWxTagID(wwTagIDList []string) (userCustomerGroup []entity.CustomerGroup, err error) {
	// 根据企微标签ID，获取系统标签ID集合
	tagIDList := orm.Service.Tag.GetTagIDListByQYID(u.enterpriseID, wwTagIDList)

	userCustomerGroup, err = u.GetUserCustomerGroupByTagID(tagIDList)
	return
}

// GetUserCustomerGroupByTagID 根据用户的标签集合，获取用户的客户分组
func (u userGroupLogic) GetUserCustomerGroupByTagID(tagIDList []int64) (userCustomerGroup []entity.CustomerGroup, err error) {
	// 获取外部联系人，可能存在的用户分组ID
	customerGroupIDList, cgiErr := orm.Service.UserGroup.GetUserMaybeCustomerGroup(u.enterpriseID, tagIDList)
	if cgiErr != nil {
		err = cgiErr
		log.Info(u.ctx).Msgf("%s::获取客户分组ID集合失败,err: %s", u.Name, cgiErr)
		return
	}

	if len(customerGroupIDList) == 0 {
		return
	}

	mappingTagID := func() map[int64]bool {
		mapping := map[int64]bool{}
		for _, id := range tagIDList {
			mapping[id] = true
		}
		return mapping
	}()
	var waitGroup sync.WaitGroup
	for _, customerGroupID := range customerGroupIDList {
		waitGroup.Add(1)
		go func(userCustomerGroup *[]entity.CustomerGroup, customerGroupID int64) {
			defer waitGroup.Done()

			// 获取客户分组详情
			customerGroupDetail, cgdErr := u.GetCustomerGroupDetail(customerGroupID)
			if cgdErr != nil {
				log.Info(u.ctx).Msgf("%s::获取客户分组详情失败,err: %s", u.Name, cgdErr)
				return
			}

			// 检测该外部联系人是否满足该客户分组的条件
			// 过滤
			for _, tagGroup := range customerGroupDetail.FilterTagGroup {
				for _, tag := range tagGroup.TagList {
					if _, ok := mappingTagID[tag.ID]; ok {
						return
					}
				}
			}

			// 检测包含; 不同组之间的标签是且的关系
			isInclude := true
			for _, tagGroup := range customerGroupDetail.IncludeTagGroup {
				if !isInclude {
					isInclude = false
					break
				}

				for _, tag := range tagGroup.TagList {
					if _, ok := mappingTagID[tag.ID]; ok {
						isInclude = true
						break
					}

					isInclude = false
				}
			}

			if isInclude {
				*userCustomerGroup = append(*userCustomerGroup, entity.CustomerGroup{
					ID:   customerGroupDetail.ID,
					Name: customerGroupDetail.Name,
					Type: customerGroupDetail.Type,
				})
			}
		}(&userCustomerGroup, customerGroupID)
	}
	waitGroup.Wait()

	return
}

// GetCustomerGroupDetail 根据客户分组ID，获取客户分组的详情
func (u userGroupLogic) GetCustomerGroupDetail(customerGroupID int64) (data entity.CustomerGroupDetail, err error) {
	// 从数据库中获取客户分组的信息集合
	var result []micro_scrm.UserGroupDetailModel
	result, err = orm.Service.UserGroup.GetCustomerGroupDetail(u.enterpriseID, customerGroupID)
	if err != nil {
		return
	}

	// 整理客户分组信息为 entity.CustomerGroupDetail
	mapping := map[int64]bool{}
	isCustomerGroupFirst := true // 客户分组信息第一次赋值
	for _, customerGroupDetail := range result {
		// 客户分组信息
		if isCustomerGroupFirst {
			data = entity.CustomerGroupDetail{
				EnterpriseID:    u.enterpriseID,
				ID:              customerGroupID,
				Name:            customerGroupDetail.UserGroup.Name,
				Type:            customerGroupDetail.UserGroup.Type,
				IncludeTagGroup: []entity.TagGroupDetail{},
				FilterTagGroup:  []entity.TagGroupDetail{},
			}

			isCustomerGroupFirst = false
		}

		// 包含
		if micro_scrm.UGD_IncludeType_Include == customerGroupDetail.IncludeType {
			if _, ok := mapping[customerGroupDetail.TagGroupID]; !ok {
				// 添加新的标签分组
				u.addTagGroupToCustomerGroup(customerGroupDetail, &data)
			}

			u.addTagToTagGroup(customerGroupDetail, &data)
		}

		// 过滤
		if micro_scrm.UGD_IncludeType_Filter == customerGroupDetail.IncludeType {
			if _, ok := mapping[customerGroupDetail.TagGroupID]; !ok {
				// 添加新的标签分组
				u.addTagGroupToCustomerGroup(customerGroupDetail, &data)
			}

			u.addTagToTagGroup(customerGroupDetail, &data)
		}
	}

	return
}

// addTagGroupToCustomerGroup 给对应的客户分组，添加标签组，该方法，专属于userGroupLogic.GetCustomerGroupDetail
func (u userGroupLogic) addTagGroupToCustomerGroup(customerGroupDetail micro_scrm.UserGroupDetailModel, data *entity.CustomerGroupDetail) {
	// 添加新的标签分组
	data.IncludeTagGroup = append(data.IncludeTagGroup, entity.TagGroupDetail{
		EnterpriseID: u.enterpriseID,
		ID:           customerGroupDetail.TagGroupID,
		Name:         customerGroupDetail.TagGroup.GroupName,
		WwID:         customerGroupDetail.TagGroup.GroupID,
		Order:        customerGroupDetail.TagGroup.Order,
		TagList:      []entity.TagDetail{},
	})
}

// addTagToTagGroup 给对应的标签分组添加标签。该方法，专属于 userGroupLogic.GetCustomerGroupDetail
func (u userGroupLogic) addTagToTagGroup(customerGroupDetail micro_scrm.UserGroupDetailModel, data *entity.CustomerGroupDetail) {
	for index, tagGroup := range data.IncludeTagGroup {
		if tagGroup.ID == customerGroupDetail.TagGroupID {
			data.IncludeTagGroup[index].TagList = append(data.IncludeTagGroup[index].TagList, entity.TagDetail{
				EnterpriseID: u.enterpriseID,
				ID:           customerGroupDetail.TagID,
				Name:         customerGroupDetail.Tag.Name,
				WwID:         customerGroupDetail.Tag.TagID,
				Type:         customerGroupDetail.Tag.Type,
				Order:        customerGroupDetail.Tag.Order,
				TagGroupID:   customerGroupDetail.TagGroupID,
			})
		}
	}
}
