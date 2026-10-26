package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/remindconfig"
)

// RemindConfigRepo
type RemindConfigRepo struct {
	*repos.SQLRepo
}

// GetRemindConfig 通过类型获取提醒配置
func (d *RemindConfigRepo) GetRemindConfig(ctx ctx.Context, config remindconfig.ConfigInfoRequest) (*remindconfig.ConfigInfoResponse, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = config
	stmt, err := d.DB.Template("t_chat_remind_config/select-config", &mapping)
	if err != nil {
		return nil, err
	}
	configInfo := &remindconfig.ConfigInfoResponse{}
	errStmt := stmt.NamedGet(ctx, configInfo, config)
	if errStmt != nil {
		return nil, errStmt
	}
	return configInfo, nil
}
