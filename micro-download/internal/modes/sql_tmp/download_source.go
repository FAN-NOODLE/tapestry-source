package sql_tmp

import (
	"fmt"
	"micro-download/internal/modules/download/delivery"
)

// 插入
func SourceInsert() string {
	return `
	INSERT INTO download_resource (
		id,
		enterprise_id,
		store_id,
		operation_name,
		operation_id,
		type,
		name,
		path,
		source_type,
		source_data,
		status,
		failed_reason,
		expiration_at,
		created_at,
		complete_at,
		retry_at,
		deleted_at,
		label
	) VALUES (
		:id,
		:enterprise_id,
		:store_id,
		:operation_name,
		:operation_id,
		:type,
		:name,
		:path,
		:source_type,
		:source_data,
		:status,
		:failed_reason,
		:expiration_at,
		:created_at,
		:complete_at,
		:retry_at,
		:deleted_at,
		:label
	)
`
}

// SourceList 列表
func SourceList(field string, req delivery.ExportListReq) string {
	sql := `
		SELECT ` + field + ` FROM download_resource 
		WHERE deleted_at = 0 AND enterprise_id = ` + fmt.Sprintf("%d", req.EnterpriseID) + `
`
	// storeIds 查询
	storeIds := ""
	for _, storeId := range req.StoreIDS {
		storeIds += fmt.Sprintf("%d,", storeId)
	}
	if len(storeIds) > 0 {
		sql += `
		AND store_id in (` + storeIds[0:len(storeIds)-1] + `) 
`
	}
	// 名称查询
	if req.Filename != "" {
		sql += `
		AND	name LIKE '%` + req.Filename + `%' 
`
	}
	// 创建时间查询
	if req.CreatedAt.Max != 0 || req.CreatedAt.Min != 0 {
		if req.CreatedAt.Min != 0 {
			sql += `
		AND created_at >= ` + fmt.Sprintf("%d", req.CreatedAt.Min) + ` 
`
		}
		if req.CreatedAt.Max != 0 {
			sql += `
		AND created_at <= ` + fmt.Sprintf("%d", req.CreatedAt.Max) + ` 
`
		}
	}
	// 创建人ID查询
	if req.OperationId != 0 {

		sql += `
		AND operation_id = ` + fmt.Sprintf("%d", req.OperationId) + ` 
`
	}
	// 创建人名称模糊
	if req.OperationName != "" {
		sql += `
		AND operation_name like '` + req.OperationName + `%' 
`
	}

	// 状态查询
	if req.Status != 0 {
		sql += `
		AND status = ` + fmt.Sprintf("%d", req.Status) + `
`
	}

	// 默认ID倒序

	sql += `
	ORDER BY created_at DESC
`

	// 分页
	if field != "count(*)" {
		offset := req.PageSize * (req.Page - 1)
		sql += `
		LIMIT ` + fmt.Sprintf("%d", req.PageSize) + ` OFFSET ` + fmt.Sprintf("%d", offset)
	}
	return sql
}

func FullDetail() string {
	return `
		SELECT * FROM download_resource WHERE id = :id and deleted_at = 0
`
}

func LatestDetailByLabel() string {
	return `
		SELECT * FROM download_resource WHERE label = :label AND deleted_at = 0 AND status in (1,2,3,4) and operation_id = :operation_id ORDER BY created_at DESC LIMIT 1
`
}

func ChangeStatus() string {
	return `
		UPDATE download_resource SET status = :status WHERE id = :id
`
}

func DeleteById() string {
	return `
		UPDATE download_resource SET deleted_at = :time WHERE id = :id
`
}

func RetryById() string {
	return `
		UPDATE download_resource SET status = :status, retry_at = :retry_at WHERE id = :id
`
}

func Failed() string {
	return `
		UPDATE download_resource SET status = :status, failed_reason = :failed_reason WHERE id = :id
`
}

func Success() string {
	return `
		UPDATE download_resource SET status = :status, path = :path, complete_at = :complete_at, failed_reason = ""  WHERE id = :id
`
}

func ChangeAuto() string {
	return `
		UPDATE download_resource SET is_auto = :is_auto WHERE id = :id
`
}

func GetNoAutoOne() string {
	return `
		SELECT * FROM download_resource WHERE status=3 and is_auto=0 and operation_id = :operation_id AND deleted_at = 0 LIMIT 1
`
}
