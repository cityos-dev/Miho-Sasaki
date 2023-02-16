package infra

import (
	"github.com/stretchr/testify/suite"
	"mime/multipart"
	"os"
	"strconv"
	"testing"
	"time"
)

type testFileServer struct {
	suite.Suite
}

func TestFileServer(t *testing.T) {
	suite.Run(t, new(testFileServer))
}

func (t *testFileServer) TestStoreFile() {
	testFilePath := "../testing/contents"
	testContentsPath := "/video"
	fileServer := NewFileServer(testFilePath, testContentsPath)
	cases := map[string]struct {
		filePath  string
		videoData Video
	}{
		"check mp4 file": {
			filePath: "../test/post_1/sample.mp4",
			videoData: Video{
				Id:       1,
				FileId:   "UQDIWWMNPPQIMGEGMGKB",
				FileName: "sample.mp4",
				Size:     2848208,
				Type:     "video/mp4",
				Created:  time.Time{},
				Updated:  time.Time{},
			},
		},
		"check mpg file": {
			filePath: "../test/post_3/sample.mpg",
			videoData: Video{
				Id:       1,
				FileId:   "UQDIWWMNPPQIDMWQNeow",
				FileName: "sample.mpg",
				Size:     6256514,
				Type:     "video/mpeg",
				Created:  time.Time{},
				Updated:  time.Time{},
			},
		},
	}

	a := t.Assert()

	for name, v := range cases {
		t.Run(name, func() {
			file, err := createFile(v.filePath)
			if err != nil {
				t.Error(err)
			}

			filePath, err := fileServer.StoreFile(v.videoData.FileName, v.videoData.Id, file)
			if err != nil {
				t.Error(err)
			}

			a.Equal(testFilePath+testContentsPath+"/"+strconv.Itoa(v.videoData.Id)+"/", filePath)

			storedFile, err := os.Open(filePath + v.videoData.FileName)
			if err != nil {
				t.Error(err)
			}

			info, err := storedFile.Stat()
			if err != nil {
				t.Error(err)
			}

			a.EqualValues(v.videoData.Size, info.Size())
			a.EqualValues(v.videoData.FileName, info.Name())

			err = removeFileDirectory(filePath, v.videoData.FileName)
			if err != nil {
				a.Errorf(err, "Failed to remove file")
			}
		})
	}
}

func createFile(filePath string) (multipart.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return f, err
}

func removeFileDirectory(filePath, fileName string) error {
	err := os.Remove(filePath + fileName)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
