package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/skyinfoblox"
	"github.com/sky-uk/skyinfoblox/api/records"
	"testing"
)

func TestAccResourceDHCPRelay(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: TestAccResourceDHCPRelayDestroy,
		Steps: []resource.TestStep{
			{},
		},
	})

}

func TestAccResourceDHCPRelayDestroy(state *terraform.State) error {
	return nil

}

func TestAccResourceDHCPRelayExists(state *terraform.State) resource.TestCheckFunc {

}

func testAccResourceDHCPRelayCreateTemplate(edgeid string) string {
	return fmt.Sprintf(`
	resource "nsx_dhcp_relay" "testrelay" {
		ipsets = ["ipset-3"]
  		fqdn = ["testdomain.paas.bskyb.com","testdomain2.paas.bskyb.com"]
  		edgeid       = "%s"
  		dhcpserverip = ["10.152.160.10"]
  		agent {
	       		vnicindex="9"
	       		giaddress="10.72.232.200"
	  	}
	}`, edgeid)
}

func testAccResourceDHCPRelayUpdateTemplate(edgeid string) string {
	return fmt.Sprintf(`
	resource "nsx_dhcp_relay" "testrelay" {
		ipsets = ["ipset-3"]
  		fqdn = ["testdomain.paas.bskyb.com"]
  		edgeid       = "%s"
  		dhcpserverip = ["10.152.160.10"]
  		agent {
	       		vnicindex="9"
	       		giaddress="10.72.232.200"
	  	}
	}`, edgeid)
}
