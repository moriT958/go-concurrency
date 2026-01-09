package worker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chromedp/chromedp"
)

// URL のスクリーンショットを撮るジョブ
type ScreenshotJob struct {
	Url string
}

func ScreenshotWorker(ctx context.Context, in <-chan ScreenshotJob, out chan<- GrayScaleJob) {
	for job := range in {
		select {
		case <-ctx.Done():
			return
		default:
			img, err := processScreenshot(ctx, job.Url)
			if err != nil {
				slog.Error(fmt.Sprintf("fail to take screenshot (%s): %v", job.Url, err))
				continue
			}
			out <- GrayScaleJob{img}
		}
	}
}

func processScreenshot(ctx context.Context, siteUrl string) ([]byte, error) {
	chromedpCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var (
		screenshot []byte
		quality    = 90
	)
	if err := chromedp.Run(chromedpCtx, chromedp.Tasks{
		chromedp.Navigate(siteUrl),
		chromedp.FullScreenshot(&screenshot, quality),
	}); err != nil {
		return nil, err
	}

	return screenshot, nil
}
