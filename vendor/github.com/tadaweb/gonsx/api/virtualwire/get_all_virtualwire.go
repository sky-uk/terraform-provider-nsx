package virtualwire

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllVirtualWiresAPI base object.
type GetAllVirtualWiresAPI struct {
	*api.BaseAPI
}

// NewGetAll returns a new object of GetAllVirtualWiresAPI.
func NewGetAll(scopeID string) *GetAllVirtualWiresAPI {
	this := new(GetAllVirtualWiresAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/vdn/scopes/"+scopeID+"/virtualwires", nil, new(VirtualWires))
	return this
}

// GetResponse returns ResponseObject of GetAllVirtualWiresAPI.
func (ga GetAllVirtualWiresAPI) GetResponse() *VirtualWires {
	return ga.ResponseObject().(*VirtualWires)
}
