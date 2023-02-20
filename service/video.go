package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"videoservice/helpers"
	"videoservice/infra"
)

const Key = "video_service_factory"

type VideoService interface {
	GetFilePathById(id string) (*infra.Video, []byte, error)
	GetFiles() ([]*infra.Video, error)
	DeleteFile(id string) error
	CreateFile(size int, name string, ct string, file multipart.File) (string, error)
}

type videoService struct {
	vd infra.VideoDatabase
}

func NewVideoService(d infra.VideoDatabase) VideoService {
	return &videoService{vd: d}
}

func (vs *videoService) GetFilePathById(id string) (*infra.Video, []byte, error) {
	v, contents, err := vs.vd.GetFile(id)
	if err != nil {
		return nil, contents, err
	}

	return v, contents, nil
}

func (vs *videoService) GetFiles() ([]*infra.Video, error) {
	files, err := vs.vd.GetFiles()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (vs *videoService) DeleteFile(id string) error {
	err := vs.vd.DeleteFile(id)
	if err != nil {
		return err
	}

	return nil
}

func (vs *videoService) CreateFile(size int, name string, ct string,
	file multipart.File) (string, error) {
	randomStr, err := helpers.MakeRandomStr(20)
	if err != nil {
		return "", err
	}
	video := &infra.Video{
		FileName: name,
		FileId:   randomStr,
		Size:     size,
		Type:     ct,
	}
	filePath, err := vs.vd.CreateFile(video, file)
	if err != nil {
		return "", err
	}

	contentLocation := filePath + video.FileId

	return contentLocation, nil
}

func SetFactoryMiddleware(svc VideoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Key, svc)
		c.Next()
	}
}
