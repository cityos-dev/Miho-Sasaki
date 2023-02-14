package model

import (
	"time"

	"videoservice/infra"
)

type GetVideo struct {
	FileId    int
	Name      string
	Size      int
	CreatedAt time.Time
}

func ConvertVideoToVideoResponse(videos []*infra.Video) []GetVideo {
	response := make([]GetVideo, 0)
	for _, v := range videos {
		v := v

		response = append(response,
			GetVideo{
				FileId:    v.Id,
				Name:      v.FileName,
				Size:      v.Size,
				CreatedAt: v.Created,
			})
	}

	return response
}
