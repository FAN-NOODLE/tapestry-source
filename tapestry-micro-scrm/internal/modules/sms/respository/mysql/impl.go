package mysql

import (
	"database/sql"
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/pkg/errorIgnore"
)

type SmsTemplateRepo struct {
	*repos.SQLRepo
}

// SmsSignFindOne 查询签名数据
func (s *SmsTemplateRepo) SmsSignFindOne(c ctx.Context, request sms.SmsSignInfoRequest) (*sms.SmsSignInfo, error) {
	stmt, stmtErr := s.DB.Template("t_sms_sign/sms-sign-find-one", s.Cond2Args(request))
	if stmtErr != nil {
		return nil, stmtErr
	}
	smsSignInfo := &sms.SmsSignInfo{}
	getErr := stmt.NamedGet(c, smsSignInfo, request)
	if getErr != nil {
		return nil, getErr
	}
	return smsSignInfo, nil
}

// CreateSmsSign 创建短信签名
func (s *SmsTemplateRepo) CreateSmsSign(ctx ctx.Context, t sms.SmsSignInfo) (int64, error) {
	stmt, err := s.DB.Template("t_sms_sign/insert")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return 0, errStmt
	}

	return res.LastInsertId()
}

// SmsTemplateFindOne 查询短信模板
func (s *SmsTemplateRepo) SmsTemplateFindOne(c ctx.Context, request sms.SmsTemplate) (*sms.SmsTemplate, error) {
	stmt, stmtErr := s.DB.Template("t_sms_template/sms-template-find-one", s.Cond2Args(request))
	if stmtErr != nil {
		return nil, stmtErr
	}
	smsTemplateInfo := &sms.SmsTemplate{}
	getErr := stmt.NamedGet(c, smsTemplateInfo, request)
	if getErr != nil && getErr != sql.ErrNoRows {
		return nil, getErr
	}
	return smsTemplateInfo, nil
}

// CreataSmsTemplate 短信模板创建
func (s *SmsTemplateRepo) CreateSmsTemplate(ctx ctx.Context, t sms.SmsTemplate) (int64, error) {
	stmt, err := s.DB.Template("t_sms_template/insert")

	if err != nil {
		return 0, nil
	}
	_, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return 0, errStmt
	}

	return t.Id, nil
}

// UpdateTemplateStatus 更新短信模板
func (s *SmsTemplateRepo) UpdateTemplateStatus(ctx ctx.Context, t sms.UpdateSmsTemplateStatus) error {
	stmt, err := s.DB.Template("t_sms_template/update-sms-template-status")

	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return errStmt
	}
	return nil
}

// SelectSmsTemplateParamsList 查询短信参数
func (s *SmsTemplateRepo) SelectSmsTemplateParamsList(c ctx.Context) (*[]sms.SmsParams, error) {
	stmt, err := s.DB.Template("t_sms_params/select-list")
	if err != nil {
		return nil, err
	}
	var smsParamsListRequest sms.SmsParams
	paramsList := &[]sms.SmsParams{}
	res, errStmt := stmt.NamedQuery(c, smsParamsListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, paramsList)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return paramsList, nil
}

// SelectSmsTemplateSignList 签名列表接口
func (s *SmsTemplateRepo) SelectSmsTemplateSignList(c ctx.Context) (*[]sms.SmsSignInfo, error) {
	stmt, err := s.DB.Template("t_sms_sign/select-list")
	if err != nil {
		return nil, err
	}
	var smsSignListRequest sms.SmsSignInfo
	ListRequest := &[]sms.SmsSignInfo{}
	res, errStmt := stmt.NamedQuery(c, smsSignListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

// SelectSmsTemplateList 模板列表接口
func (s *SmsTemplateRepo) SelectSmsTemplateList(c ctx.Context, request sms.TemplateListRequest) (*[]sms.SmsTemplate, error) {
	stmt, err := s.DB.Template("t_sms_template/select-list", s.Cond2Args(request))
	if err != nil {
		return nil, err
	}
	var smsTemplateListRequest sms.SmsSignInfo
	ListRequest := &[]sms.SmsTemplate{}
	res, errStmt := stmt.NamedQuery(c, smsTemplateListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

// SmsTemplateDelete 删除短信模板
func (s *SmsTemplateRepo) SmsTemplateDelete(ctx ctx.Context, t sms.UpdateSmsTemplateStatus) error {
	stmt, err := s.DB.Template("t_sms_template/del-sms-template")

	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return errStmt
	}
	return nil
}

// SmsTemplateDelete 编辑短信模板
func (s *SmsTemplateRepo) SmsTemplateEdit(ctx ctx.Context, request sms.UpdateSmsTemplateStatus) error {
	stmt, err := s.DB.Template("t_sms_template/edit-sms-template", s.Cond2Args(request))

	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(ctx, request)

	if errStmt != nil {
		return errStmt
	}
	return nil
}

// UpdateSmsSign 编辑短信签名
func (s *SmsTemplateRepo) UpdateSmsSign(ctx ctx.Context, request sms.SmsSignInfoRequest) error {
	stmt, err := s.DB.Template("t_sms_sign/sms-template-sign-update", s.Cond2Args(request))

	if err != nil {
		return nil
	}
	_, errStmt := stmt.NamedExec(ctx, request)

	if errStmt != nil {
		return errStmt
	}
	return nil
}

// SelectSmsTemplateSignList 签名列表接口
func (s *SmsTemplateRepo) SelectTemplateList(c ctx.Context, EnterpriseId int64) (*[]sms.SmsTemplate, error) {

	smsTemplateListRequest := sms.SmsTemplate{
		EnterpriseId: EnterpriseId,
	}
	stmt, err := s.DB.Template("t_sms_template/template-list", s.Cond2Args(smsTemplateListRequest))
	if err != nil {
		return nil, err
	}

	ListRequest := &[]sms.SmsTemplate{}
	res, errStmt := stmt.NamedQuery(c, smsTemplateListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

// CreateSmsLog 添加短信日志
func (s *SmsTemplateRepo) CreateSmsLog(ctx ctx.Context, t sms.SendLog) (int64, error) {
	stmt, err := s.DB.Template("t_sms_log/insert")

	if err != nil {
		return 0, nil
	}
	res, errStmt := stmt.NamedExec(ctx, t)

	if errStmt != nil {
		return 0, errStmt
	}

	return res.LastInsertId()
}

// CreateSmsLog 添加短信日志
func (s *SmsTemplateRepo) SmsCount(c ctx.Context, request sms.Request) (*sms.CountResponse, error) {
	stmt, stmtErr := s.DB.Template("t_sms_log/select-count")
	if stmtErr != nil {
		return nil, stmtErr
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(c, countResponseInfo, request)
	if getErr != nil {
		return nil, getErr
	}
	return countResponseInfo, nil
}

// SmsTemplateFindOne 查询短信模板
func (s *SmsTemplateRepo) CheckDefaultSmsTemplate(c ctx.Context, request sms.SmsTemplate) (*sms.SmsTemplate, error) {
	stmt, stmtErr := s.DB.Template("t_sms_template/check-default-sms-template")
	if stmtErr != nil {
		return nil, stmtErr
	}
	smsTemplateInfo := &sms.SmsTemplate{}
	getErr := stmt.NamedGet(c, smsTemplateInfo, request)
	if getErr != nil {
		return nil, getErr
	}
	return smsTemplateInfo, nil
}

// SmsTemplateSignList 短信签名列表
func (s *SmsTemplateRepo) SmsTemplateSignList(cc ctx.Context, templateSignRequest sms.TemplateListRequest) (int64, *[]sms.SmsSignInfo, error) {
	offset := (templateSignRequest.Page - 1) * templateSignRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(templateSignRequest.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_sms_sign/search-select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]sms.SmsSignInfo{}
		},
	}
	total, list, err := s.FindPage(cc, pageVo, templateSignRequest, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	resList := list.(*[]sms.SmsSignInfo)

	return total, resList, nil
}

// SmsTemplateList 短信模板列表
func (s *SmsTemplateRepo) SmsTemplateList(c ctx.Context, TemplateListRequest sms.TemplateListRequest) (int64, *[]sms.SmsTemplate, error) {
	offset := (TemplateListRequest.Page - 1) * TemplateListRequest.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(TemplateListRequest.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_sms_template/select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]sms.SmsTemplate{}
		},
	}
	total, list, err := s.FindPage(c, pageVo, TemplateListRequest, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	resList := list.(*[]sms.SmsTemplate)
	return total, resList, nil
}

// GetEnterpriseSmsTemplateCount 获取企业短信模板总条数
func (s *SmsTemplateRepo) GetEnterpriseSmsTemplateCount(c ctx.Context, request sms.SmsTemplate) int64 {
	stmt, err := s.DB.Template("t_sms_template/get-enterprise-sms-template-count", s.Cond2Args(request))
	if err != nil {
		return 0
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(c, countResponseInfo, request)
	if getErr != nil {
		return 0
	}
	return countResponseInfo.Count
}

type GetDefaultTemplateRequest struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	TemplateName string `json:"template_name" db:"template_name"`
}

// GetDefaultTemplate 查询默认的短信模板
func (s *SmsTemplateRepo) GetDefaultTemplate(c ctx.Context, EnterpriseId int64, Name string) (*sms.SmsTemplate, error) {

	GetDefaultTemplateRequest := GetDefaultTemplateRequest{
		EnterpriseId: EnterpriseId,
		TemplateName: fmt.Sprintf(Name+"-%d", EnterpriseId),
	}

	stmt, err := s.DB.Template("t_sms_template/get-enterprise-default-template", s.Cond2Args(GetDefaultTemplateRequest))
	if err != nil {
		return nil, err
	}
	smsTemplate := &sms.SmsTemplate{}
	getErr := stmt.NamedGet(c, smsTemplate, GetDefaultTemplateRequest)
	if getErr != nil {
		return nil, getErr
	}
	return smsTemplate, nil
}

// SelectTemplateListByIDs 签名列表接口
func (s *SmsTemplateRepo) SelectTemplateListByIDs(c ctx.Context, EnterpriseId int64, templateIDs []int64, strategyIDs []int64) (*[]sms.SmsTemplate, error) {
	smsTemplateListRequest := sms.SmsTemplate{
		EnterpriseId: EnterpriseId,
		TemplateIDs:  templateIDs,
		StrategyIDs:  strategyIDs,
	}
	stmt, err := s.DB.Template("t_sms_template/template-list-by-ids", s.Cond2Args(smsTemplateListRequest))
	if err != nil {
		return nil, err
	}

	ListRequest := &[]sms.SmsTemplate{}
	res, errStmt := stmt.NamedQuery(c, smsTemplateListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}
