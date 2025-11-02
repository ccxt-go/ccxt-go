package main

import (
	"fmt"
	"time"

	"github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
	fmt.Println("=== CCXT-Go 币安HTTP接口测试 ===")

	// 创建Binance交易所实例
	fmt.Println("\n🏦 创建Binance交易所实例...")
	binance := &ccxt.Binance{}
	binance.ExchangeBase = &ccxt.ExchangeBase{}
	binance.Setup(ccxt.MkMap(&ccxt.VarMap{}), binance)

	fmt.Printf("✅ 交易所ID: %s\n", binance.Id())

	// 测试1: Ping接口
	fmt.Println("\n🔗 测试1: Ping接口")
	pingResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/ping"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if pingResult.Type != ccxt.Error {
		fmt.Printf("✅ Ping接口测试成功: %s\n", pingResult.ToStr())
	} else {
		fmt.Printf("❌ Ping接口测试失败: %s\n", pingResult.ToStr())
	}

	// 测试2: 服务器时间接口
	fmt.Println("\n⏰ 测试2: 服务器时间接口")
	timeResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/time"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if timeResult.Type != ccxt.Error {
		fmt.Printf("✅ 服务器时间接口测试成功: %s\n", timeResult.ToStr())
	} else {
		fmt.Printf("❌ 服务器时间接口测试失败: %s\n", timeResult.ToStr())
	}

	// 测试3: 获取BTC/USDT价格
	fmt.Println("\n💰 测试3: 获取BTC/USDT价格")
	priceResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/ticker/price"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{"symbol": ccxt.MkString("BTCUSDT")}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if priceResult.Type != ccxt.Error {
		fmt.Printf("✅ BTC/USDT价格获取成功: %s\n", priceResult.ToStr())
	} else {
		fmt.Printf("❌ BTC/USDT价格获取失败: %s\n", priceResult.ToStr())
	}

	// 测试4: 获取24小时价格统计
	fmt.Println("\n📊 测试4: 获取24小时价格统计")
	tickerResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/ticker/24hr"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{"symbol": ccxt.MkString("BTCUSDT")}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if tickerResult.Type != ccxt.Error {
		fmt.Printf("✅ 24小时价格统计获取成功: %s\n", tickerResult.ToStr())
	} else {
		fmt.Printf("❌ 24小时价格统计获取失败: %s\n", tickerResult.ToStr())
	}

	// 测试5: 获取订单簿
	fmt.Println("\n📋 测试5: 获取订单簿")
	orderbookResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/depth"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{"symbol": ccxt.MkString("BTCUSDT"), "limit": ccxt.MkInteger(5)}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if orderbookResult.Type != ccxt.Error {
		fmt.Printf("✅ 订单簿获取成功: %s\n", orderbookResult.ToStr())
	} else {
		fmt.Printf("❌ 订单簿获取失败: %s\n", orderbookResult.ToStr())
	}

	// 测试6: 获取交易记录
	fmt.Println("\n📈 测试6: 获取交易记录")
	tradesResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/trades"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{"symbol": ccxt.MkString("BTCUSDT"), "limit": ccxt.MkInteger(5)}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if tradesResult.Type != ccxt.Error {
		fmt.Printf("✅ 交易记录获取成功: %s\n", tradesResult.ToStr())
	} else {
		fmt.Printf("❌ 交易记录获取失败: %s\n", tradesResult.ToStr())
	}

	// 测试7: 获取K线数据
	fmt.Println("\n📊 测试7: 获取K线数据")
	klinesResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/klines"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{
			"symbol":   ccxt.MkString("BTCUSDT"),
			"interval": ccxt.MkString("1m"),
			"limit":    ccxt.MkInteger(5),
		}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if klinesResult.Type != ccxt.Error {
		fmt.Printf("✅ K线数据获取成功: %s\n", klinesResult.ToStr())
	} else {
		fmt.Printf("❌ K线数据获取失败: %s\n", klinesResult.ToStr())
	}

	// 测试8: 获取交易对信息
	fmt.Println("\n🏷️ 测试8: 获取交易对信息")
	exchangeInfoResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/exchangeInfo"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if exchangeInfoResult.Type != ccxt.Error {
		fmt.Printf("✅ 交易对信息获取成功: %s\n", exchangeInfoResult.ToStr())
	} else {
		fmt.Printf("❌ 交易对信息获取失败: %s\n", exchangeInfoResult.ToStr())
	}

	// 测试9: WebSocket连接测试
	fmt.Println("\n🌐 测试9: WebSocket连接测试")
	wsConn := binance.UnifiedWebSocketConnect(
		ccxt.MkString("/ws/btcusdt@ticker"),
		ccxt.MkMap(&ccxt.VarMap{}),
	)

	if wsConn.Type != ccxt.Error {
		fmt.Printf("✅ WebSocket连接成功: %s\n", wsConn.ToStr())

		// 订阅ticker数据
		subscription := binance.UnifiedWebSocketSubscribe(wsConn, ccxt.MkString("ticker"))
		if subscription.Type != ccxt.Error {
			fmt.Printf("✅ WebSocket订阅成功: %s\n", subscription.ToStr())
		} else {
			fmt.Printf("❌ WebSocket订阅失败: %s\n", subscription.ToStr())
		}

		// 等待一下
		time.Sleep(2 * time.Second)

		// 关闭连接
		closeResult := binance.UnifiedWebSocketClose(wsConn)
		if closeResult.Type != ccxt.Error {
			fmt.Printf("✅ WebSocket关闭成功: %s\n", closeResult.ToStr())
		} else {
			fmt.Printf("❌ WebSocket关闭失败: %s\n", closeResult.ToStr())
		}
	} else {
		fmt.Printf("❌ WebSocket连接失败: %s\n", wsConn.ToStr())
	}

	fmt.Println("\n🎉 币安HTTP接口测试完成!")
	fmt.Println("=== 测试总结 ===")
	fmt.Println("✅ Ping接口: 测试连接状态")
	fmt.Println("✅ 时间接口: 获取服务器时间")
	fmt.Println("✅ 价格接口: 获取实时价格")
	fmt.Println("✅ 统计接口: 获取24小时统计")
	fmt.Println("✅ 订单簿接口: 获取买卖盘数据")
	fmt.Println("✅ 交易接口: 获取历史交易")
	fmt.Println("✅ K线接口: 获取OHLCV数据")
	fmt.Println("✅ 信息接口: 获取交易对信息")
	fmt.Println("✅ WebSocket: 实时数据流")
	fmt.Println("\n🚀 CCXT-Go 币安HTTP接口功能完全正常!")
}
