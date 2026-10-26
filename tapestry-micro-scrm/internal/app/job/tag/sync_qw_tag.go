package tagJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gorm.io/gorm"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type SyncQwTag struct {
	name string
}

func NewSyncQwTag() SyncQwTag {
	return SyncQwTag{
		name: "同步企微标签",
	}
}

func (job SyncQwTag) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	// var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		// wg.Add(1)
		// go func(enterpriseId int64) {
		//	defer wg.Done()
		err = job.handleEnterprise(ctx, enterprise.ID)
		if err != nil {
			log.Error(ctx).Msgf("处理企业%d异常,错误:%v", enterprise.ID, err)
		}
		// }(enterprise.ID)
	}
	// wg.Wait()
	return nil
}

func (job SyncQwTag) handleEnterprise(ctx ctx.Context, enterpriseId int64) error {

	var (
		workWeChatObj      *workwechat.WorkWeChat
		tagGroupInfoList   []*work.CorpTagGroupEntity
		err                error
		db                 = gormx.GetInstance()
		sql                string
		qyGroupEntry       = make(map[string]*work.CorpTagGroupEntity)
		weSocialGroup      []*entity.WorkExternalTagGroupModel
		weSocialGroupEntry = make(map[string]*entity.WorkExternalTagGroupModel)
	)
	// 查询企业微信中的标签信息
	if workWeChatObj, err = workwechat.NewWorkWeChat(ctx, enterpriseId); err != nil {
		return err
	}
	if tagGroupInfoList, err = workWeChatObj.TagList(); err != nil {
		return err
	}

	for _, item := range tagGroupInfoList {
		qyGroupEntry[item.GroupID] = item
	}
	// 查询本地数据库标签组信息
	sql = fmt.Sprintf("select  id,group_id from  work_external_tag_group  where enterprise_id = ? and deleted_at = 0")
	if err = db.Raw(sql, enterpriseId).Scan(&weSocialGroup).Error; nil != err {
		log.Error(ctx).Msgf("根据企业id查询标签组信息失败,enterprise_id:%d,error:%s", enterpriseId, err.Error())
		return err
	}
	for _, item := range weSocialGroup {
		weSocialGroupEntry[item.GroupID] = item
	}

	/**
		标签组
	        从企业微信角度出发，如果wesocial不包含企业微信标签组 新建标签组

		    从管理后台出发  如果企业微信不包含 wesocial的标签组 删除wesocial标签组
		标签
	*/
	// 交叉对比 标签组
	for _, item := range tagGroupInfoList {
		// 企业微信标签不存在
		if _, ok := weSocialGroupEntry[item.GroupID]; !ok {
			now := time.Now().Unix()
			groupID := uid.ID(ctx)
			// 创建标签组
			sql = "insert into work_external_tag_group (id,enterprise_id,group_id,group_name,`order`,`type`,create_time,created_at,updated_at)" +
				"values (?,?,?,?,?,?,?,?,?) "
			if err = db.Exec(sql, groupID, enterpriseId, item.GroupID, item.GroupName, item.Order, entity.TagGroupBasic, item.CreateTime, now, now).Error; nil != err {
				log.Error(ctx).Msgf("企业微信存在，wesocial不存在的标签组创建失败，work_group_id:%s,error:%s", item.GroupID, err.Error())
				continue
				// return err
			}
			// 创建标签
			for _, tag := range item.Tag {
				tagID := uid.ID(ctx)
				sql = "insert into work_external_tag (id,enterprise_id,group_id,tag_id,name,`order`,create_time,created_at,updated_at)" +
					" values  (?,?,?,?,?,?,?,?,?)"
				if err = db.Exec(sql, tagID, enterpriseId, groupID, tag.ID, tag.Name, tag.Order, tag.CreateTime, now, now).Error; nil != err {
					log.Error(ctx).Msgf("企业微信存在，wesocial不存在的标签组创建成功,标签创建失败.group_id:%d,work_group_id:%s,tag_id:%d,error:%s",
						groupID, item.GroupID, tagID, err.Error())
					// return err
				}
			}

		}

		/**
			企业微信标签存在wesocial 比较标签情况,企业微信存在但wesocial不存在，wesocial
		后台删除;wesocial存在，企业微信不存在，wesocial删除
		*/
		if value, ok := weSocialGroupEntry[item.GroupID]; ok {

			var (
				wesocialTagEntry = make(map[string]*entity.WorkExternalTagModel)
				workTagEntry     = make(map[string]*work.CorpTagEntity)
				tagList          []*entity.WorkExternalTagModel
			)
			sql = "select id,tag_id as work_tag_id ,`type`  from work_external_tag where group_id = ? and deleted_at =0"
			if err = db.Raw(sql, value.ID).Scan(&tagList).Error; nil != err {
				return err
			}
			for _, tag := range tagList {
				wesocialTagEntry[tag.WorkTagID] = tag
			}
			for _, tag := range item.Tag {
				workTagEntry[tag.ID] = tag
			}

			for _, tag := range workTagEntry {

				if _, ok := wesocialTagEntry[tag.ID]; !ok {
					// 同一个标签组下 wesocial不存在企业微信的标签,wesocial创建
					tagID := uid.ID(ctx)
					now := time.Now().Unix()
					sql = "insert into work_external_tag (id,enterprise_id,group_id,tag_id,`name`,`order`,create_time,created_at,updated_at)" +
						" values  (?,?,?,?,?,?,?,?,?)"

					if err = db.Exec(sql, tagID, enterpriseId, value.ID, tag.ID, tag.Name, tag.Order, tag.CreateTime, now, now).Error; nil != err {
						log.Error(ctx).Msgf("企业微信标签组存在wesocial,但不存在wesocial的标签创建失败.group_id:%d,work_tag_id:%s,error:%s",
							value.ID, tag.ID, err.Error())
						// return err
					}

				}
				// 同一个标签组下 wesocial存在企业微信的标签, 更新标签名和order
				if _, ok := wesocialTagEntry[tag.ID]; ok {
					sql = "update work_external_tag set `name` = ? , `order`= ?  where deleted_at = 0 and enterprise_id = ? and tag_id = ?"
					if err = db.Exec(sql, tag.Name, tag.Order, enterpriseId, tag.ID).Error; nil != err {
						log.Error(ctx).Msgf("更新wesocial标签信息失败:%s", err.Error())
					}

				}
			}
			for _, tag := range wesocialTagEntry {
				if _, ok := workTagEntry[tag.WorkTagID]; !ok {
					// 同一个标签组下 wesocial存在标签 企业微信不存在 wesocial删除
					if err = db.Transaction(func(tx *gorm.DB) error {
						now := time.Now().Unix()
						sql = "updated  work_external_tag  set deleted_at = ?  where id = ?"

						if err = tx.Exec(sql, now, tag.ID).Error; nil != err {
							log.Error(ctx).Msgf("企业微信标签存在wesocial，但wesocial存在企业微信不存在的标签删除失败.group_id:%d,tag_id:%d,error:%s",
								value.ID, tag.ID, err.Error())
							return err
						}
						// 标签类型为第三方标签 增加删除标签异步任务
						// if tag.Type == entity.TagThirdPart {
						sql = "insert into `work_external_tag_del_task` (enterprise_id,target_id,`type`,created_at,updated_at) VALUES (?,?,?,?,?)"
						if err = tx.Exec(sql, enterpriseId, tag.ID, entity.TaskTag, now, now).Error; nil != err {
							log.Error(ctx).Msgf("wesocial不存在企业微信的标签,增加异步删除任务失败. tag_id:%d,error:%s", tag.ID, err.Error())
							return err
						}
						// }
						return nil
					}); nil != err {
						log.Error(ctx).Msgf("企业微信标签组存在wesocial，但wesocial存在企业微信不存在的标签删除失败.group_id:%d,tag_id:%d,error:%s",
							value.ID, tag.ID, err.Error())
					}

				}

			}

		}
	}

	for _, item := range weSocialGroup {
		// wesocial标签组不存在企业微信 删除wesocial标签组
		if _, ok := qyGroupEntry[item.GroupID]; !ok {
			now := time.Now().Unix()
			if err = db.Transaction(func(tx *gorm.DB) error {

				log.Info(ctx).Msgf("删除wesocial存在但企业微信不存在的标签组，group_id:%d", item.ID)
				// 删除标签组
				sql = "update work_external_tag_group    set  deleted_at =?  where  deleted_at =0 and  id = ?"
				if err = tx.Exec(sql, now, item.ID).Error; nil != err {
					log.Error(ctx).Msgf("wesocial不存在企业微信的标签组删除失败. group_id:%d,error:%s", item.ID, err.Error())
					return err
				}
				// 删除标签组下的标签
				sql = "update  work_external_tag set deleted_at = ? where deleted_at =0 and group_id = ?"
				if err = tx.Exec(sql, now, item.ID).Error; nil != err {
					log.Error(ctx).Msgf("wesocial不存在企业微信的标签组下标签删除失败. group_id:%d,error:%s", item.ID, err.Error())
					return err
				}
				// 增加删除第三方标签组 异步任务
				// if item.Type == entity.TagGroupThird {
				sql = "insert into `work_external_tag_del_task` (enterprise_id,target_id,`type`,created_at,updated_at) VALUES (?,?,?,?,?)"
				if err = tx.Exec(sql, enterpriseId, item.ID, entity.TaskGroup, now, now).Error; nil != err {
					log.Error(ctx).Msgf("wesocial不存在企业微信的标签组,增加异步删除任务失败. group_id:%d,error:%s", item.ID, err.Error())
					return err
				}
				/*	log.Info(ctx).Msgf("增加标签删除任务成功，target_id:%d,任务类型:%d", item.ID, entity.TaskGroup)
					}*/
				return nil
			}); nil != err {
				log.Error(ctx).Err(err)
			}
		}
	}
	return nil
}
