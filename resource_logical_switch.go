package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/virtualwire"
	"net/http"
	"regexp"
)

func resourceLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalSwitchCreate,
		Read:   resourceLogicalSwitchRead,
		Update: resourceLogicalSwitchUpdate,
		Delete: resourceLogicalSwitchDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the logical switch",
			},
			"desc": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of the logical switch",
			},
			"controlplanemode": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The control plane mode to use with the logical switch. Typically this will be UNICAST_MODE. Other valid options are HYBRID_MODE and MULTICAST_MODE",
				ValidateFunc: validateLogicalSwitchControlPlaneMode,
			},
			"tenantid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Virtual Wire tenant ID. Can't be changed after creation",
			},
			"scopeid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The transport zone ID. Only required for creation.",
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The virtual wire labels",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func validateLogicalSwitchControlPlaneMode(v interface{}, k string) (ws []string, errors []error) {
	controlPlaneMode := v.(string)
	controlPlaneModeOptions := regexp.MustCompilePOSIX(`^(UNICAST_MODE|HYBRID_MODE|MULTICAST_MODE)$`)
	if !controlPlaneModeOptions.MatchString(controlPlaneMode) {
		errors = append(errors, fmt.Errorf("%q must be one of UNICAST_MODE, HYBRID_MODE or MULTICAST_MODE", k))
	}
	return
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*gonsx.NSXClient)
	var scopeID string
	var logicalSwitchCreate virtualwire.CreateSpec

	if v, ok := d.GetOk("name"); ok && v != "" {
		logicalSwitchCreate.Name = v.(string)
	} else {
		return fmt.Errorf("Error logical switch create: name attribute is required")
	}
	if v, ok := d.GetOk("desc"); ok && v != "" {
		logicalSwitchCreate.Description = v.(string)
	} else {
		return fmt.Errorf("Error logical switch create: desc attribute is required")
	}
	if v, ok := d.GetOk("controlplanemode"); ok && v != "" {
		logicalSwitchCreate.ControlPlaneMode = v.(string)
	} else {
		return fmt.Errorf("Error logical switch create: controlplanemode attribute is required")
	}
	if v, ok := d.GetOk("tenantid"); ok && v != "" {
		logicalSwitchCreate.TenantID = v.(string)
	} else {
		return fmt.Errorf("Error logical switch create: tenantid attribute is required")
	}
	if v, ok := d.GetOk("scopeid"); ok && v != "" {
		scopeID = v.(string)
	} else {
		return fmt.Errorf("Error logical switch create: scopeid attribute is required")
	}

	createAPI := virtualwire.NewCreate(logicalSwitchCreate, scopeID)
	err := nsxClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Error while creating logical switch %s: %v", logicalSwitchCreate.Name, err)
	}
	createResponseCode := createAPI.StatusCode()
	createResponse := createAPI.GetResponse()
	if createResponseCode != http.StatusCreated {
		return fmt.Errorf("Error while creating logical switch %s. Invalid HTTP response code %d received. Response: %v", logicalSwitchCreate.Name, createResponseCode, createResponse)
	}

	// NSX API returns the virtualwire ID as a string on successful creation and nothing else. E.g. virtualwire-101
	d.SetId(createResponse)
	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*gonsx.NSXClient)
	logicalSwitchID := d.Id()
	if logicalSwitchID == "" {
		return fmt.Errorf("Error obtaining logical switch ID from state during read")
	}

	getAPI := virtualwire.NewGet(logicalSwitchID)
	err := nsxClient.Do(getAPI)
	if err != nil {
		return fmt.Errorf("Error while reading logical switch ID %s. Error: %v", logicalSwitchID, err)
	}
	if getAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	logicalSwitch := getAPI.GetResponse()
	d.SetId(logicalSwitch.ObjectID)
	d.Set("name", logicalSwitch.Name)
	d.Set("desc", logicalSwitch.Description)
	d.Set("controlplanemode", logicalSwitch.ControlPlaneMode)
	d.Set("tenantid", logicalSwitch.TenantID)

	labelList := make([]string, 0)
	for _, context := range logicalSwitch.VdsContext {
		labelName := "vxw-" + context.Switch.ObjectID + "-" + logicalSwitch.ObjectID + "-sid-" + logicalSwitch.VdnID + "-" + logicalSwitch.Name
		labelList = append(labelList, labelName)
	}

	d.Set("labels", labelList)

	return nil
}

func resourceLogicalSwitchUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*gonsx.NSXClient)
	var updateVirtualWire virtualwire.VirtualWire
	hasChanges := false
	updateVirtualWire.ObjectID = d.Id()

	if updateVirtualWire.ObjectID == "" {
		return fmt.Errorf("Error obtaining logical switch ID from state during update")
	}

	// We need to send all attributes on each update, not just the changes.
	if v, ok := d.GetOk("name"); ok && v != "" {
		if d.HasChange("name") {
			hasChanges = true
		}
		updateVirtualWire.Name = v.(string)
	} else {
		return fmt.Errorf("Error logical switch update: name attribute is required")
	}
	if v, ok := d.GetOk("desc"); ok && v != "" {
		if d.HasChange("desc") {
			hasChanges = true
		}
		updateVirtualWire.Description = v.(string)
	} else {
		return fmt.Errorf("Error logical switch update: desc attribute is required")
	}
	if v, ok := d.GetOk("controlplanemode"); ok && v != "" {
		if d.HasChange("controlplanemode") {
			hasChanges = true
		}
		updateVirtualWire.ControlPlaneMode = v.(string)
	} else {
		return fmt.Errorf("Error logical switch update: controlplanemode attribute is required")
	}
	// Tenant ID can't be updated. Update will return as successful, but doesn't change attribute.
	if v, ok := d.GetOk("tenantid"); ok && v != "" {
		if d.HasChange("tenantid") {
			hasChanges = true
		}
		updateVirtualWire.TenantID = v.(string)
	} else {
		return fmt.Errorf("Error logical switch update: tenantid attribute is required")
	}

	if hasChanges {
		updateLogicalSwitchAPI := virtualwire.NewUpdate(updateVirtualWire)
		err := nsxClient.Do(updateLogicalSwitchAPI)
		if err != nil {
			return fmt.Errorf("Error while updating logical switch ID %s. Error %v", updateVirtualWire.ObjectID, err)
		}
		responseCode := updateLogicalSwitchAPI.StatusCode()
		if responseCode != http.StatusOK {
			return fmt.Errorf("Error while updating logical switch ID %s - invalid HTTP resposne code %d received", updateVirtualWire.ObjectID, responseCode)
		}
		// NSX API doesn't return any content when change is successful. Setting values read in from the template.
		d.SetId(updateVirtualWire.ObjectID)
		d.Set("name", updateVirtualWire.Name)
		d.Set("desc", updateVirtualWire.Description)
		d.Set("controlplanemode", updateVirtualWire.ControlPlaneMode)
		d.Set("tenantid", updateVirtualWire.TenantID)
	}
	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*gonsx.NSXClient)
	virtualWireID := d.Id()

	if virtualWireID == "" {
		return fmt.Errorf("Error obtaining logical switch ID from state during delete")
	}

	deleteAPI := virtualwire.NewDelete(virtualWireID)
	err := nsxClient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf("Error while deleting logical switch ID %s. Error: %v", virtualWireID, err)
	}

	responseCode := deleteAPI.StatusCode()
	if responseCode != http.StatusOK && responseCode != http.StatusNotFound {
		return fmt.Errorf("Error while deleting logical switch ID %s. Received invalid HTTP response code %d", virtualWireID, responseCode)
	}

	d.SetId("")
	return nil
}
