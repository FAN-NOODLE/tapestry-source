package source

// TypeText 文本类型
const TypeText = 1

// TypeImg 图片类型
const TypeImg = 2

// TypeLink 链接类型
const TypeLink = 3

// TypeMiniProgram 小程序类型
const TypeMiniProgram = 4

// TypeVideo 视频类型
const TypeVideo = 5

// TypeFile 文件类型
const TypeFile = 6

// TypeCollect 文件集合
const TypeCollect = 7

// TypeGroup 文件分组
var TypeFileGroup = map[string]int{
	"png": TypeImg,
	"jpg": TypeImg,
	"mp4": TypeVideo,
}

// 文件后缀常量
const (
	WordFile  = 1
	ExcelFile = 2
	PptFile   = 3
	PdfFIle   = 4
	OtherFile = 5
)

// FileExtGroup 文件后缀分组
var FileExtGroup = map[string]int{
	"doc":  WordFile,
	"docx": WordFile,
	"xls":  ExcelFile,
	"xlsx": ExcelFile,
	"ppt":  PptFile,
	"pptx": PptFile,
	"pdf":  PdfFIle,
}

// LinkUrlTypeNormal 链接类型URL
const LinkUrlTypeNormal = 1

// LinkUrlTypeSelf 链接类型自定义
const LinkUrlTypeSelf = 2

// SelfLinkUrl 链接自定义素材URL
const SelfLinkUrl = "%s/friend/%s/source-link?id="

// SourceFileType 文件类型
const SourceFileType = "file"

// SourceVideoType 视频类型
const SourceVideoType = "video"

// SourceImageType 图片类型
const SourceImageType = "image"
