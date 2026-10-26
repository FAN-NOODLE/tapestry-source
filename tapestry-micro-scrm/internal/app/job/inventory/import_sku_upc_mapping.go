package inventory

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/xjob"
)

const (
	importSkuUpcMappingBatchSize = 5000
	importSkuUpcMaxErrorSummary  = 20
)

type ImportSkuUpcMapping struct {
	name string
}

type skuUpcCsvRow struct {
	SkuItem string
	UpcCode string
}

type skuUpcImportSummary struct {
	TotalRows      int
	SuccessRows    int
	FailedRows     int
	DuplicateRows  int
	UniqueSkuCount int
	UniqueUpcCount int
	ErrorSummary   string
}

func NewImportSkuUpcMapping() *ImportSkuUpcMapping {
	return &ImportSkuUpcMapping{
		name: "库存查询 SKU UPC 映射导入",
	}
}

func (job *ImportSkuUpcMapping) Run(c ctx.Context, handler *xjob.JobHandler) error {
	specs := handler.Runner().Specs()
	if len(specs) < 3 {
		return fmt.Errorf("usage: importInventorySkuUpcMapping <enterprise_id> <csv_path> <effective_month>")
	}

	enterpriseId, err := strconv.ParseInt(specs[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid enterprise_id %s: %w", specs[0], err)
	}

	csvPath := specs[1]
	effectiveMonth := strings.TrimSpace(specs[2])
	if effectiveMonth == "" {
		return fmt.Errorf("effective_month is required")
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	batchDao := microScrmDao.NewInventorySkuUpcImportBatchDao()
	mappingDao := microScrmDao.NewInventorySkuUpcMappingDao()
	batchNo := fmt.Sprintf("inventory_sku_upc_%s_%d", strings.ReplaceAll(effectiveMonth, "-", ""), time.Now().Unix())
	batch := microScrmModel.InventorySkuUpcImportBatchModel{
		BaseModel:      models.CreateBase(app.UID()),
		BatchNo:        batchNo,
		EffectiveMonth: effectiveMonth,
		SourceFile:     filepath.Base(csvPath),
		Status:         microScrmModel.InventorySkuUpcBatchStatusProcessing,
		IsActive:       microScrmModel.InventorySkuUpcBatchInactive,
		MappingEnabled: microScrmModel.InventorySkuUpcMappingDisabled,
	}
	if err = batchDao.CreateData(c, enterpriseId, batch); err != nil {
		return err
	}

	rows, summary, err := parseSkuUpcCsvRows(file)
	if err != nil {
		job.markBatchFailed(c, enterpriseId, batch.ID, summary, err)
		return err
	}
	if len(rows) == 0 {
		err = fmt.Errorf("csv has no valid sku upc mapping rows")
		job.markBatchFailed(c, enterpriseId, batch.ID, summary, err)
		return err
	}

	for start := 0; start < len(rows); start += importSkuUpcMappingBatchSize {
		end := start + importSkuUpcMappingBatchSize
		if end > len(rows) {
			end = len(rows)
		}
		mappings := newInventorySkuUpcMappingModels(batchNo, effectiveMonth, rows[start:end])
		if err = mappingDao.CreateInBatchesData(c, enterpriseId, mappings); err != nil {
			job.markBatchFailed(c, enterpriseId, batch.ID, summary, err)
			return err
		}
	}

	if err = batchDao.DisableActiveBatches(c, enterpriseId); err != nil {
		job.markBatchFailed(c, enterpriseId, batch.ID, summary, err)
		return err
	}

	err = batchDao.UpdateData(c, enterpriseId, map[string]interface{}{"id": batch.ID}, map[string]interface{}{
		"total_rows":       summary.TotalRows,
		"success_rows":     summary.SuccessRows,
		"failed_rows":      summary.FailedRows,
		"duplicate_rows":   summary.DuplicateRows,
		"unique_sku_count": summary.UniqueSkuCount,
		"unique_upc_count": summary.UniqueUpcCount,
		"status":           microScrmModel.InventorySkuUpcBatchStatusSuccess,
		"is_active":        microScrmModel.InventorySkuUpcBatchActive,
		"mapping_enabled":  microScrmModel.InventorySkuUpcMappingEnabled,
		"error_summary":    summary.ErrorSummary,
		"updated_at":       time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	log.Info(c).Msgf("inventory sku upc mapping import success enterprise_id=%d batch_no=%s success_rows=%d failed_rows=%d duplicate_rows=%d", enterpriseId, batchNo, summary.SuccessRows, summary.FailedRows, summary.DuplicateRows)
	return nil
}

func (job *ImportSkuUpcMapping) markBatchFailed(c ctx.Context, enterpriseId, batchId int64, summary skuUpcImportSummary, importErr error) {
	if summary.ErrorSummary == "" && importErr != nil {
		summary.ErrorSummary = importErr.Error()
	}
	err := microScrmDao.NewInventorySkuUpcImportBatchDao().UpdateData(c, enterpriseId, map[string]interface{}{"id": batchId}, map[string]interface{}{
		"total_rows":       summary.TotalRows,
		"success_rows":     summary.SuccessRows,
		"failed_rows":      summary.FailedRows,
		"duplicate_rows":   summary.DuplicateRows,
		"unique_sku_count": summary.UniqueSkuCount,
		"unique_upc_count": summary.UniqueUpcCount,
		"status":           microScrmModel.InventorySkuUpcBatchStatusFailed,
		"is_active":        microScrmModel.InventorySkuUpcBatchInactive,
		"mapping_enabled":  microScrmModel.InventorySkuUpcMappingDisabled,
		"error_summary":    summary.ErrorSummary,
		"updated_at":       time.Now().Unix(),
	})
	if err != nil {
		log.Error(c).Msgf("inventory sku upc mapping import mark failed error enterprise_id=%d batch_id=%d err=%v", enterpriseId, batchId, err)
	}
}

func parseSkuUpcCsvRows(reader io.Reader) ([]skuUpcCsvRow, skuUpcImportSummary, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	rows := make([]skuUpcCsvRow, 0)
	seenPair := make(map[string]struct{})
	seenSku := make(map[string]struct{})
	seenUpc := make(map[string]struct{})
	errors := make([]string, 0, importSkuUpcMaxErrorSummary)
	summary := skuUpcImportSummary{}

	lineNo := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		lineNo++
		if err != nil {
			summary.FailedRows++
			appendSkuUpcImportError(&errors, fmt.Sprintf("line %d: %v", lineNo, err))
			continue
		}
		summary.TotalRows++
		if len(record) < 2 {
			summary.FailedRows++
			appendSkuUpcImportError(&errors, fmt.Sprintf("line %d: missing upc or sku column", lineNo))
			continue
		}

		upcCode := strings.TrimSpace(record[0])
		skuItem := strings.ToUpper(strings.TrimSpace(record[1]))
		if upcCode == "" || skuItem == "" {
			summary.FailedRows++
			appendSkuUpcImportError(&errors, fmt.Sprintf("line %d: empty upc or sku", lineNo))
			continue
		}

		pairKey := skuItem + "|" + upcCode
		if _, ok := seenPair[pairKey]; ok {
			summary.DuplicateRows++
			continue
		}
		seenPair[pairKey] = struct{}{}
		seenSku[skuItem] = struct{}{}
		seenUpc[upcCode] = struct{}{}
		rows = append(rows, skuUpcCsvRow{
			SkuItem: skuItem,
			UpcCode: upcCode,
		})
		summary.SuccessRows++
	}

	summary.UniqueSkuCount = len(seenSku)
	summary.UniqueUpcCount = len(seenUpc)
	summary.ErrorSummary = strings.Join(errors, "; ")
	return rows, summary, nil
}

func appendSkuUpcImportError(errors *[]string, message string) {
	if len(*errors) >= importSkuUpcMaxErrorSummary {
		return
	}
	*errors = append(*errors, message)
}

func newInventorySkuUpcMappingModels(batchNo, effectiveMonth string, rows []skuUpcCsvRow) []microScrmModel.InventorySkuUpcMappingModel {
	now := time.Now().Unix()
	mappings := make([]microScrmModel.InventorySkuUpcMappingModel, 0, len(rows))
	for _, row := range rows {
		mappings = append(mappings, microScrmModel.InventorySkuUpcMappingModel{
			BaseModel: models.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
			BatchNo:        batchNo,
			SkuItem:        row.SkuItem,
			UpcCode:        row.UpcCode,
			EffectiveMonth: effectiveMonth,
		})
	}
	return mappings
}
