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
| DHCP Relay              |   Y    |   Y   |    N    |   Y    |
| Edge Interface          |   Y    |   Y   |    N    |   Y    |
| Logical Switch          |   Y    |   Y   |    N    |   Y    |
| Security Group          |   Y    |   Y   |    Y    |   Y    |
| Security Policy         |   Y    |   Y   |    Y    |   Y    |
| Security Policy Rules   |   Y    |   Y   |    N    |   Y    |
| Security Tag            |   Y    |   Y   |    N    |   Y    |
| Security Tag Attachment |   Y    |   Y   |    N    |   Y    |
| Service                 |   Y    |   Y   |    N    |   Y    |

### DHCP Relay
***ADD EXAMPLE HERE***

### Edge Interface
***ADD EXAMPLE HERE***

### Logical Switch
***ADD EXAMPLE HERE***

### Security Group

> resource "nsx_security_group" "paas_test" {  
  &nbsp;&nbsp;name = "paas_test-oooo_test_security_group"  
  &nbsp;&nbsp;scopeid = "globalroot-0"  
  &nbsp;&nbsp;dynamicmembership {  
  &nbsp;&nbsp;&nbsp;&nbsp;membershipcriterialist {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;membershipoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriadetails {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriaoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;key = "VM.SECURITY_TAG"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;value = "dynamicset1_criteria1"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteria = "contains"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriadetails {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriaoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;key = "VM.SECURITY_TAG"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;value = "dynamicset1_criteria2"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteria = "contains"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;&nbsp;&nbsp;membershipcriterialist {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;          membershipoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriadetails {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriaoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;key = "VM.SECURITY_TAG"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;value = "dynamicset2_criteria1"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteria = "contains"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriadetails {  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteriaoperator = "OR"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;key = "VM.SECURITY_TAG"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;value = "dynamicset2_criteria2"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;criteria = "contains"  
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;&nbsp;&nbsp;}  
  &nbsp;&nbsp;}  
}  


### Security Policy
***ADD EXAMPLE HERE***

### Security Policy Rules
***ADD EXAMPLE HERE***

### Security Tags
***ADD EXAMPLE HERE***

### Security Tag Attachment
***ADD EXAMPLE HERE***

### Service
***ADD EXAMPLE HERE***


### Limitations

This is currently a proof of concept and only has a very limited number of
supported resources.  These resources also have a very limited number
of attributes.

We have only implemented the ability to Create, Read and Delete resources.
Currently there is no implementation of Update.

