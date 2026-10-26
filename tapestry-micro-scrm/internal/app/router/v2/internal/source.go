package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/source"
)

func Source(sourceGroup *gin.RouterGroup) {
	// 素材-产品类别
	categoryHandle := source.NewCategoryHandle()
	{
		sourceGroup.POST("/category/list", categoryHandle.List)
		sourceGroup.POST("/category/create", categoryHandle.Create)
		sourceGroup.POST("/category/update", categoryHandle.Update)
		sourceGroup.POST("/category/delete", categoryHandle.Delete)
		sourceGroup.POST("/category/source_relation", categoryHandle.SourceRelation)
		sourceGroup.POST("/remove-category", categoryHandle.RemoveCategory)
	}

	// 素材-文本
	textHandle := source.NewTextHandle()
	{
		sourceGroup.POST("/text/detail", textHandle.Detail)
		sourceGroup.POST("/text/update", textHandle.Update)
	}

	// 素材-图片
	imageHandle := source.NewImageHandle()
	{
		sourceGroup.POST("/image/detail", imageHandle.Detail)
		sourceGroup.POST("/image/update", imageHandle.Update)
	}

	// 素材-链接
	linkHandle := source.NewLinkHandle()
	{
		sourceGroup.POST("/link/detail", linkHandle.Detail)
		sourceGroup.POST("/link/update", linkHandle.Update)
	}

	// 素材-小程序
	miniProgramHandle := source.NewMiniProgramHandle()
	{
		sourceGroup.POST("/mini_program/detail", miniProgramHandle.Detail)
		sourceGroup.POST("/mini_program/update", miniProgramHandle.Update)
	}

	// 素材-视频
	videoHandle := source.NewVideoHandle()
	{
		sourceGroup.POST("/video/detail", videoHandle.Detail)
		sourceGroup.POST("/video/update", videoHandle.Update)
	}

	// 素材-文件
	fileHandle := source.NewFileHandle()
	{
		sourceGroup.POST("/file/detail", fileHandle.Detail)
		sourceGroup.POST("/file/update", fileHandle.Update)
	}

	// 素材-朋友圈
	friendHandle := source.NewFriendHandle()
	{
		sourceGroup.POST("/friend/create", friendHandle.Create)
		sourceGroup.POST("/friend/update", friendHandle.Update)
		sourceGroup.POST("/friend/detail", friendHandle.Detail)
	}

	sourceHandle := source.NewSourceHandle()
	{
		// 批量启用/禁用素材
		sourceGroup.POST("/updateStatus", sourceHandle.UpdateStatus)

		// 发送文本卡片消息
		// sourceGroup.POST("/sendTextMessage", sourceHandle.SendTextMessage)
		// 获取上传cos的凭证
		sourceGroup.POST("credential", sourceHandle.Credential)
		// 小程序端导购上传素材
		sourceGroup.POST("/upload-media", sourceHandle.UploadMedia)
		// 压缩结果的展示
		sourceGroup.POST("/compress", sourceHandle.Compress)
		// 素材能否获取企微的mediaid
		sourceGroup.POST("/get-media-info", sourceHandle.GetMediaInfo)
		// 压缩文件
		sourceGroup.POST("/compress-file", sourceHandle.CompressFile)
		// 修改cos视频创建时间
		sourceGroup.POST("/edit-video-create-time", sourceHandle.EditVideoTime)
		// 上报计次
		sourceGroup.POST("/report-times", sourceHandle.ReportTimes)
		// 验证是否用户首次打开素材交互提示
		sourceGroup.POST("/check-first", sourceHandle.CheckFirst)
		// 上报计次
		sourceGroup.POST("/make-cover", sourceHandle.MakeCover)
		// 素材库设置置顶
		sourceGroup.POST("/set-top", sourceHandle.SetTop)

	}

	// 素材分组归属
	groupRelationHandle := source.NewGroupRelationHandle()
	{
		// 根据分组id获取素材分组归属门店列表
		sourceGroup.POST("/group-relation/list", groupRelationHandle.List)
		// 根据门店获取素材分组列表
		sourceGroup.POST("/group-relation/get-list-for-store", groupRelationHandle.GetListForStore)
		// 素材分组归属门店保存
		sourceGroup.POST("/group-relation/save", groupRelationHandle.Save)
		sourceGroup.POST("/group-relation/save-all", groupRelationHandle.SaveAll)
		// 小程序端门店员工根据查看分组权限查询各素材列表
		sourceGroup.POST("/select-list", groupRelationHandle.SelectList)
		// 后台管理查询各素材列表
		sourceGroup.POST("/console-list", groupRelationHandle.ConsoleList)
		// 小程序端素材上新数量
		sourceGroup.POST("/get-new-num", groupRelationHandle.GetNewSourceNum)
		// 小程序端门店员工根据查看分组权限查询各素材列表(收藏）
		sourceGroup.POST("/collect-select-list", groupRelationHandle.CollectSelectList)
	}

	// 品牌素材库
	brandHandle := source.NewBrandHandle()
	{
		// 品牌素材库查询商品素材
		sourceGroup.POST("/brand/search-list", brandHandle.SearchList)
		// 根据upcCode获取商品素材缓存
		sourceGroup.POST("/brand/get-cache-info", brandHandle.CacheInfo)
		// 品牌素材库查询商品素材带收藏数
		sourceGroup.POST("/brand/collect/search-list", brandHandle.CollectSearchList)
		// 模糊查询skuItem
		sourceGroup.POST("/brand/search-list-group-sku", brandHandle.SearchListGroupSku)
		// Talkie图片搜索
		sourceGroup.POST("/brand/talkie-image-search", brandHandle.TalkieImageSearch)

	}

	collectHandler := source.NewCollectHandle()
	{
		// 收藏添加
		sourceGroup.POST("/collect/edit", collectHandler.Edit)
		// 收藏列表
		//sourceGroup.POST("/collect/list", collectHandler.List)
	}
}
