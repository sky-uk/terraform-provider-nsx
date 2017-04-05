package tzone

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// GetAllTransportZonesAPI base api object.
type GetAllTransportZonesAPI struct {
	*api.BaseAPI
}

// NewGetAll returns new object of GetAllTransportZonesAPI
func NewGetAll() *GetAllTransportZonesAPI {
	this := new(GetAllTransportZonesAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/2.0/vdn/scopes", nil, new(NetworkScopeList))
	return this
}

// GetResponse returns ResponseObject of GetAllTransportZonesAPI
func (ga GetAllTransportZonesAPI) GetResponse() *NetworkScopeList {
	return ga.ResponseObject().(*NetworkScopeList)
}
