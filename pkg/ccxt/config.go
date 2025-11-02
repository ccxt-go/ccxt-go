package ccxt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	mu       sync.RWMutex
	config   map[string]interface{}
	viper    *viper.Viper
	filePath string
}

// ExchangeConfig 交易所配置
type ExchangeConfig struct {
	APIKey          string            `json:"apiKey,omitempty"`
	Secret          string            `json:"secret,omitempty"`
	Password        string            `json:"password,omitempty"`
	Sandbox         bool              `json:"sandbox,omitempty"`
	RateLimit       int               `json:"rateLimit,omitempty"`
	Timeout         int               `json:"timeout,omitempty"`
	EnableRateLimit bool              `json:"enableRateLimit,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Proxy           string            `json:"proxy,omitempty"`
	UserAgent       string            `json:"userAgent,omitempty"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	DefaultTimeout   int                       `json:"defaultTimeout"`
	DefaultRateLimit int                       `json:"defaultRateLimit"`
	EnableLogging    bool                      `json:"enableLogging"`
	LogLevel         string                    `json:"logLevel"`
	LogFile          string                    `json:"logFile"`
	EnableMetrics    bool                      `json:"enableMetrics"`
	MetricsPort      int                       `json:"metricsPort"`
	Exchanges        map[string]ExchangeConfig `json:"exchanges"`
}

var (
	globalConfigManager *ConfigManager
	configOnce          sync.Once
)

// GetConfigManager 获取全局配置管理器
func GetConfigManager() *ConfigManager {
	configOnce.Do(func() {
		globalConfigManager = NewConfigManager()
	})
	return globalConfigManager
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{
		config:   make(map[string]interface{}),
		viper:    viper.New(),
		filePath: "config.yaml",
	}

	// 设置默认配置
	cm.setDefaultConfig()

	// 尝试加载配置文件
	cm.LoadConfig()

	return cm
}

// setDefaultConfig 设置默认配置
func (cm *ConfigManager) setDefaultConfig() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	defaultConfig := GlobalConfig{
		DefaultTimeout:   30000, // 30秒
		DefaultRateLimit: 1200,  // 每分钟1200次请求
		EnableLogging:    true,
		LogLevel:         "info",
		LogFile:          "ccxt-go.log",
		EnableMetrics:    false,
		MetricsPort:      9090,
		Exchanges:        make(map[string]ExchangeConfig),
	}

	cm.config["global"] = defaultConfig
}

// LoadConfig 加载配置文件
func (cm *ConfigManager) LoadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 设置配置文件路径
	cm.viper.SetConfigFile(cm.filePath)
	cm.viper.SetConfigType("yaml")

	// 设置默认值
	cm.viper.SetDefault("global.defaultTimeout", 30000)
	cm.viper.SetDefault("global.defaultRateLimit", 1200)
	cm.viper.SetDefault("global.enableLogging", true)
	cm.viper.SetDefault("global.logLevel", "info")
	cm.viper.SetDefault("global.logFile", "ccxt-go.log")
	cm.viper.SetDefault("global.enableMetrics", false)
	cm.viper.SetDefault("global.metricsPort", 9090)

	// 尝试读取配置文件
	if err := cm.viper.ReadInConfig(); err != nil {
		// 如果文件不存在，创建默认配置文件
		if os.IsNotExist(err) {
			return cm.createDefaultConfigFile()
		}
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置到内存
	var globalConfig GlobalConfig
	if err := cm.viper.Unmarshal(&globalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	cm.config["global"] = globalConfig
	return nil
}

// createDefaultConfigFile 创建默认配置文件
func (cm *ConfigManager) createDefaultConfigFile() error {
	defaultConfig := GlobalConfig{
		DefaultTimeout:   30000,
		DefaultRateLimit: 1200,
		EnableLogging:    true,
		LogLevel:         "info",
		LogFile:          "ccxt-go.log",
		EnableMetrics:    false,
		MetricsPort:      9090,
		Exchanges: map[string]ExchangeConfig{
			"binance": {
				Sandbox:         false,
				RateLimit:       1200,
				Timeout:         30000,
				EnableRateLimit: true,
				Headers: map[string]string{
					"User-Agent": "ccxt-go/1.0",
				},
			},
			"okx": {
				Sandbox:         false,
				RateLimit:       20,
				Timeout:         30000,
				EnableRateLimit: true,
				Headers: map[string]string{
					"User-Agent": "ccxt-go/1.0",
				},
			},
		},
	}

	// 确保目录存在
	dir := filepath.Dir(cm.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 写入配置文件
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	cm.config["global"] = defaultConfig
	return nil
}

// SaveConfig 保存配置到文件
func (cm *ConfigManager) SaveConfig() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	globalConfig, ok := cm.config["global"].(GlobalConfig)
	if !ok {
		return fmt.Errorf("无效的全局配置")
	}

	data, err := json.MarshalIndent(globalConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(cm.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// GetGlobalConfig 获取全局配置
func (cm *ConfigManager) GetGlobalConfig() GlobalConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if config, ok := cm.config["global"].(GlobalConfig); ok {
		return config
	}

	// 返回默认配置
	return GlobalConfig{
		DefaultTimeout:   30000,
		DefaultRateLimit: 1200,
		EnableLogging:    true,
		LogLevel:         "info",
		LogFile:          "ccxt-go.log",
		EnableMetrics:    false,
		MetricsPort:      9090,
		Exchanges:        make(map[string]ExchangeConfig),
	}
}

// GetExchangeConfig 获取交易所配置
func (cm *ConfigManager) GetExchangeConfig(exchangeId string) ExchangeConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	globalConfig := cm.GetGlobalConfig()
	if config, ok := globalConfig.Exchanges[exchangeId]; ok {
		return config
	}

	// 返回默认交易所配置
	return ExchangeConfig{
		Sandbox:         false,
		RateLimit:       globalConfig.DefaultRateLimit,
		Timeout:         globalConfig.DefaultTimeout,
		EnableRateLimit: true,
		Headers: map[string]string{
			"User-Agent": "ccxt-go/1.0",
		},
	}
}

// SetExchangeConfig 设置交易所配置
func (cm *ConfigManager) SetExchangeConfig(exchangeId string, config ExchangeConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	globalConfig := cm.GetGlobalConfig()
	globalConfig.Exchanges[exchangeId] = config
	cm.config["global"] = globalConfig

	return cm.SaveConfig()
}

// GetString 获取字符串配置
func (cm *ConfigManager) GetString(key string) string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.viper.GetString(key)
}

// GetInt 获取整数配置
func (cm *ConfigManager) GetInt(key string) int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.viper.GetInt(key)
}

// GetBool 获取布尔配置
func (cm *ConfigManager) GetBool(key string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.viper.GetBool(key)
}

// SetString 设置字符串配置
func (cm *ConfigManager) SetString(key, value string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.viper.Set(key, value)
}

// SetInt 设置整数配置
func (cm *ConfigManager) SetInt(key string, value int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.viper.Set(key, value)
}

// SetBool 设置布尔配置
func (cm *ConfigManager) SetBool(key string, value bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.viper.Set(key, value)
}

// WatchConfig 监听配置文件变化
func (cm *ConfigManager) WatchConfig() {
	cm.viper.WatchConfig()
	cm.viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已更新: %s\n", e.Name)
		cm.LoadConfig()
	})
}

// GetConfigPath 获取配置文件路径
func (cm *ConfigManager) GetConfigPath() string {
	return cm.filePath
}

// SetConfigPath 设置配置文件路径
func (cm *ConfigManager) SetConfigPath(path string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.filePath = path
	cm.viper.SetConfigFile(path)
}
