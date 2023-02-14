package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"videoservice/helpers"
	"videoservice/model"

	"videoservice/service"
)

type ServerHandler interface {

	// (GET /files)
	GetFiles(c *gin.Context)

	// (POST /files)
	PostFiles(c *gin.Context)

	// (DELETE /files/{fileid})
	DeleteFilesFileId(c *gin.Context)

	// (GET /files/{fileid})
	GetFilesFileId(c *gin.Context)

	// (GET /health)
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := model.ConvertVideoToVideoResponse(videos)
	c.JSON(http.StatusOK, res)

}

// PostFiles store video file
func (sh *serverHandler) PostFiles(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("data")
	defer file.Close()

	ft := fileHeader.Header.Get("Content-Type")
	if ft == "" || ft != "video/mp4" && ft != "video/mpeg" {
		log.Println(ft)
		c.AbortWithError(http.StatusInternalServerError, errors.New(helpers.ContentTypeIsWrong))
		return
	}

	s := c.MustGet(service.Key).(service.VideoService)
	err = s.CreateFile(c, int(fileHeader.Size), fileHeader.Filename, ft, file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteFilesFileId delete video file with id
func (sh *serverHandler) DeleteFilesFileId(c *gin.Context) {
	id := c.Param("id")
	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New(helpers.ParamIsInvalid))
		return
	}

	s := c.MustGet(service.Key).(service.VideoService)
	err = s.DeleteFile(c, fileId)
	if err != nil {
		c.AbortWithError(helpers.GetStatusCodeFromErr(err), err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFilesFileId get file with file id
func (sh *serverHandler) GetFilesFileId(c *gin.Context) {
	id := c.Param("id")
	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New(helpers.ParamIsInvalid))
		return
	}

	s := c.MustGet(service.Key).(service.VideoService)
	filePath, err := s.GetFilePathById(c, fileId)
	if err != nil {
		c.AbortWithError(helpers.GetStatusCodeFromErr(err), err)
		return
	}

	c.File(filePath)
	c.Status(http.StatusOK)
}

func (sh *serverHandler) GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "hello!")
}
