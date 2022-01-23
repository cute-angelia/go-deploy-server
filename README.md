# go-deploy

上线发布工具：更新，回滚; 支持svn、git

发布工具市面大概有这么几种：

1. 发布机模式
   单独部署一台发布机，它拉取代码再同步到各个服务器，Jenkins，瓦力发布都是这种模式
   
2. 客户端模式（本工具）
   在各个服务器放置一个客户端，接受服务端的命令，拉去对应代码

为什么会有这个工具？

自身项目因网络隔离，再部署一台Jenkins太笨重，用这个更简单适用，但是需要自动化编译，还是得上Jenkins

2. 前端代码Vue 3.0
3. 注重安全，接口过滤敏感信息

# Screenshot
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
 1. 启动前端代码
    
    ```
    cd cmd/server/vue
    npm run serve
    ```
    
 2. 启动后端服务端
 
    ```
    go run server.go -c server.json
    ```
 
 3. 启动后端客户端（可选）
 
    ```
    go run client.go -l :8093
    ```
 
### 部署到线上

上线：
    1. 强烈建议在配置文件配置 jwt_secret
    
 
 
### Server端配置文件
```
|---server
    |---server.json
```
 
1. 上传server二进制文件到服务器
2. 创建或修改server.json，添加项目以及相应的节点。
3. 获取commit log
   如果项目类型为svn则fetchlogpath设置为空，切换到应用程序执行账户下，这里假设为www账户，手动获取一次日志，目的是让工具拉日志的时候可以免密。
```bash
su www
svn log --limit 10 svn://x.x.x.x/project/
```
   如果项目类型为git则fetchlogpath需要配置一个目录，并将项目拉取一份在这个目录下。请不要将任何web站点指向这个目录，这个目录的作用仅仅是为了获取git提交日志用。

4. 使用www账号启动server
```bash
su www
chmod +x /pathto/server
/pathto/server -c /pathto/server.json
```
     
5. 打开浏览器查看web管理界面 http://ip:port 是否可以正常访问 (port为server.json里配置的listen_http 监听端口)

### Client节点配置
```
 |---client
```

1. 上传client二进制文件 到所有的节点服务器，Node节点只用部署一个单文件(以下操作每个Node节点服务器都需要执行同样的操作)

2. 首选拉取一份项目代码，假设www为你的web运行账户(使用ps aux|grep php-fpm查看)
   查看www账号是否有shell权限
```bash
cat /etc/passwd|grep www 
www:x:501:501::/home/www:/bin/bash
```
   如果没有shell权限则需要执行下面命令
```bash
usermod -s /bin/bash www #www是用户名
```
   拉取项目代码，确保下次执行可以免密更新
```bash
su www
svn co svn://1.2.3.4/app /data/wwwroot/app
cd /data/wwwroot/app
svn up
```
3. 启动client
```bash
su www
chmod +x /pathto/client
/pathto/client -l :8081 #8081为client监听端口
```

### web管理访问安全问题
   部署完后web管理界面直接暴露给所有人的，可以加一层nginx反向代理，设置一个auth认证
```bash
htpasswd -c /usr/local/openresty/nginx/conf/vhost/passwd.db yourusername
```
   配置nginx
```bash
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

### 已知问题

1. windows，用 pm2 守护的进程，没法 IPC 通信通知结束进程
   大致解决思路，采用tty远程手动结束，暂时这里不处理了

# 参考

项目代码基于[ikool-cn/go-deploy](https://github.com/ikool-cn/go-deploy/)

