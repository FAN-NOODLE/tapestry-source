package storeOperationJob

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	appModels "micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
)

type dashboardAreaSyncStats struct {
	Inserted    int
	Updated     int
	Unchanged   int
	Invalidated int
	Skipped     int
}

type dashboardAreaSyncPlan struct {
	Creates     []dashboardModel.DashboardAreaModel
	Updates     map[int64]map[string]interface{}
	Invalidates []int64
	Stats       dashboardAreaSyncStats
}

func planDashboardAreaSync(existing []dashboardModel.DashboardAreaModel, desired []dashboardAreaNode, protectedStoreCodes map[string]struct{}, closedAtByStore map[string]int64, hierarchyStoreCount int, now int64) (dashboardAreaSyncPlan, error) {
	plan := dashboardAreaSyncPlan{Updates: make(map[int64]map[string]interface{})}
	if len(desired) == 0 && (hierarchyStoreCount > 0 || len(protectedStoreCodes) == 0 && len(closedAtByStore) == 0) {
		return plan, fmt.Errorf("LMS未生成有效的中国区域组织架构节点")
	}
	existingByKey := make(map[dashboardAreaNodeKey]dashboardModel.DashboardAreaModel, len(existing))
	var maxID int64
	protectedIDs := protectedDashboardAreaIDs(existing, protectedStoreCodes)
	for _, area := range existing {
		if area.ID > maxID {
			maxID = area.ID
		}
		if area.LMSID == "" {
			continue
		}
		key := dashboardAreaNodeKey{Type: area.Type, LMSID: area.LMSID}
		if _, duplicated := existingByKey[key]; duplicated {
			return plan, fmt.Errorf("dashboard_area存在重复LMS节点：type=%d,lms_id=%s", area.Type, area.LMSID)
		}
		existingByKey[key] = area
	}

	resolved := make(map[dashboardAreaNodeKey]dashboardModel.DashboardAreaModel, len(desired))
	desiredKeys := make(map[dashboardAreaNodeKey]struct{}, len(desired))
	for _, node := range desired {
		key := dashboardAreaNodeKey{Type: node.Type, LMSID: node.LMSID}
		desiredKeys[key] = struct{}{}
		parentID, areaPath, err := resolveDashboardAreaParent(node, resolved)
		if err != nil {
			return plan, err
		}
		area, found := existingByKey[key]
		if !found {
			maxID++
			area = dashboardModel.DashboardAreaModel{
				BaseModel: appModels.BaseModel{ID: maxID, CreatedAt: now, UpdatedAt: now},
				Type:      node.Type, Name: node.Name, ParentID: parentID, Path: areaPath,
				ThirdID: node.ThirdID, IsVirtual: 0, LMSID: node.LMSID,
				ClosedAt: node.ClosedAt, Sort: node.Sort,
			}
			plan.Creates = append(plan.Creates, area)
			plan.Stats.Inserted++
		} else {
			updates := dashboardAreaChanges(area, node, parentID, areaPath, now)
			if len(updates) == 0 {
				plan.Stats.Unchanged++
			} else {
				plan.Updates[area.ID] = updates
				applyDashboardAreaChanges(&area, node, parentID, areaPath)
				plan.Stats.Updated++
			}
		}
		resolved[key] = area
	}

	for storeCode, closedAt := range closedAtByStore {
		area, found := existingByKey[dashboardAreaNodeKey{Type: 4, LMSID: storeCode}]
		if !found {
			continue
		}
		if area.ClosedAt == closedAt {
			plan.Stats.Unchanged++
			continue
		}
		updates := plan.Updates[area.ID]
		if updates == nil {
			updates = make(map[string]interface{})
			plan.Updates[area.ID] = updates
			plan.Stats.Updated++
		}
		updates["closed_at"] = closedAt
		updates["updated_at"] = now
	}

	for _, area := range existing {
		if area.LMSID == "" || area.DeletedAt > 0 {
			continue
		}
		key := dashboardAreaNodeKey{Type: area.Type, LMSID: area.LMSID}
		_, protected := protectedIDs[area.ID]
		if _, exists := desiredKeys[key]; !exists && !protected {
			plan.Invalidates = append(plan.Invalidates, area.ID)
			plan.Stats.Invalidated++
		}
	}
	sort.Slice(plan.Invalidates, func(i, j int) bool { return plan.Invalidates[i] < plan.Invalidates[j] })
	return plan, nil
}

func protectedDashboardAreaIDs(existing []dashboardModel.DashboardAreaModel, protectedStoreCodes map[string]struct{}) map[int64]struct{} {
	protectedIDs := make(map[int64]struct{})
	if len(protectedStoreCodes) == 0 {
		return protectedIDs
	}
	for _, area := range existing {
		if area.Type != 4 {
			continue
		}
		if _, protected := protectedStoreCodes[area.LMSID]; !protected {
			continue
		}
		protectedIDs[area.ID] = struct{}{}
		for _, pathID := range strings.Split(area.Path, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(pathID), 10, 64)
			if err == nil && id > 0 {
				protectedIDs[id] = struct{}{}
			}
		}
	}
	return protectedIDs
}

func resolveDashboardAreaParent(node dashboardAreaNode, resolved map[dashboardAreaNodeKey]dashboardModel.DashboardAreaModel) (int, string, error) {
	if node.Type == 1 {
		return 0, "", nil
	}
	parentKey := dashboardAreaNodeKey{Type: node.Type - 1, LMSID: node.ParentLMSID}
	parent, exists := resolved[parentKey]
	if !exists {
		return 0, "", fmt.Errorf("找不到LMS节点%s的父节点%s", node.LMSID, node.ParentLMSID)
	}
	parentPath := strconv.FormatInt(parent.ID, 10)
	if parent.Path != "" {
		parentPath = parent.Path + "," + parentPath
	}
	return int(parent.ID), parentPath, nil
}

func dashboardAreaChanges(area dashboardModel.DashboardAreaModel, node dashboardAreaNode, parentID int, areaPath string, now int64) map[string]interface{} {
	updates := make(map[string]interface{})
	if area.Name != node.Name {
		updates["name"] = node.Name
	}
	if area.ParentID != parentID {
		updates["parent_id"] = parentID
	}
	if area.Path != areaPath {
		updates["path"] = areaPath
	}
	if area.ThirdID != node.ThirdID {
		updates["third_id"] = node.ThirdID
	}
	if area.IsVirtual != 0 {
		updates["is_virtual"] = 0
	}
	if area.ClosedAt != node.ClosedAt {
		updates["closed_at"] = node.ClosedAt
	}
	if area.Sort != node.Sort {
		updates["sort"] = node.Sort
	}
	if area.DeletedAt != 0 {
		updates["deleted_at"] = int64(0)
	}
	if len(updates) > 0 {
		updates["updated_at"] = now
	}
	return updates
}

func applyDashboardAreaChanges(area *dashboardModel.DashboardAreaModel, node dashboardAreaNode, parentID int, areaPath string) {
	area.Name = node.Name
	area.ParentID = parentID
	area.Path = areaPath
	area.ThirdID = node.ThirdID
	area.IsVirtual = 0
	area.ClosedAt = node.ClosedAt
	area.Sort = node.Sort
	area.DeletedAt = 0
}
