package main

import (
	"errors"
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
			"dynamic_membership": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"set_operator": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSecurityGroupSetOperator,
						},
						"rules_operator": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSecurityGroupRulesOperator,
						},
						"rules": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateSecurityGroupRuleKey,
									},
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"criteria": &schema.Schema{
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateSecurityGroupRuleCriteria,
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

func validateSecurityGroupSetOperator(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "OR" && value != "AND" {
		errors = append(errors, fmt.Errorf("%q must be one of \"OR\" or \"AND\" ", k))
	}
	return
}

func validateSecurityGroupRulesOperator(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "OR" && value != "AND" {
		errors = append(errors, fmt.Errorf("%q must be one of \"OR\" or \"AND\" ", k))
	}
	return
}

func validateSecurityGroupRuleKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	keyTypes := map[string]bool{
		"VM.SECURITY_TAG":       true,
		"VM.GUEST_OS_FULL_NAME": true,
		"VM.GUEST_HOST_NAME":    true,
		"VM.NAME":               true,
		"ENTITY":                true,
	}
	if !keyTypes[value] {
		errors = append(errors, fmt.Errorf("%q must be a valid key, check documentation for acceptable values", k))
	}
	return
}

func validateSecurityGroupRuleCriteria(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	possibleCriteriaValues := map[string]bool{
		"belongs_to":  true,
		"starts_with": true,
		"ends_with":   true,
		"=":           true,
		"!=":          true,
		"contains":    true,
		"similar_to":  true,
	}
	if !possibleCriteriaValues[value] {
		errors = append(errors, fmt.Errorf("%q must be a valid criteria value, check documentation for acceptable values", k))
	}
	return
}

func buildDynamicMemberDefinition(m interface{}) (*securitygroup.DynamicMemberDefinition, error) {
	newDynamicMemberDefinition := &securitygroup.DynamicMemberDefinition{
		DynamicSet: make([]securitygroup.DynamicSet, 0),
	}

	dynamicSetList := make([]securitygroup.DynamicSet, len(m.([]interface{})))
	for index, v := range m.([]interface{}) {
		data := v.(map[string]interface{})
		dynamicRulesList, err := buildDynamicRules(data["rules"], data["rules_operator"].(string))
		if err != nil {
			return newDynamicMemberDefinition, err
		}
		dynamicSetList[index].Operator = data["set_operator"].(string)
		dynamicSetList[index].DynamicCriteria = dynamicRulesList
		log.Printf(fmt.Sprintf("[DEBUG] DynamicSetList: %v", dynamicSetList))

	}
	newDynamicMemberDefinition.DynamicSet = dynamicSetList
	return newDynamicMemberDefinition, nil
}

func buildDynamicRules(m interface{}, rulesOperator string) ([]securitygroup.DynamicCriteria, error) {
	newDynamicCriterion := make([]securitygroup.DynamicCriteria, len(m.(*schema.Set).List()))
	for index, value := range m.(*schema.Set).List() {
		dynamicCriterion := value.(map[string]interface{})
		newDynamicCriterion[index].Operator = rulesOperator
		newDynamicCriterion[index].Key = dynamicCriterion["key"].(string)
		newDynamicCriterion[index].Value = dynamicCriterion["value"].(string)
		newDynamicCriterion[index].Criteria = dynamicCriterion["criteria"].(string)
	}
	return newDynamicCriterion, nil
}

func resourceSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {

	nsxclient := m.(*gonsx.NSXClient)
	var scopeid, name string
	var dynamicMemberDefinition *securitygroup.DynamicMemberDefinition
	var err error

	// Gather the attributes for the resource.
	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
	}

	if v, ok := d.GetOk("dynamic_membership"); ok {
		log.Printf(fmt.Sprintf("[DEBUG] dynamic_membership create : %+v", v))
		dynamicMemberDefinition, err = buildDynamicMemberDefinition(v)
		if err != nil {
			return err
		}
		// dynamicMemberDefinition, err = getDynamicMemberDefinitionFromTemplate(v)
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewCreate(%s, %s, %v", scopeid, name, &dynamicMemberDefinition))
	createAPI := securitygroup.NewCreate(scopeid, name, dynamicMemberDefinition)
	err = nsxclient.Do(createAPI)

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
	var dynamicMembership *securitygroup.DynamicMemberDefinition
	var scopeid, name string
	var err error

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
	}

	if v, ok := d.GetOk("dynamic_membership"); ok {
		dynamicMembership, err = buildDynamicMemberDefinition(v)
		if err != nil {
			return err
		}
	} else {
		dynamicMembership = &securitygroup.DynamicMemberDefinition{
			DynamicSet: make([]securitygroup.DynamicSet, 0),
		}
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
		return nil
	}

	if securityGroupObject.DynamicMemberDefinition == nil {
		return nil
	}
	log.Printf(fmt.Sprintf("[DEBUG] dynamicMembership := %v", securityGroupObject.DynamicMemberDefinition))
	for idx, remoteDynamicSet := range securityGroupObject.DynamicMemberDefinition.DynamicSet {
		dynamicMembership.DynamicSet[idx].Operator = remoteDynamicSet.Operator
		readDynamicCriteria(dynamicMembership.DynamicSet[idx].DynamicCriteria,
			remoteDynamicSet.DynamicCriteria)
	}
	d.Set("dynamic_membership", dynamicMembership)
	return nil
}

func readDynamicCriteria(localCriteriaList, remoteCriteriaList []securitygroup.DynamicCriteria) {
	for _, localRule := range localCriteriaList {
		for _, remoteRule := range remoteCriteriaList {
			if localRule.Value == remoteRule.Value {
				localRule.Criteria = remoteRule.Criteria
				localRule.Key = remoteRule.Key
				break
			}
		}
	}
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {

	var scopeid string
	var dynamicMembership *securitygroup.DynamicMemberDefinition
	var err error

	nsxclient := m.(*gonsx.NSXClient)
	hasChanges := false

	if v, ok := d.GetOk("scopeid"); ok {
		scopeid = v.(string)
	} else {
		return errors.New("scopeid argument is required")
	}

	log.Printf(fmt.Sprintf("[DEBUG] securitygroup.NewGetAll(%s)", scopeid))
	oldName, newName := d.GetChange("name")
	securityGroupObject, err := getSingleSecurityGroup(scopeid, oldName.(string), nsxclient)
	id := securityGroupObject.ObjectID

	// TODO: change attributes other than name. Requires changes in gonsx.
	if d.HasChange("name") {
		hasChanges = true
		securityGroupObject.Name = newName.(string)
		log.Printf(fmt.Sprintf("[DEBUG] Changing name of security group from %s to %s", oldName.(string), newName.(string)))
	}

	if d.HasChange("dynamic_membership") {
		if v, ok := d.GetOk("dynamic_membership"); ok {
			dynamicMembership, err = buildDynamicMemberDefinition(v)
			if err != nil {
				return err
			}
		}
		hasChanges = true
		securityGroupObject.DynamicMemberDefinition = dynamicMembership
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
		return errors.New("scopeid argument is required")
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return errors.New("name argument is required")
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
