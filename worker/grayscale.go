package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
)

// スクリーンショットをグレスケール変換するジョブ
type GrayScaleJob struct {
	Img []byte
}

func GrayScaleWorker(ctx context.Context, in <-chan GrayScaleJob, out chan<- UploadImgJob) {
	for job := range in {
		select {
		case <-ctx.Done():
			return
		default:
			grayscaleImg, err := processGrayScale(job.Img)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to process gray scale: %v", err))
				continue
			}
			out <- UploadImgJob{grayscaleImg}
		}
	}
}

func processGrayScale(originalImg []byte) ([]byte, error) {
	// 元画像をデコード
	originalImgDecoded, _, err := image.Decode(bytes.NewReader(originalImg))
	if err != nil {
		return nil, fmt.Errorf("failed to decode original image: %v", err)
	}

	// グレースケール変換
	bounds := originalImgDecoded.Bounds()
	gray := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			gray.Set(x, y, originalImgDecoded.At(x, y))
		}
	}

	// jpeg 形式で書き込む
	var grayscaleBuf bytes.Buffer
	if err := jpeg.Encode(&grayscaleBuf, gray, nil); err != nil {
		return nil, fmt.Errorf("failed to write jpeg image: %v", err)
	}

	return grayscaleBuf.Bytes(), nil
}
