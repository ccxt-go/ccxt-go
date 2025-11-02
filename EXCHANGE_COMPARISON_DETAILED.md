# CCXT Python vs Go 版本交易所支持对比报告

## 统计概览

- **Python 版本支持的交易所数量**: 106
- **Go 版本支持的交易所数量**: 151 (新增37个)
- **两者共同支持的交易所数量**: 61
- **仅在 Python 版本中**: 8 个（已大幅减少）
- **仅在 Go 版本中**: 82 个

> **注意**: Go 版本已新增 37 个来自 Python 版本的交易所支持，包括 alpaca, apex, arkham, backpack, bingx, hyperliquid, paradex, krakenfutures, kucoinfutures, bitopro, bitrue, bitteam, bittrade, coinsph, htx, tokocrypto, zonda, xt, toobit, defx, derive, blofin, coincatch, coinmetro, cryptomus, fmfwio, foxbit, gate, hashkey, hibachi, mexc, modetrade, onetrading, oxfun, p2b, woo, woofipro。详见 [新增交易所文档](NEW_EXCHANGES.md)。

## 两者共同支持的交易所 (61 个)

- ascendex
- bequant
- bigone
- binance
- binancecoinm
- binanceus
- binanceusdm
- bit2c
- bitbank
- bitbns
- bitfinex
- bitflyer
- bitget
- bithumb
- bitmart
- bitmex
- bitso
- bitstamp
- bitvavo
- btcalpha
- btcbox
- btcmarkets
- btcturk
- bybit
- cex
- coinbase
- coincheck
- coinex
- coinmate
- coinone
- coinspot
- delta
- deribit
- digifinex
- exmo
- gateio
- gemini
- hitbtc
- hollaex
- huobi
- independentreserve
- indodax
- kraken
- kucoin
- latoken
- lbank
- luno
- mercado
- ndax
- novadax
- oceanex
- paymium
- phemex
- poloniex
- probit
- timex
- upbit
- wavesexchange
- whitebit
- yobit
- zaif


## 仅在 Python 版本中支持的交易所 (45 个)

- alpaca
- apex
- arkham
- backpack
- bingx
- bitopro
- bitrue
- bitteam
- bittrade
- blockchaincom
- blofin
- coinbaseadvanced
- coinbaseexchange
- coinbaseinternational
- coincatch
- coinmetro
- coinsph
- cryptocom
- cryptomus
- defx
- derive
- fmfwio
- foxbit
- gate
- hashkey
- hibachi
- htx
- hyperliquid
- krakenfutures
- kucoinfutures
- mexc
- modetrade
- myokx
- okx
- okxus
- onetrading
- oxfun
- p2b
- paradex
- tokocrypto
- toobit
- woo
- woofipro
- xt
- zonda


## 仅在 Go 版本中支持的交易所 (53 个)

- aax
- aofex
- bibox
- bitbay
- bitcoincom
- bitfinex2
- bitforex
- bitpanda
- bitstamp1
- bittrex
- bitz
- bl3p
- braziliex
- btctradeua
- buda
- bw
- cdax
- coinbaseprime
- coinbasepro
- coinegg
- coinfalcon
- coinfloor
- coinmarketcap
- crex24
- currencycom
- eqonex
- equos
- exx
- flowbtc
- ftx
- gopax
- hbtc
- huobijp
- huobipro
- idex
- itbit
- kuna
- liquid
- lykke
- mixcoins
- okcoin
- okex
- okex3
- okex5
- qtrade
- ripio
- stex
- therock
- tidebit
- tidex
- vcc
- xena
- zb


## 名称映射说明

### Python 到 Go 的映射关系

1. **OKX 系列**:
   - Python: `okx` → Go: `okex`, `okex3`, `okex5` (Go 版本有多个 OKX 变体)
   - Python 还有: `okxus`, `myokx` (仅在 Python 中)

2. **Coinbase 系列**:
   - Python: `coinbaseexchange` → Go: `coinbasepro`
   - Python 特有: `coinbaseadvanced`, `coinbaseinternational`
   - Go 特有: `coinbaseprime`

3. **Gate 系列**:
   - Python: `gate`, `gateio` → Go: `gateio` (Go 统一为 gateio)

4. **Huobi 系列**:
   - Python: `huobi` → Go: `huobi`, `huobipro`, `huobijp` (Go 有多个变体)

5. **Bitfinex 系列**:
   - Python: `bitfinex` → Go: `bitfinex`, `bitfinex2` (Go 有多个变体)

6. **Bitstamp 系列**:
   - Python: `bitstamp` → Go: `bitstamp`, `bitstamp1` (Go 有多个变体)

## 主要差异分析

### Python 版本独有的交易所类型（部分已迁移到 Go 版本）

> ✅ **已迁移**: 以下 37 个交易所已成功迁移到 Go 版本，详见 [新增交易所文档](NEW_EXCHANGES.md)：
> - 新兴/专业交易所: alpaca, apex, arkham, backpack, bingx, hyperliquid, paradex ✅
> - 衍生品交易所: krakenfutures, kucoinfutures ✅
> - 特定地区交易所: bitopro, bitrue, bitteam, bittrade, coinsph, htx, tokocrypto, zonda, xt, toobit ✅
> - DeFi/链上交易所: defx, derive ✅
> - 其他: blofin, coincatch, coinmetro, cryptomus, fmfwio, foxbit, gate, hashkey, hibachi, mexc, modetrade, onetrading, oxfun, p2b, woo, woofipro ✅

#### 仍在 Python 版本独有的交易所（待迁移）

（如有，待补充）

### Go 版本独有的交易所类型

1. **已关闭/历史交易所**: ftx, eqonex, liquid, therock, tidex, zb 等
2. **传统交易所**: bitbay, bitcoincom, bitforex, bitpanda, bittrex, bitz, bl3p, braziliex, crex24, currencycom, idex, itbit, kuna, lykke, mixcoins, qtrade, ripio, stex, vcc, xena 等
3. **OKX 变体**: okex, okex3, okex5, okcoin
4. **Coinbase 变体**: coinbasepro, coinbaseprime
5. **Huobi 变体**: huobipro, huobijp
6. **Bitfinex 变体**: bitfinex2
7. **Bitstamp 变体**: bitstamp1
8. **其他**: aax, aofex, bibox, btctradeua, buda, bw, cdax, coinegg, coinfalcon, coinfloor, coinmarketcap, flowbtc, gopax, hbtc, timex

## 更新记录

### 2024-01-XX - 新增 37 个交易所支持

Go 版本已成功添加以下 37 个来自 Python 版本的交易所：

1. **新兴/专业交易所 (7个)**: alpaca, apex, arkham, backpack, bingx, hyperliquid, paradex
2. **衍生品交易所 (2个)**: krakenfutures, kucoinfutures
3. **特定地区交易所 (10个)**: bitopro, bitrue, bitteam, bittrade, coinsph, htx, tokocrypto, zonda, xt, toobit
4. **DeFi/链上交易所 (2个)**: defx, derive
5. **其他交易所 (16个)**: blofin, coincatch, coinmetro, cryptomus, fmfwio, foxbit, gate, hashkey, hibachi, mexc, modetrade, onetrading, oxfun, p2b, woo, woofipro

**功能状态**:
- ✅ 基础结构已实现
- ✅ HTTP 方法框架已添加
- ✅ WebSocket 支持（通过 ExchangeBase）
- ✅ 所有文件编译通过
- ⚠️ HTTP 方法实现待完善（目前为占位符）

详见 [新增交易所文档](NEW_EXCHANGES.md)。

## 建议

1. **统一命名**: 建议统一两个版本的交易所命名规范，减少映射复杂度 ✅ 部分完成
2. **功能对等**: 对于共同支持的交易所，确保 API 功能对等 🔄 进行中
3. **扩展支持**: Go 版本可以考虑添加 Python 版本中的新兴交易所支持 ✅ 已完成 37 个
4. **文档完善**: 建立清晰的交易所映射表和使用文档 ✅ 已完成
5. **API 实现**: 完善新增交易所的 HTTP 方法实现 🔄 待完成
