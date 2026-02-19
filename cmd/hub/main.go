package main

import (
	"log"
	"time"

	"rs-vision-hub-go/pkg/bridge"
	"rs-vision-hub-go/pkg/display"
	"rs-vision-hub-go/pkg/hud"
	"rs-vision-hub-go/pkg/processor"

	"gocv.io/x/gocv"
)

func main() {
	// 1. 配置参数 (可根据你的 Realsense 型号支持的分辨率进行调整)
	const width, height, fps = 640, 480, 30
	log.Println("Initializing RealSense Middleware Client...")
	client, err := bridge.NewMiddlewareClient(width, height, fps)
	if err != nil {
		log.Fatalf("Failed to initialize middleware client: %v", err)
	}
	defer client.Close()

	log.Println("Middleware client initialized successfully. Initialize display screen...")
	screen := display.NewScreen()
	defer screen.Close()

	log.Println("Starting Main Render Loop...")
	var frameCount uint64
	var currentFPS float64
	lastTime := time.Now()

	// 核心生命周期与渲染主循环 (SKILLs.md 5.1)
	for {
		// [Step 1] 获取同步帧 (Fetch)
		batch, err := client.Fetch()
		if err != nil {
			log.Printf("Warning: Failed to fetch frame: %v", err)
			continue
		}

		// [Step 2] 数据转换 (Processor)
		// 将 []byte 转换为 GoCV Mat，注意指明正确的矩阵类型 (3通道)
		rawRGBMat := processor.ToMat(batch.RawColor, batch.Width, batch.Height, gocv.MatTypeCV8UC3)

		// 修复色彩空间错位：FormatRGB8 是 3 通道，必须使用 ColorRGBToBGR
		// 创建新 Mat 存储转换结果，避免修改原始只读内存
		colorMat := gocv.NewMat()
		gocv.CvtColor(rawRGBMat, &colorMat, gocv.ColorRGBToBGR)
		rawRGBMat.Close() // 及时释放中间变量

		rawDepthMat := processor.ToMat(batch.RawDepth, batch.Width, batch.Height, gocv.MatTypeCV16UC1)

		// 将 16-bit 深度图归一化并映射为伪彩色图
		colorizedDepthMat := processor.ColorizeDepth(rawDepthMat)

		// [Step 3] 数据提取与性能计算
		// 提取画面中心点的物理距离 (单位: 米)
		centerDistance := processor.GetCenterDistance(batch.RawDepth, batch.Width, batch.Height)

		frameCount++
		if frameCount%10 == 0 { // 每 10 帧更新一次 FPS
			now := time.Now()
			currentFPS = 10.0 / now.Sub(lastTime).Seconds()
			lastTime = now
		}
		batch.FrameIndex = frameCount // 填充硬件帧号占位

		// [Step 4] 视觉层叠加 (HUD)
		// 绘制基础信息 (时间戳、FPS等)
		hud.OverlayHUD(&colorMat, batch, currentFPS)
		hud.OverlayHUD(&colorizedDepthMat, batch, currentFPS)

		// 绘制中心测距准星 (可以同时画在彩色图和深度图上，这里我们两边都画)
		hud.DrawCenterDistance(&colorMat, batch, centerDistance)
		hud.DrawCenterDistance(&colorizedDepthMat, batch, centerDistance)

		// [Step 5] 渲染层同步显示 (Display)
		// 必须在主协程调用
		screen.Render(colorMat, colorizedDepthMat)

		// [Step 6] 清理 (Clean)
		// 极度重要：每次循环必须显式调用 Close 释放 CGO 侧分配的图像内存 (SKILLs.md 5.5)
		colorMat.Close()
		rawDepthMat.Close()
		colorizedDepthMat.Close()
	}
}
