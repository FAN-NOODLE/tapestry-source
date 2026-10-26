package miniProgram

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/errorIgnore"
)

// Service  应用服务
type Service struct {
	repo Repository
}

// List 小程序列表
func (s Service) List(c ctx.Context, params *RequestPrograms) (int64, *ProgramList, error) {

	pageVO := assist.PageVO(params.Page, params.Size)
	findPageOption := &repos.FindPageOption{
		TplPrefix: "t_mini_program/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &ProgramList{}
		},
	}

	total, list, err := s.repo.FindPage(c, pageVO, params, findPageOption)
	if err != nil {
		return 0, nil, err
	}

	return total, list.(*ProgramList), nil
}

// Add 新增小程序
func (s Service) Add(c ctx.Context, params *MiniProgram) (int64, error) {

	if params.ID == 0 {
		params.ID = int64(assist.IdGenerate(c))
	}

	if params.CreatedAt == 0 {
		params.CreatedAt = assist.Timestamp()
	}

	_, err := s.repo.Add(c, params)

	if err != nil {
		return 0, err
	}

	return params.ID, nil
}

// Edit 修改
func (s Service) Edit(c ctx.Context, params *MiniProgram) (int64, error) {
	if params.UpdatedAt == 0 {
		params.UpdatedAt = assist.Timestamp()
	}

	_, err := s.repo.Update(c, params)
	if err != nil {
		return 0, err
	}

	return params.ID, nil
}

// Del 删除素材分组
func (s Service) Del(c ctx.Context, id int64) (int64, error) {
	return s.repo.Del(c, id)
}

// Authorize 维护授权小程序
func (s Service) Authorize(c ctx.Context, params *MiniProgram) (int64, error) {
	cc := ctx.Wrap(c)
	var id int64

	// 查询该企业对应的小程序是否存在，是：更新  否：新增
	id, _ = s.repo.AuthorizeMiniProgramDetail(cc, params.EnterpriseID)

	var err error
	if id > 0 {
		params.ID = id
		_, err = s.Edit(cc, params)
	} else {
		id, err = s.Add(cc, params)
	}

	if err != nil {
		return 0, err
	}

	return params.ID, nil
}

func (s Service) IsAuthOperate(c ctx.Context, id int64) bool {
	entityFactory := func(ctx.Context) interface{} {
		return &MiniProgram{}
	}
	vo, err := s.repo.Detail(c, id, entityFactory)

	if errorIgnore.Ignore(err) != nil {
		return false
	}

	detail := vo.(*MiniProgram)
	if detail.Type == TypeAuthorize {
		return false
	}

	return true
}
