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
	XMLName          xml.Name `xml:"virtualWire"`
	Name             string   `xml:"name"`
	ObjectID         string   `xml:"objectId,omitempty"`
	ControlPlaneMode string   `xml:"controlPlaneMode"`
	Description      string   `xml:"description"`
	TenantID         string   `xml:"tenantId,omitempty"`
}

// CreateSpec is used in create call on VirtualWire api.
type CreateSpec struct {
	XMLName          xml.Name `xml:"virtualWireCreateSpec"`
	Name             string   `xml:"name"`
	ControlPlaneMode string   `xml:"controlPlaneMode"`
	Description      string   `xml:"description"`
	TenantID         string   `xml:"tenantId,omitempty"`
}
