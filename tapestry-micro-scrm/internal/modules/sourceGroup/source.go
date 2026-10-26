package sourceGroup

import (
	"errors"
	"github.com/jinzhu/copier"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	source2 "micro-scrm/internal/modules/source"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/source"
	"strconv"
	"time"
)

func AddSource(ctx ctx.Context, info NewSource) error {
	res, err := orm.Service.Source.GetMaxOrder()
	if err != nil {
		return err
	}
	//res2 := orm.Service.Source.GetSourceLevel(info.ParentID)
	//if info.ParentID != 0 && res2 != 0 {
	//	return errors.New("素材分级最多三级，请确认!")
	//}
	newinfo := source.GroupModel{
		MicroScrmBaseModel: model.MicroScrmBaseModel{
			BaseModel: model.BaseModel{
				//ID:        uid.ID(ctx),
				ID:        info.ID,
				CreatedAt: time.Now().Unix(),
			},
			EnterpriseID: info.EnterpriseID,
		},
		GroupName: info.GroupName,
		ParentId:  info.ParentID,
		Orders:    int(res) + 1,
	}
	if err := orm.Service.Source.AddSource(newinfo); err != nil {
		return err
	}
	return nil
}

func NameMap(enterpriseID, parentID int64) map[string]struct{} {
	var result = make(map[string]struct{}, 0)
	res, _ := orm.Service.Source.GetSourceName(enterpriseID, parentID)
	for _, v := range res {
		result[v] = struct{}{}
	}

	return result
}

func EditSourceInfo(info EditSource) error {
	updatesMap := make(map[string]interface{}, 0)
	updatesMap["group_name"] = info.GroupName
	updatesMap["updated_at"] = time.Now().Unix()
	if err := orm.Service.Source.EditSource(info.ID, info.EnterpriseID, updatesMap, source.GroupModel{}); err != nil {
		return nil
	}
	return nil
}

func UpdateIgnoreIndex(info UpdateIgnoreIndexReq) error {
	updatesMap := make(map[string]interface{}, 0)
	updatesMap["ignore_index"] = info.IgnoreIndex
	updatesMap["updated_at"] = time.Now().Unix()
	if err := orm.Service.Source.EditSource(info.ID, info.EnterpriseID, updatesMap, source.GroupModel{}); err != nil {
		return nil
	}
	return nil
}

func DelSource(ID int64) error {
	if err := orm.Service.Source.RemoveSource(ID, source.GroupModel{}); err != nil {
		return err
	}
	if err := orm.Service.Source.RemoveSourceMaterial(ID); err != nil {
		return err
	}
	return nil
}

func DelVerify(ID, enterpriseID int64) (result []int64, err error) {
	result, _ = source2.GetChildIDList(ID, enterpriseID)
	//result, err = orm.Service.Source.RemoveVerify(list, source.SourceModel{})
	return
}

// SourceTree 分组无限级树
func SourceTree(enterpriseId int64, paramsID int64) (result []*micro_scrm.SourceTree, err error) {
	var info source.GroupModel
	info.EnterpriseID = enterpriseId
	res, err := orm.Service.Source.SourceInfo(info)
	if err != nil || res == nil {
		return nil, err
	}
	treeMap := make(map[int64]*micro_scrm.SourceTree, 0)
	treeList := make([]*micro_scrm.SourceTree, 0)
	newList := make([]*micro_scrm.SourceTree, 0)
	copier.Copy(&treeList, &res)
	for _, v := range treeList {
		treeMap[v.ID] = v
	}
	for _, sourceVal := range treeList {
		if sourceVal.ParentId == 0 {
			newList = append(newList, sourceVal)
		} else {
			if val, ok := treeMap[sourceVal.ParentId]; ok {
				if val.Children == nil {
					val.Children = make([]*micro_scrm.SourceTree, 0)
				}
				val.Children = append(val.Children, sourceVal)
			}
		}
	}

	// 处理层级
	handleTree(newList, 1)

	// 获取某个层级的分组
	if paramsID != 0 {
		newListOne := handleTreeFindInfo(newList, paramsID)
		return newListOne, nil
	}

	return newList, nil
}

// handleTreeFindInfo 查找分类下所有子级
func handleTreeFindInfo(newList []*micro_scrm.SourceTree, paramsID int64) (result []*micro_scrm.SourceTree) {

	newListOne := make([]*micro_scrm.SourceTree, 0)
	if len(newList) == 0 {
		return nil
	}
	for _, v1 := range newList {
		if v1.ID == paramsID {
			newListOne = append(newListOne, v1)
			return newListOne
		}
		res := handleTreeFindInfo(v1.Children, paramsID)
		if res != nil {
			return res
		}

	}
	return nil
}

// handleTree 递归处理层级
func handleTree(newList []*micro_scrm.SourceTree, level int) (result []*micro_scrm.SourceTree) {

	for _, v1 := range newList {
		v1.Level = level
		if level == 1 {
			v1.Ids = strconv.FormatInt(v1.ID, 10)
		}
		if v1.Children != nil {
			for _, v2 := range v1.Children {
				v2.Ids = v1.Ids + "," + strconv.FormatInt(v2.ID, 10)

				handleTree(v1.Children, level+1)
			}
		}
	}
	return newList

}

// SourceInfoById 根据id查询分组信息
func SourceInfoById(enterpriseId int64, id int64) (result []*micro_scrm.SourceTree, err error) {
	var info source.GroupModel
	info.EnterpriseID = enterpriseId
	info.ID = id

	res, err := orm.Service.Source.SourceInfoById(info)
	if err != nil || res == nil {
		return nil, err
	}
	treeList := make([]*micro_scrm.SourceTree, 0)
	copier.Copy(&treeList, &res)

	return treeList, nil
}

// MoveSource 分组移动
func MoveSource(infos MoveReceive) error {
	ids := []int64{infos.MoverId, infos.SwaperId}
	list, err := orm.Service.Source.GetListByIds(ids, infos.EnterpriseID)
	if err != nil {
		return err
	}

	if len(list) < 2 {
		return errors.New("分组参数错误")
	}
	if list[0].ParentId != list[1].ParentId {
		return errors.New("分组移动仅能在当前层级内移动，不能跨层级移动分组")
	}
	var moveInfo, targetInfo source.GroupModel
	for _, groupModel := range list {
		if groupModel.ID == infos.MoverId {
			moveInfo = groupModel
		} else if groupModel.ID == infos.SwaperId {
			targetInfo = groupModel
		}
	}
	//先把要插入的位置之后的数据往后调整
	if err := orm.Service.Source.IncrementOrder(infos.EnterpriseID, moveInfo.ParentId, targetInfo.Orders); err != nil {
		return err
	}
	//把目标分组排序放在原始位置
	if err := orm.Service.Source.UpdateOrder(moveInfo.ID, targetInfo.Orders); err != nil {
		return err
	}
	return nil
}
