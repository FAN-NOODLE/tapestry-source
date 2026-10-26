package membersJob

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"strings"
	"sync"
	"time"

	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"

	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/pkg/db"
)

// MemberJob 组织同步
func MemberJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	dbc := db.DB()
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIds := make([]int64, 0)
	if len(specsArgs) > 0 {
		for _, arg := range specsArgs {
			enterpriseId, _ := strconv.ParseInt(arg, 10, 64)
			enterpriseIds = append(enterpriseIds, enterpriseId)
		}
	}

	total, enterprises, err := enterprise.New(enterpriseMysql.New()).ListEnterprise(c, enterpriseIds, 1, 0)
	if err != nil || total <= 0 {
		return err
	}

	departmentService := department.New(membersMysql.NewDepartment())

	bgTime := time.Now().Unix()
	enterpriseM := make(map[int64]*enterprise.Corp)
	for _, vItem := range *enterprises {
		enterpriseM[vItem.EnterpriseID] = vItem
		history, locked := lockCommit(c, dbc, vItem)
		if !locked {
			continue
		}
		// 1. 查询企业应用信息
		req := boss.WorkPlugin{
			EnterpriseID: vItem.EnterpriseID,
		}

		plugins, err := boss.WorkPluginList(c, req)

		var plugin boss.WorkPlugin
		// 过滤侧边栏和通讯录应用插件
		for _, p := range plugins {
			// 判断是否为侧边栏
			if p.Identifier == "crs-manu-customer" {
				plugin = p
			}
			// 如果是通讯录则取通讯录插件
			if p.Identifier == "communalContact" {
				plugin = p
				break
			}
		}

		if err != nil || plugin.ID < 0 {
			log.Err(c, err).Msgf("企业id:%+v,没有配置: %v", vItem.EnterpriseID)
			continue
		}

		plugins = []boss.WorkPlugin{plugin}
		workImpl := new(work.QyWorkImpl).SetDomain(vItem.Domain)
		departs, dPlugin, dPhash := fetchDepartments(c, []boss.WorkPlugin{plugin}, workImpl)
		if len(departs) == 0 {
			log.Err(c, errs.WorkPluginErr).Msgf("企业id:%+v,应用信息:%+v,没有部门: %v", vItem.EnterpriseID, dPlugin)
			continue
		}

		// 查询db中的所有部门
		dbDeparts := make(map[int64]department.Department)
		allDepartmentFromDb, dErr := departmentService.FindAll(c, vItem.EnterpriseID)
		if dErr != nil {
			log.Err(c, dErr).Msgf("企业id:%+v,查询企业部门错误: %v", vItem.EnterpriseID)
			continue
		}

		for _, depDb := range allDepartmentFromDb {
			dbDeparts[depDb.ThirdID] = depDb
		}

		if dPhash != history.DepartHash {
			if dbDeparts, err = batchRefreshDepartments(c, vItem.EnterpriseID, dbDeparts, departs); err != nil {
				log.Err(c, err).Msgf("企业id:%+v,应用信息:%+v,批量插入部门失败: %v", vItem.EnterpriseID, dPlugin, err)
				continue
			}
		}

		// 查询成员信息
		membersWork, err := workImpl.WorkMembers(c, work.WorkMemberRequest{
			QyGwRequest:  work.QyGwRequest{ReqPL: plugin.Identifier},
			DepartmentID: 1,
			FetchChild:   1,
		})

		if err != nil || len(membersWork) == 0 {
			log.Err(c, err).Msgf("企业id:%+v,应用信息:%+v,查询部门成员失败: %+v", vItem.EnterpriseID, dPlugin)
			continue
		}

		b, _ := json.Marshal(membersWork)
		memberHash := helper.Md5V(string(b))

		if memberHash != history.MemberHash {
			if err = batchInsertMemberDepartment(c, vItem.EnterpriseID, membersWork, dbDeparts, workImpl); err != nil {
				log.Err(c, err).Msgf("企业id:%+v,应用信息:%+v", vItem.EnterpriseID, dPlugin)
				continue
			}
			checkSwSaCode(c, vItem.EnterpriseID)
		}

		// 同步员工和外部联系人关系
		if memberHash != history.MemberHash {
			syncMemberRelation(c, dbc, enterpriseM, bgTime, vItem.EnterpriseID, vItem.Domain)
		}

		// 缓存部门、子部门、子子部门
		if err := getDepartmentsAndChildDepartments(c, dbc, vItem.EnterpriseID); err != nil {
			log.Err(c, err).Msgf("缓存[%+v]组织架构部门关系", vItem.EnterpriseID)
		}

		if err := completeMemberSync(
			func() error {
				return gormx.GetInstance().WithContext(c).
					Table("work_member_sync").
					Where("enterprise_id = ?", vItem.EnterpriseID).
					Updates(map[string]any{
						"succ_at":          time.Now().Unix(),
						"succ_depart_hash": dPhash,
						"succ_member_hash": memberHash,
					}).Error
			},
			func() error {
				return syncHkStoreAfterMemberSync(
					c,
					vItem.EnterpriseID,
					boss.EnterpriseConfig,
					func(c ctx.Context, enterpriseID int64) error {
						_, err := dashboardSvc.NewHkStoreService().Sync(c, dashboardReq.HkStoreSync{
							EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: enterpriseID},
						})
						return err
					},
				)
			},
		); err != nil {
			log.Err(c, err).Msgf("企业id:%+v,组织架构同步完成处理失败", vItem.EnterpriseID)
		}
	}

	return nil
}

func completeMemberSync(recordSuccess, syncStore func() error) error {
	if err := recordSuccess(); err != nil {
		return fmt.Errorf("记录组织架构同步成功状态: %w", err)
	}
	if err := syncStore(); err != nil {
		return fmt.Errorf("同步港澳门店: %w", err)
	}
	return nil
}

func syncHkStoreAfterMemberSync(
	c ctx.Context,
	enterpriseID int64,
	getConfig func(ctx.Context, interface{}) (*boss.EnterpriseConfigResponse, error),
	syncStore func(ctx.Context, int64) error,
) error {
	configData, err := getConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"key":           "hkmo_store_root_dept_id",
	})
	if err != nil {
		return err
	}
	if configData == nil || len(configData.Data) == 0 || configData.Data[0].Value <= 0 {
		return nil
	}
	return syncStore(c, enterpriseID)
}

func fetchDepartments(c ctx.Context, plugins []boss.WorkPlugin, workImpl *work.QyWorkImpl) (map[int64]*department.Department, *boss.WorkPlugin, string) {
	for _, plugin := range plugins {
		departments, err := workImpl.WorkDepartmentList(c, work.WorkDepartmentListRequest{QyGwRequest: work.QyGwRequest{ReqPL: plugin.Identifier}})
		if err != nil {
			continue
		}

		departmentMap := make(map[int64]*department.Department)
		for _, depart := range departments {
			tmpDepart := department.Department{
				ID:        depart.ID,
				Name:      depart.Name,
				NameEn:    depart.NameEn,
				ParentID:  depart.Parentid,
				Order:     depart.Order,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}
			departmentMap[depart.ID] = &tmpDepart
		}

		if _, ok := departmentMap[1]; ok {
			for _, depart := range departmentMap {
				if parent, ok := departmentMap[depart.ParentID]; ok {
					parent.ChildrenCount++
				}
			}
			b, _ := json.Marshal(departments)
			return departmentMap, &plugin, helper.Md5V(string(b))
		}
	}

	return nil, nil, ""
}

func batchRefreshDepartments(c ctx.Context, enterpriseId int64, dbDeparts map[int64]department.Department, thirdDeparts map[int64]*department.Department) (map[int64]department.Department, error) {
	wg := &sync.WaitGroup{}
	pool, _ := ants.NewPool(10)
	var allThirdIds []int64
	for thirdId, depart := range thirdDeparts {
		allThirdIds = append(allThirdIds, thirdId)
		wg.Add(1)
		_ = pool.Submit(func(thId int64, dep department.Department) func() {
			return func() {
				defer wg.Done()
				// 部门已在数据库中存在
				if dbDepartInfo, ok := dbDeparts[thId]; ok && dbDepartInfo.ID > 0 {
					// 更新部门信息
					_ = departmentsDao.NewDepartments().Update(c, map[string]any{
						"third_id":        dep.ID,
						"third_parent_id": dep.ParentID,
						"name":            dep.Name,
						"name_en":         dep.NameEn,
						"order":           dep.Order,
						"updated_at":      time.Now().Unix(),
						"deleted_at":      0,
					}, map[string]any{
						"id":            dbDepartInfo.ID,
						"enterprise_id": enterpriseId,
					})
				} else {
					// 插入部门信息
					tmpArrDepartment := departmentsModels.DepartmentsModel{
						MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
						ThirdID:        int(dep.ID),
						ThirdParentID:  int(dep.ParentID),
						Name:           dep.Name,
						NameEn:         dep.NameEn,
						Order:          int(dep.Order),
					}
					_ = departmentsDao.NewDepartments().Create(c, tmpArrDepartment)
					dbDeparts[dep.ID] = department.Department{
						ID: tmpArrDepartment.ID,
					}
				}
			}
		}(thirdId, *depart))
	}
	wg.Wait()
	if len(thirdDeparts) > 0 {
		// 删除不存在的部门
		_ = departmentsDao.NewDepartments().DeleteData(c, "enterprise_id = ? AND third_id not in ?", enterpriseId, allThirdIds)

	}

	return dbDeparts, nil
}

func batchInsertMemberDepartment(c ctx.Context, enterpriseId int64, members []work.MemberDetail, dbDepartments map[int64]department.Department, workImpl *work.QyWorkImpl) error {
	// 查询db中的所有员工
	var allMembers []membersModel.MemberModel
	gormx.GetInstance().WithContext(c).Unscoped().Where("enterprise_id = ?", enterpriseId).Find(&allMembers)

	dbMembers := make(map[string]membersModel.MemberModel)
	var dbMemberUserIds []string
	for _, member := range allMembers {
		dbMembers[member.Userid] = member
		dbMemberUserIds = append(dbMemberUserIds, member.Userid)
	}
	var workMemberUserIds []string
	for _, member := range members {
		workMemberUserIds = append(workMemberUserIds, member.Userid)
	}

	// 取出数据库员工id和企微员工id差集删除
	delUserIds, _ := helper.ArrayDifference(dbMemberUserIds, workMemberUserIds)
	if len(delUserIds) > 0 {
		gormx.GetInstance().WithContext(c).Table("members").
			Where("enterprise_id = ?", enterpriseId).
			Where("userid in (?)", delUserIds).
			Where("deleted_at = 0").
			Updates(map[string]any{
				"status":     membersModel.MembersStatusLeave,
				"updated_at": time.Now().Unix(),
				"deleted_at": time.Now().Unix(),
			})

		//移除管理后台账号
		var delMobiles []string
		for _, delUserId := range delUserIds {
			delMobiles = append(delMobiles, dbMembers[delUserId].Mobile)
		}
		//处理因为存在手机号重复的账号管理后台账号被频繁删除的问题
		if len(delMobiles) > 0 {
			var existsMobile []string
			gormx.GetInstance().WithContext(c).Where("deleted_at = 0 AND enterprise_id = ? AND mobile in ?", enterpriseId, delMobiles).Pluck("mobile", &existsMobile)
			if len(existsMobile) > 0 {
				delMobiles, _ = helper.ArrayDifference(delMobiles, existsMobile)
			}
			if len(delMobiles) > 0 {
				_ = microBossDao.NewEnterpriseAdmin().BatchRemoveByMobile(c, delMobiles, enterpriseId)
			}
		}

		//解除所有正常的好友关系
		gormx.GetInstance().WithContext(c).Table("work_external_user_relation").
			Where("enterprise_id = ?", enterpriseId).
			Where("user_id in (?)", delUserIds).
			Where("relation_status = ?", microScrmModel.WorkUserRelationStatusNormal).
			Where("deleted_at = 0").
			Updates(map[string]any{
				"relation_status": microScrmModel.WorkUserRelationStatusDeleteEachOther,
				"updated_at":      time.Now().Unix(),
			})
	}

	wg := &sync.WaitGroup{}
	pool, _ := ants.NewPool(10)

	for _, member := range members {
		wg.Add(1)
		_ = pool.Submit(func(mem work.MemberDetail) func() {
			return func() {
				defer wg.Done()
				dbMemInfo, ok := dbMembers[mem.Userid]
				if ok && dbMemInfo.ID > 0 {
					// 是否需要更新
					needUpdate := false

					tmpUpdate := map[string]any{
						"userid":       mem.Userid,
						"name":         mem.Name,
						"mobile":       mem.Mobile,
						"position":     mem.Position,
						"gender":       int32(helper.StringToInt(mem.Gender)),
						"email":        mem.Email,
						"avatar":       mem.Avatar,
						"thumb_avatar": mem.ThumbAvatar,
						"telephone":    mem.Telephone,
						"alias":        mem.Alias,
						"status":       mem.Status,
						"qr_code":      mem.QrCode,
						"address":      mem.Address,
						//"open_userid":  mem.OpenUserid,
						"updated_at": time.Now().Unix(),
						"deleted_at": 0,
					}
					if len(mem.Extattr.Attrs) > 0 && mem.Extattr.Attrs[0].Name == "sid" && mem.Extattr.Attrs[0].Value != "" {
						tmpUpdate["open_userid"] = mem.Extattr.Attrs[0].Value
						if tmpUpdate["open_userid"] != dbMemInfo.OpenUserid {
							needUpdate = true
						}
					}
					if dbMemInfo.DeletedAt > 0 {
						tmpUpdate["type"] = membersModel.MembersTypeSA
						needUpdate = true
					}

					// 针对有更新的字段才更新
					if dbMemInfo.Name != mem.Name || dbMemInfo.Mobile != mem.Mobile || dbMemInfo.Position != mem.Position || dbMemInfo.Avatar != mem.Avatar || dbMemInfo.Status != int32(mem.Status) {
						needUpdate = true
					}

					if needUpdate {
						_ = membersDao.NewMembers().UpdateDataTable(c, map[string]any{
							"id":            dbMemInfo.ID,
							"enterprise_id": enterpriseId,
						}, tmpUpdate)
					}

				} else {
					dbMemInfo = membersModel.MemberModel{
						MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
						Userid:         mem.Userid,
						Name:           mem.Name,
						Mobile:         mem.Mobile,
						Position:       mem.Position,
						Gender:         int32(helper.StringToInt(mem.Gender)),
						Email:          mem.Email,
						Avatar:         mem.Avatar,
						ThumbAvatar:    mem.ThumbAvatar,
						Telephone:      mem.Telephone,
						Alias:          mem.Alias,
						Status:         int32(mem.Status),
						QrCode:         mem.QrCode,
						Address:        mem.Address,
						//OpenUserid:     mem.OpenUserid,
					}
					if len(mem.Extattr.Attrs) > 0 && mem.Extattr.Attrs[0].Name == "sid" && mem.Extattr.Attrs[0].Value != "" {
						dbMemInfo.OpenUserid = mem.Extattr.Attrs[0].Value
					}
					_ = membersDao.NewMembers().Create(c, dbMemInfo)
					//如果新增用户职位符合要求，添加【会话存档_正式】标签
					if helper.InArray(dbMemInfo.Position, membersDao.ConversationArchiveTag) {
						storeThirdIds := getStoreThirdIds(c, enterpriseId)
						if len(storeThirdIds) > 0 && len(mem.Department) > 0 && len(helper.ArrayIntersect(storeThirdIds, mem.Department)) > 0 {
							conf := config.CfgData.CoachSaTag
							_ = workImpl.AddTagUsers(c, work.AddTagUsersRequest{
								QyGwRequest: work.QyGwRequest{ReqPL: "communalContact"},
								TagID:       conf.TagId,
								UserList:    []string{mem.Userid},
							})
						}
					}
				}

				// 获取当前员工的部门信息id
				var memberDepartmentIds []int64
				gormx.GetInstance().WithContext(c).Table("departments").
					Select("id").
					Where("enterprise_id", enterpriseId).
					Where("third_id in (?)", mem.Department).
					Scan(&memberDepartmentIds)

				if len(memberDepartmentIds) > 0 {
					// 删除不在当前员工部门关系
					gormx.GetInstance().WithContext(c).Table("member_departments").
						Where("member_id = ? and department_id not in (?)", dbMemInfo.ID, memberDepartmentIds).
						Where("deleted_at = 0").
						Updates(map[string]any{
							"updated_at": time.Now().Unix(),
							"deleted_at": time.Now().Unix(),
						})
				}

				// 循环当前员工的部门
				for _, thirdId := range mem.Department {
					// 从数据库中获取部门信息
					if dbDepartment, ok := dbDepartments[thirdId]; ok {
						// 获取员工是否存在当前部门关系
						memberDepartment, _ := membersDao.NewMemberDepartments().GetInfoByMemberDepartment(c, enterpriseId, dbMemInfo.ID, dbDepartment.ID)
						if memberDepartment.ID == 0 {
							_ = membersDao.NewMemberDepartments().Create(c, membersModel.MemberDepartmentsModel{
								MicroBaseModel: models.CreateMicroBase(enterpriseId, app.UID()),
								MemberID:       dbMemInfo.ID,
								DepartmentID:   dbDepartment.ID,
							})
						}
						if memberDepartment.DeletedAt > 0 {
							_ = membersDao.NewMemberDepartments().UpdateData(c,
								map[string]any{"id": memberDepartment.ID},
								map[string]any{
									"updated_at": time.Now().Unix(),
									"deleted_at": 0,
								})
						}
					}
				}
			}
		}(member))
	}

	wg.Wait()

	return nil
}

// 同步新增的员工通讯录外部联系人关系
func syncMemberRelation(c ctx.Context, dbc *sqlplus.ClusterDB, enterpriseM map[int64]*enterprise.Corp, bgTime int64, enterpriseID int64, domain string) {
	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	userRepo := userMysql.NewUser()
	// 查询本次新增的员工，发起查询该员工的外部联系人
	addMembers := make([]members.SyncMemberRelationJob, 0)

	err := dbc.DB.Select(c, &addMembers, "select id,userid,enterprise_id from members where enterprise_id=? AND created_at>?", enterpriseID, bgTime)
	if len(addMembers) == 0 || err != nil {
		return
	}
	// 实例化服务
	workImpl := new(work.QyWorkImpl)
	workImpl.SetDomain(domain)
	for _, syncMemberRelationJob := range addMembers {
		if _, ok := enterpriseM[enterpriseID]; !ok {
			continue
		}
		syncMemberRelationJob.CorpID = enterpriseM[enterpriseID].CorpID
		cursor := ""
		for {
			externalUsers, err := workImpl.BatchExternalContactGet(c, work.BatchExternalContactRequest{
				MemberUserID: syncMemberRelationJob.Userid,
				Cursor:       cursor,
				Limit:        100,
			})
			if err != nil {
				log.Err(c, err).Msgf("syncMemberRelationQueue-企业id:%+v,获取成员列表失败,员工:%+v", syncMemberRelationJob.EnterpriseID, syncMemberRelationJob.Userid)
				break
			}
			wg := &sync.WaitGroup{}
			ch := make(chan bool, 100) // 控制并发的数量
			for _, euser := range externalUsers.ExternalContactList {
				wg.Add(1)
				ch <- true
				go func(externalUser work.WorkExternalContact, ch chan bool, wg *sync.WaitGroup) {
					defer EndGroup(wg, ch)
					source := &user.SourceUser{
						EnterpriseID: syncMemberRelationJob.EnterpriseID,
						SourceUserBase: customer.SourceUserBase{
							AppID:        syncMemberRelationJob.CorpID,
							OpenID:       externalUser.Base.ExternalUserID,
							UnionID:      externalUser.Base.UnionID,
							CustomerType: customer.CustomerTypeExternal,
						},

						EnterpriseUserBase: customer.EnterpriseUserBase{
							Name:   externalUser.Base.Name,
							Avatar: externalUser.Base.Avatar,
							Sex:    externalUser.Base.Gender,
						},

						EnterpriseUserExternal: customer.EnterpriseUserExternal{
							Position:        externalUser.Base.Position,
							CorpName:        externalUser.Base.CorpName,
							CorpFullName:    externalUser.Base.CorpFullName,
							Type:            externalUser.Base.Type,
							ExternalProfile: string(externalUser.Base.ExternalProfile),
						},

						CreateAt: time.Now().Unix(),
					}

					_, err := userRepo.MergeUser(c, source)
					if err != nil {
						log.Err(c, err).Msgf("syncMemberRelationQueue-企业id:%+v,外部联系人插入数据库失败: %+v", syncMemberRelationJob.EnterpriseID, source)
					}
					relation, err := customerService.FindExternalRelation(c, customer.WorkUserMemberRelation{EnterpriseID: syncMemberRelationJob.EnterpriseID, MemberID: syncMemberRelationJob.ID, WorkExternalUserID: source.ID})
					if err == sql.ErrNoRows {
						relation := &customer.WorkUserMemberRelation{
							ID:                 uid.ID(c),
							WorkExternalUserID: source.ID,
							MemberID:           syncMemberRelationJob.ID,
							CreatedAt:          time.Now().Unix(),
							UpdatedAt:          time.Now().Unix(),
							EnterpriseID:       syncMemberRelationJob.EnterpriseID,
							Remark:             externalUser.FollowInfo.Remark,
							Description:        externalUser.FollowInfo.Description,
							CreateTime:         externalUser.FollowInfo.CreateTime,
							RemarkCorpName:     externalUser.FollowInfo.RemarkCorpName,
							RemarkMobiles:      strings.Join(externalUser.FollowInfo.RemarkMobiles, ","),
							AddWay:             externalUser.FollowInfo.AddWay,
							OperUserID:         externalUser.FollowInfo.OperUserID,
							RelationStatus:     customer.WorkUserRelationStatusNormal,
							UserId:             syncMemberRelationJob.Userid,
							UnionId:            externalUser.Base.UnionID,
						}
						err := customerService.AddExternalRelation(c, []customer.WorkUserMemberRelation{*relation})

						if err != nil {
							log.Err(c, err).Msgf("syncMemberRelationQueue-企业id:%+v,外部联系人关系插入数据库失败: %+v", syncMemberRelationJob.EnterpriseID, relation)
						}

					} else if err == nil {
						if relation.RelationStatus != customer.WorkUserRelationStatusNormal {
							relation.RelationStatus = customer.WorkUserRelationStatusNormal // 回调新增关系状态设置为0
						}
						relation.UpdatedAt = time.Now().Unix()
						err := customerService.UpdateExternalRelation(c, *relation)

						if err != nil {
							log.Err(c, err).Msgf("syncMemberRelationQueue-企业id:%+v,更新外部联系人关系数据库失败: %+v", syncMemberRelationJob.EnterpriseID, relation)
						}
					}
				}(euser, ch, wg)
			}
			wg.Wait()

			if externalUsers.NextCursor == "" {
				break
			}
			cursor = externalUsers.NextCursor
		}
	}

	return
}

func EndGroup(wg *sync.WaitGroup, ch chan bool) {
	<-ch
	wg.Done()
}

// 获取部门，及其子部门，子子部门数据
func getDepartmentsAndChildDepartments(c ctx.Context, mdb *sqlplus.ClusterDB, enterpriseId int64) error {
	var (
		// mdb            = db.DB()
		dps            []*Department
		thirdToDpMap   = make(map[int]int64) // thirdPid => departmentId 一个thirdId对应的一个部门
		belongToFather = make(map[int][]int) // father => []children  一个父thirdId对应的多个子thirdId
	)
	if err := mdb.Select(ctx.Wrap(c), &dps, `select id dpid,third_id tid,third_parent_id tpid from departments where enterprise_id=? and deleted_at=0`, enterpriseId); err != nil {
		return err
	}

	// 映射有子部门
	for _, v := range dps {
		thirdToDpMap[v.Tid] = v.DpId
		if v.TPid > 0 {
			belongToFather[v.TPid] = append(belongToFather[v.TPid], v.Tid)
		}
	}

	var fatherMap = make(map[int]map[int]int)
	for father, children := range belongToFather {
		var childMap = make(map[int]int)
		childMap[father] = father
		getChildren(childMap, children, belongToFather)

		fatherMap[father] = childMap // 所有子部门，子子部门
	}

	var dpMap = make(map[int64][]int64)
	for father, children := range fatherMap {
		for _, child := range children {
			dpMap[thirdToDpMap[father]] = append(dpMap[thirdToDpMap[father]], thirdToDpMap[child])
		}
	}

	// 没有子部门的也放进来
	for third, dpId := range thirdToDpMap {
		if _, ok := dpMap[thirdToDpMap[third]]; !ok {
			dpMap[thirdToDpMap[third]] = append(dpMap[thirdToDpMap[third]], dpId)
		}
	}

	wg := &sync.WaitGroup{}
	pool, _ := ants.NewPool(10)

	for father, children := range dpMap {
		wg.Add(1)
		_ = pool.Submit(func(fa int64, ch []int64) func() {
			return func() {
				defer wg.Done()
				childByte, _ := json.Marshal(ch)
				_, err := mdb.DB.Exec(c, `update departments set child_departments = ? where id = ? and enterprise_id = ?`, string(childByte), fa, enterpriseId)
				if err != nil {
					log.Error(c).Msgf("update departments child_departments error %+v", err)
				}
			}
		}(father, children))
	}

	wg.Wait()

	return nil
}

func getChildren(fatherMap map[int]int, children []int, belongToFather map[int][]int) {
	for _, third := range children {
		fatherMap[third] = third
		getChildren(fatherMap, belongToFather[third], belongToFather)
	}
}

type syncHistory struct {
	SuccAt     int64  `db:"succ_at"`
	CommitAt   int64  `db:"commit_at"`
	SyncAt     int64  `db:"sync_at"`
	DepartHash string `db:"succ_depart_hash"`
	MemberHash string `db:"succ_member_hash"`
}

type Department struct {
	DpId int64 `db:"dpid"`
	Tid  int   `db:"tid"`
	TPid int   `db:"tpid"`
}

func lockCommit(c ctx.Context, dbc *sqlplus.ClusterDB, item *enterprise.Corp) (*syncHistory, bool) {
	nowtm := time.Now().Unix()
	var history syncHistory
	if err := dbc.DB.Get(c, &history, "select succ_at, commit_at, sync_at, succ_depart_hash, succ_member_hash from work_member_sync where enterprise_id=?", item.EnterpriseID); err == sql.ErrNoRows {
		if _, err = dbc.DB.Exec(c, "insert into work_member_sync(enterprise_id, sync_at) values(?, ?)", item.EnterpriseID, nowtm); err != nil {
			return nil, false
		}
		return &history, true
	} else if err != nil {
		log.Error(c).Msgf("企业id:%+v,获取历史同步记录错误: %v", item.EnterpriseID, errs.WorkPluginErr)
		return nil, false
	} else {
		// 如果最近同步时间超过1min 并且超过20分钟没有同步 或 手动请求同步
		if nowtm-history.SyncAt > 60 && (nowtm-history.SuccAt > 20*60 || history.CommitAt > history.SuccAt) {
			res, err := dbc.DB.Exec(c, "update work_member_sync set sync_at=? where enterprise_id=? and sync_at=?", nowtm, item.EnterpriseID, history.SyncAt)
			if err != nil {
				log.Error(c).Msgf("企业id:%+v,更新历史同步记录错误: %v", item.EnterpriseID, errs.WorkPluginErr)
				return nil, false
			}

			affected, err := res.RowsAffected()
			if err != nil {
				log.Error(c).Msgf("企业id:%+v,更新历史同步记录错误: %v", item.EnterpriseID, errs.WorkPluginErr)
				return nil, false
			}

			return &history, affected > 0
		}
	}
	return nil, false
}

func checkSwSaCode(c ctx.Context, enterpriseId int64) {
	cacheRes, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if cacheRes.Name != "sw" {
		return
	}
	var allMembers []membersModel.MemberModel
	gormx.GetInstance().WithContext(c).Unscoped().Where("enterprise_id = ?", enterpriseId).Find(&allMembers)
	seen := make(map[string]bool)

	for _, item := range allMembers {
		if item.OpenUserid == "" {
			continue
		}
		if seen[item.OpenUserid] {
			notice.NewWechat().SendTextMessage(c, fmt.Sprintf("Sw员工工号有重复：%s", item.OpenUserid))
		}
		seen[item.OpenUserid] = true
	}
}
func getStoreThirdIds(c ctx.Context, enterpriseId int64) []int64 {
	var storeThirdIds []int64
	// 获取企业配置的门店父级部门id信息
	configData, _ := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"key":           "store_root_dept_id",
	})
	storeParentThirdId := int(configData.Data[0].Value) //所有门店的父级企微部门id
	departments, _ := departmentsDao.NewDepartments().GetListByThirdParentId(c, enterpriseId, storeParentThirdId)
	if len(departments) > 0 {
		for _, storeDept := range departments {
			storeThirdIds = append(storeThirdIds, int64(storeDept.ThirdID))
		}
	}
	return storeThirdIds
}
