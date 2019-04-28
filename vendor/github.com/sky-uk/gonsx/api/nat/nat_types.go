package nat

import "encoding/xml"

// Nat object
type Nat struct {
	XMLName xml.Name `xml:"nat"`
	Rules   Rules    `xml:"natRules"`
}

// Rules object within Nat
type Rules struct {
	XMLName xml.Name `xml:"natRules"`
	Rules   []Rule   `xml:"natRule"`
}

// Rule object within Rules list
type Rule struct {
	XMLName                     xml.Name `xml:"natRule"`
	RuleID                      string   `xml:"ruleId,omitempty"`
	Description                 string   `xml:"description"`
	Enabled                     bool     `xml:"enabled,omitempty"`
	LoggingEnabled              bool     `xml:"loggingEnabled,omitempty"`
	Action                      string   `xml:"action"`
	Vnic                        string   `xml:"vnic,omitempty"`
	OriginalAddress             string   `xml:"originalAddress"`
	TranslatedAddress           string   `xml:"translatedAddress"`
	DnatMatchSourceAddress      string   `xml:"dnatMatchSourceAddress,omitempty"`
	SnatMatchDestinationAddress string   `xml:"snatMatchDestinationAddress,omitempty"`
	IcmpType                    string   `xml:"icmpType,omitempty"`
	OriginalPort                string   `xml:"originalPort,omitempty"`
	TranslatedPort              string   `xml:"translatedPort,omitempty"`
	DnatMatchSourcePort         string   `xml:"dnatMatchSourcePort,omitempty"`
	SnatMatchDestinationPort    string   `xml:"snatMatchDestinationPort,omitempty"`
	Protocol                    string   `xml:"protocol,omitempty"`
}
