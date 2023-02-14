package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
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
	fs infra.FileServer
}

func NewVideoService(d infra.VideoDatabase, f infra.FileServer) VideoService {
	return &videoService{vd: d, fs: f}
}

func (vs *videoService) GetFilePathById(ctx context.Context, id int) (string, error) {
	v, err := vs.vd.GetFile(id)
	if err != nil {
		return "", err
	}

	filePath := vs.fs.GetFilePath(v.Id) + v.FileName

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
	video, err := vs.vd.DeleteFile(id)
	if err != nil {
		return err
	}

	err = vs.fs.DeleteFile(video.FileName, video.Id)
	if err != nil {
		return err
	}

	return nil
}

func (vs *videoService) CreateFile(ctx context.Context, size int, name string, ct string,
	file multipart.File) error {
	v, session, err := vs.vd.CreateFile(
		&infra.Video{
			FileName: name,
			Size:     size,
			Type:     ct,
		})
	if err != nil {
		return err
	}
	defer session.Close()

	err = vs.fs.StoreFile(name, v.Id, v.Size, file)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()

	return nil
}

func ValidateAndResponseContentType(file multipart.File) (string, error) {
	// To sniff the content type only the first, 512 bytes are used.
	buf := make([]byte, 512)

	_, err := file.Read(buf)
	if err != nil {
		return "", nil
	}

	file.Seek(0, 0)
	ct := http.DetectContentType(buf)
	return ct, nil
}

func SetFactoryMiddleware(svc VideoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Key, svc)
		c.Next()
	}
}
