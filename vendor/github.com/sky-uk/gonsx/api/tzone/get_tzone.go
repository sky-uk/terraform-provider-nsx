package tzone

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetTransportZoneAPI base api object.
type GetTransportZoneAPI struct {
	*api.BaseAPI
}

// NewGet returns new object of GetTransportZoneAPI
func NewGet(id string) *GetTransportZoneAPI {
	this := new(GetTransportZoneAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/vdn/scopes/"+id, nil, new(NetworkScope))
	return this
}

// GetResponse returns ResponseObject of GetTransportZoneAPI
func (ga GetTransportZoneAPI) GetResponse() *NetworkScope {
	return ga.ResponseObject().(*NetworkScope)
}
