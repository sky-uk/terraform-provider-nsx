package main

import (
	"log"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitypolicy"
)

func resourceSecurityPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityPolicyRuleCreate,
		Read:   resourceSecurityPolicyRuleRead,
		Delete: resourceSecurityPolicyRuleDelete,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"securitypolicyname": {
				Type: 	  schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"securitygroupids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, securitypolicyname, action, direction string
	var securitygroupids []string

	// Gather the attributes for the resource.

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("securitypolicyname"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("securitypolicyname argument is required")
	}

	if v, ok := d.GetOk("action"); ok {
		action = v.(string)
	} else {
		return fmt.Errorf("action argument is required")
	}

	if v, ok := d.GetOk("direction"); ok {
		direction = v.(string)
	} else {
		return fmt.Errorf("direction argument is required")
	}

	if v, ok := d.GetOk("securitygroupids"); ok {
		list := v.([]interface{})

		securitygroupids = make([]string, len(list))
		for i, value := range list {
			groupID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in securitygroups")
			}
			securitygroupids[i] = groupID
		}
	} else {
		return fmt.Errorf("securitygroupids argument is required")
	}


	policyToModify, err := getSingleSecurityPolicy(securitypolicyname, nsxclient)


	if err != nil {
		return err
	}

	log.Printf(fmt.Sprintf("[DEBUG] policyToModify.AddOutboundFirewallAction(%s, %s, %s, %s)",name, action, direction, securitygroupids))
	policyToModify.AddOutboundFirewallAction(name, action, direction, securitygroupids)
	updateAPI := securitypolicy.NewUpdate(policyToModify.ObjectID, policyToModify)
	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error creating security group: %v", err)
	}

	if updateAPI.StatusCode() != 200 {
		return fmt.Errorf("%s", updateAPI.ResponseObject())
	}

	d.SetId(policyToModify.VsmUUID)
	return resourceSecurityPolicyRuleRead(d, m)
}

func resourceSecurityPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSecurityPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}