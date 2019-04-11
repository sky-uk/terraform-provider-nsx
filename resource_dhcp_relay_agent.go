package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/dhcprelay"
)

func getAllDhcpRelayAgents(edgeID string, nsxclient *gonsx.NSXClient) (*dhcprelay.DhcpRelay, error) {
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

	log.Println("Get Agents Response: ", api.GetResponse())
	return api.GetResponse(), nil
}

func printDHCPRelay(relay *dhcprelay.DhcpRelay) {
	log.Printf("------------ Relay ------------")
	dhcprelay_xml, err := xml.MarshalIndent(relay, "", "  ")
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf(string(dhcprelay_xml))
}

func resourceDHCPRelayAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceDHCPRelayAgentCreate,
		Read:   resourceDHCPRelayAgentRead,
		Delete: resourceDHCPRelayAgentDelete,

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
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Gateway address of network attached to the vNic, defaults to the vNic primary address",
			},
		},
	}
}

func resourceDHCPRelayAgentCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string
	var vnicindex string
	var giaddress string

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

	var id = edgeid + ":" + vnicindex
	log.Printf("ID: %s", id)

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	DHCPRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	for _, element := range DHCPRelay.RelayAgents {
		log.Printf("VnicIndex %v", element.VnicIndex)
		log.Printf("GiAddress %v", element.GiAddress)
		if element.VnicIndex == vnicindex {
			return fmt.Errorf("Relay Agent does already exist: vnicindex: %s", vnicindex)
		}
	}

	newAgent := dhcprelay.RelayAgent{}
	newAgent.VnicIndex = vnicindex
	newAgent.GiAddress = giaddress

	DHCPRelay.RelayAgents = append(DHCPRelay.RelayAgents, newAgent)

	updateAPI := dhcprelay.NewUpdate(edgeid, *DHCPRelay)
	err := nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update the DHCP relay %s", updateAPI.GetResponse())
	}

	d.SetId(id)
	return resourceDHCPRelayRead(d, m)
}

func resourceDHCPRelayAgentRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string
	var vnicindex string
	var giaddress string

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

	var id = edgeid + ":" + vnicindex
	log.Printf("ID: %s", id)

	DHCPRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	for _, element := range DHCPRelay.RelayAgents {
		if element.VnicIndex == vnicindex || element.GiAddress == giaddress {
			d.SetId(id)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceDHCPRelayAgentDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string
	var vnicindex string
	var giaddress string

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

	var id = edgeid + ":" + vnicindex
	log.Printf("ID: %s", id)

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	DHCPRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	for i, element := range DHCPRelay.RelayAgents {
		log.Printf("VnicIndex %v", element.VnicIndex)
		log.Printf("GiAddress %v", element.GiAddress)
		if element.VnicIndex == vnicindex || element.GiAddress == giaddress {
			DHCPRelay.RelayAgents = append(DHCPRelay.RelayAgents[:i], DHCPRelay.RelayAgents[i+1:]...)
			break
		}
	}

	updateAPI := dhcprelay.NewUpdate(edgeid, *DHCPRelay)
	err := nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update the DHCP relay %s", updateAPI.GetResponse())
	}

	d.SetId("")
	return nil
}
