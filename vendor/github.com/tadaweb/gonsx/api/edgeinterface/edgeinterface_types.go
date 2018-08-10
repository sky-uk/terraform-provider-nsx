package edgeinterface

import "encoding/xml"

// EdgeInterfaces top level list object
type EdgeInterfaces struct {
	XMLName    xml.Name        `xml:"interfaces"`
	Interfaces []EdgeInterface `xml:"interface"`
}

// EdgeInterface object within EdgeInterfaces list.
type EdgeInterface struct {
	XMLName       xml.Name      `xml:"interface"`
	Name          string        `xml:"name"`
	Label         string        `xml:"label,omitempty"`
	Mtu           int           `xml:"mtu"`
	Type          string        `xml:"type"`
	IsConnected   bool          `xml:"isConnected"`
	ConnectedToID string        `xml:"connectedToId"`
	AddressGroups AddressGroups `xml:"addressGroups"`
	Index         int           `xml:"index,omitempty"`
}

// AddressGroups within EdgeInterface.
type AddressGroups struct {
	AddressGroups []AddressGroup `xml:"addressGroup"`
}

// AddressGroup object within AddressGroup list.
type AddressGroup struct {
	XMLName        xml.Name `xml:"addressGroup"`
	PrimaryAddress string   `xml:"primaryAddress"`
	SubnetMask     string   `xml:"subnetMask"`
}
