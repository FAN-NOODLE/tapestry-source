package mgo

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
)

type MgoRepo struct {
}

func New() *MgoRepo {
	return &MgoRepo{}
}

func (m *MgoRepo) Add(vo interface{}) (sql.Result, error) {
	panic("implement me")
}

func (m *MgoRepo) Del(id typex.ID) (sql.Result, error) {
	panic("implement me")
}

func (m *MgoRepo) Update(entity interface{}) (sql.Result, error) {
	panic("implement me")
}

func (m *MgoRepo) UpdateActive(entity interface{}) (sql.Result, error) {
	panic("implement me")
}

func (m *MgoRepo) Find(vo interface{}) (*sqlx.Rows, error) {
	panic("implement me")
}

func (m *MgoRepo) Detail(vo interface{}, id typex.ID) error {
	panic("implement me")
}
