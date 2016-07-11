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

            "name": {
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
    var desc, name, tenantid, scopeid string

    if v, ok := d.GetOk("desc"); ok {
        desc = v.(string)
    } else {
        return fmt.Errorf("desc argument is required")
    }

    if v, ok := d.GetOk("name"); ok {
        name = v.(string)
    } else {
        return fmt.Errorf("name argument is required")
    }

    if v, ok := d.GetOk("tenantid"); ok {
        tenantid = v.(string)
    } else {
        return fmt.Errorf("tenantid argument is required")
    }

    if v, ok := d.GetOk("scopeid"); ok {
        scopeid = v.(string)
    } else {
        return fmt.Errorf("scopeid argument is required")
    }

    log.Printf(fmt.Sprintf("[DEBUG] virtualwire.NewCreate(%s, %s, %s, %s)", name, desc, tenantid, scopeid))
    create_api := virtualwire.NewCreate(name, desc, tenantid, scopeid)
    nsxclient.Do(create_api)

    if create_api.StatusCode() != 201 {
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
