package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitytag"
	"log"
)

func getSingleSecurityTagAttached(tagid, moid string, nsxclient *gonsx.NSXClient) (*securitytag.BasicInfo, error) {

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] securitytag.NewGetAllAttached()"))
	getAllAPI := securitytag.NewGetAllAttached(tagid)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	securityTagAttached := getAllAPI.GetResponse().FilterByIDAttached(moid)

	return securityTagAttached, nil
}

func resourceSecurityTagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityTagAttachmentCreate,
		Read:   resourceSecurityTagAttachmentRead,
		Delete: resourceSecurityTagAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"tagid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"moid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSecurityTagAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var tagid, moid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("tagid"); ok {
		tagid = v.(string)
	} else {
		return fmt.Errorf("tagid argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitytag.NewAssign(%s, %s)", tagid, moid))
	createAPI := securitytag.NewAssign(tagid, moid)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return err
	}

	if createAPI.StatusCode() != 200 {
		log.Printf(fmt.Sprintf("[DEBUG] Response %v", createAPI.ResponseObject()))
		return fmt.Errorf("Failed to attach security tag %s", tagid)
	}

	id := tagid + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if tagid != "" && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return resourceSecurityTagAttachmentRead(d, m)
}

func resourceSecurityTagAttachmentRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var tagid, moid string

	if v, ok := d.GetOk("tagid"); ok {
		tagid = v.(string)
	} else {
		return fmt.Errorf("tag argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", moid))
	securityTagAttachedObject, err := getSingleSecurityTagAttached(tagid, moid, nsxclient)

	if err != nil {
		return err
	}

	if securityTagAttachedObject.ObjectID == "" {
		return fmt.Errorf("Not found %s", moid)
	}

	id := tagid + "/" + moid
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	if tagid != "" && moid != "" {
		d.SetId(id)
	} else {
		return errors.New("Can not establish the id of the created resource")
	}

	return nil
}

func resourceSecurityTagAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var tagid, moid string

	if v, ok := d.GetOk("tagid"); ok {
		tagid = v.(string)
	} else {
		return fmt.Errorf("tag argument is required")
	}

	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", moid))
	securityTagAttached, err := getSingleSecurityTagAttached(tagid, moid, nsxclient)
	id := securityTagAttached.ObjectID

	if err != nil {
		return err
	}
	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	detachAPI := securitytag.NewDetach(tagid, moid)
	err = nsxclient.Do(detachAPI)

	if err != nil {
		return err
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")
	log.Printf(fmt.Sprintf("[DEBUG] id %s detached.", id))

	return nil
}
