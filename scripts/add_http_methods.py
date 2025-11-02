#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
为新生成的交易所添加 HTTP 方法实现
"""

import os
import re
import sys

# 需要添加方法的交易所列表
EXCHANGES = [
    'alpaca', 'apex', 'arkham', 'backpack', 'bingx', 'hyperliquid', 'paradex',
    'krakenfutures', 'kucoinfutures',
    'bitopro', 'bitrue', 'bitteam', 'bittrade', 'coinsph', 'htx', 'tokocrypto', 'zonda', 'xt', 'toobit',
    'defx', 'derive',
    'blofin', 'coincatch', 'coinmetro', 'cryptomus', 'fmfwio', 'foxbit', 'gate', 'hashkey', 'hibachi', 'mexc', 'modetrade', 'onetrading', 'oxfun', 'p2b', 'woo', 'woofipro'
]

def to_camel_case(name):
    """将下划线命名转换为驼峰命名"""
    parts = name.split('_')
    return ''.join(word.capitalize() for word in parts)

def get_class_name(exchange_id):
    """获取交易所的类名"""
    # 处理特殊情况
    if exchange_id == 'gate':
        return 'Gate'
    elif exchange_id == 'krakenfutures':
        return 'Krakenfutures'
    elif exchange_id == 'kucoinfutures':
        return 'Kucoinfutures'
    elif exchange_id == 'htx':
        return 'Htx'
    elif exchange_id == 'p2b':
        return 'P2b'
    else:
        return to_camel_case(exchange_id)

def extract_api_info(py_file):
    """从 Python 文件中提取 API 信息"""
    if not os.path.exists(py_file):
        return None
    
    try:
        with open(py_file, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"错误: 无法读取文件 {py_file}: {e}")
        return None
    
    api_info = {
        'fetch_markets_method': None,
        'fetch_ticker_method': None,
        'fetch_orderbook_method': None,
    }
    
    # 查找 fetch_markets 方法
    markets_match = re.search(r'def\s+fetch_markets[^:]*:\s*(.*?)(?=\n    def|\nclass|\Z)', content, re.DOTALL)
    if markets_match:
        method_body = markets_match.group(1)
        # 查找 API 调用
        call_match = re.search(r'self\.(\w+)\(', method_body)
        if call_match:
            api_info['fetch_markets_method'] = call_match.group(1)
    
    # 查找 fetch_ticker 方法
    ticker_match = re.search(r'def\s+fetch_ticker[^:]*:\s*(.*?)(?=\n    def|\nclass|\Z)', content, re.DOTALL)
    if ticker_match:
        method_body = ticker_match.group(1)
        call_match = re.search(r'self\.(\w+)\(', method_body)
        if call_match:
            api_info['fetch_ticker_method'] = call_match.group(1)
    
    # 查找 fetch_order_book 方法
    orderbook_match = re.search(r'def\s+fetch_order_book[^:]*:\s*(.*?)(?=\n    def|\nclass|\Z)', content, re.DOTALL)
    if orderbook_match:
        method_body = orderbook_match.group(1)
        call_match = re.search(r'self\.(\w+)\(', method_body)
        if call_match:
            api_info['fetch_orderbook_method'] = call_match.group(1)
    
    return api_info

def generate_http_methods(exchange_id, output_dir, api_info=None):
    """为交易所生成 HTTP 方法"""
    class_name = get_class_name(exchange_id)
    go_file = os.path.join(output_dir, f"ex_{exchange_id}.go")
    
    if not os.path.exists(go_file):
        print(f"警告: {go_file} 不存在，跳过")
        return False
    
    # 读取现有文件
    with open(go_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 检查是否已经有方法实现
    if f'func (this *{class_name}) FetchMarkets' in content:
        print(f"  {exchange_id}: 已有 HTTP 方法实现，跳过")
        return True
    
    # 生成方法代码
    methods_code = f'''

// FetchMarkets 获取交易市场列表
func (this *{class_name}) FetchMarkets(goArgs ...*Variant) *Variant {{
	params := GoGetArg(goArgs, 0, MkMap(&VarMap{{}}))
	_ = params
	
	// TODO: 根据交易所实际 API 实现
	// 基础实现：返回空数组，需要根据交易所文档实现
	result := MkArray(&VarArray{{}})
	
	// 示例：调用公共 API
	// response := this.Call(MkString("publicGetMarkets"), params)
	// markets := this.SafeValue(response, MkString("data"), MkArray(&VarArray{{}}))
	// result = this.ParseMarkets(markets)
	
	return result
}}

// FetchTicker 获取指定交易对的价格信息
func (this *{class_name}) FetchTicker(goArgs ...*Variant) *Variant {{
	symbol := GoGetArg(goArgs, 0, MkUndefined())
	_ = symbol
	params := GoGetArg(goArgs, 1, MkMap(&VarMap{{}}))
	_ = params
	
	this.LoadMarkets()
	market := this.Market(symbol)
	_ = market
	
	// TODO: 根据交易所实际 API 实现
	// 基础实现：返回空对象，需要根据交易所文档实现
	// 示例：调用公共 API
	// response := this.Call(MkString("publicGetTicker"), this.Extend(MkMap(&VarMap{{
	// 	"symbol": *(market).At(MkString("id")),
	// }}), params))
	// return this.ParseTicker(response, market)
	
	return MkMap(&VarMap{{
		"symbol":   symbol,
		"timestamp": MkUndefined(),
		"datetime":  MkUndefined(),
		"high":      MkUndefined(),
		"low":       MkUndefined(),
		"bid":       MkUndefined(),
		"ask":       MkUndefined(),
		"last":      MkUndefined(),
		"volume":    MkUndefined(),
		"info":      MkMap(&VarMap{{}}),
	}})
}}

// FetchOrderBook 获取订单簿
func (this *{class_name}) FetchOrderBook(goArgs ...*Variant) *Variant {{
	symbol := GoGetArg(goArgs, 0, MkUndefined())
	_ = symbol
	limit := GoGetArg(goArgs, 1, MkUndefined())
	_ = limit
	params := GoGetArg(goArgs, 2, MkMap(&VarMap{{}}))
	_ = params
	
	this.LoadMarkets()
	market := this.Market(symbol)
	_ = market
	
	// TODO: 根据交易所实际 API 实现
	// 基础实现：返回空订单簿，需要根据交易所文档实现
	// 示例：调用公共 API
	// request := MkMap(&VarMap{{
	// 	"symbol": *(market).At(MkString("id")),
	// }})
	// if IsTrue(limit) {{
	// 	*(request).At(MkString("limit")) = limit.ToString()
	// }}
	// response := this.Call(MkString("publicGetOrderBook"), this.Extend(request, params))
	// return this.ParseOrderBook(response, market)
	
	return MkMap(&VarMap{{
		"symbol": symbol,
		"bids":   MkArray(&VarArray{{}}),
		"asks":   MkArray(&VarArray{{}}),
		"timestamp": MkUndefined(),
		"datetime":  MkUndefined(),
		"nonce":     MkUndefined(),
	}})
}}

// FetchTrades 获取交易历史
func (this *{class_name}) FetchTrades(goArgs ...*Variant) *Variant {{
	symbol := GoGetArg(goArgs, 0, MkUndefined())
	_ = symbol
	since := GoGetArg(goArgs, 1, MkUndefined())
	_ = since
	limit := GoGetArg(goArgs, 2, MkUndefined())
	_ = limit
	params := GoGetArg(goArgs, 3, MkMap(&VarMap{{}}))
	_ = params
	
	this.LoadMarkets()
	market := this.Market(symbol)
	_ = market
	
	// TODO: 根据交易所实际 API 实现
	return MkArray(&VarArray{{}})
	
	// 示例：调用公共 API
	// request := MkMap(&VarMap{{
	// 	"symbol": *(market).At(MkString("id")),
	// }})
	// if IsTrue(since) {{
	// 	*(request).At(MkString("since")) = since.ToString()
	// }}
	// if IsTrue(limit) {{
	// 	*(request).At(MkString("limit")) = limit.ToString()
	// }}
	// response := this.Call(MkString("publicGetTrades"), this.Extend(request, params))
	// return this.ParseTrades(response, market, since, limit)
}}

// FetchBalance 获取账户余额
func (this *{class_name}) FetchBalance(goArgs ...*Variant) *Variant {{
	params := GoGetArg(goArgs, 0, MkMap(&VarMap{{}}))
	_ = params
	
	// TODO: 根据交易所实际 API 实现
	return MkMap(&VarMap{{
		"info": MkMap(&VarMap{{}}),
	}})
	
	// 示例：调用私有 API
	// response := this.Call(MkString("privateGetAccount"), params)
	// return this.ParseBalance(response)
}}
'''
    
    # 在文件末尾添加方法
    # 找到 Describe 函数结束的位置（最后一个 }）和文件末尾的位置
    # 需要找到 Describe 函数结束的位置，然后在文件末尾添加方法
    
    # 检查是否已经有方法
    if f'func (this *{class_name}) FetchMarkets' in content:
        print(f"  {exchange_id}: 已有 HTTP 方法实现，跳过")
        return True
    
    # 找到 Describe 函数结束的位置
    describe_end = content.rfind('})')
    if describe_end == -1:
        print(f"错误: {go_file} 格式不正确，无法找到 Describe 函数结束位置")
        return False
    
    # 在 Describe 函数结束后添加方法
    # 需要找到 Describe 函数结尾的 }))
    describe_end_idx = content.find('}))', describe_end)
    if describe_end_idx == -1:
        describe_end_idx = content.find('}', describe_end)
    
    if describe_end_idx == -1:
        print(f"错误: {go_file} 格式不正确，无法找到函数结束位置")
        return False
    
    # 在 Describe 函数结束后，文件末尾之前添加方法
    new_content = content[:describe_end_idx + 3] + '\n' + methods_code
    
    # 写入文件
    with open(go_file, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    print(f"  ✓ {exchange_id}: 已添加 HTTP 方法")
    return True

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    ccxt_go_dir = os.path.dirname(script_dir)
    project_root = os.path.dirname(ccxt_go_dir)
    python_dir = os.path.join(project_root, "python", "ccxt")
    go_dir = os.path.join(ccxt_go_dir, "pkg", "ccxt")
    
    if not os.path.exists(go_dir):
        print(f"错误: Go CCXT 目录不存在: {go_dir}")
        sys.exit(1)
    
    print(f"为 {len(EXCHANGES)} 个交易所添加 HTTP 方法...")
    print("=" * 70)
    
    success_count = 0
    failed_count = 0
    
    for exchange_id in EXCHANGES:
        print(f"\n处理 {exchange_id}...")
        
        # 尝试从 Python 文件提取 API 信息
        py_file = os.path.join(python_dir, f"{exchange_id}.py")
        api_info = None
        if os.path.exists(py_file):
            api_info = extract_api_info(py_file)
        
        if generate_http_methods(exchange_id, go_dir, api_info):
            success_count += 1
        else:
            failed_count += 1
    
    print("\n" + "=" * 70)
    print(f"完成!")
    print(f"成功: {success_count}/{len(EXCHANGES)}")
    print(f"失败: {failed_count}/{len(EXCHANGES)}")

if __name__ == "__main__":
    main()

