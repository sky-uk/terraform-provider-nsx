package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitypolicy"
	"log"
)

func getSingleSecurityPolicy(name string, nsxclient *gonsx.NSXClient) (*securitypolicy.SecurityPolicy, error) {
	getAllAPI := securitypolicy.NewGetAll()
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	log.Printf(fmt.Sprintf("[DEBUG] getAllAPI.GetResponse().FilterByName(\"%s\").ObjectID", name))
	securityPolicy := getAllAPI.GetResponse().FilterByName(name)

	return securityPolicy, nil
}

func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityPolicyCreate,
		Read:   resourceSecurityPolicyRead,
		Delete: resourceSecurityPolicyDelete,
		Update: resourceSecurityPolicyUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"precedence": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"securitygroups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var name, description, precedence string
	var securitygroups []string
	var actions []securitypolicy.Action

	// Gather the attributes for the resource.

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("precedence"); ok {
		precedence = v.(string)
	} else {
		return fmt.Errorf("precedence argument is required")
	}

	if v, ok := d.GetOk("description"); ok {
		description = v.(string)
	} else {
		description = name
	}

	if v, ok := d.GetOk("securitygroups"); ok {
		list := v.([]interface{})
		securitygroups = make([]string, len(list))
		for i, value := range list {
			groupID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitygroups")
			}
			securitygroups[i] = groupID
		}
	} else {
		securitygroups = make([]string, 0)
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitypolicy.NewCreate(%s, %s, %s, %s, %v)", name, precedence, description, securitygroups, actions))
	createAPI := securitypolicy.NewCreate(name, precedence, description, securitygroups, actions)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error creating security policy: %v", err)
	}

	if createAPI.StatusCode() != 201 {
		return fmt.Errorf("%s", createAPI.ResponseObject())
	}

	d.SetId(createAPI.GetResponse())
	return resourceSecurityPolicyRead(d, meta)
}

func resourceSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	securityPolicyObject, err := getSingleSecurityPolicy(name, nsxclient)
	if err != nil {
		return err
	}
	id := securityPolicyObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}
	return nil
}

func resourceSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	securityPolicyObject, err := getSingleSecurityPolicy(name, nsxclient)
	if err != nil {
		return err
	}
	id := securityPolicyObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	// FIXME: we need to get terraform force call and pass it here.
	deleteAPI := securitypolicy.NewDelete(id, false)
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

func resourceSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	// flag if changes have to be applied
	hasChanges := false

	nsxclient := meta.(*gonsx.NSXClient)
	var name string
	var securitygroups []string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	securityPolicyToChange, err := getSingleSecurityPolicy(name, nsxclient)
	if err != nil {
		return err
	}

	id := securityPolicyToChange.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource is not found, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// Update resource properties.
	if d.HasChange("description") {
		hasChanges = true
		securityPolicyToChange.Description = d.Get("description").(string)
	}

	if d.HasChange("precedence") {
		hasChanges = true
		securityPolicyToChange.Precedence = d.Get("precedence").(string)
	}

	if d.HasChange("securitygroups") {
		hasChanges = true

		// TODO: fix this when API is updated, for now we remove everything first.
		for _, securityGroup := range securityPolicyToChange.SecurityGroupBinding {
			securityPolicyToChange.RemoveSecurityGroupBinding(securityGroup.ObjectID)
		}

		if v, ok := d.GetOk("securitygroups"); ok {
			list := v.([]interface{})

			securitygroups = make([]string, len(list))
			for i, value := range list {
				groupID, ok := value.(string)
				if !ok {
					return fmt.Errorf("empty element found in securitygroups")
				}
				securitygroups[i] = groupID
			}
		} else {
			securitygroups = make([]string, 0)
		}

		for _, securityGroupID := range securitygroups {
			securityPolicyToChange.AddSecurityGroupBinding(securityGroupID)
		}
	}

	// do nothing if there are no changes
	if !hasChanges {
		return nil
	}

	securityPolicyToChange.Revision += securityPolicyToChange.Revision
	updateAPI := securitypolicy.NewUpdate(id, securityPolicyToChange)
	err = nsxclient.Do(updateAPI)

	if err != nil {
		return err
	}

	return resourceSecurityPolicyRead(d, meta)

}
