package micro_scrm_enterprise

type NewResultSignInDetail struct {
	Avatar           string             `json:"avatar"`
	Name             string             `json:"name"`
	Phone            string             `json:"phone"`
	Day              int64              `json:"day"`
	SignInInfoDetail []SignInInfoDetail `json:"sign_in_info_detail"`
	Count            int64              `json:"count"`
	ShrName          string             `json:"shr_name"`
	ShrPhone         string             `json:"shr_phone"`
	ShrAddress       string             `json:"shr_address"`
}
type SignInInfoDetail struct {
	Days        int    `json:"days"`
	SignTime    int64  `json:"sign_time"`
	SignContent string `json:"sign_content"`
	SignPic     string `json:"sign_pic"`
}

type SignInList struct {
	UnionID      string `json:"union_id"`
	Count        int    `json:"count"`
	LastSignTime int64  `json:"last_sign_time"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
}

type NewResultList struct {
	Count      int64         `json:"count"`
	SignInList []*SignInList `json:"sign_in_list"`
}

type UserInfo struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	ShrName    string `json:"shr_name"`
	ShrPhone   string `json:"shr_phone"`
	ShrAddress string `json:"shr_address"`
}

type AvatarInfo struct {
	Avatar string `json:"avatar"`
}
