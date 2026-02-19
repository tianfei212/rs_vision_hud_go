package hud

import (
	"fmt"
	"image"
	"image/color"

	"rs-vision-hub-go/pkg/models"

	"gocv.io/x/gocv"
)

const (
	// 定义透明度 (Alpha = 120, 占比约为 0.47)
	alphaWeight = 120.0 / 255.0
	betaWeight  = 1.0 - alphaWeight

	// 字体配置
	fontFace  = gocv.FontHersheySimplex
	fontScale = 0.6
	thickness = 2
)

// 预定义颜色
var (
	textColor = color.RGBA{R: 255, G: 255, B: 255, A: 0} // 白色字体
	rectColor = color.RGBA{R: 0, G: 0, B: 0, A: 0}       // 黑色底衬
)

// OverlayHUD 在图像上实时绘制 HUD 信息
// 注意：为了满足左下角显示实时 FPS 的需求，函数签名增加了 currentFPS 参数
func OverlayHUD(img *gocv.Mat, batch *models.UnifiedFrame, currentFPS float64) {
	if img == nil || img.Empty() {
		return
	}

	// 1. 创建半透明底衬的 Overlay 图层
	// 使用 Clone() 复制一帧用于 Alpha 混合，避免直接修改底层像素引发内存竞争
	overlay := img.Clone()
	defer overlay.Close() // 必须显式释放 C++ 内存

	// 2. 准备绘制内容
	// 右上角：时间戳 (严格按照规范要求格式化)
	timeStr := batch.Timestamp.Format("2006-01-02 15:04:05.000")

	// 左下角：FPS, 分辨率, 硬件帧序号
	infoStr := fmt.Sprintf("FPS: %.1f | Res: %dx%d | Frame: %d",
		currentFPS, batch.Width, batch.Height, batch.FrameIndex)

	// 3. 计算文本尺寸与背景矩形位置
	// 时间戳矩形 (右上角)
	timeSize := gocv.GetTextSize(timeStr, fontFace, fontScale, thickness)
	timeRect := image.Rect(
		img.Cols()-timeSize.X-20, 10,
		img.Cols()-10, 10+timeSize.Y+15,
	)

	// 信息矩形 (左下角)
	infoSize := gocv.GetTextSize(infoStr, fontFace, fontScale, thickness)
	infoRect := image.Rect(
		10, img.Rows()-infoSize.Y-20,
		10+infoSize.X+10, img.Rows()-10,
	)

	// 4. 在 Overlay 图层上绘制黑色矩形
	gocv.Rectangle(&overlay, timeRect, rectColor, -1)
	gocv.Rectangle(&overlay, infoRect, rectColor, -1)

	// 5. 执行 Alpha 混合融合底图与 Overlay
	gocv.AddWeighted(overlay, alphaWeight, *img, betaWeight, 0.0, img)

	// 6. 在融合后的主图上绘制清晰文本 (确保文本本身不受 Alpha 影响)
	gocv.PutText(img, timeStr, image.Pt(timeRect.Min.X+5, timeRect.Max.Y-8), fontFace, fontScale, textColor, thickness)
	gocv.PutText(img, infoStr, image.Pt(infoRect.Min.X+5, infoRect.Max.Y-8), fontFace, fontScale, textColor, thickness)
}

// DrawCenterDistance 在图像中心绘制准星和距离信息
func DrawCenterDistance(img *gocv.Mat, batch *models.UnifiedFrame, distance float64) {
	if img == nil || img.Empty() {
		return
	}

	cx, cy := batch.Width/2, batch.Height/2

	// 绘制中心十字准星 (绿色)
	crosshairColor := color.RGBA{R: 0, G: 255, B: 0, A: 0}
	length := 20
	// 横线
	gocv.Line(img, image.Pt(cx-length, cy), image.Pt(cx+length, cy), crosshairColor, 2)
	// 竖线
	gocv.Line(img, image.Pt(cx, cy-length), image.Pt(cx, cy+length), crosshairColor, 2)

	// 绘制距离文字
	distStr := fmt.Sprintf("Dist: %.2fm", distance)
	// 计算文字大小以居中
	textSize := gocv.GetTextSize(distStr, fontFace, fontScale, thickness)
	textOrigin := image.Pt(cx-textSize.X/2, cy+length+textSize.Y+5)

	// 绘制文字 (带简单阴影以增强可见度)
	gocv.PutText(img, distStr, textOrigin, fontFace, fontScale, color.RGBA{0, 0, 0, 0}, thickness+2) // 黑色描边
	gocv.PutText(img, distStr, textOrigin, fontFace, fontScale, crosshairColor, thickness)           // 绿色主体
}
