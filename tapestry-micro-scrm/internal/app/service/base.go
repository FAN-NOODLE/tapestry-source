package service

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	"micro-scrm/internal/app/http/request"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/logger"
	"micro-scrm/internal/pkg/uid"
)

type Base struct {
}

func (s Base) Log() *logger.Logger {
	return app.Logger()
}

func (s Base) Redis() *redis.Client {
	return app.Redis()
}

func (s Base) DB() *gorm.DB {
	return app.DB()
}

// GetUserIdByMobile 获取员工信息
func (s Base) GetUserIdByMobile(c ctx.Context, enterpriseId int64, mobile string) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	// 查询员工企微id
	membersInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Mobile:         mobile,
	})
	if err != nil || membersInfo.ID == 0 {
		return result, errors.New("员工信息未查到")
	}
	return membersInfo, nil
}

func (s Base) GetUserIdByUserId(c ctx.Context, enterpriseId int64, userId string) (membersModel.MemberModel, error) {
	var result membersModel.MemberModel
	// 查询员工企微id
	membersInfo, err := membersDao.NewMembers().GetOne(c, membersModel.MemberModel{
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: enterpriseId},
		Userid:         userId,
	})
	if err != nil || membersInfo.ID == 0 {
		return result, errors.New("员工信息未查到")
	}
	return membersInfo, nil
}

func (s Base) GetChannelId() int64 {
	return uid.ChannelId.GetChannelId()
}

// GetSearchMemberIds 管理后台导购部门相关参数通用检索条件处理
func (s Base) GetSearchMemberIds(c ctx.Context, params request.MemberSearchReq) (bool, []int64) {
	searchExists := true // 用户标记是否还要查询后续检索条件
	if params.UserId != "" {
		params.MemberIdArr, _ = membersDao.NewMembers().GetMemberIdsByUserIdStr(c, params.EnterpriseId, params.UserId)
		if len(params.MemberIdArr) == 0 {
			searchExists = false
		}
	}
	if searchExists && params.UserName != "" {
		tmpMemberIdArr, _ := membersDao.NewMembers().GetMemberIdsByName(c, params.EnterpriseId, params.UserName)
		if len(params.MemberIdArr) > 0 {
			params.MemberIdArr = helper.ArrayIntersect(params.MemberIdArr, tmpMemberIdArr)
		} else {
			params.MemberIdArr = tmpMemberIdArr
		}
		if len(params.MemberIdArr) == 0 {
			searchExists = false
		}
	}
	if searchExists && params.ThirdId > 0 {
		departmentIds, _ := departmentsDao.NewDepartments().GetDepartmentIdsByThirdIds(c, params.EnterpriseId, []int64{params.ThirdId})
		if len(departmentIds) > 0 {
			tmpMemberIdArr, _ := membersDao.NewMemberDepartments().GetMemberIds(c, params.EnterpriseId, departmentIds)
			if len(params.MemberIdArr) > 0 {
				params.MemberIdArr = helper.ArrayIntersect(params.MemberIdArr, tmpMemberIdArr)
			} else {
				params.MemberIdArr = tmpMemberIdArr
			}
		}
		if len(departmentIds) == 0 || len(params.MemberIdArr) == 0 {
			searchExists = false
		}
	}
	return searchExists, params.MemberIdArr
}
