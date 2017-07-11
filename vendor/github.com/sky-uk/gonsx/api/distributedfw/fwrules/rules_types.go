package fwrules

import "encoding/xml"

// Section - Contains the rules
type Section struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Rules []Rule `xml:"rule"`
}

// Rule - The firewall rules
type Rule struct {
	XMLName       xml.Name    `xml:"rule"`
	RuleID        string      `xml:"id,attr,omitempty"`
	Name          string      `xml:"name"`
	Disabled      bool        `xml:"disabled,attr"`
	RuleType      string      `xml:"-"`
	Logged        string      `xml:"logged,attr"`
	Action        string      `xml:"action"`
	AppliedToList *AppliedLst `xml:"appliedToList,ommitempty"`
	Sources       *SourceList `xml:"sources,omitempty"`
	Destinations  *DstList    `xml:"destinations,omitempty"`
	Services      *SvcList    `xml:"services,omitempty"`
	SectionID     int         `xml:"sectionId"`
	Direction     string      `xml:"direction"`
	PacketType    string      `xml:"packetType"`
}

// SvcList List pointer destination
type SvcList struct {
	Services []Service `xml:"services,omitempty>service"`
}

// Service - Struct for the services
type Service struct {
	XMLName         xml.Name `xml:"service"`
	Name            string   `xml:"name,omitempty"`
	Value           string   `xml:"value"`
	DestinationPort int      `xml:"destinationPort"`
	Protocol        int      `xml:"protocol"`
	SubProtocol     int      `xml:"subProtocol,omitempty"`
	Type            string   `xml:"type,omitempty"`
}

// AppliedLst - Applied to list pointer points here
type AppliedLst struct {
	AppliedToList []AppliedTo `xml:"appliedToList"`
}

// AppliedTo - Objects to which the rule is applied
type AppliedTo struct {
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Type    string `xml:"type"`
	IsValid bool   `xml:"isValid"`
}

// SourceList - List of source
type SourceList struct {
	Excluded string   `xml:"excluded,attr,default:'false'"`
	Sources  []Source `xml:"source,omitempty"`
}

// Source - The source for the rule
type Source struct {
	Name    string `xml:"name,omitempty"`
	Type    string `xml:"type,omitempty"`
	Value   string `xml:"value,omitempty"`
	IsValid *bool  `xml:"isValid,omitempty"`
}

// DstList - Destination List
type DstList struct {
	Excluded     string        `xml:"excluded,attr"`
	Destinations []Destination `xml:"destination,omitempty"`
}

// Destination - The destination for the rule
type Destination struct {
	Name    string `xml:"name"`
	Type    string `xml:"type"`
	Value   string `xml:"value"`
	IsValid bool   `xml:"isValid"`
}
