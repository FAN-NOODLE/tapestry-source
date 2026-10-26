package workExternalJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io"
	"math"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	tagsSvc "micro-scrm/internal/app/service/tags"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"strings"
	"sync"
)

type SyncExternalTag struct {
	name string
}

func NewSyncExternalTag() *SyncExternalTag {
	return &SyncExternalTag{
		name: "同步第三方外部联系人和标签的关系",
	}
}

func (job SyncExternalTag) Run(c ctx.Context, _ *xjob.JobHandler) error {
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
		res, err := job.ReadFileInfo(value)
		if err != nil {
			log.Error(c).Msg(fmt.Sprintf("读取文件失败：%s", err))
			continue
		}
		if err = job.ExternalTagHandler(c, value, res); err != nil {
			log.Error(c).Msg(fmt.Sprintf("处理外部联系人和标签关系报错：%s", err))
			continue
		}
		_ = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, value.EnterpriseID, value.ID, 2)
	}

	return nil
}

func (job SyncExternalTag) DownLoadFileInfo(c ctx.Context, enterpriseId int64, ch chan crmSftpFileModels.CrmSftpDownloadedFileModel) error {
	fileInfoList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetExistsFile(c, enterpriseId, []int{crmSftpFileModels.TAG_CUSTOMER_REL})
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

func (job SyncExternalTag) ReadFileInfo(Info crmSftpFileModels.CrmSftpDownloadedFileModel) (result []tagsReq.ExternalUserThirdTag, err error) {
	fileUrl := Info.FileURL
	b, err := helper.ReadRemoteFile(ctx.New(), fileUrl)
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
		var Data tagsReq.ExternalUserThirdTag
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 7 {
			if splitInfo[4] == "" {
				continue
			}
			Data.ThirdTagID = splitInfo[1]
			Data.TagName = splitInfo[2]
			Data.Mobile = splitInfo[3]
			Data.UnionId = splitInfo[4]
			Data.ActionType, _ = strconv.Atoi(splitInfo[5])
		}
		result = append(result, Data)
	}
	return result, nil
}

func (job SyncExternalTag) ExternalTagHandler(c ctx.Context, fileInfo crmSftpFileModels.CrmSftpDownloadedFileModel, data []tagsReq.ExternalUserThirdTag) error {
	var workTagsAdd []tagsReq.ExternalUserThirdTag
	var workTagsDel []tagsReq.ExternalUserThirdTag
	err := thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, fileInfo.EnterpriseID, fileInfo.ID, 1)
	if err != nil {
		return err
	}
	for _, item := range data {
		if item.ActionType == 1 {
			workTagsAdd = append(workTagsAdd, item)
		}
		if item.ActionType == 2 {
			workTagsDel = append(workTagsDel, item)
		}
	}

	tagsAdd := make(chan []tagsReq.ExternalUserThirdTag, 5)

	size := 1000

	go func() {
		// 处理外部联系人和标签关系
		addTotal := len(workTagsAdd)
		if addTotal > 0 {
			addMaxPage := int(math.Ceil(float64(addTotal / size)))
			addMaxPage += 1
			for i := 1; i <= addMaxPage; i++ {
				start := (i - 1) * size
				end := i * size
				addData := workTagsAdd[start:]
				if end < addTotal {
					addData = workTagsAdd[start:end]
				}
				tagsAdd <- addData
			}
		}
		close(tagsAdd)
	}()

	wg := sync.WaitGroup{}
	// 控制企微接口频率 160*60=9600 企微接口一秒最多1w
	p, _ := ants.NewPool(10)
	defer p.Release()

	for v := range tagsAdd {
		wg.Add(1)
		newAddTags := v
		_ = p.Submit(func() {
			defer wg.Done()
			err = tagsSvc.NewWorkTags().SyncExternalWorkTags(c, fileInfo.EnterpriseID, newAddTags, 1)
			if err != nil {
				log.Error(c).Msgf("添加标签 异常:%v", err)
			}
		})
	}
	wg.Wait()

	tagsDel := make(chan []tagsReq.ExternalUserThirdTag, 5)
	go func() {
		delTotal := len(workTagsDel)
		if delTotal > 0 {
			delMaxPage := int(math.Ceil(float64(delTotal / size)))
			delMaxPage += 1
			for i := 1; i <= delMaxPage; i++ {
				start := (i - 1) * size
				end := i * size
				delData := workTagsDel[start:]
				if end < delTotal {
					delData = workTagsDel[start:end]
				}
				tagsDel <- delData
			}
		}
		close(tagsDel)
	}()

	delWg := sync.WaitGroup{}
	for v := range tagsDel {
		delWg.Add(1)
		newDelTags := v
		_ = p.Submit(func() {
			defer delWg.Done()
			err = tagsSvc.NewWorkTags().SyncExternalWorkTags(c, fileInfo.EnterpriseID, newDelTags, 2)
			if err != nil {
				log.Error(c).Msgf("删除标签 异常:%v", err)
			}
		})
	}
	delWg.Wait()

	return nil
}
