package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

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
				Computed:    true,
				Description: "Gateway address of network attached to the vNic, defaults to the vNic primary address",
			},
		},
	}
}

func resourceDHCPRelayAgentCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	edgeid := d.Get("edgeid").(string)
	vnicindex := d.Get("vnicindex").(string)

	var id = edgeid + ":" + vnicindex

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Get existing Configuration
	dhcpRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	// Check if already in list
	for _, element := range dhcpRelay.RelayAgents {
		if element.VnicIndex == vnicindex {
			return fmt.Errorf("Relay Agent does already exist: vnicindex: %s", vnicindex)
		}
	}

	// Add to list
	newAgent := dhcprelay.RelayAgent{}
	newAgent.VnicIndex = vnicindex
	if giaddress, ok := d.GetOk("giaddress"); ok {
		newAgent.GiAddress = giaddress.(string)
	}
	dhcpRelay.RelayAgents = append(dhcpRelay.RelayAgents, newAgent)

	// API Update
	updateAPI := dhcprelay.NewUpdate(edgeid, *dhcpRelay)
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

	s := strings.Split(d.Id(), ":")
	edgeid, vnicindex := s[0], s[1]

	dhcpRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	for _, element := range dhcpRelay.RelayAgents {
		if element.VnicIndex == vnicindex {
			// Found
			d.Set("vnicindex", element.VnicIndex)
			d.Set("giaddress", element.GiAddress)
			return nil
		}
	}

	// Not found
	d.SetId("")
	return nil
}

func resourceDHCPRelayAgentDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)

	s := strings.Split(d.Id(), ":")
	edgeid, vnicindex := s[0], s[1]

	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	// Get existing Configuration
	dhcpRelay, _ := getAllDhcpRelayAgents(edgeid, nsxclient)

	// Remove (if exists)
	for i, element := range dhcpRelay.RelayAgents {
		if element.VnicIndex == vnicindex {
			dhcpRelay.RelayAgents = append(dhcpRelay.RelayAgents[:i], dhcpRelay.RelayAgents[i+1:]...)
			break
		}
	}

	// Update
	updateAPI := dhcprelay.NewUpdate(edgeid, *dhcpRelay)
	err := nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update the DHCP relay %s", updateAPI.GetResponse())
	}

	d.SetId("")
	return nil
}
