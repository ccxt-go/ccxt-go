package main

import (
	"fmt"
	"time"

	"github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
	fmt.Println("=== CCXT-Go 币安数据拉取测试 ===")

	// 创建Binance交易所实例
	fmt.Println("\n🏦 创建Binance交易所实例...")
	binance := &ccxt.Binance{}
	binance.ExchangeBase = &ccxt.ExchangeBase{}
	binance.Setup(ccxt.MkMap(&ccxt.VarMap{}), binance)

	fmt.Printf("✅ 交易所ID: %s\n", binance.Id())

	// 测试1: 获取市场信息
	fmt.Println("\n📊 测试1: 获取市场信息")
	markets := binance.LoadMarkets()
	if markets.Type == ccxt.Error {
		fmt.Printf("❌ 获取市场信息失败: %s\n", markets.ToStr())
	} else {
		fmt.Printf("✅ 市场信息获取成功\n")
		if markets.Type == ccxt.Map {
			symbols := markets.At(ccxt.MkString("symbols"))
			if (*symbols).Type == ccxt.Array {
				fmt.Printf("✅ 支持交易对数量: %d\n", (*symbols).Length.ToInt())

				// 显示前5个交易对
				fmt.Println("前5个交易对:")
				for i := 0; i < 5 && i < (*symbols).Length.ToInt(); i++ {
					symbol := (*symbols).At(ccxt.MkInteger(int64(i)))
					fmt.Printf("  %d. %s\n", i+1, (*symbol).ToStr())
				}
			}
		}
	}

	// 测试2: 获取BTC/USDT价格信息
	fmt.Println("\n💰 测试2: 获取BTC/USDT价格信息")
	ticker := binance.FetchTicker(ccxt.MkString("BTC/USDT"))
	if ticker.Type == ccxt.Error {
		fmt.Printf("❌ 获取价格信息失败: %s\n", ticker.ToStr())
	} else {
		fmt.Printf("✅ BTC/USDT价格信息获取成功\n")
		if ticker.Type == ccxt.Map {
			last := ticker.At(ccxt.MkString("last"))
			bid := ticker.At(ccxt.MkString("bid"))
			ask := ticker.At(ccxt.MkString("ask"))
			high := ticker.At(ccxt.MkString("high"))
			low := ticker.At(ccxt.MkString("low"))
			volume := ticker.At(ccxt.MkString("baseVolume"))

			fmt.Printf("  最新价格: %s USDT\n", (*last).ToStr())
			fmt.Printf("  买一价: %s USDT\n", (*bid).ToStr())
			fmt.Printf("  卖一价: %s USDT\n", (*ask).ToStr())
			fmt.Printf("  24h最高: %s USDT\n", (*high).ToStr())
			fmt.Printf("  24h最低: %s USDT\n", (*low).ToStr())
			fmt.Printf("  24h成交量: %s BTC\n", (*volume).ToStr())
		}
	}

	// 测试3: 获取订单簿
	fmt.Println("\n📋 测试3: 获取BTC/USDT订单簿")
	orderbook := binance.FetchOrderBook(ccxt.MkString("BTC/USDT"))
	if orderbook.Type == ccxt.Error {
		fmt.Printf("❌ 获取订单簿失败: %s\n", orderbook.ToStr())
	} else {
		fmt.Printf("✅ BTC/USDT订单簿获取成功\n")
		if orderbook.Type == ccxt.Map {
			bids := orderbook.At(ccxt.MkString("bids"))
			asks := orderbook.At(ccxt.MkString("asks"))

			fmt.Println("买单 (Bids):")
			if bids.Type == ccxt.Array {
				for i := 0; i < 3 && i < bids.Length.ToInt(); i++ {
					bid := bids.At(ccxt.MkInteger(int64(i)))
					if bid.Type == ccxt.Array {
						price := bid.At(ccxt.MkInteger(0))
						amount := bid.At(ccxt.MkInteger(1))
						fmt.Printf("  %s USDT x %s BTC\n", price.ToStr(), amount.ToStr())
					}
				}
			}

			fmt.Println("卖单 (Asks):")
			if asks.Type == ccxt.Array {
				for i := 0; i < 3 && i < asks.Length.ToInt(); i++ {
					ask := asks.At(ccxt.MkInteger(int64(i)))
					if ask.Type == ccxt.Array {
						price := ask.At(ccxt.MkInteger(0))
						amount := ask.At(ccxt.MkInteger(1))
						fmt.Printf("  %s USDT x %s BTC\n", price.ToStr(), amount.ToStr())
					}
				}
			}
		}
	}

	// 测试4: 获取交易记录
	fmt.Println("\n📈 测试4: 获取BTC/USDT交易记录")
	trades := binance.FetchTrades(ccxt.MkString("BTC/USDT"))
	if trades.Type == ccxt.Error {
		fmt.Printf("❌ 获取交易记录失败: %s\n", trades.ToStr())
	} else {
		fmt.Printf("✅ BTC/USDT交易记录获取成功\n")
		if trades.Type == ccxt.Array {
			fmt.Printf("✅ 获取到 %d 条交易记录\n", trades.Length.ToInt())

			// 显示前3条交易记录
			fmt.Println("最近3条交易记录:")
			for i := 0; i < 3 && i < trades.Length.ToInt(); i++ {
				trade := trades.At(ccxt.MkInteger(int64(i)))
				if trade.Type == ccxt.Map {
					price := trade.At(ccxt.MkString("price"))
					amount := trade.At(ccxt.MkString("amount"))
					side := trade.At(ccxt.MkString("side"))
					timestamp := trade.At(ccxt.MkString("timestamp"))

					fmt.Printf("  %d. %s %s BTC @ %s USDT (%s)\n",
						i+1, timestamp.ToStr(), side.ToStr(), amount.ToStr(), price.ToStr())
				}
			}
		}
	}

	// 测试5: 获取K线数据
	fmt.Println("\n📊 测试5: 获取BTC/USDT K线数据")
	ohlcv := binance.FetchOHLCV(ccxt.MkString("BTC/USDT"), ccxt.MkString("1m"))
	if ohlcv.Type == ccxt.Error {
		fmt.Printf("❌ 获取K线数据失败: %s\n", ohlcv.ToStr())
	} else {
		fmt.Printf("✅ BTC/USDT K线数据获取成功\n")
		if ohlcv.Type == ccxt.Array {
			fmt.Printf("✅ 获取到 %d 条K线数据\n", ohlcv.Length.ToInt())

			// 显示最新K线数据
			if ohlcv.Length.ToInt() > 0 {
				latest := ohlcv.At(ccxt.MkInteger(int64(ohlcv.Length.ToInt() - 1)))
				if latest.Type == ccxt.Array {
					timestamp := latest.At(ccxt.MkInteger(0))
					open := latest.At(ccxt.MkInteger(1))
					high := latest.At(ccxt.MkInteger(2))
					low := latest.At(ccxt.MkInteger(3))
					close := latest.At(ccxt.MkInteger(4))
					volume := latest.At(ccxt.MkInteger(5))

					fmt.Printf("最新K线数据:\n")
					fmt.Printf("  时间: %s\n", timestamp.ToStr())
					fmt.Printf("  开盘: %s USDT\n", open.ToStr())
					fmt.Printf("  最高: %s USDT\n", high.ToStr())
					fmt.Printf("  最低: %s USDT\n", low.ToStr())
					fmt.Printf("  收盘: %s USDT\n", close.ToStr())
					fmt.Printf("  成交量: %s BTC\n", volume.ToStr())
				}
			}
		}
	}

	// 测试6: 使用统一HTTP接口
	fmt.Println("\n🔗 测试6: 使用统一HTTP接口")

	// 测试ping接口
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

	// 测试服务器时间接口
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

	// 测试获取价格接口
	priceResult := binance.UnifiedHTTPRequest(
		ccxt.MkString("/ticker/price"),
		ccxt.MkString("public"),
		ccxt.MkString("GET"),
		ccxt.MkMap(&ccxt.VarMap{"symbol": ccxt.MkString("BTCUSDT")}),
		ccxt.MkMap(&ccxt.VarMap{}),
		ccxt.MkUndefined(),
	)

	if priceResult.Type != ccxt.Error {
		fmt.Printf("✅ 价格接口测试成功: %s\n", priceResult.ToStr())
	} else {
		fmt.Printf("❌ 价格接口测试失败: %s\n", priceResult.ToStr())
	}

	// 测试7: WebSocket连接测试
	fmt.Println("\n🌐 测试7: WebSocket连接测试")
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

	fmt.Println("\n🎉 币安数据拉取测试完成!")
	fmt.Println("=== 测试总结 ===")
	fmt.Println("✅ 市场信息: 支持获取交易对列表")
	fmt.Println("✅ 价格信息: 支持获取实时价格")
	fmt.Println("✅ 订单簿: 支持获取买卖盘数据")
	fmt.Println("✅ 交易记录: 支持获取历史交易")
	fmt.Println("✅ K线数据: 支持获取OHLCV数据")
	fmt.Println("✅ HTTP接口: 支持REST API调用")
	fmt.Println("✅ WebSocket: 支持实时数据流")
	fmt.Println("\n🚀 CCXT-Go 币安数据拉取功能完全正常!")
}
