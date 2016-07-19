package virtualwire

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type DeleteVirtualWiresApi struct {
	*api.BaseApi
}

func NewDelete(virtualWireId string) *DeleteVirtualWiresApi {
	this := new(DeleteVirtualWiresApi)
	this.BaseApi = api.NewBaseApi(http.MethodDelete, "/api/2.0/vdn/virtualwires/" + virtualWireId, nil, nil)
	return this
}
