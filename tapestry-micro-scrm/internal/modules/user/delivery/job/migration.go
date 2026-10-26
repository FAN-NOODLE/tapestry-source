package job

import (
	"database/sql"
	"fmt"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"strings"
	"sync"
	"time"

	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"

	"strconv"

	"micro-scrm/internal/pkg/db"
)

var dbClient *sqlplus.ClusterDB
var userRepo user.Repository

type EnterpriseInfo struct {
	ShopAppID string
	CorpID    string
}

// FixExternalContact 修复外部联系人详情
func FixExternalContact(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIDs := []int64{}
	for _, param := range specsArgs {
		if enterpriseID, err := strconv.ParseInt(param, 10, 64); err == nil {
			enterpriseIDs = append(enterpriseIDs, enterpriseID)
		}
	}
	dbClient = db.DB()
	enterpriseService := enterprise.New(enterpriseMysql.New())
	_, enterprises, err := enterpriseService.ListEnterprise(c, enterpriseIDs, 1, 0)
	if err != nil {
		return err
	}
	if len(*enterprises) == 0 {
		return nil
	}

	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	memberService := members.New(membersMysql.New())
	userRepo = userMysql.NewUser()

	for _, vItem := range *enterprises {
		// 执行策略
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			workImpl := new(work.QyWorkImpl)
			workImpl.SetDomain(vItem.Domain)
			list, err := memberService.All(c, vItem.EnterpriseID)
			if err != nil {
				log.Err(c, err).Msgf("FixExternalContact-%+v获取成员列表失败: %+v", vItem.EnterpriseID, vItem)
				return
			}
			beg := time.Now()
			total := 0
			for _, m := range list {
				cursor := ""
				for {
					externalUsers, err := workImpl.BatchExternalContactGet(c, work.BatchExternalContactRequest{
						MemberUserID: m.Userid,
						Cursor:       cursor,
						Limit:        100,
					})
					if err != nil {
						log.Err(c, err).Msgf("FixExternalContact-%+v获取成员列表失败: %+v,员工:%+v", vItem.EnterpriseID, vItem, m.Userid)
						break
					}
					total += len(externalUsers.ExternalContactList)
					wg := &sync.WaitGroup{}
					ch := make(chan bool, 100) // 控制并发的数量
					for _, euser := range externalUsers.ExternalContactList {
						wg.Add(1)
						ch <- true
						go func(externalUser work.WorkExternalContact, ch chan bool, wg *sync.WaitGroup) {
							defer EndGroup(wg, ch)
							source := &user.SourceUser{
								EnterpriseID: vItem.EnterpriseID,
								SourceUserBase: customer.SourceUserBase{
									AppID:        vItem.CorpID,
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
								log.Err(c, err).Msgf("FixExternalContact-%+v外部联系人插入数据库失败: %+v", vItem.EnterpriseID, source)
							}
							relation, err := customerService.FindExternalRelation(c, customer.WorkUserMemberRelation{EnterpriseID: vItem.EnterpriseID, MemberID: m.ID, WorkExternalUserID: source.ID})
							if err == sql.ErrNoRows {
								relation := &customer.WorkUserMemberRelation{
									ID:                 uid.ID(c),
									WorkExternalUserID: source.ID,
									MemberID:           m.ID,
									CreatedAt:          time.Now().Unix(),
									UpdatedAt:          time.Now().Unix(),
									EnterpriseID:       vItem.EnterpriseID,
									Remark:             externalUser.FollowInfo.Remark,
									Description:        externalUser.FollowInfo.Description,
									CreateTime:         externalUser.FollowInfo.CreateTime,
									RemarkCorpName:     externalUser.FollowInfo.RemarkCorpName,
									RemarkMobiles:      strings.Join(externalUser.FollowInfo.RemarkMobiles, ","),
									AddWay:             externalUser.FollowInfo.AddWay,
									OperUserID:         externalUser.FollowInfo.OperUserID,
									RelationStatus:     customer.WorkUserRelationStatusNormal,
									UserId:             m.Userid,
									UnionId:            externalUser.Base.UnionID,
								}
								err := customerService.AddExternalRelation(c, []customer.WorkUserMemberRelation{*relation})

								if err != nil {
									log.Err(c, err).Msgf("FixExternalContact-%+v外部联系人关系插入数据库失败: %+v", vItem.EnterpriseID, relation)
								}

							} else if err == nil {
								if relation.RelationStatus != customer.WorkUserRelationStatusNormal {
									relation.RelationStatus = customer.WorkUserRelationStatusNormal // 回调新增关系状态设置为0
								}
								relation.UpdatedAt = time.Now().Unix()
								err := customerService.UpdateExternalRelation(c, *relation)

								if err != nil {
									log.Err(c, err).Msgf("FixExternalContact-%+v更新外部联系人关系数据库失败: %+v", vItem.EnterpriseID, relation)
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

			log.Info(c).Msgf("FixExternalContact-%+v，从企业微信获取到的外部联系人关系总数: %+v", vItem.EnterpriseID, total)
			_ = customerService.BatchDeleteUpdateExternalReleation(c, vItem.EnterpriseID, beg.Unix()-1)
		})
	}
	return nil
}

type memberS struct {
	ID           int64  `json:"mem_id" db:"mem_id"`
	Userid       string `json:"userid" db:"userid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
}

// FixExternalContactByUserID 指定员工修复其外部联系人
func FixExternalContactByUserID(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	var enterpriseIDs []int64
	for _, param := range specsArgs {
		if enterpriseID, err := strconv.ParseInt(param, 10, 64); err == nil {
			enterpriseIDs = append(enterpriseIDs, enterpriseID)
		}
	}
	dbClient = db.DB()
	enterpriseService := enterprise.New(enterpriseMysql.New())
	_, enterprises, err := enterpriseService.ListEnterprise(c, enterpriseIDs, 1, 0)
	if err != nil {
		return err
	}
	if len(*enterprises) == 0 {
		return nil
	}

	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	userRepo = userMysql.NewUser()

	for _, vItem := range *enterprises {
		// 执行策略
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			workImpl := new(work.QyWorkImpl)
			workImpl.SetDomain(vItem.Domain)
			// 查询指定的员工
			var list []memberS
			sqlStr := fmt.Sprintf("select enterprise_id,userid,mem_id from enterprise_fixed_members where enterprise_id=?")
			err := dbClient.DB.Select(c, &list, sqlStr, vItem.EnterpriseID)
			if err != nil {
				log.Err(c, err).Msgf("FixExternalContactByUserID-%+v获取成员列表失败: %+v", vItem.EnterpriseID, vItem)
				return
			}
			beg := time.Now()
			total := 0
			for _, m := range list {
				cursor := ""
				isErr := false
				for {
					externalUsers, err := workImpl.BatchExternalContactGet(c, work.BatchExternalContactRequest{
						MemberUserID: m.Userid,
						Cursor:       cursor,
						Limit:        100,
					})
					if err != nil {
						isErr = true
						errStr := fmt.Sprintf("%+v", err)
						upStr := fmt.Sprintf("update enterprise_fixed_members set err = ? where userid = ?")
						_, _ = dbClient.DB.Exec(c, upStr, errStr, m.Userid)
						log.Err(c, err).Msgf("FixExternalContact-%+v获取成员列表失败: %+v,员工:%+v", vItem.EnterpriseID, vItem, m.Userid)
						break
					}
					total += len(externalUsers.ExternalContactList)
					wg := &sync.WaitGroup{}
					ch := make(chan bool, 100) // 控制并发的数量
					for _, euser := range externalUsers.ExternalContactList {
						wg.Add(1)
						ch <- true
						go func(externalUser work.WorkExternalContact, ch chan bool, wg *sync.WaitGroup) {
							defer EndGroup(wg, ch)
							source := &user.SourceUser{
								EnterpriseID: vItem.EnterpriseID,
								SourceUserBase: customer.SourceUserBase{
									AppID:        vItem.CorpID,
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
								log.Err(c, err).Msgf("FixExternalContact-%+v外部联系人插入数据库失败: %+v", vItem.EnterpriseID, source)
							}
							relation, err := customerService.FindExternalRelation(c, customer.WorkUserMemberRelation{EnterpriseID: vItem.EnterpriseID, MemberID: m.ID, WorkExternalUserID: source.ID})
							if err == sql.ErrNoRows {
								relation := &customer.WorkUserMemberRelation{
									ID:                 uid.ID(c),
									WorkExternalUserID: source.ID,
									MemberID:           m.ID,
									CreatedAt:          time.Now().Unix(),
									UpdatedAt:          time.Now().Unix(),
									EnterpriseID:       vItem.EnterpriseID,
									Remark:             externalUser.FollowInfo.Remark,
									Description:        externalUser.FollowInfo.Description,
									CreateTime:         externalUser.FollowInfo.CreateTime,
									RemarkCorpName:     externalUser.FollowInfo.RemarkCorpName,
									RemarkMobiles:      strings.Join(externalUser.FollowInfo.RemarkMobiles, ","),
									AddWay:             externalUser.FollowInfo.AddWay,
									OperUserID:         externalUser.FollowInfo.OperUserID,
									RelationStatus:     customer.WorkUserRelationStatusNormal,
									UserID:             m.Userid,
									UnionId:            externalUser.Base.UnionID,
								}
								err := customerService.AddExternalRelation(c, []customer.WorkUserMemberRelation{*relation})

								if err != nil {
									log.Err(c, err).Msgf("FixExternalContact-%+v外部联系人关系插入数据库失败: %+v", vItem.EnterpriseID, relation)
								}

							} else if err == nil {
								//如果是顾客单方面删除则不更新，因为只有回调才能体现单方面删除
								if relation.RelationStatus != customer.WorkUserRelationStatusNormal && relation.RelationStatus != customer.WorkUserRelationStatusDeletedByExternal {
									relation.RelationStatus = customer.WorkUserRelationStatusNormal // 回调新增关系状态设置为0
								}
								relation.UpdatedAt = time.Now().Unix()
								err := customerService.UpdateExternalRelation(c, *relation)

								if err != nil {
									log.Err(c, err).Msgf("FixExternalContact-%+v更新外部联系人关系数据库失败: %+v", vItem.EnterpriseID, relation)
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
				if !isErr {
					upStr := fmt.Sprintf("update enterprise_fixed_members set err = ? where userid = ?")
					_, _ = dbClient.DB.Exec(c, upStr, "已完成", m.Userid)
					_ = microScrmDao.NewUserRelation().BatchDeleteExternalRelationByMemberIds(c, vItem.EnterpriseID, beg.Unix()-1, []int64{m.ID})
				}
			}

			log.Info(c).Msgf("FixExternalContact-%+v，从企业微信获取到的外部联系人关系总数: %+v", vItem.EnterpriseID, total)
		})
	}
	return nil
}
