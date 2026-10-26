package mysql

import (
	"fmt"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/modules/tag/entity"
)

type userGroupRepo struct {
	*repos.SQLRepo
}

var (
	_ tag.UserGroupRepository = &userGroupRepo{}
)

// 检测用户分组名称是否已经被使用
func (u *userGroupRepo) IsExist(c ctx.Context, name string) (bool, error) {
	stmt, tmpErr := u.DB.Template("t_user_group/is-exist-by-name")
	if tmpErr != nil {
		return false, tmpErr
	}

	result := struct {
		UseNum int `db:"use_num"`
	}{}
	resErr := stmt.NamedGet(c, &result, map[string]interface{}{"name": name})
	if resErr != nil {
		return false, resErr
	}

	return result.UseNum > 0, nil
}

// 新增客户分组
func (u *userGroupRepo) AddUserGroup(c ctx.Context, data entity.UserGroupModelEntity) (int64, error) {
	stmt, tmpErr := u.DB.Template("t_user_group/insert")
	if tmpErr != nil {
		return 0, tmpErr
	}

	res, resErr := stmt.NamedExec(c, data)
	if resErr != nil {
		return 0, resErr
	}

	return res.RowsAffected()
}

func (u *userGroupRepo) UpdateUserGroup(c ctx.Context, enterpriseID int64, id int64, attrs map[string]interface{}) error {
	if len(attrs) == 0 {
		return nil
	}
	sqlstr := "update user_group set %s where enterprise_id=? and id=?"
	var conds []string
	var args []interface{}
	for k, v := range attrs {
		conds = append(conds, fmt.Sprintf("%s=?", k))
		args = append(args, v)
	}
	conds = append(conds, "updated_at=?")
	args = append(args, time.Now().Unix())
	sqlstr = fmt.Sprintf(sqlstr, strings.Join(conds, ","))
	args = append(args, enterpriseID, id)
	_, err := u.DB.DB.Exec(c, sqlstr, args...)
	return err
}

// 新增客户详情
func (u *userGroupRepo) AddUserGroupDetail(c ctx.Context, data []entity.UserGroupDetailModelEntity) (int64, error) {
	stmt, tmpErr := u.DB.Template("t_user_group_detail/insert")
	if tmpErr != nil {
		return 0, tmpErr
	}

	res, resErr := stmt.NamedExec(c, data)
	if resErr != nil {
		return 0, resErr
	}

	return res.RowsAffected()
}

// 查询说有的客户分组
func (u *userGroupRepo) GetUserGroupList(c ctx.Context, enterpriseID int64) ([]entity.UserGroupInfo, error) {
	stmt, tmpErr := u.DB.Template("t_user_group/select-user-group-all")
	if tmpErr != nil {
		return nil, tmpErr
	}

	result := []entity.UserGroupInfo{}
	resErr := stmt.NamedSelect(c, &result, map[string]interface{}{"enterprise_id": enterpriseID})
	if resErr != nil {
		return nil, resErr
	}

	return result, nil
}

// 通过客户分组ID，获取外部联系人集合
func (u *userGroupRepo) GetExternalUserIDList(c ctx.Context, param entity.QueryUserGroupParam) ([]int64, error) {
	stmt, tmpErr := u.DB.Template("t_user_group_detail/select-external-user-id-list")
	if tmpErr != nil {
		return nil, tmpErr
	}

	rows, resErr := stmt.NamedQuery(c, param)
	if resErr != nil {
		return nil, resErr
	}
	// 遍历获取外部联系人集合ID
	defer rows.Close()

	IDList := []int64{}
	for rows.Next() {
		data := struct {
			ID int64 `db:"work_external_user_id"`
		}{}

		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}

		IDList = append(IDList, data.ID)
	}

	return IDList, nil
}

// 获取客户分组详情
func (u *userGroupRepo) GetUserGroupDetail(c ctx.Context, param entity.QueryUserGroupParam) ([]entity.UserGroupDetailTagInfo, error) {
	// 初始化返回值
	result := make([]entity.UserGroupDetailTagInfo, 0)

	stmt, tmpErr := u.DB.Template("t_user_group_detail/select-tag-info-list")
	if tmpErr != nil {
		return result, tmpErr
	}

	resErr := stmt.NamedSelect(c, &result, param)
	if resErr != nil {
		return result, tmpErr
	}

	return result, nil
}

// 根据标签ID集合，获取关联该标签的用户集合
func (u *userGroupRepo) GetExternalUserIDByTagIDList(c ctx.Context, tagIDList []int64, enterpriseID int64) ([]int64, error) {
	userIDList := []int64{}

	stmt, tmpErr := u.DB.Template("t_work_external_user_tag/user-group-tag-user-id-list", &map[string]interface{}{
		"IDList": tagIDList,
	})
	if tmpErr != nil {
		return userIDList, tmpErr
	}

	rows, rowErr := stmt.NamedQuery(c, map[string]interface{}{"enterprise_id": enterpriseID})
	if rowErr != nil {
		return userIDList, rowErr
	}
	defer rows.Close()

	for rows.Next() {
		data := struct {
			ID int64 `db:"work_external_user_id"`
		}{}

		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}

		userIDList = append(userIDList, data.ID)
	}

	return userIDList, nil
}
