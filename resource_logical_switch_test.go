package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/virtualwire"
	"net/http"
	"regexp"
	"testing"
)

func TestAccNSXLogicalSwitchBasic(t *testing.T) {

	randomInt := acctest.RandInt()
	switchName := fmt.Sprintf("acctest-nsx-logical-switch-%d", randomInt)
	updateSwitchName := fmt.Sprintf("acctest-nsx-logical-switch-%d-update", randomInt)
	scopeID := "vdnscope-1"
	testResourceName := "nsx_logical_switch.acctest"

	fmt.Printf("\n\nlogical switch name is %s\n\n", switchName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName, scopeID)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccInfobloxLogicalSwitchNoNameTemplate(scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxLogicalSwitchNoDescriptionTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxLogicalSwitchNoTenantIDTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxLogicalSwitchNoScopeIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccInfobloxLogicalSwitchInvalidControlPlaneModeTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`must be one of UNICAST_MODE, HYBRID_MODE or MULTICAST_MODE`),
			},
			{
				Config: testAccInfobloxLogicalSwitchCreateTemplate(switchName, scopeID),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxLogicalSwitchExists(switchName, scopeID, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "desc", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tenantid", "tf_testid"),
					resource.TestCheckResourceAttr(testResourceName, "scopeid", scopeID),
					resource.TestCheckResourceAttr(testResourceName, "controlplanemode", "UNICAST_MODE"),
				),
			},
			{
				Config: testAccInfobloxLogicalSwitchUpdateTemplate(updateSwitchName, scopeID),
				Check: resource.ComposeTestCheckFunc(
					testAccInfobloxLogicalSwitchExists(updateSwitchName, scopeID, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "desc", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tenantid", "tf_testid"),
					resource.TestCheckResourceAttr(testResourceName, "scopeid", scopeID),
					resource.TestCheckResourceAttr(testResourceName, "controlplanemode", "UNICAST_MODE"),
				),
			},
		},
	})

}

func testAccInfobloxLogicalSwitchExists(name, scopeID, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical switch resource %s not found in resources", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("NSX logical switch resource ID not set in resources ")
		}

		getAllAPI := virtualwire.NewGetAll(scopeID)
		err := nsxClient.Do(getAllAPI)
		if err != nil {
			return fmt.Errorf("Error while checking if logical switch exists %v", err)
		}
		if getAllAPI.StatusCode() != http.StatusOK {
			return fmt.Errorf("Error while checking if logical switch exists. HTTP return code was %d", getAllAPI.StatusCode())
		}
		virtualWire := getAllAPI.GetResponse().FilterByName(name)

		if name == virtualWire.Name {
			return nil
		}
		return fmt.Errorf("NSX logical switch %s wasn't found", name)
	}
}

func testAccNSXLogicalSwitchCheckDestroy(state *terraform.State, name, scopeID string) error {

	nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsx_logical_switch" {
			continue
		}
		if id, ok := rs.Primary.Attributes["id"]; ok && id != "" {
			return nil
		}

		getAllAPI := virtualwire.NewGetAll(scopeID)
		err := nsxClient.Do(getAllAPI)
		if err != nil {
			return nil
		}
		virtualWire := getAllAPI.GetResponse().FilterByName(name)

		if name == virtualWire.Name {
			return fmt.Errorf("NSX logical switch %s still exists", name)
		}
	}
	return nil
}

func testAccInfobloxLogicalSwitchNoNameTemplate(scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
desc = "Acceptance Test Update"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, scopeID)
}

func testAccInfobloxLogicalSwitchNoDescriptionTemplate(name, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, name, scopeID)
}

func testAccInfobloxLogicalSwitchNoTenantIDTemplate(name, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, name, scopeID)
}

func testAccInfobloxLogicalSwitchNoScopeIDTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
tenantid = "tf_testid"
controlplanemode = "UNICAST_MODE"
}`, name)
}

func testAccInfobloxLogicalSwitchInvalidControlPlaneModeTemplate(switchName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "INVALID_CONTROL_PLANE_MODE"
}`, switchName, scopeID)
}

func testAccInfobloxLogicalSwitchCreateTemplate(switchName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, switchName, scopeID)
}

func testAccInfobloxLogicalSwitchUpdateTemplate(switchUpdateName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, switchUpdateName, scopeID)
}
