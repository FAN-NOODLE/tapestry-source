package download_http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-download/internal/config"
	"micro-download/internal/modules/download"
	"micro-download/internal/modules/download/delivery"
	"net/http"
	"net/url"
	"strconv"
)

var service *download.Service

func createExport(c *gin.Context) {
	var serviceParams delivery.CreateExport
	if err := c.ShouldBindJSON(&serviceParams); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if serviceParams.SourceRestfulHeaders != nil {
		// 将SDK auth信息放入请求头
		if _, ok := serviceParams.SourceRestfulHeaders["auth-client"]; !ok {
			serviceParams.SourceRestfulHeaders["Auth-Client"] = c.GetHeader("Auth-Client")
		}
		if _, ok := serviceParams.SourceRestfulHeaders["auth-sign"]; !ok {
			serviceParams.SourceRestfulHeaders["Auth-Sign"] = c.GetHeader("Auth-Sign")
		}
		if _, ok := serviceParams.SourceRestfulHeaders["auth-timestamp"]; !ok {
			serviceParams.SourceRestfulHeaders["Auth-Timestamp"] = c.GetHeader("Auth-Timestamp")
		}
	}
	err, created := service.CreateExport(ctx.Wrap(c), serviceParams)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, created)
	return
}

func exportList(c *gin.Context) {
	var params delivery.ExportListReq
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err, list := service.ExportList(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if list.Count == 0 {
		list.List = []delivery.ExportList{}
	}
	ginplus.ResSuccess(c, list)
	return
}

func deleteFile(c *gin.Context) {
	var params delivery.FileID
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := service.DeleteFile(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, map[string]string{})
	return
}

func retry(c *gin.Context) {
	var params delivery.FileID
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := service.Retry(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, map[string]string{})
	return
}

func process(c *gin.Context) {
	var params delivery.FileID
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err, res := service.Process(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

func downLoad(c *gin.Context) {
	downloadId := c.Query("id")
	downloadType := c.Query("type")
	if downloadId == "" {
		ginplus.ResError(c, errors.New("文件标识[id]不能为空"))
		return
	}
	if downloadType == "" {
		ginplus.ResError(c, errors.New("下载类型[type]不能为空"))
		return
	}
	params := &delivery.DownloadReq{}
	dId, err := strconv.ParseInt(downloadId, 10, 64)
	if err != nil {
		ginplus.ResError(c, errors.New("文件标识[id]错误"))
		return
	}
	params.FileID = dId
	dType, err := strconv.ParseInt(downloadType, 10, 8)
	if err != nil {
		ginplus.ResError(c, errors.New("下载类型[type]错误"))
		return
	}
	params.DownloadType = int8(dType)
	err, res := service.Download(ctx.Wrap(c), *params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if res.DownloadType == delivery.DownloadTypePath {
		ginplus.ResSuccess(c, map[string]string{
			"host":     config.CfgData.Blob.Domain,
			"filename": res.Filename,
			"path":     res.Path,
			"fileExt":  res.FileExt,
		})
		return
	}
	if res.DownloadType == delivery.DownloadTypeBin {
		// 像浏览器输出文件
		c.Header("Content-Type", res.ContentType)
		// 兼容中文乱码
		responseFileName := res.Filename + "." + res.FileExt
		encodeFileName := url.QueryEscape(responseFileName)
		c.Header("Content-Disposition", "attachment;filename="+encodeFileName+";filename*=UTF-8''"+encodeFileName)
		c.Data(http.StatusOK, res.ContentType, res.FileData)
		return
	}
}

func detail(c *gin.Context) {
	var params delivery.FileID
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err, res := service.Detail(ctx.Wrap(c), params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
	return
}

func latestExport(c *gin.Context) {
	var params delivery.Label
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	_, res := service.LatestExport(ctx.Wrap(c), params)
	// if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	// }
	ginplus.ResSuccess(c, res)
	return
}

func autoDetail(c *gin.Context) {
	var params delivery.OperationReq
	if err := c.ShouldBindJSON(&params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	_, res := service.GetNoAutoOne(ctx.Wrap(c), params)
	// if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	// }
	ginplus.ResSuccess(c, res)
	return
}
