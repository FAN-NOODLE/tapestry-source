package inventorySw

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"time"
)

func Post(ctx ctx.Context, params ParamsReq) (Response, error) {
	var resp Response
	var opts []clients.Option

	// 获取配置
	conf := config.CfgData.InventorySw
	url := conf.Url

	// 构建 p_data 数据（根据实际业务需求调整字段）
	pData := buildPData(params)

	// 计算签名
	sign := calculateSign(pData, conf.Key)

	// 构建 URL 参数
	url += fmt.Sprintf("p_loc_no=%s&p_sign=%s", conf.PLocNo, sign)

	// 构建 POST 表单数据
	reqMap := map[string]string{
		"p_data": pData,
	}

	body := clients.Form(reqMap)
	opts = append(opts, body)

	// 返回值
	var respStr string
	opts = append(opts, clients.Header("Content-Type", "application/x-www-form-urlencoded"))
	client := clients.NewHTTPClient(clients.Timeout(120 * time.Second))
	res, err := client.Post(ctx, url, opts...)
	if err != nil {
		return resp, errors.New("HTTP请求失败: " + err.Error())
	}

	if res.Response().StatusCode != 200 {
		return resp, errors.New(fmt.Sprintf("HTTP状态码异常:%d", res.Response().StatusCode))
	}

	rt, rErr := res.Read(clients.Text)
	if rErr != nil {
		return resp, errors.New("读取响应失败: " + rErr.Error())
	}

	respStr = rt.(string)

	log.Error(ctx).Msgf("HTTP请求成功: url=%s, request:%v,resp=%s", url, reqMap, respStr)

	uIntErr := json.Unmarshal([]byte(respStr), &resp)
	if uIntErr != nil {
		return resp, errors.New("JSON解析失败: " + uIntErr.Error())
	}

	// 根据协议，返回结果应该包含状态码和状态信息
	if resp.Code != 0 { // 假设0表示成功，根据实际接口调整
		return resp, errors.New(resp.Msg)
	}

	return resp, nil
}

// buildPData 构建 p_data 参数
func buildPData(params ParamsReq) string {
	// 根据业务需求构建数据，字段之间用TAB分隔(ASCII码9)
	// 示例：用户注册格式
	fields := []string{
		"Z", // 指令
	}
	if params.SearchOtherStore {
		fields = append(fields, "KDKCCX", params.StoreCode, params.Plu)
	} else if params.StkName != "" {
		fields = append(fields, "HPMCCX", params.StoreCode, params.StkName)
	} else if params.StkCode != "" {
		fields = append(fields, "HPDMCX", params.StoreCode, params.StkCode)
	} else if params.Plu != "" {
		fields = append(fields, "HPTMCX", params.StoreCode, params.Plu)
	}
	// 使用TAB分隔符连接所有字段
	return strings.Join(fields, "\t")
}

// calculateSign 计算MD5签名
func calculateSign(pData, key string) string {
	// 拼接 p_data 和 KEY
	data := pData + key
	// 转为大写十六进制字符串
	return strings.ToUpper(helper.Md5V(data))
}
