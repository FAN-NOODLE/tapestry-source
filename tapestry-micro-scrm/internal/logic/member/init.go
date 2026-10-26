package member

// 员工信息查询的相关逻辑

type memberLogic struct {
	Name         string
	enterpriseID int64
}

func NewMemberLogic(enterpriseID int64) *memberLogic {
	return &memberLogic{
		Name:         "员工信息查询的相关逻辑",
		enterpriseID: enterpriseID,
	}
}
