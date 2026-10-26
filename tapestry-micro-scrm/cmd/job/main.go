package main

import (
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/job/clean"
	dashboardJob "micro-scrm/internal/app/job/dashboard"
	"micro-scrm/internal/app/job/display"
	drainageContactJob "micro-scrm/internal/app/job/drainage_contacts"
	"micro-scrm/internal/app/job/excitation"
	externalSearchJob "micro-scrm/internal/app/job/external_search"
	externalUserJob "micro-scrm/internal/app/job/external_user"
	inventoryJob "micro-scrm/internal/app/job/inventory"
	"micro-scrm/internal/app/job/new_cdp"
	"micro-scrm/internal/app/job/notification"
	"micro-scrm/internal/app/job/online_warehouse"
	orderAchievementJob "micro-scrm/internal/app/job/order_achievement"
	reserveJob "micro-scrm/internal/app/job/reserve"
	shenCeJob "micro-scrm/internal/app/job/shenCe"
	sourceJob "micro-scrm/internal/app/job/source"
	storeOperationJob "micro-scrm/internal/app/job/store_operation"
	thirdCrmSftpJob "micro-scrm/internal/app/job/third_crm_sftp"
	workExternalJob "micro-scrm/internal/app/job/work_external"
	drainagejob "micro-scrm/internal/modules/drainage/delivery/job"
	overtimeJob "micro-scrm/internal/modules/overtime/delivery/job"
	smsJob "micro-scrm/internal/modules/sms/delivery/job"
	"os"
	"runtime"
	"strings"

	customer "micro-scrm/internal/modules/customer/delivery/job"

	planTaskJob_V2 "micro-scrm/internal/modules/plan_task/delivery/job/v2"

	planTaskJob "micro-scrm/internal/modules/plan_task/delivery/job"
	"micro-scrm/internal/modules/tag/delivery/job"

	"micro-scrm/cmd"
	"micro-scrm/internal/config"
	chatJob "micro-scrm/internal/modules/chat/delivery/job"
	commandJob "micro-scrm/internal/modules/command/delivery/job"
	groupmsgJob "micro-scrm/internal/modules/groupmsg/delivery/job"
	marketingJob "micro-scrm/internal/modules/marketing/delivery/job"
	membersJob "micro-scrm/internal/modules/members/delivery/job"
	socialDraw "micro-scrm/internal/modules/socialDraw/delivery/job"
	transferCustomerJob "micro-scrm/internal/modules/transfer_customer/delivery/job"
	userJob "micro-scrm/internal/modules/user/delivery/job"
	"micro-scrm/internal/pkg/redis"

	"micro-scrm/internal/pkg/xjob"
	"micro-scrm/internal/pkg/xjob/store"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	planTasksJob "micro-scrm/internal/app/job/plan_tasks"
	tagJob "micro-scrm/internal/app/job/tag"

	customerJob "micro-scrm/internal/app/job/customer"

	appMembersJob "micro-scrm/internal/app/job/members"
)

// Version ...
var Version = "dev"

// Build ...
var Build = "now"

func main() {
	// 解析命令行参数
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	// 注册配置
	cmd.RegisterConfig(false)
	ctxContext := ctx.New()
	log.Info(ctxContext).Msgf("job.cmd:%+v, job.spec:%+v, job.delCmdSpec:%+v", config.CfgData.Job.SubCmd, config.CfgData.Job.Spec, config.CfgData.Job.DelCmdSpec)

	// 限制cpu,防止过大造成cpu负担
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	}

	s := store.NewSpecRedisStore(redis.Client(), lock.NewRedisLockerFactory(redis.Client()))

	var runner xjob.Runner
	if config.CfgData.Job.Spec != "" {
		specs := strings.Split(config.CfgData.Job.Spec, ",")
		runner = xjob.NewSimpleRunner(config.CfgData.App, Version, s, xjob.Specs(specs...))
	} else {
		runner = xjob.NewSimpleRunner(config.CfgData.App, Version, s)
	}
	registerStoreOperationJobs(runner)

	// 测试
	runner.AddFunc("test", appJob.Test)
	// 会话存档标签补打
	runner.AddFunc("repairConversationArchiveTag", appJob.NewConversationArchiveTagRepair().Run)

	// 拉取会话存档记录
	runner.AddFunc("getMsg", chatJob.GetMsgJob)
	// 获取会话存档媒体文件
	runner.AddFunc("getMsgMedia", chatJob.GetMsgMediaJob)
	// 会话存档超时检测
	runner.AddFunc("overtime", overtimeJob.OvertimeJob)

	// 同步短信签名短信模板状态
	runner.AddFunc("syncTemplateStatus", smsJob.SynTemplateStatusJob)

	// 执行策略可执行记录
	runner.AddFunc("marketingCron", marketingJob.MarketingCronJob)
	// 营销策略--获取群发发送结果
	runner.AddFunc("getGroupMsgResult", marketingJob.GetGroupMsgResultJob)
	// 定时开始营销策略
	runner.AddFunc("autoStartMarketing", marketingJob.AutoStartMarketing)
	// 定时结束营销策略
	runner.AddFunc("autoEndMarketing", marketingJob.AutoEndMarketing)

	// 组织架构成员定时同步
	runner.AddFunc("memberSync", membersJob.MemberJob)

	// 企微-群发消息定时任务
	runner.AddFunc("sendGroupMsgCronJob", groupmsgJob.SendGroupMsgCronJob)
	// 企微-群发消息定时任务发送结果查询
	runner.AddFunc("sendResultCronJob", groupmsgJob.SendResultCronJob)

	// 企微-拉取外部联系人数据到系统
	runner.AddFunc("FixExternalContact", userJob.FixExternalContact)
	// 修复标签数据
	runner.AddFunc("FixTag", job.FixTag)
	// 同步企微标签数据(执行参数配置同步的企业)
	runner.AddFunc("syncQwTag", tagJob.NewSyncQwTag().Run)
	// 同步外部联系人的企微标签
	runner.AddFunc("syncQwExternalTag", job.NewSyncQwExternalTag().Run)

	// 检测投放中的策略模版是否正常
	runner.AddFunc("checkGoOnStrategyTemplate", marketingJob.CheckGoOnStrategyTemplateJob)

	// 计划的定时开启
	runner.AddFunc("SyncPlanStart", planTaskJob.TurnOnPlan)
	// 计划任务结束
	runner.AddFunc("SyncPlanClose", planTaskJob.ClosurePlan)
	// 子任务详情中累计触达客户数 每天凌晨3点
	runner.AddFunc("StatisticsOrderData", planTaskJob.StatisticsOrderData)

	// v2 计划任务 start
	planTask := planTaskJob_V2.NewPlanTaskJob()
	// runner.AddFunc("PlanTask-V2-GetSFTPData", planTask.PullSFTPData.Run)         // 获取数据
	// runner.AddFunc("PlanTask-V2-MakeData", planTask.MakeData.Run)                // 预处理数据
	// runner.AddFunc("PlanTask-V2-ExternalUser-Tagging", planTask.MakeTag.Run)     // 外部联系人打标签
	runner.AddFunc("PlanTask-V2-Push", planTask.PushTask.Run)                    // 客户画像推送计划任务
	runner.AddFunc("PlanTask-V2-Community-Push", planTask.PushCommunityTask.Run) // 社群任务推送
	// runner.AddFunc("PlanTask-V2-SyncTag", planTask.SyncTag.Run)                  // 推送计划任务
	// runner.AddFunc("PlanTask-V2-BirthdayPush", planTasksJob.NewBirthdayTaskPushData().Run) // 生日类型任务推送(已废弃)
	runner.AddFunc("PlanTask-V2-ConsumerPush", planTasksJob.NewTemplateTaskData().Run) // 人群包类型的任务推送

	// v2 计划任务 end
	// 修复外部联系人信息
	runner.AddFunc("CommandFixExternalUser", commandJob.NewFixExternalUserInfoJob().Run)

	// 同步外部人联系信息上次sftp
	runner.AddFunc("CommandSyncExternalcontact", commandJob.NewSyncExternalcontact().Run)
	// 流失客户定时推送到智汇导购小程序
	runner.AddFunc("LostCustomersSendMessage", userJob.LostCustomersSendMessage)

	// 在职继承
	// 生成在职继承数据
	runner.AddFunc("OnTransferCustomerTasksStart", transferCustomerJob.NewOnTransferCustomerTasksStartJob().Run)
	// 开始对客户进行在职继承
	runner.AddFunc("OnTransferCustomerTasksRunning", transferCustomerJob.NewOnTransferCustomerTasksRunningJob().Run)
	// 获取客户在职继承的状态
	// runner.AddFunc("OnTransferResultJob", transferCustomerJob.NewOnTransferResultJob().Run)

	// 离职继承
	// 待继承人员数据入表
	runner.AddFunc("SaveCustomerInheritData", transferCustomerJob.NewSaveCustomerInheritData().Run)
	// 待继承人员数据完善
	// runner.AddFunc("FillCustomerInheritData", transferCustomerJob.FillCustomerInheritData().Run)

	// 每小时计算一次用户的分组，更新用户的标签
	// runner.AddFunc("CalculateUserGroup", userJob.CalculateUserGroup)
	// 同步用户标签到es
	runner.AddFunc("FixHistoryUserTag", userJob.FixHistoryUserTag)
	// 同步用户的备注
	runner.AddFunc("FixHistoryUserRemark", userJob.FixHistoryUserRemark)
	// 指定员工修复其外部联系人
	runner.AddFunc("FixExternalContactByUserID", userJob.FixExternalContactByUserID)
	// 同步外部联系人关系
	runner.AddFunc("SyncMemberExternal", appMembersJob.NewSyncExternal().Run)

	// 历史社群数据修复
	runner.AddFunc("FixHistoryGroupChat", membersJob.FixHistoryGroupChat)
	// 根据指定chatid修复群数据
	runner.AddFunc("FixHistoryGroupChatByChatId", membersJob.FixHistoryGroupChatByChatId)
	// 同步社群画像相关数据
	runner.AddFunc("SocialDrawJob", socialDraw.SocialDrawJob)
	// 统计社群画像活动人数统计
	runner.AddFunc("CountUserActivities", socialDraw.CountUserActivities)
	// 定时计算客户重复数,每天凌晨3点执行一次
	runner.AddFunc("RepeatCustomerCount", customer.CountRepeatCustomer)
	// 初始化流失客户
	runner.AddFunc("LossCustomerAddRelationTimeJob", customer.LossCustomerAddRelationTimeJob)
	// 修改欢迎语的内容结构。手动执行一次即可！如果未执行完，容器就退出了，需要再次执行
	runner.AddFunc("FixDrainageWelcomeContentJob", drainagejob.FixDrainageWelcomeContentJob)
	// 重复好友的导购离职了，需要处理好友与离职导购的关系
	runner.AddFunc("SyncRepeatCustomer", workExternalJob.NewRepeatCustomer().Run)
	// 处理三方crm sftp文件 start
	// 获取加密文件并解密后存入本地等待处理
	runner.AddFunc("PullCrmSftpData", thirdCrmSftpJob.NewPullSFTPData().Run)
	// 同步三方标签
	runner.AddFunc("SyncThirdTagS", tagJob.NewThirdTags().Run)

	runner.AddFunc("SyncKeySa", workExternalJob.NewKaySa().Run)
	// 同步外部联系人的三方标签
	runner.AddFunc("SyncExternalTag", workExternalJob.NewSyncExternalTag().Run)

	// 上传外部联系人标签操作记录到sftp
	runner.AddFunc("UploadImpressionTagOperateLog", thirdCrmSftpJob.NewUploadTagOperateLog().Run)

	// 同步CDP-Crm更新会员信息
	runner.AddFunc("SyncExternalUser", externalUserJob.NewSyncExternalUser().Run)

	// 预约到店邀约当天提醒
	runner.AddFunc("SendTodayReserveData", reserveJob.NewSendTodayReserveData().Run)

	// KS/SW推送前一天到店的预约单还未完成的提醒
	runner.AddFunc("SendPreTodayReserveData", reserveJob.NewSendPreTodayReserveData().Run)
	// KS 更新处理到店时间超过7天还是进行中的邀约单修改为已过期
	runner.AddFunc("AuthUpdatePreMonthOrder", reserveJob.NewAuthUpdatePreMonthOrder().Run)
	// 预约到店，给用户发送订阅消息，coach 为到店时间提前半个小时，ks&&sw企业为到店时间提前三个小时
	runner.AddFunc("SendTomorrowReserveData", reserveJob.NewSendTomorrowReserveData().Run)
	// 上传sw客户附加信息到sftp
	runner.AddFunc("UploadCustomerInfoAddLog", thirdCrmSftpJob.NewUploadCustomerInfoAddLog().Run)
	// KS/SW修改前一天还未完成的预约单 废弃
	// runner.AddFunc("PreTodayReserveOrderChange", reserveJob.NewPreTodayReserveOrderChange().Run)

	// ks的crm sftp文件的脚本
	runner.AddFunc("NewPullCrmSftpData", thirdCrmSftpJob.NewPullKsSFTPData().Run)

	runner.AddFunc("CleanTaskLog", clean.NewCleanTaskLog().Run)

	// 播种任务推送
	runner.AddFunc("sowTaskPush", customerJob.NewSowTask().Run)
	// 定时播种任务推送
	runner.AddFunc("timingSowTaskPush", customerJob.NewTimingSowTask().Run)

	// 推送前一天上新的素材给导购
	runner.AddFunc("PushDailyUpdateSourceData", sourceJob.NewPushDailyUpdateSourceData().Run)
	// 同步 COACH KOS 素材链接
	runner.AddJob("syncCoachKOSSourceLink", sourceJob.NewSyncCoachKOS())
	// 更新coach的服务配置有效期
	runner.AddFunc("UpdateReserveServiceStatus", reserveJob.NewUpdateReserveServiceStatus().Run)
	// 推送昨天进行中和邀约完成的邀约单给到安客诚（acxiom）通过sftp
	runner.AddFunc("PushReserveDataToAcxiom", reserveJob.NewPushReserveDataToAcxiom().Run)
	// 拉取安客诚的匹配邀约单修改到店邀约的预约单状态(旧脚本，只有SW支持，coach改为新的对接脚本PullReserveDataAndSalesForAcxiom)
	runner.AddFunc("PullReserveDataForAcxiom", reserveJob.NewPullReserveDataForAcxiom().Run)
	// 匹配安客诚的交易订单到邀约单(新脚本，只有coach对接)
	runner.AddFunc("PullReserveDataAndSalesForAcxiom", reserveJob.NewPullReserveDataAndSalesForAcxiom().Run)
	// 安客诚cdp会员信息的同步落表-sftp(coach,sw)
	runner.AddFunc("SyncCdpInfoForAcxiom", externalUserJob.NewSyncCdpInfo().Run)
	runner.AddFunc("SyncCdpSWInfoForAcxiom", externalUserJob.NewSyncCdpSWInfo().Run)
	// wedo会员信息的同步落表-sftp(ks)
	runner.AddFunc("SyncKSInfoForWedo", externalUserJob.NewSyncWedoKSInfo().Run)
	// 畅移 Successful_Appointment 邀约单给acxiom
	runner.AddFunc("StoreReserveOrderToAcxiom", reserveJob.NewStoreReserveOrderToAcxiom().Run)
	// 畅移 STORE_APPOINTMENTQTY 邀约单给acxiom 每个财月
	runner.AddFunc("StoreAppointMENTQTYToAcxiom", reserveJob.NewStoreAppointMENTQTYToAcxiom().Run)
	// 畅移 Report_Data 数据给acxiom 每个财月
	runner.AddFunc("StoreReportDataToAcxiom", reserveJob.NewReportDataToAcxiom().Run)
	// 发送奖章任务
	runner.AddJob("SendMedal", excitation.NewSendMedal())
	// 修复发送奖章记录
	runner.AddJob("FixSendMedalRecord", excitation.NewFixSendMedalRecord())
	// 生产数据积累并触发积分获取
	runner.AddJob("PointsProduce", excitation.NewPointsProduce())
	// 成长激励-更新导购等级(每天)
	runner.AddJob("UpdateMemberLevel", excitation.NewUpdateMemberLevel())
	// 成长激励-关联能力图谱规则处理
	runner.AddJob("HandleRelationRule", excitation.NewRelationRule())
	// 成长激励-导购金币过期（目前配置是剩余1小时内触发延迟队列，所以必须每小时执行）
	runner.AddJob("CoinsExpired", excitation.NewCoinsExpired())
	// 同步syncCs2k文件脚本
	runner.AddJob("SyncCs2k", excitation.NewSyncCs2k())
	// 同步welearning的文件脚本
	runner.AddJob("SyncWeLearning", excitation.NewSyncWeLearning())
	// 同步Cps的文件脚本
	runner.AddJob("SyncCps", excitation.NewSyncCps())

	// 计划任务，MA任务(根据任务模板创建的任务)// 用户模板类型的推送
	runner.AddJob("PlanTask-V2-MaPush", planTasksJob.NewMaTaskData())
	// 计划任务，MA任务结束之后需要将完成结果推送给三方/
	runner.AddJob("PlanTask-V2-MaEndPush", planTasksJob.NewMaTaskDataEnd())
	// 同步ks的SA订单数据
	runner.AddJob("PullOrderAchievement", orderAchievementJob.NewPullOrderAchievement())
	// 每日统计三方会员信息变化（SW）
	runner.AddJob("ExternalEverydayStatistics", externalUserJob.NewExternalEverydayStatistics())

	// 处理过期标签
	runner.AddJob("processExpiredTag", tagJob.NewProcessExpired())

	// 计划任务企业群发任务的消息触达，每半小时查询验证下触达情况
	runner.AddJob("checkPlanTaskGroupMsgResult", planTasksJob.NewGroupMsgResult())
	// 计划任务朋友圈群发消息触达验证，每10分钟执行一次
	runner.AddJob("checkPlanTaskMomentMsgResult", planTasksJob.NewMomentMsgResult())
	// 计划任务朋友圈互动数据并推送MA
	runner.AddJob("getPlanTaskMomentComments", planTasksJob.NewMomentComments())

	// 更新印象标签关系
	runner.AddJob("updateImpressionTagRelation", tagJob.NewImpressionTagRelation())

	// 神策数据推送
	runner.AddJob("ShenCePushData", shenCeJob.NewPushData())

	// 企业群发任务，顾客的群发消息当日触达上线，次日补推群发消息，每天13点推送
	runner.AddJob("supplyPushGroupMsg", planTasksJob.NewSupplyPushGroupMsg())

	// dashboard
	// 邀约数据同步
	runner.AddJob("ReserveDataSync", dashboardJob.NewReserveData())

	// 朋友圈任务脚本创建是否完成结果查询,针对于人群包的朋友圈任务标签删除 每5分钟执行一次
	runner.AddJob("checkPlanTaskMomentSendResult", planTasksJob.NewMomentSendResult())

	// 顾客列表附加检索条件信息定时任务落表 每3小时执行一次
	runner.AddJob("externalExtraSync", customerJob.NewExternalExtraSync())

	// SW顾客列表附加检索条件信息定时任务落表 每3小时执行一次
	runner.AddJob("externalExtraSwSync", customerJob.NewExternalExtraSwSync())

	// 初始化task_condition_data分表数据
	runner.AddJob("taskConditionDataShard", planTasksJob.NewTaskConditionDataShard())
	// task_condition_data分表增加字段
	runner.AddJob("alterTaskConditionDataShard", planTasksJob.NewAlterTaskConditionDataShard())

	// 神策数据推送
	runner.AddJob("ShenCePushDataExternal", shenCeJob.NewShenCePushDataExternal())

	// 全量导购关系数据同步CDP
	runner.AddJob("syncRelationDataToCdp", appMembersJob.NewSyncRelationDataToCdp())
	// 每天凌晨执行VM的状态更新脚本，因为VM任务的开始时间是每天的0点，结束时间是每天的23:59:59;所以脚本每天0点执行一次
	runner.AddJob("displayUpdateStatus", display.NewTaskUpdateStatus())
	// 每天早上10天提醒 VM任务当天开始的模板消息
	runner.AddJob("displayStartPush", display.NewTaskStartPush())
	// 每天早上10天提醒 VM任务即将结束的模板消息
	runner.AddJob("displayEndPush", display.NewTaskEndPush())
	// 每天早上10点提醒 VM需重新审核
	runner.AddJob("displayReplayPush", display.NewTaskReplayPush())

	// 生成增量数据上传至DEL
	runner.AddJob("syncDataToDEL", new_cdp.NewSyncDataToDEL())

	runner.AddJob("syncDataToDELAll", new_cdp.NewSyncDataToDELALL())
	// 任务完成明细数据给到acxiom
	runner.AddJob("syncTaskInfoToAcxiom", planTasksJob.NewTaskInfoToAcxiom())
	// 匹配acxiom给到的任务明细数据订单
	runner.AddJob("syncTaskOrderForAcxiom", planTasksJob.NewTaskOrderForAcxiom())

	// 生成dashboard增量数据
	runner.AddJob("dashboardData", dashboardJob.NewTaskData())
	// 任务数据清洗
	runner.AddJob("planTaskDataClen", dashboardJob.NewTaskDataClen())
	// 任务的已联系数据整理
	runner.AddJob("planTaskDataContact", dashboardJob.NewTaskContactData())
	// 生成dashboard全量数据
	runner.AddJob("dashboardDataAll", dashboardJob.NewTaskAllData())
	// 同步CORE管理组织架构的文件脚本
	runner.AddJob("syncArea", dashboardJob.NewSyncArea())

	// Datafabric 数据对接
	runner.AddJob("syncDataFabric", dashboardJob.NewSyncDataFabric())

	//传讯定时发布
	runner.AddJob("NotificationTimedRelease", notification.NewTimedRelease())
	// 网仓订单待发货通知（每天早上8点和下午16点通知）
	runner.AddJob("OnlineWareHouseNotice", online_warehouse.NewOrderNotice())
	// 网仓订单取消（每天0点执行前一天待填写的订单）
	runner.AddJob("OnlineWareHouseCancel", online_warehouse.NewOrderCancel())
	// 网仓订单待填写的消息通知（每天晚上8点）
	runner.AddJob("OnlineWareHouseWaitNotice", online_warehouse.NewOrderWaitNotice())

	// 顾客数据物化视图更新
	runner.AddJob("updateMaterializeView", externalSearchJob.NewUpdateMaterializeView())

	//补偿修复extra遗漏的数据
	runner.AddJob("FixExternalExtra", externalSearchJob.NewFixExternalExtra())

	// 将CK宽表中的业务扩展字段分批回填到企业MySQL扩展表
	runner.AddJob("backfillExternalExtensions", externalSearchJob.NewBackfillExternalExtensions())

	// CDP标签每天全量更新
	runner.AddJob("updateCdpTags", externalSearchJob.NewUpdateCdpTags())
	//本地执行CORA文件更新
	runner.AddJob("syncAreaLocal", dashboardJob.NewSyncAreaLocal())

	//网仓老的商品数据补充sku信息
	runner.AddJob("supplementSkuInfo", online_warehouse.NewSupplementSkuInfo())
	runner.AddJob("importInventorySkuUpcMapping", inventoryJob.NewImportSkuUpcMapping())
	runner.AddJob("compareInventorySkuUpcMapping", inventoryJob.NewCompareSkuUpcMapping())

	//sw订单及会员标签数据拉取入库
	runner.AddJob("syncCrmSWPurchase", externalUserJob.NewSyncCrmSWPurchase())

	//更换开放平台更新unionid
	//runner.AddJob("updateAllUnionid", externalUserJob.NewUpdateUnionid())

	//联系我更新共享码里的员工
	runner.AddJob("drainageContactsUpdate", drainageContactJob.NewDrainageContactsUpdate())

	//cos文件替换cdn地址
	runner.AddJob("cos_to_cdn", appJob.NewCosToCdn())
	runner.AddJob("cos_private_copy", appJob.NewCosPrivateCopy())
	runner.AddJob("cos_private_db_url_refresh", appJob.NewCosPrivateDbURLRefresh())
	runner.AddJob("cos_private_delete_old", appJob.NewCosPrivateDelete())
	runner.AddJob("cos_private_db_migrate", appJob.NewCosPrivateDbMigrate())

	if err := runner.RunCmd(config.CfgData.Job.SubCmd); err != nil {
		log.Info(ctxContext).Msgf("%s执行结果:%v", config.CfgData.Job.SubCmd, err)
	}
}

func registerStoreOperationJobs(runner xjob.Runner) {
	runner.AddJob("syncCoachLMSStoreMaster", storeOperationJob.NewSyncCoachLMSStoreMaster())
	runner.AddJob("syncCoachWorkdayEmployeeMaster", storeOperationJob.NewSyncCoachWorkdayEmployeeMaster())
	runner.AddJob("syncCoachLMSDashboardArea", storeOperationJob.NewSyncCoachLMSDashboardArea())
}
