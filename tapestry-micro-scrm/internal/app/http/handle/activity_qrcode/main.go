package activityQRCodeHandle

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	activityQRCodeSvc "micro-scrm/internal/app/service/activity_qrcode"
	"micro-scrm/internal/pkg/common"
	"strings"
)

type Handle struct {
	handle.BaseHandle
	mainSvc  *activityQRCodeSvc.MainSvc
	frontSvc *activityQRCodeSvc.FrontSvc
}

func NewHandle() Handle {
	return Handle{
		mainSvc:  activityQRCodeSvc.NewMainSvc(),
		frontSvc: activityQRCodeSvc.NewFrontSvc(),
	}
}

func (h Handle) ActivitySave(c *gin.Context) {
	var req activityQRCodeReq.ActivitySaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ActivitySave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) ActivityDetail(c *gin.Context) {
	var req activityQRCodeReq.ActivityInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.ActivityDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) ActivityList(c *gin.Context) {
	var req activityQRCodeReq.ActivityListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mainSvc.ActivityList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) ActivityCategoryList(c *gin.Context) {
	var req activityQRCodeReq.ActivityCategoryListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.mainSvc.ActivityCategoryList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) ActivityDelete(c *gin.Context) {
	var req activityQRCodeReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.ActivityDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) GetQrCodeByUnionid(c *gin.Context) {
	var req activityQRCodeReq.GetQrCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if len(req.Unionid) > 32 {
		//处理unionid解密
		key := "hN9fX4jQ2iX2sL6b"
		iv := "pH4vK3hP6uH5kG4c"
		unionIdStr, unionErr := base64.StdEncoding.DecodeString(strings.ReplaceAll(req.Unionid, " ", "+"))
		if unionErr != nil {
			h.ResError(c, errors.New("unionid解密失败"))
			return
		}
		req.Unionid = string(common.AES_CBC_Decrypt(unionIdStr, []byte(key), []byte(iv)))

	}
	list, err := h.frontSvc.GetQrCodeByUnionid(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) GetAllStoreList(c *gin.Context) {
	var req activityQRCodeReq.GetQrCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.frontSvc.GetAllStoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h Handle) WelcomeSave(c *gin.Context) {
	var req activityQRCodeReq.WelcomeInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.WelcomeSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) WelcomeDetail(c *gin.Context) {
	var req activityQRCodeReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.WelcomeDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h Handle) WelcomeList(c *gin.Context) {
	var req activityQRCodeReq.WelcomeListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mainSvc.WelcomeList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h Handle) WelcomeDelete(c *gin.Context) {
	var req activityQRCodeReq.InfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mainSvc.WelcomeDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h Handle) Statistics(c *gin.Context) {
	var req activityQRCodeReq.StatisticsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.mainSvc.Statistics(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
