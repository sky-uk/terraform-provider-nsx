package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/edgeinterface"
	"net/http"
	"strconv"
)

func resourceEdgeInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceEdgeInterfaceCreate,
		Read:   resourceEdgeInterfaceRead,
		Delete: resourceEdgeInterfaceDelete,
		Update: resourceEdgeInterfaceUpdate,

		Schema: map[string]*schema.Schema{
			"edgeid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"interfacetype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"isconnected": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"index": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connectedtoid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"addressgroups": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primaryaddress": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"subnetmask": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildAddressGroups(addressGroups []interface{}) []edgeinterface.AddressGroup {
	var addressGroupList []edgeinterface.AddressGroup
	for _, address := range addressGroups {
		data := address.(map[string]interface{})
		addr := edgeinterface.AddressGroup{
			PrimaryAddress: data["primaryaddress"].(string),
			SubnetMask:     data["subnetmask"].(string),
		}
		addressGroupList = append(addressGroupList, addr)
	}
	return addressGroupList
}

func resourceEdgeInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	var edge edgeinterface.EdgeInterface

	edgeid := d.Get("edgeid").(string)

	edge.Name = d.Get("name").(string)
	edge.Mtu = d.Get("mtu").(int)
	edge.Type = d.Get("interfacetype").(string)

	if v, ok := d.GetOk("isconnected"); ok {
		edge.IsConnected = v.(bool)
	} else {
		edge.IsConnected = false
	}

	if v, ok := d.GetOk("connectedtoid"); ok {
		edge.ConnectedToID = v.(string)
	}

	if v, ok := d.GetOk("addressgroups"); ok {
		edge.AddressGroups.AddressGroups = buildAddressGroups(v.([]interface{}))
	}

	requestPayload := new(edgeinterface.EdgeInterfaces)
	requestPayload.Interfaces = append(requestPayload.Interfaces, edge)

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	createAPI := edgeinterface.NewCreate(requestPayload, edgeid)
	err := nsxclient.Do(createAPI)
	if err != nil {
		return err
	}

	if createAPI.StatusCode() != http.StatusOK {
		return fmt.Errorf("Failed to create edge interface, StatusCode:%d Response:%s",
			createAPI.StatusCode(),
			createAPI.RawResponse())
	}

	edges := createAPI.GetResponse()
	setEdge(d, edges.Interfaces[0])
	d.SetId(edgeid + "_" + strconv.Itoa(edges.Interfaces[0].Index))
	return nil
}

func resourceEdgeInterfaceRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	edgeid := d.Get("edgeid").(string)
	index := d.Get("index").(int)

	api := edgeinterface.NewGet(edgeid, index)
	err := nsxclient.Do(api)
	if err != nil {
		d.SetId("")
		return nil
	}

	if api.StatusCode() != http.StatusOK {
		d.SetId("")
		return fmt.Errorf("Error getting all interfaces: Status code: %d", api.StatusCode())
	}

	edge := api.GetResponse()
	setEdge(d, edge)

	return nil
}

func setEdge(d *schema.ResourceData, edge edgeinterface.EdgeInterface) {
	d.Set("name", edge.Name)
	d.Set("label", edge.Label)
	d.Set("mtu", edge.Mtu)
	d.Set("interfacetype", edge.Type)
	d.Set("isconnected", edge.IsConnected)
	d.Set("connectedtoid", edge.ConnectedToID)
	var adrGroupList []map[string]string
	for _, adrGroup := range edge.AddressGroups.AddressGroups {
		adrgrpResource := make(map[string]string)
		adrgrpResource["primaryaddress"] = adrGroup.PrimaryAddress
		adrgrpResource["subnetmask"] = adrGroup.SubnetMask
		adrGroupList = append(adrGroupList, adrgrpResource)
	}
	d.Set("addressgroups", adrGroupList)
	d.Set("index", edge.Index)
}

func resourceEdgeInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	edgeid := d.Get("edgeid").(string)
	if index, ok := d.GetOk("index"); ok {
		nsxMutexKV.Lock(edgeid)
		defer nsxMutexKV.Unlock(edgeid)
		deleteAPI := edgeinterface.NewDelete(edgeid, index.(int))
		err := nsxclient.Do(deleteAPI)
		if err != nil {
			return err
		}
		if deleteAPI.StatusCode() == http.StatusNoContent {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error deleting interface from NSX, Status code: %d", deleteAPI.StatusCode())
	}
	return fmt.Errorf("Error deleting resource %s, index not set", d.Get("Name").(string))
}

func resourceEdgeInterfaceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	hasChanges := false

	var updatedEdge edgeinterface.EdgeInterface

	edgeid := d.Get("edgeid").(string)
	index := d.Get("index").(int)

	oldName, newName := d.GetChange("name")
	if d.HasChange("name") {
		hasChanges = true
		updatedEdge.Name = newName.(string)
	} else {
		updatedEdge.Name = oldName.(string)
	}

	oldMtu, newMtu := d.GetChange("mtu")
	if d.HasChange("mtu") {
		hasChanges = true
		updatedEdge.Mtu = newMtu.(int)
	} else {
		updatedEdge.Mtu = oldMtu.(int)
	}

	oldType, newType := d.GetChange("interfacetype")
	if d.HasChange("interfacetype") {
		hasChanges = true
		updatedEdge.Type = newType.(string)
	} else {
		updatedEdge.Type = oldType.(string)
	}

	oldIsConn, newIsConn := d.GetChange("isconnected")
	if d.HasChange("isconnected") {
		hasChanges = true
		updatedEdge.IsConnected = newIsConn.(bool)
	} else {
		updatedEdge.IsConnected = oldIsConn.(bool)
	}

	oldConnID, newConnID := d.GetChange("connectedtoid")
	if d.HasChange("connectedtoid") {
		hasChanges = true
		updatedEdge.ConnectedToID = newConnID.(string)
	} else {
		updatedEdge.ConnectedToID = oldConnID.(string)
	}

	if d.HasChange("addressgroups") {
		hasChanges = true
		v := d.Get("addressgroups")
		updatedEdge.AddressGroups.AddressGroups = buildAddressGroups(v.([]interface{}))
	}

	if hasChanges {
		updateAPI := edgeinterface.NewUpdate(edgeid, index, updatedEdge)

		nsxMutexKV.Lock(edgeid)
		defer nsxMutexKV.Unlock(edgeid)

		err := nsxclient.Do(updateAPI)
		if err != nil {
			return err
		}

		if updateAPI.StatusCode() != http.StatusNoContent {
			return fmt.Errorf("Error updating resource: status code: %d", updateAPI.StatusCode())
		}
	}

	return nil
}
