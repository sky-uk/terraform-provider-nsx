package firewallexclusion

import (
	"fmt"
)

func (a FirewallExclusions) String() string {
	return fmt.Sprintf("FirewallExclusions object, contains member objects.")
}

func (a Member) String() string {
	return fmt.Sprintf("MOID: %-20s name: %-20s.", a.MOID, a.Name)
}

// FilterByMOID returns a single member object if it matches the moid in FirewallExclusions
func (a FirewallExclusions) FilterByMOID(moid string) *Member {
	var memberFound Member
	for _, member := range a.Members {
		if member.MOID == moid {
			memberFound = member
			break
		}
	}
	return &memberFound
}

// CheckByMOID - Returns true or false depending if moid is in FirewallExclusions
func (a FirewallExclusions) CheckByMOID(moid string) bool {
	for _, member := range a.Members {
		if member.MOID == moid {
			return true
		}
	}
	return false
}
