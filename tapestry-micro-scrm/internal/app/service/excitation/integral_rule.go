package excitationSvc

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
	"sort"
	"time"
)

type IntegralRuleService struct {
	service.Base
	integralRuleDao       *excitationDao.IntegralRule
	integralRuleDetailDao *excitationDao.IntegralRuleDetail
}

func NewIntegralRuleService() *IntegralRuleService {
	return &IntegralRuleService{
		integralRuleDao:       excitationDao.NewIntegralRule(),
		integralRuleDetailDao: excitationDao.NewIntegralRuleDetail(),
	}
}

func (s *IntegralRuleService) GetList(c ctx.Context, params excitationReq.RuleList) ([]excitationResp.IntegralRuleList, int64, error) {
	var (
		resp   []excitationResp.IntegralRuleList
		ruleId []int64
		nums   int64
	)

	//规则信息
	result, total, err := s.integralRuleDao.GetList(c, params)

	if err != nil {
		return resp, total, err
	}

	ruleNum := make(map[int64]int64)
	for _, value := range result {
		ruleId = append(ruleId, value.ID)
	}
	//规则详情信息
	ruleDetail, _ := s.integralRuleDetailDao.GetDetail(c, ruleId)
	for _, v := range ruleDetail {
		if _, ok := ruleNum[v.RuleID]; ok {
			ruleNum[v.RuleID] = ruleNum[v.RuleID] + 1
		} else {
			ruleNum[v.RuleID] = 1
		}
	}

	for _, value := range result {
		if _, ok := ruleNum[value.ID]; ok {
			nums = ruleNum[value.ID]
		}
		//如果没有查询到数据默认就是关闭状态
		status := int64(value.IncentiveRelation.Status)
		if status == 0 {
			status = 2
		}
		resp = append(resp, excitationResp.IntegralRuleList{
			Id:           value.ID,
			EnterpriseId: value.EnterpriseID,
			Type:         value.Type,
			TypeName:     value.IncentiveConfig.RuleName,
			Status:       status,
			RelationId:   value.IncentiveRelation.ID,
			RuleName:     value.RuleName,
			RuleNum:      nums,
			CreatedAt:    value.CreatedAt,
			UpdatedAt:    value.UpdatedAt,
		})
	}

	return resp, total, err

}

//ParameterConfiguration 积分规则条件配置信息
func (s *IntegralRuleService) ParameterConfiguration(params excitationReq.ParameterConfiguration) []excitationResp.Configuration {
	ruleIfDefault := []excitationResp.RuleIfValue{
		{
			Value: excitationModels.RuleIfTimes,
			Name:  "次数条件",
		},
	}
	ruleIf := []excitationResp.RuleIfValue{
		{
			Value: excitationModels.RuleIfGT,
			Name:  "高于",
		}, {
			Value: excitationModels.RuleIfEgt,
			Name:  "达到",
		},
	}
	ruleIf1 := []excitationResp.RuleIfValue{
		{
			Value: excitationModels.RuleIf1Value,
			Name:  "百分比",
		}, {
			Value: excitationModels.RuleIf1Average,
			Name:  "平均值",
		},
	}
	respMap := map[int][]excitationResp.Configuration{
		excitationModels.IncentiveConfigTypeInvite: {
			{
				Id:   excitationModels.RuleCategoryId1,
				Name: "邀约数(每次)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId2,
				Name: "邀约到店率(每周)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  ruleIf1,
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId3,
				Name: "邀约成功率(每月)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  ruleIf1,
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId4,
				Name: "邀约成功(每次)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId5,
				Name: "邀约占比(每周)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  ruleIf1,
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId6,
				Name: "邀约到店确认(每次)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
		},
		excitationModels.IncentiveConfigTypeTask: {
			{
				Id:   excitationModels.RuleCategoryId7,
				Name: "任务完成(每次)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId8,
				Name: "任务完成率(每月)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId9,
				Name: "任务触达(每人)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
		},
		excitationModels.IncentiveConfigTypeInteract: {
			{
				Id:   excitationModels.RuleCategoryId10,
				Name: "新增好友数(每人)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId11,
				Name: "客户流失数(每人)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId12,
				Name: "聊天总数(每月)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  ruleIf1,
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId13,
				Name: "已回复聊天占比(每月)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  ruleIf1,
					Selected: 2,
				},
			},
			{
				Id:   excitationModels.RuleCategoryId14,
				Name: "完成播种(按次)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
		},
		excitationModels.IncentiveConfigTypeTrain: {
			{
				Id:   excitationModels.RuleCategoryId15,
				Name: "积分等比兑换",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIfDefault,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 1,
				},
			},
		},
		excitationModels.IncentiveConfigTypeOnline: {
			{
				Id:   excitationModels.RuleCategoryId16,
				Name: "线上销售业绩率(每周)",
				Rule: excitationResp.ConfigurationRule{
					RuleIf:   ruleIf,
					RuleIf1:  []excitationResp.RuleIfValue{},
					Selected: 2,
				},
			},
		},
	}

	var resp []excitationResp.Configuration
	resp = respMap[params.Type]
	return resp
}

func (s *IntegralRuleService) Create(c ctx.Context, params excitationReq.RuleCreate) (err error) {

	ruleData := excitationModels.ExcitationIntegralRuleModel{
		EnterpriseID: params.EnterpriseId,
		Type:         params.Type,
		RuleName:     params.RuleName,
		UpdatedAt:    time.Now().Unix(),
	}

	// 维度不能重复条件
	exits, err := s.integralRuleDao.GetOne(c, "enterprise_id = ? and type = ?", params.EnterpriseId, params.Type)
	if err != nil {
		return err
	}
	if exits.ID > 0 {
		return errors.New("维度不能重复条件")
	}

	err = app.DB().Transaction(func(tx *gorm.DB) error {
		// 写入数据
		_, ruleId := excitationDao.NewIntegralRule().WithTx(tx).CreateData(c, ruleData)

		// 规则信息
		var ruleCategoryId []int64
		var ruleDetailData []excitationModels.ExcitationIntegralRuleDetailModel
		for _, detail := range params.RuleDetail {
			ruleCategoryId = append(ruleCategoryId, detail.RuleCategoryID)
			for _, value := range detail.RuleValue {
				ruleDetailData = append(ruleDetailData, excitationModels.ExcitationIntegralRuleDetailModel{
					EnterpriseID:   params.EnterpriseId,
					RuleID:         ruleId,
					RuleCategoryID: detail.RuleCategoryID,
					ValidityType:   detail.ValidityType,
					StartTime:      detail.StartTime,
					EndTime:        detail.EndTime,
					IsDiffStore:    detail.IsDiffStore, //是否区分奥莱零售
					StoreType:      value.StoreType,    //0 通用 1 奥莱 2 零售
					Integral:       value.Integral,
					IntegralType:   value.IntegralType,
					RuleIf:         value.RuleIf,
					RuleIf1:        value.RuleIf1,
					RuleValue:      value.RuleValue,
					IsAstrict:      detail.IsAstrict,
					AstrictType:    detail.AstrictType,
					MaxIntegral:    detail.MaxIntegral,
				})
			}
		}

		// 判断不能重复
		ruleCategoryMap := make(map[int64]int)
		for _, v := range ruleCategoryId {
			if _, ok := ruleCategoryMap[v]; ok {
				//return errors.New("规则类目不能重复")
				return fmt.Errorf("规则类目不能重复")
			} else {
				ruleCategoryMap[v] = 1
			}
		}

		//写入规则
		err = excitationDao.NewIntegralRuleDetail().WithTx(tx).CreateData(c, ruleDetailData)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *IntegralRuleService) Detail(c ctx.Context, params excitationReq.RuleDetail) (excitationResp.IntegralRule, error) {
	var resp excitationResp.IntegralRule

	result, err := s.integralRuleDao.GetOne(c, "enterprise_id = ? and id = ?", params.EnterpriseId, params.Id)
	if err != nil {
		return resp, err
	}

	detailMap := make(map[int64]*excitationResp.IntegralRuleDetail)
	for _, value := range result.Details {
		if _, ok := detailMap[value.RuleCategoryID]; ok {
			//追加数据
			detailMap[value.RuleCategoryID].DetailValue = append(detailMap[value.RuleCategoryID].DetailValue, excitationResp.IntegralRuleDetailValue{
				Integral:     value.Integral,
				IntegralType: value.IntegralType,
				RuleIf:       value.RuleIf,
				RuleIf1:      value.RuleIf1,
				RuleValue:    value.RuleValue,
				StoreType:    value.StoreType,
			})
		} else {
			//第一次赋值
			detailMap[value.RuleCategoryID] = &excitationResp.IntegralRuleDetail{
				ID:             value.ID,
				EnterpriseID:   value.EnterpriseID,
				RuleID:         value.RuleID,
				RuleCategoryID: value.RuleCategoryID,
				ValidityType:   value.ValidityType,
				StartTime:      value.StartTime,
				EndTime:        value.EndTime,
				IsAstrict:      value.IsAstrict,
				AstrictType:    value.AstrictType,
				MaxIntegral:    value.MaxIntegral,
				IsDiffStore:    value.IsDiffStore,
				DetailValue: []excitationResp.IntegralRuleDetailValue{{
					Integral:     value.Integral,
					IntegralType: value.IntegralType,
					RuleIf:       value.RuleIf,
					RuleIf1:      value.RuleIf1,
					RuleValue:    value.RuleValue,
					StoreType:    value.StoreType,
				}},
			}
		}
	}

	//规则信息
	detail := resp.Detail
	for _, item := range detailMap {
		detail = append(detail, excitationResp.IntegralRuleDetail{
			ID:             item.ID,
			EnterpriseID:   item.EnterpriseID,
			RuleID:         item.RuleID,
			RuleCategoryID: item.RuleCategoryID,
			ValidityType:   item.ValidityType,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
			IsAstrict:      item.IsAstrict,
			AstrictType:    item.AstrictType,
			MaxIntegral:    item.MaxIntegral,
			IsDiffStore:    item.IsDiffStore,
			StoreType:      item.StoreType,
			DetailValue:    item.DetailValue,
		})
	}

	sort.Slice(detail, func(i, j int) bool {
		return detail[i].ID < detail[j].ID
	})

	resp = excitationResp.IntegralRule{
		ID:           result.ID,
		EnterpriseID: result.EnterpriseID,
		Type:         result.Type,
		RuleName:     result.RuleName,
		Detail:       detail,
	}

	return resp, err
}

func (s *IntegralRuleService) Destroy(c ctx.Context, params excitationReq.RuleDestroy) error {

	// where条件
	var ruleWhere = map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"id":            params.Id,
	}
	err := s.integralRuleDao.DeleteData(c, ruleWhere)
	fmt.Println(err)
	if err != nil {
		return err
	}

	//规则信息where条件
	var ruleDetailWhere = map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"rule_id":       params.Id,
	}
	err = s.integralRuleDetailDao.DeleteData(c, ruleDetailWhere)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *IntegralRuleService) Update(c ctx.Context, params excitationReq.RuleUpdate) error {

	err := app.DB().Transaction(func(tx *gorm.DB) error {
		// 更新数据
		var ruleWhere = map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"id":            params.ID,
		}
		ruleData := excitationModels.ExcitationIntegralRuleModel{
			RuleName:  params.RuleName,
			UpdatedAt: time.Now().Unix(),
		}

		err := excitationDao.NewIntegralRule().WithTx(tx).UpdateData(c, ruleWhere, ruleData)
		if err != nil {
			return errors.New("修改失败")
		}

		// 规则信息
		var ruleCategoryId []int64
		var ruleDetailData []excitationModels.ExcitationIntegralRuleDetailModel
		for _, detail := range params.RuleDetail {
			ruleCategoryId = append(ruleCategoryId, detail.RuleCategoryID)
			for _, value := range detail.RuleValue {
				ruleDetailData = append(ruleDetailData, excitationModels.ExcitationIntegralRuleDetailModel{
					EnterpriseID:   params.EnterpriseId,
					RuleID:         params.ID,
					RuleCategoryID: detail.RuleCategoryID,
					ValidityType:   detail.ValidityType,
					StartTime:      detail.StartTime,
					EndTime:        detail.EndTime,
					IsAstrict:      detail.IsAstrict,
					AstrictType:    detail.AstrictType,
					MaxIntegral:    detail.MaxIntegral,
					IsDiffStore:    detail.IsDiffStore,
					StoreType:      value.StoreType,
					Integral:       value.Integral,
					IntegralType:   value.IntegralType,
					RuleIf:         value.RuleIf,
					RuleIf1:        value.RuleIf1,
					RuleValue:      value.RuleValue,
				})
			}
		}

		// 判断累不能重复
		ruleCategoryMap := make(map[int64]int)
		for _, v := range ruleCategoryId {
			if _, ok := ruleCategoryMap[v]; ok {
				//return errors.New("规则类目不能重复")
				return fmt.Errorf("规则类目不能重复")
			} else {
				ruleCategoryMap[v] = 1
			}
		}

		//先删除旧规则信息
		var ruleDetailWhere = map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"rule_id":       params.ID,
		}
		err = excitationDao.NewIntegralRuleDetail().WithTx(tx).DeleteData(c, ruleDetailWhere)
		if err != nil {
			return err
		}

		//写入规则
		err = excitationDao.NewIntegralRuleDetail().WithTx(tx).CreateData(c, ruleDetailData)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *IntegralRuleService) GetMany(c ctx.Context, params excitationReq.RuleMany) ([]excitationModels.ExcitationIntegralRuleModel, int64, error) {
	result, total, err := s.integralRuleDao.GetMany(c, params)
	if err != nil {
		return result, total, err
	}

	return result, total, err
}
