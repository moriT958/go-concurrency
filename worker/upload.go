package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// 画像をストレージにアップロードするジョブ
type UploadImgJob struct {
	GrayScaleImg []byte
}

func (c *S3Client) UploadWorker(ctx context.Context, in <-chan UploadImgJob) {
	for job := range in {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.processUploadImg(ctx, job.GrayScaleImg); err != nil {
				slog.Error(fmt.Sprintf("fail to upload image: %v", err))
				continue
			}
		}
	}
}

// オブジェクトストレージにアップロードする
func (c *S3Client) processUploadImg(ctx context.Context, img []byte) error {
	objKey := fmt.Sprintf("grayscale-%s.jpg", uuid.New().String())

	// S3 にアップロード
	_, err := c.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objKey),
		Body:   bytes.NewReader(img),
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}

	// アップロード完了を待機
	if err := s3.NewObjectExistsWaiter(&c.Client).Wait(
		ctx,
		&s3.HeadObjectInput{Bucket: aws.String(c.bucketName), Key: aws.String(objKey)},
		time.Minute,
	); err != nil {
		return fmt.Errorf("failed to wait uploading finished: %v", err)
	}

	return nil
}

// AWS S3 クライアント
type S3Client struct {
	s3.Client
	bucketName string
}

// Constructor
func NewS3Client(
	s3c *s3.Client,
	bn string,
) *S3Client {
	client := new(S3Client)
	client.Client = *s3c
	client.bucketName = bn
	return client
}

// Bucket を作成 (初回時のみ実行)
func (c *S3Client) CreateBucket(ctx context.Context, name string, region string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
		LocationConstraint: types.BucketLocationConstraint(region),
	}

	if _, err := c.Client.CreateBucket(ctx, input); err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			slog.Info("You already own bucket %s.",
				"BucketName", name,
			)
			err = owned
		} else if errors.As(err, &exists) {
			slog.Info("Bucket already exists",
				"BucketName", name,
			)
			err = exists
		} else {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	if err := s3.NewBucketExistsWaiter(c).Wait(
		ctx, &s3.HeadBucketInput{Bucket: aws.String(name)},
		time.Minute,
	); err != nil {
		return fmt.Errorf("failed attempt to wait for bucket %s to exitst: %v", name, err)
	}

	return nil
}
