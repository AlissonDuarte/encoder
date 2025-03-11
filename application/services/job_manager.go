package services

import (
	repositories "encoder/application/repository"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
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

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRRENCY"))

	if err != nil {
		concurrency = 1
		log.Fatalf("Error parsing CONCURRRENCY env var: %v", err)
	}

	for i := 0; i < concurrency; i++ {
		go JobWorker(j.MessageChannel, j.JobReturn, jobService, i, j.Domain)
	}

	for jobResult := range j.JobReturn {
		if jobResult.Error != nil {
			log.Printf("Error processing job: %v", jobResult.Error)
			j.notifyError(jobResult)

		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			log.Printf("Error processing job: %v", err)
			jobResult.Message.Reject(false)
		}
	}
}

func (j *JobManager) notifyError(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Fatalf("Error processing job: %v and tag %v", jobResult.Error, jobResult.Message.DeliveryTag)
	} else {
		log.Fatalf("Error processing job: %v", jobResult.Error)
	}

	jobNotificationError := JobNotificationError{
		Message: "Error processing job",
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(jobNotificationError)

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
		os.Getenv("RABBIT_NOTIFICATION_EXCHANGE"),
		os.Getenv("RABBIT_NOTIFICATION_ROUTING_KEY"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {

	jobJson, err := json.Marshal(jobResult.Job)

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
