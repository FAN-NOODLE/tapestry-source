package qyWork

type MemberListReq struct {
	DepartmentID int64 `json:"department_id" desc:"获取的部门id"`
}

// MemberDetailResp 成员详情
type MemberDetailResp struct {
	Userid         string  `json:"userid"`
	Name           string  `json:"name"`
	Department     []int64 `json:"department"`
	Position       string  `json:"position"`
	Mobile         string  `json:"mobile"`
	Gender         string  `json:"gender"`
	Email          string  `json:"email"`
	IsLeaderInDept []int8  `json:"is_leader_in_dept"`
	Avatar         string  `json:"avatar"`
	ThumbAvatar    string  `json:"thumb_avatar"`
	Telephone      string  `json:"telephone"`
	Alias          string  `json:"alias"`
	Address        string  `json:"address"`
	OpenUserid     string  `json:"open_userid"`
	MainDepartment int32   `json:"main_department"`
	Status         int8    `json:"status"`
	QrCode         string  `json:"qr_code"`
	Extattr        Extattr `json:"extattr"`
}

type Attrs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Extattr struct {
	Attrs []Attrs `json:"attrs"`
}
