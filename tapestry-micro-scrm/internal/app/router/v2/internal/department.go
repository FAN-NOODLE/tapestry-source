package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/department"
)

func Department(dempartmentGroup *gin.RouterGroup) {

	departmentHandle := department.NewHandle()

	dempartmentGroup.POST("/sa-store-lists", departmentHandle.SaStoreList)

	dempartmentGroup.POST("/store-list", departmentHandle.StoreList)
}
