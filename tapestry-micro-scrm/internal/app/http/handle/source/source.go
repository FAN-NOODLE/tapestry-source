package source

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/cosroute"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/microwork"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"time"
)

type SourceHandle struct {
	handle.BaseHandle
	sourceService     sourceService.Service
	customerService   customer.MessageService
	enterpriseService enterprise.Service
}

func NewSourceHandle() SourceHandle {
	return SourceHandle{
		sourceService: sourceService.NewService(),
	}
}

type CredentialInfo struct {
	Bucket        string `json:"bucket"`
	PrivateBucket string `json:"private_bucket"`
	Region        string `json:"region"`
	SecretID      string `json:"secret_id"`
	SecretKey     string `json:"secret_key"`
	SessionToken  string `json:"session_token"`
	StartTime     int    `json:"start_time"`
	ExpiredTime   int    `json:"expired_time"`
}

type CoverInfo struct {
	FileUrl  string `json:"file_url"`
	CoverUrl string `json:"cover_url"`
}

// UpdateStatus 批量启用/禁用素材
func (h SourceHandle) UpdateStatus(c *gin.Context) {
	var request sourceRequest.UpdateStatusRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.sourceService.UpdateStatus(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

// SendTextMessage 发送文本卡片消息
func (h SourceHandle) SendTextMessage(c *gin.Context) {
	var request sourceRequest.SendRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	corpReq := enterprise.Corp{
		EnterpriseID: request.EnterpriseID,
	}
	// 查询企业信息
	enterpriseData, enterpriseErr := h.enterpriseService.Detail(h.WarpContext(c), corpReq)
	if enterpriseErr != nil {
		h.ResError(c, enterpriseErr)
		return
	}

	// 服务发送
	message, err := h.sourceService.SendTextMessage(h.WarpContext(c), request, enterpriseData)
	if err != nil {
		h.ResError(c, err)
		return
	}
	fmt.Println(message)

	h.ResOk(c)
	return
}

// Credential 获取cos上传的凭证
func (h SourceHandle) Credential(c *gin.Context) {
	bucket := config.CfgData.Storage.COS.Bucket
	privateBucket := cosroute.PrivateBucket()
	region := config.CfgData.Storage.COS.Origin
	ctx := sts.NewClient(
		config.CfgData.Storage.COS.SecretID,
		config.CfgData.Storage.COS.SecretKey,
		nil,
	)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"*",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + region + ":uid/" + config.CfgData.Storage.COS.AppID + ":" + bucket + "/*",
						"qcs::cos:" + region + ":uid/" + config.CfgData.Storage.COS.AppID + ":" + privateBucket + "/*",
					},
				},
			},
		},
	}

	res, err := ctx.GetCredential(opt)
	if err != nil {
		panic(err)
	}

	credential := CredentialInfo{
		Region:        region,
		Bucket:        bucket,
		PrivateBucket: privateBucket,
		SessionToken:  res.Credentials.SessionToken,
		SecretID:      res.Credentials.TmpSecretID,
		SecretKey:     res.Credentials.TmpSecretKey,
		StartTime:     res.StartTime,
		ExpiredTime:   res.ExpiredTime,
	}
	h.ResSuccess(c, credential)
	return
}

// UploadMedia 小程序端导购上传素材
func (h SourceHandle) UploadMedia(c *gin.Context) {
	var request sourceRequest.UploadMediaRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	// 查询企业信息
	enterpriseData, enterpriseErr := h.enterpriseService.Detail(cc, enterprise.Corp{EnterpriseID: int64(request.EnterpriseID)})
	if enterpriseErr != nil {
		h.ResError(c, enterpriseErr)
		return
	}

	sourceIds, compressMediaInfo, compressId, err := h.sourceService.UploadMedia(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if len(compressMediaInfo) > 0 && compressId > 0 {
		threading.GoSafe(func() error {
			h.sourceService.MediaCompress(cc, request, compressMediaInfo, enterpriseData, compressId)
			return nil
		})
	}

	threading.GoSafe(func() error {
		// 神策收到素材上报
		shenCeClient := microwork.NewShenCe()
		for _, id := range sourceIds {
			_, _ = shenCeClient.AssetReceive(cc, int64(request.EnterpriseID), helper.Int64ToString(id), id)
		}
		return nil
	})

	h.ResSuccess(c, sourceIds)
	return
}

// Compress 获取压缩结果
func (h SourceHandle) Compress(c *gin.Context) {
	var request sourceRequest.CompressReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.sourceService.GetCompress(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

// GetMediaInfo 验证cos资源信息是否可以获取mediaId
func (h SourceHandle) GetMediaInfo(c *gin.Context) {
	var request sourceRequest.GetMediaInfoRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := request.EnterpriseID
	mediaType := "image"
	if request.SourceType == 5 {
		mediaType = "video"
	}
	var mediaId string
	if request.Type == 1 {
		qyNewWork := qyWork.NewQyWork(enterpriseId)
		mediaId, err = qyNewWork.GetUploadAttachmentMediaId(h.WarpContext(c), request.Url, mediaType, 1)
		if err != nil {
			h.ResError(c, err)
			return
		}
	} else {
		info, err := common.GetEnterpriseInfo(h.WarpContext(c), enterpriseId)
		if err != nil {
			h.ResError(c, err)
			return
		}
		qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
		mediaId, err = qyWorkApi.GetMediaId(h.WarpContext(c), enterpriseId, request.Url, mediaType)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}
	h.ResSuccess(c, mediaId)
}

// CompressFile 压缩文件
func (h SourceHandle) CompressFile(c *gin.Context) {
	var request struct {
		FilePath string `json:"file_path" binding:"required" required:"图片地址必填"`
		FileType string `json:"file_type" binding:"required" required:"资源类型必填"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	cosFileAddr, err := h.sourceService.CompressFile(h.WarpContext(c), request.FilePath, request.FileType)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, cosFileAddr)
	return
}

// EditVideoTime 修改视频创建时间
func (h SourceHandle) EditVideoTime(c *gin.Context) {
	var request struct {
		FilePath string `json:"file_path" binding:"required" required:"视频地址必填"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	cosFileAddr, err := h.sourceService.EditVideoTime(request.FilePath)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, cosFileAddr)
	return
}

// ReportTimes 上报计次
func (h SourceHandle) ReportTimes(c *gin.Context) {
	var req sourceRequest.ReportTimesReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.sourceService.ReportTimes(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

// CheckFirst 验证是否用户首次打开素材交互提示
func (h SourceHandle) CheckFirst(c *gin.Context) {
	var req struct {
		EnterpriseID int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
		UserId       string `json:"userid" binding:"required" comment:"导购的userid"`
	}
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	hKey := fmt.Sprintf("%s:%d", "source", req.EnterpriseID)
	firstVal, err := app.Redis().HGet(c, hKey, req.UserId).Int()
	if !errors.Is(err, redis.Nil) && err != nil {
		h.ResError(c, err)
		return
	}
	if firstVal == 0 {
		app.Redis().HIncrBy(c, hKey, req.UserId, 1)
	}
	h.ResSuccess(c, firstVal)
}

func (h SourceHandle) MakeCover(c *gin.Context) {
	var request struct {
		FilePath string `json:"file_path" binding:"required" required:"资源地址必填"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	videoUrl, coverUrl, err := h.sourceService.MakeCover(h.WarpContext(c), request.FilePath)
	if err != nil {
		h.ResError(c, err)
		return
	}
	coverInfo := CoverInfo{
		FileUrl:  videoUrl,
		CoverUrl: coverUrl,
	}
	h.ResSuccess(c, coverInfo)
}

func (h SourceHandle) SetTop(c *gin.Context) {
	var req sourceRequest.SourceSetTop
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.sourceService.SetTop(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
