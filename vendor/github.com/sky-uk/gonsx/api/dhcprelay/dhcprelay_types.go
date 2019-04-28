package dhcprelay

import "encoding/xml"

// DhcpRelay - top level struct.
type DhcpRelay struct {
	XMLName     xml.Name     `xml:"relay"`
	RelayServer RelayServer  `xml:"relayServer"`
	RelayAgents []RelayAgent `xml:"relayAgents>relayAgent"`
}

// RelayServer - relayserver within DhcpRelay object.
type RelayServer struct {
	IPSets     []string `xml:"groupingObjectId"`
	IPAddress  []string `xml:"ipAddress"`
	DomainName []string `xml:"fqdn"`
}

// RelayAgent - relayagent within DhcpRelay object.
type RelayAgent struct {
	XMLName   xml.Name `xml:"relayAgent"`
	VnicIndex string   `xml:"vnicIndex"`
	GiAddress string   `xml:"giAddress,omitempty"`
}
