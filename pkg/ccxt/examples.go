package ccxt

import (
	"fmt"
	"time"
)

// ExampleUsage 展示统一调用接口的使用方法
func ExampleUsage() {
	// 创建Binance交易所实例
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	// 示例1: HTTP请求 - 获取市场信息
	fmt.Println("=== HTTP请求示例 ===")

	// 使用统一HTTP请求
	markets := binance.UnifiedHTTPRequest(
		MkString("/api/v3/exchangeInfo"),
		MkString("public"),
		MkString("GET"),
		MkMap(&VarMap{}),
		MkMap(&VarMap{}),
		MkUndefined(),
	)

	if markets.Type != Error {
		fmt.Printf("获取到 %d 个市场\n", markets.Length.ToInt())
	}

	// 示例2: WebSocket连接 - 实时价格订阅
	fmt.Println("\n=== WebSocket连接示例 ===")

	// 建立WebSocket连接
	wsConn := binance.UnifiedWebSocketConnect(
		MkString("/ws/btcusdt@ticker"),
		MkMap(&VarMap{}),
	)

	if wsConn.Type != Error {
		fmt.Printf("WebSocket连接成功: %s\n", wsConn.ToStr())

		// 订阅消息
		subscription := binance.UnifiedWebSocketSubscribe(
			wsConn,
			MkString("ticker"),
		)

		fmt.Printf("订阅状态: %s\n", subscription.ToStr())

		// 模拟接收消息
		time.Sleep(2 * time.Second)

		// 关闭连接
		closeResult := binance.UnifiedWebSocketClose(wsConn)
		fmt.Printf("关闭连接: %s\n", closeResult.ToStr())
	}

	// 示例3: 带认证的私有API请求
	fmt.Println("\n=== 私有API请求示例 ===")

	// 设置API密钥（实际使用中应该从配置文件读取）
	binance.Setup(MkMap(&VarMap{
		"apiKey": MkString("your_api_key"),
		"secret": MkString("your_secret"),
	}), binance)

	// 获取账户余额
	balance := binance.UnifiedHTTPRequest(
		MkString("/api/v3/account"),
		MkString("private"),
		MkString("GET"),
		MkMap(&VarMap{}),
		MkMap(&VarMap{}),
		MkUndefined(),
	)

	if balance.Type != Error {
		fmt.Println("账户余额获取成功")
	} else {
		fmt.Printf("账户余额获取失败: %s\n", balance.ToStr())
	}
}

// AdvancedWebSocketExample 高级WebSocket使用示例
func AdvancedWebSocketExample() {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	fmt.Println("=== 高级WebSocket示例 ===")

	// 建立多个WebSocket连接
	connections := []string{
		"/ws/btcusdt@ticker",
		"/ws/ethusdt@ticker",
		"/ws/btcusdt@depth",
	}

	for _, path := range connections {
		conn := binance.UnifiedWebSocketConnect(
			MkString(path),
			MkMap(&VarMap{}),
		)

		if conn.Type != Error {
			fmt.Printf("连接 %s 成功\n", path)

			// 订阅不同的主题
			topic := "ticker"
			if len(path) > 20 && path[len(path)-5:] == "depth" {
				topic = "depth"
			}

			binance.UnifiedWebSocketSubscribe(conn, MkString(topic))
		}
	}

	// 模拟运行一段时间
	time.Sleep(5 * time.Second)

	// 关闭所有连接
	GlobalNetworkManager.CloseAll()
	fmt.Println("所有WebSocket连接已关闭")
}

// RateLimitExample 速率限制示例
func RateLimitExample() {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	fmt.Println("=== 速率限制示例 ===")

	// 设置速率限制
	GlobalNetworkManager.rateLimiter.SetRateLimit("binance", 10) // 每分钟10个请求

	// 快速发送多个请求
	for i := 0; i < 15; i++ {
		result := binance.UnifiedHTTPRequest(
			MkString("/api/v3/ping"),
			MkString("public"),
			MkString("GET"),
			MkMap(&VarMap{}),
			MkMap(&VarMap{}),
			MkUndefined(),
		)

		if result.Type == Error {
			fmt.Printf("请求 %d 被限制: %s\n", i+1, result.ToStr())
		} else {
			fmt.Printf("请求 %d 成功\n", i+1)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// ErrorHandlingExample 错误处理示例
func ErrorHandlingExample() {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	fmt.Println("=== 错误处理示例 ===")

	// 测试各种错误情况

	// 1. 无效URL
	result1 := binance.UnifiedHTTPRequest(
		MkString("/invalid/endpoint"),
		MkString("public"),
		MkString("GET"),
		MkMap(&VarMap{}),
		MkMap(&VarMap{}),
		MkUndefined(),
	)

	if result1.Type == Error {
		fmt.Printf("无效URL错误: %s\n", result1.ToStr())
	}

	// 2. 无效的WebSocket连接
	wsConn := binance.UnifiedWebSocketConnect(
		MkString("/invalid/ws/path"),
		MkMap(&VarMap{}),
	)

	if wsConn.Type == Error {
		fmt.Printf("WebSocket连接错误: %s\n", wsConn.ToStr())
	}

	// 3. 向不存在的连接发送消息
	sendResult := binance.UnifiedWebSocketSend(
		MkString("nonexistent_connection"),
		MkString("test message"),
	)

	if sendResult.Type == Error {
		fmt.Printf("发送消息错误: %s\n", sendResult.ToStr())
	}
}

// ConcurrentExample 并发使用示例
func ConcurrentExample() {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	fmt.Println("=== 并发使用示例 ===")

	// 并发发送多个HTTP请求
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
			} else {
				fmt.Printf("协程 %d 请求失败: %s\n", index, result.ToStr())
			}

			done <- true
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < 5; i++ {
		<-done
	}

	fmt.Println("所有并发请求完成")
}
