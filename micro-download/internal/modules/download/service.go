package download

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io/ioutil"
	"micro-download/internal/common/source"
	"micro-download/internal/config"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
	"micro-download/internal/pkg/cos"
	"micro-download/internal/pkg/redis"
	"micro-download/internal/queue"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 应用服务
type Service struct{}

const downloadPresignedURLExpire = 30 * time.Minute

func normalizeDownloadObjectKey(filePath string) (string, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", errors.New("下载文件路径不能为空")
	}

	if parsedURL, err := url.Parse(filePath); err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		filePath = parsedURL.Path
	}

	objectKey := strings.TrimLeft(filePath, "/")
	if !strings.HasPrefix(objectKey, "download/") {
		return "", errors.New("下载文件路径不合法")
	}

	return objectKey, nil
}

func presignedDownloadURL(c ctx.Context, filePath string) (string, error) {
	objectKey, err := normalizeDownloadObjectKey(filePath)
	if err != nil {
		return "", err
	}

	presignedURL, err := cos.Cos().Object.GetPresignedURL(
		c,
		http.MethodGet,
		objectKey,
		config.CfgData.Cos.SecretId,
		config.CfgData.Cos.SecretKey,
		downloadPresignedURLExpire,
		nil,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// 创建
func (s Service) CreateExport(c ctx.Context, params delivery.CreateExport) (error, *delivery.CreatedExport) {
	err, fileID := source.Create(c, params)
	if err != nil {
		return err, nil
	}
	err = queue.Dispatch(c, queue.QueueFileGenerate, strconv.FormatInt(fileID, 10))
	if err != nil {
		return err, nil
	}
	return nil, &delivery.CreatedExport{FileID: fileID}
}

// 列表
func (s Service) ExportList(c ctx.Context, params delivery.ExportListReq) (error, *delivery.ExportListRes) {
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	if params.Page == 0 {
		params.Page = 1
	}
	model := &modes.DownloadSource{}
	model.SetParams(params)
	return model.ExportList(c)
}

// 删除
func (s Service) DeleteFile(c ctx.Context, id delivery.FileID) error {
	model := &modes.DownloadSource{}
	model.SetParams(id)
	err, detail := model.Detail(c)
	if err != nil {
		return err
	}
	_, _ = cos.Cos().Object.Delete(c, detail.Path)
	func(path string) {}(detail.Path)
	return model.Delete(c)
}

// 重试
func (s Service) Retry(c ctx.Context, id delivery.FileID) error {
	model := &modes.DownloadSource{}
	model.SetParams(id)
	err := model.Retry(c)
	if err != nil {
		return err
	}
	return queue.Dispatch(c, queue.QueueFileGenerate, strconv.FormatInt(id.FileID, 10))
}

// 查看生成进度
func (s Service) Process(c ctx.Context, id delivery.FileID) (error, *delivery.ProcessRes) {
	res := &delivery.ProcessRes{}
	model := &modes.DownloadSource{}
	model.SetParams(id)
	err, detail := model.Detail(c)
	if err != nil {
		return err, nil
	}
	if detail.Status != delivery.StatusGenerating {
		if detail.Status == delivery.StatusCreated {
			res.Process = 0
		} else {
			res.Process = 100
		}
		return nil, res
	}
	process, err := redis.Client().Get(c, fmt.Sprintf(delivery.SourceRedisProcess, id.FileID)).Int()
	if err != nil {
		return err, nil
	}
	res.Process = process
	return nil, res
}

// 下载
func (s Service) Download(c ctx.Context, param delivery.DownloadReq) (error, *delivery.DownloadData) {
	res := &delivery.DownloadData{}
	modelParam := delivery.FileID{FileID: param.FileID}
	model := &modes.DownloadSource{}
	model.SetParams(modelParam)
	err, detail := model.Detail(c)
	if err != nil {
		return errors.New("数据不存在"), nil
	}
	res.DownloadType = param.DownloadType
	res.Path = detail.Path
	res.Filename = detail.Name
	if detail.Type == delivery.TypeExcel {
		res.FileExt = "xlsx"
		res.ContentType = "application/excel"
	}
	if detail.Type == delivery.TypeCsv {
		res.FileExt = "csv"
		res.ContentType = "application/text"
	}
	if detail.Type == delivery.TypeZip {
		res.FileExt = "zip"
		res.ContentType = "application/octet-stream"
	}

	presignedURL, err := presignedDownloadURL(c, detail.Path)
	if err != nil {
		return err, nil
	}
	res.Path = presignedURL

	if param.DownloadType == delivery.DownloadTypePath {
		return nil, res
	}
	if param.DownloadType == delivery.DownloadTypeBin {
		rep, err := http.Get(presignedURL)
		if err != nil {
			return err, nil
		}
		data, err := ioutil.ReadAll(rep.Body)
		defer func() {
			_ = rep.Body.Close()
		}()
		if err != nil {
			return err, nil
		}
		res.FileData = data
		return nil, res
	}
	return errors.New("不支持的下载类型"), nil
}

// 查看任务详情
func (s Service) Detail(c ctx.Context, id delivery.FileID) (error, *delivery.ExportList) {
	out := &delivery.ExportList{}
	model := &modes.DownloadSource{}
	model.SetParams(id)
	err, res := model.Detail(c)
	if err != nil {
		return err, nil
	}
	out.ID = res.ID
	out.CreatedAt = res.CreatedAt
	out.FailedReason = res.FailedReason
	out.Status = res.Status
	out.Type = res.Type
	out.OperationName = res.OperationName
	out.Filename = res.Name
	out.CompleteAt = res.CompleteAt
	return nil, out
}

// 查看任务详情
func (s Service) LatestExport(c ctx.Context, label delivery.Label) (error, *delivery.ExportList) {
	model := &modes.DownloadSource{}
	model.SetParams(label)
	err, res := model.LatestDetailByLabel(c)
	if err != nil {
		return err, nil
	}
	out := &delivery.ExportList{
		ID:            res.ID,
		Filename:      res.Name,
		Type:          res.Type,
		OperationName: res.OperationName,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt,
		FailedReason:  res.FailedReason,
		CompleteAt:    res.CompleteAt,
		Host:          config.CfgData.Blob.Domain,
		Path:          res.Path,
	}
	if res.Type == delivery.TypeExcel {
		out.FileExt = "xlsx"
	}
	if res.Type == delivery.TypeCsv {
		out.FileExt = "csv"
	}
	if res.Type == delivery.TypeZip {
		out.FileExt = "zip"
	}

	return nil, out
}

// 获取未自动下载的资源
func (s Service) GetNoAutoOne(c ctx.Context, auto delivery.OperationReq) (error, *delivery.ExportList) {
	model := &modes.DownloadSource{}
	model.SetParams(auto)
	err, res := model.GetNoAutoOne(c)
	if err != nil {
		return err, nil
	}
	out := &delivery.ExportList{
		ID:            res.ID,
		Filename:      res.Name,
		Type:          res.Type,
		OperationName: res.OperationName,
		Status:        res.Status,
		CreatedAt:     res.CreatedAt,
		FailedReason:  res.FailedReason,
		CompleteAt:    res.CompleteAt,
		Host:          config.CfgData.Blob.Domain,
		Path:          res.Path,
	}
	if res.Type == delivery.TypeExcel {
		out.FileExt = "xlsx"
	}
	if res.Type == delivery.TypeCsv {
		out.FileExt = "csv"
	}
	if res.Type == delivery.TypeZip {
		out.FileExt = "zip"
	}
	if res.ID > 0 {
		changAuto := delivery.FileID{}
		changAuto.FileID = res.ID
		changModel := &modes.DownloadSource{}
		changModel.SetParams(changAuto)
		err = changModel.ChangeAuto(c, 1)
	}
	return err, out
}
