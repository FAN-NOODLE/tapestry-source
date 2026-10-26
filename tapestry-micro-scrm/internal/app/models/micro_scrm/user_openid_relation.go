package microScrmModel

import "micro-scrm/internal/app/models"

type WxUserOpenidRelationModel struct {
	models.BaseModel
	WxOpenid string `json:"wx_openid" gorm:"column:wx_openid;type:varchar(100);not null;uniqueIndex:unionid_index;comment:微信端用户的openid"`
	Unionid  string `json:"unionid" gorm:"column:unionid;type:varchar(100);not null;uniqueIndex:unionid_index;comment:用户的unionid"`
}

// TableName returns the table name of the WxUserOpenidRelation model
func (WxUserOpenidRelationModel) TableName() string {
	return "wx_user_openid_relation"
}
