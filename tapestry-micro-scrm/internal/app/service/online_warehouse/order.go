package onlineWarehouseSvc

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"golang.org/x/image/draw"
	"gorm.io/gorm"
	"image"
	"image/png"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	onlineWarehouseDao "micro-scrm/internal/app/dao/online_warehouse"
	"micro-scrm/internal/app/http/request"
	onlineWarehouseReq "micro-scrm/internal/app/http/request/online_warehouse"
	onlineWarehouseResp "micro-scrm/internal/app/http/response/online_warehouse"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/app/service"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type OrderSvc struct {
	service.Base
	orderDao    *onlineWarehouseDao.OrderDao
	orderSkuDao *onlineWarehouseDao.OrderSkuDao
}

func NewOrderSvc() *OrderSvc {
	return &OrderSvc{
		orderDao:    onlineWarehouseDao.NewOrderDao(),
		orderSkuDao: onlineWarehouseDao.NewOrderSkuDao(),
	}
}

func readOnlineWarehouseImportFile(c ctx.Context, fileURL string) ([]byte, error) {
	return helper.ReadRemoteFile(c, fileURL)
}

func (s *OrderSvc) Create(c ctx.Context, req onlineWarehouseReq.CreateReq) (int64, error) {
	var err error
	orderSn := strings.TrimSpace(req.OrderSn)
	//one, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
	//	"order_sn": orderSn,
	//})
	//if err != nil {
	//	return 0, err
	//}
	//if one.ID > 0 {
	//	return 0, errors.New("销售单号重复")
	//}
	var skuList []onlineWarehouseModel.OnlineWarehouseOrderSkuModel
	orderId := app.UID()
	var totalQuantity int
	for _, sku := range req.SkuList {
		skuList = append(skuList, onlineWarehouseModel.OnlineWarehouseOrderSkuModel{
			OrderId:    orderId,
			OrderSn:    orderSn,
			UpcCode:    sku.UpcCode,
			SkuName:    sku.SkuName,
			SkuInfo:    sku.SkuInfo,
			Quantity:   sku.Quantity,
			ImageUrl:   sku.ImageUrl,
			Dept:       sku.Dept,
			StyleColor: sku.StyleColor,
			SkuNumber:  sku.SkuNumber,
			BaseModel:  models.CreateBase(app.UID()),
			Status:     onlineWarehouseModel.OrderStatusWait,
		})
		totalQuantity += sku.Quantity
	}
	err = s.orderSkuDao.CreateInBatchesData(c, req.EnterpriseId, skuList)
	if err != nil {
		return 0, err
	}
	info := onlineWarehouseModel.OnlineWarehouseOrderModel{
		OrderSn:        orderSn,
		ThirdId:        req.ThirdId,
		MemberID:       req.MemberId,
		TotalQuantity:  totalQuantity,
		BaseModel:      models.CreateBase(orderId),
		Status:         onlineWarehouseModel.OrderStatusWait,
		ExternalUserId: req.ExternalUserId,
		UnionId:        req.UnionId,
		UserName:       req.UserName,
	}
	err = s.orderDao.CreateData(c, req.EnterpriseId, info)
	_, _ = s.createSharePoster(c, req.EnterpriseId, orderId)
	return orderId, err
}

// 创建海报
func (s *OrderSvc) createSharePoster(c ctx.Context, enterpriseId int64, orderId int64) (string, error) {
	fileName := fmt.Sprintf("%s_%d.png", "online_poster", orderId)
	defer os.Remove(fileName)
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(c, map[string]interface{}{"enterprise_id": enterpriseId, "key": "web_h5_domain"})
	if err != nil {
		return fileName, err
	}
	if res.Name == "" {
		return fileName, err
	}
	domain := res.Name
	qrcodeUrl := fmt.Sprintf("%s/link?order_id=%d", domain, orderId)
	qrCode, err := qrcode.New(qrcodeUrl, qrcode.Medium)
	if err != nil {
		return fileName, err
	}
	qrCode.DisableBorder = true
	qrCodeSize := 360
	cs := storage.Init(cos.COS)
	cs.GetDomain()
	backUrl := fmt.Sprintf("%s/poster/online_order_back.png", cs.GetDomain())
	// 打开背景图片
	backgroundImageData, err := helper.ReadRemoteFile(c, backUrl)
	if err != nil {
		return fileName, err
	}

	backgroundImage, _, err := image.Decode(bytes.NewReader(backgroundImageData))
	if err != nil {
		return fileName, err
	}

	// 创建一个新的图像，用于合并二维码和背景图片
	mergedImage := image.NewRGBA(backgroundImage.Bounds())

	// 将背景图片绘制到新图像上
	draw.Draw(mergedImage, mergedImage.Bounds(), backgroundImage, image.Point{}, draw.Src)

	// 计算二维码在背景图片上的位置
	qrCodePosition := image.Point{
		X: 339,
		Y: 285,
	}

	// 将二维码绘制到新图像上
	qrCodeImage := qrCode.Image(qrCodeSize)
	draw.Draw(mergedImage, qrCodeImage.Bounds().Add(qrCodePosition), qrCodeImage, image.Point{}, draw.Over)
	// 创建一个新的文件，用于保存合并后的图片
	mergedImageFile, err := os.Create(fileName)
	if err != nil {
		return fileName, err
	}
	defer mergedImageFile.Close()

	// 将合并后的图片保存为 JPEG 格式wo
	err = png.Encode(mergedImageFile, mergedImage)
	if err != nil {
		return fileName, err
	}
	cosUrl, err := helper.UploadFileToCos(c, fileName, false)
	if err != nil {
		return fileName, err
	}
	err = s.orderDao.UpdateData(c, enterpriseId, map[string]interface{}{
		"id": orderId,
	}, map[string]any{
		"is_poster":        1,
		"share_poster_url": cosUrl,
	})
	return cosUrl, err
}

func (s *OrderSvc) SaveReceiveInfo(c ctx.Context, req onlineWarehouseReq.SaveReceiveInfoReq) error {
	var err error
	one, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id":     req.OrderId,
		"status": onlineWarehouseModel.OrderStatusWait,
	})
	if err != nil {
		return err
	}
	if one.ID == 0 {
		return errors.New("销售单号不存在")
	}
	if one.Status != onlineWarehouseModel.OrderStatusWait {
		return errors.New("订单状态不正确")
	}
	err = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id":     one.ID,
		"status": onlineWarehouseModel.OrderStatusWait,
	}, map[string]any{
		"receiver_name":    req.ReceiverName,
		"receiver_phone":   req.ReceiverPhone,
		"receiver_address": req.ReceiverAddress,
		"province":         req.Province,
		"city":             req.City,
		"county":           req.County,
		"status":           onlineWarehouseModel.OrderStatusPending,
	})
	err = s.orderSkuDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"order_id": one.ID,
	}, map[string]any{
		"status": onlineWarehouseModel.OrderStatusPending,
	})
	if err != nil {
		return err
	}
	return err
}

func (s *OrderSvc) List(c ctx.Context, req onlineWarehouseReq.ListReq) ([]onlineWarehouseResp.ListResp, int64, error) {
	resList := make([]onlineWarehouseResp.ListResp, 0)
	var (
		total int64
		err   error
	)
	storeMap := make(map[int]microWorkModel.StoreDepartmentModel)
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	departments, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	if err != nil {
		return nil, 0, err
	}
	departmentMap := make(map[int]string)
	for _, department := range departments {
		departmentMap[department.ThirdID] = department.Name
		if (req.StoreCode != "" && department.Name == req.StoreCode) || (req.StoreName != "" && strings.Contains(department.Name, req.StoreName)) {
			req.ThirdIds = append(req.ThirdIds, department.ThirdID)
		}
	}
	for _, storeInfo := range storeList {
		storeMap[int(storeInfo.WechatDepartmentID)] = storeInfo
		if (req.StoreCode != "" && storeInfo.ThirdStoreCode == req.StoreCode) || (req.StoreName != "" && strings.Contains(storeInfo.StoreName, req.StoreName)) {
			req.ThirdIds = append(req.ThirdIds, int(storeInfo.WechatDepartmentID))
		}
	}
	if (req.StoreCode != "" || req.StoreName != "") && len(req.ThirdIds) == 0 { // 根据门店名称和门店code没有找到门店，则返回空数据
		return resList, 0, err
	}
	// 获取数据
	list, total, err := s.orderDao.GetPagination(c, req, false)
	if err != nil {
		return resList, 0, err
	}
	//获取数据中的导购id集合
	var listMemberIds []int64

	for _, order := range list {
		listMemberIds = append(listMemberIds, order.MemberID)
	}
	if len(listMemberIds) == 0 {
		return resList, 0, err
	}
	//获取导购姓名
	memberList, err := membersDao.NewMembers().GetUnscopedMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"id":            listMemberIds,
	})
	if err != nil {
		return resList, total, err
	}
	memberMap := make(map[int64]string)
	for _, model := range memberList {
		memberMap[model.ID] = model.Name
	}

	for _, order := range list {
		var transferFile onlineWarehouseModel.TransferFileInfo
		helper.Unmarshal(order.TransferFile, &transferFile)
		transferStatus := helper.RuleThree(transferFile.Url == "", onlineWarehouseModel.TransferNO, onlineWarehouseModel.TransferYES)
		resList = append(resList, onlineWarehouseResp.ListResp{
			Id:              order.ID,
			OrderSn:         order.OrderSn,
			MemberId:        order.MemberID,
			ThirdId:         order.ThirdId,
			MemberName:      helper.RuleThreeForMap(memberMap, order.MemberID, memberMap[order.MemberID], ""),
			StoreCode:       helper.RuleThreeForMap(storeMap, order.ThirdId, storeMap[order.ThirdId].ThirdStoreCode, departmentMap[order.ThirdId]),
			StoreName:       helper.RuleThreeForMap(storeMap, order.ThirdId, storeMap[order.ThirdId].StoreName, departmentMap[order.ThirdId]),
			ReceiverName:    order.ReceiverName,
			TrackingNumber:  order.TrackingNumber,
			ShippedAt:       order.ShippedAt,
			ReceiverPhone:   order.ReceiverPhone,
			ReceiverAddress: order.ReceiverAddress,
			TransferFile:    transferFile,
			Status:          int(order.Status),
			TransferStatus:  transferStatus,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
			TotalQuantity:   order.TotalQuantity,
			SendQuantity:    order.SendQuantity,
			UserName:        order.UserName,
			ExternalUserId:  order.ExternalUserId,
		})
	}

	return resList, total, err
}

// Submit 订单操作
func (s *OrderSvc) Submit(c ctx.Context, req onlineWarehouseReq.SubmitOrderReq) error {
	refusalList := req.RefusalInfo
	orderInfo, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.OrderId,
	})
	if err != nil {
		return err
	}
	skuList, err := s.orderSkuDao.GetMany(c, req.EnterpriseId, map[string]interface{}{
		"order_id": req.OrderId,
	})
	if err != nil {
		return err
	}
	refusalSkuMap := make(map[string]onlineWarehouseReq.SkuList)
	for _, sku := range refusalList {
		refusalSkuMap[sku.SkuInfo] = sku
	}
	skuUpdateMap := make(map[string]onlineWarehouseModel.OnlineWarehouseOrderSkuModel) // 更新数据
	skuInsertMap := make(map[string]onlineWarehouseModel.OnlineWarehouseOrderSkuModel) // 插入数据
	refusalNum := 0
	allNum := 0
	sendQuantity := 0
	for _, dbSku := range skuList {
		allNum += dbSku.Quantity
		refusalSku, ok := refusalSkuMap[dbSku.SkuInfo]
		if !ok { // 不在拒单中的商品全部更新为已发货
			dbSku.Status = onlineWarehouseModel.OrderStatusShipped
			skuUpdateMap[dbSku.SkuInfo] = dbSku
			continue
		}
		if refusalSku.Quantity > dbSku.Quantity {
			return errors.New(fmt.Sprintf("【%s】拒单数量大于下单数量", refusalSku.SkuInfo))
		}
		if refusalSku.Quantity == dbSku.Quantity { // 拒单数量和下单数量一至，全部更新为已拒单
			dbSku.Status = onlineWarehouseModel.OrderStatusRefusal
			skuUpdateMap[dbSku.SkuInfo] = dbSku
			refusalNum += refusalSku.Quantity
		}
		if refusalSku.Quantity < dbSku.Quantity { // 拒单数量小于下单数据，部分更新，部分插入
			insertQuantity := dbSku.Quantity - refusalSku.Quantity
			updateQuantity := refusalSku.Quantity
			insertSku := dbSku
			updateSku := dbSku
			insertSku.Quantity = insertQuantity
			insertSku.Status = onlineWarehouseModel.OrderStatusShipped
			updateSku.Quantity = updateQuantity
			updateSku.Status = onlineWarehouseModel.OrderStatusRefusal
			skuUpdateMap[dbSku.SkuInfo] = updateSku
			insertSku.ID = app.UID()
			insertSku.UpdatedAt = time.Now().Unix()
			insertSku.CreatedAt = time.Now().Unix()
			skuInsertMap[dbSku.SkuInfo] = insertSku

		}
	}
	insertList := make([]onlineWarehouseModel.OnlineWarehouseOrderSkuModel, 0)
	for _, sku := range skuInsertMap {
		insertList = append(insertList, sku)
	}
	if len(insertList) > 0 {
		err = s.orderSkuDao.CreateInBatchesData(c, req.EnterpriseId, insertList)
		if err != nil {
			return err
		}
	}
	for _, sku := range skuUpdateMap {
		if sku.Status == onlineWarehouseModel.OrderStatusShipped {
			sendQuantity += sku.Quantity
		}
		err = s.orderSkuDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
			"id": sku.ID,
		}, map[string]interface{}{
			"status":   sku.Status,
			"quantity": sku.Quantity,
		})
		if err != nil {
			return err
		}
	}
	orderStatus := onlineWarehouseModel.OrderStatusShipped
	if refusalNum == allNum {
		orderStatus = onlineWarehouseModel.OrderStatusRefusal
	}
	nowTime := time.Now().Unix()
	updateMap := map[string]interface{}{
		"status":          orderStatus,
		"tracking_number": req.TrackingNumber,
		"shipped_at":      time.Now().Unix(),
		"send_quantity":   sendQuantity,
		"operate_time":    nowTime,
		"admin_id":        req.AdminID,
	}
	if req.TransferFile.Url != "" {
		updateMap["transfer_file"] = helper.MarshalToString(req.TransferFile)
	}
	err = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": req.OrderId,
	}, updateMap)
	sendOrderSn := make([]int64, 0)
	sendOrderSn = append(sendOrderSn, req.OrderId)
	if orderInfo.Status == onlineWarehouseModel.OrderStatusPending { // 这里订单只有是待发货的才提醒消息
		go s.sendUpdateNotice(c, req.EnterpriseId, sendOrderSn, onlineWarehouseModel.OrderStatusShipped) // 发送消息通知
	}
	return err
}

// Info 查看订单详情
func (s *OrderSvc) Info(c ctx.Context, req onlineWarehouseReq.InfoReq) (info onlineWarehouseResp.InfoResp, err error) {
	orderInfo, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.OrderId,
	})
	if err != nil {
		return info, err
	}
	info.Id = orderInfo.ID
	info.TrackingNumber = orderInfo.TrackingNumber
	info.OrderSn = orderInfo.OrderSn
	info.ShippedAt = orderInfo.ShippedAt
	info.CreatedAt = orderInfo.CreatedAt
	info.UpdatedAt = orderInfo.UpdatedAt
	info.Status = orderInfo.Status
	info.UserName = orderInfo.UserName
	info.ExternalUserId = orderInfo.ExternalUserId
	info.ReceiverPhone = orderInfo.ReceiverPhone
	info.ReceiverName = orderInfo.ReceiverName
	info.ReceiverAddress = orderInfo.Province + orderInfo.City + orderInfo.County + orderInfo.ReceiverAddress
	info.SharePosterUrl = orderInfo.SharePosterUrl
	info.OperateTime = orderInfo.OperateTime
	address := fmt.Sprintf("%s%s%s%s", orderInfo.Province, orderInfo.City, orderInfo.County, orderInfo.ReceiverAddress)
	if req.Origin == 1 {
		info.ReceiverPhone = helper.HandleMobile(orderInfo.ReceiverPhone)
		info.ReceiverName = helper.DesensitizeNickname(orderInfo.ReceiverName)
		info.ReceiverAddress = helper.DesensitizeAddress(address, 11)
	}
	skuMany, err := s.orderSkuDao.GetMany(c, req.EnterpriseId, map[string]interface{}{
		"order_id": orderInfo.ID,
	})
	if err != nil {
		return info, err
	}
	skuMap := make(map[string]onlineWarehouseResp.InfoSku)
	for _, sku := range skuMany {
		key := sku.SkuInfo
		if req.Origin == 1 {
			key = sku.SkuInfo + "_" + strconv.Itoa(int(sku.Status))
		}
		item, ok := skuMap[key]
		info.AllNum += sku.Quantity
		var infoSku onlineWarehouseResp.InfoSku
		if !ok {
			infoSku.UpcCode = sku.UpcCode
			infoSku.SkuName = sku.SkuName
			infoSku.SkuInfo = sku.SkuInfo
			infoSku.ImageUrl = sku.ImageUrl
			infoSku.Quantity = sku.Quantity
		} else {
			infoSku = item
			infoSku.Quantity += sku.Quantity
		}
		if sku.Status == onlineWarehouseModel.OrderStatusRefusal {
			infoSku.RefusalQuantity += sku.Quantity
			info.RefusalNum += sku.Quantity
		}
		infoSku.Status = int(sku.Status)
		skuMap[key] = infoSku
	}
	for _, sku := range skuMap {
		info.SkuList = append(info.SkuList, sku)
	}
	storeInfo, err := microWorkDao.NewStoreDepartment().GetInfo(c, map[string]interface{}{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": orderInfo.ThirdId,
	})
	if err != nil {
		return info, err
	}
	info.StoreCode = storeInfo.ThirdStoreCode
	info.StoreName = storeInfo.StoreName

	if storeInfo.ThirdStoreCode == "" {
		department, err := departmentsDao.NewDepartments().GetOne(c, map[string]interface{}{
			"enterprise_id": req.EnterpriseId,
			"third_id":      orderInfo.ThirdId,
		})
		if err != nil {
			return info, err
		}
		info.StoreCode = department.Name
		info.StoreName = department.Name
	}

	memberInfo, err := membersDao.NewMembers().GetOneByUnscoped(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"id":            orderInfo.MemberID,
	})
	if err != nil {
		return info, err
	}
	info.MemberName = memberInfo.Name
	info.MemberId = orderInfo.MemberID
	helper.Unmarshal(orderInfo.TransferFile, &info.TransferFile)
	return
}

// CustomerOrderInfo 顾客查看订单
func (s *OrderSvc) CustomerOrderInfo(c ctx.Context, req onlineWarehouseReq.OrderInfoReq) (result onlineWarehouseResp.CustomerOrderInfoResp, err error) {
	var info onlineWarehouseResp.InfoResp
	orderInfo, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.Id,
	})
	if err != nil {
		return result, err
	}
	result.OrderSn = orderInfo.OrderSn
	if helper.InArray(orderInfo.Status, []int32{onlineWarehouseModel.OrderStatusPending, onlineWarehouseModel.OrderStatusShipped, onlineWarehouseModel.OrderStatusRefusal}) {
		result.OrderStatus = 1
		return result, err
	}
	if orderInfo.Status == onlineWarehouseModel.OrderStatusCancel {
		result.OrderStatus = 2
		return result, err
	}
	info.Id = orderInfo.ID
	info.TrackingNumber = orderInfo.TrackingNumber
	info.OrderSn = orderInfo.OrderSn
	info.ReceiverPhone = orderInfo.ReceiverPhone
	info.ReceiverName = orderInfo.ReceiverName
	info.ReceiverAddress = orderInfo.ReceiverAddress
	info.ShippedAt = orderInfo.ShippedAt
	info.CreatedAt = orderInfo.CreatedAt
	info.UpdatedAt = orderInfo.UpdatedAt
	info.Status = orderInfo.Status
	skuMany, err := s.orderSkuDao.GetMany(c, req.EnterpriseId, map[string]interface{}{
		"order_id": orderInfo.ID,
	})
	if err != nil {
		return result, err
	}
	skuMap := make(map[string]onlineWarehouseResp.InfoSku)
	for _, sku := range skuMany {
		item, ok := skuMap[sku.SkuInfo]
		info.AllNum += sku.Quantity
		var infoSku onlineWarehouseResp.InfoSku
		if !ok {
			infoSku.UpcCode = sku.UpcCode
			infoSku.SkuName = sku.SkuName
			infoSku.SkuInfo = sku.SkuInfo
			infoSku.ImageUrl = sku.ImageUrl
			infoSku.Quantity = sku.Quantity
		} else {
			infoSku = item
			infoSku.Quantity += sku.Quantity
		}
		if sku.Status == onlineWarehouseModel.OrderStatusRefusal {
			infoSku.RefusalQuantity += sku.Quantity
			info.RefusalNum += sku.Quantity
		}
		infoSku.Status = int(sku.Status)
		skuMap[sku.SkuInfo] = infoSku
	}
	for _, sku := range skuMap {
		info.SkuList = append(info.SkuList, sku)
	}
	result.OrderInfo = info
	return
}

func (s *OrderSvc) Edit(c ctx.Context, req onlineWarehouseReq.SubmitOrderReq) error {
	return s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": req.OrderId,
	}, map[string]interface{}{
		"tracking_number": req.TrackingNumber,
		"transfer_file":   helper.MarshalToString(req.TransferFile),
	})
}

func (s *OrderSvc) StoreList(c ctx.Context, req onlineWarehouseReq.StoreListReq) ([]onlineWarehouseResp.StoreListResp, int64, error) {
	resList := make([]onlineWarehouseResp.StoreListResp, 0)
	var total int64
	var thirdIds []int
	storeMap := make(map[int]microWorkModel.StoreDepartmentModel)
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	for _, storeInfo := range storeList {
		storeMap[int(storeInfo.WechatDepartmentID)] = storeInfo
		if (req.StoreCode != "" && storeInfo.ThirdStoreCode == req.StoreCode) || (req.StoreName != "" && strings.Contains(storeInfo.StoreName, req.StoreName)) {
			thirdIds = append(thirdIds, int(storeInfo.WechatDepartmentID))
		}
	}
	if (req.StoreCode != "" || req.StoreName != "") && len(thirdIds) == 0 { // 根据门店名称和门店code没有找到门店，则返回空数据
		return resList, 0, err
	}
	resList, total, err = s.orderDao.GetPaginationGroupBYStore(c, req, thirdIds)
	if err != nil {
		return resList, total, err
	}
	for i, item := range resList {
		item.StoreCode = storeMap[item.ThirdId].ThirdStoreCode
		item.StoreName = storeMap[item.ThirdId].StoreName
		resList[i] = item
	}
	return resList, total, err
}

func (s *OrderSvc) StoreOrderList(c ctx.Context, req onlineWarehouseReq.StoreOrderListReq) ([]onlineWarehouseResp.StoreOrderListResp, int64, error) {
	resList := make([]onlineWarehouseResp.StoreOrderListResp, 0)
	memberList, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	if err != nil {
		return resList, 0, err
	}
	memberMap := make(map[int64]string)
	for _, model := range memberList {
		memberMap[model.ID] = model.Name
	}
	list, total, err := s.orderDao.GetStoreOrderPagination(c, req)
	if err != nil {
		return resList, 0, err
	}
	for _, item := range list {
		resList = append(resList, onlineWarehouseResp.StoreOrderListResp{
			OrderSn:       item.OrderSn,
			MemberId:      item.MemberID,
			MemberName:    memberMap[item.MemberID],
			ShippedAt:     item.ShippedAt,
			CreatedAt:     item.CreatedAt,
			TotalQuantity: item.TotalQuantity,
		})
	}
	return resList, total, err
}

func (s *OrderSvc) ExportFile(c ctx.Context, req onlineWarehouseReq.ListReq) (string, error) {
	fileName := "订单列表明细" + time.Now().Format("20060102150405") + ".xlsx"
	resList := make([]onlineWarehouseResp.ExportListResp, 0)
	storeMap := make(map[int]microWorkModel.StoreDepartmentModel)
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	for _, storeInfo := range storeList {
		storeMap[int(storeInfo.WechatDepartmentID)] = storeInfo
		if (req.StoreCode != "" && storeInfo.ThirdStoreCode == req.StoreCode) || (req.StoreName != "" && strings.Contains(storeInfo.StoreName, req.StoreName)) {
			req.ThirdIds = append(req.ThirdIds, int(storeInfo.WechatDepartmentID))
		}
	}
	departments, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	if err != nil {
		return fileName, err
	}
	departmentMap := make(map[int]string)
	for _, department := range departments {
		departmentMap[department.ThirdID] = department.Name
		if (req.StoreCode != "" && department.Name == req.StoreCode) || (req.StoreName != "" && strings.Contains(department.Name, req.StoreName)) {
			req.ThirdIds = append(req.ThirdIds, department.ThirdID)
		}
	}
	//获取导购姓名
	memberList, err := membersDao.NewMembers().GetUnscopedMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
	})
	if err != nil {
		return fileName, err
	}
	memberMap := make(map[int64]string)
	for _, model := range memberList {
		memberMap[model.ID] = model.Name
	}

	isImport := false
	if req.ExportType == "import" {
		isImport = true
	}
	excelFile := excelize.NewFile()
	sheetName := "订单列表明细"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	err = excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return fileName, err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return fileName, err
	}
	row := make([]any, 0)
	row = []any{
		"1、是否拒单：1代表拒单，0代表不拒单\r\n2、导入订单物流号时若是否拒单为1，则需填写SKU对应的拒单数量\r\n3、导入订单物流号时，当SKU拒单数量小于商品数量时，系统默认部分发货，需填写顺丰单号",
	}
	styleId, err := excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true},
	})
	if err != nil {
		return fileName, err
	}
	err = stream.SetRow("A1", row, excelize.RowOpts{
		Height:  50,
		StyleID: styleId,
	})
	hightLineStyle, err := excelFile.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"FFFF00"},
		},
	})
	linkStyle, err := excelFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color:     "0563C1",
			Underline: "single",
		},
	})
	if err != nil {
		return fileName, err
	}
	if req.IsH5 {
		row = []any{
			"销售单号", excelize.Cell{Value: "Dept", StyleID: hightLineStyle}, excelize.Cell{Value: "Style color", StyleID: hightLineStyle}, excelize.Cell{Value: "new format SKU", StyleID: hightLineStyle}, "SKU", "是否拒单", "拒单数量", "顺丰单号", "转货单", "商品数量", "门店代码", "门店名称", "订单状态", "创建人", "创建时间", "更新时间",
		}
	} else {
		row = []any{
			"销售单号", "SKU", "UPC_code", "是否拒单", "拒单数量", "顺丰单号", "转货单", "商品数量",
			"收件人", "收件电话", "省", "市", "区", "详细地址",
			"门店代码", "门店名称", "订单状态", "创建人", "创建时间", "更新时间",
		}
	}
	if isImport {
		row = []any{
			"销售单号", "SKU", "UPC_code", excelize.Cell{Value: "是否拒单", StyleID: hightLineStyle}, excelize.Cell{Value: "拒单数量", StyleID: hightLineStyle}, excelize.Cell{Value: "顺丰单号", StyleID: hightLineStyle}, "商品数量",
			"收件人", "收件电话", "省", "市", "区", "详细地址",
			"门店代码", "门店名称", "订单状态", "创建人", "创建时间", "更新时间", "后台编号(用于相同销售单号区分发货,导入时勿删)",
		}
	}

	err = stream.SetRow("A2", row)
	if err != nil {
		return fileName, err
	}
	startRow := int64(3)
	wg := threading.RoutineGroup{}
	orderStatusMap := map[int]string{
		onlineWarehouseModel.OrderStatusWait:    "待客户填写",
		onlineWarehouseModel.OrderStatusPending: "待发货",
		onlineWarehouseModel.OrderStatusShipped: "已发货",
		onlineWarehouseModel.OrderStatusRefusal: "已拒单",
		onlineWarehouseModel.OrderStatusCancel:  "已作废",
	}
	req.Page = 0
	req.Size = 10000
	cs := new(cos.Storage)

	for {
		req.Page++
		// 获取数据
		var updateExportIds []int64
		list, _, err := s.orderDao.GetPagination(c, req, true)
		if err != nil || (len(list) == 0) {
			break
		}
		for _, order := range list {
			var transferFile onlineWarehouseModel.TransferFileInfo
			transferStatus := helper.RuleThree(order.TransferFile == "", onlineWarehouseModel.TransferNO, onlineWarehouseModel.TransferYES)
			helper.Unmarshal(order.TransferFile, &transferFile)
			if len(order.SkuLists) > 0 {
				//skuInfoMap := make(map[string]onlineWarehouseResp.SkuExportList)
				for _, sku := range order.SkuLists {
					//value, ok := skuInfoMap[sku.UpcCode]
					//if !ok {
					exportInfo := onlineWarehouseResp.SkuExportList{
						UpcCode:      sku.UpcCode,
						SkuName:      sku.SkuName,
						SkuInfo:      sku.SkuInfo,
						ImageUrl:     sku.ImageUrl,
						Quantity:     sku.Quantity,
						SkuStatus:    int(sku.Status),
						Dept:         sku.Dept,
						StyleColor:   sku.StyleColor,
						NewFormatSku: sku.SkuNumber,
					}
					if sku.Status == onlineWarehouseModel.OrderStatusRefusal {
						exportInfo.RefusalQuantity = sku.Quantity
					}
					//skuInfoMap[sku.UpcCode] = exportInfo
					//} else {
					//	value.Quantity += sku.Quantity
					//	if sku.Status == onlineWarehouseModel.OrderStatusRefusal {
					//		value.RefusalQuantity += sku.Quantity
					//	}
					//	skuInfoMap[sku.UpcCode] = value
					//}
					resList = append(resList, onlineWarehouseResp.ExportListResp{
						ListResp: onlineWarehouseResp.ListResp{
							Id:              order.ID,
							OrderSn:         order.OrderSn,
							MemberId:        order.MemberID,
							MemberName:      helper.RuleThreeForMap(memberMap, order.MemberID, memberMap[order.MemberID], ""),
							StoreCode:       helper.RuleThreeForMap(storeMap, order.ThirdId, storeMap[order.ThirdId].ThirdStoreCode, departmentMap[order.ThirdId]),
							StoreName:       helper.RuleThreeForMap(storeMap, order.ThirdId, storeMap[order.ThirdId].StoreName, departmentMap[order.ThirdId]),
							ReceiverName:    order.ReceiverName,
							TrackingNumber:  order.TrackingNumber,
							ShippedAt:       order.ShippedAt,
							ReceiverPhone:   order.ReceiverPhone,
							ReceiverAddress: order.ReceiverAddress,
							TransferFile:    transferFile,
							Status:          int(order.Status),
							TransferStatus:  transferStatus,
							CreatedAt:       order.CreatedAt,
							UpdatedAt:       order.UpdatedAt,
							TotalQuantity:   order.TotalQuantity,
							Province:        order.Province,
							City:            order.City,
							County:          order.County,
						},
						SkuExportList: exportInfo,
					})
				}
			}
			if isImport {
				if order.IsExport == 0 {
					updateExportIds = append(updateExportIds, order.ID)
				}
			}
		}
		wg.RunSafe(func() {
			for _, item := range resList {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				refusalStatus := 0
				if item.RefusalQuantity > 0 {
					refusalStatus = 1
				}
				createTime := helper.GetFormatDateByParams(item.CreatedAt)
				updateTime := helper.GetFormatDateByParams(item.UpdatedAt)
				linkUri := ""
				if isImport {
					id := fmt.Sprintf("'%d", item.Id)
					row = []any{
						item.OrderSn, item.SkuInfo, item.UpcCode, refusalStatus, item.RefusalQuantity, item.TrackingNumber, item.Quantity,
						item.ReceiverName, item.ReceiverPhone, item.Province, item.City, item.County, item.ReceiverAddress,
						item.StoreCode, item.StoreName, orderStatusMap[item.Status], item.MemberName, createTime, updateTime, id,
					}
				} else {
					transferInfo := ""
					var year, month, day string
					linkUrl := item.TransferFile.Url
					if linkUrl != "" {
						parsedURL, err := url.Parse(linkUrl)
						if err != nil {
							fmt.Println("Error parsing URL:", err)
							return
						}

						// 获取路径部分
						path := parsedURL.Path

						// 使用正则表达式匹配日期
						re := regexp.MustCompile(`/(\d{4})/(\d{2})/(\d{2})/`)
						matches := re.FindStringSubmatch(path)

						if len(matches) == 4 {
							year, month, day = matches[1], matches[2], matches[3]
						}
						linkUri, err = cs.GetPresignedURL(c, linkUrl)
						if err == nil {
							transferInfo = fmt.Sprintf("转货单_%s_%s_%s_%s", item.StoreCode, year, month, day)
						}

					}
					if req.IsH5 {
						row = []any{
							item.OrderSn, excelize.Cell{Value: item.Dept, StyleID: hightLineStyle}, excelize.Cell{Value: item.StyleColor, StyleID: hightLineStyle}, excelize.Cell{Value: item.NewFormatSku, StyleID: hightLineStyle}, item.SkuInfo, refusalStatus, item.RefusalQuantity, item.TrackingNumber, excelize.Cell{Value: transferInfo, StyleID: linkStyle}, item.Quantity, item.StoreCode, item.StoreName,
							orderStatusMap[item.Status], item.MemberName, createTime, updateTime,
						}
						excelFile.SetCellHyperLink(sheetName, "I"+rowStartStr, linkUri, "External")
					} else {
						row = []any{
							item.OrderSn, item.SkuInfo, item.UpcCode, refusalStatus, item.RefusalQuantity, item.TrackingNumber, excelize.Cell{Value: transferInfo, StyleID: linkStyle}, item.Quantity,
							item.ReceiverName, item.ReceiverPhone, item.Province, item.City, item.County, item.ReceiverAddress,
							item.StoreCode, item.StoreName, orderStatusMap[item.Status], item.MemberName, createTime, updateTime,
						}
						excelFile.SetCellHyperLink(sheetName, "G"+rowStartStr, linkUri, "External")
					}
				}
				err = stream.SetRow("A"+rowStartStr, row)
				atomic.AddInt64(&startRow, 1)
			}
		})

		wg.RunSafe(func() {
			if len(updateExportIds) > 0 {
				_ = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]any{
					"id": updateExportIds,
				}, map[string]interface{}{
					"is_export": onlineWarehouseModel.OrderExportYes,
				})
			}
		})

	}
	wg.Wait()
	err = stream.MergeCell("A1", "T1")
	if err != nil {
		return fileName, err
	}
	err = stream.Flush()
	if err != nil {
		return fileName, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return fileName, err
	}
	return fileName, err
}

func (s *OrderSvc) ImportOrder(c ctx.Context, req onlineWarehouseReq.ImportOrderReq) (onlineWarehouseResp.ImportOrderResp, error) {
	nowTime := time.Now().Unix()
	var result onlineWarehouseResp.ImportOrderResp
	fileInfo, err := readOnlineWarehouseImportFile(c, req.FileUrl)
	if err != nil {
		return result, err
	}
	var data map[int64][]onlineWarehouseReq.ImportOrderData
	zr, err := zip.NewReader(bytes.NewReader(fileInfo), int64(len(fileInfo)))
	if err != nil {
		return result, err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return result, err
	}
	enterpriseId := req.EnterpriseId
	var sendOrderSn []int64
	data, result, err = s.readXlsxData(c, req.EnterpriseId, fileData.Sheets)
	if err != nil {
		return result, err
	}
	fileName := fmt.Sprintf("网仓订单发货失败_%d", time.Now().Unix()) + ".xlsx"
	excelFile := excelize.NewFile()
	sheetName := "网仓订单发货明细"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	err = excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return result, err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return result, err
	}
	row := make([]any, 0)
	row = []any{
		"1、是否拒单：1代表拒单，0代表不拒单\r\n2、导入订单物流号时若是否拒单为1，则需填写SKU对应的拒单数量\r\n3、导入订单物流号时，当SKU拒单数量小于商品数量时，系统默认部分发货，需填写顺丰单号",
	}
	styleId, err := excelFile.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true},
	})
	if err != nil {
		return result, err
	}
	err = stream.SetRow("A1", row, excelize.RowOpts{
		Height:  50,
		StyleID: styleId,
	})
	row = []any{
		"销售单号", "SKU", "UPC_code", "是否拒单", "拒单数量", "顺丰单号", "商品数量",
		"收件人", "收件电话", "省", "市", "区", "详细地址",
		"门店代码", "门店名称", "订单状态", "创建人", "创建时间", "更新时间", "后台编号(用于相同销售单号区分发货,导入时勿删)", "数据失败原因",
	}
	err = stream.SetRow("A2", row)
	if err != nil {
		return result, err
	}
	startRow := int64(3)
	for k, item := range data {
		isFail := false
		failReason := ""
		for _, value := range item {
			if value.FailReason != "" { // 如果有一条是失败的，则该ordersn下的所有单子都是失败
				isFail = true
				failReason = value.FailReason
			}
		}
		if isFail == true { // ordersn下的所有数都是失败
			for _, value := range item {
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				id := fmt.Sprintf("'%d", value.Id)
				row = []any{
					value.OrderSn, value.Sku, value.UpcCode, strconv.Itoa(value.IsRefusal), strconv.Itoa(value.RefusalQuantity), value.TrackingNumber,
					strconv.Itoa(value.Quantity),
					value.ReceiverName, value.ReceiverPhone, value.Province, value.City, value.County, value.ReceiverAddress,
					value.StoreCode, value.StoreName, value.OrderStatus, value.MemberName, value.CreateTime, value.UpdateTime, id, failReason,
				}
				err = stream.SetRow("A"+rowStartStr, row)
				atomic.AddInt64(&startRow, 1)
			}
			result.FailedNum += len(item)
		} else {
			orderId := k
			one, _ := s.orderDao.GetOne(c, enterpriseId, map[string]interface{}{
				"id": orderId,
			})
			sendQuantity := one.SendQuantity
			var trackingNumber string
			for _, value := range item {
				id := fmt.Sprintf("'%d", value.Id)
				skuInfo, _ := s.orderSkuDao.GetOne(c, enterpriseId, map[string]interface{}{
					"order_id": orderId,
					"sku_info": value.Sku,
				})
				if skuInfo.Quantity == value.RefusalQuantity {
					value.OrderStatus = "已拒单"
					value.UpdateTime = helper.GetFormatDateByParams(time.Now().Unix())
					_ = s.orderSkuDao.UpdateData(c, enterpriseId, map[string]interface{}{
						"id": skuInfo.ID,
					}, map[string]interface{}{
						"status": onlineWarehouseModel.OrderStatusRefusal,
					})
				}
				if skuInfo.Quantity > value.RefusalQuantity {
					insertQuantity := skuInfo.Quantity - value.RefusalQuantity
					if value.RefusalQuantity == 0 { // 没有退单,全部更新为已发货
						_ = s.orderSkuDao.UpdateData(c, enterpriseId, map[string]interface{}{
							"id": skuInfo.ID,
						}, map[string]interface{}{
							"quantity": insertQuantity,
							"status":   onlineWarehouseModel.OrderStatusShipped,
						})
						value.OrderStatus = "已发货"
						value.UpdateTime = helper.GetFormatDateByParams(time.Now().Unix())
					}
					if value.RefusalQuantity > 0 { // 退单数据小于商品订单数，则一部分更改为退单，一部分插入发货
						_ = s.orderSkuDao.UpdateData(c, enterpriseId, map[string]interface{}{
							"id": skuInfo.ID,
						}, map[string]interface{}{
							"quantity": value.RefusalQuantity,
							"status":   onlineWarehouseModel.OrderStatusRefusal,
						})
						if insertQuantity > 0 {
							insertSkuInfo := skuInfo
							insertSkuInfo.Quantity = insertQuantity
							insertSkuInfo.ID = app.UID()
							insertSkuInfo.CreatedAt = time.Now().Unix()
							insertSkuInfo.UpdatedAt = time.Now().Unix()
							insertSkuInfo.Status = onlineWarehouseModel.OrderStatusShipped
							value.OrderStatus = "已发货"
							value.UpdateTime = helper.GetFormatDateByParams(insertSkuInfo.UpdatedAt)
							_ = s.orderSkuDao.CreateData(c, enterpriseId, insertSkuInfo)
						}
					}
					sendQuantity += insertQuantity
					trackingNumber = value.TrackingNumber
				}
				rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
				row = []any{
					value.OrderSn, value.Sku, value.UpcCode, strconv.Itoa(value.IsRefusal), strconv.Itoa(value.RefusalQuantity), value.TrackingNumber,
					strconv.Itoa(value.Quantity),
					value.ReceiverName, value.ReceiverPhone, value.Province, value.City, value.County, value.ReceiverAddress,
					value.StoreCode, value.StoreName, value.OrderStatus, value.MemberName, value.CreateTime, value.UpdateTime, id, failReason,
				}
				err = stream.SetRow("A"+rowStartStr, row)
				atomic.AddInt64(&startRow, 1)
			}
			updateMap := map[string]interface{}{
				"status":       onlineWarehouseModel.OrderStatusRefusal,
				"admin_id":     req.AdminID,
				"operate_time": nowTime,
			}
			if sendQuantity > 0 { // 发货
				updateMap["tracking_number"] = trackingNumber
				updateMap["send_quantity"] = sendQuantity
				updateMap["status"] = onlineWarehouseModel.OrderStatusShipped
				updateMap["shipped_at"] = time.Now().Unix()
			}
			_ = s.orderDao.UpdateData(c, enterpriseId, map[string]interface{}{
				"id": one.ID,
			}, updateMap)
			result.Success += len(item)
			sendOrderSn = append(sendOrderSn, orderId)
		}
	}
	err = stream.MergeCell("A1", "T1")
	if err != nil {
		return result, err
	}
	err = stream.Flush()
	if err != nil {
		return result, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return result, err
	}
	if result.FailedNum > 0 {
		cs := storage.Init(cos.COS)
		remotePath := fmt.Sprintf("online_warehouse/%d/%d/%d/%s", time.Now().Year(), time.Now().Month(), time.Now().Day(), fileName)
		cosPath, err := cs.Upload(c, fileName, remotePath)
		if err != nil {
			return result, err
		}
		result.FailUrl = cs.GetURL(cosPath)
	}
	if len(sendOrderSn) > 0 {
		go s.sendUpdateNotice(c, enterpriseId, sendOrderSn, onlineWarehouseModel.OrderStatusShipped) // 发送消息通知
	}
	return result, err
}

func (s *OrderSvc) readXlsxData(c ctx.Context, enterpriseId int64, sheets []*xlsx.Sheet) (map[int64][]onlineWarehouseReq.ImportOrderData, onlineWarehouseResp.ImportOrderResp, error) {
	var resp onlineWarehouseResp.ImportOrderResp
	riMin := 2
	var settings []onlineWarehouseReq.ImportOrderData
	orderSnSkuMap := make(map[string]string)
	orderSnInfoMap := make(map[int64][]onlineWarehouseReq.ImportOrderData)
	for index, sheet := range sheets {
		if index > 0 {
			return nil, resp, errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri == 1 {
				if len(row.Cells) < 20 {

					return nil, resp, errors.New("上传的文件错误")
				}
				title := []string{
					"销售单号", "SKU", "UPC_code", "是否拒单", "拒单数量", "顺丰单号", "商品数量",
					"收件人", "收件电话", "省", "市", "区", "详细地址",
					"门店代码", "门店名称", "订单状态", "创建人", "创建时间", "更新时间", "后台编号(用于相同销售单号区分发货,导入时勿删)",
				}
				for k, cell := range row.Cells {
					if k < 20 && cell.String() != title[k] {
						return nil, resp, errors.New("上传的文件错误")
					}
				}
			}
			if ri < riMin {
				continue
			}
			if len(row.Cells) < 20 {
				continue
			}
			resp.UploadNum++
			var setting = onlineWarehouseReq.ImportOrderData{}
			isRefusal := strings.TrimSpace(row.Cells[3].String())
			setting.OrderSn = strings.TrimSpace(row.Cells[0].String())
			setting.Sku = strings.TrimSpace(row.Cells[1].String())
			setting.UpcCode = strings.TrimSpace(row.Cells[2].String())
			setting.IsRefusal, _ = strconv.Atoi(strings.TrimSpace(row.Cells[3].String()))
			setting.RefusalQuantity, _ = strconv.Atoi(strings.TrimSpace(row.Cells[4].String()))
			setting.TrackingNumber = strings.TrimSpace(row.Cells[5].String())
			setting.Quantity, _ = strconv.Atoi(strings.TrimSpace(row.Cells[6].String()))
			setting.ReceiverName = strings.TrimSpace(row.Cells[7].String())
			setting.ReceiverPhone = strings.TrimSpace(row.Cells[8].String())
			setting.Province = strings.TrimSpace(row.Cells[9].String())
			setting.City = strings.TrimSpace(row.Cells[10].String())
			setting.County = strings.TrimSpace(row.Cells[11].String())
			setting.ReceiverAddress = strings.TrimSpace(row.Cells[12].String())
			setting.StoreCode = strings.TrimSpace(row.Cells[13].String())
			setting.StoreName = strings.TrimSpace(row.Cells[14].String())
			setting.OrderStatus = strings.TrimSpace(row.Cells[15].String())
			setting.MemberName = strings.TrimSpace(row.Cells[16].String())
			setting.CreateTime = strings.TrimSpace(row.Cells[17].String())
			setting.UpdateTime = strings.TrimSpace(row.Cells[18].String())
			strId := strings.TrimSpace(row.Cells[19].String())
			if strId != "" {
				value := strings.TrimPrefix(strId, "'")
				setting.Id, _ = strconv.ParseInt(value, 10, 64)
			}
			checkData := true
			if setting.OrderSn == "" {
				setting.FailReason = "销售单号为空"
				checkData = false
			}
			if isRefusal == "" {
				setting.FailReason = "是否拒单为空"
				checkData = false
			}
			if setting.Id == 0 {
				setting.FailReason = "无效的后台编号"
				checkData = false
			}
			if checkData {
				one, err := s.orderDao.GetOne(c, enterpriseId, map[string]interface{}{
					"id": setting.Id,
				})
				uniqueKey := fmt.Sprintf("%d%s", setting.Id, setting.Sku)
				if err != nil || one.ID == 0 {
					setting.FailReason = "无效的销售单号"
				} else {
					if one.Status != onlineWarehouseModel.OrderStatusPending {
						setting.FailReason = "订单状态错误"
					}
					if one.Status == onlineWarehouseModel.OrderStatusShipped {
						setting.FailReason = "该订单已发货"
					}
					if one.Status == onlineWarehouseModel.OrderStatusRefusal {
						setting.FailReason = "该订单已拒单"
					}
					if _, ok := orderSnSkuMap[uniqueKey]; ok {
						resp.DuplicateNum++
						setting.FailReason = "数据重复"
					} else {
						orderSnSkuMap[uniqueKey] = uniqueKey
					}
				}
				skuInfo, err := s.orderSkuDao.GetOne(c, enterpriseId, map[string]interface{}{
					"order_id": setting.Id,
					"sku_info": setting.Sku,
				})

				if err != nil || skuInfo.ID == 0 {
					setting.FailReason = "销售单号或SKU错误"
				} else {
					if skuInfo.Quantity < setting.RefusalQuantity {
						setting.FailReason = "拒单数量超出该商品下单数量"
					}
					if skuInfo.Quantity > setting.RefusalQuantity {
						if setting.TrackingNumber == "" {
							setting.FailReason = "顺丰单号为空"
						}
					}
				}
			}
			settings = append(settings, setting)
			orderSnInfoMap[setting.Id] = append(orderSnInfoMap[setting.Id], setting)
		}
	}
	return orderSnInfoMap, resp, nil
}

func (s *OrderSvc) BatchTransfer(c ctx.Context, req onlineWarehouseReq.BatchTransferReq) error {
	orderList, err := s.orderDao.GetMany(c, req.EnterpriseId, "status=? and transfer_file=? and third_id = ? and shipped_at >= ?  and shipped_at < ? ", onlineWarehouseModel.OrderStatusShipped, "", req.ThirdId, req.ShippedTime, req.ShippedTime+86400)
	if err != nil {
		return err
	}
	var orderIds []int64
	for _, model := range orderList {
		orderIds = append(orderIds, model.ID)
	}
	transferFileUrl := helper.MarshalToString(req.TransferFile)
	err = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": orderIds,
	}, map[string]interface{}{
		"transfer_file": transferFileUrl,
	})
	return err
}

func (s *OrderSvc) Refund(c ctx.Context, req onlineWarehouseReq.InfoReq) error {
	one, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.OrderId,
	})
	if err != nil {
		return err
	}
	if one.ID == 0 {
		return errors.New("无效的销售单号")
	}
	if one.Status != onlineWarehouseModel.OrderStatusPending {
		return errors.New("订单状态不正确")
	}
	if err = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": one.ID,
	}, map[string]interface{}{
		"status":       onlineWarehouseModel.OrderStatusRefusal,
		"operate_time": time.Now().Unix(),
	}); err != nil {
		return err
	}
	if err = s.orderSkuDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"order_id": req.OrderId,
	}, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusRefusal,
	}); err != nil {
		return err
	}
	go s.sendUpdateNotice(c, req.EnterpriseId, []int64{req.OrderId}, onlineWarehouseModel.OrderStatusRefusal) // 发送消息通知
	return nil
}

func (s *OrderSvc) GetTipsNum(c ctx.Context, req onlineWarehouseReq.DashboardReq) (onlineWarehouseResp.TipsStatusNums, error) {
	var result onlineWarehouseResp.TipsStatusNums
	res, err := s.orderDao.GetTipsNum(c, req)
	if err != nil {
		return onlineWarehouseResp.TipsStatusNums{}, err
	}
	for _, re := range res {
		result.AllNum += re.Num
		switch re.Status {
		case onlineWarehouseModel.OrderStatusWait:
			result.WaitNum += re.Num
		case onlineWarehouseModel.OrderStatusPending:
			result.PendingNum = re.Num
		case onlineWarehouseModel.OrderStatusRefusal:
			result.RefuseNum = re.Num
		case onlineWarehouseModel.OrderStatusShipped:
			result.SendNum = re.Num
		case onlineWarehouseModel.OrderStatusCancel:
			result.CancelNum = re.Num

		default:

		}
	}
	return result, err
}

func (s *OrderSvc) StoreRank(c ctx.Context, req onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.StoreRankResp, error) {
	result := make([]onlineWarehouseResp.StoreRankResp, 0)
	list, err := s.orderDao.GetCountGroupThirdId(c, req)
	if err != nil {
		return result, err
	}
	var thirdIds []int
	for _, resp := range list {
		thirdIds = append(thirdIds, resp.ThirdId)
	}
	storeMap := make(map[int]microWorkModel.StoreDepartmentModel)
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": thirdIds,
	})
	for _, storeInfo := range storeList {
		storeMap[int(storeInfo.WechatDepartmentID)] = storeInfo
	}
	departmentMap := make(map[int]string, 0)
	departments, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"third_id":      thirdIds,
	})
	for _, dempt := range departments {
		departmentMap[dempt.ThirdID] = dempt.Name
	}
	for _, resp := range list {
		if item, ok := storeMap[resp.ThirdId]; ok {
			resp.StoreName = item.StoreName
			resp.StoreCode = item.ThirdStoreCode
		} else {
			resp.StoreName = departmentMap[resp.ThirdId]
			resp.StoreCode = departmentMap[resp.ThirdId]
		}
		result = append(result, resp)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SkuNum == result[j].SkuNum {
			return result[i].Num > result[j].Num
		} else {
			return result[i].SkuNum > result[j].SkuNum
		}
	})
	return result, err
}

func (s *OrderSvc) SaRank(c ctx.Context, req onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.SaRankResp, error) {
	result := make([]onlineWarehouseResp.SaRankResp, 0)
	list, err := s.orderDao.GetCountGroupMemberId(c, req)
	if err != nil {
		return result, err
	}
	storeMap := make(map[int]microWorkModel.StoreDepartmentModel)
	storeList, err := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
		"enterprise_id":        req.EnterpriseId,
		"wechat_department_id": req.ThirdIds,
	})
	for _, storeInfo := range storeList {
		storeMap[int(storeInfo.WechatDepartmentID)] = storeInfo
	}
	departmentMap := make(map[int]string, 0)
	departments, err := departmentsDao.NewDepartments().GetMany(c, map[string]interface{}{
		"enterprise_id": req.EnterpriseId,
		"third_id":      req.ThirdIds,
	})
	for _, dempt := range departments {
		departmentMap[dempt.ThirdID] = dempt.Name
	}
	memberIds := make([]int64, 0)
	for _, resp := range list {
		memberIds = append(memberIds, resp.MemberId)
		if item, ok := storeMap[resp.ThirdId]; ok {
			resp.StoreName = item.StoreName
			resp.StoreCode = item.ThirdStoreCode
		} else {
			resp.StoreName = departmentMap[resp.ThirdId]
			resp.StoreCode = departmentMap[resp.ThirdId]
		}
		result = append(result, resp)
	}
	memberList, err := membersDao.NewMembers().GetUnscopedMany(c, map[string]interface{}{
		"id": memberIds,
	})
	if err != nil {
		return result, err
	}
	memberMap := make(map[int64]membersModel.MemberModel, 0)
	for _, member := range memberList {
		memberMap[member.ID] = member
	}
	for k, item := range result {
		item.MemberName = memberMap[item.MemberId].Name
		item.Avatar = memberMap[item.MemberId].Avatar
		result[k] = item
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SkuNum == result[j].SkuNum {
			return result[i].Num > result[j].Num
		} else {
			return result[i].SkuNum > result[j].SkuNum
		}
	})
	return result, err
}

func (s *OrderSvc) GetUpdateNumByMember(c ctx.Context, enterpriseId int64, startTime int64, endTime int64) ([]onlineWarehouseResp.MemberNumResp, error) {
	result := make([]onlineWarehouseResp.MemberNumResp, 0)
	list, err := s.orderDao.GetUpdateNumByMember(c, enterpriseId, startTime, endTime)
	if err != nil {
		return result, err
	}
	memberIds := make([]int64, 0)
	for _, resp := range list {
		memberIds = append(memberIds, resp.MemberId)
	}
	many, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"id": memberIds,
	})
	if err != nil {
		return result, err
	}
	memberMap := make(map[int64]string, 0)
	for _, member := range many {
		memberMap[member.ID] = member.Userid
	}
	for _, item := range list {
		item.UserId = memberMap[item.MemberId]
		result = append(result, item)
	}
	return result, err
}

func (s *OrderSvc) sendUpdateNotice(c ctx.Context, enterpriseId int64, orderIdArr []int64, status int) error {
	list, err := s.orderDao.GetUpdateByOrderId(c, enterpriseId, orderIdArr)
	if err != nil {
		return err
	}
	memberIds := make([]int64, 0)
	for _, resp := range list {
		memberIds = append(memberIds, resp.MemberId)
	}
	many, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"id": memberIds,
	})
	if err != nil {
		return err
	}
	memberMap := make(map[int64]string, 0)
	for _, member := range many {
		memberMap[member.ID] = member.Userid
	}

	for _, model := range list {
		userId := memberMap[model.MemberId]
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "更新订单数",
			"value":   strconv.FormatInt(model.Num, 10),
		})
		content = append(content, map[string]string{
			"keyname": "更新时间",
			"value":   helper.GetFormatDateByParams(model.UpdateAt),
		})
		uri := fmt.Sprintf("online-order-list?status=%d", status)
		_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			Title:   "订单状态更新",
			Uri:     uri,
			Party:   nil,
			UserIds: []string{userId},
			Content: content,
			AppSign: qyWork.OnlineWareHouse,
		})
		if err != nil {
			fmt.Printf("导购%s:网仓订单消息通知失败", model.UserId)
		}
	}
	return err
}

func (s *OrderSvc) GetAllMemberId(c ctx.Context, enterpriseId int64) ([]int64, error) {
	var thirdIds = make([]int, 0)
	many, err := onlineWarehouseDao.NewStoreDao().GetMany(c, enterpriseId, map[string]interface{}{
		"deleted_at": 0,
	})
	if err != nil {
		return nil, err
	}
	for _, model := range many {
		thirdIds = append(thirdIds, model.ThirdId)
	}
	return s.orderDao.GetAllMemberId(c, enterpriseId, thirdIds)

}

func (s *OrderSvc) GetAllThirdId(c ctx.Context, enterpriseId int64) ([]int, error) {
	var thirdIds = make([]int, 0)
	many, err := onlineWarehouseDao.NewStoreDao().GetMany(c, enterpriseId, map[string]interface{}{
		"deleted_at": 0,
	})
	if err != nil {
		return nil, err
	}
	for _, model := range many {
		thirdIds = append(thirdIds, model.ThirdId)
	}
	return s.orderDao.GetAllThirdId(c, enterpriseId, thirdIds)

}

func (s *OrderSvc) Cancel(c ctx.Context, req onlineWarehouseReq.InfoReq) error {
	orderInfo, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id":        req.OrderId,
		"member_id": req.MemberId,
	})
	if err != nil {
		return err
	}
	if !helper.InArray(orderInfo.Status, []int32{onlineWarehouseModel.OrderStatusWait, onlineWarehouseModel.OrderStatusPending}) {
		return errors.New("只可以取消待填写，待发货的订单")
	}
	if orderInfo.Status == onlineWarehouseModel.OrderStatusPending && orderInfo.IsExport == onlineWarehouseModel.OrderExportYes {
		return errors.New("仓库正在处理当前订单，不能进行取消")
	}
	err = s.orderDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"id": orderInfo.ID,
	}, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusCancel,
	})
	err = s.orderSkuDao.UpdateData(c, req.EnterpriseId, map[string]interface{}{
		"order_id": orderInfo.ID,
	}, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusCancel,
	})
	return err
}

func (s *OrderSvc) TransferExport(c ctx.Context, req onlineWarehouseReq.TransferExportReq) (string, error) {
	fileName := "门店待上传转货单列表" + time.Now().Format("20060102150405") + ".xlsx"
	transferInfo, err := s.orderDao.GetAllTransferInfoByThird(c, req.EnterpriseId, req.ThirdId)
	if err != nil {
		return fileName, err
	}
	excelFile := excelize.NewFile()
	sheetName := "门店待上传转货单列表"
	// 创建sheet
	excelFile.SetSheetName("Sheet1", sheetName)
	err = excelFile.SaveAs(fileName)
	defer excelFile.Close()
	if err != nil {
		return fileName, err
	}
	stream, err := excelFile.NewStreamWriter(sheetName)
	if err != nil {
		return fileName, err
	}
	row := make([]any, 0)
	row = []any{
		"销售单号", "后台编号(用于相同销售单号区分发货,导入时勿删)",
	}
	err = stream.SetRow("A1", row)
	startRow := int64(2)
	for _, item := range transferInfo {
		rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
		id := fmt.Sprintf("'%d", item.ID)
		row = []any{item.OrderSn, id}
		err = stream.SetRow("A"+rowStartStr, row)
		atomic.AddInt64(&startRow, 1)
	}
	if err != nil {
		return fileName, err
	}
	err = stream.Flush()
	if err != nil {
		return fileName, err
	}
	err = excelFile.SaveAs(fileName)
	if err != nil {
		return fileName, err
	}
	return fileName, err
}

func (s *OrderSvc) readTransferXlsxData(c ctx.Context, enterpriseId int64, thirdId int, sheets []*xlsx.Sheet) ([]onlineWarehouseReq.TransferImportData, onlineWarehouseResp.TransferImportResp, error) {

	var resp onlineWarehouseResp.TransferImportResp
	var settings []onlineWarehouseReq.TransferImportData
	transferInfo, err := s.orderDao.GetAllTransferInfoByThird(c, enterpriseId, thirdId)
	if err != nil {
		return settings, resp, err
	}
	if len(transferInfo) == 0 {
		return settings, resp, errors.New("该门店没有待上传的转货单订单")
	}
	ids := make([]int64, 0)
	for _, item := range transferInfo {
		ids = append(ids, item.ID)
	}
	for index, sheet := range sheets {
		if index > 0 {
			return settings, resp, errors.New("文件格式错误")
		}
		for ri, row := range sheet.Rows {
			if ri == 0 {
				if len(row.Cells) < 2 {
					return settings, resp, errors.New("上传的文件错误")
				}
				title := []string{
					"销售单号", "后台编号(用于相同销售单号区分发货,导入时勿删)",
				}
				for k, cell := range row.Cells {
					if cell.String() != title[k] {
						return settings, resp, errors.New("上传的文件错误")
					}
				}
				continue
			}
			if len(row.Cells) < 2 {
				continue
			}
			orderSn := strings.TrimSpace(row.Cells[0].String())
			if orderSn == "" {
				continue
			}
			resp.UploadNum++
			failReason := ""
			var orderId int64
			strId := strings.TrimSpace(row.Cells[1].String())
			if strId != "" {
				value := strings.TrimPrefix(strId, "'")
				orderId, _ = strconv.ParseInt(value, 10, 64)
			}
			if !helper.InArray(orderId, ids) {
				resp.FailedNum++
				failReason = "后台编号不正确"
			} else {
				resp.Success++
			}
			settings = append(settings, onlineWarehouseReq.TransferImportData{
				OrderId:    orderId,
				OrderSn:    orderSn,
				FailReason: failReason,
			})
		}
	}
	return settings, resp, nil
}

func (s *OrderSvc) TransferImport(c ctx.Context, req onlineWarehouseReq.TransferImportReq) (onlineWarehouseResp.TransferImportResp, error) {
	var result onlineWarehouseResp.TransferImportResp
	var data []onlineWarehouseReq.TransferImportData
	fileInfo, err := readOnlineWarehouseImportFile(c, req.FileUrl)
	if err != nil {
		return result, err
	}

	zr, err := zip.NewReader(bytes.NewReader(fileInfo), int64(len(fileInfo)))
	if err != nil {
		return result, err
	}
	fileData, err := xlsx.ReadZipReader(zr)
	if err != nil {
		return result, err
	}
	enterpriseId := req.EnterpriseId
	var sendOrderId []int64

	data, result, err = s.readTransferXlsxData(c, enterpriseId, req.ThirdId, fileData.Sheets)
	if err != nil {
		return result, err
	}
	if result.FailedNum > 0 { // 如果有错误的销售单号，则全都不更新转货单，展示错误信息
		fileName := "门店待上传转货单错误列表" + time.Now().Format("20060102150405") + ".xlsx"
		excelFile := excelize.NewFile()
		sheetName := "门店待上传转货单错误列表"
		// 创建sheet
		excelFile.SetSheetName("Sheet1", sheetName)
		err = excelFile.SaveAs(fileName)
		defer excelFile.Close()
		if err != nil {
			return result, err
		}
		stream, err := excelFile.NewStreamWriter(sheetName)
		if err != nil {
			return result, err
		}
		row := make([]any, 0)
		row = []any{
			"销售单号", "后台编号(用于相同销售单号区分发货,导入时勿删)", "错误原因",
		}
		err = stream.SetRow("A1", row)
		startRow := int64(2)
		for _, item := range data {
			rowStartStr := strconv.Itoa(int(atomic.LoadInt64(&startRow)))
			id := fmt.Sprintf("'%d", item.OrderId)
			row = []any{item.OrderSn, id, item.FailReason}
			err = stream.SetRow("A"+rowStartStr, row)
			atomic.AddInt64(&startRow, 1)
		}
		if err != nil {
			return result, err
		}
		err = stream.Flush()
		if err != nil {
			return result, err
		}
		err = excelFile.SaveAs(fileName)
		if err != nil {
			return result, err
		}
		failUrl, err := helper.UploadFileToCos(c, fileName, false)
		if err != nil {
			return result, err
		}
		result.FailUrl = failUrl
		return result, err
	}
	for _, v := range data {
		sendOrderId = append(sendOrderId, v.OrderId)
	}
	if len(sendOrderId) == 0 {
		return result, errors.New("销售单号数据无效")
	}
	err = s.orderDao.UpdateData(c, enterpriseId, map[string]interface{}{
		"id":       sendOrderId,
		"third_id": req.ThirdId,
	}, map[string]interface{}{
		"transfer_file": helper.MarshalToString(req.TransferFile),
	})
	return result, err
}

func (s *OrderSvc) GetMediaId(c ctx.Context, req onlineWarehouseReq.OrderInfoReq) (string, error) {
	one, err := s.orderDao.GetOne(c, req.EnterpriseId, map[string]interface{}{
		"id": req.Id,
	})
	if err != nil {
		return "", err
	}
	posterUrl := one.SharePosterUrl
	if one.IsPoster != 1 || posterUrl == "" {
		posterUrl, err = s.createSharePoster(c, req.EnterpriseId, req.Id)
		if err != nil {
			return "", err
		}
	}
	mediaId, err := qyWork.NewQyWork(req.EnterpriseId).GetUploadMediaId(c, posterUrl, "image")
	return mediaId, err
}

func (s *OrderSvc) CheckOrderSn(c ctx.Context, req onlineWarehouseReq.CheckOrderSnReq) (onlineWarehouseModel.OnlineWarehouseOrderModel, error) {
	orderSn := strings.TrimSpace(req.OrderSn)
	res, err := s.orderDao.GetLastOrderInfo(c, req.EnterpriseId, orderSn)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return res, err
}

func (s *OrderSvc) Dashboard(c ctx.Context, req onlineWarehouseReq.DashboardReq) (onlineWarehouseResp.DashboardResp, error) {
	var resp onlineWarehouseResp.DashboardResp
	res, err := s.orderDao.GetTipsNum(c, req)
	if err != nil {
		return resp, err
	}
	result := resp.Order
	for _, re := range res {
		result.AllNum += re.Num
		switch re.Status {
		case onlineWarehouseModel.OrderStatusWait:
			result.WaitNum += re.Num
		case onlineWarehouseModel.OrderStatusPending:
			result.PendingNum = re.Num
		case onlineWarehouseModel.OrderStatusRefusal:
			result.RefuseNum = re.Num
		case onlineWarehouseModel.OrderStatusShipped:
			result.SendNum = re.Num
		case onlineWarehouseModel.OrderStatusCancel:
			result.CancelNum = re.Num
		default:
		}
	}
	resp.Order = result
	ids, err := s.orderDao.GetOrderIds(c, req)
	if err != nil {
		return resp, err
	}
	goodsTipsNum, err := s.orderSkuDao.GetDashboardNum(c, req.EnterpriseId, ids)
	if err != nil {
		return resp, err
	}
	goodsResult := resp.Goods
	for _, re := range goodsTipsNum {
		goodsResult.AllNum += re.Num
		switch re.Status {
		case onlineWarehouseModel.OrderStatusWait:
			goodsResult.WaitNum += re.Num
		case onlineWarehouseModel.OrderStatusPending:
			goodsResult.PendingNum = re.Num
		case onlineWarehouseModel.OrderStatusRefusal:
			goodsResult.RefuseNum = re.Num
		case onlineWarehouseModel.OrderStatusShipped:
			goodsResult.SendNum = re.Num
		case onlineWarehouseModel.OrderStatusCancel:
			goodsResult.CancelNum = re.Num
		default:
		}
	}
	resp.Goods = goodsResult
	return resp, err
}

func (s *OrderSvc) GoodsRank(c ctx.Context, req onlineWarehouseReq.DashboardReq) ([]onlineWarehouseResp.GoodsRankResp, error) {
	return s.orderSkuDao.GoodsRank(c, req)
}
