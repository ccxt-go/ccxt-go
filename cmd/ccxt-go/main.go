package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ccxt-go/ccxt-go/pkg/ccxt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	exchangeId string
	symbol     string
	apiKey     string
	secret     string
	sandbox    bool
	verbose    bool
	output     string
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "ccxt-go",
	Short: "CCXT-Go: 加密货币交易所统一API",
	Long: `CCXT-Go 是一个用Go语言实现的加密货币交易所统一API库。
支持100+个交易所，提供统一的接口进行交易、查询市场数据等操作。

示例:
  ccxt-go markets --exchange binance
  ccxt-go ticker --exchange binance --symbol BTC/USDT
  ccxt-go balance --exchange binance --api-key YOUR_KEY --secret YOUR_SECRET`,
	Version: "1.0.0",
}

// marketsCmd 获取市场信息命令
var marketsCmd = &cobra.Command{
	Use:   "markets [exchange]",
	Short: "获取交易所的市场信息",
	Long:  `获取指定交易所的所有交易对信息`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if exchangeId == "" && len(args) > 0 {
			exchangeId = args[0]
		}

		if exchangeId == "" {
			fmt.Println("错误: 请指定交易所ID")
			cmd.Help()
			return
		}

		// 创建交易所实例
		exchange := createExchange(exchangeId)
		if exchange == nil {
			fmt.Printf("错误: 不支持的交易所: %s\n", exchangeId)
			return
		}

		// 获取市场信息
		markets := exchange.LoadMarkets()
		if markets.Type == ccxt.Error {
			fmt.Printf("错误: 获取市场信息失败: %s\n", markets.ToStr())
			return
		}

		// 输出结果
		if output == "json" {
			jsonStr, err := ccxt.JSONUtils.ToPrettyJSON(markets.ToMap())
			if err != nil {
				fmt.Printf("错误: JSON序列化失败: %v\n", err)
				return
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Printf("交易所: %s\n", exchangeId)
			fmt.Printf("交易对数量: %d\n", markets.Length.ToInt())
			fmt.Println("支持的交易对:")

			if markets.Type == ccxt.Map {
				symbols := markets.At(ccxt.MkString("symbols"))
				if symbols.Type == ccxt.Array {
					for i := 0; i < symbols.Length.ToInt(); i++ {
						symbol := symbols.At(ccxt.MkInteger(int64(i)))
						fmt.Printf("  %s\n", symbol.ToStr())
					}
				}
			}
		}
	},
}

// tickerCmd 获取价格信息命令
var tickerCmd = &cobra.Command{
	Use:   "ticker [symbol]",
	Short: "获取交易对的价格信息",
	Long:  `获取指定交易对的24小时价格统计信息`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if symbol == "" && len(args) > 0 {
			symbol = args[0]
		}

		if exchangeId == "" {
			fmt.Println("错误: 请指定交易所ID")
			cmd.Help()
			return
		}

		if symbol == "" {
			fmt.Println("错误: 请指定交易对")
			cmd.Help()
			return
		}

		// 创建交易所实例
		exchange := createExchange(exchangeId)
		if exchange == nil {
			fmt.Printf("错误: 不支持的交易所: %s\n", exchangeId)
			return
		}

		// 获取价格信息
		ticker := exchange.FetchTicker(ccxt.MkString(symbol))
		if ticker.Type == ccxt.Error {
			fmt.Printf("错误: 获取价格信息失败: %s\n", ticker.ToStr())
			return
		}

		// 输出结果
		if output == "json" {
			jsonStr, err := ccxt.JSONUtils.ToPrettyJSON(ticker.ToMap())
			if err != nil {
				fmt.Printf("错误: JSON序列化失败: %v\n", err)
				return
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Printf("交易所: %s\n", exchangeId)
			fmt.Printf("交易对: %s\n", symbol)
			if ticker.Type == ccxt.Map {
				last := ticker.At(ccxt.MkString("last"))
				bid := ticker.At(ccxt.MkString("bid"))
				ask := ticker.At(ccxt.MkString("ask"))
				high := ticker.At(ccxt.MkString("high"))
				low := ticker.At(ccxt.MkString("low"))
				volume := ticker.At(ccxt.MkString("baseVolume"))

				fmt.Printf("最新价格: %s\n", last.ToStr())
				fmt.Printf("买一价: %s\n", bid.ToStr())
				fmt.Printf("卖一价: %s\n", ask.ToStr())
				fmt.Printf("24h最高: %s\n", high.ToStr())
				fmt.Printf("24h最低: %s\n", low.ToStr())
				fmt.Printf("24h成交量: %s\n", volume.ToStr())
			}
		}
	},
}

// orderbookCmd 获取订单簿命令
var orderbookCmd = &cobra.Command{
	Use:   "orderbook [symbol]",
	Short: "获取交易对的订单簿",
	Long:  `获取指定交易对的订单簿信息`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if symbol == "" && len(args) > 0 {
			symbol = args[0]
		}

		if exchangeId == "" {
			fmt.Println("错误: 请指定交易所ID")
			cmd.Help()
			return
		}

		if symbol == "" {
			fmt.Println("错误: 请指定交易对")
			cmd.Help()
			return
		}

		// 创建交易所实例
		exchange := createExchange(exchangeId)
		if exchange == nil {
			fmt.Printf("错误: 不支持的交易所: %s\n", exchangeId)
			return
		}

		// 获取订单簿
		orderbook := exchange.FetchOrderBook(ccxt.MkString(symbol))
		if orderbook.Type == ccxt.Error {
			fmt.Printf("错误: 获取订单簿失败: %s\n", orderbook.ToStr())
			return
		}

		// 输出结果
		if output == "json" {
			jsonStr, err := ccxt.JSONUtils.ToPrettyJSON(orderbook.ToMap())
			if err != nil {
				fmt.Printf("错误: JSON序列化失败: %v\n", err)
				return
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Printf("交易所: %s\n", exchangeId)
			fmt.Printf("交易对: %s\n", symbol)

			if orderbook.Type == ccxt.Map {
				bids := orderbook.At(ccxt.MkString("bids"))
				asks := orderbook.At(ccxt.MkString("asks"))

				fmt.Println("\n买单 (Bids):")
				if bids.Type == ccxt.Array {
					for i := 0; i < bids.Length.ToInt() && i < 5; i++ {
						bid := bids.At(ccxt.MkInteger(int64(i)))
						if bid.Type == ccxt.Array {
							price := bid.At(ccxt.MkInteger(0))
							amount := bid.At(ccxt.MkInteger(1))
							fmt.Printf("  %s x %s\n", price.ToStr(), amount.ToStr())
						}
					}
				}

				fmt.Println("\n卖单 (Asks):")
				if asks.Type == ccxt.Array {
					for i := 0; i < asks.Length.ToInt() && i < 5; i++ {
						ask := asks.At(ccxt.MkInteger(int64(i)))
						if ask.Type == ccxt.Array {
							price := ask.At(ccxt.MkInteger(0))
							amount := ask.At(ccxt.MkInteger(1))
							fmt.Printf("  %s x %s\n", price.ToStr(), amount.ToStr())
						}
					}
				}
			}
		}
	},
}

// balanceCmd 获取余额命令
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "获取账户余额",
	Long:  `获取指定交易所的账户余额信息`,
	Run: func(cmd *cobra.Command, args []string) {
		if exchangeId == "" {
			fmt.Println("错误: 请指定交易所ID")
			cmd.Help()
			return
		}

		if apiKey == "" || secret == "" {
			fmt.Println("错误: 请提供API Key和Secret")
			cmd.Help()
			return
		}

		// 创建交易所实例
		exchange := createExchange(exchangeId)
		if exchange == nil {
			fmt.Printf("错误: 不支持的交易所: %s\n", exchangeId)
			return
		}

		// 设置API凭据
		exchange.At(ccxt.MkString("apiKey")).Value = apiKey
		exchange.At(ccxt.MkString("secret")).Value = secret

		// 获取余额
		balance := exchange.FetchBalance()
		if balance.Type == ccxt.Error {
			fmt.Printf("错误: 获取余额失败: %s\n", balance.ToStr())
			return
		}

		// 输出结果
		if output == "json" {
			jsonStr, err := ccxt.JSONUtils.ToPrettyJSON(balance.ToMap())
			if err != nil {
				fmt.Printf("错误: JSON序列化失败: %v\n", err)
				return
			}
			fmt.Println(jsonStr)
		} else {
			fmt.Printf("交易所: %s\n", exchangeId)
			fmt.Println("账户余额:")

			if balance.Type == ccxt.Map {
				// 遍历余额信息
				balanceMap := balance.ToMap()
				for currency, account := range balanceMap {
					if account.Type == ccxt.Map {
						free := account.At(ccxt.MkString("free"))
						used := account.At(ccxt.MkString("used"))
						total := account.At(ccxt.MkString("total"))

						if total.ToFloat() > 0 {
							fmt.Printf("  %s: 可用=%s, 冻结=%s, 总计=%s\n",
								currency, free.ToStr(), used.ToStr(), total.ToStr())
						}
					}
				}
			}
		}
	},
}

// exchangesCmd 列出支持的交易所
var exchangesCmd = &cobra.Command{
	Use:   "exchanges",
	Short: "列出所有支持的交易所",
	Long:  `列出CCXT-Go支持的所有交易所`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CCXT-Go 支持的交易所:")
		fmt.Println("===================")

		// 这里应该从实际的交易所列表中获取
		exchanges := []string{
			"binance", "okx", "coinbase", "kraken", "huobi", "kucoin",
			"gateio", "bybit", "bitfinex", "bitmex", "deribit", "ftx",
			"poloniex", "bittrex", "bitstamp", "gemini", "coinbasepro",
			"binanceus", "binanceusdm", "binancecoinm", "okex", "okex3",
			"okex5", "huobipro", "huobijp", "upbit", "bithumb", "coinone",
			"krakenfutures", "kucoinfutures", "phemex", "ascendex",
			"bequant", "bigone", "bitbank", "bitflyer", "bitget",
			"bitmart", "bitso", "bitvavo", "btcalpha", "btcbox",
			"btcmarkets", "btcturk", "cex", "coinex", "coinmate",
			"coinspot", "digifinex", "exmo", "gate", "hitbtc",
			"hollaex", "independentreserve", "indodax", "latoken",
			"lbank", "luno", "mercado", "ndax", "novadax", "oceanex",
			"paymium", "probit", "timex", "whitebit", "yobit", "zaif",
		}

		for i, exchange := range exchanges {
			fmt.Printf("%-20s", exchange)
			if (i+1)%4 == 0 {
				fmt.Println()
			}
		}
		if len(exchanges)%4 != 0 {
			fmt.Println()
		}

		fmt.Printf("\n总计: %d 个交易所\n", len(exchanges))
	},
}

// configCmd 配置管理命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  `管理CCXT-Go的配置`,
}

// configSetCmd 设置配置
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置配置项",
	Long:  `设置指定的配置项`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		configManager := ccxt.GetConfigManager()
		configManager.SetString(key, value)

		fmt.Printf("配置已设置: %s = %s\n", key, value)
	},
}

// configGetCmd 获取配置
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "获取配置项",
	Long:  `获取指定的配置项`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		configManager := ccxt.GetConfigManager()
		value := configManager.GetString(key)

		fmt.Printf("%s = %s\n", key, value)
	},
}

// configListCmd 列出所有配置
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Long:  `列出所有配置项`,
	Run: func(cmd *cobra.Command, args []string) {
		configManager := ccxt.GetConfigManager()
		globalConfig := configManager.GetGlobalConfig()

		fmt.Println("全局配置:")
		fmt.Printf("  默认超时: %d ms\n", globalConfig.DefaultTimeout)
		fmt.Printf("  默认速率限制: %d req/min\n", globalConfig.DefaultRateLimit)
		fmt.Printf("  启用日志: %t\n", globalConfig.EnableLogging)
		fmt.Printf("  日志级别: %s\n", globalConfig.LogLevel)
		fmt.Printf("  日志文件: %s\n", globalConfig.LogFile)
		fmt.Printf("  启用指标: %t\n", globalConfig.EnableMetrics)
		fmt.Printf("  指标端口: %d\n", globalConfig.MetricsPort)

		fmt.Println("\n交易所配置:")
		for exchangeId, config := range globalConfig.Exchanges {
			fmt.Printf("  %s:\n", exchangeId)
			fmt.Printf("    沙盒模式: %t\n", config.Sandbox)
			fmt.Printf("    速率限制: %d req/min\n", config.RateLimit)
			fmt.Printf("    超时时间: %d ms\n", config.Timeout)
			fmt.Printf("    启用速率限制: %t\n", config.EnableRateLimit)
		}
	},
}

// createExchange 创建交易所实例
func createExchange(exchangeId string) ccxt.Exchange {
	switch strings.ToLower(exchangeId) {
	case "binance":
		exchange := &ccxt.Binance{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "okx":
		exchange := &ccxt.Okex{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "coinbase":
		exchange := &ccxt.Coinbase{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "kraken":
		exchange := &ccxt.Kraken{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "huobi":
		exchange := &ccxt.Huobi{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "kucoin":
		exchange := &ccxt.Kucoin{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "gateio":
		exchange := &ccxt.Gateio{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	case "bybit":
		exchange := &ccxt.Bybit{}
		exchange.ExchangeBase = &ccxt.ExchangeBase{}
		exchange.Setup(ccxt.MkMap(&ccxt.VarMap{}), exchange)
		return exchange
	default:
		return nil
	}
}

// init 初始化命令
func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: config.yaml)")
	rootCmd.PersistentFlags().StringVar(&exchangeId, "exchange", "", "交易所ID")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "详细输出")
	rootCmd.PersistentFlags().StringVar(&output, "output", "table", "输出格式 (table|json)")

	// 交易所相关标志
	marketsCmd.Flags().StringVar(&exchangeId, "exchange", "", "交易所ID")
	tickerCmd.Flags().StringVar(&exchangeId, "exchange", "", "交易所ID")
	tickerCmd.Flags().StringVar(&symbol, "symbol", "", "交易对符号")
	orderbookCmd.Flags().StringVar(&exchangeId, "exchange", "", "交易所ID")
	orderbookCmd.Flags().StringVar(&symbol, "symbol", "", "交易对符号")

	// 账户相关标志
	balanceCmd.Flags().StringVar(&exchangeId, "exchange", "", "交易所ID")
	balanceCmd.Flags().StringVar(&apiKey, "api-key", "", "API Key")
	balanceCmd.Flags().StringVar(&secret, "secret", "", "API Secret")
	balanceCmd.Flags().BoolVar(&sandbox, "sandbox", false, "使用沙盒环境")

	// 添加子命令
	rootCmd.AddCommand(marketsCmd)
	rootCmd.AddCommand(tickerCmd)
	rootCmd.AddCommand(orderbookCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(exchangesCmd)
	rootCmd.AddCommand(configCmd)

	// 配置子命令
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Println("使用配置文件:", viper.ConfigFileUsed())
		}
	}
}

// main 主函数
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
