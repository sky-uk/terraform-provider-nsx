package edgeinterface

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
	"strconv"
)

// DeleteEdgeInterfaceAPI struct
type DeleteEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewDelete returns a new delete method object of DeleteEdgeInterfaceAPI
func NewDelete(edgeID string, index int) *DeleteEdgeInterfaceAPI {
	this := new(DeleteEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/interfaces/"+strconv.Itoa(index), nil, nil)
	return this
}
