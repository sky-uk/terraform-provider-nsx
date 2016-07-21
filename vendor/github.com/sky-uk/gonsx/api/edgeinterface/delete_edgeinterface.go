package edgeinterface

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteEdgeInterfaceAPI struct
type DeleteEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new delete method object of DeleteEdgeInterfaceAPI
func NewDelete(interfaceIndex, edgeID string) *DeleteEdgeInterfaceAPI {
	this := new(DeleteEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/interfaces/?index="+interfaceIndex, nil, nil)
	return this
}
