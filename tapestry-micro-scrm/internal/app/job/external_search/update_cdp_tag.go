package externalSearchJob

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	"io/ioutil"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	planService "micro-scrm/internal/app/service/plan"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type UpdateCdpTags struct {
	sftpClient          *sftp.SFTP
	Name                string
	SftpRootPath        string
	LocalDir            string
	InertCdpTagPageSize int
	UpdateExtraPageSize int
	TagSlice            []string
}

func NewUpdateCdpTags() UpdateCdpTags {
	return UpdateCdpTags{
		Name:                "每天拉取CDP标签并全量刷新数据",
		sftpClient:          nil,
		SftpRootPath:        "/outbound/cdp/",
		LocalDir:            "./storage/cdp_tag/",
		InertCdpTagPageSize: 100000,
		UpdateExtraPageSize: 100000,
		TagSlice: []string{
			externalSearchModel.CdpTagAmountTotal,           //一年总金额
			externalSearchModel.CdpTagHistoryAmountTotal,    //历史总金额
			externalSearchModel.CdpTagLastPurchaseDate,      //最后一次消费日期
			externalSearchModel.CdpTagLastPurchaseStoreCode, //最后一次消费店铺
			externalSearchModel.CdpTagCategoryPreference,    //品类偏好
			externalSearchModel.CdpTagPurchasedDate,         //历史消费日期
			externalSearchModel.CdpTagHighWilling,           //高意愿
			externalSearchModel.CdpTagPriceSensitivity,      //价格敏感度
			externalSearchModel.CdpTagBuyCategoryPreference, //购买品类偏好
			externalSearchModel.CdpTagPurchasedMonth,        //购买时间分布
			externalSearchModel.CdpTagManagerStoreCode,      //管理柜台
		},
	}
}

func (job UpdateCdpTags) Run(c ctx.Context, handler *xjob.JobHandler) error {
	//根据每天全量CDP标签更新顾客属性表
	sftpConfigList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList(micro_scrm.SftpTypeCdpTag)
	if err != nil {
		fmt.Printf("%s::获取SFTP配置失败: %v", job.Name, err)
	}
	for _, config := range sftpConfigList {
		cacheRes, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, config.EnterpriseID, "store_name")
		if cacheRes.Name != "coach" {
			continue
		}
		if client, clientErr := sftp.NewSFTPClient(config.ConfigDetail); clientErr != nil {
			fmt.Printf("%s::企业ID:%d, 连接SFTP失败: %v", job.Name, config.EnterpriseID, clientErr)
			return nil
		} else {
			job.sftpClient = client
			err = job.handle(c, config.EnterpriseID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (job UpdateCdpTags) handle(c ctx.Context, enterpriseId int64) error {
	// 由于每天都是全量所以须要先清空之前的表
	err := customerSvc.NewExternalSearch().TruncateCdpTagTable(c, enterpriseId)
	if err != nil {
		return err
	}
	// 获取昨天的日期
	today := time.Now().Format("20060102")
	//today := "20241127"
	var hasErr bool
	for _, tagStr := range job.TagSlice {
		tagFileName := fmt.Sprintf("tag_%s_%s.zip", tagStr, today)
		err := job.pullFileAndSaveData(c, enterpriseId, tagFileName, tagStr)
		if err != nil {
			notice.NewWechat().SendTextMessage(c, fmt.Sprintf("%s异常：%s，%s", job.Name, tagFileName, err.Error()))
			hasErr = true
			continue
		}
	}
	if hasErr {
		//有任何异常则不再执行数据刷新
		return nil
	}
	//删除临时文件
	_ = helper.DeleteFilesInDirectory(job.LocalDir)
	//CDP tag表数据合并
	err = customerSvc.NewExternalSearch().CdpTagDataMerge(c, enterpriseId, false)
	//根据每天全量CDP标签更新顾客属性表
	err = job.UpdateExtraCk(c, enterpriseId)
	if err != nil {
		return err
	}
	err = customerSvc.NewExternalSearch().ExtraDataMerge(c, enterpriseId, false)
	if err != nil {
		return err
	}
	planService.NewTaskNotice().SftpPullNotice(c, planService.SFTPMessageData{
		EnterpriseName: "coach",
		SFTPName:       job.Name,
		UploadRes:      fmt.Sprintf("%s,执行成功", job.Name),
	})
	return nil
}

func (job UpdateCdpTags) pullFileAndSaveData(c ctx.Context, enterpriseId int64, tagFileName string, tagStr string) error {
	//远端文件路径
	sftpFilePath := job.SftpRootPath + tagFileName
	//本地存储文件路径
	localFilePath := job.LocalDir + tagFileName
	sftpFile, err := job.sftpClient.Client.Open(sftpFilePath)
	defer sftpFile.Close()
	if err != nil {
		fmt.Printf("%s::获取目录下文件失败: %v", job.Name, err)
		return err
	}
	dstFile, err := os.Create(localFilePath)
	defer dstFile.Close()
	if _, err := sftpFile.WriteTo(dstFile); err != nil {
		fmt.Printf("%s::写入文件失败: %v", localFilePath, err)
		return err
	}

	fileList, err := helper.UnZipToDir(localFilePath, job.LocalDir)
	if err != nil {
		fmt.Printf("%s::文件解压失败: %v", localFilePath, err)
		return err
	}
	if len(fileList) == 0 {
		fmt.Printf("%s::zip中没有文件: %v", localFilePath, err)
		return err
	}
	re := regexp.MustCompile(`^tag_([a-zA-Z0-9_]+)_(\d{8})__.*\.csv$`)
	fileName := filepath.Base(fileList[0])
	match := re.FindStringSubmatch(fileName)
	if match == nil || tagStr != match[1] {
		fmt.Printf("%s::csv文件名匹配失败", fileList[0])
		return errors.New("csv文件名匹配失败")
	}
	dataByte, _ := ioutil.ReadFile(fileList[0])
	reader := csv.NewReader(bytes.NewBuffer(dataByte))
	// 存在本地数据
	var fileChan = make(chan string, 10000)
	go func(reader *csv.Reader) {
		var flag = true
		for {
			row, err := reader.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				fmt.Println(err.Error())
				continue
			}
			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader)
	var dataList []externalSearchModel.CdpExternalTag
	for v := range fileChan {
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) < 2 || splitInfo[0] == "" || splitInfo[1] == "" {
			continue
		}
		var tmpFile externalSearchModel.CdpExternalTag
		tmpFile.Unionid = strings.TrimSpace(splitInfo[0])
		tmpFile.TagValue = strings.TrimSpace(splitInfo[1])

		switch tagStr {
		case externalSearchModel.CdpTagAmountTotal,
			externalSearchModel.CdpTagHistoryAmountTotal:
			tmpFile.TagValue = helper.StringFloatToStringInt(tmpFile.TagValue)
			//最后一次消费日期
		case externalSearchModel.CdpTagLastPurchaseDate:
			tmpTime, err := time.Parse("2006-01-02", tmpFile.TagValue)
			if err != nil {
				continue
			}
			tmpFile.TagValue = helper.Int64ToString(tmpTime.Unix())
			//历史消费日期
		case externalSearchModel.CdpTagPurchasedDate:
			tmpTimeSlice := strings.Split(tmpFile.TagValue, "|")
			var purchasedDateErr = false
			for i, timeStr := range tmpTimeSlice {
				tmpTime, err := time.Parse("2006-01-02", timeStr)
				if err != nil {
					purchasedDateErr = true
				}
				tmpTimeSlice[i] = helper.Int64ToString(tmpTime.Unix())
			}
			if purchasedDateErr {
				continue
			}
			tmpFile.TagValue = strings.Join(tmpTimeSlice, "|")
			/** 这些类型不用处理
			case externalSearchModel.CdpTagCategoryPreference, //品类偏好
				externalSearchModel.CdpTagHighWilling,           //高意愿
				externalSearchModel.CdpTagPriceSensitivity,      //价格敏感度
				externalSearchModel.CdpTagPurchasedMonth,		//购买时间分布
				externalSearchModel.CdpTagBuyCategoryPreference,	//购买品类偏好
				externalSearchModel.CdpTagLastPurchaseStoreCode: //最近消费门店
			*/
		}
		tmpFile.TagType = tagStr
		dataList = append(dataList, tmpFile)
		if len(dataList) > job.InertCdpTagPageSize {
			err = externalSearchDao.NewCdpExternalTagCk().BatchSaveData(c, enterpriseId, dataList)
			if err != nil {
				return err
			}
			dataList = []externalSearchModel.CdpExternalTag{}
		}
	}
	// 在管道处理完毕后，如果还有剩余的数据，执行最后一次保存操作
	if len(dataList) > 0 {
		err = externalSearchDao.NewCdpExternalTagCk().BatchSaveData(c, enterpriseId, dataList)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateExtraCk 根据每天全量CDP标签更新顾客属性表
func (job UpdateCdpTags) UpdateExtraCk(c ctx.Context, enterpriseId int64) error {
	err := externalSearchDao.NewExternalExtraCk().UpdateExternalExtraCkByCdpTag(c, enterpriseId, job.UpdateExtraPageSize)
	if err != nil {
		return err
	}
	return nil
}
