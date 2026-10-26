package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/tags"
)

func Tags(tagGroupRouter *gin.RouterGroup) {
	// 印象标签
	impressionGroup := tagGroupRouter.Group("impression")
	{
		impressionTagHandle := tags.NewImpressionHandle()
		impressionGroup.POST("/group/list", impressionTagHandle.GroupList)
		//impressionGroup.POST("/group/create", impressionTagHandle.GroupCreate)
		impressionGroup.POST("/tag/list", impressionTagHandle.TagList)
		//impressionGroup.POST("/tag/update", impressionTagHandle.TagUpdate)
		//impressionGroup.POST("/tag/create", impressionTagHandle.TagCreate)
		// 客户印象标签分组和标签列表
		impressionGroup.POST("/user-tags-list", impressionTagHandle.UserTagsList)
		// 印象标签列表
		impressionGroup.POST("/get-list", impressionTagHandle.GetList)

	}
	// 客户标签
	workTagsGroup := tagGroupRouter.Group("detail")
	{
		workTagsHandle := tags.NewWorkTagHandle()
		// 客户分组的企微标签用户
		workTagsGroup.POST("/list", workTagsHandle.List)
		// 客户标签分组和标签列表
		workTagsGroup.POST("/user-tags", workTagsHandle.UserTags)
		//标签是否具有创建权限
		workTagsGroup.POST("/auth", workTagsHandle.Auth)

		// 企微标签列表
		workTagsGroup.POST("/get-list", workTagsHandle.GetList)
	}
}
