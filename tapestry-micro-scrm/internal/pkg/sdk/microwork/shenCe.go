package microwork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"strconv"
)

type ShenCe struct {
}

func NewShenCe() ShenCe {
	return ShenCe{}
}

// External 顾客
func (s ShenCe) External(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "external",
		ItemId:       itemId,
		Properties:   properties,
	})
}

// ExternalPay 企微对外收款订单
func (s ShenCe) ExternalPay(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "externalpay",
		ItemId:       itemId,
		Properties:   properties,
	})
}

// ReserveService 邀约服务设置
func (s ShenCe) ReserveService(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "reserve_service",
		ItemId:       itemId,
		Properties:   properties,
	})
}

// Task 任务信息维度
func (s ShenCe) Task(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "task",
		ItemId:       itemId,
		Properties:   properties,
	})
}

// Material 素材
func (s ShenCe) Material(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "material",
		ItemId:       itemId,
		Properties:   properties,
	})
}

// InviteTemplate 邀请函模板
func (s ShenCe) InviteTemplate(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.itemSet(c, ItemSetReq{
		EnterpriseId: enterpriseId,
		ItemType:     "invite_template",
		ItemId:       itemId,
		Properties:   properties,
	})
}

func (s ShenCe) itemSet(c ctx.Context, request ItemSetReq) (bool, error) {
	var response struct {
		ResponseBase
		Data struct {
			Result bool `json:"result"`
		} `json:"data"`
	}
	err := httpClient.Post(c, GetURL("%s/service/v1/shenCe/itemSet"), request, &response)
	if err != nil {
		return response.Data.Result, err
	}

	return response.Data.Result, nil
}

func (s ShenCe) ProfileSet(c ctx.Context, request ProfileSetReq) (bool, error) {
	var response struct {
		ResponseBase
		Data struct {
			Result bool `json:"result"`
		} `json:"data"`
	}
	err := httpClient.Post(c, GetURL("%s/service/v1/shenCe/profileSet"), request, &response)
	if err != nil {
		return response.Data.Result, err
	}
	return response.Data.Result, nil
}

func (s ShenCe) track(c ctx.Context, request TrackReq) (bool, error) {
	var response struct {
		ResponseBase
		Data struct {
			Result bool `json:"result"`
		} `json:"data"`
	}
	err := httpClient.Post(c, GetURL("%s/service/v1/shenCe/track"), request, &response)
	if err != nil {
		return response.Data.Result, err
	}
	return response.Data.Result, nil
}

// TaskReceive 收到任务
// profiles.task_id 任务ID
// profiles.task_name 任务名称
// profiles.task_coverage 任务覆盖顾客数
// profiles.task_starttime 任务发布时间
func (s ShenCe) TaskReceive(c ctx.Context, enterpriseId int64, itemId string, profiles any) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "TaskReceive",
		DistinctId:   itemId,
		Profiles:     profiles,
	})
}

// TaskComplete 完成任务
func (s ShenCe) TaskComplete(c ctx.Context, enterpriseId int64, itemId string, taskId int64) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "TaskComplete",
		DistinctId:   itemId,
		Profiles: map[string]any{
			"task_id": strconv.FormatInt(taskId, 10),
		},
	})
}

// TaskExpire 任务过期
func (s ShenCe) TaskExpire(c ctx.Context, enterpriseId int64, itemId string, taskId int64) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "TaskExpire",
		DistinctId:   itemId,
		Profiles: map[string]any{
			"task_id": strconv.FormatInt(taskId, 10),
		},
	})
}

// AssetReceive 收到素材
func (s ShenCe) AssetReceive(c ctx.Context, enterpriseId int64, itemId string, sourceId int64) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "AssetReceive",
		DistinctId:   itemId,
		Profiles: map[string]any{
			"asset_id": strconv.FormatInt(sourceId, 10),
		},
	})
}

// InviteSuccess 邀约成功
func (s ShenCe) InviteSuccess(c ctx.Context, enterpriseId int64, itemId string, properties any) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "AppointmentComplete",
		DistinctId:   itemId,
		Profiles:     properties,
	})
}

// AssignCustomer 分配顾客（SM）
func (s ShenCe) AssignCustomer(c ctx.Context, enterpriseId int64, itemId string, profiles any) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "AssignCustomer",
		DistinctId:   itemId,
		Profiles:     profiles,
	})
}

// CustomerAssigned 顾客被分配（SA）
func (s ShenCe) CustomerAssigned(c ctx.Context, enterpriseId int64, itemId string, profiles any) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "CustomerAssigned",
		DistinctId:   itemId,
		Profiles:     profiles,
	})
}

// CompleteDisplayTask VM任务完成
func (s ShenCe) CompleteDisplayTask(c ctx.Context, enterpriseId int64, itemId, title, resSubmit string) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "CompleteDisplayTask",
		DistinctId:   itemId,
		Profiles: map[string]any{
			"app_name":          "WeStoreVM_sa",
			"dispalytask_title": title,
			"is_resubmit":       resSubmit,
		},
	})
}

// AddWecomSA  添加活动好友
func (s ShenCe) AddWecomSA(c ctx.Context, enterpriseId int64, itemId string, profiles any) (bool, error) {
	return s.track(c, TrackReq{
		EnterpriseId: enterpriseId,
		EventName:    "AddWecomSA",
		DistinctId:   itemId,
		Profiles:     profiles,
	})
}
