package virtualwire

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type CreateVirtualWireApi struct {
	*api.BaseApi
}

func NewCreate(name, desc, tenantID, scopeId string) *CreateVirtualWireApi {
	this := new(CreateVirtualWireApi)
	requestPayload := new(VirtualWireCreateSpec)
	requestPayload.Name = name
	requestPayload.TenantID = tenantID
	requestPayload.Description = desc
	// TODO: need to make it argument
	requestPayload.ControlPlaneMode = "UNICAST_MODE"

	this.BaseApi = api.NewBaseApi(http.MethodPost, "/api/2.0/vdn/scopes/" + scopeId +"/virtualwires", requestPayload, new(string))
	return this
}

func (this CreateVirtualWireApi) GetResponse() string{
	return this.ResponseObject().(string)
}
