package main

import (
	"errors"
	"fmt"
        "github.com/sky-uk/gonsx"
        "github.com/sky-uk/gonsx/api/dhcprelay"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func getAllDhcpRelays(edgeId string, nsxclient *gonsx.NSXClient) (*dhcprelay.DhcpRelay, error) {
        //
        // Get All DHCP Relay agents.
        //
        api := dhcprelay.NewGetAll(edgeId)
        // make the api call with nsxclient
        err := nsxclient.Do(api)
        // check if we err otherwise read response.
        if err != nil {
                fmt.Println("Error:", err)
                return nil, err
        } else {
                fmt.Println("Get All Response: ", api.GetResponse())
                return api.GetResponse(), nil
        }
}

func resourceDHCPRelay() *schema.Resource {
	return &schema.Resource{
		Create: resourceDHCPRelayCreate,
		Read:   resourceDHCPRelayRead,
		Delete: resourceDHCPRelayDelete,

		Schema: map[string]*schema.Schema{
			"edgeid": {
				  Type:     schema.TypeString,
				  Required: true,
				  ForceNew: true,
			},

			"vnicindex": {
				     Type:     schema.TypeString,
				     Required: true,
				     ForceNew: true,
			},

			"giaddress": {
				     Type:     schema.TypeString,
				     Required: true,
				     ForceNew: true,
			},
		},
	}
}

func resourceDHCPRelayCreate(d *schema.ResourceData, m interface{}) error {
        nsxclient := m.(*gonsx.NSXClient)
        var edgeid, vnicindex, giaddress string

        // Gather the attributes for the resource.
        if v, ok := d.GetOk("edgeid"); ok {
                edgeid = v.(string)
        } else {
                return fmt.Errorf("edgeid argument is required")
        }

        if v, ok := d.GetOk("vnicindex"); ok {
                vnicindex = v.(string)
        } else {
                return fmt.Errorf("vnicindex argument is required")
        }

        if v, ok := d.GetOk("giaddress"); ok {
                giaddress = v.(string)
        } else {
                return fmt.Errorf("giaddress argument is required")
        }

        // Create the API, use it and check for errors.
        currentDHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)

        if err != nil {
                return err
        }

        newRelayAgent := dhcprelay.RelayAgent{VnicIndex: vnicindex, GiAddress: giaddress}

        // If we go here, everything is OK.  Set the ID for the Terraform state
        // and return the response from the READ method.
        d.SetId(createAPI.GetResponse())
        return resourceDHCPRelayRead(d, m)
}

func resourceDHCPRelayRead(d *schema.ResourceData, m interface{}) error {
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
	id := (api.GetResponse().FilterByName(name).ObjectID)
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}

	return nil
}

func resourceDHCPRelayDelete(d *schema.ResourceData, m interface{}) error {
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
	id := (api.GetResponse().FilterByName(name).ObjectID)
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
