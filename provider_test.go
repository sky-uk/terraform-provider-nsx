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
