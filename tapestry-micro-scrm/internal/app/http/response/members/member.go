package membersResp

type MemberForShenCeResp struct {
	ID           int64  `json:"id"`
	Userid       string `json:"userid"`
	Mobile       string `json:"mobile"`
	Position     string `json:"position"`
	Type         int    `json:"type"`
	Gender       int    `json:"gender"`
	Avatar       string `json:"avatar"`
	Name         string `json:"name"`
	SVDepartment string `json:"sv_department" gorm:"column:sv_department"`
	SVManager    int    `json:"sv_manager" gorm:"column:sv_manager"`
}
