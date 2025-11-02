# 快速开始指南

## 安装

```bash
go get github.com/ccxt-go/ccxt-go
```

## 基础示例

### 1. 使用已有交易所（如 Binance）

```go
package main

import (
    "fmt"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 创建 Binance 交易所实例
    binance := &ccxt.Binance{}
    binance.ExchangeBase = &ccxt.ExchangeBase{}
    binance.Setup(ccxt.MkMap(&ccxt.VarMap{}), binance)
    
    // 获取价格信息
    ticker := binance.FetchTicker(ccxt.MkString("BTC/USDT"))
    fmt.Printf("BTC/USDT 价格: %s\n", ticker.ToStr())
}
```

### 2. 使用新增交易所（如 Alpaca）

```go
package main

import (
    "fmt"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 创建 Alpaca 交易所实例
    alpaca := &ccxt.Alpaca{}
    alpaca.ExchangeBase = &ccxt.ExchangeBase{}
    
    // 配置交易所（可选）
    config := ccxt.MkMap(&ccxt.VarMap{
        "apiKey": ccxt.MkString("your_api_key"),
        "secret": ccxt.MkString("your_secret"),
    })
    alpaca.Setup(config, alpaca)
    
    // 加载市场信息
    markets := alpaca.LoadMarkets()
    fmt.Printf("加载了 %d 个交易市场\n", markets.Length.ToInt())
    
    // 获取价格信息（待实现）
    ticker := alpaca.FetchTicker(ccxt.MkString("BTC/USDT"))
    fmt.Printf("价格信息: %s\n", ticker.ToStr())
}
```

### 3. 使用 WebSocket（所有交易所都支持）

```go
package main

import (
    "fmt"
    "time"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 创建交易所实例
    exchange := &ccxt.Alpaca{}  // 或任何其他交易所
    exchange.ExchangeBase = &ccxt.ExchangeBase{}
    exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
    
    // 建立 WebSocket 连接
    wsConn := exchange.UnifiedWebSocketConnect(
        ccxt.MkString("/ws/market"),  // WebSocket 路径
        ccxt.MkMap(&ccxt.VarMap{}),
    )
    
    if wsConn.Type != ccxt.Error {
        fmt.Printf("WebSocket 连接成功\n")
        
        // 订阅消息
        subscription := exchange.UnifiedWebSocketSubscribe(
            wsConn,
            ccxt.MkString("ticker:BTC/USDT"),
        )
        
        fmt.Printf("订阅结果: %s\n", subscription.ToStr())
        
        // 等待一段时间
        time.Sleep(5 * time.Second)
        
        // 关闭连接
        exchange.UnifiedWebSocketClose(wsConn)
        fmt.Printf("连接已关闭\n")
    }
}
```

### 4. 统一 HTTP 接口

```go
package main

import (
    "fmt"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 创建交易所实例
    exchange := &ccxt.Alpaca{}
    exchange.ExchangeBase = &ccxt.ExchangeBase{}
    exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
    
    // 使用统一 HTTP 接口调用 API
    result := exchange.UnifiedHTTPRequest(
        ccxt.MkString("/api/v1/markets"),  // 路径
        ccxt.MkString("public"),           // API 类型
        ccxt.MkString("GET"),              // HTTP 方法
        ccxt.MkMap(&ccxt.VarMap{}),        // 查询参数
        ccxt.MkMap(&ccxt.VarMap{}),        // 请求头
        ccxt.MkUndefined(),                // 请求体
    )
    
    if result.Type != ccxt.Error {
        fmt.Printf("API 响应: %s\n", result.ToStr())
    } else {
        fmt.Printf("请求失败: %s\n", result.ToStr())
    }
}
```

## 支持的新增交易所

以下 37 个交易所都可以直接使用（HTTP 方法实现待完善）：

### 新兴/专业交易所
- `ccxt.Alpaca`
- `ccxt.Apex`
- `ccxt.Arkham`
- `ccxt.Backpack`
- `ccxt.Bingx`
- `ccxt.Hyperliquid`
- `ccxt.Paradex`

### 衍生品交易所
- `ccxt.Krakenfutures`
- `ccxt.Kucoinfutures`

### 特定地区交易所
- `ccxt.Bitopro`
- `ccxt.Bitrue`
- `ccxt.Bitteam`
- `ccxt.Bittrade`
- `ccxt.Coinsph`
- `ccxt.Htx`
- `ccxt.Tokocrypto`
- `ccxt.Zonda`
- `ccxt.Xt`
- `ccxt.Toobit`

### DeFi/链上交易所
- `ccxt.Defx`
- `ccxt.Derive`

### 其他交易所
- `ccxt.Blofin`
- `ccxt.Coincatch`
- `ccxt.Coinmetro`
- `ccxt.Cryptomus`
- `ccxt.Fmfwio`
- `ccxt.Foxbit`
- `ccxt.Gate`
- `ccxt.Hashkey`
- `ccxt.Hibachi`
- `ccxt.Mexc`
- `ccxt.Modetrade`
- `ccxt.Onetrading`
- `ccxt.Oxfun`
- `ccxt.P2b`
- `ccxt.Woo`
- `ccxt.Woofipro`

## 查看所有支持的交易所

```go
package main

import (
    "fmt"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 获取所有交易所
    exchanges := ccxt.Exchanges
    
    fmt.Printf("支持的交易所数量: %d\n", len(exchanges))
    
    for i, exchange := range exchanges {
        describe := exchange.Describe()
        id := *(describe).At(ccxt.MkString("id"))
        name := *(describe).At(ccxt.MkString("name"))
        fmt.Printf("%d. %s (%s)\n", i+1, name.ToStr(), id.ToStr())
    }
}
```

## 注意事项

1. **HTTP 方法实现**: 新增交易所的 HTTP 方法目前是占位符实现，返回空数据。需要根据实际 API 完善实现。

2. **WebSocket 支持**: 所有交易所都自动支持 WebSocket，但需要根据交易所的实际 WebSocket 端点配置路径。

3. **认证**: 对于需要认证的 API 调用，需要在 `Setup()` 中配置 API Key 和 Secret。

4. **错误处理**: 始终检查返回结果的错误状态：
   ```go
   result := exchange.FetchTicker(...)
   if result.Type == ccxt.Error {
       // 处理错误
   }
   ```

## 更多资源

- [完整文档](README.md)
- [新增交易所文档](NEW_EXCHANGES.md)
- [统一 API 文档](UNIFIED_API.md)
- [交易所对比](EXCHANGE_COMPARISON_DETAILED.md)

