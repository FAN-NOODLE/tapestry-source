package onlineWarehouseSvc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func TestReadOnlineWarehouseImportFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("online warehouse import"))
	}))
	defer server.Close()

	data, err := readOnlineWarehouseImportFile(ctx.New(), server.URL+"/order.xlsx")
	if err != nil {
		t.Fatalf("readOnlineWarehouseImportFile() error = %v", err)
	}

	if string(data) != "online warehouse import" {
		t.Fatalf("readOnlineWarehouseImportFile() = %q", string(data))
	}
}
