package source

type CategoryList struct {
	Name         string `json:"name"`
	EnterpriseID int64  `json:"enterprise_id"`
	AdminId      int64  `json:"admin_id"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

type CategoryCreate struct {
	Category
	Name string `json:"name"`
}

type CategoryUpdate struct {
	Category
	Name string `json:"name"`
}

type Category struct {
	Id           int64 `json:"id" comment:"Id"`
	EnterpriseID int64 `json:"enterprise_id"`
	AdminId      int64 `json:"admin_id"`
}
