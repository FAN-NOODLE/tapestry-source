package source

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/helper"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"
)

const coachKOSImageCompressParam = "imageMogr2/quality/80/format/jpg"

func (s *CoachKOSSyncService) prepareCoachKOSCoverURL(c app.Ctx, enterpriseID int64, post coachkos.PostResult) string {
	originalURL := post.CoverURL
	if strings.TrimSpace(originalURL) == "" {
		return originalURL
	}

	if s.getKOSByPostID != nil {
		existingKOS, err := s.getKOSByPostID(c, enterpriseID, post.PostID)
		if err != nil {
			s.logSyncInfo(c, "KOS封面复用检查失败，将重新转存：enterprise_id=%d post_id=%s err=%v", enterpriseID, post.PostID, err)
		} else if existingKOS.ID != 0 && existingKOS.SourceID != 0 && existingKOS.CoverURL == originalURL && s.getLinkByID != nil {
			existingLink, linkErr := s.getLinkByID(c, enterpriseID, existingKOS.SourceID)
			if linkErr != nil {
				s.logSyncInfo(c, "KOS封面复用查询失败，将重新转存：enterprise_id=%d post_id=%s err=%v", enterpriseID, post.PostID, linkErr)
			} else if isCoachKOSProcessedCoverURL(existingLink.CoverURL, s.cdnHost) {
				return existingLink.CoverURL
			}
		}
	}

	if s.uploadCover == nil {
		return originalURL
	}
	cdnURL, err := s.uploadCover(c, originalURL)
	if err != nil || strings.TrimSpace(cdnURL) == "" {
		if err == nil {
			err = fmt.Errorf("COS返回地址为空")
		}
		s.logSyncInfo(c, "KOS封面转存失败，使用原始地址：enterprise_id=%d post_id=%s err=%v", enterpriseID, post.PostID, err)
		return originalURL
	}
	return appendCoachKOSImageCompressParam(cdnURL)
}

func uploadCoachKOSCover(c app.Ctx, rawURL string) (string, error) {
	data, err := helper.ReadRemoteFile(c, rawURL)
	if err != nil {
		return "", fmt.Errorf("下载封面失败: %w", err)
	}
	if len(data) == 0 {
		return "", fmt.Errorf("下载封面失败: 文件内容为空")
	}

	extension := coachKOSCoverExtension(rawURL, data)
	tmpFile, err := os.CreateTemp("", "coach-kos-cover-*"+extension)
	if err != nil {
		return "", err
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err = tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return "", err
	}
	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	cdnURL, err := helper.UploadFileToCos(c, tmpPath, false)
	if err != nil {
		return "", fmt.Errorf("上传封面到COS失败: %w", err)
	}
	return cdnURL, nil
}

func coachKOSCoverExtension(_ string, _ []byte) string {
	// COS 地址会继续携带 format/jpg 参数，统一后缀可确保下游上传时文件名也按 JPG 处理。
	return ".jpg"
}

func appendCoachKOSImageCompressParam(rawURL string) string {
	separator := "?"
	if strings.Contains(rawURL, "?") {
		separator = "&"
	}
	return rawURL + separator + coachKOSImageCompressParam
}

func isCoachKOSProcessedCoverURL(rawURL, cdnHost string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || strings.TrimSpace(cdnHost) == "" {
		return false
	}
	return strings.EqualFold(parsedURL.Host, cdnHost) && strings.Contains(parsedURL.RawQuery, coachKOSImageCompressParam)
}
