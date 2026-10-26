package membersDao

import (
	"bytes"
	"strings"
	"testing"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao"
)

func TestGetUnscopedMemberByMobilesIncludesSoftDeletedMembers(t *testing.T) {
	var sqlLog bytes.Buffer
	db := newMemberDepartmentDryRunDB(t, &sqlLog)
	d := &Members{
		dao:       dao.Dao{DB: db},
		tableName: "members",
	}

	if _, err := d.GetUnscopedMemberByMobiles(toolboxCtx.New(), 1001, []string{"13800000000"}); err != nil {
		t.Fatalf("query uploader members: %v", err)
	}

	emittedSQL := strings.ToLower(sqlLog.String())
	if !strings.Contains(emittedSQL, "`enterprise_id` = 1001") || !strings.Contains(emittedSQL, "`mobile` = '13800000000'") {
		t.Fatalf("query does not constrain enterprise and mobile: %s", emittedSQL)
	}
	if strings.Contains(emittedSQL, "deleted_at") {
		t.Fatalf("uploader lookup unexpectedly excludes soft-deleted members: %s", emittedSQL)
	}
}
