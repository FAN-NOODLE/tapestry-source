package marketing

type RecordConditionRequest struct {
	EnterpriseID   int64 `json:"enterprise_id" db:"enterprise_id"`
	Condition      int8  `json:"condition" db:"condition"`
	InnerCondition int8  `json:"inner_condition" db:"inner_condition"`
	IsFinished     int8  `json:"is_finished" db:"is_finished"`
	CurrentTime    int64 `json:"current_time" db:"current_time"`
}
