package sourceGroup

type ID struct {
	ID int64 `json:"id"`
}

type Rows struct {
	Rows int64 `json:"rows"`
}

type SourceGroup struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	GroupName    string `json:"group_name" db:"group_name"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

type NewSource struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	GroupName    string `json:"group_name" db:"group_name"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
	ParentID     int64  `json:"parent_id"`
}

type EditSource struct {
	ID           int64  `json:"id" db:"id"  validate:"required" required:"素材分组ID不能为空"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	GroupName    string `json:"group_name" db:"group_name" validate:"required" required:"素材名称不能为空"`
	IgnoreIndex  int    `json:"ignore_index"`
}

type UpdateIgnoreIndexReq struct {
	ID           int64 `json:"id" db:"id"  validate:"required" required:"素材分组ID不能为空"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	IgnoreIndex  int   `json:"ignore_index"`
}

type RemoveSource struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type MoveReceive struct {
	MoverId      int64 `json:"mover_id"`
	SwaperId     int64 `json:"swaper_id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type FindSource struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}
