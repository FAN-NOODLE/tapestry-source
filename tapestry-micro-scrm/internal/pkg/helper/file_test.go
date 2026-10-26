package helper

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"micro-scrm/internal/config"
)

func TestReadRemoteFileUsesHTTPForNonCOSURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("id,name\n1,test\n"))
	}))
	defer server.Close()

	data, err := ReadRemoteFile(context.Background(), server.URL+"/data.csv")
	if err != nil {
		t.Fatalf("ReadRemoteFile() error = %v", err)
	}

	if string(data) != "id,name\n1,test\n" {
		t.Fatalf("ReadRemoteFile() = %q", string(data))
	}
}

func TestReadRemoteFileReturnsErrorForHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer server.Close()

	_, err := ReadRemoteFile(context.Background(), server.URL+"/data.csv")
	if err == nil {
		t.Fatal("ReadRemoteFile() error = nil, want non-nil")
	}
}

func TestReadRemoteFileUsesCOSDownloadForCOSDomain(t *testing.T) {
	oldCOS := config.CfgData.Storage.COS
	config.CfgData.Storage.COS.Schema = "https"
	config.CfgData.Storage.COS.CdnHost = "cdn.example.com"
	config.CfgData.Storage.COS.Bucket = "bucket"
	config.CfgData.Storage.COS.Origin = "ap-shanghai"
	defer func() {
		config.CfgData.Storage.COS = oldCOS
	}()

	oldReadCOSRemoteFile := readCOSRemoteFile
	defer func() {
		readCOSRemoteFile = oldReadCOSRemoteFile
	}()

	var gotURL string
	readCOSRemoteFile = func(ctx context.Context, fileURL string) ([]byte, error) {
		gotURL = fileURL
		return []byte("id,name\n1,cos\n"), nil
	}

	data, err := ReadRemoteFile(context.Background(), "https://cdn.example.com/third_sftp/file.csv")
	if err != nil {
		t.Fatalf("ReadRemoteFile() error = %v", err)
	}

	if gotURL != "https://cdn.example.com/third_sftp/file.csv" {
		t.Fatalf("readCOSRemoteFile URL = %q", gotURL)
	}
	if string(data) != "id,name\n1,cos\n" {
		t.Fatalf("ReadRemoteFile() = %q", string(data))
	}
}

func TestReadRemoteFileReturnsCOSDownloadError(t *testing.T) {
	oldCOS := config.CfgData.Storage.COS
	config.CfgData.Storage.COS.Schema = "https"
	config.CfgData.Storage.COS.CdnHost = "cdn.example.com"
	defer func() {
		config.CfgData.Storage.COS = oldCOS
	}()

	oldReadCOSRemoteFile := readCOSRemoteFile
	defer func() {
		readCOSRemoteFile = oldReadCOSRemoteFile
	}()

	readCOSRemoteFile = func(ctx context.Context, fileURL string) ([]byte, error) {
		return nil, errors.New("download failed")
	}

	_, err := ReadRemoteFile(context.Background(), "https://cdn.example.com/third_sftp/file.csv")
	if err == nil {
		t.Fatal("ReadRemoteFile() error = nil, want non-nil")
	}
}

func TestDownloadFileToLocalReturnsErrorForHTTPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer server.Close()

	filePath := filepath.Join(t.TempDir(), "data.csv")
	err := DownloadFileToLocal(server.URL+"/data.csv", filePath)
	if err == nil {
		t.Fatal("DownloadFileToLocal() error = nil, want non-nil")
	}

	if _, statErr := os.Stat(filePath); !os.IsNotExist(statErr) {
		t.Fatalf("DownloadFileToLocal() wrote file on failed response, statErr = %v", statErr)
	}
}

func TestDownloadFileToLocalUsesCOSDownloadForCOSDomain(t *testing.T) {
	oldCOS := config.CfgData.Storage.COS
	config.CfgData.Storage.COS.Schema = "https"
	config.CfgData.Storage.COS.CdnHost = "cdn.example.com"
	config.CfgData.Storage.COS.Bucket = "bucket"
	config.CfgData.Storage.COS.Origin = "ap-shanghai"
	defer func() {
		config.CfgData.Storage.COS = oldCOS
	}()

	oldReadCOSRemoteFile := readCOSRemoteFile
	defer func() {
		readCOSRemoteFile = oldReadCOSRemoteFile
	}()

	var gotURL string
	readCOSRemoteFile = func(ctx context.Context, fileURL string) ([]byte, error) {
		gotURL = fileURL
		return []byte("id,name\n1,cos\n"), nil
	}

	filePath := filepath.Join(t.TempDir(), "data.csv")
	err := DownloadFileToLocal("https://cdn.example.com/private/file.csv", filePath)
	if err != nil {
		t.Fatalf("DownloadFileToLocal() error = %v", err)
	}

	if gotURL != "https://cdn.example.com/private/file.csv" {
		t.Fatalf("readCOSRemoteFile URL = %q", gotURL)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "id,name\n1,cos\n" {
		t.Fatalf("DownloadFileToLocal() wrote %q", string(data))
	}
}
