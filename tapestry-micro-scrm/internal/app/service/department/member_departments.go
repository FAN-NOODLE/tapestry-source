package department

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type MemberDepartmentsService struct {
	service.Base
	memberDepartmentsDao *membersDao.MemberDepartments
}

func NewMemberDepartmentsService() MemberDepartmentsService {
	return MemberDepartmentsService{
		memberDepartmentsDao: membersDao.NewMemberDepartments(),
	}
}

func (s MemberDepartmentsService) GetOneByMemberID(c app.Ctx, enterpriseId int64, memberId int64) (error, membersModel.MemberDepartmentsModel) {
	return s.memberDepartmentsDao.GetOneByMemberID(c, enterpriseId, memberId)
}

func (s MemberDepartmentsService) GetListsByMemberID(c app.Ctx, enterpriseId int64, memberId int64) (error, []membersModel.MemberDepartmentsModel) {
	return s.memberDepartmentsDao.GetListsByMemberID(c, enterpriseId, memberId)
}

// GetMemberIdsByStoreType 获取奥莱或者零售的导购id集合
func (s MemberDepartmentsService) GetMemberIdsByStoreType(c app.Ctx, enterpriseId int64, storeType int64) ([]int64, error) {
	//获取奥莱与零售的用户memberId集合
	memberIdRedisKey := fmt.Sprintf("excitation:memberId:%d-%d", enterpriseId, storeType)
	memberIdsJson, err := s.Redis().Get(c, memberIdRedisKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	var memberIds []int64 //奥莱或零售的所有memberId集合（因为是定时任务全企业批量处理 避免每个用户都查一次是奥莱还是零售）
	if memberIdsJson != nil {
		//有缓存直接拿缓存
		err := json.Unmarshal(memberIdsJson, &memberIds)
		if err != nil {
			return nil, err
		}
	} else {
		//无缓存查库获取
		storeIds, err := microWorkDao.NewStoreDepartment().GetDpIdsByType(c, enterpriseId, storeType)
		if err != nil {
			return nil, err
		}
		// 通过企微部门id查询本地部门详情
		departmentInfoList, _ := departmentsDao.NewDepartments().GetDepartmentsThirdIds(c, enterpriseId, storeIds)
		var dpIds []int64 //departments表id集合
		for _, dpInfo := range departmentInfoList {
			dpIds = append(dpIds, dpInfo.ID)
		}
		memberIds, err = membersDao.NewMemberDepartments().GetMemberIds(c, enterpriseId, dpIds)
		if err != nil {
			return nil, err
		}
		memberIdsByte, err := json.Marshal(memberIds)
		if err != nil {
			return nil, err
		}
		s.Redis().Set(c, memberIdRedisKey, string(memberIdsByte), time.Hour*12)
	}
	return memberIds, nil
}

// GetMemberStoreType 判断用户是奥莱或者零售的导购
func (s MemberDepartmentsService) GetMemberStoreType(c app.Ctx, enterpriseId int64, memberId int64) (int64, error) {
	var storeTypeArr = []int64{microWorkDao.STORE_TYPE_OUTLETS, microWorkDao.STORE_TYPE_RETAIL}
	for _, storeType := range storeTypeArr {
		memberIds, err := s.GetMemberIdsByStoreType(c, enterpriseId, storeType)
		if err != nil {
			return 0, err
		}
		// 判断用户是否是奥莱用户或者零售用户
		if helper.InArray(memberId, memberIds) {
			return storeType, nil
		}
	}
	return 0, nil
}

// GetMemberCountByStoreType  获取奥莱或者零售的员工人数
func (s MemberDepartmentsService) GetMemberCountByStoreType(c app.Ctx, enterpriseId, storeType int64) (int64, error) {
	memberIds, err := s.GetMemberIdsByStoreType(c, enterpriseId, storeType)
	if err != nil {
		return 0, err
	}

	return int64(len(memberIds)), nil
}
