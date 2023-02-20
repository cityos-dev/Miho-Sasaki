package infra

import (
	"xorm.io/xorm"
	core "xorm.io/xorm/core"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
)

func getDBMock() (*xorm.Engine, sqlMock.Sqlmock, error) {
	db, mock, err := sqlMock.New()
	if err != nil {
		return nil, nil, err
	}

	engine, err := xorm.NewEngineWithDB("mysql", "mock", core.FromDB(db))
	if err != nil {
		return nil, nil, err
	}

	return engine, mock, nil
}
