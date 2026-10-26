package source

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"micro-scrm/internal/app"
	sourcedao "micro-scrm/internal/app/dao/source"
	sourcerequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models"
	sourcemodel "micro-scrm/internal/app/models/source"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"

	"gorm.io/gorm"
)

const coachKOSRootGroupName = "FY27小红书素材"

type CoachKOSSyncService struct {
	client           coachkos.Client
	pageSize         int
	ensureGroup      func(app.Ctx, int64, string, map[string]int64) (int64, error)
	getOrCreateGroup func(app.Ctx, int64, int64, string) (int64, error)
	grantGroupAccess func(app.Ctx, int64, int64) error
	getKOSByPostID   func(app.Ctx, int64, string) (sourcemodel.LinkKOSModel, error)
	getLinkByID      func(app.Ctx, int64, int64) (sourcemodel.LinkModel, error)
	resolveUploaders func(app.Ctx, int64, []coachkos.PostResult) (map[string]coachKOSUploader, error)
	uploadCover      func(app.Ctx, string) (string, error)
	cdnHost          string
	upsertPost       func(app.Ctx, int64, int64, coachkos.PostResult, string, coachKOSUploader) (CoachKOSGoodsTask, error)
	enrichGoods      func(app.Ctx, []CoachKOSGoodsTask)
	logInfo          func(app.Ctx, string, ...any)
}

func NewCoachKOSSyncService() *CoachKOSSyncService {
	config := app.Cfg().CoachKOS
	service := &CoachKOSSyncService{
		client:   coachkos.NewClient(config.BaseURL, config.Token, config.PageSize),
		pageSize: config.PageSize,
	}
	service.ensureGroup = service.ensureCoachKOSGroup
	service.getOrCreateGroup = service.getOrCreateSourceGroup
	service.grantGroupAccess = service.grantSourceGroupAccess
	service.getKOSByPostID = sourcedao.NewLinkKOSDao().GetByPostIDUnscoped
	service.getLinkByID = sourcedao.NewLinkDao().GetOneUnscoped
	service.resolveUploaders = resolveCoachKOSUploaders
	service.uploadCover = uploadCoachKOSCover
	service.cdnHost = app.Cfg().Storage.COS.CdnHost
	service.upsertPost = service.upsertCoachKOSPost
	service.enrichGoods = EnrichCoachKOSRelationGoods
	service.logInfo = func(c app.Ctx, format string, args ...any) {
		app.Logger().Info(c).Msgf(format, args...)
	}
	return service
}

func (s *CoachKOSSyncService) logSyncInfo(c app.Ctx, format string, args ...any) {
	if s.logInfo != nil {
		s.logInfo(c, format, args...)
	}
}

func (s *CoachKOSSyncService) SyncEnterprise(c app.Ctx, enterpriseID int64, updatedAt time.Time) error {
	errorMessages := make([]string, 0)
	ensuredGroups := make(map[string]int64)
	cursor := updatedAt.Format("2006-01-02 15:04:05")
	var (
		received            int64
		reportedTotal       int64
		successCount        int
		skippedInvalidCount int
		skippedEmptyTitle   int
		failedCount         int
		goodsTaskCount      int
	)
	s.logSyncInfo(c, "KOS素材同步开始：enterprise_id=%d updated_at=%s page_size=%d", enterpriseID, cursor, s.pageSize)

	for page := 1; ; page++ {
		pageTasks := make([]CoachKOSGoodsTask, 0)
		data, err := s.client.List(c, coachkos.ListRequest{UpdatedAt: cursor, Page: page, Size: s.pageSize})
		if err != nil {
			failedCount++
			errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d page=%d: %v", enterpriseID, page, err))
			break
		}
		reportedTotal = data.Total
		s.logSyncInfo(c, "KOS素材分页返回：enterprise_id=%d page=%d response_page=%d size=%d total=%d list_count=%d", enterpriseID, page, data.Page, data.Size, data.Total, len(data.List))
		nextReceived, err := validateCoachKOSPage(page, received, data)
		if err != nil {
			failedCount++
			errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d page=%d: %v", enterpriseID, page, err))
			break
		}
		uploaders := make(map[string]coachKOSUploader)
		if s.resolveUploaders != nil {
			uploaders, err = s.resolveUploaders(c, enterpriseID, data.List)
			if err != nil {
				failedCount++
				errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d page=%d 上传人查询失败: %v", enterpriseID, page, err))
				break
			}
		}

		for _, post := range data.List {
			if strings.TrimSpace(post.Name) == "" {
				skippedEmptyTitle++
				continue
			}
			if post.IsValid == "N" {
				existingKOS, err := s.getKOSByPostID(c, enterpriseID, post.PostID)
				if err != nil {
					failedCount++
					errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d post_id=%s: %v", enterpriseID, post.PostID, err))
					continue
				}
				if existingKOS.ID == 0 {
					skippedInvalidCount++
					continue
				}
			}
			groupID, err := s.ensureGroup(c, enterpriseID, post.DataDate, ensuredGroups)
			if err != nil {
				failedCount++
				errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d post_id=%s: %v", enterpriseID, post.PostID, err))
				continue
			}
			materialCoverURL := s.prepareCoachKOSCoverURL(c, enterpriseID, post)
			employeeID := strings.TrimSpace(post.EmployeeID)
			uploader, matched := uploaders[employeeID]
			if employeeID != "" && !matched {
				s.logSyncInfo(c, "KOS上传人未匹配：enterprise_id=%d post_id=%s employee_id=%s", enterpriseID, post.PostID, employeeID)
			}
			task, err := s.upsertPost(c, enterpriseID, groupID, post, materialCoverURL, uploader)
			if err != nil {
				failedCount++
				errorMessages = append(errorMessages, fmt.Sprintf("enterprise_id=%d post_id=%s: %v", enterpriseID, post.PostID, err))
				continue
			}
			if len(task.ProductCodes) > 0 {
				pageTasks = append(pageTasks, task)
			}
			successCount++
		}
		if len(pageTasks) > 0 {
			goodsTaskCount += len(pageTasks)
			s.logSyncInfo(c, "KOS商品补全开始：enterprise_id=%d page=%d goods_tasks=%d", enterpriseID, page, len(pageTasks))
			s.enrichGoods(c, pageTasks)
			s.logSyncInfo(c, "KOS商品补全结束：enterprise_id=%d page=%d goods_tasks=%d", enterpriseID, page, len(pageTasks))
		}

		received = nextReceived
		if received >= data.Total {
			break
		}
	}
	s.logSyncInfo(c, "KOS素材同步结束：enterprise_id=%d total=%d received=%d success=%d skipped_invalid=%d skipped_empty_title=%d failed=%d goods_tasks=%d", enterpriseID, reportedTotal, received, successCount, skippedInvalidCount, skippedEmptyTitle, failedCount, goodsTaskCount)

	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return nil
}

func validateCoachKOSPage(requestPage int, received int64, data coachkos.PageData) (int64, error) {
	if data.Page <= 0 {
		return received, fmt.Errorf("分页响应 page=%d 非法", data.Page)
	}
	if data.Size <= 0 {
		return received, fmt.Errorf("分页响应 size=%d 非法", data.Size)
	}
	if data.Total < 0 {
		return received, fmt.Errorf("分页响应 total=%d 非法", data.Total)
	}
	if data.Page != requestPage {
		return received, fmt.Errorf("分页响应 page=%d 与请求 page=%d 不一致", data.Page, requestPage)
	}
	listCount := int64(len(data.List))
	if listCount > int64(data.Size) {
		return received, fmt.Errorf("分页返回数量=%d 超过 size=%d", listCount, data.Size)
	}
	if received > data.Total || listCount > data.Total-received {
		return received, fmt.Errorf("分页累计接收数量会超过 total=%d：已接收=%d，本页=%d", data.Total, received, listCount)
	}
	nextReceived := received + listCount
	if nextReceived < data.Total && listCount < int64(data.Size) {
		return received, fmt.Errorf("分页返回短页无法继续可靠拉取：本页数量=%d，size=%d，累计=%d，total=%d", listCount, data.Size, nextReceived, data.Total)
	}
	return nextReceived, nil
}

func (s *CoachKOSSyncService) ensureCoachKOSGroup(c app.Ctx, enterpriseID int64, dataDate string, ensuredGroups map[string]int64) (int64, error) {
	monthGroupName, err := coachKOSMonthGroupName(dataDate)
	if err != nil {
		return 0, fmt.Errorf("解析月份分组失败: %w", err)
	}

	rootID, err := s.ensureSourceGroup(c, enterpriseID, 0, coachKOSRootGroupName, ensuredGroups)
	if err != nil {
		return 0, err
	}
	return s.ensureSourceGroup(c, enterpriseID, rootID, monthGroupName, ensuredGroups)
}

func (s *CoachKOSSyncService) ensureSourceGroup(c app.Ctx, enterpriseID, parentID int64, groupName string, ensuredGroups map[string]int64) (int64, error) {
	cacheKey := fmt.Sprintf("%d:%s", parentID, groupName)
	if groupID, ok := ensuredGroups[cacheKey]; ok {
		return groupID, nil
	}
	groupID, err := s.getOrCreateGroup(c, enterpriseID, parentID, groupName)
	if err != nil {
		return 0, err
	}
	if err := s.grantGroupAccess(c, enterpriseID, groupID); err != nil {
		return 0, err
	}
	ensuredGroups[cacheKey] = groupID
	return groupID, nil
}

func (s *CoachKOSSyncService) getOrCreateSourceGroup(c app.Ctx, enterpriseID, parentID int64, groupName string) (int64, error) {
	groupDao := sourcedao.NewSourceGroupDao()
	group, err := groupDao.GetOneByName(c, enterpriseID, parentID, groupName)
	if err != nil {
		return 0, err
	}
	if group.ID != 0 {
		return group.ID, nil
	}

	orders, err := groupDao.GetMaxOrders(c, enterpriseID, parentID)
	if err != nil {
		return 0, err
	}
	now := int(time.Now().Unix())
	group = sourcemodel.SourceGroupModel{
		ID:           app.UID(),
		EnterpriseID: enterpriseID,
		GroupName:    groupName,
		CreatedAt:    now,
		UpdatedAt:    now,
		ParentID:     parentID,
		Orders:       orders + 1,
	}
	if err := groupDao.CreateData(c, group); err != nil {
		return 0, err
	}
	return group.ID, nil
}

func (s *CoachKOSSyncService) grantSourceGroupAccess(c app.Ctx, enterpriseID, groupID int64) error {
	request := sourcerequest.GroupRelationSaveRequest{
		EnterpriseID:   enterpriseID,
		SourceGroupIds: []int64{groupID},
	}
	if err := NewGroupRelationService().GroupRelationSaveAll(c, request); err != nil {
		return err
	}
	return nil
}

func mergeExistingCoachKOSModels(
	postID string,
	existingKOS sourcemodel.LinkKOSModel,
	existingSource sourcemodel.SourceModel,
	existingLink sourcemodel.LinkModel,
	generatedSource sourcemodel.SourceModel,
	generatedLink sourcemodel.LinkModel,
	generatedKOS sourcemodel.LinkKOSModel,
) (sourcemodel.SourceModel, sourcemodel.LinkModel, sourcemodel.LinkKOSModel, error) {
	if existingKOS.SourceID == 0 || existingKOS.SourceLinkID == 0 ||
		existingKOS.SourceID != existingKOS.SourceLinkID ||
		existingSource.ID != existingKOS.SourceID || existingLink.ID != existingKOS.SourceID {
		return sourcemodel.SourceModel{}, sourcemodel.LinkModel{}, sourcemodel.LinkKOSModel{},
			fmt.Errorf("post_id=%s KOS关联异常: source_id=%d source_link_id=%d", postID, existingKOS.SourceID, existingKOS.SourceLinkID)
	}

	sourceModel := existingSource
	sourceModel.GroupID = generatedSource.GroupID
	sourceModel.Status = generatedSource.Status
	sourceModel.AdminID = generatedSource.AdminID
	if generatedSource.MemberID > 0 {
		sourceModel.MemberID = generatedSource.MemberID
	}
	sourceModel.CreatedAt = generatedSource.CreatedAt
	sourceModel.EndTime = generatedSource.EndTime
	sourceModel.UpdatedAt = generatedSource.UpdatedAt
	sourceModel.DeletedAt = models.DeletedAt(0)

	linkModel := existingLink
	linkModel.Title = generatedLink.Title
	linkModel.Description = generatedLink.Description
	linkModel.CoverURL = generatedLink.CoverURL
	linkModel.LinkType = generatedLink.LinkType
	linkModel.LinkURL = generatedLink.LinkURL
	linkModel.Author = generatedLink.Author
	linkModel.Content = generatedLink.Content
	linkModel.CreatedAt = generatedLink.CreatedAt
	linkModel.UpdatedAt = generatedLink.UpdatedAt
	linkModel.DeletedAt = models.DeletedAt(0)

	kosModel := generatedKOS
	kosModel.ID = existingKOS.ID
	kosModel.SourceID = existingKOS.SourceID
	kosModel.SourceLinkID = existingKOS.SourceLinkID
	kosModel.CreatedAt = existingKOS.CreatedAt
	kosModel.DeletedAt = models.DeletedAt(0)
	return sourceModel, linkModel, kosModel, nil
}

func buildCoachKOSGoodsTask(enterpriseID, sourceID int64, productCodes []string) CoachKOSGoodsTask {
	return CoachKOSGoodsTask{
		EnterpriseID: enterpriseID,
		SourceID:     sourceID,
		ProductCodes: append([]string(nil), productCodes...),
	}
}

func (s *CoachKOSSyncService) upsertCoachKOSPost(c app.Ctx, enterpriseID, groupID int64, post coachkos.PostResult, materialCoverURL string, uploader coachKOSUploader) (CoachKOSGoodsTask, error) {
	var task CoachKOSGoodsTask
	err := app.DB().WithContext(c).Transaction(func(tx *gorm.DB) error {
		now := time.Now().Unix()
		kosDao := sourcedao.NewLinkKOSDao().WithTx(tx)
		existingKOS, err := kosDao.GetByPostIDUnscoped(c, enterpriseID, post.PostID)
		if err != nil {
			return err
		}

		sourceID := existingKOS.SourceID
		if existingKOS.ID == 0 {
			sourceID = app.UID()
		}

		sourceModel, linkModel, kosModel := buildCoachKOSModels(enterpriseID, sourceID, groupID, now, post, materialCoverURL)
		sourceModel.MemberID = uploader.memberID
		sourceModel.AdminID = uploader.adminID
		linkModel.ID = sourceID
		linkModel.EnterpriseID = enterpriseID
		kosModel.SourceLinkID = sourceID

		sourceDao := sourcedao.NewSourceDao().WithTx(tx)
		linkDao := sourcedao.NewLinkDao().WithTx(tx)
		if existingKOS.ID == 0 {
			kosModel.ID = app.UID()
		} else {
			if existingKOS.SourceID == 0 || existingKOS.SourceLinkID == 0 {
				return fmt.Errorf("post_id=%s KOS关联异常: source_id=%d source_link_id=%d", post.PostID, existingKOS.SourceID, existingKOS.SourceLinkID)
			}

			existingSource, err := sourceDao.GetOneUnscoped(c, enterpriseID, sourceID)
			if err != nil {
				return err
			}
			existingLink, err := linkDao.GetOneUnscoped(c, enterpriseID, sourceID)
			if err != nil {
				return err
			}
			sourceModel, linkModel, kosModel, err = mergeExistingCoachKOSModels(
				post.PostID, existingKOS, existingSource, existingLink, sourceModel, linkModel, kosModel,
			)
			if err != nil {
				return err
			}
		}

		sourceModel.DeletedAt = models.DeletedAt(0)
		linkModel.DeletedAt = models.DeletedAt(0)
		kosModel.DeletedAt = models.DeletedAt(0)
		if err := sourceDao.SaveData(c, sourceModel); err != nil {
			return err
		}
		if err := linkDao.SaveData(c, linkModel); err != nil {
			return err
		}
		if err := kosDao.SaveData(c, kosModel); err != nil {
			return err
		}

		goodsDao := sourcedao.NewRelationGoodsDao()
		goodsDao.SetDb(tx)
		existingGoods, err := goodsDao.GetMany(c,
			"enterprise_id = ? AND source_id = ? AND type = ? AND deleted_at = 0",
			enterpriseID, sourceID, sourcemodel.TypeLink,
		)
		if err != nil {
			return err
		}
		productCodes := expandCoachKOSProductCodes(post.ProductID)
		goodsPlan := planCoachKOSRelationGoods(existingGoods, productCodes)
		if len(goodsPlan.DeleteIDs) > 0 {
			if err := goodsDao.ForceDelete(c,
				"enterprise_id = ? AND source_id = ? AND type = ? AND id IN ?",
				enterpriseID, sourceID, sourcemodel.TypeLink, goodsPlan.DeleteIDs,
			); err != nil {
				return err
			}
		}

		goods := buildCoachKOSRelationGoods(enterpriseID, sourceID, now, goodsPlan.CreateProductCodes)
		for index := range goods {
			goods[index].ID = app.UID()
		}
		if len(goods) > 0 {
			if err := goodsDao.CreateInBatchesData(c, goods); err != nil {
				return err
			}
		}

		task = buildCoachKOSGoodsTask(enterpriseID, sourceID, goodsPlan.EnrichProductCodes)
		return nil
	})
	if err != nil {
		return CoachKOSGoodsTask{}, err
	}
	return task, nil
}
