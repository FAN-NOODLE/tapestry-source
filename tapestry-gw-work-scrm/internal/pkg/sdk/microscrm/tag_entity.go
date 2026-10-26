package microScrm

type TagRequestList struct {
	RequestBase
	RelationID int64 `json:"relation_id" validate:"required" required:"客户关系ID不能为空"`
}
