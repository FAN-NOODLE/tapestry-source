package drainagehttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/drainage"
	response "micro-scrm/internal/pkg/resoponse"
)

// 外部条件数据欢迎语

type externalDataHandle struct {
	Name string
}

func NewMemberLevelHandle() *externalDataHandle {
	return &externalDataHandle{
		Name: "外部条件数据欢迎语",
	}
}

// 获取外部条件数据
// 1、目前每天从DL同步一次卡等信息，入库
// 2、数据信息缓存中有效存储1天
func (handle externalDataHandle) externalDataList(c *gin.Context) {

	// err := drainage.UpdateEnterpriseKdIinfo(ctx.Wrap(c), )
	// if err != nil {
	//	ginplus.ResError(c, err)
	//	return
	// }
	// ginplus.ResSuccess(c, true)
}

/**
 * 根据typeId查找相应type详细信息
 */
func (handle externalDataHandle) externalTypeById(ctx *gin.Context) {
	var TypeRequest drainage.ExternalTypeIdRequest
	if err := ctx.ShouldBindJSON(&TypeRequest); err != nil {
		ginplus.ResError(ctx, err)
		return
	}

	typeId, err := drainage.GetTypeId(TypeRequest.EnterpriseId)
	if err != nil {
		ginplus.ResError(ctx, errors.New("尚未接入会员系统"))
		return
	}

	TypeRequest.TypeId = typeId

	res, err := drainage.GetDataByTypeId(TypeRequest)
	if err != nil {
		ginplus.ResError(ctx, err)
		return
	}
	ginplus.ResSuccess(ctx, res)
}

/**
 * 获取所有的欢迎语类型type
 */
func (handle externalDataHandle) externalWelcomeTypeList(ctx *gin.Context) {
	var TypeAll drainage.ExternalTypeAllRequest
	if err := ctx.ShouldBindJSON(&TypeAll); err != nil {
		ginplus.ResError(ctx, err)
		return
	}

	res, err := drainage.GetAllType(TypeAll)
	if err != nil {
		ginplus.ResError(ctx, err)
		return
	}
	ginplus.ResSuccess(ctx, res)
}

func (handle externalDataHandle) welcomeDetailNew(c *gin.Context) {
	var welcomeDetail drainage.WelcomeDetailNewRequest
	err := c.ShouldBindJSON(&welcomeDetail)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	data, DetailErr := drainage.WelcomeDetail(welcomeDetail)
	if DetailErr != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

func (handle externalDataHandle) welcomeDel(c *gin.Context) {
	var welcomeDel drainage.WelcomeDelRequest
	err := c.ShouldBindJSON(&welcomeDel)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if err := drainage.WelcomeDel(welcomeDel); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

/**
 *新建&编辑欢迎语
 */
func (handle externalDataHandle) welcomeNew(c *gin.Context) {
	var welcomeRequest drainage.ExternalWelcomeRequest
	err := c.ShouldBindJSON(&welcomeRequest)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	res, err := drainage.WelcomeNewEdit(ctx.Wrap(c), welcomeRequest)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// 获取欢迎语的列表
func (handle externalDataHandle) welcomeList(c *gin.Context) {
	var welcomeAllRequest drainage.WelcomeAllRequest
	if err := c.ShouldBindJSON(&welcomeAllRequest); err != nil {
		ginplus.ResError(c, err)
		return
	}

	data, total, err := drainage.GetAllWelcome(welcomeAllRequest)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	response.SuccessList(c, total, welcomeAllRequest.Page, welcomeAllRequest.Size, data)
}
