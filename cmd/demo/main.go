package main

import (
	"fmt"
	"time"

	"github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
	fmt.Println("=== CCXT-Go 统一调用接口演示 ===")

	// 创建Binance交易所实例
	binance := &ccxt.Binance{}
	binance.ExchangeBase = &ccxt.ExchangeBase{}
	binance.Setup(ccxt.MkMap(&ccxt.VarMap{}), binance)

	// 演示1: HTTP请求
	fmt.Println("\n1. HTTP请求演示")
	result := binance.UnifiedHTTPRequest(
		ccxt.MkString("/ping"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if result.Type != ccxt.Error {
		fmt.Printf("✅ HTTP请求成功: %s\n", result.ToStr())
	} else {
		fmt.Printf("❌ HTTP请求失败: %s\n", result.ToStr())
	}

	// 演示2: WebSocket连接
	fmt.Println("\n2. WebSocket连接演示")
	wsConn := binance.UnifiedWebSocketConnect(
		ccxt.MkString("/ws/btcusdt@ticker"),
		ccxt.MkMap(&ccxt.VarMap{}),
	)

	if wsConn.Type != ccxt.Error {
		fmt.Printf("✅ WebSocket连接成功: %s\n", wsConn.ToStr())

		// 订阅消息
		subscription := binance.UnifiedWebSocketSubscribe(
			wsConn,
			ccxt.MkString("ticker"),
		)
		fmt.Printf("✅ 订阅状态: %s\n", subscription.ToStr())

		// 等待一段时间
		time.Sleep(2 * time.Second)

		// 关闭连接
		closeResult := binance.UnifiedWebSocketClose(wsConn)
		fmt.Printf("✅ 关闭连接: %s\n", closeResult.ToStr())
	} else {
		fmt.Printf("❌ WebSocket连接失败: %s\n", wsConn.ToStr())
	}

	// 演示3: 速率限制
	fmt.Println("\n3. 速率限制演示")
	rateLimiter := ccxt.NewRateLimiter()
	rateLimiter.SetRateLimit("demo", 2) // 每分钟2个请求

	for i := 0; i < 4; i++ {
		if rateLimiter.Allow("demo") {
			fmt.Printf("✅ 请求 %d: 允许\n", i+1)
		} else {
			fmt.Printf("❌ 请求 %d: 被限制\n", i+1)
		}
	}

	// 演示4: 网络管理器
	fmt.Println("\n4. 网络管理器演示")
	nm := ccxt.NewNetworkManager()
	fmt.Printf("✅ 网络管理器创建成功\n")

	// 关闭所有连接
	nm.CloseAll()
	fmt.Printf("✅ 所有连接已关闭\n")

	fmt.Println("\n=== 演示完成 ===")
}
