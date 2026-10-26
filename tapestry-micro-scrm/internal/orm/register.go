package orm

import (
	"micro-scrm/internal/orm/service/data_compass"
	"micro-scrm/internal/orm/service/event"
	"micro-scrm/internal/orm/service/work"
	"sync"

	"micro-scrm/internal/orm/service/activity_campaign"

	"gorm.io/gorm"
	"micro-scrm/internal/orm/service/activity"
	"micro-scrm/internal/orm/service/auth"
	"micro-scrm/internal/orm/service/ba_performance"
	enterprise_boss "micro-scrm/internal/orm/service/boss/enterprise"
	"micro-scrm/internal/orm/service/brand_config"
	"micro-scrm/internal/orm/service/departments"
	"micro-scrm/internal/orm/service/enterprise"
	"micro-scrm/internal/orm/service/fission"
	"micro-scrm/internal/orm/service/group_chat"
	"micro-scrm/internal/orm/service/marketing"
	"micro-scrm/internal/orm/service/member"
	"micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/orm/service/source"
	"micro-scrm/internal/orm/service/structure"
	"micro-scrm/internal/orm/service/tag"
	"micro-scrm/internal/orm/service/transfer_customer"
	"micro-scrm/internal/orm/service/user_event"
	"micro-scrm/internal/orm/service/user_group"
	"micro-scrm/internal/orm/service/version"
	"micro-scrm/internal/orm/service/welcome"
	"micro-scrm/internal/orm/service/work_external_user"
	"micro-scrm/internal/pkg/db/gormx"
)

var Service *ormService

var BossService *ormBossService

var WorkService *ormWorkService

// 增加服务时，需要如下几步
// 1、在ormService添加服务
// 2、增加创建服务的ormCreate函数，例如newOrmUserGroup() 创建了UserGroup服务
// 3、在service中增加对应的服务

type ormService struct {
	db               *gorm.DB
	UserGroup        *user_group.UserGroupService
	Tag              *tag.TagService
	Marketing        *marketing.MarketingService
	Member           *member.MemberService
	WorkExternalUser *work_external_user.WorkExternalUserService
	Fission          *fission.FissionService
	Enterprise       *enterprise.EnterpriseService
	UserEvent        *user_event.UserEventService
	Auth             *auth.AuthService
	Version          *version.VersionService
	Structure        *structure.StructureService
	GroupChat        *group_chat.GroupChatService
	PlanTask         *plan_task.PlanTaskService
	TransferCustomer *transfer_customer.TransferCustomerService
	Source           *source.SourceService
	BaPerformance    *ba_performance.BaPerformanceService
	Departments      *departments.DepartmentsService
	Welcome          *welcome.WelcomeService
	BrandConfig      *brand_config.BrandConfigService
	Activity         *activity.ActivityService
	ActivityCampaign *activity_campaign.ActivityCampaignService
	DataCompass      *data_compass.DataCompassService
	Event            *event.EventService

}

type ormBossService struct {
	db         *gorm.DB
	Enterprise *enterprise_boss.EnterpriseBossService
}

type ormWorkService struct {
	db   *gorm.DB
	Work *work.WorkDBService
}

var once sync.Once

func Init() {
	once.Do(func() {
		Service = newOrmService(
			newOrmUserGroup(),
			newOrmTag(),
			NewOrmMarketing(),
			NewOrmMember(),
			NewWorkExternalUser(),
			NewFission(),
			NewEnterprise(),
			NewUserEvent(),
			NewVersion(),
			NewAuth(),
			NewStructure(),
			NewGroupChat(),
			NewPlanTask(),
			NewSource(),
			NewBAPerformance(),
			NewDepartments(),
			NewWelcome(),
			NewBrandConfig(),
			NewTransferCustomer(),
			NewActivity(),
			NewActivityCampaignService(),
			NewDataCompass(),
			NewEvent(),
		)

		BossService = newOrmBossService(
			NewBossEnterprise(),
		)

		WorkService = newOrmWorkService(
			NewWorkDB(),
		)
	})
}

type ormCreate struct {
	f func(*ormService)
}

// newOrmService 构建全局Orm
func newOrmService(serviceList ...ormCreate) *ormService {
	orm := &ormService{
		db: gormx.GetInstance(),
	}

	for _, service := range serviceList {
		service.f(orm)
	}

	return orm
}

// newOrmUserGroup 构建UserGroupService
func newOrmUserGroup() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.UserGroup = user_group.NewUserGroupService(service.db)
		},
	}
}

// newOrmTag 构建TagService
func newOrmTag() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Tag = tag.NewTagService(service.db)
		},
	}
}

// NewOrmMarketing 构建策略服务
func NewOrmMarketing() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Marketing = marketing.NewMarketingService(service.db)
		},
	}
}

// NewOrmMember 构建员工服务
func NewOrmMember() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Member = member.NewMemberService(service.db)
		},
	}
}

// NewWorkExternalUser 外部联系人服务
func NewWorkExternalUser() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.WorkExternalUser = work_external_user.NewWorkExternalUserService(service.db)
		},
	}
}

// NewFission 裂变活动服务
func NewFission() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Fission = fission.NewFissionService(service.db)
		},
	}
}

// NewEnterprise 企业相关服务
func NewEnterprise() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Enterprise = enterprise.NewEnterpriseService(service.db)
		},
	}
}

func NewDepartments() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Departments = departments.NewDepartmentsService(service.db)
		},
	}
}

func NewAuth() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Auth = auth.NewAuthService(service.db)
		},
	}
}

// NewUserEvent 用户事件服务
func NewUserEvent() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.UserEvent = user_event.NewUserEventService(service.db)
		},
	}
}

func NewStructure() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Structure = structure.NewStructureService(service.db)
		},
	}
}

// NewVersion 版本信息相关服务
func NewVersion() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Version = version.NewVersionService(service.db)
		},
	}
}

// NewGroupChat 版本信息相关服务
func NewGroupChat() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.GroupChat = group_chat.NewGroupChatService(service.db)
		},
	}
}

func NewPlanTask() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.PlanTask = plan_task.NewPlanTaskService(service.db)
		},
	}
}

func NewActivity() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Activity = activity.NewActivityService(service.db)
		},
	}
}

func NewEvent() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Event = event.NewEventService(service.db)
		},
	}
}

func NewActivityCampaignService() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.ActivityCampaign = activity_campaign.NewActivityCampaignService(service.db)
		},
	}
}
func NewSource() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Source = source.NewSourceService(service.db)
		},
	}
}

func NewBAPerformance() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.BaPerformance = ba_performance.NewBaPerformanceService(service.db)
		},
	}
}

func NewWelcome() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.Welcome = welcome.NewWelcomeService(service.db)
		},
	}
}

func NewBrandConfig() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.BrandConfig = brand_config.NewBrandConfigService(service.db)
		},
	}
}
func NewTransferCustomer() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.TransferCustomer = transfer_customer.NewTransferCustomerService(service.db)
		},
	}
}

func NewDataCompass() ormCreate {
	return ormCreate{
		func(service *ormService) {
			service.DataCompass = data_compass.NewDataCompassService(service.db)
		},
	}
}

// boss
type ormBossCreate struct {
	f func(*ormBossService)
}

func NewBossEnterprise() ormBossCreate {
	return ormBossCreate{
		func(service *ormBossService) {
			service.Enterprise = enterprise_boss.NewEnterpriseBossService(service.db)
		},
	}
}

func newOrmBossService(serviceList ...ormBossCreate) *ormBossService {
	orm := &ormBossService{
		db: gormx.GetBossInstance(),
	}

	for _, service := range serviceList {
		service.f(orm)
	}

	return orm
}

// work
type ormWorkCreate struct {
	f func(service *ormWorkService)
}

func newOrmWorkService(serviceList ...ormWorkCreate) *ormWorkService {
	orm := &ormWorkService{
		db: gormx.GetWorkInstance(),
	}

	for _, service := range serviceList {
		service.f(orm)
	}

	return orm
}

func NewWorkDB() ormWorkCreate {
	return ormWorkCreate{
		func(service *ormWorkService) {
			service.Work = work.NewWorkDBService(service.db)
		},
	}
}

