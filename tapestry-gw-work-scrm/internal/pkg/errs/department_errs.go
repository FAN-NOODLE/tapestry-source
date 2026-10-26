package errs

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

var (
	// DepartmentIdIsEmpty 组织ID不能为空
	DepartmentIdIsEmpty = errorx.NewErrorMsg(110001, "组织ID不能为空")
	// DepartmentNotFund 组织信息不存在
	DepartmentNotFund = errorx.NewErrorMsg(110002, "组织信息不存在")
	// SyncContactRecordFail 同步通讯录记录失败
	SyncContactRecordFail = errorx.NewErrorMsg(110003, "同步通讯录记录失败")
)
