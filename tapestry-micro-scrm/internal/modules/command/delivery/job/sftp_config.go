package job

import (
	"encoding/json"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	planTaskSftp "micro-scrm/internal/orm/service/plan_task"
	"micro-scrm/internal/pkg/sftp"
)

type sftpConfig struct {
	JobName string
}

func NewSftpConfig(JobName string) *sftpConfig {
	return &sftpConfig{
		JobName: JobName,
	}
}

func (sftpconfig *sftpConfig) GetSFTPConfig(ctx ctx.Context, EnterpriseIds []int64) (configList []micro_scrm.BrandConfigSftpDetail, err error) {
	dataList, err := orm.Service.BrandConfig.GetSFTPConfig(EnterpriseIds, planTaskSftp.SFTP_TYPE_COMMUNITY_DATE)
	if err != nil {
		log.Info(ctx).Msgf("获取企业::%v ,SFTP配置信息错误%v", EnterpriseIds, err)
	}
	// enterpriseID := int64(3323495180054592) 巴黎欧莱雅
	// enterpriseID := 3306120225798016 //测试
	// dataList := []micro_scrm.BrandConfigSftpModel{
	// 	micro_scrm.BrandConfigSftpModel{
	// 		MicroScrmBaseModel: model.MicroScrmBaseModel{
	// 			EnterpriseID: 3323495180054592,
	// 		},
	// 		BrandCode:  "lrl",
	// 		BrandName:  "巴黎欧莱雅",
	// 		SftpConfig: `{"host":"mft.lorealchina.com","port":22,"user_name":"LRLvchangyi","password":"REDACTED_FOR_SOURCE_SHARING"}`,
	// 		SftpPath:   `{"upload_file":"/"}`,
	// 	},
	// }
	for _, data := range dataList {
		var sftpConfig sftp.SSHConfig
		var pathConfig micro_scrm.SFTPPathConfig
		_ = json.Unmarshal([]byte(data.SftpConfig), &sftpConfig)
		_ = json.Unmarshal([]byte(data.SftpPath), &pathConfig)

		configList = append(configList, micro_scrm.BrandConfigSftpDetail{
			ConfigDetail: sftpConfig,
			PathConfig:   pathConfig,
			BrandConfigSftpModel: micro_scrm.BrandConfigSftpModel{
				BrandCode:  data.BrandCode,
				BrandName:  data.BrandName,
				SftpConfig: data.SftpConfig,
				SftpPath:   data.SftpPath,
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					EnterpriseID: data.EnterpriseID,
					BaseModel: model.BaseModel{
						ID:        data.ID,
						CreatedAt: data.CreatedAt,
						UpdatedAt: data.UpdatedAt,
					},
				},
			},
		})
	}

	return
}
