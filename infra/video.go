package infra

import (
	"errors"
	"log"
	"mime/multipart"
	"time"
	"xorm.io/xorm"

	"videoservice/helpers"
)

const tableName = "video"

type Video struct {
	Id       int    `xorm:"'id' pk autoincr"`
	FileId   string `xorm:"unique"`
	FileName string
	Size     int
	Type     string
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

type VideoDatabase interface {
	GetFiles() ([]*Video, error)
	GetFile(id string) (*Video, error)
	CreateFile(video *Video, file multipart.File) (string, error)
	DeleteFile(id string) error
	GetFilePathBy(v *Video) string
}

type videoDatabase struct {
	engine     *xorm.Engine
	fileServer FileServer
}

func NewVideDatabase(en *xorm.Engine, fs FileServer) VideoDatabase {
	return &videoDatabase{
		engine:     en,
		fileServer: fs,
	}
}

func (vd *videoDatabase) GetFiles() ([]*Video, error) {
	var video []*Video
	if err := vd.engine.Table(tableName).Find(&video); err != nil {
		return nil, err
	}

	return video, nil
}

func (vd *videoDatabase) GetFile(id string) (*Video, error) {
	var video Video
	found, err := vd.engine.Table(tableName).Where("file_id=?", id).Get(&video)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !found {
		return nil, errors.New(helpers.FileNotFound)
	}

	return &video, nil
}

func (vd *videoDatabase) CreateFile(video *Video, file multipart.File) (string, error) {
	session := vd.engine.NewSession()
	err := session.Begin()
	if err != nil {
		return "", err
	}
	defer session.Close()

	if _, err := vd.engine.Table(tableName).Insert(video); err != nil {
		return "", err
	}

	filePath, err := vd.fileServer.StoreFile(video.FileName, video.Id, file)
	if err != nil {
		session.Rollback()
		return "", err
	}

	session.Commit()

	return filePath, nil
}

func (vd *videoDatabase) DeleteFile(id string) error {
	session := vd.engine.NewSession()
	err := session.Begin()
	if err != nil {
		return err
	}
	defer session.Close()

	var video Video
	found, err := vd.engine.Table(tableName).Where("file_id=?", id).Get(&video)
	if err != nil {
		log.Println(err)
		return err
	}

	if !found {
		return errors.New(helpers.FileNotFound)
	}

	affected, err := vd.engine.Table(tableName).ID(video.Id).Delete(&Video{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(helpers.FileNotFound)
	}

	err = vd.fileServer.DeleteFile(video.FileName, video.Id)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

func (vd *videoDatabase) GetFilePathBy(v *Video) string {
	return vd.fileServer.GetFilePath(v.Id) + v.FileName
}
