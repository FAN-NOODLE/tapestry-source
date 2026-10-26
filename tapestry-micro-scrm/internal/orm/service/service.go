package service

import "gorm.io/gorm"

// 增加查询条件scope
// 函数的开头为ScopeCond_

// PaginateParam ScopeCond_Paginate 分页器
// 分页参数
type PaginateParam struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

const (
	MaxPageSize = 99999999
	MinPageSize = 5
	DefaultPage = 1
)

func ScopeCond_Paginate(param PaginateParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 获取页码，如果没有设置则默认为第一页
		page := DefaultPage
		if param.Page != 0 {
			page = param.Page
		}

		pageSize := MinPageSize
		switch {
		case param.PageSize > MaxPageSize:
			pageSize = MaxPageSize
		case param.PageSize <= 0:
			pageSize = MinPageSize
		default:
			pageSize = param.PageSize
		}

		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

// ScopeCond_EQDelete 过滤删除掉的数据
func ScopeCond_EQDelete(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at = ?", 0)
}

func ScopeLeftCond_EQDelete(db *gorm.DB) *gorm.DB {
	return db.Where("drainage_welcomes.deleted_at = ?", 0)
}

// ScopeCond_EQDelete 过滤删除掉的数据
func ScopeCond_EQDeleteAndGroup(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at = ? group by plan_id", 0)
}
