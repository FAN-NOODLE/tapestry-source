package uc

// Member 成员详情
type Member struct {
	MemUid      string         `json:"memUid"`
	MemUserID   string         `json:"memUserid"`
	MemMobile   string         `json:"memMobile"`
	MemActive   int            `json:"memActive"`
	MemUsername string         `json:"memUsername"`
	MemIndex    string         `json:"memIndex"`
	MemFace     string         `json:"memFace"`
	MemQrCode   string         `json:"memQrCode"`
	DpName      []MemberDpName `json:"dpName"`
}

type MemberDpName struct {
	DpId   string `json:"dpId"`
	EpId   string `json:"epId"`
	DpName string `json:"dpName"`
}

// MemberByUserIdRequest
type MemberByUserIdRequest struct {
	MemUserID string `json:"memUserid"`
}

type ListResp struct {
	Total    int           `json:"total"`
	PageSize int           `json:"pageSize"`
	PageNum  int           `json:"pageNum"`
	ExtList  []interface{} `json:"extList"`
}
type DepartmentListResp struct {
	ListResp
	List []DepartmentInfoResp `json:"list"`
}

type DepartmentInfoResp struct {
	DpId                     string  `json:"dpId"`
	EpId                     string  `json:"epId"`
	DpName                   string  `json:"dpName"`
	DpParentid               string  `json:"dpParentid"`
	DpLevel                  int     `json:"dpLevel"`
	DpDisplayorder           float64 `json:"dpDisplayorder"`
	DpCreated                int64   `json:"dpCreated"`
	DpUpdated                int64   `json:"dpUpdated"`
	DpThirdid                string  `json:"dpThirdid"`
	DptId                    string  `json:"dptId"`
	IsChildDepartment        int     `json:"isChildDepartment"`
	DepartmentPath           string  `json:"departmentPath"`
	ChildrensDepartmentCount int     `json:"childrensDepartmentCount"`
	DepartmentMemberCount    int     `json:"departmentMemberCount"`
	DpSerialNum              string  `json:"dpSerialNum"`
	DpLevelCount             int     `json:"dpLevelCount"`
}
