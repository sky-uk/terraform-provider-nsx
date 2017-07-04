package edgeinterface

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

// DeleteAllEdgeInterfaceAPI struct
type DeleteAllEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewDeleteAll returns a new delete method object of DeleteEdgeInterfaceAPI
func NewDeleteAll(edgeID string) *DeleteAllEdgeInterfaceAPI {
	this := new(DeleteAllEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/interfaces", nil, nil)
	return this
}
