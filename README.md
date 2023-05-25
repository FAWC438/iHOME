# iHome

> 目前进度：Redis 功能开发 - 手机验证码微服务化

一个以微服务架构为主体，租房业务为流程的项目案例

本项目主要关注后端开发内容，前端相关（即 view ）内容不在本仓库中展示，但在测试和 demo 中会出现，请知悉。

## 图片验证码

采用 <www.github.com/afocus/captcha > 库实现图片验证码微服务

### 微服务创建步骤

1. 在 [Service](src/service) 文件夹中新建 go-micro 项目文件 [getCaptcha](src/service/getCaptcha)。并配置将其注册到 Consul
   服务发现。
2. 在微服务端的 [handle](src/service/getCaptcha/handler/getCaptcha.go) 目录中编写业务代码
3. [控制器中](src/controller/)通过 consul 调用微服务实现具体业务

> 在单机上运行测试，需要 gin 服务前台获得 go-micro 服务后台的 protobuf 。为了方便本地导包，需要在 **项目** `go.mod`
> 文件中，添加 `replace getCaptcha => ./src/service/getCaptcha` 然后在 gin 的控制器中调用

### 运行测试

1. 先运行 consul 服务
2. 再启动 go-micro 微服务
3. 启动 gin ，打开 `http://127.0.0.1:8080/home/register.html` 测试验证码功能是否实现

## Redis

> 教程中提到的 Redis
> 常见面试知识点：[Redis 持久化方式 RDB 和 AOF 的对比](https://blog.csdn.net/Aa112233aA1/article/details/124245231)

通过 WSL 在 [Windows 下安装 Redis](https://redis.io/docs/getting-started/installation/install-redis-on-windows/)

### 客户端

#### 可视化客户端

然后可以选择安装一个 Redis 可视化客户端。这里推荐一个开源好用的国人开发的 Redis
客户端 [AnotherRedisDesktopManager](https://github.com/qishibo/AnotherRedisDesktopManager) ，在 Windows 下（不是在
WSL 里）使用 `winget` ：

```shell
winget install qishibo.AnotherRedisDesktopManager -i
```

如果在 Windows 环境下，独立显卡会导致 AnotherRedisDesktopManager 白屏。需要在 AnotherRedisDesktopManager
快捷方式中添加启动项 `--disable-gpu`
，详见 [github-issue](https://github.com/qishibo/AnotherRedisDesktopManager/issues/887) 。

#### go 客户端

使用 [redigo](https://github.com/gomodule/redigo) 作为 go 的 redis 客户端。

> 关于 redis 连接池的操作不同客户端处理不同
>
> - go-redis 是自动管理，类似 go/sql 包的方式，在真正执行的时候从连接池取一个连接，执行完毕后放回去，对调用者透明。
    调用者如果手动关闭连接，连接不能被复用，表现上看就是 redis 服务器的 tcp 新建连接数特别多，而业务机器的 timewait 数量大。
>
> - redigo 是手动管理，调用者需要明确获取一个连接，执行完毕再手动关闭。不及时关闭，会造成连接池泄露，表现上看就是 redis
    的连接数持续增长
>
> 总结就是 go-redis 不要调用 close ，而 redigo 需要调用 close ，正好相反。

### WSL 中启动 redis

在 WSL 中使用命令 `sudo service redis-server start` 启动 Redis。

如果需要，可以在 `./etc/redis/redis.conf`
中修改配置，在了解这方面的时候，我看到一篇比较重要的文章贴在这里以供参考：[ Redis 的 bind 的误区](https://blog.csdn.net/cw_hello1/article/details/83444013)

## 手机验证码

利用[阿里云平台短信服务](https://dysms.console.aliyun.com/overview)来实现短信验证码校验

本仓库提供的阿里云用户信息权限仅能使用个人用户短信服务的测试功能。

### 接入 redis

引入 redis 连接池，redigo 相关[文档](https://pkg.go.dev/github.com/gomodule/redigo/redis@v1.8.9#Pool)。

需要注意的小细节是，即使是调用从连接池 `Get()` 来的连接，使用时也必须使用 `defer` 关闭连接。

### 敏感数据存储

手机号码属于敏感隐私数据，因此不建议明文存储在 redis 中。此处采用了 AES 的 CRT 模式加密了手机号码并将其存储到 redis 里

## Gorm 与 MySQL

Gorm [官网](https://gorm.io/)，[中文文档](https://gorm.io/zh_CN/docs/index.html)。

> 注意，Gorm 的判断语句（where、select等）必须写在查询语句（take、find等）的前面

### 接下来的工作

- 数据 redis 存储
- 微服务化
