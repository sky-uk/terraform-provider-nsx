package edgeinterface

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateEdgeInterfaceAPI struct
type CreateEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateEdgeInterfaceAPI
func NewCreate(edgeID, interfaceName, virtualWireID, gateway,
	subnetMask, interfaceType string, mtu int) *CreateEdgeInterfaceAPI {
	this := new(CreateEdgeInterfaceAPI)

	addressGroup := AddressGroup{PrimaryAddress: gateway, SubnetMask: subnetMask}
	addressGroupList := []AddressGroup{addressGroup}

	edgeInterface := EdgeInterface{
		Name:          interfaceName,
		ConnectedToID: virtualWireID,
		Type:          interfaceType,
		Mtu:           mtu,
		IsConnected:   true,
		AddressGroups: AddressGroups{addressGroupList},
	}
	requestPayload := &EdgeInterfaces{}
	requestPayload.Interfaces = []EdgeInterface{edgeInterface}

	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/4.0/edges/"+edgeID+"/interfaces/?action=patch", requestPayload, new(EdgeInterfaces))
	return this
}

// GetResponse returns the ResponseObject.
func (createAPI CreateEdgeInterfaceAPI) GetResponse() *EdgeInterfaces {
	return createAPI.ResponseObject().(*EdgeInterfaces)
}
