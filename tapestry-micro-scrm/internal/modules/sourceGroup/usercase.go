package sourceGroup

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/assist"
)

// Service  应用服务
type Service struct {
	repo Repository
}

// List 素材分组列表
func (s Service) List(c ctx.Context, params *SourceGroup) (*[]SourceGroup, error) {

	listOption := &repos.ListOption{
		TplPrefix: "t_source_group/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]SourceGroup{}
		},
	}

	list, err := s.repo.List(c, params, listOption)
	if err != nil {
		return nil, err
	}

	return list.(*[]SourceGroup), nil
}

// Add 新增素材分组
func (s Service) Add(c ctx.Context, params *SourceGroup) (int64, error) {

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

// IsRepeatName 素材分组是否重名
func (s Service) IsRepeatName(c ctx.Context, params *SourceGroup) (bool, error) {
	return s.repo.IsRepeatName(c, params)
}

// Edit 修改素材分组
func (s Service) Edit(c ctx.Context, params *SourceGroup) (int64, error) {
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

// BatchDel 批量删除素材分组
func (s Service) BatchDel(c ctx.Context, ids []int64) (int64, error) {
	return s.repo.BatchDel(c, ids)
}
