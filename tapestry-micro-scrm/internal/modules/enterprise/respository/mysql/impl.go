package mysql

import (
	"database/sql"
	"encoding/json"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/enterprise"
)

// enterpriseRepo
type enterpriseRepo struct {
	*repos.SQLRepo
}

// FindEnterpriseDetail 查询企业信息
func (e *enterpriseRepo) FindEnterpriseDetail(ctx ctx.Context, c enterprise.Corp) (*enterprise.Corp, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = c
	stmt, err := e.DB.Template("t_enterprise_corp/detail", &mapping)

	if err != nil {
		return nil, err
	}
	dest := &enterprise.Corp{}
	rErr := stmt.NamedGet(ctx, dest, c)
	if rErr != nil {
		return nil, rErr
	}
	return dest, nil
}

// ListEnterprise 查询企业信息
func (e *enterpriseRepo) ListEnterprise(cc ctx.Context, enterpriseIds []int64, page int, size int) (int64, *[]*enterprise.Corp, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_enterprise_corp/select",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*enterprise.Corp{}
		},
	}
	mapping := make(map[string]interface{})
	type listEnterprise struct {
		Page             int `json:"page"`
		Size             int `json:"size"`
		EnterpriseIds    []int64
		EnterpriseIdsCnt int
	}
	inter := listEnterprise{
		Page:             page,
		Size:             size,
		EnterpriseIds:    enterpriseIds,
		EnterpriseIdsCnt: len(enterpriseIds),
	}
	offset := (page - 1) * size
	pageOption := repos.PageVO{
		Offset: offset,
		Limit:  size,
	}
	mapping["cond"] = inter
	total, list, err := e.FindPage(cc, pageOption, inter, option)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}

	return total, list.(*[]*enterprise.Corp), nil
}

// UpdateEnterprise 更新企业信息
func (e *enterpriseRepo) UpdateEnterprise(cc ctx.Context, c enterprise.Corp) error {
	mapping := make(map[string]interface{})
	mapping["cond"] = c
	stmt, err := e.DB.Template("t_enterprise_corp/update", &mapping)
	if err != nil {
		return nil
	}
	res, sErr := stmt.NamedExec(cc, c)
	if sErr != nil {
		return sErr
	}
	_, resErr := res.RowsAffected()
	return resErr
}

// OfficialInfo 获取企业的公众号配置信息
func (e *enterpriseRepo) OfficialInfo(context ctx.Context, id int64) (*enterprise.OfficialInfo, error) {
	req := &enterprise.ConfigReq{
		EnterpriseID: id,
		Key:          "official_info",
	}
	data, err := e.EnterpriseConfig(context, req)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, nil
	}

	result := new(enterprise.OfficialInfo)
	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// EnterpriseConfig 获取配置信息
func (e *enterpriseRepo) EnterpriseConfig(ctx ctx.Context, req *enterprise.ConfigReq) (string, error) {
	stmt, err := e.DB.Template("t_enterprise_config/get_config")
	if err != nil {
		return "", err
	}
	dest := new(enterprise.ConfigValue)
	rErr := stmt.NamedGet(ctx, dest, req)
	switch rErr {
	case sql.ErrNoRows:
		return "", nil
	default:
		return "", rErr
	case nil:
		return dest.Value, nil
	}
}

func (e *enterpriseRepo) OfficialUser(context ctx.Context, req *enterprise.OfficialUserReq) (*enterprise.OfficialUser, error) {
	stmt, err := e.DB.Template("t_fission_mp_user/find-user", e.Cond2Args(req))
	if err != nil {
		return nil, err
	}
	dest := new(enterprise.OfficialUser)
	rErr := stmt.NamedGet(context, dest, req)
	switch rErr {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, rErr
	case nil:
		return dest, nil
	}
}

func (e *enterpriseRepo) CreateOfficialUser(context ctx.Context, user *enterprise.OfficialUser) error {
	stmt, err := e.DB.Template("t_fission_mp_user/insert", e.Cond2Args(user))
	if err != nil {
		return nil
	}
	res, sErr := stmt.NamedExec(context, user)
	if sErr != nil {
		return sErr
	}
	_, resErr := res.RowsAffected()
	return resErr
}

func (e enterpriseRepo) CreateOfficialInfo(context ctx.Context, data *enterprise.EpConfig) error {
	stmt, err := e.DB.Template("t_enterprise_config/insert", e.Cond2Args(data))
	if err != nil {
		return nil
	}
	res, sErr := stmt.NamedExec(context, data)
	if sErr != nil {
		return sErr
	}
	_, resErr := res.RowsAffected()
	return resErr
}

// QueryAppTypeCount 查询公众号授权状态
func (e *enterpriseRepo) QueryAppTypeCount(c ctx.Context, enterpriseID int64) (map[int]int, error) {
	var cnts []struct {
		AppType int `db:"app_type"`
		Count   int `db:"cnt"`
	}

	stmt, err := e.DB.Template("t_wx_app_config/select", &map[string]interface{}{"enterprise_id": enterpriseID})
	if err != nil {
		return nil, err
	}

	if err = stmt.Select(c, &cnts); err != nil {
		return nil, err
	}

	cntmap := make(map[int]int)
	for _, cnt := range cnts {
		cntmap[cnt.AppType] = cnt.Count
	}
	return cntmap, nil
}

// AddWxAppConfig 添加公众号账号信息
func (e *enterpriseRepo) AddWxAppConfig(c ctx.Context, data *enterprise.WxAppConfigEntity) (int64, error) {
	stmt, err := e.DB.Template("t_wx_app_config/insert", &map[string]interface{}{"enterprise_id": data.EnterpriseID})
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (e *enterpriseRepo) QueryOneApp(c ctx.Context, enterpriseID int64, apptype int) (string, error) {
	stmt, err := e.DB.Template("t_wx_app_config/select-one-appid", &map[string]interface{}{"enterprise_id": enterpriseID})
	if err != nil {
		return "", err
	}
	var appid string
	if err := stmt.Get(c, &appid, apptype); err != nil {
		return "", err
	}
	return appid, nil
}
