package activityCampaign

import (
	"errors"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/logic"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
)

const (
	ACTIVITIE_STATUS_ZERO = iota
	ACTIVITIE_STATUS_NEED_START
	ACTIVITIE_STATUS_RUNNING
	ACTIVITIE_STATUS_FINISHED
)

func ActivitiesH5Info(params ActivitiesRequest) micro_scrm_enterprise.OperationActivityModel {
	activitie, _ := orm.Service.ActivityCampaign.ActivitieInfo(params.EnterpriseId, params.ActivityId)
	if activitie.Status == ACTIVITIE_STATUS_RUNNING {
		activitieendTime := helper.GetFormatDate(activitie.EndTime, 0)
		if helper.GetDateMorningTimestamp(activitieendTime) < helper.GetDateMorningTimestamp(helper.GetNowDate()) {
			activitie.Status = 3
		}
	}
	return activitie
}
func H5UsersignInList(params UsersignInRequest) UsersignInResponse {
	result := UsersignInResponse{
		UnionId:     params.UnionId,
		SingInDates: []SingInDate{},
	}
	activitie, err := orm.Service.ActivityCampaign.ActivitieInfo(params.EnterpriseId, params.ActivityId)
	if err != nil || activitie.ID == 0 {
		return result
	}

	if activitie.StartTime == 0 || activitie.EndTime == 0 {
		return result
	}
	activitieStartTime := activitie.StartTime
	// activitieEndtime := activitie.EndTime
	result.ActivitieDate = helper.GetFormatDate(activitie.StartTime, 1) + "-" + helper.GetFormatDate(activitie.EndTime, 1)
	singInDates := []SingInDate{}
	singInDates = append(singInDates, SingInDate{
		Date:       helper.GetFormatDate(activitieStartTime, 0),
		FormatDate: helper.GetFormatDate(activitieStartTime, 3),
	})
	// 查询签到日期
	userSignInDates := orm.Service.ActivityCampaign.GetUserSignInDate(params.EnterpriseId, params.ActivityId, params.UnionId)
	for {
		activitieStartTime += 86400
		if activitieStartTime > activitie.EndTime && helper.GetFormatDate(activitieStartTime, 2) != helper.GetFormatDate(activitie.EndTime, 2) {
			break
		}
		singInDates = append(singInDates, SingInDate{
			Date:       helper.GetFormatDate(activitieStartTime, 0),
			FormatDate: helper.GetFormatDate(activitieStartTime, 3),
		})
	}
	newDate := helper.GetNowDate()
	for _, signdate := range singInDates {
		flag := false
		for _, userSignInDate := range userSignInDates {
			if signdate.Date == userSignInDate {
				flag = true
				break
			}
		}
		if flag {
			signdate.ActivitieStatus = 2
			signdate.Remark = "用户已签到"
		} else {
			qiandao := helper.GetDateMorningTimestamp(signdate.Date)
			now := helper.GetDateMorningTimestamp(newDate)
			if signdate.Date == helper.GetFormatDate(activitie.StartTime, 0) {
				qiandao = activitie.StartTime
				now = time.Now().Unix()
				if qiandao <= now {
					qiandao = helper.GetDateMorningTimestamp(signdate.Date)
					now = helper.GetDateMorningTimestamp(newDate)
				}
			}
			if signdate.Date == helper.GetFormatDate(activitie.EndTime, 0) {
				qiandao = activitie.EndTime
				now = time.Now().Unix()
				if qiandao > now {
					qiandao = helper.GetDateMorningTimestamp(signdate.Date)
					now = helper.GetDateMorningTimestamp(newDate)
				}
			}
			if qiandao < now {
				signdate.ActivitieStatus = 4
				signdate.Remark = "日期已过,用户未签到"
			} else if qiandao == now {
				signdate.ActivitieStatus = 1
				signdate.Remark = "可签到"
			} else {
				signdate.ActivitieStatus = 3
				signdate.Remark = "签到时间未到"
			}
		}
		result.SingInDates = append(result.SingInDates, signdate)
	}
	// 查询用户是否报名
	// 查询用户是否报名过依赖的活动
	parentActivitie, err := orm.Service.ActivityCampaign.FindUserSignup(params.EnterpriseId, activitie.ParentId, params.UnionId)
	if err != nil || parentActivitie.ID == 0 {
		result.Message = "很抱歉！您暂无资格参与本次活动，谢谢！"
	}
	return result
}
func UsersignInH5(ctx ctx.Context, params UsersignInRequest) error {
	// 查询用户是否存在如果不存在插入用户信息
	if params.SignInContent == "" || params.MediaId == "" {
		return errors.New("签到信息不齐全")
	}
	activitie, err := orm.Service.ActivityCampaign.ActivitieInfo(params.EnterpriseId, params.ActivityId)
	if err != nil {
		return errors.New("未找到对应的活动信息")
	}
	if activitie.ID == 0 {
		return errors.New("未找到对应的活动信息")
	}
	if activitie.Type != 2 {
		return errors.New("该活动不支持签到")
	}
	if activitie.Status != ACTIVITIE_STATUS_RUNNING {
		return errors.New("活动已结束或未开始")
	}
	// newTime := helper.GetDateMorningTimestamp(helper.GetNowDate())
	// activitieStartTime := helper.GetFormatDate(activitie.StartTime, 0)
	newTime := time.Now().Unix()
	if activitie.StartTime > newTime {
		return errors.New("活动暂未开始")
	}
	// activitieendTime := helper.GetFormatDate(activitie.EndTime, 0)
	if activitie.EndTime < newTime {
		return errors.New("活动期限已过")
	}
	if activitie.ParentId == 0 {
		return errors.New("很抱歉！您暂无资格参与本次活动，谢谢！")
	}
	// 查询用户是否报名过依赖的活动
	parentActivitie, err := orm.Service.ActivityCampaign.FindUserSignup(params.EnterpriseId, activitie.ParentId, params.UnionId)
	if err != nil || parentActivitie.ID == 0 {
		return errors.New("很抱歉！您暂无资格参与本次活动，谢谢！")
	}
	// 查询用户今日是否签到过
	startTime, endTime := helper.PeriodOfTimestamp()
	if startTime == 0 || endTime == 0 {
		return errors.New("服务器错误,系统获取时间错误")
	}
	result, err := orm.Service.ActivityCampaign.Findactivitysigninlog(params.EnterpriseId, params.ActivityId, startTime, endTime, params.UnionId)
	if err != nil {
		return errors.New("用户今日已签到")
	}
	if result.ID > 0 {
		return errors.New("用户今日已签到")
	}
	o, err := logic.NewOfficialAccount(ctx, params.EnterpriseId)
	if err != nil {
		log.Info(ctx).Msgf("活动签到图片上传失败:%v", err)
		return errors.New("很抱歉，网络繁忙中，请稍后再试，谢谢！")
	}
	picUrl, err := o.GetMedia(ctx, params.MediaId)
	if err != nil {
		log.Info(ctx).Msgf("活动签到图片上传失败:%v", err)
		return errors.New("很抱歉，网络繁忙中，请稍后再试，谢谢！")
	}
	// 查询用户是否存在存在就不添加
	activityusers, err := orm.Service.ActivityCampaign.FindActivityUserinfo(params.EnterpriseId, params.ActivityId, params.UnionId)
	activityUserId := activityusers.ID
	if err != nil || activityusers.ID == 0 { // 用户不存需要添加用户
		data := micro_scrm_enterprise.UserModel{

			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(ctx),
				},
			},
			ActivityId: params.ActivityId,
			UnionId:    params.UnionId,
		}
		if parentActivitie.ID != 0 { // 表示参加过该活动直接取信息
			data.WxName = parentActivitie.WxName
			data.MemberCode = parentActivitie.MemberCode
			data.WxPhone = parentActivitie.WxPhone
			data.ReceiverName = parentActivitie.ReceiverName
			data.ReceiverPhone = parentActivitie.ReceiverPhone
			data.ReceiverAddress = parentActivitie.ReceiverAddress
		} else {
			user, err := orm.Service.Activity.GetUserName(params.EnterpriseId, params.UnionId)
			if err == nil && user.ID != 0 {
				data.WxName = user.Name
				data.MemberCode = user.UnionCode
				Eud, err := orm.Service.Activity.GetUserPhone(params.EnterpriseId, user.Euid)
				if err == nil && Eud.ID != 0 {
					data.WxPhone = Eud.Mobile
				}
			}
		}
		if _, err := orm.Service.ActivityCampaign.InsertactivityUser(params.EnterpriseId, data); err == nil {
			// 活动加一人
			activityUserId = data.ID
			orm.Service.ActivityCampaign.ActivitieAddNum(params.EnterpriseId, params.ActivityId)
		}

	}
	// 签到
	resultlog := micro_scrm_enterprise.OperationActivitySignInLog{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: uid.ID(ctx),
			},
		},
		ActivityId:     params.ActivityId,
		UnionId:        params.UnionId,
		ActivityUserId: activityUserId,
		SignInContent:  params.SignInContent,
		SignInPic:      picUrl,
	}

	err = orm.Service.ActivityCampaign.Createactivitysigninlog(params.EnterpriseId, resultlog)
	return err
}
func UsersignInInfoH5(params UsersignInRequest) UsersignInInfoResponse {
	startTime := helper.GetDateMorningTimestamp(params.SignInDate)
	endTime := helper.GetDateNightTimestamp(params.SignInDate)
	result := UsersignInInfoResponse{
		UnionId:         params.UnionId,
		SingInDate:      params.SignInDate,
		FormatDate:      strings.ReplaceAll(params.SignInDate, "-", "/"),
		ActivitieStatus: 0,
		Remark:          "暂无用户签到信息",
	}
	if startTime == 0 || endTime == 0 {
		return result
	}
	activitie, err := orm.Service.ActivityCampaign.ActivitieInfo(params.EnterpriseId, params.ActivityId)
	if err != nil || activitie.ID == 0 {
		return result
	}

	if activitie.StartTime == 0 || activitie.EndTime == 0 {
		return result
	}
	result.ActivitieDate = helper.GetFormatDate(activitie.StartTime, 1) + "-" + helper.GetFormatDate(activitie.EndTime, 1)

	result.SingInInfo, _ = orm.Service.ActivityCampaign.Findactivitysigninlog(params.EnterpriseId, params.ActivityId, startTime, endTime, params.UnionId)
	if result.SingInInfo.ID != 0 {
		result.ActivitieStatus = 2
		result.Remark = "用户已签到"
	} else if result.SingInInfo.ID == 0 {
		newDate := helper.GetNowDate()
		qiandao := helper.GetDateMorningTimestamp(newDate)
		if startTime > qiandao {
			result.ActivitieStatus = 3
			result.Remark = "签到时间未到"
		} else if startTime == qiandao {
			result.ActivitieStatus = 1
			result.Remark = "用户可签到"
		} else {
			result.ActivitieStatus = 4
			result.Remark = "用户未签到"
		}
	}
	return result
}
