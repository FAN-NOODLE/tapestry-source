package appJob

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"sync"
	"testing"

	toolboxctx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCosPrivateBuildScanQueryUsesCursor(t *testing.T) {
	field := cosPrivateField{"micro_scrm", "download_resource", "source_data"}
	query, args, err := cosPrivateBuildScanQuery(field, 12, 500, "old.example.com", "cdn.example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := "SELECT id, `source_data` AS value FROM `micro_scrm`.`download_resource` WHERE id > ? AND (`source_data` LIKE ? OR `source_data` LIKE ?) ORDER BY id ASC LIMIT ?"
	if query != want || len(args) != 4 {
		t.Fatalf("query=%q args=%#v", query, args)
	}
}

func TestCosPrivateBuildUpdateQueryUsesOriginalValue(t *testing.T) {
	query, args, err := cosPrivateBuildUpdateQuery(cosPrivateField{"micro_scrm", "download_resource", "source_data"}, cosPrivateRow{ID: 7, Value: "old"}, "new")
	if err != nil {
		t.Fatal(err)
	}
	want := "UPDATE `micro_scrm`.`download_resource` SET `source_data` = ? WHERE id = ? AND `source_data` = ?"
	if query != want || !reflect.DeepEqual(args, []interface{}{"new", int64(7), "old"}) {
		t.Fatalf("query=%q args=%#v", query, args)
	}
}

func TestCosPrivateIdentifierRejectsSQL(t *testing.T) {
	_, _, err := cosPrivateBuildScanQuery(cosPrivateField{"micro_scrm", "x`;DROP_TABLE_x", "url"}, 0, 10, "old", "cdn")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCosPrivateBuildMaxIDQueryUsesValidatedField(t *testing.T) {
	query, err := cosPrivateBuildMaxIDQuery(cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	if err != nil {
		t.Fatal(err)
	}
	if want := "SELECT MAX(id) AS max_id FROM `micro_scrm`.`download_resource`"; query != want {
		t.Fatalf("query=%q want=%q", query, want)
	}
}

func TestCosPrivateRowStoreMaxIDReturnsZeroForEmptyTable(t *testing.T) {
	db, state := newCosPrivateMaxIDTestDB(t)
	state.values = [][]driver.Value{{nil}}

	got, err := (cosPrivateRowStore{db: db}).MaxID(toolboxctx.New(), cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("max id=%d want=0", got)
	}
}

func TestCosPrivateRowStoreMaxIDPropagatesDatabaseError(t *testing.T) {
	db, state := newCosPrivateMaxIDTestDB(t)
	state.err = errors.New("database unavailable")

	_, err := (cosPrivateRowStore{db: db}).MaxID(toolboxctx.New(), cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	if !errors.Is(err, state.err) {
		t.Fatalf("error=%v want=%v", err, state.err)
	}
}

func TestCosPrivateRowStoreMaxIDPropagatesParseError(t *testing.T) {
	db, state := newCosPrivateMaxIDTestDB(t)
	state.values = [][]driver.Value{{"not-an-int64"}}

	_, err := (cosPrivateRowStore{db: db}).MaxID(toolboxctx.New(), cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestCosPrivateRowStoreMaxIDPreservesLargestInt64(t *testing.T) {
	db, state := newCosPrivateMaxIDTestDB(t)
	state.values = [][]driver.Value{{int64(9223372036854775807)}}

	got, err := (cosPrivateRowStore{db: db}).MaxID(toolboxctx.New(), cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(9223372036854775807); got != want {
		t.Fatalf("max id=%d want=%d", got, want)
	}
}

func TestCosPrivateRowStoreMaxIDRejectsMaliciousIdentifier(t *testing.T) {
	db, _ := newCosPrivateMaxIDTestDB(t)

	_, err := (cosPrivateRowStore{db: db}).MaxID(toolboxctx.New(), cosPrivateField{"micro_scrm", "download_resource;DROP_TABLE_x", "source_data"})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

const cosPrivateMaxIDTestDriverName = "cos_private_max_id_test"

var (
	cosPrivateMaxIDTestDriverOnce sync.Once
	cosPrivateMaxIDTestStateMu    sync.Mutex
	cosPrivateMaxIDTestState      *cosPrivateMaxIDDriverState
)

type cosPrivateMaxIDDriverState struct {
	err    error
	query  string
	values [][]driver.Value
}

type cosPrivateMaxIDTestDriver struct{}

func (cosPrivateMaxIDTestDriver) Open(string) (driver.Conn, error) {
	return cosPrivateMaxIDTestConn{}, nil
}

type cosPrivateMaxIDTestConn struct{}

func (cosPrivateMaxIDTestConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is unsupported")
}
func (cosPrivateMaxIDTestConn) Close() error { return nil }
func (cosPrivateMaxIDTestConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are unsupported")
}

func (cosPrivateMaxIDTestConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	cosPrivateMaxIDTestStateMu.Lock()
	defer cosPrivateMaxIDTestStateMu.Unlock()
	cosPrivateMaxIDTestState.query = query
	if cosPrivateMaxIDTestState.err != nil {
		return nil, cosPrivateMaxIDTestState.err
	}
	return &cosPrivateMaxIDTestRows{values: cosPrivateMaxIDTestState.values}, nil
}

type cosPrivateMaxIDTestRows struct {
	values [][]driver.Value
	index  int
}

func (r *cosPrivateMaxIDTestRows) Columns() []string { return []string{"max_id"} }
func (r *cosPrivateMaxIDTestRows) Close() error      { return nil }
func (r *cosPrivateMaxIDTestRows) Next(dest []driver.Value) error {
	if r.index == len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}

func newCosPrivateMaxIDTestDB(t *testing.T) (*gorm.DB, *cosPrivateMaxIDDriverState) {
	t.Helper()
	cosPrivateMaxIDTestDriverOnce.Do(func() {
		sql.Register(cosPrivateMaxIDTestDriverName, cosPrivateMaxIDTestDriver{})
	})
	state := &cosPrivateMaxIDDriverState{}
	cosPrivateMaxIDTestStateMu.Lock()
	cosPrivateMaxIDTestState = state
	cosPrivateMaxIDTestStateMu.Unlock()

	sqlDB, err := sql.Open(cosPrivateMaxIDTestDriverName, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db, state
}
