# RS Vision Hub Go

`rs-vision-hub-go` is a high-performance, real-time vision monitoring application designed for the NVIDIA Jetson platform. It serves as the official reference implementation for the `jetson-rs-middleware`, demonstrating how to efficiently acquire aligned RealSense camera streams and render them with a Heads-Up Display (HUD) using pure Go.

## 🌟 Features

*   **Real-time Stream Alignment**: Seamlessly captures and aligns RGB and Depth streams from RealSense cameras.
*   **Zero-Copy Processing**: Optimized data handling using `GoCV` to minimize memory overhead on embedded devices.
*   **Heads-Up Display (HUD)**:
    *   Real-time FPS monitoring.
    *   Timestamp overlay.
    *   Resolution and frame index tracking.
    *   **Center Distance Measurement**: Live depth distance reading at the center of the frame with a visual crosshair.
*   **Depth Visualization**: Converts 16-bit raw depth data into intuitive pseudo-color heatmaps (Jet colormap).
*   **Dual-Window Rendering**: Synchronized display of RGB and Depth streams.

## 🛠 Prerequisites

*   **NVIDIA Jetson Device** (Or any Linux machine with RealSense SDK installed)
*   **Go 1.21+**
*   **OpenCV 4.x**
*   **Intel RealSense SDK 2.0** (`librealsense2`)

## 📦 Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/your-username/rs-vision-hub-go.git
    cd rs-vision-hub-go
    ```

2.  **Download dependencies**
    ```bash
    go mod download
    ```

3.  **Build the application**
    ```bash
    go build -o rs-vision-hub ./cmd/hub/main.go
    ```

## 🚀 Usage

Run the compiled binary:

```bash
./rs-vision-hub
```

**Note**: Ensure your RealSense camera is connected via USB 3.0.

## 📂 Project Structure

```text
rs-vision-hub-go/
├── cmd/hub/main.go          # Application entry point, main loop
├── pkg/
│   ├── bridge/              # Middleware adapter (CGO encapsulation)
│   ├── processor/           # Image processing (Conversion, Colorization, Extraction)
│   ├── hud/                 # HUD rendering (Text, Crosshair)
│   └── display/             # Window management
├── go.mod                   # Dependency management
└── Makefile                 # Build automation
```

## 🧩 Architecture

The application follows a strict pipeline:

1.  **Fetch**: Acquires synchronized frames from the middleware.
2.  **Process**: Converts raw bytes to `gocv.Mat`, colorizes depth maps, and extracts distance data.
3.  **Overlay**: Draws HUD elements (FPS, Timestamp, Crosshair) on the frames.
4.  **Render**: Displays the processed frames in synchronized windows.
5.  **Clean**: Explicitly releases CGO memory resources.

## 📄 License

This project is licensed under the MIT License.
