package main

import (
	"concurrency/internal/worker"
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

const (
	s3BucketName = "screenshots"
	s3RegionName = "us-east-1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// AWS S3 SDK の初期設定
	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s3RegionName),
	)
	if err != nil {
		log.Fatalf("failed to load aws sdk default config: %v", err)
		return
	}
	// SeaweedFS 用に設定
	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true // SeaweedFS では path-style が必要
	})
	s3c := worker.NewS3Client(s3Client, s3BucketName)

	// Bucket を作成
	if err := s3c.CreateBucket(ctx, s3BucketName, s3RegionName); err != nil {
		log.Fatalf("failed to initialized screenshots bucket: %v", err)
	}

	// process 開始
	screenshotCh := make(chan worker.ScreenshotJob, 1000)
	grayscaleCh := make(chan worker.GrayScaleJob, 1000)
	uploadCh := make(chan worker.UploadImgJob, 1000)

	var wg sync.WaitGroup

	// NOTE:
	// Go 1.25 から追加された WaitGroup.Go
	// wg.Add, wg.Done を省略できる

	// Screenshot Worker
	wg.Go(func() {
		worker.ScreenshotWorker(ctx, screenshotCh, grayscaleCh)
		close(grayscaleCh) // 処理完了後に次のチャネルをクローズ
	})

	// GrayScale Worker
	wg.Go(func() {
		worker.GrayScaleWorker(ctx, grayscaleCh, uploadCh)
		close(uploadCh) // 処理完了後に次のチャネルをクローズ
	})

	// Upload Worker
	wg.Go(func() {
		s3c.UploadWorker(ctx, uploadCh)
	})

	// 10個のジョブを投入
	for range 10 {
		screenshotCh <- worker.ScreenshotJob{Url: "https://morit958.com"}
	}
	close(screenshotCh) // 最初のチャネルだけクローズ

	// すべての処理が完了するまで待機
	wg.Wait()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
