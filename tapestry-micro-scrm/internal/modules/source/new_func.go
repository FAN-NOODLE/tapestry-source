package source

import (
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/source"
)

// GetChildIDList 获取分组下的所有子id
func GetChildIDList(parentID, enterpriseID int64) ([]int64, error) {
	var info source.GroupModel
	info.EnterpriseID = enterpriseID
	res, err := orm.Service.Source.GetGroupID(info)
	if err != nil {
		return nil, err
	}
	groupIDList := []int64{parentID}
	groupIDList = append(groupIDList, getLevelAllGroupID(groupIDList, res)...)
	if parentID == 0 {
		groupIDList = append(groupIDList, -1)
	}
	if len(groupIDList) == 0 {
		return nil, err
	}
	return groupIDList, err
}

func getLevelAllGroupID(groupIDList []int64, groupInfo []source.GroupModel) (nodeIDList []int64) {
	for _, item := range groupInfo {
		for _, groupID := range groupIDList {
			if groupID == item.ParentId {
				nodeIDList = append(nodeIDList, item.ID)
			}
		}
	}

	if len(nodeIDList) == 0 {
		return
	}

	return append(nodeIDList, getLevelAllGroupID(nodeIDList, groupInfo)...)
}

func GetGroupName(groupID int64) (result micro_scrm.ReceiveGroupName) {
	result = orm.Service.Source.GetGroupName(groupID)
	return
}
