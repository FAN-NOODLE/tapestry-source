package source

type FriendCreateReq struct {
	EnterpriseID int64   `json:"enterprise_id" comment:"企业ID"`
	AdminId      int64   `json:"admin_id"`
	ShowType     int     `json:"show_type" comment:"展示类型"`
	FileType     int     `json:"file_type" comment:"文件类型 1 图片 2 视频"`
	SourceType   int     `json:"source_type" comment:"资源类型"`
	GroupId      int     `json:"group_id" comment:"群组id"`
	CategoryIds  []int64 `json:"category_ids" comment:"分类id"`
	StartTime    int     `json:"start_time"`
	EndTime      int     `json:"end_time"`
	TextContent  string  `json:"text_content" comment:"文本内容"`
	Sources      []struct {
		Id            int64        `json:"id"`
		Url           string       `json:"url"`
		VideoCoverUrl string       `json:"video_cover_url" comment:"封面图"`
		IsOrigin      int          `json:"is_origin"`
		Name          string       `json:"name"`
		Ext           string       `json:"ext"`
		Size          int          `json:"size"`
		IsTranscode   int          `json:"is_transcode"`
		RelationGoods []BrandGoods `json:"relation_goods"`
	} `json:"sources"`
}

type FriendDetailReq struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type FriendUpdateReq struct {
	ID int64 `json:"id"`
	FriendCreateReq
}
