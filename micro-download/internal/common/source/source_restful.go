package source

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"micro-download/internal/common"
	"micro-download/internal/config"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
)

// id 分段方式的一些默认配置
var (
	defaultRequestSortField           = "sort_key"              // 默认的排序key
	defaultSortKey                    = "id"                    // 排序字段
	defaultRequestSortTypeField       = "sort_type"             // 默认的排序方式的key
	defaultSortType                   = "ASC"                   // 这个是我们的默认排序方式，如果数据源不传数据的话，我们会在我们默认的排序字段上加上这个玩意
	defaultRequestSelectField         = "will_select"           // 这个是我们反查数据时会发送给我们的数据源我们要获取的数据，我们的数据源可以通过这个字段拿到key，只查询这些数据会减少一个网络io
	defaultRequestScrollPageInfoField = "will_scroll_page_info" // 这个是我们当前是否在获取我们的分页时传给数据源让他来确定我们在区间内的最远id
	defaultRequestMinID               = "will_min_id"
	defaultRequestMaxID               = "will_max_id"
	defaultSplitSize                  = common.DefaultSplitSize // 这个是我们的默认一页数据大小，默认跟普通模式是一样的，后续可以调整
	key                               = "REDACTED_FOR_SOURCE_SHARING"
	client                            = "download"
)

func (c *restfulSource) getField(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

// 接口调用类型
type restfulSource struct {
}

func (c *restfulSource) Create(ctx ctx.Context, export delivery.CreateExport) (error, int64) {
	if export.SourceRestfulHeaders == nil || export.SourceRestfulDataField == "" || export.SourceRestfulMethod == "" || export.SourceRestfulCountField == "" || export.SourceRestfulPath == "" || export.SourceRestfulRequestSizeField == "" || export.SourceRestfulRequestPageField == "" {
		return errors.New("无效的restful导出参数"), 0
	}
	model := &modes.DownloadSource{}
	model.SetParams(export)
	return model.SourceCreate(
		ctx, func(d *modes.DownloadSource) (err error, downloadSource *modes.DownloadSource) {
			// 分页类型默认为 1：页码
			if export.SourceRestfulPaginationType == 0 {
				export.SourceRestfulPaginationType = 1
			}
			data := &delivery.SourceTypeRestfulData{
				ExcelSourceDataCommonOption: delivery.ExcelSourceDataCommonOption{
					SourceTitleAlias:     export.SourceTitleAlias,
					SourceTitle:          export.SourceTitle,
					SourceTitleValue:     export.SourceTitleValue,
					SourceSuffix:         export.SourceSuffix,
					SourceExcelCellMerge: export.SourceExcelCellMerge,
					SourceHeaderPrefix:   export.SourceHeaderPrefix,
				},
				Path:                          export.SourceRestfulPath,
				Method:                        export.SourceRestfulMethod,
				Headers:                       export.SourceRestfulHeaders,
				ResponseCountField:            export.SourceRestfulCountField,
				ResponseDataField:             export.SourceRestfulDataField,
				RequestPageField:              export.SourceRestfulRequestPageField,
				RequestSizeField:              export.SourceRestfulRequestSizeField,
				Params:                        export.SourceRestfulParams,
				SourceRestfulIDPaginationType: export.SourceRestfulIDPaginationType,
			}
			byteData, err := json.Marshal(data)
			if err != nil {
				return err, d
			}
			d.SourceData = string(byteData)
			return nil, d
		},
	)
}

func (c *restfulSource) Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error) {
	data := &delivery.SourceTypeRestfulData{}
	err := json.Unmarshal([]byte(detail.SourceData), data)
	if err != nil {
		return nil, nil, err
	}
	result := &delivery.SourceOverview{
		Title:      data.SourceTitle,
		TitleAlias: data.SourceTitleAlias,
		Value:      data.SourceTitleValue,
		Suffix:     data.SourceSuffix,
		HeadPrefix: data.SourceHeaderPrefix,
		CellMerge:  data.SourceExcelCellMerge,
	}
	// 如果分页类型是ID，则根据总数动态计算页数
	if data.SourceRestfulPaginationType == delivery.RestfulPaginationTypeID {
		data.Params["source_restful_pagination_type"] = delivery.RestfulPaginationTypeID
		data.Params[c.getField(data.SourceRestfulScrollPageInfoField, defaultRequestScrollPageInfoField)] = true
		defer func() {
			data.Params[c.getField(data.SourceRestfulScrollPageInfoField, defaultRequestScrollPageInfoField)] = false
		}()
		// 拿到我们的miniid，往后不断迭代走我们的信息
		pfrom := make([]delivery.PageFromInfo, 0)
		var err error
		var id int64
		err, id = c.doHttpRequestGetID(ctx, id, 0, data)
		if err != nil {
			return nil, nil, err
		}
		if id != 0 {
			pfrom = append(
				pfrom, delivery.PageFromInfo{
					MinID: id,
				},
			)
		} else {
			result.PageNum = 0
			result.PageFrom = pfrom
			data.RequestFrom = result.PageFrom
			return result, data, nil
		}
		for {
			var nid int64
			err, nid = c.doHttpRequestGetID(ctx, id, int64(defaultSplitSize), data)
			if err != nil {
				return nil, nil, err
			}
			if nid == 0 {
				break
			}
			if nid == id {
				if len(pfrom) >= 2 {
					// 1. 推断我们之前的规律看看是减小的还是增大的
					p2MinID := pfrom[len(pfrom)-2].MinID
					p1MinID := pfrom[len(pfrom)-1].MinID
					if p2MinID == p1MinID {
						return nil, nil, errors.New("递归排序键失败")
					} else if p2MinID > p1MinID {
						nid--
					} else {
						nid++
					}
				} else {
					// 判断我们有没有传递排序，如果有就按照传递的排序来走
					if strings.ToUpper(data.SourceRestfulSortType) == "DESC" {
						nid--
					} else if strings.ToUpper(data.SourceRestfulSortType) == "ASC" {
						nid++
					} else {
						return nil, nil, errors.New("接口或者传递参数不满足id分段方式")
					}
				}
			}
			pfrom[len(pfrom)-1].MaxID = nid
			id = nid
			pfrom = append(
				pfrom, delivery.PageFromInfo{
					MinID: nid,
				},
			)
		}
		result.PageNum = int64(len(pfrom))
		result.PageFrom = pfrom
	} else {
		err, count, _ := c.doHttpRequest(ctx, 1, 1, nil, data, nil)
		if err != nil {
			return nil, nil, err
		}
		result.Count = count
		result.PageNum = common.PageCount(count)
	}
	data.RequestFrom = result.PageFrom
	return result, data, nil
}

func (c *restfulSource) Read(ctx ctx.Context, page int64, title []string, params interface{}, value []string) (error, [][]any) {
	data := params.(*delivery.SourceTypeRestfulData)
	// 进行首次请求 本次请求将确定 接口数据总数、表格title顺序
	err, _, list := c.doHttpRequest(ctx, page, common.DefaultSplitSize, title, data, value)
	return err, list
}

// pageOrFrom 普通模式，这就是page，id分段这个就是from
func (c *restfulSource) doRequest(ctx ctx.Context, pageOrFrom, size int64, data *delivery.SourceTypeRestfulData) (error, *map[string]interface{}) {
	retryDuration := make([]time.Duration, 0)
	for _, v := range config.CfgData.Downloader.RetrySleep {
		retryDuration = append(retryDuration, time.Duration(v)*time.Millisecond)
	}
	timeOut := time.Duration(config.CfgData.Downloader.Timeout) * time.Second
	if timeOut == 0 {
		timeOut = time.Second * 60 * 20
	}
	if len(retryDuration) == 0 {
		retryDuration = append(retryDuration, time.Millisecond*50, time.Millisecond*500, time.Millisecond*3000)
	}
	retryTimes := 0
RETRY:
	client := &http.Client{
		Timeout: timeOut,
	}
	err, request := c.getHttpRequest(ctx, pageOrFrom, size, data)
	if err != nil {
		return err, nil
	}
	startRequestTime := time.Now()
	retryCount := len(retryDuration)

	log.Info(ctx).Msgf("开始请求数据%d,%d", pageOrFrom, size)
	response, err := client.Do(request)
	if err != nil {
		log.Error(ctx).Msgf(
			"n-%d-请求数据失败%d,%d,耗时%s, %v", retryTimes, pageOrFrom, size, time.Since(startRequestTime), err,
		)
		if retryTimes < retryCount {
			time.Sleep(retryDuration[retryTimes])
			retryTimes++
			goto RETRY
		}
		return err, nil
	}
	defer response.Body.Close()

	log.Info(ctx).Msgf("请求数据成功%d,%d 耗时%s", pageOrFrom, time.Since(startRequestTime))
	reader, err := common.GetResponseReader(response.Header.Get("Content-Encoding"), response.Body)
	if err != nil {
		return err, nil
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err, nil
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("接口返回%d,返回body为:%s", response.StatusCode, string(body))), nil
	}

	var bodyM = &map[string]interface{}{}
	err = json.Unmarshal(body, bodyM)
	if err != nil {
		return err, nil
	}
	return nil, bodyM
}

// 返回我们要求的id
func (c *restfulSource) doHttpRequestGetID(ctx ctx.Context, from, offset int64, data *delivery.SourceTypeRestfulData) (error, int64) {
	err, info := c.doRequest(ctx, from, offset, data)
	if err != nil {
		return err, 0
	}

	if from == 0 && offset == 0 {
		for k, v := range data.SourceRestfulExt {
			err, i := c.getKeyFromData(k, *info)
			if err != nil {
				return err, 0
			}
			if i != nil {
				data.Params[v] = i
			}
		}
	}

	// 拿到我们的数据字段
	dataField := data.ResponseDataField
	// 从返回数据中解析我们的数据字段
	err, d := c.getKeyFromData(dataField, *info)
	if err != nil {
		return err, 0
	}
	// 转换为interface
	dataInterface := d.([]interface{})
	// 如果长度=0，证明符合要求的数据已经被拿完了
	if len(dataInterface) == 0 {
		return nil, 0
	}
	// 如果长度！=0，那我们只拿我们第一个就行，因为长度应该为1，而且我们是有排序的
	newData, ok := dataInterface[0].(map[string]interface{})
	if !ok {
		return errors.New("返回的数据不能解析为map[string]interface{}，不符合该方式的预期"), 0
	}
	// 1. 如果sortKey没有写，就按id排序
	// 2. 如果sortKey排序了，我们就要按照sortKey的顺序，但是sortKey可能包含一些.之类的，或者在查询时候有别名，所以
	// 最终我们提取的以来id应该是
	// 1. defaultSortKey
	// 2. sortKey
	// 3. data.SourceRestfulSortKeyAs
	sortField := ""
	if data.SourceRestfulSortKeyAs != "" {
		sortField = data.SourceRestfulSortKeyAs
	} else if data.SourceRestfulSortKey != "" {
		sortField = data.SourceRestfulSortKey
	} else {
		sortField = defaultSortKey
	}
	if idInterface, ok := newData[sortField]; ok {
		switch idT := idInterface.(type) {
		case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
			return nil, reflect.ValueOf(idInterface).Int()
		case float32, float64:
			return nil, int64(reflect.ValueOf(idInterface).Float())
		case string:
			id, err := strconv.ParseInt(idInterface.(string), 10, 64)
			return err, id
		default:
			return fmt.Errorf(
				"%s是一个%s，%v，不符合该方式的预期", sortField, reflect.TypeOf(idInterface).String(), idT,
			), 0
		}
	}
	return errors.New("没有找到我们想要的排序以来，不符合该方式的预期"), 0
}

// doHttpRequest 读取数据
func (c *restfulSource) doHttpRequest(ctx ctx.Context, page, size int64, title []string, data *delivery.SourceTypeRestfulData, value []string) (error, int64, [][]any) {
	err, info := c.doRequest(ctx, page, size, data)
	if err != nil {
		return err, 0, nil
	}
	return c.analysis(ctx, title, data.ResponseCountField, data.ResponseDataField, *info, value)
}

// 解析数据响应
func (c *restfulSource) analysis(ctx ctx.Context, title []string, countField, dataField string, body map[string]interface{}, value []string) (error, int64, [][]any) {
	// 处理count
	err, count := c.getKeyFromData(countField, body)
	var countC int64 = 0
	if err != nil {
		return err, 0, nil
	}
	// count 可能返回两种类型 string, int
	switch count.(type) {
	case string:
		countC, err = strconv.ParseInt(count.(string), 10, 64)
		if err != nil {
			return err, 0, nil
		}
		break
	case int64:
		countC = count.(int64)
		break
	case int:
		countC = int64(count.(int))
		break
	case float64:
		countC = int64(count.(float64))
		break
	case float32:
		countC = int64(float64(count.(float32)))
		break
	}
	// 处理data
	err, data := c.getKeyFromData(dataField, body)
	if err != nil {
		return err, 0, nil
	}
	var dataR [][]any
	dataInterface := data.([]interface{})
	for _, simpleData := range dataInterface {
		simpleDataC := simpleData.(map[string]interface{})
		var itemC []any
		// 导出列表可能返回三种类型 string, int, float
		for _, mapKey := range title {
			item := simpleDataC[mapKey]
			if common.StringInArray(mapKey, value) {
				itemC = append(itemC, item)
				continue
			}
			switch item.(type) {
			case string:
				itemC = append(itemC, item.(string))
				break
			case int:
			case int64:
			case int8:
				itemC = append(itemC, fmt.Sprintf("%d", item))
				break
			case float32:
				itemC = append(itemC, strconv.FormatFloat(item.(float64), 'f', -1, 32))
				break
			case float64:
				itemC = append(itemC, strconv.FormatFloat(item.(float64), 'f', -1, 64))
				break
			default:
				return errors.New(fmt.Sprintf("无法解析字段【%s】,因为列表返回无此字段", mapKey)), countC, nil
			}
		}
		dataR = append(dataR, itemC)
	}
	return nil, countC, dataR
}

// 解析response中的key
func (c *restfulSource) getKeyFromData(key string, data map[string]interface{}) (error, interface{}) {
	fieldKeys := strings.Split(key, ".")
	var fieldOperationObject = data
	var length = len(fieldKeys)
	var index = 0
	var dataR interface{}
start:
	value, ok := fieldOperationObject[fieldKeys[index]]
	if ok {
		if index == length-1 {
			dataR = value
		} else {
			fieldOperationObject = value.(map[string]interface{})
			index++
			goto start
		}
	} else {
		return errors.New(fmt.Sprintf("restful数据key[%s]解析失败", key)), nil
	}
	if dataR == nil {
		return errors.New(fmt.Sprintf("restful数据key[%s]解析失败", key)), nil
	}
	return nil, dataR
}

// getHttpRequest 构造分页请求
func (c *restfulSource) getHttpRequest(ctx ctx.Context, page, size int64, data *delivery.SourceTypeRestfulData) (error, *http.Request) {
	// 默认设置请求类型为json
	// 此处不可直接使用data.Headers 因为会存在并发写入map的问题导致painc
	headers := map[string]string{}
	for key, value := range data.Headers {
		headers[key] = value
	}
	headers["Content-Type"] = "application/json"
	// 添加分页参数
	params := map[string]interface{}{}
	for key, value := range data.Params {
		params[key] = value
	}
	// 如果是按照ID分页
	if data.SourceRestfulPaginationType == delivery.RestfulPaginationTypeID {
		params[c.getField(data.SourceRestfulSortField, defaultRequestSortField)] = c.getField(
			data.SourceRestfulSortKey, defaultSortKey,
		)
		params[c.getField(
			data.SourceRestfulSortTypeField, defaultRequestSortTypeField,
		)] = c.getField(data.SourceRestfulSortType, defaultSortType)
		params[data.RequestSizeField] = int(size)
		if data.SourceRestfulReceiveSelectSwitch {
			selectInfo := data.SourceRestfulSelectInfo
			if len(selectInfo) == 0 {
				selectInfo = data.SourceTitle
			}
			params[c.getField(data.SourceRestfulReceiveSelectField, defaultRequestSelectField)] = selectInfo
		}
		if data.RequestFrom == nil {
			params[c.getField(data.SourceRestfulMinIDField, defaultRequestMinID)] = page
			params[c.getField(data.SourceRestfulMaxIDField, defaultRequestMaxID)] = 0
		} else {
			params[c.getField(data.SourceRestfulMinIDField, defaultRequestMinID)] = data.RequestFrom[page].MinID
			params[c.getField(data.SourceRestfulMaxIDField, defaultRequestMaxID)] = data.RequestFrom[page].MaxID
		}
	} else {
		params[c.getField(data.SourceRestfulSortField, defaultRequestSortField)] = c.getField(
			data.SourceRestfulSortKey, defaultSortKey,
		)
		params[c.getField(
			data.SourceRestfulSortTypeField, defaultRequestSortTypeField,
		)] = c.getField(data.SourceRestfulSortType, defaultSortType)
		params[data.RequestPageField] = int(page)
		params[data.RequestSizeField] = int(size)
	}
	byteParam, err := json.Marshal(params)
	if err != nil {
		return err, nil
	}
	request, err := http.NewRequest(strings.ToUpper(data.Method), data.Path, bytes.NewBuffer(byteParam))
	if err != nil {
		return err, nil
	}
	log.Info(ctx).Msgf(
		"请求信息，url:%s, method: %s, params: %s", request.URL.String(), request.Method, string(byteParam),
	)
	// 设置请求头
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// todo key配置
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	m := md5.New()
	m.Write([]byte(client + key + curTime))
	signStr := fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
	request.Header.Set("Auth-Client", client)
	request.Header.Set("Auth-Sign", signStr)
	request.Header.Set("Auth-Timestamp", curTime)

	return nil, request
}

// calculateIDRangePageNum 计算ID分段页数
func (c *restfulSource) calculateIDRangePageNum(total int64) int64 {
	if total <= 100000 {
		return 1
	}
	// if total <= 500000 {
	// 	return 10
	// }
	if total <= 1500000 {
		return 5
	}
	return 20
}
