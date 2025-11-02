package ccxt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// StringUtils 字符串工具
type StringUtils struct{}

// IsEmpty 检查字符串是否为空
func (su *StringUtils) IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 检查字符串是否不为空
func (su *StringUtils) IsNotEmpty(s string) bool {
	return !su.IsEmpty(s)
}

// Trim 去除字符串两端空白
func (su *StringUtils) Trim(s string) string {
	return strings.TrimSpace(s)
}

// ToLower 转换为小写
func (su *StringUtils) ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper 转换为大写
func (su *StringUtils) ToUpper(s string) string {
	return strings.ToUpper(s)
}

// Contains 检查是否包含子字符串
func (su *StringUtils) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// StartsWith 检查是否以指定字符串开头
func (su *StringUtils) StartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// EndsWith 检查是否以指定字符串结尾
func (su *StringUtils) EndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Replace 替换字符串
func (su *StringUtils) Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Split 分割字符串
func (su *StringUtils) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Join 连接字符串
func (su *StringUtils) Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// PadLeft 左填充
func (su *StringUtils) PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return padding + s
}

// PadRight 右填充
func (su *StringUtils) PadRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return s + padding
}

// CamelCase 转换为驼峰命名
func (su *StringUtils) CamelCase(s string) string {
	if su.IsEmpty(s) {
		return s
	}

	words := strings.FieldsFunc(s, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	})

	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := words[i]
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	return result
}

// SnakeCase 转换为蛇形命名
func (su *StringUtils) SnakeCase(s string) string {
	if su.IsEmpty(s) {
		return s
	}

	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}

// KebabCase 转换为短横线命名
func (su *StringUtils) KebabCase(s string) string {
	return strings.ReplaceAll(su.SnakeCase(s), "_", "-")
}

// NumberUtils 数字工具
type NumberUtils struct{}

// IsInteger 检查是否为整数
func (nu *NumberUtils) IsInteger(n float64) bool {
	return n == float64(int64(n))
}

// Round 四舍五入
func (nu *NumberUtils) Round(n float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(n*multiplier) / multiplier
}

// Ceil 向上取整
func (nu *NumberUtils) Ceil(n float64) float64 {
	return math.Ceil(n)
}

// Floor 向下取整
func (nu *NumberUtils) Floor(n float64) float64 {
	return math.Floor(n)
}

// Min 获取最小值
func (nu *NumberUtils) Min(a, b float64) float64 {
	return math.Min(a, b)
}

// Max 获取最大值
func (nu *NumberUtils) Max(a, b float64) float64 {
	return math.Max(a, b)
}

// Abs 获取绝对值
func (nu *NumberUtils) Abs(n float64) float64 {
	return math.Abs(n)
}

// Clamp 限制数值范围
func (nu *NumberUtils) Clamp(n, min, max float64) float64 {
	return math.Max(min, math.Min(max, n))
}

// Lerp 线性插值
func (nu *NumberUtils) Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// FormatNumber 格式化数字
func (nu *NumberUtils) FormatNumber(n float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, n)
}

// ParseFloat 解析浮点数
func (nu *NumberUtils) ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseInt 解析整数
func (nu *NumberUtils) ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// CryptoUtils 加密工具
type CryptoUtils struct{}

// MD5 计算MD5哈希
func (cu *CryptoUtils) MD5(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA1 计算SHA1哈希
func (cu *CryptoUtils) SHA1(data string) string {
	hash := sha1.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA256 计算SHA256哈希
func (cu *CryptoUtils) SHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SHA512 计算SHA512哈希
func (cu *CryptoUtils) SHA512(data string) string {
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Base64Encode Base64编码
func (cu *CryptoUtils) Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode Base64解码
func (cu *CryptoUtils) Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// HexEncode 十六进制编码
func (cu *CryptoUtils) HexEncode(data string) string {
	return hex.EncodeToString([]byte(data))
}

// HexDecode 十六进制解码
func (cu *CryptoUtils) HexDecode(data string) (string, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// TimeUtils 时间工具
type TimeUtils struct{}

// Now 获取当前时间
func (tu *TimeUtils) Now() time.Time {
	return time.Now()
}

// Unix 获取Unix时间戳
func (tu *TimeUtils) Unix() int64 {
	return time.Now().Unix()
}

// UnixMilli 获取Unix毫秒时间戳
func (tu *TimeUtils) UnixMilli() int64 {
	return time.Now().UnixMilli()
}

// UnixMicro 获取Unix微秒时间戳
func (tu *TimeUtils) UnixMicro() int64 {
	return time.Now().UnixMicro()
}

// UnixNano 获取Unix纳秒时间戳
func (tu *TimeUtils) UnixNano() int64 {
	return time.Now().UnixNano()
}

// FormatTime 格式化时间
func (tu *TimeUtils) FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 解析时间
func (tu *TimeUtils) ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// AddDays 添加天数
func (tu *TimeUtils) AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddHours 添加小时
func (tu *TimeUtils) AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddMinutes 添加分钟
func (tu *TimeUtils) AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddSeconds 添加秒数
func (tu *TimeUtils) AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// DiffDays 计算天数差
func (tu *TimeUtils) DiffDays(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Hours() / 24)
}

// DiffHours 计算小时差
func (tu *TimeUtils) DiffHours(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Hours())
}

// DiffMinutes 计算分钟差
func (tu *TimeUtils) DiffMinutes(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Minutes())
}

// DiffSeconds 计算秒数差
func (tu *TimeUtils) DiffSeconds(t1, t2 time.Time) int {
	return int(t1.Sub(t2).Seconds())
}

// IsWeekend 检查是否为周末
func (tu *TimeUtils) IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday 检查是否为工作日
func (tu *TimeUtils) IsWeekday(t time.Time) bool {
	return !tu.IsWeekend(t)
}

// ArrayUtils 数组工具
type ArrayUtils struct{}

// Contains 检查数组是否包含元素
func (au *ArrayUtils) Contains(arr []interface{}, item interface{}) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

// IndexOf 获取元素在数组中的索引
func (au *ArrayUtils) IndexOf(arr []interface{}, item interface{}) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

// Remove 移除数组中的元素
func (au *ArrayUtils) Remove(arr []interface{}, item interface{}) []interface{} {
	var result []interface{}
	for _, v := range arr {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

// Unique 去重
func (au *ArrayUtils) Unique(arr []interface{}) []interface{} {
	keys := make(map[interface{}]bool)
	var result []interface{}

	for _, item := range arr {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Sort 排序
func (au *ArrayUtils) Sort(arr []interface{}, ascending bool) []interface{} {
	result := make([]interface{}, len(arr))
	copy(result, arr)

	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return fmt.Sprintf("%v", result[i]) < fmt.Sprintf("%v", result[j])
		}
		return fmt.Sprintf("%v", result[i]) > fmt.Sprintf("%v", result[j])
	})

	return result
}

// Reverse 反转数组
func (au *ArrayUtils) Reverse(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, v := range arr {
		result[len(arr)-1-i] = v
	}
	return result
}

// Chunk 分块
func (au *ArrayUtils) Chunk(arr []interface{}, size int) [][]interface{} {
	var result [][]interface{}
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}
	return result
}

// MapUtils 映射工具
type MapUtils struct{}

// Get 获取映射值
func (mu *MapUtils) Get(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// Set 设置映射值
func (mu *MapUtils) Set(m map[string]interface{}, key string, value interface{}) {
	m[key] = value
}

// Has 检查映射是否包含键
func (mu *MapUtils) Has(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// Delete 删除映射键
func (mu *MapUtils) Delete(m map[string]interface{}, key string) {
	delete(m, key)
}

// Keys 获取所有键
func (mu *MapUtils) Keys(m map[string]interface{}) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values 获取所有值
func (mu *MapUtils) Values(m map[string]interface{}) []interface{} {
	var values []interface{}
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Merge 合并映射
func (mu *MapUtils) Merge(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Clone 克隆映射
func (mu *MapUtils) Clone(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

// ValidationUtils 验证工具
type ValidationUtils struct{}

// IsEmail 验证邮箱
func (vu *ValidationUtils) IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsURL 验证URL
func (vu *ValidationUtils) IsURL(url string) bool {
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// IsIP 验证IP地址
func (vu *ValidationUtils) IsIP(ip string) bool {
	pattern := `^(\d{1,3}\.){3}\d{1,3}$`
	matched, _ := regexp.MatchString(pattern, ip)
	if !matched {
		return false
	}

	parts := strings.Split(ip, ".")
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// IsNumeric 验证数字
func (vu *ValidationUtils) IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsInteger 验证整数
func (vu *ValidationUtils) IsInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsAlpha 验证字母
func (vu *ValidationUtils) IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric 验证字母数字
func (vu *ValidationUtils) IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsUUID 验证UUID
func (vu *ValidationUtils) IsUUID(s string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, _ := regexp.MatchString(pattern, strings.ToLower(s))
	return matched
}

// JSONUtils JSON工具
type JSONUtils struct{}

// ToJSON 转换为JSON
func (ju *JSONUtils) ToJSON(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToPrettyJSON 转换为格式化JSON
func (ju *JSONUtils) ToPrettyJSON(obj interface{}) (string, error) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从JSON解析
func (ju *JSONUtils) FromJSON(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

// IsValidJSON 验证JSON格式
func (ju *JSONUtils) IsValidJSON(jsonStr string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// MathUtils 数学工具
type MathUtils struct{}

// Factorial 阶乘
func (mu *MathUtils) Factorial(n int) *big.Int {
	if n < 0 {
		return big.NewInt(0)
	}
	if n == 0 || n == 1 {
		return big.NewInt(1)
	}

	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

// GCD 最大公约数
func (mu *MathUtils) GCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM 最小公倍数
func (mu *MathUtils) LCM(a, b int64) int64 {
	return (a * b) / mu.GCD(a, b)
}

// IsPrime 检查是否为质数
func (mu *MathUtils) IsPrime(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	for i := int64(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Fibonacci 斐波那契数列
func (mu *MathUtils) Fibonacci(n int) []int64 {
	if n <= 0 {
		return []int64{}
	}
	if n == 1 {
		return []int64{0}
	}
	if n == 2 {
		return []int64{0, 1}
	}

	fib := make([]int64, n)
	fib[0] = 0
	fib[1] = 1

	for i := 2; i < n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}

	return fib
}
