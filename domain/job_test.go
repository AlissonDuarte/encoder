package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	uuidv4 := uuid.New()
	uuidv4String := uuidv4.String()

	video := domain.NewVideo()
	video.ID = uuidv4String
	video.Path = "/videos/sample.mp4"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "queued", video)

	require.Nil(t, err)
	require.NotNil(t, job)
}
