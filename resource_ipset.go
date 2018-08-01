package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/ipset"
	"log"
)

func getSingleIPSet(scopeID, name string, nsxclient *gonsx.NSXClient) (*ipset.IPSet, error) {
	getAllAPI := ipset.NewGetAll(scopeID)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	listIPSets := getAllAPI.GetResponse()
	for _, service := range listIPSets.IPSets {
		if service.Name == name {
			return &service, nil
		}
	}
	return nil, nil
}

func resourceIPSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPSetCreate,
		Read:   resourceIPSetRead,
		Update: resourceIPSetUpdate,
		Delete: resourceIPSetDelete,

		Schema: map[string]*schema.Schema{
			"scopeid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIPSetCreate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name, description string
	var ipSet ipset.IPSet
	var err error

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	} else {
		description = ""
	}

	if v, ok := d.GetOk("value"); ok {
		log.Printf(fmt.Sprintf("[DEBUG] IPSet create : %+v", v))
		ipSet = ipset.IPSet{Value: v.(string), Name: name, Description: description}
		if err != nil {
			return err
		}
	} else {
		return errors.New("value list is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] ipset.NewCreate(%s, %s, %v", scopeid, name, ipSet))
	createAPI := ipset.NewCreate(scopeid, &ipSet)
	err = nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error creating ipset: %v", err)
	}

	if createAPI.StatusCode() != 201 {
		return fmt.Errorf("%s", createAPI.ResponseObject())
	}

	d.SetId(createAPI.GetResponse())
	return resourceIPSetRead(d, m)
}

func resourceIPSetRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name string
	var err error

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	ipSetObject, err := getSingleIPSet(scopeid, name, nsxclient)
	if err != nil {
		return err
	}

	id := ipSetObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceIPSetUpdate(d *schema.ResourceData, m interface{}) error {

	var scopeid string
	var err error

	nsxclient := m.(*gonsx.NSXClient)
	hasChanges := false

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] ipset.NewGetAll(%s)", scopeid))
	oldName, newName := d.GetChange("name")
	ipSetObject, err := getSingleIPSet(scopeid, oldName.(string), nsxclient)
	id := ipSetObject.ObjectID

	if d.HasChange("name") {
		hasChanges = true
		ipSetObject.Name = newName.(string)
		log.Printf(fmt.Sprintf("[DEBUG] Changing name of ipset from %s to %s", oldName.(string), newName.(string)))
	}

	if d.HasChange("description") {
		oldDescription, newDescription := d.GetChange("description")
		hasChanges = true
		ipSetObject.Description = newDescription.(string)
		log.Printf(fmt.Sprintf("[DEBUG] Changing description of ipset from %s to %s", oldDescription.(string), newDescription.(string)))
	}

	if d.HasChange("value") {
		if v, ok := d.GetOk("value"); ok {
			ipSetObject.Value = v.(string)
			if err != nil {
				return err
			}
		}
		hasChanges = true
	}

	if hasChanges {
		updateAPI := ipset.NewUpdate(id, ipSetObject)
		err = nsxclient.Do(updateAPI)
		if err != nil {
			log.Printf(fmt.Sprintf("[DEBUG] Error updating ipset: %s", err))
		}
	}
	return resourceIPSetRead(d, m)
}

func resourceIPSetDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
	}

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] ipset.NewGetAll(%s)", scopeid))
	api := ipset.NewGetAll(scopeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	ipSetObject, err := getSingleIPSet(scopeid, name, nsxclient)
	id := ipSetObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := ipset.NewDelete(id)
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
