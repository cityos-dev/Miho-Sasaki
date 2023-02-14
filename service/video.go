package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"videoservice/infra"
)

const Key = "video_service_factory"

type VideoService interface {
	GetFilePathById(ctx context.Context, id int) (string, error)
	GetFiles(ctx context.Context) ([]*infra.Video, error)
	DeleteFile(ctx context.Context, id int) error
	CreateFile(ctx context.Context, size int, name string, ct string, file multipart.File) error
}

type videoService struct {
	vd infra.VideoDatabase
}

func NewVideoService(d infra.VideoDatabase) VideoService {
	return &videoService{vd: d}
}

func (vs *videoService) GetFilePathById(ctx context.Context, id int) (string, error) {
	v, err := vs.vd.GetFile(id)
	if err != nil {
		return "", err
	}

	filePath := vs.vd.GetFilePathBy(v)

	return filePath, nil
}

func (vs *videoService) GetFiles(ctx context.Context) ([]*infra.Video, error) {
	files, err := vs.vd.GetFiles()
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (vs *videoService) DeleteFile(ctx context.Context, id int) error {
	err := vs.vd.DeleteFile(id)
	if err != nil {
		return err
	}

	return nil
}

func (vs *videoService) CreateFile(ctx context.Context, size int, name string, ct string,
	file multipart.File) error {
	err := vs.vd.CreateFile(
		&infra.Video{
			FileName: name,
			Size:     size,
			Type:     ct,
		},
		file,
	)
	if err != nil {
		return err
	}

	return nil
}

func SetFactoryMiddleware(svc VideoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Key, svc)
		c.Next()
	}
}
