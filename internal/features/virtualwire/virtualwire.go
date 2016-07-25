package virtualwire

import (
	"github.com/sky-uk/gonsx/api/virtualwire"
	"github.com/sky-uk/terraform-provider-nsx/internal"
)

func deleteAllVirtualWiresWithNameInScope(name string, scope string) {
	for {
		getApi := virtualwire.NewGetAll(scope)
		err := Client.Do(getApi)
		internal.CheckError(err)
		objectId := getApi.GetResponse().FilterByName(name).ObjectID
		if objectId == "" {
			break
		}
		deleteApi := virtualwire.NewDelete(objectId)
		err = Client.Do(deleteApi)
		internal.CheckError(err)
	}
}

func getVirtualWire(name string, scope string) *virtualwire.VirtualWire {
	getApi := virtualwire.NewGetAll(scope)
	err := Client.Do(getApi)
	internal.CheckError(err)
	return getApi.GetResponse().FilterByName(name)
}

func createVirtualWire(name string, description string, tenant string, scope string) {
	createApi := virtualwire.NewCreate(name, description, tenant, scope)
	err := Client.Do(createApi)
	internal.CheckError(err)
}
