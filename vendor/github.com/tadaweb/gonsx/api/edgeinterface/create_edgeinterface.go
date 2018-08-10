package edgeinterface

import (
	"github.com/tadaweb/gonsx/api"
	"net/http"
)

// CreateEdgeInterfaceAPI struct
type CreateEdgeInterfaceAPI struct {
	*api.BaseAPI
}

// NewCreate returns a new object of CreateEdgeInterfaceAPI
func NewCreate(edgeInterfaceList *EdgeInterfaces, edgeID string) *CreateEdgeInterfaceAPI {

	this := new(CreateEdgeInterfaceAPI)

	this.BaseAPI = api.NewBaseAPI(http.MethodPost, "/api/4.0/edges/"+edgeID+"/interfaces/?action=patch", edgeInterfaceList, new(EdgeInterfaces))
	return this
}

// GetResponse returns the ResponseObject.
func (createAPI CreateEdgeInterfaceAPI) GetResponse() *EdgeInterfaces {
	return createAPI.ResponseObject().(*EdgeInterfaces)
}
