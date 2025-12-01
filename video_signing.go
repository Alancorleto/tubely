package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}
	urlSplit := strings.Split(*video.VideoURL, ",")
	bucket := urlSplit[0]
	key := urlSplit[1]

	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, 1*time.Hour)
	if err != nil {
		return database.Video{}, err
	}

	video.VideoURL = &presignedURL

	return video, nil
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignGetObjectInput := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	presignedHTTPRequest, err := presignClient.PresignGetObject(context.Background(), &presignGetObjectInput, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return presignedHTTPRequest.URL, nil
}
