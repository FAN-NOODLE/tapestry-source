package cdp

var (
	// UserTagList 标签分组列表
	UserTagList = "%s/scrm/user/tag/list"
	// ManualGroupCreate 创建手动标签分组
	ManualGroupCreate = "%s/scrm/user/tag/manual/group/create"
	// CreateGroup 创建客户分组
	CreateGroup = "%s/scrm/user/group/create"
	// GroupUserListURL 客户分组列表
	GroupUserListURL = "%s/scrm/user/group/list"
)

// URLEnum cdp接口url配置
var URLEnum = struct {
	UserTagList       string // 标签分组列表
	ManualGroupCreate string // 创建手动标签分组
	ManualTagCreate   string // 创建手动标签
	AutoTagUpdate     string // 更新自动标签
}{
	UserTagList:       "%s/scrm/user/tag/list",
	ManualGroupCreate: "%s/scrm/user/tag/manual/group/create",
	ManualTagCreate:   "%s/scrm/user/tag/manual/tag/create",
	AutoTagUpdate:     "%s/scrm/user/tag/update",
}
