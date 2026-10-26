package activityQRCodeSvc

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	activityQRCodeDao "micro-scrm/internal/app/dao/activity_qrcode"
	dashboardDao "micro-scrm/internal/app/dao/dashboard"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/app/http/request"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	activityQRCodeResp "micro-scrm/internal/app/http/response/activity_qrcode"
	"micro-scrm/internal/app/models"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type FrontSvc struct {
	service.Base
	mainDao        *activityQRCodeDao.MainDao
	memberScopeDao *activityQRCodeDao.MemberScopeDao
	pointDao       *activityQRCodeDao.PointDao
	contactWayDao  *activityQRCodeDao.ContactWayDao
}

func NewFrontSvc() *FrontSvc {
	return &FrontSvc{
		mainDao:        activityQRCodeDao.NewMainDao(),
		memberScopeDao: activityQRCodeDao.NewMemberScopeDao(),
		pointDao:       activityQRCodeDao.NewPointDao(),
		contactWayDao:  activityQRCodeDao.NewContactWayDao(),
	}
}

func (s *FrontSvc) GetQrCodeByUnionid(c ctx.Context, req activityQRCodeReq.GetQrCodeReq) (activityQRCodeResp.QrCodeResp, error) {
	resp := activityQRCodeResp.QrCodeResp{
		StoreList: []activityQRCodeResp.StoreInfo{},
	}
	activityInfo, errCode := s.getActivityInfo(c, req.EnterpriseId, req.State)
	if errCode > 0 {
		return activityQRCodeResp.QrCodeResp{ErrorCode: errCode}, nil
	}
	resp.ActivityName = activityInfo.AliasName
	resp.BannerUrl = activityInfo.BannerUrl

	if req.ThirdId > 0 {
		//如果选了门店 直接返回选的门店二维码
		return s.GetQrCodeByThirdId(c, req, activityInfo)
	} else {
		scopeMemberList := s.formatScopeMemberList(c, req, activityInfo)
		if len(scopeMemberList) == 0 {
			return activityQRCodeResp.QrCodeResp{ErrorCode: activityQRCodeResp.ActivityErrNoSa}, nil
		}
		relMemberInfo, err := s.getRelMemberInfo(c, req.EnterpriseId, req.Unionid, scopeMemberList)
		if err != nil {
			return activityQRCodeResp.QrCodeResp{}, err
		}
		if relMemberInfo.MemberId > 0 {
			//范围内有已经添加且符合条件的sa，返回sa信息
			resp.QrCodeMemberInfo = relMemberInfo
			return resp, nil
		}
		//获取导购管理柜台
		extraModel, _ := externalSearchDao.NewCdpExternalTagCk().GetOne(c, req.EnterpriseId,
			"unionid = ? AND tag_type = ?", req.Unionid, externalSearchModel.CdpTagManagerStoreCode)
		if extraModel.TagValue != "" {
			storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]any{
				"enterprise_id":    req.EnterpriseId,
				"third_store_code": extraModel.TagValue,
			})
			//如果有管理柜台则返回管理柜台相关门店联系我二维码
			req.ThirdId = int(storeDepartmentInfo.WechatDepartmentID)
			//检查管理柜台是否在活动范围内
			newQrUserIds, _, _ := s.getQrUserIdsAndStoreName(req, scopeMemberList)
			if len(newQrUserIds) > 0 {
				//有导购存在则说明在范围内，返回管理柜台二维码
				return s.GetQrCodeByThirdId(c, req, activityInfo)
			}
		}
		//不是会员或者管理柜台不存在或不在范围内则返回门店列表
		resp.StoreList, errCode = s.getStoreList(c, req, scopeMemberList)
		return resp, nil
	}

}

func (s *FrontSvc) GetAllStoreList(c ctx.Context, req activityQRCodeReq.GetQrCodeReq) (activityQRCodeResp.StoreListResp, error) {
	resp := activityQRCodeResp.StoreListResp{
		StoreList: []activityQRCodeResp.StoreInfo{},
	}
	activityInfo, errCode := s.getActivityInfo(c, req.EnterpriseId, req.State)
	if errCode > 0 {
		return activityQRCodeResp.StoreListResp{ErrorCode: errCode}, nil
	}
	scopeMemberList := s.formatScopeMemberList(c, req, activityInfo)
	if len(scopeMemberList) == 0 {
		return activityQRCodeResp.StoreListResp{ErrorCode: activityQRCodeResp.ActivityErrNoSa}, nil
	}
	//不是会员或者管理柜台不存在或不在范围内则返回门店列表
	resp.StoreList, errCode = s.getStoreList(c, req, scopeMemberList)
	return resp, nil
}
func (s *FrontSvc) getStoreList(
	c ctx.Context,
	req activityQRCodeReq.GetQrCodeReq,
	scopeMemberList []activityQRCodeResp.QrCodeMemberInfo,
) ([]activityQRCodeResp.StoreInfo, int) {
	var resp = make([]activityQRCodeResp.StoreInfo, 0)
	var thirdIds []int
	for _, info := range scopeMemberList {
		thirdIds = append(thirdIds, info.ThirdId)
	}
	storeDepartmentList, err := microWorkDao.NewStoreDepartment().GetAllByWechatDepartmentIds(c, req.EnterpriseId, helper.ArrayUniq(thirdIds))
	if err != nil || len(storeDepartmentList) == 0 {
		return resp, activityQRCodeResp.ActivityErrNoSa
	}
	for _, storeInfo := range storeDepartmentList {
		var cityStr = helper.RuleThree(strings.Contains(storeInfo.ProvinceZh, "市"), storeInfo.ProvinceZh, storeInfo.CityZh)
		tmp := activityQRCodeResp.StoreInfo{
			StoreCode:    storeInfo.ThirdStoreCode,
			ThirdId:      int(storeInfo.WechatDepartmentID),
			StoreName:    storeInfo.StoreName,
			StoreAddress: storeInfo.Address,
			City:         cityStr,
			CityPinYin:   strings.ToUpper(helper.GetPinYinFirstLetter(helper.CovertPinYin(cityStr))),
			Distance:     0,
		}
		if req.Lat > 0 && req.Lng > 0 {
			if storeInfo.Lat != "" && storeInfo.Lng != "" {
				lat, _ := strconv.ParseFloat(storeInfo.Lat, 64)
				lng, _ := strconv.ParseFloat(storeInfo.Lng, 64)
				tmp.Distance = int(helper.Haversine(
					req.Lat, req.Lng,
					lat, lng,
				))
			}
		}
		resp = append(resp, tmp)
	}
	sort.Slice(resp, func(i, j int) bool {
		// 如果 i 的距离为 0，应该排在后面
		if resp[i].Distance == 0 {
			return false
		}
		// 如果 j 的距离为 0，i 应该排在前面
		if resp[j].Distance == 0 {
			return true
		}
		// 两个距离都不为 0 时，距离小的排在前面
		return resp[i].Distance < resp[j].Distance
	})
	return resp, 0
}
func (s *FrontSvc) GetQrCodeByThirdId(c ctx.Context, req activityQRCodeReq.GetQrCodeReq, activityInfo activityQRCodeModel.ActivityQRCodeMainModel) (activityQRCodeResp.QrCodeResp, error) {
	resp := activityQRCodeResp.QrCodeResp{
		ActivityName: activityInfo.AliasName,
		BannerUrl:    activityInfo.BannerUrl,
	}
	contactWayModel, _ := s.contactWayDao.GetOne(c, req.EnterpriseId, "state = ? AND third_id = ?", req.State, req.ThirdId)
	if contactWayModel.ID > 0 {
		//找到了相关联系我二维码
		resp.QrCodeMemberInfo = activityQRCodeResp.QrCodeMemberInfo{
			QrCode:    contactWayModel.QrCode,
			StoreName: contactWayModel.StoreName,
			StoreCode: contactWayModel.StoreCode,
		}
		go func() {
			newQrUserIds, storeCode, storeName := s.getQrUserIdsAndStoreName(req, s.formatScopeMemberList(c, req, activityInfo))
			if len(newQrUserIds) > 0 {
				//如果二维码更新时间大于12小时则更新一次范围
				if time.Now().Unix()-contactWayModel.UpdatedAt > 3600*6 {
					_, err := qyWork.NewQyWork(req.EnterpriseId).UpdateContactWay(c, qyWork.ContactUpdateRequest{
						ConfigId:   contactWayModel.ConfigId,
						SkipVerify: helper.RuleThree(activityInfo.SkipVerify > 0, true, false),
						User:       newQrUserIds,
						State:      req.State,
					})
					if err != nil {
						log.Info(c).Msgf("更新联系我二维码失败,错误:%v", err)
					}
					_ = s.contactWayDao.UpdateData(c, req.EnterpriseId, fmt.Sprintf("id = %d", contactWayModel.ID), map[string]any{
						"store_code": storeCode,
						"store_name": storeName,
						"updated_at": time.Now().Unix(),
					})
				}
			}
		}()
		return resp, nil
	}
	scopeMemberList := s.formatScopeMemberList(c, req, activityInfo)
	if len(scopeMemberList) == 0 {
		return activityQRCodeResp.QrCodeResp{ErrorCode: activityQRCodeResp.ActivityErrNoSa}, nil
	}
	//没找到联系我二维码则重新生成
	contactWayResp, errCode := s.saveContactWay(c, req, activityInfo, scopeMemberList)
	if errCode > 0 {
		return activityQRCodeResp.QrCodeResp{ErrorCode: errCode}, nil
	}
	resp.QrCodeMemberInfo = activityQRCodeResp.QrCodeMemberInfo{
		QrCode:    contactWayResp.QrCode,
		StoreCode: contactWayResp.StoreCode,
		StoreName: contactWayResp.StoreName,
	}
	return resp, nil

}

func (s *FrontSvc) formatScopeMemberList(
	c ctx.Context,
	req activityQRCodeReq.GetQrCodeReq,
	activityInfo activityQRCodeModel.ActivityQRCodeMainModel,
) (scopeMemberList []activityQRCodeResp.QrCodeMemberInfo) {
	// 获取平均值缓存
	activityMemberScopeRedisKey := fmt.Sprintf("activity-qrcode:member-scope:%d-%d", req.EnterpriseId, activityInfo.ID)
	activityMemberScopeJson, _ := s.Redis().Get(c, activityMemberScopeRedisKey).Bytes()
	if activityMemberScopeJson != nil {
		//有缓存直接拿缓存
		_ = json.Unmarshal(activityMemberScopeJson, &scopeMemberList)
	}
	if scopeMemberList == nil || len(scopeMemberList) == 0 {
		scopeMemberList = s.getMemberListByScope(c, req.EnterpriseId, activityInfo.ID)
		// 写缓存
		activityMemberScopeByte, _ := json.Marshal(scopeMemberList)
		s.Redis().Set(c, activityMemberScopeRedisKey, string(activityMemberScopeByte), time.Hour*6)
	}
	return scopeMemberList
}
func (s *FrontSvc) getQrUserIdsAndStoreName(req activityQRCodeReq.GetQrCodeReq, scopeMemberList []activityQRCodeResp.QrCodeMemberInfo) (newQrUserIds []string, storeCode string, storeName string) {
	if len(scopeMemberList) == 0 {
		return []string{}, "", ""
	}
	var qrUserIds []string
	//用户选择了门店则直接生成联系我二维码
	for _, info := range scopeMemberList {
		if req.ThirdId > 0 && info.ThirdId == req.ThirdId {
			qrUserIds = append(qrUserIds, info.UserId)
			storeName = info.StoreName
			storeCode = info.StoreCode
		}
	}
	userLen := len(qrUserIds)
	if userLen > 100 {
		newQrUserIds = qrUserIds[0:99]
	} else {
		newQrUserIds = qrUserIds
	}

	return newQrUserIds, storeCode, storeName
}
func (s *FrontSvc) saveContactWay(c ctx.Context,
	req activityQRCodeReq.GetQrCodeReq,
	activityInfo activityQRCodeModel.ActivityQRCodeMainModel,
	scopeMemberList []activityQRCodeResp.QrCodeMemberInfo,
) (
	contactWayResp activityQRCodeModel.ActivityQRCodeContactWayModel,
	errCode int,
) {
	var qrUserIds []string
	var storeName string
	var storeCode string
	//用户选择了门店则直接生成联系我二维码
	for _, info := range scopeMemberList {
		if info.ThirdId == req.ThirdId {
			qrUserIds = append(qrUserIds, info.UserId)
			storeName = info.StoreName
			storeCode = info.StoreCode
		}
	}
	userLen := len(qrUserIds)
	if userLen == 0 {
		return contactWayResp, activityQRCodeResp.ActivityErrNoSa
	}
	var newQrUserIds []string
	if userLen > 100 {
		newQrUserIds = qrUserIds[0:99]
	} else {
		newQrUserIds = qrUserIds
	}
	qyResp, err := qyWork.NewQyWork(req.EnterpriseId).AddContactWay(c, qyWork.ContactCreateRequest{
		SkipVerify: helper.RuleThree(activityInfo.SkipVerify > 0, true, false),
		Type:       2,
		Scene:      2,
		State:      req.State,
		User:       newQrUserIds,
	})
	if err != nil {
		realUserIds := s.checkUserReal(c, req.EnterpriseId, newQrUserIds, req.ThirdId)
		if len(realUserIds) > 0 {
			qyResp, err = qyWork.NewQyWork(req.EnterpriseId).AddContactWay(c, qyWork.ContactCreateRequest{
				SkipVerify: helper.RuleThree(activityInfo.SkipVerify > 0, true, false),
				Type:       2,
				Scene:      2,
				State:      req.State,
				User:       realUserIds,
			})
			if err != nil {
				return contactWayResp, activityQRCodeResp.ActivityErrWeComErr
			}
		} else {
			return contactWayResp, activityQRCodeResp.ActivityErrWeComErr
		}

	}
	contactWayResp = activityQRCodeModel.ActivityQRCodeContactWayModel{
		BaseModel:  models.CreateBase(app.UID()),
		ActivityID: activityInfo.ID,
		State:      req.State,
		ThirdId:    req.ThirdId,
		StoreCode:  storeCode,
		StoreName:  storeName,
		ConfigId:   qyResp.ConfigId,
		QrCode:     qyResp.QrCode,
	}
	_ = s.contactWayDao.CreateData(c, req.EnterpriseId, contactWayResp)
	return contactWayResp, 0
}

func (s *FrontSvc) getActivityInfo(c ctx.Context, enterpriseId int64, state string) (info activityQRCodeModel.ActivityQRCodeMainModel, errCode int) {
	unixTime := int(time.Now().Unix())
	//获取活动信息
	pointInfo, _ := s.pointDao.GetOne(c, enterpriseId, "id = ?", state)
	if pointInfo.ID == 0 {
		return info, activityQRCodeResp.ActivityErrNotExists
	}
	activityId := pointInfo.ActivityID
	activityInfo, _ := s.mainDao.GetOne(c, enterpriseId, "id = ?", activityId)
	if activityInfo.ID == 0 {
		return info, activityQRCodeResp.ActivityErrNotExists
	}
	if activityInfo.TimeType != activityQRCodeModel.TimeTypeLongTerm {
		if activityInfo.StartTime > unixTime {
			return info, activityQRCodeResp.ActivityErrNotStart
		}
		if activityInfo.EndTime < unixTime {
			return info, activityQRCodeResp.ActivityErrExpired
		}
	}
	return activityInfo, 0
}

func (s *FrontSvc) getRelMemberInfo(c ctx.Context, enterpriseId int64, unionid string, scopeMemberList []activityQRCodeResp.QrCodeMemberInfo) (activityQRCodeResp.QrCodeMemberInfo, error) {

	var resp = activityQRCodeResp.QrCodeMemberInfo{}
	var scopeMemberIds []int64
	var memberMap = make(map[int64]activityQRCodeResp.QrCodeMemberInfo)
	for _, info := range scopeMemberList {
		scopeMemberIds = append(scopeMemberIds, info.MemberId)
		memberMap[info.MemberId] = info
	}
	relList, _ := microScrmDao.NewUserRelation().GetRelInfoByUnionIdAndMemberIdsWithOrder(c, enterpriseId, unionid, scopeMemberIds)
	if len(relList) == 0 {
		return activityQRCodeResp.QrCodeMemberInfo{}, nil
	}
	for _, rel := range relList {
		if rel.RelationStatus == microScrmModel.WorkUserRelationStatusNormal {
			resp = memberMap[rel.MemberID]
			break
		}
	}

	return resp, nil
}

func (s *FrontSvc) getMemberListByScope(c ctx.Context, enterpriseId, activityId int64) []activityQRCodeResp.QrCodeMemberInfo {
	var resp = make([]activityQRCodeResp.QrCodeMemberInfo, 0)
	//获取活动投放范围的人员
	memberScopes, _ := s.memberScopeDao.GetMany(c, enterpriseId, "activity_id = ?", activityId)
	// 转换成员范围
	var areaIds []int64
	var memberIds []int64
	for _, ms := range memberScopes {
		if ms.AreaId > 0 {
			areaIds = append(areaIds, int64(ms.AreaId))
		}
		if ms.MemberId > 0 {
			memberIds = append(memberIds, ms.MemberId)
		}
	}
	if len(areaIds) > 0 {
		thirdIds, _ := dashboardDao.NewArea().GetThirdIdsByAreaIds(c, dashboardReq.ManagerAuthAreaIds{
			EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: enterpriseId},
			AreaIds:           areaIds,
		})
		if len(thirdIds) > 0 {
			memberIdsByThirdIds, err := membersDao.NewMembers().GetMemberIdsByThirdIds(c, enterpriseId, thirdIds)
			if err != nil {
				return resp
			}
			memberIds = append(memberIds, memberIdsByThirdIds...)
		}
	}
	memberIds = helper.ArrayUniq(memberIds)
	if len(memberIds) == 0 {
		return resp
	}
	memberDepartmentList, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByMemberIds(c, enterpriseId, memberIds)
	if err != nil || len(memberDepartmentList) == 0 {
		return resp
	}
	allStoreDepartment, _ := microWorkDao.NewStoreDepartment().GetAll(c, enterpriseId)
	var storeDepartmentMap = make(map[int]microWorkModel.StoreDepartmentModel)
	for _, model := range allStoreDepartment {
		storeDepartmentMap[int(model.WechatDepartmentID)] = model
	}
	var tmp activityQRCodeResp.QrCodeMemberInfo
	// 整理导购id和导购信息的key=>value
	for _, info := range memberDepartmentList {
		if info.DeletedAt > 0 || info.Status != membersModel.MembersStatusEnable {
			continue
		}
		if !helper.InArray(strings.ToLower(info.Position), helper.GetSalesPosition()) {
			//过滤no sales
			continue
		}
		tmp.MemberId = info.ID
		tmp.UserId = info.Userid
		tmp.SaName = info.Name
		tmp.QrCode = info.QrCode
		tmp.Avatar = info.Avatar

		for _, department := range info.MemberDepartments {
			if department.DeletedAt > 0 {
				continue
			}
			if department.Department.ID > 0 && department.Department.DeletedAt == 0 {
				tmp.StoreName = department.Department.Name
				tmp.ThirdId = department.Department.ThirdID
				break
			}
		}
		if storeInfo, ok := storeDepartmentMap[tmp.ThirdId]; ok {
			tmp.StoreName = storeInfo.StoreName
			tmp.StoreCode = storeInfo.ThirdStoreCode
		}
		resp = append(resp, tmp)
	}
	return resp
}

func (s *FrontSvc) checkUserReal(c ctx.Context, enterpriseId int64, userIds []string, thirdId int) []string {
	var realUserIds []string
	activityCheckUserRealRedisKey := fmt.Sprintf("activity-qrcode:check-user-real:%d-%d", enterpriseId, thirdId)
	activityCheckUserRealJson, _ := s.Redis().Get(c, activityCheckUserRealRedisKey).Bytes()
	if activityCheckUserRealJson != nil {
		//有缓存直接拿缓存
		_ = json.Unmarshal(activityCheckUserRealJson, &realUserIds)
	}
	if realUserIds == nil || len(realUserIds) == 0 {
		resultChan := make(chan string, len(userIds))

		for _, userId := range userIds {
			go func(userId string) {
				qyResp, err := qyWork.NewQyWork(enterpriseId).AddContactWay(c, qyWork.ContactCreateRequest{
					SkipVerify: true,
					Type:       1,
					Scene:      2,
					State:      "test",
					User:       []string{userId},
				})

				if err == nil && qyResp.ConfigId != "" {
					// 直接删除，不等待结果
					go func() {
						_, _ = qyWork.NewQyWork(enterpriseId).DeleteContactWay(c, qyWork.ContactDeleteRequest{
							ConfigId: qyResp.ConfigId,
						})
					}()
					resultChan <- userId
				} else {
					resultChan <- "" // 无效用户也发送空字符串
				}
			}(userId)
		}

		for i := 0; i < len(userIds); i++ {
			if userId := <-resultChan; userId != "" {
				realUserIds = append(realUserIds, userId)
			}
		}
		// 写缓存
		realUserIdsBytes, _ := json.Marshal(realUserIds)
		s.Redis().Set(c, activityCheckUserRealRedisKey, string(realUserIdsBytes), time.Minute*5)
	}

	return realUserIds
}
