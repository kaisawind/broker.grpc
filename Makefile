.PHONY: grpc build test

all: build

build: test
	$(GO) build $(FLAGS)broker .

test:
	$(CGO) test -v ./...

grpc:
	protoc -I ./pb/. --go_out=plugins=grpc,paths=source_relative:./pb ./pb/broker.proto

docker:
	docker buildx build \
		--platform=linux/arm64,linux/amd64 \
		-f scripts/Dockerfile \
		-t kaisawind/broker.grpc . \
		--push

##### ^^^^^^ EDIT ABOVE ^^^^^^ #####

##### =====> Internals <===== #####

# 版本号 v1.0.3-6-g0c2b1cf-dev
# 1、6:表示自打tag v1.0.3以来有6次提交（commit）
# 2、g0c2b1cf：g 为git的缩写，在多种管理工具并存的环境中很有用处
# 3、0c2b1cf：7位字符表示为最新提交的commit id 前7位
# 4、如果本地仓库有修改，则认为是dirty的，则追加-dev，表示是开发版：v1.0.3-6-g0c2b1cf-dev
VERSION          := $(shell git describe --tags --always --dirty="-dev")

# 时间
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')

# 版本标志  -s -w 缩小可执行文件大小
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)" -s -w'

# 输出文件夹
OUTPUT_DIR       := -o ./bin/

# 标志
FLAGS            := $(VERSION_FLAGS) $(OUTPUT_DIR)

GO        		 := CGO_ENABLED=0 GO111MODULE=on go
CGO        		 := CGO_ENABLED=1 GO111MODULE=on go