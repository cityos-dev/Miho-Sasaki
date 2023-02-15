package model

import (
	"time"

	"videoservice/infra"
)

type GetVideo struct {
	FileId    string    `json:"fileid"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

func ConvertVideoToVideoResponse(videos []*infra.Video) []GetVideo {
	response := make([]GetVideo, 0)
	for _, v := range videos {
		v := v

		response = append(response,
			GetVideo{
				FileId:    v.FileId,
				Name:      v.FileName,
				Size:      v.Size,
				CreatedAt: v.Created,
			})
	}

	return response
}
