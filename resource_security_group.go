package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitygroup"
	"log"
)

func getSingleSecurityGroup(scopeID, name string, nsxclient *gonsx.NSXClient) (*securitygroup.SecurityGroup, error) {
	getAllAPI := securitygroup.NewGetAll(scopeID)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	securityGroup := getAllAPI.GetResponse().FilterByName(name)
	return securityGroup, nil
}

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,

		Schema: map[string]*schema.Schema{

			"scopeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"setoperator": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"criteriaoperator": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"criteriakey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"criteriavalue": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"criteria": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name, setoperator, criteriaoperator, criteriakey, criteriavalue, criteria string

	// Gather the attributes for the resource.

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("setoperator"); ok {
		setoperator = v.(string)
	} else {
		return fmt.Errorf("setoperator argument is required")
	}

	if v, ok := d.GetOk("criteriaoperator"); ok {
		criteriaoperator = v.(string)
	} else {
		return fmt.Errorf("criteriaoperator argument is required")
	}

	if v, ok := d.GetOk("criteriakey"); ok {
		criteriakey = v.(string)
	} else {
		return fmt.Errorf("criteriakey argument is required")
	}

	if v, ok := d.GetOk("criteriavalue"); ok {
		criteriavalue = v.(string)
	} else {
		return fmt.Errorf("criteriavalue argument is required")
	}

	if v, ok := d.GetOk("criteria"); ok {
		criteria = v.(string)
	} else {
		return fmt.Errorf("criteria argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewCreate(%s, %s, %s, %s, %s, %s, %s)", scopeid, name, setoperator, criteriaoperator, criteriakey, criteriavalue, criteria))
	createAPI := securitygroup.NewCreate(scopeid, name, setoperator, criteriaoperator, criteriakey, criteriavalue, criteria)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error creating security group: %v", err)
	}

	if createAPI.StatusCode() != 201 {
		return fmt.Errorf("%s", createAPI.ResponseObject())
	}

	d.SetId(createAPI.GetResponse())
	return resourceSecurityGroupRead(d, m)
}

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name string

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	securityGroupObject, err := getSingleSecurityGroup(scopeid, name, nsxclient)
	if err != nil {
		return err
	}

	id := securityGroupObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}

	return nil
}


func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	//var scopeid, name, setoperator, criteriaoperator, criteriakey, criteriavalue, criteria string
	//hasChanges := false
	securityGroupToUpdate := "TEST"

	getAllAPI := securitygroup.NewGetAll("globalroot-0")
	err := nsxclient.Do(getAllAPI)
	if err != nil {
		return fmt.Errorf("Error: ", err)
	}

	if getAllAPI.StatusCode() != 200 {
		return fmt.Errorf("Status code: %v, Response: %v\n", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	if getAllAPI.StatusCode() == 200 {
		securityGroupToSearch := d.Get("name").(string)
		securityGroupToUpdate := getAllAPI.GetResponse().FilterByName(securityGroupToSearch)
		if securityGroupToUpdate.ObjectID == "" {
			return fmt.Errorf("ERROR: Service update - service %s not found", d.Get("name"))
		}
	} else {
		return fmt.Errorf("ERROR: Service update returned bad http response code for %s.", d.Get("name"))
	}


	fmt.Printf("Security Group to update is: %s", securityGroupToUpdate)

	return nil
	//return resourceSecurityGroupRead(d, m)
}


func resourceSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewGetAll(%s)", scopeid))
	api := securitygroup.NewGetAll(scopeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	securityGroupObject, err := getSingleSecurityGroup(scopeid, name, nsxclient)
	id := securityGroupObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := securitygroup.NewDelete(id)
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
