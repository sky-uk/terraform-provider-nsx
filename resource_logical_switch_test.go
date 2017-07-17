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
			return testAccNSXLogicalSwitchCheckDestroy(state, switchName)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccNSXLogicalSwitchNoNameTemplate(scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccNSXLogicalSwitchNoDescriptionTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccNSXLogicalSwitchNoTenantIDTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccNSXLogicalSwitchNoScopeIDTemplate(switchName),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccNSXLogicalSwitchNoControlPlaneModeTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`required field is not set`),
			},
			{
				Config:      testAccNSXLogicalSwitchInvalidControlPlaneModeTemplate(switchName, scopeID),
				ExpectError: regexp.MustCompile(`must be one of UNICAST_MODE, HYBRID_MODE or MULTICAST_MODE`),
			},
			{
				Config: testAccNSXLogicalSwitchCreateTemplate(switchName, scopeID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(switchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "name", switchName),
					resource.TestCheckResourceAttr(testResourceName, "desc", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "tenantid", "tf_testid"),
					resource.TestCheckResourceAttr(testResourceName, "scopeid", scopeID),
					resource.TestCheckResourceAttr(testResourceName, "controlplanemode", "UNICAST_MODE"),
					resource.TestCheckResourceAttrSet(testResourceName, "labels.0"),
				),
			},
			{
				Config: testAccNSXLogicalSwitchUpdateTemplate(updateSwitchName, scopeID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLogicalSwitchExists(updateSwitchName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "name", updateSwitchName),
					resource.TestCheckResourceAttr(testResourceName, "desc", "Acceptance Test Update"),
					resource.TestCheckResourceAttr(testResourceName, "tenantid", "tf_testid"),
					resource.TestCheckResourceAttr(testResourceName, "scopeid", scopeID),
					resource.TestCheckResourceAttr(testResourceName, "controlplanemode", "UNICAST_MODE"),
					resource.TestCheckResourceAttrSet(testResourceName, "labels.0"),
				),
			},
		},
	})

}

func testAccNSXLogicalSwitchExists(name, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX logical switch resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX logical switch resource ID not set in resources ")
		}

		getLogicalSwitchAPI := virtualwire.NewGet(resourceID)
		err := nsxClient.Do(getLogicalSwitchAPI)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}
		responseCode := getLogicalSwitchAPI.StatusCode()
		logicalSwitch := getLogicalSwitchAPI.GetResponse()
		if responseCode != http.StatusOK {
			return fmt.Errorf("Error while checking if logical switch %s exists. HTTP return code was %d", resourceID, responseCode)
		}

		if name == logicalSwitch.Name {
			return nil
		}
		return fmt.Errorf("NSX logical switch %s wasn't found", name)
	}
}

func testAccNSXLogicalSwitchCheckDestroy(state *terraform.State, name string) error {

	nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsx_logical_switch" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		getLogicalSwitchAPI := virtualwire.NewGet(resourceID)
		err := nsxClient.Do(getLogicalSwitchAPI)
		if err != nil {
			return fmt.Errorf("Error while retrieving logical switch ID %s. Error: %v", resourceID, err)
		}

		if getLogicalSwitchAPI.StatusCode() == http.StatusNotFound {
			return nil
		}

		logicalSwitch := getLogicalSwitchAPI.GetResponse()
		if name == logicalSwitch.Name {
			return fmt.Errorf("NSX logical switch %s still exists", name)
		}
	}
	return nil
}

func testAccNSXLogicalSwitchNoNameTemplate(scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
desc = "Acceptance Test No Name Template"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, scopeID)
}

func testAccNSXLogicalSwitchNoDescriptionTemplate(name, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, name, scopeID)
}

func testAccNSXLogicalSwitchNoTenantIDTemplate(name, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, name, scopeID)
}

func testAccNSXLogicalSwitchNoScopeIDTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
tenantid = "tf_testid"
controlplanemode = "UNICAST_MODE"
}`, name)
}

func testAccNSXLogicalSwitchNoControlPlaneModeTemplate(name, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
tenantid = "tf_testid"
scopeid = "%s"
}`, name, scopeID)
}

func testAccNSXLogicalSwitchInvalidControlPlaneModeTemplate(switchName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "INVALID_CONTROL_PLANE_MODE"
}`, switchName, scopeID)
}

func testAccNSXLogicalSwitchCreateTemplate(switchName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, switchName, scopeID)
}

func testAccNSXLogicalSwitchUpdateTemplate(switchUpdateName, scopeID string) string {
	return fmt.Sprintf(`
resource "nsx_logical_switch" "acctest" {
name = "%s"
desc = "Acceptance Test Update"
tenantid = "tf_testid"
scopeid = "%s"
controlplanemode = "UNICAST_MODE"
}`, switchUpdateName, scopeID)
}
