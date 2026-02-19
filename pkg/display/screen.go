package display

import (
	"fmt"

	"gocv.io/x/gocv"
)

// Screen 封装 GoCV 窗口同步显示逻辑
type Screen struct {
	colorWin *gocv.Window
	depthWin *gocv.Window
}

// NewScreen 初始化双窗口渲染引擎
func NewScreen() *Screen {
	// 创建独立窗口，如果不希望用户调整大小，可以使用 gocv.WindowAutoSize
	colorWin := gocv.NewWindow("RealSense RGB Stream - Vision Hub")
	depthWin := gocv.NewWindow("RealSense Depth Stream - Vision Hub")

	return &Screen{
		colorWin: colorWin,
		depthWin: depthWin,
	}
}

// Render 在两个独立窗口同步刷新 RGB 和深度图像 (SKILLs.md 4.4)
// 注意：该方法必须在 main 主协程中调用，否则在 Jetson/Linux 环境下可能引发 X11 线程错误
func (s *Screen) Render(color, depth gocv.Mat) {
	if color.Empty() || depth.Empty() {
		return
	}

	// 同步将图像帧送入显示缓冲区
	s.colorWin.IMShow(color)
	s.depthWin.IMShow(depth)

	// WaitKey(1) 是必须的：它不仅用于捕获键盘事件，更是驱动 OpenCV GUI 消息循环刷新画面的核心机制
	// 1ms 延迟对实时 FPS 的影响微乎其微
	s.colorWin.WaitKey(1)
}

// Close 安全释放底层 X11/Wayland 窗口句柄资源
func (s *Screen) Close() error {
	var errs []error

	if s.colorWin != nil {
		if err := s.colorWin.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close color window: %w", err))
		}
	}

	if s.depthWin != nil {
		if err := s.depthWin.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close depth window: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("screen close errors: %v", errs)
	}
	return nil
}
