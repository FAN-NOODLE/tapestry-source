package membersJob

import (
	"fmt"
	"strconv"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// 历史社群数据修复，只处理数据库中不存在的群的数据的新增
func FixHistoryGroupChat(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIDs := []int64{}
	isfix := false
	for _, param := range specsArgs {
		if enterpriseID, err := strconv.ParseInt(param, 10, 64); err == nil {
			enterpriseIDs = append(enterpriseIDs, enterpriseID)
		}
	}
	if len(enterpriseIDs) == 0 {
		log.Info(c).Msgf("企业id不能为空")
		return nil
	}
	if specsArgs[len(specsArgs)-1] == "true" {
		isfix = true
	}

	var wg sync.WaitGroup
	for _, enterpriseID := range enterpriseIDs {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()
			service := members.New(membersMysql.New())
			enterpriseService := enterprise.New(enterpriseMysql.New())
			domain, corpID, _, err := enterpriseService.GetInfo(c, enterpriseID)
			if err != nil {
				log.Info(c).Msgf("获取企业信息失败", err)
			}
			if isfix {
				err = service.FixAllGroupChatAndMember(c, enterpriseID, domain, corpID)
				if err != nil {
					log.Info(c).Msgf("同步企业群和成员信息失败", err)
				}
			} else {
				err = service.FixAllMemberChat(c, enterpriseID, domain, corpID)
				if err != nil {
					log.Info(c).Msgf("同步企业群信息失败", err)
				}
			}

		}(enterpriseID)
	}
	wg.Wait()
	return nil
}

// 修复指定chat_id的历史社群数据
func FixHistoryGroupChatByChatId(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业id，社群id
	specsArgs := jobHandler.Runner().Specs()
	if len(specsArgs) < 2 {
		fmt.Println("参数有误，格式为--job.spec=enterprise_id,chat_id")
		return nil
	}
	if specsArgs[0] == "" {
		fmt.Println("企业id不能为空，格式为--job.spec=enterprise_id,chat_id")
		return nil
	}
	if specsArgs[1] == "" {
		fmt.Println("社群id不能为空，格式为--job.spec=enterprise_id,chat_id")
		return nil
	}
	var enterpriseID int64
	enterpriseID, _ = strconv.ParseInt(specsArgs[0], 10, 64)
	service := members.New(membersMysql.New())
	err := service.FixByChatId(c, specsArgs[1], enterpriseID)
	if err != nil {
		log.Info(c).Msgf("同步企业群信息失败", err)
		return err
	}
	return nil
}
