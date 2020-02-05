// Code generated by `ggen -ent EventReceivers -conf -coll -mods Select,Filter,Top,OrderBy`; DO NOT EDIT.

package api

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (eventReceivers *EventReceivers) Conf(config *RequestConfig) *EventReceivers {
	eventReceivers.config = config
	return eventReceivers
}

// Select adds $select OData modifier
func (eventReceivers *EventReceivers) Select(oDataSelect string) *EventReceivers {
	eventReceivers.modifiers.AddSelect(oDataSelect)
	return eventReceivers
}

// Filter adds $filter OData modifier
func (eventReceivers *EventReceivers) Filter(oDataFilter string) *EventReceivers {
	eventReceivers.modifiers.AddFilter(oDataFilter)
	return eventReceivers
}

// Top adds $top OData modifier
func (eventReceivers *EventReceivers) Top(oDataTop int) *EventReceivers {
	eventReceivers.modifiers.AddTop(oDataTop)
	return eventReceivers
}

// OrderBy adds $orderby OData modifier
func (eventReceivers *EventReceivers) OrderBy(oDataOrderBy string, ascending bool) *EventReceivers {
	eventReceivers.modifiers.AddOrderBy(oDataOrderBy, ascending)
	return eventReceivers
}