package edgeinterface

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// GetAllEdgeInterfacesAPI base object.
type GetAllEdgeInterfacesAPI struct {
	*api.BaseAPI
}

// NewGetAll returns the api object of GetAllEdgeInterfacesAPI
func NewGetAll(edgeID string) *GetAllEdgeInterfacesAPI {
	this := new(GetAllEdgeInterfacesAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/4.0/edges/"+edgeID+"/interfaces", nil, new(EdgeInterfaces))
	return this
}

// GetResponse returns ResponseObject of GetAllEdgeInterfacesAPI
func (getAllAPI GetAllEdgeInterfacesAPI) GetResponse() *EdgeInterfaces {
	return getAllAPI.ResponseObject().(*EdgeInterfaces)
}
