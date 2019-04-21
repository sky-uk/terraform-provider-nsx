package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/nat"
	"regexp"
	"testing"
)

func TestAccResourceNatRuleFormat(t *testing.T) {
	edgeid := loadESGId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(s *terraform.State) error { return testAccResourceNatRulesAreEmpty(edgeid) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0"
                                        action              = "dnat"
                                        original_address    = "1.1.1.3222"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				ExpectError: regexp.MustCompile(`expected original_address to contain a valid IP`),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0-address"
                                        action              = "dnat"
                                        original_address    = "1.1.1.22"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccNatRuleWithDescriptionExists(edgeid, "rule0-address"),
				),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0"
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "any2"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				ExpectError: regexp.MustCompile(`expected original_address to contain a valid IP`),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0-any"
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "any"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccNatRuleWithDescriptionExists(edgeid, "rule0-any"),
				),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0" 
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "192.168.1.0/244"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				ExpectError: regexp.MustCompile(`expected original_address to contain a valid IP`),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0-cidr" 
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "192.168.1.0/24"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = 80
                                        translated_port     = 80
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccNatRuleWithDescriptionExists(edgeid, "rule0-cidr"),
				),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0" 
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "192.168.1.1"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = "any2"
                                        translated_port     = 80
                                    }`, edgeid),
				ExpectError: regexp.MustCompile(`Specify any, a port(.*) or port range(.*)"`),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0-port" 
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "192.168.1.1"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = "any"
                                        translated_port     = 80
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccNatRuleWithDescriptionExists(edgeid, "rule0-port"),
				),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule0" {
                                        edgeid              = "%s"
                                        description         = "rule0-portrange" 
                                        action              = "dnat"
                                        vnic                = 0
                                        original_address    = "192.168.1.1"
                                        translated_address  = "192.168.1.1"
                                        protocol            = "tcp"
                                        original_port       = "1234-1239"
                                        translated_port     = 80
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					testAccNatRuleWithDescriptionExists(edgeid, "rule0-portrange"),
				),
			},
		},
	})
}

func TestAccResourceDNatRule(t *testing.T) {
	edgeid := loadESGId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(s *terraform.State) error { return testAccResourceNatRulesAreEmpty(edgeid) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule1" {
	                                     edgeid              = "%s"
	                                     action              = "dnat"
	                                     description         = "rule1"
	                                     original_address    = "1.1.1.3"
	                                     translated_address  = "192.168.1.1"
	                                     protocol            = "tcp"
	                                     translated_port     = 80
                                         dnat_match_source_address  = "192.168.10.1/24"
                                         dnat_match_source_port     = "1234-1239"
	                                 }
	                                 resource "nsx_nat_rule" "rule2" {
	                                     edgeid              = "%s"
	                                     action              = "dnat"
	                                     description         = "rule2"
	                                     original_address    = "1.1.1.3"
	                                     translated_address  = "192.168.1.1"
	                                     protocol            = "tcp"
                                         original_port       = 88
	                                     translated_port     = 80
	                                 }`, edgeid, edgeid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nsx_nat_rule.rule1", "edgeid", edgeid),
					resource.TestCheckResourceAttrSet("nsx_nat_rule.rule1", "id"),
					resource.TestCheckResourceAttr("nsx_nat_rule.rule1", "original_port", "any"),
					testAccNatRuleWithDescriptionExists(edgeid, "rule1"),
					resource.TestCheckResourceAttr("nsx_nat_rule.rule2", "edgeid", edgeid),
					resource.TestCheckResourceAttrSet("nsx_nat_rule.rule2", "id"),
					resource.TestCheckResourceAttr("nsx_nat_rule.rule2", "original_port", "88"),
					testAccNatRuleWithDescriptionExists(edgeid, "rule2"),
				),
			},
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule1" {
			                                     edgeid              = "%s"
			                                     action              = "dnat"
			                                     description         = "rule1"
			                                     original_address    = "1.1.1.3"
			                                     translated_address  = "192.168.1.1"
			                                     protocol            = "tcp"
			                                     original_port       = 81
			                                     translated_port     = 80
			                                 }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nsx_nat_rule.rule1", "edgeid", edgeid),
					resource.TestCheckResourceAttrSet("nsx_nat_rule.rule1", "id"),
					resource.TestCheckResourceAttr("nsx_nat_rule.rule1", "original_port", "81"),
					testAccNatRuleWithDescriptionExists(edgeid, "rule1"),
					testAccNatRuleWithDescriptionNotExists(edgeid, "rule2"),
				),
			},
		},
	})
}

func TestAccResourceSNatRule(t *testing.T) {
	edgeid := loadESGId(t)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(s *terraform.State) error { return testAccResourceNatRulesAreEmpty(edgeid) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "nsx_nat_rule" "rule3" {
                                        edgeid              = "%s"
                                        action              = "snat"
                                        description         = "rule3"
                                        original_address    = "1.1.1.3"
                                        translated_address  = "192.168.1.1"
                                    }`, edgeid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nsx_nat_rule.rule3", "edgeid", edgeid),
					resource.TestCheckResourceAttrSet("nsx_nat_rule.rule3", "id"),
					testAccNatRuleWithDescriptionExists(edgeid, "rule3"),
				),
			},
		},
	})
}

func testAccResourceNatRulesAreEmpty(edgeid string) error {
	nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

	api := nat.NewGetAll(edgeid)
	err := nsxClient.Do(api)
	if err != nil {
		return err
	}

	natconfig := api.GetResponse()

	if len(natconfig.Rules.Rules) > 0 {
		return fmt.Errorf("There are still NAT Rules")
	}

	return nil
}

func testAccNatRuleWithDescriptionExists(edgeid string, description string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

		api := nat.NewGetAll(edgeid)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		natconfig := api.GetResponse()

		for _, natrule := range natconfig.Rules.Rules {
			if natrule.Description == description {
				return nil
			}
		}

		return fmt.Errorf("NSX nat rule with description %s wasn't found", description)
	}
}

func testAccNatRuleWithDescriptionNotExists(edgeid string, description string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*gonsx.NSXClient)

		api := nat.NewGetAll(edgeid)
		err := nsxClient.Do(api)
		if err != nil {
			return err
		}

		natconfig := api.GetResponse()

		for _, natrule := range natconfig.Rules.Rules {
			if natrule.Description == description {
				return fmt.Errorf("NSX nat rule with description %s was found, when it shouldn't.", description)
			}
		}

		return nil
	}
}
