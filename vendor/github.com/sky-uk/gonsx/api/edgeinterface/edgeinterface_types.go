package edgeinterface

import "encoding/xml"

type EdgeInterfaces struct {
	XMLName    xml.Name        `xml:"interfaces"`
	Interfaces []EdgeInterface `xml:"interface"`
}

type EdgeInterface struct {
	XMLName       xml.Name      `xml:"interface"`
	Name          string        `xml:"name"`
	Label         string        `xml:"label,omitempty"`
	Mtu           int           `xml:"mtu"`
	Type          string        `xml:"type"`
	IsConnected   bool          `xml:"isConnected"`
	ConnectedToId string        `xml:"connectedToId"`
	AddressGroups AddressGroups `xml:"addressGroups"`
	Index         string        `xml:"index,omitempty"`
}

type AddressGroups struct {
	AddressGroups []AddressGroup `xml:"addressGroup"`
}

type AddressGroup struct {
	XMLName        xml.Name `xml:"addressGroup"`
	PrimaryAddress string   `xml:"primaryAddress"`
	SubnetMask     string   `xml:"subnetMask"`
}
