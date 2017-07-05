package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/fwrules"
	"log"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Schema: map[string]*schema.Schema{

		},
	}

}

func resourceFirewallRuleCreate(d *schema.ResourceData, m interface{}) error {
	return nil

}

func resourceFirewallRuleRead(d *schema.ResourceData, m interface{}) error {
	return nil

}

func resourceFirewallRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil

}

func resourceFirewallRuleDelete(d *schema.ResourceData, m interface{}) error {
	return nil

}
