package virtualwire

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type GetVirtualWireApi struct {
	*api.BaseApi
}

func NewGet(id string) *GetVirtualWireApi {
	this := new(GetVirtualWireApi)
	this.BaseApi = api.NewBaseApi(http.MethodGet, "/api/2.0/vdn/virtualwires/" + id, nil, new(VirtualWire))
	return this
}

func (this GetVirtualWireApi) GetResponse() *VirtualWire {
	return this.ResponseObject().(*VirtualWire)
}
