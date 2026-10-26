package sourceGroupHttp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	source2 "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/modules/source"
	"micro-scrm/internal/modules/sourceGroup"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	response "micro-scrm/internal/pkg/resoponse"
	"micro-scrm/internal/pkg/uid"
)

var s *sourceGroup.Service
var sou *source.Service

// list 列表
func list(c *gin.Context) {
	params := &sourceGroup.SourceGroup{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	list, err := s.List(ctx.Wrap(c), params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, list)
	return
}

// del 删除
func del(c *gin.Context) {
	params := &struct {
		IDs []int64 `json:"ids"  validate:"required" required:"素材分组IDs不能为空"`
	}{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	var rows int64
	var err error
	if len(params.IDs) == 1 {
		rows, err = s.Del(ctx.Wrap(c), params.IDs[0])
	} else {
		rows, err = s.BatchDel(ctx.Wrap(c), params.IDs)
	}
	// 将当前分组的素材更新成未分组
	sourceUpdate := &source.RequestGroupMove{
		GroupID: params.IDs[0],
	}
	_, err = sou.UpdateGroup(ctx.Wrap(c), sourceUpdate)

	if err != nil {
		response.Notice(c, errs.SourceGroupDelErr)
		return
	}

	response.Success(c, &sourceGroup.Rows{Rows: rows})
	return
}

// add 新增
func add(c *gin.Context) {
	params := &sourceGroup.SourceGroup{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)

	// 重名验证
	exists, err := s.IsRepeatName(cc, params)
	if err != nil {
		response.Notice(c, errs.SourceGroupNameRepeated)
		return
	}
	if exists {
		response.Notice(c, errs.SourceGroupNameRepeated)
		return
	}

	params.CreatedAt = assist.Timestamp()
	id, err := s.Add(cc, params)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &sourceGroup.ID{ID: id})
	return
}

// edit
func edit(c *gin.Context) {
	params := &struct {
		ID           int64  `json:"id" db:"id"  validate:"required" required:"素材分组ID不能为空"`
		EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"  validate:"required" required:"企业ID不能为空"`
		GroupName    string `json:"group_name" db:"group_name" validate:"required" required:"素材名称不能为空"`
	}{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if params.ID == 0 {
		response.Notice(c, errs.SourceGroupIDEmpty)
		return
	}
	cc := ctx.Wrap(c)

	sourceGroupParams := &sourceGroup.SourceGroup{
		ID:           params.ID,
		EnterpriseID: params.EnterpriseID,
		GroupName:    params.GroupName,
		UpdatedAt:    assist.Timestamp(),
	}
	// 重名验证
	exists, err := s.IsRepeatName(cc, sourceGroupParams)
	if err != nil {
		response.Notice(c, errs.SourceGroupIDEmpty)
		return
	}
	if exists {
		response.Notice(c, errs.SourceGroupNameRepeated)
		return
	}

	id, err := s.Edit(cc, sourceGroupParams)

	if errorIgnore.Ignore(err) != nil {
		response.Notice(c, errorx.NewError(errs.DefaultErrorCode, err), err)
		return
	}

	response.Success(c, &sourceGroup.ID{ID: id})
	return
}

// 新增 addSource
func addSource(c *gin.Context) {
	var addSource sourceGroup.NewSource
	if err := c.ShouldBindJSON(&addSource); err != nil {
		log.Info(ctx.Wrap(c)).Msg("addSource: Data parsing error!")
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	// 重名验证
	nameMap := sourceGroup.NameMap(addSource.EnterpriseID, addSource.ParentID)
	if len(nameMap) != 0 {
		if _, ok := nameMap[addSource.GroupName]; ok {
			ginplus.ResError(c, errs.SourceGroupNameRepeated)
			return
		}
	}

	ID := uid.ID(ctx.Wrap(c))
	// 一级分组默认全部部门可见；子分组继承直接父级当前权限。
	if err := source2.NewGroupRelationService().GroupRelationSaveInitial(ctx.Wrap(c), addSource.EnterpriseID, ID, addSource.ParentID); err != nil {
		ginplus.ResError(c, err)
		return
	}

	addSource.ID = ID
	if err := sourceGroup.AddSource(ctx.Wrap(c), addSource); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

//编辑 editSource
func editSource(c *gin.Context) {
	var editInfo sourceGroup.EditSource
	if err := c.ShouldBindJSON(&editInfo); err != nil {
		log.Info(ctx.Wrap(c)).Msg("editSource: Data parsing error!")
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	if err := sourceGroup.EditSourceInfo(editInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// UpdateIgnoreIndex 分组首页不可见
func updateIgnoreIndex(c *gin.Context) {
	var editInfo sourceGroup.UpdateIgnoreIndexReq
	if err := c.ShouldBindJSON(&editInfo); err != nil {
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	if err := sourceGroup.UpdateIgnoreIndex(editInfo); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// 删除 removeSource
func removeSource(c *gin.Context) {
	var delId sourceGroup.RemoveSource
	if err := c.BindJSON(&delId); err != nil {
		log.Info(ctx.Wrap(c)).Msg("removeSource : Data parsing error!")
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	res, _ := sourceGroup.DelVerify(delId.ID, delId.EnterpriseID)
	if len(res) > 1 {
		ginplus.ResError(c, errors.New("该组下有数据，该组不能删除！"))
		return
	}
	if err := sourceGroup.DelSource(delId.ID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// SourceTree 分组树生成
func SourceTree(c *gin.Context) {
	var findParams sourceGroup.FindSource
	if err := c.ShouldBind(&findParams); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if findParams.EnterpriseID == 0 {
		ginplus.ResError(c, errors.New("企业id不能为空！"))
		return
	}

	if result, err := sourceGroup.SourceTree(findParams.EnterpriseID, findParams.ID); err != nil {
		log.Info(ctx.Wrap(c)).Msg(fmt.Sprintf("SourceTree : err->%v", err))
		ginplus.ResError(c, err)
		return
	} else {
		ginplus.ResSuccess(c, result)
	}

}

// SourceMove 分组位置移动
func SourceMove(c *gin.Context) {
	var receive sourceGroup.MoveReceive
	if err := c.ShouldBindJSON(&receive); err != nil {
		log.Info(ctx.Wrap(c)).Msg("SourceMove : Data Move failed!")
		ginplus.ResError(c, errors.New("Data parsing error!"))
		return
	}
	if err := sourceGroup.MoveSource(receive); err != nil {
		log.Info(ctx.Wrap(c)).Msg("SourceMove : Data Move failed!")
		ginplus.ResError(c, errors.New("SourceMove : "+err.Error()))
		return
	}
	ginplus.ResOK(c)
}
