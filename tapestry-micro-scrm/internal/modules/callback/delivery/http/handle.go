// Package callback_http 回调
package callback_http

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	callbackSvc "micro-scrm/internal/app/service/callback"
	"micro-scrm/internal/pkg/helper"

	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/config"
	"micro-scrm/internal/events"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/modules/callback"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/drainage"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/modules/scripts"
	taghttp "micro-scrm/internal/modules/tag/delivery/http"
	userModule "micro-scrm/internal/modules/user"
	userhttp "micro-scrm/internal/modules/user/delivery/http"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk/work"
)

var service *callback.Service
var enterpriseService *enterprise.Service
var memberService *members.Service
var customerService *customer.Service
var userService *userModule.Service
var messageService *customer.MessageService
var externalTagService *customer.ExternalTagService
var drainageService *drainage.DrainageService
var contactService *drainage.ContactService
var scriptService *scripts.Service
var fissionService *fission.Service

// enterpriseAdd 新建企业
func enterpriseAdd(c *gin.Context) {
	req := &enterprise.Corp{}
	if err := c.ShouldBind(req); err != nil {
		ginplus.ResError(c, err)
		return
	}
	// 查询企业是否已经存在
	eReq := &enterprise.Corp{
		EnterpriseID: req.EnterpriseID,
	}
	detail, err := enterpriseService.FindEnterpriseFromDb(ctx.Wrap(c), *eReq)
	if err != nil && err != sql.ErrNoRows {
		ginplus.ResError(c, errs.EnterpriseNotFund)
		return
	}
	c.Set(common.ContextKeyEnum.EnterpriseID, req.EnterpriseID)

	group := threading.NewRoutineGroup()
	// 企业信息已经存在
	if detail != nil {
		_ = enterpriseService.UpdateEnterprise(ctx.Wrap(c), *req)
		// 更新缓存
		redis.Client().Del(ctx.Wrap(c), strconv.FormatInt(req.EnterpriseID, 10))
	} else {
		enterpriseService.Add(ctx.Wrap(c), req)
	}

	// 初始化标签
	group.RunSafe(func() {
		initErr := taghttp.EnterpriseInit(c)
		if initErr != nil {
			log.Error(ctx.Wrap(c)).Err(initErr).Interface("req", req).Msg("初始化企业标签信息出错")
			return
		}
	})
	// 初始化回调
	if req.CorpID != "" {
		workImpl := new(work.QyWorkImpl)
		group.RunSafe(func() {
			configUrl := fmt.Sprintf(callback.CallbackUrl, config.CfgData.Service.MicroScrm.Domain)
			list, err := workImpl.CallbackList(ctx.Wrap(c), work.SCallbackList{
				CallbackType: 1,
				CorpID:       req.CorpID,
			})
			if err != nil {
				log.Error(ctx.Wrap(c)).Err(err).Interface("req", req).Msg("查询企业外部联系人回调地址")
			}
			isConfig := false
			if len(*list) > 0 {
				for _, vItem := range *list {
					url := strings.Split(vItem.Uri, "?")
					if url[0] == configUrl {
						isConfig = true
						break
					}
				}
			}
			// 未配置，配置
			if !isConfig {
				req := work.SCallbackAdd{
					CorpID:   req.CorpID,
					CallType: 1,
					CreateBy: common.SystemAdminName,
					Uri:      configUrl,
				}
				aErr := workImpl.AddCallbackUrl(ctx.Wrap(c), req)

				if aErr != nil {
					<-time.After(1 * time.Second)
					err = workImpl.AddCallbackUrl(ctx.Wrap(c), req)
					if err != nil {
						log.Error(ctx.Wrap(c)).Err(err).Interface("req", req).Msg("新增企业外部联系人回调地址")
					}
				}
			}
		})
	}
	// 企业信息不存在时
	if detail == nil {
		// 初始化分表
		group.RunSafe(func() {
			_ = scriptService.SpiltTable(ctx.Wrap(c), req.EnterpriseID, "")
		})

		group.RunSafe(func() {
			data := enterprise.OfficialInfo{
				AppID:  config.CfgData.Mp.AppID,
				Secret: config.CfgData.Mp.Secret,
			}
			mpInfo, err := json.Marshal(data)
			if err != nil {
				return
			}

			err = enterpriseService.CreateOfficialInfo(
				ctx.Wrap(c),
				&enterprise.EpConfig{
					EnterpriseID: req.EnterpriseID,
					Hash:         req.Hash,
					Key:          "official_info",
					Value:        string(mpInfo),
					CreatedAt:    time.Now().Unix(),
					UpdatedAt:    time.Now().Unix(),
					DeletedAt:    0,
				},
			)
			if err != nil {
				log.Error(ctx.Wrap(c)).Err(err).Interface("req", req).Msg("新增企业公众号信息")
			}

		})
		// 初始化企业裂变活动兑换提醒状态
		group.RunSafe(func() {
			err := fissionService.OpenExchangeWarn(ctx.Wrap(c), fission.ExchangeWarnRequest{
				Status:       1, // 默认开启状态
				EnterpriseID: req.EnterpriseID,
				MemberIds:    make([]int64, 0),
			})
			if nil != err {
				log.Error(ctx.Wrap(c)).Err(err).Interface("req", req).Msg("初始化企业裂变活动兑换提醒状态失败")
				return
			}

		})
	}
	group.Wait()

	// 给企业分库
	_ = scriptService.SpiltDatabase(ctx.Wrap(c), req.EnterpriseID, req.Domain)
	// 增加企业的配置
	InsertWxAppConfig(ctx.Wrap(c), req.EnterpriseID)
	ginplus.ResSuccess(c, nil)
}

// tagInit 标签初始化结果
func tagInit(c *gin.Context) {
	reqCallback := &enterprise.Corp{}
	rErr := c.BindJSON(reqCallback)
	ctxContext := ctx.Wrap(c)
	if rErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	c.Set(common.ContextKeyEnum.EnterpriseID, reqCallback.EnterpriseID)
	taghttp.TagSync(c)
	_ = enterpriseService.UpdateEnterprise(ctxContext, *reqCallback)
	ginplus.ResSuccess(c, nil)
}

// externalCallback 外部联系人回调
func externalCallback(c *gin.Context) {
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	log.Info(ctx.Wrap(c)).Msgf("新增联系人验证回调,外部联系人回调信息：%s", string(reqBody))
	// 先获取企业微信的回调参数
	qwCallbackMsg := &callback.QwExternalCallback{}
	// 先取Query数据
	queryErr := c.BindQuery(qwCallbackMsg)
	if queryErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,从企业微信获取到的信息获取Query参数失败，错误信息：%+v", queryErr)
		ginplus.ResError(c, errs.CallbackQueryErr)
		return
	}
	qwErr := xml.Unmarshal(reqBody, qwCallbackMsg)
	if qwErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,从企业微信获取到的回调参数：%+v, 错误信息：%+v", qwCallbackMsg, qwErr)
		ginplus.ResError(c, errs.CallbackXmlErr)
		return
	}

	// 解密回调参数
	decryptReq := work.QyDecryptRequest{
		MsgSignature:  qwCallbackMsg.MsgSignature,
		Nonce:         qwCallbackMsg.Nonce,
		Timestamp:     qwCallbackMsg.Timestamp,
		CorpID:        qwCallbackMsg.ToUserName,
		EncryptString: string(reqBody),
	}

	workImpl := new(work.QyWorkImpl)

	callbackData, err := workImpl.QyDecrypt(ctx.Wrap(c), decryptReq)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,解密企业微信回调信息报错:%+v", err)
		ginplus.ResError(c, err)
		return
	}

	// var callbackData *work.MemberExternalCallback
	// gormx.GetInstance().Model(micro_scrm.RelationCallbackRecordsModel{}).Where(map[string]interface{}{
	//	"msg_type": "event",
	//	"event": "change_external_contact",
	//	"change_type": "add_external_contact",
	//	//"user_id": "",
	//	//"external_user_id": "",
	//	"state":"3356423205106816",
	// }).First(&callbackData)

	log.Info(ctx.Wrap(c)).Msgf("新增联系人验证回调,外部联系人解密数据：%+v", callbackData)
	service.MemberService = memberService
	service.UserService = userService
	service.CustomerService = customerService
	service.MessageService = messageService
	service.ExternalTagService = externalTagService
	service.DrainageService = drainageService
	service.ContactService = contactService
	service.FissionService = fissionService

	// 外部联系人回调
	corpReq := enterprise.Corp{
		CorpID: callbackData.ToUserName,
	}
	// 查询企业信息
	enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	if enterpriseErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,查询企业信息失败，错误信息：%+v", enterpriseErr)
		ginplus.ResError(c, enterpriseErr)
		return
	}
	workImpl.SetDomain(enterpriseData.Domain)
	var eErr error

	// 加分布式锁保证唯一性
	key := helper.Md5V(callbackData.ToUserName +
		callbackData.Event +
		// callbackData.ChangeType +
		callbackData.UserID +
		callbackData.ExternalUserID +
		callbackData.ChatID,
	)

	mtx, err := lock.Mutex(
		ctx.Wrap(c),
		lock.Key(key),
		//lock.TTL(5*time.Second),
		//lock.Retry(lock.NoRetry()),
		lock.TTL(10*time.Second),
		lock.Retry(lock.LimitRetry(lock.LinearBackoff(100*time.Millisecond), 110)),
	)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,外部联系人数据生产锁失败,错误信息：%+v", err)
		return
	}

	err = mtx.Lock()
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,外部联系人数据加锁失败,错误信息：%+v", err)
		return
	}

	channelID, err := customerService.AddCustomerCallbackRecord(ctx.Wrap(c), callbackData)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("记录回调消息记录出错, err:%v", err)
		return
	}

	_ = mtx.Unlock()

	group := threading.NewRoutineGroup()
	switch callbackData.Event {

	case callback.ExternalCallbackEventOfChangeExternalContact:
		// 外部联系人通讯录变更事件
		// 协程去处理同步日志到kafka
		group.RunSafe(func() {
			log.Debug(ctx.Wrap(c)).Msg("新增联系人验证回调,处理外部联系人回调信息到kafka")
			handleKafka := callback.NewSyncToKafka(service)
			err = callback.HandleCallbackMessage(ctx.Wrap(c), callbackData, handleKafka)
			if err != nil {
				log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,处理入 kafka 消息失败，%s", err)
			}
		})

		// 转发到新服务处理回调
		threading.GoSafe(func() error {
			if err = callbackSvc.NewExternalContact().Handel(ctx.Wrap(c), callbackData, workImpl); err != nil {
				log.Error(ctx.Wrap(c)).Msgf("处理外部联系人回调异常：%v", err)
			}
			return err
		})

		// 外部联系人通讯录变更事件处理
		eErr = service.ChangeExternalContact(ctx.Wrap(c), callbackData, enterpriseData, workImpl, channelID)

		//ok, runErr := service.IsRunFission(ctx.Wrap(c), callbackData, enterpriseData.EnterpriseID)
		//if runErr != nil {
		//	log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,判断是否能裂变活动失败%s", runErr)
		//}
		//if ok {
		//	group.RunSafe(func() {
		//		// 添加好友回调裂变活动处理 (删除kafka队列处理)
		//		err = service.AddFriendCallbackActivityHandle(ctx.Wrap(c), callbackData)
		//		if err != nil {
		//			log.Error(ctx.Wrap(c)).Msgf("新增联系人验证回调,添加好友回调裂变活动处理失败%s", err)
		//		}
		//	})
		//}
		break
	case callback.ExternalCallbackEventOfChangeExternalChat:
		// 外部联系人群变更通知
		eErr = service.ChangeExternalChat(ctx.Wrap(c), callbackData, enterpriseData, workImpl)
		break
	}
	group.Wait()

	// 事件回调，不把所有的业务逻辑都堆到这里
	events.ScrmEvent.ExternalUserCallback.Dispatch(ctx.Wrap(c), entity.ExternalUserCallbackEntity{
		EnterpriseID:           enterpriseData.EnterpriseID,
		MemberExternalCallback: *callbackData,
	})

	if eErr != nil && eErr != sql.ErrNoRows {
		log.Error(ctx.Wrap(c)).Err(err).Interface("callback", callbackData).Msg("ExternalUserCallback.Dispatch")
	}

	ginplus.ResOK(c)
}

// externalManual 手动处理外部联系人回调接口,删除外部联系人和删员工需要手动指定渠道ID
func externalManual(c *gin.Context) {
	callbackData := &work.MemberExternalCallback{}
	rErr := c.ShouldBindJSON(callbackData)
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	// 外部联系人回调
	corpReq := enterprise.Corp{
		CorpID: callbackData.ToUserName,
	}
	// 查询企业信息
	enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), corpReq)
	if enterpriseErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("查询企业信息失败，错误信息：%+v", enterpriseErr)
		ginplus.ResError(c, enterpriseErr)
		return
	}
	workImpl := (&work.QyWorkImpl{}).SetDomain(enterpriseData.Domain)
	service.MemberService = memberService
	service.UserService = userService
	service.CustomerService = customerService
	service.MessageService = messageService
	service.ExternalTagService = externalTagService
	service.DrainageService = drainageService
	service.ContactService = contactService
	service.FissionService = fissionService
	var eErr error
	switch callbackData.Event {
	case callback.ExternalCallbackEventOfChangeExternalContact:
		eErr = service.ChangeExternalContact(ctx.Wrap(c), callbackData, enterpriseData, workImpl, callbackData.ChannelID)
		break
	case callback.ExternalCallbackEventOfChangeExternalChat:
		eErr = service.ChangeExternalChat(ctx.Wrap(c), callbackData, enterpriseData, workImpl)
		break
	}
	if eErr != nil {
		log.Info(ctx.Wrap(c)).Msgf("处理外部联系人出错，错误信息：%+v, 回调信息：%+v", eErr, callbackData)
	}
	ginplus.ResOK(c)
}

// userAddressSync 同步商城用户收货地址数据
func userAddressSync(c *gin.Context) {
	ginplus.ResOK(c)
}

// userEvent 用户事件同步
func userEvent(c *gin.Context) {
	userhttp.AddEvent(c)
}

func enterpriseOfficialAccount(c *gin.Context) {
	var request callback.EnterpriseOfficialRequest
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	officialAccount := enterprise.WxAppConfigEntity{
		EnterpriseID: request.EnterpriseID,
		AppID:        request.AppID,
		AppType:      customer.CustomerTypeOpen,
		Name:         request.Name,
		BindStatus:   request.Type,
	}

	if err = enterpriseService.AddWxAppConfig(ctx.Wrap(c), &officialAccount); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

var dbClient *sqlplus.ClusterDB
var bossDBClient *sqlplus.ClusterDB
var MicroScrmDbOldName string
var MicroScrmDbName string

type OpenPlatformInfo struct {
	Name                      string         `db:"name" desc:"小程序名称"`
	ShopAppID                 sql.NullString `db:"shop_app_id" desc:"小程序appid"`
	ShopBindPlatformAppID     string         `desc:"小程序绑定的虚拟开放平台appid"`
	ExternalAppID             string         `db:"corp_id" desc:"企业微信corp_id"`
	ExternalBindPlatformAppID string         `db:"open_app_id" desc:"企业微信绑定的开放平台id"`
	CreatedAt                 int            `db:"created_at" desc:"创建时间"`
	UpdatedAt                 int            `db:"updated_at" desc:"更新时间"`
	GzhBindStatus             int            `db:"bind_status" desc:"公众号解绑状态：绑定状态：0-未知，1-绑定，2-解绑"`
	GzhName                   string         `db:"open_name" desc:"公众号名称"`
	GzhAppID                  string         `db:"open_appid" desc:"公众号appid"`
	GzhBindCreatedAt          int            `db:"gzh_bind_created_at" desc:"公众号绑定创建时间"`
	GzhBindUpdatedAt          int            `db:"gzh_bind_updated_at" desc:"公众号绑定更新时间"`
}

// WxAppConfig wx_app_config 当前企业接入小程序、公众号、企业微信等
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

// 插入企业信息
func InsertWxAppConfig(c ctx.Context, enterpriseID int64) {
	dbClient = db.DB()
	var err error
	bossDBClient, err = initBossDB()
	MicroScrmDbOldName = "micro_scrm"
	MicroScrmDbName = initMicroScrmSplitDbName()

	config2 := &WxAppConfig{}
	// 初始化
	op, err := GetOpenPlatformInfo(c, enterpriseID)
	fmt.Println(op)
	if err != nil && err != sql.ErrNoRows {
		log.Error(c).Msgf("GetOpenPlatformInfo出错，err=%+v", err)
	}
	if op == nil {
		return
	}
	config2.AppID = op.ShopAppID.String
	config2.AppType = customer.CustomerTypeShop
	config2.Name = op.Name
	config2.OpenPlatformID = op.ShopBindPlatformAppID
	config2.CreatedAt = op.CreatedAt
	config2.UpdatedAt = op.UpdatedAt
	config2.BindStatus = BindStatusEnum.Yes

	appID, err := GetWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID))
	if err != nil && err != sql.ErrNoRows {
		log.Error(c).Msgf("查询商城app_config出错，err=%+v", err)
		return
	}

	if appID != "" {
		if err = UpdateWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID)); err != nil {
			log.Error(c).Msgf("更新商城app_config出错，err=%+v", err)
			return
		}
		return
	}
	if op.ShopAppID.String != "" {
		if err = AddWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID)); err != nil {
			log.Error(c).Msgf("添加商城app_config出错，err=%+v", err)
			return
		}
	}

	if op.ExternalAppID != "" {
		config2.AppID = op.ExternalAppID
		config2.AppType = customer.CustomerTypeExternal
		config2.Name = op.Name
		config2.OpenPlatformID = op.ExternalBindPlatformAppID
		config2.CreatedAt = op.CreatedAt
		config2.UpdatedAt = op.UpdatedAt
		config2.BindStatus = BindStatusEnum.Yes

		appID, err := GetWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil && err != sql.ErrNoRows {
			log.Error(c).Msgf("查询企业微信app_config出错，err=%+v", err)
			return
		}

		if appID != "" {
			if err = UpdateWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID)); err != nil {
				log.Error(c).Msgf("更新商城app_config出错，err=%+v", err)
				return
			}
			return
		}

		err = AddWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil {
			log.Error(c).Msgf("添加企业微信app_config出错，err=%+v", err)
		}
	}

	if op.GzhAppID != "" {
		config2.AppID = op.GzhAppID
		config2.AppType = customer.CustomerTypeOpen
		config2.Name = op.GzhName
		config2.OpenPlatformID = op.ShopBindPlatformAppID // 企业公众号绑定的虚拟开放平台appid,暂时绑定相同的开放平台
		config2.CreatedAt = op.GzhBindCreatedAt
		config2.UpdatedAt = op.GzhBindUpdatedAt
		config2.BindStatus = op.GzhBindStatus

		appID, err := GetWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID))
		if err != nil && err != sql.ErrNoRows {
			log.Error(c).Msgf("查询公众号app_config出错，err=%+v", err)
			return
		}

		if appID != "" {
			if err = UpdateWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID)); err != nil {
				log.Error(c).Msgf("更新商城app_config出错，err=%+v", err)
				return
			}
			return
		}

		err = AddWxAppConfig(c, config2, fmt.Sprintf(MicroScrmDbName, enterpriseID))

		if err != nil {
			log.Error(c).Msgf("添加企业公众号app_config出错，err=%+v", err)
		}
	}
}

func initBossDB() (*sqlplus.ClusterDB, error) {
	config2 := config.CfgData.DB
	config2.MasterDSN = config.CfgData.DB.BossDSN
	return db.Init(config2)
}

func initMicroScrmSplitDbName() string {
	return "micro_scrm_enterprise_%d"
}

func GetOpenPlatformInfo(c ctx.Context, enterpriseID int64) (*OpenPlatformInfo, error) {
	openPlatformInfo := &OpenPlatformInfo{}
	openPlatformInfo.ExternalBindPlatformAppID = "default"
	openPlatformInfo.Name = "unknown"
	openPlatformInfo.GzhName = "unknown"
	// 查询小程序的虚拟开放平台appid
	shopQuery := fmt.Sprintf("select open_app_id from %s where enterprise_id = ?", "enterprise_virtual_open_platform")

	externalQuery := fmt.Sprintf("select a.corp_id,b.app_id as shop_app_id,a.name,a.created_at,a.updated_at from %s a left join %s b on a.id = b.enterprise_id where a.id = ?  ORDER BY b.created_at DESC", "enterprise", "enterprise_third_part_account")
	// 查询小程序企业开放平台信息
	if err := bossDBClient.DB.Get(c, openPlatformInfo, shopQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		return openPlatformInfo, err
	}
	// 查询企业微信开放平台信息
	if err := bossDBClient.DB.Get(c, openPlatformInfo, externalQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		return openPlatformInfo, err
	}

	err := GetOfficialOpenPlatformInfo(c, enterpriseID, MicroScrmDbOldName, openPlatformInfo)
	if err != nil {
		return openPlatformInfo, err
	}
	openPlatformInfo.ShopBindPlatformAppID = openPlatformInfo.ExternalBindPlatformAppID
	return openPlatformInfo, nil
}

// GetOfficialOpenPlatformInfo 获取公众号的信息
func GetOfficialOpenPlatformInfo(c ctx.Context, enterpriseID int64, dbname string, openPlatformInfo *OpenPlatformInfo) error {
	openQuery := fmt.Sprintf("select name as open_name,appid as open_appid,bind_status,created_at as gzh_bind_created_at, updated_at as gzh_bind_updated_at from %s.%s where enterprise_id= ? order by created_at desc limit 1", dbname, "official_account")

	// 查询公众号开放平台信息
	if err := dbClient.DB.Get(c, openPlatformInfo, openQuery, enterpriseID); err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return err
	}
	return nil
}

// AddWxAppConfig 添加企业初始化信息
func AddWxAppConfig(c ctx.Context, wxAppConfig *WxAppConfig, dbname string) error {
	wxAppConfigQuery := fmt.Sprintf("insert into %s.%s (appid,app_type,open_platform_id,bind_status,name,created_at,updated_at) values (?,?,?,?,?,?,?)", dbname, "wx_app_config")
	_, err := dbClient.DB.Exec(c, wxAppConfigQuery, wxAppConfig.AppID, wxAppConfig.AppType, wxAppConfig.OpenPlatformID, wxAppConfig.BindStatus, wxAppConfig.Name, wxAppConfig.CreatedAt, wxAppConfig.UpdatedAt)
	if err != nil {
		return err
	}
	return err
}

// UpdateWxAppConfig 更新企业初始化信息
func UpdateWxAppConfig(c ctx.Context, wxAppConfig *WxAppConfig, dbname string) error {
	wxAppConfigQuery := fmt.Sprintf("update %s.%s set app_type=?,open_platform_id=?,bind_status=?,name=?,updated_at=? where appid = ?", dbname, "wx_app_config")
	_, err := dbClient.DB.Exec(c, wxAppConfigQuery, wxAppConfig.AppType, wxAppConfig.OpenPlatformID, wxAppConfig.BindStatus, wxAppConfig.Name, wxAppConfig.UpdatedAt, wxAppConfig.AppID)
	if err != nil {
		return err
	}
	return err
}

// GetWxAppConfig 查询企业是否已插入过
func GetWxAppConfig(c ctx.Context, wxAppConfig *WxAppConfig, dbname string) (string, error) {
	appid := ""
	wxAppConfigQuery := fmt.Sprintf("select appid from %s.%s where appid = ?", dbname, "wx_app_config")
	err := dbClient.DB.Get(c, &appid, wxAppConfigQuery, wxAppConfig.AppID)
	if err != nil {
		return appid, err
	}
	return appid, err
}
