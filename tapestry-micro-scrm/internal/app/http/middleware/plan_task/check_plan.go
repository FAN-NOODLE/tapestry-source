package plan_task

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

const MemberId = "member_id"
const UserId = "user_id"

type checkMemberInfo struct {
	EnterpriseId   int64  `json:"enterprise_id"`
	Mobile         string `json:"mobile"`
	TaskAllotLogId int64  `json:"task_allot_log_id"`
}

// CheckConditionAuth sa操作task_condition_data表的数据验证
func CheckConditionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据Authorization判断是否需要走授权
		var request checkMemberInfo
		// 用shouldbindbodywith的原因是因为要验证请求body的数据,处理err为EOF的报错
		if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
			ginplus.ResError(c, errors.New("无效的token信息"))
			c.Abort()
			return
		}
		memberInfo, err := membersDao.NewMembers().GetMemberByMobile(ctx.Wrap(c), request.EnterpriseId, request.Mobile)
		if err != nil {
			ginplus.ResError(c, errors.New("导购身份信息有误"))
			c.Abort()
			return
		}
		taskAllotLogInfo, err := planDao.NewPlanTaskAllotLogDao().GetInfoById(ctx.Wrap(c), request.EnterpriseId, request.TaskAllotLogId)
		if err != nil {
			ginplus.ResError(c, errors.New("无效的参数"))
			c.Abort()
			return
		}
		taskInfo, err := planDao.NewPlanTaskDao().GetOne(ctx.Wrap(c), request.EnterpriseId, map[string]interface{}{
			"id": taskAllotLogInfo.TaskId,
		})
		if err != nil {
			ginplus.ResError(c, errors.New("子任务不存在"))
			c.Abort()
			return
		}
		if taskInfo.Status == microScrmModel.PlanEnd { // 子任务已过期
			ginplus.ResError(c, errors.New("子任务已结束"))
			c.Abort()
			return
		}
		c.Set(MemberId, memberInfo.Id)
		c.Set(UserId, memberInfo.Userid)
		c.Next()
	}
}
