package virtualwire

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// DeleteVirtualWiresAPI base object.
type DeleteVirtualWiresAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new object of DeleteVirtualWiresAPI.
func NewDelete(virtualWireID string) *DeleteVirtualWiresAPI {
	this := new(DeleteVirtualWiresAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/2.0/vdn/virtualwires/"+virtualWireID, nil, nil)
	return this
}
