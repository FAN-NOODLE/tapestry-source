package departmentshttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gorm.io/gorm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/modules/department"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
	"net/http"
	"time"
)

var departmentService *department.Service

// listDept 组织列表
func listDept(c *gin.Context) {
	dept := &department.DepartmentListRequest{}
	err := c.BindJSON(dept)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	res := departmentService.List(ctx.Wrap(c), dept)
	ginplus.ResSuccess(c, res)
	return
}

// detailDept 组织详情
func detailDept(c *gin.Context) {
	dept := &department.DepartmentDetailRequest{}
	err := c.BindJSON(dept)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if dept.EnterpriseID == 0 {
		ginplus.ResError(c, errs.DepartmentIdIsEmpty)
		return
	}
	res := departmentService.Detail(ctx.Wrap(c), dept.EnterpriseID)
	ginplus.ResSuccess(c, res)
	return
}

func updateRecord(c *gin.Context) {
	dept := &department.DepartmentDetailRequest{}
	err := c.BindJSON(dept)
	if err != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if dept.EnterpriseID == 0 {
		ginplus.ResError(c, errs.DepartmentIdIsEmpty)
		return
	}
	res := departmentService.FindUpdateRecordRecent(ctx.Wrap(c), dept.EnterpriseID)
	ginplus.ResSuccess(c, res)
	return
}

// departmentMomentList 选择部门
func departmentMomentList(c *gin.Context) {
	request := &department.MomentListRequest{}
	jsonError := c.BindJSON(request)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	// 判断企业，coach计算零售和奥莱的逻辑，与ks&&sw的逻辑不一样
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx.Wrap(c), request.EnterpriseId, "store_name")
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}
	isCoach := false
	if config.Name == "coach" {
		isCoach = true
	}
	result, logicError := departmentService.DepartmentMomentList(ctx.Wrap(c), request)
	var res []department.MomentDepartment
	res = *result
	var outletsIdsArr, retailIdsArr []int64
	outletsInfo := make(map[int64]string, 0)
	retailInfo := make(map[int64]string, 0)
	outletsList, err := microWorkDao.NewStoreDepartment().GetMany(ctx.Wrap(c), map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"store_type":    microWorkDao.STORE_TYPE_OUTLETS,
	})
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
	}
	for _, v := range outletsList {
		item := v
		parseInt := int64(item.WechatDepartmentID)
		if err != nil {
			ginplus.ResError(c, err, http.StatusOK)
			return
		}
		outletsIdsArr = append(outletsIdsArr, parseInt)
		if isCoach {
			outletsInfo[parseInt] = item.StoreName
		}

	}
	retailList, err := microWorkDao.NewStoreDepartment().GetMany(ctx.Wrap(c), map[string]interface{}{
		"enterprise_id": request.EnterpriseId,
		"store_type":    microWorkDao.STORE_TYPE_RETAIL,
	})
	if err != nil {
		ginplus.ResError(c, err, http.StatusOK)
		return
	}

	for _, v := range retailList {
		item := v
		parseInt := int64(item.WechatDepartmentID)
		retailIdsArr = append(retailIdsArr, parseInt)
		if isCoach {
			retailInfo[parseInt] = item.StoreName
		}
	}

	if logicError != nil {
		ginplus.ResError(c, logicError, http.StatusOK)
	}
	if len(res) > 0 {
		for k, momentDepartment := range res {
			v := momentDepartment
			if helper.InArray(v.ThirdID, outletsIdsArr) { // 奥莱
				res[k].StoreType = 2
				if isCoach {
					if item, ok := outletsInfo[v.ThirdID]; ok {
						res[k].DpName = v.DpName + "-" + item
					}
				}
			}
			if helper.InArray(v.ThirdID, retailIdsArr) { // 零售
				res[k].StoreType = 1
				if isCoach {
					if item, ok := retailInfo[v.ThirdID]; ok {
						res[k].DpName = v.DpName + "-" + item
					}
				}
			}
		}
	}

	ginplus.ResSuccess(c, res)
}

// userLostConfig 导购丢失客户推送消息配置
func userLostConfig(c *gin.Context) {
	var request struct {
		EnterpriseId int64 `json:"enterprise_id"  binding:"required" comment:"企业Id"`
		DpId         int64 `json:"dp_id" binding:"required" comment:"部门Id"`
		UserList     []struct {
			UserId   string `json:"user_id" binding:"required" comment:"用户Id"`
			SendType uint32 `json:"send_type" binding:"required" comment:"推送类型"`
		} `json:"user_ids" binding:"required" comment:"推送消息用户集合"`
	}
	err := c.Bind(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var userIds []string
	for _, userInfo := range request.UserList {
		userIds = append(userIds, userInfo.UserId)
	}

	// 根据部门ID，获取所有下级部门id集合
	departmentIds := orm.Service.Departments.GetDepartmentsChildrenIds(request.EnterpriseId, request.DpId)

	// 删除已设置的员工信息
	gormx.GetInstance().
		Where("enterprise_id", request.EnterpriseId).
		Where("dp_id in (?)", departmentIds).
		Delete(&micro_scrm.CustomersDelUserMessageModel{})

	// 已存在的员工id和要被删除的员工id
	var existUsers, deleteUsers []string
	// 查询部门已添加的推送人
	gormx.GetInstance().Model(&micro_scrm.CustomersDelUserMessageModel{}).
		Where("enterprise_id", request.EnterpriseId).
		Where("dp_id", request.DpId).
		Pluck("user_id", &existUsers)
	for _, user := range existUsers {
		if !helper.InArray[string](user, userIds) {
			deleteUsers = append(deleteUsers, user)
		}
	}
	if len(deleteUsers) != 0 {
		gormx.GetInstance().Where("enterprise_id", request.EnterpriseId).
			Where("dp_id in (?)", departmentIds).
			Where("user_id in (?)", deleteUsers).
			Delete(&micro_scrm.CustomersDelUserMessageModel{})
	}

	if len(request.UserList) == 0 {
		ginplus.ResOK(c)
		return
	}

	var customersDelUserMessages []micro_scrm.CustomersDelUserMessageModel
	// 组装部门和被推送人的配置数据
	for _, departmentId := range departmentIds {
		for _, userInfo := range request.UserList {
			customersDelUserData := micro_scrm.CustomersDelUserMessageModel{
				ID:           uid.ID(ctx.Wrap(c)),
				DpID:         departmentId,
				UserID:       userInfo.UserId,
				EnterpriseID: request.EnterpriseId,
				SendType:     userInfo.SendType,
				CreatedAt:    int(time.Now().Unix()),
				UpdatedAt:    int(time.Now().Unix()),
			}
			customersDelUserMessages = append(customersDelUserMessages, customersDelUserData)
		}
	}

	err = gormx.GetInstance().Model(&micro_scrm.CustomersDelUserMessageModel{}).Create(&customersDelUserMessages).Error
	if err != nil {
		ginplus.ResError(c, errors.New("添加配置失败,err:"+err.Error()))
		return
	}

	ginplus.ResOK(c)
	return
}

// userLostInfo 导购丢失客户推送消息详情
func userLostInfo(c *gin.Context) {
	var request struct {
		DpId         int64 `json:"dp_id" comment:"部门Id"`
		EnterpriseId int64 `json:"enterprise_id"  binding:"required" comment:"企业Id"`
	}

	err := c.Bind(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	var responseData []struct {
		DpName     string `json:"dp_name"`
		MemberName string `json:"member_name"`
		SendType   uint   `json:"send_type"`
		UserId     string `json:"user_id"`
	}
	err = gormx.GetInstance().Model(micro_scrm.CustomersDelUserMessageModel{}).
		Select("customers_del_user_message.send_type,members.name as member_name,departments.name as dp_name,customers_del_user_message.user_id").
		Joins("left join members on members.userid = customers_del_user_message.user_id").
		Joins("left join departments on departments.id = customers_del_user_message.dp_id").
		Where("members.enterprise_id", request.EnterpriseId).
		Where("departments.enterprise_id", request.EnterpriseId).
		Where("customers_del_user_message.enterprise_id", request.EnterpriseId).
		Where("customers_del_user_message.dp_id", request.DpId).
		Find(&responseData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ginplus.ResError(c, errors.New("查询用户错误,err:"+err.Error()))
		return
	}

	ginplus.ResSuccess(c, responseData)
	return
}
