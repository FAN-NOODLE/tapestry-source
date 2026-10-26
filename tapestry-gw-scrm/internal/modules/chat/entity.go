package chat

// File 文件信息
type File struct {
	Url     string `json:"url" desc:"文件url"`
	Content string `json:"content" desc:"文件流"`
}
