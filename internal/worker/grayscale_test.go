package worker

import (
	"bytes"
	"image"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// グレースケール変換テスト
func TestProcessGrayScale(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// 元画像読み込み
	originalImgPath := path.Join("testdata", "effects", "out-screenshot.jpg")
	originalImg, err := os.ReadFile(originalImgPath)
	if err != nil {
		t.Fatalf("変換前画像の読み込みに失敗: %v", err)
	}

	// グレースケール変換
	grayScaleImg, err := processGrayScale(originalImg)
	if err != nil {
		t.Fatalf("グレースケール変換失敗: %v", err)
	}

	// 空でないことを確認
	if len(grayScaleImg) == 0 {
		t.Fatal("変換後画像が空です")
	}

	// デコードできることを確認
	img, _, err := image.Decode(bytes.NewReader(grayScaleImg))
	if err != nil {
		t.Fatalf("画像デコード失敗: %v", err)
	}

	// グレースケール判定
	if !isGrayscale(img, 1000) {
		t.Fatal("画像がグレースケールではありません")
	}

	// 画像をファイルに保存
	outImgPath := "testdata/effects/out-grayscale.jpg"
	if err := os.MkdirAll(filepath.Dir(outImgPath), 0o755); err != nil {
		t.Fatalf("保存失敗: %v", err)
	}
	if err := os.WriteFile(outImgPath, grayScaleImg, 0o644); err != nil {
		t.Fatalf("保存失敗: %v", err)
	}
}

// グレースケール判定処理 (誤差許容)
func isGrayscale(img image.Image, tolerance uint32) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if diff(r, g) > tolerance ||
				diff(g, b) > tolerance ||
				diff(r, b) > tolerance {
				return false
			}
		}
	}
	return true
}

func diff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}
