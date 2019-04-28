package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/nat"

	uuid "github.com/hashicorp/go-uuid"
)

var validateNSXIPAddress = validation.Any(validation.SingleIP(), validation.IPRange(), validation.CIDRNetwork(0, 32), validation.StringMatch(regexp.MustCompile(`^any$`), "value must be any"))

func resourceNatRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatRuleCreate,
		Read:   resourceNatRuleRead,
		Update: resourceNatRuleUpdate,
		Delete: resourceNatRuleDelete,

		Schema: map[string]*schema.Schema{
			"edgeid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					if value != "dnat" && value != "snat" {
						errs = append(errs, fmt.Errorf("%q must be one of \"dnat\" or \"snat\" ", key))
					}
					return
				},
				Description: "Type of NAT: snat or dnat",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If Rule is Enabled",
			},

			"logging_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If Logging is Enabled for this Rule",
			},

			"vnic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Interface on which the translating is applied.",
			},

			"original_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "any",
				Description:  "Original address or address range. Destination address for DNAT rules.",
				ValidateFunc: validateNSXIPAddress,
			},

			"translated_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Translated address or address range.",
			},

			"dnat_match_source_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Source address to match for DNAT rules.",
			},

			"snat_match_destination_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Destionation address to match for SNAT rules.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "any",
				Description: "Protocol.",
			},

			"icmp_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ICMP type. Only supported when protocol is icmp.",
			},

			"original_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "any",
				Description:  "Original Destination port.",
				ValidateFunc: validateNSXPortOrAny(),
			},

			"translated_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "any",
				Description:  "Translated Port.",
				ValidateFunc: validateNSXPortOrAny(),
			},

			"dnat_match_source_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "any",
				Description:  "Source port for DNAT rules.",
				ValidateFunc: validateNSXPortOrAny(),
			},

			"snat_match_destination_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "any",
				Description:  "Destionation port for SNAT rules.",
				ValidateFunc: validateNSXPortOrAny(),
			},
		},
	}
}

func printRule(rule *nat.Rule) {
	rule_xml, err := xml.MarshalIndent(rule, "", "  ")
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf(string(rule_xml))
}

func composeNatRuleId(edgeid string, natid string) string {
	var id = edgeid + ":" + natid
	return id
}

func decomposeNatRuleId(ruleid string) (string, string) {
	s := strings.Split(ruleid, ":")
	return s[0], s[1]
}

func getAllNatRules(nsxclient *gonsx.NSXClient, edgeID string) (*nat.Rules, error) {
	api := nat.NewGetAll(edgeID)
	err := nsxclient.Do(api)

	if err != nil {
		return nil, fmt.Errorf("Could not fetch NAT Rules: %s", err)
	}

	natconfig := api.GetResponse()

	return &natconfig.Rules, nil
}

func getNatRule(nsxclient *gonsx.NSXClient, id string) (*nat.Rule, error) {
	edgeid, ruleid := decomposeNatRuleId(id)
	natrules, err := getAllNatRules(nsxclient, edgeid)
	if err != nil {
		return nil, err
	}

	for _, natRule := range natrules.Rules {
		if natRule.RuleID == ruleid {
			return &natRule, nil
		}
	}

	// not finding the rule is not an error
	return nil, nil
}

func expandNatRule(d *schema.ResourceData, rule *nat.Rule) error {
	rule.Action = d.Get("action").(string)

	// Common Between dnat and snat
	if description, ok := d.GetOk("description"); ok {
		rule.Description = description.(string)
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		rule.Enabled = enabled.(bool)
	}

	if loggingEnabled, ok := d.GetOk("logging_enabled"); ok {
		rule.LoggingEnabled = loggingEnabled.(bool)
	}

	if value, ok := d.GetOk("vnic"); ok {
		rule.Vnic = value.(string)
	}

	if value, ok := d.GetOk("original_address"); ok {
		rule.OriginalAddress = value.(string)
	}

	if value, ok := d.GetOk("translated_address"); ok {
		rule.TranslatedAddress = value.(string)
	}

	if value, ok := d.GetOk("protocol"); ok {
		rule.Protocol = value.(string)
	}

	// ICMP only
	if rule.Protocol == "icmp" {
		if value, ok := d.GetOk("icmp_type"); ok {
			rule.IcmpType = value.(string)
		}
	}

	// DNAT only
	if rule.Action == "dnat" {
		if value, ok := d.GetOk("dnat_match_source_address"); ok {
			rule.DnatMatchSourceAddress = value.(string)
		}

		if value, ok := d.GetOk("dnat_match_source_port"); ok {
			rule.DnatMatchSourcePort = value.(string)
		}
	}

	// SNAT only
	if rule.Action == "snat" {
		if value, ok := d.GetOk("snat_match_destination_address"); ok {
			rule.SnatMatchDestinationAddress = value.(string)
		}

		if value, ok := d.GetOk("snat_match_destination_port"); ok {
			rule.SnatMatchDestinationPort = value.(string)
		}
	}

	if value, ok := d.GetOk("original_port"); ok {
		rule.OriginalPort = value.(string)
	}

	if value, ok := d.GetOk("translated_port"); ok {
		rule.TranslatedPort = value.(string)
	}

	return nil
}

func resourceNatRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	edgeid := d.Get("edgeid").(string)

	rule_uuid, uuid_err := uuid.GenerateUUID()
	if uuid_err != nil {
		panic(fmt.Errorf("Failed to generate uuid for NatRule: %v", uuid_err))
	}

	var rule nat.Rule
	err := expandNatRule(d, &rule)
	if err != nil {
		return err
	}
	rule.Description = "Terraform: " + rule_uuid

	// API Create
	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)

	createAPI := nat.NewCreateRule(edgeid, rule)
	err = nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if createAPI.StatusCode() != 201 {
		return fmt.Errorf("Failed to create DNAT rule: %s", createAPI.GetResponse())
	}

	// Read Back
	rules, err := getAllNatRules(nsxclient, edgeid)
	if err != nil {
		return fmt.Errorf("Could not fetch NAT Rules: %s", err)
	}

	var ruleid string = ""
	for _, natRule := range rules.Rules {
		if rule.Description == natRule.Description {
			ruleid = natRule.RuleID
			break
		}
	}

	if ruleid == "" {
		return fmt.Errorf("Failed to find created DNAT rule: %s", rule.Description)
	}

	id := composeNatRuleId(edgeid, ruleid)
	d.SetId(id)

	//TODO Update with proper description
	err = expandNatRule(d, &rule)
	if err != nil {
		return err
	}

	updateAPI := nat.NewUpdate(edgeid, ruleid, rule)
	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update NAT rule: %s", updateAPI.GetResponse())
	}

	return resourceNatRuleRead(d, m)
}

func resourceNatRuleRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var id = d.Id()
	log.Printf(fmt.Sprintf("[DEBUG] Reading NATID: |%s|", id))

	// Read Back
	rule, err := getNatRule(nsxclient, id)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	// printRule(rule)
	d.Set("action", rule.Action)
	d.Set("description", rule.Description)
	d.Set("enabled", rule.Enabled)
	d.Set("logging_enabled", rule.LoggingEnabled)
	d.Set("vnic", rule.Vnic)
	d.Set("original_address", rule.OriginalAddress)
	d.Set("translated_address", rule.TranslatedAddress)
	d.Set("protocol", rule.Protocol)
	d.Set("icmp_type", rule.IcmpType)
	d.Set("original_port", rule.OriginalPort)
	d.Set("translated_port", rule.TranslatedPort)

	// DNAT only
	if rule.Action == "dnat" {
		d.Set("dnat_match_source_address", rule.DnatMatchSourceAddress)
		d.Set("dnat_match_source_port", rule.DnatMatchSourcePort)
	}

	// SNAT only
	if rule.Action == "snat" {
		d.Set("snat_match_destination_address", rule.SnatMatchDestinationAddress)
		d.Set("snat_match_destination_port", rule.SnatMatchDestinationPort)
	}

	return nil
}

func resourceNatRuleUpdate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var id = d.Id()
	edgeid, ruleid := decomposeNatRuleId(id)

	var rule nat.Rule
	err := expandNatRule(d, &rule)
	if err != nil {
		return err
	}
	// API Update
	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)
	updateAPI := nat.NewUpdate(edgeid, ruleid, rule)
	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if updateAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to update NAT rule: %s", updateAPI.GetResponse())
	}

	return resourceNatRuleRead(d, m)
}

func resourceNatRuleDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var id = d.Id()
	edgeid, ruleid := decomposeNatRuleId(id)

	// Delete
	nsxMutexKV.Lock(edgeid)
	defer nsxMutexKV.Unlock(edgeid)
	deleteAPI := nat.NewDelete(edgeid, ruleid)
	err := nsxclient.Do(deleteAPI)

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	} else if deleteAPI.StatusCode() != 204 {
		return fmt.Errorf("Failed to delete %s", id)
	}

	d.SetId("")
	return nil
}

func validateNSXPortOrAny() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, errs []error) {
		porterror := fmt.Errorf("%q: Specify any, a port(e.g. 1234) or port range(1234-1239)", k)
		value := i.(string)

		if value == "any" {
			return
		}

		// Create list of Numbers (one or two)
		var numberstrings []string
		if strings.Contains(value, "-") {
			parts := strings.Split(value, "-")
			if len(parts) == 2 {
				numberstrings = parts
			} else {
				errs = append(errs, porterror)
				return
			}
		} else {
			numberstrings = append(numberstrings, value)
		}

		// Check all Numbers
		for _, numberstring := range numberstrings {
			number, err := strconv.Atoi(numberstring)
			if err != nil {
				errs = append(errs, porterror)
				return
			}

			if (number < 0) || (number > 65535) {
				errs = append(errs, porterror)
				return
			}
		}

		// No error
		return
	}
}
