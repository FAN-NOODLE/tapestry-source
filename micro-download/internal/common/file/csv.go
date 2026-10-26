package file

import (
	"encoding/csv"
	"os"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-download/internal/common"
	"micro-download/internal/modules/download/delivery"
)

type CsvGenerate struct{}

func (c *CsvGenerate) Generate(ctx ctx.Context, filename string, startLine int, data [][]any, overview *delivery.SourceOverview, complete bool) error {
	// filename = "1.csv"
	var newData [][]string
	var file *os.File
	var err error
	if common.FileExists(filename) {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		// 写入utf8的bom头
		_, err = file.WriteString("\xEF\xBB\xBF")
		var prefix [][]string
		for _, item := range overview.HeadPrefix {
			for _, iitem := range item {
				prefix = append(prefix, []string{iitem})
			}
		}
		var title []string
		for _, item := range overview.Title {
			title = append(title, overview.TitleAlias[item])
		}
		prefix = append(prefix, title)
		for _, item := range data {
			var info []string
			for _, a := range item {
				info = append(info, a.(string))
			}
			newData = append(newData)
		}
		newData = append(prefix, newData...)
		if err != nil {
			return err
		}
	}
	w := csv.NewWriter(file)

	for _, item := range newData {
		err = w.Write(item)
		if err != nil {
			return err
		}
	}
	if complete {
		for _, item := range overview.Suffix {
			err = w.Write(item)
			if err != nil {
				return err
			}
		}
	}
	w.Flush()
	return nil
}
