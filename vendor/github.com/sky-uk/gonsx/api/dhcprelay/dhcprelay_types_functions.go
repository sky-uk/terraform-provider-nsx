package dhcprelay

import "fmt"

func (s RelayServer) String() string {
	return fmt.Sprintf("DhcpRelayServer ipAddress:%s", s.IpAddress)
}

func (d RelayAgent) String() string {
	return fmt.Sprintf("DhcpRelayAgent VnicIndex:%s, GiAddress:%s", d.VnicIndex, d.GiAddress)
}

func (v DhcpRelay) FilterByIpAddress(ip_address string) *RelayAgent {
	var relayAgentFound RelayAgent
	for _, relay_agent := range v.RelayAgents{
		if relay_agent.GiAddress == ip_address {
			relayAgentFound = relay_agent
			break
		}
	}
	return &relayAgentFound
}

func (v DhcpRelay) FilterByVnicIndex(vnicIndex string) *RelayAgent {
	var relayAgentFound RelayAgent
	for _, relay_agent := range v.RelayAgents{
		if relay_agent.VnicIndex == vnicIndex{
			relayAgentFound = relay_agent
			break
		}
	}
	return &relayAgentFound
}

func (v DhcpRelay) RemoveByVnicIndex(vnicIndex string) *DhcpRelay{
	for idx, relay_agent := range v.RelayAgents{
		if relay_agent.VnicIndex == vnicIndex{
			v.RelayAgents = append(v.RelayAgents[:idx], v.RelayAgents[idx+1:]...)
			break
		}
	}
	return &v
}


func (v DhcpRelay) CheckByVnicIndex(vnicIndex string) bool {
	for _, relay_agent := range v.RelayAgents{
		if relay_agent.VnicIndex == vnicIndex{
			return true
		}
	}
	return false
}
