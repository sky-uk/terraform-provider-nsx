package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/service"
	"log"
)

func getAllServices(scopeID string, nsxclient *gonsx.NSXClient) (*service.Service, error) {
	getAllAPI := service.NewGetAll(scopeid)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	return getAllAPI, nil
}

func resourceService() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceCreate,
		Read:   resourceServiceRead,
		Delete: resourceServiceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scopeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"desc": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"proto": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ports": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceServiceCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid, desc, proto, ports string

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

	if v, ok := d.GetOk("desc"); ok {
		desc = v.(string)
	} else {
		return fmt.Errorf("desc argument is required")
	}

	if v, ok := d.GetOk("proto"); ok {
		proto = v.(string)
	} else {
		return fmt.Errorf("proto argument is required")
	}

	if v, ok := d.GetOk("ports"); ok {
		ports = v.(string)
	} else {
		return fmt.Errorf("ports argument is required")
	}

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] service.NewGetAll(%s)", scopeid))
	getAllAPI := service.NewGetAll(scopeid)
        err := nsxclient.Do(getAllAPI)

	if err != nil {
		return fmt.Errorf("Error:", err)
	}

        // check the status code and proceed accordingly.
        if getAllAPI.StatusCode() != 200 {
                fmt.Println("Status code:", getAllAPI.StatusCode())
                fmt.Println("Response: ", getAllAPI.ResponseObject())
        }

        

	// If we get here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	d.SetId(name)
	return resourceServiceRead(d, m)
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid, vnicindex string

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

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %s)", edgeid, nsxclient))
	currentService, err := getAllDhcpRelays(edgeid, nsxclient)

	if err != nil {
		return fmt.Errorf("Error:", err)
	}

	if !currentService.CheckByVnicIndex(vnicindex) {
		d.SetId("")
	}

	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid, vnicindex string

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

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %s)", edgeid, nsxclient))
	currentService, err := getAllDhcpRelays(edgeid, nsxclient)

	if err != nil {
		return fmt.Errorf("Error:", err)
	}

	// Check to see if an entry with the vnicindex exists at all.  If
	// not, assume it has been deleted manually and notify Terraform
	// and exit gracefully.
	if !currentService.CheckByVnicIndex(vnicindex) {
		d.SetId("")
		return nil
	}

	if currentService.CheckByVnicIndex(vnicindex) && (len(currentService.RelayAgents) == 1) {
		deleteAPI := dhcprelay.NewDelete(edgeid)
		err = nsxclient.Do(deleteAPI)
		if err != nil {
			return fmt.Errorf("Error:", err)
		}

		log.Println("DHCP Relay agent deleted.")
	} else {
		// if we got more than one relay agents, then we have to call update after removing
		// the entry we want to remove.
		log.Println("There are other DHCP Relay agents, only removing single entry with update.")
		newRelayAgentsList := currentService.RemoveByVnicIndex(vnicindex).RelayAgents

		updateAPI := dhcprelay.NewUpdate(currentService.RelayServer.IPAddress, edgeid, newRelayAgentsList)
		err = nsxclient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("Error:", err)
		} else if updateAPI.StatusCode() != 204 {
			return fmt.Errorf(updateAPI.GetResponse())
		} else {
			log.Printf("Updated DHCP Relay - %s", updateAPI.GetResponse())
		}
	}

	return nil
}
