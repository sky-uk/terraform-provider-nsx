package virtualwire

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// UpdateVirtualWireAPI base api object.
type UpdateVirtualWireAPI struct {
	*api.BaseAPI
}

// NewUpdate returns a new object of UpdateVirtualWireAPI. Returns response code 200 with no content.
func NewUpdate(virtualWire VirtualWire) *UpdateVirtualWireAPI {
	this := new(UpdateVirtualWireAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/2.0/vdn/virtualwires/"+virtualWire.ObjectID, virtualWire, nil)
	return this
}
