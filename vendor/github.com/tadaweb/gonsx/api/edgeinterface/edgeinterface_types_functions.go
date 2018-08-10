package edgeinterface

import "fmt"

func (s EdgeInterfaces) String() string {
	return fmt.Sprintf("%v", s.Interfaces)
}

//func (s EdgeInterface) String() string {
//	return fmt.Sprintf("index: %s, name: %s", s.Index, s.Name)
//}

// FilterByName filters EdgeInterfaces list by given name of interface and returns
// the found interface object.
func (s EdgeInterfaces) FilterByName(name string) *EdgeInterface {
	var edgeInterfaceFound EdgeInterface
	for _, edgeInterface := range s.Interfaces {
		if edgeInterface.Name == name {
			edgeInterfaceFound = edgeInterface
			break
		}
	}
	return &edgeInterfaceFound
}
