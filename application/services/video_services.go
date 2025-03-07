package services

import (
	"context"
	repositories "encoder/application/repository"
	"encoder/domain"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video            *domain.Video
	VideoReposeitory repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	envPath := os.Getenv("localStoragePath")
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	fmt.Printf("Downloading %s to %s\n", v.Video.Path, envPath)
	obj := bkt.Object(v.Video.Path)
	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	f, err := os.Create(envPath + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	created, err := f.Write(body)

	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("Downloaded %d bytes to %s.", created, v.Video.ID+".mp4")

	return nil
}
