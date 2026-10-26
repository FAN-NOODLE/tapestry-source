package excitationSvc

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	membersDao "micro-scrm/internal/app/dao/members"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"time"
)

type MedalService struct {
	service.Base
	medalDao *excitationDao.Medal
}

func NewMedalService() *MedalService {
	return &MedalService{
		medalDao: excitationDao.NewMedal(),
	}
}

func (s *MedalService) GetMedalList(c ctx.Context, params excitationReq.MedalList) ([]excitationResp.MedalList, int64, error) {
	var resp []excitationResp.MedalList
	result, total, err := s.medalDao.GetPagination(c, params)
	if err != nil {
		return resp, total, err
	}

	var medalRuleIds []int
	var medalIds []int64
	for _, model := range result {
		medalRuleIds = append(medalRuleIds, model.RuleDimensionID, model.RuleCategoryID)
		medalIds = append(medalIds, model.ID)
	}

	medalRuleNameMap := make(map[int64]string)
	medalRule, err := excitationDao.NewMedalRule().GetMany(c, "id in (?) and status = ?", medalRuleIds, excitationModels.MedalRuleStatusEnable)
	if err != nil {
		return resp, total, err
	}
	for _, rule := range medalRule {
		medalRuleNameMap[rule.ID] = rule.RuleName
	}
	//获取是否已经有奖章记录
	existsMedalIds, err := excitationDao.NewMedalRecord().GetExistsByMedalIds(c, params.EnterpriseId, medalIds)
	if err != nil {
		return nil, 0, err
	}
	for _, model := range result {
		existsRecord := false
		if helper.InArray(model.ID, existsMedalIds) {
			existsRecord = true
		}
		tmp := excitationResp.MedalList{
			MedalBase: excitationResp.MedalBase{
				Id:              model.ID,
				Name:            model.Name,
				Source:          model.Source,
				RuleType:        model.RuleType,
				RuleStartTime:   model.RuleStartTime,
				RuleEndTime:     model.RuleEndTime,
				MedaType:        model.MedaType,
				MedaDescription: model.MedaDescription,
				OtherRewards:    model.OtherRewards,
				Integral:        model.Integral,
				ItemName:        model.ItemName,
				RuleDimensionID: model.RuleDimensionID,
				RuleCategoryID:  model.RuleCategoryID,
				InviteTop:       model.InviteTop,
				CreateTime:      model.CreatedAt,
			},
			Status:            model.Status,
			RuleDimensionName: medalRuleNameMap[int64(model.RuleDimensionID)],
			RuleCategoryName:  medalRuleNameMap[int64(model.RuleCategoryID)],
			ExistsRecord:      existsRecord,
		}

		// 判断有效 1是有效 2是无效
		nowTime := time.Now().Unix()
		if (int64(model.RuleStartTime) <= nowTime && int64(model.RuleEndTime) > nowTime) || model.RuleType == excitationModels.MedalRuleTypeLong {
			tmp.IsEffective = 1
		}
		if int64(model.RuleStartTime) >= nowTime && int64(model.RuleEndTime) < nowTime {
			tmp.IsEffective = 2
		}
		resp = append(resp, tmp)
	}

	return resp, total, err
}

func (s *MedalService) Create(c ctx.Context, params excitationReq.MedalCreate) (err error) {
	// 其它奖励参数验证
	if params.OtherRewards == 2 && params.Integral == 0 {
		return errors.New("其它奖励积分不能为空")
	}
	if params.OtherRewards == 3 && params.ItemName == "" {
		return errors.New("其它奖励奖品名称不能为空")
	}
	if params.OtherRewards == 4 && params.Integral == 0 && params.ItemName == "" {
		return errors.New("其它奖励积分和奖品名称不能为空")
	}

	// 防止前端切换乱传
	if params.OtherRewards == 1 {
		params.Integral = 0
		params.ItemName = ""
	}
	if params.OtherRewards == 2 {
		params.ItemName = ""
	}
	if params.OtherRewards == 3 {
		params.Integral = 0
	}

	// 名称重复校验
	exits, err := s.medalDao.GetOne(c, "enterprise_id = ? and name = ?", params.EnterpriseId, params.Name)
	if err != nil {
		return err
	}
	if exits.ID > 0 {
		return errors.New("规则名称重复")
	}

	if params.Source == excitationModels.MedalSourceManualUpload {
		err = s.manualUploadCreate(c, params)
	}

	if params.Source == excitationModels.MedalSourceRuleStatistics {
		err = s.ruleStatisticsCreate(c, params)
	}

	return err
}

// 手动上传创建
func (s *MedalService) manualUploadCreate(c ctx.Context, params excitationReq.MedalCreate) error {
	if params.RuleType != excitationModels.MedalRuleTypeOnlyOne {
		return errors.New("手动上传-只能为仅生效一次")
	}

	if len(params.UserIds) == 0 {
		return errors.New("人群包不能为空")
	}

	params.UserIds = lo.Uniq(params.UserIds)

	medalData := excitationModels.ExcitationMedalModel{
		MicroBaseModel:  models.CreateMicroBase(params.EnterpriseId, app.UID()),
		Name:            params.Name,
		Source:          params.Source,
		RuleType:        params.RuleType,
		MedaType:        params.MedaType,
		MedaDescription: params.MedaDescription,
		OtherRewards:    params.OtherRewards,
		Integral:        params.Integral,
		ItemName:        params.ItemName,
		Status:          excitationModels.MedalStatusEnable,
	}

	err := app.DB().Transaction(func(tx *gorm.DB) error {
		// 创建规则
		err := excitationDao.NewMedal().WithTx(tx).CreateData(c, medalData)
		if err != nil {
			return err
		}

		medalStyle := excitationModels.ExcitationMedalStyleModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
			MedalID:        medalData.ID,
			MedalName:      params.Styles[0].MedalName,
			LightIcon:      params.Styles[0].LightIcon,
			UnlitIcon:      params.Styles[0].UnlitIcon,
		}
		// 创建规则样式
		err = excitationDao.NewMedalStyle().WithTx(tx).CreateData(c, medalStyle)
		if err != nil {
			return err
		}

		// 批量添加人群包奖章
		members, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"userid":        params.UserIds,
		})
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return fmt.Errorf("人群包数据为空")
		}

		var medalRecords []excitationModels.ExcitationMedalRecordModel
		for _, member := range members {
			medalRecords = append(medalRecords, excitationModels.ExcitationMedalRecordModel{
				MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
				MedalStyleID:   medalStyle.ID,
				MedalID:        medalData.ID,
				MemberID:       member.ID,
			})
		}
		err = excitationDao.NewMedalRecord().WithTx(tx).CreateInBatchesData(c, medalRecords)
		if err != nil {
			return err
		}

		// 发送积分
		if params.Integral > 0 {
			var addMedalPoints []excitationReq.ManualPointsSave
			for _, member := range members {
				tmp := excitationReq.ManualPointsSave{
					MemberId:     member.ID,
					RuleName:     params.Name,
					RuleDetailId: 0,
					PointsRemark: fmt.Sprintf("奖章-%s", medalStyle.MedalName),
					Points:       int64(params.Integral),
				}
				addMedalPoints = append(addMedalPoints, tmp)
			}
			err = NewPointsProduceService().MedalPoints(c, params.EnterpriseId, addMedalPoints)
			return err
		}

		return nil
	})

	return err
}

// 规则统计创建
func (s *MedalService) ruleStatisticsCreate(c ctx.Context, params excitationReq.MedalCreate) error {
	if !helper.InArray(params.RuleType, []int32{excitationModels.MedalRuleTypeLong, excitationModels.MedalRuleTypeSpecify}) {
		return errors.New("规则统计-只能为长期有效或指定时间")
	}

	if params.RuleType == excitationModels.MedalRuleTypeSpecify && params.RuleStartTime == 0 && params.RuleEndTime == 0 {
		return errors.New("规则统计-指定时间不能为空")
	}

	err := s.Base.DB().Transaction(func(tx *gorm.DB) error {
		medalData := excitationModels.ExcitationMedalModel{
			MicroBaseModel:  models.CreateMicroBase(params.EnterpriseId, app.UID()),
			Name:            params.Name,
			Source:          params.Source,
			RuleType:        params.RuleType,
			RuleStartTime:   params.RuleStartTime,
			RuleEndTime:     params.RuleEndTime,
			MedaType:        params.MedaType,
			MedaDescription: params.MedaDescription,
			OtherRewards:    params.OtherRewards,
			Integral:        params.Integral,
			ItemName:        params.ItemName,
			RuleDimensionID: params.RuleDimensionID,
			RuleCategoryID:  params.RuleCategoryID,
			InviteTop:       params.InviteTop,
			Status:          excitationModels.MedalStatusEnable,
		}
		err := excitationDao.NewMedal().WithTx(tx).CreateData(c, medalData)
		if err != nil {
			return err
		}

		var medalStyles []excitationModels.ExcitationMedalStyleModel
		for _, style := range params.Styles {
			medalStyles = append(medalStyles, excitationModels.ExcitationMedalStyleModel{
				MicroBaseModel: models.CreateMicroBase(params.EnterpriseId, app.UID()),
				MedalID:        medalData.ID,
				TimeValue:      style.TimeValue,
				TimeName:       style.TimeName,
				MedalName:      style.MedalName,
				LightIcon:      style.LightIcon,
				UnlitIcon:      style.UnlitIcon,
			})
		}
		// 创建规则样式
		err = excitationDao.NewMedalStyle().WithTx(tx).CreateInBatchesData(c, medalStyles)
		return err
	})

	return err
}

func (s *MedalService) Detail(c ctx.Context, id int64) (excitationResp.MedalDetail, error) {
	var resp excitationResp.MedalDetail
	result, err := s.medalDao.GetOne(c, "id = ?", id)
	if err != nil {
		return resp, err
	}

	medalRuleNameMap := make(map[int64]string)
	if result.Source == excitationModels.MedalSourceRuleStatistics {
		medalRule, err := excitationDao.NewMedalRule().GetMany(c, "id in (?) and  status = 1", []int{result.RuleDimensionID, result.RuleCategoryID})
		if err != nil {
			return resp, err
		}
		for _, rule := range medalRule {
			medalRuleNameMap[rule.ID] = rule.RuleName
		}
	}

	medalStyle, err := excitationDao.NewMedalStyle().GetMany(c, "medal_id = ?", id)
	if err != nil {
		return resp, err
	}

	var styles []excitationResp.MedalDetailStyle
	for _, style := range medalStyle {
		styles = append(styles, excitationResp.MedalDetailStyle{
			Id:        style.ID,
			TimeValue: style.TimeValue,
			TimeName:  style.TimeName,
			MedalName: style.MedalName,
			LightIcon: style.LightIcon,
			UnlitIcon: style.UnlitIcon,
		})
	}

	resp = excitationResp.MedalDetail{
		MedalBase: excitationResp.MedalBase{
			Id:              result.ID,
			Name:            result.Name,
			Source:          result.Source,
			RuleType:        result.RuleType,
			RuleStartTime:   result.RuleStartTime,
			RuleEndTime:     result.RuleEndTime,
			MedaType:        result.MedaType,
			MedaDescription: result.MedaDescription,
			OtherRewards:    result.OtherRewards,
			Integral:        result.Integral,
			ItemName:        result.ItemName,
			RuleDimensionID: result.RuleDimensionID,
			RuleCategoryID:  result.RuleCategoryID,
			InviteTop:       result.InviteTop,
			CreateTime:      result.CreatedAt,
		},
		RuleDimensionName: medalRuleNameMap[int64(result.RuleDimensionID)],
		RuleCategoryName:  medalRuleNameMap[int64(result.RuleCategoryID)],
		Styles:            styles,
	}

	return resp, nil
}

func (s *MedalService) ChangeStatus(c ctx.Context, id int64, status int) error {
	return excitationDao.NewMedal().UpdateData(c, map[string]interface{}{"id": id}, map[string]interface{}{"status": status})
}

func (s *MedalService) RemoveRule(c ctx.Context, id int64) error {
	err := s.medalDao.DeleteData(c, id)
	if err != nil {
		return err
	}

	return excitationDao.NewMedalStyle().DeleteData(c, map[string]interface{}{"medal_id": id})
}

func (s *MedalService) Update(c ctx.Context, params excitationReq.MedalUpdate) error {
	result, err := s.medalDao.GetOne(c, "id = ?", params.Id)
	if err != nil {
		return err
	}
	if result.ID == 0 {
		return errors.New("数据为空")
	}

	if result.Source == excitationModels.MedalSourceManualUpload {
		return errors.New("手动上传不能编辑")
	}

	// 名称重复校验
	exits, err := s.medalDao.GetOne(c, "enterprise_id = ? and name = ?", params.EnterpriseId, params.Name)
	if err != nil {
		return err
	}
	if exits.ID > 0 && exits.ID != result.ID {
		return errors.New("规则名称重复")
	}

	medalData := excitationModels.ExcitationMedalModel{
		Name:            params.Name,
		Source:          params.Source,
		RuleType:        params.RuleType,
		RuleStartTime:   params.RuleStartTime,
		RuleEndTime:     params.RuleEndTime,
		MedaType:        params.MedaType,
		MedaDescription: params.MedaDescription,
		OtherRewards:    params.OtherRewards,
		Integral:        params.Integral,
		ItemName:        params.ItemName,
		RuleDimensionID: params.RuleDimensionID,
		RuleCategoryID:  params.RuleCategoryID,
		InviteTop:       params.InviteTop,
	}

	err = excitationDao.NewMedal().UpdateData(c, map[string]any{"id": params.Id}, medalData)
	if err != nil {
		return err
	}
	// 更新奖章样式
	for _, style := range params.Styles {
		_ = excitationDao.NewMedalStyle().UpdateData(c, map[string]any{"id": style.Id}, map[string]any{
			"medal_name": style.MedalName,
			"light_icon": style.LightIcon,
			"unlit_icon": style.UnlitIcon,
		})
	}

	return err
}

// MedalWall 奖章墙
func (s *MedalService) MedalWall(c ctx.Context, enterpriseId int64, mobile string) (excitationResp.MedalWall, error) {
	var resp excitationResp.MedalWall
	// 查询员工企微id
	membersUser, err := s.GetUserIdByMobile(c, enterpriseId, mobile)
	if err != nil {
		return resp, err
	}

	// 获取已获得的奖章
	memberMedalRecord, err := excitationDao.NewMedalRecord().GetGained(c, enterpriseId, membersUser.ID)
	if err != nil {
		return resp, err
	}

	// 员工以获取奖章的规则id集合
	var memberMedalIds []int64
	for _, record := range memberMedalRecord {
		memberMedalIds = append(memberMedalIds, record.MedalID)
	}

	// 获取创建规则的所有数据
	medalList, err := excitationDao.NewMedal().GetManyUnscoped(c, "enterprise_id = ? and (id in ? or meda_type = ?)", enterpriseId, memberMedalIds, excitationModels.MedalMedaTypeOrdinary)
	if err != nil {
		return resp, err
	}
	if len(medalList) == 0 {
		return resp, err
	}

	//获取所有财月
	allMonth, _ := microWorkDao.NewFiscalMonth().GetAllMonth(c, enterpriseId)
	fiscalYearMap := make(map[int64]excitationResp.FiscalYear)
	for _, month := range allMonth {
		tmpYear := fiscalYearMap[month.FiscalYear]
		tmpYear.Year = month.FiscalYear
		if month.Period == 1 {
			tmpYear.StartTime = month.StartTime
		}
		if month.Period == 12 {
			tmpYear.EndTime = month.EndTime
		}
		fiscalYearMap[month.FiscalYear] = tmpYear
	}
	// 取当前时间所在财年
	fiscalMonth, err := microWorkDao.NewFiscalMonth().GetDataByYearAndTime(c, enterpriseId, time.Now().Unix())
	resp.NowYear = fiscalMonth.FiscalYear
	if err != nil {
		return resp, err
	}
	for year, yearInfo := range fiscalYearMap {
		if year > resp.NowYear {
			continue
		}
		memberMedalStyleMap := make(map[int64]excitationModels.ExcitationMedalRecordModel)
		yearData := excitationResp.YearMedalWall{
			Year: year,
		}
		for _, model := range medalList {
			if model.RuleType == excitationModels.MedalRuleTypeLong {
				if model.CreatedAt > yearInfo.EndTime {
					//长期规则 只判断规则的创建时间是否在财年之后
					continue
				}
			}
			// 如果规则有效期与当前财年没有交集则跳过(由于奖章发送时间一般会有延迟，如果交集小于3天就不算有交集）
			if model.RuleType == excitationModels.MedalRuleTypeSpecify {
				if (int64(model.RuleEndTime) - yearInfo.StartTime) < 86400*3 {
					//规则结束时间与当前财年开始时间小于3天则当前财年忽略该规则
					continue
				}
				if (yearInfo.EndTime - int64(model.RuleStartTime)) < 86400*3 {
					//当前财年结束时间与规则开始时间小于3天则当前财年忽略该规则
					continue
				}
			}
			if model.RuleType == excitationModels.MedalRuleTypeOnlyOne {
				if yearInfo.StartTime > model.CreatedAt || yearInfo.EndTime < model.CreatedAt {
					//一次性规则不在当前财年则忽略
					continue
				}
			}
			// 判断是否领取过当前奖章
			var exist bool
			for _, record := range memberMedalRecord {
				if record.MedalID == model.ID {
					if (record.CreatedAt - yearInfo.EndTime) > 86400 {
						//如果奖章领取时间已经超出财年结束时间一天 则该财年忽略该奖章
						continue
					}
					if (record.CreatedAt - yearInfo.StartTime) < 86400 {
						//如果奖章领取时未超过新财年的第一天则忽略（算到上一个财年)
						continue
					}

					memberMedalStyleMap[record.MedalStyleID] = record
					exist = true
				}
			}

			// 已删除并且未领取不统计和展示
			if model.DeletedAt > 0 && !exist {
				continue
			}

			// 手动奖章未领取不展示
			if model.Source == excitationModels.MedalSourceManualUpload && !exist {
				continue
			}

			wallList := excitationResp.MedalWallList{RuleId: model.ID}
			// 判断当前规则类型
			// 取出当前时间上月，上季度，上年，奖章
			// 对应 medal.MedalStyle map坐标
			var timeValue int64
			switch model.MedalRule.TimeType {
			case excitationModels.MedalRuleTimeTypeMonth:
				monthInt := fiscalMonth.Period - 1 //默认展示上个月的奖章
				timeValue = monthInt - 1           //奖章下标
				if monthInt <= 0 {
					timeValue = 11
				}
			case excitationModels.MedalRuleTimeTypeQuarter:
				// 判断当前是第几季度,
				if helper.InArray(fiscalMonth.Period, []int64{1, 2, 3}) {
					// 当前是1季度则展示去年4季度（奖章下标3）
					timeValue = 3
				}
				if helper.InArray(fiscalMonth.Period, []int64{4, 5, 6}) {
					// 当前是2季度则展示1季度（奖章下标0）
					timeValue = 0
				}
				if helper.InArray(fiscalMonth.Period, []int64{7, 8, 9}) {
					// 当前是3季度则展示2季度（奖章下标1）
					timeValue = 1
				}
				if helper.InArray(fiscalMonth.Period, []int64{10, 11, 12}) {
					// 当前是4季度则展示1季度（奖章下标2）
					timeValue = 2
				}
			case excitationModels.MedalRuleTimeTypeYear:
				// 年份默认为1
				timeValue = 0
			}
			if year == fiscalMonth.FiscalYear && fiscalMonth.Period == 1 {
				//如果是当前年且月份是P1,则默认展示P1的财月和P1财季
				timeValue = 0
			}

			if model.MedaType == excitationModels.MedalMedaTypeOrdinary {
				// 返回当前时间坐标给前端
				wallList.TimeIndex = timeValue
			}

			// 默认为未领取的奖章信息
			wallList.RuleName = model.MedalStyle[timeValue].MedalName
			wallList.Icon = model.MedalStyle[timeValue].UnlitIcon
			var isShow = false
			for key, style := range model.MedalStyle {
				val, ok := memberMedalStyleMap[style.ID]
				if ok {
					// 针对隐藏奖章只展示已领取的最新一个
					if model.MedaType == excitationModels.MedalMedaTypeRare {
						wallList.RuleName = style.MedalName
						wallList.Icon = style.LightIcon
						wallList.IsShow = val.IsShow
					}
					// 判断时间取点亮图标
					if key == int(timeValue) && ok {
						wallList.Icon = style.LightIcon
						wallList.GetTime = val.CreatedAt
					}
				}

				// 总奖章数量 = 普通+隐藏已领取+手动规则
				if model.MedaType == excitationModels.MedalMedaTypeOrdinary ||
					(model.MedaType == excitationModels.MedalMedaTypeRare && ok) ||
					(model.Source == excitationModels.MedalSourceManualUpload) {
					resp.Total++
				}

				// 普通奖章总数
				if model.MedaType == excitationModels.MedalMedaTypeOrdinary {
					wallList.RuleTotal++
				}
				if model.MedaType == excitationModels.MedalMedaTypeRare && ok {
					wallList.RuleTotal++
					wallList.TimeIndex++
				}

				// 奖章new
				if ok {
					// 已获得总奖章数
					resp.Gained++
					// 已获得奖章数
					wallList.RuleGained++

					// 是否展示new标识 如何有有一个未查看 则显示new
					// 如果有未查看则给 循环外isShow赋值标记
					if val.IsShow == excitationModels.MedalRecordIsShowFalse {
						isShow = true
						continue
					}
					// 如果已存在未查看并且当前记录未查看则跳过
					if isShow && val.IsShow == excitationModels.MedalRecordIsShowTrue {
						continue
					}
					wallList.IsShow = excitationModels.MedalRecordIsShowTrue
				}
			}

			if model.MedaType == excitationModels.MedalMedaTypeOrdinary {
				yearData.Ordinary = append(yearData.Ordinary, wallList)
			} else {
				if wallList.TimeIndex > 0 {
					wallList.TimeIndex--
				}
				if wallList.GetTime > 0 {
					yearData.Rare = append(yearData.Rare, wallList)
				}

			}
		}
		resp.List = append(resp.List, yearData)
	}

	return resp, nil
}

// MedalWallDetail 奖章墙详情
func (s *MedalService) MedalWallDetail(c ctx.Context, params excitationReq.MedalWallDetail) (excitationResp.MedalWallDetail, error) {
	var resp excitationResp.MedalWallDetail

	membersUser, err := s.GetUserIdByMobile(c, params.EnterpriseId, params.Mobile)
	if err != nil {
		return resp, err
	}

	enterpriseId := params.EnterpriseId
	ruleId := params.RuleId

	medalDetail, err := excitationDao.NewMedal().GetOneDetail(c, "enterprise_id = ? and id = ?", enterpriseId, ruleId)
	if err != nil {
		return resp, err
	}
	if medalDetail.ID == 0 {
		return resp, fmt.Errorf("数据为空")
	}

	// 获取所有导购总数
	memberTotal, err := membersDao.NewMembers().GetAllTotal(c, map[string]any{"enterprise_id": enterpriseId})
	if err != nil {
		return resp, err
	}

	var startTime, endTime int64
	// 通过财月数据获财年的所有数据
	fiscalMonthYearAll, err := microWorkDao.NewFiscalMonth().GetByYear(c, enterpriseId, params.Year)
	if err != nil || len(fiscalMonthYearAll) == 0 {
		return resp, err
	}
	for _, model := range fiscalMonthYearAll {
		if model.Period == 1 {
			startTime = model.StartTime
		}
		if model.Period == 12 {
			endTime = model.EndTime
		}
	}
	//新的财年第一天发的奖章其实是上个财年的
	startTime += 3600 * 5
	endTime += 3600 * 5

	// 获取已领取的奖章集合
	medalRecordGroup, err := excitationDao.NewMedalRecord().GetReceivedNumGroupStyleId(c, enterpriseId, startTime, endTime)
	if err != nil {
		return resp, err
	}

	medalRecordGroupMap := make(map[int64]int64)
	for _, groupStyleId := range medalRecordGroup {
		medalRecordGroupMap[groupStyleId.MedalStyleId] = groupStyleId.Total
	}

	resp = excitationResp.MedalWallDetail{
		RuleId:          ruleId,
		RuleName:        medalDetail.Name,
		RuleDescription: medalDetail.MedaDescription,
		TopNum:          medalDetail.InviteTop,
		Integral:        medalDetail.Integral,
		ItemName:        medalDetail.ItemName,
	}

	// 获取当前导购次规则已获得的奖章
	memberMedalRecords, err := excitationDao.NewMedalRecord().GetRecordByMemberIdAndMedalId(c, enterpriseId, membersUser.ID, ruleId, startTime, endTime)
	if err != nil {
		return resp, err
	}

	memberMedalRecordsMap := make(map[int64]int64)
	var unReceiveRecordIds []int64 // 未领取的奖章的记录
	for _, record := range memberMedalRecords {
		memberMedalRecordsMap[record.MedalStyleID] = record.CreatedAt
		if record.IsReceive == excitationModels.MedalRecordIsReceiveFalse {
			unReceiveRecordIds = append(unReceiveRecordIds, record.ID)
		}
	}

	var medalStyle []excitationResp.MedalWallDetailStyle
	for _, style := range medalDetail.MedalStyle {
		hasExceeded := (float64(memberTotal) - float64(medalRecordGroupMap[style.ID])) / float64(memberTotal)
		tmp := excitationResp.MedalWallDetailStyle{
			StyleId:             style.ID,
			LightIcon:           style.LightIcon,
			UnlitIcon:           style.UnlitIcon,
			ReceiveNumberPeople: medalRecordGroupMap[style.ID],
			MedalName:           style.MedalName,
		}

		val, ok := memberMedalRecordsMap[style.ID]
		if ok {
			tmp.ReceiveTime = helper.GetFormatDate(val, 0)
			tmp.CreateTime = val
			// 已超过xxx%人，算法：（导购数 - 获取人数）/ 导购数 * 100%
			tmp.HasExceeded = decimal.NewFromFloat(hasExceeded * 100).Round(2).InexactFloat64()
		}

		// 隐藏奖章不展示未领取的其他奖章
		if medalDetail.MedaType == excitationModels.MedalMedaTypeRare && ok {
			medalStyle = append(medalStyle, tmp)
		}
		// 普通奖章全部展示
		if medalDetail.MedaType == excitationModels.MedalMedaTypeOrdinary {
			medalStyle = append(medalStyle, tmp)
		}
	}
	resp.MedalStyle = medalStyle

	// 更新领取记录关联的规则为已查看，如果有未领取的则自动领取，未弹出的标记为已弹出
	err = excitationDao.NewMedalRecord().UpdateData(c,
		map[string]any{
			"enterprise_id": enterpriseId,
			"medal_id":      ruleId,
			"member_id":     membersUser.ID,
		},
		map[string]any{"is_show": excitationModels.MedalRecordIsShowTrue,
			"is_receive": excitationModels.MedalRecordIsReceiveTrue,
			"is_pop":     excitationModels.MedalRecordIsPopTrue})
	if err != nil {
		return resp, err
	}
	return resp, err
}
