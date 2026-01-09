package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	testS3BucketName = "test-screenshots"
	testS3RegionName = "us-east-1"
)

func setupFixture() (*S3Client, error) {
	// テスト用の環境変数
	os.Setenv("AWS_ENDPOINT_URL_S3", "http://localhost:8333")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	ctx := context.Background()

	// AWS S3 SDK の初期設定
	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(testS3RegionName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws sdk default config: %v", err)
	}

	// SeaweedFS 用に設定
	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true // SeaweedFS では path-style が必要
	})
	s3c := NewS3Client(s3Client, testS3BucketName)

	// Bucket を作成
	if err := s3c.CreateBucket(ctx, testS3BucketName, testS3RegionName); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %v", err)
	}

	return s3c, nil
}

// グレースケール変換テスト
func TestProcessUploadImg(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// テストフィクスチャのセットアップ
	client, err := setupFixture()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// ダミー画像データを生成
	img := image.NewGray(image.Rect(0, 0, 100, 100))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to encode image: %v", err)
	}

	// 画像アップロード処理の実行
	if err := client.processUploadImg(ctx, buf.Bytes()); err != nil {
		t.Fatalf("processUploadImg failed: %v", err)
	}

	// バケット内のオブジェクト確認
	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(testS3BucketName),
	})
	if err != nil {
		t.Fatalf("failed to list objects: %v", err)
	}

	if result.KeyCount == nil || *result.KeyCount == 0 {
		t.Fatal("no objects found in bucket")
	}
}
