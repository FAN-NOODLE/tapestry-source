package plan_task

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/uid"
	"time"
)

func GetTypeList(enterpirse_id int64, param BaseReq) (int64, []micro_scrm_enterprise.OperationActivityModel, error) {
	total, data, err := orm.Service.Activity.GetTypeList(enterpirse_id, param.Page, param.Size)
	if err != nil {
		return 0, nil, err
	}
	for index, item := range data {
		s, ok := micro_scrm_enterprise.STATUSTEXT[item.Status]
		if !ok {
			s = "进行中"
		}
		data[index].StatusText = s
		data[index].StartTimeText = time.Unix(item.StartTime, 0).Format("2006-01-02")
		data[index].EndTimeText = time.Unix(item.EndTime, 0).Format("2006-01-02")
	}
	return total, data, err
}

func GetActivityId(enterpriseID int64) (data int64, err error) {
	data, err = orm.Service.Activity.GetActivityId(enterpriseID)
	return
}

func CheckActivityStatus(param CheckStatusReq) (data micro_scrm_enterprise.OperationActivityUserModel, err error) {
	data, err = orm.Service.Activity.GetInfoByUnionid(param.EnterpriseID, param.UnionId, param.ActivityId)
	return
}

func InsertUserData(c ctx.Context, param InsertUserReq) error {
	data, err := orm.Service.Activity.GetInfoByUnionid(param.EnterpriseID, param.UnionId, param.ActivityId)
	if err != nil {
		return err
	}
	if data.ID == 0 {
		var NewUser micro_scrm_enterprise.UserModel
		jsonStr, _ := json.Marshal(param)
		UnmarshalErr := json.Unmarshal(jsonStr, &NewUser)
		if UnmarshalErr != nil {
			return UnmarshalErr
		}
		NewUser.Status = 1
		NewUser.ID = uid.ID(c)
		NewUser.WxName = param.WxName
		AesEncryptName, encryptErr := common.AesEncrypt(NewUser.ReceiverName)
		if encryptErr != nil {
			return encryptErr
		}
		NewUser.ReceiverName = AesEncryptName
		AesEncryptPhone, encryptPhoneErr := common.AesEncrypt(NewUser.ReceiverPhone)
		if encryptPhoneErr != nil {
			return encryptPhoneErr
		}
		NewUser.ReceiverPhone = AesEncryptPhone

		user, GetNameErr := orm.Service.Activity.GetUserName(param.EnterpriseID, param.UnionId)
		if GetNameErr != nil {
			return GetNameErr
		}
		if user.ID != 0 {
			EnterpriseUser, GetPhoneErr := orm.Service.Activity.GetUserPhone(param.EnterpriseID, user.Euid)
			if GetPhoneErr != nil {
				return GetPhoneErr
			}
			if EnterpriseUser.ID != 0 {
				NewUser.WxPhone = EnterpriseUser.Mobile
				NewUser.MemberCode = EnterpriseUser.UnionCode
			}
		}

		insertErr := orm.Service.Activity.InsertUser(param.EnterpriseID, NewUser)
		if insertErr != nil {
			return insertErr
		}
		UpdateErr := orm.Service.Activity.UpdateActivityNums(param.EnterpriseID, param.ActivityId)
		if UpdateErr != nil {
			return UpdateErr
		}
	}
	return nil
}

func ExportUserList(interpriseId int64, param ExportUsreReq) ([]micro_scrm_enterprise.OperationActivityUserModel, error) {
	data, err := orm.Service.Activity.ExportUserList(interpriseId, param.ActivityId, param.StartTime, param.EndTime, param.Status)
	if err != nil {
		return []micro_scrm_enterprise.OperationActivityUserModel{}, err
	}
	for index, item := range data {
		data[index].CreatedTimeText = time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04")
		str, ok := micro_scrm_enterprise.USERSTATUSTEXT[item.Status]
		if !ok {
			str = "待审核"
		}
		if len(item.ReceiverPhone) > 0 {
			decodePhone, decodePhoneErr := common.AesDecrypt(item.ReceiverPhone)
			if decodePhoneErr != nil {
				return []micro_scrm_enterprise.OperationActivityUserModel{}, decodePhoneErr
			}
			data[index].ReceiverPhone = decodePhone
		}
		if len(item.ReceiverName) > 0 {
			decodeName, decodeNameErr := common.AesDecrypt(item.ReceiverName)
			if decodeNameErr != nil {
				return []micro_scrm_enterprise.OperationActivityUserModel{}, decodeNameErr
			}
			data[index].ReceiverName = decodeName
		}
		data[index].StatusText = str
	}
	return data, nil
}

func GetUserListByTypeId(interpriseId int64, param UserListReq) (int64, []micro_scrm_enterprise.OperationActivityUserListModel, error) {
	total, data, err := orm.Service.Activity.GetUserList(interpriseId, param.ActivityId, param.StartTime, param.EndTime, param.Status, param.Page, param.Size)
	if err != nil {
		return 0, []micro_scrm_enterprise.OperationActivityUserListModel{}, err
	}
	for index, item := range data {
		if len(item.ReceiverPhone) > 0 {
			decodePhone, decodePhoneErr := common.AesDecrypt(item.ReceiverPhone)
			if decodePhoneErr != nil {
				return 0, []micro_scrm_enterprise.OperationActivityUserListModel{}, decodePhoneErr
			}
			data[index].ReceiverPhone = HiddenPhone(decodePhone)
		}
		if len(item.ReceiverName) > 0 {
			decodeName, decodeNameErr := common.AesDecrypt(item.ReceiverName)
			if decodeNameErr != nil {
				return 0, []micro_scrm_enterprise.OperationActivityUserListModel{}, decodeNameErr
			}
			data[index].ReceiverName = HiddenName(decodeName)
		}

		data[index].CreatedTimeText = time.Unix(item.CreatedAt, 0).Format("2006-01-02 15:04")
		str, ok := micro_scrm_enterprise.USERSTATUSTEXT[item.Status]
		if !ok {
			str = "待审核"
		}
		data[index].StatusText = str
	}
	return total, data, nil
}

func GetUserInfo(enterpriseId int64, activityUserId int64) (micro_scrm_enterprise.OperationActivityUserModel, error) {
	data, err := orm.Service.Activity.GetUserInfo(enterpriseId, activityUserId)
	if err != nil {
		return micro_scrm_enterprise.OperationActivityUserModel{}, err
	}
	if len(data.WxPhone) > 0 {
		data.WxPhone = HiddenPhone(data.WxPhone)
	}
	if len(data.ReceiverName) > 0 {
		decodeName, decodeNameErr := common.AesDecrypt(data.ReceiverName)
		if decodeNameErr != nil {
			return micro_scrm_enterprise.OperationActivityUserModel{}, decodeNameErr
		}
		data.ReceiverName = HiddenName(decodeName)
	}
	if len(data.ReceiverPhone) > 0 {
		decodePhone, decodePhoneErr := common.AesDecrypt(data.ReceiverPhone)
		if decodePhoneErr != nil {
			return micro_scrm_enterprise.OperationActivityUserModel{}, decodePhoneErr
		}
		data.ReceiverPhone = HiddenPhone(decodePhone)
	}
	str, ok := micro_scrm_enterprise.USERSTATUSTEXT[data.Status]
	if !ok {
		str = "待审核"
	}
	data.StatusText = str
	return data, nil
}

func HiddenName(name string) string {
	nameRune := []rune(name)
	var result string
	lens := len(nameRune)
	if lens <= 1 {
		result = "***"
	} else if lens == 2 {
		result = string(nameRune[:1]) + "*"
	} else if lens == 3 {
		result = string(nameRune[:1]) + "**"
	} else if lens == 4 {
		result = string(nameRune[:1]) + "***"
	} else if lens > 4 {
		result = string(nameRune[:2]) + "***"
	}
	return result
}

func HiddenPhone(phone string) string {
	str := []rune(phone)
	return string(str[0:3]) + "****" + string(str[7:11])
}

func VerifyUser(enterpriseId int64, activityUserId int64) error {
	return orm.Service.Activity.VerifyUser(enterpriseId, activityUserId)
}
