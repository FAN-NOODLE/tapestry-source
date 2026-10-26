package smart_shopping_guide

type smartShoppingGuideLogic struct {
	name         string
	enterpriseID int64
}

func NewSmartShoppingGuideLogic(enterpriseID int64) *smartShoppingGuideLogic {
	return &smartShoppingGuideLogic{
		name:         "智慧导购相关逻辑",
		enterpriseID: enterpriseID,
	}
}
