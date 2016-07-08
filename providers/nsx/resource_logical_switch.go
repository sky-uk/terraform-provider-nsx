package main

import (
    "github.com/hashicorp/terraform/helper/schema"
)

func resourceLogicalSwitch() *schema.Resource {
    return &schema.Resource{
        Create: resourceLogicalSwitchCreate,
        Read:   resourceLogicalSwitchRead,
        Delete: resourceLogicalSwitchDelete,

        Schema: map[string]*schema.Schema{
        },
    }
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {
    return nil
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {
    return nil
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {
    return nil
}
