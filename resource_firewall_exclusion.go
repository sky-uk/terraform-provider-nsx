package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/firewallexclusion"
	"log"
)

func getMember(moid string, nsxclient *gonsx.NSXClient) (*firewallexclusion.Member, error) {
	getAllAPI := firewallexclusion.NewGetAll()
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	member := getAllAPI.GetResponse().FilterByMOID(moid)

	if member.MOID == "" {
		return nil, nil
	}

	return member, nil
}

func resourceFirewallExclusion() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallExclusionCreate,
		Read:   resourceFirewallExclusionRead,
		Delete: resourceFirewallExclusionDelete,

		Schema: map[string]*schema.Schema{
			"moid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFirewallExclusionCreate(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var moid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] member.NewCreate(%s)", moid))
	createAPI := firewallexclusion.NewCreate(moid)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if createAPI.StatusCode() != 200 {
		return fmt.Errorf("%s", createAPI.ResponseObject())
	}

	// If we get here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	d.SetId(moid)
	return resourceFirewallExclusionRead(d, meta)
}

func resourceFirewallExclusionRead(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var moid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// of firewall exclusions
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByMOID(\"%s\")", moid))
	memberObject, err := getMember(moid, nsxclient)
	if err != nil {
		return err
	}

	// If the resource has been removed manually, notify Terraform of this fact.
	if memberObject == nil {
		d.SetId("")
	}

	return nil
}

func resourceFirewallExclusionDelete(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	var moid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("moid"); ok {
		moid = v.(string)
	} else {
		return fmt.Errorf("moid argument is required")
	}

	// Gather all the resources that are associated with the specified
	// moid.
	log.Printf(fmt.Sprintf("[DEBUG] member.NewGetAll(%s)", moid))
	api := firewallexclusion.NewGetAll()
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the moid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByMOID(\"%s\").MOID", moid))
	memberObject, err := getMember(moid, nsxclient)
	id := memberObject.MOID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := firewallexclusion.NewDelete(id)
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
