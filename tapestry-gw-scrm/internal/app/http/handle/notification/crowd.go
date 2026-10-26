package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	notificationReq "gw-scrm/internal/app/http/request/notification"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type CrowdHandle struct {
	handle.BaseHandle
	notificationSdk microscrm.NotificationApi
}

func NewCrowdHandle() CrowdHandle {
	return CrowdHandle{
		notificationSdk: microscrm.NewNotificationApi(),
	}
}

func (h CrowdHandle) Save(c *gin.Context) {
	var req notificationReq.CrowdReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, mobile := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminId = adminId
	req.Mobile = mobile
	res, err := h.notificationSdk.CrowdSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h CrowdHandle) Detail(c *gin.Context) {
	var req notificationReq.CrowdDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId

	res, err := h.notificationSdk.CrowdDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h CrowdHandle) Delete(c *gin.Context) {
	var req notificationReq.CrowdDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId

	res, err := h.notificationSdk.CrowdDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h CrowdHandle) List(c *gin.Context) {
	var req notificationReq.CrowdListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}
	res, err := h.notificationSdk.CrowdList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

type crowdImportRequest struct {
	FileURL string `json:"file_url" binding:"required"`
}

type crowdImportSummary struct {
	Import         int      `json:"import"`
	Success        int      `json:"success"`
	Repeat         int      `json:"repeat"`
	Fail           int      `json:"fail"`
	FailStoreCodes []string `json:"fail_store_code"`
}

type crowdImportMember struct {
	MemberID     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	DepartmentID int64  `json:"department_id"`
}

type crowdMemberLookup struct {
	StoreCodes []string            `json:"store_codes"`
	MemberList []crowdImportMember `json:"member_list"`
}

type crowdImportResponse struct {
	Result     crowdImportSummary  `json:"result"`
	MemberList []crowdImportMember `json:"member_list"`
}

func (h CrowdHandle) Import(c *gin.Context) {
	var req crowdImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	body, err := downloadCrowdWorkbook(req.FileURL)
	if err != nil {
		h.ResError(c, err)
		return
	}
	storeCodes, err := parseCrowdWorkbook(bytes.NewReader(body))
	if err != nil {
		h.ResError(c, err)
		return
	}
	uniqueCodes := uniqueStrings(storeCodes)
	enterpriseID, _ := common.GetEnterprise(c)
	lookupData, err := microscrm.GetMemberListByStoreCodes(h.WarpContext(c), map[string]interface{}{
		"enterprise_id": enterpriseID,
		"store_codes":   uniqueCodes,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	encoded, err := json.Marshal(lookupData)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var lookup crowdMemberLookup
	if err := json.Unmarshal(encoded, &lookup); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, crowdImportResponse{
		Result:     summarizeCrowdImport(storeCodes, lookup.StoreCodes),
		MemberList: lookup.MemberList,
	})
}

func (h CrowdHandle) Template(c *gin.Context) {
	writeCrowdWorkbook(c, "传讯_白名单_导入模版.xlsx", [][]interface{}{{"*门店code"}})
}

func (h CrowdHandle) ExportFailures(c *gin.Context) {
	var req struct {
		FailStoreCodes []string `json:"fail_store_codes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	rows := [][]interface{}{{"*门店code", "失败数据原因"}}
	for _, code := range req.FailStoreCodes {
		rows = append(rows, []interface{}{code, "未匹配到门店"})
	}
	writeCrowdWorkbook(c, "传讯_白名单_导入失败数据.xlsx", rows)
}

func parseCrowdWorkbook(reader io.Reader) ([]string, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, errors.New("文件解析失败，请使用下载的 xlsx 模板")
	}
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("文件中没有工作表")
	}
	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, errors.New("文件内容读取失败")
	}
	if len(rows) == 0 || len(rows[0]) == 0 || strings.TrimSpace(rows[0][0]) != "*门店code" {
		return nil, errors.New("文件格式错误，请使用下载的模板")
	}
	storeCodes := make([]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}
		if code := strings.TrimSpace(row[0]); code != "" {
			storeCodes = append(storeCodes, code)
		}
	}
	if len(storeCodes) == 0 {
		return nil, errors.New("导入数据为空")
	}
	return storeCodes, nil
}

func summarizeCrowdImport(imported, matched []string) crowdImportSummary {
	counts := make(map[string]int, len(imported))
	for _, code := range imported {
		counts[code]++
	}
	matchedSet := make(map[string]struct{}, len(matched))
	for _, code := range matched {
		matchedSet[code] = struct{}{}
	}
	failures := make([]string, 0)
	for _, code := range uniqueStrings(imported) {
		if _, exists := matchedSet[code]; !exists {
			failures = append(failures, code)
		}
	}
	sort.Strings(failures)
	repeats := 0
	for _, count := range counts {
		if count > 1 {
			repeats++
		}
	}
	return crowdImportSummary{
		Import: len(imported), Success: len(matchedSet), Repeat: repeats,
		Fail: len(failures), FailStoreCodes: failures,
	}
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func downloadCrowdWorkbook(rawURL string) ([]byte, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" {
		return nil, errors.New("文件地址无效")
	}
	host := strings.ToLower(parsed.Hostname())
	if host != "cdn-cos.coach.com.cn" && !strings.HasSuffix(host, ".myqcloud.com") {
		return nil, errors.New("文件地址不属于允许的对象存储")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	response, err := client.Get(parsed.String())
	if err != nil {
		return nil, errors.New("文件获取失败")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("文件获取失败：%s", response.Status)
	}
	return io.ReadAll(io.LimitReader(response.Body, 10<<20))
}

func writeCrowdWorkbook(c *gin.Context, filename string, rows [][]interface{}) {
	file := excelize.NewFile()
	for index, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, index+1)
		if err := file.SetSheetRow("Sheet1", cell, &row); err != nil {
			c.Error(err)
			return
		}
	}
	buffer, err := file.WriteToBuffer()
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(filename))
	_, _ = c.Writer.Write(buffer.Bytes())
}
