# Terraform-Provider-NSX

A Terraform provider for VMware NSX.  The NSX provider is used to interact
with resources supported by VMware NSX.  The provider needs to be configured
with the proper credentials before it can be used.

## Wiki Pages
* [Home](https://github.com/sky-uk/terraform-provider-nsx/wiki)
* [Getting Started Guide](https://github.com/sky-uk/terraform-provider-nsx/wiki/Getting-Started-Guide)

## Features
| Feature                 | Create | Read  | Update  | Delete |
|-------------------------|--------|-------|---------|--------|
| DHCP Relay              |   Y    |   Y   |    Y    |   Y    |
| Edge Interface          |   Y    |   Y   |    N    |   Y    |
| Logical Switch          |   Y    |   Y   |    N    |   Y    |
| Security Group          |   Y    |   Y   |    Y    |   Y    |
| Security Policy         |   Y    |   Y   |    Y    |   Y    |
| Security Policy Rules   |   Y    |   Y   |    Y    |   Y    |
| Security Tag            |   Y    |   Y   |    Y    |   Y    |
| Security Tag Attachment |   Y    |   Y   |    Y    |   Y    |
| Service                 |   Y    |   Y   |    Y    |   Y    |


### Limitations

* Security-tag resource requires vsphere-provider with moid parameter implemented. ([branch](https://github.com/sky-uk/terraform/tree/OREP-176) not yet pushed to upstream). Docker image link with already built vsphere-provider available in getting started link above.
* At the moment only a very limited number of vSphere NSX resources have been implemented.  These resources also have the basic attributes implemented, look at wiki link above to find more details about each of these resources.



Examples for resources 
======================


DHCP Relay
-----------

This is an example terraform template for the DHCP Relay resource
 ```
 resource "nsx_dhcp_relay" "name_for_resource" {
   ipsets = ["ipset-3"]
   fqdn = ["testdomain.paas.bskyb.com","testdomain2.paas.bskyb.com"]
   edgeid       = "edge-5"
   dhcpserverip = ["10.152.160.11","10.152.160.12"]
   agent {
        vnicindex="9"
        giaddress="10.72.232.200"
   }
 }
 ```
 
 There is a limitation around ipsets and dhcpserverip, the api only supports 16 ip addresses, and this limit applies to both dhcpserver ip addresss and ipsets, 
 also the fqdn can only have 2 domain names in the list. 
 