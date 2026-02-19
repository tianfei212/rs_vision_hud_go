package processor

import (
	"gocv.io/x/gocv"
)

// ToMat 实现从 byte slice 到 gocv.Mat 的转换
// 严格遵循规格书要求的函数签名与 Zero-copy 原则
func ToMat(data []byte, w, h int, t gocv.MatType) gocv.Mat {
	// 必须使用 gocv.NewMatFromBytes 确保零拷贝
	// 此处省略错误处理以匹配规格书定义的返回值类型
	mat, _ := gocv.NewMatFromBytes(h, w, t, data)
	return mat
}
