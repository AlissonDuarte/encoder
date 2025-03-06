package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuidv4" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"in(queued|processing|finished|failed)"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `valid:"-" gorm:"column:video_id; type:uuid;notnull"`
	Error            string    `json:"error" valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (j *Job) prepare() {
	u := uuid.New()
	uString := u.String()

	j.ID = uString
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	j := Job{
		OutputBucketPath: output,
		Status:           status,
		Video:            video,
		VideoID:          video.ID,
	}

	j.prepare()

	if err := j.Validate(); err != nil {
		return nil, err
	}

	return &j, nil
}

func (j *Job) Validate() error {
	_, err := govalidator.ValidateStruct(j)

	if err != nil {
		return err
	}

	return nil
}
