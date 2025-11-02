package ccxt

import (
	"time"
)

// UnifiedClient 统一客户端接口
type UnifiedClient struct {
	exchange *ExchangeBase
	network  *NetworkManager
}

// NewUnifiedClient 创建统一客户端
func NewUnifiedClient(exchange *ExchangeBase) *UnifiedClient {
	return &UnifiedClient{
		exchange: exchange,
		network:  GlobalNetworkManager,
	}
}

// HTTPRequest 统一的HTTP请求方法
func (uc *UnifiedClient) HTTPRequest(path *Variant, apiType *Variant, method *Variant, params *Variant, headers *Variant, body *Variant) *Variant {
	// 构建完整URL
	baseURL := uc.getBaseURL(apiType)
	fullURL := baseURL + path.ToStr()

	// 构建请求配置
	config := &RequestConfig{
		URL:     fullURL,
		Method:  method.ToStr(),
		Headers: uc.buildHeaders(headers),
		Body:    uc.buildBody(body),
		Timeout: time.Duration((*uc.exchange.At(MkString("timeout"))).ToInt()) * time.Millisecond,
		Retry:   true,
	}

	// 添加查询参数
	if params.Type == Map {
		queryParams := uc.buildQueryParams(params)
		if queryParams != "" {
			config.URL += "?" + queryParams
		}
	}

	// 执行请求
	result, err := uc.network.HTTPRequest(config)
	if err != nil {
		return uc.handleRequestError(err)
	}

	return result
}

// WebSocketConnect 建立WebSocket连接
func (uc *UnifiedClient) WebSocketConnect(path *Variant, params *Variant) *Variant {
	// 构建WebSocket URL
	wsURL := uc.getWebSocketURL(path)

	// 构建WebSocket配置
	config := &WebSocketConfig{
		URL:          wsURL,
		Headers:      uc.buildWSHeaders(),
		Reconnect:    true,
		PingInterval: 30 * time.Second,
		PongTimeout:  10 * time.Second,
	}

	// 建立连接
	_, err := uc.network.WebSocketConnect(config)
	if err != nil {
		return uc.handleRequestError(err)
	}

	// 返回连接标识符
	return MkString(wsURL)
}

// WebSocketSubscribe 订阅WebSocket消息
func (uc *UnifiedClient) WebSocketSubscribe(connectionID *Variant, topic *Variant) *Variant {
	conn, exists := uc.network.GetConnection(connectionID.ToStr())
	if !exists {
		return NewNetworkError(MkString("WebSocket connection not found"))
	}

	// 订阅消息
	conn.Subscribe(topic.ToStr())

	// 返回订阅通道（这里简化处理）
	return MkString("subscribed")
}

// WebSocketSend 发送WebSocket消息
func (uc *UnifiedClient) WebSocketSend(connectionID *Variant, message *Variant) *Variant {
	conn, exists := uc.network.GetConnection(connectionID.ToStr())
	if !exists {
		return NewNetworkError(MkString("WebSocket connection not found"))
	}

	// 发送消息
	err := conn.Send(VariantToItf(message))
	if err != nil {
		return uc.handleRequestError(err)
	}

	return MkBool(true)
}

// WebSocketClose 关闭WebSocket连接
func (uc *UnifiedClient) WebSocketClose(connectionID *Variant) *Variant {
	conn, exists := uc.network.GetConnection(connectionID.ToStr())
	if !exists {
		return NewNetworkError(MkString("WebSocket connection not found"))
	}

	err := conn.Close()
	if err != nil {
		return uc.handleRequestError(err)
	}

	return MkBool(true)
}

// getBaseURL 获取基础URL
func (uc *UnifiedClient) getBaseURL(apiType *Variant) string {
	// 获取交易所ID，用于判断交易所类型
	exchangeID := ""
	if idPtr := uc.exchange.At(MkString("id")); idPtr != nil {
		id := *idPtr
		if id.Type != Error {
			exchangeID = id.ToStr()
		}
	}

	if apiType.Type == Undefined {
		apiType = MkString("public")
	}

	apiTypeStr := apiType.ToStr()

	// 针对Binance交易所，直接使用API URL
	switch exchangeID {
	case "binance":
		// Binance 现货 API
		if apiTypeStr == "public" {
			return "https://api.binance.com"
		}
	case "binanceusdm":
		// Binance USD-M 合约 API
		if apiTypeStr == "fapiPublic" || apiTypeStr == "public" {
			return "https://fapi.binance.com"
		}
	case "binancecoinm":
		// Binance COIN-M 合约 API
		if apiTypeStr == "dapiPublic" || apiTypeStr == "public" {
			return "https://dapi.binance.com"
		}
	}

	// 其他交易所：尝试从配置中获取
	urls := *uc.exchange.At(MkString("urls"))
	api := *urls.At(MkString("api"))

	// 检查api是否为map类型
	if api.Type != Map {
		return ""
	}

	apiMap := api.ToMap()
	if apiURL, exists := (*apiMap)[apiTypeStr]; exists {
		return (*apiURL).ToStr()
	}

	// 默认使用public API
	if publicAPI, exists := (*apiMap)["public"]; exists {
		return (*publicAPI).ToStr()
	}

	return ""
}

// getWebSocketURL 获取WebSocket URL
func (uc *UnifiedClient) getWebSocketURL(path *Variant) string {
	// 获取交易所ID，用于判断交易所类型
	exchangeID := ""
	if idPtr := uc.exchange.At(MkString("id")); idPtr != nil {
		id := *idPtr
		if id.Type != Error {
			exchangeID = id.ToStr()
		}
	}

	// 针对Binance交易所，直接使用WebSocket URL
	switch exchangeID {
	case "binance":
		// Binance现货WebSocket: wss://stream.binance.com
		return "wss://stream.binance.com" + path.ToStr()
	case "binanceusdm":
		// Binance USD-M合约WebSocket: wss://fstream.binance.com
		return "wss://fstream.binance.com" + path.ToStr()
	case "binancecoinm":
		// Binance COIN-M合约WebSocket: wss://dstream.binance.com
		return "wss://dstream.binance.com" + path.ToStr()
	}

	// 其他交易所：尝试从配置中获取
	urls := *uc.exchange.At(MkString("urls"))
	api := *urls.At(MkString("api"))

	// 检查api是否为map类型
	if api.Type != Map {
		return ""
	}

	apiMap := api.ToMap()

	// 查找WebSocket API
	if wsAPI, exists := (*apiMap)["ws"]; exists {
		baseURL := (*wsAPI).ToStr()
		return baseURL + path.ToStr()
	}

	// 如果没有专门的WebSocket API，尝试使用public API
	if publicAPI, exists := (*apiMap)["public"]; exists {
		baseURL := (*publicAPI).ToStr()
		// 将http替换为ws
		if len(baseURL) > 4 && baseURL[:4] == "http" {
			baseURL = "ws" + baseURL[4:]
		}
		return baseURL + path.ToStr()
	}

	return ""
}

// buildHeaders 构建请求头
func (uc *UnifiedClient) buildHeaders(headers *Variant) map[string]string {
	result := make(map[string]string)

	// 添加默认请求头
	defaultHeaders := *uc.exchange.At(MkString("headers"))
	if defaultHeaders.Type == Map {
		for key, value := range *defaultHeaders.ToMap() {
			result[key] = (*value).ToStr()
		}
	}

	// 添加用户提供的请求头
	if headers.Type == Map {
		for key, value := range *headers.ToMap() {
			result[key] = (*value).ToStr()
		}
	}

	// 添加User-Agent
	if userAgent := uc.exchange.At(MkString("userAgent")); (*userAgent).Type != Undefined {
		result["User-Agent"] = (*userAgent).ToStr()
	}

	return result
}

// buildWSHeaders 构建WebSocket请求头
func (uc *UnifiedClient) buildWSHeaders() map[string]string {
	result := make(map[string]string)

	// 添加默认请求头
	defaultHeaders := *uc.exchange.At(MkString("headers"))
	if defaultHeaders.Type == Map {
		for key, value := range *defaultHeaders.ToMap() {
			result[key] = (*value).ToStr()
		}
	}

	// 添加User-Agent
	if userAgent := uc.exchange.At(MkString("userAgent")); (*userAgent).Type != Undefined {
		result["User-Agent"] = (*userAgent).ToStr()
	}

	return result
}

// buildBody 构建请求体
func (uc *UnifiedClient) buildBody(body *Variant) interface{} {
	if body.Type == Undefined {
		return nil
	}

	return VariantToItf(body)
}

// buildQueryParams 构建查询参数
func (uc *UnifiedClient) buildQueryParams(params *Variant) string {
	if params.Type != Map {
		return ""
	}

	var result string
	first := true

	for key, value := range *params.ToMap() {
		if !first {
			result += "&"
		}
		result += key + "=" + (*value).ToStr()
		first = false
	}

	return result
}

// handleRequestError 处理请求错误
func (uc *UnifiedClient) handleRequestError(err error) *Variant {
	errorMsg := err.Error()

	// 根据错误类型返回相应的异常
	if len(errorMsg) > 0 {
		return NewNetworkError(MkString(errorMsg))
	}

	return NewNetworkError(MkString("Unknown network error"))
}

// 扩展ExchangeBase以支持统一客户端
func (this *ExchangeBase) GetUnifiedClient() *UnifiedClient {
	return NewUnifiedClient(this)
}

// UnifiedHTTPRequest 统一的HTTP请求方法（添加到ExchangeBase）
func (this *ExchangeBase) UnifiedHTTPRequest(path *Variant, apiType *Variant, method *Variant, params *Variant, headers *Variant, body *Variant) *Variant {
	client := this.GetUnifiedClient()
	return client.HTTPRequest(path, apiType, method, params, headers, body)
}

// UnifiedWebSocketConnect 统一的WebSocket连接方法
func (this *ExchangeBase) UnifiedWebSocketConnect(path *Variant, params *Variant) *Variant {
	client := this.GetUnifiedClient()
	return client.WebSocketConnect(path, params)
}

// UnifiedWebSocketSubscribe 统一的WebSocket订阅方法
func (this *ExchangeBase) UnifiedWebSocketSubscribe(connectionID *Variant, topic *Variant) *Variant {
	client := this.GetUnifiedClient()
	return client.WebSocketSubscribe(connectionID, topic)
}

// UnifiedWebSocketSend 统一的WebSocket发送方法
func (this *ExchangeBase) UnifiedWebSocketSend(connectionID *Variant, message *Variant) *Variant {
	client := this.GetUnifiedClient()
	return client.WebSocketSend(connectionID, message)
}

// UnifiedWebSocketClose 统一的WebSocket关闭方法
func (this *ExchangeBase) UnifiedWebSocketClose(connectionID *Variant) *Variant {
	client := this.GetUnifiedClient()
	return client.WebSocketClose(connectionID)
}
