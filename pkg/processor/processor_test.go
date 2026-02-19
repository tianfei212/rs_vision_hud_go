package processor

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestToMat(t *testing.T) {
	// 准备测试数据
	width := 2
	height := 2
	channels := 1
	data := []byte{0, 1, 2, 3} // 2x2 单通道图像数据

	// 调用被测函数
	mat := ToMat(data, width, height, gocv.MatTypeCV8UC1)
	defer mat.Close()

	// 验证结果
	if mat.Empty() {
		t.Error("Expected non-empty Mat")
	}
	if mat.Rows() != height {
		t.Errorf("Expected rows %d, got %d", height, mat.Rows())
	}
	if mat.Cols() != width {
		t.Errorf("Expected cols %d, got %d", width, mat.Cols())
	}
	if mat.Type() != gocv.MatTypeCV8UC1 {
		t.Errorf("Expected type %v, got %v", gocv.MatTypeCV8UC1, mat.Type())
	}
	if mat.Channels() != channels {
		t.Errorf("Expected channels %d, got %d", channels, mat.Channels())
	}

	// 验证数据内容
	// 注意：gocv.Mat 获取数据可能比较麻烦，这里简单验证基本属性即可
	// 如果需要验证数据，可以使用 mat.ToBytes()
	resultData := mat.ToBytes()
	if len(resultData) != len(data) {
		t.Errorf("Expected data length %d, got %d", len(data), len(resultData))
	}
	for i, v := range data {
		if resultData[i] != v {
			t.Errorf("Expected data at %d to be %d, got %d", i, v, resultData[i])
		}
	}
}

func TestColorizeDepth(t *testing.T) {
	// 准备测试数据：创建一个 16-bit 单通道深度图
	width := 10
	height := 10
	rawDepth := gocv.NewMatWithSize(height, width, gocv.MatTypeCV16UC1)
	defer rawDepth.Close()

	// 填充一些模拟深度数据
	// 注意：这里为了简单起见，不填充具体值，或者填充全0/全1
	// 真实场景下可以填充梯度值来验证伪彩色映射

	// 调用被测函数
	colorized := ColorizeDepth(rawDepth)
	defer colorized.Close()

	// 验证结果
	if colorized.Empty() {
		t.Error("Expected non-empty colorized Mat")
	}
	if colorized.Rows() != height {
		t.Errorf("Expected rows %d, got %d", height, colorized.Rows())
	}
	if colorized.Cols() != width {
		t.Errorf("Expected cols %d, got %d", width, colorized.Cols())
	}
	// ColorizeDepth 输出应该是 8-bit 3通道 (BGR)
	if colorized.Type() != gocv.MatTypeCV8UC3 {
		t.Errorf("Expected type %v, got %v", gocv.MatTypeCV8UC3, colorized.Type())
	}
	if colorized.Channels() != 3 {
		t.Errorf("Expected channels 3, got %d", colorized.Channels())
	}
}
