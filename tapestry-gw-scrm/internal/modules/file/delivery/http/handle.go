package file_http

import (
	"errors"
	"gw-scrm/internal/modules/file/storage/cos"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/file"
	fileEntity "gw-scrm/internal/modules/file"
	"gw-scrm/internal/pkg/common"
	scrmErrs "gw-scrm/internal/pkg/errs"
)

// upload 文件上传
func upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, scrmErrs.NotFoundFileUpload)

		return
	}

	// 文件后缀验证
	ext := path.Ext(fh.Filename)
	//var allowExtArr = []string{"JPG", "jpg", "JPEG", "jpeg", "PNG", "png"}
	var allowExtArr = []string{"JPG", "jpg", "PNG", "png"}
	if !common.InArray(ext[1:], allowExtArr) {
		ginplus.ResError(c, scrmErrs.FileUploadExtErr)
		return
	}

	// 确保临时目录存在
	err = os.MkdirAll(config.CfgData.Storage.TmpDir, 0755)
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	// 文件名称 纳秒的md5 + 后缀名
	fileName := common.Md5(strconv.FormatInt(time.Now().UnixNano(), 10)) + ext
	tmpFilePath := config.CfgData.Storage.TmpDir + "/" + fileName

	// 临时保存文件
	if err := c.SaveUploadedFile(fh, tmpFilePath); err != nil {
		ginplus.ResError(c, scrmErrs.FileUploadErr)

		return
	}

	f, err := file.UploadCloud(*c, cos.COS, tmpFilePath, "/uploads")
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	ginplus.ResSuccess(c, f)

	return
}

// credential 生成临时密钥
func credential(c *gin.Context) {

	signature, err := file.Credential()
	if err != nil {
		ginplus.ResError(c, err)

		return
	}

	ginplus.ResSuccess(c, signature)

	return
}

// credential 生成临时密钥
func privateByKey(c *gin.Context) {
	var requestInfo file.PrivateFileKeyReq
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	url, err := file.AesDecrypt(requestInfo.Key)
	if err != nil {
		return
	}
	signature, err := file.Credential()
	if err != nil {
		ginplus.ResError(c, err)

		return
	}
	ginplus.ResSuccess(c, file.PrivateFileInfo{
		Url:        url,
		Credential: signature,
	})

	return
}

func privateURL(c *gin.Context) {
	var requestInfo file.PrivateFileURLReq
	if err := c.ShouldBindJSON(&requestInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	url, err := file.GetPrivateFileURL(c.Request.Context(), requestInfo.Url)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, file.PrivateFileURLRes{
		Url: url,
	})

	return
}

func badatassupload(c *gin.Context) {

	enterpriseId, err := strconv.ParseInt(c.Query("enterpriseId"), 10, 64)
	year, _ := strconv.Atoi(c.Query("year"))
	if year == 0 {
		ginplus.ResError(c, errors.New("传入的年份错误"))
		return
	}
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	if tokenInfo.EnterpriseId != enterpriseId {
		ginplus.ResError(c, errors.New("无权限上传该企业的数据"))
		return
	}
	baKanBandata := new(fileEntity.BaKanBandata)
	baKanBandata.EnterpriseId = enterpriseId
	baKanBandata.Year = year
	baKanBandata.File, err = c.FormFile("file")
	if err != nil {
		ginplus.ResError(c, scrmErrs.NotFoundFileUpload)
		return
	}
	if err := baKanBandata.OperationBaKanBanFile(ctx.Wrap(c)); err == nil {
		ginplus.ResSuccess(c, "上传成功")
	} else {
		ginplus.ResError(c, err)
	}
}
