package dhcprelay


import (
	"github.com/sky-uk/gonsx/api"
	"net/http"
)

type UpdateDhcpRelayApi struct {
	*api.BaseApi
}

func NewUpdate(dhcpIpAddress, edgeId string, relayAgentslist []RelayAgent) *UpdateDhcpRelayApi {
	this := new(UpdateDhcpRelayApi)
	requestPayload := new(DhcpRelay)
	requestPayload.RelayServer.IpAddress = dhcpIpAddress
	requestPayload.RelayAgents = relayAgentslist

	this.BaseApi = api.NewBaseApi(http.MethodPut, "/api/4.0/edges/" + edgeId +"/dhcp/config/relay", requestPayload, new(string))
	return this
}


func (this UpdateDhcpRelayApi) GetResponse() string{
	return this.ResponseObject().(string)
}

