package index

// 首页公告
type NoticeParam struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

// 综合指标查询
type ShopFeeParam struct {
	Day          int64 `json:"day" desc:"0 今日 1昨天 近n天"`
	Type         int64 `json:"type"`
	StoreId      int64 `json:"store_id"`
	EnterpriseId int64 `json:"enterprise_id"`
}

type ShopFeeBody struct {
	Indices  []string `json:"indices" desc:"指标"`
	Sign     string   `json:"sign"`
	StoreId  int64    `json:"store_id"`
	SignTime int64    `json:"sign_time"`
}

// 商城变现响应
type ShopFeeRes struct {
	Data []struct {
		SsFee   float64 `json:"ss_fee" desc:"实收金额"`
		ReFee   float64 `json:"re_fee" desc:"退款金额"`
		Uv      int64   `json:"uv" desc:"访客数"`
		OrderUv int64   `json:"order_uv" desc:"下单买家数"`
		PayUv   int64   `json:"pay_uv" desc:"支付买家数"`
		PayNum  int64   `json:"pay_num" desc:"支付订单数"`
	} `json:"data"`
}

// 商城变现真实响应
type ShopFeeRealRes struct {
	Uv           int64   `json:"uv" desc:"访客数"`
	SsFee        float64 `json:"ss_fee" desc:"实收金额"`
	ReFee        float64 `json:"re_fee" desc:"退款金额"`
	OrderUv      int64   `json:"order_uv" desc:"下单买家数"`
	PayUv        int64   `json:"pay_uv" desc:"支付买家数"`
	PayNum       int64   `json:"pay_num" desc:"支付订单数"`
	UnitPrice    float64 `json:"unit_price" desc:"客单价-实收金额/支付买家数"`
	UvOrderThan  string  `json:"uv_order_than" desc:"访客下单转化-下单买家数/访客数*100%"`
	OrderPayThan string  `json:"order_pay_than" desc:"下单支付转化-支付买家数/下单买家数*100%"`
	UvPayThan    string  `json:"uv_pay_than" desc:"访客支付转化-支付买家数/访客数*100%"`
}

// 客户总转化响应
type ShopConversion struct {
	Data []struct {
		MinTime int64   `json:"min_time" desc:"起始时间"`
		AllUv   int64   `json:"all_uv" desc:"总访客数"`
		PaysUv  int64   `json:"pays_uv" desc:"总支付人数"`
		FgUv    int64   `json:"fg_uv" desc:"总复购人数"`
		FgTime  float64 `json:"fg_time" desc:"复购周期"`
		ZhTime  float64 `json:"zh_time" desc:"转化周期"`
	} `json:"data"`
}

type ShopConversionRealRes struct {
	MinTime       int64   `json:"min_time" desc:"起始时间"`
	AllUv         int64   `json:"all_uv" desc:"总访客数"`
	PaysUv        int64   `json:"pays_uv" desc:"总支付人数"`
	FgUv          int64   `json:"fg_uv" desc:"总复购人数"`
	FgTime        float64 `json:"fg_time" desc:"复购周期"`
	ZhTime        float64 `json:"zh_time" desc:"转化周期"`
	UvPayThan     string  `json:"uv_pay_than" desc:"访问支付转化率-总支付买家数/总访客数*100%"`
	RepeatPayThan string  `json:"repeat_pay_than" desc:"复购支付转化率-总复购客户/总支付买家数*100%"`
}

// 企业好友和群成员响应
type Friend struct {
	ReUv       int64 `json:"re_uv" desc:"留存好友数"`
	ReUvThan   int64 `json:"re_uv_than" desc:"留存好友环比数"`
	AddUv      int64 `json:"add_uv" desc:"新增好友数"`
	AddUvThan  int64 `json:"add_uv_than" desc:"新增好友环比数"`
	LossUv     int64 `json:"loss_uv" desc:"流失好友数"`
	LossUvThan int64 `json:"loss_uv_than" desc:"流失好友环比数"`
}

type Group struct {
	ReGroupUv       int64 `json:"re_group_uv" desc:"留存群成员数"`
	ReGroupUvThan   int64 `json:"re_group_uv_than" desc:"留存群成员环比数"`
	AddGroupUv      int64 `json:"add_group_uv" desc:"新增群成员数"`
	AddGroupUvThan  int64 `json:"add_group_uv_than" desc:"新增群成员环比数"`
	LossGroupUv     int64 `json:"loss_group_uv" desc:"流失群成员"`
	LossGroupUvThan int64 `json:"loss_group_uv_than" desc:"流失群环比成员"`
	GroupNum        int64 `json:"group_num" desc:"群数"`
	GroupNumThan    int64 `json:"group_num_than" desc:"群数环比"`
}

// 公众号粉丝相关
type FansParam struct {
	Day          int64 `json:"day" desc:"0 今日 1昨天 近n天"`
	EnterpriseId int64 `json:"enterprise_id"`
}

type FansBody struct {
	Sign     string `json:"sign"`
	SignTime int64  `json:"sign_time"`
}

// 粉丝
type FansCDPRes struct {
	Data []struct {
		Appid string `json:"appid"`
		Name  string `json:"name"`
		V     struct {
			Now []float64 `json:"now"`
			New []float64 `json:"new"`
			Del []float64 `json:"del"`
		} `json:"v"`
	}
}

type FansRealCDPRes []struct {
	Appid string `json:"appid"`
	Name  string `json:"name"`
	V     struct {
		Now []float64 `json:"now"`
		New []float64 `json:"new"`
		Del []float64 `json:"del"`
	} `json:"v"`
}

type Fans struct {
	Appid string `json:"appid"`
	Name  string `json:"name"`
	V     struct {
		Now struct {
			Uv        float64 `json:"uv"`
			UvThan    float64 `json:"uv_than"`
			Direction int64   `json:"direction"`
			Percent   string  `json:"percent"`
		} `json:"now"`
		New struct {
			Uv        float64 `json:"uv"`
			UvThan    float64 `json:"uv_than"`
			Direction int64   `json:"direction"`
			Percent   string  `json:"percent"`
		} `json:"new"`
		Del struct {
			Uv        float64 `json:"uv"`
			UvThan    float64 `json:"uv_than"`
			Direction int64   `json:"direction"`
			Percent   string  `json:"percent"`
		} `json:"del"`
	} `json:"v"`
}

type PrivateDomain struct {
	Fans   []Fans `json:"fans"`
	Friend struct {
		ReUv            int64  `json:"re_uv" desc:"留存好友数"`
		ReUvThan        int64  `json:"re_uv_than" desc:"留存好友数"`
		ReUvDirection   int64  `json:"re_uv_direction" desc:"留存好友增减"`
		ReUvPercent     string `json:"re_uv_percent" desc:"留存好友百分比"`
		AddUv           int64  `json:"add_uv" desc:"新增好友数"`
		AddUvThan       int64  `json:"add_uv_than" desc:"新增好友数"`
		AddUvDirection  int64  `json:"add_uv_direction" desc:"新增好友增减"`
		AddUvPercent    string `json:"add_uv_percent" desc:"新增好友百分比"`
		LossUv          int64  `json:"loss_uv" desc:"流失好友数"`
		LossUvThan      int64  `json:"loss_uv_than" desc:"流失好友数"`
		LossUvDirection int64  `json:"loss_uv_direction" desc:"流失好友数增减"`
		LossUvPercent   string `json:"loss_uv_percent" desc:"流失好友数百分比"`
	} `json:"friend"`
	Group struct {
		ReGroupUv            int64  `json:"re_group_uv" desc:"留存群成员数"`
		ReGroupUvThan        int64  `json:"re_group_uv_than" desc:"留存群成员数"`
		ReGroupUvDirection   int64  `json:"re_group_uv_direction" desc:"留存群成员增减"`
		ReGroupUvPercent     string `json:"re_group_uv_percent" desc:"留存群成员百分比"`
		AddGroupUv           int64  `json:"add_group_uv" desc:"新增群成员数"`
		AddGroupUvThan       int64  `json:"add_group_uv_than" desc:"新增群成员数"`
		AddGroupUvDirection  int64  `json:"add_group_uv_direction" desc:"新增群成员数增减"`
		AddGroupUvPercent    string `json:"add_group_uv_percent" desc:"新增群成员数百分比"`
		LossGroupUv          int64  `json:"loss_group_uv" desc:"流失群成员"`
		LossGroupUvThan      int64  `json:"loss_group_uv_than" desc:"流失群成员"`
		LossGroupUvDirection int64  `json:"loss_group_uv_direction" desc:"流失好友数增减"`
		LossGroupUvPercent   string `json:"loss_group_uv_percent" desc:"流失好友数百分比"`
		GroupNum             int64  `json:"group_num" desc:"群数"`
		GroupNumThan         int64  `json:"group_num_than" desc:"群数"`
		GroupNumUvDirection  int64  `json:"group_num_uv_direction" desc:"群数增减"`
		GroupNumUvPercent    string `json:"group_num_uv_percent" desc:"群数百分比"`
	} `json:"group"`
}

// 导购相关
type ShopperListParam struct {
	YjSort       int64 `json:"yj_sort" desc:"排序"`
	Day          int64 `json:"day" desc:"0 今日 1昨天 近n天"`
	EnterpriseId int64 `json:"enterprise_id"`
}

type ShopperListBody struct {
	Indices  []string `json:"indices" desc:"指标"`
	Sign     string   `json:"sign"`
	SignTime int64    `json:"sign_time"`
}

type ShopperListRes struct {
	Data []map[int64]struct {
		Id           int64   `json:"id"`
		YjFee        float64 `json:"yj_fee"`
		UV           float64 `json:"uv"`
		UV2          float64 `json:"uv2"`
		QwUv         float64 `json:"qw_uv"`
		QwLayUv      float64 `json:"qw_lay_uv"`
		CouponFee    float64 `json:"coupon_fee"`
		OrderFee     float64 `json:"order_fee"`
		Name         string  `json:"name"`
		Contribution float64 `json:"contribution" desc:"贡献占比-导购业绩 / total里面的数字"`
		Discount     float64 `json:"discount" desc:"折扣率-优惠券面额 / 含券支付金额"`
	} `json:"data"`
	Total struct {
		YjFee        float64 `json:"yj_fee"`
		UV           float64 `json:"uv"`
		UV2          float64 `json:"uv2"`
		QwUv         float64 `json:"qw_uv"`
		QwLayUv      float64 `json:"qw_lay_uv"`
		CouponFee    float64 `json:"coupon_fee"`
		OrderFee     float64 `json:"order_fee"`
		Contribution string  `json:"contribution" desc:"贡献占比-导购业绩 / total里面的数字"`
		Discount     string  `json:"discount" desc:"折扣率-优惠券面额 / 含券支付金额"`
	} `json:"total"`
}

type ShopperListRealRes struct {
	List []struct {
		Id           int64   `json:"id"`
		YjFee        float64 `json:"yj_fee"`
		UV           float64 `json:"uv"`
		UV2          float64 `json:"uv2"`
		QwUv         float64 `json:"qw_uv"`
		QwLayUv      float64 `json:"qw_lay_uv"`
		CouponFee    float64 `json:"coupon_fee"`
		OrderFee     float64 `json:"order_fee"`
		Name         string  `json:"name"`
		Contribution float64 `json:"contribution" desc:"贡献占比-导购业绩 / total里面的数字"`
		Discount     float64 `json:"discount" desc:"折扣率-优惠券面额 / 含券支付金额"`
	} `json:"list"`
	Total struct {
		YjFee        float64 `json:"yj_fee"`
		UV           float64 `json:"uv"`
		UV2          float64 `json:"uv2"`
		QwUv         float64 `json:"qw_uv"`
		QwLayUv      float64 `json:"qw_lay_uv"`
		CouponFee    float64 `json:"coupon_fee"`
		OrderFee     float64 `json:"order_fee"`
		Contribution string  `json:"contribution" desc:"贡献占比-导购业绩 / total里面的数字"`
		Discount     string  `json:"discount" desc:"折扣率-优惠券面额 / 含券支付金额"`
	} `json:"total"`
}

// 门店列表
type StoreListRes struct {
	Data []struct {
		Id        int64  `json:"id"`
		ShortName string `json:"short_name"`
		IsServer  int64  `json:"is_server"`
	} `json:"data"`
}

// 基础数据检查
type BaseInfoCheck struct {
	Data struct {
		PaymentConfig int64 `json:"payment_config"`
		QwConfig      int64 `json:"qw"`
		ShopAuth      int64 `json:"shop_auth"`
	} `json:"data"`
}

type HelpListRes struct {
	Data struct {
		List []struct {
			Id      int64  `json:"id"`
			Key     string `json:"key"`
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"list"`
	} `json:"data"`
}

func (h HelpListRes) CheckCode() error {
	return nil
}

type NoticeHelpRes struct {
	Data struct {
		List []struct {
			Id      int64  `json:"id"`
			Route   string `json:"route"`
			Content string `json:"content"`
		}
	} `json:"data"`
}
