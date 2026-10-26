package drainageContactsSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	"micro-scrm/internal/app/dao/drainage"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	drainageContactGroupReq "micro-scrm/internal/app/http/request/drainage_contact"
	drainageContactResp "micro-scrm/internal/app/http/response/drainage_contact"
	"micro-scrm/internal/app/models"
	drainageContactModel "micro-scrm/internal/app/models/drainage"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MainSvc struct {
	service.Base
	contactGroupDao  *drainageContactDao.DrainageContactGroupDao
	contactsDao      *drainageContactDao.DrainageContactsDao
	contactMemberDao *drainageContactDao.DrainageContactMemberDao
	welcomeDao       *drainageContactDao.DrainageContactGroupWelcomesDao
	memberDao        *membersDao.Members
}

func NewMainSvc() *MainSvc {
	return &MainSvc{
		contactGroupDao:  drainageContactDao.NewDrainageContactGroupDao(),
		contactsDao:      drainageContactDao.NewDrainageContactsDao(),
		contactMemberDao: drainageContactDao.NewDrainageContactMemberDao(),
		welcomeDao:       drainageContactDao.NewDrainageContactGroupWelcomesDao(),
		memberDao:        membersDao.NewMembers(),
	}
}

func (s *MainSvc) checkConsoleContactGroup(groupInfo drainageContactModel.DrainageContactGroupModel) error {
	if groupInfo.ID == 0 {
		return errors.New("联系我活动不存在")
	}
	if groupInfo.SourceType != drainageContactModel.DrainageContactGroupSourceConsole {
		return errors.New("特殊联系我活动不允许操作")
	}
	return nil
}

func (s *MainSvc) ContactGroupDetail(c ctx.Context, req drainageContactGroupReq.InfoReq) (res drainageContactResp.GroupDetailResp, err error) {
	groupInfo, err := s.contactGroupDao.GetOnePreloadWelcomes(c, req.EnterpriseId, req.ID)
	if err != nil {
		return res, err
	}
	if err = s.checkConsoleContactGroup(groupInfo); err != nil {
		return res, err
	}
	res.ID = groupInfo.ID
	res.Name = groupInfo.Name
	res.Type = groupInfo.Type
	res.ChannelId = groupInfo.ChannelID
	res.AutoFriend = groupInfo.AutoFriend
	if groupInfo.Welcomes != nil && len(groupInfo.Welcomes) > 0 {
		for _, welcome := range groupInfo.Welcomes {
			var contentJson []map[string]interface{}
			_ = json.Unmarshal([]byte(welcome.Content), &contentJson)
			res.WelcomeList = append(res.WelcomeList, drainageContactResp.WelcomeInfo{
				Type:    welcome.Type,
				Content: contentJson,
			})
		}
	}
	res.HasWelcome = helper.RuleThree(len(res.WelcomeList) > 0, 1, 0)
	var memberIds = strings.Split(groupInfo.MemberIds, `,`)
	var thirdIds = strings.Split(groupInfo.ThirdIds, `,`)
	res.MemberData = make([]drainageContactResp.MemberData, 0)
	if len(memberIds) > 0 {
		memberIds = helper.ArrayUniq(memberIds)
		memberList, _ := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND id in ? AND deleted_at = 0", req.EnterpriseId, memberIds)
		var memberNameMap = make(map[int64]string)
		for _, model := range memberList {
			memberNameMap[model.ID] = model.Name
		}
		for _, id := range memberIds {
			tmpId := int64(helper.StringToInt(id))
			if name, ok := memberNameMap[tmpId]; ok {
				res.MemberData = append(res.MemberData, drainageContactResp.MemberData{
					MemberId:   tmpId,
					MemberName: name,
				})
			}
		}
	}
	res.DepartmentData = make([]drainageContactResp.DepartmentData, 0)
	if len(thirdIds) > 0 {
		thirdIds = helper.ArrayUniq(thirdIds)
		deptList, _ := departmentsDao.NewDepartments().GetMany(c, "enterprise_id = ? AND third_id in ? AND deleted_at = 0", req.EnterpriseId, thirdIds)
		var deptNameMap = make(map[int]string)
		for _, model := range deptList {
			deptNameMap[model.ThirdID] = model.Name
		}
		for _, id := range thirdIds {
			tmpId := helper.StringToInt(id)
			if name, ok := deptNameMap[tmpId]; ok {
				res.DepartmentData = append(res.DepartmentData, drainageContactResp.DepartmentData{
					ThirdId:        tmpId,
					DepartmentName: name,
				})
			}
		}
	}
	return res, nil
}
func (s *MainSvc) ContactGroupDetailList(c ctx.Context, req drainageContactGroupReq.DetailListReq) (res drainageContactResp.GroupInfoListResp, err error) {
	groupInfo, err := s.contactGroupDao.GetOnePreloadWelcomes(c, req.EnterpriseId, req.ID)
	if err != nil {
		return res, err
	}
	if err = s.checkConsoleContactGroup(groupInfo); err != nil {
		return res, err
	}
	res.Name = groupInfo.Name
	res.Type = groupInfo.Type
	res.ID = groupInfo.ID
	res.HasWelcome = helper.RuleThree(groupInfo.Welcomes != nil && len(groupInfo.Welcomes) > 0, 1, 0)
	channelInfo, _ := drainageContactDao.NewDrainageContactChannelDao().GetOnePreload(c, "enterprise_id= ? AND id = ?", req.EnterpriseId, groupInfo.ChannelID)
	res.ChannelName = channelInfo.Name
	res.Tags = make([]string, 0)
	if channelInfo.Labels != nil && len(channelInfo.Labels) > 0 {
		for _, label := range channelInfo.Labels {
			if label.DeletedAt == 0 && label.Tags.ID > 0 && label.Tags.DeletedAt == 0 {
				res.Tags = append(res.Tags, label.Tags.Name)
			}
		}
	}

	list, total, err := s.contactsDao.GetPage(c, req)
	if err != nil || total == 0 {
		return res, err
	}
	res.Total = total
	var stateSlice []string
	var memberIds []int64
	var thirdIds []int
	for _, contact := range list {
		stateSlice = append(stateSlice, helper.Int64ToString(contact.ID))
		if len(contact.Members) > 0 {
			for _, member := range contact.Members {
				if member.MemberId > 0 {
					memberIds = append(memberIds, member.MemberId)
				}
				if member.ThirdId > 0 {
					thirdIds = append(thirdIds, member.ThirdId)
				}
			}
		}
	}
	//查询用户及门店编码
	var memberNameMap = make(map[int64]string)
	if len(memberIds) > 0 {
		memberList, _ := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND id in ? AND deleted_at = 0", req.EnterpriseId, memberIds)
		for _, model := range memberList {
			memberNameMap[model.ID] = model.Name
		}
	}
	var deptNameMap = make(map[int]string)
	if len(thirdIds) > 0 {
		deptList, _ := departmentsDao.NewDepartments().GetMany(c, "enterprise_id = ? AND third_id in ? AND deleted_at = 0", req.EnterpriseId, thirdIds)
		for _, model := range deptList {
			deptNameMap[model.ThirdID] = model.Name
		}
	}
	//查询添加好友统计
	addCustomerTotalMap, _ := microScrmDao.NewUserRelation().GetCountMapByStateSlice(c, req.EnterpriseId, stateSlice)
	for _, contact := range list {
		var tmpMemberNames = make([]string, 0)
		var tmpDepartmentData = make([]drainageContactResp.DepartmentData, 0)
		var tmpRealMembers = make([]drainageContactResp.RealMember, 0)
		if len(contact.Members) > 0 {
			for _, member := range contact.Members {
				if member.ThirdId > 0 {
					if name, ok := deptNameMap[member.ThirdId]; ok {
						tmpDepartmentData = append(tmpDepartmentData, drainageContactResp.DepartmentData{
							ThirdId:        member.ThirdId,
							DepartmentName: name,
						})
					}
				}
				if member.MemberId > 0 {
					if name, ok := memberNameMap[member.MemberId]; ok {
						tmpMemberNames = append(tmpMemberNames, name)
					}
				}
			}

			if contact.RealUserIds != "" {
				realUserIds := strings.Split(contact.RealUserIds, ",")
				memberDepartmentList, _ := membersDao.NewMembers().GetMemberInfoAndDepartmentByUserIds(c, req.EnterpriseId, realUserIds)
				for _, info := range memberDepartmentList {
					if len(info.MemberDepartments) > 0 {
						for _, departmentRelation := range info.MemberDepartments {
							if departmentRelation.DeletedAt == 0 && departmentRelation.Department.DeletedAt == 0 {
								tmpRealMembers = append(tmpRealMembers, drainageContactResp.RealMember{
									MemberId:   info.ID,
									MemberName: info.Name,
									ThirdId:    departmentRelation.Department.ThirdID,
								})
								break
							}
						}
					}
				}
			}
		}
		var addCustomerTotal int64
		if tmpTotal, ok := addCustomerTotalMap[helper.Int64ToString(contact.ID)]; ok {
			addCustomerTotal = tmpTotal
		}
		res.ContactList = append(res.ContactList, drainageContactResp.ContactInfo{
			ID:                contact.ID,
			Name:              contact.Name,
			CreatedAt:         contact.CreatedAt,
			Members:           tmpMemberNames,
			Departments:       tmpDepartmentData,
			RealMembers:       tmpRealMembers,
			AddCustomerNumber: addCustomerTotal,
			QrCodeUrl:         contact.QrCodeURL,
		})
	}
	return res, nil
}

// ContactGroupList   活动列表
func (s *MainSvc) ContactGroupList(c ctx.Context, req drainageContactGroupReq.GroupListReq) ([]drainageContactResp.GroupBaseResp, int64, error) {
	var resList = make([]drainageContactResp.GroupBaseResp, 0)
	list, total, err := s.contactGroupDao.GetPage(c, req)
	if err != nil || total == 0 {
		return resList, 0, err
	}
	var groupIds []int64
	var channelIds []int64
	var stateSlice []string
	for _, model := range list {
		groupIds = append(groupIds, model.ID)
		channelIds = append(channelIds, model.ChannelID)
		if len(model.Contacts) > 0 {
			for _, contact := range model.Contacts {
				stateSlice = append(stateSlice, helper.Int64ToString(contact.ID))
			}
		}
	}
	welcomeList, _ := s.welcomeDao.GetMany(c, req.EnterpriseId, "contact_group_id in ?", groupIds)
	var welcomeMap = make(map[int64]bool)
	for _, model := range welcomeList {
		welcomeMap[model.ContactGroupId] = true
	}
	channelList, _ := drainageContactDao.NewDrainageContactChannelDao().GetManyPreload(c, "enterprise_id= ? AND id in ?", req.EnterpriseId, channelIds)
	var channelMap = make(map[int64]drainageContactModel.DrainageContactChannelsModel)
	for _, model := range channelList {
		channelMap[model.ID] = model
	}
	addCustomerTotalMap, _ := microScrmDao.NewUserRelation().GetCountMapByStateSlice(c, req.EnterpriseId, stateSlice)
	for _, model := range list {
		codeNumber := len(model.Contacts)
		channelName := ""
		var tags = make([]string, 0)
		if channel, ok := channelMap[model.ChannelID]; ok {
			channelName = channel.Name
			if channel.Labels != nil && len(channel.Labels) > 0 {
				for _, label := range channel.Labels {
					if label.DeletedAt == 0 && label.Tags.ID > 0 && label.Tags.DeletedAt == 0 {
						tags = append(tags, label.Tags.Name)
					}
				}

			}

		}
		var addCustomerTotal int64
		for _, contact := range model.Contacts {
			if contact.DeletedAt == 0 {
				if addTotal, ok := addCustomerTotalMap[helper.Int64ToString(contact.ID)]; ok {
					addCustomerTotal += addTotal
				}
			}
		}
		resList = append(resList, drainageContactResp.GroupBaseResp{
			ID:                model.ID,
			Name:              model.Name,
			Type:              model.Type,
			CodeNumber:        codeNumber,
			ChannelName:       channelName,
			Tags:              tags,
			HasWelcome:        helper.RuleThree(model.Welcomes != nil && len(model.Welcomes) > 0, 1, 0),
			AddCustomerNumber: addCustomerTotal,
			CreatedAt:         model.CreatedAt,
		})
	}
	return resList, total, nil
}

// ContactGroupDelete    活动删除
func (s *MainSvc) ContactGroupDelete(c ctx.Context, req drainageContactGroupReq.InfoReq) error {
	nowUnix := time.Now().Unix()
	contactsGroupInfo, err := s.contactGroupDao.GetOnePreloadContacts(c, req.EnterpriseId, req.ID)
	if err != nil {
		return err
	}
	if err = s.checkConsoleContactGroup(contactsGroupInfo); err != nil {
		return err
	}
	var delContactIds []int64
	var delConfigIds []string
	if contactsGroupInfo.Contacts != nil && len(contactsGroupInfo.Contacts) > 0 {
		for _, contact := range contactsGroupInfo.Contacts {
			if contact.DeletedAt == 0 {
				delContactIds = append(delContactIds, contact.ID)
				delConfigIds = append(delConfigIds, contact.ConfigID)
			}
		}
	}
	qyWorkClient := qyWork.NewQyWork(req.EnterpriseId)
	for _, id := range delConfigIds {
		_, _ = qyWorkClient.DeleteContactWay(c, qyWork.ContactDeleteRequest{
			ConfigId: id,
		})
	}
	err = s.contactGroupDao.UpdateData(c, req.EnterpriseId, map[string]any{
		"id": req.ID,
	}, map[string]any{
		"deleted_at": nowUnix,
	})
	if err != nil {
		return err
	}
	_ = s.contactsDao.UpdateData(c, map[string]any{
		"contact_group_id": req.ID,
		"enterprise_id":    req.EnterpriseId,
	}, map[string]any{
		"deleted_at": nowUnix,
	})

	_ = s.contactMemberDao.UpdateData(c, map[string]any{
		"contact_id":    delContactIds,
		"enterprise_id": req.EnterpriseId,
	}, map[string]any{
		"deleted_at": nowUnix,
	})
	return nil
}

// ContactGroupSave  活动保存
func (s *MainSvc) ContactGroupSave(c ctx.Context, req drainageContactGroupReq.GroupSaveReq) (drainageContactResp.GroupSaveResp, error) {
	result := drainageContactResp.GroupSaveResp{}
	qyWorkClient := qyWork.NewQyWork(req.EnterpriseId)
	var contactsGroupInfo drainageContactModel.DrainageContactGroupModel
	if req.ID > 0 {
		//获取原始信息
		var err error
		contactsGroupInfo, err = s.contactGroupDao.GetOnePreloadContactAndMember(c, req.EnterpriseId, req.ID)
		if err != nil || contactsGroupInfo.ID == 0 {
			return result, err
		}
		if err = s.checkConsoleContactGroup(contactsGroupInfo); err != nil {
			return result, err
		}
	}
	req.MemberIds = helper.ArrayUniq(req.MemberIds)
	var paramsAllMemberIds = req.MemberIds
	var deptName string
	if len(req.ThirdIds) > 0 {
		req.ThirdIds = helper.ArrayUniq(req.ThirdIds)
		//获取部门人员
		dpMemberList, _ := s.memberDao.GetMemberByThirdIds(c, req.EnterpriseId, req.ThirdIds, []int64{})
		storeDeptList, err := microWorkDao.NewStoreDepartment().GetMany(c, "enterprise_id = ? AND wechat_department_id in ?", req.EnterpriseId, req.ThirdIds)
		if err != nil {
			return result, err
		}
		if len(dpMemberList) > 0 {
			deptName = dpMemberList[0].StoreName
			for _, storeDept := range storeDeptList {
				deptName = storeDept.StoreName
				break
			}
			for _, member := range dpMemberList {
				paramsAllMemberIds = append(paramsAllMemberIds, member.MembersId)
			}
		}
	}
	members, _ := s.memberDao.GetMany(c, "enterprise_id = ? AND id in ? AND deleted_at=0", req.EnterpriseId, paramsAllMemberIds)
	if len(members) == 0 {
		return result, errors.New("您选择的部门/员工中没有在岗人员")
	}
	var trueMembers []membersModel.MemberModel
	for _, member := range members {
		if helper.InArray(member.ID, req.MemberIds) || helper.InArray(strings.ToLower(member.Position), helper.GetSalesPosition()) {
			//共享码过滤no sales
			trueMembers = append(trueMembers, member)
		}
	}
	if len(trueMembers) == 0 {
		return result, errors.New("您选择的门店/人员中没有sales人员")
	}
	//去重
	if len(trueMembers) > 100 {
		return result, errors.New("共享码总的人员数量不能超过100")
	}

	var existSingleContacts = make([]drainageContactModel.DrainageContactsModel, 0)
	var err error

	if req.ID > 0 {
		//处理共享码更新逻辑以及个人码的删除逻辑
		existSingleContacts, err = s.contactGroupUpdate(c, req, contactsGroupInfo, trueMembers, deptName, qyWorkClient)
		if err != nil {
			unverifiedUsers, err := s.handleUnverifiedError(c, req.EnterpriseId, err, trueMembers)
			if len(unverifiedUsers) > 0 {
				result.UnverifiedUsers = unverifiedUsers
			}
			return result, err
		}
	} else {
		//新增联系我活动
		req.ID = app.UID()
		//处理活动及共享码添加逻辑
		err = s.contactGroupAdd(c, req, trueMembers, deptName, qyWorkClient)
		if err != nil {
			_ = s.contactGroupDao.UpdateData(c, req.EnterpriseId, map[string]any{
				"id": req.ID,
			}, map[string]any{
				"deleted_at": time.Now().Unix(),
			})
			unverifiedUsers, err := s.handleUnverifiedError(c, req.EnterpriseId, err, trueMembers)
			if len(unverifiedUsers) > 0 {
				result.UnverifiedUsers = unverifiedUsers
			}
			return result, err
		}
	}

	if req.Type == drainageContactModel.DrainageContactGroupStore {
		//处理个人码
		err = s.contactSaAdd(c, req, contactsGroupInfo, existSingleContacts, trueMembers, qyWorkClient)
		if err != nil {
			unverifiedUsers, err := s.handleUnverifiedError(c, req.EnterpriseId, err, trueMembers)
			if len(unverifiedUsers) > 0 {
				result.UnverifiedUsers = unverifiedUsers
			}
			return result, err
		}
	}

	_ = s.welcomeSave(c, req)
	result.ID = req.ID
	return result, nil
}

func (s *MainSvc) ZipUrl(c ctx.Context, req drainageContactGroupReq.ZipUrlReq) (string, error) {
	// 生成ZIP文件名（包含宽度信息）
	nowTime := time.Now().Format(helper.YYYYMMDDHHMMSS)
	if !helper.InArray(req.Width, []int{258, 430, 1417}) {
		return "", nil
	}
	var widthMap = make(map[int]string, 0)
	widthMap[258] = "小尺寸"
	widthMap[430] = "中尺寸"
	widthMap[1417] = "大尺寸"
	contactsGroupInfo, err := s.contactGroupDao.GetOnePreloadContacts(c, req.EnterpriseId, req.ID)
	if err != nil {
		return "", err
	}
	if err = s.checkConsoleContactGroup(contactsGroupInfo); err != nil {
		return "", err
	}

	var imageMap = make(map[string]string)
	if contactsGroupInfo.Contacts != nil && len(contactsGroupInfo.Contacts) > 0 {
		for i, contact := range contactsGroupInfo.Contacts {
			if contact.DeletedAt == 0 && contact.QrCodeURL != "" {
				imageMap[fmt.Sprintf("%s_%s_%d.png", contact.Name, widthMap[req.Width], i+1)] = contact.QrCodeURL
			}
		}
	}

	if len(imageMap) == 0 {
		return "", fmt.Errorf("没有可用的二维码图片")
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "zip_images_*")
	if err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建ZIP文件映射
	files := make(map[string]string)

	for fileName, imageURL := range imageMap {
		if fileName == "" || imageURL == "" {
			continue
		}
		tmpFileName := helper.Md5V(fileName)
		// 下载原图到临时文件
		originalPath := filepath.Join(tempDir, tmpFileName)
		if err := helper.DownloadImage(imageURL, originalPath); err != nil {
			return "", fmt.Errorf("下载图片失败 (URL: %s): %w", imageURL, err)
		}

		// 上传原图到COS
		imageCosUrl, err := helper.UploadFileToCos(c, originalPath, false)
		if err != nil {
			return "", fmt.Errorf("上传图片失败: %w", err)
		}

		// 生成指定宽度的图片URL
		resizedImageUrl := fmt.Sprintf("%s?imageMogr2/thumbnail/%dx", imageCosUrl, req.Width)

		// 生成缩放后的文件名（可选，保留原文件名也可以）
		fileExt := filepath.Ext(fileName)
		fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
		resizedFileName := fmt.Sprintf("%s_%s_%s", fileNameWithoutExt, nowTime, fileExt)
		resizedLocalPath := filepath.Join(tempDir, resizedFileName)

		// 下载缩放后的图片
		if err := helper.DownloadFileToLocal(resizedImageUrl, resizedLocalPath); err != nil {
			return "", fmt.Errorf("下载缩放图片失败 (URL: %s): %w", resizedImageUrl, err)
		}

		// 添加到文件映射，使用缩放文件名
		files[resizedLocalPath] = resizedFileName
	}

	if len(files) == 0 {
		return "", fmt.Errorf("没有可打包的图片文件")
	}

	// 确保ZIP目录存在
	zipDir := "./storage/drainage_contact/"
	if err := os.MkdirAll(zipDir, 0755); err != nil {
		return "", fmt.Errorf("创建ZIP目录失败: %w", err)
	}

	zipName := fmt.Sprintf("%s_%s_%s", contactsGroupInfo.Name, widthMap[req.Width], nowTime)
	zipFileName := filepath.Join(zipDir, zipName+".zip")

	// 创建ZIP
	if err := helper.CreateZip(zipFileName, files); err != nil {
		return "", fmt.Errorf("创建ZIP失败: %w", err)
	}

	// 确保清理本地ZIP文件
	defer func() {
		if err := os.Remove(zipFileName); err != nil {
			log.Info(c).Msgf("清理ZIP文件失败, path:%s, 错误:%v", zipFileName, err)
		}
	}()

	// 上传ZIP到COS
	zipURL, err := helper.UploadFileToCos(c, zipFileName, false)
	if err != nil {
		return "", fmt.Errorf("上传ZIP失败: %w", err)
	}

	return zipURL, nil
}

func (s *MainSvc) UpdateContact(c ctx.Context, enterpriseId int64) error {
	qyWorkClient := qyWork.NewQyWork(enterpriseId)
	contactsDao := drainageContactDao.NewDrainageContactsDao()
	nowUnix := time.Now().Unix()
	//上线前的数据不更新
	contactsList, err := s.contactsDao.GetManyPreloadThirdMembers(c, "enterprise_id = ? AND created_at > 1765938600 AND contact_group_id > 0 AND is_contact_group_single=0", enterpriseId)
	if err != nil || len(contactsList) == 0 {
		return err
	}
	var thirdIds []int
	for _, model := range contactsList {
		if model.Members != nil && len(model.Members) > 0 {
			for _, member := range model.Members {
				if member.ThirdId > 0 {
					thirdIds = append(thirdIds, member.ThirdId)
				}
			}
		}
	}
	thirdIds = helper.ArrayUniq(thirdIds)
	dpMemberList, _ := s.memberDao.GetMemberByThirdIds(c, enterpriseId, thirdIds, []int64{})
	var usersByThirdIdMap = make(map[int][]string)
	for _, member := range dpMemberList {
		if !helper.InArray(strings.ToLower(member.Position), helper.GetSalesPosition()) {
			//过滤no sales
			continue
		}
		if len(usersByThirdIdMap[member.StoreId]) > 0 {
			usersByThirdIdMap[member.StoreId] = append(usersByThirdIdMap[member.StoreId], member.MembersUserid)
		} else {
			usersByThirdIdMap[member.StoreId] = []string{member.MembersUserid}
		}
	}
	memberList, _ := s.memberDao.GetMany(c, "enterprise_id=? AND deleted_at = 0", enterpriseId)
	var userIdMap = make(map[int64]string)
	for _, model := range memberList {
		userIdMap[model.ID] = model.Userid
	}
	for _, model := range contactsList {
		var tmpThirdIds []int
		var tmpUserIds []string
		if model.Members != nil && len(model.Members) > 0 {
			for _, member := range model.Members {
				if member.DeletedAt == 0 {
					if member.ThirdId > 0 {
						tmpThirdIds = append(tmpThirdIds, member.ThirdId)
					} else if member.MemberId > 0 {
						if userId, ok := userIdMap[member.MemberId]; ok {
							tmpUserIds = append(tmpUserIds, userId)
						}
					}
				}
			}
		}
		if len(tmpThirdIds) > 0 {
			var allUserIds []string
			for _, thirdId := range tmpThirdIds {
				if userIds, ok := usersByThirdIdMap[thirdId]; ok {
					allUserIds = append(allUserIds, userIds...)
				}
			}
			allUserIds = append(allUserIds, tmpUserIds...)
			allUserIds = helper.ArrayUniq(allUserIds)
			if len(allUserIds) > 0 {
				//如果
				_, err := qyWorkClient.UpdateContactWay(c, qyWork.ContactUpdateRequest{
					ConfigId:   model.ConfigID,
					SkipVerify: helper.RuleThree(model.AutoFriend > 0, true, false),
					User:       allUserIds,
					State:      helper.Int64ToString(model.ID),
				})
				if err != nil {
					log.Info(c).Msgf("更新联系我活动共享码失败,%d, 错误:%v", model.ID, err)
				}
				_ = contactsDao.UpdateData(c, map[string]any{
					"id": model.ID,
				}, map[string]interface{}{
					"real_user_ids": strings.Join(allUserIds, ","),
					"updated_at":    nowUnix,
				})
			}
		}

	}

	return nil
}
func (s *MainSvc) contactGroupUpdate(
	c ctx.Context,
	req drainageContactGroupReq.GroupSaveReq,
	contactsGroupInfo drainageContactModel.DrainageContactGroupModel,
	members []membersModel.MemberModel,
	deptName string,
	qyWorkClient *qyWork.QyWork,
) ([]drainageContactModel.DrainageContactsModel, error) {
	var singleContacts = make([]drainageContactModel.DrainageContactsModel, 0)
	nowUnix := time.Now().Unix()
	var newUserIds []string
	var memberMap = make(map[int64]membersModel.MemberModel)
	for _, member := range members {
		newUserIds = append(newUserIds, member.Userid)
		memberMap[member.ID] = member
	}
	if contactsGroupInfo.Type != req.Type {
		return singleContacts, errors.New("编辑时不能修改活动类型")
	}
	var storeContact drainageContactModel.DrainageContactsModel
	for _, contact := range contactsGroupInfo.Contacts {
		if contact.DeletedAt == 0 && len(contact.Members) > 0 {
			if contact.IsContactGroupSingle == 0 {
				storeContact = contact
			} else {
				singleContacts = append(singleContacts, contact)
			}
		}
	}

	if storeContact.ID > 0 {
		if len(newUserIds) > 0 {
			_, err := qyWorkClient.UpdateContactWay(c, qyWork.ContactUpdateRequest{
				ConfigId:   storeContact.ConfigID,
				SkipVerify: helper.RuleThree(req.AutoFriend > 0, true, false),
				User:       newUserIds,
				State:      helper.Int64ToString(storeContact.ID),
			})
			if err != nil {
				log.Info(c).Msgf("更新联系我二维码失败,userid:%s,错误:%v", strings.Join(newUserIds, ","), err)
				return nil, err
			}
			_ = s.contactMemberDao.UpdateData(c, map[string]any{
				"enterprise_id": req.EnterpriseId,
				"contact_id":    storeContact.ID,
			}, map[string]interface{}{
				"deleted_at": nowUnix,
			})
			_ = s.contactsDao.UpdateData(c, map[string]any{
				"enterprise_id": req.EnterpriseId,
				"id":            storeContact.ID,
			}, map[string]interface{}{
				"real_user_ids": strings.Join(newUserIds, ","),
				"name":          helper.RuleThree(req.Type == drainageContactModel.DrainageContactGroupStore, helper.RuleThree(deptName != "", deptName, "门店共享码")+"-"+req.Name, req.Name),
				"channel_id":    req.ChannelId,
				"auto_friend":   req.AutoFriend,
				"updated_at":    nowUnix,
			})
		} else {
			_, _ = qyWorkClient.DeleteContactWay(c, qyWork.ContactDeleteRequest{
				ConfigId: storeContact.ConfigID,
			})
			_ = s.contactsDao.UpdateData(c, map[string]any{
				"enterprise_id": req.EnterpriseId,
				"id":            storeContact.ID,
			}, map[string]interface{}{
				"deleted_at": nowUnix,
			})
			storeContact.ID = 0
		}
	} else {
		if len(newUserIds) > 0 {
			storeContactID := app.UID()
			contactWayRes, err := qyWorkClient.AddContactWay(c, qyWork.ContactCreateRequest{
				SkipVerify: helper.RuleThree(req.AutoFriend > 0, true, false),
				Type:       2,
				Scene:      2,
				State:      helper.Int64ToString(storeContactID),
				User:       newUserIds,
			})
			if err != nil {
				return singleContacts, err
			}
			_ = s.contactsDao.CreateData(c, drainageContactModel.DrainageContactsModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, storeContactID),
				Name:           helper.RuleThree(deptName != "", deptName, "门店共享码") + "-" + req.Name,
				RealUserIds:    strings.Join(newUserIds, ","),
				ContactGroupId: req.ID,
				ChannelID:      req.ChannelId,
				AutoFriend:     req.AutoFriend,
				QrCodeURL:      contactWayRes.QrCode,
				ConfigID:       contactWayRes.ConfigId,
				Type:           1,
			})
			storeContact = drainageContactModel.DrainageContactsModel{MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, storeContactID)}
		}
	}

	if storeContact.ID > 0 {
		var contactMemberSlice = make([]drainageContactModel.DrainageContactMemberModel, 0)
		for _, id := range req.ThirdIds {
			contactMemberSlice = append(contactMemberSlice, drainageContactModel.DrainageContactMemberModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
				ContactId:      storeContact.ID,
				ThirdId:        id,
			})
		}
		if len(req.MemberIds) > 0 {
			for _, memberId := range req.MemberIds {
				contactMemberSlice = append(contactMemberSlice, drainageContactModel.DrainageContactMemberModel{
					MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
					ContactId:      storeContact.ID,
					MemberId:       memberId,
				})
			}
		}
		_ = s.contactMemberDao.CreateInBatchesData(c, contactMemberSlice)
	}
	//更新联系我活动表信息
	_ = s.contactGroupDao.UpdateData(c, req.EnterpriseId, map[string]any{
		"id": req.ID,
	}, map[string]interface{}{
		"member_ids":  helper.SliceToString(req.MemberIds),
		"third_ids":   helper.SliceToString(req.ThirdIds),
		"name":        req.Name,
		"channel_id":  req.ChannelId,
		"auto_friend": req.AutoFriend,
		"updated_at":  nowUnix,
	})

	return singleContacts, nil
}

func (s *MainSvc) contactGroupAdd(
	c ctx.Context,
	req drainageContactGroupReq.GroupSaveReq,
	members []membersModel.MemberModel,
	deptName string,
	qyWorkClient *qyWork.QyWork,
) error {
	//活动添加
	_ = s.contactGroupDao.CreateData(c, req.EnterpriseId, drainageContactModel.DrainageContactGroupModel{
		BaseModel:  models.CreateBase(req.ID),
		MemberIds:  helper.SliceToString(req.MemberIds),
		ThirdIds:   helper.SliceToString(req.ThirdIds),
		Name:       req.Name,
		ChannelID:  req.ChannelId,
		AutoFriend: req.AutoFriend,
		Type:       req.Type,
	})
	var newUserIds []string
	var memberMap = make(map[int64]membersModel.MemberModel)
	for _, member := range members {
		newUserIds = append(newUserIds, member.Userid)
		memberMap[member.ID] = member
	}
	if len(newUserIds) > 0 {
		storeContactId := app.UID()
		//如果有门店或者多个员工则创建共享码
		contactWayRes, err := qyWorkClient.AddContactWay(c, qyWork.ContactCreateRequest{
			SkipVerify: helper.RuleThree(req.AutoFriend > 0, true, false),
			Type:       2,
			Scene:      2,
			State:      helper.Int64ToString(storeContactId),
			User:       newUserIds,
		})
		if err != nil {
			return err
		}
		//活动共享码入库
		err = s.contactsDao.CreateData(c, drainageContactModel.DrainageContactsModel{
			MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, storeContactId),
			Name:           helper.RuleThree(req.Type == drainageContactModel.DrainageContactGroupStore, helper.RuleThree(deptName != "", deptName, "门店共享码")+"-"+req.Name, req.Name),
			RealUserIds:    strings.Join(newUserIds, ","),
			ContactGroupId: req.ID,
			ChannelID:      req.ChannelId,
			AutoFriend:     req.AutoFriend,
			QrCodeURL:      contactWayRes.QrCode,
			ConfigID:       contactWayRes.ConfigId,
			Type:           1,
		})
		var addMembers []drainageContactModel.DrainageContactMemberModel
		for _, memberId := range req.MemberIds {
			addMembers = append(addMembers, drainageContactModel.DrainageContactMemberModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
				ContactId:      storeContactId,
				MemberId:       memberId,
			})
		}
		for _, thirdId := range req.ThirdIds {
			addMembers = append(addMembers, drainageContactModel.DrainageContactMemberModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
				ContactId:      storeContactId,
				ThirdId:        thirdId,
			})
		}
		_ = s.contactMemberDao.CreateInBatchesData(c, addMembers)
	}
	return nil
}

func (s *MainSvc) contactSaAdd(
	c ctx.Context,
	req drainageContactGroupReq.GroupSaveReq,
	contactsGroupInfo drainageContactModel.DrainageContactGroupModel,
	existSingleContacts []drainageContactModel.DrainageContactsModel,
	members []membersModel.MemberModel,
	qyWorkClient *qyWork.QyWork,
) error {
	var nowUnix = time.Now().Unix()
	var reqMemberIds []int64 //提交的所有用户（含部门中的）
	var memberMap = make(map[int64]membersModel.MemberModel, 0)
	for _, member := range members {
		reqMemberIds = append(reqMemberIds, member.ID)
		memberMap[member.ID] = member
	}
	var newMemberIds = reqMemberIds //默认是所有用户都要创建个人码，如果是编辑则要过滤需要新增的

	// 收集企微API调用结果和待处理的DB数据
	var (
		delContactIds     []int64
		delConfigIds      []string
		updateContactIds  []int64
		updateContactMap  = make(map[int64]drainageContactModel.DrainageContactsModel) // contactId -> contact
		updateTmpUserIds  = make(map[int64]string)                                     // contactId -> userid
		addContacts       []drainageContactModel.DrainageContactsModel
		addContactMembers []drainageContactModel.DrainageContactMemberModel
		apiErr            error // 记录创建新个人码时的40098错误
	)

	if len(existSingleContacts) > 0 {
		//如果有已经存在的个人码 则处理这些个人码
		var existsMemberIds []int64
		for _, contact := range existSingleContacts {
			if contact.Members != nil && len(contact.Members) > 0 {
				existsMemberIds = append(existsMemberIds, contact.Members[0].MemberId)
			}
		}
		//需要删除的联系我员工码
		delMemberIds := helper.ArrayDiff(existsMemberIds, reqMemberIds)
		var updateMemberIds []int64
		if req.Name != contactsGroupInfo.Name || req.ChannelId != contactsGroupInfo.ChannelID || req.AutoFriend != contactsGroupInfo.AutoFriend {
			//如果信息变了，则更新已存在的个人码信息，这里收集需要更新的个人码的人
			updateMemberIds = helper.ArrayIntersect(existsMemberIds, reqMemberIds)
		}
		for _, memberContact := range existSingleContacts {
			if helper.InArray(memberContact.Members[0].MemberId, delMemberIds) {
				//这里收集需要删除的个人码
				delContactIds = append(delContactIds, memberContact.ID)
				delConfigIds = append(delConfigIds, memberContact.ConfigID)
			}
			if len(updateMemberIds) > 0 && helper.InArray(memberContact.Members[0].MemberId, updateMemberIds) {
				//这里收集需要更新名称的个人码ID
				updateContactIds = append(updateContactIds, memberContact.ID)
				updateContactMap[memberContact.ID] = memberContact
				for _, member := range memberContact.Members {
					if m, ok := memberMap[member.MemberId]; ok {
						updateTmpUserIds[memberContact.ID] = m.Userid
					}
				}
			}
		}
		//需要新增的联系我员工码
		newMemberIds = helper.ArrayDiff(reqMemberIds, existsMemberIds)
	}

	if len(newMemberIds) > 0 {
		//先调企微API创建个人码，成功后收集待入库数据
		for _, memberId := range newMemberIds {
			var contactName = helper.Int64ToString(memberId)
			if member, ok := memberMap[memberId]; ok {
				contactName = member.Name
			}
			contactId := app.UID()
			contactWay, err := qyWorkClient.AddContactWay(c, qyWork.ContactCreateRequest{
				SkipVerify: helper.RuleThree(req.AutoFriend > 0, true, false),
				Type:       2,
				Scene:      2,
				State:      helper.Int64ToString(contactId),
				User:       []string{memberMap[memberId].Userid},
			})
			if err != nil {
				log.Info(c).Msgf("生成联系我活动个人码失败,userid:%s,错误:%v", memberMap[memberId].Userid, err)
				// 如果是未实名/未激活错误，记录并中断返回
				if e, ok := err.(*errorx.Error); ok && (e.Code == 40098 || e.Code == 41054) {
					apiErr = err
					break
				}
				continue
			}
			addContacts = append(addContacts, drainageContactModel.DrainageContactsModel{
				MicroBaseModel:       models.CreateMicroBase(req.EnterpriseId, contactId),
				Name:                 fmt.Sprintf("%s-%s", contactName, req.Name),
				ContactGroupId:       req.ID,
				QrCodeURL:            contactWay.QrCode,
				ChannelID:            req.ChannelId,
				AutoFriend:           req.AutoFriend,
				ConfigID:             contactWay.ConfigId,
				IsContactGroupSingle: drainageContactModel.DrainageContactGroupSingle,
				Type:                 1,
			})

			addContactMembers = append(addContactMembers, drainageContactModel.DrainageContactMemberModel{
				MicroBaseModel: models.CreateMicroBase(req.EnterpriseId, app.UID()),
				ContactId:      contactId,
				MemberId:       memberId,
			})
		}
	}

	if apiErr != nil {
		return apiErr
	}

	// DB写操作
	if len(delContactIds) > 0 {
		_ = s.contactsDao.UpdateData(c, map[string]any{
			"id": delContactIds,
		}, map[string]interface{}{
			"deleted_at": nowUnix,
		})
		_ = s.contactMemberDao.UpdateData(c, map[string]any{
			"contact_id": delContactIds,
		}, map[string]interface{}{
			"deleted_at": nowUnix,
		})
	}

	if len(updateContactIds) > 0 {
		//如果名字变了，则更新个人码名称
		for _, contactId := range updateContactIds {
			if contact, ok := updateContactMap[contactId]; ok {
				_ = s.contactsDao.UpdateData(c, map[string]any{
					"id": contact.ID,
				}, map[string]interface{}{
					"name":        strings.Replace(contact.Name, contactsGroupInfo.Name, req.Name, 1),
					"channel_id":  req.ChannelId,
					"auto_friend": req.AutoFriend,
					"updated_at":  nowUnix,
				})
			}
		}
	}

	if len(addContacts) > 0 {
		_ = s.contactsDao.CreateInBatchesData(c, addContacts)
	}
	if len(addContactMembers) > 0 {
		_ = s.contactMemberDao.CreateInBatchesData(c, addContactMembers)
	}

	// 异步执行企微API调用（删除和更新个人码）
	if len(delConfigIds) > 0 {
		go func() {
			for _, id := range delConfigIds {
				_, _ = qyWorkClient.DeleteContactWay(c, qyWork.ContactDeleteRequest{
					ConfigId: id,
				})
			}
		}()
	}
	if len(updateContactIds) > 0 {
		for _, contactId := range updateContactIds {
			if tmpUserId, ok := updateTmpUserIds[contactId]; ok && tmpUserId != "" {
				if contact, ok2 := updateContactMap[contactId]; ok2 {
					contactCopy := contact
					go func() {
						_, _ = qyWorkClient.UpdateContactWay(c, qyWork.ContactUpdateRequest{
							ConfigId:   contactCopy.ConfigID,
							SkipVerify: helper.RuleThree(req.AutoFriend > 0, true, false),
							User:       []string{tmpUserId},
							State:      helper.Int64ToString(contactCopy.ID),
						})
					}()
				}
			}
		}
	}

	return nil
}

// checkMembersReal 逐个检测员工是否未实名认证或未激活
// 当 AddContactWay 返回 40098(未实名) 或 41054(未激活) 错误时调用此方法，逐个员工创建临时二维码来定位具体问题员工
// 并发控制最大10个goroutine
func (s *MainSvc) checkMembersReal(c ctx.Context, enterpriseId int64, userIds []string) []string {
	if len(userIds) == 0 {
		return nil
	}
	// 如果只有一个员工，直接返回
	if len(userIds) == 1 {
		return userIds
	}

	var (
		unverifiedUserIds []string
		mu                sync.Mutex
		wg                sync.WaitGroup
		sem               = make(chan struct{}, 10) // 并发控制最大10
	)

	for _, uid := range userIds {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量
		go func(userId string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			qyResp, err := qyWork.NewQyWork(enterpriseId).AddContactWay(c, qyWork.ContactCreateRequest{
				SkipVerify: true,
				Type:       1,
				Scene:      2,
				State:      "check_real",
				User:       []string{userId},
			})

			if err != nil {
				// 判断是否是40098未实名或41054未激活错误（未激活一定未实名）
				if e, ok := err.(*errorx.Error); ok && (e.Code == 40098 || e.Code == 41054) {
					mu.Lock()
					unverifiedUserIds = append(unverifiedUserIds, userId)
					mu.Unlock()
				}
				return
			}

			// 创建成功说明该员工已实名，异步删除临时二维码
			if qyResp.ConfigId != "" {
				go func() {
					_, _ = qyWork.NewQyWork(enterpriseId).DeleteContactWay(c, qyWork.ContactDeleteRequest{
						ConfigId: qyResp.ConfigId,
					})
				}()
			}
		}(uid)
	}
	wg.Wait()

	return unverifiedUserIds
}

// handleUnverifiedError 处理企微返回的未实名/未激活错误（40098/41054）
// 如果是40098或41054错误，调用checkMembersReal定位具体问题员工，返回未实名员工信息和error
// 其他错误原样返回
func (s *MainSvc) handleUnverifiedError(c ctx.Context, enterpriseId int64, err error, members []membersModel.MemberModel) ([]drainageContactResp.UnverifiedUser, error) {
	if e, ok := err.(*errorx.Error); ok && (e.Code == 40098 || e.Code == 41054) {
		// 收集所有员工的userid
		var userIds []string
		useridToMember := make(map[string]membersModel.MemberModel)
		for _, member := range members {
			userIds = append(userIds, member.Userid)
			useridToMember[member.Userid] = member
		}
		// 逐个检测定位未实名/未激活员工
		unverifiedUserIds := s.checkMembersReal(c, enterpriseId, userIds)
		// 构建(userid+姓名)列表
		var unverifiedUsers []drainageContactResp.UnverifiedUser
		for _, uid := range unverifiedUserIds {
			name := uid
			if m, ok := useridToMember[uid]; ok {
				name = m.Name
			}
			unverifiedUsers = append(unverifiedUsers, drainageContactResp.UnverifiedUser{
				Userid: uid,
				Name:   name,
			})
		}
		return unverifiedUsers, errorx.NewError(40098, fmt.Errorf("以下员工未实名认证: %s", strings.Join(unverifiedUserIds, ",")))
	}
	return nil, err
}

func (s *MainSvc) welcomeSave(
	c ctx.Context,
	req drainageContactGroupReq.GroupSaveReq,
) error {
	//删除老的欢迎语
	if err := s.welcomeDao.UpdateData(c, req.EnterpriseId, map[string]any{
		"contact_group_id": req.ID,
	}, map[string]any{
		"deleted_at": time.Now().Unix(),
	}); err != nil {
		return err
	}
	if len(req.WelcomeList) > 0 {
		var welcomeModelList = make([]drainageContactModel.DrainageContactsWelcomesModel, 0)
		for _, info := range req.WelcomeList {
			contentStr, _ := json.Marshal(info.Content)
			welcomeModelList = append(welcomeModelList, drainageContactModel.DrainageContactsWelcomesModel{
				BaseModel:      models.CreateBase(app.UID()),
				ContactGroupId: req.ID,
				Content:        string(contentStr),
				Type:           info.Type,
			})
		}
		if err := s.welcomeDao.CreateInBatchesData(c, req.EnterpriseId, welcomeModelList); err != nil {
			return err
		}
	}
	return nil
}
