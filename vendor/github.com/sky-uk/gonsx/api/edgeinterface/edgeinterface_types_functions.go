package edgeinterface

import "fmt"

func (s EdgeInterfaces) String() string {
	return fmt.Sprintf("%s", s.Interfaces)
}

//func (s EdgeInterface) String() string {
//	return fmt.Sprintf("index: %s, name: %s", s.Index, s.Name)
//}

func (v EdgeInterfaces) FilterByName(name string) *EdgeInterface {
	var edgeInterfaceFound EdgeInterface
	for _, edgeInterface := range v.Interfaces {
		if edgeInterface.Name == name {
			edgeInterfaceFound = edgeInterface
			break
		}
	}
	return &edgeInterfaceFound
}
