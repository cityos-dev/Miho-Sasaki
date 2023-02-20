package infra

import (
	"fmt"
	"regexp"
	"testing"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/core"
	xormLog "xorm.io/xorm/log"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
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
	engine.Logger().SetLevel(xormLog.LOG_INFO)

	return engine, mock, nil
}

type testVideoDB struct {
	suite.Suite
}

func TestVideoDB(t *testing.T) {
	suite.Run(t, new(testVideoDB))
}

func (t *testVideoDB) TestVideoGetFiles() {
	fs := TestFileServerMock{}
	a := t.Assert()
	engine, mock, err := getDBMock()
	if err != nil {
		a.Error(err)
	}

	videoDB := NewVideoDatabase(engine, fs)

	t.Run("test GetFiles()", func() {
		rows := sqlMock.NewRows([]string{"id", "file_id",
			"file_name", "size", "type", "created_at", "updated_at"}).AddRow(
			1, "UQDIWWMNPPQIDMWQNeow", "sample.mpg", 6256514, "video/mpeg", time.Time{}, time.Time{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM video`)).WillReturnRows(rows)

		videos, err := videoDB.GetFiles()
		if err != nil {
			a.Error(err)
		}

		fmt.Println("videos")
		fmt.Println(videos)

		if err := mock.ExpectationsWereMet(); err != nil {
			a.Errorf(err, "Test Get Videos: %v")
		}

	})
}

func (t *testVideoDB) TestCreateFile() {
	fmt.Println("hello")
	fs := TestFileServerMock{}
	a := t.Assert()
	engine, mock, err := getDBMock()
	if err != nil {
		a.Error(err)
	}

	videoDB := NewVideoDatabase(engine, fs)

	t.Run("test CreateFile()", func() {
		id := 1
		fileId := "UQDIWWMNPPdnosmDNF"
		expectQuery := "INSERT INTO video (id, file_id, file_name, size, type, created_at, updated_at) " +
			"VALUES ($1,$2,$3,$4,$5,$6,$7)"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(expectQuery)).WithArgs(
			id, fileId, 6256514, "video/mpeg", time.Time{}, time.Time{}).WillReturnResult(
			sqlMock.NewResult(int64(id), 1))
		mock.ExpectCommit()

		_, err := videoDB.CreateFile(&Video{
			Id:       id,
			FileId:   fileId,
			FileName: "sample.mpg",
			Size:     6256514,
			Type:     "video/mpeg",
			Created:  time.Time{},
			Updated:  time.Time{},
		}, nil)
		if err != nil {
			a.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			a.Errorf(err, "Test GetFile(): %v")
		}
	})
}

func (t *testVideoDB) TestVideoGetWithId() {
	fs := TestFileServerMock{}
	a := t.Assert()
	engine, mock, err := getDBMock()
	if err != nil {
		a.Error(err)
	}

	videoDB := NewVideoDatabase(engine, fs)

	t.Run("test GetFile() with fileId", func() {
		fileId := "UQDIWWMNPPQIDMWQNeow"
		rows := sqlMock.NewRows([]string{"id", "file_id",
			"file_name", "size", "type", "created_at", "updated_at"}).AddRow(
			1, fileId, "sample.mpg", 6256514, "video/mpeg", time.Time{}, time.Time{})
		expectQuery := "SELECT `id`, `file_id`, `file_name`, `size`, `type`, `created`," +
			" `updated` FROM `video` WHERE (file_id=?) LIMIT 1"

		mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WithArgs(fileId).WillReturnRows(rows)

		video, _, err := videoDB.GetFile(fileId)
		if err != nil {
			a.Error(err)
		}

		a.Equal(fileId, video.FileId)

		if err := mock.ExpectationsWereMet(); err != nil {
			a.Errorf(err, "Test GetFile(): %v")
		}
	})
}

func (t *testVideoDB) TestVideoDeleteWithId() {
	fs := TestFileServerMock{}
	a := t.Assert()
	engine, mock, err := getDBMock()
	if err != nil {
		a.Error(err)
	}

	videoDB := NewVideoDatabase(engine, fs)

	t.Run("test DeleteFile() with fileId", func() {
		fileId := "UQDIWWMNPPQIDMWQNeow"
		expectQuery := "DELETE FROM video WHERE (file_id=?)"

		mock.ExpectExec(regexp.QuoteMeta(expectQuery)).WithArgs(fileId).WillReturnResult(
			sqlMock.NewResult(1, 1))

		err := videoDB.DeleteFile(fileId)
		if err != nil {
			a.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			a.Errorf(err, "Test DeleteFile(): %v")
		}
	})
}
