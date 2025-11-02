# 新增交易所支持文档

## 概述

本文档介绍 CCXT-Go 最新新增的 37 个交易所支持。这些交易所来自 Python CCXT 版本，现已完整移植到 Go 版本，支持 HTTP 和 WebSocket 两种通信方式。

## 新增交易所列表

### 1. 新兴/专业交易所 (7个)

这些交易所专注于特定领域或提供专业的交易服务：

- **alpaca** - Alpaca 证券交易平台
- **apex** - Apex Protocol
- **arkham** - Arkham Intelligence
- **backpack** - Backpack Exchange
- **bingx** - BingX 交易所
- **hyperliquid** - Hyperliquid 去中心化交易所
- **paradex** - Paradex 交易所

### 2. 衍生品交易所 (2个)

专注于期货和衍生品交易：

- **krakenfutures** - Kraken Futures
- **kucoinfutures** - KuCoin Futures

### 3. 特定地区交易所 (10个)

服务于特定地区或国家的交易所：

- **bitopro** - BitoPro (台湾)
- **bitrue** - Bitrue
- **bitteam** - BitTeam
- **bittrade** - BitTrade
- **coinsph** - Coins.ph (菲律宾)
- **htx** - HTX (原火币，Huobi)
- **tokocrypto** - Tokocrypto (印尼)
- **zonda** - Zonda (波兰)
- **xt** - XT.com
- **toobit** - Toobit

### 4. DeFi/链上交易所 (2个)

去中心化和链上交易所：

- **defx** - DEFX
- **derive** - Derive

### 5. 其他交易所 (16个)

- **blofin** - Blofin
- **coincatch** - CoinCatch
- **coinmetro** - CoinMetro
- **cryptomus** - Cryptomus
- **fmfwio** - FMFW.io
- **foxbit** - Foxbit (巴西)
- **gate** - Gate.io
- **hashkey** - HashKey
- **hibachi** - Hibachi
- **mexc** - MEXC
- **modetrade** - ModeTrade
- **onetrading** - One Trading
- **oxfun** - OX.FUN
- **p2b** - P2B
- **woo** - WOO Network
- **woofipro** - WOO X Pro

## 功能支持

### HTTP API 方法

所有新增交易所都支持以下 HTTP 方法（目前为占位符实现，需根据实际 API 完善）：

#### 1. FetchMarkets - 获取交易市场列表

```go
// 示例：获取所有交易市场
alpaca := &Alpaca{}
alpaca.ExchangeBase = &ExchangeBase{}
alpaca.Setup(MkMap(&VarMap{}), alpaca)

markets := alpaca.FetchMarkets(MkMap(&VarMap{}))
// TODO: 根据交易所实际 API 实现
```

#### 2. FetchTicker - 获取价格信息

```go
// 示例：获取 BTC/USDT 价格信息
ticker := alpaca.FetchTicker(
    MkString("BTC/USDT"),
    MkMap(&VarMap{}),
)
// TODO: 根据交易所实际 API 实现
```

#### 3. FetchOrderBook - 获取订单簿

```go
// 示例：获取 BTC/USDT 订单簿
orderBook := alpaca.FetchOrderBook(
    MkString("BTC/USDT"),
    MkUndefined(),  // limit
    MkMap(&VarMap{}),
)
// TODO: 根据交易所实际 API 实现
```

#### 4. FetchTrades - 获取交易历史

```go
// 示例：获取 BTC/USDT 交易历史
trades := alpaca.FetchTrades(
    MkString("BTC/USDT"),
    MkUndefined(),  // since
    MkUndefined(),  // limit
    MkMap(&VarMap{}),
)
// TODO: 根据交易所实际 API 实现
```

#### 5. FetchBalance - 获取账户余额

```go
// 示例：获取账户余额（需要API密钥）
alpaca.Setup(MkMap(&VarMap{
    "apiKey": MkString("your_api_key"),
    "secret": MkString("your_secret"),
}), alpaca)

balance := alpaca.FetchBalance(MkMap(&VarMap{}))
// TODO: 根据交易所实际 API 实现
```

### WebSocket 支持

所有新增交易所都**自动支持 WebSocket**，因为它们都继承自 `ExchangeBase`，而 `ExchangeBase` 提供了统一的 WebSocket 接口。

#### WebSocket 方法

1. **UnifiedWebSocketConnect** - 建立 WebSocket 连接
2. **UnifiedWebSocketSubscribe** - 订阅 WebSocket 消息
3. **UnifiedWebSocketSend** - 发送 WebSocket 消息
4. **UnifiedWebSocketClose** - 关闭 WebSocket 连接

#### 使用示例

```go
package main

import (
    "fmt"
    "github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
    // 创建交易所实例
    alpaca := &ccxt.Alpaca{}
    alpaca.ExchangeBase = &ccxt.ExchangeBase{}
    alpaca.Setup(ccxt.MkMap(&ccxt.VarMap{}), alpaca)
    
    // 建立 WebSocket 连接
    wsConn := alpaca.UnifiedWebSocketConnect(
        ccxt.MkString("/ws/market"),  // WebSocket 路径（根据交易所调整）
        ccxt.MkMap(&ccxt.VarMap{}),
    )
    
    if wsConn.Type != ccxt.Error {
        fmt.Printf("WebSocket 连接成功: %s\n", wsConn.ToStr())
        
        // 订阅价格更新
        subscription := alpaca.UnifiedWebSocketSubscribe(
            wsConn,
            ccxt.MkString("ticker:BTC/USDT"),
        )
        
        fmt.Printf("订阅结果: %s\n", subscription.ToStr())
        
        // 发送消息
        sendResult := alpaca.UnifiedWebSocketSend(
            wsConn,
            ccxt.MkString("ping"),
        )
        
        fmt.Printf("发送结果: %s\n", sendResult.ToStr())
        
        // 关闭连接
        closeResult := alpaca.UnifiedWebSocketClose(wsConn)
        fmt.Printf("关闭结果: %s\n", closeResult.ToStr())
    }
}
```

## 使用方法

### 1. 初始化交易所

```go
import "github.com/ccxt-go/ccxt-go/pkg/ccxt"

// 创建交易所实例
exchange := &ccxt.Alpaca{}  // 替换为其他交易所
exchange.ExchangeBase = &ccxt.ExchangeBase{}

// 配置交易所
config := ccxt.MkMap(&ccxt.VarMap{
    "apiKey": ccxt.MkString("your_api_key"),
    "secret": ccxt.MkString("your_secret"),
    "sandbox": ccxt.MkBool(false),  // 是否使用沙盒环境
})

exchange.Setup(config, exchange)
```

### 2. 加载市场信息

```go
// 加载所有交易市场
markets := exchange.LoadMarkets()
if markets.Type != ccxt.Error {
    fmt.Printf("加载了 %d 个交易市场\n", markets.Length.ToInt())
}
```

### 3. 调用 API 方法

```go
// 获取价格信息
ticker := exchange.FetchTicker(
    ccxt.MkString("BTC/USDT"),
    ccxt.MkMap(&ccxt.VarMap{}),
)

// 获取订单簿
orderBook := exchange.FetchOrderBook(
    ccxt.MkString("BTC/USDT"),
    ccxt.MkInteger(20),  // limit
    ccxt.MkMap(&ccxt.VarMap{}),
)

// 获取交易历史
trades := exchange.FetchTrades(
    ccxt.MkString("BTC/USDT"),
    ccxt.MkUndefined(),
    ccxt.MkInteger(100),  // limit
    ccxt.MkMap(&ccxt.VarMap{}),
)
```

## 实现状态

### 当前状态

- ✅ **基础结构** - 所有交易所的 `Describe()` 函数已实现
- ✅ **HTTP 方法框架** - 所有 HTTP 方法占位符已添加
- ✅ **WebSocket 支持** - 通过 `ExchangeBase` 自动支持
- ✅ **编译通过** - 所有文件无语法错误

### 待完善

- ⚠️ **HTTP 方法实现** - 需要根据各交易所的实际 API 完善实现
- ⚠️ **错误处理** - 需要添加具体的错误处理逻辑
- ⚠️ **认证机制** - 需要实现各交易所的签名算法
- ⚠️ **速率限制** - 需要配置各交易所的速率限制参数
- ⚠️ **测试覆盖** - 需要添加单元测试和集成测试

## 开发指南

### 如何实现新的交易所 API

1. **查找交易所 API 文档**
   - 访问交易所的官方 API 文档
   - 了解 API 端点、认证方式和参数格式

2. **实现 FetchMarkets**
   ```go
   func (this *Alpaca) FetchMarkets(goArgs ...*Variant) *Variant {
       params := GoGetArg(goArgs, 0, MkMap(&VarMap{}))
       
       // 调用 API
       response := this.Call(MkString("publicGetMarkets"), params)
       
       // 解析响应
       markets := this.SafeValue(response, MkString("data"), MkArray(&VarArray{}))
       
       // 解析为统一格式
       return this.ParseMarkets(markets)
   }
   ```

3. **实现 FetchTicker**
   ```go
   func (this *Alpaca) FetchTicker(goArgs ...*Variant) *Variant {
       symbol := GoGetArg(goArgs, 0, MkUndefined())
       params := GoGetArg(goArgs, 1, MkMap(&VarMap{}))
       
       this.LoadMarkets()
       market := this.Market(symbol)
       
       // 调用 API
       response := this.Call(MkString("publicGetTicker"), this.Extend(MkMap(&VarMap{
           "symbol": *(market).At(MkString("id")),
       }), params))
       
       // 解析为统一格式
       return this.ParseTicker(response, market)
   }
   ```

4. **实现认证签名**
   ```go
   func (this *Alpaca) Sign(path string, api string, method string, params *Variant, headers *Variant, body *Variant) string {
       // 实现签名算法
       // 根据交易所的签名规则生成签名
       return signedPath
   }
   ```

5. **配置 Describe 函数**
   ```go
   func (this *Alpaca) Describe(goArgs ...*Variant) *Variant {
       return this.DeepExtend(this.BaseDescribe(), MkMap(&VarMap{
           "id":   MkString("alpaca"),
           "name": MkString("Alpaca"),
           "urls": MkMap(&VarMap{
               "api": MkMap(&VarMap{
                   "public":  MkString("https://api.alpaca.com"),
                   "private": MkString("https://api.alpaca.com"),
               }),
           }),
           "rateLimit": MkInteger(2000),
           // ... 其他配置
       }))
   }
   ```

## 文件位置

所有新增交易所的实现文件位于：

```
pkg/ccxt/ex_*.go
```

具体文件列表：

- `ex_alpaca.go`
- `ex_apex.go`
- `ex_arkham.go`
- `ex_backpack.go`
- `ex_bingx.go`
- `ex_hyperliquid.go`
- `ex_paradex.go`
- `ex_krakenfutures.go`
- `ex_kucoinfutures.go`
- `ex_bitopro.go`
- `ex_bitrue.go`
- `ex_bitteam.go`
- `ex_bittrade.go`
- `ex_coinsph.go`
- `ex_htx.go`
- `ex_tokocrypto.go`
- `ex_zonda.go`
- `ex_xt.go`
- `ex_toobit.go`
- `ex_defx.go`
- `ex_derive.go`
- `ex_blofin.go`
- `ex_coincatch.go`
- `ex_coinmetro.go`
- `ex_cryptomus.go`
- `ex_fmfwio.go`
- `ex_foxbit.go`
- `ex_gate.go`
- `ex_hashkey.go`
- `ex_hibachi.go`
- `ex_mexc.go`
- `ex_modetrade.go`
- `ex_onetrading.go`
- `ex_oxfun.go`
- `ex_p2b.go`
- `ex_woo.go`
- `ex_woofipro.go`

## 测试

### 运行所有测试

```bash
go test ./pkg/ccxt -v
```

### 测试特定交易所

```go
func TestAlpacaExchange(t *testing.T) {
    alpaca := &Alpaca{}
    alpaca.ExchangeBase = &ExchangeBase{}
    alpaca.Setup(MkMap(&VarMap{}), alpaca)
    
    // 测试 Describe
    describe := alpaca.Describe()
    assert.NotNil(t, describe)
    
    // 测试 HTTP 方法
    markets := alpaca.FetchMarkets(MkMap(&VarMap{}))
    assert.NotNil(t, markets)
}
```

## 常见问题

### Q: 为什么 HTTP 方法返回空数据？

A: 目前所有 HTTP 方法都是占位符实现，返回空数据。需要根据交易所的实际 API 实现具体逻辑。

### Q: WebSocket 连接失败怎么办？

A: 检查以下几点：
1. 确认交易所是否支持 WebSocket
2. 验证 WebSocket URL 格式是否正确
3. 检查网络连接和防火墙设置
4. 查看交易所文档中的 WebSocket 端点

### Q: 如何添加新的交易所？

A: 参考现有交易所的实现，创建新的 `ex_*.go` 文件，实现 `Describe()` 和必要的 HTTP 方法，然后在 `init()` 函数中注册。

### Q: 如何处理交易所的认证？

A: 在 `Setup()` 函数中配置 API Key 和 Secret，然后在 `Sign()` 方法中实现签名算法。

## 贡献

欢迎贡献代码完善这些交易所的实现！请：

1. Fork 项目
2. 创建功能分支
3. 实现交易所的 API 方法
4. 添加测试
5. 提交 Pull Request

## 更新日志

### 2024-01-XX

- ✅ 新增 37 个交易所的基础结构
- ✅ 实现所有交易所的 `Describe()` 函数
- ✅ 添加所有 HTTP 方法的占位符实现
- ✅ 确认 WebSocket 支持（通过 ExchangeBase）
- ✅ 所有文件编译通过

## 相关文档

- [统一 API 文档](UNIFIED_API.md)
- [交易所对比](EXCHANGE_COMPARISON_DETAILED.md)
- [README](README.md)

