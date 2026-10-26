package source

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceModel "micro-scrm/internal/app/models/source"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/errs"
	"time"
)

type CategoryService struct {
	categoryDao         *sourceDao.CategoryDao
	CategoryRelationDao *sourceDao.CategoryRelationDao
}

func NewCategoryService() CategoryService {
	return CategoryService{
		categoryDao:         sourceDao.NewCategoryDao(),
		CategoryRelationDao: sourceDao.NewCategoryRelationDao(),
	}
}

func (s CategoryService) GetList(c app.Ctx, request sourceRequest.CategoryListRequest) (error, []sourceModel.CategoryModel, int64) {
	categoryIds := make([]int64, 0)
	sourceIds := make([]int64, 0)
	if request.IsCollect > 0 {
		collectSourceId, err := sourceDao.NewCollectRelationDao().GetCollectAllSourceIds(c, request.EnterpriseID, request.MemberId)
		if err != nil {
			return err, nil, 0
		}
		if !request.OriginH5 {
			collectFriendSourceIds, err := sourceDao.NewCollectRelationDao().GetFriendSourceId(c, request.EnterpriseID, request.MemberId)
			if err != nil {
				return err, nil, 0
			}
			collectSourceId = append(collectSourceId, collectFriendSourceIds...)
		}
		// 所属门店下 分组的sourceGroupIDs
		var whereGroupRelation = map[string]interface{}{
			"third_id": request.ThirdIds,
		}
		sourceGroupIds, err := sourceDao.NewGroupRelationDao().GetManySourceGroupIds(c, request.EnterpriseID, whereGroupRelation)
		if err != nil {
			return err, nil, 0
		}
		sourceIds, err = sourceDao.NewSourceDao().GetExpireSourceIds(c, request.EnterpriseID, collectSourceId, sourceGroupIds)
		if err != nil {
			return err, nil, 0
		}
		if len(sourceIds) == 0 {
			return nil, make([]sourceModel.CategoryModel, 0), 0
		}
		categoryIds, err = sourceDao.NewCategoryRelationDao().GetGetSourceIdsByCategoryIdsBySourceId(c, request.EnterpriseID, sourceIds)
		if err != nil || len(categoryIds) == 0 {
			return err, nil, 0
		}
	}

	return s.categoryDao.GetList(c, request.Name, request.Page, request.Size, request.EnterpriseID, categoryIds)
}

func (s CategoryService) Create(c app.Ctx, request sourceRequest.CategoryCreateRequest) error {
	// 查询是否存在同名
	err, result := s.categoryDao.GetOneByName(c, request.Name, request.EnterpriseID)
	if err != nil {
		return err
	}

	if result.ID != 0 {
		return errs.DataExist
	}

	data := sourceModel.CategoryModel{
		ID:           uint64(app.UIDContext(c)),
		Name:         request.Name,
		EnterpriseID: uint64(request.EnterpriseID),
		CreatedAt:    uint(time.Now().Unix()),
		UpdatedAt:    uint(time.Now().Unix()),
	}

	err = s.categoryDao.Create(c, data)
	if err != nil {
		return err
	}

	return nil
}

func (s CategoryService) Update(c app.Ctx, request sourceRequest.CategoryUpdateRequest) error {
	err, result := s.categoryDao.GetOneByName(c, request.Name, request.EnterpriseID)
	if err != nil {
		return err
	}
	if request.Name == result.Name && uint64(request.Id) != result.ID {
		return errs.DataExist
	}
	query := map[string]interface{}{
		"id":            request.Id,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{"name": request.Name}
	err = s.categoryDao.Update(c, query, column)
	if err != nil {
		return err
	}
	return nil
}

func (s CategoryService) SourceRelation(c app.Ctx, request sourceRequest.CategoryRequest) ([]int64, error) {
	return s.CategoryRelationDao.GetSourceIdsByCategoryId(c, request.Id)
}

func (s CategoryService) Delete(c app.Ctx, request sourceRequest.CategoryRequest) error {
	err, result := s.categoryDao.GetOneById(c, request.Id, request.EnterpriseID)
	if err != nil {
		return err
	}
	if result.ID == 0 {
		return errs.SourceNotFund
	}

	return gormx.GetInstance().Transaction(func(tx *gorm.DB) error {
		if err = s.categoryDao.Delete(c, map[string]interface{}{
			"id": request.Id,
		}); err != nil {
			return err
		}

		if err = s.CategoryRelationDao.ForceDelete(c, map[string]interface{}{
			"category_id": result.ID,
		}); err != nil {
			return err
		}
		return err
	})
}

func (s CategoryService) RemoveCategory(c app.Ctx, request sourceRequest.RemoveCategoryRequest) error {
	return s.CategoryRelationDao.ForceDelete(c, map[string]interface{}{
		"enterprise_id": request.EnterpriseID,
		"category_id":   request.CategoryId,
		"source_id":     request.SourceID,
	})
}
