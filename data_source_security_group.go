package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/securitygroup"
)

func dataSourceSecurityGroup() *schema.Resource {

	return &schema.Resource{

		Read: dataSourceSecurityGroupRead,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"scopeid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "globalroot-0",
			},
		},
	}
}

func dataSourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	nsxclient := meta.(*gonsx.NSXClient)
	getAllAPI := securitygroup.NewGetAll(d.Get("scopeid").(string))

	err := nsxclient.Do(getAllAPI)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	for _, secGroup := range getAllAPI.GetResponse().SecurityGroups {
		if secGroup.Name == name {
			d.SetId(secGroup.ObjectID)
			return nil
		}
	}
	return fmt.Errorf("Security group %s not found", name)
}
