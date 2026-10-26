package service

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao"
)

type DemoService struct {
	Base
	demoDao *dao.Demo
}

func NewDemoService() *DemoService {
	return &DemoService{
		demoDao: dao.NewDemo(),
	}
}

func (s *DemoService) GetList(c ctx.Context) {

}
