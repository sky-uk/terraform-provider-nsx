package virtualwire

import "fmt"

func (s VirtualWires) String() string {
	return fmt.Sprintf("%s", s.DataPage)
}

func (s VirtualWire) String() string {
	return fmt.Sprintf("id: %s, name: %s", s.ObjectID, s.Name)
}

// FilterByName returns a single virtualWire object if it matches the name from VirtualWires
func (s VirtualWires) FilterByName(name string) *VirtualWire {
	var virtualWireFound VirtualWire
	for _, virtualWire := range s.DataPage.VirtualWires {
		if virtualWire.Name == name {
			virtualWireFound = virtualWire
			break
		}
	}
	return &virtualWireFound
}
