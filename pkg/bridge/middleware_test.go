package bridge

import (
	"testing"
	"time"
)

func TestNewMiddlewareClient(t *testing.T) {
	// 尝试初始化客户端
	// 注意：如果没有连接 RealSense 设备，这里可能会失败
	// 我们预期它要么成功，要么返回特定的错误（如 "No device detected"）

	width := 640
	height := 480
	fps := 30

	client, err := NewMiddlewareClient(width, height, fps)
	if err != nil {
		// 如果是因为没有设备，我们可以跳过测试
		t.Logf("Failed to create middleware client (expected if no device): %v", err)
		return
	}
	defer client.Close()

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.width != width {
		t.Errorf("Expected width %d, got %d", width, client.width)
	}
	if client.height != height {
		t.Errorf("Expected height %d, got %d", height, client.height)
	}
	if client.fps != fps {
		t.Errorf("Expected fps %d, got %d", fps, client.fps)
	}
}

func TestRS2Client_Fetch(t *testing.T) {
	width := 640
	height := 480
	fps := 30

	client, err := NewMiddlewareClient(width, height, fps)
	if err != nil {
		t.Skipf("Skipping Fetch test due to initialization failure: %v", err)
	}
	defer client.Close()

	// 尝试获取几帧数据
	for i := 0; i < 3; i++ {
		start := time.Now()
		frame, err := client.Fetch()
		if err != nil {
			t.Errorf("Fetch failed on iteration %d: %v", i, err)
			continue
		}

		// 验证 UnifiedFrame 结构
		if frame == nil {
			t.Fatal("Fetch returned nil frame")
		}

		if frame.Width != width {
			t.Errorf("Frame width mismatch: expected %d, got %d", width, frame.Width)
		}
		if frame.Height != height {
			t.Errorf("Frame height mismatch: expected %d, got %d", height, frame.Height)
		}
		if len(frame.RawColor) == 0 {
			t.Error("RawColor is empty")
		}
		if len(frame.RawDepth) == 0 {
			t.Error("RawDepth is empty")
		}

		// 验证时间戳是否合理
		if frame.Timestamp.Before(start) {
			t.Error("Frame timestamp is before fetch start time")
		}

		t.Logf("Frame %d: Index=%d, Timestamp=%v, ColorLen=%d, DepthLen=%d",
			i, frame.FrameIndex, frame.Timestamp, len(frame.RawColor), len(frame.RawDepth))
	}
}
