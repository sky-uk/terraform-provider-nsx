package virtualwire

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateVirtualWireAPI base api object.
type CreateVirtualWireAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateVirtualWireAPI.
func NewCreate(virtualWireSpec CreateSpec, scopeID string) *CreateVirtualWireAPI {
	this := new(CreateVirtualWireAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/vdn/scopes/"+scopeID+"/virtualwires", virtualWireSpec, new(string))
	return this
}

// GetResponse returns ResponseObject of CreateVirtualWireAPI.
func (ca CreateVirtualWireAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
