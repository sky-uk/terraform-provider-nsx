# Terraform-Provider-NSX

A Terraform provider for VMware NSX.  The NSX provider is used to interact
with resources supported by VMware NSX.  The provider needs to be configured
with the proper credentials before it can be used.

## Getting Started

```terra
# Configure the VMware NSX Provider
provider "nsx" {
    nsxusername = "${var.nsx_username}"
    nsxpassword = "${var.nsx_password}"
    nsxserver = "${var.nsx_server}"
}

# Create a logical switch
resource "nsx_logic_switch" "web" {
    ...
}
```
