package source

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io/ioutil"
	"micro-download/internal/common"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
	"net/http"
	"os"
	"regexp"
	"sync"
)

var lock sync.Mutex

type packSource struct{}

func (s packSource) Create(ctx ctx.Context, export delivery.CreateExport) (error, int64) {
	if export.SourceTypePackPath == "" || export.SourceTypePackParams == nil {
		return errors.New("无效的pack导出参数"), 0
	}
	model := &modes.DownloadSource{}
	model.SetParams(export)
	return model.SourceCreate(ctx, func(d *modes.DownloadSource) (err error, downloadSource *modes.DownloadSource) {
		data := delivery.SourceTypePackData{
			Path:                     export.SourceTypePackPath,
			Method:                   export.SourceTypePackMethod,
			Headers:                  export.SourceTypePackHeaders,
			SourceTypePackDataParams: *export.SourceTypePackParams,
		}
		byteData, err := json.Marshal(data)
		if err != nil {
			return err, d
		}
		d.SourceData = string(byteData)
		return nil, d
	})
}

func (s packSource) Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error) {
	data := &delivery.SourceTypePackData{}
	ow := &delivery.SourceOverview{}
	ow.ID = detail.ID
	err := json.Unmarshal([]byte(detail.SourceData), data)
	if err != nil {
		return nil, nil, err
	}
	data.ID = detail.ID
	ow.Count = s.countTotal(delivery.SourceTypePackDataParams{
		ParamsEach:   data.ParamsEach,
		ParamsNumber: data.ParamsNumber,
		ParamsRule:   data.ParamsRule,
	})
	ow.PageNum = ow.Count
	return ow, data, nil
}

func (s packSource) Read(ctx ctx.Context, page int64, title []string, params interface{}, value []string) (error, [][]any) {
	client := &http.Client{}
	data := params.(*delivery.SourceTypePackData)
	req, err := s.getHttpRequest(page, data)
	if err != nil {
		return err, nil
	}
	response, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("接口返回%d", response.StatusCode)), nil
	}
	reader, err := common.GetResponseReader(response.Header.Get("Content-Encoding"), response.Body)
	defer response.Body.Close()
	if err != nil {
		return errors.New("读取response-body失败" + err.Error()), nil
	}
	// 创建名称为任务ID的打包目录
	folder := fmt.Sprintf("/tmp/%d", data.ID)
	lock.Lock()
	if !common.FileExists(folder) {
		err := os.Mkdir(folder, os.ModePerm)
		if err != nil {
			return errors.New("创建打包目录失败" + err.Error()), nil
		}
	}
	lock.Unlock()
	// 获取文件名称
	fileDesc := response.Header.Get("Content-Disposition")
	if fileDesc == "" {
		return errors.New("无法获取资源Disposition"), nil
	}
	reg := regexp.MustCompile(`filename=(.*\.\w+)`)
	matched := reg.FindAllStringSubmatch(fileDesc, 1)
	fileName := ""
	if len(matched) > 0 && len(matched[0]) >= 1 {
		fileName = fmt.Sprintf("%d-%s", page, matched[0][1])
	} else {
		return errors.New("无法从资源Disposition中获取文件名称"), nil
	}
	// 保存文件
	file, err := os.Create(folder + "/" + fileName)
	if err != nil {
		return errors.New("创建临时文件失败" + err.Error()), nil
	}
	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("读取数据流失败" + err.Error()), nil
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		return errors.New("写入临时文件失败" + err.Error()), nil
	}
	err = file.Close()
	if err != nil {
		return errors.New("关闭临时文件句柄失败" + err.Error()), nil
	}
	return nil, [][]any{}
}

func (s packSource) getHttpRequest(page int64, param *delivery.SourceTypePackData) (*http.Request, error) {
	// 请求头
	headers := map[string]string{}
	for key, value := range param.Headers {
		headers[key] = value
	}
	headers["Content-Type"] = "application/json"
	// 参数
	byteParams, err := json.Marshal(s.getParams(page, param))
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(param.Method, param.Path, bytes.NewBuffer(byteParams))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request, nil
}

func (s packSource) getParams(page int64, param *delivery.SourceTypePackData) map[string]interface{} {
	params := map[string]interface{}{}
	if param.ParamsRule == delivery.PackParamsRuleEach {
		for name, values := range param.ParamsEach {
			params[name] = values[page-1]
		}
	}
	if param.ParamsRule == delivery.PackParamsRuleIncrement {
		for name, value := range param.ParamsNumber {
			params[name] = value.Start + int64(value.Step)*(page-1)
		}
	}
	if param.ParamsRule == delivery.PackParamsRuleDecrement {
		for name, value := range param.ParamsNumber {
			params[name] = value.Start - int64(value.Step)*(page-1)
		}
	}
	return params
}

// countTotal 计算打包资源长度
func (s packSource) countTotal(params delivery.SourceTypePackDataParams) int64 {
	// 参数遍历形式
	if params.ParamsRule == delivery.PackParamsRuleEach {
		var item []string
		for _, itemC := range params.ParamsEach {
			item = itemC
			break
		}
		return int64(len(item))
	}
	var paramParams delivery.SourceTypePackDataParamsNumberItem
	for _, i2 := range params.ParamsNumber {
		paramParams = i2
		break
	}
	start := paramParams.Start
	count := 0
	// 自减形式
	if params.ParamsRule == delivery.PackParamsRuleDecrement {
		for start <= paramParams.End {
			start -= int64(paramParams.Step)
			count += paramParams.Step
		}
		return int64(count)
	}
	// 自加形式
	if params.ParamsRule == delivery.PackParamsRuleIncrement {
		for start >= paramParams.End {
			start += int64(paramParams.Step)
			count += paramParams.Step
		}
		return int64(count)
	}
	return 0
}
