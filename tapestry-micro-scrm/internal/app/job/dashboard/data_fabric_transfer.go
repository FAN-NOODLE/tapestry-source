package dashboardJob

import (
	"fmt"
	"path/filepath"
	"strings"
)

// dataFabricTarget keeps upload and completion state independent for each SFTP.
type dataFabricTarget struct {
	kind     int
	brand    string
	upload   func(string) error
	close    func() error
	uploaded map[string]bool
	failed   bool
}

type dataFabricTransfer struct {
	targets          []*dataFabricTarget
	notice           func(int, error)
	err              error
	generationFailed bool
	selectedKind     int // 0 uploads to all targets; otherwise only to this SFTP type.
}

// A nil target with no error means that this SFTP is not configured.
func newDataFabricTransfer(kinds []int, connect func(int) (*dataFabricTarget, error), notice func(int, error)) *dataFabricTransfer {
	batch := &dataFabricTransfer{notice: notice}
	for _, kind := range kinds {
		target, err := connect(kind)
		if err != nil {
			batch.fail(kind, err)
			continue
		}
		if target == nil {
			continue
		}
		target.kind = kind
		target.uploaded = make(map[string]bool)
		batch.targets = append(batch.targets, target)
	}
	return batch
}

func (batch *dataFabricTransfer) fail(kind int, err error) {
	batch.err = fmt.Errorf("datafabric sftp_type=%d: %w", kind, err)
	batch.notice(kind, err)
}

func (batch *dataFabricTransfer) recordGenerationError(err error) {
	if err != nil {
		batch.generationFailed = true
		batch.err = err
	}
}

func (batch *dataFabricTransfer) upload(local string) {
	name := filepath.Base(local)
	ext := filepath.Ext(name)
	for _, target := range batch.targets {
		if batch.selectedKind != 0 && target.kind != batch.selectedKind {
			continue
		}
		if target.failed || ext == ".rdy" && batch.generationFailed {
			continue
		}
		if ext == ".done" && !target.uploaded[strings.TrimSuffix(name, ext)+".zip"] {
			target.failed = true
			batch.fail(target.kind, fmt.Errorf("%s 对应的数据文件未上传成功", name))
			continue
		}
		if err := target.upload(local); err != nil {
			target.failed = true
			batch.fail(target.kind, fmt.Errorf("上传 %s 失败: %w", name, err))
			continue
		}
		target.uploaded[name] = true
	}
}

func (batch *dataFabricTransfer) close() {
	for _, target := range batch.targets {
		if target.close != nil {
			_ = target.close()
		}
	}
}
