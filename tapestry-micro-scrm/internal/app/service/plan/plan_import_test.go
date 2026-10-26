package planService

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadTaskImportRemoteFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("task import file"))
	}))
	defer server.Close()

	data, err := readTaskImportRemoteFile(context.Background(), server.URL+"/import.xlsx")
	if err != nil {
		t.Fatalf("readTaskImportRemoteFile() error = %v", err)
	}

	if string(data) != "task import file" {
		t.Fatalf("readTaskImportRemoteFile() = %q", string(data))
	}
}
