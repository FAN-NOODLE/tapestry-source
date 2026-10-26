package display

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/http/request/display"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type AllotDao struct {
	dao       dao.Dao
	model     *microScrmModel.DisplayTaskAllotModel
	tableName string
}

func NewAllotDao() *AllotDao {
	return &AllotDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.DisplayTaskAllotModel),
		tableName: microScrmModel.DisplayTaskAllotModel{}.TableName(),
	}
}

func (d *AllotDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.DisplayTaskAllotModel, error) {
	var result microScrmModel.DisplayTaskAllotModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *AllotDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.DisplayTaskAllotModel, error) {
	var result []microScrmModel.DisplayTaskAllotModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

// GetImageList 获取门店下VM陈列的图片
func (d *AllotDao) GetImageList(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.DisplayTaskAllotDataModel, error) {
	var result []microScrmModel.DisplayTaskAllotDataModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).
		Where(query, args...).Preload("Remark", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId))
	}).Find(&result).Error
	return result, err
}

func (d *AllotDao) GetListPagination(c ctx.Context, request display.MemberTaskListReq, storeId []int) ([]display.MemberTaskListRes, error) {
	enterpriseId := request.EnterpriseId
	result := make([]display.MemberTaskListRes, 0)
	whereMap := map[string]interface{}{
		"display_task_allot.store_id":   storeId,
		"display_task_allot.deleted_at": 0,
		"a.show_type":                   microScrmModel.ShowStatusUp,
	}
	if request.FinishStatus >= 0 {
		whereMap["display_task_allot.finish_status"] = request.FinishStatus
	}
	taskTable := d.dao.GetTable(microScrmModel.DisplayTaskModel{}.TableName(), request.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN " + taskTable + " as a ON a.id=display_task_allot.display_task_id").
		Where(whereMap)
	if request.Name != "" {
		query.Where("display_task_allot.display_name LIKE ?", "%"+request.Name+"%")
	}
	offset := (request.Page - 1) * request.Size
	selectStr := "display_task_allot.id as allot_id,display_task_allot.finish_status,display_task_allot.store_id,display_task_allot.grade,a.id,a.display_name,a.start_time,a.end_time,a.status,a.content"
	err := query.Select(selectStr).Order("a.start_time desc").Offset(offset).Limit(request.Size).Scan(&result).Error
	return result, err
}

// CreateTaskByStoreIds 创建任务，根据任务创建有关门店的任务内容
func (d *AllotDao) CreateTaskByStoreIds(c ctx.Context, enterpriseId, taskId int64, displayName string, storeInfo []display.StoreBaseInfo) error {
	allotLists := make([]microScrmModel.DisplayTaskAllotModel, 0)
	nowTime := time.Now().Unix()
	for _, info := range storeInfo {
		value := microScrmModel.DisplayTaskAllotModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: nowTime,
				UpdatedAt: nowTime,
			},
			DisplayName:   displayName,
			DisplayTaskID: taskId,
			StoreCode:     info.StoreCode,
			StoreID:       info.StoreId,
			StoreName:     info.StoreName,
			StoreType:     int32(info.StoreType),
			StoreCity:     info.StoreCity,
			SpecialUserId: info.SpecialUserId,
		}
		allotLists = append(allotLists, value)
	}
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(allotLists).Error
}

func (d *AllotDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) error {
	var model microScrmModel.DisplayTaskAllotModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Delete(&model).Error
	return err
}

func (d *AllotDao) DeleteImageList(c ctx.Context, enterpriseId int64, ids []int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"id": ids,
	}).Delete(&microScrmModel.DisplayTaskAllotDataModel{}).Error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"allot_data_id": ids,
	}).Delete(&microScrmModel.DisplayTaskAllotRemarkModel{}).Error
	return err
}

// SaveAllotData VM陈列任务门店数据保存
func (d *AllotDao) SaveAllotData(c ctx.Context, enterpriseId, taskId, allotId int64, dataMap map[int64]display.ImageInfo, oldIds []int64) error {
	var err error
	oldAllotData := make([]display.ImageInfo, 0)
	saveImageData := make([]microScrmModel.DisplayTaskAllotDataModel, 0)
	saveCommentMap := make(map[int64][]display.CommentInfo, 0)
	saveComment := make([]microScrmModel.DisplayTaskAllotRemarkModel, 0)
	nowTime := time.Now().Unix()
	for k, info := range dataMap {
		if helper.InArray(k, oldIds) {
			oldAllotData = append(oldAllotData, info)
		} else {
			imageData := microScrmModel.DisplayTaskAllotDataModel{
				BaseModel: models.BaseModel{
					ID:        k,
					CreatedAt: nowTime,
					UpdatedAt: nowTime,
				},
				DisplayTaskID: taskId,
				AllotID:       allotId,
				Identifier:    info.Identifier,
				ImageURL:      info.Url,
				Status:        0,
			}
			saveImageData = append(saveImageData, imageData)
			saveCommentMap[k] = info.Comment
		}
	}
	if len(saveCommentMap) > 0 {
		err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).Create(saveImageData).Error
		for k, item := range saveCommentMap {
			for _, info := range item {
				if info.ID == 0 {
					commentData := microScrmModel.DisplayTaskAllotRemarkModel{
						BaseModel: models.BaseModel{
							ID:        app.UID(),
							CreatedAt: nowTime,
							UpdatedAt: nowTime,
						},
						DisplayTaskID: taskId,
						AllotID:       allotId,
						AllotDataID:   k,
						MemberID:      info.MemberId,
						MemberUserID:  info.MemberUserId,
						MemberName:    info.MemberName,
						Avatar:        info.Avatar,
						CommentTime:   info.CommentTime,
						IsVm:          info.IsVm,
						StoreCode:     info.StoreCode,
						Content:       helper.MarshalToString(info.Content),
						ReplyName:     info.ReplyName,
					}
					saveComment = append(saveComment, commentData)
				}
			}
		}
		if len(saveComment) > 0 {
			err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Create(saveComment).Error
		}
	}
	err = d.UpdateCommentInfo(c, enterpriseId, taskId, allotId, oldAllotData)
	return err
}

// UpdateCommentInfo 图片的备注内容更新
func (d *AllotDao) UpdateCommentInfo(c ctx.Context, enterpriseId, taskId, allotId int64, oldAllotData []display.ImageInfo) error {
	// 数据更新
	saveComment := make([]microScrmModel.DisplayTaskAllotRemarkModel, 0)
	nowTime := time.Now().Unix()
	commentDelId := make([]int64, 0)
	commentEditMap := make(map[int64]string, 0)
	var err error
	if len(oldAllotData) > 0 {
		for _, info := range oldAllotData {
			for _, value := range info.Comment {
				if value.Action == microScrmModel.RemarkDel {
					commentDelId = append(commentDelId, value.ID)
				}
				if value.Action == microScrmModel.RemarkEdit {
					commentEditMap[value.ID] = helper.MarshalToString(value.Content)
				}
				if value.ID == 0 {
					commentData := microScrmModel.DisplayTaskAllotRemarkModel{
						BaseModel: models.BaseModel{
							ID:        app.UID(),
							CreatedAt: nowTime,
							UpdatedAt: nowTime,
						},
						DisplayTaskID: taskId,
						AllotID:       allotId,
						AllotDataID:   info.ID,
						MemberID:      value.MemberId,
						MemberUserID:  value.MemberUserId,
						MemberName:    value.MemberName,
						Avatar:        value.Avatar,
						Content:       helper.MarshalToString(value.Content),
						StoreCode:     value.StoreCode,
						CommentTime:   value.CommentTime,
						IsVm:          value.IsVm,
						ReplyName:     value.ReplyName,
					}
					saveComment = append(saveComment, commentData)
				}
			}
		}
	}
	if len(saveComment) > 0 {
		err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Create(saveComment).Error
	}
	if len(commentEditMap) > 0 {
		for k, info := range commentEditMap {
			err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).
				Where("id=?", k).Updates(map[string]interface{}{
				"content": info,
			}).Error
		}
	}
	if len(commentDelId) > 0 {
		err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
			"id": commentDelId,
		}).Delete(&microScrmModel.DisplayTaskAllotRemarkModel{}).Error
	}
	return err
}

// SaveAllotSubmitDataLog VM陈列任务门店数据保存日志
func (d *AllotDao) SaveAllotSubmitDataLog(c ctx.Context, enterpriseId int64, submitDataLog microScrmModel.DisplayTaskAllotSubmitModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotSubmitModel{}.TableName(), enterpriseId)).Create(&submitDataLog).Error

	return err
}

// SaveAllotReviewDataLog VM审核门店任务日志
func (d *AllotDao) SaveAllotReviewDataLog(c ctx.Context, enterpriseId int64, submitDataLog microScrmModel.DisplayTaskAllotReviewModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotReviewModel{}.TableName(), enterpriseId)).Create(&submitDataLog).Error

	return err
}

// UpdateData 更新数据
func (d *AllotDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

// CreateLookLog 导购查看图片的VM备注
func (d *AllotDao) CreateLookLog(c ctx.Context, enterpriseId int64, lookLog microScrmModel.DisplayTaskLookLogModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskLookLogModel{}.TableName(), enterpriseId)).Create(&lookLog).Error
}

// GetLooKLog 获取导购查看的图片id
func (d *AllotDao) GetLooKLog(c ctx.Context, enterpriseId, allotId, memberId int64) ([]display.LookRes, error) {
	var result []display.LookRes
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskLookLogModel{}.TableName(), enterpriseId)).
		Where(map[string]interface{}{
			"allot_id":  allotId,
			"member_id": memberId,
		}).Select("allot_data_id,comment_id").Order("comment_id asc").Scan(&result).Error
	return result, err
}

// AuthInfo 获取VM的身份信息
func (d *AllotDao) AuthInfo(c ctx.Context, enterpriseId int64, memberId int64) (microScrmModel.DisplayManagerModel, error) {
	var info microScrmModel.DisplayManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayManagerModel{}.TableName(), enterpriseId)).Where("member_id=?", memberId).Find(&info).Error
	return info, err
}

// VmList 获取任务下VM导购id
func (d *AllotDao) VmList(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var res []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayManagerModel{}.TableName(), enterpriseId)).
		Where(query, args).Pluck("member_id", &res).Error
	return res, err
}

func (d *AllotDao) GetCount(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Count(&total).Error
	return total, err
}

func (d *AllotDao) GetCheckList(c ctx.Context, req display.MemberTaskListReq, isSpecial bool) ([]display.MemberTaskListRes, error) {
	enterpriseId := req.EnterpriseId
	result := make([]display.MemberTaskListRes, 0)
	whereMap := map[string]interface{}{}
	whereMap["a.show_type"] = microScrmModel.ShowStatusUp
	if req.FinishStatus >= 0 {
		whereMap["display_task_allot.finish_status"] = req.FinishStatus
		whereMap["display_task_allot.deleted_at"] = 0
	}
	if isSpecial {
		whereMap["display_task_allot.special_user_id"] = req.MemberUserId
	} else {
		whereMap["display_task_allot.owner_member_id"] = req.MemberId
	}
	taskTable := d.dao.GetTable(microScrmModel.DisplayTaskModel{}.TableName(), req.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN " + taskTable + " as a ON a.id=display_task_allot.display_task_id").
		Where(whereMap)
	if req.Name != "" {
		query.Where("display_task_allot.display_name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status >= 0 {
		query.Where("a.status=?", req.Status)
	}
	if isSpecial {
		query.Where("display_task_allot.special_check>0")
	} else {
		query.Where("(display_task_allot.special_check=2 or display_task_allot.special_user_id='')")
	}
	offset := (req.Page - 1) * req.Size
	selectStr := "display_task_allot.id as allot_id,display_task_allot.finish_status,display_task_allot.store_id,display_task_allot.grade,a.id,a.display_name,a.start_time,a.end_time,a.status,a.content,a.class_url"
	err := query.Select(selectStr).Group("display_task_id").Order("a.start_time desc").Offset(offset).Limit(req.Size).Scan(&result).Error
	return result, err
}

func (d *AllotDao) GetCheckTaskCount(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]display.ReviewStatus, error) {
	var result []display.ReviewStatus
	taskTable := d.dao.GetTable(microScrmModel.DisplayTaskModel{}.TableName(), enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("display_task_id,finish_status,count(1) as num").
		Joins("LEFT JOIN "+taskTable+" as a ON a.id=display_task_allot.display_task_id").
		Where("a.show_type=? and display_task_allot.deleted_at=0", microScrmModel.ShowStatusUp).
		Where(query, args...).Group("display_task_id,finish_status").Scan(&result).Error
	return result, err
}

func (d *AllotDao) GetVmStoreInfo(c ctx.Context, req display.MemberStoreTaskInfoReq, storeIds []int) (microScrmModel.DisplayTaskAllotModel, error) {
	var result microScrmModel.DisplayTaskAllotModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Where(map[string]interface{}{
			"display_task_id": req.ID,
			"store_id":        storeIds,
			"deleted_at":      0,
		}).Preload("ImageList", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), req.EnterpriseId)).Where(map[string]interface{}{
			"deleted_at": 0,
		}).Order("created_at asc")
	}).Preload("ImageList.Remark", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), req.EnterpriseId)).Where(map[string]interface{}{
			"deleted_at": 0,
		}).Order("comment_time asc")
	}).Find(&result).Error
	return result, err
}

func (d *AllotDao) UpdateImageData(c ctx.Context, enterpriseId int64, imageIds []int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).
		Where(map[string]interface{}{
			"id": imageIds,
		}).UpdateColumns(map[string]interface{}{"status": microScrmModel.VmImageFail}).Error
}

// GetAllotList 管理后台获取门店任务列表接口
func (d *AllotDao) GetAllotList(c ctx.Context, req display.SearchAllotReq) (int64, []microScrmModel.DisplayTaskAllotModel, error) {
	var result []microScrmModel.DisplayTaskAllotModel
	var total int64
	whereMap := make(map[string]interface{}, 0)
	whereMap["deleted_at"] = 0
	if req.OwnerMemberID > 0 {
		whereMap["owner_member_id"] = req.OwnerMemberID
	}
	if req.FinishStatus >= 0 {
		whereMap["finish_status"] = req.FinishStatus
	}
	if req.DisplayTaskId > 0 {
		whereMap["display_task_id"] = req.DisplayTaskId
	}
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where(whereMap)
	if req.DisplayName != "" {
		query.Where(" display_name like ? ", "%"+req.DisplayName+"%")
	}
	if req.StoreName != "" {
		query.Where(" (store_name like ?  or store_code like ? ) ", "%"+req.StoreName+"%", "%"+req.StoreName+"%")
	}
	if req.SubmitStartTime > 0 && req.SubmitEndTime > 0 {
		query.Where(" submit_time between ? and ? ", req.SubmitStartTime, req.SubmitEndTime)
	}
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).Preload("Submit", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotSubmitModel{}.TableName(), req.EnterpriseId)).Order("id,submit_type desc")
	}).Preload("Review", func(db *gorm.DB) *gorm.DB {
		return db.Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotReviewModel{}.TableName(), req.EnterpriseId)).Order("id desc")
	}).Order("submit_time desc").Offset(offset).Limit(req.Size).Find(&result).Error
	return total, result, err
}

// GetAllotListCount 获取VM任务列表的统计数
func (d *AllotDao) GetAllotListCount(c ctx.Context, req display.SearchAllotBaseReq) ([]display.AllotListCount, error) {
	var result []display.AllotListCount
	whereMap := make(map[string]interface{}, 0)
	whereMap["deleted_at"] = 0
	if req.OwnerMemberID > 0 {
		whereMap["owner_member_id"] = req.OwnerMemberID
	}
	if req.FinishStatus >= 0 {
		whereMap["finish_status"] = req.FinishStatus
	}
	if req.DisplayTaskId > 0 {
		whereMap["display_task_id"] = req.DisplayTaskId
	}
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).Where(whereMap)
	if req.DisplayName != "" {
		query.Where("display_name like ? ", "%"+req.DisplayName+"%")
	}
	if req.StoreName != "" {
		query.Where(" (store_name like ?  or store_code like ? ) ", "%"+req.StoreName+"%", "%"+req.StoreName+"%")
	}
	if req.SubmitStartTime > 0 && req.SubmitEndTime > 0 {
		query.Where("submit_time between ? and ? ", req.SubmitStartTime, req.SubmitEndTime)
	}
	err := query.Select("finish_status,count(1) as num").Group("finish_status").Scan(&result).Error
	return result, err
}

func (d *AllotDao) GetAllotSendNoticeList(c ctx.Context, req display.SearchAllotBaseReq) ([]display.AllotSendNoticeInfo, error) {
	var result []display.AllotSendNoticeInfo
	whereMap := make(map[string]interface{}, 0)
	whereMap = map[string]interface{}{
		"display_task_allot.deleted_at":    0,
		"a.show_type":                      microScrmModel.ShowStatusUp,
		"a.status":                         microScrmModel.DisplayStatusIng,
		"display_task_allot.finish_status": []int{microScrmModel.VMFinishStatusIng, microScrmModel.VMFinishStatusFail},
	}
	if req.FinishStatus >= 0 {
		whereMap["display_task_allot.finish_status"] = req.FinishStatus
	}
	taskTable := d.dao.GetTable(microScrmModel.DisplayTaskModel{}.TableName(), req.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Joins("LEFT JOIN " + taskTable + " as a ON a.id=display_task_allot.display_task_id").
		Where(whereMap)
	if req.DisplayName != "" {
		query.Where("display_task_allot.display_name LIKE ?", "%"+req.DisplayName+"%")
	}
	if req.StoreName != "" {
		query.Where("display_task_allot.store_name like ? ", "%"+req.StoreName+"%")
	}
	if req.OwnerMemberID > 0 {
		whereMap["display_task_allot.owner_member_id"] = req.OwnerMemberID
	}
	if req.DisplayTaskId > 0 {
		whereMap["display_task_allot.display_task_id"] = req.DisplayTaskId
	}
	if req.SubmitStartTime > 0 && req.SubmitEndTime > 0 {
		query.Where("display_task_allot.submit_time between ? and ? ", req.SubmitStartTime, req.SubmitEndTime)
	}
	selectStr := "display_task_allot.display_task_id,display_task_allot.store_id,display_task_allot.finish_status"
	err := query.Select(selectStr).Order("a.start_time desc").Scan(&result).Error

	return result, err
}

func (d *AllotDao) AllotImageList(c ctx.Context, req display.SaveImageReq) ([]string, error) {
	result := make([]string, 0)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), req.EnterpriseId)).
		Where(map[string]interface{}{
			"allot_id": req.AllotId,
			"status":   microScrmModel.VmImageSuccess,
		}).Pluck("image_url", &result).Error
	return result, err
}

func (d *AllotDao) GetAllotNoticeNum(c ctx.Context, req display.SearchAllotBaseReq) (int64, error) {
	var total int64
	whereMap := make(map[string]interface{}, 0)
	whereMap = map[string]interface{}{
		"display_task_allot.deleted_at":    0,
		"a.show_type":                      microScrmModel.ShowStatusUp,
		"a.status":                         microScrmModel.DisplayStatusIng,
		"display_task_allot.finish_status": []int{microScrmModel.VMFinishStatusIng, microScrmModel.VMFinishStatusFail},
	}
	if req.FinishStatus >= 0 {
		whereMap["display_task_allot.finish_status"] = req.FinishStatus
	}
	taskTable := d.dao.GetTable(microScrmModel.DisplayTaskModel{}.TableName(), req.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, req.EnterpriseId)).
		Joins("LEFT JOIN " + taskTable + " as a ON a.id=display_task_allot.display_task_id").
		Where(whereMap)
	if req.DisplayName != "" {
		query.Where("display_task_allot.display_name LIKE ?", "%"+req.DisplayName+"%")
	}
	if req.StoreName != "" {
		query.Where("display_task_allot.store_name like ? ", "%"+req.StoreName+"%")
	}
	if req.SubmitStartTime > 0 && req.SubmitEndTime > 0 {
		query.Where("display_task_allot.submit_time between ? and ? ", req.SubmitStartTime, req.SubmitEndTime)
	}
	if req.OwnerMemberID > 0 {
		whereMap["display_task_allot.owner_member_id"] = req.OwnerMemberID
	}
	if req.DisplayTaskId > 0 {
		whereMap["display_task_allot.display_task_id"] = req.DisplayTaskId
	}

	err := query.Count(&total).Error

	return total, err
}

func (d *AllotDao) SaveAllotDataInfo(c ctx.Context, enterpriseId int64, data microScrmModel.DisplayTaskAllotDataModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).Create(&data).Error
}

func (d *AllotDao) SaveAllotComment(c ctx.Context, enterpriseId int64, data microScrmModel.DisplayTaskAllotRemarkModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Create(&data).Error
}

func (d *AllotDao) DelAllotComment(c ctx.Context, enterpriseId int64, id int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"id": id,
	}).Delete(&microScrmModel.DisplayTaskAllotRemarkModel{}).Error
}

func (d *AllotDao) UpdateAllotCommentInfo(c ctx.Context, enterpriseId, id int64, info string) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotRemarkModel{}.TableName(), enterpriseId)).
		Where("id=?", id).Updates(map[string]interface{}{
		"content": info,
	}).Error
}

func (d *AllotDao) GetStoreCount(c ctx.Context, enterpriseId int64, status int, taskIds []int64) ([]display.StoreCount, error) {
	var result []display.StoreCount
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("store_id,count(1) as num").
		Where(map[string]interface{}{
			"display_task_id": taskIds,
			"finish_status":   status,
			"deleted_at":      0,
		}).Group("store_id").Scan(&result).Error
	return result, err
}

func (d *AllotDao) GetVmDashboardCount(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]display.StoreDashBoard, error) {
	var result []display.StoreDashBoard
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("store_id,finish_status,count(1) as num").
		Where(query, args...).Group("store_id,finish_status").Scan(&result).Error
	return result, err
}

func (d *AllotDao) GetSubmitInfo(c ctx.Context, enterpriseId, id int64) (microScrmModel.DisplayTaskAllotSubmitModel, int64, error) {
	var result microScrmModel.DisplayTaskAllotSubmitModel
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotSubmitModel{}.TableName(), enterpriseId)).
		Where(map[string]interface{}{
			"allot_id":    id,
			"submit_type": microScrmModel.VMSubmitCheck,
		}).Order("created_at desc").First(&result).Error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotSubmitModel{}.TableName(), enterpriseId)).
		Where(map[string]interface{}{
			"allot_id":    id,
			"submit_type": microScrmModel.VMSubmitCheck,
		}).Count(&total).Error
	return result, total, err
}

// GetManagerList 获取vm的配置人员表的列表
func (d *AllotDao) GetManagerList(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.DisplayManagerModel, error) {
	var info []microScrmModel.DisplayManagerModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayManagerModel{}.TableName(), enterpriseId)).Where(query, args...).Find(&info).Error
	return info, err
}

// SpecialUpdateImageData special审核数据的保存
func (d *AllotDao) SpecialUpdateImageData(c ctx.Context, enterpriseId int64, imageIds []int64, updateData interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotDataModel{}.TableName(), enterpriseId)).
		Where(map[string]interface{}{
			"id": imageIds,
		}).UpdateColumns(updateData).Error
}

// GetAllotReviewDataLog VM审核门店任务日志
func (d *AllotDao) GetAllotReviewDataLog(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.DisplayTaskAllotReviewModel, error) {
	var data []microScrmModel.DisplayTaskAllotReviewModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.DisplayTaskAllotReviewModel{}.TableName(), enterpriseId)).
		Where(query, args...).Order("created_at desc").Group("review_user_id").Find(&data).Error

	return data, err
}
