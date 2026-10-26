package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/download"
	filemodule "gw-scrm/internal/modules/file"
	"gw-scrm/internal/modules/right_rule"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/downloadcenter"
	"gw-scrm/internal/pkg/sdk/microboss"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"strconv"
)

type downloadHandle struct {
}

func newDownloadHandle() downloadHandle {
	return downloadHandle{}
}

// 下载中心列表
func (handle downloadHandle) exportList(c *gin.Context) {
	var requestInfo download.ExportListReq
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	var req struct {
		EnterpriseId int64 `json:"enterprise_id" comment:"企业ID"`
		AdminId      int64 `json:"admin_id"`
	}
	cc := ctx.Wrap(c)
	req.EnterpriseId = tokenInfo.EnterpriseId
	req.AdminId = tokenInfo.AdminId
	isWhiteUser := false
	res, err := microscrm.DownloadWhiteLists(cc, req)
	if err == nil {
		isWhiteUser = res.(bool)
	}
	operationId := tokenInfo.AdminId
	if isWhiteUser {
		operationId = 0
	}
	request := downloadcenter.ExportListRequest{
		Page:          requestInfo.Page,
		PageSize:      requestInfo.PageSize,
		Status:        requestInfo.Status,
		Filename:      requestInfo.Filename,
		EnterpriseId:  tokenInfo.EnterpriseId,
		OperationIds:  requestInfo.OperationIds,
		OperationId:   operationId,
		CreatedAt:     requestInfo.CreatedAt,
		OperationName: requestInfo.OperationName,
	}
	lists, err := downloadcenter.DC.ExportList(cc, request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var adminIds []int64
	for _, v := range lists.Data.List {
		adminIds = append(adminIds, v.OperationId)
	}
	adminListReq := right_rule.AdminDetailListRequest{
		EnterpriseId: tokenInfo.EnterpriseId,
		Role:         []int{0, 1},
		AdminIds:     adminIds,
		Size:         request.PageSize,
	}
	adminList, err := microboss.GetAdminDetailList(cc, adminListReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	adminInfoMap := make(map[int64]string, 0)
	for _, info := range adminList.List {
		adminInfoMap[info.Id] = info.Username
	}
	for k, v := range lists.Data.List {
		lists.Data.List[k].OperationName = adminInfoMap[v.OperationId]
	}
	ginplus.ResSuccess(c, lists.Data)
}

// 下载文件
func (handle downloadHandle) downloadFile(c *gin.Context) {
	var requestInfo downloadcenter.DownloadRequest
	requestInfo.FileId, _ = strconv.ParseInt(c.Query("id"), 10, 64)
	downloadType, _ := strconv.Atoi(c.Query("type"))
	requestInfo.DownloadType = uint8(downloadType)

	if res, err := downloadcenter.DC.Download(ctx.Wrap(c), requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res.Data)
		return
	}
}

// 删除文件
func (handle downloadHandle) deleteFile(c *gin.Context) {
	var requestInfo downloadcenter.DeleteRequest
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if _, err := downloadcenter.DC.Delete(ctx.Wrap(c), requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, nil)
		return
	}
}

// 查询详情
func (handle downloadHandle) detail(c *gin.Context) {
	var requestInfo downloadcenter.DetailRequest
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if res, err := downloadcenter.DC.ExportDetail(ctx.Wrap(c), requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res.Data)
		return
	}
}

// 重新生成
func (handle downloadHandle) retry(c *gin.Context) {
	var requestInfo downloadcenter.RetryRequest
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if _, err := downloadcenter.DC.Retry(ctx.Wrap(c), requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, nil)
		return
	}
}

// 获取指定业务标签的最近一次导出
func (handle downloadHandle) latestExportDetail(c *gin.Context) {
	var requestInfo downloadcenter.LatestExportDetailRequest
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	requestInfo.OperationId = tokenInfo.AdminId
	if res, err := downloadcenter.DC.LatestExportDetail(ctx.Wrap(c), requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res.Data)
		return
	}
}

// 获取指定业务标签的最近一次导出
func (handle downloadHandle) getAutoDetail(c *gin.Context) {
	var requestInfo downloadcenter.OperationReq
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	requestInfo.OperationId = tokenInfo.AdminId
	cc := ctx.Wrap(c)
	if res, err := downloadcenter.DC.GetAutoDetail(cc, requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		if res.Data.Id > 0 && res.Data.Path != "" {
			res.Data.Path, err = signAutoDownloadPath(cc, res.Data.Path)
			if err != nil {
				ginplus.ResError(c, err)
				return
			}
		}
		ginplus.ResSuccess(c, res.Data)
		return
	}
}

func signAutoDownloadPath(ctx context.Context, path string) (string, error) {
	return filemodule.GetPrivateFileURL(ctx, path)
}

func (handle downloadHandle) whiteList(c *gin.Context) {
	var request struct {
		EnterpriseId int64 `json:"enterprise_id" comment:"企业ID"`
		AdminId      int64 `json:"admin_id"`
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	request.EnterpriseId = tokenInfo.AdminId
	request.AdminId = tokenInfo.AdminId
	if res, err := microscrm.DownloadWhiteLists(ctx.Wrap(c), request); err != nil {
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, res)
		return
	}
}
