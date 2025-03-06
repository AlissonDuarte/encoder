package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {

	tests := []struct {
		name    string
		video   domain.Video
		isValid bool
	}{
		{
			name: "Valid Video",
			video: domain.Video{
				ID:         "54b5c6a7-8c9d-4e3f-a4b5-c6a7d8e9f0a1",
				ResourceID: "123",
				Path:       "/videos/sample.mp4",
				CreatedAt:  time.Now(),
			},
			isValid: true,
		},
		{
			name: "Invalid UUID",
			video: domain.Video{
				ID:         "invalid-uuid",
				ResourceID: "123",
				Path:       "/videos/sample.mp4",
				CreatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "Empty ResourceID",
			video: domain.Video{
				ID:         "54b5c6a7-8c9d-4e3f-a4b5-c6a7d8e9f0a1",
				ResourceID: "",
				Path:       "/videos/sample.mp4",
				CreatedAt:  time.Now(),
			},
			isValid: false,
		},
		{
			name: "Empty Path",
			video: domain.Video{
				ID:         "54b5c6a7-8c9d-4e3f-a4b5-c6a7d8e9f0a1",
				ResourceID: "123",
				Path:       "",
				CreatedAt:  time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := govalidator.ValidateStruct(tt.video)
			if tt.isValid {
				assert.NoError(t, err, "Expected valid struct, but got an error")
			} else {
				assert.Error(t, err, "Expected an error, but struct is valid")
			}
		})
	}
}
