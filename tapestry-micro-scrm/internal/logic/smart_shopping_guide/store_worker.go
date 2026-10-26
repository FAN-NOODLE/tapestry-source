package smart_shopping_guide

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/work"
)

// 功能列表
// 1、获取导购集合，可以分类型，也可以获取全部
// 2、根据智慧导购的ID，获取该导购企微的user_id
// 3、根据智慧导购的ID集合，获取这些导购的企业的user_id集合

// GetBaUserIDList 1、获取导购集合，可以分类型，也可以获取全部
func (guide smartShoppingGuideLogic) GetBaUserIDList(ctx ctx.Context, baType []uint8) (userIDList []string, err error) {
	userInfoList, infoErr := orm.WorkService.Work.GetAllBAInfoList(guide.enterpriseID, baType)
	if infoErr != nil {
		log.Info(ctx).Msgf("%s::获取类型:%d导购异常，错误为:%v", guide.name, baType, infoErr)
		err = infoErr
		return
	}

	for _, item := range userInfoList {
		if item.MemUserid != "" {
			userIDList = append(userIDList, item.MemUserid)
		}
	}

	return
}

// GetBaUserID 2、根据智慧导购的ID，获取该导购企微的user_id
func (guide smartShoppingGuideLogic) GetBaUserID(ctx ctx.Context, memUserID string) (userID string, err error) {
	userInfo, infoErr := orm.WorkService.Work.GetBAInfo(guide.enterpriseID, memUserID)
	if infoErr != nil {
		log.Info(ctx).Msgf("%s::获取导购ID为%s的导购信息异常，错误为:%v", guide.name, memUserID, infoErr)
		err = infoErr
		return
	}

	userID = userInfo.MemUserid
	return
}

// GetBaUserIDListPart 3、根据智慧导购的ID集合，获取这些导购的企业的user_id集合
func (guide smartShoppingGuideLogic) GetBaUserIDListPart(ctx ctx.Context, baType []uint8, memUserIDList []string) (userIDList []string, err error) {
	userInfoList, infoErr := orm.WorkService.Work.GetBAInfoListPart(guide.enterpriseID, baType, memUserIDList)
	if infoErr != nil {
		log.Info(ctx).Msgf("%s::获取类型:%v导购异常，错误为:%v", guide.name, memUserIDList, infoErr)
		err = infoErr
		return
	}

	for _, item := range userInfoList {
		if item.MemUserid != "" {
			userIDList = append(userIDList, item.MemUserid)
		}
	}

	return
}

// 分页查询BA的userId
func (guide smartShoppingGuideLogic) GetPagingBaUserIDList(ctx ctx.Context, baType []uint8, memUserIDList []string, page, size int) ([]*work.StoreWorkerModel, int64, error) {
	return orm.WorkService.Work.GetPagingAllBAInfoList(guide.enterpriseID, baType, memUserIDList, page, size)
}

// GetBaList 查询在职BA
func (guide smartShoppingGuideLogic) GetBaList(ctx ctx.Context, baType []uint8, memUserIDList []string) ([]*work.StoreWorkerModel, int64, error) {
	return orm.WorkService.Work.GetAllBAList(guide.enterpriseID, baType, memUserIDList)
}