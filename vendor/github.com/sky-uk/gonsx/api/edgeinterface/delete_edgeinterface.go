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
func NewDelete(interfaceIndex int, edgeID string) *DeleteEdgeInterfaceAPI {
	this := new(DeleteEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodDelete, "/api/4.0/edges/"+edgeID+"/interfaces/?index="+strconv.Itoa(interfaceIndex), nil, nil)
	return this
}
