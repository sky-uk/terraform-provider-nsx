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
    name = "paas_test-oooo_test_security_group"
    scopeid = "globalroot-0"
    dynamicmembership {
        membershipcriterialist {
            membershipoperator = "OR"
            criteriadetails {
                criteriaoperator = "OR"
                key = "VM.SECURITY_TAG"
                value = "dynamicset1_criteria1"
                criteria = "contains"
            }
            criteriadetails {
                criteriaoperator = "OR"
                key = "VM.SECURITY_TAG"
                value = "dynamicset1_criteria2"
                criteria = "contains"
            }
        }
        membershipcriterialist {
            membershipoperator = "OR"
            criteriadetails {
                criteriaoperator = "OR"
                key = "VM.SECURITY_TAG"
                value = "dynamicset2_criteria1"
                criteria = "contains"
            }
            criteriadetails {
                criteriaoperator = "OR"
                key = "VM.SECURITY_TAG"
                value = "dynamicset2_criteria2"
                criteria = "contains"
            }
        }
    }
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

