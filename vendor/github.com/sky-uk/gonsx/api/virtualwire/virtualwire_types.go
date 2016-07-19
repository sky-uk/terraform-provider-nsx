package virtualwire

import "encoding/xml"

type VirtualWires struct {
	DataPage		DataPage	`xml:"dataPage"`
}

type DataPage struct {
	VirtualWires		[]VirtualWire	`xml:"virtualWire"`
}

type VirtualWire struct {
	XMLName			xml.Name	`xml:"virtualWire"`
	Name			string		`xml:"name"`
	ObjectID		string		`xml:"objectId,omitempty"`
	ControlPlaneMode	string		`xml:"controlPlaneMode"`
	Description		string		`xml:"description"`
	TenantID		string		`xml:"tenantId,omitempty"`
}

type VirtualWireCreateSpec struct {
	XMLName			xml.Name	`xml:"virtualWireCreateSpec"`
	Name			string		`xml:"name"`
	ControlPlaneMode	string		`xml:"controlPlaneMode"`
	Description		string		`xml:"description"`
	TenantID		string		`xml:"tenantId,omitempty"`
}
