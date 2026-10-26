package externalUserService

import (
	"micro-scrm/internal/app"
	externalUserDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/http/request/external"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/sdk/microwork"
)

type ExternalUser struct {
	wxUserCorpExternalDao *externalUserDao.WxUserCorpExternal
}

func NewWxUserCorpExternal() *ExternalUser {
	return &ExternalUser{
		wxUserCorpExternalDao: externalUserDao.NewWxUserCorpExternal(),
	}
}

func (d *ExternalUser) GetWxUserCorpExternalCount(c app.Ctx, enterpriseID int64) (total int64, err error) {
	return d.wxUserCorpExternalDao.GetCount(c, enterpriseID)
}

func (d *ExternalUser) GetWxUserCorpExternalList(c app.Ctx, enterpriseID int64, Page, Size int) ([]microScrmModel.WxUserCorpExternalModel, error) {
	return d.wxUserCorpExternalDao.GetList(c, enterpriseID, Page, Size)
}

func (d *ExternalUser) SyncExternalUserJob(c app.Ctx, enterpriseID int64) error {

	pageLimit := 500
	beginID := int64(0)

	for {
		var rcResult []microScrmModel.WxUserCorpExternalModel

		rcResult, err := d.wxUserCorpExternalDao.GetListLimit(c, enterpriseID, pageLimit, beginID)
		if err != nil {
			return err
		}

		if len(rcResult) > 0 {
			for _, row := range rcResult {

				// 查询work请求crm
				crmParams := microwork.UserInfoRequest{}
				crmParams.EnterpriseID = enterpriseID
				crmParams.UnionId = row.Unionid
				UserInfoResponse, errs := microwork.NewCrm().GetUserInfoByUnionId(c, crmParams)

				if errs != nil {
					return errs
				}

				if UserInfoResponse.Code == 200 && UserInfoResponse.Data.Status == true {

					UserInfoData := UserInfoResponse.Data
					// 修改数据库字段 wx_user_corp_external
					err1 := orm.Service.WorkExternalUser.UpdateUserCorpExternalByEuId(enterpriseID, row.Euid, UserInfoData.Membership, UserInfoData.Birthday, UserInfoData.Name, UserInfoData.Mobile, UserInfoData.BirthYear, UserInfoData.BirthMonth, UserInfoData.BirthDay, UserInfoData.CustomerLevel, UserInfoData.ClientType, UserInfoData.Last12mPurchaseAmount, UserInfoData.CustomerCode)
					if err1 != nil {
						return err1
					}
					//err2 := orm.Service.WorkExternalUser.UpdateEnterpriseUserById(enterpriseID, row.Euid, UserInfoData.Membership, UserInfoData.Birthday, UserInfoData.Name, UserInfoData.Mobile)
					//if err2 != nil {
					//	return err2
					//}

					// 处理更新
					err3 := orm.Service.WorkExternalUser.CHHandleExternalUserLog(c, enterpriseID, row.Unionid, UserInfoData.Membership, UserInfoData.Birthday, UserInfoData.Name, UserInfoData.Mobile)
					if err3 != nil {
						return err3
					}
				}
			}
			beginID = rcResult[len(rcResult)-1].ID
		}

		if len(rcResult) < pageLimit {
			break
		}
	}
	return nil
}

// SyncCdpInfo 同步cdp的用户会员信息
func (d *ExternalUser) SyncCdpInfo(c app.Ctx, enterpriseID int64, info external.CdpExternalInfo, updateData interface{}) error {
	return d.wxUserCorpExternalDao.SyncCdpInfo(c, enterpriseID, info, updateData)
}

// GetWxUserCorpExternalByUnionIds 根据unionids获取外部联系人信息
func (d *ExternalUser) GetWxUserCorpExternalByUnionIds(c app.Ctx, enterpriseID int64, unionIds []string) ([]microScrmModel.WxUserCorpExternalModel, error) {
	return d.wxUserCorpExternalDao.GetWxUserCorpExternalByUnionIds(c, enterpriseID, unionIds)
}

// GetNickNameAndRemarkByUnionIds   unionIds
func (d *ExternalUser) GetNickNameAndRemarkByUnionIds(c app.Ctx, enterpriseID int64, memberId int64, unionIds []string) ([]external.CoachNameInfo, error) {
	var list = make([]external.CoachNameInfo, 0)
	externalList, err := d.wxUserCorpExternalDao.GetWxUserCorpExternalByUnionIds(c, enterpriseID, unionIds)
	if err != nil {
		return nil, err
	}
	relationList, _ := externalUserDao.NewUserRelation().GetMany(c, "enterprise_id = ? and member_id = ? AND unionid in ?", enterpriseID, memberId, unionIds)
	var remarkMap = make(map[string]string)
	for _, relation := range relationList {
		remarkMap[relation.Unionid] = relation.Remark
	}
	for _, info := range externalList {
		tmp := external.CoachNameInfo{
			Unionid: info.Unionid,
			Nick:    info.Nick,
			Remark:  "",
		}
		if remark, ok := remarkMap[info.Unionid]; ok {
			tmp.Remark = remark
		}
		list = append(list, tmp)
	}
	return list, nil
}
