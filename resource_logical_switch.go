package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/virtualwire"
	"log"
)

func getSingleLogicalSwitch(scopeid, name string, nsxclient *gonsx.NSXClient) (*virtualwire.VirtualWire, error) {
	getAllAPI := virtualwire.NewGetAll(scopeid)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	service := getAllAPI.GetResponse().FilterByName(name)

	if service.ObjectID == "" {
		return nil, fmt.Errorf("Not found %s", name)
	}

	return service, nil
}

func resourceLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalSwitchCreate,
		Read:   resourceLogicalSwitchRead,
		Delete: resourceLogicalSwitchDelete,

		Schema: map[string]*schema.Schema{
			"desc": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tenantid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scopeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var desc, name, tenantid, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("desc"); ok {
		desc = v.(string)
	} else {
		return fmt.Errorf("desc argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("tenantid"); ok {
		tenantid = v.(string)
	} else {
		return fmt.Errorf("tenantid argument is required")
	}

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] virtualwire.NewCreate(%s, %s, %s, %s)", name, desc, tenantid, scopeid))
	createAPI := virtualwire.NewCreate(name, desc, tenantid, scopeid)
	nsxclient.Do(createAPI)

	if createAPI.StatusCode() != 201 {
		return errors.New(createAPI.GetResponse())
	}

	// If we go here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	d.SetId(createAPI.GetResponse())
	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] virtualwire.NewGetAll(%s)", scopeid))
	api := virtualwire.NewGetAll(scopeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	logicalSwitchObject, err := getSingleLogicalSwitch(scopeid, name, nsxclient)

	if err != nil {
		return err
	}

	id := logicalSwitchObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}

	return nil
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] virtualwire.NewGetAll(%s)", scopeid))
	api := virtualwire.NewGetAll(scopeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	logicalSwitchObject, err := getSingleLogicalSwitch(scopeid, name, nsxclient)
	id := logicalSwitchObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := virtualwire.NewDelete(id)
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
