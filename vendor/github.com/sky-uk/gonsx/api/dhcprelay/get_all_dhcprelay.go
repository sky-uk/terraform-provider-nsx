package dhcprelay

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type GetAllDhcpRelaysApi struct {
	*api.BaseApi
}

func NewGetAll(edgeId string) *GetAllDhcpRelaysApi {
	this := new(GetAllDhcpRelaysApi)
	this.BaseApi = api.NewBaseApi(http.MethodGet, "/api/4.0/edges/"+edgeId+"/dhcp/config/relay", nil, new(DhcpRelay))
	return this
}

func (this GetAllDhcpRelaysApi) GetResponse() *DhcpRelay {
	return this.ResponseObject().(*DhcpRelay)
}
