package baozun

type TokenReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// TokenResponse   三方商品图片鉴权固定返回
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Response    三方商品图片接口返回结果固定返回
type Response struct {
	AccessToken string `json:"access_token"`
}

type SearchParams struct {
	EnterpriseId            int64    `json:"enterprise_id"`
	Upc                     string   `json:"upc"`
	UpcSlice                []string `json:"upc_slice"`
	ProductName             string   `json:"product_name"`
	ProductNo               string   `json:"product_no"`
	SkuColor                string   `json:"sku_color"`
	StoreCode               string   `json:"store_code"`
	IsLike                  bool     `json:"is_like"`
	Size                    int      `json:"size"`
	LastCode                string   `json:"last_code"`
	ProductCode             []string `json:"product_code"`
	SkuItem                 string   `json:"sku_item"`
	IncludeNoneImageProduct bool     `json:"include_none_image_product"`
}

type ParamsUpcReq struct {
	Filter struct {
		IncludeNoneImageProduct bool `json:"includeNoneImageProduct"`
		SkfParam                struct {
			Properties struct {
				Upc string `json:"Upc,omitempty"`
			} `json:"properties"`
			VariantCode string `json:"variantCode,omitempty"`
		} `json:"skfParam"`
		ResourceParam struct {
			Name string `json:"name"`
		} `json:"resourceParam"`
	} `json:"filter"`
	ShopCode string `json:"shopCode,omitempty"`
	Pages    struct {
		PageSize int `json:"pageSize,omitempty"`
	} `json:"pages"`
	Sort struct {
		SortType string `json:"sortType"`
		Column   string `json:"column"`
	} `json:"sort"`
}

type ParamsSpuReq struct {
	Filter struct {
		IncludeNoneImageProduct bool `json:"includeNoneImageProduct"`
		SpuParam                struct {
			Properties struct {
				ProductNo   string `json:"productNo,omitempty"`
				ProductName string `json:"productName,omitempty"`
			} `json:"properties"`
			ProductCode string `json:"productCode,omitempty"`
		} `json:"spuParam"`
		ResourceParam struct {
			Name string `json:"name"`
		} `json:"resourceParam"`
	} `json:"filter"`
	ShopCode string `json:"shopCode,omitempty"`
	Pages    struct {
		PageSize int    `json:"pageSize,omitempty"`
		LastCode string `json:"lastCode,omitempty"`
	} `json:"pages"`
	Sort struct {
		SortType string `json:"sortType"`
		Column   string `json:"column"`
	} `json:"sort"`
}

type DataResponse struct {
	Data []struct {
		ExtData struct {
			Resource map[string]struct {
				Name string `json:"name"`
				Path string `json:"path"`
				Type int    `json:"type"`
				Url  string `json:"url"`
				Meta struct {
					Extension     string  `json:"extension"`
					ArtNo         string  `json:"artNo"`
					FullName      string  `json:"fullName"`
					Type          string  `json:"type"`
					ExtName       string  `json:"extName"`
					HtmlMediaType string  `json:"htmlMediaType"`
					ScaleX        int     `json:"scaleX"`
					ScaleY        int     `json:"scaleY"`
					Size          float64 `json:"size"`
					RelativePath  string  `json:"relativePath"`
					Name          string  `json:"name"`
					Width         int     `json:"width"`
					ColorNo       string  `json:"colorNo"`
					FileType      string  `json:"fileType"`
					Height        int     `json:"height"`
				} `json:"meta"`
				Properties struct {
					OriginalFileName string `json:"originalFileName"`
					OpDomains        []struct {
						Code string `json:"code"`
						Name string `json:"name"`
					} `json:"opDomains"`
					FileName        string `json:"fileName"`
					StorageProvider string `json:"storageProvider"`
					Creator         struct {
						Code string `json:"code"`
						Name string `json:"name"`
					} `json:"creator"`
					AssetKind        string `json:"assetKind"`
					LastModifiedDate string `json:"lastModifiedDate"`
					ArtNo            string `json:"artNo"`
					Format           string `json:"format"`
					StoragePath      string `json:"storagePath"`
					Source           string `json:"source"`
					Shares           []struct {
						Code     string `json:"code"`
						ShopName string `json:"shopName"`
						ShopId   string `json:"shopId"`
					} `json:"shares"`
					SortNo    string `json:"sortNo"`
					CreatedAt string `json:"createdAt"`
					Size      string `json:"size"`
					Supplier  string `json:"supplier"`
					Etag      string `json:"etag"`
					Id        string `json:"id"`
					ColorNo   string `json:"colorNo"`
					Dimension struct {
						X string `json:"x"`
						Y string `json:"y"`
					} `json:"dimension"`
					StorageUrl string `json:"storageUrl"`
					Tasks      string `json:"tasks"`
					Properties string `json:"properties"`
					Status     string `json:"status"`
				} `json:"properties"`
			} `json:"resource"`
			Prices []struct {
				VariantCode string `json:"variantCode"`
				Prices      []struct {
					Key      string  `json:"key"`
					Value    float64 `json:"value"`
					ShopCode string  `json:"shopCode"`
				} `json:"prices"`
			} `json:"prices"`
		} `json:"extData"`
		Master struct {
			ProductCode  string   `json:"productCode"`
			ID           string   `json:"id"`
			UniqueCode   string   `json:"uniqueCode"`
			SchemaCode   string   `json:"schemaCode"`
			FailLabels   []string `json:"failLabels"`
			ResourceList []string `json:"resourceList"`
			ReleaseFlag  int      `json:"releaseFlag"`
			Properties   struct {
				OriginalLocation struct {
					Country    string `json:"country"`
					AreaExtend string `json:"areaExtend"`
				} `json:"originalLocation"`
				Gender     string `json:"gender"`
				YearSeason string `json:"year_season"`
				Title      string `json:"title"`
				Brand      string `json:"brand"`
				LaunchInfo struct {
					Month  string `json:"month"`
					Year   string `json:"year"`
					Season string `json:"season"`
				} `json:"launchInfo"`
				ProductName string `json:"productName"`
				ProductNo   string `json:"productNo"`
			} `json:"properties"`
			Fragments []struct {
				ID         string `json:"id"`
				SchemaCode string `json:"schemaCode"`
				Properties struct {
					PDPImageForEC                        string `json:"PDPImageForEC"`
					CollectionCode                       string `json:"CollectionCode"`
					SilhouetteDesc                       string `json:"SilhouetteDesc"`
					DisplayNameEnUS                      string `json:"DisplayName#en_US"`
					BasecolorName                        string `json:"basecolor_name"`
					SubCollectionDesc                    string `json:"SubCollectionDesc"`
					ProductNameEn                        string `json:"productName_en"`
					EICarbonReductionCardValue           string `json:"EICarbonReductionCardValue"`
					IsCoachtopia                         string `json:"IsCoachtopia"`
					HandbagSizeZH                        string `json:"HandbagSizeZH"`
					Department                           string `json:"Department"`
					CountryOfOriginDesc                  string `json:"CountryOfOriginDesc"`
					Manufacturer                         string `json:"Manufacturer"`
					StyleCode                            string `json:"StyleCode"`
					ColorGroupDesc                       string `json:"ColorGroupDesc"`
					EnvironmentImpactCarbonReductionCard string `json:"EnvironmentImpactCarbonReductionCard"`
					ProductLifecycle                     string `json:"ProductLifecycle"`
					SubClassDescZH                       string `json:"SubClassDescZH"`
					CollectionDesc                       string `json:"CollectionDesc"`
					SwatchImage                          string `json:"SwatchImage"`
					SignatureType                        string `json:"SignatureType"`
					ShortDescriptionZhCN                 string `json:"ShortDescription#zh_CN"`
					SilhouetteCode                       string `json:"SilhouetteCode"`
					ClassCode                            string `json:"ClassCode"`
					CountryOfOriginCode                  string `json:"CountryOfOriginCode"`
					StyleGroup                           string `json:"StyleGroup"`
					StyleIntroDate                       string `json:"StyleIntroDate"`
					MaterialDescZH                       string `json:"MaterialDescZH"`
					ColorBuckets                         string `json:"ColorBuckets"`
					ProductOriginEN                      string `json:"ProductOriginEN"`
					EIkgsCO2eCardValue                   string `json:"EIkgsCO2eCardValue"`
					ShellMaterial                        string `json:"ShellMaterial"`
					HangtagMaterial                      string `json:"HangtagMaterial"`
					EIGSWasteDivertedCardValue           string `json:"EIGSWasteDivertedCardValue"`
					EnvironmentImpactGSWasteDivertedCard string `json:"EnvironmentImpactGSWasteDivertedCard"`
					Gender                               string `json:"gender"`
					Color                                string `json:"color"`
					IntroDate                            string `json:"IntroDate"`
					Video                                string `json:"video"`
					MaterialDesc                         string `json:"MaterialDesc"`
					LabelMaterial                        string `json:"LabelMaterial"`
					LongDescriptionZhCN                  string `json:"LongDescription#zh_CN"`
					HandleMaterial                       string `json:"HandleMaterial"`
					EIGSWasteDivertedCardIcon            string `json:"EIGSWasteDivertedCardIcon"`
					SubClassCode                         string `json:"SubClassCode"`
					LongDescriptionEnUS                  string `json:"LongDescription#en_US"`
					PDPImageForOW                        string `json:"PDPImageForOW"`
					EnvironmentImpactKgsCO2e             string `json:"EnvironmentImpactkgsCO2eCard"`
					EICarbonReductionCardIcon            string `json:"EICarbonReductionCardIcon"`
					ClassDesc                            string `json:"ClassDesc"`
					StyleDeleteDate                      string `json:"StyleDeleteDate"`
					RefinementColor                      string `json:"RefinementColor"`
					SubCollectionCode                    string `json:"SubCollectionCode"`
					EIKgsCO2eCardIcon                    string `json:"EIkgsCO2eCardIcon"`
					BaseColorGroup                       string `json:"basecolor_group"`
					FactoryStoreType                     string `json:"FactoryStoreType"`
					SubClassDesc                         string `json:"SubClassDesc"`
					MaterialCode                         string `json:"MaterialCode"`
					NRFColorCode                         string `json:"NRFColorCode"`
					SignatureTypeDesc                    string `json:"SignatureTypeDesc"`
				} `json:"properties"`
				ResourceList []string `json:"resourceList"`
			} `json:"fragments"`
		} `json:"master"`
		Variants []struct {
			VariantCode  string   `json:"variantCode"`
			ID           string   `json:"id"`
			ExtCode      string   `json:"extCode"`
			UniqueCode   string   `json:"uniqueCode"`
			SchemaCode   string   `json:"schemaCode"`
			ResourceList []string `json:"resourceList"`
			ReleaseFlag  int      `json:"releaseFlag"`
			OperatorID   string   `json:"operatorId"`
			Properties   struct {
				CustomColor struct {
					CustomColorID   string `json:"customColorID"`
					CustomColorName string `json:"customColorName"`
				} `json:"customColor"`
				SkuTagPrice  string `json:"skuTagPrice"`
				CustomSize   string `json:"customSize"`
				PlatformCode string `json:"platformCode"`
			} `json:"properties"`
			Fragments []struct {
				ID         string `json:"id"`
				SchemaCode string `json:"schemaCode"`
				Properties struct {
					Upc       string `json:"Upc"`
					SkuNumber string `json:"skuNumber"`
				} `json:"properties"`
				ResourceList []string `json:"resourceList"`
			} `json:"fragments"`
		} `json:"variants"`
	} `json:"data"`
	Success    bool        `json:"success"`
	Code       string      `json:"code"`
	Msg        string      `json:"msg"`
	Message    string      `json:"message"`
	RandomCode interface{} `json:"randomCode"`
}

// ProductBriefCache 宝尊商品简要信息缓存结构
type ProductBriefCache struct {
	ProductNo string `json:"product_no"` // 商品款号
	Name      string `json:"name"`       // 商品名称
	Url       string `json:"url"`        // a0主图地址
	Color     string `json:"color"`      // 商品颜色
	Upc       string `json:"upc"`        // upc码
}
