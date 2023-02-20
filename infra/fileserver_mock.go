package infra

import (
	"mime/multipart"
)

type TestFileServerMock struct {
}

func (fs TestFileServerMock) StoreFile(name string, id int, content multipart.File) (string, error) {
	return "", nil
}

func (fs TestFileServerMock) DeleteFile(name string, id int) error {
	return nil
}
func (fs TestFileServerMock) GetFileContent(name string, id int) ([]byte, error) {
	return []byte{}, nil
}
func (fs TestFileServerMock) GetFilePath(id int) string {
	return ""
}
