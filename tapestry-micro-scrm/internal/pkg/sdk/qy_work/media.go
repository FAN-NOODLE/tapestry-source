package qyWork

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/helper"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// UploadAttachMent 上传附件资源
	uploadAttachMent = "cgi-bin/media/upload_attachment"
	uploadMedia      = "cgi-bin/media/upload"
)

func (w *QyWork) UploadAttachment(c app.Ctx, req map[string]string, filePath string) (UploadAttachmentResp, error) {
	var resp UploadAttachmentResp
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	var files []*os.File
	if filePath != "" {
		file, fileErr := os.Open(filePath)
		if fileErr != nil {
			return resp, err
		}
		files = append(files, file)
	}

	err = w.requestPostForm(c, w.getTokenUrl(uploadAttachMent, token), req, files, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUploadAttachmentMediaId 获取上传附件资源获取media_id
func (w *QyWork) GetUploadAttachmentMediaId(c ctx.Context, mediaUrl string, mediaType string, attachmentTye int64) (string, error) {
	redisKey := helper.Md5V(mediaUrl)
	mediaId := app.Redis().Get(c, redisKey).Val()
	if mediaId != "" {
		return mediaId, nil
	}
	ext := path.Ext(mediaUrl)
	// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
	if u, err := url.Parse(mediaUrl); err == nil && u.Path != "" {
		ext = path.Ext(u.Path)
	}
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	reqData := map[string]string{
		"media_type":      mediaType,
		"filename":        fileName,
		"attachment_type": helper.Int64ToString(attachmentTye),
	}

	fileData, err := readMediaFile(c, mediaUrl)
	if err != nil {
		return mediaId, err
	}
	tempFilepath := fmt.Sprintf("storage/upload/%s", fileName)
	defer os.Remove(tempFilepath)
	if err = os.WriteFile(tempFilepath, fileData, 0755); err != nil {
		return mediaId, err
	}

	attachment, err := w.UploadAttachment(c, reqData, tempFilepath)
	if err != nil {
		return mediaId, err
	}
	mediaId = attachment.MediaId
	app.Redis().Set(c, redisKey, mediaId, 3*24*time.Hour) // 缓存素材ID 三天
	return mediaId, err
}

func (w *QyWork) UploadMedia(c app.Ctx, req map[string]string, filePath string) (UploadMediaResp, error) {
	var resp UploadMediaResp
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	var files []*os.File
	if filePath != "" {
		file, fileErr := os.Open(filePath)
		if fileErr != nil {
			return resp, err
		}
		files = append(files, file)
	}

	err = w.requestPostForm(c, w.getTokenUrl(uploadMedia, token), req, files, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUploadMediaId 上传临时素材获取media_id
func (w *QyWork) GetUploadMediaId(c ctx.Context, mediaUrl string, mediaType string) (string, error) {
	redisKey := helper.Md5V(mediaUrl)
	mediaId := app.Redis().Get(c, redisKey).Val()
	if mediaId != "" {
		return mediaId, nil
	}
	ext := path.Ext(mediaUrl)
	// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
	if u, err := url.Parse(mediaUrl); err == nil && u.Path != "" {
		ext = path.Ext(u.Path)
	}
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	reqData := map[string]string{
		"type":     mediaType,
		"filename": fileName,
	}

	fileData, err := readMediaFile(c, mediaUrl)
	if err != nil {
		return mediaId, err
	}
	tempFilepath := fmt.Sprintf("storage/upload/%s", fileName)
	defer os.Remove(tempFilepath)
	if err = os.WriteFile(tempFilepath, fileData, 0755); err != nil {
		return mediaId, err
	}

	attachment, err := w.UploadMedia(c, reqData, tempFilepath)
	if err != nil {
		return mediaId, err
	}
	mediaId = attachment.MediaId
	app.Redis().Set(c, redisKey, mediaId, 3*24*time.Hour) // 缓存素材ID 三天
	return mediaId, err
}

func readMediaFile(c ctx.Context, mediaUrl string) ([]byte, error) {
	parsed, err := url.Parse(mediaUrl)
	if err != nil || parsed.RawQuery == "" {
		return helper.ReadRemoteFile(c, mediaUrl)
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, mediaUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download media file failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
