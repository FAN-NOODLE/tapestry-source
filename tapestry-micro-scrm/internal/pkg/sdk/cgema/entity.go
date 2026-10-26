package cgema

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SendResultRequest struct {
	TaskID       string            `json:"task_id"`
	CrowdID      string            `json:"crowd_id"`
	TaskCategory int               `json:"task_category"`
	CrowdResult  []CrowdResultItem `json:"crowd_result"`
}

type CrowdResultItem struct {
	Date                string `json:"date"`
	SaName              string `json:"sa_name"`
	SaCode              string `json:"sa_code"`
	FlowID              string `json:"flowId"`
	Campaign            string `json:"campaign"`
	CTaskRemark         string `json:"c_task_remark"`
	CTaskClassification string `json:"c_task_classification"`
	UnionID             string `json:"unionid"`
	Result              string `json:"result"`
	ResultMessage       string `json:"result_message"`
	TaskRemark          string `json:"task_remark"`
}

type SendFinishRequest struct {
	Events      []EventItem `json:"events"`
	CommonProps CommonProps `json:"commonProps"`
}

type EventItem struct {
	IdentityType        string `json:"identityType"`
	IdentityValue       string `json:"identityValue"`
	Date                string `json:"date"`
	CResult             string `json:"c_result"`
	CSaName             string `json:"c_sa_name"`
	CSaCode             string `json:"c_sa_code"`
	CTaskID             string `json:"c_task_id"`
	CTaskName           string `json:"c_task_name"`
	CCrowdID            string `json:"c_crowd_id"`
	FlowID              string `json:"flowId"`
	Campaign            string `json:"campaign"`
	CTaskRemark         string `json:"c_task_remark"`
	CTaskClassification string `json:"c_task_classification"`
}

type CommonProps struct {
	Event string `json:"event"`
}

type SendMomentResultRequest struct {
	SaCode              string   `json:"sa_code"`
	TaskID              string   `json:"task_id"`
	CrowdID             string   `json:"crowd_id"`
	Campaign            string   `json:"campaign"`
	CustomerUnionIDList []string `json:"customer_unionid_list"`
}

type SendMomentCommentsRequest struct {
	SaCode      string              `json:"sa_code"`
	TaskID      string              `json:"task_id"`
	CrowdID     string              `json:"crowd_id"`
	Campaign    string              `json:"campaign"`
	CommentList []MomentCommentItem `json:"comment_list"`
	LikeList    []MomentCommentItem `json:"like_list"`
}

type MomentCommentItem struct {
	UnionID    string `json:"unionid"`
	CreateTime int64  `json:"create_time"`
}

type RelationEventRequest struct {
	EventType    string           `json:"event_type"`
	EventTime    string           `json:"event_time"`
	Customer     RelationCustomer `json:"customer"`
	Sa           RelationSa       `json:"sa"`
	Relation     RelationInfo     `json:"relation"`
	UnbindReason string           `json:"unbind_reason,omitempty"`
}

type RelationCustomer struct {
	UnionID        string `json:"unionid"`
	ExternalUserID string `json:"external_userid,omitempty"`
}

type RelationSa struct {
	UserID    string `json:"userid"`
	StoreCode string `json:"store_code,omitempty"`
}

type RelationInfo struct {
	FirstBindTime  string `json:"first_bind_time,omitempty"`
	LatestBindTime string `json:"latest_bind_time,omitempty"`
	UnbindTime     string `json:"unbind_time,omitempty"`
}
