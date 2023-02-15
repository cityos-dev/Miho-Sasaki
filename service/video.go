package service

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"videoservice/helpers"
	"videoservice/infra"
)

const Key = "video_service_factory"

type VideoService interface {
	GetFilePathById(ctx context.Context, id string) (*infra.Video, string, error)
	GetFiles(ctx context.Context) ([]*infra.Video, error)
	DeleteFile(ctx context.Context, id string) error
	CreateFile(ctx context.Context, size int, name string, ct string, file multipart.File) (string, error)
}

type videoService struct {
	vd infra.VideoDatabase
}

func NewVideoService(d infra.VideoDatabase) VideoService {
	return &videoService{vd: d}
}

func (vs *videoService) GetFilePathById(ctx context.Context, id string) (*infra.Video, string, error) {
	v, err := vs.vd.GetFile(id)
	if err != nil {
		return nil, "", err
	}

	filePath := vs.vd.GetFilePathBy(v) + "/" + v.FileName

	return v, filePath, nil
}

func (vs *videoService) GetFiles(ctx context.Context) ([]*infra.Video, error) {
	files, err := vs.vd.GetFiles()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (vs *videoService) DeleteFile(ctx context.Context, id string) error {
	err := vs.vd.DeleteFile(id)
	if err != nil {
		return err
	}

	return nil
}

func (vs *videoService) CreateFile(ctx context.Context, size int, name string, ct string,
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
