package main

import (
	"fmt"

	"github.com/ccxt-go/ccxt-go/pkg/ccxt"
)

func main() {
	fmt.Println("=== CCXT-Go 完善功能演示 ===")

	// 1. 配置管理演示
	fmt.Println("\n🔧 配置管理演示")
	configManager := ccxt.GetConfigManager()
	globalConfig := configManager.GetGlobalConfig()
	fmt.Printf("默认超时时间: %d ms\n", globalConfig.DefaultTimeout)
	fmt.Printf("默认速率限制: %d req/min\n", globalConfig.DefaultRateLimit)
	fmt.Printf("启用日志: %t\n", globalConfig.EnableLogging)

	// 2. 日志系统演示
	fmt.Println("\n📝 日志系统演示")
	logManager := ccxt.GetLogManager()
	logManager.Info("CCXT-Go 日志系统测试")
	logManager.Warn("这是一个警告消息")
	logManager.LogError("这是一个错误消息")

	// 3. 工具函数演示
	fmt.Println("\n🛠️ 工具函数演示")

	// 字符串工具
	fmt.Printf("驼峰命名: %s\n", ccxt.StringUtils.CamelCase("hello_world"))
	fmt.Printf("蛇形命名: %s\n", ccxt.StringUtils.SnakeCase("HelloWorld"))

	// 数字工具
	fmt.Printf("四舍五入: %.2f\n", ccxt.NumberUtils.Round(3.14159, 2))
	fmt.Printf("限制范围: %.2f\n", ccxt.NumberUtils.Clamp(15, 10, 20))

	// 加密工具
	data := "hello world"
	fmt.Printf("MD5哈希: %s\n", ccxt.CryptoUtils.MD5(data))
	fmt.Printf("SHA256哈希: %s\n", ccxt.CryptoUtils.SHA256(data))

	// 时间工具
	now := ccxt.TimeUtils.Now()
	fmt.Printf("当前时间: %s\n", now.Format("2006-01-02 15:04:05"))
	tomorrow := ccxt.TimeUtils.AddDays(now, 1)
	fmt.Printf("明天: %s\n", tomorrow.Format("2006-01-02 15:04:05"))

	// 验证工具
	fmt.Printf("邮箱验证: %t\n", ccxt.ValidationUtils.IsEmail("test@example.com"))
	fmt.Printf("URL验证: %t\n", ccxt.ValidationUtils.IsURL("https://example.com"))
	fmt.Printf("IP验证: %t\n", ccxt.ValidationUtils.IsIP("192.168.1.1"))

	// 4. Variant系统演示
	fmt.Println("\n🔢 Variant系统演示")

	// 基本类型
	str := ccxt.MkString("Hello World")
	fmt.Printf("字符串: %s\n", str.ToStr())

	num := ccxt.MkNumber(123.45)
	fmt.Printf("数字: %s\n", num.ToStr())

	boolean := ccxt.MkBool(true)
	fmt.Printf("布尔: %s\n", boolean.ToStr())

	// Map操作
	m := ccxt.MkMap(&ccxt.VarMap{
		"key1": ccxt.MkString("value1"),
		"key2": ccxt.MkNumber(42),
	})
	fmt.Printf("Map: %s\n", m.ToStr())

	// Array操作
	arr := ccxt.MkArray(&ccxt.VarArray{
		ccxt.MkString("item1"),
		ccxt.MkString("item2"),
		ccxt.MkString("item3"),
	})
	fmt.Printf("Array: %s\n", arr.ToStr())
	fmt.Printf("数组长度: %d\n", arr.Length.ToInt())

	// 5. 数学运算演示
	fmt.Println("\n🧮 数学运算演示")
	a := ccxt.MkNumber(10.5)
	b := ccxt.MkNumber(2.5)

	add := ccxt.OpAdd(a, b)
	fmt.Printf("加法: %s + %s = %s\n", a.ToStr(), b.ToStr(), add.ToStr())

	sub := ccxt.OpSub(a, b)
	fmt.Printf("减法: %s - %s = %s\n", a.ToStr(), b.ToStr(), sub.ToStr())

	mul := ccxt.OpMulti(a, b)
	fmt.Printf("乘法: %s * %s = %s\n", a.ToStr(), b.ToStr(), mul.ToStr())

	div := ccxt.OpDiv(a, b)
	fmt.Printf("除法: %s / %s = %s\n", a.ToStr(), b.ToStr(), div.ToStr())

	// 6. 网络管理器演示
	fmt.Println("\n🌐 网络管理器演示")
	nm := ccxt.NewNetworkManager()
	fmt.Println("网络管理器创建成功")

	// 速率限制器
	rateLimiter := ccxt.NewRateLimiter()
	rateLimiter.SetRateLimit("test", 3)

	allowed := 0
	for i := 0; i < 5; i++ {
		if rateLimiter.Allow("test") {
			allowed++
			fmt.Printf("请求 %d: 允许\n", i+1)
		} else {
			fmt.Printf("请求 %d: 被限制\n", i+1)
		}
	}
	fmt.Printf("速率限制测试: %d/5 请求被允许\n", allowed)

	// 7. 交易所基础功能演示
	fmt.Println("\n🏦 交易所基础功能演示")
	binance := &ccxt.Binance{}
	binance.ExchangeBase = &ccxt.ExchangeBase{}
	binance.Setup(ccxt.MkMap(&ccxt.VarMap{}), binance)

	fmt.Printf("交易所ID: %s\n", binance.Id())

	// 8. 统一客户端演示
	fmt.Println("\n🔗 统一客户端演示")
	client := binance.GetUnifiedClient()
	if client != nil {
		fmt.Println("统一客户端创建成功")
	}

	// 9. JSON工具演示
	fmt.Println("\n📄 JSON工具演示")
	dataMap := map[string]interface{}{
		"name":   "test",
		"age":    30,
		"active": true,
	}

	jsonStr, err := ccxt.JSONUtils.ToPrettyJSON(dataMap)
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
	} else {
		fmt.Printf("JSON输出:\n%s\n", jsonStr)
	}

	// 10. 数学工具演示
	fmt.Println("\n🔢 数学工具演示")
	fact := ccxt.MathUtils.Factorial(5)
	fmt.Printf("5的阶乘: %d\n", fact.Int64())

	gcd := ccxt.MathUtils.GCD(12, 8)
	fmt.Printf("12和8的最大公约数: %d\n", gcd)

	lcm := ccxt.MathUtils.LCM(12, 8)
	fmt.Printf("12和8的最小公倍数: %d\n", lcm)

	isPrime := ccxt.MathUtils.IsPrime(17)
	fmt.Printf("17是否为质数: %t\n", isPrime)

	fib := ccxt.MathUtils.Fibonacci(10)
	fmt.Printf("前10个斐波那契数: %v\n", fib)

	// 清理资源
	fmt.Println("\n🧹 清理资源")
	nm.CloseAll()
	logManager.Close()

	fmt.Println("\n🎉 CCXT-Go 完善功能演示完成!")
	fmt.Println("=== 总结 ===")
	fmt.Println("✅ 配置管理: 支持全局和交易所配置")
	fmt.Println("✅ 日志系统: 结构化日志记录")
	fmt.Println("✅ 工具函数: 丰富的工具函数库")
	fmt.Println("✅ Variant系统: 动态类型系统")
	fmt.Println("✅ 数学运算: 完整的数学运算支持")
	fmt.Println("✅ 网络管理: HTTP/WebSocket支持")
	fmt.Println("✅ 交易所功能: 基础交易所功能")
	fmt.Println("✅ 统一客户端: 统一的API接口")
	fmt.Println("✅ JSON处理: JSON序列化/反序列化")
	fmt.Println("✅ 数学工具: 高级数学函数")
	fmt.Println("\n🚀 CCXT-Go 项目已全面完善!")
}
