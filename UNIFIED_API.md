# CCXT-Go 统一调用接口文档

## 概述

CCXT-Go 现在提供了统一的 HTTP 和 WebSocket 调用接口，支持所有加密货币交易所的 API 调用。该接口提供了以下功能：

- 🔄 **统一HTTP客户端** - 支持所有REST API调用
- 🌐 **WebSocket支持** - 实时数据流订阅
- ⚡ **异步处理** - 并发请求和消息处理
- 🔒 **速率限制** - 自动管理API调用频率
- 🔄 **自动重试** - 网络错误自动重试机制
- 🔌 **连接池** - 高效的连接管理
- 🛡️ **错误处理** - 完善的错误处理机制

## 核心组件

### 1. NetworkManager (网络管理器)
全局网络管理器，负责管理所有HTTP和WebSocket连接。

### 2. UnifiedClient (统一客户端)
为每个交易所实例提供统一的调用接口。

### 3. WebSocketConnection (WebSocket连接)
封装WebSocket连接，支持订阅、发送消息和自动重连。

### 4. RateLimiter (速率限制器)
管理API调用频率，防止超出交易所限制。

## 使用方法

### HTTP请求

```go
// 创建交易所实例
binance := &Binance{}
binance.ExchangeBase = &ExchangeBase{}
binance.Setup(MkMap(&VarMap{}), binance)

// 发送HTTP请求
result := binance.UnifiedHTTPRequest(
    MkString("/api/v3/exchangeInfo"),  // 路径
    MkString("public"),                // API类型 (public/private)
    MkString("GET"),                   // HTTP方法
    MkMap(&VarMap{}),                  // 查询参数
    MkMap(&VarMap{}),                  // 请求头
    MkUndefined(),                     // 请求体
)

if result.Type != Error {
    fmt.Printf("请求成功: %s\n", result.ToStr())
}
```

### WebSocket连接

```go
// 建立WebSocket连接
wsConn := binance.UnifiedWebSocketConnect(
    MkString("/ws/btcusdt@ticker"),  // WebSocket路径
    MkMap(&VarMap{}),                // 连接参数
)

if wsConn.Type != Error {
    // 订阅消息
    subscription := binance.UnifiedWebSocketSubscribe(
        wsConn,                      // 连接ID
        MkString("ticker"),          // 订阅主题
    )
    
    // 发送消息
    sendResult := binance.UnifiedWebSocketSend(
        wsConn,                      // 连接ID
        MkString("ping"),           // 消息内容
    )
    
    // 关闭连接
    closeResult := binance.UnifiedWebSocketClose(wsConn)
}
```

### 私有API调用

```go
// 设置API密钥
binance.Setup(MkMap(&VarMap{
    "apiKey": MkString("your_api_key"),
    "secret": MkString("your_secret"),
}), binance)

// 调用私有API
balance := binance.UnifiedHTTPRequest(
    MkString("/api/v3/account"),
    MkString("private"),
    MkString("GET"),
    MkMap(&VarMap{}),
    MkMap(&VarMap{}),
    MkUndefined(),
)
```

## 高级功能

### 速率限制

```go
// 设置自定义速率限制
GlobalNetworkManager.rateLimiter.SetRateLimit("binance", 10) // 每分钟10个请求

// 检查速率限制
if GlobalNetworkManager.rateLimiter.Allow("binance") {
    // 执行请求
} else {
    // 请求被限制
}
```

### 并发请求

```go
// 并发发送多个请求
done := make(chan bool, 5)

for i := 0; i < 5; i++ {
    go func(index int) {
        result := binance.UnifiedHTTPRequest(
            MkString("/api/v3/ping"),
            MkString("public"),
            MkString("GET"),
            MkMap(&VarMap{}),
            MkMap(&VarMap{}),
            MkUndefined(),
        )
        
        if result.Type != Error {
            fmt.Printf("协程 %d 请求成功\n", index)
        }
        
        done <- true
    }(i)
}

// 等待所有请求完成
for i := 0; i < 5; i++ {
    <-done
}
```

### 错误处理

```go
result := binance.UnifiedHTTPRequest(
    MkString("/api/v3/invalid"),
    MkString("public"),
    MkString("GET"),
    MkMap(&VarMap{}),
    MkMap(&VarMap{}),
    MkUndefined(),
)

if result.Type == Error {
    fmt.Printf("请求失败: %s\n", result.ToStr())
    // 处理错误
}
```

## 配置选项

### HTTP请求配置

```go
config := &RequestConfig{
    URL:       "https://api.binance.com/api/v3/ping",
    Method:    "GET",
    Headers:   map[string]string{"User-Agent": "ccxt-go"},
    Body:      nil,
    Timeout:   30 * time.Second,
    Retry:     true,
    RateLimit: "binance",
    Proxy:     "",
    UserAgent: "ccxt-go/1.0",
}
```

### WebSocket配置

```go
config := &WebSocketConfig{
    URL:          "wss://stream.binance.com:9443/ws/btcusdt@ticker",
    Headers:      map[string]string{"User-Agent": "ccxt-go"},
    Reconnect:    true,
    PingInterval: 30 * time.Second,
    PongTimeout:  10 * time.Second,
    ReadTimeout:  60 * time.Second,
    WriteTimeout: 10 * time.Second,
    Subscriptions: []string{"ticker"},
}
```

## 性能优化

### 连接池
网络管理器自动管理HTTP连接池，提高请求效率。

### 自动重试
网络错误时自动重试，支持指数退避算法。

### 内存管理
WebSocket连接使用缓冲通道，避免内存泄漏。

## 最佳实践

1. **合理设置速率限制** - 根据交易所API限制设置合适的速率
2. **及时关闭连接** - 使用完毕后及时关闭WebSocket连接
3. **错误处理** - 始终检查返回结果的错误状态
4. **并发控制** - 避免过多并发请求导致API限制
5. **资源清理** - 程序退出时调用 `GlobalNetworkManager.CloseAll()`

## 示例代码

完整的使用示例请参考 `examples.go` 文件，包含：

- 基础HTTP请求示例
- WebSocket连接和订阅示例
- 私有API调用示例
- 并发请求示例
- 错误处理示例
- 性能测试示例

## 测试

运行测试：

```bash
go test ./pkg/ccxt -v
```

运行性能测试：

```bash
go test ./pkg/ccxt -bench=.
```

## 故障排除

### 常见问题

1. **WebSocket连接失败**
   - 检查网络连接
   - 验证WebSocket URL格式
   - 确认交易所支持WebSocket

2. **HTTP请求超时**
   - 增加超时时间
   - 检查网络稳定性
   - 验证API端点

3. **速率限制**
   - 调整速率限制设置
   - 使用指数退避重试
   - 分散请求时间

4. **认证失败**
   - 验证API密钥
   - 检查签名算法
   - 确认权限设置

## 更新日志

- **v1.0.0** - 初始版本，支持基础HTTP和WebSocket功能
- **v1.1.0** - 添加速率限制和自动重试
- **v1.2.0** - 优化连接池和并发处理
- **v1.3.0** - 完善错误处理和日志记录
