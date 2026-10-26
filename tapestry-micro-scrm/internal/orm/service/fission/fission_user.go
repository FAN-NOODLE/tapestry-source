package fission

import (
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
)

func (f *FissionService) UpdateFissionUserInfo(enterpriseID, activityID, userID int64, newAddNum int) error {
	return f.Model(&fission.FissionUserModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"activity_id":   activityID,
		"user_id":       userID,
	}).Updates(map[string]interface{}{
		"cur_assistant_num": gorm.Expr("cur_assistant_num + ?", newAddNum),
	}).Error
}

func (f *FissionService) GetStyleConfigList(eid int64) (dataList []FissionStyleConfig, err error) {
	err = f.Model(&fission.FissionStyleConfigModel{}).
		Select("id", "styles", "rgb_colors", "is_check").
		Where("enterprise_id=?", eid).
		Find(&dataList).Error
	return
}

func (f *FissionService) SetCheckStyle(enterpriseID, id int64) error {
	err := f.Model(&fission.FissionStyleConfigModel{}).
		Where("enterprise_id", enterpriseID).
		Update("is_check", 0).Error
	if err != nil {
		return err
	}
	err = f.Model(&fission.FissionStyleConfigModel{}).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "id": id}).
		Update("is_check", 1).Error
	return err
}

func (f *FissionService) GetBgConfig(enterpriseID int64) (dataList fission.FissionBgConfigModel, err error) {
	err = f.Model(&fission.FissionBgConfigModel{}).
		Where("enterprise_id", enterpriseID).
		First(&dataList).Error
	return
}

func (f *FissionService) GetBgConfigList(enterpriseID int64) (dataList []FissionBgConfig, err error) {
	err = f.Model(&fission.FissionBgConfigModel{}).
		Where("enterprise_id", enterpriseID).
		Find(&dataList).Error
	return
}

func (f *FissionService) SetBgConfig(enterpriseID int64, bgImg, shareImg string) error {
	err := f.Model(&fission.FissionBgConfigModel{}).
		Where("enterprise_id", enterpriseID).
		Updates(map[string]interface{}{"bg_img": bgImg, "share_img": shareImg}).Error
	if err != nil {
		return err
	}
	return err
}

func (f *FissionService) GetCheckStyle(enterpriseID int64) (dataList FissionStyleConfig, err error) {
	err = f.Model(&fission.FissionStyleConfigModel{}).
		Select("id", "is_check", "styles", "rgb_colors").
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "is_check": 1}).
		First(&dataList).Error
	return
}

func (f *FissionService) AddStyleConfig(model fission.FissionStyleConfigModel) error {
	return f.Create(&model).Error
}

func (f *FissionService) AddBgConfig(model fission.FissionBgConfigModel) error {
	return f.Create(&model).Error
}

func (f *FissionService) DelStyleConfig(enterpriseID int64, id []int64) error {
	return f.Delete(fission.FissionStyleConfigModel{}, "id in ? and enterprise_id=?", id, enterpriseID).Error
}

//func (f *FissionService) GetJoinActivityByGroup(enterpriseID int64, unionIds []string) (data []fission.UserJoinNums, err error) {
//	var userIds []string
//	f.Select("id").Where("union_id in ?", unionIds).Model(&fission.FissionMpUserModel{}).Pluck("id", &userIds)
//	err = f.Exec("SELECT u.user_id,mp.nick_name,mp.img_url as `avatar`,count( u.user_id ) AS join_num "+
//		"FROM fission_user AS u LEFT JOIN fission_mp_user AS mp ON u.user_id = mp.id"+
//		" WHERE u.user_id IN ? and u.enterprise_id=?"+
//		" GROUP BY u.user_id ORDER BY join_num DESC", userIds, enterpriseID).Find(&data).Error
//	return
//}
