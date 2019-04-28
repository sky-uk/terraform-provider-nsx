package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/dhcprelay"
	"testing"
)

func TestAccResourceDHCPRelay(t *testing.T) {
	edgeid := loadDLRId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceDHCPRelayCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDHCPRelayCreateTemplate(edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceDHCPRelayExists(edgeid, "nsx_dhcp_relay.testrelay"),
					resource.TestCheckResourceAttrSet("nsx_dhcp_relay.testrelay", "edgeid"),
					resource.TestCheckResourceAttr("nsx_dhcp_relay.testrelay", "edgeid", edgeid),
				),
			}, {
				Config: testAccResourceDHCPRelayUpdateTemplate(edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceDHCPRelayExists(edgeid, "nsx_dhcp_relay.testrelay"),
					resource.TestCheckResourceAttrSet("nsx_dhcp_relay.testrelay", "edgeid"),
					resource.TestCheckResourceAttr("nsx_dhcp_relay.testrelay", "edgeid", edgeid),
				),
			},
		},
	})
}

func testAccResourceDHCPRelayCheckDestroy(state *terraform.State) error {
	nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsx_dhcp_relay" {
			continue
		}

		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}
		api := dhcprelay.NewGetAll(rs.Primary.Attributes["id"])
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		// Check if RelayServer is not set
		if len(api.GetResponse().RelayServer.IPAddress) > 0 {
			return fmt.Errorf("DHCP Relay still exists")

		}
	}
	return nil

}

func testAccResourceDHCPRelayExists(edgeid, resourcename string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourcename]

		if !ok {
			return fmt.Errorf("DHCPRelay resource does not exist")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DHCPRelay resource ID not set")
		}
		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
		api := dhcprelay.NewGetAll(edgeid)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		// Check if RelayServer is set
		if len(api.GetResponse().RelayServer.IPAddress) > 0 {
			return nil
		}

		return fmt.Errorf("Relay does not exist")
	}

}

func testAccResourceDHCPRelayCreateTemplate(edgeid string) string {
	return fmt.Sprintf(`
	resource "nsx_dhcp_relay" "testrelay" {
  		edgeid       = "%s"
  		dhcpserverip = ["10.152.160.10"]
  		agent {
	       		vnicindex="10"
	       		giaddress="192.168.1.1"
	  	}
	}`, edgeid)
}

func testAccResourceDHCPRelayUpdateTemplate(edgeid string) string {
	return fmt.Sprintf(`
	resource "nsx_dhcp_relay" "testrelay" {
  		edgeid       = "%s"
  		dhcpserverip = ["10.152.160.10","10.152.160.11"]
  		agent {
	       		vnicindex="10"
	       		giaddress="192.168.1.1"
	  	}
	}`, edgeid)
}
