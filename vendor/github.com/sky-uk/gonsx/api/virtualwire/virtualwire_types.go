package virtualwire

import "encoding/xml"

// VirtualWires - top level xml element
type VirtualWires struct {
	DataPage DataPage `xml:"dataPage"`
}

// DataPage within VirtualWires
type DataPage struct {
	VirtualWires []VirtualWire `xml:"virtualWire"`
}

// VirtualWire is a single virtual wire object within virtualWire list.
type VirtualWire struct {
	XMLName          xml.Name     `xml:"virtualWire"`
	Name             string       `xml:"name"`
	ObjectID         string       `xml:"objectId,omitempty"`
	ControlPlaneMode string       `xml:"controlPlaneMode"`
	Description      string       `xml:"description"`
	TenantID         string       `xml:"tenantId,omitempty"`
	VdnID            string       `xml:"vdnId,omitempty"`
	VdsContext       []VdsContext `xml:"vdsContextWithBacking,omitempty"`
}

// VdsContext represents a port group for which a VirtualWire is provisioned; a VirtualWire can have several VdsContexts
type VdsContext struct {
	Switch Switch `xml:"switch"`
}

// Switch is the logical switch into which a VirtualWire is "plugged"
type Switch struct {
	ObjectID string `xml:"objectId"`
}

// CreateSpec is used in create call on VirtualWire api.
type CreateSpec struct {
	XMLName          xml.Name `xml:"virtualWireCreateSpec"`
	Name             string   `xml:"name"`
	ControlPlaneMode string   `xml:"controlPlaneMode"`
	Description      string   `xml:"description"`
	TenantID         string   `xml:"tenantId,omitempty"`
}
