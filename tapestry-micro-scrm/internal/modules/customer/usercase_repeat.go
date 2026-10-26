package customer

import (
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/errorIgnore"
)

//计算重复客户
func (s Service) RepeatCustomerCount(c ctx.Context, enterpriserId int64) {
	log.Info(c).Msgf("RepeatCustomerCountBegin %d", enterpriserId)
	err := s.repo.RepeatCustomerCountJob(c, enterpriserId)
	if err != nil {
		log.Error(c).Msgf("RepeatCustomerCount err %+v", err)
	}
}

// RepeatCustomerList 重复客户列表
func (s Service) RepeatCustomerList(c ctx.Context, params *RepeatListParams) (int64, *RepeatListResponse, error) {
	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_customer/repeat-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &RepeatListResponse{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*RepeatListResponse)

	return total, resList, nil
}

// FindExternalRelationByUserId 根据外部联系人userid 查询外部联系人
func (s Service) FindExternalRelationByUserId(c ctx.Context, params *ExternalRelationByUserIdRequest) (int64, *WorkUserMemberRelationMemberIds, error) {
	total, resList, err := s.repo.FindExternalBindMembers(c, params)
	if err != nil {
		return 0, nil, err
	}
	if total > 0 {
		// 查询客户绑定的导购
		guideDetail, _ := GuideService.FindGuideByUserIDForCustomer(c, params.EnterpriseID, params.UserID)
		for _, li := range *resList {
			// 成员组织查询
			departmentParams := &DepartmentListByMemberRequest{li.MemberID}
			departments, _ := s.DepartmentListByMember(c, departmentParams)
			if departments != nil {
				var membersMap []string
				for _, department := range *departments {
					membersMap = append(membersMap, department.DpName)
				}
				if len(membersMap) > 0 {
					li.DepartmentName = strings.Join(membersMap, ";")
				}
			}
			// 成员是否为此客户的绑定导购
			if guideDetail != nil && li.MemberID == guideDetail.BindMemberID {
				li.IsBindGuide = 1
			}
		}
	}
	return total, resList, nil
}

// DepartmentListByMember 根据用户查部门
func (s Service) DepartmentListByMember(c ctx.Context, params *DepartmentListByMemberRequest) (*[]Department, error) {
	option := &repos.ListOption{
		TplPrefix: "t_departments/select-dept-by-memberid",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]Department{}
		},
	}
	list, err := s.repo.List(c, params, option)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}
	resList := list.(*[]Department)
	return resList, nil
}
