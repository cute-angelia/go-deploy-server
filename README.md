# go-deploy 上线发布工具

上线发布工具：更新，回滚; 支持svn、git

为什么会有这个项目？

发布工具大概有这么几种：

1. 发布机模式
   单独部署一台发布机，它拉取代码再通过配置ssh授信分发到各个服务器，如：Jenkins，瓦力等
   
2. 客户端模式（本工具）
   在各个服务器放置一个客户端，服务端发送更新指令，客户端去更新代码

项目因安全性考虑，无法相互授信ssh，也就无法使用发布机

###  Screenshot
![](https://github.com/cute-angelia/go-deploy-server/blob/master/screenshot1.png)
![](https://github.com/cute-angelia/go-deploy-server/blob/master/screenshot2.png)

### 特性
- 平台支持 windows, linux
- 支持 svn 和 git
- 前端代码Vue 3.0
- 更新、回滚
- server 和 client 采用 tcp 通讯 + 心跳保活 节点在线状态实时监控
- 支持 befor_deploy、after_deploy 部署前和部署后的hook命令，清理缓存、执行重启等操作
 
 ### 本机开发
 1. 编译前端vue
    
    ```
    cd cmd/server/vue
    pnpm run serve
    pnpm run build
    ```
    
 2. 启动后端服务
 
    ```
    cp server.example.json server.json
    go run main.go -c server.json
    ```
 
 3. 启动客户端（可选）
 
    ```
    go run client.go -l :8093
    ```
 
### 部署到线上

上线：
    1. 上传server二进制文件和配置文件到服务器，启动
    2. 上传client二进制文件到需要业务服务器，启动
    3. 更新配置文件节点，强烈建议在配置文件配置 jwt_secret
    4. 打开浏览器查看web管理界面 http://ip:port 是否可以正常访问
    5. 修改一下server vue 里面 isLocal 里面一些授权逻辑，这里接入了企业微信

### 配置文件说明

```json
{
  // 站点监听端口，可以通过 http 打开
  "listen_http": ":8082",
  // 日志输出模式
  "debug": true,
  // jwt 的 secret
  "jwt_secret" : "",
  "apps": [
    {
      "name": "douke", // 名称
      "type": "svn",   // 类型：svn git
      "url": "svn://192.168.1.207/xiaohua", // 地址
      "fetchlogpath": "/data/fetch_log/xiaohua/trunk", // version < 1.0.1 服务端机器目录， 提交日志文件地址：svn 留空；  git 请填一个地址存储提交日志，这个目录的作用仅仅是为了获取git提交日志用
      "node": [ // 机器节点
        {
          "alias": "api-1", // 名称
          "addr": "192.168.1.11:8081", // 客户端监听端口
          "path": "/data/wwwroot/xiaohua/trunk",  // 项目地址
          "befor_deploy": "", // 执行前操作
          "after_deploy": "" // 执行后操作
        },
        {
          "alias": "api-2",
          "addr": "192.168.1.207:8081",
          "path": "/data/wwwroot/xiaohua/trunk",
          "befor_deploy": "ls /data/wwwroot",
          "after_deploy": "ls /data/wwwroot/xiaohua/trunk"
        }
      ]
    } 
  ]
}
```
 

### 一些问题

1. 无法拉取日志：

可以用 `svn log --limit 10 svn://x.x.x.x/project/` 执行下看下错误信息


2. windows，用 pm2 守护的进程，没法 IPC 通信通知结束进程
   大致解决思路，采用tty远程手动结束，暂时这里不处理了


3. web管理访问安全问题

4. 部署完后web管理界面直接暴露给所有人的，可以加一层nginx反向代理，设置一个auth认证

```bash
htpasswd -c /usr/local/openresty/nginx/conf/vhost/passwd.db yourusername

server {
   listen 80;
   server_name yourdomain;
   auth_basic "User Authentication";
   auth_basic_user_file /usr/local/openresty/nginx/conf/vhost/passwd.db;
   location / {
       proxy_pass http://127.0.0.1:http_port;
   }
}
```



### 协议设计

```
 *  |----------Len--------|------------------------------------MetaInfo------------------------------------|
 *  |---------4Byte-------|--------4Byte--------|----1Byte----|----------------X---------------------------|
 *	+------------------------------------------------------------------------------------------------------+
 *	|      Magic          |     PayloadLen      |   MsgType   |                Body                        |
 *	+------------------------------------------------------------------------------------------------------+
 
PayloadLen - length of body in byte, 3 bytes big-endian integer.
Magic - magic number
MsgType - package type, 1 byte
    0x01: heartbeat package
    0x02: data package
    0x03: disconnect message from server
body - binary payload.
```