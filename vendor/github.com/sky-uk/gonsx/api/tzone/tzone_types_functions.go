package tzone

import "fmt"

func (s NetworkScopeList) String() string {
	return fmt.Sprintf("%s", s.NetworkScopeList)
}

func (s NetworkScope) String() string {
	return fmt.Sprintf("id: %s, name: %s", s.ObjectID, s.Name)
}

// FilterByName filters NetworkScopeList object given a name and returns the found object.
func (s NetworkScopeList) FilterByName(name string) *NetworkScope {
	var networkScopeFound NetworkScope
	for _, networkScope := range s.NetworkScopeList {
		if networkScope.Name == name {
			networkScopeFound = networkScope
			break
		}
	}
	return &networkScopeFound
}
