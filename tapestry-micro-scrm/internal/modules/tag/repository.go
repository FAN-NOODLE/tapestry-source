package tag

import (
	"micro-scrm/internal/modules/tag/entity"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 仓库接口
type Repository interface {
	repos.IRepository
	QueryTags(ctx.Context, []int64) ([]*entity.TagEntity, error)
	QueryEnterpriseTags(ctx.Context, int64) ([]*entity.TagEntity, error)
	QueryTagGroups(ctx.Context, []int64) ([]*entity.TagGroupEntity, error)
	QueryUserTagGroups(ctx.Context, []int64) ([]*entity.UserGroupEntity, error)
	QueryEnterpriseTagGroups(ctx.Context, int64) ([]*entity.TagGroupEntity, error)
	AddTagGroup(ctx.Context, []*entity.TagGroupEntity) (int64, error)
	GetTagList(ctx.Context, int64, string) ([]*entity.TagEntity, error)
	GetUserGroupList(ctx.Context, int64, string) ([]*entity.UserGroupEntity, error)
	AddTag(ctx.Context, []*entity.TagEntity) (int64, error)
	AddUserGroup(ctx.Context, []*entity.UserGroupEntity) (int64, error)
	GetTagListPage(ctx.Context, *entity.GetTagListRequest) (int64, []*entity.TagEntity, error)
	GetUserGroupListPage(ctx.Context, *entity.GetTagListRequest) (int64, []*entity.UserGroupEntity, error)
	UpdateTagGroup(ctx.Context, *entity.TagGroupEntity) (int64, error)
	DeleteTagGroup(c ctx.Context, enterpriseID int64, tagGroupIDList []int64) (int64, error)
	UpdateTag(c ctx.Context, data *entity.TagEntity) (int64, error)
	DeleteTag(c ctx.Context, enterpriseID int64, tagIDList []int64) (int64, error)
	UpdateUserGroup(c ctx.Context, data *entity.UserGroupEntity) (int64, error)
	DeleteUserGroup(c ctx.Context, enterpriseID int64, userGroupIDList []int64) (int64, error)
	WorkExternalTagList(cc ctx.Context, enterpriseID int64, relationID int64) (*[]entity.WorkExternalTag, error)
	// 获取企业对应的标签信息，每条标签信息包含组信息
	GetTagInfoList(c ctx.Context, enterpriseID int64, tagType, isThirdTag int) ([]entity.WorkTagInfoEntity, error)
	// 根据标签查询客户分组
	UserGroupByTag(ctx.Context, int64, []int64) ([]entity.UserGroupDetailTagInfo, error)
	// 获取所有在使用的企业微信tag，tagid集合
	GetAllUsedExternalTags(cc ctx.Context, enterpriseID int64) ([]int64, error)
	QueryWorkExternalTagList(c ctx.Context, tagIDs []int64, relationID int64) ([]entity.WorkExternalTag, error)
}

// work_external_tag_group 相关
type WorkTagGroupRepository interface {
	repos.IRepository
	// 根据企业ID和标签名称，检测标签信息是否存在
	GetTagGroupInfoByName(ctx.Context, entity.CreateManualTagGroupParams) (bool, error)
	// 添加标签组信息
	AddTagGroup(ctx.Context, entity.WorkExternalTagGroupModel) (int64, error)
	// 添加标签
	AddTag(ctx.Context, entity.WorkExternalTagModel) error
	// 检测标签组中，是否有外部联系人使用的标签
	IsExistExternalUser(ctx.Context, entity.DeleteManualTagGroupRequestParam) (*entity.WorkExternalTagGroupEntity, error)
	// 删除标签组/标签
	DeleteTagGroup(ctx.Context, entity.DeleteTagGroup) (int64, error)
	// 检测标签名称是是否已经存在
	IsExistTagName(ctx.Context, int64, []string) (bool, error)
}

// work_external_tag 相关
type WorkRepository interface {
	repos.IRepository
	// 根据标签组ID，获取标签组内所有的标签
	QueryTagListByGroupId(ctx.Context, entity.TagParam) ([]entity.WorkExternalTagEntity, error)
	// 获取企业下的所有的标签组信息
	GetTagGroupListByEnterpriseId(ctx.Context, int64) ([]entity.WorkExternalTagGroupEntity, error)
	// 根据标签组名称，模糊匹配标签组信息
	MatchTagGroupByGroupName(ctx.Context, entity.SearchTagParam) ([]entity.WorkExternalTagGroupEntity, error)
	// 检测标签ID集合下是否存在外部联系人
	CheckTagHasExternalUser(ctx.Context, entity.TagOperationParam) ([]entity.WorkExternalTagEntity, error)
	// 根据企业微信标签ID，删除标签
	DeleteTag(ctx.Context, entity.TagOperationParam) (int64, error)
	// 根据企业ID和标签ID，获取标签信息
	GetTagInfoByID(ctx.Context, entity.UpdateManualTagParam) (*entity.WorkExternalTagModel, error)
	// 根据企业ID和标签名称，判断标签是否存在
	IsExistTagByName(ctx.Context, entity.UpdateManualTagParam) (bool, error)
	// 根据企业ID和标签ID，更新标签的名称
	UpdateTagNameByID(ctx.Context, entity.UpdateManualTagParam) error
	// 根据企业ID和标签组ID，获取标签组信息
	GetTagGroupInfoByID(ctx.Context, entity.CreateManualTagParam) (*entity.WorkExternalTagGroupModel, error)
	// 新增标签
	AddTag(ctx.Context, entity.WorkExternalTagModel) (int64, error)
}

// 客户分组接口
type UserGroupRepository interface {
	repos.IRepository
	// 检测分组名是否被使用
	IsExist(ctx.Context, string) (bool, error)
	// 新增客户分组
	AddUserGroup(ctx.Context, entity.UserGroupModelEntity) (int64, error)
	// 新增客户详情
	AddUserGroupDetail(ctx.Context, []entity.UserGroupDetailModelEntity) (int64, error)
	UpdateUserGroup(c ctx.Context, enterpriseID int64, id int64, attrs map[string]interface{}) error
	// 查询所有的客户分组
	GetUserGroupList(ctx.Context, int64) ([]entity.UserGroupInfo, error)
	// 通过客户分组ID，获取外部联系人集合
	GetExternalUserIDList(ctx.Context, entity.QueryUserGroupParam) ([]int64, error)
	// 获取客户分组详情
	GetUserGroupDetail(ctx.Context, entity.QueryUserGroupParam) ([]entity.UserGroupDetailTagInfo, error)
	// 根据标签ID集合，获取关联该标签的用户集合
	GetExternalUserIDByTagIDList(ctx.Context, []int64, int64) ([]int64, error)
	// 通过客户分组ID，查询客户分组关联的外部联系人

}
