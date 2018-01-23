package firewallexclusion

import "encoding/xml"

// FirewallExclusions - top level xml element
type FirewallExclusions struct {
	XMLName xml.Name `xml:"VshieldAppConfiguration"`
	Members []Member `xml:"excludeListConfiguration>excludeMember>member"`
}

// Member object
type Member struct {
	XMLName xml.Name `xml:"member"`
	MOID    string   `xml:"objectId"`
	Name    string   `xml:"name"`
}
