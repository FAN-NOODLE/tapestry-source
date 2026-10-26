package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/go/cy-sdk-go/exception"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/user"
)

// WxUser 微信用户
type WxUser interface {
	Exists() bool
	NeedMerge() bool
	NeedInsertUnion() bool
	Insert(ctx.Context, *sqlplus.ClusterDB) error          // 只有openid,且在系统中没有用户
	InsertUnion(ctx.Context, *sqlplus.ClusterDB) error     // 已有openid，需要插入unionid
	Update(ctx.Context, *sqlplus.ClusterDB) error          // 只有openid，更新即可
	UpdateWithUnion(ctx.Context, *sqlplus.ClusterDB) error // unionEuid=UserEuid 和 unionEuid！=UserEuid两种情况
	SetSource(ctx.Context, *sqlplus.ClusterDB) error
}

type wxUserBase struct {
	source    *user.SourceUser
	dbname    string
	table     string
	unionEuid int64
}

func (u *wxUserBase) Exists() bool {
	return u.source.ID != 0
}

func (u *wxUserBase) NeedMerge() bool {
	return u.unionEuid != u.source.EuID
}

func (u *wxUserBase) NeedInsertUnion() bool {
	return u.unionEuid == 0 && u.source.EuID != 0 && u.source.UnionID != ""
}

func (u *wxUserBase) SetSource(c ctx.Context, db *sqlplus.ClusterDB) error {
	u.unionEuid = 0
	u.source.ID = 0
	u.source.EuID = 0

	if u.source.CreateAt == 0 {
		u.source.CreateAt = time.Now().Unix()
	}
	u.dbname = fmt.Sprintf("micro_scrm_enterprise_%d", u.source.EnterpriseID)
	if u.source.UnionID != "" {
		if err := db.DB.Get(c, &u.unionEuid, fmt.Sprintf("select euid from %s.wx_user_union where unionid=?", u.dbname), u.source.UnionID); err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	var query string
	if u.source.CustomerType == customer.CustomerTypeExternal && u.source.Name == "" {
		query = fmt.Sprintf("select id, euid, name from %s.%s where appid=? and openid=?", u.dbname, u.table)
	} else {
		query = fmt.Sprintf("select id, euid from %s.%s where appid=? and openid=?", u.dbname, u.table)
	}
	if err := db.DB.Get(c, u.source, query, u.source.AppID, u.source.OpenID); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil

}

func (u *wxUserBase) updateBase(c ctx.Context, db *sqlplus.ClusterDB) error {
	var cols []string
	var args []interface{}
	if u.source.Name != "" {
		cols = append(cols, "name=?")
		args = append(args, u.source.Name)
	}
	if u.source.Mobile != "" {
		cols = append(cols, "mobile=?")
		args = append(args, u.source.Mobile)
	}

	if u.source.DeliveryMobile != "" {
		cols = append(cols, "delivery_mobile=?")
		args = append(args, u.source.DeliveryMobile)
	}
	if u.source.Nick != "" || u.source.Avatar != "" || u.source.Sex != 0 {
		cols = append(cols, "nick=?", "avatar=?", "sex=?")
		args = append(args, u.source.Nick, u.source.Avatar, u.source.Sex)
	}

	if u.source.Country != "" || u.source.Province != "" || u.source.City != "" {
		cols = append(cols, "country=?", "province=?", "city=?")
		args = append(args, u.source.Country, u.source.Province, u.source.City)
	}

	args = append(args, u.source.EuID)

	if len(cols) > 0 {
		_, err := db.DB.Exec(c, fmt.Sprintf("update %s.enterprise_user set %s where id=?", u.dbname, strings.Join(cols, ",")), args...)
		return err
	}
	return nil
}

func (u *wxUserBase) updateDateBaseWithUnion(c ctx.Context, db *sqlplus.ClusterDB, cols1 []string, args1 []interface{}) error {
	return db.Tx(c, func(tx *sqlplus.Tx) error {
		openPlatformID := "default" //没有开放平台时的默认值
		query := fmt.Sprintf("select open_platform_id from %s.wx_app_config where appid=? and app_type = ?", u.dbname)
		if err := tx.Get(c, &openPlatformID, query, u.source.AppID, u.source.CustomerType); err != nil {
			exception.SendWarn(fmt.Errorf("%v %w", u.source, err))
		}
		unionUser := map[string]string{openPlatformID: u.source.UnionID}
		unionb, _ := json.Marshal(unionUser)

		var cols []string
		var args []interface{}
		if u.source.Name != "" {
			cols = append(cols, "name=?")
			args = append(args, u.source.Name)
		}
		if u.source.Mobile != "" {
			cols = append(cols, "mobile=?")
			args = append(args, u.source.Mobile)
		}
		if u.source.Nick != "" || u.source.Avatar != "" || u.source.Sex != 0 {
			cols = append(cols, "nick=?", "avatar=?", "sex=?")
			args = append(args, u.source.Nick, u.source.Avatar, u.source.Sex)
		}

		if u.source.Country != "" || u.source.Province != "" || u.source.City != "" {
			cols = append(cols, "country=?", "province=?", "city=?")
			args = append(args, u.source.Country, u.source.Province, u.source.City)
		}

		cols = append(cols, "wx_union=?")
		args = append(args, unionb)

		args = append(args, u.source.EuID)

		if len(cols) > 0 {
			_, err := tx.Exec(c, fmt.Sprintf("update %s.enterprise_user set %s where id=?", u.dbname, strings.Join(cols, ",")), args...)
			if err != nil {
				return err
			}
		}

		wxUnionQuery := fmt.Sprintf("insert into %s.wx_user_union (unionid,open_platform_id,open_platform_source,euid) values (?,?,?,?)", u.dbname)
		if _, err := tx.Exec(c, wxUnionQuery, u.source.UnionID, openPlatformID, 1, u.source.EuID); err != nil {
			return err
		}

		wxUserOpenQuery := fmt.Sprintf("UPDATE %s.%s SET %s WHERE id = ?", u.dbname, u.table, strings.Join(cols1, ","))
		if _, err := tx.Exec(c, wxUserOpenQuery, args1...); err != nil {
			return err
		}

		u.source.UnionUID = u.source.EuID
		return nil

	})
}

func (u *wxUserBase) insertEnterpriseUser(c ctx.Context, tx *sqlplus.Tx, appUserKey string) (int64, int64, error) {
	unionUser := make(map[string]string)
	euid, appUserID := u.source.UnionUID, u.source.OpenUID
	if euid == 0 {
		euid = uid.ID(c)
	}
	if appUserID == 0 {
		appUserID = uid.ID(c)
	}
	if u.source.UnionID != "" {
		openPlatformID := "default" //没有开放平台时的默认值
		query := fmt.Sprintf("select open_platform_id from %s.wx_app_config where appid=? and app_type = ?", u.dbname)
		if err := tx.Get(c, &openPlatformID, query, u.source.AppID, u.source.CustomerType); err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		}
		unionUser[openPlatformID] = u.source.UnionID

		//给wx_user_union表中插入数据
		wxUnionQuery := fmt.Sprintf("insert into %s.wx_user_union (unionid,open_platform_id,open_platform_source,euid) values (?,?,?,?)", u.dbname)

		if _, err := tx.Exec(c, wxUnionQuery, u.source.UnionID, openPlatformID, 1, euid); err != nil {
			return 0, 0, err
		}

	}

	appUser := make(map[string]*customer.ThirdUserID)
	appUser[u.source.AppID] = &customer.ThirdUserID{ID: appUserID, ThirdID: u.source.OpenID}
	unionb, _ := json.Marshal(unionUser)
	b, _ := json.Marshal(appUser)
	query := fmt.Sprintf("insert into %s.enterprise_user (id, name, nick, avatar, mobile, birthday, sex, language, country, province, city, remark, %s, wx_union,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", u.dbname, appUserKey)
	_, err := tx.Exec(c, query, euid, u.source.Name, u.source.Nick, u.source.Avatar, u.source.Mobile, u.source.Birthday, u.source.Sex, u.source.Language, u.source.Country, u.source.Province, u.source.City, u.source.Remark, string(b), string(unionb), u.source.CreateAt, time.Now().Unix())
	return euid, appUserID, err
}

func (u *wxUserBase) updateEnterpriseUser(c ctx.Context, tx *sqlplus.Tx, appUserKey string) (int64, int64, error) {
	appUserID := u.source.OpenUID
	if appUserID == 0 {
		appUserID = uid.ID(c)
	}
	appUser := make(map[string]*customer.ThirdUserID)
	appUser[u.source.AppID] = &customer.ThirdUserID{ID: appUserID, ThirdID: u.source.OpenID}
	b, _ := json.Marshal(appUser)
	var cols []string
	var args []interface{}
	cols = append(cols, fmt.Sprintf("%s=?", appUserKey))
	args = append(args, string(b))
	if u.source.Name != "" {
		cols = append(cols, "name=?")
		args = append(args, u.source.Name)
	}

	if u.source.Mobile != "" {
		cols = append(cols, "mobile=?")
		args = append(args, u.source.Mobile)
	}

	if u.source.DeliveryMobile != "" {
		cols = append(cols, "delivery_mobile=?")
		args = append(args, u.source.DeliveryMobile)
	}

	if u.source.Birthday != "" {
		cols = append(cols, "birthday=?")
		args = append(args, u.source.Birthday)
	}

	if u.source.Language != "" {
		cols = append(cols, "language=?")
		args = append(args, u.source.Language)
	}

	if u.source.Remark != "" {
		cols = append(cols, "remark=?")
		args = append(args, u.source.Remark)
	}

	if u.source.Nick != "" || u.source.Avatar != "" {
		cols = append(cols, "nick=?")
		args = append(args, u.source.Nick)
		cols = append(cols, "avatar=?")
		args = append(args, u.source.Avatar)
		cols = append(cols, "sex=?")
		args = append(args, u.source.Sex)
	}

	if u.source.Country != "" || u.source.Province != "" || u.source.City != "" {
		cols = append(cols, "country=?")
		args = append(args, u.source.Country)
		cols = append(cols, "province=?")
		args = append(args, u.source.Province)
		cols = append(cols, "city=?")
		args = append(args, u.source.City)
	}

	args = append(args, u.unionEuid)

	query := fmt.Sprintf("update %s.enterprise_user set %s where id=?", u.dbname, strings.Join(cols, ","))
	_, err := tx.Exec(c, query, args...)
	return u.unionEuid, appUserID, err
}

func (u *wxUserBase) insertOrUpdateEnterpriseUser(c ctx.Context, tx *sqlplus.Tx, appUserKey string) (int64, int64, error) {
	if u.unionEuid != 0 {
		return u.updateEnterpriseUser(c, tx, appUserKey)
	}
	return u.insertEnterpriseUser(c, tx, appUserKey)
}

func (u *wxUserBase) mergeEnterpriseUser(c ctx.Context, db *sqlplus.ClusterDB, unionEuid, userEuid int64) (int64, error) {
	var users []*struct {
		customer.EnterpriseUserBindDB
		Version int   `db:"version"`
		LevelID int64 `db:"level_id"`
	}
	if err := db.Select(c, &users,
		fmt.Sprintf("select id, wx_open, wx_shop, wx_external, wx_union, bind_to, version, level_id from %s.enterprise_user where id in (?, ?)", u.dbname),
		unionEuid, userEuid); err != nil {
		return 0, err
	}

	// 此处为异常情况，数据丢失
	if len(users) != 2 {
		return 0, nil
	}

	u1, u2 := users[0], users[1]
	if users[1].ID == int64(unionEuid) {
		u1, u2 = u2, u1
	}

	if u2.BindTo != 0 {
		if u2.BindTo == u1.ID || u1.BindTo != 0 {
			return 0, nil
		}
		return u.mergeEnterpriseUser(c, db, unionEuid, u1.BindTo)
	}

	uu1, uu2 := u1.ToUser(), u2.ToUser()
	if uu1.WxExternal == nil {
		uu1.WxExternal = make(map[string]*customer.ThirdUserID)
	}
	for k, v := range uu2.WxExternal {
		uu1.WxExternal[k] = v
	}

	if uu1.WxOpen == nil {
		uu1.WxOpen = make(map[string]*customer.ThirdUserID)
	}
	for k, v := range uu2.WxOpen {
		uu1.WxOpen[k] = v
	}

	if uu1.WxShop == nil {
		uu1.WxShop = make(map[string]*customer.ThirdUserID)
	}
	for k, v := range uu2.WxShop {
		uu1.WxShop[k] = v
	}

	if uu1.WxUnion == nil {
		uu1.WxUnion = make(map[string]string)
	}
	for k, v := range uu2.WxUnion {
		uu1.WxUnion[k] = v
	}

	err := db.Tx(c, func(tx *sqlplus.Tx) error {
		if _, err := tx.Exec(c, fmt.Sprintf("update %s.%s set euid=? where id=?", u.dbname, u.table), uu1.ID, u.source.ID); err != nil {
			return err
		}
		res, err := tx.Exec(c, fmt.Sprintf("update %s.enterprise_user set bind_to=?, version=version+1 where id=? and version=?", u.dbname), u1.ID, u2.ID, u2.Version)
		if err != nil {
			return err
		}
		count, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}

		bWxExternal, _ := json.Marshal(uu1.WxExternal)
		bWxShop, _ := json.Marshal(uu1.WxShop)
		bWxOpen, _ := json.Marshal(uu1.WxOpen)
		bWxUnion, _ := json.Marshal(uu1.WxUnion)
		res, err = tx.Exec(c,
			fmt.Sprintf("update %s.enterprise_user set wx_external=?, wx_shop=?, wx_open=?, wx_union=?, level_id=?, version=version+1 where id=? and version=?", u.dbname),
			string(bWxExternal), string(bWxShop), string(bWxOpen), string(bWxUnion), u1.LevelID, u1.ID, u1.Version)
		if err != nil {
			return err
		}

		count, err = res.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return uu1.ID, err
}

type wxUserOpen struct {
	wxUserBase
}

//新增公众号用户
func (u *wxUserOpen) Insert(c ctx.Context, db *sqlplus.ClusterDB) error {
	var euid, appUserID int64
	if err := db.Tx(c, func(tx *sqlplus.Tx) error {
		var err error
		euid, appUserID, err = u.insertOrUpdateEnterpriseUser(c, tx, "wx_open")
		if err != nil {
			return err
		}
		wxUserOpenQuery := fmt.Sprintf("insert into %s.%s (id,appid,unionid,openid,subscribe_at,subscribe_status,subscribe_scene,remark,qr_scene,qr_scene_str,sex,avatar,language,city,province,country,nick_name,euid,created_at) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", u.dbname, u.table)

		_, err = tx.Exec(c, wxUserOpenQuery, appUserID, u.source.AppID, u.source.UnionID,
			u.source.OpenID, u.source.SubscribeTime, u.source.SubscribeStatus,
			u.source.SubscribeScene, u.source.Remark, u.source.QRScene, u.source.QRSceneStr,
			u.source.Sex, u.source.Avatar, u.source.Language, u.source.City, u.source.Province, u.source.Country, u.source.Nick, euid, u.source.CreateAt)
		if err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	u.source.ID = appUserID
	u.source.EuID = euid

	return nil
}

func (u *wxUserOpen) Update(c ctx.Context, db *sqlplus.ClusterDB) error {
	if err := u.updateBase(c, db); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}

	cols, args := u.makeUpdateArgs()
	if len(cols) > 0 {
		if _, err := db.DB.Exec(c, fmt.Sprintf("UPDATE %s.%s SET %s WHERE id = ?", u.dbname, u.table, strings.Join(cols, ",")), args...); err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		}
	}

	return nil
}

func (u *wxUserOpen) makeUpdateArgs() ([]string, []interface{}) {
	var cols []string
	var args []interface{}

	if u.source.UnionID != "" {
		cols = append(cols, "unionid=?")
		args = append(args, u.source.UnionID)
	}

	if u.source.SubscribeTime != 0 {
		cols = append(cols, "subscribe_at=?", "subscribe_status=?", "subscribe_scene=?")
		args = append(args, u.source.SubscribeTime, u.source.SubscribeStatus, u.source.SubscribeScene)
	}

	if u.source.Remark != "" {
		cols = append(cols, "remark=?")
		args = append(args, u.source.Remark)
	}

	if u.source.QRScene != "" {
		cols = append(cols, "qr_scene=?", "qr_scene_str=?")
		args = append(args, u.source.QRScene, u.source.QRSceneStr)
	}

	if u.source.Sex != 0 {
		cols = append(cols, "sex=?")
		args = append(args, u.source.Sex)
	}

	if u.source.Language != "" {
		cols = append(cols, "language=?")
		args = append(args, u.source.Language)
	}

	if u.source.Nick != "" || u.source.Avatar != "" {
		cols = append(cols, "nick_name=?")
		args = append(args, u.source.Nick)
		cols = append(cols, "avatar=?")
		args = append(args, u.source.Avatar)
	}

	if u.source.Country != "" || u.source.Province != "" || u.source.City != "" {
		cols = append(cols, "country=?")
		args = append(args, u.source.Country)
		cols = append(cols, "province=?")
		args = append(args, u.source.Province)
		cols = append(cols, "city=?")
		args = append(args, u.source.City)
	}

	if len(cols) > 0 {
		cols = append(cols, "updated_at=?")
		args = append(args, time.Now().Unix(), u.source.ID)
	}
	return cols, args
}

func (u *wxUserOpen) InsertUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	cols, args := u.makeUpdateArgs()
	if err := u.updateDateBaseWithUnion(c, db, cols, args); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}
	return nil
}

func (u *wxUserOpen) UpdateWithUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	euid, err := u.mergeEnterpriseUser(c, db, u.unionEuid, u.source.EuID)
	if err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		return err
	}
	if euid != 0 {
		u.unionEuid = euid
		u.source.EuID = euid
		return u.Update(c, db)
	}
	return nil
}

type wxUserShop struct {
	wxUserBase
}

func (u *wxUserShop) Insert(c ctx.Context, db *sqlplus.ClusterDB) error {
	var euid, appUserID int64
	if err := db.Tx(c, func(tx *sqlplus.Tx) error {
		var err error
		euid, appUserID, err = u.insertOrUpdateEnterpriseUser(c, tx, "wx_shop")
		//给企业用户表增加数据
		if err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
			return err
		}
		wxUserShopQuery := fmt.Sprintf("insert into %s.%s (id,appid,unionid,openid,third_user_id,euid,nick_name,name,avatar,mobile,birthday,sex,country,province,city,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", u.dbname, u.table)

		_, err = tx.Exec(c, wxUserShopQuery, appUserID, u.source.AppID, u.source.UnionID,
			u.source.OpenID, u.source.ThirdUserId, euid,
			u.source.Nick, u.source.Name, u.source.Avatar, u.source.Mobile, u.source.Birthday,
			u.source.Sex, u.source.Country, u.source.Province, u.source.City, u.source.CreateAt, time.Now().Unix())
		if err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	u.source.ID = appUserID
	u.source.EuID = euid

	return nil
}

func (u *wxUserShop) Update(c ctx.Context, db *sqlplus.ClusterDB) error {
	if err := u.updateBase(c, db); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}

	cols, args := u.makeUpdateArgs()
	if len(cols) > 0 {
		if _, err := db.DB.Exec(c, fmt.Sprintf("UPDATE %s.%s SET %s WHERE id = ?", u.dbname, u.table, strings.Join(cols, ",")), args...); err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		}
	}

	return nil
}

func (u *wxUserShop) makeUpdateArgs() ([]string, []interface{}) {
	var cols []string
	var args []interface{}

	if u.source.UnionID != "" {
		cols = append(cols, "unionid=?")
		args = append(args, u.source.UnionID)
	}

	if u.source.Sex != 0 {
		cols = append(cols, "sex=?")
		args = append(args, u.source.Sex)
	}

	if u.source.Nick != "" || u.source.Avatar != "" {
		cols = append(cols, "nick_name=?")
		args = append(args, u.source.Nick)
		cols = append(cols, "avatar=?")
		args = append(args, u.source.Avatar)
	}

	if u.source.Country != "" || u.source.Province != "" || u.source.City != "" {
		cols = append(cols, "country=?")
		args = append(args, u.source.Country)
		cols = append(cols, "province=?")
		args = append(args, u.source.Province)
		cols = append(cols, "city=?")
		args = append(args, u.source.City)
	}

	if u.source.Mobile != "" {
		cols = append(cols, "mobile=?")
		args = append(args, u.source.Mobile)
	}

	if len(cols) > 0 {
		cols = append(cols, "updated_at=?")
		args = append(args, time.Now().Unix(), u.source.ID)
	}
	return cols, args
}

func (u *wxUserShop) InsertUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	cols, args := u.makeUpdateArgs()
	if err := u.updateDateBaseWithUnion(c, db, cols, args); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}
	return nil
}

func (u *wxUserShop) UpdateWithUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	euid, err := u.mergeEnterpriseUser(c, db, u.unionEuid, u.source.EuID)
	if err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		return err
	}
	if euid != 0 {
		u.unionEuid = euid
		u.source.EuID = euid
		return u.Update(c, db)
	}
	return nil
}

type wxUserExternal struct {
	wxUserBase
}

func (u *wxUserExternal) Insert(c ctx.Context, db *sqlplus.ClusterDB) error {
	var euid, appUserID int64
	if err := db.Tx(c, func(tx *sqlplus.Tx) error {
		var err error
		euid, appUserID, err = u.insertOrUpdateEnterpriseUser(c, tx, "wx_external")
		if err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
			return err
		}
		wxUserCorpQuery := fmt.Sprintf("insert into %s.%s (id,appid,unionid,openid,euid,nick,name,position,avatar,corp_name,corp_full_name,type,sex,external_profile,created_at,updated_at) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", u.dbname, u.table)

		_, err = tx.Exec(c, wxUserCorpQuery, appUserID, u.source.AppID, u.source.UnionID,
			u.source.OpenID, euid, u.source.Nick, u.source.Name, u.source.Position, u.source.Avatar, u.source.CorpName, u.source.CorpFullName,
			u.source.Type, u.source.Sex, u.source.ExternalProfile, u.source.CreateAt, time.Now().Unix())
		if err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	u.source.ID = appUserID
	u.source.EuID = euid

	return nil
}

func (u *wxUserExternal) makeUpdateArgs() ([]string, []interface{}) {
	var cols []string
	var args []interface{}

	if u.source.UnionID != "" {
		cols = append(cols, "unionid=?")
		args = append(args, u.source.UnionID)
	}

	if u.source.Sex != 0 {
		cols = append(cols, "sex=?")
		args = append(args, u.source.Sex)
	}

	if u.source.Avatar != "" {
		cols = append(cols, "avatar=?")
		args = append(args, u.source.Avatar)
	}

	if u.source.Name != "" {
		cols = append(cols, "name=?")
		args = append(args, u.source.Name)
	}

	if u.source.Nick != "" {
		cols = append(cols, "nick=?")
		args = append(args, u.source.Nick)
	}

	if u.source.Position != "" {
		cols = append(cols, "position=?")
		args = append(args, u.source.Position)
	}

	if u.source.CorpName != "" {
		cols = append(cols, "corp_name=?", "corp_full_name=?")
		args = append(args, u.source.CorpName, u.source.CorpFullName)
	}

	if u.source.Type != 0 {
		cols = append(cols, "type=?")
		args = append(args, u.source.Type)
	}

	if u.source.ExternalProfile != "" {
		cols = append(cols, "external_profile=?")
		args = append(args, u.source.ExternalProfile)
	}

	if len(cols) > 0 {
		cols = append(cols, "updated_at=?")
		args = append(args, time.Now().Unix(), u.source.ID)
	}
	return cols, args
}

func (u *wxUserExternal) Update(c ctx.Context, db *sqlplus.ClusterDB) error {
	if err := u.updateBase(c, db); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}

	cols, args := u.makeUpdateArgs()
	if len(cols) > 0 {
		if _, err := db.DB.Exec(c, fmt.Sprintf("UPDATE %s.%s SET %s WHERE id = ?", u.dbname, u.table, strings.Join(cols, ",")), args...); err != nil {
			exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		}
	}

	return nil
}

func (u *wxUserExternal) InsertUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	cols, args := u.makeUpdateArgs()
	if err := u.updateDateBaseWithUnion(c, db, cols, args); err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
	}

	return nil
}

func (u *wxUserExternal) UpdateWithUnion(c ctx.Context, db *sqlplus.ClusterDB) error {
	euid, err := u.mergeEnterpriseUser(c, db, u.unionEuid, u.source.EuID)
	if err != nil {
		exception.SendError(err, map[string]interface{}{"MergeUser": u.source})
		return err
	}
	if euid != 0 {
		u.unionEuid = euid
		u.source.EuID = euid
		return u.Update(c, db)
	}
	return nil
}

// NewWxUser 初始化用户
func NewWxUser(source *user.SourceUser) WxUser {
	var wxuser WxUser
	switch source.CustomerType {
	case customer.CustomerTypeOpen:
		wxuser = &wxUserOpen{wxUserBase: wxUserBase{table: "wx_user_open", source: source}}
	case customer.CustomerTypeShop:
		wxuser = &wxUserShop{wxUserBase: wxUserBase{table: "wx_user_shop", source: source}}
	case customer.CustomerTypeExternal:
		wxuser = &wxUserExternal{wxUserBase: wxUserBase{table: "wx_user_corp_external", source: source}}
	}
	return wxuser
}
