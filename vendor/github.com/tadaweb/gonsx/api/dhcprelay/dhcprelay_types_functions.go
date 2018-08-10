package dhcprelay

import "fmt"

func (s RelayServer) String() string {
	return fmt.Sprintf("DhcpRelayServer ipAddress:%s.", s.IPAddress)
}

func (d RelayAgent) String() string {
	return fmt.Sprintf("DhcpRelayAgent VnicIndex:%s, GiAddress:%s.", d.VnicIndex, d.GiAddress)
}

// FilterByIPAddress - Filters the DhcpRelay->RelayAgents with provided IP Address.
func (v DhcpRelay) FilterByIPAddress(ipAddress string) *RelayAgent {
	var relayAgentFound RelayAgent
	for _, relayAgent := range v.RelayAgents {
		if relayAgent.GiAddress == ipAddress {
			relayAgentFound = relayAgent
			break
		}
	}
	return &relayAgentFound
}

// FilterByVnicIndex - Filters the DhcpRelay->RelayAgents with provided vnicIndex.
func (v DhcpRelay) FilterByVnicIndex(vnicIndex string) *RelayAgent {
	var relayAgentFound RelayAgent
	for _, relayAgent := range v.RelayAgents {
		if relayAgent.VnicIndex == vnicIndex {
			relayAgentFound = relayAgent
			break
		}
	}
	return &relayAgentFound
}

// RemoveByVnicIndex - Removes the relayagent from DhcpRelay->RelayAgents which relates to provided vnicIndex.
func (v DhcpRelay) RemoveByVnicIndex(vnicIndex string) *DhcpRelay {
	for idx, relayAgent := range v.RelayAgents {
		if relayAgent.VnicIndex == vnicIndex {
			v.RelayAgents = append(v.RelayAgents[:idx], v.RelayAgents[idx+1:]...)
			break
		}
	}
	return &v
}

// CheckByVnicIndex - Returns true/false depending on if vnicIndex exists in DhcpRelay->RelayAgents list.
func (v DhcpRelay) CheckByVnicIndex(vnicIndex string) bool {
	for _, relayAgent := range v.RelayAgents {
		if relayAgent.VnicIndex == vnicIndex {
			return true
		}
	}
	return false
}
