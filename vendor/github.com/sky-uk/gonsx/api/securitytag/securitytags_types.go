package securitytag

import "encoding/xml"

// SecurityTags top level struct
type SecurityTags struct {
	SecurityTags []SecurityTag `xml:"securityTag"`
}

// SecurityTag object struct
type SecurityTag struct {
	XMLName     xml.Name `xml:"securityTag"`
	ObjectID    string   `xml:"objectId,omitempty"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	TypeName    string   `xml:"type>typeName"`
}
