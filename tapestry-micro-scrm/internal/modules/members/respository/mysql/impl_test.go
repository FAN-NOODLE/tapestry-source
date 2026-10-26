package mysql

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
)

func TestDeleteMemberDepartmentRelationTemplateUpdatesBothTimestamps(t *testing.T) {
	templates := sqlplus.ScanSQLTemplate(filepath.Join("..", "..", "..", "..", "..", "db", "templates"))
	query, err := templates.Execute("t_member_departments/delete_relation")
	if err != nil {
		t.Fatalf("load delete_relation template: %v", err)
	}
	query, args, err := sqlx.Named(query, map[string]interface{}{
		"enterprise_id": int64(123),
		"deleted_at":    int64(456),
	})
	if err != nil {
		t.Fatalf("bind delete_relation template: %v", err)
	}
	normalized := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(query, " "))
	want := "UPDATE `member_departments` SET `deleted_at` = ?, `updated_at` = ? WHERE `enterprise_id` = ?"
	if normalized != want {
		t.Fatalf("delete_relation SQL = %q, want %q", normalized, want)
	}
	if len(args) != 3 || args[0] != int64(456) || args[1] != int64(456) || args[2] != int64(123) {
		t.Fatalf("delete_relation args = %#v, want [456 456 123]", args)
	}
}
