package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/dhcprelay"
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
		Update: resourceDHCPRelayUpdate,

		Schema: map[string]*schema.Schema{

			"ipsets": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IPsets (at least one of ipsets, fqdn, dhcpserverip must be provided)",
			},
			"fqdn": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Comma separated lists of Domain names, max: 2 (at least one of ipsets, fqdn, dhcpserverip must be provided)",
			},

			"edgeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dhcpserverip": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Comma separated lists of IP Addresses (at least one of ipsets, fqdn, dhcpserverip must be provided)",
			},
			"agent": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vnicindex": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},

						"giaddress": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    false,
							Description: "Gateway address of network attached to the vNic, defaults to the vNic primary address",
						},
					},
				},
			},
		},
	}
}

func resourceDHCPRelayCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string
	var agentList []dhcprelay.RelayAgent
	var dhcpRelay dhcprelay.DhcpRelay

	if v, ok := d.GetOk("ipsets"); ok {
		for _, ipset := range v.([]interface{}) {
			dhcpRelay.RelayServer.IPSets = append(dhcpRelay.RelayServer.IPSets, ipset.(string))
		}
	}

	if v, ok := d.GetOk("fqdn"); ok {
		for _, name := range v.([]interface{}) {
			dhcpRelay.RelayServer.DomainName = append(dhcpRelay.RelayServer.DomainName, name.(string))
		}
	}

	if len(dhcpRelay.RelayServer.DomainName) > 2 {
		return fmt.Errorf("Error: Field fqdn only Supports 2 domains")
	}

	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}

	if v, ok := d.GetOk("agent"); ok {
		if agents, ok := v.(*schema.Set); ok {
			for _, value := range agents.List() {
				agentObject := value.(map[string]interface{})
				newAgent := dhcprelay.RelayAgent{}
				if vnicIndexValue, ok := agentObject["vnicindex"]; ok {
					newAgent.VnicIndex = vnicIndexValue.(string)
				}

				if giAddressValue, ok := agentObject["giaddress"]; ok {
					newAgent.GiAddress = giAddressValue.(string)
				}

				if dhcpServerIPValue, ok := agentObject["dhcpserverip"]; ok {
					newAgent.GiAddress = dhcpServerIPValue.(string)
				}

				//finally add to the list
				agentList = append(agentList, newAgent)
			}

		}

		dhcpRelay.RelayAgents = agentList
	} else {
		return fmt.Errorf("agent parameter is required (DHCP Relay Agent) ")
	}

	if v, ok := d.GetOk("dhcpserverip"); ok {
		for _, ipaddr := range v.([]interface{}) {
			dhcpRelay.RelayServer.IPAddress = append(dhcpRelay.RelayServer.IPAddress, ipaddr.(string))
		}
	} else {
		return fmt.Errorf("dhcpserverip argument is required")
	}
	if len(dhcpRelay.RelayServer.IPAddress) > 16 {
		return fmt.Errorf("Error: Field IP addresses  only Supports 16 records")
	}
	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)
	updateAPI := dhcprelay.NewCreate(edgeid, dhcpRelay)
	err := nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update the DHCP relay %s", updateAPI.GetResponse())
	}

	// If we get here, everything is OK.  Set the ID for the Terraform state
	// and return the response from the READ method.
	d.SetId(edgeid)
	return resourceDHCPRelayRead(d, m)
	//return nil
}

func resourceDHCPRelayRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string
	var agentList []dhcprelay.RelayAgent
	// Gather the attributes for the resource.
	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}

	DHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if len(DHCPRelay.RelayServer.IPSets) > 0 {
		d.Set("ipsets", DHCPRelay.RelayServer.IPSets)
	}

	if len(DHCPRelay.RelayServer.DomainName) > 0 {
		d.Set("fqdn", DHCPRelay.RelayServer.DomainName)
	}

	if len(DHCPRelay.RelayServer.IPAddress) > 0 {
		d.Set("dhcpserverip", DHCPRelay.RelayServer.IPAddress)
	}

	if len(DHCPRelay.RelayAgents) > 0 {
		if v, ok := d.GetOk("agent"); ok {
			if agents, ok := v.(*schema.Set); ok {
				for _, value := range agents.List() {
					agentObject := value.(map[string]interface{})
					newAgent := dhcprelay.RelayAgent{}
					if vnicIndexValue, ok := agentObject["vnicindex"]; ok {
						newAgent.VnicIndex = vnicIndexValue.(string)
					}

					if giAddressValue, ok := agentObject["giaddress"]; ok {
						newAgent.GiAddress = giAddressValue.(string)
					}

					if dhcpServerIPValue, ok := agentObject["dhcpserverip"]; ok {
						newAgent.GiAddress = dhcpServerIPValue.(string)
					}

					//finally add to the list
					agentList = append(agentList, newAgent)
				}

			}

			DHCPRelay.RelayAgents = agentList
		}
	} else {
		return fmt.Errorf("RelayAgents should not be empty since this is not allowed by api")
	}

	return nil
}

func resourceDHCPRelayUpdate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var agentList []dhcprelay.RelayAgent
	var currentRelay *dhcprelay.DhcpRelay
	var hasChanges bool
	var updateObject dhcprelay.DhcpRelay
	currentRelay, getAllErr := getAllDhcpRelays(d.Id(), nsxclient)
	if getAllErr != nil {
		return fmt.Errorf("Error: %v", getAllErr)
	}

	if d.HasChange("ipsets") {
		hasChanges = true
		if v, ok := d.GetOk("ipsets"); ok {
			for _, ipset := range v.([]interface{}) {
				updateObject.RelayServer.IPSets = append(updateObject.RelayServer.IPSets, ipset.(string))
			}
		}

	} else {
		updateObject.RelayServer.IPSets = currentRelay.RelayServer.IPSets
	}

	if d.HasChange("fqdn") {
		hasChanges = true
		if v, ok := d.GetOk("fqdn"); ok {
			for _, name := range v.([]interface{}) {
				updateObject.RelayServer.DomainName = append(updateObject.RelayServer.DomainName, name.(string))
			}
		}
	} else {
		updateObject.RelayServer.DomainName = currentRelay.RelayServer.DomainName
	}

	if len(updateObject.RelayServer.DomainName) > 2 {
		return fmt.Errorf("Error: Field fqdn only Supports 2 domains")
	}

	if d.HasChange("dhcpserverip") {
		hasChanges = true
		if v, ok := d.GetOk("dhcpserverip"); ok {
			for _, ipaddr := range v.([]interface{}) {
				updateObject.RelayServer.IPAddress = append(updateObject.RelayServer.IPAddress, ipaddr.(string))
			}
		}
	} else {
		updateObject.RelayServer.IPAddress = currentRelay.RelayServer.IPAddress
	}

	if len(updateObject.RelayServer.IPAddress) > 16 {
		return fmt.Errorf("Error: Field IP addresses  only Supports 16 records")
	}
	if d.HasChange("agent") {
		hasChanges = true
		if v, ok := d.GetOk("agent"); ok {
			if agents, ok := v.(*schema.Set); ok {
				for _, value := range agents.List() {
					agentObject := value.(map[string]interface{})
					newAgent := dhcprelay.RelayAgent{}
					if vnicIndexValue, ok := agentObject["vnicindex"]; ok {
						newAgent.VnicIndex = vnicIndexValue.(string)
					}

					if giAddressValue, ok := agentObject["giaddress"]; ok {
						newAgent.GiAddress = giAddressValue.(string)
					}

					if dhcpServerIPValue, ok := agentObject["dhcpserverip"]; ok {
						newAgent.GiAddress = dhcpServerIPValue.(string)
					}

					//finally add to the list
					agentList = append(agentList, newAgent)
				}

			}

			updateObject.RelayAgents = agentList
		}
	} else {
		updateObject.RelayAgents = currentRelay.RelayAgents
	}

	if hasChanges {
		nsxMutexKV.Lock(d.Id())
		defer nsxMutexKV.Unlock(d.Id())
		updateAPI := dhcprelay.NewUpdate(d.Id(), updateObject)
		err := nsxclient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Could not update the resource : %s", err)
		}
		if updateAPI.StatusCode() != 204 {
			d.SetId("")
			return fmt.Errorf("Error updating record : %s", updateAPI.GetResponse())
		}
		return resourceDHCPRelayRead(d, m)
	}

	return nil

}

func resourceDHCPRelayDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var edgeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("edgeid"); ok {
		edgeid = v.(string)
	} else {
		return fmt.Errorf("edgeid argument is required")
	}
	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)
	deleteAPI := dhcprelay.NewDelete(edgeid)
	err := nsxclient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	log.Println("DHCP Relay agent deleted.")
	d.SetId("")
	return nil
}
