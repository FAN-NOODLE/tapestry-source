package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/callback"
)

func Callback(callbackGroup *gin.RouterGroup) {
	callbackHandle := callback.NewCallbackHandle()
	// 回调处理
	callbackGroup.GET("/callback", callbackHandle.Handle)
	callbackGroup.POST("/callback", callbackHandle.Handle)

	contactsHandle := callback.NewContactsHandle()
	callbackGroup.POST("/contacts", contactsHandle.AutoHandle)
	callbackGroup.POST("/contacts/manual", contactsHandle.ManualHandle)
}
