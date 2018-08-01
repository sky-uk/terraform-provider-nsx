package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitypolicy"
	"log"
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
				Type:     schema.TypeString,
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
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"serviceids": {
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
	var securitygroupids, serviceids []string

	// Gather the attributes for the resource.

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("securitypolicyname"); ok {
		securitypolicyname = v.(string)
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
		if action == "outbound" {
			return fmt.Errorf("securitygroupids argument is required")
		}
	}

	if v, ok := d.GetOk("serviceids"); ok {
		list := v.([]interface{})

		serviceids = make([]string, len(list))
		for i, value := range list {
			serviceID, ok := value.(string)
			if !ok {
				return fmt.Errorf("empty element found in services")
			}
			serviceids[i] = serviceID
		}
	} else {
		return fmt.Errorf("serviceids argument is required")
	}

	log.Print("Getting policy object to modify")
	policyToModify, err := getSingleSecurityPolicy(securitypolicyname, nsxclient)
	log.Printf("[DEBUG] - policyTOModify :%s", policyToModify)

	if err != nil {
		return err
	}

	existingAction := policyToModify.GetFirewallRuleByName(name)
	if existingAction.Name != "" {
		return fmt.Errorf("Firewall rule with same name already exists in this security policy")
	}

	if direction == "inbound" {
		log.Printf(fmt.Sprintf("[DEBUG] policyToModify.AddInboundFirewallAction(%s, %s, %s, %s)", name, action, direction, serviceids))
		err = policyToModify.AddInboundFirewallAction(name, action, direction, serviceids)
		if err != nil {
			return fmt.Errorf("Error in adding the rule to policy object: %s", err)
		}
	} else {
		log.Printf(fmt.Sprintf("[DEBUG] policyToModify.AddOutboundFirewallAction(%s, %s, %s, %s, %s)", name, action, direction, securitygroupids, serviceids))
		err = policyToModify.AddOutboundFirewallAction(name, action, direction, securitygroupids, serviceids)
		if err != nil {
			return fmt.Errorf("Error in adding the rule to policy object: %s", err)
		}
	}

	log.Printf("[DEBUG] - policyTOModify :%s", policyToModify)
	policyToModify.Revision += policyToModify.Revision
	updateAPI := securitypolicy.NewUpdate(policyToModify.ObjectID, policyToModify)

	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error creating security group: %v", err)
	}

	if updateAPI.StatusCode() != 200 {
		return fmt.Errorf("%s", updateAPI.ResponseObject())
	}

	d.SetId(name)
	return resourceSecurityPolicyRuleRead(d, m)
}

func resourceSecurityPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name string
	var securitypolicyname string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}
	if v, ok := d.GetOk("securitypolicyname"); ok {
		securitypolicyname = v.(string)
	} else {
		return fmt.Errorf("securitypolicyname argument is required")
	}

	policyToRead, err := getSingleSecurityPolicy(securitypolicyname, nsxclient)
	log.Printf("[DEBUG] - policyToRead :%s", policyToRead)

	if err != nil {
		return err
	}

	existingAction := policyToRead.GetFirewallRuleByName(name)
	id := existingAction.VsmUUID
	log.Printf(fmt.Sprintf("[DEBUG] VsmUUID := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}
	return nil
}

func resourceSecurityPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name string
	var securityPolicyName string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("securitypolicyname"); ok {
		securityPolicyName = v.(string)
	} else {
		return fmt.Errorf("securitypolicyname argument is required")
	}

	log.Print("Getting policy object to modify")
	policyToModify, err := getSingleSecurityPolicy(securityPolicyName, nsxclient)
	log.Printf("[DEBUG] - policyTOModify :%s", policyToModify)

	if err != nil {
		return err
	}

	log.Printf(fmt.Sprintf("[DEBUG] policyToModify.Remove(%s)", name))
	// FIXME:  RemoveFirewallActionByName probably return a error for consistency
	policyToModify.RemoveFirewallActionByName(name)
	log.Printf("[DEBUG] - policyTOModify :%s", policyToModify)
	updateAPI := securitypolicy.NewUpdate(policyToModify.ObjectID, policyToModify)

	err = nsxclient.Do(updateAPI)

	if err != nil {
		return fmt.Errorf("Error creating security group: %v", err)
	}

	if updateAPI.StatusCode() != 200 {
		return fmt.Errorf("%s", updateAPI.ResponseObject())
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")
	log.Printf(fmt.Sprintf("[DEBUG] firewall rule with name %s from securitypolicy %s deleted.", name, securityPolicyName))

	return nil
}
