package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"nsx": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NSXUSERNAME"); v == "" {
		t.Fatal("NSXUSERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("NSXPASSWORD"); v == "" {
		t.Fatal("NSXPASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("NSXSERVER"); v == "" {
		t.Fatal("NSXSERVER must be set for acceptance tests")
	}
}

func loadDLRId(t *testing.T) string {
	edgeid := os.Getenv("NSX_TESTING_DLR_ID")
	if edgeid == "" {
		t.Skip("skipping test; NSX_TESTING_DLR_ID not set")
	}
	return edgeid
}

func loadESGId(t *testing.T) string {
	edgeid := os.Getenv("NSX_TESTING_ESG_ID")
	if edgeid == "" {
		t.Skip("skipping test; NSX_TESTING_ESG_ID not set")
	}
	return edgeid
}

func loadScopeId(t *testing.T) string {
	scopeid := os.Getenv("NSX_TESTING_SCOPE_ID")
	if scopeid == "" {
		t.Skip("skipping test; NSX_TESTING_SCOPE_ID not set")
	}
	return scopeid
}

func loadVirtualwireId(t *testing.T) string {
	virtualwire := os.Getenv("NSX_TESTING_VIRTUALWIRE_ID")
	if virtualwire == "" {
		t.Skip("skipping test; NSX_TESTING_VIRTUALWIRE_ID not set")
	}
	return virtualwire
}
