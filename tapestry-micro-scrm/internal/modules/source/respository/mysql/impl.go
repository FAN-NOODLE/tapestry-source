package mysql

import (
	"database/sql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/app"
	sourceSvc "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/modules/source"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/errorIgnore"
	helper2 "micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
	"time"
)

type SourceRepo struct {
	*repos.SQLRepo
}

// TextSave 保存文本素材
func (r *SourceRepo) TextSave(cc ctx.Context, data source.RequestTextSave) (id int64, err error) {

	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			cond := r.Cond2Args(data)
			// 主表数据
			if data.UpdatedAt == 0 {
				data.UpdatedAt = time.Now().Unix()
				stmt, err := r.DB.Template("t_source/source-insert", cond)
				if err != nil {
					return err
				}
				res, err := stmt.NamedExec(cc, data)
				if err != nil || res == nil {
					return err
				}
			}

			timeNow := time.Now().Unix()
			// 添加素材类别和素材关系
			for _, categoryId := range data.CategoryIds {
				query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, relationInsertErr := tx.Exec(cc, query, uid.ID(cc), data.EnterpriseId, categoryId, data.ID, source.TypeText, timeNow, timeNow)
				if relationInsertErr != nil {
					return relationInsertErr
				}
			}

			// 文本素材数据
			stmtSource, err := r.DB.Template("t_source/source-text-insert", cond)
			if err != nil {
				return err
			}
			resSource, err := stmtSource.NamedExec(cc, data)
			if err != nil || resSource == nil {
				return err
			}

			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return
}

// LinkSave 保存链接素材
func (r *SourceRepo) LinkSave(cc ctx.Context, data source.RequestLinkSave) (id int64, err error) {

	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			cond := r.Cond2Args(data)

			// 主表数据
			if data.UpdatedAt == 0 {
				data.UpdatedAt = time.Now().Unix()
				stmt, err := r.DB.Template("t_source/source-insert", cond)
				if err != nil {
					return err
				}
				res, err := stmt.NamedExec(cc, data)
				if err != nil || res == nil {
					return err
				}
			}

			timeNow := time.Now().Unix()
			// 添加素材类别和素材关系
			for _, categoryId := range data.CategoryIds {
				query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, relationInsertErr := tx.Exec(cc, query, uid.ID(cc), data.EnterpriseId, categoryId, data.ID, source.TypeLink, timeNow, timeNow)
				if relationInsertErr != nil {
					return relationInsertErr
				}
			}

			// 链接素材数据
			stmtSource, err := r.DB.Template("t_source/source-link-insert", cond)
			if err != nil {
				return err
			}
			resSource, err := stmtSource.NamedExec(cc, data)

			if err != nil || resSource == nil {
				return err
			}

			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return
}

// TextDetail 文本素材详情
func (r *SourceRepo) TextDetail(cc ctx.Context, data *source.RequestSourceDetail) (*source.ResponseTextDetail, error) {
	cond := r.Cond2Args(data)
	stmt, err := r.DB.Template("t_source/source-text-detail", cond)
	if err != nil {
		return nil, err
	}
	res := &source.ResponseTextDetail{}

	err = stmt.Get(cc, res)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	return res, nil
}

// LinkDetail 链接素材详情
func (r *SourceRepo) LinkDetail(cc ctx.Context, data *source.RequestSourceDetail) (*source.ResponseLinkDetail, error) {
	cond := r.Cond2Args(data)
	stmt, err := r.DB.Template("t_source/source-link-detail", cond)
	if err != nil {
		return nil, err
	}
	res := &source.ResponseLinkDetail{}

	err = stmt.Get(cc, res)
	if errorIgnore.Ignore(err) != nil {
		return nil, err
	}

	return res, nil
}

// AddSourcesImage 新增图片是素材
func (r *SourceRepo) AddSourcesImage(cc ctx.Context, data *source.SourcesImageRequest) (id int64, err error) {
	// 判断是否为组合展示,组合展示需要新增合并id
	var mergeId int64
	if data.ShowType == 2 {
		mergeId = app.UID()
	}

	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			if !helper.IsNil(data.Images) && len(data.Images) > 0 {
				timeNow := assist.Timestamp()
				for _, item := range data.Images {
					// 插入主表
					insertSource := &source.Source{
						ID:           item.ID,
						EnterpriseID: item.EnterpriseID,
						GroupID:      data.GroupID,
						Type:         source.TypeImg,
						AdminID:      data.AdminID,
						MemberID:     data.MemberID,
						StartTime:    data.StartTime,
						EndTime:      data.EndTime,
						CreatedAt:    timeNow,
						UpdatedAt:    timeNow,
						DeletedAt:    0,
						MergeID:      mergeId,
						Sort:         helper2.RuleThree(data.ShowType == 1, 0, item.Sort),
					}
					stmt, err := tx.Template("t_source/source-insert")
					if err != nil {
						return err
					}
					_, insertErr := stmt.NamedExec(cc, insertSource)
					if insertErr != nil {
						return insertErr
					}

					// 添加素材类别和素材关系
					for _, categoryId := range data.CategoryIds {
						query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
						_, relationInsertErr := tx.Exec(cc, query, app.UID(), data.EnterpriseID, categoryId, item.ID, source.TypeImg, timeNow, timeNow)
						if relationInsertErr != nil {
							return relationInsertErr
						}
					}

					// 插入Img
					insertImgSource := &source.SourcesImage{
						ID:           item.ID,
						EnterpriseID: item.EnterpriseID,
						Name:         item.Name,
						URL:          item.URL,
						Ext:          item.Ext,
						Size:         item.Size,
						CreatedAt:    timeNow,
						UpdatedAt:    timeNow,
						DeletedAt:    0,
					}
					stmtImg, errImg := tx.Template("t_source/source-image-insert")
					if errImg != nil {
						return err
					}
					_, insertImgErr := stmtImg.NamedExec(cc, insertImgSource)
					if insertImgErr != nil {
						return insertImgErr
					}
					// 更新关联商品
					_ = sourceSvc.NewImageService().UpdateRelationGoods(cc, uint64(item.ID), uint64(data.EnterpriseID), item.RelationGoods, source.TypeImg)
				}
			}
			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})
	return
}

// AddSourcesFile 新增文件是素材
func (r *SourceRepo) AddSourcesFile(cc ctx.Context, data *source.SourcesFileRequest) (id int64, err error) {

	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {

			timeNow := assist.Timestamp()
			// 插入主表
			insertSource := &source.Source{
				ID:           data.ID,
				EnterpriseID: data.EnterpriseID,
				GroupID:      data.GroupID,
				Type:         source.TypeFile,
				AdminID:      data.AdminID,
				MemberID:     data.MemberID,
				StartTime:    data.StartTime,
				EndTime:      data.EndTime,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
				DeletedAt:    0,
			}
			stmt, err := tx.Template("t_source/source-insert")
			if err != nil {
				return err
			}
			_, insertErr := stmt.NamedExec(cc, insertSource)
			if insertErr != nil {
				return insertErr
			}

			// 添加素材类别和素材关系
			for _, categoryId := range data.CategoryIds {
				query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, relationInsertErr := tx.Exec(cc, query, uid.ID(cc), data.EnterpriseID, categoryId, data.ID, source.TypeFile, timeNow, timeNow)
				if relationInsertErr != nil {
					return relationInsertErr
				}
			}

			// 插入文件
			insertImgSource := &source.SourcesFile{
				ID:           data.ID,
				EnterpriseID: data.EnterpriseID,
				Name:         data.Name,
				URL:          data.URL,
				AttachName:   data.AttachName,
				AttachID:     data.AttachID,
				FileType:     data.FileType,
				Ext:          data.Ext,
				Size:         data.Size,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
				DeletedAt:    0,
			}
			stmtFile, errFile := tx.Template("t_source/source-file-insert")
			if errFile != nil {
				return err
			}
			_, insertFileErr := stmtFile.NamedExec(cc, insertImgSource)
			if insertFileErr != nil {
				return insertFileErr
			}

			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})
	return
}

// SourcesImageList 图片列表
func (r *SourceRepo) SourcesImageList(c ctx.Context, data *source.RequestList) (int64, *[]source.ResponseImgList, error) {

	offset := (data.Page - 1) * data.Size
	pageVo := repos.PageVO{Offset: offset, Limit: data.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-image-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]source.ResponseImgList{}
		},
	}
	total, list, err := r.FindPage(c, pageVo, data, option)
	if err != nil && err != sql.ErrNoRows {
		return 0, nil, err
	}
	resList := list.(*[]source.ResponseImgList)
	return total, resList, nil
}

// SourcesFileList 文件列表
func (r *SourceRepo) SourcesFileList(c ctx.Context, data *source.RequestList) (int64, *[]source.ResponseFileList, error) {

	offset := (data.Page - 1) * data.Size
	pageVo := repos.PageVO{Offset: offset, Limit: data.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-file-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]source.ResponseFileList{}
		},
	}
	total, list, err := r.FindPage(c, pageVo, data, option)
	if err != nil {
		return 0, nil, err
	}
	resList := list.(*[]source.ResponseFileList)
	return total, resList, nil
}

// SourceMiniProgramAdd 小程序素材--新增
func (r *SourceRepo) SourceMiniProgramAdd(cc ctx.Context, data *source.RequestSourceMini) (id int64, err error) {
	timeNow := assist.Timestamp()
	id = int64(assist.IdGenerate(cc))
	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			// 插入主表
			insertSource := &source.Source{
				ID:           id,
				EnterpriseID: data.EnterpriseID,
				GroupID:      data.GroupID,
				Type:         source.TypeMiniProgram,
				AdminID:      data.AdminID,
				MemberID:     data.MemberID,
				StartTime:    data.StartTime,
				EndTime:      data.EndTime,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			}
			stmt, err := tx.Template("t_source/source-insert")
			if err != nil {
				return err
			}
			_, insertErr := stmt.NamedExec(cc, insertSource)
			if insertErr != nil {
				return insertErr
			}

			// 添加素材类别和素材关系
			for _, categoryId := range data.CategoryIds {
				query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, relationInsertErr := tx.Exec(cc, query, uid.ID(cc), data.EnterpriseID, categoryId, data.ID, source.TypeMiniProgram, timeNow, timeNow)
				if relationInsertErr != nil {
					return relationInsertErr
				}
			}

			// 插入素材小程序子表
			insertSub := &source.SourceMini{
				ID:            id,
				EnterpriseID:  data.EnterpriseID,
				Title:         data.Title,
				MiniProgramID: data.MiniProgramID,
				Path:          data.Path,
				CoverUrl:      data.CoverUrl,
				CreatedAt:     timeNow,
				UpdatedAt:     timeNow,
			}
			stmtMini, errMini := tx.Template("t_source/source-mini-insert")
			if errMini != nil {
				return errMini
			}
			_, insertSubErr := stmtMini.NamedExec(cc, insertSub)
			if insertSubErr != nil {
				return insertSubErr
			}

			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return
}

// SourceMiniProgramUpdate 小程序素材--更新
func (r *SourceRepo) SourceMiniProgramUpdate(cc ctx.Context, data *source.RequestSourceMini) (int64, error) {
	data.UpdatedAt = time.Now().Unix()
	//list, _ := source.GetChildIDList(data.GroupID, data.EnterpriseID)
	//if len(list) > 1 {
	//	return 0, errors.New("该节点下有分组，无法添加素材")
	//}
	cond := r.Cond2Args(data)
	stmt, err := r.DB.Template("t_source/source-mini-update", cond)
	if err != nil {
		return 0, err
	}
	_, err = stmt.NamedExec(cc, data)
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

// SourceVideoAdd 视频素材--新增
func (r *SourceRepo) SourceVideoAdd(cc ctx.Context, data *source.RequestSourceVideo) (id int64, err error) {
	timeNow := assist.Timestamp()
	id = int64(assist.IdGenerate(cc))
	err = r.DB.Tx(cc, func(tx *sqlplus.Tx) error {
		// 防止 crash
		e, crash := fx.RunSafe(func() error {
			// 插入主表
			insertSource := &source.Source{
				ID:           id,
				EnterpriseID: data.EnterpriseID,
				GroupID:      data.GroupID,
				Type:         source.TypeVideo,
				AdminID:      data.AdminID,
				MemberID:     data.MemberID,
				StartTime:    data.StartTime,
				EndTime:      data.EndTime,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			}
			stmt, err := tx.Template("t_source/source-insert")
			if err != nil {
				return err
			}
			_, insertErr := stmt.NamedExec(cc, insertSource)
			if insertErr != nil {
				return insertErr
			}

			// 添加素材类别和素材关系
			for _, categoryId := range data.CategoryIds {
				query := "INSERT INTO source_category_relation (`id`, `enterprise_id`, `category_id`, `source_id`, `type`, `created_at`, `updated_at`) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, relationInsertErr := tx.Exec(cc, query, uid.ID(cc), data.EnterpriseID, categoryId, id, source.TypeVideo, timeNow, timeNow)
				if relationInsertErr != nil {
					return relationInsertErr
				}
			}

			// 插入视频子表
			insertSub := &source.SourceVideo{
				ID:           id,
				EnterpriseID: data.EnterpriseID,
				Name:         data.Name,
				Url:          data.Url,
				CoverUrl:     data.CoverUrl,
				Size:         data.Size,
				Ext:          data.Ext,
				CreatedAt:    timeNow,
				UpdatedAt:    timeNow,
			}
			stmtMini, errMini := tx.Template("t_source/source-video-insert")
			if errMini != nil {
				return errMini
			}
			_, insertSubErr := stmtMini.NamedExec(cc, insertSub)
			if insertSubErr != nil {
				return insertSubErr
			}

			return err
		})

		if crash != nil {
			return crash
		}

		return e
	})

	return
}

// SourceVideoUpdate 视频素材--更新
func (r *SourceRepo) SourceVideoUpdate(cc ctx.Context, data *source.RequestSourceVideo) (int64, error) {
	data.UpdatedAt = time.Now().Unix()
	//list, _ := source.GetChildIDList(data.GroupID, data.EnterpriseID)
	//if len(list) > 1 {
	//	return 0, errors.New("该节点下有分组，无法添加素材")
	//}
	cond := r.Cond2Args(data)
	stmt, err := r.DB.Template("t_source/source-video-update", cond)
	if err != nil {
		return 0, err
	}
	_, err = stmt.NamedExec(cc, data)
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

// SourceMiniProgramList 小程序素材--列表
func (r *SourceRepo) SourceMiniProgramList(cc ctx.Context, search *source.RequestSourceSearch) (int64, *[]source.SourceMini, error) {
	pageVo := assist.PageVO(search.Page, search.Size)
	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-mini-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]source.SourceMini{}
		},
	}

	total, list, err := r.FindPage(cc, pageVo, search, option)
	if err != nil {
		return 0, nil, err
	}

	return total, list.(*[]source.SourceMini), nil
}

// SourceVideoList 视频--列表
func (r *SourceRepo) SourceVideoList(cc ctx.Context, search *source.RequestSourceSearch) (int64, *[]source.SourceVideo, error) {
	pageVo := assist.PageVO(search.Page, search.Size)
	option := &repos.FindPageOption{
		TplPrefix: "t_source/source-video-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]source.SourceVideo{}
		},
	}

	total, list, err := r.FindPage(cc, pageVo, search, option)
	if err != nil {
		return 0, nil, err
	}

	return total, list.(*[]source.SourceVideo), nil
}

// SourceDelete 删除素材
func (r *SourceRepo) SourceDelete(c ctx.Context, data *source.RequestSourceDelete) (int64, error) {

	cond := &map[string]interface{}{
		"IDs": data.IDs,
	}

	stmt, err := r.DB.Template("t_source/source-delete", cond)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(c)
	if err != nil {
		return 0, err
	}

	deletedAt := time.Now().Unix()
	updates := map[string]interface{}{
		"deleted_at": deletedAt,
		"updated_at": deletedAt,
	}
	for _, tableName := range []string{
		"source_text",
		"source_image",
		"source_link",
		"source_mini_program",
		"source_video",
		"source_file",
		"source_friend",
	} {
		if err := app.DB().WithContext(c).Table(tableName).
			Where("id IN ?", data.IDs).
			Updates(updates).Error; err != nil {
			return 0, err
		}
	}

	return res.RowsAffected()
}

// MoveGroup 移动分组
func (r *SourceRepo) MoveGroup(c ctx.Context, data *source.RequestGroupMove) (int64, error) {
	cond := &map[string]interface{}{
		"GroupID": data.GroupID,
		"IDs":     data.SourceIDs,
	}

	stmt, err := r.DB.Template("t_source/source-move-group", cond)
	if err != nil {
		return 0, err
	}

	res, err := stmt.NamedExec(c, data)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// UpdateGroup 更细分组
func (r *SourceRepo) UpdateGroup(c ctx.Context, data *source.RequestGroupMove) (int64, error) {
	stmt, err := r.DB.Template("t_source/source-update-group")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(c, data.GroupID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpdateSourceMedia 更新素材信息
func (r *SourceRepo) UpdateSourceMedia(c ctx.Context, req source.UpdateSourceMediaRequest) error {
	mapping := make(map[string]interface{})
	mapping["source"] = req
	stmt, err := r.DB.Template("t_source/source-update-media", &mapping)
	if err != nil {
		return err
	}

	res, sErr := stmt.NamedExec(c, req)
	if sErr != nil {
		return sErr
	}

	_, rowErr := res.RowsAffected()

	if rowErr != nil {
		return rowErr
	}

	return nil
}

// InitSourceMedia 初始化素材media更新状态
func (r *SourceRepo) InitSourceMedia(c ctx.Context, enterpriseID int64, id []int64, sourceType string) error {
	mapping := make(map[string]interface{})
	mapping["sourceType"] = sourceType
	mapping["ids"] = id
	stmt, err := r.DB.Template("t_source/source-init-media-status", &mapping)
	if err != nil {
		return err
	}

	res, sErr := stmt.Exec(c, enterpriseID)
	if sErr != nil {
		return sErr
	}

	_, rowErr := res.RowsAffected()

	if rowErr != nil {
		return rowErr
	}

	return nil
}

// SourceMediaDetail 查询素材详情
func (r *SourceRepo) SourceMediaDetail(ctx ctx.Context, enterpriseID int64, sourceID int64, sourceType string) (*source.MediaItemDetail, error) {
	mapping := make(map[string]interface{})
	mapping["sourceType"] = sourceType
	stmt, err := r.DB.Template("t_source/source-sub-detail", &mapping)
	if err != nil {
		return nil, err
	}
	dest := &source.MediaItemDetail{}
	sErr := stmt.Get(ctx, dest, sourceID, enterpriseID)
	if sErr != nil && sErr != sql.ErrNoRows {
		return nil, sErr
	}

	return dest, nil
}
