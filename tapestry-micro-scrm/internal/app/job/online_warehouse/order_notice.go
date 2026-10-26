package online_warehouse

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	onlineWarehouseDao "micro-scrm/internal/app/dao/online_warehouse"
	emailReq "micro-scrm/internal/app/http/request/email"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/service/email"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type OrderNotice struct {
	name      string
	startTime int64
	endTime   int64
}

func NewOrderNotice() *OrderNotice {
	return &OrderNotice{
		name: "网仓订单的通知",
	}
}

func (job *OrderNotice) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err := job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *OrderNotice) Handle(c ctx.Context, enterpriseId int64) error {
	nowTimeHour := time.Now().Hour()
	job.endTime = time.Now().Unix()
	todayZeroTime, _ := helper.PeriodOfTimestamp()
	if nowTimeHour < 12 { // 早上8点
		job.startTime = todayZeroTime - (10 * 3600)
	} else { // 下午16点
		job.startTime = todayZeroTime + (10 * 3600)
	}

	count, err := onlineWarehouseDao.NewOrderDao().GetPendingNums(c, enterpriseId, job.startTime, job.endTime)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	mailList, err := onlineWarehouseDao.NewOrderDao().GetSendMail(c, enterpriseId, 1)
	if err != nil {
		return err
	}
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(c, map[string]interface{}{"enterprise_id": enterpriseId, "key": "web_h5_domain"})
	if err != nil {
		return err
	}
	if res.Name == "" {
		return err
	}
	domain := res.Name
	host := fmt.Sprintf("%s/login", domain)
	var sendMail []string
	for _, mail := range mailList {
		sendMail = append(sendMail, mail.Email)
	}
	err = email.NewSendEmail().SendEmailByDiy(emailReq.SendMailReq{
		Title:     "网仓订单提醒",
		Content:   fmt.Sprintf("有新的网仓订单提交，请登录后台查看处理 <a href=\"%s\">Weclient(coach.com.cn)</a>", host),
		ToAddress: sendMail,
		FileList:  nil,
	}, config.VchangyiEmailConfig{
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		Username: "storeoperation@vchangyi.com",
		Password: "REDACTED_FOR_SOURCE_SHARING",
	})
	if err != nil {
		fmt.Printf("发送失败:%v\n", err)
	}
	return nil
}
