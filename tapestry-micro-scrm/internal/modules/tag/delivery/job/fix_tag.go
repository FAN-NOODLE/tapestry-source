package job

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/xjob"
	"sync"
)

func FixTag(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
	// 赫莲娜企业id（生产执行需要改）
	var enterpriseID int64 = 3314240094228544
	// 打上指定标签（生产执行需要改）
	var tagId int64 = 3332149564102720
	page := 1
	size := 1000
	db := gormx.GetInstance()
	var waitGroup sync.WaitGroup
	for {
		var relationList []int64
		var newTagData []work_external_user.WorkExternalUserTagModel
		db.Debug().Raw("SELECT id FROM work_external_user_relation WHERE relation_status = 0 and created_at>=1627056000  and enterprise_id=? LIMIT ? OFFSET ?", enterpriseID, size, (page-1)*size).Scan(&relationList)
		if len(relationList) == 0 {
			break
		}
		for _, relation := range relationList {
			waitGroup.Add(1)
			go func(rID int64) {
				defer waitGroup.Done()
				isExistParam := work_external_user.WorkExternalUserTagModel{
					MicroScrmBaseModel: model.MicroScrmBaseModel{
						EnterpriseID: enterpriseID,
					},
					RelationID:    rID,
					ExternalTagID: tagId,
				}
				isExistParam.ID = uid.ID(ctx)
				newTagData = append(newTagData, isExistParam)

			}(relation)
			waitGroup.Wait()
		}
		if len(newTagData) != 0 {
			orm.Service.WorkExternalUser.InsertUserTag(newTagData)
		}
		page++
	}

	return nil
}
