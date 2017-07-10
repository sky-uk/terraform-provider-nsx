package sections

type firewallConfiguration struct {
	ContextID string    `xml:"contextId"`
	Sections  []Section `xml:"section,omitempty"`
}

// Section - Contains the rules
type Section struct {
	ID        string `xml:"id,attr,omitempty"`
	Name      string `xml:"name,attr,omitempty"`
	Type      string `xml:"type,attr,omitempty"`
	Timestamp string `xml:"timestamp,attr,omitempty"`
}
