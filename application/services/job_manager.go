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

	if err != nil {
		return err
	}

	err = j.Rabbit.Publish(
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jobJson,
		},
		"job_error",
	)

	if err != nil {
		return err
	}

	return nil
}
