package tag

import (
	"archive/zip"
	"bytes"
	"errors"
	"math"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"time"

	"github.com/tealeg/xlsx"

	"micro-scrm/internal/events"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	EventData "micro-scrm/internal/events/entity"
)

// 客户分组服务

type UserGroupService struct {
	Repo UserGroupRepository
}

// 根据客户分组ID，删除客户分组
func (u *UserGroupService) DeleteCustomerGroup(c ctx.Context, enterpriseID, customerGroupID int64) error {
	if err := orm.Service.UserGroup.DelCustomerGroupByID(enterpriseID, customerGroupID); err != nil {
		return err
	}

	go events.ScrmEvent.DeleteCustomerGroup.Dispatch(EventData.Data{
		Ctx:             c,
		EnterpriseID:    enterpriseID,
		CustomerGroupID: customerGroupID,
	})

	return nil
}

// 客户分组预计算人数
func (u *UserGroupService) PreCalculateCount(c ctx.Context, request entity.UserGroupPrecalculateRequest, enterpriseID int64) (int, error) {
	includeList := make([]int64, 0)
	filterList := make([]int64, 0)
	for _, condition := range request.Conditions {
		// 包含标签
		if 1 == condition.IncludeType {
			for _, tagGroup := range condition.TagGroups {
				includeOne, includeErr := u.Repo.GetExternalUserIDByTagIDList(c, tagGroup.TagIds, enterpriseID)
				if includeErr != nil {
					return 0, includeErr
				}
				if len(includeOne) == 0 {
					break
				}
				if len(includeList) > 0 {
					includeList = u.intersection(includeList, includeOne)
					if len(includeList) == 0 {
						break
					}
				} else {
					includeList = includeOne
				}
			}
		}
		// 过滤标签
		if 2 == condition.IncludeType {
			for _, tagGroup := range condition.TagGroups {
				filterOne, filterErr := u.Repo.GetExternalUserIDByTagIDList(c, tagGroup.TagIds, enterpriseID)
				if filterErr != nil {
					return 0, filterErr
				}
				if len(filterOne) == 0 {
					break
				}
				if len(filterList) > 0 {
					filterList = u.intersection(filterList, filterOne)
					if len(filterList) == 0 {
						break
					}
				} else {
					filterList = filterOne
				}
			}
		}
	}

	return len(u.aNotB(includeList, filterList)), nil
}

func (u *UserGroupService) intersection(a, b []int64) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return make([]int64, 0)
	}
	var tmp = make(map[int64]int, len(a))
	for _, v := range a {
		tmp[v] = 1
	}
	for _, arg := range b {
		if tmp[arg] != 0 {
			tmp[arg]++
		}
	}
	var res = make([]int64, 0)
	for k, v := range tmp {
		if v > 1 {
			res = append(res, k)
		}
	}
	return res
}

// 客户分组详情
func (u *UserGroupService) UserGroupDetail(c ctx.Context, enterpriseID, userGroupID int64) (entity.UserGroupDetail, error) {
	result := entity.UserGroupDetail{} // 初始化返回值

	tagInfoList, listErr := u.Repo.GetUserGroupDetail(c, entity.QueryUserGroupParam{
		EnterpriseID: enterpriseID,
		UserGroupID:  userGroupID,
	})
	if listErr != nil {
		return result, listErr
	}
	var DepartmentIDs []int64
	for _, userGroup := range tagInfoList {
		result.ID = userGroup.UserGroupID
		result.Name = userGroup.UserGroupName
		result.IsExtract = userGroup.IsExtract
		result.FileID = userGroup.FileID
		result.UserGroupType = userGroup.UserGroupType
		if userGroup.DepartmentID != 0 {
			DepartmentIDs = append(DepartmentIDs, userGroup.DepartmentID)
		}
	}

	oneList := make(map[int64][]entity.UserGroupDetailTagInfo, 0)
	twoList := make(map[int64][]entity.UserGroupDetailTagInfo, 0)
	for _, userGroup := range tagInfoList {
		if 1 == userGroup.IncludeType {
			oneList[userGroup.TagGroupID] = append(oneList[userGroup.TagGroupID], userGroup)
		} else if 2 == userGroup.IncludeType {
			twoList[userGroup.TagGroupID] = append(twoList[userGroup.TagGroupID], userGroup)
		}
	}

	result.Conditions = append(result.Conditions, &entity.UserGroupDetailCondition{
		IncludeType: 1,
		TagGroups:   u.changeUserGroupDetailTagGroup(oneList),
	})
	result.Conditions = append(result.Conditions, &entity.UserGroupDetailCondition{
		IncludeType: 2,
		TagGroups:   u.changeUserGroupDetailTagGroup(twoList),
	})
	var DepartmentInfo []entity.Departments
	if len(DepartmentIDs) > 0 {
		var Department entity.Departments
		Departments, err := orm.Service.Departments.GetDepartmentInfoBYIds(enterpriseID, DepartmentIDs)
		if err != nil {
			return result, err
		}
		for _, department := range Departments {
			Department.ID = department.ID
			Department.Name = department.Name
			DepartmentInfo = append(DepartmentInfo, Department)
		}
	}
	result.Departments = DepartmentInfo

	// 判断该客户分组是否可以编辑、删除；
	// 条件为：没有被除了已结束的营销策略使用，则可以修改
	canEdit, editErr := orm.Service.Marketing.IsEditByUserGroupID(userGroupID)
	if editErr != nil {
		return result, editErr
	}

	result.CanEdit = canEdit

	return result, nil
}

func (u *UserGroupService) changeUserGroupDetailTagGroup(a map[int64][]entity.UserGroupDetailTagInfo) []*entity.UserGroupDetailTagGroup {
	dataList := make([]*entity.UserGroupDetailTagGroup, 0)
	for _, item := range a {
		isHas := false
		data := entity.UserGroupDetailTagGroup{}
		for _, tag := range item {
			if !isHas {
				data.TagGroupID = tag.TagGroupID
				data.TagGroupName = tag.TagGroupName
				isHas = true
			}
			data.Tags = append(data.Tags, &entity.UserGroupDetailTag{
				TagID:   tag.TagID,
				TagName: tag.TagName,
			})
		}

		dataList = append(dataList, &data)
	}

	return dataList
}

// 获取客户分组列表
func (u *UserGroupService) GetUserGroupList(c ctx.Context, enterpriseID int64) ([]entity.UserGroupInfo, error) {
	dataList, err := u.Repo.GetUserGroupList(c, enterpriseID)
	if err != nil {
		return nil, err
	}

	return dataList, nil
}

// A - (A ∩ B)
func (u *UserGroupService) aNotB(a, b []int64) []int64 {
	if 0 == len(a) || 0 == len(b) {
		return a
	}

	bb := make(map[int64]int8, 0)
	for _, v := range b {
		bb[v] = 0
	}

	aa := make([]int64, len(a))
	for _, v := range a {
		if _, ok := bb[v]; !ok {
			aa = append(aa, v)
		}
	}

	return aa
}

// 新增客户分组
// 分为两步：1、创建用户分组；2、建立用户分组和用户分组的关系
func (u *UserGroupService) AddUserGroup(c ctx.Context, enterpriseID int64, body entity.CreateUserGroupBody) (int64, error) {
	// 校验用户分组名称是否已经存在
	isExist, existErr := u.Repo.IsExist(c, body.Name)
	if existErr != nil {
		return 0, existErr
	}

	if isExist {
		return 0, errors.New("该客户分组名称已经被使用")
	}

	// 1、创建用户分组
	userGroupData := entity.UserGroupModelEntity{
		ID:           uid.ID(c),
		EnterpriseID: enterpriseID,
		Name:         body.Name,
		Type:         entity.UserGroupTypeSelf,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	_, addUserGroupErr := u.Repo.AddUserGroup(c, userGroupData)
	if addUserGroupErr != nil {
		return 0, addUserGroupErr
	}

	if body.UserGroupType == 2 && body.FileID == 0 {
		return 0, errors.New("客户导入的文件未上传")
	}
	if body.UserGroupType == 3 && len(body.DepartmentIDs) < 1 {
		return 0, errors.New("请选择门店")
	}
	var dataList []entity.UserGroupDetailModelEntity

	if body.UserGroupType == 0 {
		body.UserGroupType = 1
	}
	if body.UserGroupType == 1 {
		// 2、建立用户分组与标签的关系
		// 构建批量添加对象
		for _, item := range body.Conditions {
			for _, tagGroup := range item.TagGroups {
				for _, tagID := range tagGroup.TagIds {
					dataList = append(dataList, entity.UserGroupDetailModelEntity{
						ID:            uid.ID(c),
						EnterpriseID:  enterpriseID,
						UserGroupID:   userGroupData.ID,
						TagGroupID:    tagGroup.TagGroupID,
						TagID:         int64(tagID),
						IncludeType:   int8(item.IncludeType),
						CreatedAt:     time.Now().Unix(),
						UpdatedAt:     time.Now().Unix(),
						UserGroupType: 1,
						FileID:        int64(0),
						DepartmentID:  int64(0),
					})
				}
			}
		}
	}

	if body.UserGroupType == 2 {
		dataList = append(dataList, entity.UserGroupDetailModelEntity{
			ID:            uid.ID(c),
			EnterpriseID:  enterpriseID,
			UserGroupID:   userGroupData.ID,
			TagGroupID:    int64(0),
			TagID:         int64(0),
			IncludeType:   int8(0),
			CreatedAt:     time.Now().Unix(),
			UpdatedAt:     time.Now().Unix(),
			UserGroupType: body.UserGroupType,
			FileID:        body.FileID,
			DepartmentID:  int64(0),
		})
	}

	if body.UserGroupType == 3 {
		for _, DepartmentID := range body.DepartmentIDs {
			dataList = append(dataList, entity.UserGroupDetailModelEntity{
				ID:            uid.ID(c),
				EnterpriseID:  enterpriseID,
				UserGroupID:   userGroupData.ID,
				TagGroupID:    int64(0),
				TagID:         int64(0),
				IncludeType:   int8(0),
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
				UserGroupType: body.UserGroupType,
				FileID:        0,
				DepartmentID:  DepartmentID,
			})
		}
	}

	_, addUserGroupDetailErr := u.Repo.AddUserGroupDetail(c, dataList)
	if addUserGroupDetailErr != nil {
		return 0, addUserGroupDetailErr
	}

	return userGroupData.ID, nil
}

func (u *UserGroupService) SyncGroupImportData(c ctx.Context, enterpriseID, fileID, GroupID int64) error {
	fileInfo, errs := orm.Service.UserGroup.GetImportFile(enterpriseID, fileID)
	if errs != nil {
		return errs
	}
	if fileInfo.IsExtract == 1 {
		return errs
	}
	b, err := helper.ReadRemoteFile(c, fileInfo.FileURL)
	if err != nil {
		return err
	}
	var data []*entity.UnionIDData
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return err
	}
	data = readXlsxData(fileData.Sheets)
	limit := 1000
	dataLen := len(data)
	if dataLen > limit {
		dataPage := int(math.Ceil(float64(dataLen / limit)))
		for i := 1; i < dataPage; i++ {
			start := (i - 1) * limit
			end := i * limit
			dataList := data[start:end]
			if end >= dataLen {
				dataList = data[start:dataLen]
			}
			if err := CreateImportData(enterpriseID, GroupID, dataList, fileInfo); err != nil {
				return err
			}
		}
	} else {
		if err := CreateImportData(enterpriseID, GroupID, data, fileInfo); err != nil {
			return err
		}
	}
	if err = orm.Service.UserGroup.UpdateISExtract(enterpriseID, fileID); err != nil {
		return err
	}
	return nil

}

func readXlsxData(sheets []*xlsx.Sheet) []*entity.UnionIDData {
	var settings []*entity.UnionIDData
	for index, sheet := range sheets {
		if index > 0 {
			continue
		}
		for ri, row := range sheet.Rows {
			if ri < 1 {
				continue
			}
			if len(row.Cells) == 0 {
				continue
			}
			var setting = &entity.UnionIDData{}
			if row.Cells[0].String() != "" {
				setting.UnionID = row.Cells[0].String()
				settings = append(settings, setting)
			}
		}
	}
	return settings
}

// CreateImportData 导入的unionid入库
func CreateImportData(enterpriseID, GroupID int64, data []*entity.UnionIDData, fileInfo micro_scrm_enterprise.UserGroupFileModel) error {
	nowTime := time.Now().Unix()
	var createData []micro_scrm_enterprise.UserGroupImportDataModel
	for _, item := range data {
		var newData micro_scrm_enterprise.UserGroupImportDataModel
		newData.CreatedAt = nowTime
		newData.UpdatedAt = nowTime
		newData.FileID = fileInfo.ID
		newData.FileName = fileInfo.FileName
		newData.UnionID = item.UnionID
		createData = append(createData, newData)
	}
	if err := orm.Service.UserGroup.CreateCSVInfo(enterpriseID, createData); err != nil {
		return err
	}
	return nil

}

func (u *UserGroupService) EditCustomerGroup(c ctx.Context, enterpriseID int64, body entity.CreateUserGroupBody) (int64, error) {
	// 2、建立用户分组与标签的关系
	err := u.Repo.UpdateUserGroup(c, enterpriseID, body.ID, map[string]interface{}{"name": body.Name})
	if err != nil {
		return 0, err
	}
	// 构建批量添加对象
	var dataList []entity.UserGroupDetailModelEntity
	if body.UserGroupType == 0 {
		body.UserGroupType = 1
	}
	if body.UserGroupType == 1 {
		for _, item := range body.Conditions {
			for _, tagGroup := range item.TagGroups {
				for _, tagID := range tagGroup.TagIds {
					dataList = append(dataList, entity.UserGroupDetailModelEntity{
						ID:           uid.ID(c),
						EnterpriseID: enterpriseID,
						UserGroupID:  body.ID,
						TagGroupID:   tagGroup.TagGroupID,
						TagID:        int64(tagID),
						IncludeType:  int8(item.IncludeType),
						CreatedAt:    time.Now().Unix(),
						UpdatedAt:    time.Now().Unix(),
					})
				}
			}
		}
	}
	if body.UserGroupType == 3 {
		for _, DepartmentID := range body.DepartmentIDs {
			dataList = append(dataList, entity.UserGroupDetailModelEntity{
				ID:            uid.ID(c),
				EnterpriseID:  enterpriseID,
				UserGroupID:   body.ID,
				TagGroupID:    int64(0),
				TagID:         int64(0),
				IncludeType:   int8(0),
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
				UserGroupType: body.UserGroupType,
				FileID:        0,
				DepartmentID:  DepartmentID,
			})
		}
	}

	_, addUserGroupDetailErr := u.Repo.AddUserGroupDetail(c, dataList)
	if addUserGroupDetailErr != nil {
		return 0, addUserGroupDetailErr
	}
	return body.ID, nil
}

func (u *UserGroupService) syncEnterpriseUserGroups(c ctx.Context, enterpriseID, userGroupID int64) {
	// 1、根据客户分组ID，获取所有和该分组相关的用户ID集合
	groupDetail, err := u.UserGroupDetail(c, enterpriseID, userGroupID)
	if err != nil {
		log.Info(c).Msgf("syncEnterpriseUserGroup同步用户客户分组信息::获取客户分组详情失败,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID, userGroupID)
		return
	}

	// 2、获取客户分组对应的客户信息
	includeList := make([]int64, 0)
	filterList := make([]int64, 0)
	for _, condition := range groupDetail.Conditions {
		// 包含标签
		if 1 == condition.IncludeType {
			for _, tagGroup := range condition.TagGroups {
				tagIDList := []int64{}
				for _, tagInfo := range tagGroup.Tags {
					tagIDList = append(tagIDList, tagInfo.TagID)
				}
				includeOne, includeErr := u.Repo.GetExternalUserIDByTagIDList(c, tagIDList, enterpriseID)
				if includeErr != nil {
					log.Info(c).Msgf("syncEnterpriseUserGroup同步用户客户分组信息::删选客户分组用户时失败,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID, userGroupID)
					return
				}

				includeList = u.intersection(includeList, includeOne)
			}
		}
		// 过滤标签
		if 2 == condition.IncludeType {
			for _, tagGroup := range condition.TagGroups {
				tagIDList := []int64{}
				for _, tagInfo := range tagGroup.Tags {
					tagIDList = append(tagIDList, tagInfo.TagID)
				}
				filterOne, filterErr := u.Repo.GetExternalUserIDByTagIDList(c, tagIDList, enterpriseID)
				if filterErr != nil {
					log.Info(c).Msgf("syncEnterpriseUserGroup同步用户客户分组信息::删选客户分组用户时失败,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID, userGroupID)
					return
				}

				filterList = u.intersection(filterList, filterOne)
			}
		}
	}
	externalUserIDList := u.aNotB(includeList, filterList)

	// 3、获取外部联系人的聚合ID enterprise_user.id
	useIDList, userErr := orm.Service.WorkExternalUser.GetUserIDByExternalUserIDList(enterpriseID, externalUserIDList)
	if userErr != nil {
		log.Info(c).Msgf("syncEnterpriseUserGroup同步用户客户分组信息::获取外部联系人的聚合ID集合失败,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID, userGroupID)
		return
	}

	// 更新字段enterprise_user.groups
	updateErr := orm.Service.WorkExternalUser.UpdateExternalUserUserGroup(enterpriseID, userGroupID, useIDList)
	if updateErr != nil {
		log.Info(c).Msgf("syncEnterpriseUserGroup同步用户客户分组信息::更新外部联系人的groups集合失败,err:%s,enterpriseID:%d,userGroupID:%d", err, enterpriseID, userGroupID)
		return
	}
}

// GetFileNameDetection 根据文件名验证是否已经导入过
func (u *UserGroupService) GetFileNameDetection(enterpriseID int64, fileName string) (result int64, err error) {
	return orm.Service.UserGroup.GetFileNameDetection(enterpriseID, fileName)
}

// SaveUploadInfo 保存客户导入模板里的数据信息
func (u *UserGroupService) SaveUploadInfo(data entity.DownLoadFileInfo, enterpriseID int64) (ID int64, err error) {
	var newInfo micro_scrm_enterprise.UserGroupFileModel
	newInfo.CreatedAt = time.Now().Unix()
	newInfo.UpdatedAt = time.Now().Unix()
	newInfo.FileName = data.FileName
	newInfo.FileType = 1
	newInfo.FileURL = data.FileUrl

	id, err := orm.Service.UserGroup.CreateUploadInfo(enterpriseID, newInfo)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetUpdateInfo 更新数据状态
func (u *UserGroupService) GetUpdateInfo(enterpriseID, ID int64) (err error) {
	return orm.Service.UserGroup.UpdateISExtract(enterpriseID, ID)
}

func (u *UserGroupService) GetImportData(enterpriseID, FileID int64, Page, PageSize int) (data []micro_scrm_enterprise.UserGroupImportDataModel, err error) {
	return orm.Service.UserGroup.GetImportData(enterpriseID, FileID, Page, PageSize)
}

func (u *UserGroupService) GetImportDataCount(enterpriseID, FileID int64) (total int64, err error) {
	return orm.Service.UserGroup.GetImportDataCount(enterpriseID, FileID)
}

func (u *UserGroupService) GetUserGroupLogCount(enterpriseID, GroupID int64) (total int64, err error) {
	return orm.Service.UserGroup.GetUserGroupLogCount(enterpriseID, GroupID)
}

func (u *UserGroupService) GetUserGroupLogList(enterpriseID, GroupID int64, Page, PageSize int) (data []entity.ExportDataInfo, err error) {
	return orm.Service.UserGroup.GetUserGroupLogList(enterpriseID, GroupID, Page, PageSize)
}
