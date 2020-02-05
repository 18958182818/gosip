// Code generated by `ggen -ent Group -conf -mods Select,Expand -helpers Data,Normalized`; DO NOT EDIT.

package api

import "encoding/json"

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (group *Group) Conf(config *RequestConfig) *Group {
	group.config = config
	return group
}

// Select adds $select OData modifier
func (group *Group) Select(oDataSelect string) *Group {
	group.modifiers.AddSelect(oDataSelect)
	return group
}

// Expand adds $expand OData modifier
func (group *Group) Expand(oDataExpand string) *Group {
	group.modifiers.AddExpand(oDataExpand)
	return group
}

/* Response helpers */

// Data response helper
func (groupResp *GroupResp) Data() *GroupInfo {
	data := NormalizeODataItem(*groupResp)
	res := &GroupInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (groupResp *GroupResp) Normalized() []byte {
	return NormalizeODataItem(*groupResp)
}