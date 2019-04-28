# Terraform-Provider-NSX

A Terraform provider for VMware NSX.  The NSX provider is used to interact
with resources supported by VMware NSX.  The provider needs to be configured
with the proper credentials before it can be used.

## Wiki Pages
* [Home](https://github.com/sky-uk/terraform-provider-nsx/wiki)
* [Authentication](https://github.com/sky-uk/terraform-provider-nsx/wiki/Authentication)
* [Getting Started Guide](https://github.com/sky-uk/terraform-provider-nsx/wiki/Getting-Started-Guide)
* [NSX DHCP Relay Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-DHCP-Relay-Resource)
* [NSX DHCP Relay Agent Resource](https://github.com/sgdigital-devops/terraform-provider-nsx/wiki/NSX-DHCP-Relay-Agent-Resource)
* [NSX Edge Interface Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Edge-Interface-Resource)
* [NSX Logical Switch Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Logical-Switch-Resource)
* [NSX Security Group Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Security-Group-Resource)
* [NSX Security Policy Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Security-Policy-Resource)
* [NSX Security Policy Rule resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Security-Policy-Resource#nsx_security_policy_rule-resource)
* [NSX Security Tag Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Security-Tag-Resource)
* [NSX Security Tag Attachment Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Security-Tag-Resource#nsx_security_tag_attachment-resource)
* [NSX Service Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Service-Resource)
* [NSX Firewall Exclusion Resource](https://github.com/sky-uk/terraform-provider-nsx/wiki/NSX-Firewall-Exclusion)
* [NSX Nat Rules Resource] (https://github.com/sgdigital-devops/terraform-provider-nsx/wiki/NSX-Nat-Rules-Resource)


## Features
| Feature                 | Create | Read | Update | Delete |
|:------------------------|:-------|:-----|:-------|:-------|
| DHCP Relay              | Y      | Y    | Y      | Y      |
| DHCP Relay Agent        | Y      | Y    | N      | Y      |
| Edge Interface          | Y      | Y    | N      | Y      |
| Logical Switch          | Y      | Y    | Y      | Y      |
| Security Group          | Y      | Y    | Y      | Y      |
| Security Policy         | Y      | Y    | Y      | Y      |
| Security Policy Rules   | Y      | Y    | Y      | Y      |
| Security Tag            | Y      | Y    | Y      | Y      |
| Security Tag Attachment | Y      | Y    | Y      | Y      |
| Service                 | Y      | Y    | Y      | Y      |
| Firewall Exclusion      | Y      | Y    | N      | Y      |
| Nat Rule                | Y      | Y    | Y      | Y      |


### Limitations

* There are two ways of providing Agents to the DHCP Relay Configuration. It has its own DHCP Relay Agent Resource and it can be provided inline inside of the DHCP Relay. Those two methods can not be mixed, and doing so will cause conflicts.

* Security-tag resource requires vsphere-provider with moid parameter implemented. ([branch](https://github.com/sky-uk/terraform/tree/OREP-176) not yet pushed to upstream). Docker image link with already built vsphere-provider available in getting started link above. - This issue was actually solved on terraform v0.9.6 - pull request here  (https://github.com/hashicorp/terraform/pull/14793)


* At the moment only a very limited number of vSphere NSX resources have been implemented.  These resources also have the basic attributes implemented, look at wiki link above to find more details about each of these resources.



### Resources to consider

 - Transport Zones
 - Distributed Switch
 - Distributed Firewall (L2 / L3 rules)
 - Edge Device Nat config and rules
 - Edge Device Routing config


### Acceptence Testing

To allow to run Acceptence tests a number of base resources have to be made available. The resources are selected by providing ENV variables. If these variables are not provided the tests that rely on them are simply skipped.

Resources needed:
* DLR with two interfaces (vnic 10 and vnic 11)
* ESG with one interface (vnic 0 (uplink))
* Virtualwire

The following reosurces are required (with example values):

```
export NSX_TESTING_DLR_ID=edge-16
export NSX_TESTING_ESG_ID=edge-17
export NSX_TESTING_VIRTUALWIRE_ID=virtualwire-48
```