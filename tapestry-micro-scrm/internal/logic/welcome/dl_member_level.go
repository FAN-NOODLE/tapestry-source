package welcome

import (
	"errors"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/welcome"
	"micro-scrm/internal/pkg/dl"
)

type DLMemberLevelParam struct {
	WorkExternalUserID string // 企微外部联系人的user_id
	UnionID            string // 同一主体下的粉丝的union_id
}

// DLMemberLevel 根据DL的会员等级推送欢迎语
func (logic welcomeUserLogic) DLMemberLevel(ctx ctx.Context, param DLMemberLevelParam) (data welcome.DrainageWelcomesModel, err error) {
	ev := log.Info(ctx).Interface("param", param)
	defer func() {
		if err != nil {
			ev.Interface("err", err)
		}
		ev.Msg("DLMemberLevel")
	}()
	// 如果union_id为空，则通过外部联系人的user_id。user_id一定存在
	if param.UnionID == "" {
		userInfo, userInfoErr := orm.Service.WorkExternalUser.GetWorkExternalUserInfoByWorkExternalUserID(logic.enterpriseID, param.WorkExternalUserID)
		if userInfoErr != nil {
			err = userInfoErr
			return
		}

		param.UnionID = userInfo.Unionid
	}

	// 根据union_id, 获取用户DL的卡等信息
	userGrade, userGradeErr := dl.NewDL().MemberQueryGrade(ctx, dl.MemberQueryGradeRequest{
		EnterpriseID: logic.enterpriseID,
		Type:         "unionId",
		Value:        param.UnionID,
	})
	ev.Interface("MemberQueryGradeOut", userGrade)

	if userGradeErr != nil {
		err = userGradeErr
		return
	}

	if userGrade.Error != "" {
		err = errors.New(userGrade.Error)
		return
	}

	if userGrade.Data.Grade == "" && userGrade.Data.GradeCode == "" {
		err = errors.New("查询不到用户等级")
		return
	}
	data, err = logic.getBrandGradeWelcome(ctx, userGrade, logic.enterpriseID)
	return
}

// 获取品牌对应卡等的欢迎互
func (logic welcomeUserLogic) getBrandGradeWelcome(ctx ctx.Context, param dl.MemberQueryGradeResponse, enterpriseId int64) (data welcome.DrainageWelcomesModel, err error) {
	ev := log.Info(ctx).Interface("param", param).Interface("enterpriseId", enterpriseId)
	defer func() {
		if err != nil {
			ev.Interface("err", err)
		}
		ev.Msg("getBrandGradeWelcome")
	}()
	inParam := welcome.WelcomesExternalTypeDataModel{
		Code: param.Data.GradeCode,
		Name: param.Data.Grade,
	}
	inParam.EnterpriseID = enterpriseId
	cond, condErr := orm.Service.Welcome.QueryWelcomeExternalData(inParam)
	if condErr != nil {
		err = condErr
		return
	}
	ev.Interface("cond", cond)
	data, err = orm.Service.Welcome.QueryExternalWelcome(welcome.DrainageWelcomesModel{
		// TODO 三方欢迎语
		//WelcomesExternalTypeDataID: uint64(cond.ID),
	})
	ev.Interface("data", data)
	return
}
