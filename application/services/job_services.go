package services

import (
	repositories "encoder/application/repository"
	"encoder/domain"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	err := j.changeJobStatus("queued")

	if err != nil {
		return err
	}

	err = j.VideoService.Download(os.Getenv("inputBucketName"))

	if err != nil {
		return err
	}

	err = j.changeJobStatus("processing")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()

	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus("encoding")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("uploading")

	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("finishing")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finalize()
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("finished")
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.changeJobStatus("uploading")

	if err != nil {
		return err
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutPutBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.Job.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload
	if uploadResult != "uploaded completed" {
		return j.failJob(err)
	}
	return err
}
func (j *JobService) changeJobStatus(status string) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {

	j.Job.Status = "failed"
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return error
}
