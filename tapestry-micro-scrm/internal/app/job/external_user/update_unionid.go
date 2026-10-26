package externalUserJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/xjob"
)

type UpdateUnionid struct {
	name         string
	enterpriseId int64
	domain       string
	batchSize    int
}

func NewUpdateUnionid() *UpdateUnionid {
	return &UpdateUnionid{
		name:         "更新外部联系人unionid",
		enterpriseId: 3374939108972352,
		domain:       "AFA7D6C3AC1101464206C1A15C31C214",
		batchSize:    20000,
	}
}
func (job *UpdateUnionid) Run(c ctx.Context, handler *xjob.JobHandler) error {
	//workImpl := new(work.QyWorkImpl)
	//externalDao := microScrmDao.NewWxUserCorpExternal()
	//unionidMapDao := microScrmDao.NewUnionidMap()
	//
	//workImpl.SetDomain(job.domain)
	//
	//// 查询所有员工
	//memberList, err := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND deleted_at = 0", job.enterpriseId)
	//if err != nil || len(memberList) == 0 {
	//	return err
	//}
	//
	//// 进度统计
	//var (
	//	processedMemberTotal  int           // 已处理的导购数量
	//	processedUnionidTotal int           // 已处理的unionid数量
	//	apiCallCount          int           // API调用次数
	//	totalAPICallDuration  time.Duration // 累计API调用耗时
	//	startTime             = time.Now()  // 开始时间
	//)
	//
	//// 使用WaitGroup和channel来处理并发数据处理
	//var wg sync.WaitGroup
	//resultsChan := make(chan int, len(memberList)*10) // 缓冲通道接收处理结果
	//
	//// 启动结果收集goroutine
	//go func() {
	//	for count := range resultsChan {
	//		processedUnionidTotal += count
	//		wg.Done()
	//	}
	//}()
	//
	//// 处理每个员工的外部联系人（顺序调用企微API）
	//for _, member := range memberList {
	//	memberUnionidCount := 0
	//	cursor := ""
	//
	//	for {
	//		// 调用企微API（顺序调用）
	//		apiCallStart := time.Now()
	//		externalUsers, err := workImpl.BatchExternalContactGet(c, work.BatchExternalContactRequest{
	//			MemberUserID: member.Userid,
	//			Cursor:       cursor,
	//			Limit:        100,
	//		})
	//		apiCallDuration := time.Since(apiCallStart)
	//		apiCallCount++
	//		totalAPICallDuration += apiCallDuration
	//
	//		if err != nil {
	//			// 记录错误但不中断
	//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("导购 %s API调用失败: %s", member.Userid, err.Error()))
	//			break
	//		}
	//
	//		if len(externalUsers.ExternalContactList) == 0 {
	//			break
	//		}
	//
	//		// 并发处理数据
	//		wg.Add(1)
	//		go func(contacts []work.WorkExternalContact, memberUserid string) {
	//			defer func() {
	//				if r := recover(); r != nil {
	//					go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("导购 %s 数据处理发生panic: %v\n", memberUserid, r))
	//					wg.Done()
	//				}
	//			}()
	//
	//			processedCount, err := job.processExternalContacts(c, contacts, externalDao, unionidMapDao, job.enterpriseId)
	//			if err != nil {
	//				// 记录错误但不中断
	//				go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("导购 %s 数据处理失败: %s", memberUserid, err.Error()))
	//			}
	//
	//			// 发送处理结果到通道
	//			resultsChan <- processedCount
	//		}(externalUsers.ExternalContactList, member.Userid)
	//
	//		if externalUsers.NextCursor == "" {
	//			break
	//		}
	//		cursor = externalUsers.NextCursor
	//	}
	//
	//	// 更新导购统计
	//	processedMemberTotal++
	//
	//	// 每5个导购打一个日志
	//	if processedMemberTotal%20 == 0 {
	//		elapsed := time.Since(startTime)
	//		avgRate := float64(apiCallCount) / elapsed.Seconds()
	//		avgCallTime := totalAPICallDuration / time.Duration(apiCallCount)
	//		go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("已处理 %d 个导购, API调用 %d 次, 平均速率: %.2f 次/秒, 平均调用耗时: %v\n",
	//			processedMemberTotal, apiCallCount, avgRate, avgCallTime))
	//
	//	}
	//
	//	// 每20个unionid打一个日志
	//	if processedUnionidTotal/1000 > (processedUnionidTotal-memberUnionidCount)/1000 {
	//		elapsed := time.Since(startTime)
	//		avgRate := float64(apiCallCount) / elapsed.Seconds()
	//		avgCallTime := totalAPICallDuration / time.Duration(apiCallCount)
	//		go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("已处理 %d 个unionid, API调用 %d 次, 平均速率: %.2f 次/秒, 平均调用耗时: %v\n",
	//			processedUnionidTotal, apiCallCount, avgRate, avgCallTime))
	//	}
	//}
	//
	//// 等待所有数据处理完成
	//wg.Wait()
	//close(resultsChan)
	//
	//// 最终统计
	//totalElapsed := time.Since(startTime)
	//finalRate := float64(apiCallCount) / totalElapsed.Seconds()
	//avgCallTime := totalAPICallDuration / time.Duration(apiCallCount)
	//fmt.Printf("任务完成: 总共处理 %d 个导购, %d 个unionid\n", processedMemberTotal, processedUnionidTotal)
	//fmt.Printf("API统计: 调用 %d 次, 总耗时: %v, 平均速率: %.2f 次/秒, 平均调用耗时: %v\n",
	//	apiCallCount, totalElapsed.Round(time.Second), finalRate, avgCallTime)
	//job.pushCrm(c)

	return nil
}

//
//// 修改processExternalContacts方法，返回处理的unionid数量
//func (job *UpdateUnionid) processExternalContacts(c ctx.Context,
//	externalContacts []work.WorkExternalContact,
//	externalDao *microScrmDao.WxUserCorpExternal,
//	unionidMapDao *microScrmDao.UnionidMapDao,
//	enterpriseId int64) (int, error) {
//
//	var openIds []string
//	for _, external := range externalContacts {
//		if external.Base.UnionID != "" {
//			openIds = append(openIds, external.Base.ExternalUserID)
//		}
//	}
//
//	if len(openIds) == 0 {
//		return 0, nil
//	}
//
//	// 查询历史数据
//	oldList, err := externalDao.GetMany(c, enterpriseId, "openid in ?", openIds)
//	if err != nil {
//		return 0, fmt.Errorf("查询历史数据失败: %w", err)
//	}
//
//	// 构建openid到unionid的映射
//	openIdMap := make(map[string]string)
//	for _, model := range oldList {
//		openIdMap[model.Openid] = model.Unionid
//	}
//
//	// 准备批量插入数据
//	nowUnix := time.Now().Unix()
//	var addSlice []microScrmModel.UnionidMapModel
//
//	for _, external := range externalContacts {
//		if external.Base.UnionID == "" {
//			continue
//		}
//
//		oldUnionid, exists := openIdMap[external.Base.ExternalUserID]
//		if !exists {
//			continue
//		}
//
//		addSlice = append(addSlice, microScrmModel.UnionidMapModel{
//			BaseModel: models.BaseModel{
//				CreatedAt: nowUnix,
//				UpdatedAt: nowUnix,
//			},
//			Nick:       external.Base.Name,
//			Openid:     external.Base.ExternalUserID,
//			OldUnionid: oldUnionid,
//			NewUnionid: external.Base.UnionID,
//		})
//	}
//
//	// 批量插入
//	if len(addSlice) > 0 {
//		err := unionidMapDao.CreateInBatches(c, enterpriseId, addSlice)
//		if err != nil {
//			return 0, fmt.Errorf("批量插入失败: %w", err)
//		}
//		return len(addSlice), nil
//	}
//
//	return 0, nil
//}
//
//func (job *UpdateUnionid) batchUpdate(c ctx.Context) error {
//
//	taskList, _ := microScrmDao.NewPlanTaskDao().GetMany(c, job.enterpriseId, "task_end_time > ?", time.Now().Unix())
//	var planIds []int64
//	for _, task := range taskList {
//		planIds = append(planIds, task.PlanID)
//	}
//	planIds = helper.ArrayUniq(planIds)
//	lastID := int64(0)
//	totalProcessed := 0
//	batchNumber := 0
//	unionidMapDao := microScrmDao.NewUnionidMap()
//	for {
//		batchNumber++
//		go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("开始处理第 %d 批, 从ID %d 开始, 限制 %d 条",
//			batchNumber, lastID, job.batchSize))
//
//		startTime := time.Now()
//
//		// 1. 获取当前批次的unionid映射关系
//		unionidList, err := unionidMapDao.GetBatch(c, job.enterpriseId, lastID, job.batchSize)
//		if err != nil {
//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("获取第 %d 批unionid映射失败: %v", batchNumber, err))
//			return err
//		}
//
//		if len(unionidList) == 0 {
//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("所有unionid处理完成，总计: %d 条", totalProcessed))
//			break
//		}
//
//		// 2. 提取当前批次的old_unionid列表
//		var oldUnionIds []string
//		for _, mapping := range unionidList {
//			if mapping.OldUnionid != "" {
//				oldUnionIds = append(oldUnionIds, mapping.OldUnionid)
//			}
//		}
//
//		if len(oldUnionIds) == 0 {
//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("第 %d 批没有需要更新的unionid，跳过", batchNumber))
//			lastID = unionidList[len(unionidList)-1].ID
//			continue
//		}
//
//		// 3. 批量更新各个表
//		affected, err := job.batchUpdateTables(c, job.enterpriseId, unionidList, planIds)
//		if err != nil {
//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("第 %d 批更新失败: %v", batchNumber, err))
//			return err
//		}
//
//		// 4. 更新lastID为当前批次的最后一个ID
//		lastID = unionidList[len(unionidList)-1].ID
//		totalProcessed += len(unionidList)
//
//		duration := time.Since(startTime)
//		go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("第 %d 批完成, 处理 %d 个unionid, 影响 %d 行, 累计 %d, 耗时: %v",
//			batchNumber, len(unionidList), affected, totalProcessed, duration))
//
//		// 5. 批次间休息
//		if len(unionidList) >= job.batchSize {
//			sleepTime := calculateSleepTime(duration, affected)
//			go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("等待 %v 后继续下一批...", sleepTime))
//			time.Sleep(sleepTime)
//		}
//	}
//
//	go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("任务完成! 总共处理 %d 个unionid", totalProcessed))
//	return nil
//}
//
//// 批量更新所有相关表
//func (job *UpdateUnionid) batchUpdateOther(c ctx.Context) error {
//	enterpriseUserDao := microScrmDao.NewEnterpriseUser()
//	unionidMapDao := microScrmDao.NewUnionidMap()
//	wxUserUnionDao := microScrmDao.NewWxUserUnion()
//	lastID := int64(3562705685500160)
//	for {
//		enterpriseUserList, err := enterpriseUserDao.GetBatch(c, job.enterpriseId, lastID, 1000)
//		if err != nil {
//			return err
//		}
//		if len(enterpriseUserList) < 2 {
//			return nil
//		}
//		var unionIds []string
//		for _, user := range enterpriseUserList {
//			var tmpMap map[string]string
//			helper.Unmarshal(user.WxUnion, &tmpMap)
//			unionIds = append(unionIds, tmpMap["default"])
//		}
//
//		// 1. 获取当前批次的unionid映射关系
//		unionidList, _ := unionidMapDao.GetMany(c, job.enterpriseId, "old_unionid in ?", unionIds)
//		var unionidMap = make(map[string]string)
//		for _, model := range unionidList {
//			unionidMap[model.OldUnionid] = model.NewUnionid
//		}
//		var delIds []int64
//		var addUser []microScrmModel.EnterpriseUserModel
//		for _, user := range enterpriseUserList {
//			var tmpMap map[string]string
//			helper.Unmarshal(user.WxUnion, &tmpMap)
//			if newUnionid, ok := unionidMap[tmpMap["default"]]; ok {
//				user.WxUnion = `{"default":"` + newUnionid + `"}`
//				// 修复：使用IsZero()判断零值时间
//				if user.EsUpdate.IsZero() {
//					user.EsUpdate = time.Now()
//				}
//				if user.DtUpdate.IsZero() {
//					user.DtUpdate = time.Now()
//				}
//				addUser = append(addUser, user)
//				delIds = append(delIds, user.ID)
//			}
//		}
//		if len(delIds) > 0 {
//			err = enterpriseUserDao.DeleteData(c, job.enterpriseId, map[string]interface{}{
//				"id": delIds,
//			})
//			if err != nil {
//				return err
//			}
//			err = enterpriseUserDao.CreateInBatchesData(c, job.enterpriseId, addUser)
//			if err != nil {
//				return err
//			}
//		}
//
//		//处理wxUnion表
//		wxUnionList, err := wxUserUnionDao.GetMany(c, job.enterpriseId, "unionid in ?", unionIds)
//		if err != nil {
//			return err
//		}
//		var delWxIds []string
//		var addWx []micro_scrm_enterprise.WxUserUnionModel
//		for _, wxUnion := range wxUnionList {
//			if newUnionid, ok := unionidMap[wxUnion.Unionid]; ok {
//				wxUnion.Unionid = newUnionid
//				addWx = append(addWx, wxUnion)
//				delWxIds = append(delWxIds, wxUnion.Unionid)
//			}
//		}
//		if len(delWxIds) > 0 {
//			err = wxUserUnionDao.DeleteData(c, job.enterpriseId, map[string]interface{}{
//				"unionid": delWxIds,
//			})
//			if err != nil {
//				return err
//			}
//			err = wxUserUnionDao.CreateInBatchesData(c, job.enterpriseId, addWx)
//			if err != nil {
//				return err
//			}
//		}
//
//		// 4. 更新lastID为当前批次的最后一个ID
//		lastID = enterpriseUserList[len(enterpriseUserList)-1].ID
//	}
//
//}
//
//// 批量更新所有相关表
//func (job *UpdateUnionid) batchUpdateTables(c ctx.Context, enterpriseId int64, mappings []microScrmModel.UnionidMapModel, planIds []int64) (int64, error) {
//	nowUnix := time.Now().Unix()
//	totalAffected := int64(0)
//	externalDao := microScrmDao.NewWxUserCorpExternal()
//	relationDao := microScrmDao.NewUserRelation()
//	//enterpriseUserDao := microScrmDao.NewEnterpriseUser()
//	orderDao := microScrmDao.NewReservesOrder()
//	taskDataDao := microScrmDao.NewPlanTaskConditionDataDao()
//	//wxUserUnionDao := microScrmDao.NewWxUserUnion()
//
//	// 1. 更新 wx_user_corp_external 表
//	affected, err := externalDao.BatchUpdateByUnionIds(c, enterpriseId, mappings, nowUnix)
//	if err != nil {
//		return 0, fmt.Errorf("更新wx_user_corp_external失败: %v", err)
//	}
//	totalAffected += affected
//	//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - wx_user_corp_external 表更新 %d 行", affected))
//
//	// 2. 更新 user_relation 表
//	affected, err = relationDao.BatchUpdateByUnionIds(c, enterpriseId, mappings, nowUnix)
//	if err != nil {
//		return 0, fmt.Errorf("更新user_relation失败: %v", err)
//	}
//	totalAffected += affected
//	//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - user_relation 表更新 %d 行", affected))
//
//	// 3. 更新 enterprise_user 表
//	//affected, err = enterpriseUserDao.BatchUpdateByUnionIds(c, enterpriseId, mappings, nowUnix)
//	//if err != nil {
//	//	return 0, fmt.Errorf("更新enterprise_user失败: %v", err)
//	//}
//	//totalAffected += affected
//	//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - enterprise_user 表更新 %d 行", affected))
//
//	// 4. 更新 reserves_order 表
//	affected, err = orderDao.BatchUpdateByUnionIds(c, enterpriseId, mappings, nowUnix)
//	if err != nil {
//		return 0, fmt.Errorf("更新reserves_order失败: %v", err)
//	}
//	totalAffected += affected
//	//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - reserves_order 表更新 %d 行", affected))
//
//	// 5. 更新 plan_task_condition_data 表
//	for _, planId := range planIds {
//		affected, err = taskDataDao.BatchUpdateByUnionIds(c, job.enterpriseId, planId, mappings)
//		if err != nil {
//			return 0, fmt.Errorf("更新plan_task_condition_data失败: %v", err)
//		}
//		totalAffected += affected
//		//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - plan_task_condition_data(plan_id=%d) 表更新 %d 行", planId, affected))
//	}
//	//// 6. 更新 wx_user_union 表
//	//affected, err = wxUserUnionDao.BatchUpdateByUnionIds(c, enterpriseId, mappings, nowUnix)
//	//if err != nil {
//	//	return 0, fmt.Errorf("更新wx_user_union失败: %v", err)
//	//}
//	//totalAffected += affected
//	//go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("  - reserves_order 表更新 %d 行", affected))
//	return totalAffected, nil
//}
//
//// 计算休息时间
//func calculateSleepTime(duration time.Duration, affected int64) time.Duration {
//	baseSleep := 2 * time.Second
//
//	// 如果处理很快，减少休息时间
//	if duration < 3*time.Second {
//		return baseSleep
//	}
//
//	// 如果处理较慢或影响行数多，增加休息时间
//	if duration > 10*time.Second || affected > 50000 {
//		return 5 * time.Second
//	}
//
//	return baseSleep
//}
//
//func (job *UpdateUnionid) pushCrm(c ctx.Context) {
//	newCdp := microwork.NewCdp()
//	// 获取所有导购所在的门店信息
//	memberDepartmentMap := make(map[string]string)
//	storeDepartmentInfo, _ := microWorkDao.NewStoreDepartment().GetMany(c, map[string]interface{}{
//		"enterprise_id": job.enterpriseId,
//		"deleted_at":    0,
//	})
//	thirdIdNameMap := make(map[int]string)
//	for _, info := range storeDepartmentInfo {
//		if _, ok := thirdIdNameMap[int(info.WechatDepartmentID)]; !ok {
//			thirdIdNameMap[int(info.WechatDepartmentID)] = info.ThirdStoreCode
//		}
//	}
//	memberLists, _ := membersDao.NewMembers().GetAllMemberInfo(c, job.enterpriseId)
//	for _, info := range memberLists {
//		if len(info.MemberDepartments) >= 1 {
//			for _, department := range info.MemberDepartments {
//				if department.DeletedAt == 0 { // 导购门店信息初始化
//					departmentName := department.Department.Name
//					memberDepartmentMap[info.Userid] = departmentName
//					thirdId := department.Department.ThirdID
//					if name, ok := thirdIdNameMap[thirdId]; ok {
//						memberDepartmentMap[info.Userid] = name
//					}
//				}
//			}
//		}
//	}
//	stateList, _ := microScrmDao.NewDrainageContactsDao().GetMany(c, map[string]interface{}{
//		"enterprise_id": job.enterpriseId,
//	})
//	var stateNameMap = make(map[int64]string)
//	for _, model := range stateList {
//		stateNameMap[model.ID] = model.Name
//	}
//	lastId := int64(3487331194673536)
//	var totalProcessed int64
//	for {
//		list, err := microScrmDao.NewUserRelation().GetManyLimit(c, job.enterpriseId, lastId, 1000)
//		if err != nil || len(list) == 0 {
//			break
//		}
//		lastId = list[len(list)-1].ID
//
//		// 创建等待组，用于等待所有协程完成
//		var wg sync.WaitGroup
//		// 创建通道用于控制并发数量
//		sem := make(chan struct{}, 5)
//
//		for _, model := range list {
//			if storeCode, ok := memberDepartmentMap[model.UserID]; ok {
//				wg.Add(1)
//				sem <- struct{}{} // 获取信号量
//
//				go func(m microwork.SaAssignmentReq, storeCode string) {
//					defer wg.Done()
//					defer func() { <-sem }() // 释放信号量
//
//					_, _ = newCdp.SaAssignment(c, m)
//					// 原子操作增加计数器
//					atomic.AddInt64(&totalProcessed, 1)
//					processed := totalProcessed
//					// 每处理10000个打印日志
//					if processed%10000 == 0 {
//						go notice.NewWechat().SendTextMessage(c, fmt.Sprintf("SW-CRM数据推送进度:%d,lastID:%d", processed, lastId))
//					}
//				}(microwork.SaAssignmentReq{
//					EnterpriseID: job.enterpriseId,
//					Params: microwork.SaAssignmentParamsReq{
//						SaAccount:      model.UserID,
//						FirstBindDate:  helper.GetFormatDateByParams(model.CreatedAt),
//						LatestBindDate: helper.GetFormatDateByParams(int64(model.Createtime)),
//						UnionId:        model.Unionid,
//						StoreCode:      storeCode,
//						Source:         getStateName(model.State, stateNameMap),
//						BindTime:       helper.GetFormatDateByParams(model.CreatedAt),
//					},
//				}, storeCode)
//			}
//		}
//
//		// 等待当前批次的所有协程完成
//		wg.Wait()
//	}
//
//}
//
//func getStateName(state string, stateNameMap map[int64]string) string {
//	if state == "" {
//		return ""
//	}
//	stateId, _ := strconv.ParseInt(state, 10, 64)
//	return stateNameMap[stateId]
//}
