package openDao

import (
	"fmt"
	"sort"

	"micro-scrm/internal/app"
	openReq "micro-scrm/internal/app/http/request/open"
	"micro-scrm/internal/orm/model"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type EmployeeMasterRow struct {
	enterpriseModel.EmployeeMasterModel
	OrganizationCode string `gorm:"column:organization_code"`
}

type EmployeeWeComData struct {
	WecomUserid        string
	WecomJoinedAt      int64
	VerificationStatus int
	WecomLeftAt        int64
	DepartmentChanges  []EmployeeDepartmentChange
}

type EmployeeDepartmentChange struct {
	EventType          string
	EventTime          int64
	WechatDepartmentID int64
	DepartmentName     string
}

type employeeMemberRow struct {
	ID         int64  `gorm:"column:id"`
	Userid     string `gorm:"column:userid"`
	OpenUserid string `gorm:"column:open_userid"`
	Status     int    `gorm:"column:status"`
	CreatedAt  int64  `gorm:"column:created_at"`
	DeletedAt  int64  `gorm:"column:deleted_at"`
}

type employeeDepartmentRow struct {
	MemberID           int64  `gorm:"column:member_id"`
	CreatedAt          int64  `gorm:"column:created_at"`
	UpdatedAt          int64  `gorm:"column:updated_at"`
	DeletedAt          int64  `gorm:"column:deleted_at"`
	WechatDepartmentID int64  `gorm:"column:wechat_department_id"`
	DepartmentName     string `gorm:"column:department_name"`
}

type StoreOperationDao struct {
	db     *gorm.DB
	workDB *gorm.DB
}

func NewStoreOperationDao() *StoreOperationDao {
	return &StoreOperationDao{db: gormx.GetInstance(), workDB: gormx.GetWorkInstance()}
}

func (d *StoreOperationDao) StoreMasterPage(c app.Ctx, enterpriseID int64, req openReq.StoreMasterListReq) ([]enterpriseModel.StoreMasterModel, int64, error) {
	var result []enterpriseModel.StoreMasterModel
	var total int64
	query := buildStoreMasterQuery(d.db.WithContext(c), enterpriseID, req)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return result, 0, nil
	}
	offset := (req.Page - 1) * req.Size
	err := query.Order("updated_at ASC, id ASC").Offset(offset).Limit(req.Size).Find(&result).Error
	return result, total, err
}

func buildStoreMasterQuery(db *gorm.DB, enterpriseID int64, req openReq.StoreMasterListReq) *gorm.DB {
	query := db.Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
		Where("deleted_at = ?", 0).
		Where("country_code IN ?", req.StoreCountryCodes)
	if req.UpdatedAt > 0 {
		query = query.Where("updated_at >= ?", req.UpdatedAt)
	}
	return query
}

func (d *StoreOperationDao) WechatDepartmentIDs(c app.Ctx, enterpriseID int64, storeCodes []string) (map[string]int64, error) {
	result := make(map[string]int64, len(storeCodes))
	if len(storeCodes) == 0 {
		return result, nil
	}
	type storeDepartment struct {
		StoreCode        string `gorm:"column:third_store_code"`
		WechatDepartment uint   `gorm:"column:wechat_department_id"`
	}
	var rows []storeDepartment
	err := d.workDB.WithContext(c).Table("store_department").
		Select("third_store_code, wechat_department_id").
		Where("enterprise_id = ? AND deleted_at = ?", enterpriseID, 0).
		Where("third_store_code IN ?", storeCodes).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.StoreCode] = int64(row.WechatDepartment)
	}
	return result, nil
}

func (d *StoreOperationDao) EmployeeMasterPage(c app.Ctx, enterpriseID int64, req openReq.EmployeeMasterListReq) ([]EmployeeMasterRow, int64, error) {
	var result []EmployeeMasterRow
	var total int64
	query := buildEmployeeMasterQuery(d.db.WithContext(c), enterpriseID, req)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return result, 0, nil
	}
	offset := (req.Page - 1) * req.Size
	err := query.Select("em.*, sm.organization_code AS organization_code").
		Order("em.updated_at ASC, em.id ASC").Offset(offset).Limit(req.Size).Scan(&result).Error
	return result, total, err
}

func buildEmployeeMasterQuery(db *gorm.DB, enterpriseID int64, req openReq.EmployeeMasterListReq) *gorm.DB {
	employeeModel := enterpriseModel.EmployeeMasterModel{}
	storeModel := enterpriseModel.StoreMasterModel{}
	employeeTable := fmt.Sprintf("%s%d.%s", employeeModel.DBPrefix(), enterpriseID, employeeModel.TableName())
	storeTable := fmt.Sprintf("%s%d.%s", storeModel.DBPrefix(), enterpriseID, storeModel.TableName())
	query := db.Table(employeeTable+" AS em").
		Joins("LEFT JOIN "+storeTable+" AS sm ON CONVERT(sm.store_code USING utf8mb4) COLLATE utf8mb4_unicode_ci = CONVERT(em.store_code USING utf8mb4) COLLATE utf8mb4_unicode_ci AND sm.deleted_at = 0").
		Where("em.deleted_at = ?", 0)
	switch req.EmploymentScope {
	case "active":
		query = query.Where("em.status = ?", "1")
	case "inactive":
		query = query.Where("em.status <> ?", "1")
	}
	switch req.StoreScope {
	case "store":
		query = query.Where("sm.id IS NOT NULL")
	case "open_store":
		query = query.Where("sm.id IS NOT NULL AND sm.store_status = ?", "OPEN")
	}
	if len(req.StoreCountryCodes) > 0 {
		query = query.Where("sm.country_code IN ?", req.StoreCountryCodes)
	}
	if req.UpdatedAt > 0 {
		query = query.Where("em.updated_at >= ?", req.UpdatedAt)
	}
	return query
}

func (d *StoreOperationDao) EmployeeWeComData(c app.Ctx, enterpriseID int64, employeeCodes []string) (map[string]EmployeeWeComData, error) {
	result := make(map[string]EmployeeWeComData, len(employeeCodes))
	if len(employeeCodes) == 0 {
		return result, nil
	}
	var members []employeeMemberRow
	err := d.db.WithContext(c).Table("members").
		Select("id, userid, open_userid, status, created_at, deleted_at").
		Where("enterprise_id = ?", enterpriseID).
		Where("(userid IN ? OR open_userid IN ?)", employeeCodes, employeeCodes).
		Order("id ASC").Find(&members).Error
	if err != nil {
		return nil, err
	}
	matched := matchEmployeeMembers(employeeCodes, members)
	if len(matched) == 0 {
		return result, nil
	}
	memberIDs := make([]int64, 0, len(matched))
	for _, member := range matched {
		memberIDs = append(memberIDs, member.ID)
	}
	var departmentRows []employeeDepartmentRow
	err = d.db.WithContext(c).Table("member_departments AS md").
		Select("md.member_id, md.created_at, md.updated_at, md.deleted_at, COALESCE(d.third_id, 0) AS wechat_department_id, COALESCE(d.name, '') AS department_name").
		Joins("LEFT JOIN departments AS d ON d.id = md.department_id AND d.enterprise_id = ?", enterpriseID).
		Where("md.enterprise_id = ? AND md.member_id IN ?", enterpriseID, memberIDs).
		Find(&departmentRows).Error
	if err != nil {
		return nil, err
	}
	return assembleEmployeeWeComData(matched, departmentRows), nil
}

func assembleEmployeeWeComData(matched map[string]employeeMemberRow, departmentRows []employeeDepartmentRow) map[string]EmployeeWeComData {
	result := make(map[string]EmployeeWeComData, len(matched))
	changesByMember := buildDepartmentChanges(departmentRows)
	for employeeCode, member := range matched {
		verificationStatus := member.Status
		wecomLeftAt := int64(0)
		if member.DeletedAt > 0 {
			verificationStatus = 5
			wecomLeftAt = member.DeletedAt
		}
		changes := changesByMember[member.ID]
		if changes == nil {
			changes = make([]EmployeeDepartmentChange, 0)
		}
		result[employeeCode] = EmployeeWeComData{
			WecomUserid:        member.Userid,
			WecomJoinedAt:      member.CreatedAt,
			VerificationStatus: verificationStatus,
			WecomLeftAt:        wecomLeftAt,
			DepartmentChanges:  changes,
		}
	}
	return result
}

func matchEmployeeMembers(employeeCodes []string, members []employeeMemberRow) map[string]employeeMemberRow {
	result := make(map[string]employeeMemberRow, len(employeeCodes))
	employeeCodeSet := make(map[string]struct{}, len(employeeCodes))
	for _, employeeCode := range employeeCodes {
		employeeCodeSet[employeeCode] = struct{}{}
	}
	for _, member := range members {
		if _, wanted := employeeCodeSet[member.Userid]; wanted {
			if _, exists := result[member.Userid]; !exists {
				result[member.Userid] = member
			}
		}
	}
	for _, member := range members {
		if _, wanted := employeeCodeSet[member.OpenUserid]; !wanted {
			continue
		}
		if _, exists := result[member.OpenUserid]; !exists {
			result[member.OpenUserid] = member
		}
	}
	return result
}

func buildDepartmentChanges(rows []employeeDepartmentRow) map[int64][]EmployeeDepartmentChange {
	result := make(map[int64][]EmployeeDepartmentChange)
	for _, row := range rows {
		if row.CreatedAt > 0 {
			result[row.MemberID] = append(result[row.MemberID], employeeDepartmentChange(row, "join", row.CreatedAt))
		}
		if row.UpdatedAt > row.CreatedAt {
			result[row.MemberID] = append(result[row.MemberID], employeeDepartmentChange(row, "join", row.UpdatedAt))
		}
		if row.DeletedAt > 0 {
			result[row.MemberID] = append(result[row.MemberID], employeeDepartmentChange(row, "leave", row.DeletedAt))
		}
	}
	for memberID := range result {
		sort.SliceStable(result[memberID], func(i, j int) bool {
			left, right := result[memberID][i], result[memberID][j]
			if left.EventTime != right.EventTime {
				return left.EventTime < right.EventTime
			}
			if left.EventType != right.EventType {
				return left.EventType == "join"
			}
			return left.WechatDepartmentID < right.WechatDepartmentID
		})
	}
	return result
}

func employeeDepartmentChange(row employeeDepartmentRow, eventType string, eventTime int64) EmployeeDepartmentChange {
	return EmployeeDepartmentChange{
		EventType:          eventType,
		EventTime:          eventTime,
		WechatDepartmentID: row.WechatDepartmentID,
		DepartmentName:     row.DepartmentName,
	}
}
