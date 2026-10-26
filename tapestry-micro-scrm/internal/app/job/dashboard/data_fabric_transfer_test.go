package dashboardJob

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDataFabricIndependentTargets(t *testing.T) {
	for _, failure := range []string{"", "missing6", "missing12", "missingBoth", "connect6", "connect12", "upload6", "upload12"} {
		t.Run(failure, func(t *testing.T) {
			root := t.TempDir()
			var notices []int
			batch := newDataFabricTransfer([]int{6, 12}, func(kind int) (*dataFabricTarget, error) {
				if failure == "missingBoth" || failure == "missing6" && kind == 6 || failure == "missing12" && kind == 12 {
					return nil, nil
				}
				if failure == "connect6" && kind == 6 || failure == "connect12" && kind == 12 {
					return nil, errors.New("connection refused")
				}
				dir := filepath.Join(root, string(rune('A'+kind)))
				if err := os.Mkdir(dir, 0700); err != nil {
					t.Fatal(err)
				}
				return &dataFabricTarget{upload: func(local string) error {
					if failure == "upload6" && kind == 6 || failure == "upload12" && kind == 12 {
						return errors.New("write failed")
					}
					data, err := os.ReadFile(local)
					if err != nil {
						return err
					}
					return os.WriteFile(filepath.Join(dir, filepath.Base(local)), data, 0600)
				}}, nil
			}, func(kind int, err error) { notices = append(notices, kind) })
			for _, name := range []string{"members_20260914.zip", "members_20260914.done", "20260914.rdy"} {
				local := filepath.Join(root, name)
				if err := os.WriteFile(local, []byte(name), 0600); err != nil {
					t.Fatal(err)
				}
				batch.upload(local)
			}
			for _, kind := range []int{6, 12} {
				want := failure != "missingBoth" && !(kind == 6 && (failure == "missing6" || failure == "connect6" || failure == "upload6")) && !(kind == 12 && (failure == "missing12" || failure == "connect12" || failure == "upload12"))
				for _, name := range []string{"members_20260914.zip", "members_20260914.done", "20260914.rdy"} {
					data, err := os.ReadFile(filepath.Join(root, string(rune('A'+kind)), name))
					if want && (err != nil || string(data) != name) {
						t.Fatalf("target %d %s: %q, %v", kind, name, data, err)
					}
					if !want && !os.IsNotExist(err) {
						t.Fatalf("unexpected file for target %d: %s", kind, name)
					}
				}
			}
			wantFailure := failure == "connect6" || failure == "connect12" || failure == "upload6" || failure == "upload12"
			if (batch.err != nil) != wantFailure {
				t.Fatalf("batch error: %v", batch.err)
			}
			if wantFailure && len(notices) == 0 {
				t.Fatal("missing failure notice")
			}
			if !wantFailure && len(notices) != 0 {
				t.Fatal("unexpected notice")
			}
		})
	}
}

func TestDataFabricCompletionMarkersRequireData(t *testing.T) {
	var uploaded []string
	batch := newDataFabricTransfer([]int{6}, func(int) (*dataFabricTarget, error) {
		return &dataFabricTarget{upload: func(local string) error { uploaded = append(uploaded, local); return nil }}, nil
	}, func(int, error) {})
	batch.upload("missing.done")
	batch.upload("20260914.rdy")
	if len(uploaded) != 0 {
		t.Fatalf("premature completion markers: %v", uploaded)
	}
}

func TestDataFabricGenerationFailurePreventsReady(t *testing.T) {
	var uploaded []string
	batch := newDataFabricTransfer([]int{6, 12}, func(int) (*dataFabricTarget, error) {
		return &dataFabricTarget{upload: func(local string) error { uploaded = append(uploaded, local); return nil }}, nil
	}, func(int, error) {})
	batch.recordGenerationError(errors.New("query failed"))
	batch.upload("20260914.rdy")
	if len(uploaded) != 0 || batch.err == nil {
		t.Fatalf("generation failure ignored: %v, %v", uploaded, batch.err)
	}
}

func TestDataFabricMarkerFailureIsIsolated(t *testing.T) {
	for _, failedFile := range []string{"members.done", "batch.rdy"} {
		t.Run(failedFile, func(t *testing.T) {
			root := t.TempDir()
			batch := newDataFabricTransfer([]int{6, 12}, func(kind int) (*dataFabricTarget, error) {
				dir := filepath.Join(root, string(rune('A'+kind)))
				if err := os.Mkdir(dir, 0700); err != nil {
					t.Fatal(err)
				}
				return &dataFabricTarget{upload: func(local string) error {
					if kind == 6 && local == failedFile {
						return errors.New("upload failed")
					}
					return os.WriteFile(filepath.Join(dir, local), []byte("complete"), 0600)
				}}, nil
			}, func(int, error) {})
			for _, name := range []string{"members.zip", "members.done", "batch.rdy"} {
				batch.upload(name)
			}
			if _, err := os.Stat(filepath.Join(root, "G", "batch.rdy")); !os.IsNotExist(err) {
				t.Fatal("failed target published ready")
			}
			if data, err := os.ReadFile(filepath.Join(root, "M", "batch.rdy")); err != nil || string(data) != "complete" {
				t.Fatalf("healthy target did not finish: %q, %v", data, err)
			}
			if batch.err == nil {
				t.Fatal("failure was not returned")
			}
		})
	}
}

func TestDataFabricRoutesDifferentTableVersionsToEachTarget(t *testing.T) {
	root := t.TempDir()
	batch := newDataFabricTransfer([]int{6, 12}, func(kind int) (*dataFabricTarget, error) {
		dir := filepath.Join(root, string(rune('A'+kind)))
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
		return &dataFabricTarget{upload: func(local string) error {
			data, err := os.ReadFile(local)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dir, filepath.Base(local)), data, 0600)
		}}, nil
	}, func(int, error) {})

	for _, version := range []struct {
		kind int
		data string
	}{{6, "incremental"}, {12, "full"}} {
		batch.selectedKind = version.kind
		for _, name := range []string{"members.zip", "members.done"} {
			local := filepath.Join(root, name)
			if err := os.WriteFile(local, []byte(version.data), 0600); err != nil {
				t.Fatal(err)
			}
			batch.upload(local)
		}
	}
	batch.selectedKind = 0
	shared := filepath.Join(root, "source.zip")
	if err := os.WriteFile(shared, []byte("shared"), 0600); err != nil {
		t.Fatal(err)
	}
	batch.upload(shared)

	for _, target := range []struct {
		kind int
		want string
	}{{6, "incremental"}, {12, "full"}} {
		dir := filepath.Join(root, string(rune('A'+target.kind)))
		for _, name := range []string{"members.zip", "members.done"} {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil || string(data) != target.want {
				t.Fatalf("target %d %s: %q, %v", target.kind, name, data, err)
			}
		}
		data, err := os.ReadFile(filepath.Join(dir, "source.zip"))
		if err != nil || string(data) != "shared" {
			t.Fatalf("shared file on target %d: %q, %v", target.kind, data, err)
		}
	}
}
