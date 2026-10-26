package workExternalJob

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	workExternalDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	planService "micro-scrm/internal/app/service/plan"
	tagsSvc "micro-scrm/internal/app/service/tags"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"strings"

	"micro-scrm/internal/pkg/xjob"
)

type KeySa struct {
	name string
}

func NewKaySa() *KeySa {
	return &KeySa{
		name: "同步三方专属导购",
	}
}

type FileInfo struct {
	SaCode  string `json:"sa_code"`
	UnionId string `json:"union_id"`
}

func (job KeySa) Run(c ctx.Context, _ *xjob.JobHandler) error {
	configList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetBrandSFTPConfig([]int{4, 5})
	if err != nil {
		return err
	}
	for _, config := range configList {
		fileList, err := thirdCrmSftpSvc.NewThirdCrmSftp().GetNotExtractFile(c, config.EnterpriseID, []int{crmSftpFileModels.KAY_SA})
		if err != nil || len(fileList) == 0 {
			continue
		}
		swUseridMap := make(map[string]string)
		if config.BrandName == "sw" {
			swMembers, _ := membersDao.NewMembers().GetMany(c, "enterprise_id= ?", config.EnterpriseID)
			if len(swMembers) > 0 {
				for _, member := range swMembers {
					if member.OpenUserid != "" {
						swUseridMap[member.OpenUserid] = member.Userid
					}
				}
			}
		}
		for _, fileInfo := range fileList {
			res, err := job.ReadFileInfo(fileInfo)
			if err != nil {
				app.Logger().Info(c).Msg(fmt.Sprintf("读取文件失败：%s", err))
				continue
				// return err
			}
			err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, fileInfo.EnterpriseID, fileInfo.ID, 1)
			var updateOldByUnionid []string
			var updateNew [][]string
			for _, v := range res {
				if config.BrandName == "ks" {
					v.SaCode = strings.Trim(v.SaCode, " ")
					// KS传过来的文件 导购编码多出了前缀 000,需要去掉
					if v.SaCode != "000489350" { // 只有这一个用户确实有000，排除她
						v.SaCode = strings.Replace(v.SaCode, "000", "", 1)
					}
				} else if config.BrandName == "sw" {
					if swUserid, ok := swUseridMap[v.SaCode]; ok {
						v.SaCode = swUserid
					}
				}
				updateOldByUnionid = append(updateOldByUnionid, v.UnionID)
				updateNew = append(updateNew, []string{
					v.UnionID, v.SaCode,
				})
				if len(updateOldByUnionid) >= 1000 {
					err := job.UpdateDb(c, config, updateOldByUnionid, updateNew)
					if err != nil {
						return err
					}
					updateOldByUnionid = []string{}
					updateNew = [][]string{}
				}
			}
			if len(updateOldByUnionid) > 0 {
				err = job.UpdateDb(c, config, updateOldByUnionid, updateNew)
			}
			err = thirdCrmSftpSvc.NewThirdCrmSftp().UpdateFileExTract(c, config.EnterpriseID, fileInfo.ID, 2)
			planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
				EnterpriseName: config.BrandName,
				SFTPName:       "SFTP同步顾客Key_sa",
				UploadRes:      helper.RuleThree(err != nil, "异常", "成功"),
			})
		}
	}

	return nil
}

func (job KeySa) UpdateDb(c ctx.Context, config micro_scrm.BrandConfigSftpModel, updateOldByUnionid []string, updateNew [][]string) error {
	return workExternalDao.NewUserRelation().UpdateKeySa(c, config.EnterpriseID, updateOldByUnionid, updateNew)
}

func (job KeySa) ReadFileInfo(Info crmSftpFileModels.CrmSftpDownloadedFileModel) (result []workExternalDao.BaCodeAndUnionid, err error) {
	fileUrl := Info.FileURL
	b, err := helper.ReadRemoteFile(ctx.New(), fileUrl)
	if err != nil {
		return nil, err
	}
	fileInfo := bytes.NewBuffer(b)
	reader := csv.NewReader(fileInfo)
	// 存在本地数据
	var fileChan = make(chan string, 2000)
	// var fileChan = make(chan FileInfo, 2000)
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

	var tmpFile workExternalDao.BaCodeAndUnionid
	for v := range fileChan {
		splitInfo := strings.Split(v, ",")
		if len(splitInfo[2]) < 28 {
			// if len(splitInfo[1]) < 6 || len(splitInfo[2]) < 28 {
			continue
		}
		tmpFile.SaCode = splitInfo[1]
		tmpFile.UnionID = splitInfo[2]

		result = append(result, tmpFile)
	}
	return result, nil
}

func (job KeySa) TagGroupHandler(c ctx.Context, enterpriseId int64, data []tagsReq.ThirdTagsInfo) error {
	var RequestWorkTags []tagsReq.ThirdTagsInfo
	var RequestImpression []tagsReq.ThirdTagsInfo
	for _, item := range data {
		if item.ActionType == 1 {
			RequestWorkTags = append(RequestWorkTags, item)
		}
		if item.ActionType == 2 {
			RequestImpression = append(RequestImpression, item)
		}
	}
	err := tagsSvc.NewWorkTags().SyncThirdTagsToWorkTags(c, enterpriseId, data)
	if err != nil {
		return err
	}
	return nil
}
