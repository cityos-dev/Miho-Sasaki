package infra

import (
	"errors"
	"fmt"
	"log"
	"time"
	"videoservice/helpers"

	"xorm.io/xorm"
)

const tableName = "video"

type Video struct {
	Id       int `xorm:"'id' pk autoincr"`
	FileName string
	Size     int
	Type     string
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

type VideoDatabase interface {
	GetFiles() ([]*Video, error)
	GetFile(id int) (*Video, error)
	CreateFile(video *Video) (*Video, *xorm.Session, error)
	DeleteFile(id int) (*Video, error)
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

func (vd *videoDatabase) GetFile(id int) (*Video, error) {
	var video Video
	found, err := vd.engine.Table(tableName).ID(id).Get(&video)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !found {
		return nil, errors.New(helpers.FileNotFound)
	}

	return &video, nil
}

func (vd *videoDatabase) CreateFile(video *Video) (*Video, *xorm.Session, error) {
	session := vd.engine.NewSession()
	err := session.Begin()
	if err != nil {
		return nil, nil, err
	}

	if _, err := vd.engine.Table(tableName).Insert(video); err != nil {
		session.Close()
		return nil, nil, err
	}

	return video, session, nil
}

func (vd *videoDatabase) DeleteFile(id int) (*Video, error) {
	session := vd.engine.NewSession()
	err := session.Begin()
	if err != nil {
		return nil, err
	}

	video := Video{}
	affected, err := vd.engine.Table(tableName).ID(id).Delete(&video)
	if err != nil {
		return nil, err
	}

	fmt.Println("affected")
	fmt.Println(affected)
	fmt.Println(video)
	if affected == 0 {
		return nil, errors.New(helpers.FileNotFound)
	}

	return &video, nil
}
