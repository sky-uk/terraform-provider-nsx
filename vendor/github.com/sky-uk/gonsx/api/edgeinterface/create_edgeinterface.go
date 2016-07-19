package edgeinterface

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type CreateEdgeInterfaceApi struct {
	*api.BaseApi
}

func NewCreate(edgeId, interfaceName, virtualWireId, gateway,
		subnetMask, interfaceType string, mtu int) *CreateEdgeInterfaceApi {
	this := new(CreateEdgeInterfaceApi)

	address_group := AddressGroup{PrimaryAddress: gateway, SubnetMask: subnetMask}
	address_group_list := []AddressGroup{address_group}

	edge_interface := EdgeInterface{
		Name: interfaceName,
		ConnectedToId: virtualWireId,
		Type: interfaceType,
		Mtu: mtu,
		IsConnected: true,
		AddressGroups:	AddressGroups{address_group_list},
	}
	requestPayload := &EdgeInterfaces{}
	requestPayload.Interfaces = []EdgeInterface{edge_interface}

	this.BaseApi = api.NewBaseApi(http.MethodPost, "/api/4.0/edges/" + edgeId + "/interfaces/?action=patch", requestPayload, new(EdgeInterfaces))
	return this
}

func (this CreateEdgeInterfaceApi) GetResponse() *EdgeInterfaces {
	return this.ResponseObject().(*EdgeInterfaces)
}
