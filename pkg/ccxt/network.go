package ccxt

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NetworkManager 统一网络管理器
type NetworkManager struct {
	httpClient    *http.Client
	wsConnections map[string]*WebSocketConnection
	mu            sync.RWMutex
	rateLimiter   *RateLimiter
	retryConfig   *RetryConfig
}

// WebSocketConnection WebSocket连接封装
type WebSocketConnection struct {
	conn        *websocket.Conn
	url         string
	subscribers map[string]chan *Variant
	mu          sync.RWMutex
	reconnect   bool
	lastPing    time.Time
	ctx         context.Context
	cancel      context.CancelFunc
}

// RateLimiter 速率限制器
type RateLimiter struct {
	requests map[string][]time.Time
	limits   map[string]int
	mu       sync.Mutex
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int
	BaseDelay     time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
}

// RequestConfig HTTP请求配置
type RequestConfig struct {
	URL       string
	Method    string
	Headers   map[string]string
	Body      interface{}
	Timeout   time.Duration
	Retry     bool
	RateLimit string
	Proxy     string
	UserAgent string
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	URL           string
	Headers       map[string]string
	Reconnect     bool
	PingInterval  time.Duration
	PongTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	Subscriptions []string
}

// NewNetworkManager 创建网络管理器
func NewNetworkManager() *NetworkManager {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &NetworkManager{
		httpClient:    client,
		wsConnections: make(map[string]*WebSocketConnection),
		rateLimiter:   NewRateLimiter(),
		retryConfig: &RetryConfig{
			MaxRetries:    3,
			BaseDelay:     time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
		},
	}
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limits: map[string]int{
			"default": 100, // 每分钟100个请求
		},
	}
}

// HTTPRequest 统一的HTTP请求方法
func (nm *NetworkManager) HTTPRequest(config *RequestConfig) (*Variant, error) {
	// 速率限制检查
	if config.RateLimit != "" {
		if !nm.rateLimiter.Allow(config.RateLimit) {
			return nil, fmt.Errorf("rate limit exceeded for %s", config.RateLimit)
		}
	}

	var body io.Reader
	if config.Body != nil {
		switch v := config.Body.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case []byte:
			body = bytes.NewBuffer(v)
		case map[string]interface{}:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(jsonData)
		default:
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequest(config.Method, config.URL, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	if config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}

	// 设置代理
	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
		if err == nil {
			transport := nm.httpClient.Transport.(*http.Transport)
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	// 设置超时
	ctx := context.Background()
	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}
	req = req.WithContext(ctx)

	// 执行请求（带重试）
	return nm.executeWithRetry(req, config.Retry)
}

// executeWithRetry 带重试的请求执行
func (nm *NetworkManager) executeWithRetry(req *http.Request, retry bool) (*Variant, error) {
	var lastErr error

	for attempt := 0; attempt <= nm.retryConfig.MaxRetries; attempt++ {
		if attempt > 0 && retry {
			delay := nm.calculateDelay(attempt)
			time.Sleep(delay)
		}

		resp, err := nm.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if !retry || attempt == nm.retryConfig.MaxRetries {
				return nil, err
			}
			continue
		}

		defer resp.Body.Close()

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			if !retry || attempt == nm.retryConfig.MaxRetries {
				return nil, err
			}
			continue
		}

		// 检查HTTP状态码
		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			if !retry || attempt == nm.retryConfig.MaxRetries {
				return nil, lastErr
			}
			continue
		}

		// 解析JSON响应
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			// 如果不是JSON，返回字符串
			return MkString(string(body)), nil
		}

		return ItfToVariant(jsonData), nil
	}

	return nil, lastErr
}

// calculateDelay 计算重试延迟
func (nm *NetworkManager) calculateDelay(attempt int) time.Duration {
	delay := time.Duration(float64(nm.retryConfig.BaseDelay) *
		(nm.retryConfig.BackoffFactor * float64(attempt)))

	if delay > nm.retryConfig.MaxDelay {
		delay = nm.retryConfig.MaxDelay
	}

	return delay
}

// Allow 检查速率限制
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limit := rl.limits["default"]
	if customLimit, exists := rl.limits[key]; exists {
		limit = customLimit
	}

	// 清理过期请求
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < time.Minute {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// 检查是否超过限制
	if len(rl.requests[key]) >= limit {
		return false
	}

	// 添加当前请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// SetRateLimit 设置速率限制
func (rl *RateLimiter) SetRateLimit(key string, limit int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits[key] = limit
}

// WebSocketConnect 建立WebSocket连接
func (nm *NetworkManager) WebSocketConnect(config *WebSocketConfig) (*WebSocketConnection, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// 检查是否已存在连接
	if conn, exists := nm.wsConnections[config.URL]; exists {
		return conn, nil
	}

	// 建立WebSocket连接
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: false},
	}

	headers := http.Header{}
	for key, value := range config.Headers {
		headers.Set(key, value)
	}

	wsConn, _, err := dialer.Dial(config.URL, headers)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	wsConnection := &WebSocketConnection{
		conn:        wsConn,
		url:         config.URL,
		subscribers: make(map[string]chan *Variant),
		reconnect:   config.Reconnect,
		lastPing:    time.Now(),
		ctx:         ctx,
		cancel:      cancel,
	}

	// 设置Ping/Pong处理器
	wsConn.SetPingHandler(func(message string) error {
		wsConnection.lastPing = time.Now()
		return wsConn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	})

	wsConn.SetPongHandler(func(message string) error {
		wsConnection.lastPing = time.Now()
		return nil
	})

	// 启动消息处理协程
	go wsConnection.handleMessages()

	// 启动心跳协程
	if config.PingInterval > 0 {
		go wsConnection.startHeartbeat(config.PingInterval)
	}

	nm.wsConnections[config.URL] = wsConnection
	return wsConnection, nil
}

// handleMessages 处理WebSocket消息
func (ws *WebSocketConnection) handleMessages() {
	defer ws.Close()

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				if ws.reconnect {
					// 尝试重连
					go ws.reconnectWebSocket()
				}
				return
			}

			// 解析消息
			var data interface{}
			if err := json.Unmarshal(message, &data); err != nil {
				// 如果不是JSON，作为字符串处理
				data = string(message)
			}

			variant := ItfToVariant(data)

			// 广播给所有订阅者
			ws.mu.RLock()
			for _, ch := range ws.subscribers {
				select {
				case ch <- variant:
				default:
					// 如果通道满了，跳过
				}
			}
			ws.mu.RUnlock()
		}
	}
}

// Subscribe 订阅WebSocket消息
func (ws *WebSocketConnection) Subscribe(topic string) <-chan *Variant {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ch := make(chan *Variant, 100) // 缓冲通道
	ws.subscribers[topic] = ch
	return ch
}

// Unsubscribe 取消订阅
func (ws *WebSocketConnection) Unsubscribe(topic string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ch, exists := ws.subscribers[topic]; exists {
		close(ch)
		delete(ws.subscribers, topic)
	}
}

// Send 发送WebSocket消息
func (ws *WebSocketConnection) Send(message interface{}) error {
	var data []byte
	var err error

	switch v := message.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

// Close 关闭WebSocket连接
func (ws *WebSocketConnection) Close() error {
	ws.cancel()

	// 关闭所有订阅通道
	ws.mu.Lock()
	for _, ch := range ws.subscribers {
		close(ch)
	}
	ws.subscribers = make(map[string]chan *Variant)
	ws.mu.Unlock()

	return ws.conn.Close()
}

// startHeartbeat 启动心跳
func (ws *WebSocketConnection) startHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-ticker.C:
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if ws.reconnect {
					go ws.reconnectWebSocket()
				}
				return
			}
		}
	}
}

// reconnectWebSocket 重连WebSocket
func (ws *WebSocketConnection) reconnectWebSocket() {
	// 实现重连逻辑
	time.Sleep(5 * time.Second) // 等待5秒后重连

	// 这里可以重新建立连接
	// 为了简化，这里只是示例
}

// CloseAll 关闭所有WebSocket连接
func (nm *NetworkManager) CloseAll() {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for _, conn := range nm.wsConnections {
		conn.Close()
	}
	nm.wsConnections = make(map[string]*WebSocketConnection)
}

// GetConnection 获取WebSocket连接
func (nm *NetworkManager) GetConnection(url string) (*WebSocketConnection, bool) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	conn, exists := nm.wsConnections[url]
	return conn, exists
}

// 全局网络管理器实例
var GlobalNetworkManager = NewNetworkManager()
