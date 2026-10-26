package thirdCrmSftpDao

import (
	"reflect"
	"testing"
)

func TestFileIdentityConditionsIncludeEnterpriseTypePathAndName(t *testing.T) {
	got := fileIdentityConditions(
		3360610327843968,
		11,
		"/outbound/mstr_location/20260819/mstr_location_20260819.csv",
		"mstr_location_20260819.csv",
	)
	want := map[string]interface{}{
		"enterprise_id":  int64(3360610327843968),
		"file_type":      11,
		"file_full_path": "/outbound/mstr_location/20260819/mstr_location_20260819.csv",
		"file_name":      "mstr_location_20260819.csv",
		"deleted_at":     0,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fileIdentityConditions() = %#v, want %#v", got, want)
	}
}
