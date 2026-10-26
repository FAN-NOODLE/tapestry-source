package remindconfig

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/pkg/errorIgnore"
	"strings"
)

// Service 应用服务
type Service struct {
	repo Repository
}

// ConfigInfo 获取提醒设置
func (s Service) ConfigInfo(c ctx.Context, params *ConfigInfoRequest, memberService *members.Service) (*ConfigInfoResponse, error) {
	resDetail, err := s.repo.GetRemindConfig(c, *params)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	if resDetail != nil {
		//查询userid的对应的名称与ID
		var userIds []interface{}
		userStrIds := strings.Split(resDetail.StatTxt, ",")
		for _, v := range userStrIds {
			userIds = append(userIds, v)
		}
		var request = members.BatchGetMemberByMemberId{
			EnterpriseID: params.EnterpriseID,
			DeletedAt:    0,
			IDs:          userIds,
		}
		memberLists, err := memberService.GetMemberByMemberIds(c, &request)
		if err != nil {
			return nil, err
		}

		for _, member := range *memberLists {
			value := StatTxt{
				Id:     member.ID,
				Name:   member.Name,
				UserId: member.Userid,
			}
			resDetail.StatList = append(resDetail.StatList, value)
		}
	}
	return resDetail, nil
}

// ConfigAdd 添加超时提醒设置
func (s Service) ConfigAdd(c ctx.Context, data *ConfigSaveRequest) (int64, error) {
	_, err := s.repo.Add(c, data)
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

// ConfigAdd 编辑超时提醒设置
func (s Service) ConfigSave(c ctx.Context, data *ConfigSaveRequest) (int64, error) {
	_, err := s.repo.Update(c, data)
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}
