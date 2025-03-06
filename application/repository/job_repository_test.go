package repositories_test

import (
	repositories "encoder/application/repository"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRespositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.Path = "/videos/sample.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "queued", video)

	require.Nil(t, err)
	require.NotNil(t, job)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)

	require.Nil(t, err)
	require.NotNil(t, j)
	require.Equal(t, job.ID, j.ID)
}

func TestJobRespositoryDbUpdate(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.New().String()
	video.Path = "/videos/sample.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", "queued", video)

	require.Nil(t, err)
	require.NotNil(t, job)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "processing"
	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)

	require.Nil(t, err)
	require.NotNil(t, j)
	require.Equal(t, job.Status, j.Status)

}
