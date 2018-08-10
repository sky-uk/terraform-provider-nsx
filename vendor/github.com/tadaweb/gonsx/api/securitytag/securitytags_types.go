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
	Revision    int      `xml:"revision,omitempty"`
}

// AttachmentList struct
type AttachmentList struct {
	XMLName                xml.Name     `xml:"securityTags"`
	SecurityTagAttachments []Attachment `xml:"securityTag"`
}

// Attachment object struct
type Attachment struct {
	//XMLName	xml.Name `xml:"securityTag"`
	ObjectID string `xml:"objectId"`
}

// BasicInfoList struct to get info of vms attached to tags
type BasicInfoList struct {
	BasicInfoList []BasicInfo `xml:"basicinfo"`
}

// BasicInfo gives info of list
type BasicInfo struct {
	ObjectID string `xml:"objectId"`
	Name     string `xml:"name"`
}
