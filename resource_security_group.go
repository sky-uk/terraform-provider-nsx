package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitygroup"
	"log"
)

func getSingleSecurityGroup(scopeID, name string, nsxclient *gonsx.NSXClient) (*securitygroup.SecurityGroup, error) {
	getAllAPI := securitygroup.NewGetAll(scopeID)
	err := nsxclient.Do(getAllAPI)

	if err != nil {
		return nil, err
	}

	if getAllAPI.StatusCode() != 200 {
		return nil, fmt.Errorf("Status code: %d, Response: %s", getAllAPI.StatusCode(), getAllAPI.ResponseObject())
	}

	securityGroup := getAllAPI.GetResponse().FilterByName(name)
	return securityGroup, nil
}

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,

		Schema: map[string]*schema.Schema{

			"scopeid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"dynamicmembership": &schema.Schema{
				Type: schema.TypeList,
				// When adding exclusion / inclusion lists we'll want to make this optional.
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"membershipcriterialist": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"membershipoperator": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"criteriadetails": &schema.Schema{
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"match": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"key": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"value": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"criteria": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getDynamicMemberDefinitionFromTemplate(templateDynamicMemberDefinition interface{}) (securitygroup.DynamicMemberDefinition, string) {

	var error string = ""

	var newDynamicMemberDefinition securitygroup.DynamicMemberDefinition
	for _, value := range templateDynamicMemberDefinition.([]interface{}) {
		dynamicMemberDefinition := value.(map[string]interface{})
		if v, ok := dynamicMemberDefinition["membershipcriterialist"]; ok {
			dynamicSetList, error := getDynamicSetsFromTemplate(v, len(v.([]interface{})))
			if error != "" {
				return newDynamicMemberDefinition, error
			}
			newDynamicMemberDefinition.DynamicSet = dynamicSetList
		} else {
			error = "dynamicmembership requires membershipcriterialist list"
			return newDynamicMemberDefinition, error
		}
	}
	return newDynamicMemberDefinition, error
}

func getDynamicSetsFromTemplate(templateDynamicSets interface{}, numberDynamicSets int) ([]securitygroup.DynamicSet, string) {

	var error string = ""

	newDynamicSets := make([]securitygroup.DynamicSet, numberDynamicSets)
	for index, value := range templateDynamicSets.([]interface{}) {
		dynamicSets := value.(map[string]interface{})
		if v, ok := dynamicSets["membershipoperator"].(string); ok && v != "" {
			newDynamicSets[index].Operator = v
		} else {
			error = "dynamicmembership -> membershipcriterialist requires membershipoperator arguement"
			return newDynamicSets, error
		}
		if v, ok := dynamicSets["criteriadetails"]; ok {
			newDynamicCriteriaList, error := getDynamicCriterionFromTemplate(v, len(v.([]interface{})))
			if error != "" {
				return newDynamicSets, error
			}
			newDynamicSets[index].DynamicCriteria = newDynamicCriteriaList
		} else {
			error = "dynamicmembership -> membershipcriterialist requires criteria list"
			return newDynamicSets, error
		}
	}
	return newDynamicSets, error
}

func getDynamicCriterionFromTemplate(templateDynamicCriterion interface{}, numberDynamicCriteria int) ([]securitygroup.DynamicCriteria, string) {

	var error string = ""

	newDynamicCriterion := make([]securitygroup.DynamicCriteria, numberDynamicCriteria)
	for index, value := range templateDynamicCriterion.([]interface{}) {
		dynamicCriterion := value.(map[string]interface{})
		if v, ok := dynamicCriterion["match"].(string); ok && v != "" {
			newDynamicCriterion[index].Operator = v
		} else {
			error = "dynamicmembership -> membershipcriterialist -> criteriadetails requires match arguement"
			return newDynamicCriterion, error
		}
		if v, ok := dynamicCriterion["key"].(string); ok && v != "" {
			newDynamicCriterion[index].Key = v
		} else {
			error = "dynamicmembership -> membershipcriterialist -> criteriadetails requires key arguement"
			return newDynamicCriterion, error
		}
		if v, ok := dynamicCriterion["value"].(string); ok && v != "" {
			newDynamicCriterion[index].Value = v
		} else {
			error = "dynamicmembership -> membershipcriterialist -> criteriadetails requires value arguement"
			return newDynamicCriterion, error
		}
		if v, ok := dynamicCriterion["criteria"].(string); ok && v != "" {
			newDynamicCriterion[index].Criteria = v
		} else {
			error = "dynamicmembership -> membershipcriterialist -> criteriadetails requires criteria arguement"
			return newDynamicCriterion, error
		}
	}
	return newDynamicCriterion, error
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name, error string
	var dynamicMemberDefinition securitygroup.DynamicMemberDefinition

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("dynamicmembership"); ok {
		dynamicMemberDefinition, error = getDynamicMemberDefinitionFromTemplate(v)
		log.Printf(fmt.Sprintf("The errors are %s", error))
		if error != "" {
			return fmt.Errorf(error)
		}
	} else {
		return fmt.Errorf("dynamicmembership list is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewCreate(%s, %s, %v", scopeid, name, dynamicMemberDefinition))
	createAPI := securitygroup.NewCreate(scopeid, name, &dynamicMemberDefinition)
	err := nsxclient.Do(createAPI)

	if err != nil {
		return fmt.Errorf("Error creating security group: %v", err)
	}

	if createAPI.StatusCode() != 201 {
		return fmt.Errorf("%s", createAPI.ResponseObject())
	}

	d.SetId(createAPI.GetResponse())
	return resourceSecurityGroupRead(d, m)
}

func resourceSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name string

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	securityGroupObject, err := getSingleSecurityGroup(scopeid, name, nsxclient)
	if err != nil {
		return err
	}

	id := securityGroupObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
	}

	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {

	var scopeid string
	nsxclient := m.(*gonsx.NSXClient)
	hasChanges := false

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewGetAll(%s)", scopeid))
	oldName, newName := d.GetChange("name")
	securityGroupObject, err := getSingleSecurityGroup(scopeid, oldName.(string), nsxclient)
	id := securityGroupObject.ObjectID

	//TODO: change attributes other than name. Requires changes in gonsx.
	if d.HasChange("name") {
		hasChanges = true
		securityGroupObject.Name = newName.(string)
		log.Printf(fmt.Sprintf("[DEBUG] Changing name of security group from %s to %s", oldName.(string), newName.(string)))
	}

	if hasChanges {
		updateAPI := securitygroup.NewUpdate(id, securityGroupObject)
		err = nsxclient.Do(updateAPI)
		if err != nil {
			log.Printf(fmt.Sprintf("[DEBUG] Error updating security group: %s", err))
		}
	}
	return resourceSecurityGroupRead(d, m)
}

func resourceSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var name, scopeid string

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return fmt.Errorf("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	// Gather all the resources that are associated with the specified
	// scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewGetAll(%s)", scopeid))
	api := securitygroup.NewGetAll(scopeid)
	err := nsxclient.Do(api)

	if err != nil {
		return err
	}

	// See if we can find our specifically named resource within the list of
	// resources associated with the scopeid.
	log.Printf(fmt.Sprintf("[DEBUG] api.GetResponse().FilterByName(\"%s\").ObjectID", name))
	securityGroupObject, err := getSingleSecurityGroup(scopeid, name, nsxclient)
	id := securityGroupObject.ObjectID
	log.Printf(fmt.Sprintf("[DEBUG] id := %s", id))

	// If the resource has been removed manually, notify Terraform of this fact.
	if id == "" {
		d.SetId("")
		return nil
	}

	// If we got here, the resource exists, so we attempt to delete it.
	deleteAPI := securitygroup.NewDelete(id)
	err = nsxclient.Do(deleteAPI)

	if err != nil {
		return err
	}

	// If we got here, the resource had existed, we deleted it and there was
	// no error.  Notify Terraform of this fact and return successful
	// completion.
	d.SetId("")
	log.Printf(fmt.Sprintf("[DEBUG] id %s deleted.", id))

	return nil
}
