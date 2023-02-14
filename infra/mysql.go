package infra

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"xorm.io/xorm"
)

func Init() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", "root:pass@tcp(storage-db:3306)/storage_database?charset=utf8mb4&parseTime=true")
	if err != nil {
		return nil, err
	}

	log.Println("call sync happen")
	err = engine.Sync(new(Video))
	if err != nil {
		log.Println(err)
		log.Println("error is happen")
		return nil, err
	}

	return engine, nil
}
