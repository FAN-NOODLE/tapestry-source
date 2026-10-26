package talkie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

const (
	TokenRedisKey    = "talkie:token:%d"
	CosCompressParam = "imageMogr2/quality/80/format/jpg" // COS 图片压缩参数，保持原分辨率，只压缩质量
)

// getToken 获取Token
func getToken(ctx ctx.Context) (string, error) {
	enterpriseId := int64(0) // Talkie 不需要企业ID，但为了缓存key统一使用0

	// 从Redis获取缓存的token
	cachedToken := redis.Client().Get(ctx, fmt.Sprintf(TokenRedisKey, enterpriseId)).Val()
	if cachedToken != "" {
		return cachedToken, nil
	}

	// 获取配置
	conf := config.CfgData.Talkie
	if conf.BaseURL == "" {
		return "", errorx.NewErrorMsg(500, "Talkie配置未设置")
	}

	tokenURL := conf.BaseURL + "/auth/v2/api/token"

	var req TokenRequest
	req.ApiKey = conf.ApiKey
	req.ApiSecret = conf.ApiSecret

	var resp TokenResponse
	reqStr := helper.MarshalToString(req)

	err := helper.WithRetry(1, time.Second, func() error {
		body := clients.JSONBody(req)
		client := clients.NewHTTPClient(clients.Timeout(30 * time.Second))
		res, err := client.Post(ctx, tokenURL, body)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie Token获取失败：%+v", err))
		}

		if res.Response().StatusCode != 200 {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("Talkie Token请求状态码异常：%s", res.Response().Status))
		}

		rt, rErr := res.Read(clients.Text)
		if rErr != nil {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("Talkie Token响应解析失败：%+v", rErr))
		}

		log.Info(ctx).Msgf("Talkie Token请求, 地址：%s, 参数：%s, 返回：%s", tokenURL, reqStr, rt.(string))

		uIntErr := json.Unmarshal([]byte(rt.(string)), &resp)
		if uIntErr != nil {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("Talkie Token响应JSON解析失败：%+v", uIntErr))
		}

		if !resp.Success {
			errMsg := "Token获取失败"
			if resp.Detail != nil {
				errMsg = fmt.Sprintf("%s: %s", resp.Detail.Msg, resp.Detail.StatusCode)
			}
			return errorx.NewErrorMsg(res.Response().StatusCode, errMsg)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if resp.Data.AccessToken == "" {
		return "", errorx.NewErrorMsg(500, "Talkie Token为空")
	}

	// 缓存token，提前5分钟过期
	expires := resp.Data.Expires - 300
	if expires <= 0 {
		expires = 3300 // 默认 55 分钟
	}
	redis.Client().Set(ctx, fmt.Sprintf(TokenRedisKey, enterpriseId), resp.Data.AccessToken, time.Duration(expires)*time.Second)

	return resp.Data.AccessToken, nil
}

// ImageSearch 图片搜索（通过URL下载图片后上传）
func ImageSearch(c ctx.Context, req ImageSearchRequest) (*ImageSearchResponse, error) {
	// 获取token
	accessToken, err := getToken(c)
	if err != nil {
		return nil, err
	}

	// 获取配置
	conf := config.CfgData.Talkie
	if conf.BaseURL == "" {
		return nil, errorx.NewErrorMsg(500, "Talkie配置未设置")
	}

	searchURL := conf.BaseURL + "/product/v2/imagesearch/api/external"

	// 设置默认值
	num := req.Num
	if num <= 0 {
		num = 5
	}
	if num > 10 {
		num = 10
	}

	// 下载图片（先用原始URL下载，超过4MB再尝试压缩）
	log.Info(c).Msgf("Talkie图片搜索开始下载图片: %s", req.ImageURL)

	httpClient := &http.Client{Timeout: 60 * time.Second}

	downloadImage := func(downloadURL string) ([]byte, int, error) {
		imageReq, err := http.NewRequestWithContext(c, "GET", downloadURL, nil)
		if err != nil {
			return nil, 0, errorx.NewErrorMsg(500, fmt.Sprintf("创建图片下载请求失败：%+v", err))
		}
		imageResp, err := httpClient.Do(imageReq)
		if err != nil {
			return nil, 0, errorx.NewErrorMsg(500, fmt.Sprintf("下载图片失败：%+v", err))
		}
		defer imageResp.Body.Close()
		if imageResp.StatusCode != 200 {
			return nil, imageResp.StatusCode, errorx.NewErrorMsg(imageResp.StatusCode, fmt.Sprintf("下载图片状态码异常：%d", imageResp.StatusCode))
		}
		data, err := io.ReadAll(imageResp.Body)
		if err != nil {
			return nil, imageResp.StatusCode, errorx.NewErrorMsg(500, fmt.Sprintf("读取图片数据失败：%+v", err))
		}
		return data, imageResp.StatusCode, nil
	}

	// 第一次：使用原始URL下载
	imageData, _, err := downloadImage(req.ImageURL)
	if err != nil {
		return nil, err
	}

	log.Info(c).Msgf("Talkie图片搜索图片下载成功, 大小: %d bytes, URL: %s", len(imageData), req.ImageURL)

	// 如果图片超过4MB，尝试通过COS循环压缩后重新下载
	const maxImageSize = 4 * 1024 * 1024
	if int64(len(imageData)) > maxImageSize {
		log.Info(c).Msgf("Talkie图片搜索图片超过4MB(%d bytes)，尝试COS循环压缩", len(imageData))

		// 定义质量级别，从高到低尝试
		qualityLevels := []int{80, 60, 40, 20}
		compressed := false

		for _, quality := range qualityLevels {
			compressedURL, buildErr := buildCosImageURL(c, req.ImageURL, quality)
			if buildErr != nil {
				log.Info(c).Msgf("Talkie图片搜索无法构建压缩URL（质量%d），跳过: %v", quality, buildErr)
				break
			}

			compressedData, _, compressErr := downloadImage(compressedURL)
			if compressErr != nil {
				log.Info(c).Msgf("Talkie图片搜索质量%d压缩下载失败: %v", quality, compressErr)
				continue
			}

			log.Info(c).Msgf("Talkie图片搜索质量%d压缩后大小: %d bytes, URL: %s", quality, len(compressedData), compressedURL)
			imageData = compressedData
			compressed = true

			// 如果压缩后小于等于4MB，满足要求，跳出循环
			if int64(len(imageData)) <= maxImageSize {
				log.Info(c).Msgf("Talkie图片搜索压缩后大小(%d bytes)满足要求(≤4MB)，停止压缩", len(imageData))
				break
			}

			log.Info(c).Msgf("Talkie图片搜索质量%d压缩后仍超过4MB(%d bytes)，继续降低质量", quality, len(imageData))
		}

		if !compressed {
			log.Info(c).Msgf("Talkie图片搜索所有压缩尝试均失败")
			return nil, errorx.NewErrorMsg(400, "图片文件过大，无法处理，请更换图片后重试")
		} else if int64(len(imageData)) > maxImageSize {
			log.Info(c).Msgf("Talkie图片搜索经过所有质量级别压缩后仍超过4MB(%d bytes)", len(imageData))
			return nil, errorx.NewErrorMsg(400, "图片文件过大，无法处理，请更换图片后重试")
		}
	}

	// 获取图片文件名（从URL中提取）
	filename := "image.jpg"
	if urlFilename := getFilenameFromURL(req.ImageURL); urlFilename != "" {
		filename = urlFilename
	}

	var resp ImageSearchResponse

	err = helper.WithRetry(1, time.Second, func() error {
		// 构建 multipart/form-data 请求
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 添加 file 字段（图片文件），手动设置 Content-Type 避免 CreateFormFile 固定写死 application/octet-stream
		contentType := mime.TypeByExtension(filepath.Ext(filename))
		if contentType == "" {
			contentType = "image/jpeg" // 默认 image/jpeg
		}
		fileHeader := make(textproto.MIMEHeader)
		fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
		fileHeader.Set("Content-Type", contentType)
		part, err := writer.CreatePart(fileHeader)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索创建表单文件失败：%+v", err))
		}
		if _, err := part.Write(imageData); err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索写入图片数据失败：%+v", err))
		}

		// 添加 num 字段
		if err := writer.WriteField("num", fmt.Sprintf("%d", num)); err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索表单字段写入失败：%+v", err))
		}

		// 添加 categoryId 字段（可选）
		if req.CategoryID > 0 {
			if err := writer.WriteField("categoryId", fmt.Sprintf("%d", req.CategoryID)); err != nil {
				return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索表单字段写入失败：%+v", err))
			}
		}

		// 关闭 writer
		if err := writer.Close(); err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索表单关闭失败：%+v", err))
		}

		httpReq, err := http.NewRequestWithContext(c, "POST", searchURL, body)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索请求创建失败：%+v", err))
		}

		httpReq.Header.Set("Content-Type", writer.FormDataContentType())
		httpReq.Header.Set("Authorization", "Bearer "+accessToken)

		log.Info(c).Msgf("Talkie图片搜索请求, 地址：%s, num：%d, categoryId：%d, filename：%s, 文件大小：%d", searchURL, num, req.CategoryID, filename, len(imageData))

		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("Talkie图片搜索请求失败：%+v", err))
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != 200 {
			return errorx.NewErrorMsg(httpResp.StatusCode, fmt.Sprintf("Talkie图片搜索状态码异常：%d", httpResp.StatusCode))
		}

		respBody, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return errorx.NewErrorMsg(httpResp.StatusCode, fmt.Sprintf("Talkie图片搜索响应读取失败：%+v", err))
		}

		bodyStr := string(respBody)
		log.Info(c).Msgf("Talkie图片搜索响应, 地址：%s, 返回：%s", searchURL, bodyStr)

		uIntErr := json.Unmarshal(respBody, &resp)
		if uIntErr != nil {
			return errorx.NewErrorMsg(httpResp.StatusCode, fmt.Sprintf("Talkie图片搜索响应JSON解析失败：%+v", uIntErr))
		}

		if !resp.Success {
			errMsg := "图片搜索失败"
			if resp.Detail != nil {
				errMsg = fmt.Sprintf("%s: %s", resp.Detail.Msg, resp.Detail.StatusCode)
			}
			return errorx.NewErrorMsg(httpResp.StatusCode, errMsg)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// buildCosImageURL 构建COS图片访问URL，添加压缩参数保持原分辨率只压缩质量
// COS 图片处理参数格式：?imageMogr2/quality/80
// 注意：对于CDN域名，需要CDN配置开启图片处理；对于COS源站域名可直接使用
func buildCosImageURL(c ctx.Context, imageURL string, quality int) (string, error) {
	// 检查是否是COS相关域名
	cosHost := config.CfgData.Storage.COS.CdnHost
	if cosHost == "" {
		cosHost = config.CfgData.Storage.COS.Host
	}

	if !strings.Contains(imageURL, cosHost) && !strings.Contains(imageURL, "cos.") {
		return "", fmt.Errorf("非COS域名图片，跳过压缩")
	}

	// 检查URL中是否已经有图片处理参数
	if strings.Contains(imageURL, "imageMogr2") {
		return "", fmt.Errorf("图片URL已包含处理参数，跳过压缩")
	}

	// 动态构建COS图片处理参数
	// 格式：?imageMogr2/quality/{quality}/format/jpg
	// - quality/{quality}: 设置绝对质量
	// - format/jpg: 转换为jpg格式（通常能显著减小体积）
	// 注意：不设置任何缩放参数，保持原分辨率
	compressParam := fmt.Sprintf("imageMogr2/quality/%d/format/jpg", quality)

	// 添加参数到URL
	separator := "?"
	if strings.Contains(imageURL, "?") {
		separator = "&"
	}

	compressedURL := imageURL + separator + compressParam

	log.Info(c).Msgf("Talkie图片搜索COS压缩URL(质量%d): %s", quality, compressedURL)

	return compressedURL, nil
}

// getFilenameFromURL 从URL中提取文件名
func getFilenameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	// 获取路径最后一部分作为文件名
	parts := strings.Split(u.Path, "/")
	if len(parts) > 0 {
		filename := parts[len(parts)-1]
		// 检查是否包含文件扩展名
		if strings.Contains(filename, ".") {
			return filename
		}
	}
	return ""
}
