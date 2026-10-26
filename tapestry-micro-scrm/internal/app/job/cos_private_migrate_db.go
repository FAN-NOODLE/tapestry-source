package appJob

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

var cosPrivateIdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

type cosPrivateRow struct {
	ID    int64
	Value string
}

type cosPrivateRowRepository interface {
	ScanAfter(ctx.Context, cosPrivateField, int64, int) ([]cosPrivateRow, error)
	MaxID(ctx.Context, cosPrivateField) (int64, error)
	LoadByID(ctx.Context, cosPrivateField, int64) (cosPrivateRow, bool, error)
	UpdateIfUnchanged(ctx.Context, cosPrivateField, cosPrivateRow, string) (bool, error)
}

type cosPrivateRowStore struct {
	db                        *gorm.DB
	oldOriginHost, oldCdnHost string
}

func (s cosPrivateRowStore) ScanAfter(c ctx.Context, field cosPrivateField, afterID int64, limit int) ([]cosPrivateRow, error) {
	query, args, err := cosPrivateBuildScanQuery(field, afterID, limit, s.oldOriginHost, s.oldCdnHost)
	if err != nil {
		return nil, err
	}

	rows := make([]cosPrivateRow, 0)
	if err := s.db.WithContext(c).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s cosPrivateRowStore) MaxID(c ctx.Context, field cosPrivateField) (int64, error) {
	query, err := cosPrivateBuildMaxIDQuery(field)
	if err != nil {
		return 0, err
	}

	var result struct {
		MaxID sql.NullInt64 `gorm:"column:max_id"`
	}
	if err := s.db.WithContext(c).Raw(query).Scan(&result).Error; err != nil {
		return 0, err
	}
	if !result.MaxID.Valid {
		return 0, nil
	}
	return result.MaxID.Int64, nil
}

func (s cosPrivateRowStore) LoadByID(c ctx.Context, field cosPrivateField, id int64) (cosPrivateRow, bool, error) {
	database, table, column, err := cosPrivateValidatedFieldIdentifiers(field)
	if err != nil {
		return cosPrivateRow{}, false, err
	}

	query := fmt.Sprintf("SELECT id, %s AS value FROM %s.%s WHERE id = ? LIMIT 1", column, database, table)
	rows := make([]cosPrivateRow, 0, 1)
	if err := s.db.WithContext(c).Raw(query, id).Scan(&rows).Error; err != nil {
		return cosPrivateRow{}, false, err
	}
	if len(rows) == 0 {
		return cosPrivateRow{}, false, nil
	}
	return rows[0], true, nil
}

func (s cosPrivateRowStore) UpdateIfUnchanged(c ctx.Context, field cosPrivateField, row cosPrivateRow, newValue string) (bool, error) {
	query, args, err := cosPrivateBuildUpdateQuery(field, row, newValue)
	if err != nil {
		return false, err
	}

	result := s.db.WithContext(c).Exec(query, args...)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

func cosPrivateBuildScanQuery(field cosPrivateField, afterID int64, limit int, oldOriginHost, oldCdnHost string) (string, []interface{}, error) {
	database, table, column, err := cosPrivateValidatedFieldIdentifiers(field)
	if err != nil {
		return "", nil, err
	}
	if limit <= 0 {
		return "", nil, fmt.Errorf("scan limit must be positive")
	}

	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 4)
	args = append(args, afterID)
	for _, host := range []string{oldOriginHost, oldCdnHost} {
		if host == "" {
			continue
		}
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
		args = append(args, "%"+host+"%")
	}
	if len(conditions) == 0 {
		return "", nil, fmt.Errorf("at least one old host is required")
	}
	args = append(args, limit)

	query := fmt.Sprintf(
		"SELECT id, %s AS value FROM %s.%s WHERE id > ? AND (%s) ORDER BY id ASC LIMIT ?",
		column,
		database,
		table,
		strings.Join(conditions, " OR "),
	)
	return query, args, nil
}

func cosPrivateBuildUpdateQuery(field cosPrivateField, row cosPrivateRow, newValue string) (string, []interface{}, error) {
	database, table, column, err := cosPrivateValidatedFieldIdentifiers(field)
	if err != nil {
		return "", nil, err
	}

	query := fmt.Sprintf("UPDATE %s.%s SET %s = ? WHERE id = ? AND %s = ?", database, table, column, column)
	return query, []interface{}{newValue, row.ID, row.Value}, nil
}

func cosPrivateBuildMaxIDQuery(field cosPrivateField) (string, error) {
	database, table, _, err := cosPrivateValidatedFieldIdentifiers(field)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SELECT MAX(id) AS max_id FROM %s.%s", database, table), nil
}

func cosPrivateValidatedFieldIdentifiers(field cosPrivateField) (string, string, string, error) {
	database, err := cosPrivateQuoteIdentifier(field.Database)
	if err != nil {
		return "", "", "", err
	}
	table, err := cosPrivateQuoteIdentifier(field.Table)
	if err != nil {
		return "", "", "", err
	}
	column, err := cosPrivateQuoteIdentifier(field.Field)
	if err != nil {
		return "", "", "", err
	}
	return database, table, column, nil
}

func cosPrivateQuoteIdentifier(value string) (string, error) {
	if !cosPrivateIdentifierPattern.MatchString(value) {
		return "", fmt.Errorf("invalid SQL identifier %q", value)
	}
	return "`" + value + "`", nil
}
