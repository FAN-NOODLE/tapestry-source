package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type MiniProgramRepo struct {
	*repos.SQLRepo
}

// AuthorizeMiniProgramDetail 授权小程序详情
func (r *MiniProgramRepo) AuthorizeMiniProgramDetail(cc ctx.Context, enterpriseId int64) (int64, error) {
	stmt, err := r.DB.Template("t_mini_program/authorize-mini-program-detail")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Query(cc, enterpriseId)
	if err != nil {
		return 0, err
	}

	id := new(int64)
	defer res.Close()
	if res.Next() {
		err = res.Scan(id)
	}
	if err != nil {
		return 0, err
	}

	return *id, err
}
