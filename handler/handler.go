package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"videoservice/helpers"
	"videoservice/model"
	"videoservice/service"
)

type ServerHandler interface {

	//GetFiles (GET /files)
	GetFiles(c *gin.Context)

	//PostFiles (POST /files)
	PostFiles(c *gin.Context)

	//DeleteFilesFileId (DELETE /files/{fileid})
	DeleteFilesFileId(c *gin.Context)

	//GetFilesFileId (GET /files/{fileid})
	GetFilesFileId(c *gin.Context)

	//GetHealth (GET /health)
	GetHealth(c *gin.Context)
}

type serverHandler struct {
	errorHandler func(*gin.Context, error, int)
}

func NewServerHandler() ServerHandler {
	return &serverHandler{
		errorHandler: func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		},
	}
}

// GetFiles get video files
func (sh *serverHandler) GetFiles(c *gin.Context) {
	s := c.MustGet(service.Key).(service.VideoService)
	videos, err := s.GetFiles(c)
	if err != nil {
		sh.errorHandler(c, err, http.StatusInternalServerError)
		return
	}

	res := model.ConvertVideoToVideoResponse(videos)
	c.JSON(http.StatusOK, res)
}

// PostFiles store video file
func (sh *serverHandler) PostFiles(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("data")
	if err != nil {
		log.Println(err)
		sh.errorHandler(c, errors.New(helpers.BadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ft := fileHeader.Header.Get("Content-Type")
	if ft == "" || ft != "video/mp4" && ft != "video/mpeg" {
		log.Println(ft)
		sh.errorHandler(c, errors.New(helpers.UnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	s := c.MustGet(service.Key).(service.VideoService)
	contentLocation, err := s.CreateFile(c, int(fileHeader.Size), fileHeader.Filename, ft, file)
	if err != nil {
		sh.errorHandler(c, err, helpers.GetStatusCodeFromErr(err))
		return
	}

	c.Header("Location", contentLocation)
	c.Status(http.StatusCreated)
}

// DeleteFilesFileId delete video file with id
func (sh *serverHandler) DeleteFilesFileId(c *gin.Context) {
	fileId := c.Param("fileid")

	s := c.MustGet(service.Key).(service.VideoService)
	err := s.DeleteFile(c, fileId)
	if err != nil {
		sh.errorHandler(c, err, helpers.GetStatusCodeFromErr(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFilesFileId get file with file id
func (sh *serverHandler) GetFilesFileId(c *gin.Context) {
	id := c.Param("fileid")

	s := c.MustGet(service.Key).(service.VideoService)
	v, filePath, err := s.GetFilePathById(c, id)
	if err != nil {
		sh.errorHandler(c, err, helpers.GetStatusCodeFromErr(err))
		return
	}

	c.Header("Content-Type", v.Type)
	c.Writer.Header().Set("Content-Disposition", `attachment; filename="`+v.FileName+`"`)
	c.File(filePath + v.FileName)
}

func (sh *serverHandler) GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
