### 

user-rpc 监听 9001 
user-api 监听 8001

social-rpc 9002
social-apt 8002

ws 8003
mq 8004



mysql 和 redis 密码  
root123456

根据 proto 生成 

goctl rpc protoc apps/social/rpc/social.proto --go_out=./apps/social/rpc --go-grpc_out=./apps/social/rpc --zrpc_out=./apps/social/rpc

根据 sql 生成

goctl model mysql ddl -src="./sql/user.sql" -dir="./apps/user/models/" -c

根据 api 生成 
goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero



ack 机制
客户端向服务端的 ack
ack 次数过多的报错暂停  done