package models

import (
	"time"
)

// UnifiedFrame 作为跨模块传递的原子数据单元，确保彩色流与深度流的绝对对齐。
type UnifiedFrame struct {
	RawColor   []byte    // 来自中间件的原始 RGB 字节流 (Format: BGR/RGB)
	RawDepth   []byte    // 来自中间件的原始 16-bit 深度字节流
	Width      int       // 图像宽度
	Height     int       // 图像高度
	Timestamp  time.Time // 数据捕获时的系统精确时间戳
	FrameIndex uint64    // RealSense 硬件帧序列号
}
