# webkit

### 简介

> webkit是用于快捷生成web项目的小工具，集成了一些常见功能，如：优雅关闭、日志切分、参数校验等

### 安装&使用

- 安装

```shell
go install github.com/LorinHan/webkit@latest
```

- 使用

在指定路径下生成模板项目

```shell
webkit
#> webkit v1.0.0
#> 请输入项目名称（默认'test_webkit'）：
#> 请输入项目路径（默认'./'）：
```

### 功能说明

> 在生成的项目中直接内置了一些常见功能

##### 1.优雅关闭

```go
main.go

...
// 创建一个优雅关闭的超时上下文
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// 关闭 HTTP 服务器
if err := server.Shutdown(ctx); err != nil {
    zap.S().Fatal("Server shutdown:", err)
}
zap.S().Info("Server exited")
...
```

##### 2.集成zap日志和lumberjack切分

- 初始化

```go
main.go

...
// logger.DefaultLog() 默认只输出到终端，debug级别
logger.Init(logger.DefaultLog())
defer logger.Sync()
...
```

可以根据需要调整配置项

```go
kit/logger/logger.go

type Option struct {
    Path       string // 日志文件路径
    Level      string // 日志级别，debug info warn error panic fatal
    MaxSize    int    // 文件多大开始切分
    MaxBackups int  // 保留文件个数
    MaxAge     int  // 文件保存多少天，maxBackups和maxAge都设置为0，则不会删除任何日志文件，全部保留
    Json       bool // 是否用json格式
    Std        bool   // 是否输出到控制台
}
```

##### 3.配置项

- 初始化

```go
main.go

...
config.InitByEnv() // 从环境变量加载配置
// config.InitByFile("config.yaml") // 从配置文件加载配置（viper库），支持yaml、json、toml等多种格式
...

```

##### 4.参数校验

> 集成了[validator库](https://github.com/go-playground/validator)

- 初始化

```go
main.go

if err := validator.Init(); err != nil {
    zap.S().Fatal(err)
}
```

- 使用

在参数字段上写tag即可，更多tag参考[validator库](https://github.com/go-playground/validator)

```go
type TestValidatorReq struct {
    Name     string `json:"name" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email" errMsg:"自定义错误信息:邮箱错误咯"`
    Password string `json:"password" binding:"required,min=6"`
}
```

- 集成了中文翻译器

```
{
    "code": 400001,
    "message": "参数错误, name长度必须至少为3个字符",
    "data": null
}
```

- 自定义字段名

```go
// tag中加入 name="xxx"
type TestValidatorReq struct {
    Password     string `json:"password" binding:"required,min=3,max=50" name="密码"`
}

// 若该字段校验失败，响应的message会将字段"password"替换为"密码"
{
    "code": 400001,
    "message": "参数错误, 密码长度必须至少为3个字符",
    "data": null
}
```

- 自定义错误信息

```go
// tag中加入 errMsg="xxx"
type TestValidatorReq struct {
    Email    string `json:"email" binding:"required,email" errMsg:"自定义错误信息-邮箱错误咯"`
}

// 若该字段校验失败，响应的message会替换为自定义的errMsg
{
    "code": 400001,
    "message": "参数错误, 自定义错误信息-邮箱错误咯",
    "data": null
}
```

##### 5.gorm相关
- 结合了zap，gorm日志将由zap代理
- 错误处理，controller层使用Fail方法进行响应时，Fail方法内对gorm的err做了处理，避免数据库相关错误信息直接返回
```go
data, err := service.HelloSvc.SayHi()
if err != nil {
    Fail(ctx, enum.FailedGetData, err)
    return
}
```

##### 6.cache
- 集成go-redis
- aop封装，以切面形式进行缓存操作
```go
kit/cache/redis_aop_test.go

func TestCacheable(t *testing.T) {
    type User struct {
        ID    int     `json:"id"`
        Name  string  `json:"name"`
        Money float64 `json:"money"`
        OK    bool    `json:"ok"`
    }
    // 初始化
    Init(&redis.Options{Addr: "127.0.0.1:6379"})
    
    ctx := context.Background()
    var user User
    
    // 若key不存在，回调函数执行，并将结果放入缓存，下次再执行key已存在，不执行回调，而是将数据映射到user指针
    if err := Cacheable(ctx, "test", &user, func() (data interface{}, err error) {
        log.Println("数据库查询等操作，执行了...")
        return User{
            ID:   1,
            Name: "test",
        }, nil
    }, time.Minute*10); err != nil {
        log.Fatal(err)
    }
    
    log.Println("res", user)
}
```