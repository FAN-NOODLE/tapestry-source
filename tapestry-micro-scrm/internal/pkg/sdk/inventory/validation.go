package inventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ContractError means the provider returned HTTP 200 but the response cannot
// be safely consumed by the inventory business flow.
type ContractError struct {
	Reason string
}

func (e *ContractError) Error() string {
	return "三方库存接口-返回数据契约异常：" + e.Reason
}

func contractError(format string, args ...interface{}) error {
	return &ContractError{Reason: fmt.Sprintf(format, args...)}
}

func IsContractError(err error) bool {
	var target *ContractError
	return errors.As(err, &target)
}

func isJSONNull(raw json.RawMessage) bool {
	return len(raw) == 0 || string(raw) == "null"
}

func requiredString(data map[string]json.RawMessage, field string) (string, error) {
	raw, ok := data[field]
	if !ok || isJSONNull(raw) {
		return "", contractError("字段 %s 缺失或为 null", field)
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", contractError("字段 %s 类型错误", field)
	}
	if strings.TrimSpace(value) == "" {
		return "", contractError("字段 %s 为空", field)
	}
	return value, nil
}

func requiredInt(data map[string]json.RawMessage, field string) (int, error) {
	raw, ok := data[field]
	if !ok || isJSONNull(raw) {
		return 0, contractError("字段 %s 缺失或为 null", field)
	}
	var value int
	if err := json.Unmarshal(raw, &value); err != nil {
		return 0, contractError("字段 %s 类型错误", field)
	}
	return value, nil
}

func decodeAndValidateResponse(isDcl bool, request ParamsReq, body string, resp *Response) error {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return contractError("JSON 无法解析：%v", err)
	}
	if root == nil {
		return contractError("响应根结构不是 JSON 对象")
	}

	headerRaw, ok := root["header"]
	if !ok || isJSONNull(headerRaw) {
		return contractError("header 缺失或为 null")
	}
	var header map[string]json.RawMessage
	if err := json.Unmarshal(headerRaw, &header); err != nil || header == nil {
		return contractError("header 类型错误")
	}
	resultCode, err := requiredString(header, "resultCode")
	if err != nil {
		return err
	}
	if resultCode != "SUCCESS" {
		return fmt.Errorf("三方库存接口返回失败结果：%s", resultCode)
	}

	detailRaw, ok := root["detailList"]
	if !ok || isJSONNull(detailRaw) {
		return contractError("detailList 缺失或为 null")
	}
	var details []map[string]json.RawMessage
	if err := json.Unmarshal(detailRaw, &details); err != nil {
		return contractError("detailList 类型错误")
	}

	requestedUPCs := make(map[string]struct{}, len(request.UpcCode))
	for _, upc := range request.UpcCode {
		requestedUPCs[upc] = struct{}{}
	}
	requestedLocations := make(map[string]struct{}, len(request.LocCode))
	for _, location := range request.LocCode {
		requestedLocations[location] = struct{}{}
	}
	validDetails := make([]map[string]json.RawMessage, 0, len(details))
	for index, detail := range details {
		upcCode, err := requiredString(detail, "upcCode")
		if err != nil {
			return contractError("detailList[%d] %v", index, err)
		}
		if len(requestedUPCs) > 0 {
			if _, exists := requestedUPCs[upcCode]; !exists {
				return contractError("detailList[%d] upcCode %s 不在本次请求范围", index, upcCode)
			}
		}
		locCode, err := requiredString(detail, "locCode")
		if err != nil {
			if isDcl {
				return contractError("detailList[%d] %v", index, err)
			}
			// HHT occasionally appends an unusable summary/detail row without locCode.
			// Ignore that row when other usable store rows exist, and reject the
			// response below when every non-empty detail row is unusable.
			continue
		}
		if _, err = requiredInt(detail, "onHand"); err != nil {
			return contractError("detailList[%d] %v", index, err)
		}

		if isDcl {
			_, segmentErr := requiredString(detail, "segment")
			if segmentErr != nil {
				return contractError("detailList[%d] %v", index, segmentErr)
			}
		} else {
			if len(requestedLocations) > 0 {
				if _, exists := requestedLocations[locCode]; !exists {
					return contractError("detailList[%d] locCode %s 不在本次请求范围", index, locCode)
				}
			}
			if _, err = requiredInt(detail, "inTransit"); err != nil {
				return contractError("detailList[%d] %v", index, err)
			}
		}
		validDetails = append(validDetails, detail)
	}

	if !isDcl && len(details) > 1 && len(validDetails) == 0 {
		return contractError("detailList 非空但所有记录的 locCode 均缺失或不可用")
	}
	if !isDcl && len(validDetails) != len(details) {
		filteredDetailRaw, marshalErr := json.Marshal(validDetails)
		if marshalErr != nil {
			return contractError("过滤无效 detailList 失败：%v", marshalErr)
		}
		root["detailList"] = filteredDetailRaw
		filteredBody, marshalErr := json.Marshal(root)
		if marshalErr != nil {
			return contractError("重建响应失败：%v", marshalErr)
		}
		body = string(filteredBody)
	}

	if err := json.Unmarshal([]byte(body), resp); err != nil {
		return contractError("响应转换失败：%v", err)
	}
	return nil
}
