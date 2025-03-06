package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuidv4" gorm:"type:uuid; primary_key"`
	ResourceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255);notnull"`
	Path       string    `json:"path" valid:"notnull"  gorm:"type:varchar(255);notnull"`
	CreatedAt  time.Time `valid:"-"`
	Jobs       []*Job    `valid:"-" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

func (v *Video) Validate() error {
	_, err := govalidator.ValidateStruct(v)

	if err != nil {
		return err
	}

	return nil
}
