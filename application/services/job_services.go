package services

import (
	repositories "encoder/application/repository"
	"encoder/domain"
	"fmt"
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
	fmt.Println("job services - Processing")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment()
	fmt.Println("job services - Fragment")

	if err != nil {
		return j.failJob(err)
	}
	err = j.changeJobStatus("encoding")
	fmt.Println("job services - Encoding")

	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode()
	fmt.Println("job services - Encoding")

	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("uploading")
	fmt.Println("job services - Uploading")

	if err != nil {
		fmt.Printf("Error changing job status: %v", err)
		return j.failJob(err)
	}

	err = j.performUpload()
	fmt.Println("job services - Uploading")
	if err != nil {
		fmt.Println("job services - Uploading")
		return j.failJob(err)
	}

	err = j.changeJobStatus("finishing")
	fmt.Println("job services - Finishing")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finalize()
	fmt.Println("job services - Finishing")
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("finished")
	fmt.Println("job services - Finished")
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
	fmt.Println("criando new video upload")
	videoUpload.OutPutBucket = os.Getenv("outputBucketName")
	fmt.Printf("OutPutBucket setado como : %s\n", videoUpload.OutPutBucket)
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.Job.Video.ID
	fmt.Printf("VideoPath setado como : %s\n", videoUpload.VideoPath)
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY"))
	fmt.Printf("concurrency setado como : %d\n", concurrency)
	doneUpload := make(chan string)

	fmt.Println("criando go rotine")
	go videoUpload.ProcessUpload(concurrency, doneUpload)
	fmt.Println("go rotine criada")

	fmt.Println("esperando resultado")
	uploadResult := <-doneUpload
	fmt.Printf("Resultado: %s\n", uploadResult)
	if uploadResult != "uploaded completed" {
		return j.failJob(err)
	}
	fmt.Println("finalizando processo de upload")
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
