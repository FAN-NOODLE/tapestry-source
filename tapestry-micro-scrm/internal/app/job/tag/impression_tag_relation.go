package tagJob

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	appJob "micro-scrm/internal/app/job"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/xjob"
	"sync"
)

type ImpressionTagRelation struct {
	name string
}

func NewImpressionTagRelation() ImpressionTagRelation {
	return ImpressionTagRelation{
		name: "印象标签关系",
	}
}

func (job ImpressionTagRelation) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	wg := threading.NewRoutineGroup()

	for _, enterpriseId := range enterpriseIds {
		id := enterpriseId
		wg.RunSafe(func() {
			err = job.handle(c, id)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job ImpressionTagRelation) handle(c ctx.Context, enterpriseId int64) error {
	// 获取印象标签关系
	impressionRelation, err := tagsDao.NewImpressionRelation().GetMany(c, map[string]any{"enterprise_id": enterpriseId})
	if err != nil {
		return err
	}

	fmt.Println(len(impressionRelation))

	wg := &sync.WaitGroup{}

	pool, _ := ants.NewPool(10)

	for _, relation := range impressionRelation {
		wg.Add(1)
		_ = pool.Submit(func(data tagModels.ImpressionTagRelationModel) func() {
			return func() {
				defer wg.Done()
				external, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"euid": data.ExternalUserID}, enterpriseId)
				if err != nil || external.ID == 0 {
					app.Logger().Error(c).Msgf("%s 获取外部联系人表异常")
					return
				}
				data.ExternalUserID = external.ID
				err = tagsDao.NewImpressionRelation().SaveData(c, data)
				if err != nil {
					app.Logger().Error(c).Msgf("%s 更新关系异常")
					return
				}
			}
		}(relation))
	}

	wg.Wait()

	return nil
}
