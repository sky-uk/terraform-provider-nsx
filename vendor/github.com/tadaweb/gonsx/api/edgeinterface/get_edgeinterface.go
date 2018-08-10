package edgeinterface

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
	"strconv"
)

// GetEdgeInterfaceAPI base object.
type GetEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewGet returns the api object of GetEdgeInterfaceAPI
func NewGet(edgeID string, index int) *GetEdgeInterfaceAPI {
	this := new(GetEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodGet, "/api/4.0/edges/"+edgeID+"/interfaces/"+strconv.Itoa(index), nil, new(EdgeInterface))
	return this
}

// GetResponse returns ResponseObject of GetEdgeInterfaceAPI
func (getAPI GetEdgeInterfaceAPI) GetResponse() EdgeInterface {
	return *getAPI.ResponseObject().(*EdgeInterface)
}
