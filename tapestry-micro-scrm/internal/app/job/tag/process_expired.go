package tagJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	tagsDao "micro-scrm/internal/app/dao/tags"
	appJob "micro-scrm/internal/app/job"
	workModels "micro-scrm/internal/app/models/work"
	workSDK "micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type ProcessExpired struct {
	name string
}

func NewProcessExpired() ProcessExpired {
	return ProcessExpired{
		name: "处理过期标签",
	}
}

func (job ProcessExpired) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, enterpriseId := range enterpriseIds {
		wg.Run(func(id int64) func() {
			return func() {
				err = job.handle(c, id)
				if err != nil {
					app.Logger().Error(c).Msgf("error : %v", err)
				}
			}
		}(enterpriseId))
	}

	wg.Wait()

	return nil
}

func (job ProcessExpired) handle(c ctx.Context, enterpriseId int64) error {
	nowTime := time.Now().Unix()
	// 查询已过期的标签
	workExternalTags, err := tagsDao.NewWorkExternalTag().
		GetMany(c, "enterprise_id = ? and expired_time > 0 and expired_time < ? and type in (?)", enterpriseId, nowTime, []int{workModels.WorkExternalTagTypeQy, workModels.MomentTaskTempTag})
	if err != nil {
		return err
	}

	if len(workExternalTags) == 0 {
		return nil
	}

	for _, tag := range workExternalTags {
		workWeChatObj, workErr := workwechat.NewWorkWeChat(c, enterpriseId)
		if workErr != nil {
			continue
		}
		// 删除企业微信标签
		workErr = workWeChatObj.DeleteTag(workSDK.DelCorpTagParam{TagID: []string{tag.TagID}})
		if workErr != nil {
			app.Logger().Error(c).Msgf("删除企业微信标签失败，错误信息：%+v", workErr)
		}

		// 删除数据库标签
		workErr = tagsDao.NewWorkExternalTag().DeleteData(c, map[string]any{
			"enterprise_id": enterpriseId,
			"id":            tag.ID,
		})
		if workErr != nil {
			app.Logger().Error(c).Msgf("同步标签删除标签失败，错误信息：%+v", workErr)
			continue
		}

		// 删除外部联系人标签关系
		workErr = tagsDao.NewWorkExternalUserTag().DeleteData(c, map[string]any{
			"enterprise_id":   enterpriseId,
			"external_tag_id": tag.ID,
		})
		if workErr != nil {
			app.Logger().Error(c).Msgf("删除外部联系人标签关系失败，错误信息：%+v", workErr)
			continue
		}

		// 获取当前标签的标签组的标签
		groupTags, workErr := tagsDao.NewWorkExternalTag().GetMany(c, "enterprise_id = ? and group_id = ?", enterpriseId, tag.GroupID)
		if workErr != nil {
			app.Logger().Error(c).Msgf("获取当前标签的标签组的标签失败，错误信息：%+v", workErr)
			continue
		}

		// 如果标签组下不存在标签则删除标签组(企微:标签组下所有的标签均被删除，则标签组会被自动删除)
		if len(groupTags) == 0 {
			workErr = tagsDao.NewWorkExternalTagGroup().DeleteData(c, map[string]any{
				"enterprise_id": enterpriseId,
				"id":            tag.GroupID,
			})
			app.Logger().Error(c).Msgf("删除标签组失败，错误信息：%+v", workErr)
			continue
		}
	}

	return nil
}
