package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/dhcprelay"
	"testing"
)

func TestAccResourceDHCPRelayAgent(t *testing.T) {
	edgeid := loadDLRId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceDHCPRelayCheckDestroy,
		Steps: []resource.TestStep{
			// Create DHCP Relay and Agents
			// Using the edgeid attribute from the testrelay is important to get dependencies right
			{
				Config: fmt.Sprintf(`
                                     resource "nsx_dhcp_relay_agent" "agent1" {
                                         edgeid       = "${nsx_dhcp_relay.testrelay.edgeid}"
                                         vnicindex="10"
                                         giaddress="192.168.1.1"
                                     }
                                     resource "nsx_dhcp_relay" "testrelay" {
                                         edgeid       = "%s"
                                         dhcpserverip = ["10.152.160.10","10.152.160.11"]
                                     }
                                     resource "nsx_dhcp_relay_agent" "agent2" {
                                         edgeid       = "${nsx_dhcp_relay.testrelay.edgeid}"
                                         vnicindex="11"
                                         giaddress="192.168.2.1"
                                     }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceDHCPRelayExists(edgeid, "nsx_dhcp_relay.testrelay"),
					resource.TestCheckResourceAttrSet("nsx_dhcp_relay.testrelay", "edgeid"),
					resource.TestCheckResourceAttr("nsx_dhcp_relay.testrelay", "edgeid", edgeid),
					testAccResourceDHCPRelayAgentExists(edgeid, "10", "192.168.1.1"),
					testAccResourceDHCPRelayAgentExists(edgeid, "11", "192.168.2.1"),
				),
			},
			// Update DHCP Relay and Remove agent
			{
				Config: fmt.Sprintf(`
                                     resource "nsx_dhcp_relay" "testrelay" {
                                         edgeid       = "%s"
                                         dhcpserverip = ["10.152.160.10"]
                                     }
                                     resource "nsx_dhcp_relay_agent" "agent1" {
                                         edgeid       = "${nsx_dhcp_relay.testrelay.edgeid}"
                                         vnicindex="11"
                                         giaddress="192.168.2.1"
                                     }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceDHCPRelayExists(edgeid, "nsx_dhcp_relay.testrelay"),
					resource.TestCheckResourceAttrSet("nsx_dhcp_relay.testrelay", "edgeid"),
					resource.TestCheckResourceAttr("nsx_dhcp_relay.testrelay", "edgeid", edgeid),
					testAccResourceDHCPRelayAgentExists(edgeid, "11", "192.168.2.1"),
					testAccResourceDHCPRelayAgentDoesNotExists(edgeid, "10"),
				),
			},
			// Update DHCP Agent
			{
				Config: fmt.Sprintf(`
                                     resource "nsx_dhcp_relay" "testrelay" {
                                         edgeid       = "%s"
                                         dhcpserverip = ["10.152.160.10"]
                                     }
                                     resource "nsx_dhcp_relay_agent" "agent1" {
                                         edgeid       = "${nsx_dhcp_relay.testrelay.edgeid}"
                                         vnicindex="10"
                                         giaddress="192.168.1.1"
                                     }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceDHCPRelayExists(edgeid, "nsx_dhcp_relay.testrelay"),
					resource.TestCheckResourceAttrSet("nsx_dhcp_relay.testrelay", "edgeid"),
					resource.TestCheckResourceAttr("nsx_dhcp_relay.testrelay", "edgeid", edgeid),
					testAccResourceDHCPRelayAgentExists(edgeid, "10", "192.168.1.1"),
					testAccResourceDHCPRelayAgentDoesNotExists(edgeid, "11"),
				),
			},
		},
	})
}

func TestAccResourceDHCPRelayAgentGIAddressDefault(t *testing.T) {
	edgeid := loadDLRId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceDHCPRelayCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
                                     resource "nsx_dhcp_relay_agent" "agent1" {
                                         edgeid = "%s"
                                         vnicindex="10"
                                     }
                                     `, edgeid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nsx_dhcp_relay_agent.agent1", "edgeid", edgeid),
					resource.TestCheckResourceAttr("nsx_dhcp_relay_agent.agent1", "vnicindex", "10"),
					testAccResourceDHCPRelayAgentExists(edgeid, "10", "192.168.1.1"),
				),
			},
		},
	})
}

func testAccResourceDHCPRelayAgentExists(edgeid, vnicindex string, giaddress string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
		api := dhcprelay.NewGetAll(edgeid)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		for _, element := range api.GetResponse().RelayAgents {
			if element.VnicIndex == vnicindex && element.GiAddress == giaddress {
				return nil
			}
		}

		return fmt.Errorf("Agent %s:%s does not exist", vnicindex, giaddress)
	}
}

func testAccResourceDHCPRelayAgentDoesNotExists(edgeid, vnicindex string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
		api := dhcprelay.NewGetAll(edgeid)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		for _, element := range api.GetResponse().RelayAgents {
			if element.VnicIndex == vnicindex {
				return fmt.Errorf("Agent %s does still exist.", vnicindex)
			}
		}

		return nil
	}
}
