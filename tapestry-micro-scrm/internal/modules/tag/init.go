package tag

// NewTag 初始化标签服务
func NewTag(repository Repository) *Service {
	return &Service{
		Repo: repository,
	}
}

// NewTagGroup 初始化标签组服务
func NewTagGroup(repository Repository) *GroupService {
	return &GroupService{
		Repo: repository,
	}
}

// work系列标签服务
func NewWorkExternalTag(repository WorkRepository) *WorkExternalTagService {
	return &WorkExternalTagService{
		Repo: repository,
	}
}

// work系列标签组服务
func NewWorkExternalTagGroup(repository WorkTagGroupRepository) *WorkExternalTagGroupService {
	return &WorkExternalTagGroupService{
		Repo: repository,
	}
}

// 用户分组服务
func NewUserGroup(repository UserGroupRepository) *UserGroupService {
	return &UserGroupService{
		Repo: repository,
	}
}
