package emailReq

type SendMailReq struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ToAddress []string   `json:"to_address"`
	FileList  []FileInfo `json:"file_info"`
}

type FileInfo struct {
	Url       string `json:"url"`
	Name      string `json:"name"`
	LocalPath string `json:"localPath"`
}
