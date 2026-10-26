package scriptsjob

import (
	"database/sql"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	config2 "micro-scrm/internal/config"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/scripts"
	scriptsMysql "micro-scrm/internal/modules/scripts/respository/mysql"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/xjob"
)

var dbClient *sqlplus.ClusterDB
var bossDBClient *sqlplus.ClusterDB
var MicroScrmDbOldName string
var MicroScrmDbName string

var EnterpriseAppInfo map[int64]string

type OpenPlatformInfo struct {
	Name                      string `db:"name" desc:"小程序名称"`
	ShopAppID                 string `db:"shop_app_id" desc:"小程序appid"`
	ShopBindPlatformAppID     string `desc:"小程序绑定的虚拟开放平台appid"`
	ExternalAppID             string `db:"corp_id" desc:"企业微信corp_id"`
	ExternalBindPlatformAppID string `db:"open_app_id" desc:"企业微信绑定的开放平台id"`
	CreatedAt                 int    `db:"created_at" desc:"创建时间"`
	UpdatedAt                 int    `db:"updated_at" desc:"更新时间"`
	GzhBindStatus             int    `db:"bind_status" desc:"公众号解绑状态：绑定状态：0-未知，1-绑定，2-解绑"`
	GzhName                   string `db:"open_name" desc:"公众号名称"`
	GzhAppID                  string `db:"open_appid" desc:"公众号appid"`
	GzhBindCreatedAt          int    `db:"gzh_bind_created_at" desc:"公众号绑定创建时间"`
	GzhBindUpdatedAt          int    `db:"gzh_bind_updated_at" desc:"公众号绑定更新时间"`
}

// wx_app_config 当前企业接入小程序、公众号、企业微信等
type WxAppConfig struct {
	AppID          string `json:"appid" db:"appid" desc:"小程序appid"`
	AppType        int    `json:"app_type" db:"appid" desc:"小程序app_type，1小程序 2企业微信 3公众号"`
	OpenPlatformID string `json:"open_platform_id" db:"open_platform_id" desc:"虚拟开放平台appid"`
	BindStatus     int    `json:"bind_status" db:"bind_status" desc:"绑定状态 0 未绑定 1已绑定"`
	Name           string `json:"name" db:"name" desc:"虚拟开放平台名称name"`
	CreatedAt      int    `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt      int    `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

var BindStatusEnum = struct {
	No  int
	Yes int
}{
	No:  0,
	Yes: 1,
}

// SplitTable 用户分类表分表脚本
func SplitTable(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	// 查询所有的企业信息
	total, list, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		log.Error(c).Msgf("查询企业列表失败，错误信息：%+v", err)
	}
	if total > 0 {
		service := scripts.New(scriptsMysql.New())
		for _, enterprise := range *list {
			jobHandler.SpecDo(c, fmt.Sprintf("%d", enterprise.EnterpriseID), func() {
				service.SpiltTable(c, enterprise.EnterpriseID, enterprise.Domain)
			})
		}
	}
	return nil
}

// SplitDatabase 企业分库迁移脚本
func SplitDatabase(c ctx.Context, jobHandler *xjob.JobHandler) error {
	dbClient = db.DB()
	var err error
	bossDBClient, err = initBossDB()
	MicroScrmDbOldName = initMicroScrmDbOldName()
	MicroScrmDbName = initMicroScrmSplitDbName()
	EnterpriseAppInfo = initEnterpriseOfficial()
	enterpriseService := enterprise.New(enterpriseMysql.New())
	// 查询所有的企业信息
	total, list, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		log.Error(c).Msgf("查询企业列表失败，错误信息：%+v", err)
	}
	if total <= 0 {
		return err
	}

	service := scripts.New(scriptsMysql.New())

	for _, enterpriseEntity := range *list {
		jobHandler.SpecDo(c, fmt.Sprintf("%d", enterpriseEntity.EnterpriseID), func() {
			err = service.SpiltDatabase(c, enterpriseEntity.EnterpriseID, enterpriseEntity.Domain)
			fmt.Println(err)
			InsertWxAppConfig(c, enterpriseEntity.EnterpriseID)
		})
	}

	return nil
}

func initBossDB() (*sqlplus.ClusterDB, error) {
	config := config2.CfgData.DB
	config.MasterDSN = config.BossDSN
	return db.Init(config)
}

func initEnterpriseOfficial() map[int64]string {
	enterpriseInfo := make(map[int64]string)
	enterpriseInfo[96519191699584] = ""
	enterpriseInfo[131372809013376] = ""
	enterpriseInfo[131546657148032] = ""
	return enterpriseInfo
}

func initMicroScrmDbOldName() string {
	switch config2.CfgData.Env {
	case "test":
		return "micro_scrm"
	case "uat":
		return "micro_scrm"
	case "prod":
		return "micro_scrm"
	default:
		return ""
	}
}
func initMicroScrmSplitDbName() string {
	return "micro_scrm_enterprise_%d"
}

// 插入企业信息
func InsertWxAppConfig(c ctx.Context, enterpriseID int64) {
	config := &WxAppConfig{}
	// 初始化
	op, err := GetOpenPlatformInfo(c, enterpriseID)
	if err != nil {
		log.Error(c).Msgf("GetOpenPlatformInfo出错，err=%+v", err)
	}
	config.AppID = op.ShopAppID
	config.AppType = customer.CustomerTypeShop
	config.Name = op.Name
	config.OpenPlatformID = op.ShopBindPlatformAppID
	config.CreatedAt = op.CreatedAt
	config.UpdatedAt = op.UpdatedAt
	config.BindStatus = BindStatusEnum.Yes

	appID, err := GetWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID))
	if err != nil && err != sql.ErrNoRows {
		log.Error(c).Msgf("查询商城app_config出错，err=%+v", err)
		return
	}

	if appID != "" {
		log.Info(c).Msgf("商城app_config已存在,不再插入，appid=%+v", appID)
		return
	}

	if err = AddWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID)); err != nil {
		log.Error(c).Msgf("添加商城app_config出错，err=%+v", err)
		return
	}

	if op.ExternalAppID != "" {
		config.AppID = op.ExternalAppID
		config.AppType = customer.CustomerTypeExternal
		config.Name = op.Name
		config.OpenPlatformID = op.ExternalBindPlatformAppID
		config.CreatedAt = op.CreatedAt
		config.UpdatedAt = op.UpdatedAt
		config.BindStatus = BindStatusEnum.Yes

		appID, err := GetWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil && err != sql.ErrNoRows {
			log.Error(c).Msgf("查询企业微信app_config出错，err=%+v", err)
			return
		}

		if appID != "" {
			log.Info(c).Msgf("企业微信app_config已存在,不再插入，appid=%+v", appID)
			return
		}

		err = AddWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil {
			log.Error(c).Msgf("添加企业微信app_config出错，err=%+v", err)
		}
	}

	if op.GzhAppID != "" {
		config.AppID = op.GzhAppID
		config.AppType = customer.CustomerTypeOpen
		config.Name = op.GzhName
		config.OpenPlatformID = op.ShopBindPlatformAppID // 企业公众号绑定的虚拟开放平台appid,暂时绑定相同的开放平台
		config.CreatedAt = op.GzhBindCreatedAt
		config.UpdatedAt = op.GzhBindUpdatedAt
		config.BindStatus = op.GzhBindStatus

		appID, err := GetWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil && err != sql.ErrNoRows {
			log.Error(c).Msgf("查询公众号app_config出错，err=%+v", err)
			return
		}

		if appID != "" {
			log.Info(c).Msgf("公众号app_config已存在,不再插入，appid=%+v", appID)
			return
		}

		err = AddWxAppConfig(c, config, fmt.Sprintf(MicroScrmDbName, enterpriseID))

		if err != nil {
			log.Error(c).Msgf("添加企业公众号app_config出错，err=%+v", err)
		}
	}
}

func GetOpenPlatformInfo(c ctx.Context, enterpriseID int64) (*OpenPlatformInfo, error) {
	openPlatformInfo := &OpenPlatformInfo{}
	// 查询小程序的虚拟开放平台appid
	shopQuery := fmt.Sprintf("select open_app_id from %s where enterprise_id = ?", "enterprise_virtual_open_platform")

	externalQuery := fmt.Sprintf("select a.corp_id,b.app_id as shop_app_id,a.name,a.created_at,a.updated_at from %s a, %s b where a.id = b.enterprise_id and a.id = ?", "enterprise", "enterprise_third_part_account")
	// 查询小程序企业开放平台信息
	if err := bossDBClient.DB.Get(c, openPlatformInfo, shopQuery, enterpriseID); err != nil {
		fmt.Println(err)
		return openPlatformInfo, err
	}
	// 查询企业微信开放平台信息
	if err := bossDBClient.DB.Get(c, openPlatformInfo, externalQuery, enterpriseID); err != nil {
		fmt.Println(err)
		return openPlatformInfo, err
	}
	err := GetOfficialOpenPlatformInfo(c, enterpriseID, MicroScrmDbOldName, openPlatformInfo)
	if err != nil {
		fmt.Println(err)
	}
	openPlatformInfo.ShopBindPlatformAppID = openPlatformInfo.ExternalBindPlatformAppID
	return openPlatformInfo, nil
}

// 获取公众号的信息
func GetOfficialOpenPlatformInfo(c ctx.Context, enterpriseID int64, dbname string, openPlatformInfo *OpenPlatformInfo) error {
	openQuery := fmt.Sprintf("select name as open_name,appid as open_appid,bind_status,created_at as gzh_bind_created_at, updated_at as gzh_bind_updated_at from %s.%s where enterprise_id= ? order by created_at desc limit 1", dbname, "official_account")

	// 查询公众号开放平台信息
	if err := dbClient.DB.Get(c, openPlatformInfo, openQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return err
	}
	return nil
}

// 添加企业初始化信息
func AddWxAppConfig(c ctx.Context, wxAppConfig *WxAppConfig, dbname string) error {
	wxAppConfigQuery := fmt.Sprintf("insert into %s.%s (appid,app_type,open_platform_id,bind_status,name,created_at,updated_at) values (?,?,?,?,?,?,?)", dbname, "wx_app_config")
	_, err := dbClient.DB.Exec(c, wxAppConfigQuery, wxAppConfig.AppID, wxAppConfig.AppType, wxAppConfig.OpenPlatformID, wxAppConfig.BindStatus, wxAppConfig.Name, wxAppConfig.CreatedAt, wxAppConfig.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

// 查询企业是否已插入过
func GetWxAppConfig(c ctx.Context, wxAppConfig *WxAppConfig, dbname string) (string, error) {
	appid := ""
	wxAppConfigQuery := fmt.Sprintf("select appid from %s.%s where appid = ?", dbname, "wx_app_config")
	err := dbClient.DB.Get(c, &appid, wxAppConfigQuery, wxAppConfig.AppID)
	if err != nil {
		return appid, err
	}
	return appid, err
}
