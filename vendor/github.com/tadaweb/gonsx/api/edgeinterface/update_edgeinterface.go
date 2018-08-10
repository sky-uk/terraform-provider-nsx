package edgeinterface

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
	"strconv"
)

// UpdateEdgeInterfaceAPI struct
type UpdateEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewUpdate returns a new delete method object of UpdateEdgeInterfaceAPI
func NewUpdate(edgeID string, index int, edge EdgeInterface) *UpdateEdgeInterfaceAPI {
	this := new(UpdateEdgeInterfaceAPI)
	this.BaseAPI = api.NewBaseAPI(http.MethodPut, "/api/4.0/edges/"+edgeID+"/interfaces/"+strconv.Itoa(index), edge, nil)
	return this
}
