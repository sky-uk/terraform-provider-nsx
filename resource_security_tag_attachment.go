package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitytag"
	"log"
)

func resourceSecurityTagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityTagAttachmentCreate,
		Read:   resourceSecurityTagAttachmentRead,
		Delete: resourceSecurityTagAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"moid": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
			},
		},
	}
}

func resourceSecurityTagAttachmentCreate (d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, moid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("name"); ok{
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("moid"); ok{
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] getSingleSecurityTag(%s)", name))
	securityTagObject, err := getSingleSecurityTag(name, nsxclient)
	tagID := securityTagObject .ObjectID

	if err != nil {
		return err
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitytag.NewAssign(%s, %s)", tagID, moid))
	createAPI := securitytag.NewAssign(name, moid)
	err = nsxclient.Do(createAPI)

	if err != nil {
		return err
	}

	if createAPI.StatusCode() != 200 {
		return fmt.Errorf("Failed to attach security tag %s", name)
	}

	id := tagID + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if tagID != "" || moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return resourceSecurityTagAttachmentRead(d, m)
}

func resourceSecurityTagAttachmentRead (d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityTagAttachmentDelete (d *schema.ResourceData, m interface{}) error {
	return nil
}