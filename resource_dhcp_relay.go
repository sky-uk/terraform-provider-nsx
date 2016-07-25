package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/dhcprelay"
	"log"
)

func getAllDhcpRelays(edgeID string, nsxclient *gonsx.NSXClient) (*dhcprelay.DhcpRelay, error) {
	//
	// Get All DHCP Relay agents.
	//
	api := dhcprelay.NewGetAll(edgeID)
	// make the api call with nsxclient
	err := nsxclient.Do(api)
	// check if we err otherwise read response.
	if err != nil {
		return nil, err
	}

	log.Println("Get All Response: ", api.GetResponse())
	return api.GetResponse(), nil
}

func resourceDHCPRelay() *schema.Resource {
	return &schema.Resource{
		Create: resourceDHCPRelayCreate,
		Read:   resourceDHCPRelayRead,
		Delete: resourceDHCPRelayDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

			"dhcpserverip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDHCPRelayCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, edgeid, vnicindex, giaddress, dhcpserverip string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

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

	if v, ok := d.GetOk("dhcpserverip"); ok {
		dhcpserverip = v.(string)
	} else {
		return fmt.Errorf("dhcpserverip argument is required")
	}

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Create the API, use it and check for errors.
	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %v)", edgeid, nsxclient))
	currentDHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.RelayAgent(%s, %s)", vnicindex, giaddress))
	newRelayAgent := dhcprelay.RelayAgent{VnicIndex: vnicindex, GiAddress: giaddress}

	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.append(%s, %s)", currentDHCPRelay.RelayAgents, newRelayAgent))
	newRelayAgentsList := append(currentDHCPRelay.RelayAgents, newRelayAgent)

	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.NewUpdate(%s, %s, %s)", dhcpserverip, edgeid, newRelayAgentsList))
	updateAPI := dhcprelay.NewUpdate(dhcpserverip, edgeid, newRelayAgentsList)

	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update the DHCP relay %s", updateAPI.GetResponse())
	}

	// If we get here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	d.SetId(name)
	return resourceDHCPRelayRead(d, m)
}

func resourceDHCPRelayRead(d *schema.ResourceData, m interface{}) error {
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
	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %v)", edgeid, nsxclient))
	currentDHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if !currentDHCPRelay.CheckByVnicIndex(vnicindex) {
		d.SetId("")
	}

	return nil
}

func resourceDHCPRelayDelete(d *schema.ResourceData, m interface{}) error {
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
	log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %v)", edgeid, nsxclient))
	currentDHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	// Check to see if an entry with the vnicindex exists at all.  If
	// not, assume it has been deleted manually and notify Terraform
	// and exit gracefully.
	if !currentDHCPRelay.CheckByVnicIndex(vnicindex) {
		d.SetId("")
		return nil
	}

	if currentDHCPRelay.CheckByVnicIndex(vnicindex) && (len(currentDHCPRelay.RelayAgents) == 1) {
		deleteAPI := dhcprelay.NewDelete(edgeid)
		err = nsxclient.Do(deleteAPI)
		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		log.Println("DHCP Relay agent deleted.")
	} else {
		// if we got more than one relay agents, then we have to call update after removing
		// the entry we want to remove.
		log.Println("There are other DHCP Relay agents, only removing single entry with update.")
		newRelayAgentsList := currentDHCPRelay.RemoveByVnicIndex(vnicindex).RelayAgents

		updateAPI := dhcprelay.NewUpdate(currentDHCPRelay.RelayServer.IPAddress, edgeid, newRelayAgentsList)
		err = nsxclient.Do(updateAPI)

		if err != nil {
			return fmt.Errorf("Error: %v", err)
		} else if updateAPI.StatusCode() != 204 {
			return fmt.Errorf(updateAPI.GetResponse())
		} else {
			log.Printf("Updated DHCP Relay - %s", updateAPI.GetResponse())
		}
	}

	return nil
}
