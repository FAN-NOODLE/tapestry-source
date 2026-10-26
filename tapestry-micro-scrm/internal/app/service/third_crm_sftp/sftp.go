package thirdCrmSftpSvc

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"micro-scrm/internal/app"
	thirdCrmSftpDao "micro-scrm/internal/app/dao/third_crm_sftp"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/service"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/sftp"
	"os"
	"strconv"
)

type ThirdCrmSftp struct {
	db          *gorm.DB
	ServiceName string
}

func NewThirdCrmSftp() *ThirdCrmSftp {
	return &ThirdCrmSftp{
		db:          gormx.GetInstance(),
		ServiceName: "third_crm_sftp",
	}
}

func (p *ThirdCrmSftp) GetSFTPConfigList(sftpType int) (configList []micro_scrm.BrandConfigSftpDetail, err error) {
	dataList, dataErr := p.GetBrandSFTPConfig([]int{sftpType})
	if dataErr != nil {
		err = dataErr
		return
	}

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

func (p *ThirdCrmSftp) GetSFTPConfig(c app.Ctx, cond interface{}) (config micro_scrm.BrandConfigSftpDetail, err error) {
	data, dataErr := p.GetBrandSFTPConfigOne(c, cond)
	if dataErr != nil {
		err = dataErr
		return
	}

	var sftpConfig sftp.SSHConfig
	var pathConfig micro_scrm.SFTPPathConfig
	_ = json.Unmarshal([]byte(data.SftpConfig), &sftpConfig)
	_ = json.Unmarshal([]byte(data.SftpPath), &pathConfig)
	config = micro_scrm.BrandConfigSftpDetail{
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
	}

	return
}

// GetBrandSFTPConfig 获取配置了计划任务功能的SFTP信息
func (p *ThirdCrmSftp) GetBrandSFTPConfig(sftpType []int) (data []micro_scrm.BrandConfigSftpModel, err error) {
	err = p.db.Scopes(service.ScopeCond_EQDelete).Where("sftp_type in ?", sftpType).
		Find(&data).Error

	return
}

func (p *ThirdCrmSftp) GetBrandSFTPConfigOne(c app.Ctx, cond interface{}) (data micro_scrm.BrandConfigSftpModel, err error) {
	err = p.db.WithContext(c).Scopes(service.ScopeCond_EQDelete).Where(cond).
		Find(&data).Error

	return
}

func (p *ThirdCrmSftp) GetExistsFile(c app.Ctx, enterpriseID int64, fileType []int) (data []crmSftpFileModels.CrmSftpDownloadedFileModel, err error) {
	return thirdCrmSftpDao.NewThirdCrmSftpFile().GetExistsFile(c, enterpriseID, fileType)
}

func (p *ThirdCrmSftp) GetNotExtractFile(c app.Ctx, enterpriseID int64, fileType []int) (data []crmSftpFileModels.CrmSftpDownloadedFileModel, err error) {
	return thirdCrmSftpDao.NewThirdCrmSftpFile().GetNotExtractFile(c, enterpriseID, fileType)
}

func (p *ThirdCrmSftp) DecryptGpg(c app.Ctx, srcFilePath string, sftp *sftp.SFTP) (LocalFilePath string, err error) {
	srcFile, err := sftp.Client.Open(srcFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::读取远程SFTP文件时异常: %v", srcFilePath, err)
		return
	}
	data, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Info(c).Msgf("%s::读取GPG加密文件内容时异常: %v", srcFilePath, err)
		return
	}
	publicKey, _ := ioutil.ReadFile(`./source/crmSftpKey/` + app.Cfg().Env + `/public.asc`)
	publicKeyObj, _ := crypto.NewKeyFromArmored(string(publicKey))
	publicKeyRing, _ := crypto.NewKeyRing(publicKeyObj)

	privateKey, _ := ioutil.ReadFile(`./source/crmSftpKey/` + app.Cfg().Env + `/private.asc`)
	// pgpMessage, _ := publicKeyRing.Encrypt(binMessage, privateKeyRing)

	// Armored message in pgpMessage.GetArmored()
	// pgpMessage can be obtained from NewPGPMessageFromArmored(ciphertext)

	// pgpMessage can be obtained from a byte array
	var pgpMessage = crypto.NewPGPMessage(data)

	privateKeyObj, _ := crypto.NewKeyFromArmored(string(privateKey))
	unlockedKeyObj, _ := privateKeyObj.Unlock(nil)
	privateKeyRing, _ := crypto.NewKeyRing(unlockedKeyObj)

	message, decryptErr := privateKeyRing.Decrypt(pgpMessage, publicKeyRing, crypto.GetUnixTime())
	errs := decryptErr.(crypto.SignatureVerificationError)
	if decryptErr != nil && errs.Status != constants.SIGNATURE_NOT_SIGNED {
		log.Info(c).Msgf("%s::GPG文件解密异常: %v", srcFilePath, decryptErr)
		return
	}
	// 文件解密后写入本地
	LocalFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath = `./storage/crm_csv/` + LocalFileName
	err = ioutil.WriteFile(LocalFilePath, message.Data, 0777)
	if err != nil {
		return
	}
	// 上传cos
	cs := storage.Init(cos.COS)
	path, err := cs.Upload(c, LocalFilePath, "third_sftp/decrypt_gpg"+LocalFileName)
	if err != nil {
		log.Info(c).Msgf("uploadCloud Error: %v", err)
		return "", err
	}
	// 删除本地文件
	err = os.Remove(LocalFilePath)
	if err != nil {
		log.Info(c).Msgf("delete local file error：%s ,%v", LocalFilePath, err)
	}

	return cs.GetURL(path), err
}

// UpdateFileExTract 更新文件状态
func (p *ThirdCrmSftp) UpdateFileExTract(c app.Ctx, enterpriseID, fileId int64, status int32) error {
	return thirdCrmSftpDao.NewThirdCrmSftpFile().UpdateFileExTract(c, enterpriseID, fileId, status)
}

// NoDecryptGpg 不加密的文件处理
func (p *ThirdCrmSftp) NoDecryptGpg(c app.Ctx, srcFilePath string, sftp *sftp.SFTP, fileType int) (LocalFilePath string, err error) {
	srcFile, err := sftp.Client.Open(srcFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::读取远程SFTP文件时异常: %v", srcFilePath, err)
		return
	}
	data, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Info(c).Msgf("%s::读取文件内容时异常: %v", srcFilePath, err)
		return
	}
	// 文件解密后写入本地
	LocalFileName := strconv.FormatInt(app.UID(), 10) + `.csv`
	LocalFilePath = `./storage/crm_csv/` + LocalFileName
	localFile, err := os.Create(LocalFilePath)
	if err != nil {
		log.Info(c).Msgf("%s::文件创建失败: %v", LocalFilePath, err)
		return
	}
	defer func(localFile *os.File, name string) {
		if err = localFile.Close(); err != nil {
			log.Info(c).Msgf("close local file error：%s ,%v", name, err)
		}
		if err = os.Remove(name); err != nil {
			log.Info(c).Msgf("delete local file error：%s ,%v", name, err)
		}
	}(localFile, LocalFilePath)
	if fileType == 4 {
		fileInfo := bytes.NewBuffer(data)
		reader := csv.NewReader(fileInfo)
		w := csv.NewWriter(localFile)
		func(reader *csv.Reader, w *csv.Writer) {
			for {
				newRow := []string{"id"} // 因为ks上传的文件内容是没有id那一列，coach和sw是有id一列，故需要给ks的文件多加一列
				row, err := reader.Read()
				if err != nil && err != io.EOF {
					return
				}
				if err == io.EOF {
					break
				}
				newRow = append(newRow, row...)
				err = w.Write(newRow)
				if err != nil {
					return
				}
			}

		}(reader, w)
		w.Flush()
	} else {
		err = ioutil.WriteFile(LocalFilePath, data, 0777)
		if err != nil {
			return
		}
	}
	// 上传cos
	cs := storage.Init(cos.COS)
	path, err := cs.Upload(c, LocalFilePath, "third_sftp/no_decrypt_gpg"+LocalFileName)
	if err != nil {
		log.Error(c).Msgf("uploadCloud Error: %v", err)
		return "", err
	}

	return cs.GetURL(path), err
}
