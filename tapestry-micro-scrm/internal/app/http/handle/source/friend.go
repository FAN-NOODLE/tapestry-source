package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/service/source"
)

type FriendHandle struct {
	handle.BaseHandle
}

func NewFriendHandle() FriendHandle {
	return FriendHandle{}
}

func (h FriendHandle) Create(c *gin.Context) {
	var request sourceRequest.FriendCreateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := source.NewFriendService().Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
	return
}

func (h FriendHandle) Detail(c *gin.Context) {
	var request sourceRequest.FriendDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := source.NewFriendService().Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

func (h FriendHandle) Update(c *gin.Context) {
	var request sourceRequest.FriendUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := source.NewFriendService().Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
	return
}
