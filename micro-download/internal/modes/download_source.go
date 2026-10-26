package modes

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/uid"
	"micro-download/internal/common/file"
	"micro-download/internal/modes/sql_tmp"
	"micro-download/internal/modules/download/delivery"
	"micro-download/internal/pkg/db"
	"time"
)

type DownloadSource struct {
	Mode
	ID            int64  `json:"id" db:"id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	StoreID       int64  `json:"store_id" db:"store_id"`
	OperationName string `json:"operation_name" db:"operation_name"`
	OperationID   int64  `json:"operation_id" db:"operation_id"`
	Type          int8   `json:"type" db:"type"`
	Name          string `json:"name" db:"name"`
	SourceType    int8   `json:"source_type" db:"source_type"`
	Path          string `json:"path" db:"path"`
	SourceData    string `json:"source_data" db:"source_data"`
	Status        int8   `json:"status" db:"status"`
	FailedReason  string `json:"failed_reason" db:"failed_reason"`
	ExpirationAt  int64  `json:"expiration_at" db:"expiration_at"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	CompleteAt    int64  `json:"complete_at" db:"complete_at"`
	RetryAt       int64  `json:"retry_at" db:"retry_at"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	Label         string `json:"label" db:"label"`
	IsAuto        int8   `json:"is_auto" db:"is_auto"`
}

// SourceCreate 创建资源
// @param whenInsert 入库前回调返回模型本事 返回错误将不继续执行
func (d *DownloadSource) SourceCreate(ctx ctx.Context, whenInsert func(d *DownloadSource) (error, *DownloadSource)) (error, int64) {
	params := d.Params.(delivery.CreateExport)
	if d.ID == 0 {
		d.ID = uid.ID(ctx)
	}
	d.EnterpriseID = params.EnterpriseID
	d.StoreID = params.StoreID
	d.OperationName = params.OperationName
	d.OperationID = params.OperationID
	d.Type = params.Type
	d.Name = file.GenerateFileName(params.Filename)
	d.Path = ""
	d.SourceType = params.SourceType
	d.SourceData = ""
	d.Status = delivery.StatusCreated
	d.FailedReason = ""
	d.ExpirationAt = 0
	d.CreatedAt = time.Now().Unix()
	d.CompleteAt = 0
	d.RetryAt = 0
	d.DeletedAt = 0
	d.Label = params.Label
	d.IsAuto = 0
	err, object := whenInsert(d)
	if err != nil {
		return err, 0
	}
	_, err = db.DB("microScrm").NamedExec(ctx, sql_tmp.SourceInsert(), object)
	if err != nil {
		return err, 0
	}
	return nil, d.ID
}

func (d *DownloadSource) ExportList(ctx ctx.Context) (error, *delivery.ExportListRes) {
	res := &delivery.ExportListRes{}
	params := d.Params.(delivery.ExportListReq)
	row := db.DB("microScrm").QueryRow(ctx, sql_tmp.SourceList("count(*)", params))
	count := new(int64)
	err := row.Scan(count)
	if err != nil {
		return err, nil
	}
	res.Count = *count
	sql := sql_tmp.SourceList("id, name, type, operation_name, operation_id,status, created_at, failed_reason,path", params)
	log.Info(ctx).Msgf("列表SQL:%s", sql)
	query, err := db.DB("microScrm").Query(ctx, sql)
	if err != nil {
		return err, nil
	}
	defer query.Close()
	list := new([]delivery.ExportList)
	err = sqlplus.ScanAll(query, list)
	if err != nil {
		return err, nil
	}
	res.List = *list
	return nil, res
}

func (d *DownloadSource) Detail(ctx ctx.Context) (error, *DownloadSource) {
	params := d.Params.(delivery.FileID)
	err := db.DB("microScrm").NamedGet(ctx, d, sql_tmp.FullDetail(), params)
	if err != nil {
		return err, nil
	}
	return nil, d
}

func (d *DownloadSource) LatestDetailByLabel(ctx ctx.Context) (error, *DownloadSource) {
	params := d.Params.(delivery.Label)
	err := db.DB("microScrm").NamedGet(ctx, d, sql_tmp.LatestDetailByLabel(), params)
	if err != nil {
		return err, nil
	}
	return nil, d
}

func (d *DownloadSource) ChangeStatus(ctx ctx.Context, status int8) error {
	param := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":     param.FileID,
		"status": status,
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.ChangeStatus(), newParam)
	return err
}

func (d *DownloadSource) Delete(ctx ctx.Context) error {
	params := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":   params.FileID,
		"time": time.Now().Unix(),
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.DeleteById(), newParam)
	return err
}

func (d *DownloadSource) Retry(ctx ctx.Context) error {
	params := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":       params.FileID,
		"status":   delivery.StatusCreated,
		"retry_at": time.Now().Unix(),
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.RetryById(), newParam)
	return err
}

func (d *DownloadSource) Failed(ctx ctx.Context, reason string) error {
	params := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":            params.FileID,
		"status":        delivery.StatusFailed,
		"failed_reason": reason,
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.Failed(), newParam)
	return err
}

func (d *DownloadSource) Success(ctx ctx.Context, path string) error {
	params := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":          params.FileID,
		"status":      delivery.StatusComplete,
		"path":        path,
		"complete_at": time.Now().Unix(),
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.Success(), newParam)
	return err
}

func (d *DownloadSource) ChangeAuto(ctx ctx.Context, isAuto int8) error {
	param := d.Params.(delivery.FileID)
	newParam := map[string]interface{}{
		"id":      param.FileID,
		"is_auto": isAuto,
	}
	_, err := db.DB("microScrm").NamedExec(ctx, sql_tmp.ChangeAuto(), newParam)
	return err
}

func (d *DownloadSource) GetNoAutoOne(ctx ctx.Context) (error, *DownloadSource) {
	params := d.Params.(delivery.OperationReq)
	err := db.DB("microScrm").NamedGet(ctx, d, sql_tmp.GetNoAutoOne(), params)
	if err != nil {
		return err, nil
	}
	return nil, d
}
