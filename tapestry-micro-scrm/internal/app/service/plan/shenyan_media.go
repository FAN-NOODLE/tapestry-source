package planService

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/cosroute"
	"micro-scrm/internal/pkg/helper"
)

const (
	shenyanImageCompressParam = "imageMogr2/thumbnail/810x"
	shenyanMaxMediaSize       = 10 * 1024 * 1024
	shenyanMomentLongSide     = 10800
	shenyanMomentShortSide    = 1080
)

var shenyanCoachCDNHosts = map[string]struct{}{
	"cdn-cos.coach.com.cn":     {},
	"cdn-cos-uat.coach.com.cn": {},
}

func normalizeShenyanImageURL(rawURL string) (string, bool) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, false
	}
	if _, ok := shenyanCoachCDNHosts[strings.ToLower(parsed.Host)]; !ok {
		return rawURL, false
	}
	if parsed.RawQuery != "" {
		return rawURL, false
	}
	parsed.RawQuery = shenyanImageCompressParam
	return parsed.String(), true
}

func validateShenyanImageMedia(c ctx.Context, rawURL string, isMoment bool) error {
	if isShenyanCoachCDNURL(rawURL) {
		return nil
	}
	data, contentType, err := fetchShenyanMedia(c, rawURL)
	if err != nil {
		return fmt.Errorf("image url cannot be read: %w", err)
	}
	if len(data) > shenyanMaxMediaSize {
		return errors.New("image size cannot exceed 10MB")
	}
	format, err := validateShenyanImageFormat(contentType, data)
	if err != nil {
		return err
	}
	config, decodedFormat, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return errors.New("image content is invalid")
	}
	if format != "" && decodedFormat != format {
		return errors.New("image format must be jpg or png")
	}
	if isMoment && !isValidShenyanMomentImageResolution(config.Width, config.Height) {
		return errors.New("moment image long side cannot exceed 10800 and short side cannot exceed 1080")
	}
	return nil
}

func validateShenyanVideoMedia(c ctx.Context, rawURL string) error {
	data, contentType, err := fetchShenyanMedia(c, rawURL)
	if err != nil {
		return fmt.Errorf("video url cannot be read: %w", err)
	}
	if len(data) > shenyanMaxMediaSize {
		return errors.New("video size cannot exceed 10MB")
	}
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	ext := strings.ToLower(path.Ext(rawURL))
	if parsed, err := url.Parse(rawURL); err == nil && parsed.Path != "" {
		ext = strings.ToLower(path.Ext(parsed.Path))
	}
	if contentType != "" && contentType != "video/mp4" && contentType != "application/octet-stream" {
		return errors.New("video must be mp4")
	}
	if ext != "" && ext != ".mp4" {
		return errors.New("video must be mp4")
	}
	return nil
}

func isShenyanCoachCDNURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	_, ok := shenyanCoachCDNHosts[strings.ToLower(parsed.Host)]
	return ok
}

func fetchShenyanMedia(c ctx.Context, rawURL string) ([]byte, string, error) {
	if _, err := cosroute.ParseObjectURL(rawURL); err == nil {
		data, err := helper.ReadRemoteFile(c, rawURL)
		if err != nil {
			return nil, "", err
		}
		return data, http.DetectContentType(data), nil
	}

	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(c, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, shenyanMaxMediaSize+1))
	if err != nil {
		return nil, "", err
	}
	return data, resp.Header.Get("Content-Type"), nil
}

func validateShenyanImageFormat(contentType string, data []byte) (string, error) {
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	switch contentType {
	case "", "application/octet-stream":
	case "image/jpeg", "image/jpg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	default:
		return "", errors.New("image format must be jpg or png")
	}
	switch http.DetectContentType(data) {
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	default:
		return "", errors.New("image format must be jpg or png")
	}
}

func isValidShenyanMomentImageResolution(width, height int) bool {
	longSide := width
	shortSide := height
	if height > width {
		longSide = height
		shortSide = width
	}
	return longSide <= shenyanMomentLongSide && shortSide <= shenyanMomentShortSide
}
