package virtualwire

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetVirtualWireAPI base object.
type GetVirtualWireAPI struct {
	*api.BaseAPI
}

// NewGet returns new object of GetVirtualWireAPI.
func NewGet(id string) *GetVirtualWireAPI {
	this := new(GetVirtualWireAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/vdn/virtualwires/"+id, nil, new(VirtualWire))
	return this
}

// GetResponse returns ResponseObject of GetVirtualWireAPI.
func (ga GetVirtualWireAPI) GetResponse() *VirtualWire {
	return ga.ResponseObject().(*VirtualWire)
}
