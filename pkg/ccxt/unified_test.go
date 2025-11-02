package ccxt

import (
	"fmt"
	"testing"
	"time"
)

// TestUnifiedHTTPRequest 测试统一HTTP请求
func TestUnifiedHTTPRequest(t *testing.T) {
	// 创建Binance交易所实例
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	// 测试ping请求
	result := binance.UnifiedHTTPRequest(
		MkString("/api/v3/ping"),
		MkString("public"),
		MkString("GET"),
		MkMap(&VarMap{}),
		MkMap(&VarMap{}),
		MkUndefined(),
	)

	if result.Type == Error {
		t.Errorf("HTTP请求失败: %s", result.ToStr())
	}
}

// TestUnifiedWebSocketConnect 测试统一WebSocket连接
func TestUnifiedWebSocketConnect(t *testing.T) {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	// 测试WebSocket连接
	conn := binance.UnifiedWebSocketConnect(
		MkString("/ws/btcusdt@ticker"),
		MkMap(&VarMap{}),
	)

	if conn.Type == Error {
		t.Errorf("WebSocket连接失败: %s", conn.ToStr())
		return
	}

	// 测试订阅
	subscription := binance.UnifiedWebSocketSubscribe(
		conn,
		MkString("ticker"),
	)

	if subscription.Type == Error {
		t.Errorf("WebSocket订阅失败: %s", subscription.ToStr())
	}

	// 等待一段时间接收消息
	time.Sleep(2 * time.Second)

	// 关闭连接
	closeResult := binance.UnifiedWebSocketClose(conn)
	if closeResult.Type == Error {
		t.Errorf("WebSocket关闭失败: %s", closeResult.ToStr())
	}
}

// TestRateLimiter 测试速率限制器
func TestRateLimiter(t *testing.T) {
	rateLimiter := NewRateLimiter()
	rateLimiter.SetRateLimit("test", 2) // 每分钟2个请求

	// 前两个请求应该成功
	if !rateLimiter.Allow("test") {
		t.Error("第一个请求应该被允许")
	}
	if !rateLimiter.Allow("test") {
		t.Error("第二个请求应该被允许")
	}

	// 第三个请求应该被限制
	if rateLimiter.Allow("test") {
		t.Error("第三个请求应该被限制")
	}
}

// TestNetworkManager 测试网络管理器
func TestNetworkManager(t *testing.T) {
	nm := NewNetworkManager()

	// 测试HTTP请求配置
	config := &RequestConfig{
		URL:     "https://httpbin.org/get",
		Method:  "GET",
		Headers: map[string]string{"User-Agent": "ccxt-go-test"},
		Timeout: 10 * time.Second,
		Retry:   false,
	}

	result, err := nm.HTTPRequest(config)
	if err != nil {
		t.Errorf("HTTP请求失败: %v", err)
		return
	}

	if result.Type == Error {
		t.Errorf("HTTP请求返回错误: %s", result.ToStr())
	}
}

// TestWebSocketConnection 测试WebSocket连接
func TestWebSocketConnection(t *testing.T) {
	nm := NewNetworkManager()

	// 测试WebSocket连接配置
	config := &WebSocketConfig{
		URL:          "wss://echo.websocket.org",
		Headers:      map[string]string{"User-Agent": "ccxt-go-test"},
		Reconnect:    false,
		PingInterval: 30 * time.Second,
		PongTimeout:  10 * time.Second,
	}

	conn, err := nm.WebSocketConnect(config)
	if err != nil {
		t.Errorf("WebSocket连接失败: %v", err)
		return
	}

	// 测试发送消息
	err = conn.Send("test message")
	if err != nil {
		t.Errorf("发送WebSocket消息失败: %v", err)
	}

	// 测试订阅
	ch := conn.Subscribe("test")
	if ch == nil {
		t.Error("订阅失败")
	}

	// 等待一段时间
	time.Sleep(1 * time.Second)

	// 取消订阅
	conn.Unsubscribe("test")

	// 关闭连接
	err = conn.Close()
	if err != nil {
		t.Errorf("关闭WebSocket连接失败: %v", err)
	}
}

// TestConcurrentRequests 测试并发请求
func TestConcurrentRequests(t *testing.T) {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	// 并发发送多个请求
	done := make(chan bool, 5)
	errors := make(chan error, 5)

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

			if result.Type == Error {
				errors <- fmt.Errorf("协程 %d 请求失败: %s", index, result.ToStr())
			} else {
				done <- true
			}
		}(i)
	}

	// 等待所有协程完成
	successCount := 0
	errorCount := 0
	for i := 0; i < 5; i++ {
		select {
		case <-done:
			successCount++
		case err := <-errors:
			t.Logf("请求错误: %v", err)
			errorCount++
		case <-time.After(10 * time.Second):
			t.Error("请求超时")
			return
		}
	}

	t.Logf("成功请求: %d, 失败请求: %d", successCount, errorCount)
}

// BenchmarkHTTPRequest HTTP请求性能测试
func BenchmarkHTTPRequest(b *testing.B) {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := binance.UnifiedHTTPRequest(
			MkString("/api/v3/ping"),
			MkString("public"),
			MkString("GET"),
			MkMap(&VarMap{}),
			MkMap(&VarMap{}),
			MkUndefined(),
		)

		if result.Type == Error {
			b.Errorf("请求失败: %s", result.ToStr())
		}
	}
}

// BenchmarkWebSocketConnection WebSocket连接性能测试
func BenchmarkWebSocketConnection(b *testing.B) {
	binance := &Binance{}
	binance.ExchangeBase = &ExchangeBase{}
	binance.Setup(MkMap(&VarMap{}), binance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn := binance.UnifiedWebSocketConnect(
			MkString("/ws/btcusdt@ticker"),
			MkMap(&VarMap{}),
		)

		if conn.Type != Error {
			binance.UnifiedWebSocketClose(conn)
		}
	}
}
