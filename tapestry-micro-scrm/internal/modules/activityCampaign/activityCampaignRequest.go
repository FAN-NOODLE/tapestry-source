package activityCampaign

type ActivitiesRequest struct {
	EnterpriseId int64 `json:"enterprise_id"  binding:"required" desc:"企业编号"`
	ActivityId   int64 `json:"activity_id" binding:"required" desc:"活动编号"`
}
type UsersignInRequest struct {
	EnterpriseId  int64  `json:"enterprise_id"  binding:"required" desc:"企业编号"`
	UnionId       string `json:"union_id"  binding:"required" desc:"用户union_id"`
	ActivityId    int64  `json:"activity_id" binding:"required" desc:"活动编号"`
	MediaId       string `json:"media_id" desc:"图片素材编号"`
	SignInPic     string `json:"sign_in_pic" desc:"签到图片"`
	SignInContent string `json:"sign_in_content" desc:"签到图片"`
	SignInDate    string `json:"sign_in_date" desc:"签到日期"`
}
