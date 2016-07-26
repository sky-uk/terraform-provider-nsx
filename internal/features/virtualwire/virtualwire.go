package virtualwire

import (
	"github.com/sky-uk/gonsx/api/virtualwire"
	"github.com/sky-uk/terraform-provider-nsx/internal"
)

func deleteAllVirtualWiresWithNameInScope(name string, scope string) {
	for {
		getAPI := virtualwire.NewGetAll(scope)
		err := Client.Do(getAPI)
		internal.CheckError(err)
		objectID := getAPI.GetResponse().FilterByName(name).ObjectID
		if objectID == "" {
			break
		}
		deleteAPI := virtualwire.NewDelete(objectID)
		err = Client.Do(deleteAPI)
		internal.CheckError(err)
	}
}

func getVirtualWire(name string, scope string) *virtualwire.VirtualWire {
	getAPI := virtualwire.NewGetAll(scope)
	err := Client.Do(getAPI)
	internal.CheckError(err)
	return getAPI.GetResponse().FilterByName(name)
}

func createVirtualWire(name string, description string, tenant string, scope string) {
	createAPI := virtualwire.NewCreate(name, description, tenant, scope)
	err := Client.Do(createAPI)
	internal.CheckError(err)
}
