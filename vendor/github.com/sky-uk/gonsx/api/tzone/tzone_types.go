package tzone

// NetworkScopeList top level xml element
type NetworkScopeList struct {
	NetworkScopeList []NetworkScope `xml:"vdnScope"`
}

// NetworkScope object within NetworkScopeList.
type NetworkScope struct {
	ObjectID string `xml:"objectId"`
	Name     string `xml:"name"`
}
