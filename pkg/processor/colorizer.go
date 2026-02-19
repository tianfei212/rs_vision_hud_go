package processor

import (
	"gocv.io/x/gocv"
)

// ColorizeDepth 将 16-bit 原始深度值归一化至 8-bit，并应用伪彩色映射
func ColorizeDepth(rawDepth gocv.Mat) gocv.Mat {
	// 1. 创建用于存储 8-bit 数据的容器
	normalized := gocv.NewMat()
	defer normalized.Close()

	// 2. 归一化处理：将 CV_16U (0-65535) 转换为 CV_8U (0-255)
	// 这里通常采用 alpha = 255.0 / 10000.0 (假设 10米为最大范围) 或自动缩放
	// Alpha = 255.0 / 4000.0
	alpha := 255.0 / 4000.0
	rawDepth.ConvertToWithParams(&normalized, gocv.MatTypeCV8U, float32(alpha), 0)
	rawDepth.ConvertTo(&normalized, gocv.MatTypeCV8U)

	// 3. 应用 gocv.ColorMapJet 转换为伪彩色图
	colorized := gocv.NewMat()
	gocv.ApplyColorMap(normalized, &colorized, gocv.ColormapJet)

	return colorized
}
