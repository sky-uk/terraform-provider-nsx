package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/fwrules"
	"github.com/sky-uk/gonsx/api/distributedfw/sections"
	"log"
)

func resourceGetSectionTimestamp(sectionID int, sectionType string, m interface{}) (*sections.Section, error) {
	sectionNSXClient := m.(*gonsx.NSXClient)
	sectionTimestamp := sections.GetSectionTimestamp(sectionID, sectionType)
	sectErr := sectionNSXClient.Do(sectionTimestamp)
	if sectErr != nil {
		log.Println("could not get timestamp")
	}
	if sectionTimestamp.StatusCode() != 200 {
		return nil, fmt.Errorf("Could not find the timestamp for section %d", sectionID)
	}
	log.Println(sectionTimestamp.ResponseObject())
	return sectionTimestamp.GetResponse(), nil
}

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
				Optional:    true,
				ForceNew:    false,
				Description: "Should this rule be logged",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "What to do with the packets that match this rule, allow,drop, etc",
			},
			"appliedto": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "Where this rule is to be applied",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the applied to",
						},
						"type": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of applied to",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the applied to",
						},
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the applied to valid",
						},
					},
				},
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
							Optional:    true,
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
							Optional:    true,
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
						"type": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of  service",
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
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the source valid",
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
				Optional:    true,
				ForceNew:    false,
				Description: "Type of packets",
			},
		},
	}

}

func resourceFirewallRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var fwRule fwrules.Rule
	var sourceList fwrules.SourceList
	var dstList fwrules.DstList
	var svcList fwrules.SvcList

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
		return fmt.Errorf("Action needs to be set")
	}

	if v, ok := d.GetOk("source"); ok {
		if sources, ok := v.(*schema.Set); ok {
			for _, source := range sources.List() {
				sourceObject := source.(map[string]interface{})
				newSource := fwrules.Source{}

				if sourceNameValue, ok := sourceObject["name"]; ok {
					newSource.Name = sourceNameValue.(string)
				}

				if sourceTypeValue, ok := sourceObject["type"]; ok {
					newSource.Type = sourceTypeValue.(string)
				}

				if sourceValue, ok := sourceObject["value"]; ok {
					newSource.Value = sourceValue.(string)
				}

				if sourceValidValue, ok := sourceObject["isvalid"]; ok {
					isValid := sourceValidValue.(bool)
					newSource.IsValid = &isValid
				}

				sourceList.Sources = append(sourceList.Sources, newSource)
				fwRule.Sources = &sourceList
			}

		}
	} else {
		return fmt.Errorf("Source  is required")
	}

	if v, ok := d.GetOk("destination"); ok {
		if destinations, ok := v.(*schema.Set); ok {
			for _, destination := range destinations.List() {
				destinationObject := destination.(map[string]interface{})
				newDestination := fwrules.Destination{}

				if destinationNameValue, ok := destinationObject["name"]; ok {
					newDestination.Name = destinationNameValue.(string)
				}

				if destinationTypeValue, ok := destinationObject["type"]; ok {
					newDestination.Type = destinationTypeValue.(string)
				}

				if destinationValue, ok := destinationObject["value"]; ok {
					newDestination.Value = destinationValue.(string)
				}

				if destinationValidValue, ok := destinationObject["isvalid"]; ok {
					newDestination.IsValid = destinationValidValue.(bool)
				}

				dstList.Destinations = append(dstList.Destinations, newDestination)
				fwRule.Destinations = &dstList
			}
		}
	}

	if v, ok := d.GetOk("service"); ok {
		if services, ok := v.(*schema.Set); ok {
			for _, service := range services.List() {
				serviceObject := service.(map[string]interface{})
				newService := fwrules.Service{}

				if serviceNameValue := serviceObject["name"]; ok {
					newService.Name = serviceNameValue.(string)

				}

				if serviceValue := serviceObject["value"]; ok {
					newService.Value = serviceValue.(string)
				}

				if serviceDstPortValue := serviceObject["dstport"]; ok {
					newService.DestinationPort = serviceDstPortValue.(int)
				}

				if serviceProtocolValue := serviceObject["protocol"]; ok {
					newService.Protocol = serviceProtocolValue.(int)
				}

				if serviceSubProtocolValue := serviceObject["subprotocol"]; ok {
					newService.SubProtocol = serviceSubProtocolValue.(int)
				}
				svcList.Services = append(svcList.Services, newService)
				fwRule.Services = &svcList
			}
		}

	}
	/*if v, ok := d.GetOk("appliedto"); ok {


	} else {
		return fmt.Errorf("no applied to")
	}*/
	if v, ok := d.GetOk("sectionid"); ok {

		fwRule.SectionID = v.(int)

	} else {
		return fmt.Errorf("Section ID is required")
	}

	if v, ok := d.GetOk("direction"); ok {
		fwRule.Direction = v.(string)
	} else {
		return fmt.Errorf("Direction is required")
	}
	if v, ok := d.GetOk("packettype"); ok {
		fwRule.PacketType = v.(string)
	} else {
		return fmt.Errorf("PacketType is required")
	}

	nsxMutexKV.Lock(fwRule.Name)
	defer nsxMutexKV.Unlock(fwRule.Name)
	createFWRuleAPI := fwrules.NewCreate(fwRule)
	timeStampCall, _ := resourceGetSectionTimestamp(fwRule.SectionID, fwRule.RuleType, m)
	log.Println(len(timeStampCall.Timestamp))
	nsxclient.SetHeader("If-Match", timeStampCall.Timestamp)
	createErr := nsxclient.Do(createFWRuleAPI)
	if createErr != nil {
		return fmt.Errorf("Could not create firewall rule")
	}

	if createFWRuleAPI.StatusCode() != 201 {
		log.Println("Response Status Code")
		log.Println(createFWRuleAPI.StatusCode())
		log.Println("Response Object")
		log.Println(createFWRuleAPI.ResponseObject())
		log.Println("Response Endpoint")
		log.Println(createFWRuleAPI.Endpoint())
		log.Println("Trying to create this rule")
		log.Println(fwRule)
		log.Println(fwRule.Services)
		return fmt.Errorf("Unable to create firewall rule")
	}
	log.Println(createFWRuleAPI.GetResponse().RuleID)

	d.SetId(createFWRuleAPI.GetResponse().RuleID)
	return resourceFirewallRuleRead(d, m)

}

func resourceFirewallRuleRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var ruleType string
	var ruleSection int
	if v, ok := d.GetOk("ruletype"); ok {
		ruleType = v.(string)
	}

	if v, ok := d.GetOk("sectionid"); ok {
		ruleSection = v.(int)

	}
	readRuleAPI := fwrules.NewGetSingle(d.Id(), ruleType, ruleSection)
	readErr := nsxclient.Do(readRuleAPI)
	log.Println(readRuleAPI.ResponseObject())
	log.Println(readRuleAPI.Endpoint())
	if readRuleAPI.StatusCode() != 200 {
		return fmt.Errorf("Error Reading Firewall rule")
	}

	if readErr != nil {
		return fmt.Errorf("Error Reading Firewall rule")
	}

	ReadRule := readRuleAPI.GetResponse()
	d.Set("name", ReadRule.Name)
	d.Set("disabled", ReadRule.Disabled)
	d.Set("action", ReadRule.Action)
	//d.Set("appliedto", ReadRule.AppliedToList)
	if err := readSources(d, ReadRule); err != nil {
		return err
	}

	/*if err := readDestinations(d, ReadRule); err != nil {
		return err
	}

	if err := readServices(d, ReadRule); err != nil {
		return err
	}*/
	d.Set("sectionid", ReadRule.SectionID)
	d.Set("direction", ReadRule.Direction)
	d.Set("packettype", ReadRule.PacketType)
	return nil

}

func resourceFirewallRuleUpdate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	//var updateRule fwrules.Rule
	var hasChanges bool
	ruleType := d.Get("ruletype").(string)
	ruleSection := d.Get("sectionid").(int)
	updateRuleAPI := fwrules.NewGetSingle(d.Id(), ruleType, ruleSection)
	getRuleErr := nsxclient.Do(updateRuleAPI)

	if getRuleErr != nil {
		return fmt.Errorf("error getting the rule for edit")
	}
	if updateRuleAPI.StatusCode() == 200 {
		//log.Println(updateRuleAPI.GetResponse())
		log.Println(updateRuleAPI.ResponseObject())
		log.Println(updateRuleAPI.Endpoint())
	}

	updateRule := updateRuleAPI.GetResponse()
	if d.HasChange("name") {
		hasChanges = true
		_, nameValue := d.GetChange("name")
		updateRule.Name = nameValue.(string)

	}
	if d.HasChange("packettype") {
		hasChanges = true
		_, packetValue := d.GetChange("packettype")
		updateRule.PacketType = packetValue.(string)
	}

	if d.HasChange("disabled") {
		hasChanges = true
		_, disabledValue := d.GetChange("disabled")
		updateRule.Disabled = disabledValue.(bool)
	}

	if d.HasChange("logged") {
		hasChanges = true
		_, loggedValue := d.GetChange("logged")
		updateRule.Logged = loggedValue.(string)
	}

	if d.HasChange("action") {
		hasChanges = true
		_, actionValue := d.GetChange("action")
		updateRule.Action = actionValue.(string)
	}

	updateRule.RuleType = d.Get("ruletype").(string)

	if hasChanges {
		nsxMutexKV.Lock(updateRule.Name)
		defer nsxMutexKV.Unlock(updateRule.Name)
		timeStampCall, _ := resourceGetSectionTimestamp(updateRule.SectionID, updateRule.RuleType, m)
		log.Println(len(timeStampCall.Timestamp))
		nsxclient.SetHeader("If-Match", timeStampCall.Timestamp)
		updateRuleAPI := fwrules.NewUpdate(*updateRule)
		updateErr := nsxclient.Do(updateRuleAPI)
		if updateErr != nil {
			return fmt.Errorf("Error updating the firewall rule ")
		}
		if updateRuleAPI.StatusCode() == 400 {
			log.Println(updateRuleAPI.StatusCode())
			log.Println(updateRuleAPI.ResponseObject())
			return fmt.Errorf("could not update the firewall rule")
		}
		log.Println(updateRuleAPI.StatusCode())
		log.Println(updateRuleAPI.ResponseObject())
		return resourceFirewallRuleRead(d, m)
	}
	return nil
}

func resourceFirewallRuleDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var ruleType string
	var sectionID int
	var deleteRule fwrules.Rule
	deleteRule.RuleID = d.Id()
	if v, ok := d.GetOk("ruletype"); ok {
		ruleType = v.(string)
		deleteRule.RuleType = ruleType
	} else {
		return fmt.Errorf("Rule Type is required")
	}

	if v, ok := d.GetOk("sectionid"); ok {
		sectionID = v.(int)
		deleteRule.SectionID = sectionID
	} else {
		return fmt.Errorf("section id is required")
	}
	timeStampCall, _ := resourceGetSectionTimestamp(deleteRule.SectionID, deleteRule.RuleType, m)
	log.Println(len(timeStampCall.Timestamp))
	nsxclient.SetHeader("If-Match", timeStampCall.Timestamp)
	deleteAPI := fwrules.NewDelete(deleteRule)
	deleteError := nsxclient.Do(deleteAPI)
	if deleteError != nil {
		return fmt.Errorf("Could not delete the rule %d", deleteAPI.StatusCode())
	}
	log.Println(deleteAPI.Endpoint())
	log.Println(deleteAPI.StatusCode())
	log.Println(deleteAPI.ResponseObject())

	d.SetId("")
	return nil
}

func readSources(d *schema.ResourceData, rule *fwrules.Rule) error {
	sources := make([]map[string]interface{}, 0)
	if v := d.Get("source"); v != nil {
		for _, sourceItem := range rule.Sources.Sources {
			source := make(map[string]interface{})
			source["name"] = sourceItem.Name
			source["type"] = sourceItem.Type
			source["value"] = sourceItem.Value
			source["isvalid"] = sourceItem.IsValid
			sources = append(sources, source)
		}
		d.Set("source", sources)
	}

	return nil
}

func readDestinations(d *schema.ResourceData, rule *fwrules.Rule) error {
	destinations := make([]map[string]interface{}, 0)
	if v := d.Get("destination"); v != nil {
		for _, destinationItem := range rule.Destinations.Destinations {
			destination := make(map[string]interface{})
			destination["name"] = destinationItem.Name
			destination["type"] = destinationItem.Type
			destination["value"] = destinationItem.Value
			destination["isvalid"] = destinationItem.IsValid
			destinations = append(destinations, destination)
			log.Println(destination)
		}
		d.Set("destination", destinations)
	}

	return nil
}

func readServices(d *schema.ResourceData, rule *fwrules.Rule) error {
	services := make([]map[string]interface{}, 0)
	if v := d.Get("service"); v != nil {
		for _, serviceItem := range rule.Services.Services {
			service := make(map[string]interface{})
			service["name"] = serviceItem.Name
			service["type"] = serviceItem.Type
			service["protocol"] = serviceItem.Protocol
			service["value"] = serviceItem.Value
			service["subprotocol"] = serviceItem.SubProtocol
			services = append(services, service)
		}
		d.Set("service", services)
	}
	return nil
}
