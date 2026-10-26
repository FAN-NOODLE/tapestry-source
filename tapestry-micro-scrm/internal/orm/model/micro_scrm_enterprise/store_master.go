package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

// StoreMasterModel Coach LMS 门店主数据。
type StoreMasterModel struct {
	model.BranchModel
	ParentOrganizationCode    string `gorm:"column:parent_organization_code" json:"parent_organization_code"`
	PrimaryEnterpriseKey      string `gorm:"column:primary_enterprise_key" json:"primary_enterprise_key"`
	InventoryOrganizationCode string `gorm:"column:inventory_organization_code" json:"inventory_organization_code"`
	Operation                 string `gorm:"column:operation" json:"operation"`
	CountryCode               string `gorm:"column:country_code" json:"country_code"`
	RegionDesc                string `gorm:"column:region_desc" json:"region_desc"`
	OrganizationCode          string `gorm:"column:organization_code" json:"organization_code"`
	BrandCode                 string `gorm:"column:brand_code" json:"brand_code"`
	LocaleCode                string `gorm:"column:locale_code" json:"locale_code"`
	OrganizationName          string `gorm:"column:organization_name" json:"organization_name"`
	DayPhone                  string `gorm:"column:day_phone" json:"day_phone"`
	EmailID                   string `gorm:"column:email_id" json:"email_id"`
	AddressLine1              string `gorm:"column:address_line1" json:"address_line1"`
	AddressLine2              string `gorm:"column:address_line2" json:"address_line2"`
	AddressLine3              string `gorm:"column:address_line3" json:"address_line3"`
	AddressLine4              string `gorm:"column:address_line4" json:"address_line4"`
	ZipCode                   string `gorm:"column:zip_code" json:"zip_code"`
	City                      string `gorm:"column:city" json:"city"`
	State                     string `gorm:"column:state" json:"state"`
	AllowGiftWrap             string `gorm:"column:allow_gift_wrap" json:"allow_gift_wrap"`
	Longitude                 string `gorm:"column:longitude" json:"longitude"`
	Latitude                  string `gorm:"column:latitude" json:"latitude"`
	ActivateFlag              string `gorm:"column:activate_flag" json:"activate_flag"`
	NodeType                  string `gorm:"column:node_type" json:"node_type"`
	ShippingNode              string `gorm:"column:shipping_node" json:"shipping_node"`
	StoreRegionalMgr          string `gorm:"column:store_regional_mgr" json:"store_regional_mgr"`
	AreaManager               string `gorm:"column:area_manager" json:"area_manager"`
	StoreDistrictMgr          string `gorm:"column:store_district_mgr" json:"store_district_mgr"`
	IsStorePickUpSupported    string `gorm:"column:is_store_pick_up_supported" json:"is_store_pick_up_supported"`
	StoreWorkingHours         string `gorm:"column:store_working_hours" json:"store_working_hours"`
	IsNewStore                string `gorm:"column:is_new_store" json:"is_new_store"`
	Monogramming              string `gorm:"column:monogramming" json:"monogramming"`
	Customization             string `gorm:"column:customization" json:"customization"`
	StoreType                 string `gorm:"column:store_type" json:"store_type"`
	StoreCode                 string `gorm:"column:store_code" json:"store_code"`
	CS2000StoreID             string `gorm:"column:cs2000_store_id" json:"cs2000_store_id"`
	VirtualStoreID            string `gorm:"column:virtual_store_id" json:"virtual_store_id"`
	CompanyCode               string `gorm:"column:company_code" json:"company_code"`
	VirtualCompanyCode        string `gorm:"column:virtual_company_code" json:"virtual_company_code"`
	SalesOrgCode              string `gorm:"column:sales_org_code" json:"sales_org_code"`
	StoreClass                string `gorm:"column:store_class" json:"store_class"`
	StoreStatus               string `gorm:"column:store_status" json:"store_status"`
	OpenDate                  string `gorm:"column:open_date" json:"open_date"`
	ClosedDate                string `gorm:"column:closed_date" json:"closed_date"`
	SalesChannel              string `gorm:"column:sales_channel" json:"sales_channel"`
	StoreFormat               string `gorm:"column:store_format" json:"store_format"`
	OrgChannelCode            string `gorm:"column:org_channel_code" json:"org_channel_code"`
	OrgChannelName            string `gorm:"column:org_channel_name" json:"org_channel_name"`
	OrgRegionCode             string `gorm:"column:org_region_code" json:"org_region_code"`
	OrgRegionName             string `gorm:"column:org_region_name" json:"org_region_name"`
	OrgDistrictCode           string `gorm:"column:org_district_code" json:"org_district_code"`
	OrgDistrictName           string `gorm:"column:org_district_name" json:"org_district_name"`
	OrgAreaCode               string `gorm:"column:org_area_code" json:"org_area_code"`
	OrgAreaName               string `gorm:"column:org_area_name" json:"org_area_name"`
	StoreNameZh               string `gorm:"column:store_name_zh" json:"store_name_zh"`
	AddressZh                 string `gorm:"column:address_zh" json:"address_zh"`
	StoreArea                 string `gorm:"column:store_area" json:"store_area"`
	AppointmentServiceTypes   string `gorm:"column:appointment_service_types" json:"appointment_service_types"`
	PaymentMethods            string `gorm:"column:payment_methods" json:"payment_methods"`
	StoreDesignConcept        string `gorm:"column:store_design_concept" json:"store_design_concept"`
	DPAP                      string `gorm:"column:dp_ap" json:"dp_ap"`
}

func (StoreMasterModel) TableName() string { return "store_master" }
