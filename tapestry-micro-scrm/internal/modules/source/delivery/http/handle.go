package sourceHttp

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	baozunCacheDao "micro-scrm/internal/app/dao/source"
	source2 "micro-scrm/internal/app/http/request/source"
	sourceSvc "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/pkg/helper"
	"net/http"
	"strings"
	"time"

	"micro-scrm/internal/modules/chat"

	"micro-scrm/internal/modules/customer"

	"micro-scrm/internal/modules/ws_trace"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/source"
	"micro-scrm/internal/modules/source/delivery"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
	"micro-scrm/internal/pkg/uid"
)

var s *source.Service
var enterpriseService *enterprise.Service
var wsTraceService *ws_trace.Service
var customerService *customer.Service
var groupChatService *chat.GroupChatService

// SourcesImageList 图片素材List接口
func SourcesImageList(c *gin.Context) {
	params := &source.RequestList{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseId)
	params.Keyword = helper.FilterEmoji(params.Keyword)
	total, list, err := s.SourceImageList(cc, params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
}

// SourcesFileList 文件素材List接口
func SourcesFileList(c *gin.Context) {
	params := &source.RequestList{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	params.Keyword = helper.FilterEmoji(params.Keyword)
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseId)
	total, list, err := s.SourceFileList(cc, params)
	if err != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
}

// SourceImageAdd 图片素材新增
func SourceImageAdd(c *gin.Context) {
	params := &source.SourcesImageRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	var imageName []string
	var imageUrl []string

	if len(params.Images) > 0 {
		for _, item := range params.Images {
			item.ID = app.UID()
			item.EnterpriseID = params.EnterpriseID

			imageName = append(imageName, item.Name)
			imageUrl = append(imageUrl, item.URL)
		}
	}

	_, err := s.SourceImageAdd(cc, params)
	if err != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, params.Images)
}

// SourceFileAdd 文件素材新增
func SourceFileAdd(c *gin.Context) {
	params := &source.SourcesFileRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)
	params.ID = uid.ID(cc)
	paramExt := strings.ToLower(params.Ext)
	var err error
	params.FileType = source.FileExtGroup[paramExt]
	if params.FileType == 0 {
		params.FileType = source.OtherFile
	}

	switch source.TypeFileGroup[paramExt] {
	case source.TypeImg:
		_, err = s.SourceImageAdd(cc, delivery.TransImgParams(params))
		break
	case source.TypeVideo:
		_, err = s.SourceVideoAdd(cc, delivery.TransVideoParams(params))
		break
	default:
		_, err = s.SourceFileAdd(cc, params)
		break
	}

	if err != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, params)
}

// TextSave 文本素材保存接口
func TextSave(c *gin.Context) {
	params := &source.RequestTextSave{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	id, err := s.TextSave(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	// 返回结构体
	res := &source.ResponseTextDetail{
		ID:           id,
		EnterpriseID: params.EnterpriseId,
		GroupID:      params.GroupID,
		Content:      params.Content,
		Type:         params.Type,
		AdminID:      params.AdminId,
		CreatedAt:    assist.Timestamp(),
	}

	response.Success(c, res)
	return
}

// LinkSave 链接素材保存接口
func LinkSave(c *gin.Context) {
	params := &source.RequestLinkSave{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	id, err := s.LinkSave(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	// 链接分组特殊情况处理
	paramsInfo := &source.RequestSourceDetail{
		ID: id,
	}
	// 返回详情结构体
	res, err := s.SourceLinkDetail(ctx.Wrap(c), paramsInfo)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	if res.LinkType == source.LinkUrlTypeSelf {
		corpRequest := enterprise.Corp{
			EnterpriseID: res.EnterpriseID,
		}
		corp, _ := enterpriseService.Detail(ctx.Wrap(c), corpRequest)
		corpHash := corp.Hash

		domain := config.CfgData.Service.WebH5.Domain
		linkUrl := fmt.Sprintf(source.SelfLinkUrl, domain, corpHash) + assist.StrVal(res.ID)
		res.LinkUrl = linkUrl
	}
	res.RelationGoods = params.RelationGoods

	response.Success(c, res)
	return
}

// TextList 文本素材List接口
func TextList(c *gin.Context) {
	var params = &source.RequestTextList{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseID)
	total, list, err := s.TextList(cc, params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
}

// LinkList 链接素材List接口
func LinkList(c *gin.Context) {
	params := &source.RequestLinkList{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	params.Keyword = helper.FilterEmoji(params.Keyword)
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseID)
	total, list, err := s.LinkList(cc, params)
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	if total > 0 {
		corpRequest := enterprise.Corp{
			EnterpriseID: (*list)[0].EnterpriseID,
		}
		corp, _ := enterpriseService.Detail(ctx.Wrap(c), corpRequest)
		corpHash := corp.Hash
		var sourceIds []int64
		for _, item := range *list {
			if item.LinkType == source.LinkUrlTypeSelf {
				domain := config.CfgData.Service.WebH5.Domain
				linkUrl := fmt.Sprintf(source.SelfLinkUrl, domain, corpHash) + assist.StrVal(item.ID)
				item.LinkUrl = linkUrl
			}
			sourceIds = append(sourceIds, item.ID)
		}
		var sourceSkuMap = make(map[int64][]string)
		goods, _ := baozunCacheDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id in ?", sourceIds)
		if len(goods) > 0 {
			for _, model := range goods {
				if _, ok := sourceSkuMap[model.SourceID]; !ok {
					sourceSkuMap[model.SourceID] = []string{model.SKU}
				} else {
					sourceSkuMap[model.SourceID] = append(sourceSkuMap[model.SourceID], model.SKU)
				}
			}
			for _, item := range *list {
				if _, ok := sourceSkuMap[item.ID]; ok {
					item.SkuList = sourceSkuMap[item.ID]
				}
			}
		}
	}
	response.SuccessList(c, total, params.Page, params.Size, list)
}

// TextDetail 文本素材详情接口
func TextDetail(c *gin.Context) {
	params := &source.RequestSourceDetail{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	res, err := s.SourceTextDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, res)
	return
}

// LinkDetail 链接素材详情接口
func LinkDetail(c *gin.Context) {
	params := &source.RequestSourceDetail{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	res, err := s.SourceLinkDetail(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	goods, _ := baozunCacheDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id = ?", params.ID)
	if len(goods) > 0 {
		return
	}
	for _, model := range goods {
		res.RelationGoods = append(res.RelationGoods, source2.BrandGoods{
			Name:    model.Name,
			Sku:     model.SKU,
			Upc:     model.UPC,
			Url:     model.Url,
			Color:   model.Color,
			AddType: model.AddType,
			Score:   model.Score,
		})
	}
	response.Success(c, res)
	return
}

// sourceMiniProgramAdd 小程序素材-新增
func sourceMiniProgramAdd(c *gin.Context) {
	params := &source.RequestSourceMini{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	id, err := s.SourceMiniProgramAdd(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	// 小程序详情
	detail, err := s.SourceMiniProgramDetail(ctx.Wrap(c), id)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// sourceMiniProgramUpdate 小程序素材-编辑
func sourceMiniProgramUpdate(c *gin.Context) {
	params := &source.RequestSourceMini{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	id, err := s.SourceMiniProgramUpdate(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &source.ID{ID: id})
	return
}

// sourceMiniProgramList 小程序素材列表
func sourceMiniProgramList(c *gin.Context) {
	params := &source.RequestSourceSearch{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseID)
	params.Keyword = helper.FilterEmoji(params.Keyword)
	total, list, err := s.SourceMiniProgramList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	var sourceType, wxExternalID, groupChatID int64
	var wsTraceID, wxExternalUnion string

	if params.ExternalUserid != "" {
		sourceType = ws_trace.SourceTypeOneToOne
		// 查询外部联系人详情
		user, err := customerService.FindExternal(ctx.Wrap(c), customer.WorkExternalUser{EnterpriseID: params.EnterpriseID, ExternalUserID: params.ExternalUserid})
		if err != nil {
			ginplus.ResError(c, errs.ParamsErr)
			return
		}
		wxExternalID = user.ID
		wxExternalUnion = user.UnionID
	}
	if params.GroupChatID != "" {
		sourceType = ws_trace.SourceTypeGroupChat
		// 查询群详情
		group, err := groupChatService.GetGroupInfo(ctx.Wrap(c), chat.WorkGroupChat{ChatID: params.GroupChatID, EnterpriseID: params.EnterpriseID})
		if errorIgnore.Ignore(err) != nil {
			ginplus.ResError(c, errs.ParamsErr)
		}
		groupChatID = group.ID
	}

	if sourceType > 0 {
		wst := ws_trace.WsTrace{
			MemberID:     params.MemberID,
			WxExternalID: wxExternalID,
			Unionid:      wxExternalUnion,
			GroupChatID:  groupChatID,
			SourceType:   sourceType,
		}
		wsTrace, err := wsTraceService.WsAdd(ctx.Wrap(c), params.EnterpriseID, wst)
		if err != nil {
			log.Err(ctx.Wrap(c), err).Msgf("生成ws_trace_id错误:%+v", wst)
		}
		if wsTrace != nil {
			wsTraceID = wsTrace.WsTraceID
		}
	}

	var responselist []source.SourceMini
	if len(*list) > 0 {
		param := "ws_trace_id=" + wsTraceID + "&fromOrigin=ws"
		for _, value := range *list {
			if params.MemberUserId != "" {
				value.Path = common.URLAddQuery(value.Path, param)
			}
			responselist = append(responselist, value)
		}
	}

	response.SuccessList(c, total, params.Page, params.Size, responselist)
	return
}

// sourceMiniProgramDetail 素材小程序--详情
func sourceMiniProgramDetail(c *gin.Context) {
	params := &source.ID{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	detail, err := s.SourceMiniProgramDetail(ctx.Wrap(c), params.ID)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// sourceVideoAdd 视频素材-新增
func sourceVideoAdd(c *gin.Context) {
	params := &source.RequestSourceVideo{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	id, err := s.SourceVideoAdd(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	_ = sourceSvc.NewImageService().UpdateRelationGoods(ctx.Wrap(c), uint64(id), uint64(params.EnterpriseID), params.RelationGoods, source.TypeVideo)
	// 视频素材详情
	detail, err := s.SourceVideoDetail(ctx.Wrap(c), id)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// sourceVideoUpdate 视频素材-编辑
func sourceVideoUpdate(c *gin.Context) {
	params := &source.RequestSourceVideo{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	id, err := s.SourceVideoUpdate(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	// 更新关联商品
	_ = sourceSvc.NewImageService().UpdateRelationGoods(ctx.Wrap(c), uint64(id), uint64(params.EnterpriseID), params.RelationGoods, source.TypeVideo)
	response.Success(c, &source.ID{ID: id})
	return
}

// sourceVideoList 视频素材--列表
func sourceVideoList(c *gin.Context) {
	params := &source.RequestSourceSearch{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	params.GroupIDList, _ = source.GetChildIDList(params.GroupID, params.EnterpriseID)
	params.Keyword = helper.FilterEmoji(params.Keyword)
	total, list, err := s.SourceVideoList(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.SuccessList(c, total, params.Page, params.Size, list)
	return
}

// sourceVideoDetail 视频素材--详情
func sourceVideoDetail(c *gin.Context) {
	params := &source.ID{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	detail, err := s.SourceVideoDetail(ctx.Wrap(c), params.ID)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, detail)
	return
}

// SourceDelete 素材删除
func SourceDelete(c *gin.Context) {
	params := &source.RequestSourceDelete{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if len(params.IDs) < 1 {
		response.Success(c, nil)
		return
	}
	rows, err := s.SourceDelete(ctx.Wrap(c), params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	response.Success(c, source.ResponseSourceDelete{Rows: rows})
	return
}

// moveGroup 移动分组
func moveGroup(c *gin.Context) {
	params := &source.RequestGroupMove{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	list, _ := source.GetChildIDList(params.GroupID, params.EnterpriseID)
	if len(list) > 1 {
		ginplus.ResError(c, errors.New("该节点下有分组，无法移动素材"))
		return
	}
	cc := ctx.Wrap(c)
	rows, err := s.MoveGroup(cc, params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}
	if params.GroupID == baozunCacheDao.SOURCE_OTHER_GROUP { // 如果进入回收站，则需要清除置顶
		_ = baozunCacheDao.NewSourceDao().CleanTop(cc, params.SourceIDs)
	}
	response.Success(c, source.ResponseSourceDelete{Rows: rows})
	return
}

// makeSourceMediaId 生成素材ID
func makeSourceMediaId(c *gin.Context) {
	req := &source.MakeMediaRequest{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	domain, _, _, _ := enterpriseService.GetInfo(ctx.Wrap(c), req.EnterpriseID)
	if domain == "" {
		ginplus.ResSuccess(c, errs.EnterpriseNotFund)
		return
	}
	ctxContext := ctx.Wrap(c)
	var detail *source.MediaItemDetail
	var tmpDetail source.MediaItemDetail
	// 先查询此素材的临时素材
	if req.ResourceKey != "" {
		brandDetail, dErr := baozunCacheDao.NewBaozunCacheDao().GetOne(ctxContext, req.EnterpriseID, "resource_key=?", req.ResourceKey)
		if dErr != nil {
			ginplus.ResError(c, errs.SourceNotFund)
			return
		}
		req.SourceID = brandDetail.ID
		tmpDetail.ResourceKey = brandDetail.ResourceKey
		tmpDetail.MediaID = brandDetail.MediaId
		tmpDetail.MediaExpireTime = brandDetail.MediaExpireTime
		detail = &tmpDetail
	} else {
		tDetail, dErr := s.SourceMediaDetail(ctxContext, req.EnterpriseID, req.SourceID, req.SourceType)
		if dErr != nil {
			ginplus.ResError(c, errs.SourceNotFund)
			return
		}
		detail = tDetail
	}
	// 未生成素材或者素材过期了
	if detail.MediaID == "" || ((detail.MediaExpireTime - 10) <= time.Now().Unix()) {
		lockKey := fmt.Sprintf("lock:source_%d", req.SourceID)
		sourceLock, lErr := lock.Mutex(
			ctxContext,
			lock.Key(lockKey),
			lock.TTL(5*time.Second),
			lock.Retry(lock.NoRetry()),
		)

		if lErr != nil {
			log.Info(ctxContext).Msgf("生成临时素材枷锁失败，错误信息：%+v", lErr)
			ginplus.ResError(c, errs.SourceMediaFail)
			return
		}

		if err := sourceLock.Lock(); err != nil {
			log.Debug(ctxContext).Msgf("加锁失败，错误原因:%s, 时间是：%d", err.Error(), time.Now().UnixNano())
			ginplus.ResError(c, errs.SourceMediaFail)
			return
		}
		defer func() {
			_ = sourceLock.Unlock()
		}()

		update := s.MakeMediaForSource(ctxContext, req.EnterpriseID, req.SourceID, req.SourceType, req.SourceUrl, req.SourceName, domain)
		detail.MediaID = update.MediaID
		if req.ResourceKey != "" {
			_ = baozunCacheDao.NewBaozunCacheDao().UpdateData(ctxContext, map[string]interface{}{
				"resource_key":  req.ResourceKey,
				"enterprise_id": req.EnterpriseID,
			}, map[string]interface{}{
				"media_id":          update.MediaID,
				"media_expire_time": update.MediaExpireTime,
				"media_status":      update.MediaStatus,
				"media_fail_reason": update.MediaFailReason,
				"updated_at":        time.Now().Unix(),
			})
		} else {
			_ = s.UpdateSourceMedia(ctxContext, update)
		}
	}
	ginplus.ResSuccess(c, detail)
	return
}

// 获取分组名称
func getSourceName(c *gin.Context) {
	receive := &struct {
		GroupID      int64 `json:"group_id"`
		EnterpriseID int64 `json:"enterprise_id"`
	}{}
	if err := c.ShouldBindJSON(&receive); err != nil {
		ginplus.ResError(c, err)
		return
	}
	domain, _, _, _ := enterpriseService.GetInfo(ctx.Wrap(c), receive.EnterpriseID)
	if domain == "" {
		ginplus.ResSuccess(c, errs.EnterpriseNotFund)
		return
	}
	res := source.GetGroupName(receive.GroupID)
	ginplus.ResSuccess(c, res)
}

// sourceFriendAdd 朋友圈素材-新增
func sourceFriendAdd(c *gin.Context) {
	var params source.RequestSourceFriend
	if err := c.ShouldBind(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)

	err := s.SourceFriendAdd(cc, &params)
	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, map[string]interface{}{})
	return
}
