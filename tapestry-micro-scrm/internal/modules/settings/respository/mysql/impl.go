package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/settings"
)

type SettingsRepo struct {
	*repos.SQLRepo
}

// FindSensitive 查询敏感词信息
func (s *SettingsRepo) GetSettingsInfo(ctx ctx.Context, request settings.SettingsInfoRequest) (*settings.ChatSettings, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_chat_settings/find-settings", &mapping)
	if err != nil {
		return nil, err
	}
	detail := &settings.ChatSettings{}
	errStmt := stmt.NamedGet(ctx, detail, request)
	if errStmt != nil {
		return nil, errStmt
	}
	return detail, nil
}

// EnterpriseList 所有企业列表
func (repo *SettingsRepo) EnterpriseList(c ctx.Context) ([]settings.ChatSettings, error) {
	option := &repos.ListOption{
		TplPrefix: "t_chat_settings/all-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]settings.ChatSettings{}
		},
	}
	res, err := repo.List(c, struct {
	}{}, option)

	return *res.(*[]settings.ChatSettings), err
}
