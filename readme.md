

发送消息请求的入口 serverWs

获取到 ws 发送的消息 > ack部分 > conn.message > handle write > chat route > kafka > 消费，记录数据 > push route >  真正的写到对方的连接中 write conn


ack 部分

第一次消息 > appendMsgMq  >  conn.readmsg  > readack  > 将 ack 写到具体 conn 连接中


消息已读

markchat route -> consume -> 更新数据库里的 bitmap -> 创建好要推送的已读消息，也就是 会话id + 消息id 对应 消息的 bitmap 
如果是私聊 -> 推送到 push 的管道 -> baseMsgTransfer转发 -> push route
 真正的写到对方的连接中 write conn
 
如果是群聊，不开启消息已读的缓存，如上
如果是群聊, 开启消息已读缓存 -> 如果存在groupReadmsg -> 合并 
                           -> 如果不存在 groupReadmag -> 创建 -> 开启携程定时或者定量推送到 push 管道



破案

websocket 的 writemessage 是将消息写到客户端展示
             readmessgae 是将请求发出的消息写出来



封装websocket，实现了连接的心跳检测
实现了消息成功接收的ack以及多种ack模式
使用kafka对消息收发进行异步处理
手写 bitmap ，实现了消息的已读未读
使用dockercompose进行部署，使用 makefile 进行管理

客户端自动 ack 没有写


## 消息已读设计 
存储
以群列表做接受列表，再做已读列表，记录用户id
使用 bit 记录用户是否储存 ， 类似于布隆过滤器
自己实现 bitmap 

## 手搓bitmap 
bitmap 是一个存储 01 的数组，1读 0未读。uid 确定该在哪，来去顶已读未读 
勉强把 byte 当 bit 来使用 