package job

import (
	"database/sql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"math"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/errorIgnore"
	"strconv"
	"strings"
	"sync"
	"time"

	"micro-scrm/internal/modules/scripts"
	scriptsMysql "micro-scrm/internal/modules/scripts/respository/mysql"

	"micro-scrm/internal/orm/model/micro_scrm"

	"micro-scrm/internal/modules/tag/entity"

	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/respository/mysql"

	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"

	"micro-scrm/internal/modules/customer"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/xjob"
)

// CalculateUserGroup 用于同步enterprise_user表,计算用户分组
func CalculateUserGroup(c ctx.Context, jobHandler *xjob.JobHandler) error {
	dbClient = db.DB()

	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	if len(specsArgs) <= 0 {
		// 获取所有企业
		enterpriseService := enterprise.New(enterpriseMysql.New())
		total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
		if err != nil || total < 0 {
			return err
		}
		for _, e := range *enterprises {
			specsArgs = append(specsArgs, strconv.FormatInt(e.EnterpriseID, 10))
		}
	}
	for _, spec := range specsArgs {
		enterpriseID, _ := strconv.ParseInt(spec, 10, 64)
		if enterpriseID == int64(3353921300665408) {
			continue
		}
		if err := CalculateGroup(c, enterpriseID); err != nil {
			continue
		}
	}
	return nil
}

type UserTagAndGroup struct {
	ID     int64
	Tags   []int64
	Groups []int64
}
type TmpCustomerGroupDetail struct {
	IncludeTagGroupKey map[string]bool `json:"include_tag_group" desc:"包含的客户标签和标签组的联合key"`
	FilterTagGroupKey  map[string]bool `json:"include_tag_group" desc:"过滤的客户标签和标签组的联合key"`
	FilterTagGroup     map[int64]int64 `json:"filter_tag_group" desc:"过滤的客户标签组"`
	IncludeTagGroup    map[int64]int64 `json:"include_tag_group" desc:"包含的客户标签组"`
}
type UserGroupInfo struct {
	EUid    int64
	GroupID int64
}

type UserGroupsMap struct {
	UserGroupList []UserGroupInfo
}

func CalculateGroup(c ctx.Context, enterpriseID int64) error {
	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	userGroupService := tag.NewUserGroup(mysql.NewUserGroupRepo())
	userService := user.NewUserService(userMysql.NewUser(), userMysql.NewEvent(), userEs.NewEsRepo(), enterpriseMysql.New())
	log.Info(c).Msgf("当前正在执行企业: %v，开始分组", enterpriseID)
	// 查询企业的所有分组
	groups, err := userGroupService.Repo.GetUserGroupList(c, enterpriseID)
	if err != nil {
		return err
	}
	// 对分组进行格式化
	groupList := make(map[int64]TmpCustomerGroupDetail)
	groupTags := make(map[int64][]int64)
	userGroups := make(map[int64][]int64)
	// 用户分组集合检查，避免重复
	checkUserGroup := make(map[string]int64)
	groupDepartments := make(map[int64][]int64)
	needAddGroupMap := make(map[int64]struct{}, 0)
	for _, group := range groups {

		checkTag := make(map[int64]bool)
		checkDempartment := make(map[int64]bool)
		// 根据分组id获取分组详情
		tagInfoList, listErr := userGroupService.Repo.GetUserGroupDetail(c, entity.QueryUserGroupParam{
			EnterpriseID: enterpriseID,
			UserGroupID:  group.ID,
		})
		if listErr != nil {
			continue
		}
		for _, v := range tagInfoList {
			if v.UserGroupType == 2 { // 客户名单导入的数据同步
				nowTime := time.Now().Unix()
				count := 0
				cond := &customer.ListCond{EnterpriseID: enterpriseID}
				cond.CustomerGroupCond.Include = true
				cond.CustomerGroupCond.GroupIDList = []int64{group.ID}
				totalRes, _ := userService.QueryUserCount(c, cond)
				if totalRes != nil {
					count = int(totalRes.MemberCount)
				}
				if count > 0 { // 判断这个分组的数据已经统计过了，就不统计了
					needAddGroupMap[group.ID] = struct{}{}
					continue
				}
				total, countErr := userGroupService.GetImportDataCount(enterpriseID, v.FileID)
				if countErr != nil {
					continue
				}
				pageSize := 100
				maxPage := int(math.Ceil(float64(total / int64(pageSize))))
				if maxPage == 0 {
					maxPage += 1
				}
				for i := 1; i <= maxPage; i++ {
					offset := (i - 1) * pageSize
					datalist, dataErr := userGroupService.GetImportData(enterpriseID, v.FileID, offset, pageSize)
					if dataErr != nil {
						continue
					}
					var unionID []string
					for _, data := range datalist {
						unionID = append(unionID, data.UnionID)
					}
					var unionIdMap = make(map[string]struct{}, 0)
					var newDataLogArr []micro_scrm_enterprise.UserGroupImportDataLogModel
					for _, value := range unionID {
						unionIdMap[value] = struct{}{}
						var newDataLog micro_scrm_enterprise.UserGroupImportDataLogModel
						newDataLog.CreatedAt = nowTime
						newDataLog.UpdatedAt = nowTime
						newDataLog.UnionID = value
						newDataLog.UserGroupID = group.ID
						newDataLog.ImportRes = "匹配失败"
						newDataLogArr = append(newDataLogArr, newDataLog)
					}
					_ = orm.Service.UserGroup.CreateExportLog(enterpriseID, newDataLogArr)
					condUnionId := &customer.CondUnionID{EnterpriseID: enterpriseID}
					condUnionId.UnionId = unionID
					_, users, userErr := customerService.QueryByUnionID(c, condUnionId)
					if userErr != nil {
						continue
					}
					for _, esUser := range users {
						if _, ok := unionIdMap[esUser.UnionId]; !ok {
							continue
						}
						var groupsMap = make(map[int64]struct{}, 0)
						for _, item := range esUser.Groups {
							groupsMap[item] = struct{}{}
						}
						if _, ok := groupsMap[group.ID]; !ok {
							esUser.Groups = append(esUser.Groups, group.ID)
							esErr := customerService.UpdateEnterpriseUserAttr(ctx.Wrap(c), enterpriseID, esUser.ID, &customer.EnterpriseUserAttr{Groups: esUser.Groups})
							if esErr != nil {
								log.Info(c).Msgf("更新用户的分组:%+v, 错误：%v", esUser.ID, err)
							} else {
								_ = orm.Service.UserGroup.UpdateExportLogStatus(enterpriseID, group.ID, esUser.UnionId)
								userGroups[esUser.ID] = append(userGroups[esUser.ID], group.ID)
							}
						}
					}
				}
				_ = orm.Service.UserGroup.UpdateDetailISExtract(enterpriseID, group.ID)
				continue
			}
			if v.UserGroupType == 1 { // 标签分组的同步
				if _, ok := checkTag[v.TagID]; !ok {
					checkTag[v.TagID] = true
					groupTags[group.ID] = append(groupTags[group.ID], v.TagID)
				}

				key := strconv.FormatInt(v.TagGroupID, 10) + "-" + strconv.FormatInt(v.TagID, 10)
				if _, ok := groupList[group.ID]; ok {
					if micro_scrm.UGD_IncludeType_Include == v.IncludeType {
						groupList[group.ID].IncludeTagGroup[v.TagGroupID] = v.TagGroupID
						groupList[group.ID].IncludeTagGroupKey[key] = true
					} else {
						groupList[group.ID].FilterTagGroup[v.TagGroupID] = v.TagGroupID
						groupList[group.ID].FilterTagGroupKey[key] = true
					}
				} else {
					tmp := TmpCustomerGroupDetail{
						IncludeTagGroupKey: make(map[string]bool),
						FilterTagGroupKey:  make(map[string]bool),
						FilterTagGroup:     make(map[int64]int64),
						IncludeTagGroup:    make(map[int64]int64),
					}
					if micro_scrm.UGD_IncludeType_Include == v.IncludeType {
						tmp.IncludeTagGroup[v.TagGroupID] = v.TagGroupID
						tmp.IncludeTagGroupKey[key] = true
					} else {
						tmp.FilterTagGroup[v.TagGroupID] = v.TagGroupID
						tmp.FilterTagGroupKey[key] = true
					}
					groupList[group.ID] = tmp
				}
			}

			if v.UserGroupType == 3 { // 门店维度的同步
				if _, ok := checkDempartment[v.DepartmentID]; !ok {
					checkDempartment[v.DepartmentID] = true
					groupDepartments[group.ID] = append(groupDepartments[group.ID], v.DepartmentID)
				}
			}

		}
	}
	if len(groupDepartments) > 0 {
		for gid, value := range groupDepartments {
			log.Info(c).Msgf("当前Department_id: %v，Department_id", value)
			pageSize := 100
			offset := 0
			pageVO := repos.PageVO{
				Offset: offset, Limit: pageSize,
			}
			total, datalist, dataErr := customerService.SelectEUidBYDepartment(ctx.Wrap(c), enterpriseID, value, pageVO)
			if errorIgnore.Ignore(dataErr) != nil || total == 0 {
				log.Info(c).Msgf("当前Department_iderr: %v，Department_iderr", dataErr)
				continue
			}

			if total > int64(pageSize) {
				maxPage := int(math.Ceil(float64(total / int64(pageSize))))
				if maxPage == 0 {
					maxPage += 1
				}
				for i := 1; i <= maxPage; i++ {
					offset = i * pageSize
					pageVO = repos.PageVO{
						Offset: offset, Limit: pageSize,
					}
					count, list, errs := customerService.SelectEUidBYDepartment(ctx.Wrap(c), enterpriseID, value, pageVO)
					if errorIgnore.Ignore(errs) != nil || count == 0 {
						break
					}
					userGroupsInfo, errs := UpdateGroupsBYEuid(c, list, enterpriseID, gid, customerService)
					if errs != nil {
						continue
					}
					for _, v := range userGroupsInfo.UserGroupList {
						userGroups[v.EUid] = append(userGroups[v.EUid], v.GroupID)
					}
				}
			} else {
				userGroupsInfo, upErrs := UpdateGroupsBYEuid(c, datalist, enterpriseID, gid, customerService)
				if upErrs != nil {
					continue
				}
				for _, v := range userGroupsInfo.UserGroupList {
					userGroups[v.EUid] = append(userGroups[v.EUid], v.GroupID)
				}
			}

		}
	}
	if len(groupList) > 0 {
		for gid, groupArr := range groupList {
			if groupTags[gid] == nil {
				continue
			}
			// 根据标签查找所有有该标签的人
			request := customer.ListRequest{
				ListCond: customer.ListCond{
					EnterpriseID: enterpriseID,
					TagCond: struct {
						Include   bool    `json:"include"`
						TagIDList []int64 `json:"tag_id_list"`
					}{
						true, groupTags[gid],
					},
				},
				Page: 1,
				Size: 1,
			}
			// 先查询出总数
			total, _, err := customerService.SimpleQueryList(c, &request)
			if err != nil {
				log.Error(c).Msgf("根据标签id集合，查找enterprise_user, 错误：%+v", err)
				continue
			}
			// 因es查询每页最大10000，所以需要分页
			maxPage := int(math.Ceil(float64(total / 10000)))
			if maxPage == 0 {
				maxPage += 1
			}
			for i := 1; i <= maxPage; i++ {
				req := customer.ListRequest{
					ListCond: customer.ListCond{
						EnterpriseID: enterpriseID,
						TagCond: struct {
							Include   bool    `json:"include"`
							TagIDList []int64 `json:"tag_id_list"`
						}{
							true, groupTags[gid],
						},
					},
					Page: i,
					Size: 10000,
				}
				_, userList, err := customerService.SimpleQueryList(c, &req)
				if err != nil {
					log.Error(c).Msgf("根据标签id集合，查找enterprise_user, 错误：%+v", err)
					continue
				}

				for _, users := range userList {
					for _, v := range users.Groups {
						if _, ok := needAddGroupMap[v]; ok {
							userGroups[users.ID] = append(userGroups[users.ID], gid)
						}
					}
					includeTagGroup := make(map[int64]int64)
					FilterTagGroup := make(map[int64]int64)
					for _, t := range users.Tags {
						for y := range groupArr.FilterTagGroupKey {
							// key中是否包含tagid，如果包含,则退出循环，进行下一个
							isContain := strings.Contains(y, strconv.FormatInt(t, 10))
							if !isContain {
								continue
							}
							arr := strings.Split(y, "-")
							tagGroupID, _ := strconv.ParseInt(arr[0], 10, 64)
							FilterTagGroup[tagGroupID] = tagGroupID
						}
						if len(FilterTagGroup) > 0 && len(FilterTagGroup) == len(groupArr.FilterTagGroup) {
							break
						}

						for y := range groupArr.IncludeTagGroupKey {
							// key中是否包含tagid，如果包含,则退出循环，进行下一个
							isContain := strings.Contains(y, strconv.FormatInt(t, 10))
							if !isContain {
								continue
							}
							arr := strings.Split(y, "-")
							tagGroupID, _ := strconv.ParseInt(arr[0], 10, 64)
							includeTagGroup[tagGroupID] = tagGroupID
						}
					}

					if len(includeTagGroup) > 0 && len(includeTagGroup) == len(groupArr.IncludeTagGroup) {
						key := strconv.FormatInt(users.ID, 10) + "-" + strconv.FormatInt(gid, 10)
						if _, ok := checkUserGroup[key]; !ok {
							userGroups[users.ID] = append(userGroups[users.ID], gid)
							checkUserGroup[key] = gid
						}
					}
				}

			}
		}

		// 循环用户的分组，更新用户分组
		wg := &sync.WaitGroup{}
		ch := make(chan bool, 3) // 控制并发的数量
		for euid, userGroup := range userGroups {
			wg.Add(1)
			ch <- true

			go func(euid int64, userGroup []int64, ch chan bool, wg *sync.WaitGroup) {
				uniqueUserGroup := make(map[int64]int64, 0)
				for _, v := range userGroup {
					if _, ok := uniqueUserGroup[v]; !ok {
						uniqueUserGroup[v] = v
					}
				}
				Group := make([]int64, 0)
				for _, v := range uniqueUserGroup {
					Group = append(Group, v)
				}
				defer EndGroup(wg, ch)
				if err := customerService.UpdateEnterpriseUserAttr(ctx.Wrap(c), enterpriseID, euid, &customer.EnterpriseUserAttr{Groups: Group}); err != nil {
					log.Info(c).Msgf("更新用户的分组:%+v, 错误：%v", euid, err)
				}
			}(euid, userGroup, ch, wg)
		}
		wg.Wait()
	}
	// 给用户打上分组
	log.Info(c).Msgf("执行企业: %+v，分组结束", enterpriseID)
	return nil
}

// UpdateGroupsBYEuid 更新用户的group通过euid
func UpdateGroupsBYEuid(c ctx.Context, datalist *[]*customer.EUidList, enterpriseID, gid int64, customerService *customer.Service) (list UserGroupsMap, err error) {
	var EUidArr []int64
	EUidMap := make(map[int64]struct{}, 0)
	for _, item := range *datalist {
		if item.EUid == 0 {
			continue
		}
		EUidArr = append(EUidArr, item.EUid)
		if _, ok := EUidMap[item.EUid]; !ok {
			EUidMap[item.EUid] = struct{}{}
		}
	}
	condUnionId := &customer.CondUnionID{EnterpriseID: enterpriseID}
	condUnionId.EUid = EUidArr
	_, users, userErr := customerService.QueryByUnionID(c, condUnionId)
	if userErr != nil {
		return UserGroupsMap{}, userErr
	}
	for _, esUser := range users {
		if _, ok := EUidMap[esUser.ID]; !ok {
			continue
		}
		var groupsMap = make(map[int64]struct{}, 0)
		for _, item := range esUser.Groups {
			groupsMap[item] = struct{}{}
		}
		if _, ok := groupsMap[gid]; !ok {
			esUser.Groups = append(esUser.Groups, gid)
			if err := customerService.UpdateEnterpriseUserAttr(ctx.Wrap(c), enterpriseID, esUser.ID, &customer.EnterpriseUserAttr{Groups: esUser.Groups}); err != nil {
				log.Info(c).Msgf("更新用户的分组:%+v, 错误：%v", esUser.ID, err)
				return UserGroupsMap{}, err
			}
			var userGroupInfo UserGroupInfo
			userGroupInfo.EUid = esUser.ID
			userGroupInfo.GroupID = gid
			list.UserGroupList = append(list.UserGroupList, userGroupInfo)

		}
	}
	return list, nil
}

func EndGroup(wg *sync.WaitGroup, ch chan bool) {
	<-ch
	wg.Done()
}

// FixHistoryUserTag  用于同步enterprise_user表,历史用户的标签
func FixHistoryUserTag(c ctx.Context, jobHandler *xjob.JobHandler) error {
	dbClient = db.DB()

	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	if len(specsArgs) <= 0 {
		// 获取所有企业
		enterpriseService := enterprise.New(enterpriseMysql.New())
		total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
		if err != nil || total < 0 {
			return err
		}
		for _, e := range *enterprises {
			specsArgs = append(specsArgs, strconv.FormatInt(e.EnterpriseID, 10))
		}
	}
	for _, spec := range specsArgs {
		enterpriseID, _ := strconv.ParseInt(spec, 10, 64)
		if err := UpdateHistoryUserTags(c, enterpriseID); err != nil {
			continue
		}
	}
	return nil
}

// UpdateHistoryUserTags 批量更新一次历史用户的标签
func UpdateHistoryUserTags(c ctx.Context, enterpriseID int64) error {
	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	externalTagService := customer.NewExternalTag(customerMysql.NewCustomer())
	tagService := tag.NewTag(mysql.NewTagRepo())
	// 用户tag集合
	checkUserTag := make(map[string]int64)
	usersTag := make(map[int64][]int64)
	log.Info(c).Msgf("当前正在执行企业: %+v，开始修复用户标签", enterpriseID)
	// 查询企业的所有标签
	tagIDs, _ := tagService.GetAllUsedExternalTags(c, enterpriseID)

	// 查询每个标签中的人
	for _, tagID := range tagIDs {
		uList, err := externalTagService.GetExternalUserUsedTag(c, enterpriseID, tagID)
		if err != nil {
			continue
		}
		for _, u := range uList {
			key := strconv.FormatInt(u.WorkExternalUserId, 10) + "-" + strconv.FormatInt(tagID, 10)
			if _, ok := checkUserTag[key]; !ok {
				usersTag[u.WorkExternalUserId] = append(usersTag[u.WorkExternalUserId], tagID)
				checkUserTag[key] = tagID
			}
		}
	}
	// 更新用户标签
	wg := &sync.WaitGroup{}

	ch := make(chan bool, 3) // 控制并发的数量
	for id, uTag := range usersTag {
		wg.Add(1)
		ch <- true
		go func(id int64, uTag []int64, ch chan bool) {
			defer EndTag(wg, ch)
			enterpriseUser, err := customerService.QueryEnterpriseUser(c, enterpriseID, &customer.SourceUserBase{
				ID:           id,
				CustomerType: customer.CustomerTypeExternal,
			})
			if err != nil || enterpriseUser == nil {
				if err == sql.ErrNoRows {
					// 对不存在的用户的用户关系记录进行删除
					externalTagService.DeleteUserTagsByUserID(ctx.Wrap(c), enterpriseID, id)
				}
				log.Error(c).Msgf("企业id：%+v，查找用户wx_external_user:%+v, 错误：%+v", enterpriseID, id, err)
				return
			}

			if err := customerService.UpdateEnterpriseUserAttr(ctx.Wrap(c), enterpriseID, enterpriseUser.ID, &customer.EnterpriseUserAttr{Tags: uTag}); err != nil {
				log.Info(c).Msgf("更新用户:%+v的标签, 错误：%+v", enterpriseUser.ID, err)
			}
		}(id, uTag, ch)

	}
	wg.Wait()
	log.Info(c).Msgf("执行企业: %+v，修复标签结束", enterpriseID)
	return nil
}

func EndTag(wg *sync.WaitGroup, ch chan bool) {
	<-ch
	wg.Done()
}

// FixHistoryUserRemark 用于同步enterprise_user表,历史用户的备注
func FixHistoryUserRemark(c ctx.Context, jobHandler *xjob.JobHandler) error {
	dbClient = db.DB()

	// 获取命令行参数配置的企业
	if err := SyncRemarkData(c); err != nil {
		return err
	}
	return nil
}

func SyncRemarkData(c ctx.Context) error {
	service := scripts.New(scriptsMysql.New())
	service.CustomerService = customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	// 查询关系表中用户有备注的人
	// 把备注更新到enterprise_user中
	return service.SyncRemarkData(c)
}
