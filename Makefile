# main.go 文件列表
GO_FILES := $(wildcard cmd/*/*.go)

# 项目名字
PROJECT_NAME=micro-frame

# 项目包s
PKG = "github.com/rshulabs/$(PROJECT_NAME)"

# 版本所需变量值
BUILD_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT = ${shell git rev-parse HEAD}
BUILD_TIME = ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION = $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH = "${PKG}/pkg/version"

# 编译所有 main文件 -ldflags 编译时在version包里绑定变量
.PHONY: build
build: $(GO_FILES)
	@go build -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/ $(GO_FILES)

# Cleaning all build output
.PHONY: clean
clean:
	@-rm -vrf dist

# 错误码生成
.PHONY: codegen
codegen: 
	@mkdir -p docs/code
	@codegen -type=int ./internal/pkg/code
	@codegen -type=int -doc -output ./docs/code/error_code_generated.md ./internal/pkg/code

# go mod tidy
.PHONY: tidy
tidy:
	@go mod tidy

# complie proto file
# @protoc -I=. -I=common/pb --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto
.PHONY: proto
proto: 
	@protoc -I=. --go_out=. --go-grpc_out=. internal/*/impl/pb/*proto

# 运行demo
.PHONY: run.demo
run.demo:
	@go run cmd/demo/demo.go --config="configs/demo.yaml"

# 打印demo帮助信息
.PHONY: help.demo
help.demo:
	@go run cmd/demo/demo.go --help

