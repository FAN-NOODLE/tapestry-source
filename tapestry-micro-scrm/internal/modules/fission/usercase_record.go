package fission

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"strings"
	"time"
)

// ActivityDetail 查询活动详情
func (s *Service) ActivityDetail(c ctx.Context, enterpriseID int64, activityID int64) (activity *ActivityEntity, err error) {
	return s.getActivityDetail(c, enterpriseID, activityID, false)
}

// ActivityDetailWithDeleted 查询活动详情，查询已删除数据
func (s *Service) ActivityDetailWithDeleted(c ctx.Context, enterpriseID int64, activityID int64) (activity *ActivityEntity, err error) {
	return s.getActivityDetail(c, enterpriseID, activityID, true)
}

// getActivityDetail 查询活动详情
func (s *Service) getActivityDetail(c ctx.Context, enterpriseID int64, activityID int64, withDeleted bool) (activity *ActivityEntity, err error) {
	if withDeleted {
		activity, err = s.repo.GetActivityDetailWithDeleted(c, activityID, enterpriseID)
	} else {
		activity, err = s.repo.ActivityDetail(c, activityID, enterpriseID)
	}

	if err != nil {
		return activity, err
	}
	if activity == nil {
		return nil, errs.FissionActivityNotExist
	}
	if activity.DeletedAt > 0 {
		activity.Status = StatusDeleted
	} else {
		activity.Status = GetActivityStatus(activity.StartTime, activity.EndTime, activity.Status)
	}
	return activity, nil
}

// InviteDetail 发起人邀请明细列表
func (s *Service) InviteDetail(c ctx.Context, req InviteDetailRequest) (total int64, list []*AssistantRecordInfo, err error) {
	return s.repo.InviteDetail(c, req)
}

// ExchangeRecord 发起人兑换记录
func (s *Service) ExchangeRecord(c ctx.Context, req ExchangeDetailRequest) (int64, *[]ExchangeRecordInfo, error) {

	return s.repo.ExchangeRecord(c, req)
}

// ActivityDetailList 活动详情 列表
func (s *Service) ActivityDetailList(c ctx.Context, req ActivityDetailRequest) (int64, *[]ActivityUserEntity, error) {
	return s.repo.ActivityDetailList(c, req)
}

// OpenExchangeWarn 兑奖提醒配置
func (s *Service) OpenExchangeWarn(c ctx.Context, req ExchangeWarnRequest) error {
	var (
		err error
	)
	err = s.repo.StoreExchangeWarnConfig(c, req)
	if nil != err {
		return err
	}
	return nil
}

// ExchangeTotal 兑换记录总数
func (s *Service) ExchangeTotal(c ctx.Context, req ExchangeDetailRequest) (int64, error) {
	var (
		total int64
		err   error
	)
	total, err = s.repo.ExchangeTotal(c, req.EnterpriseID)
	if nil != err {
		return 0, err
	}
	return total, nil
}

// SendMessage发送兑奖提醒消息
func (s *Service) SendMessage(c ctx.Context, enterpriseID int64, activityID int64, activityName string, nickName string) error {
	var (
		configStr  string
		err        error
		config     ExchangeWarnRequest
		members    = make([]*Member, 0)
		ids        = make([]string, 0)
		pluginInfo *boss.WorkPluginResponse
		toUser     string
		idsArr     [][]string
		domain     string
		length     int
	)
	configStr, err = s.repo.GetExchangeWarnConfig(c, enterpriseID, ExchangeWarn)
	if nil != errorIgnore.Ignore(err) {
		return err
	}
	if "" != configStr {
		err = json.Unmarshal([]byte(configStr), &config)
		if nil != err {
			return err
		}
	}
	// 提醒是开启状态
	if config.Status == WarnOn {
		// 查询所有活动下的裂变员工去重
		members, err = s.repo.FissionMemberList(c, enterpriseID, activityID)
		if nil != err {
			log.Info(c).Msgf("获取活动下的裂变员工出错:%v", err.Error())
			return err
		}
		// 合并列表员工和额外员工
		for _, member := range members {
			ids = append(ids, member.UserID)
		}
		length = len(config.MemberIds)

		if 0 != length {

			members, err = s.repo.MemberList(c, config.MemberIds, enterpriseID)
			if nil != err {
				log.Info(c).Msgf("查询员工信息出错:%v", err.Error())
				return err
			}
			for _, member := range members {
				ids = append(ids, member.UserID)
			}
		}
		domain, _, _, err = s.Enterprise.GetInfo(c, enterpriseID)
		if nil != err {
			log.Info(c).Msgf("获取企业信息错误:%v", err.Error())
			return err
		}
		req := boss.WorkPlugin{
			EnterpriseID: enterpriseID,
		}
		// 获取自建应用信息
		pluginInfo, err = boss.WorkPluginInfo(c, req)
		if err != nil {
			return err
		}

		if pluginInfo.Data.AgentID == 0 {
			return err
		}
		workImpl := new(work.QyWorkImpl).SetDomain(domain)

		// messContent := work.MessageTypeTextCardContent{
		//	Title:       activityName + WarnTitle,
		//	Description: nickName + WarnText,
		//	URL:         customer.CrsDomain,
		// }
		messContent := work.TextData{
			Content: nickName + WarnText,
		}
		// 因为toUser最大限制1000,故分批处理
		for {
			length = len(ids)
			if length > 990 {
				idsArr = append(idsArr, ids[0:990])
				ids = ids[990:length]
			} else {
				idsArr = append(idsArr, ids)
				break
			}
		}
		for _, ids = range idsArr {
			toUser = strings.Join(ids, "|")
			sendReq := work.SendMessageRequest{
				Touser:                 toUser,
				Msgtype:                "text",
				Text:                   &messContent,
				Agentid:                pluginInfo.Data.AgentID,
				DuplicateCheckInterval: 1800,
				QyGwRequest: work.QyGwRequest{
					ReqPL: pluginInfo.Data.Identifier,
				},
			}
			_, err = workImpl.SendMessage(c, sendReq)
			if nil != err {
				// 2s后重试一次
				<-time.After(2 * time.Second)
				_, err = workImpl.SendMessage(c, sendReq)
				if nil != err {
					log.Info(c).Msgf("企业%d下的发送奖品兑换成功提示信息员出错, 错误信息：%+v", enterpriseID, err)
					return err
				}
			}
		}

		return nil
	}
	return nil
}

// ActivityListSimple 活动列表- 全部有效活动
func (s *Service) ActivityListSimple(c ctx.Context, request ActivityListRequest) ([]*ActivityEntity, error) {

	var (
		list = make([]*ActivityEntity, 0)
		err  error
	)
	list, err = s.repo.ActivityListSimple(c, request.EnterpriseID)
	if nil != err {
		return list, err
	}
	return list, nil
}

// getExchangeConfigUser 获取兑换提醒的额外员工
func (s *Service) GetExchangeConfigMember(c ctx.Context, enterpriseID int64) (*ExchangeWarnReponse, error) {

	var (
		members   = make([]*Member, 0)
		configStr string
		err       error
		config    ExchangeWarnRequest
		response  = &ExchangeWarnReponse{}
		length    int
	)

	configStr, err = s.repo.GetExchangeWarnConfig(c, enterpriseID, ExchangeWarn)
	if nil != errorIgnore.Ignore(err) {
		return response, err
	}
	if "" != configStr {
		err = json.Unmarshal([]byte(configStr), &config)
		if nil != err {
			return response, err
		}
	}

	length = len(config.MemberIds)
	if length > 0 {
		members, err = s.repo.MemberList(c, config.MemberIds, enterpriseID)
		if nil != errorIgnore.Ignore(err) {
			return response, err
		}
	}
	response.EnterpriseID = enterpriseID
	response.Status = config.Status
	response.Members = members
	return response, nil
}
