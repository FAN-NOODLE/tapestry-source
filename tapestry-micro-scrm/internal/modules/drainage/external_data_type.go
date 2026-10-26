package drainage

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/welcome"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
	"strings"
	"time"
)

// GetDataByTypeId 获取指定类型下的所有类型信息
func GetDataByTypeId(r ExternalTypeIdRequest) ([]welcome.WelcomesExternalTypeDataModel, error) {
	data, err := orm.Service.Welcome.GetDataByTypeId(r.TypeId, r.EnterpriseId)
	return data, err
}

func GetTypeId(enterpriseId int64) (int64, error) {
	data, err := orm.Service.Welcome.GetTypeId(enterpriseId)
	return data, err
}

// GetAllType 获取所有的类型
func GetAllType(r ExternalTypeAllRequest) ([]welcome.WelcomesExternalTypeModel, error) {
	data, err := orm.Service.Welcome.GetAllType(r.EnterpriseId)
	return data, err
}

func WelcomeDel(r WelcomeDelRequest) error {
	err := orm.Service.Welcome.DelWelcome(r.ID)
	return err
}

func WelcomeDetail(r WelcomeDetailNewRequest) (WelcomeAllResponse, error) {
	drainageWelcome, err := orm.Service.Welcome.GetWelcomeDetailById(r.ID)
	var dataKeyValue map[string]interface{}
	helper.Unmarshal(drainageWelcome.ConditionContent, &dataKeyValue)
	resItem := translator(dataKeyValue, drainageWelcome)
	return resItem, err
}

// WelcomeNewEdit 欢迎语的新增和修改
func WelcomeNewEdit(ctx ctx.Context, r ExternalWelcomeRequest) (int64, error) {
	data, _ := orm.Service.Welcome.GetWelcomeDetailById(r.ID)
	if data.ID == 0 {
		contentStr, _ := json.Marshal(r.Content)
		m := welcome.DrainageWelcomesModel{
			Content:          string(contentStr),
			IsDefault:        welcome.ISDEFAULT_KD,
			ChannelID:        welcome.CHANNLEID_KD,
			Type:             welcome.TypeSelfHoo,
			ConditionContent: helper.MarshalToString(r.ConditionContent),
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				EnterpriseID: r.EnterpriseId,
				BaseModel: model.BaseModel{
					ID: uid.ID(ctx),
				},
			},
		}
		res, err := orm.Service.Welcome.CreateWelcomeDetail(m)
		if err != nil {
			return 0, err
		}
		return res, nil
	}

	contentStr, _ := json.Marshal(r.Content)
	m := welcome.DrainageWelcomesModel{
		Content: string(contentStr),
		MicroScrmBaseModel: model.MicroScrmBaseModel{
			BaseModel: model.BaseModel{
				ID:        data.ID,
				UpdatedAt: time.Now().Unix(),
			},
		},
		ConditionContent: helper.MarshalToString(r.ConditionContent),
	}
	res, err := orm.Service.Welcome.EditWelcomeDetail(m)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func GetAllWelcome(r WelcomeAllRequest) ([]WelcomeAllResponse, int64, error) {
	drainageWelcomes, total, err := orm.Service.Welcome.GetAllWelcome(r.EnterpriseId, r.Size, r.Page, r.KeyWords)
	result := make([]WelcomeAllResponse, 0)
	welcomesExternalTypeDataKeyValue := getWelcomesExternalTypeDataAllKeyValue(r.EnterpriseId)
	for _, drainageWelcome := range drainageWelcomes {
		resItem := translator(welcomesExternalTypeDataKeyValue, drainageWelcome)
		result = append(result, resItem)
	}
	return result, total, err
}

func translator(dataKeyValue map[string]interface{}, data welcome.DrainageWelcomesModel) WelcomeAllResponse {
	var content []map[string]interface{}
	helper.Unmarshal(data.Content, &content)

	var conditionContent map[string][]any
	helper.Unmarshal(data.ConditionContent, &conditionContent)

	var codeName []string
	for _, value := range conditionContent {
		if len(value) > 0 {
			for _, v := range value {
				if val, ok := dataKeyValue[v.(string)]; ok {
					codeName = append(codeName, val.(string))
				}
			}
		}
	}

	var resItem WelcomeAllResponse
	resItem.TypeName = strings.Join(codeName, ",")

	resItem.ContentText = content
	resItem.ConditionContent = conditionContent

	resItem.ID = data.ID
	resItem.EnterpriseID = data.EnterpriseID
	resItem.IsDefault = data.IsDefault
	resItem.ChannelID = data.ChannelID
	resItem.Type = data.Type
	resItem.CreatedAt = data.CreatedAt
	resItem.UpdatedAt = data.UpdatedAt

	return resItem
}

func getWelcomesExternalTypeDataAllKeyValue(enterpriseId int64) map[string]interface{} {
	welcomesExternalTypeData, _ := orm.Service.Welcome.GetWelcomesExternalTypeDataAll(welcome.WelcomesExternalTypeDataModel{
		MicroScrmBaseModel: model.MicroScrmBaseModel{
			EnterpriseID: enterpriseId,
		},
	})
	allType, _ := orm.Service.Welcome.GetAllType(enterpriseId)
	welcomesType := make(map[int64]string)
	for _, typeModel := range allType {
		welcomesType[typeModel.ID] = typeModel.Code
	}
	welcomesExternalTypeDataKeyValue := make(map[string]interface{})
	for _, item := range welcomesExternalTypeData {
		if typeCode, ok := welcomesType[item.WelcomesExternalTypeID]; ok {
			if helper.CheckCoachTime() {
				if typeCode == welcome.MemberLevel {
					continue
				}
			} else {
				if typeCode == welcome.MemberLevelFY26 {
					continue
				}
			}
			welcomesExternalTypeDataKeyValue[item.Code] = item.Name
		}
	}
	return welcomesExternalTypeDataKeyValue
}
