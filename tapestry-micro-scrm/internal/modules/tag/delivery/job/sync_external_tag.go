package job

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRepository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"sync"
)

type syncQWExternalTag struct {
	name      string
	wg        sync.WaitGroup
	memberIds sync.Map
}

func NewSyncQwExternalTag() *syncQWExternalTag {
	return &syncQWExternalTag{
		name: "同步外部联系人企微标签",
	}
}

// Run 脚本任务处理外部联系人的企微标签
func (job *syncQWExternalTag) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	specsArgs := handler.Runner().Specs()
	if len(specsArgs) != 1 {
		// fmt.Println("参数有误，格式为--job.spec=enterprise_id")
		return nil
	}
	if specsArgs[0] == "" {
		// fmt.Println("企业id不能为空，格式为--job.spec=enterprise_id")
		return nil
	}

	var enterpriseID int64
	enterpriseID, _ = strconv.ParseInt(specsArgs[0], 10, 64)

	// fmt.Printf("enterpriseID：%v",enterpriseID)

	err := job.memberHandle(ctx, enterpriseID)
	if err != nil {
		return err
	}

	return nil
}

func (job syncQWExternalTag) memberHandle(ctx ctx.Context, enterpriseId int64) error {
	countNum, err := orm.Service.Member.GetMemberCount(enterpriseId)
	if err != nil {
		return err
	}
	memberListChannel := make(chan micro_scrm.MemberModel, 100)

	if countNum > 100 {
		rangeNum := countNum / 100
		for i := 0; i < int(rangeNum); i++ {
			offset := i * 100
			job.handleMemberExternal(ctx, enterpriseId, offset, 100, memberListChannel)
		}
		close(memberListChannel)

	} else {
		job.handleMemberExternal(ctx, enterpriseId, 0, 100, memberListChannel)
		close(memberListChannel)

	}
	job.wg.Wait()
	fmt.Printf("任务完成\n")
	return nil
}

func (job *syncQWExternalTag) handleMemberExternal(ctx ctx.Context, enterpriseId int64, offset, limit int, memberListChannel chan micro_scrm.MemberModel) error {
	memberList, err := orm.Service.Member.GetMemberListPage(enterpriseId, offset, limit)
	if err != nil {
		return err
	}
	for _, memberInfo := range memberList {
		memberListChannel <- memberInfo
	}
	var errorInfo *error
	job.wg.Add(1)
	go func(errorInfo *error) {
		err := job.handleExternalTag(ctx, enterpriseId, memberListChannel)
		if err != nil {
			errorInfo = &err
			return
		}
		return
	}(errorInfo)
	var index int
	if offset == 0 {
		index = 1
	} else {
		index = offset / 100
	}
	fmt.Printf("启动了%v个groutine\n", index)
	return nil
}

// handleExternalTag  处理导购成员的外部联系人的企微标签
func (job *syncQWExternalTag) handleExternalTag(ctx ctx.Context, enterpriseId int64, memberListChannel chan micro_scrm.MemberModel) error {
	// fmt.Printf("处理企业%d----hanle", enterpriseId)
	defer job.wg.Done()
	QyWorkImpl, err := job.getWorkImpl(ctx, enterpriseId)
	if err != nil {
		// fmt.Printf("处理企业%d获取企业微信Impl异常,错误:%v", enterpriseId, err)
		// log.Error(ctx).Msgf("处理企业%d获取企业微信Impl异常,错误:%v", enterpriseId, err)
		return err
	}
	// memberId := 3360651868056384
	for value := range memberListChannel {
		memberId := value.ID
		if _, ok := job.memberIds.LoadOrStore(memberId, memberId); ok {
			continue
		}

		countNum, errCount := orm.Service.WorkExternalUser.GetExternalUserCountByMemberId(enterpriseId, memberId)
		if errCount != nil {
			return errCount
		}
		if countNum > 10 {
			rangeNum := countNum / 10
			for i := 0; i < int(rangeNum); i++ {
				handleTag(ctx, enterpriseId, QyWorkImpl, memberId, i, countNum)
			}
		} else {
			handleTag(ctx, enterpriseId, QyWorkImpl, memberId, 0, countNum)
		}
		syncTagData := work_external_user.MemberExternalSyncTagLogModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(ctx),
				},
				EnterpriseID: enterpriseId,
			},
			MemberID:    memberId,
			ExternalNum: int(countNum),
		}
		orm.Service.WorkExternalUser.CreateSyncLogData(syncTagData)
		fmt.Printf("导购%v的数据完成处理完成", memberId)
	}
	return nil
}

// getWorkImpl 获取企业微信的impl，调用企业微信接口
func (job *syncQWExternalTag) getWorkImpl(ctx ctx.Context, enterpriseId int64) (*work.QyWorkImpl, error) {
	enterpriseService := enterprise.New(enterpriseRepository.New())
	enterpriseData, enterpriseInfoErr := enterpriseService.Detail(ctx, enterprise.Corp{
		EnterpriseID: enterpriseId,
	})
	if enterpriseInfoErr != nil {
		// fmt.Printf("获取企业信息错误:" + enterpriseInfoErr.Error())
		// log.Info(ctx).Msgf("获取企业信息错误:" + enterpriseInfoErr.Error())
		return nil, enterpriseInfoErr
	}
	QyWorkImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	return QyWorkImpl, nil
}

// 执行处理查询用户信息企微，开始更新标签
func handleTag(ctx ctx.Context, enterpriseId int64, workImpl *work.QyWorkImpl, memberId int64, index int, countNum int64) error {
	offset := index * 10
	list, err := orm.Service.WorkExternalUser.GetExternalUserInfoByMemberId(enterpriseId, memberId, offset, 10)
	if err != nil {
		// fmt.Printf("获取企业信息错误SQL:" + err.Error())
		return err
	}
	type result struct {
		Openid string
	}
	fmt.Printf("导购%v的数据开始处理\n", memberId)
	for _, item := range list {
		workExternalUserId := item.WorkExternalUserID
		RelationID := item.ID
		var resultInfo result
		sql := fmt.Sprintf("select openid from micro_scrm_enterprise_%d.wx_user_corp_external where id=?", enterpriseId)
		orm.Service.WorkExternalUser.DB.Raw(sql, workExternalUserId).Scan(&resultInfo)
		userInfo, err := workImpl.ExternalDetail(ctx, work.ExternalContactRequest{
			ExternalUserID: resultInfo.Openid,
		})
		if err != nil {
			// fmt.Printf("获取外部联系人关联错误SQL:" + err.Error())
			// errorInfo = &err
			return err
		}
		followInfo := userInfo.FollowUser
		var tagIds []string
		for _, value := range followInfo {
			if len(value.Tags) > 0 {
				for _, v := range value.Tags {
					tagIds = append(tagIds, v.TagID)
				}
			}

		}
		tagIds = removeDuplicatesAndEmpty(tagIds)
		if len(tagIds) == 0 {
			// fmt.Printf("1--获取外部联系人%v的标签数据为空\n", workExternalUserId)
			continue
		}
		tagIdList := orm.Service.Tag.GetTagIDListByQYID(enterpriseId, tagIds)
		if len(tagIdList) == 0 {
			// fmt.Printf("3-----获取外部联系人%v的标签id为空\n", workExternalUserId)
			continue
		}
		for _, id := range tagIdList {
			_, row, _ := orm.Service.Tag.GetEnterpriseTagUserInfo(enterpriseId, RelationID, id)
			if row == 0 {
				userTagData := work_external_user.WorkExternalUserTagModel{
					MicroScrmBaseModel: model.MicroScrmBaseModel{
						BaseModel: model.BaseModel{
							ID: uid.ID(ctx),
						},
						EnterpriseID: enterpriseId,
					},
					RelationID:    RelationID,
					ExternalTagID: id,
				}
				_, _, errInsert := orm.Service.Tag.CreateUserTag(userTagData)
				if errInsert != nil {
					return errInsert
				}
			}

		}

	}
	return nil
}

func removeDuplicatesAndEmpty(a []string) (ret []string) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
