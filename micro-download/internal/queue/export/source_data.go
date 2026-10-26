package export

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"micro-download/internal/common"
	"micro-download/internal/common/file"
	"micro-download/internal/common/source"
	"micro-download/internal/config"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
	"micro-download/internal/pkg/coroutine"
	"micro-download/internal/pkg/redis"
)

type SourceData struct {
}

func (s *SourceData) Recv(c ctx.Context, msg *string) error {
	// 获取消息数据模型
	model, err := getOperationModel(msg)
	if err != nil {
		return err
	}
	// 获取任务详情
	err, detail := model.Detail(c)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("获取任务成功...%s", *msg)
	// 更改状态为生成中
	err = model.ChangeStatus(c, delivery.StatusGenerating)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("更改任务状态成功...%s", *msg)
	// 获取任务概览
	overview, params, err := source.Overview(c, detail)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("获取任务概览成功...%s,需要page:%d", *msg, overview.PageNum)
	// 请求并写入文件
	tmpPath, err := chunkReadAndWrite(c, detail, overview, params)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("获取全部数据成功...%s", *msg)
	// 上传到cos
	err, cosPath := file.UploadCos(c, tmpPath)
	// err, cosPath := file.UploadBlob(c, tmpPath)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("上传文件成功...%s", *msg)
	// 更改文件状态
	err = model.Success(c, cosPath)
	if err != nil {
		_ = model.Failed(c, err.Error())
		return err
	}
	log.Info(c).Msgf("更改任务状态完成...%s", *msg)
	// 清除进度信息
	redis.Client().Del(c, fmt.Sprintf(delivery.SourceRedisProcess, detail.ID))
	log.Info(c).Msgf("清除进度信息成功...%s", *msg)
	return nil
}

// chunkReadAndWrite 分片读取数据并写入文件
func chunkReadAndWrite(ctx ctx.Context, detail *modes.DownloadSource, overview *delivery.SourceOverview, requestParams interface{}) (string, error) {
	task := coroutine.NewSimpleTask(config.CfgData.Downloader.ParallelNumber, ctx)
	// 推任务
	var page int64 = 1
	for ; page <= overview.PageNum; page++ {
		// for ; page <= 10000; page++ {
		task.PushTask(func(params ...interface{}) interface{} {
			tPage := params[0].([]interface{})[0].(int64)
			// tPage := int64(1)
			isPaginationTypeID := false
			// 如果导出类型是Restful，分页类型是ID，则需要替换请求参数中的最小ID和最大ID，否则按照原逻辑获取数据
			if detail.SourceType == delivery.SourceTypeRestful {
				tParams := requestParams.(*delivery.SourceTypeRestfulData)
				isPaginationTypeID = tParams.SourceRestfulPaginationType == delivery.RestfulPaginationTypeID
			}
			if !isPaginationTypeID {
				err, list := source.Read(ctx, tPage, detail.ID, detail.SourceType, overview.Title, requestParams, overview.Value)
				if err != nil {
					return err
				}
				return list
			}
			// 从下载详情中获取原始请求参数
			tParams := &delivery.SourceTypeRestfulData{}
			err := json.Unmarshal([]byte(detail.SourceData), tParams)
			if err != nil {
				return err
			}
			log.Info(ctx).Msgf("request_min_and_max_id, page: %v, from: %d, page_size: %v", tPage-1, overview.PageFrom[tPage-1], common.DefaultSplitSize)
			err, list := source.Read(ctx, tPage-1, detail.ID, detail.SourceType, overview.Title, requestParams, overview.Value)
			if err != nil {
				return err
			}
			return list
		}, page)
	}
	// 删除每次批量的结果
	task.SetBatchExecutedDeleteRecord(true)
	// 每次批量后进行进度更新、文件追加
	filename := file.GeneratePathFileName(detail.Name, detail.Type)
	log.Info(ctx).Msgf("临时文件将被保存在:%s", filename)
	task.SetBatchExecuted(func(executedTaskLength int, executedResult []interface{}) error {
		var data [][]any
		for _, pageList := range executedResult {
			data = append(data, pageList.([][]any)...)
		}
		log.Info(ctx).Msgf("批次执行完毕%d,%d...", executedTaskLength, detail.ID)
		// 进度更新
		// process(ctx, detail.ID, int((float64(executedTaskLength * common.DefaultSplitSize) / float64(overview.Count)) * 100))
		startLine := (executedTaskLength-config.CfgData.Downloader.ParallelNumber)*common.DefaultSplitSize + 1
		if startLine <= 0 {
			startLine = 1
		}
		log.Info(ctx).Msgf("进度更新完毕开始写入批次文件%d,%d...", executedTaskLength, detail.ID)
		startWriteFileTime := time.Now()
		err := file.Generate(ctx, detail.Type, filename, startLine, data, overview, int64(executedTaskLength) == overview.PageNum)
		// err := file.Generate(ctx, detail.Type, filename, startLine, data, overview, int64(executedTaskLength) == 10000)
		log.Info(ctx).Msgf("文件写入完毕%d,%d...,耗时:%s", executedTaskLength, detail.ID, time.Since(startWriteFileTime))
		return err
	})
	log.Info(ctx).Msgf("投递协程池任务成功%d...", detail.ID)
	if overview.PageNum == 0 {
		file.Generate(ctx, detail.Type, filename, 1, [][]any{}, overview, true)
	}
	_, err := task.Run()
	if err != nil {
		return "", err
	}
	return filename, nil
}

// process 更新生成进度
func process(c ctx.Context, fileId int64, process int) {
	redis.Client().Set(c, fmt.Sprintf(delivery.SourceRedisProcess, fileId), process, time.Second*1800)
}

// getOperationModel 获取消息数据操作模型
func getOperationModel(msg *string) (*modes.DownloadSource, error) {
	// msg 就是file_id
	fileId, err := strconv.ParseInt(*msg, 10, 64)
	if err != nil {
		return nil, err
	}
	// 查询文件详情
	params := &delivery.FileID{FileID: fileId}
	model := &modes.DownloadSource{}
	model.SetParams(*params)
	return model, nil
}
