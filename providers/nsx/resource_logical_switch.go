package main

import (
    "errors"
    "fmt"
    "git.devops.int.ovp.bskyb.com/paas/gonsx/client"
    "git.devops.int.ovp.bskyb.com/paas/gonsx/client/api/virtualwire"
    "github.com/hashicorp/terraform/helper/schema"
    "log"
)

func resourceLogicalSwitch() *schema.Resource {
    return &schema.Resource{
        Create: resourceLogicalSwitchCreate,
        Read:   resourceLogicalSwitchRead,
        Delete: resourceLogicalSwitchDelete,

        Schema: map[string]*schema.Schema{
            "desc": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },

            "tenantid": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },

            "scopeid": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
        },
    }
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {
    nsxclient := m.(*client.NSXClient)

    if v, ok := d.GetOk("desc"); ok {
        desc := v.(string)
        log.Printf("[DEBUG] Create Logical Switch - desc := ", desc)
    } else {
        return fmt.Errorf("desc argument is required")
    }

    create_api := virtualwire.NewCreate("test", "test desc", "test", "vdnscope-19")
    nsxclient.Do(create_api)

    if create_api.StatusCode() != 201 {
        log.Printf(create_api.GetResponse())
        return errors.New(create_api.GetResponse())
    }

    d.SetId(create_api.GetResponse())
    return nil
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {
    return nil
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {
    return nil
}
