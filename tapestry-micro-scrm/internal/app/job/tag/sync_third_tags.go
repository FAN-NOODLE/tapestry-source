package tagJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	tagsSvc "micro-scrm/internal/app/service/tags"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"strings"

	"micro-scrm/internal/pkg/xjob"
)

type ThirdTags struct {
	name string
}

func NewThirdTags() *ThirdTags {
	return &ThirdTags{
		name: "同步第三方的标签组",
	}
}

func (job ThirdTags) Run(c ctx.Context, handler *xjob.JobHandler) error {
	configList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetBrandSFTPConfig([]int{4})
	if err != nil {
		return err
	}
	ch := make(chan crmSftpFileModels.CrmSftpDownloadedFileModel, 2)
	go func() {
		for _, config := range configList {
			if err = job.DownLoadFileInfo(c, config.EnterpriseID, ch); err != nil {
				continue
			}
		}
		close(ch)
	}()

	for value := range ch {
		res, err := job.ReadFileInfo(c, value)
		if err != nil {
			log.Error(c).Msg(fmt.Sprintf("读取文件失败：%s", err))
			continue
		}
		if err = job.TagGroupHandler(c, value, res); err != nil {
			log.Error(c).Msg(fmt.Sprintf("三方标签同步报错：%s", err))
			continue
		}
		_ = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, value.EnterpriseID, value.ID, 2)
	}
	return nil
}

// DownLoadFileInfo 同步拉取要处理的文件内容
func (job ThirdTags) DownLoadFileInfo(c ctx.Context, enterpriseId int64, ch chan crmSftpFileModels.CrmSftpDownloadedFileModel) error {
	fileInfoList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetExistsFile(c, enterpriseId, []int{crmSftpFileModels.TAG})
	if err != nil {
		return err
	}
	for _, v := range fileInfoList {
		if v.IsExtract == 0 {
			ch <- v
		}
	}
	return nil
}

// ReadFileInfo csv 读取远程文件内容
func (job ThirdTags) ReadFileInfo(c ctx.Context, Info crmSftpFileModels.CrmSftpDownloadedFileModel) (result []tagsReq.ThirdTagsInfo, err error) {
	fileUrl := Info.FileURL
	b, err := helper.ReadRemoteFile(c, fileUrl)
	if err != nil {
		return nil, err
	}
	fileInfo := bytes.NewBuffer(b)
	reader := csv.NewReader(fileInfo)
	// 存在本地数据
	var fileChan = make(chan string, 2000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				return
			}
			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	for v := range fileChan {
		var TagData tagsReq.ThirdTagsInfo
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 5 {
			TagData.ID = splitInfo[0]
			TagData.TagGroup = splitInfo[1]
			TagData.TagName = splitInfo[2]
			TagData.ActionType, _ = strconv.Atoi(splitInfo[3])
			TagData.TagDescription = splitInfo[4]
		}
		result = append(result, TagData)
	}
	return result, nil
}

// TagGroupHandler 标签和标签组的功能同步
func (job ThirdTags) TagGroupHandler(c ctx.Context, fileInfo crmSftpFileModels.CrmSftpDownloadedFileModel, data []tagsReq.ThirdTagsInfo) error {
	var requestWorkTags []tagsReq.ThirdTagsInfo
	err := thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, fileInfo.EnterpriseID, fileInfo.ID, 1)
	if err != nil {
		return err
	}
	for _, item := range data {
		if item.ActionType == 1 {
			requestWorkTags = append(requestWorkTags, item)
		}
		if item.ActionType == 2 {
			var requestImpression tagsReq.CreateGroupAndTag
			requestImpression.GroupName = item.TagGroup
			requestImpression.Name = item.TagName
			requestImpression.EnterpriseID = fileInfo.EnterpriseID
			requestImpression.ThirdTagId = item.ID
			requestImpression.Description = item.TagDescription
			err := tagsSvc.NewImpression().CreateGroupAndTag(c, requestImpression)
			if err != nil {
				return err
			}
		}
	}
	err = tagsSvc.NewWorkTags().SyncThirdTagsToWorkTags(c, fileInfo.EnterpriseID, requestWorkTags)
	if err != nil {
		return err
	}
	return nil
}
