package dhcprelay

import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type DeleteDhcpRelayApi struct {
	*api.BaseApi
}

func NewDelete(edgeId string) *DeleteDhcpRelayApi {
	this := new(DeleteDhcpRelayApi)
	this.BaseApi = api.NewBaseApi(http.MethodDelete, "/api/4.0/edges/"+ edgeId + "/dhcp/config/relay", nil, nil)
	return this
}
