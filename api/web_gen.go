// Code generated by `ggen -ent Web -conf -mods Select,Expand -helpers Data,Normalized`; DO NOT EDIT.

package api

import "encoding/json"

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (web *Web) Conf(config *RequestConfig) *Web {
	web.config = config
	return web
}

// Select adds $select OData modifier
func (web *Web) Select(oDataSelect string) *Web {
	web.modifiers.AddSelect(oDataSelect)
	return web
}

// Expand adds $expand OData modifier
func (web *Web) Expand(oDataExpand string) *Web {
	web.modifiers.AddExpand(oDataExpand)
	return web
}

/* Response helpers */

// Data response helper
func (webResp *WebResp) Data() *WebInfo {
	data := NormalizeODataItem(*webResp)
	res := &WebInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (webResp *WebResp) Normalized() []byte {
	return NormalizeODataItem(*webResp)
}