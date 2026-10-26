package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/drainage_contact"
)

// DrainageContact  联系我活动
func DrainageContact(groupRouter *gin.RouterGroup) {
	contactGroupHandle := drainageContactHandle.NewHandle()
	// 活动保存
	groupRouter.POST("/save", contactGroupHandle.GroupSave)
	//活动列表
	groupRouter.POST("/list", contactGroupHandle.GroupList)
	// 活动详情
	groupRouter.POST("/detail", contactGroupHandle.GroupDetail)
	// 活动子列表
	groupRouter.POST("/detail-list", contactGroupHandle.GroupDetailList)
	//活动删除
	groupRouter.POST("/delete", contactGroupHandle.GroupDelete)
	//下载zip
	groupRouter.POST("/zip-url", contactGroupHandle.ZipUrl)
}
