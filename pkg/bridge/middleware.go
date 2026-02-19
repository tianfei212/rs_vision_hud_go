package bridge

import (
	"fmt"
	"time"

	rs2 "github.com/tianfei212/jetson-rs-middleware"
)

// UnifiedFrame 所有模块间传递的核心对象 (SKILLs.md 4.1)
type UnifiedFrame struct {
	RawColor   []byte    // 来自中间件的原始 RGB 字节流
	RawDepth   []byte    // 来自中间件的原始 16bit 深度字节流
	Width      int       // 图像宽度
	Height     int       // 图像高度
	Timestamp  time.Time // 系统捕获时刻
	FrameIndex uint64    // 硬件帧序列号
}

// RS2Client RealSense 中间件客户端
type RS2Client struct {
	context  *rs2.Context
	pipeline *rs2.Pipeline
	config   *rs2.Config
	align    *rs2.Align
	width    int
	height   int
	fps      int
}

// NewMiddlewareClient 初始化中间件客户端 (SKILLs.md 5.1)
// 必须在此处显式开启中间件的 Align 功能
func NewMiddlewareClient(w, h, fps int) (*RS2Client, error) {
	// 1. 创建上下文
	ctx, err := rs2.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	// 2. 创建管道
	pipeline, err := rs2.NewPipeline(ctx)
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("failed to create pipeline: %w", err)
	}

	// 3. 创建配置
	config, err := rs2.NewConfig()
	if err != nil {
		pipeline.Close()
		ctx.Close()
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	// 配置流参数 (StreamType, Width, Height, FPS, Format)
	// 注意：根据 jetson-rs-middleware 的实际 API 签名调整参数顺序
	config.EnableStream(rs2.StreamColor, w, h, fps, rs2.FormatRGB8)
	config.EnableStream(rs2.StreamDepth, w, h, fps, rs2.FormatZ16)

	// 4. 初始化对齐器 (对齐到彩色流)
	align, err := rs2.NewAlign(rs2.StreamColor)
	if err != nil {
		config.Close()
		pipeline.Close()
		ctx.Close()
		return nil, fmt.Errorf("failed to create align: %w", err)
	}

	// 启动 Pipeline
	if err := pipeline.Start(config); err != nil {
		align.Close()
		config.Close()
		pipeline.Close()
		ctx.Close()
		return nil, fmt.Errorf("failed to start pipeline: %w", err)
	}

	return &RS2Client{
		context:  ctx,
		pipeline: pipeline,
		config:   config,
		align:    align,
		width:    w,
		height:   h,
		fps:      fps,
	}, nil
}

// Fetch 获取一帧对齐后的数据 (SKILLs.md 5.1)
// 内部调用线上库的 WaitForFrames()，并记录系统到达时间
func (c *RS2Client) Fetch() (*UnifiedFrame, error) {
	// 等待帧数据 (阻塞调用, 1000ms 超时)
	// 注意：WaitForFrames 接受 uint 类型的超时时间
	frames, err := c.pipeline.WaitForFrames(uint(1000))
	if err != nil {
		return nil, fmt.Errorf("wait for frames failed: %w", err)
	}
	// 确保释放原始 FrameSet
	defer frames.Close()

	// 执行对齐处理
	alignedFrames, err := c.align.Process(frames)
	if err != nil {
		return nil, fmt.Errorf("align process failed: %w", err)
	}
	// 确保释放对齐后的 FrameSet
	defer alignedFrames.Close()

	// 获取颜色帧
	colorFrame, err := alignedFrames.GetFrame(rs2.StreamColor)
	if err != nil {
		return nil, fmt.Errorf("failed to get color frame: %w", err)
	}
	// 确保释放 ColorFrame
	defer colorFrame.Close()

	// 获取深度帧
	depthFrame, err := alignedFrames.GetFrame(rs2.StreamDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to get depth frame: %w", err)
	}
	// 确保释放 DepthFrame
	defer depthFrame.Close()

	// 提取原始字节流 (⚠️ 严禁修改中间件返回的 []byte 内容)
	rawColor := colorFrame.GetRawData()
	rawDepth := depthFrame.GetRawData()

	// 获取帧序号 (暂不支持 GetFrameNumber，设为 0)
	frameIndex := uint64(0)

	// 记录系统捕获时刻
	timestamp := time.Now()

	return &UnifiedFrame{
		RawColor:   rawColor,
		RawDepth:   rawDepth,
		Width:      c.width,
		Height:     c.height,
		Timestamp:  timestamp,
		FrameIndex: frameIndex,
	}, nil
}

// Close 关闭 Pipeline 释放资源 (SKILLs.md 6.2 硬件容错)
func (c *RS2Client) Close() error {
	var errs []error

	if c.pipeline != nil {
		c.pipeline.Stop()
		c.pipeline.Close()
	}
	if c.config != nil {
		c.config.Close()
	}
	if c.align != nil {
		c.align.Close()
	}
	if c.context != nil {
		c.context.Close()
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing client: %v", errs)
	}
	return nil
}
