package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/sections"
	"log"

)

func resourceFirewallSection() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallSectionCreate,
		Read:   resourceFirewallSectionRead,
		Update: resourceFirewallSectionUpdate,
		Delete: resourceFirewallSectionDelete,
		Schema: map[string]*schema.Schema{
			"sectionid" : {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name" : {
				Type: schema.TypeString,
				Required: true,
				ForceNew: false,
				Description: "A name for the section",
			},
			"type": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Type of section, can be LAYER2 or LAYER3",
			},

		},
	}
}

func resourceFirewallSectionCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var createSection sections.Section
	if v,ok := d.GetOk("name"); ok {
		createSection.Name = v.(string)
	} else {
		return fmt.Errorf("Section name is required")
	}
	if v, ok := d.GetOk("type"); ok {
		createSection.Type = v.(string)
	} else {
		return fmt.Errorf("Section Type is requried LAYER2/3")
	}

	createSectionAPI := sections.NewCreate(createSection)
	createErr := nsxclient.Do(createSectionAPI)
	if createErr != nil {
		return fmt.Errorf("could not create section")
	}

	if createSectionAPI.StatusCode() != 201 {
		log.Println(createSectionAPI.ResponseObject())
		log.Println(createSectionAPI.Endpoint())
		return fmt.Errorf("Error creating section")
	}
	d.SetId(createSectionAPI.GetResponse().ID)
	return nil
}

func resourceFirewallSectionRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFirewallSectionUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceFirewallSectionDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
