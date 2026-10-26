package assist

import "fmt"

func FormatFileSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%.2fB", float64(n)/float64(1))
	} else if n < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(n)/float64(1024))
	} else if n < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(n)/float64(1024*1024))
	} else if n < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(n)/float64(1024*1024*1024))
	} else if n < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(n)/float64(1024*1024*1024*1024))
	} else { //if n < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(n)/float64(1024*1024*1024*1024*1024))
	}
}
