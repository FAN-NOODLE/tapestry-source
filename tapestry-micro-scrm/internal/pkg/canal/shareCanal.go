package canal

import (
	"crypto/md5"
	"encoding/hex"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"os"

	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/canal"
	"micro-scrm/internal/config"
)

var ShareCanal *canal.Canal

func InitCanal() error {
	dcf := canal.NewDefaultConfig()
	dcf.DSN = config.CfgData.DB.MasterDSN

	host, _ := os.Hostname()
	m := md5.Sum([]byte(host))
	hostMd5 := hex.EncodeToString(m[:])
	dcf.LocalHost = hostMd5

	dcf.IncludeTableRegex = []string{"enterprise_user"}

	c, err := canal.NewCanal(dcf)
	if err != nil {
		log.Error(ctx.New()).Msg("canal fail to NewCanal ")
		return err
	}
	ShareCanal = c
	return nil
}
func Close() {
	if ShareCanal != nil {
		ShareCanal.Close()
	}
}
