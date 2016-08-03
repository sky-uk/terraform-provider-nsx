package virtualwire

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// UpdateVirtualWireAPI base api object.
type UpdateVirtualWireAPI struct {
	*api.BaseAPI
}

// NewUpdate returns a new object of UpdateVirtualWireAPI.
func NewUpdate(name, desc, virtualwireID string) *UpdateVirtualWireAPI {
	this := new(UpdateVirtualWireAPI)
	requestPayload := new(VirtualWire)
	requestPayload.Name = name
	requestPayload.ControlPlaneMode = "UNICAST_MODE"
	requestPayload.Description = desc

	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/2.0/vdn/virtualwires/"+virtualwireID, requestPayload, nil)
	return this
}
