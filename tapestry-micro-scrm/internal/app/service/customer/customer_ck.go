package customerSvc

import (
	"micro-scrm/internal/app"
	customerReq "micro-scrm/internal/app/http/request/customer"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

func (s *Customer) GetExternalUserListCk(c app.Ctx, request customerReq.ExternalUserListsReq, enterpriseConfigName string) ([]microScrmModel.ExternalUserList, int64, error) {

	result := make([]microScrmModel.ExternalUserList, 0)
	// 判断当前员工类型 店员或店长
	memberList, userids, err := s.GetMembersByManageType(c, request.EnterpriseID, request.MemberId, request.Userid, request.MemberType, request.ThirdIds, request.MemberList)
	if err != nil {
		return result, 0, err
	}
	request.MemberLists = memberList
	request.Userids = userids

	total, result, err := s.workExternalDao.GetUserListCk(c, request)
	if err != nil {
		return result, 0, err
	}
	if total == 0 {
		return result, total, nil
	}

	// 处理数据
	if enterpriseConfigName == "coach" {
		result, err = s.HandleUserListCoach(c, request, result, enterpriseConfigName)
	}
	if enterpriseConfigName == "sw" {
		result, err = s.HandleUserListSW(c, request, result, enterpriseConfigName)
	}
	if enterpriseConfigName == "ks" {
		result, err = s.HandleUserListKS(c, request, result, enterpriseConfigName)
	}
	if err != nil {
		return result, 0, err
	}

	return result, total, nil
}
