package openReq

import (
	"micro-scrm/internal/modules/department"
)

type GenerateReserveRequest struct {
	EnterpriseId    int64  `json:"enterprise_id"`
	StoreCode       string `json:"store_code" binding:"required" comment:"门店编码"`
	ReserveNumberId int64  `json:"reserve_number_id" binding:"required" comment:"预约单Id"`
	ServiceName     string `json:"service_name" binding:"required" comment:"服务名称"`
	UnionId         string `json:"unionid" binding:"required" comment:"用户unionid"`
}

type AssignGuideRequest struct {
	EnterpriseId int64  `json:"enterprise_id"`
	StoreCode    string `json:"store_code" binding:"required" required:"门店编码"`
	UnionId      string `json:"unionid" binding:"required" required:"用户unionid"`
	Type         int    `json:"type" comment:"二维码场景：0-门店导购码，1-服务预约"`
}

type GetRelationRequest struct {
	EnterpriseId    int64  `json:"enterprise_id"`
	ReserveNumberId int64  `json:"reserve_number_id" binding:"required" comment:"预约单Id"`
	StoreCode       string `json:"store_code" binding:"required" comment:"门店编码"`
}

type DepartmentListRequest struct {
	department.DepartmentListRequest
}

type UserReserveListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
	Status       int    `json:"status" comment:"邀约单状态"`
	ReserveTime  int64  `json:"reserve_time" comment:"到店时间"`
	UnionId      string `json:"unionid" comment:"用户的unionid"`
	Month        int    `json:"month" comment:"是否按月查询"`
}

type UserReserveListRes struct {
	ID              int64  `json:"id" comment:"预约单id"`
	CreatedAt       int64  `json:"created_at" comment:"创建时间"`
	ReserveNumber   string `json:"reserve_number" comment:"预约单的编号"`
	ReserveTimeStr  string `json:"reserve_time_str" comment:"到店时间"`
	StoreName       string `json:"store_name" comment:"邀约门店"`
	StoreAddress    string `json:"store_address" comment:"门店地址"`
	Status          int32  `json:"status" comment:"邀约状态 0进行中,1邀约成功,2邀约完成,3已取消,4已过期"`
	UserName        string `json:"user_name" comment:"顾客姓名"`
	UserUnionid     string `json:"user_unionid" comment:"顾客unionid"`
	MemberID        int64  `json:"member_id" comment:"邀约导购ID"`
	MemberName      string `json:"member_name" comment:"邀约导购"`
	ReserveRemark   string `json:"reserve_remark" comment:"导购备注"`
	UserRemark      string `json:"user_remark" comment:"顾客备注"`
	Remark          string `json:"remark" comment:"评价内容"`
	OrderOrigin     int32  `json:"order_origin" comment:"服务类型 1为到点邀约单，2为邀请函"`
	ServiceDuration int    `json:"service_duration" comment:"服务时长"`
}
