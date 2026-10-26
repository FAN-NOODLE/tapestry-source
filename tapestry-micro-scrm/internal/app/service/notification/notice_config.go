package notificationSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao/notification"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	"micro-scrm/internal/app/service"
	"time"
)

type NoticeConfigSvc struct {
	service.Base
	noticeConfigDao *notificationDao.NoticeConfigDao
}

func NewNoticeConfigSvc() *NoticeConfigSvc {
	return &NoticeConfigSvc{
		noticeConfigDao: notificationDao.NewNoticeConfigDao(),
	}
}

func (s *NoticeConfigSvc) Save(c ctx.Context, params notificationReq.NoticeConfigSaveReq) error {
	for _, config := range params.NoticeConfigList {
		err := s.noticeConfigDao.UpdateData(c, params.EnterpriseId, map[string]interface{}{
			"info_type": config.InfoType,
		}, map[string]interface{}{
			"deadline_week": config.DeadlineWeek,
			"deadline_time": config.DeadlineTime,
			"release_week":  config.ReleaseWeek,
			"release_time":  config.ReleaseTime,
			"updated_at":    time.Now().Unix(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NoticeConfigSvc) Detail(c ctx.Context, enterpriseId int64) (list []notificationReq.NoticeConfigReq, err error) {
	configList, err := s.noticeConfigDao.GetAll(c, enterpriseId)
	if err != nil {
		return list, err
	}
	for _, config := range configList {
		list = append(list, notificationReq.NoticeConfigReq{
			InfoType:     config.InfoType,
			DeadlineWeek: config.DeadlineWeek,
			DeadlineTime: config.DeadlineTime,
			ReleaseWeek:  config.ReleaseWeek,
			ReleaseTime:  config.ReleaseTime,
		})
	}
	return
}
