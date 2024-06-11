f f f f f f### 



ack 机制
客户端向服务端的 ack


消息请求的入口 serverWs

获取到 ws 发送的消息 > ack部分 > conn.message > handle write > chat route > kafka > 消费，记录数据 > push route >  真正的写到对方的连接中 write conn


ack 部分

第一次消息 > appendMsgMq  >  conn.readmsg  > readack  > 将 ack 写到具体 conn 连接中


破案

websocket 的 writemessage 是将消息写到客户端展示
             readmessgae 是将请求发出的消息写出来



封装websocket，实现了连接的心跳检测
实现了消息成功接收的ack以及多种ack模式
使用kafka对消息收发进行异步处理
使用实现了消息的已读未读显示
使用dockercompose进行部署，使用 makefile 进行管理

客户端自动 ack 没有写