package work

type DLResultBaseData struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type RepUserPointExpireData struct {
	DLResultBaseData
	Data UserPointExpire `json:"data"`
}

type RepUserGradeCalculatorData struct {
	DLResultBaseData
	Data UserGradeCalculator `json:"data"`
}

type RepAllGradeData struct {
	DLResultBaseData
	Data []AllGrade `json:"data"`
}

type UserPointExpire struct {
	ExpireDate string `json:"expire_date"`
	ExpireNum  int64  `json:"expire_num"`
}

type UserGradeCalculator struct {
	RequiredAmount float64 `json:"requiredAmount"`
	RequiredTimes  int64   `json:"requiredTimes"`
	CurrentAmount  float64 `json:"currentAmount"`
	OnceAmount     int64   `json:"onceAmount"`
}

type AllGrade struct {
	TierName string `json:"tierName"`
	TierCode string `json:"tierCode"`
}

type CommonResponseData struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}
