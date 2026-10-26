package source

import (
	"database/sql"
	"testing"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	sourceModel "micro-scrm/internal/app/models/source"
)

func TestLinkKOSDaoSupportsTransactionDB(t *testing.T) {
	d := NewLinkKOSDao()
	tx := &gorm.DB{}
	if got := d.SetDb(tx); got != &d.Repo {
		t.Fatal("SetDb 应返回嵌入的 Repo")
	}
	if d.DB != tx {
		t.Fatal("SetDb 未切换到事务 DB")
	}
}

func TestLinkKOSDaoSaveDataDoesNotPanicWhenBeforeUpdateHookRuns(t *testing.T) {
	sqlDB, err := sql.Open("mysql", "")
	if err != nil {
		t.Fatalf("创建测试 SQL DB 失败：%v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("创建测试 GORM DB 失败：%v", err)
	}

	d := (&LinkKOSDao{}).WithTx(db)
	data := sourceModel.LinkKOSModel{}
	data.ID = 1

	if err := d.SaveData(toolboxCtx.New(), data); err != nil {
		t.Fatalf("SaveData() error = %v", err)
	}
}
