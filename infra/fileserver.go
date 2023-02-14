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
	GetFile(name string, id int, size int) ([]byte, error)
	StoreFile(name string, id int, size int, content multipart.File) error
	DeleteFile(name string, id int) error
	GetFilePath(id int) string
}

type fileServer struct {
	filePath string // /video
}

func NewFileServer(path string) FileServer {
	return &fileServer{filePath: path}
}

// GetFile filepath example: ../contents + /video + /id/ + name
func (fs *fileServer) GetFile(name string, id int, size int) ([]byte, error) {
	f, err := os.Open(fs.GetFilePath(id) + name)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data := make([]byte, size)
	length, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	log.Println("length  ")
	log.Println(length)

	return data, nil
}

func (fs *fileServer) StoreFile(name string, id int, size int, content multipart.File) error {
	err := os.MkdirAll(fs.GetFilePath(id), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fs.GetFilePath(id) + name)
	if err != nil {
		return err
	}
	log.Println("path")
	log.Println(fs.GetFilePath(id) + name)

	defer f.Close()
	written, err := io.Copy(f, content)
	if err != nil {
		return err
	}
	log.Println("written  ")
	log.Println(written)

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
