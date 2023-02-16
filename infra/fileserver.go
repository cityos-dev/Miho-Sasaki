package infra

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"videoservice/helpers"
)

const contentsPath = "../contents"

type FileServer interface {
	StoreFile(name string, id int, content multipart.File) (string, error)
	DeleteFile(name string, id int) error
	GetFileContent(name string, id int) ([]byte, error)
	GetFilePath(id int) string
}

type fileServer struct {
	contentsPath string
	filePath     string
}

func NewFileServer(cp, fp string) FileServer {
	return &fileServer{contentsPath: cp, filePath: fp}
}

func (fs *fileServer) StoreFile(name string, id int, content multipart.File) (string, error) {
	filePath := fs.GetFilePath(id)
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filePath + name)
	if err != nil {
		return "", err
	}

	defer f.Close()
	written, err := io.Copy(f, content)
	if err != nil {
		return "", err
	}
	log.Printf("%d bytes are written", written)

	return filePath, nil
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

func (fs *fileServer) GetFileContent(name string, id int) ([]byte, error) {
	f, err := os.Open(fs.GetFilePath(id) + name)
	if err != nil {
		return nil, errors.New(helpers.FileNotFound)
	}
	defer f.Close()
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fstat.Size())
	reads, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	log.Printf("%d bytes are readed", reads)

	return buf, nil
}

func (fs *fileServer) GetFilePath(id int) string {
	return fs.contentsPath + fs.filePath + "/" + strconv.Itoa(id) + "/"
}
