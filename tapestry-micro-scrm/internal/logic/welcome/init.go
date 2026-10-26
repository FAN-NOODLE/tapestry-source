package welcome

// 获取指定用户的欢迎语

type welcomeUserLogic struct {
	Name         string
	enterpriseID int64
}

func NewWelcomeUserLogic(enterpriseID int64) *welcomeUserLogic {
	return &welcomeUserLogic{
		Name:         "用户欢迎语逻辑层",
		enterpriseID: enterpriseID,
	}
}
