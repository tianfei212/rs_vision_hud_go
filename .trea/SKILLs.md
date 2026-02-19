# **🛠 SKILLS.md: rs-vision-hub-go 开发规格说明书**

## **1\. 项目定义 (Project Definition)**

**rs-vision-hub-go** 是一个基于 NVIDIA Jetson 平台的实时视觉测试终端。它通过集成线上仓库 jetson-rs-middleware，实现对 Intel RealSense 相机的双路（RGB & Depth）同步捕获，并通过 Go 语言在本地实现 HUD 信息叠加与实时双窗口渲染。

## **2\. 目录结构与模块职责 (Project Layout)**

rs-vision-hub-go/  
├── cmd/  
│   └── hub/  
│       └── main.go          \# 应用入口：管控程序生命周期与渲染主循环  
├── pkg/  
│   ├── bridge/  
│   │   └── middleware.go    \# 适配层：封装对线上中间件接口的调用  
│   ├── processor/  
│   │   ├── converter.go     \# 转换层：byte slice \-\> gocv.Mat (实现 Zero-copy)  
│   │   └── colorizer.go     \# 处理层：16-bit 深度数据归一化与伪彩色映射  
│   ├── hud/  
│   │   └── painter.go       \# 视觉层：在图像上实时绘制 HUD 信息（时间戳、FPS、元数据）  
│   └── display/  
│       └── screen.go        \# 渲染层：封装 GoCV 窗口同步显示逻辑  
├── go.mod                   \# 依赖管理：锁定 \[github.com/tianfei212/jetson-rs-middleware\](https://github.com/tianfei212/jetson-rs-middleware)  
└── Makefile                 \# 构建脚本：自动化处理特定版本 .so 的 CGO 链接与 RPATH 嵌入

## **3\. 核心数据模型 (Data Models)**

### **3.1 UnifiedFrame**

作为跨模块传递的原子数据单元，确保彩色流与深度流的绝对对齐。

type UnifiedFrame struct {  
    RawColor   \[\]byte            // 来自中间件的原始 RGB 字节流 (Format: BGR/RGB)  
    RawDepth   \[\]byte            // 来自中间件的原始 16-bit 深度字节流  
    Width      int                 
    Height     int                 
    Timestamp  time.Time         // 数据捕获时的系统精确时间戳  
    FrameIndex uint64            // RealSense 硬件帧序列号  
}

## **4\. 函数定义与接口规范 (Function Specifications)**

### **📂 pkg/bridge/middleware.go**

* **func NewMiddlewareClient(w, h, fps int) (\*RS2Client, error)**  
  * **职责**: 初始化 Pipeline。  
  * **指导**: 必须调用中间件接口显式开启 Align(RS2\_STREAM\_COLOR)，确保深度点云与彩色像素在空间坐标上完全重合。  
* **func (c \*RS2Client) Fetch() (\*UnifiedFrame, error)**  
  * **职责**: 封装中间件的 WaitForFrames()，将底层数据包转换为 UnifiedFrame。

### **📂 pkg/processor/**

* **func ToMat(data \[\]byte, w, h int, t gocv.MatType) gocv.Mat**  
  * **指导**: **必须**使用 gocv.NewMatFromBytes。严禁在处理循环中产生不必要的内存拷贝（Zero-copy 原则）。  
* **func ColorizeDepth(rawDepth gocv.Mat) gocv.Mat**  
  * **职责**: 将 16-bit 原始深度值归一化至 8-bit，并应用 gocv.ColorMapJet 转换为直观的伪彩色图。

### **📂 pkg/hud/painter.go**

* **func OverlayHUD(img \*gocv.Mat, batch \*UnifiedFrame)**  
  * **职责**:  
    1. 在右上角绘制 2006-01-02 15:04:05.000 格式时间戳。  
    2. 在左下角叠加实时 FPS、分辨率及硬件帧序号。  
  * **要求**: 所有文字必须带有半透明（Alpha=120）黑色矩形底衬，确保在复杂光影环境下 HUD 信息清晰可见。

### **📂 pkg/display/screen.go**

* **func (s \*Screen) Render(color, depth gocv.Mat)**  
  * **职责**: 在 RGB Stream 和 Depth Stream 两个独立窗口同步刷新图像。

## **5\. 交互逻辑要求 (Interaction Requirements)**

### **🔄 核心执行流**

1. **Init**: main.go 初始化 middleware.RS2Client。  
2. **Loop**: 开启 for 循环，调用 client.Fetch() 获取同步帧。  
3. **Async Process (Optional)**: 可以在 Goroutine 中完成数据转换与 HUD 叠加以提升吞吐。  
4. **Sync Display**: 必须在**主线程**调用 screen.Render。  
5. **Clean**: 每次循环结束前，必须显式调用所有 gocv.Mat 的 Close() 方法，防止 Jetson 显存溢出。

## **6\. 特定驱动库链接规范 (CGO & Shared Library Rules)**

### **⚠️ 重要依赖说明**

由于 jetson-rs-middleware 随包分发了特定版本的 librealsense2.so（位于其 /lib 目录），本项目**严禁**链接系统全局库。

1. **编译时链接**: 必须通过 go list 动态定位中间件在 $GOPATH/pkg/mod 的缓存位置。  
2. **运行时寻址 (RPATH)**: 构建时必须使用 \-ldflags "-r \[path\]" 参数，将中间件 lib 目录的绝对路径嵌入到生成的二进制文件中。  
3. **版本隔离**: 确保运行时加载的是中间件自带的 .so，以保证 CGO 调用时的结构体对齐完全一致。

## **7\. 构建与运行指南 (Build & Run)**

\# 1\. 下载/更新线上依赖  
go mod download \[github.com/tianfei212/jetson-rs-middleware\](https://github.com/tianfei212/jetson-rs-middleware)

\# 2\. 自动化定位中间件缓存路径  
MDW\_DIR=$(go list \-m \-f '{{.Dir}}' \[github.com/tianfei212/jetson-rs-middleware\](https://github.com/tianfei212/jetson-rs-middleware))

\# 3\. 执行链接并嵌入运行时路径 (RPATH)  
go build \-ldflags="-r $MDW\_DIR/lib" \-o rs-vision-hub ./cmd/hub/main.go

\# 4\. 运行应用  
./rs-vision-hub  
