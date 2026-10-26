package qyWork

type UploadAttachmentResp struct {
	QyResponse
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int    `json:"created_at"`
}

type UploadMediaResp struct {
	QyResponse
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}
