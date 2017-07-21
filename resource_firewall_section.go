package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/sections"
	"log"
	"strconv"
)

func resourceFirewallSection() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallSectionCreate,
		Read:   resourceFirewallSectionRead,
		Update: resourceFirewallSectionUpdate,
		Delete: resourceFirewallSectionDelete,
		Schema: map[string]*schema.Schema{
			"sectionid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "A name for the section",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of section, can be LAYER2 or LAYER3",
			},
		},
	}
}

func resourceFirewallSectionCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var createSection sections.Section
	if v, ok := d.GetOk("name"); ok {
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
	return resourceFirewallSectionRead(d, m)
}

func resourceFirewallSectionRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var readSection *sections.Section
	var sectionType string
	if v, ok := d.GetOk("type"); ok {
		sectionType = v.(string)
	}
	readSectionAPI := sections.NewGetSingle(d.Id(), sectionType)
	readSectionErr := nsxclient.Do(readSectionAPI)
	if readSectionErr != nil {
		return fmt.Errorf("could not get the section")
	}
	if readSectionAPI.StatusCode() != 200 {
		log.Println(readSectionAPI.StatusCode())
		return fmt.Errorf("could not find the section")
	}
	log.Println("RESPONSE")
	log.Println(readSectionAPI.GetResponse())
	log.Println("SECTION")
	log.Println(readSection)
	readSection = readSectionAPI.GetResponse()
	d.Set("name", readSection.Name)
	d.Set("type", readSection.Type)
	d.Set("sectionid", readSection.ID)
	return nil
}

func resourceFirewallSectionUpdate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var updateSection sections.Section
	var hasChanges bool
	updateSection.ID = d.Id()
	updateSection.Type = d.Get("type").(string)
	if d.HasChange("name") {
		hasChanges = true
		_, nameValue := d.GetChange("name")
		updateSection.Name = nameValue.(string)
	}
	if hasChanges {
		sectIDINT, _ := strconv.Atoi(updateSection.ID)
		timeStampCall, _ := resourceGetSectionTimestamp(sectIDINT, updateSection.Type, m)
		log.Println(len(timeStampCall.Timestamp))
		nsxclient.SetHeader("If-Match", timeStampCall.Timestamp)
		updateSectionAPI := sections.NewUpdate(updateSection)
		updateErr := nsxclient.Do(updateSectionAPI)
		if updateErr != nil {
			return fmt.Errorf("could not update section")
		}
		log.Println("UPDATE")
		log.Println(updateSectionAPI.ResponseObject())
		return resourceFirewallSectionRead(d, m)
	}

	return nil
}

func resourceFirewallSectionDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var deleteSection sections.Section
	if v, ok := d.GetOk("sectionid"); ok {
		deleteSection.ID = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		deleteSection.Type = v.(string)
	}
	deleteSectionAPI := sections.NewDelete(deleteSection)
	deleteError := nsxclient.Do(deleteSectionAPI)
	if deleteError != nil {
		return fmt.Errorf("Could not delete the section")
	}
	if deleteSectionAPI.StatusCode() != 204 {
		log.Println(deleteSectionAPI.Endpoint())
		log.Println(deleteSectionAPI.ResponseObject())
		return fmt.Errorf("Could not delete the Section")
	}

	d.SetId("")
	return nil
}
