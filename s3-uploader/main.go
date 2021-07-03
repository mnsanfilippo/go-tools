package s3_uploader

import (
	"bytes"
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

var BUCKET = flag.String("b", os.Getenv("BUCKET"), "S3 Bucket Name")
var KEY = flag.String("k", os.Getenv("KEY"), "S3 Key")
var FILE = flag.String("f", os.Getenv("FILE"), "File to be uploaded")

func Upload() error {

	file, err := os.Open(*FILE)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
		return err
	}

	uploader := manager.NewUploader(s3.NewFromConfig(cfg), func(u *manager.Uploader) {
		// Define a strategy that will buffer 500 MiB in memory
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(500 * 1024 * 1024)
	})

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(*BUCKET),
		Key:    aws.String(*KEY),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}
