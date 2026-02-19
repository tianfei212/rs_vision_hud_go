package processor

// GetCenterDistance 提取画面中心点的深度距离（单位：米）
func GetCenterDistance(rawDepth []byte, width, height int) float64 {
	if len(rawDepth) == 0 {
		return 0.0
	}

	cx, cy := width/2, height/2
	// 深度图是 16-bit (CV16UC1)，每个像素占 2 字节
	pixelIndex := (cy*width + cx) * 2

	// 防止数组越界越界
	if pixelIndex+1 >= len(rawDepth) {
		return 0.0
	}

	// RealSense 采用小端序 (Little-Endian)，合并两个 byte 为 uint16 (毫米)
	rawZ := uint16(rawDepth[pixelIndex]) | (uint16(rawDepth[pixelIndex+1]) << 8)

	// 转换为米
	return float64(rawZ) / 1000.0
}
