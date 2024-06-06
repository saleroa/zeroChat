

SERVER_NAME=social
SERVER_TYPE=rpc

# 镜像 和 二进制文件 的名字都是 ${SERVER_NAME}-${SERVER_TYPE}
IMAGE_NAME=${SERVER_NAME}-${SERVER_TYPE}
BINARY_NAME=${SERVER_NAME}-${SERVER_TYPE}

# dockerfile 的名字 
DOCKER_FILE=./deploy/dockerfile/${SERVER_NAME}_${SERVER_TYPE}_dockerfile

# 测试环境的编译发布
build:
	# 将服务构建成二进制文件，放在 bin 目录下
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${BINARY_NAME} ./apps/${SERVER_NAME}/${SERVER_TYPE}/${SERVER_NAME}.go
	# 将根据 dockerfile 构建成一个镜像
	docker build . -f ${DOCKER_FILE} --no-cache -t ${IMAGE_NAME}
