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

func (t *testFileServer) TestFileOperations() {
	testFilePath := "../testing/contents"
	testContentsPath := "/video"
	fs := NewFileServer(testFilePath, testContentsPath)
	cases := testCasesForFiles
	a := t.Assert()

	for name, v := range cases {
		t.Run(name, func() {
			file, err := openFile(v.filePath)
			if err != nil {
				t.Error(err)
			}

			filePath, err := fs.StoreFile(v.videoData.FileName, v.videoData.Id, file)
			if err != nil {
				t.Error(err)
			}

			a.Equal(testFilePath+testContentsPath+"/"+strconv.Itoa(v.videoData.Id)+"/", filePath)

			contents, err := fs.GetFileContent(v.videoData.FileName, v.videoData.Id)
			if err != nil {
				t.Error(err)
			}

			a.EqualValues(v.videoData.Size, len(contents))

			err = fs.DeleteFile(v.videoData.FileName, v.videoData.Id)
			a.Nil(err)
		})
	}
}

func openFile(filePath string) (multipart.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return f, err
}

var testCasesForFiles = map[string]struct {
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
