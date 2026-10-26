package activity_campaign

import "gorm.io/gorm"

type ActivityCampaignService struct {
	*gorm.DB
	ServiceName string
}

// NewPlanTaskService 创建服务
func NewActivityCampaignService(g *gorm.DB) *ActivityCampaignService {
	return &ActivityCampaignService{
		DB:          g,
		ServiceName: "activity_campaign_service",
	}
}
