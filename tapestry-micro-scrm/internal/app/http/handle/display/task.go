package display

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	display2 "micro-scrm/internal/app/http/request/display"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
	"micro-scrm/internal/app/service/department"
	"micro-scrm/internal/app/service/display"
	"micro-scrm/internal/pkg/helper"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
)

type TaskHandle struct {
	handle.BaseHandle
	taskSvc  *display.TaskSvc
	allotSvc *display.AllotSvc
}

func NewTaskHandle() TaskHandle {
	return TaskHandle{
		taskSvc:  display.NewTaskSvc(),
		allotSvc: display.NewAllotSvc(),
	}
}

// AddTask 任务添加
func (h TaskHandle) AddTask(c *gin.Context) {
	var req display2.AddTaskReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)

	_, importData, _, err := h.taskSvc.ReadReviewFile(cc, req.EnterpriseId, req.ReviewFile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if req.ReviewType == microScrmModel.VMTaskReviewTypeFile && req.ReviewFile == "" {
		h.ResError(c, errors.New("定向分配的文件不能为空"))
		return
	}
	taskId, err := h.taskSvc.CreateTask(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.allotSvc.CreateAllotInfo(cc, req, taskId, importData)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// List 任务列表
func (h TaskHandle) List(c *gin.Context) {
	var req display2.SearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	total, res, err := h.taskSvc.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, res)
}

// Detail 任务详情
func (h TaskHandle) Detail(c *gin.Context) {
	var req display2.TaskBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.taskSvc.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

// Copy 任务复制
func (h TaskHandle) Copy(c *gin.Context) {
	var req display2.TaskBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	copyId, err := h.taskSvc.Copy(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.taskSvc.Detail(cc, display2.TaskBaseReq{
		ID:                copyId,
		EnterpriseBaseReq: req.EnterpriseBaseReq,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	_, importData, _, err := h.taskSvc.ReadReviewFile(cc, req.EnterpriseId, res.ReviewFile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.allotSvc.CreateAllotInfo(cc, display2.AddTaskReq{
		EnterpriseBaseReq: req.EnterpriseBaseReq,
		TaskInfoRes:       res,
	}, copyId, importData)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// UpdateShow 任务上下架
func (h TaskHandle) UpdateShow(c *gin.Context) {
	var req display2.TaskBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.taskSvc.UpdateShow(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// UpdateTask 编辑任务
func (h TaskHandle) UpdateTask(c *gin.Context) {
	var req display2.AddTaskReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	if req.ReviewType == microScrmModel.VMTaskReviewTypeFile && req.ReviewFile == "" {
		h.ResError(c, errors.New("定向分配的文件不能为空"))
		return
	}
	_, importData, _, err := h.taskSvc.ReadReviewFile(cc, req.EnterpriseId, req.ReviewFile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.taskSvc.UpdateTask(cc, req)
	if err != nil {
		h.ResError(c, err)
	}
	res, err := h.taskSvc.Detail(cc, display2.TaskBaseReq{
		ID:                req.ID,
		EnterpriseBaseReq: req.EnterpriseBaseReq,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.allotSvc.CreateAllotInfo(cc, display2.AddTaskReq{
		EnterpriseBaseReq: req.EnterpriseBaseReq,
		TaskInfoRes:       res,
	}, req.ID, importData)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// Save 导购保存草稿|提交任务数据
func (h TaskHandle) Save(c *gin.Context) {
	var req display2.StoreTaskSaveInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.SaveTaskInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// TaskInfo H5端导购获取任务详情
func (h TaskHandle) TaskInfo(c *gin.Context) {
	var req display2.MemberStoreTaskInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := h.taskSvc.Detail(h.WarpContext(c), req.TaskBaseReq)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, detail)
}

// TaskList H5端导购获取任务列表
func (h TaskHandle) TaskList(c *gin.Context) {
	var req display2.MemberTaskListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.allotSvc.GetTaskList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, 0, req.Page, req.Size, list)
}

// Look 导购查看备注的上报
func (h TaskHandle) Look(c *gin.Context) {
	var req display2.LookReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.CreateLookLog(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// Auth VM身份信息
func (h TaskHandle) Auth(c *gin.Context) {
	var req display2.AuthReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result := h.allotSvc.AuthInfo(h.WarpContext(c), req)
	h.ResSuccess(c, result)
}

// VmList 获取VM列表
func (h TaskHandle) VmList(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.allotSvc.VmList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// StoreList 预览他店，获取门店列表
func (h TaskHandle) StoreList(c *gin.Context) {
	var req display2.StoreSearchReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.allotSvc.StoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// ListTips 获取导购列表的任务状态数
func (h TaskHandle) ListTips(c *gin.Context) {
	var req display2.ListTipsReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.allotSvc.GetListTips(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// CheckList VM的任务审核列表
func (h TaskHandle) CheckList(c *gin.Context) {
	var req display2.MemberTaskListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.allotSvc.GetCheckList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, 0, req.Page, req.Size, list)
}

// CheckStoreList 核销的门店列表
func (h TaskHandle) CheckStoreList(c *gin.Context) {
	var req display2.StoreListInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.allotSvc.GetCheckStoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) StoreTaskInfo(c *gin.Context) {
	var req display2.MemberStoreTaskInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.allotSvc.GetStoreTaskInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

// VmCheck VM审核门店任务
func (h TaskHandle) VmCheck(c *gin.Context) {
	var req display2.StoreTaskReviewReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.VmCheckStoreTask(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) SendNotice(c *gin.Context) {
	var req display2.SendNoticeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.VmSendNotice(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) AllotList(c *gin.Context) {
	var req display2.SearchAllotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	total, list, err := h.allotSvc.GetAllotList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h TaskHandle) AllotNum(c *gin.Context) {
	var req display2.SearchAllotBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.allotSvc.GetAllotNum(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) AllotListExport(c *gin.Context) {
	var req display2.SearchAllotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	total, list, err := h.allotSvc.GetAllotListExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h TaskHandle) NoticeSetting(c *gin.Context) {
	var req display2.NoticeSettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.SetNoticeSetting(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) AllotInfo(c *gin.Context) {
	var req display2.MemberStoreTaskInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.allotSvc.GetStoreTaskInfo(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h TaskHandle) AllotCheck(c *gin.Context) {
	var req display2.StoreTaskReviewWebReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.AllotCheck(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) NoticeSettingInfo(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.allotSvc.GetNoticeSetting(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{
		"value": info,
	})
}

func (h TaskHandle) AllotSendNotice(c *gin.Context) {
	var req display2.SearchAllotBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.AllotSendNotice(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) ImageList(c *gin.Context) {
	var req display2.SaveImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.allotSvc.AllotImageList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) AllotNoticeNum(c *gin.Context) {
	var req display2.SearchAllotBaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	total, err := h.allotSvc.AllotNoticeNum(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"num": total})
}

func (h TaskHandle) CosAvatar(c *gin.Context) {
	var req struct {
		Url string `json:"url" binding:"required" required:"图片地址必填"`
	}
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	cosFileAddr, err := h.allotSvc.CosAvatar(h.WarpContext(c), req.Url)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, gin.H{
		"cos_avatar": cosFileAddr,
	})

}

func (h TaskHandle) SaveImage(c *gin.Context) {
	var req display2.SaveStoreImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	id, err := h.allotSvc.SaveStoreImage(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"id": id})
}

func (h TaskHandle) DelImage(c *gin.Context) {
	var req display2.DelImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.DelImage(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) SaveComment(c *gin.Context) {
	var req display2.SaveStoreCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	id, err := h.allotSvc.SaveStoreComment(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{"id": id})
}

func (h TaskHandle) DelComment(c *gin.Context) {
	var req display2.DelImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.DelComment(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h TaskHandle) AllList(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.taskSvc.GetAllAllotList(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) DashBoard(c *gin.Context) {
	var req display2.DashBoard
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	dashboard, err := h.allotSvc.VmDashboard(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, dashboard)
}

func (h TaskHandle) FailList(c *gin.Context) {
	var req display2.MemberTaskReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.allotSvc.FailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) DashBoardList(c *gin.Context) {
	var req display2.DashBoardListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	result := make([]display2.DashBoardListRes, 0)
	list, err := dashboardSvc.NewAreaService().GetThirdIdsByAreaList(cc, req.ManagerAuthAreaIds)
	if err != nil {
		h.ResError(c, err)
		return
	}
	storeCount, err := h.allotSvc.GetDashBoardCount(cc, req.EnterpriseId, req.Ids, req.FinishStatus)
	if err != nil {
		h.ResError(c, err)
		return
	}
	storeNameMap, _ := department.NewDepartmentsService().GetStoreNameMap(cc, req.EnterpriseId)

	for _, info := range list {
		value := display2.DashBoardListRes{
			AreaInfo: info,
			StoreNum: 0,
			TaskNum:  0,
		}
		storeCode := info.Name
		if info.Type == dashboardModel.AreaStore && len(info.ChildrenThirdIds) > 0 { // 如果是门店的话，那么门店的名称需要调整 名称_storeCode
			if v, ok := storeNameMap[info.ChildrenThirdIds[0]]; ok {
				value.Name = fmt.Sprintf("%s_%s", v, storeCode)
			}
		}
		for _, item := range storeCount {
			if helper.InArray(item.StoreId, info.ChildrenThirdIds) && item.Num > 0 {
				value.StoreNum++
				value.TaskNum += item.Num
			}
		}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	h.ResSuccess(c, result)
}

func (h TaskHandle) StoreTaskList(c *gin.Context) {
	var req display2.StoreListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.taskSvc.GetStoreTaskList(h.WarpContext(c), req.EnterpriseId, req.StoreId, req.FinishStatus, req.Ids)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h TaskHandle) ReviewFileCheck(c *gin.Context) {
	var req display2.ReviewFileReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.taskSvc.CheckReviewFile(h.WarpContext(c), req.EnterpriseId, req.FileUrl)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h TaskHandle) SpecialSave(c *gin.Context) {
	var req display2.StoreTaskReviewReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.allotSvc.SpecialSaveStoreTask(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

// ReviewExport 下载分享详情数据
func (h TaskHandle) ReviewExport(c *gin.Context) {
	var req display2.TaskBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	fileName, err := h.allotSvc.ReviewExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	url, err := helper.UploadFileToCos(c, fileName, false)
	if err != nil {
		h.ResError(c, err)
		return
	}
	defer os.Remove(fileName)
	h.ResSuccess(c, gin.H{"url": url})
}
