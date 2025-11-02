#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Python 交易所文件转 Go 文件的转换工具
"""

import os
import re
import sys

# 需要支持的交易所列表
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

def parse_python_exchange(py_file):
    """解析 Python 交易所文件，提取基本信息"""
    try:
        with open(py_file, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"错误: 无法读取文件 {py_file}: {e}")
        return None
    
    # 提取 class 名
    class_match = re.search(r'class\s+(\w+)\s*\([^\)]+\):', content)
    class_name = class_match.group(1) if class_match else None
    
    # 提取 describe 函数中的基本信息 - 使用更宽松的匹配
    # 匹配从 describe 函数开始到返回字典结束
    describe_match = re.search(
        r"def\s+describe\(self\)[^:]*:\s*return\s+self\.deep_extend\([^,]*describe\(\)[^,]*,\s*\{(.*?)\}\s*\)",
        content,
        re.DOTALL
    )
    
    # 如果失败，尝试匹配到函数结束
    if not describe_match:
        describe_match = re.search(
            r"def\s+describe\(self\)[^:]*:\s*return\s+self\.deep_extend\([^,]*describe\(\)[^,]*,\s*\{(.*)\}",
            content,
            re.DOTALL
        )
    
    if not describe_match:
        # 即使无法解析，也创建基本信息
        exchange_id = os.path.basename(py_file).replace('.py', '')
        return {
            'id': exchange_id,
            'name': class_name or exchange_id.capitalize(),
            'countries': [],
            'rateLimit': '2000',
            'hostname': None,
            'pro': False,
            'certified': False,
            'class_name': class_name
        }
    
    describe_content = describe_match.group(1)
    
    # 提取 id
    id_match = re.search(r"['\"]id['\"]:\s*['\"]([^'\"]+)['\"]", describe_content)
    exchange_id = id_match.group(1) if id_match else os.path.basename(py_file).replace('.py', '')
    
    # 提取 name
    name_match = re.search(r"['\"]name['\"]:\s*['\"]([^'\"]+)['\"]", describe_content)
    exchange_name = name_match.group(1) if name_match else exchange_id
    
    # 提取 countries
    countries_match = re.search(r"['\"]countries['\"]:\s*(\[[^\]]*\])", describe_content)
    countries = []
    if countries_match:
        countries_str = countries_match.group(1)
        countries = re.findall(r"['\"]([^'\"]+)['\"]", countries_str)
    
    # 提取 rateLimit
    rate_limit_match = re.search(r"['\"]rateLimit['\"]:\s*(\d+)", describe_content)
    rate_limit = rate_limit_match.group(1) if rate_limit_match else "2000"
    
    # 提取 hostname
    hostname_match = re.search(r"['\"]hostname['\"]:\s*['\"]([^'\"]+)['\"]", describe_content)
    hostname = hostname_match.group(1) if hostname_match else None
    
    # 检查是否有 pro 标志
    pro_match = re.search(r"['\"]pro['\"]:\s*(True|False)", describe_content)
    pro = pro_match and pro_match.group(1) == 'True'
    
    # 检查是否有 certified 标志
    certified_match = re.search(r"['\"]certified['\"]:\s*(True|False)", describe_content)
    certified = certified_match and certified_match.group(1) == 'True'
    
    return {
        'id': exchange_id,
        'name': exchange_name,
        'countries': countries,
        'rateLimit': rate_limit,
        'hostname': hostname,
        'pro': pro,
        'certified': certified,
        'class_name': class_name
    }

def generate_go_file(exchange_info, output_dir):
    """生成 Go 交易所文件"""
    exchange_id = exchange_info['id']
    exchange_name = exchange_info['name']
    go_class_name = to_camel_case(exchange_id)
    
    # 如果 class_name 是复数形式，使用原 class_name
    if exchange_info.get('class_name'):
        go_class_name = exchange_info['class_name'].capitalize()
    
    # 处理特殊情况
    if exchange_id == 'gate':
        go_class_name = 'Gate'
    elif exchange_id == 'krakenfutures':
        go_class_name = 'Krakenfutures'
    elif exchange_id == 'kucoinfutures':
        go_class_name = 'Kucoinfutures'
    elif exchange_id == 'htx':
        go_class_name = 'Htx'
    
    # 生成 countries 数组
    countries_code = "MkArray(&VarArray{"
    for country in exchange_info['countries']:
        countries_code += f"\n\t\t\tMkString(\"{country}\"),"
    if exchange_info['countries']:
        countries_code += "\n\t\t"
    countries_code += "})"
    
    if not exchange_info['countries']:
        countries_code = "MkArray(&VarArray{})"
    
    # 生成 status 字段
    status_code = '''
		"status": MkMap(&VarMap{
			"status": MkString("ok"),
		}),'''
    
    go_code = f"""package ccxt

type {go_class_name} struct {{
	*ExchangeBase
}}

var _ Exchange = (*{go_class_name})(nil)

func init() {{
	exchange := &{go_class_name}{{}}
	Exchanges = append(Exchanges, exchange)
}}

func (this *{go_class_name}) Describe(goArgs ...*Variant) *Variant {{
	return this.DeepExtend(this.BaseDescribe(), MkMap(&VarMap{{
		"id":   MkString("{exchange_id}"),
		"name": MkString("{exchange_name}"),
		"countries": {countries_code},
		"rateLimit": MkInteger({exchange_info['rateLimit']}),
		{status_code}
		"has": MkMap(&VarMap{{
			"cancelOrder":   MkBool(true),
			"createOrder":   MkBool(true),
			"fetchBalance":  MkBool(true),
			"fetchMarkets":  MkBool(true),
			"fetchOrderBook": MkBool(true),
			"fetchTicker":   MkBool(true),
			"fetchTrades":   MkBool(true),
		}}),
		"urls": MkMap(&VarMap{{
			"api": MkMap(&VarMap{{
				"public":  MkString("https://api.{{hostname}}"),
				"private": MkString("https://api.{{hostname}}"),
			}}),
			"www": MkString("https://www.{{hostname}}"),
		}}),
	}}))
}}
"""
    
    # 替换 hostname 占位符
    if exchange_info.get('hostname'):
        hostname_value = exchange_info['hostname']
    else:
        hostname_value = exchange_id + '.com'
    
    go_code = go_code.replace('{hostname}', hostname_value)
    
    # 写入文件
    output_file = os.path.join(output_dir, f"ex_{exchange_id}.go")
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(go_code)
    
    print(f"生成: {output_file}")
    return output_file

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    # scripts 目录在 ccxt-go 下，需要向上两级到项目根目录
    ccxt_go_dir = os.path.dirname(script_dir)
    project_root = os.path.dirname(ccxt_go_dir)
    python_dir = os.path.join(project_root, "python", "ccxt")
    go_dir = os.path.join(ccxt_go_dir, "pkg", "ccxt")
    
    if not os.path.exists(python_dir):
        print(f"错误: Python CCXT 目录不存在: {python_dir}")
        sys.exit(1)
    
    if not os.path.exists(go_dir):
        print(f"错误: Go CCXT 目录不存在: {go_dir}")
        sys.exit(1)
    
    print(f"Python CCXT 目录: {python_dir}")
    print(f"Go CCXT 输出目录: {go_dir}")
    print(f"\n开始转换 {len(EXCHANGES)} 个交易所...\n")
    
    converted = []
    failed = []
    
    for exchange_id in EXCHANGES:
        py_file = os.path.join(python_dir, f"{exchange_id}.py")
        
        if not os.path.exists(py_file):
            print(f"✗ {exchange_id}: Python 文件不存在")
            failed.append(exchange_id)
            continue
        
        print(f"处理 {exchange_id}...")
        exchange_info = parse_python_exchange(py_file)
        
        if exchange_info:
            try:
                generate_go_file(exchange_info, go_dir)
                converted.append(exchange_id)
            except Exception as e:
                print(f"✗ {exchange_id}: 生成 Go 文件失败: {e}")
                failed.append(exchange_id)
        else:
            print(f"✗ {exchange_id}: 解析失败")
            failed.append(exchange_id)
    
    print(f"\n转换完成!")
    print(f"成功: {len(converted)}/{len(EXCHANGES)}")
    print(f"失败: {len(failed)}/{len(EXCHANGES)}")
    
    if failed:
        print(f"\n失败的交易所: {', '.join(failed)}")

if __name__ == "__main__":
    main()

