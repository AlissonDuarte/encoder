package services

import (
	repositories "encoder/application/repository"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db             *gorm.DB
	Domain         domain.Job
	MessageChannel chan amqp.Delivery
	JobReturn      chan JobWorkerResult
	Rabbit         *queue.Rabbit
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, messageChannel chan amqp.Delivery, jobReturn chan JobWorkerResult, rabbit *queue.Rabbit) *JobManager {
	return &JobManager{
		Db:             db,
		Domain:         domain.Job{},
		MessageChannel: messageChannel,
		JobReturn:      jobReturn,
		Rabbit:         rabbit,
	}

}

func (j *JobManager) Start(ch *amqp.Channel) {

	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))

	if err != nil {
		concurrency = 1
		log.Fatalf("Error parsing CONCURRENCY env var: %v", err)
	}

	for i := 0; i < concurrency; i++ {
		go JobWorker(j.MessageChannel, j.JobReturn, jobService, i, j.Domain)
	}

	for jobResult := range j.JobReturn {
		fmt.Println("Job result: ", jobResult)
		fmt.Println("Job result: ", jobResult.Error)
		if jobResult.Error != nil {
			log.Printf("Error processing job1: %v", jobResult.Error)
			j.notifyError(jobResult)

		} else {
			err = j.notifySuccess(jobResult)
		}

		if err != nil {
			log.Printf("Error processing job2: %v", err)
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) notifyError(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID, %v, Error During the job %v with video %v, error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.VideoID, jobResult.Error.Error())
	} else {
		log.Printf("Error During the job %v, error: %v", jobResult.Job.ID, jobResult.Error.Error())
	}

	jobNotificationError := JobNotificationError{
		Message: "Error processing job5",
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(jobNotificationError)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.Rabbit.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBIT_NOTIFICATION_EX"),
		os.Getenv("RABBIT_NOTIFICATION_ROUTING_KEY"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult) error {

	Mutex.Lock()
	jobJson, err := json.Marshal(jobResult.Job)
	Mutex.Unlock()

	if err != nil {
		return err
	}

	err = j.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}
