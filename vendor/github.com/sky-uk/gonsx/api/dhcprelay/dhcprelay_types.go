package dhcprelay

import "encoding/xml"

type DhcpRelay struct {
	XMLName     xml.Name     `xml:"relay"`
	RelayServer RelayServer  `xml:"relayServer"`
	RelayAgents []RelayAgent `xml:"relayAgents>relayAgent"`
}

type RelayServer struct {
	IpAddress string `xml:"ipAddress"`
}

type RelayAgent struct {
	XMLName   xml.Name `xml:"relayAgent"`
	VnicIndex string   `xml:"vnicIndex"`
	GiAddress string   `xml:"giAddress"`
}
