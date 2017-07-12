package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/edgeinterface"
	"net/http"
	"strconv"
	"testing"
)

func TestAccResourceEdgeInterface(t *testing.T) {
	edgeid := "edge-7"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceEdgeInterfaceCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEdgeInterfaceCreateTemplate(edgeid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsx_edge_interface.testAccInterface", "edgeid"),
					resource.TestCheckResourceAttr("nsx_edge_interface.testAccInterface", "edgeid", edgeid),
				),
			},
		},
	})

}

func testAccResourceEdgeInterfaceCheckDestroy(state *terraform.State) error {
	nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
	// TODO  this seems not having effect on client behaviour...
	nsxClient.IgnoreSSL = true

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsx_edge_interface" {
			continue
		}

		var edgeid, indexStr string
		var index int
		var ok bool
		if edgeid, ok := rs.Primary.Attributes["edgeid"]; ok && edgeid != "" {
			return nil
		}
		if indexStr, ok = rs.Primary.Attributes["index"]; ok {
			return nil
		}

		index, _ = strconv.Atoi(indexStr)

		api := edgeinterface.NewGet(edgeid, index)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}
		if api.StatusCode() == http.StatusOK {
			return fmt.Errorf("Resource still exists")
		}
	}
	return nil

}

func testAccResourceEdgeInterfaceExists(edgeid string, index int, name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("nsx_edge_interface resource does not exist")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("nsx_edge_interface resource ID not set")
		}

		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)
		nsxClient.IgnoreSSL = true

		api := edgeinterface.NewGet(edgeid, index)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}
		if api.StatusCode() == http.StatusOK {
			return nil
		}
		return fmt.Errorf("Resource does not exists")
	}
}

func testAccResourceEdgeInterfaceCreateTemplate(edgeid string) string {

	return fmt.Sprintf(`resource "nsx_edge_interface" "testAccInterface" {
    edgeid = "edge-7"
    name = "edge_interface_testacc"
    isconnected = true
    connectedtoid = "virtualwire-182"
    interfacetype = "internal"
    mtu = 1500
    addressgroups = [
    /*
        {
            "primaryaddress" = "10.152.172.1"
            "subnetmask"     = "255.255.255.0"
        },
        */
    ]
}`)
}
