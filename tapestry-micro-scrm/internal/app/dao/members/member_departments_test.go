package membersDao

import (
	"bytes"
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"micro-scrm/internal/app/dao"
)

func TestMemberDepartmentRemoveMethodsUpdateDeletedAtAndUpdatedAt(t *testing.T) {
	tests := []struct {
		name string
		run  func(*MemberDepartments) error
	}{
		{
			name: "remove by member id",
			run: func(d *MemberDepartments) error {
				return d.RemoveByMemberId(toolboxCtx.New(), 11, 22)
			},
		},
		{
			name: "remove by member id and department ids",
			run: func(d *MemberDepartments) error {
				return d.RemoveByMemberIdAndDepIds(toolboxCtx.New(), 11, 22, []int64{33, 44})
			},
		},
		{
			name: "remove by department id",
			run: func(d *MemberDepartments) error {
				return d.RemoveByDepartmentId(toolboxCtx.New(), 33, 22)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sqlLog bytes.Buffer
			db := newMemberDepartmentDryRunDB(t, &sqlLog)
			d := &MemberDepartments{
				dao:       dao.Dao{DB: db},
				tableName: "member_departments",
			}

			if err := tt.run(d); err != nil {
				t.Fatalf("remove member department relation: %v", err)
			}

			emittedSQL := strings.ToLower(sqlLog.String())
			deletedAt := assignmentValue(t, emittedSQL, "deleted_at")
			updatedAt := assignmentValue(t, emittedSQL, "updated_at")
			if deletedAt != updatedAt {
				t.Fatalf("deleted_at = %d, updated_at = %d; want the same timestamp: %s", deletedAt, updatedAt, emittedSQL)
			}
		})
	}
}

func assignmentValue(t *testing.T, sql, column string) int64 {
	t.Helper()
	match := regexp.MustCompile("`" + column + "`=([0-9]+)").FindStringSubmatch(sql)
	if len(match) != 2 {
		t.Fatalf("soft-delete SQL does not assign %s: %s", column, sql)
	}
	value, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		t.Fatalf("parse %s assignment: %v", column, err)
	}
	return value
}

func newMemberDepartmentDryRunDB(t *testing.T, sqlLog *bytes.Buffer) *gorm.DB {
	t.Helper()
	sqlDB, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatalf("create test SQL DB: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		Logger: logger.New(log.New(sqlLog, "", 0), logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
		}),
	})
	if err != nil {
		t.Fatalf("create test GORM DB: %v", err)
	}
	return db
}
