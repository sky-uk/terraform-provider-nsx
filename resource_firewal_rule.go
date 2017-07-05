package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/fwrules"
	//"log"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Schema: map[string]*schema.Schema{
			"ruleid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "A name for the  rule",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				ForceNew:    false,
				Description: "determines if the rule is disabled or not",
			},
			"ruletype": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of rule, valid values are LAYER2 and LAYER3",
			},
			"logged": {
				Type:        schema.TypeString,
				Optional:    false,
				ForceNew:    false,
				Description: "Should this rule be logged",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "What to do with the packets that match this rule, allow,drop, etc",
			},
			"edgeid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge Device id",
			},
			"source": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "Source of traffic for the firewall rule, it could be, CDIR, IP Set, IPv4 addresses, Virtual Machine, Vnic, Security Group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the source",
						},
						"type": {
							Optional:    false,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of source",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the source",
						},
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the source valid",
						},
					},
				},
			},
			"destination": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "Source of traffic for the firewall rule, it could be, CDIR, IP Set, IPv4 addresses, Virtual Machine, Vnic, Security Group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the source",
						},
						"type": {
							Optional:    false,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of source",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the source",
						},
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the source valid",
						},
					},
				},
			},
			"service": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the service",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the service",
						},
						"dstport": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "Destination port for the service",
						},
						"protocol": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "Protocol id ",
						},
						"subprotocol": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "SubProtocol id ",
						},
					},
				},
			},
			"sectionid": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Id for the section where the rule bellongs",
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Direction for the traffic",
			},
			"packettype": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Type of packets",
			},
		},
	}

}

func resourceFirewallRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var fwRule fwrules.Rule

	if v, ok := d.GetOk("name"); ok {
		fwRule.Name = v.(string)
	} else {
		return fmt.Errorf("Name parameter is required")
	}

	if v, ok := d.GetOk("disabled"); ok {
		fwRule.Disabled = v.(bool)
	}

	if v, ok := d.GetOk("ruletype"); ok {
		fwRule.RuleType = v.(string)
	} else {
		return fmt.Errorf("Rule Type is required")
	}

	if v, ok := d.GetOk("logged"); ok {
		fwRule.Logged = v.(string)
	}

	if v, ok := d.GetOk("action"); ok {
		fwRule.Action = v.(string)
	} else {
		return fmt.Errorf("Error needs to be set")
	}

	if v, ok := d.GetOk("edgeid"); ok {
		fwRule.EdgeID = v.(string)
	} else {
		return fmt.Errorf("Edge ID required")
	}

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
