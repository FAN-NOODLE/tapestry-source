package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"micro-scrm/internal/pkg/db"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/ws_trace"
)

type WsTraceRepo struct {
	*repos.SQLRepo
}

func (wr *WsTraceRepo) InsertWsTrace(c ctx.Context, enterpriseID int64, wst ws_trace.WsTrace) (wstr *ws_trace.WsTrace, err error) {
	wst.WsTraceID = strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(9)
	now := time.Now().Unix()
	queryStr := db.MakeUpsertSQL(fmt.Sprintf("micro_scrm_enterprise_%d.ws_trace", enterpriseID), []db.UpsertColumn{
		{Name: "ws_trace_id", Insert: true},
		{Name: "uk_voucher", Insert: true},
		{Name: "member_id", Insert: true},
		{Name: "wx_external_id", Insert: true},
		{Name: "unionid", Insert: true, Update: true},
		{Name: "group_chat_id", Insert: true},
		{Name: "main_task_id", Insert: true},
		{Name: "child_task_id", Insert: true},
		{Name: "sub_task_id", Insert: true},
		{Name: "qw_group_send_id", Insert: true},
		{Name: "source_type", Insert: true},
		{Name: "created_at", Insert: true},
		{Name: "updated_at", Insert: true, Update: true},
	})
	_, err = wr.DB.DB.Exec(c, queryStr, wst.WsTraceID, wst.UkVoucher, wst.MemberID, wst.WxExternalID, wst.Unionid, wst.GroupChatID, wst.MainTaskID, wst.ChildTaskID, wst.SubTaskID, wst.QwGroupSendID, wst.SourceType, now, now, wst.Unionid, now)
	if err != nil {
		return nil, err
	}
	return wr.GetWsTraceByPrimary(c, enterpriseID, wst.UkVoucher)
}

func (wr *WsTraceRepo) GetWsTraceByPrimary(c ctx.Context, enterpriseID int64, ukVoucher string) (*ws_trace.WsTrace, error) {
	wstr := &ws_trace.WsTrace{}
	getstr := fmt.Sprintf("select ws_trace_id,uk_voucher,member_id,wx_external_id,unionid,group_chat_id,main_task_id,child_task_id,qw_group_send_id,source_type,created_at,updated_at from micro_scrm_enterprise_%d.ws_trace where uk_voucher=?", enterpriseID)
	err := wr.DB.DB.Get(c, wstr, getstr, ukVoucher)
	if err == sql.ErrNoRows {
		wstr = nil
	}
	return wstr, err
}

func (wr *WsTraceRepo) GetWsTraceByID(c ctx.Context, enterpriseID int64, wsTraceID int64) (*ws_trace.WsTrace, error) {
	wstr := &ws_trace.WsTrace{}
	getstr := fmt.Sprintf("select ws_trace_id,uk_voucher,member_id,wx_external_id,unionid,group_chat_id,main_task_id,child_task_id,qw_group_send_id,source_type,created_at,updated_at from micro_scrm_enterprise_%d.ws_trace where ws_trace_id=?", enterpriseID)
	err := wr.DB.DB.Get(c, wstr, getstr, wsTraceID)
	if err != nil {
		return nil, err
	}
	return wstr, err
}

func (wr *WsTraceRepo) SelectWsTraces(c ctx.Context, enterpriseID int64, sourceType int64) ([]ws_trace.WsTrace, error) {
	var list []ws_trace.WsTrace
	getstr := fmt.Sprintf("select ws_trace_id,uk_voucher,member_id,wx_external_id,unionid,group_chat_id,main_task_id,child_task_id,sub_task_id,qw_group_send_id,source_type,created_at,updated_at from micro_scrm_enterprise_%d.ws_trace where source_type=?", enterpriseID)
	err := wr.DB.Select(c, &list, getstr, sourceType)
	if err != nil {
		return nil, err
	}
	return list, err
}
