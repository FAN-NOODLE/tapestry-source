package ws_trace

import (
	"micro-scrm/internal/pkg/helper"
	"strconv"

	"micro-scrm/internal/pkg/errs"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/assist"
)

func (s *Service) WsAdd(c ctx.Context, enterpriseID int64, wst WsTrace) (wstr *WsTrace, err error) {
	keyM := assist.StructToMap(wst)
	wst.UkVoucher, err = GetWsTracePrimary(enterpriseID, wst.SourceType, keyM)
	if err != nil {
		return nil, err
	}
	wstr, err = s.repo.InsertWsTrace(c, enterpriseID, wst)
	if err != nil {
		return nil, errs.ActionErr
	}
	return wstr, err
}

func (s *Service) SelectWsTracesBySourceType(c ctx.Context, enterpriseID int64, sourceType int64) ([]WsTrace, error) {
	return s.repo.SelectWsTraces(c, enterpriseID, sourceType)
}

// 组装生成聚合主键
func GetWsTracePrimary(enterpriseID int64, sourceType int64, m map[string]interface{}) (string, error) {
	primaryKey := strconv.FormatInt(enterpriseID, 10) + strconv.FormatInt(sourceType, 10)
	switch sourceType {
	case SourceTypeOneToOne:
		a := assist.StrVal(m["member_id"])
		b := assist.StrVal(m["wx_external_id"])
		primaryKey += a + b
		if a == "0" || b == "0" {
			primaryKey = ""
		}
	case SourceTypeGroupChat:
		a := assist.StrVal(m["member_id"])
		b := assist.StrVal(m["group_chat_id"])
		primaryKey += a + b
		if a == "0" || b == "0" {
			primaryKey = ""
		}
	case SourceTypeTask:
		a := assist.StrVal(m["member_id"])
		b := assist.StrVal(m["main_task_id"])
		c := assist.StrVal(m["child_task_id"])
		d := assist.StrVal(m["sub_task_id"])
		e := assist.StrVal(m["group_chat_id"])
		primaryKey += a + b + c + d + e
		if a == "0" || b == "0" || c == "0" || d == "0" {
			primaryKey = ""
		}
	case SourceTypeGroupChatSend:
		a := assist.StrVal(m["member_id"])
		b := assist.StrVal(m["qw_group_send_id"])
		primaryKey += a + b
		if a == "0" || b == "0" {
			primaryKey = ""
		}
	default:
		primaryKey = ""
	}
	if primaryKey == "" {
		return "", errs.DataTransformErr
	}
	return helper.Md5V(primaryKey), nil
}
