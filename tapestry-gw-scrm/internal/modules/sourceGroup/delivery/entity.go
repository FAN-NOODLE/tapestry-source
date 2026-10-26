package delivery

type MoveReceive struct {
	MoverId      int64 `json:"mover_id"`
	SwaperId     int64 `json:"swaper_id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
}

type NewSourceGroup struct {
	ID           int64             `json:"id" db:"id"`
	EnterpriseID int64             `json:"enterprise_id" db:"enterprise_id"`
	GroupName    string            `json:"group_name" db:"group_name"`
	CreatedAt    int64             `json:"created_at" db:"created_at"`
	UpdatedAt    int64             `json:"updated_at" db:"updated_at"`
	DeletedAt    int64             `json:"deleted_at" db:"deleted_at"`
	ParentID     int64             `json:"parent_id" db:"parent_id"`
	Orders       int               `gorm:"column:orders;NOT NULL" json:"orders"`
	Children     []*NewSourceGroup `json:"children"`
}
