package worker

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// 結合テスト: スクリーンショット
func TestProcessScreenshot(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	siteUrl := "https://morit958.com"
	outImgPath := "testdata/effects/out-screenshot.jpg"

	// スクリーンショット実行
	ctx := context.Background()
	screenshot, err := processScreenshot(ctx, siteUrl)
	if err != nil {
		t.Fatalf("スクリーンショット失敗: %v", err)
	}

	// 画像をファイルに保存
	if err := os.MkdirAll(filepath.Dir(outImgPath), 0o755); err != nil {
		t.Fatalf("保存失敗: %v", err)
	}
	if err := os.WriteFile(outImgPath, screenshot, 0o644); err != nil {
		t.Fatalf("保存失敗: %v", err)
	}
}
