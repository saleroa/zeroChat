# 所有操作都会在这个文件里被执行
# 运行具体的 makefile 将服务代码打包成二进制，二进制制作成docker镜像
# 删除 bin 目录下的所有二进制文件，方便修改
# docker compose 启动
# 进入 mysql 创建 数据库 和 数据表



run: server compose-start

stop: compose-stop rm-components rm-image



 
compose-start:
	docker-compose up -d

compose-stop:
	docker-compose down

rm-image:
	$(shell docker images | grep "rpc\|api" | awk '{print$3}' | xargs -r docker rmi)

rm-components:
	rm -r ./components/mysql


# server part 

server: user-server 

user-server: user-rpc user-api

user-rpc:
    # 执行 build 目标，构建二进制文件
	@make -f deploy/makefile/user_rpc.mk build

user-api:
	@make -f deploy/makefile/user_api.mk build









