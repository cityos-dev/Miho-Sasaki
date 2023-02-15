package infra

import (
	"log"
	"time"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
)

var retry = 0

func Init() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", "root:pass@tcp(storage-db:3306)/storage_database?charset=utf8mb4&parseTime=true")
	if err != nil {
		return nil, err
	}

	log.Println("call sync happen")
	err = engine.Sync(new(Video))
	if err != nil {
		log.Println(err)
		log.Println("db is not ready. Retry connecting...")
		time.Sleep(time.Second * 3)
		retry++
		if retry > 10 {
			return nil, err
		}

		return Init()
	}

	return engine, nil
}
