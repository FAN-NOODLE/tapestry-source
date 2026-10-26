package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/pkg/db"
)

// NewTagRepo 初始化仓库
func NewTagRepo() tag.Repository {
	return &tagRepo{repos.NewSQLRepo("t_tag", db.DB())}
}

// 标签
func NewWorkTagRepo() tag.WorkRepository {
	return &workTagRepo{repos.NewSQLRepo("t_work_external_tag", db.DB())}
}

// 标签组
func NewWorkTagGroupRepo() tag.WorkTagGroupRepository {
	return &workTagGroupRepo{repos.NewSQLRepo("t_work_external_tag_group", db.DB())}
}

// 客户分组
func NewUserGroupRepo() tag.UserGroupRepository {
	return &userGroupRepo{repos.NewSQLRepo("t_user_group", db.DB())}
}
