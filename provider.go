package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider is a basic structure that describes a provider: the configuration
// keys it takes, the resources it supports, a callback to configure, etc.
func Provider() terraform.ResourceProvider {
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"debug": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_ALLOW_UNVERIFIED_SSL", false),
			},
			"nsxusername": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXUSERNAME", nil),
			},
			"nsxpassword": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXPASSWORD", nil),
			},
			"nsxserver": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXSERVER", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"nsx_logical_switch":          resourceLogicalSwitch(),
			"nsx_edge_interface":          resourceEdgeInterface(),
			"nsx_dhcp_relay":              resourceDHCPRelay(),
			"nsx_service":                 resourceService(),
			"nsx_security_group":          resourceSecurityGroup(),
			"nsx_security_tag":            resourceSecurityTag(),
			"nsx_security_tag_attachment": resourceSecurityTagAttachment(),
			"nsx_security_policy":         resourceSecurityPolicy(),
			"nsx_security_policy_rule":    resourceSecurityPolicyRule(),
			"nsx_firewall_exclusion":      resourceFirewallExclusion(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"nsx_security_group": dataSourceSecurityGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	debug := d.Get("debug").(bool)
	insecure := d.Get("insecure").(bool)
	nsxusername := d.Get("nsxusername").(string)

	if nsxusername == "" {
		return nil, fmt.Errorf("nsxusername must be provided")
	}

	nsxpassword := d.Get("nsxpassword").(string)

	if nsxpassword == "" {
		return nil, fmt.Errorf("nsxpassword must be provided")
	}

	nsxserver := d.Get("nsxserver").(string)

	if nsxserver == "" {
		return nil, fmt.Errorf("nsxserver must be provided")
	}

	config := Config{
		Debug:       debug,
		Insecure:    insecure,
		NSXUserName: nsxusername,
		NSXPassword: nsxpassword,
		NSXServer:   nsxserver,
	}

	return config.Client()
}

// This is a global MutexKV for use within this plugin.
var nsxMutexKV = mutexkv.NewMutexKV()
