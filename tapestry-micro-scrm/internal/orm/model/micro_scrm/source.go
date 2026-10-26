package micro_scrm

import "micro-scrm/internal/orm/model/micro_scrm/source"

type SourceTree struct {
	source.GroupModel
	Level    int           `json:"level"` // 层级
	Ids      string        `json:"ids"`
	Children []*SourceTree `json:"children"`
}

type ReceiveGroupName struct {
	GroupName string `json:"group_name"`
}

type Ids struct {
	Id int64 `json:"id"`
}
