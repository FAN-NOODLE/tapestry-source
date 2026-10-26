package userGroupDao

import (
	"micro-scrm/internal/app/dao"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/pkg/db/gormx"
)

type GroupFile struct {
	dao       dao.Dao
	tableName string
}

func NewGroupFile() *GroupFile {
	return &GroupFile{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: userGroupModel.UserGroupFileModel{}.TableName(),
	}
}
