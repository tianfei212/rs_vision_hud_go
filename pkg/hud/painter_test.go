package hud

import (
	"image"
	"image/color"
	"testing"
	"time"

	"rs-vision-hub-go/pkg/models"

	"gocv.io/x/gocv"
)

func TestOverlayHUD(t *testing.T) {
	// 1. Prepare a blank image
	width := 640
	height := 480
	// Create a black image (CV_8UC3)
	img := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC3)
	defer img.Close()

	// Fill with black to ensure we have a consistent starting point
	// Although NewMatWithSize allocates, it might contain garbage or zeros depending on implementation,
	// usually zeros but setting explicitly is safer for test logic.
	// Actually NewMatWithSize is usually zeroed in GoCV/OpenCV context? No, it's uninitialized usually.
	// So let's zero it out.
	black := color.RGBA{0, 0, 0, 0}
	gocv.Rectangle(&img, image.Rect(0, 0, width, height), black, -1)

	// Keep a copy of the original (black) image for comparison
	original := img.Clone()
	defer original.Close()

	// 2. Prepare test data
	now := time.Now()
	frame := &models.UnifiedFrame{
		Width:      width,
		Height:     height,
		Timestamp:  now,
		FrameIndex: 12345,
	}
	currentFPS := 30.5

	// 3. Call the function
	OverlayHUD(&img, frame, currentFPS)

	// 4. Basic validation
	if img.Empty() {
		t.Error("Result image should not be empty")
	}
	if img.Rows() != height {
		t.Errorf("Image height changed. Got %d, want %d", img.Rows(), height)
	}
	if img.Cols() != width {
		t.Errorf("Image width changed. Got %d, want %d", img.Cols(), width)
	}

	// 5. Check if image content changed (HUD drawn)
	// Compute absolute difference between original (black) and result
	diff := gocv.NewMat()
	defer diff.Close()
	gocv.AbsDiff(original, img, &diff)

	// Convert to grayscale to count non-zero pixels easily
	grayDiff := gocv.NewMat()
	defer grayDiff.Close()
	gocv.CvtColor(diff, &grayDiff, gocv.ColorBGRToGray)

	nonZero := gocv.CountNonZero(grayDiff)
	if nonZero == 0 {
		t.Error("Image content did not change; HUD was not drawn")
	} else {
		t.Logf("HUD drawn successfully, changed %d pixels", nonZero)
	}
}

func TestOverlayHUD_NilImage(t *testing.T) {
	// Should not panic when passed nil
	frame := &models.UnifiedFrame{
		Timestamp: time.Now(),
	}
	OverlayHUD(nil, frame, 0)
}

func TestOverlayHUD_EmptyImage(t *testing.T) {
	// Should not panic when passed empty Mat
	img := gocv.NewMat()
	defer img.Close()
	
	frame := &models.UnifiedFrame{
		Timestamp: time.Now(),
	}
	
	OverlayHUD(&img, frame, 0)
	
	if !img.Empty() {
		t.Error("Empty image should remain empty")
	}
}
