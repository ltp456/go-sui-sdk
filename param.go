package go_sui_sdk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type JsonRpc struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func NewJsonRpc(method string, param interface{}) JsonRpc {
	if param == nil {
		return JsonRpc{
			Jsonrpc: "2.0",
			ID:      time.Now().Nanosecond(),
			Method:  method,
			Params:  []interface{}{},
		}
	} else {
		return JsonRpc{
			Jsonrpc: "2.0",
			ID:      time.Now().Nanosecond(),
			Method:  method,
			Params:  param,
		}
	}

}

type JsonResp struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   Error           `json:"error"`
	Id      int             `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MapParams map[string]interface{}

func (mp *MapParams) SetKey(key string, value interface{}) {
	(*mp)[key] = value
}

type Params []interface{}

func (p *Params) AddValue(value interface{}) {
	*p = append(*p, value)
}

type UrlParams map[string]interface{}

func (p UrlParams) SetValue(key string, value interface{}) {
	p[key] = value
}

func (p UrlParams) Encode() string {
	if p == nil {
		return ""
	}

	var buf strings.Builder
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		vs := p[k]
		keyEscaped := url.QueryEscape(k)
		if i == len(keys)-1 {
			buf.WriteString(fmt.Sprintf("%v=%v", keyEscaped, vs))
		} else {
			buf.WriteString(fmt.Sprintf("%v=%v&", keyEscaped, vs))
		}

	}
	return buf.String()
}

type Option struct {
	Key   string
	Value interface{}
}

type Options []Option

func (ov *Options) Add(key string, value interface{}) {
	*ov = append(*ov, Option{key, value})
}
