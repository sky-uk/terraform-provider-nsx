package virtualwire

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// CreateVirtualWireAPI base api object.
type CreateVirtualWireAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateVirtualWireAPI.
func NewCreate(name, desc, tenantID, scopeID string) *CreateVirtualWireAPI {
	this := new(CreateVirtualWireAPI)
	requestPayload := new(CreateSpec)
	requestPayload.Name = name
	requestPayload.TenantID = tenantID
	requestPayload.Description = desc
	// TODO: need to make it argument
	requestPayload.ControlPlaneMode = "UNICAST_MODE"

	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/2.0/vdn/scopes/"+scopeID+"/virtualwires", requestPayload, new(string))
	return this
}

// GetResponse returns ResponseObject of CreateVirtualWireAPI.
func (ca CreateVirtualWireAPI) GetResponse() string {
	return ca.ResponseObject().(string)
}
