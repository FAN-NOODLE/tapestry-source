package fission

type UpdateVipFlag struct {
	ID      int64  `gorm:"column:id"`
	UnionID string `gorm:"column:union_id"`
}

type FissionStyleConfig struct {
	ID        int64  `json:"id"`
	RgbColors string `json:"rgb_colors"`
	Styles    string `json:"styles"`
	Is_Check  int64  `json:"is_check"`
}

type FissionBgConfig struct {
	ID       int64  `json:"id"`
	BgImg    string `json:"bg_img"`
	ShareImg string `json:"share_img"`
}
