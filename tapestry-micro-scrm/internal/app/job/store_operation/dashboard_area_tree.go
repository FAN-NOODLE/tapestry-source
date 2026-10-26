package storeOperationJob

import (
	"fmt"
	"sort"
	"strings"
	"time"

	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

const maxDashboardAreaTimestamp int64 = 1<<31 - 1

var dashboardDistrictCodes = map[string]struct{}{
	"CAFP_CN_CD":  {},
	"CAFP_CN_ED":  {},
	"CAFP_CN_HD":  {},
	"CAFP_CN_ND":  {},
	"CAFP_CN_SWD": {},
	"CAFS_CN_CD":  {},
	"CAFS_CN_ESD": {},
	"CAFS_CN_HD":  {},
	"CAFS_CN_ND":  {},
	"CAFS_CN_WD":  {},
	"CAFP_CN_CSD": {},
}

type dashboardAreaNode struct {
	Type        int
	LMSID       string
	Name        string
	ParentLMSID string
	ThirdID     int
	ClosedAt    int64
	Sort        int64
}

type dashboardAreaNodeKey struct {
	Type  int
	LMSID string
}

func partitionDashboardAreaStores(stores []enterpriseModel.StoreMasterModel, now time.Time) ([]enterpriseModel.StoreMasterModel, map[string]int64, map[string]struct{}) {
	hierarchyStores := make([]enterpriseModel.StoreMasterModel, 0, len(stores))
	closedAtByStore := make(map[string]int64)
	protectedStoreCodes := make(map[string]struct{})
	for _, store := range stores {
		status := strings.ToUpper(strings.TrimSpace(store.StoreStatus))
		storeCode := strings.TrimSpace(store.StoreCode)
		switch status {
		case "OPEN", "PENDING":
			if isDashboardDistrictCode(store.OrgDistrictCode) {
				hierarchyStores = append(hierarchyStores, store)
			}
		case "CLOSED":
			if isDashboardDistrictCode(store.OrgDistrictCode) {
				hierarchyStores = append(hierarchyStores, store)
			}
			if storeCode != "" {
				closedAtByStore[storeCode] = dashboardAreaClosedAt(store.StoreStatus, store.ClosedDate, now)
				protectedStoreCodes[storeCode] = struct{}{}
			}
		default:
			if storeCode != "" {
				protectedStoreCodes[storeCode] = struct{}{}
			}
		}
	}
	return hierarchyStores, closedAtByStore, protectedStoreCodes
}

func buildDashboardAreaNodes(stores []enterpriseModel.StoreMasterModel, storeThirdIDs map[string]int, now time.Time) ([]dashboardAreaNode, []string, error) {
	nodeMap := make(map[dashboardAreaNodeKey]dashboardAreaNode)
	skipped := make([]string, 0)
	for _, store := range stores {
		if !isDashboardDistrictCode(store.OrgDistrictCode) {
			continue
		}
		channelCode := strings.TrimSpace(store.OrgChannelCode)
		districtCode := strings.TrimSpace(store.OrgDistrictCode)
		areaCode := strings.TrimSpace(store.OrgAreaCode)
		storeCode := strings.TrimSpace(store.StoreCode)
		if invalidAreaCode(channelCode) || invalidAreaCode(districtCode) || invalidAreaCode(areaCode) || invalidAreaCode(storeCode) {
			skipped = append(skipped, storeCode)
			continue
		}

		storeNodes := []dashboardAreaNode{
			{Type: 1, LMSID: channelCode, Name: areaNodeName(store.SalesChannel, channelCode)},
			{Type: 2, LMSID: districtCode, Name: areaNodeName(store.OrgDistrictName, districtCode), ParentLMSID: channelCode},
			{Type: 3, LMSID: areaCode, Name: areaNodeName(store.OrgAreaName, areaCode), ParentLMSID: districtCode},
			{
				Type:        4,
				LMSID:       storeCode,
				Name:        storeCode,
				ParentLMSID: areaCode,
				ThirdID:     storeThirdIDs[storeCode],
				ClosedAt:    dashboardAreaClosedAt(store.StoreStatus, store.ClosedDate, now),
			},
		}
		for _, node := range storeNodes {
			if err := putDashboardAreaNode(nodeMap, node); err != nil {
				return nil, skipped, err
			}
		}
	}

	nodes := make([]dashboardAreaNode, 0, len(nodeMap))
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Type != nodes[j].Type {
			return nodes[i].Type < nodes[j].Type
		}
		if nodes[i].ParentLMSID != nodes[j].ParentLMSID {
			return nodes[i].ParentLMSID < nodes[j].ParentLMSID
		}
		return nodes[i].LMSID < nodes[j].LMSID
	})
	var previousParent string
	var previousType int
	var siblingSort int64
	for index := range nodes {
		if nodes[index].Type != previousType || nodes[index].ParentLMSID != previousParent {
			siblingSort = 0
		}
		siblingSort++
		nodes[index].Sort = siblingSort
		previousType = nodes[index].Type
		previousParent = nodes[index].ParentLMSID
	}
	sort.Strings(skipped)
	return nodes, skipped, nil
}

func isDashboardDistrictCode(districtCode string) bool {
	districtCode = strings.ToUpper(strings.TrimSpace(districtCode))
	_, exists := dashboardDistrictCodes[districtCode]
	return exists
}

func putDashboardAreaNode(nodes map[dashboardAreaNodeKey]dashboardAreaNode, node dashboardAreaNode) error {
	key := dashboardAreaNodeKey{Type: node.Type, LMSID: node.LMSID}
	existing, exists := nodes[key]
	if exists && existing != node {
		return fmt.Errorf("LMS组织节点存在冲突：type=%d,lms_id=%s", node.Type, node.LMSID)
	}
	nodes[key] = node
	return nil
}

func invalidAreaCode(value string) bool {
	return value == "" || strings.EqualFold(value, "N/A")
}

func areaNodeName(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "N/A") {
		return fallback
	}
	return value
}

func dashboardAreaClosedAt(status, closeDate string, now time.Time) int64 {
	status = strings.TrimSpace(status)
	if strings.EqualFold(status, "OPEN") || strings.EqualFold(status, "PENDING") {
		return 0
	}
	closeDate = strings.TrimSpace(closeDate)
	for _, layout := range []string{"01/02/2006", "2006/01/02", "2006-01-02", "01/02/2006 15:04:05", "2006/01/02 15:04:05", "2006-01-02 15:04:05"} {
		if parsed, err := time.ParseInLocation(layout, closeDate, now.Location()); err == nil {
			closedAt := parsed.Unix()
			if closedAt > 0 && closedAt <= maxDashboardAreaTimestamp {
				return closedAt
			}
			return now.Unix()
		}
	}
	return now.Unix()
}
