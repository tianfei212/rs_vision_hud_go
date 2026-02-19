.PHONY: build run clean

# Jetson 常见的 OpenCV 和 RealSense 路径，可根据实际环境调整
CGO_CFLAGS := -I/usr/include/opencv4 -I/usr/include/librealsense2
CGO_LDFLAGS := -L/usr/lib/aarch64-linux-gnu -lopencv_core -lopencv_imgproc -lopencv_highgui -lrealsense2

build:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -o rs-vision-hub ./cmd/hub/main.go

run: build
	./rs-vision-hub

clean:
	rm -f rs-vision-hub
	go clean