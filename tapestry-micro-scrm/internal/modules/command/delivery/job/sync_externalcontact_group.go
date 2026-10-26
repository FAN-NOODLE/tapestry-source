package job

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
)

const (
	receiveCHan int = 2000
	limit       int = 2000
	witerLimit  int = 200
)

type externalContactGroup struct {
	name string
}
type UserExternal struct {
	WorkExternalUserId int64  `json:"work_external_user_id"`
	UnionId            string `json:"union_id"`
	OpenId             string `json:"open_id"`
	Euid               int64  `json:"euid"`
	UserType           int    `json:"user_type"`
}
type userExternalBa struct {
	UserExternal

	AddWay         int    `json:"add_way"`
	RelationStatus int    `json:"relation_status"`
	MemberId       int64  `json:"member_id"`
	BaUserId       string `json:"ba_user_id"`
	BaUserName     string `json:"ba_user_name"`
	QwJoinTime     int64  `json:"qw_join_time"`
	QwLeaveTime    int64  `json:"qw_leave_time"`
}
type userExternalChat struct {
	UserExternal
	ChatId           string `json:"chat_id"`
	ChatName         string `json:"chat_name"`
	ChatOwner        string `json:"chat_owner"`
	ChatMemberStatus int    `json:"chat_member_status"`
	ChatJoinTime     int64  `json:"chat_join_time"`
	ChatLeaveTime    int64  `json:"chat_leave_time"`
	ChatJoinScene    int    `json:"chat_join_scene"`
	ChatMemberType   int    `json:"chat_member_type"`
}
type UserFootprint struct {
	Event     string `JSON:"event"`
	EventDesc string `JSON:"event_desc"`
	CreatedAt int64  `JSON:"created_at"`
}

func NewSyncExternalcontact() *externalContactGroup {
	return &externalContactGroup{
		name: fmt.Sprintf("%s::外部联系人统计JOB", config.CfgData.App),
	}
}

func (job *externalContactGroup) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) (err error) {
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIds := make([]int64, 0)
	for _, args := range specsArgs {
		enterpriseId, err := strconv.ParseInt(args, 10, 64)
		if err != nil {
			log.Info(ctx).Msgf("外部联系人统计JOB未捕获到相应的参数:%v", err)
		}
		enterpriseIds = append(enterpriseIds, enterpriseId)
	}
	sftpConfigList, err := NewSftpConfig(job.name).GetSFTPConfig(ctx, enterpriseIds)
	if err != nil {
		log.Info(ctx).Msgf("计划任务-获取品牌SFTP配置失败： %v", err)
		return err
	}
	var sync sync.WaitGroup
	for _, config := range sftpConfigList {
		sync.Add(1)
		go func(config micro_scrm.BrandConfigSftpDetail) {
			defer sync.Done()
			job.createEnterpriseIdSFTPcsv(ctx, config)
		}(config)
	}
	sync.Wait()
	return err
}
func (job *externalContactGroup) createEnterpriseIdSFTPcsv(ctx ctx.Context, config micro_scrm.BrandConfigSftpDetail) error {
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx).Msgf("创建SFTP文件panic错误： %v", err)
		}
	}()

	fileName := time.Now().Format("20060102150405")
	syncExternalUserSftp, err := orm.Service.WorkExternalUser.FindSyncExternalUserSftpLastTime(config.EnterpriseID)
	startTime := int64(0)
	if err != nil {
		startTime = int64(0)
	} else {
		startTime = syncExternalUserSftp.FileTime + 1
	}
	endTime := time.Now().Unix()
	var sync sync.WaitGroup
	sync.Add(2)
	go func() {
		defer sync.Done()
		job.createSFTPBacsv(ctx, config, fileName, startTime, endTime)
	}()
	go func() {
		defer sync.Done()
		job.createSFTPGroupChatcsv(ctx, config, fileName, startTime, endTime)
	}()
	sync.Wait()
	params := micro_scrm_enterprise.SyncExternalUserSftpModel{
		FileTime: endTime,
		FileName: config.BrandCode + "_external_user_data_" + fileName + ".csv," + config.BrandCode + "_external_user_chat_data_" + fileName + ".csv",
	}
	orm.Service.WorkExternalUser.CreateSyncExternalUserSftp(config.EnterpriseID, params)
	return nil
}
func (job *externalContactGroup) createSFTPBacsv(ctx ctx.Context, config micro_scrm.BrandConfigSftpDetail, name string, startTime, endTime int64) error {
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx).Msgf("创建SFTPCSV导购文件panic错误： %v", err)
		}
	}()
	fileName := config.BrandCode + "_external_user_data_" + name + ".csv"
	locFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	locFile.Write([]byte("\xEF\xBB\xBF"))
	w := csv.NewWriter(locFile)
	title := [][]string{{
		"external_userid", "unionid", "qw_source", "is_join_qw", "ba_userid",
		"join_time", "is_leave_qw", "leave_time", "stay_time"}}
	w.WriteAll(title)
	// 文件写入
	receiveExternalUsertBaCHan := make(chan userExternalBa, receiveCHan)
	go findExternalUserBa(receiveExternalUsertBaCHan, config.EnterpriseID, startTime, endTime)
	num := 0
	data := make([][]string, 0)
	sftpFlag := false
	for user := range receiveExternalUsertBaCHan {
		if !sftpFlag {
			sftpFlag = true
		}
		num++
		unionId := base64.StdEncoding.EncodeToString(common.AES_CBC_Encrypt([]byte(user.UnionId), []byte("rAk2PnNL17rd2utEx4zOmnP2ZfI4McSA"), []byte("7ZfLrFhjAx2iyCBh")))
		qwsource := ""
		joinQW := "是"
		if user.MemberId == 0 {
			joinQW = "否"
		}
		baUser := ""

		joinTime := ""
		leavQW := "否"
		leavTime := ""
		stopTime := ""
		if user.MemberId != 0 {
			switch user.AddWay {
			case 0:
				qwsource = "其他"
			case 1:
				qwsource = "扫描二维码"
			case 2:
				qwsource = "搜索手机号"
			case 3:
				qwsource = "名片分享"
			case 4:
				qwsource = "群聊"
			case 5:
				qwsource = "手机通讯录"
			case 6:
				qwsource = "微信联系人"
			case 7:
				qwsource = "来自微信的添加好友申请"
			case 8:
				qwsource = "安装第三方应用时自动添加的客服人员"
			case 9:
				qwsource = "搜索邮箱"
			case 201:
				qwsource = "内部成员共享"
			case 202:
				qwsource = "管理员/负责人分配"
			}
			baUser = user.BaUserId
			joinTime = helper.GetFormatDate(user.QwJoinTime, 0)
			if user.RelationStatus != 0 {
				leavQW = "是"
				leavTime = helper.GetFormatDate(user.QwLeaveTime, 0)
			}
			stopTime = helper.GetFormatDate(user.QwJoinTime, 0) + "~" + helper.GetFormatDate(user.QwLeaveTime, 0)
		}
		user := []string{
			user.OpenId,
			unionId,
			qwsource,
			joinQW,
			baUser,
			joinTime,
			leavQW,
			leavTime,
			stopTime,
		}
		data = append(data, user)
		if num == witerLimit {
			w.WriteAll(data)
			w.Flush()
			data = make([][]string, 0)
			num = 0
		}
	}
	if len(data) > 0 {
		w.WriteAll(data)
	}
	w.Flush()
	defer func() {
		locFile.Close()
		if sftpFlag {
			job.SFTPUpload(ctx, config.ConfigDetail, fileName)
		} else {
			os.Remove(fileName)
		}
	}()
	return nil
}
func (job *externalContactGroup) createSFTPGroupChatcsv(ctx ctx.Context, config micro_scrm.BrandConfigSftpDetail, name string, startTime, endTime int64) error {
	defer func() {
		if err := recover(); err != nil {
			log.Info(ctx).Msgf("创建SFTPCSV群聊文件panic错误： %v", err)
		}
	}()
	fileName := config.BrandCode + "_external_user_chat_data_" + name + ".csv"
	locFile, err := os.Create(fileName)
	locFile.Write([]byte("\xEF\xBB\xBF"))
	w := csv.NewWriter(locFile) // 创建一个新的写入文件流
	receiveExternalUserGroupChatCHan := make(chan userExternalChat, receiveCHan)
	go findExternalUserGroupChat(receiveExternalUserGroupChatCHan, config.EnterpriseID, startTime, endTime)
	title := [][]string{{
		"external_userid", "unionid", "is_join_chat", "chat_name", "chat_id",
		"chat_owner", "join_chat_time", "join_scene", "is_leave_chat", "leave_time", "stay_time"}}

	w.WriteAll(title)
	num := 0
	data := make([][]string, 0)
	sftpFlag := false
	for user := range receiveExternalUserGroupChatCHan {
		if !sftpFlag {
			sftpFlag = true
		}
		num++
		unionId := base64.StdEncoding.EncodeToString(common.AES_CBC_Encrypt([]byte(user.UnionId), []byte("rAk2PnNL17rd2utEx4zOmnP2ZfI4McSA"), []byte("7ZfLrFhjAx2iyCBh")))
		chatJoinBoole := "否"
		chatJoinTime := ""
		if user.ChatId != "" {
			chatJoinBoole = "是"
			chatJoinTime = helper.GetFormatDate(user.ChatJoinTime, 0)
		}
		joinScene := ""
		if user.ChatJoinScene == 1 {
			joinScene = "直接邀请入群"
		} else if user.ChatJoinScene == 2 {
			joinScene = "邀请链接入群"
		} else if user.ChatJoinScene == 3 {
			joinScene = "描群二维码入群"
		} else if user.ChatId != "" {
			joinScene = "其他"
		}
		chatleveBool := ""
		chatLeveTime := ""
		if user.ChatMemberStatus == 2 {
			chatleveBool = "是"
			chatLeveTime = helper.GetFormatDate(user.ChatLeaveTime, 0)
		} else if user.ChatMemberStatus == 1 {
			chatleveBool = "否"
		}
		residenceTime := ""
		if user.ChatMemberStatus == 2 {
			residenceTime = helper.GetFormatDate(user.ChatJoinTime, 0) + "~" + helper.GetFormatDate(user.ChatLeaveTime, 0)
		} else if user.ChatMemberStatus == 1 {
			residenceTime = helper.GetFormatDate(user.ChatJoinTime, 0) + "~" + helper.GetFormatDate(0, 0)
		}
		chatName := strings.Replace(user.ChatName, "\n", "", -1)
		user := []string{
			user.OpenId,
			unionId,
			chatJoinBoole,
			chatName,
			user.ChatId,
			user.ChatOwner,
			chatJoinTime,
			joinScene,
			chatleveBool,
			chatLeveTime,
			residenceTime,
		}
		data = append(data, user)
		if num == witerLimit {
			w.WriteAll(data)
			w.Flush()
			data = make([][]string, 0)
			num = 0
		}
	}
	if len(data) > 0 {
		w.WriteAll(data)
	}
	w.Flush()
	defer func() {
		locFile.Close()
		if sftpFlag {
			job.SFTPUpload(ctx, config.ConfigDetail, fileName)
		} else {
			os.Remove(fileName)
		}
	}()
	return err
}
func (job *externalContactGroup) SFTPUpload(ctx ctx.Context, config sftp.SSHConfig, fileName string) error {
	locFile, err := os.Open(fileName)
	if err != nil {
		log.Info(ctx).Msgf("%v获取失败: %v", fileName, err)
	}
	sftp, clientErr := sftp.NewSFTPClient(config)
	if clientErr != nil {
		log.Info(ctx).Msgf("%v上传失败: %v", fileName, clientErr)
		return clientErr
	}
	file, err := sftp.Client.Create("/" + fileName)
	if err != nil {
		log.Info(ctx).Msgf("%v上传失败: %v", fileName, clientErr)
		return clientErr
	}
	buf := make([]byte, 1024)
	for {
		n, err := locFile.Read(buf)
		if err != nil {
			break
		}
		_, _ = file.Write(buf[:n])
	}
	defer func() {
		file.Close()
		sftp.Client.Close()
		locFile.Close()
		os.Remove(fileName)
	}()
	return nil
}
func findExternalUserBa(receiveExternalUsertBaCHan chan userExternalBa, enterpriseID int64, startTime, endTime int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("国士无双外部联系人数据查询错误")
			close(receiveExternalUsertBaCHan)
		}
	}()
	userCorpExternal := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "
	for page := 1; ; page++ {
		params := []interface{}{startTime, (page - 1) * limit, limit}
		userCorpExternalSql := `SELECT a.id work_external_user_id,a.openid open_id,a.unionid union_id,a.euid,IF(a.updated_at>=?,1,2) user_type FROM ` + userCorpExternal + " a WHERE a.deleted_at=0 AND a.unionid!='' ORDER BY id DESC LIMIT ?,?"
		externalUserResult := []userExternalBa{}
		err := orm.Service.WorkExternalUser.Raw(userCorpExternalSql, params...).Scan(&externalUserResult).Error
		if err != nil {
			close(receiveExternalUsertBaCHan)
			return
		}
		flag := false
		workExternalUserIds := make([]int64, 0)
		for _, value := range externalUserResult {
			if !flag {
				flag = true
			}
			workExternalUserIds = append(workExternalUserIds, value.WorkExternalUserId)
		}
		if !flag {
			close(receiveExternalUsertBaCHan)
			return
		}

		memberRelationParams := []interface{}{enterpriseID, enterpriseID, startTime, endTime, workExternalUserIds}
		memberRelationSql := `
		SELECT 
		b.work_external_user_id,b.add_way,b.relation_status,b.member_id,c.userid ba_user_id,c.name ba_user_name,b.createtime qw_join_time,
		CASE WHEN b.member_deleted_time > b.external_deleted_time THEN b.member_deleted_time 
		  WHEN b.member_deleted_time < b.external_deleted_time THEN b.external_deleted_time 
		  ELSE b.updated_at END qw_leave_time FROM work_external_user_relation b 
		LEFT JOIN members c ON b.member_id=c.id AND c.enterprise_id=? AND c.deleted_at=0 
		WHERE b.enterprise_id=? AND b.deleted_at=0 AND c.id IS NOT NULL AND b.updated_at BETWEEN ? AND ?  AND b.work_external_user_id IN (?)
		ORDER BY b.work_external_user_id,b.member_id
		`
		memberRelationResult := []userExternalBa{}
		orm.Service.WorkExternalUser.Raw(memberRelationSql, memberRelationParams...).Scan(&memberRelationResult)
		for _, externalUser := range externalUserResult {
			userFlag := false
			for _, memberRelation := range memberRelationResult {
				if externalUser.WorkExternalUserId == memberRelation.WorkExternalUserId {
					userFlag = true
					externalUser.AddWay = memberRelation.AddWay
					externalUser.RelationStatus = memberRelation.RelationStatus
					externalUser.MemberId = memberRelation.MemberId
					externalUser.BaUserId = memberRelation.BaUserId
					externalUser.BaUserName = memberRelation.BaUserName
					externalUser.QwJoinTime = memberRelation.QwJoinTime
					externalUser.QwLeaveTime = memberRelation.QwLeaveTime
					footprintList, err := findBArelationLog(enterpriseID, externalUser.WorkExternalUserId, externalUser.Euid, externalUser.MemberId, externalUser.BaUserId, startTime, endTime)
					if err != nil || len(footprintList) == 0 {
						receiveExternalUsertBaCHan <- externalUser
					} else {
						joinTime := int64(0)
						leavTime := int64(0)
						relationStatus := 0
						for _, footprint := range footprintList {
							if footprint.Event == "add_external_user" {
								joinTime = footprint.CreatedAt
							}
							if footprint.Event == "delete_external_user" || footprint.Event == "delete_external_user_by_member" {
								if leavTime == 0 {
									leavTime = footprint.CreatedAt
								}
							}
							if joinTime != 0 {
								if leavTime == 0 {
									relationStatus = 0
									leavTime = helper.GetNowDateTimestamp()
								} else if leavTime > 0 {
									relationStatus = 1
								}
								externalUser.QwJoinTime = joinTime
								externalUser.QwLeaveTime = leavTime
								externalUser.RelationStatus = relationStatus

								receiveExternalUsertBaCHan <- externalUser

								joinTime = 0
								leavTime = 0
							}
						}
					}

				} else {
					if userFlag {
						break
					}
				}
			}
			if !userFlag && externalUser.UserType == 1 {
				receiveExternalUsertBaCHan <- externalUser
			}
		}
	}
}
func findBArelationLog(enterpriseID, workExternalUserId, euid, memberId int64, baUserId string, startTime, endTime int64) (result []UserFootprint, err error) {
	sql := `SELECT event,created_at  FROM user_events_v1 WHERE enterprise_id=? AND user_id=?
	 AND event IN("add_external_user","delete_external_user","delete_external_user_by_member")
	 AND event_desc LIKE ?  AND deleted_at=0 AND created_at BETWEEN ? AND ? ORDER BY created_at DESC
	 `
	params := []interface{}{enterpriseID, euid, "%" + baUserId + "%", startTime, endTime}
	err = orm.Service.WorkExternalUser.Raw(sql, params...).Scan(&result).Error
	return
}

func findExternalUserGroupChat(receiveExternalUserGroupChatCHan chan userExternalChat, enterpriseID int64, startTime, endTime int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("数据查询错误")
			close(receiveExternalUserGroupChatCHan)
		}
	}()
	userCorpExternal := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "
	for page := 1; ; page++ {
		params := []interface{}{startTime, (page - 1) * limit, limit}
		userCorpExternalSql := `SELECT a.id work_external_user_id,a.openid open_id,a.unionid union_id,a.euid,IF(a.updated_at>=?,1,2) user_type FROM ` + userCorpExternal + " a WHERE a.deleted_at=0 AND a.unionid!='' ORDER BY id DESC LIMIT ?,?"
		externalUserResult := []userExternalChat{}
		err := orm.Service.WorkExternalUser.Raw(userCorpExternalSql, params...).Scan(&externalUserResult).Error
		if err != nil {
			close(receiveExternalUserGroupChatCHan)
			return
		}
		flag := false
		unionids := make([]string, 0)
		for _, value := range externalUserResult {
			if !flag {
				flag = true
			}
			unionids = append(unionids, value.UnionId)
		}
		if !flag {
			close(receiveExternalUserGroupChatCHan)
			return
		}
		chatResult := []userExternalChat{}
		chatSql := `SELECT c.unionid union_id,d.chat_id,d.name chat_name,e.name chat_owner,
                  c.status chat_member_status,c.join_time chat_join_time,c.leave_time chat_leave_time,c.join_scene chat_join_scene,c.type chat_member_type 
		FROM work_group_chat_member c LEFT JOIN work_group_chat d ON c.group_chat_id = d.id  AND d.enterprise_id = ? AND d.deleted_at = 0
        LEFT JOIN members e ON  d.owner_member_id =e.id AND e.enterprise_id = ? AND e.deleted_at = 0 
        WHERE c.enterprise_id=? AND c.deleted_at=0 AND c.updated_at BETWEEN ? AND ?  AND c.unionid IN (?) ORDER BY c.unionid
		`
		chatParams := []interface{}{enterpriseID, enterpriseID, enterpriseID, startTime, endTime, unionids}
		orm.Service.WorkExternalUser.Raw(chatSql, chatParams...).Scan(&chatResult)

		for _, externalUser := range externalUserResult {
			userFlag := false
			for _, chat := range chatResult {
				if externalUser.UnionId == chat.UnionId {
					userFlag = true
					externalUser.ChatId = chat.ChatId
					externalUser.ChatName = chat.ChatName
					externalUser.ChatOwner = chat.ChatOwner
					externalUser.ChatMemberStatus = chat.ChatMemberStatus
					externalUser.ChatJoinTime = chat.ChatJoinTime
					externalUser.ChatLeaveTime = chat.ChatLeaveTime
					externalUser.ChatJoinScene = chat.ChatJoinScene
					externalUser.ChatMemberType = chat.ChatMemberType
					receiveExternalUserGroupChatCHan <- externalUser
				} else {
					if userFlag {
						break
					}
				}
			}
			if !userFlag && externalUser.UserType == 1 {
				receiveExternalUserGroupChatCHan <- externalUser
			}
		}
	}
}
