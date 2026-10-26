package source

import (
	"micro-scrm/internal/app/http/request"
)

type DetailRequest struct {
	SourceID     uint64 `json:"source_id" binding:"required" comment:"source_id"`
	EnterpriseID uint64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type Request struct {
	ID            uint64       `json:"id" binding:"required" comment:"id"`
	EnterpriseID  uint64       `json:"enterprise_id" binding:"required" comment:"企业ID"`
	StartTime     int64        `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime       int64        `json:"end_time" binding:"required" comment:"结束时间"`
	CategoryIds   []uint64     `json:"category_ids" binding:"required" comment:"产品类别"`
	SourceType    int32        `json:"source_type"`
	GroupId       int64        `json:"group_id"`
	Ids           []int64      `json:"ids" comment:"ids"`
	RelationGoods []BrandGoods `json:"relation_goods"`
}
type BrandGoods struct {
	Name    string `json:"name"`
	Sku     string `json:"sku"`
	Upc     string `json:"upc"`
	Url     string `json:"url"`
	Color   string `json:"color"`
	AddType int8   `json:"add_type" comment:"添加方式 1 手动添加 2 AI识别"`
	Score   string `json:"score" comment:"置信度（AI识别）"`
}

type TextUpdateRequest struct {
	Request
	Content string `json:"content" binding:"required" comment:"文本内容"`
}

type ImageUpdateRequest struct {
	Request
	Sources []MergeImage `json:"sources" binding:"required" required:"图片集合不能为空"`
}
type MergeImage struct {
	Id              int64        `json:"id"`
	Name            string       `json:"name"`
	URL             string       `json:"url"`
	Sort            int32        `json:"sort"`
	DownloadTimes   int          `json:"download_times"`
	ForwardingTimes int          `json:"forwarding_times"`
	IsCollect       int          `json:"is_collect"`
	CollectNum      int          `json:"collect_num"`
	RelationGoods   []BrandGoods `json:"relation_goods"`
}

type LinkUpdateRequest struct {
	Request
	Title       string `json:"title" binding:"required" comment:"标题"`
	CoverUrl    string `json:"cover_url" binding:"required" comment:"封面图"`
	Description string `json:"description" binding:"required" comment:"链接描述"`
	LinkType    int    `json:"link_type" binding:"required" comment:"链接类型"`
	LinkUrl     string `json:"link_url" binding:"required" comment:"跳转地址"`
}

type MiniProgramUpdateRequest struct {
	Request
	Title         string `json:"title" binding:"required" comment:"标题"`
	CoverUrl      string `json:"cover_url" binding:"required" comment:"封面图"`
	Path          string `json:"path" binding:"required" comment:"页面路径"`
	MiniProgramID int64  `json:"mini_program_id" binding:"required" comment:"小程序ID"`
}

type VideoUpdateRequest struct {
	Request
	Name          string       `json:"name" binding:"required" comment:"文件名"`
	RelationGoods []BrandGoods `json:"relation_goods" comment:"关联商品"`
}

type FileUpdateRequest struct {
	Request
	Name string `json:"name" binding:"required" comment:"文件名"`
}

type UpdateStatusRequest struct {
	IDs          []int64 `json:"ids" binding:"required" comment:"ids"`
	EnterpriseID int64   `json:"enterprise_id" binding:"required" comment:"企业ID"`
	SourceType   int     `json:"source_type" comment:"source_type"`
	SourceStatus int     `json:"source_status" comment:"source_status"`
}

type SendRequest struct {
	EnterpriseID  int64    `json:"enterprise_id" binding:"required" comment:"企业ID"`
	MemberUserIDs []string `json:"member_user_ids" `
}

// MessageMember 发送消息成员
type MessageMember struct {
	Department []string `json:"department"`
	Member     []string `json:"member"`
}

type SendMessageMemberGroup struct {
	MemberUserIDs []string `json:"member_user_ids" `
	CreateAt      int64    `json:"create_at"`
	IdCount       int64    `json:"id_count"`
}

type SendMessageThirdGroup struct {
	ThirdIds  []int `json:"third_ids" `
	StartTime int64 `json:"start_time"`
	IdCount   int64 `json:"id_count"`
}

type UploadMediaRequest struct {
	ID           uint64      `json:"id" comment:"id"`
	EnterpriseID uint64      `json:"enterprise_id" binding:"required" comment:"企业ID"`
	StartTime    int64       `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime      int64       `json:"end_time" binding:"required" comment:"结束时间"`
	CategoryIds  []uint64    `json:"category_ids" comment:"产品类别"`
	SourceType   int32       `json:"source_type"`
	GroupId      int64       `json:"group_id"`
	MemberId     int64       `json:"member_id"`
	UserId       string      `json:"userid"`
	Content      string      `json:"content" comment:"文本内容"`
	Sources      []MediaInfo `json:"sources" comment:"资源信息"`
	Mobile       string      `json:"mobile"`
	ShowType     int         `json:"show_type" comment:"展示类型:1单张,2合并"`
	FileType     int         `json:"file_type" comment:"文件类型: 1图片,2视频"`
}
type MediaInfo struct {
	Url           string       `json:"url"`
	Name          string       `json:"name"`
	Ext           string       `json:"ext"`
	Size          int          `json:"size"`
	Sort          int32        `json:"sort" comment:"排序"`
	IsOrigin      int          `json:"is_origin"`
	IsTranscode   int          `json:"is_transcode" comment:"是否转码,0为否,1为是"`
	CoverUrl      string       `json:"cover_url" comment:"视频第一帧的封面图"`
	Type          int          `json:"type" comment:"1:图片 2:视频"`
	RelationGoods []BrandGoods `json:"relation_goods" comment:"关联商品"`
}

type SourcesInfo struct {
	SourceId   int64 `json:"source_id"`
	SourceType int32 `json:"source_type"`
}

type SourcesMediaInfo struct {
	SourcesInfo
	MediaInfo
}

type CompressReq struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
	SourceType   int   `json:"source_type"`
	MemberId     int64 `json:"member_id"`
	CompressId   int64 `json:"compress_id"`
}

type CompressRes struct {
	Url    string `json:"url"`
	Name   string `json:"name"`
	IsDeal int32  `json:"is_deal" comment:"压缩结果,2为成功，3为失败"`
	Reason string `json:"reason" comment:"压缩失败原因"`
}

type FfmpegProbe struct {
	Streams []struct {
		Index            int    `json:"index"`
		CodecName        string `json:"codec_name"`
		CodecLongName    string `json:"codec_long_name"`
		Profile          string `json:"profile,omitempty"`
		CodecType        string `json:"codec_type"`
		CodecTagString   string `json:"codec_tag_string"`
		CodecTag         string `json:"codec_tag"`
		Width            int    `json:"width,omitempty"`
		Height           int    `json:"height,omitempty"`
		CodedWidth       int    `json:"coded_width,omitempty"`
		CodedHeight      int    `json:"coded_height,omitempty"`
		ClosedCaptions   int    `json:"closed_captions,omitempty"`
		FilmGrain        int    `json:"film_grain,omitempty"`
		HasBFrames       int    `json:"has_b_frames,omitempty"`
		PixFmt           string `json:"pix_fmt,omitempty"`
		Level            int    `json:"level,omitempty"`
		ChromaLocation   string `json:"chroma_location,omitempty"`
		FieldOrder       string `json:"field_order,omitempty"`
		Refs             int    `json:"refs,omitempty"`
		IsAvc            string `json:"is_avc,omitempty"`
		NalLengthSize    string `json:"nal_length_size,omitempty"`
		ID               string `json:"id"`
		RFrameRate       string `json:"r_frame_rate"`
		AvgFrameRate     string `json:"avg_frame_rate"`
		TimeBase         string `json:"time_base"`
		StartPts         int    `json:"start_pts"`
		StartTime        string `json:"start_time"`
		DurationTs       int    `json:"duration_ts"`
		Duration         string `json:"duration"`
		BitRate          string `json:"bit_rate"`
		BitsPerRawSample string `json:"bits_per_raw_sample,omitempty"`
		NbFrames         string `json:"nb_frames"`
		ExtradataSize    int    `json:"extradata_size,omitempty"`
		Disposition      struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
			Captions        int `json:"captions"`
			Descriptions    int `json:"descriptions"`
			Metadata        int `json:"metadata"`
			Dependent       int `json:"dependent"`
			StillImage      int `json:"still_image"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
			VendorID    string `json:"vendor_id"`
		} `json:"tags"`
		SampleFmt     string `json:"sample_fmt,omitempty"`
		SampleRate    string `json:"sample_rate,omitempty"`
		Channels      int    `json:"channels,omitempty"`
		ChannelLayout string `json:"channel_layout,omitempty"`
		BitsPerSample int    `json:"bits_per_sample,omitempty"`
	} `json:"streams"`
	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		ProbeScore     int    `json:"probe_score"`
		Tags           struct {
			MajorBrand       string `json:"major_brand"`
			MinorVersion     string `json:"minor_version"`
			CompatibleBrands string `json:"compatible_brands"`
			Encoder          string `json:"encoder"`
		} `json:"tags"`
	} `json:"format"`
}

type GetMediaInfoRequest struct {
	EnterpriseID int64  `json:"enterprise_id" comment:"企业ID"`
	SourceID     int64  `json:"source_id"`
	SourceType   int64  `json:"source_type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件)"`
	Url          string `json:"url" comment:"素材的cos地址"`
	Type         int64  `json:"type" comment:"0为临时素材，1为附件资源(用户朋友圈)"`
}

type FriendCreateReq struct {
	EnterpriseID uint64   `json:"enterprise_id" binding:"required" required:"企业ID"`
	StartTime    int64    `json:"start_time" binding:"required" required:"开始时间"`
	EndTime      int64    `json:"end_time" binding:"required" required:"结束时间"`
	CategoryIds  []uint64 `json:"category_ids" binding:"required" required:"产品类别"`
	SourceType   int32    `json:"source_type" binding:"required" required:"资源类型"`
	GroupId      int64    `json:"group_id" binding:"required" required:"分组"`
	FileType     int      `json:"file_type" comment:"文件类型: 1图片,2视频"`
	AdminId      int64    `json:"admin_id"`
	ShowType     int      `json:"show_type" comment:"展示类型"`
	TextContent  string   `json:"text_content" comment:"文本内容"`
	Sources      []struct {
		Id            int64        `json:"id"`
		Url           string       `json:"url"`
		VideoCoverUrl string       `json:"video_cover_url" comment:"封面图"`
		IsOrigin      int          `json:"is_origin"`
		Name          string       `json:"name"`
		Ext           string       `json:"ext"`
		Size          int          `json:"size"`
		RelationGoods []BrandGoods `json:"relation_goods" comment:"关联商品"`
	} `json:"sources"`
}

type FriendDetailReq struct {
	DetailRequest
}

type FriendUpdateReq struct {
	ID uint64 `json:"id" binding:"required" comment:"id"`
	FriendCreateReq
}

type BrandListReq struct {
	BrandCacheReq
	Size     int    `json:"size"`
	LastCode string `json:"last_code"`
}
type BrandCacheReq struct {
	EnterpriseId            int64    `json:"enterprise_id" binding:"required"`
	SkuItem                 string   `json:"sku_item"`
	SkuColor                string   `json:"sku_color"`
	ProductName             string   `json:"product_name"`
	StoreCode               string   `json:"store_code"`
	Mobile                  string   `json:"mobile"`
	UpcCode                 string   `json:"upc_code"`
	UpcCodes                []string `json:"upc_codes"`
	ProductCode             []string `json:"product_code"`
	Sign                    bool     `json:"sign"`
	Accurate                bool     `json:"accurate"`
	IncludeNoneImageProduct bool     `json:"include_none_image_product"`
}

// MakeBrandMediaReq  给某个素材生成临时素材ID
type MakeBrandMediaReq struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id"`
	SourceType   string `json:"source_type"`
	SourceUrl    string `json:"source_url"`
	SourceName   string `json:"source_name"`
}

type ReportTimesReq struct {
	Ids          []int    `json:"ids"`
	Type         int      `json:"type" binding:"required" required:"类型"`
	EnterpriseID int64    `json:"enterprise_id" comment:"企业ID"`
	ResourceKeys []string `json:"resource_keys"`
}

// CollectListReq 收藏请求
type CollectListReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	MemberId    int64   `json:"member_id" binding:"required"`
	Type        int32   `json:"type" binding:"required" required:"类型"`
	CategoryIds []int64 `json:"category_ids" comment:"产品类别"`
	Keyword     string  `json:"keyword"`
}

type CollectReq struct {
	request.EnterpriseBaseReq
	SourceInfo []CollectSourceInfoReq `json:"source_info"`
	MemberId   int64                  `json:"member_id" binding:"required"`
	IsCollect  int32                  `json:"is_collect" comment:"是否收藏,0取消收藏，1添加收藏"`
}
type CollectSourceInfoReq struct {
	Type        int32    `json:"type" required:"类型"`
	SourceId    []int64  `json:"source_id"`
	MergeId     []int64  `json:"merge_id"`
	ProductCode []string `json:"product_code"`
}

type SourceSetTop struct {
	request.EnterpriseBaseReq
	SourcesInfo
	IsTop int32 `json:"is_top" comment:"是否置顶,1是，2否"`
}

// TalkieImageSearchReq Talkie图片搜索请求
type TalkieImageSearchReq struct {
	request.EnterpriseBaseReq
	ImageURL string `json:"image_url" binding:"required" comment:"图片URL地址"`
	Num      int    `json:"num" comment:"每个区域返回产品数量，默认5，最多10"`
}
