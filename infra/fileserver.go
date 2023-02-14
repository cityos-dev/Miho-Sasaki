package infra

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

const contentsPath = "../contents"

type FileServer interface {
	StoreFile(name string, id int, content multipart.File) error
	DeleteFile(name string, id int) error
	GetFilePath(id int) string
}

type fileServer struct {
	filePath string // /video
}

func NewFileServer(path string) FileServer {
	return &fileServer{filePath: path}
}

func (fs *fileServer) StoreFile(name string, id int, content multipart.File) error {
	err := os.MkdirAll(fs.GetFilePath(id), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fs.GetFilePath(id) + name)
	if err != nil {
		return err
	}

	defer f.Close()
	written, err := io.Copy(f, content)
	if err != nil {
		return err
	}
	log.Printf("%d bytes are written", written)

	return nil
}

func (fs *fileServer) DeleteFile(name string, id int) error {
	err := os.Remove(fs.GetFilePath(id) + name)
	if err != nil {
		return err
	}

	err = os.Remove(fs.GetFilePath(id))
	if err != nil {
		return err
	}

	return nil
}

func (fs *fileServer) GetFilePath(id int) string {
	return contentsPath + fs.filePath + "/" + strconv.Itoa(id) + "/"
}
