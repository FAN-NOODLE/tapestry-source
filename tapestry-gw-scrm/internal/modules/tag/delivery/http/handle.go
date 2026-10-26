package tag_http

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/logic"
	"gw-scrm/internal/modules/enterprise"
	"gw-scrm/internal/modules/file/storage"
	"gw-scrm/internal/modules/file/storage/cos"
	"gw-scrm/internal/modules/tag"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microgoods"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/sdk/microshop"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 获取企业的标签详情
func getTagGroupDetailList(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagGroupDetail(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// userCreateGroup 创建用户分组
func userCreateGroup(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.UserCreateGroup(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// userGroupList 获取客户分组列表
func userGroupList(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.UserGroupList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagList 标签分组列表
func tagList(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagSearch 标签/标签分组搜索
func tagSearch(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagSearch(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// createGroup 创建手动标签组
func createGroup(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.CreateGroup(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// createTag 新增手动标签
func createTag(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.CreateTag(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// updateAutoTag 更新自动标签设置
func updateAutoTag(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.UpdateAutoTag(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// goodsSpu 商品Spu列表
func goodsSpu(c *gin.Context) {
	var request tag.GoodsSpuRequest
	c.BindJSON(&request)

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	request.EnterpriseId = tokenInfo.EnterpriseId

	response, err := microgoods.GoodsSpuList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// wxCouponList 微信优惠券列表
func wxCouponList(c *gin.Context) {
	var request tag.CouponRequest
	c.BindJSON(&request)
	request.SearchIsAll = 1 // 查询所有

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	request.EnterpriseId = tokenInfo.EnterpriseId

	response, err := microshop.CouponList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagUpdate 编辑子标签
func tagUpdate(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagUpdate(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagInfo 标签组详情
func tagInfo(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagInfo(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagRand 标签组排序
func tagRand(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)
	response, err := microscrm.TagRand(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagGroupDelete 标签组删除
func tagGroupDelete(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.TagGroupDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagUpdateTime 刷新标签更新时间
func tagUpdateTime(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.TagUpdateTime(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagDelete 刷新标签更新时间
func tagDelete(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.TagDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// userGroupImport 客户分组模板导入
func userGroupImport(c *gin.Context) {
	//上传文件流
	formFile, _ := c.FormFile("file")
	openFile, err := formFile.Open()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var fileUrl string
	prefix := logic.GetBolbPrefix()
	prefix += "/micro-scrm/user-group/manual"
	//获取enterpriseID
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	//获取其他参数
	var dataMap = make(map[string]string, 0)
	d, _ := c.MultipartForm()
	for k, v := range d.Value {
		dataMap[k] = v[0]
	}
	for _, v2 := range d.File {
		dataMap["file_name"] = v2[0].Filename
	}
	ext := filepath.Ext(formFile.Filename)
	if ext != ".xlsx" {
		ginplus.ResError(c, errors.New("只允许上传以xlsx,xls后缀的文件"))
		return
	}
	zr, err := zip.NewReader(openFile, formFile.Size)
	if err != nil {
		ginplus.ResError(c, errors.New("文件格式错误"))
		return
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err = checkXlsxData(fileData.Sheets); err != nil {
		ginplus.ResError(c, err)
		return
	}

	//上传数据到blob
	var fileName = fmt.Sprintf("%d", time.Now().Unix()) + "_ConsumerFile" + ext

	var fileDetection tag.UserGroupFileDetection
	fileDetection.FileName = dataMap["file_name"]
	fileDetection.EnterpriseID = tokenInfo.EnterpriseId
	init := storage.Init(cos.COS)
	newFile, err := os.Create(fileName)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, err := ioutil.ReadAll(openFile)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if _, err0 := newFile.Write(bytes); err0 != nil {
		ginplus.ResError(c, err0)
		return
	}
	if err := newFile.Close(); err != nil {
		ginplus.ResError(c, err)
		return
	}

	upload, err := init.Upload(ctx.Wrap(c), fileName, fileName, prefix)
	if err != nil {
		os.Remove(fileName)
		ginplus.ResError(c, err)
		return
	}
	fileUrl = init.GetURL(upload)
	if err := os.Remove(fileName); err != nil {
		ginplus.ResError(c, err)
		return
	}
	//写入文件路径信息到数据库，获取id
	var tagDownloadInfo tag.DownLoadFileInfo
	tagDownloadInfo.EnterpriseID = tokenInfo.EnterpriseId
	tagDownloadInfo.FileUrl = fileUrl
	tagDownloadInfo.FileName = dataMap["file_name"]

	IDInfo, err := microscrm.UserGroupDownloadFileInfo(ctx.Wrap(c), &tagDownloadInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newID tag.IDInfo
	marshal2, err := json.Marshal(IDInfo)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(marshal2, &newID); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, newID)
	return
}

// 检查文件内容
func checkXlsxData(sheets []*xlsx.Sheet) error {
	for index, sheet := range sheets {
		if index > 0 {
			return errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri == 0 {
				if row.Cells[0].String() != "unionid" {
					return errors.New("上传的文件错误")
				}
			}
			if len(row.Cells) == 0 {
				continue
			}
		}
	}
	return nil
}

func GetUserGroupExport(c *gin.Context, request tag.ExportReq) (data []tag.ExportData, err error) {
	list, err := microscrm.UserGroupExportData(ctx.Wrap(c), request)
	if err != nil {
		return nil, err
	}
	var newList tag.ExportList
	byte, errJson := json.Marshal(list)
	if errJson != nil {
		return nil, errJson
	}
	if errs := json.Unmarshal(byte, &newList); err != nil {
		return nil, errs
	}
	return newList.ExportData, nil
}

// userGroupExport 客户分组客户导入的记录模板导出
func userGroupExport(c *gin.Context) {
	var request tag.ExportReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	request.Page = 1
	request.PageSize = 1000
	request.Total = int64(0)
	request.EnterpriseID = tokenInfo.EnterpriseId
	list, err := microscrm.UserGroupExportData(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newList tag.ExportList
	byte, err := json.Marshal(list)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(byte, &newList); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var datalist []tag.ExportData
	for _, item := range newList.ExportData {
		datalist = append(datalist, item)
	}
	limitEnd := request.Page * request.PageSize
	if int(newList.Total) > limitEnd {
		request.Total = newList.Total
		maxPage := int(math.Ceil(float64(newList.Total / int64(request.PageSize))))
		for i := 2; i <= maxPage; i++ {
			request.Page = i
			data, errs := GetUserGroupExport(c, request)
			if errs != nil {
				ginplus.ResError(c, err)
				return
			}
			for _, v := range data {
				datalist = append(datalist, v)
			}
		}
	}
	excelFile := excelize.NewFile()
	sheetName := "客户导入结果"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "B", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "B1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "unionID")
	_ = excelFile.SetCellStr(sheetName, "B1", "导入结果")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range datalist {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.UnionID)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.ImportRes)
		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "P"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())

	return
}

// tagDelete 刷新标签更新时间
func userGroupDetail(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.UserGroupDetail(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// userGroupEdit 刷新标签更新时间
func userGroupEdit(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.UserGroupEdit(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// userGroupDelete 刷新标签更新时间
func userGroupDelete(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.UserGroupDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagSyncQw 同步企业微信
func tagSyncQw(c *gin.Context) {
	var request interface{}
	c.BindJSON(&request)

	response, err := microscrm.TagSyncQw(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// tagExport 导出标签列表
func tagExport(c *gin.Context) {
	var request tag.TagExportRequest
	c.BindJSON(&request)
	tagGroups, err := microscrm.TagExport(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	excelFile := excelize.NewFile()
	excelFile.SetActiveSheet(0)
	sheetName := "标签统计"
	// 创建sheet
	excelFile.SetSheetName(excelFile.GetSheetName(0), sheetName)
	// 写入表头
	excelFile.SetCellStr(sheetName, "A1", "标签组名称")
	excelFile.SetCellStr(sheetName, "B1", "标签名称")
	excelFile.SetCellStr(sheetName, "C1", "标签人数")
	// 行开始
	rowStart := 2

	// 导出手动标签
	if request.TagType == tag.TAG_EXPORT_TYPE_MANUAL {
		for _, tagGroup := range tagGroups.Manual {
			tags, err := microscrm.TagInfo(ctx.Wrap(c), map[string]any{
				"tag_group_id": tagGroup.TagGroupId,
				"page":         1,
				"size":         9999,
			})
			if err != nil || tags.Total == 0 {
				continue
			}

			for _, tag := range tags.List {
				rowStartStr := strconv.Itoa(rowStart)
				excelFile.SetCellStr(sheetName, "A"+rowStartStr, tagGroup.TagGroupName)
				excelFile.SetCellStr(sheetName, "B"+rowStartStr, tag.TagName)
				excelFile.SetCellStr(sheetName, "C"+rowStartStr, strconv.Itoa(tag.UseNum))
				rowStart = rowStart + 1
			}
		}
	}

	excelFile.SetActiveSheet(0)
	// 文件名称 纳秒的md5 + 后缀名
	fileName := "标签统计" + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

func tagDetailList(c *gin.Context) {
	var request tag.InfoReq

	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.TagDetailList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func tagDetailExport(c *gin.Context) {
	var request tag.InfoReq
	if err := c.ShouldBindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	request.EnterpriseID = tokenInfo.EnterpriseId
	request.Page = 1
	request.Size = 1000
	list, err := microscrm.TagDetailList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	var newList tag.ExportInfo
	byte, err := json.Marshal(list)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if err := json.Unmarshal(byte, &newList); err != nil {
		ginplus.ResError(c, err)
		return
	}
	var datalist []tag.InfoUserList
	datalist = newList.UserList

	if len(newList.UserList) > 0 {
		for i := 2; ; i++ {
			request.Page = i
			data, errs := GetTagUserList(c, request)
			if errs != nil {
				ginplus.ResError(c, err)
				return
			}
			if len(data) == 0 {
				break
			}
			for _, item := range data {
				datalist = append(datalist, item)
			}
		}
	}
	excelFile := excelize.NewFile()
	sheetName := "客户名单"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	//设置宽度
	err = excelFile.SetColWidth(sheetName, "A", "C", 15)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	//设置样式
	style, styleErr := excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`)
	if nil != styleErr {
		ginplus.ResError(c, styleErr)
		return
	}
	err = excelFile.SetCellStyle(sheetName, "A1", "C1", style)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "客户昵称")
	_ = excelFile.SetCellStr(sheetName, "B1", "客户手机号")
	_ = excelFile.SetCellStr(sheetName, "C1", "unionID")
	startRow := 2
	//填充数据
	style, err = excelFile.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`)
	if nil != err {
		ginplus.ResError(c, err)
		return
	}
	for _, item := range datalist {
		rowStartStr := strconv.Itoa(startRow)
		_ = excelFile.SetCellStr(sheetName, "A"+rowStartStr, item.Name)
		_ = excelFile.SetCellStr(sheetName, "B"+rowStartStr, item.Mobile)
		_ = excelFile.SetCellStr(sheetName, "C"+rowStartStr, item.UnionID)
		startRow = startRow + 1
		err = excelFile.SetCellStyle(sheetName, "A"+rowStartStr, "C"+rowStartStr, style)
		if nil != err {
			ginplus.ResError(c, err)
			return
		}
	}
	excelFile.SetActiveSheet(index)
	fileName := time.Now().Format("200601021504") + ".xlsx"

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())

	return

}

func tagAuthCreate(c *gin.Context) {
	type request struct {
		EnterpriseID int64 `json:"enterprise_id"`
	}
	var req request
	tokenInfo := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	req.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.TagAuthCreate(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func GetTagUserList(c *gin.Context, request tag.InfoReq) (data []tag.InfoUserList, err error) {
	list, err := microscrm.TagDetailList(ctx.Wrap(c), request)
	if err != nil {
		return nil, err
	}

	var newList tag.ExportInfo
	byte, errJson := json.Marshal(list)
	if errJson != nil {
		return nil, errJson
	}
	if errs := json.Unmarshal(byte, &newList); err != nil {
		return nil, errs
	}
	data = newList.UserList
	return data, nil
}

// precalculateUserCount 客户分组预计算人数
func precalculateUserCount(c *gin.Context) {
	var request interface{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.GroupUserCount(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// test 测试获取企业信息
func test(c *gin.Context) {
	response := enterprise.Info(ctx.Wrap(c), 2)
	ginplus.ResSuccess(c, response)
}
