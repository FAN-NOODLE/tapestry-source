package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/orm"
)

//MemberRepo 定义成员存储
type MarketingRepo struct {
	*repos.SQLRepo
}

var txshare *sqlplus.Tx

//MarketingAdd 营销策略添加
func (m *MarketingRepo) MarketingAdd(c ctx.Context, request *marketing.MarketingStrategies) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-add")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/marketing-strategy-add")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, &request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//MarketingStrategyConditionAdd 营销策略条件添加
func (m *MarketingRepo) MarketingStrategyConditionAdd(c ctx.Context, request *marketing.MarketingStrategyCondition) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-add-condition")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/marketing-strategy-add-condition")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, &request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//MarketingStrategyChannelAdd 营销策略渠道添加
func (m *MarketingRepo) MarketingStrategyChannelAdd(c ctx.Context, request []marketing.ChannelInfo) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-add-channels")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/marketing-strategy-add-channels")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//MarketingStrategyTriggerRelationAdd 营销策略触发关系
func (m *MarketingRepo) MarketingStrategyTriggerRelationAdd(c ctx.Context, request []marketing.RelationInfo) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-add-relations")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/marketing-strategy-add-relations")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//MarketingDetail 策略详情
func (m *MarketingRepo) MarketingDetail(c ctx.Context, request int64) (*marketing.MarketingStrategy, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-detail")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	response := &marketing.MarketingStrategy{}
	sTmtErr := stmt.Get(c, response, request)
	if sTmtErr != nil {
		return nil, sTmtErr
	}
	return response, err
}

func (m *MarketingRepo) MarketingDetailChannels(c ctx.Context, request marketing.MarketingDetailRequest) (*[]marketing.ChannelInfo, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-detail-channels")
	if err != nil {
		return nil, err
	}
	response := &[]marketing.ChannelInfo{}
	res, sTmtErr := stmt.NamedQuery(c, request)
	scErr := sqlplus.ScanAll(res, response)
	if scErr != nil {
		return nil, sTmtErr
	}
	if sTmtErr != nil {
		return nil, sTmtErr
	}
	return response, err
}
func (m *MarketingRepo) MarketingDetailRelations(c ctx.Context, request marketing.MarketingDetailRequest) (*[]marketing.RelationInfo, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-detail-relations")
	if err != nil {
		return nil, err
	}
	response := &[]marketing.RelationInfo{}
	res, sTmtErr := stmt.NamedQuery(c, request)
	scErr := sqlplus.ScanAll(res, response)
	if scErr != nil {
		return nil, sTmtErr
	}
	if sTmtErr != nil {
		return nil, sTmtErr
	}
	return response, err
}

//MarketingUpdate 策略更新
func (m *MarketingRepo) MarketingUpdate(c ctx.Context, request marketing.MarketingStrategy) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-update")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

func (m *MarketingRepo) MarketingUpdateTriggerCount2(c ctx.Context, request marketing.MarketingStrategy) error {

	stmt, err := m.DB.Template("t_marketings/marketing-strategy-trigger-count-update2")
	if err != nil {
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

//MarketingDelete 策略删除
func (m *MarketingRepo) MarketingDelete(c ctx.Context, request marketing.DeleteRequest) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-delete")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//MarketingStrategyRecordAdd 营销策略触发记录添加
func (m *MarketingRepo) MarketingStrategyRecordAdd(c ctx.Context, request marketing.MarketingStrategy) (int64, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-record-add")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	lastID, lastErr := result.LastInsertId()
	if lastErr != nil {
		return 0, lastErr
	}

	return lastID, err
}

//MarketingStrategyRecordUpdate 营销策略触发记录修改
func (m *MarketingRepo) MarketingStrategyRecordUpdate(c ctx.Context, request marketing.MarketingStrategy) (int64, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-record-add")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	lastID, lastErr := result.LastInsertId()
	if lastErr != nil {
		return 0, lastErr
	}

	return lastID, err
}

//DeleteMarketingStrategy 删除策略
func (m *MarketingRepo) DeleteMarketingStrategy(c ctx.Context, request marketing.DeleteRequest) error {
	stmt, err := m.DB.Template("t_marketings/delete-marketing-strategy")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/delete-marketing-strategy")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//DeleteMarketingStrategyCondition 删除策略条件
func (m *MarketingRepo) DeleteMarketingStrategyCondition(c ctx.Context, request marketing.DeleteRequest) error {
	stmt, err := m.DB.Template("t_marketings/delete-marketing-strategy-condition")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/delete-marketing-strategy-condition")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//DeleteMarketingStrategyChannels 删除策略渠道
func (m *MarketingRepo) DeleteMarketingStrategyChannels(c ctx.Context, request marketing.DeleteRequest) error {
	stmt, err := m.DB.Template("t_marketings/delete-marketing-strategy-channels")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/delete-marketing-strategy-channels")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

//DeleteMarketingStrategyRelations 删除策略关联
func (m *MarketingRepo) DeleteMarketingStrategyRelations(c ctx.Context, request marketing.DeleteRequest) error {
	stmt, err := m.DB.Template("t_marketings/delete-marketing-strategy-relations")
	if txshare != nil {
		stmt, err = txshare.Template("t_marketings/delete-marketing-strategy-relations")
	}
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

func (m *MarketingRepo) GetSmsNumber(c ctx.Context, args marketing.MarketingDataDetailRequest) ([]marketing.SmsSendLog, error) {
	stmt, err := m.DB.Template("t_marketings/select-sms-send-log", m.Cond2Args(args))
	var list []marketing.SmsSendLog
	if err != nil {
		helper.Must(err)
	}
	res, err := stmt.NamedQuery(c, args)
	helper.Must(err)
	err = sqlx.StructScan(res, &list)
	helper.Must(err)
	return list, err
}

func (m *MarketingRepo) AddMarketStrategyRecord(c ctx.Context, request *marketing.MarketStrategyRecord) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-add")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// BatchAddMarketStrategyRecord 批量生成策略记录
func (m *MarketingRepo) BatchAddMarketStrategyRecord(c ctx.Context, request []*marketing.MarketStrategyRecord) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-add")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

func (m *MarketingRepo) GeSmsCount(c ctx.Context, smsId int64) (int64, error) {
	mapping := make(map[string]interface{})
	cond := &struct {
		SmsId int64
	}{
		SmsId: smsId,
	}
	mapping["cond"] = cond
	stmt, err := m.DB.Template("t_marketings/select-marketing-strategy-channels-count", &mapping)

	if err != nil {
		return 0, err
	}
	type selectFromDb struct {
		Count int64 `json:"count" db:"count"`
	}
	dest := &selectFromDb{}
	sErr := stmt.Get(c, dest)
	if sErr != nil {
		return 0, sErr
	}
	return dest.Count, nil
}

func (m *MarketingRepo) UpdateMarketStrategyRecord(c ctx.Context, request *marketing.MarketStrategyRecord) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-update")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	log.Info(c).Msgf("更新策略记录参数：%+v", request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

func (m *MarketingRepo) CloseMarketingInnerStatus(c ctx.Context, request *marketing.MarketingDataDetailRequest) error {
	stmt, err := m.DB.Template("t_marketings/close-maketing-innert-status")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// MarketingCron 查询可执行的策略
func (m *MarketingRepo) MarketingCron(c ctx.Context, request *marketing.RecordCronRequest) ([]*marketing.MarketTrigExecRequest, error) {
	mapping := m.Cond2Args(request)
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-find-list", mapping)
	if err != nil {
		return nil, err
	}
	var response []*marketing.MarketTrigExecRequest
	res, sTmtErr := stmt.NamedQuery(c, request)
	scErr := sqlplus.ScanAll(res, &response)
	if scErr != nil {
		return nil, scErr
	}
	if sTmtErr != nil {
		return nil, sTmtErr
	}
	return response, err
}

// MarketingUpdateStatus
func (m *MarketingRepo) MarketingUpdateStatus(c ctx.Context, request *marketing.UpdateStatusRequest) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-update-status", m.Cond2Args(request))
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}
func (m *MarketingRepo) UpdateStrategyTime(c ctx.Context, args *marketing.UpdateStatusRequest) error {
	stmt, err := m.DB.Template("t_marketings/update-strategy-setting-day", m.Cond2Args(args))
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, args)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// SelectMarketingList 获取策略列表接口
func (m *MarketingRepo) SelectMarketingList(c ctx.Context) (*[]marketing.MarketingStrategies, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-list")
	if err != nil {
		return nil, err
	}
	var marketingListRequest marketing.MarketingStrategies
	ListRequest := &[]marketing.MarketingStrategies{}
	res, errStmt := stmt.NamedQuery(c, marketingListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

func (m *MarketingRepo) MarketTrig(cc ctx.Context, data *marketing.MarketTrigRequest) (*marketing.MarketTrigResponse, error) {

	option := &repos.ListOption{
		TplPrefix: "t_marketings/marketing-trigger-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &marketing.MarketTrigResponse{}
		},
	}
	params := &marketing.MarketTrigRequest{
		EnterpriseID:   data.EnterpriseID,
		Condition:      data.Condition,
		StrategyStatus: marketing.StrategyStatusStart,
		GroupIds:       data.GroupIds,
	}
	strategyList, err := m.List(cc, params, option)
	if err != nil {
		return nil, err
	}
	ResStrategyList := strategyList.(*marketing.MarketTrigResponse)
	return ResStrategyList, nil
}

// MarketTrig 获取客户分组下可执行的营销策略
func (m *MarketingRepo) MarketTrigList(cc ctx.Context, data *marketing.MarketTrigRequest) (*marketing.MarketTrigResponse, error) {

	option := &repos.ListOption{
		TplPrefix: "t_marketings/customers_marketing-trigger-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &marketing.MarketTrigResponse{}
		},
	}
	params := &marketing.MarketTrigRequest{
		EnterpriseID:   data.EnterpriseID,
		Condition:      data.Condition,
		StrategyStatus: marketing.StrategyStatusStart,
		GroupIds:       data.GroupIds,
		Count:          data.Count,
	}
	strategyList, err := m.List(cc, params, option)
	if err != nil {
		return nil, err
	}
	ResStrategyList := strategyList.(*marketing.MarketTrigResponse)
	return ResStrategyList, nil
}

// GetUserGroup 获取客户标签
func (m *MarketingRepo) GetUserGroup(cc ctx.Context, request marketing.UserGroup) ([]marketing.UserGroup, error) {
	stmt, err := m.DB.Template("t_marketings/select-user-group-list")
	if err != nil {
		return nil, err
	}

	userGroupList := []marketing.UserGroup{}
	resErr := stmt.NamedSelect(cc, &userGroupList, request)
	if resErr != nil {
		return nil, resErr
	}

	return userGroupList, nil
}

// GetUserGroup 获取标签中客户总数
func (m *MarketingRepo) GetUserCategory(cc ctx.Context, request marketing.UserCategory) int64 {
	stmt, err := m.DB.Template("t_marketings/select-user-category-count", m.Cond2Args(request))
	if err != nil {
		return 0
	}
	countResponseInfo := &sms.CountResponse{}
	getErr := stmt.NamedGet(cc, countResponseInfo, request)
	if getErr != nil {
		return 0
	}
	return countResponseInfo.Count
}

// GetUserGroup 获取标签用户中存在手机号的用户总数
func (m *MarketingRepo) GetMobileUserCount(cc ctx.Context, request marketing.UserCategory) (userList []*marketing.UserInfo, err error) {
	stmt, err := m.DB.Template("t_marketings/select-moblie-count", m.Cond2Args(request))
	if err != nil {
		return
	}
	rows, err := stmt.NamedQuery(cc, request)
	if nil != err && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	for rows.Next() {
		userInfo := &marketing.UserInfo{}
		rows.Scan(&userInfo.UnionID)
		if "" != userInfo.UnionID {
			userList = append(userList, userInfo)
		}
	}
	return
}

// GetShopUserCount 获取标签用户中商城用户总数
func (m *MarketingRepo) GetShopUserCount(cc ctx.Context, request marketing.UserCategory) (userList []*marketing.UserInfo, err error) {
	stmt, err := m.DB.Template("t_marketings/select-shop-count", m.Cond2Args(request))
	if err != nil {
		return
	}
	rows, err := stmt.NamedQuery(cc, request)
	if nil != err && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	for rows.Next() {
		userInfo := &marketing.UserInfo{}
		rows.Scan(&userInfo.UnionID)
		if "" != userInfo.UnionID {
			userList = append(userList, userInfo)
		}
	}
	return
}

// GetQyUserCount 获取标签用户中企微用户总数
func (m *MarketingRepo) GetQyUserCount(cc ctx.Context, request marketing.UserCategory) (userList []*marketing.UserInfo, err error) {
	stmt, err := m.DB.Template("t_marketings/select-qw-count", m.Cond2Args(request))
	if err != nil {
		return
	}
	rows, err := stmt.NamedQuery(cc, request)
	if nil != err && err != sql.ErrNoRows {
		return
	}
	defer rows.Close()
	for rows.Next() {
		userInfo := &marketing.UserInfo{}
		rows.Scan(&userInfo.UnionID)
		if "" != userInfo.UnionID {
			userList = append(userList, userInfo)
		}
	}
	return
}

// GetQyUserCount 获取标签用户中企微用户总数
func (m *MarketingRepo) GetWeChatOfficialUserCount(cc ctx.Context, request marketing.UserCategory) (userList []*marketing.UserInfo, err error) {
	stmt, err := m.DB.Template("t_marketings/select-we-chat-official-count", m.Cond2Args(request))
	if err != nil {
		return
	}
	rows, err := stmt.NamedQuery(cc, request)
	if nil != err && err != sql.ErrNoRows {
		return
	}
	defer rows.Close() //需要关闭
	for rows.Next() {
		userInfo := &marketing.UserInfo{}
		rows.Scan(&userInfo.UnionID)
		if "" != userInfo.UnionID {
			userList = append(userList, userInfo)
		}
	}
	return
}

// GetUserStatisticData 获取用户统计数据
func (m *MarketingRepo) GetUserStatisticData(ctx ctx.Context, enterpriseID int64, userID int64) (*marketing.UserStatisticEntity, error) {
	var entity marketing.UserStatisticEntity
	stmt, err := m.DB.Template("t_marketing_user_statistics/query-user-id")
	if err != nil {
		return nil, err
	}
	errStmt := stmt.Get(ctx, &entity, enterpriseID, userID, 0)
	if errStmt != nil {
		log.Info(ctx).Msgf("查询用户统计数据失败, 企业id:%d, userID:%d, err:%v", enterpriseID, userID, errStmt)
		return nil, errStmt
	}
	return &entity, nil
}

// AddUserStatisticData 添加用户统计数据
func (m *MarketingRepo) AddUserStatisticData(ctx ctx.Context, data *marketing.UserStatisticEntity) (int64, error) {
	stmt, err := m.DB.Template("t_marketing_user_statistics/insert")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UpdateUserStatisticData 更新用户统计数据
func (m *MarketingRepo) UpdateUserStatisticData(ctx ctx.Context, data *marketing.UserStatisticEntity) (int64, error) {
	stmt, err := m.DB.Template("t_marketing_user_statistics/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// MarketTrigChannelDetailByWhere 查询渠道信息
func (m *MarketingRepo) MarketTrigChannelDetailByWhere(ctx ctx.Context, req marketing.ChannelInfo) (*marketing.ChannelInfo, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = req
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-channel-by-where", &mapping)
	if err != nil {
		return nil, err
	}
	channelInfo := &marketing.ChannelInfo{}
	sErr := stmt.NamedGet(ctx, channelInfo, req)
	if sErr != nil {
		return nil, sErr
	}
	return channelInfo, nil
}

// 批量修改策略记录
func (m *MarketingRepo) BatchUpdateMarketStrategyRecord(ctx ctx.Context, request *marketing.BatchUpdateRecordsRequest) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = request.IDs
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-batch-update", &mapping)
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(ctx, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// SelectMsgIds 获取未发送消息列表
func (m *MarketingRepo) SelectMsgIds(ctx ctx.Context) ([]marketing.QyMsgSendResult, error) {
	mapping := make(map[string]interface{})
	stmt, err := m.DB.Template("t_qy_msg_send_result/select-msg-ids", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]marketing.QyMsgSendResult{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return *list, nil
}

// UpdateQyMsgSendResult 更新群发发送记录
func (m *MarketingRepo) UpdateQyMsgSendResult(ctx ctx.Context, data *marketing.QyMsgSendResult) (int64, error) {
	stmt, err := m.DB.Template("t_qy_msg_send_result/update")
	if err != nil {
		return 0, err
	}
	res, err := stmt.NamedExec(ctx, data)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//UpdateMarketStrategyRecordByEcuserid 营销策略触发记录修改-ecuserid
func (m *MarketingRepo) UpdateMarketStrategyRecordByEcuserid(ctx ctx.Context, data *marketing.QyMsgSendResult) (int64, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-records-update-ecuserid")
	if err != nil {
		return 0, err
	}
	result, sTmtErr := stmt.NamedExec(ctx, data)
	if sTmtErr != nil {
		return 0, sTmtErr
	}
	return result.RowsAffected()

}

// 添加企业群发结果
func (m *MarketingRepo) AddQySendMsgResult(c ctx.Context, request *marketing.QyMsgSendResult) error {
	stmt, err := m.DB.Template("t_qy_msg_send_result/add")
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// IsUserGroupUserInMarket 查询客户分组是否有被营销策略引用
func (m *MarketingRepo) IsUserGroupUserInMarket(ctx ctx.Context, groupID int64) (bool, error) {
	var entity marketing.MarketingStrategies
	stmt, err := m.DB.Template("t_marketings/is-user-group-in-market")
	if err != nil {
		return false, err
	}
	errStmt := stmt.Get(ctx, &entity, groupID, 0)
	if errStmt != nil && errStmt != sql.ErrNoRows {
		log.Info(ctx).Msgf("查询使用客户分组的营销策略失败, 客户分组id:%d, err:%v", groupID, errStmt)
		return false, errStmt
	}
	if errStmt == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

// MarketingUpdateStatusByTime 更新策略根据时间
func (m *MarketingRepo) MarketingUpdateStatusByTime(c ctx.Context, request *marketing.UpdateStatusRequest) error {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-update-status-by-time", m.Cond2Args(request))
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// UpdateStrategyStatusByIds 通过策略IDs修改策略状态
func (m *MarketingRepo) UpdateStrategyStatusByIds(c ctx.Context, status int, ids []int64) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := m.DB.Template("t_marketings/update-strategy-status-by-ids", &mapping)
	if err != nil {
		return err
	}
	updateData := &struct {
		StrategyStatus int   `db:"strategy_status"`
		UpdatedAt      int64 `db:"updated_at"`
		HistoryStatus  int   `db:"history_status"`
	}{
		StrategyStatus: status,
		UpdatedAt:      time.Now().Unix(),
	}
	_, sTmtErr := stmt.NamedExec(c, updateData)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// UpdateStrategyHistoryStatusByIds 通过策略IDs修改策略状态
func (m *MarketingRepo) UpdateStrategyHistoryStatusByIds(c ctx.Context, historyStatus int, ids []int64) error {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := m.DB.Template("t_marketings/update-strategy-history-status-by-ids", &mapping)
	if err != nil {
		return err
	}
	updateData := &struct {
		UpdatedAt     int64 `db:"updated_at"`
		HistoryStatus int   `db:"history_status"`
	}{
		UpdatedAt:     time.Now().Unix(),
		HistoryStatus: historyStatus,
	}
	_, sTmtErr := stmt.NamedExec(c, updateData)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// SimpleStrategyList 查询策略总数
func (m *MarketingRepo) SimpleStrategyList(cc ctx.Context, request marketing.ListRequest) ([]*marketing.MarketingStrategy, error) {
	request.StrategyStatusMultiCnt = len(request.StrategyStatusMulti)
	mapping := m.Cond2Args(request)
	offset := (request.Page - 1) * request.Size
	(*mapping)["page"] = repos.PageVO{
		Offset: offset,
		Limit:  request.Size,
	}
	stmt, err := m.DB.Template("t_marketings/select", mapping)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var dest []*marketing.MarketingStrategy
	sErr := stmt.NamedSelect(cc, &dest, request)
	if sErr != nil {
		return nil, sErr
	}
	return dest, nil
}

// ListChannelByStrategyIds 根据策略ID查询策略渠道
func (m *MarketingRepo) ListChannelByStrategyIds(cc ctx.Context, req marketing.StrategyChannelListRequest) (*[]*marketing.ChannelInfo, error) {
	stmt, err := m.DB.Template("t_marketings/select-channel-by-ids", m.Cond2Args(req))
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	dest := &[]*marketing.ChannelInfo{}
	sErr := stmt.NamedSelect(cc, dest, req)
	if sErr != nil {
		return nil, sErr
	}
	return dest, nil
}

// ListConditionByStrategyIds 根据策略ID查询策略条件
func (m *MarketingRepo) ListConditionByStrategyIds(cc ctx.Context, ids []int64) (*[]*marketing.StrategyConditionItem, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := m.DB.Template("t_marketings/select-condition-by-ids", &mapping)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	dest := &[]*marketing.StrategyConditionItem{}
	type selectStruct struct {
		Ids []int64 `json:"ids" db:"ids"`
	}
	req := &selectStruct{
		Ids: ids,
	}
	sErr := stmt.NamedSelect(cc, dest, req)
	if sErr != nil {
		return nil, sErr
	}
	return dest, nil
}
func (m *MarketingRepo) GetEnterList(cc ctx.Context) (enterList []int64, err error) {

	stmt, err := m.DB.Template("t_marketings/select-enterprise-id")
	if nil != err {
		return
	}
	err = stmt.Select(cc, &enterList)
	if nil != err {
		return
	}
	return
}

//获取商品详情
func (m *MarketingRepo) GetSpuDetail(c ctx.Context, SpuId int64) (marketing.SpuDetailRes, error) {
	spuDetail := marketing.SpuDetailRes{}
	stmt, err := m.DB.Template("t_marketings/select-spu-detail")
	if err != nil {
		log.Info(c).Msgf("GetSpuDetailError:%v", err)
		return spuDetail, err
	}
	err = stmt.Get(c, &spuDetail, SpuId)
	return spuDetail, err
}

//ShortUrlData(c ctx.Context, SpuId ShortUrlDataReq) error
func (m *MarketingRepo) ShortUrlData(c ctx.Context, request marketing.ShortUrlDataReq) (marketing.ShortUrlDataRes, error) {
	shortDetail := marketing.ShortUrlDataRes{}
	stmt, err := m.DB.Template("t_marketings/select-short-url-data")
	if err != nil {
		log.Info(c).Msgf("ShortUrlDataError:%v", err)
		return shortDetail, err
	}
	err = stmt.Get(c, &shortDetail, request.Id)
	return shortDetail, nil
}

func (m *MarketingRepo) ShortUrlDel(c ctx.Context, request marketing.ShortUrlDelReq) error {
	stmt, err := m.DB.Template("t_marketings/short-url-del")
	if err != nil {
		log.Info(c).Msgf("ShortUrlDelError:%v", err)
		return err
	}
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		return err
	}
	return nil
}

func (m *MarketingRepo) TwoLevelIds(c ctx.Context, request marketing.ShortUrlAddReq) (IdsRes []marketing.ShortUrlDataReq, err error) {
	stmt, err := m.DB.Template("t_marketings/short-url-row")
	if err != nil {
		log.Info(c).Msgf("TwoLevelIdsError:%v", err)
		return nil, err
	}
	err = stmt.Select(c, &IdsRes)
	if err != nil {
		log.Info(c).Msgf("sql select TwoLevelIdsError:%v", err)
		return nil, err
	}
	return IdsRes, nil
}

//添加短链
func (m *MarketingRepo) ShortUrlAdd(c ctx.Context, request marketing.ShortUrlAddReq) error {
	stmt, err := m.DB.Template("t_marketings/short-url-update")
	if request.Id == 0 {
		stmt, err = m.DB.Template("t_marketings/short-url-add")
		request.CreatedAt = time.Now().Unix()
	}
	if err != nil {
		log.Info(c).Msgf("sql模板报错:%v", err)
		return err
	}
	configParams, err := json.Marshal(request.MiddleData)
	if err != nil {
		return err
	}
	var shortId marketing.ShortUrlDataRes

	if request.Id == 0 {
		stmt, err = m.DB.Template("t_marketings/short-url-one")
		err := stmt.Get(c, &shortId, request.EnterpriseId, request.TwoLevelId)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if shortId.Id != 0 {
			return errors.New("该二级渠道已添加过")
		}
		stmt, err = m.DB.Template("t_marketings/short-url-add")
		request.CreatedAt = time.Now().Unix()
	}
	request.TmpMiddleData = configParams
	_, err = stmt.NamedExec(c, request)
	if err != nil {
		log.Info(c).Msgf("sql执行短链添加失败:%v", err)

		return err
	}
	return nil
}

// GetWxChannelList 获取场景值列表的信息
// @param c ctx.Context 上下文对象
func (m *MarketingRepo) GetWxChannelList(c ctx.Context) ([]*marketing.WxChannelInfo, error) {
	var wxChannelList = make([]*marketing.WxChannelInfo, 0)
	stmt, err := m.DB.Template("t_marketings/select-enumerate-list")
	if err != nil {
		log.Info(c).Msgf("GetWxChannelListError:%v", err)
		return nil, err
	}
	err = stmt.Select(c, &wxChannelList)
	return wxChannelList, err
}

func (m *MarketingRepo) GetMarketingStrategyIsUsed(cc ctx.Context, strategyId int64) (cnum int64, err error) {
	log.Info(cc).Msgf("GetMarketingStrategyIsUsed <-:strategyId = %v", strategyId)
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-is-used")
	if nil != err {
		log.Info(cc).Msgf("GetMarketingStrategyIsUsed Error:%v", err)
		return cnum, err
	}
	err = stmt.Get(cc, &cnum, strategyId)
	log.Info(cc).Msgf("GetMarketingStrategyIsUsed <-:strategyId = %v,cnum=%v", strategyId, cnum)
	if nil != err {
		log.Error(cc).Msgf("GetMarketingStrategyIsUsed Select Error:%v", err)
	}
	return cnum, err
}

func (m *MarketingRepo) GetWxCouponDetail(c ctx.Context, id int64) (*marketing.WxCouponDatail, error) {
	wxcd := &marketing.WxCouponDatail{}
	stmt, err := m.DB.Template("t_coupon/select-field-coupon")
	if err != nil {
		log.Info(c).Msgf("GetWxCouponDetail:%v", err)
		return wxcd, err
	}
	err = stmt.Get(c, wxcd, id)
	return wxcd, err
}

// SelectMarketingWaitStartList 获取待开始的策略列表
func (m *MarketingRepo) SelectMarketingWaitStartList(c ctx.Context, enterpriseID int64) (*[]marketing.MarketingStrategies, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-wait-start-list")
	if err != nil {
		return nil, err
	}
	var marketingListRequest marketing.MarketingStrategies
	marketingListRequest.EnterpriseID = enterpriseID
	ListRequest := &[]marketing.MarketingStrategies{}
	res, errStmt := stmt.NamedQuery(c, marketingListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

// SelectMarketingWaitStartList 获取待开始的策略列表
func (m *MarketingRepo) SelectMarketingWaitEndList(c ctx.Context, enterpriseID int64) (*[]marketing.MarketingStrategies, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-wait-end-list")
	if err != nil {
		return nil, err
	}
	var marketingListRequest marketing.MarketingStrategies
	marketingListRequest.EnterpriseID = enterpriseID
	ListRequest := &[]marketing.MarketingStrategies{}
	res, errStmt := stmt.NamedQuery(c, marketingListRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ListRequest)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ListRequest, nil
}

//
func (m *MarketingRepo) MarketingDetailChannelsByID(c ctx.Context, strategyId int64) ([]marketing.ChannelInfo, error) {
	stmt, err := m.DB.Template("t_marketings/marketing-strategy-detail-channels-by-id")
	if err != nil {
		return nil, err
	}
	var response []marketing.ChannelInfo
	sTmtErr := stmt.Select(c, &response, strategyId)
	if sTmtErr != nil {
		return nil, sTmtErr
	}
	return response, err
}

//获取策略的优惠券列表
func (m *MarketingRepo) GetStrategyCouponList(c ctx.Context, strategyId int64) (couponList []marketing.StrategyCoupon, err error) {
	stmt, err := m.DB.Template("t_marketings/select-strategy-coupons")
	if err != nil {
		return couponList, err
	}

	sTmtErr := stmt.Select(c, &couponList, strategyId)
	if sTmtErr != nil {
		return couponList, sTmtErr
	}
	return couponList, err
}

func (m *MarketingRepo) TxUpdateMarketingStrategy(c ctx.Context, md marketing.DeleteRequest, strategyInfo *marketing.MarketingStrategies, strategyCondition *marketing.MarketingStrategyCondition, relationInfos []marketing.RelationInfo, channelInfos []marketing.ChannelInfo, couponIdList []marketing.MarketingStrategyCoupons) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	txshare = tx
	defer func() {
		txshare = nil
	}()
	errboolChan := make(chan bool)
	go func(tx *sqlplus.Tx) {
		if err := m.DeleteMarketingStrategy(c, md); err != nil {
			errboolChan <- false
		}
		if err := m.DeleteMarketingStrategyCondition(c, md); err != nil {
			errboolChan <- false
		}
		if err := m.DeleteMarketingStrategyChannels(c, md); err != nil {
			errboolChan <- false
		}
		if err := m.DeleteMarketingStrategyRelations(c, md); err != nil {
			errboolChan <- false
		}
		if err := m.MarketingAdd(c, strategyInfo); err != nil {
			errboolChan <- false
		}
		if err := orm.Service.Marketing.UpdateRecordTriggerCount(strategyInfo.ID, strategyInfo.TriggerCount); err != nil {
			errboolChan <- false
		}
		if err := m.MarketingStrategyConditionAdd(c, strategyCondition); err != nil {
			errboolChan <- false
		}
		if len(relationInfos) > 0 {
			txErr := m.MarketingStrategyTriggerRelationAdd(c, relationInfos)
			if txErr != nil {
				errboolChan <- false
			}
		}
		if len(channelInfos) > 0 {
			txErr := m.MarketingStrategyChannelAdd(c, channelInfos)
			if txErr != nil {
				errboolChan <- false
			}
		}

		errboolChan <- true
	}(tx)

	if errnool := <-errboolChan; errnool {
		err := tx.Commit()
		return err
	} else {
		err := tx.Rollback()
		if err == nil {
			return errorx.New("保存失败")
		}
		return err
	}
}

// BatchUpdateStrategy 批量更新策略状态
func (m *MarketingRepo) BatchUpdateStrategy(c ctx.Context, strategyIDs []int64, strategyStatus int8, AbnormalPause int8) error {
	request := marketing.BatchUpdateStrategy{
		StrategyIDs:    strategyIDs,
		StrategyStatus: strategyStatus,
		AbnormalPause:  AbnormalPause,
		UpdatedAt:      time.Now().Unix(),
	}
	stmt, err := m.DB.Template("t_marketings/batch-update-marketing-strategy-status", m.Cond2Args(request))
	if err != nil {
		return err
	}
	_, sTmtErr := stmt.NamedExec(c, request)
	if sTmtErr != nil {
		return sTmtErr
	}
	return err
}

// GoOnStrategyTemplateIDs 获取投放中策略的模版ID
func (m *MarketingRepo) GoOnStrategyTemplateIDs(c ctx.Context, enterpriseID int64) (*[]*marketing.GetGoOnStrategyResp, error) {
	dest := &[]*marketing.GetGoOnStrategyResp{}
	mapping := make(map[string]interface{})
	mapping["enterpriseID"] = enterpriseID
	stmt, err := m.DB.Template("t_marketings/select-goon-strategy-template-ids", &mapping)
	if err != nil && err != sql.ErrNoRows {
		return dest, err
	}
	sErr := stmt.NamedSelect(c, dest, mapping)
	if sErr != nil {
		return dest, sErr
	}
	return dest, nil
}

func (m *MarketingRepo) AllShortUrlList(c ctx.Context, request marketing.ShortUrlListReq) (*[]*marketing.ShortUrlDataRes, error) {
	var data []*marketing.ShortUrlDataRes
	err := m.DB.DB.Select(c, &data, marketing.GetAllShortUrlSql, request.EnterpriseId)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

//根据企微群发ID获取策略ID和外部联系人union_id
func (m *MarketingRepo) GetQyMsgSendLogWithIdr(c ctx.Context, EnterpriseID int64, MsgID, externalUserid string) (int64, error) {
	var userID int64
	sqlv2 := "SELECT `user_id` FROM `marketing_strategy_records_v2` " +
		"WHERE `external_userid` = ? AND `msg_id` = ? AND `enterprise_id` = ? LIMIT 1"
	err := m.DB.DB.Get(c, &userID, sqlv2, externalUserid, MsgID, EnterpriseID)
	return userID, err
}
