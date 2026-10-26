package http

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/config"
	"gw-scrm/internal/logic"
	"gw-scrm/internal/modules/file/storage"
	"gw-scrm/internal/modules/file/storage/blob"
	"gw-scrm/internal/modules/tag_csv_upload/delivery"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"io"
	"os"
	"strings"
	"time"
)

// CSVFileUpload vsc标签文件上传
func CSVFileUpload(c *gin.Context) {
	//上传文件流
	formFile, _ := c.FormFile("file")
	//保存数据流
	formFile2, _ := c.FormFile("file")
	openFile, _ := formFile.Open()
	openFile2, _ := formFile2.Open()
	buffer := bytes.NewBuffer(nil)
	buffer2 := bytes.NewBuffer(nil)
	io.Copy(buffer, openFile)
	io.Copy(buffer2, openFile2)
	reader := csv.NewReader(buffer)
	reader2 := csv.NewReader(buffer2)
	reader.FieldsPerRecord = -1
	reader2.FieldsPerRecord = -1

	var fileUrl string
	prefix := logic.GetBolbPrefix()
	prefix += "/micro-scrm/plan-task/manual"
	//获取enterpriseID
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	//获取其他参数
	var dataMap = make(map[string]string, 0)
	d, _ := c.MultipartForm()
	for k, v := range d.Value {
		dataMap[k] = v[0]
	}
	for _, v2 := range d.File {
		dataMap["file_name"] = v2[0].Filename
	}
	//上传数据到blob
	var fileName = fmt.Sprintf("%d", time.Now().Unix()) + "_tagFile.csv"
	var filedetection delivery.FileDetection
	filedetection.FileName = dataMap["file_name"]
	filedetection.EnterpriseID = tokenInfo.EnterpriseId
	init := storage.Init(blob.BLOB)
	newFile, err := os.Create(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	//保证utf-8不会中文数据乱码
	newFile.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(newFile)
	for {
		rowInfo, err0 := reader.Read()
		if err0 == io.EOF {
			break
		}
		if err0 := writer.Write(rowInfo); err0 != nil {
			ginplus.ResError(c, err)
			return
		}
		//刷新流文件
		writer.Flush()
	}

	if err := newFile.Close(); err != nil {
		ginplus.ResError(c, err)
		return
	}
	detection, err := microscrm.FileDetection(ctx.Wrap(c), filedetection)
	if err != nil {
		os.Remove(fileName)
		ginplus.ResError(c, err)
		return
	}
	var countInfo delivery.CountInfo
	marshal, err := json.Marshal(&detection)
	if err != nil {
		os.Remove(fileName)
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(marshal, &countInfo); err != nil {
		os.Remove(fileName)
		ginplus.ResError(c, err)
		return
	}
	if err == nil && countInfo.Count == 0 {
		upload, err := init.Upload(ctx.Wrap(c), fileName, fileName, prefix)
		if err != nil {
			os.Remove(fileName)
			ginplus.ResError(c, err)
			return
		}
		fileUrl = config.CfgData.Storage.Blob.Domain + upload
	} else {
		os.Remove(fileName)
		ginplus.ResError(c, errors.New("该文件已经上传或解析失败！"))
		return
	}
	if err := os.Remove(fileName); err != nil {
		ginplus.ResError(c, err)
		return
	}

	//写入文件路径信息到数据库，获取id
	var tagDownloadInfo delivery.DownLoadFileInfo
	tagDownloadInfo.EnterpriseID = tokenInfo.EnterpriseId
	tagDownloadInfo.TagName = dataMap["tag_name"]
	tagDownloadInfo.FileUrl = fileUrl
	tagDownloadInfo.FileName = dataMap["file_name"]
	tagDownloadInfo.BAType = dataMap["ba_type"]
	tagDownloadInfo.TagGroupName = dataMap["tag_group"]
	IDInfo, err := microscrm.DownloadFileInfo(ctx.Wrap(c), &tagDownloadInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newID delivery.IDInfo
	marshal2, err := json.Marshal(IDInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(marshal2, &newID); err != nil {
		ginplus.ResError(c, err)
		return
	}

	//存在本地数据
	var fileChan = make(chan string, 2000)
	go func(reader2 *csv.Reader) {
		var flag = true
		for {
			row, err := reader2.Read()
			if flag {
				flag = false
				continue
			}
			if err != nil && err != io.EOF {
				ginplus.ResError(c, err)
				return
			}

			if err == io.EOF {
				close(fileChan)
				break
			}
			fileChan <- strings.Join(row, ",")
		}
	}(reader2)

	var data = make([]delivery.FileField, 0)
	var request delivery.GetFileInfo
	request.EnterpriseID = tokenInfo.EnterpriseId
	request.ID = newID.ID
	request.UnionIDEncrypt = dataMap["union_id_encrypt"]
	request.FileName = dataMap["file_name"]
	request.TagName = dataMap["tag_name"]
	for v := range fileChan {
		var newData delivery.FileField
		splitInfo := strings.Split(v, ",")
		if len(splitInfo) >= 4 {
			newData.Ccode = splitInfo[0]
			newData.Unionid = splitInfo[1]
			newData.GradeCode = splitInfo[2]
			newData.AdminBA = splitInfo[3]
		} else {
			ginplus.ResError(c, errors.New("文件格式不符合规范 ，请重新上传！"))
			return
		}
		data = append(data, newData)
		if len(data) == 1000 {
			request.File = data
			_, err := microscrm.UploadCsvFile(ctx.Wrap(c), &request)
			if err != nil {
				ginplus.ResError(c, err)
				return
			}
			data = make([]delivery.FileField, 0)
		}
	}
	if len(data) != 0 {
		request.File = data
		_, err := microscrm.UploadCsvFile(ctx.Wrap(c), &request)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
	}
	//更新完成状态
	var newRequest delivery.NewIDInfo
	newRequest.ID = newID.ID
	newRequest.EnterpriseID = tokenInfo.EnterpriseId
	_, newErr := microscrm.UpdateISExtract(ctx.Wrap(c), &newRequest)
	if newErr != nil {
		ginplus.ResError(c, newErr)
		return
	}
	ginplus.ResOK(c)
}
