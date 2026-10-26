package enterprise_sso_config

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type EnterpriseSSOConfig struct {
	dao   dao.Dao
	model *micro_scrm.SSOConfigModel
}

func NewEnterpriseSSOConfig() *EnterpriseSSOConfig {
	return &EnterpriseSSOConfig{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(micro_scrm.SSOConfigModel),
	}
}

// GetConfigByEnterpriseIdAndType 根据企业ID和SSO类型获取配置
func (d *EnterpriseSSOConfig) GetConfigByEnterpriseIdAndType(c app.Ctx, enterpriseId int64, ssoType string) (micro_scrm.SSOConfigModel, error) {
	var configModel micro_scrm.SSOConfigModel
	err := d.dao.WithContext(c).Model(d.model).Where("enterprise_id = ? AND sso_type = ?", enterpriseId, ssoType).Find(&configModel).Error
	if err != nil {
		return configModel, err
	}

	return configModel, nil
}
