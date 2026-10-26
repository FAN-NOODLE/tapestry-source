package mysql

import (
	"fmt"

	"micro-scrm/internal/modules/scripts"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type SysRepo struct {
	*repos.SQLRepo
}

// SpiteTable 创建表
func (s *SysRepo) SpiltTable(cc ctx.Context, enterpriseID int64) error {
	mapping := make(map[string]interface{})
	mapping["EnterpriseID"] = enterpriseID
	stmt, err := s.DB.Template("t_split_table/split-table", &mapping)
	if err != nil {
		log.Error(cc).Msgf("切表 分表SQL：%s", stmt.SQLRaw())
		return err
	}
	_, nErr := stmt.Query(cc)
	if nErr != nil {
		return nErr
	}
	return nil
}

// SpiteTable 创建数据库
func (s *SysRepo) SpiltDatabase(cc ctx.Context, enterpriseID int64) error {
	mapping := make(map[string]interface{})
	mapping["EnterpriseID"] = enterpriseID
	stmt, err := s.DB.Template("t_split_database/split-database", &mapping)
	if err != nil {
		log.Error(cc).Msgf("切库 分库SQL：%s", stmt.SQLRaw())
		return err
	}
	_, nErr := stmt.Query(cc)
	if nErr != nil {
		return nErr
	}
	return nil
}

func (s *SysRepo) GetRepeatData(cc ctx.Context) ([]scripts.Data, error) {
	list := make([]scripts.Data, 0)
	sqlStr := fmt.Sprintf("select id,count(1) as num from work_external_user_relation where id>0 group by id")
	err := s.DB.Select(cc, &list, sqlStr)
	if err != nil {
		log.Error(cc).Msgf("查询重复数据：%+v", err)
		return list, err
	}
	return list, nil
}

func (s *SysRepo) DelRepeatData(cc ctx.Context, id int64, number int64) error {
	sqlStr := fmt.Sprintf("delete from work_external_user_relation where id=? limit ?")
	_, err := s.DB.DB.Exec(cc, sqlStr, id, number)
	if err != nil {
		log.Error(cc).Msgf("删除重复数据：%+v", err)
		return err
	}
	return nil
}

// 查询有备注的外部联系人关系
func (s *SysRepo) GetHasRemarkData(cc ctx.Context) ([]scripts.ExternalUserRelationData, error) {
	list := make([]scripts.ExternalUserRelationData, 0)
	sqlStr := fmt.Sprintf("select enterprise_id,work_external_user_id,description from work_external_user_relation where relation_status=0 AND description!=''")
	err := s.DB.Select(cc, &list, sqlStr)
	if err != nil {
		log.Error(cc).Msgf("查询有备注的外部联系人关系数据：%+v", err)
		return list, err
	}
	return list, nil
}

func (s *SysRepo) UpdateRemarkData(cc ctx.Context, enterpriseID int64, id int64, remark string) error {
	sqlStr := fmt.Sprintf("update micro_scrm_enterprise_%d.enterprise_user set remark=? where id = ?", enterpriseID)
	_, err := s.DB.DB.Exec(cc, sqlStr, remark, id)
	if err != nil {
		log.Error(cc).Msgf("更新备注数据：%+v", err)
		return err
	}
	return nil
}
