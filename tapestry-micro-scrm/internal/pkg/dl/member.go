package dl

import (
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type MemberQueryInfoRequest struct {
	//BaseMemberQuery
	EnterpriseID string `json:"enterprise_id"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}

type MemberQueryInfoResponse struct {
	BaseResponse
	Data struct {
		Address []struct {
			AddressType  string `json:"address_type"`
			CityCode     string `json:"cityCode"`
			ContactName  string `json:"contactName"`
			CountryCode  string `json:"countryCode"`
			IsDefault    string `json:"isDefault"`
			Location     string `json:"location"`
			MemberCode   string `json:"memberCode"`
			Phone        string `json:"phone"`
			ProgramCode  string `json:"programCode"`
			ProvinceCode string `json:"province_code"`
			SeqNum       string `json:"seqNum"`
			ZipCode      string `json:"zipcode"`
		} `json:"address"`
		Certificate []struct {
			CardNo string `json:"cardNo"`
			Type   string `json:"type"`
		} `json:"certificate"`
		MediaAccount []struct {
			AccountNo      string `json:"accountNo"`
			BindDate       string `json:"bindDate"`
			BindStatus     string `json:"bindStatus"`
			ConsentStatus  string `json:"consentStatus"`
			ConsentTime    string `json:"consentTime"`
			FollowDate     string `json:"followDate"`
			FollowStatus   string `json:"followStatus"`
			FromSocialCode string `json:"fromSocialCode"`
			MemberCode     string `json:"memberCode"`
			Nickname       string `json:"nickname"`
			Remark         string `json:"remark"`
			SeqNum         string `json:"seqNum"`
			Status         string `json:"status"`
			Type           string `json:"type"`
			UnbindDate     string `json:"unbindDate"`
			UnfollowDate   string `json:"unfollowDate"`
			UnionID        string `json:"unionId"`
		} `json:"media_account"`
		ChannelCode    string `json:"chanel_code"` // DL 的文档就是这个字段
		ConsentStatus  string `json:"consentStatus"`
		ConsentTime    string `json:"consentTime"`
		Email          string `json:"email"`
		MarsMemberNum  string `json:"mars_member_num"`
		AdminStoreCode string `json:"admin_store_code"`
		BaCode         string `json:"ba_code"`
		Birthday       string `json:"birthday"`
		BrandCode      string `json:"brand_code"`
		CardNo         string `json:"cardno"`
		Password       string `json:"password"`
		Mobile         string `json:"mobile"`
		Name           string `json:"name"`
		OfficePhone    string `json:"officePhone"`
		ProgramCode    string `json:"program_code"`
		RBMember       string `json:"rB_Member"`
		Sex            string `json:"sex"`
		Status         string `json:"status"`
		StoreCode      string `json:"store_code"`
		UnionCode      string `json:"union_code"`
		PrmtPerson     string `json:"prmtPerson"`
		EnrollmentDate string `json:"enrollment_date"`
	} `json:"data"`
}

// MemberQueryInfo 会员信息查询
func (dl *DataLake) MemberQueryInfo(ctx ctx.Context, request MemberQueryInfoRequest) (data MemberQueryInfoResponse, err error) {
	url := fmt.Sprintf("%s%s", dl.prefixURL, MEMBER_INFO_QUERY_URL)
	resp, err := dl.httpClient.Post(ctx, url, clients.JSONBody(request))

	err = dl.makeResponseData(resp, &data)

	return
}

type MemberQueryGradeRequest struct {
	//BaseMemberQuery
	//Category string `json:"category"`
	EnterpriseID int64  `json:"enterprise_id"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}

type BAListRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	RankLength   int    `json:"rank_length"`
	QyDomain     string `json:"qy_domain"`
	AdminId      string `json:"admin_id"`
	Type         int    `json:"type"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
	RankColumn   string `json:"rank_column"`
}

type BAListRespone struct {
	AreaLevelconfig int      `json:"area_level_config"`
	List            []BAinfo `json:"list"`
}
type BAinfo struct {
	AreaLevelList []string `json:"area_level_list"`
	DpName        string   `json:"dp_name"`
	DpNamePath    string   `json:"dp_name_path"`
	MemUID        string   `json:"mem_uid"`
	WorkerNo      string   `json:"worker_no"`
	Avatar        string   `json:"avatar"`
	Num           int64    `json:"num"`
	Name          string   `json:"name"`
	WorkerId      int64    `json:"worker_id"`
}

type BAListResponeExt struct {
	AreaLevelconfig int         `json:"area_level_config"`
	List            []BAinfoExt `json:"list"`
}
type BAinfoExt struct {
	BAinfo
	MemUserID string `json:"mem_userid"`
}

type MemberQueryGradeResponse struct {
	BaseResponse
	Data struct {
		Grade          string `json:"grade"`
		GradeCode      string `json:"grade_code"`
		GradeEndDate   string `json:"grade_end_date"`
		GradeStartDate string `json:"grade_start_date"`
	} `json:"data"`
}

// MemberQueryGrade 会员等级信息查询
func (dl *DataLake) MemberQueryGrade(ctx ctx.Context, request MemberQueryGradeRequest) (data MemberQueryGradeResponse, err error) {
	url := fmt.Sprintf("%s%s", dl.prefixURL, MEMBER_GRADE_QUERY_URL)
	resp, err := dl.httpClient.Post(ctx, url, clients.JSONBody(request))
	if err != nil {
		log.Error(ctx).Interface("MemberQueryGrade", err).Msgf("MemberQueryGrade %+v", request)
	}
	err = dl.makeResponseData(resp, &data)

	return
}

// FindBaList 查询BA信息
func (dl *DataLake) FindBaList(ctx ctx.Context, request BAListRequest) (data BAListRespone, err error) {
	url := fmt.Sprintf("%s%s", dl.prefixURL, BA_INFO_URL)
	resp, err := dl.httpClient.Post(ctx, url, clients.JSONBody(request))

	err = dl.makeResponseData(resp, &data)

	return
}

// FindBaListWithUserId 查询BA信息
func (dl *DataLake) FindBaListWithUserId(ctx ctx.Context, request BAListRequest) (data BAListResponeExt, err error) {
	url := fmt.Sprintf("%s%s", dl.prefixURL, BA_INFO_URL)
	resp, err := dl.httpClient.Post(ctx, url, clients.JSONBody(request))

	err = dl.makeResponseData(resp, &data)

	return
}
