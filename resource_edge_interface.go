package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/edgeinterface"
	"log"
)

func getSingleEdgeInterface(edgeid, name string, nsxclient *gonsx.NSXClient) (*edgeinterface.EdgeInterface, error) {
	getAllAPI := edgeinterface.NewGetAll(edgeid)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	edgeinterface := getAllAPI.GetResponse().FilterByName(name)

	if edgeinterface.Index == "" {
		return nil, fmt.Errorf("Not found %s", name)
	}

	return edgeinterface, nil
}

func resourceEdgeInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceEdgeInterfaceCreate,
		Read:   resourceEdgeInterfaceRead,
		Delete: resourceEdgeInterfaceDelete,

		Schema: map[string]*schema.Schema{
			"edgeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtualwireid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"gateway": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnetmask": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interfacetype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mtu": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEdgeInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid, name, virtualwireid, gateway, subnetmask, interfacetype string
	var mtu int

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("virtualwireid"); ok {
		virtualwireid = v.(string)
	} else {
		return fmt.Errorf("virtualwireid argument is required")
	}

	if v, ok := d.GetOk("gateway"); ok {
		gateway = v.(string)
	} else {
		return fmt.Errorf("gateway argument is required")
	}

	if v, ok := d.GetOk("subnetmask"); ok {
		subnetmask = v.(string)
	} else {
		return fmt.Errorf("subnetmask argument is required")
	}

	if v, ok := d.GetOk("interfacetype"); ok {
		interfacetype = v.(string)
	} else {
		return fmt.Errorf("interfacetype argument is required")
	}

	if v, ok := d.GetOk("mtu"); ok {
		mtu = v.(int)
	} else {
		return fmt.Errorf("mtu argument is required")
	}

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] edgeinterface.NewCreate(%s, %s, %s, %s, %s, %s, %d)", edgeid, name, virtualwireid, gateway, subnetmask, interfacetype, mtu))
	createAPI := edgeinterface.NewCreate(edgeid, name, virtualwireid, gateway, subnetmask, interfacetype, mtu)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return err
	}

	if createAPI.StatusCode() != 200 {
		return fmt.Errorf("Failed to create edge interface:%s StatusCode:%s Response:%s", err, createAPI.StatusCode(), createAPI.RawResponse())
	}

	// If we go here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	id := (createAPI.GetResponse().FilterByName(name).Index)

	if id != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return resourceEdgeInterfaceRead(d, m)
}

func resourceEdgeInterfaceRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid, name string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	// Gather all the resources that are associated with the specified
	// edgeid.
	log.Printf(fmt.Sprintf("[DEBUG] edgeinterface.NewGetAll(%s)", edgeid))
	api := edgeinterface.NewGetAll(edgeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the edgeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").Index", name))
	edgeInterfaceObject, err := getSingleEdgeInterface(edgeid, name, nsxclient)
	id := edgeInterfaceObject.Index
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}

	return nil
}

func resourceEdgeInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid, name string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Gather all the resources that are associated with the specified
	// edgeid.
	log.Printf(fmt.Sprintf("[DEBUG] edgeinterface.NewGetAll(%s)", edgeid))
	api := edgeinterface.NewGetAll(edgeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the edgeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").Index", name))
	edgeInterfaceObject, err := getSingleEdgeInterface(edgeid, name, nsxclient)
	id := edgeInterfaceObject.Index
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := edgeinterface.NewDelete(id, edgeid)
	err = nsxclient.Do(deleteAPI)

	if err != nil {
		return err
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")
	log.Printf(fmt.Sprintf("[DEBUG] id %s deleted.", id))
	return nil
}
