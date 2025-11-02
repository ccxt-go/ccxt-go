package ccxt

import (
	"strings"
)

//
// This file implements the rest api request mechanisms
//

func (this *ExchangeBase) BuildApiPath(name string, prefix string, level int, v *Variant) (path string, method string, type_ string) {
	if v.Type == Map {
		for key, value := range *v.ToMap() {
			newPrefix := prefix + key
			if len(name) < len(newPrefix) || strings.ToLower(name[:len(newPrefix)]) != strings.ToLower(newPrefix) {
				continue
			}

			if level == 2 {
				return key, "", ""
			}

			tmp1, tmp2, _ := this.BuildApiPath(name, newPrefix, level+1, *value)
			if tmp2 == "" {
				tmp2 = strings.ToUpper(key)
			}
			return tmp1, tmp2, newPrefix
		}
	}
	/*if v.Type == Array {
		pos := v.IndexOf(MkString(strings.ToLower(name[len(prefix):])))
		if pos.ToInt() == -1 {
			return "", "", ""
		}
		return (*v.GetRef(pos)).ToStr(), "", ""
	}*/
	return "", "", ""
}

func (this *ExchangeBase) TryCallAPI(name string, args ...*Variant) *Variant {

	api := *this.At(MkString("api"))
	subPath, method, type_ := this.BuildApiPath(name, "", 0, api)
	if subPath == "" {
		return nil
	}

	params := MkMap(&VarMap{})
	if len(args) >= 1 {
		params = args[0]
	}

	//return this.Request(MkString(subPath), MkString(FixJsName(callItems[0])), MkString(method), params)
	return this.Request(MkString(subPath), MkString(type_), MkString(method), params)
}

func (this *ExchangeBase) ExecuteRestRequest(vUrl *Variant, vMethod *Variant, vHeaders *Variant, vBody *Variant) *Variant {
	// 使用新的统一客户端
	return this.UnifiedHTTPRequest(vUrl, MkString("public"), vMethod, MkMap(&VarMap{}), vHeaders, vBody)
}

func (this *ExchangeBase) Fetch(v ...*Variant) *Variant {
	// transpiled function
	url := GoGetArg(v, 0, MkUndefined())
	method := GoGetArg(v, 1, MkString("GET"))
	headers := GoGetArg(v, 2, MkUndefined())
	body := GoGetArg(v, 3, MkUndefined())
	isNode := true

	if isNode && IsTrue(*this.At(MkString("userAgent"))) {
		if OpType(*this.At(MkString("userAgent"))).ToStr() == "string" {
			headers = this.Extend(MkMap(&VarMap{"User-Agent": *this.At(MkString("userAgent"))}), headers)
		} else {
			if OpType(*this.At(MkString("userAgent"))).ToStr() == "object" && (*this.At(MkString("userAgent"))).Has("User-Agent") {
				headers = this.Extend(*this.At(MkString("userAgent")), headers)
			}
		}
	}

	// todo:
	/*if IsTrue(OpEq2(OpType(*this.At(MkString("proxy"))), MkString("function"))) {
		url = this.Call(MkString("proxy"), url)
		if isNode {
			headers = this.Extend(MkMap(&VarMap{"Origin": *this.At(MkString("origin"))}), headers)
		}
	} else {
		if IsTrue(OpEq2(OpType(*this.At(MkString("proxy"))), MkString("string"))) {
			if isNode && IsTrue((*this.At(MkString("proxy"))).Length) {
				headers = this.Extend(MkMap(&VarMap{"Origin": *this.At(MkString("origin"))}), headers)
			}
			url = OpAdd(*this.At(MkString("proxy")), url)
		}
	}*/
	headers = this.Extend(*this.At(MkString("headers")), headers)
	headers = this.SetHeaders(headers)
	//if IsTrue(*this.At(MkString("verbose")) ) {
	//	this.Call(MkString("print"), MkString("fetch:\n") , this.Id , method , url , MkString("\nRequest:\n") , headers , MkString("\n") , body , MkString("\n") )
	//}

	return this.ExecuteRestRequest(url, method, headers, body)
}

func (this *ExchangeBase) SetHeaders(goArgs ...*Variant) *Variant {
	ret := this.VCall("SetHeaders", goArgs...)
	if ret != nil {
		return ret
	}
	return goArgs[0]
}

func (this *ExchangeBase) Fetch2(v ...*Variant) *Variant {

	// transpiled function
	//path := GoGetArg(v, 0, MkUndefined());
	//type_ := GoGetArg(v, 1, MkString("public"));
	//method := GoGetArg(v, 2, MkString("GET"));
	//params := GoGetArg(v, 3, MkMap(&VarMap{}));
	//headers := GoGetArg(v, 4, MkUndefined());
	//body := GoGetArg(v, 5, MkUndefined());

	// todo
	//if IsTrue(*this.At(MkString("enableRateLimit"))) {
	//	this.Throttle(*this.At(MkString("rateLimit")))
	//}

	request := this.Sign(v...)

	return this.Fetch(*request.At(MkString("url")), *request.At(MkString("method")), *request.At(MkString("headers")), *request.At(MkString("body")))
}

func (this *ExchangeBase) Sign(goArgs ...*Variant) *Variant {
	ret := this.VCall("Sign", goArgs...)
	if ret != nil {
		return ret
	}
	panic("Sign is a mandatory function and must be implemented by the exchange")
}

func (this *ExchangeBase) HandleErrors(goArgs ...*Variant) *Variant {
	ret := this.VCall("HandleErrors", goArgs...)
	if ret != nil {
		return ret
	}
	return MkUndefined()
}

func (this *ExchangeBase) Request(goArgs ...*Variant) *Variant {
	ret := this.VCall("Request", goArgs...)
	if ret != nil {
		return ret
	}
	return this.Fetch2(goArgs...)
}
