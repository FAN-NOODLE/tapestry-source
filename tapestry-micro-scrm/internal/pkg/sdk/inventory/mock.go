package inventory

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	redisClient "micro-scrm/internal/pkg/redis"
)

const inventoryMockKeyPrefix = "inventory:mock:uat"

type inventoryMock struct {
	StatusCode int
	Body       string
	Delay      time.Duration
}

type inventoryMockGetter func(key string) (map[string]string, error)

func findInventoryMock(env string, isDcl bool, request ParamsReq, get inventoryMockGetter) (inventoryMock, bool, error) {
	if strings.ToLower(env) != "uat" {
		return inventoryMock{}, false, nil
	}
	source := "hht"
	if isDcl {
		source = "dcl"
	}
	for _, upc := range request.UpcCode {
		key := fmt.Sprintf("%s:%s:%s", inventoryMockKeyPrefix, source, upc)
		data, err := get(key)
		if err != nil {
			return inventoryMock{}, false, err
		}
		if len(data) == 0 {
			continue
		}
		mock := inventoryMock{StatusCode: 200, Body: data["body"]}
		if value := data["status_code"]; value != "" {
			mock.StatusCode, err = strconv.Atoi(value)
			if err != nil {
				return inventoryMock{}, false, fmt.Errorf("库存 mock status_code 无效：%s", value)
			}
		}
		if value := data["delay_ms"]; value != "" {
			delayMs, parseErr := strconv.Atoi(value)
			if parseErr != nil || delayMs < 0 || delayMs > 120000 {
				return inventoryMock{}, false, fmt.Errorf("库存 mock delay_ms 无效：%s", value)
			}
			mock.Delay = time.Duration(delayMs) * time.Millisecond
		}
		return mock, true, nil
	}
	return inventoryMock{}, false, nil
}

func applyInventoryMock(isDcl bool, request ParamsReq, mock inventoryMock, resp *Response) error {
	if mock.StatusCode != 200 {
		return errorx.NewErrorMsg(mock.StatusCode, fmt.Sprintf("三方库存接口-接口请求状态码异常：mock HTTP %d", mock.StatusCode))
	}
	if mock.Body == "" {
		return fmt.Errorf("库存 mock body 不能为空")
	}
	return decodeAndValidateResponse(isDcl, request, mock.Body, resp)
}

func tryInventoryMock(c ctx.Context, isDcl bool, request ParamsReq, resp *Response) (bool, string, error) {
	mock, found, err := findInventoryMock(config.CfgData.Env, isDcl, request, func(key string) (map[string]string, error) {
		return redisClient.Client().HGetAll(c, key).Result()
	})
	if err != nil {
		log.Error(c).Msgf("库存 mock 读取失败，将请求真实接口：%v", err)
		return false, "", nil
	}
	if !found {
		return false, "", nil
	}
	if mock.Delay > 0 {
		time.Sleep(mock.Delay)
	}
	log.Info(c).Msgf("库存 mock 命中, is_dcl:%v, upc:%v, status_code:%d", isDcl, request.UpcCode, mock.StatusCode)
	return true, mock.Body, applyInventoryMock(isDcl, request, mock, resp)
}
