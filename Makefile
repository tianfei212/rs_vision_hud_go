# Makefile for rs-vision-hub-go

# 定义中间件模块名称
MIDDLEWARE_MOD := github.com/tianfei212/jetson-rs-middleware
# 目标二进制文件名
APP_NAME := rs-vision-hub

.PHONY: all deps build run clean

all: deps build

deps:
	@echo "=> 下载并更新 Go 依赖..."
	go mod tidy
	go mod download $(MIDDLEWARE_MOD)

build:
	@echo "=> 自动化定位中间件缓存路径..."
	$(eval MDW_DIR := $(shell go list -m -f '{{.Dir}}' $(MIDDLEWARE_MOD)))
	@echo "   中间件路径: $(MDW_DIR)"
	@echo "=> 编译并嵌入 RPATH..."
	go build -ldflags="-r $(MDW_DIR)/lib" -o $(APP_NAME) ./cmd/hub/main.go
	@echo "=> 构建完成: $(APP_NAME)"

run: build
	@echo "=> 启动应用程序..."
	./$(APP_NAME)

clean:
	@echo "=> 清理构建产物..."
	rm -f $(APP_NAME)