package uc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
)

type Department struct {
}

func NewDepartment() *Department {
	return &Department{}
}

// List 部门列表
func (d Department) List(ctx ctx.Context, request interface{}, identifier string, domain string) ([]DepartmentInfoResp, error) {
	url := fmt.Sprintf("%s/a/%s/%s/qy/department/ignore-permissions/page-list", config.CfgData.Service.Work.Domain, domain, identifier)
	var response struct {
		Data DepartmentListResp `json:"data"`
	}
	err := Post(ctx, url, request, &response)
	if err != nil {
		return response.Data.List, err
	}
	return response.Data.List, nil
}
