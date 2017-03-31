# Terraform-Provider-NSX

A Terraform provider for VMware NSX.  The NSX provider is used to interact
with resources supported by VMware NSX.  The provider needs to be configured
with the proper credentials before it can be used.

## Contents

* [Installation](#installation)
* [Features] (#features)
* [Getting Started](#getting-started)
* [Authentication](#authentication)
* [NSX_LOGICAL_SWITCH Resource](#nsx_logical_switch-resource)
* [NSX_EDGE_INTERFACE Resource](#nsx_edge_interface-resource)
* [NSX_DHCP_RELAY Resource](#nsx_dhcp_relay-resource)
* [NSX_SERVICE Resource](#nsx_service-resource)
* [NSX_SECURITY_GROUP Resource](#nsx_security_group-resource)
* [NSX_SECURITY_TAG Resource](#nsx_security_tag-resource)
* [NSX_SECURITY_TAG_ATTACHMENT Resource](#nsx_security_tag_attachment-resource)
* [Limitations](#limitations)

## Installation

These instructions were tested against version 1.6.2 of Go and Terraform
version 0.6.16.

Assuming that you have already got GOPATH setup
(see https://golang.org/doc/code.html for details). Do the following:

```bash
go get github.com/sky-uk/terraform-provider-nsx
```

This will also build the binary and add the `terraform-provider-nsx`
plugin into the `$GOPATH/bin`.

## Features
| Feature                 | Create | Read  | Update  | Delete |
|-------------------------|--------|-------|---------|--------|
| DHCP Relay              |   Y    |   Y   |    N    |   Y    |
| Edge Interface          |   Y    |   Y   |    N    |   Y    |
| Logical Switch          |   Y    |   Y   |    N    |   Y    |
| Security Group          |   Y    |   Y   |    N    |   Y    |
| Security Policy         |   Y    |   Y   |    N    |   Y    |
| Security Policy Rules   |   Y    |   Y   |    N    |   Y    |
| Security Tag            |   Y    |   Y   |    N    |   Y    |
| Security Tag Attachment |   Y    |   Y   |    N    |   Y    |
| Service                 |   Y    |   Y   |    N    |   Y    |


## Getting Started

```terra
# Configure the VMware NSX Provider
provider "nsx" {
    nsxusername = "${var.nsx_username}"
    nsxpassword = "${var.nsx_password}"
    nsxserver = "${var.nsx_server}"
}

# Create a logical switch
resource "nsx_logic_switch" "OVP_1" {
    ...
}
```

## Authentication

The NSX provider offers flexible means of providing credentials for
authentication.  The following methods are supported, in this order and
explained below:

* Static credentials
* Environment variables

### Static credentials

Static credentials can be provided by adding `nsxusername`, `nsxpassword`
and `nsxserver` in-line in the nsx provider block:

Usage:

```terra
provider "nsx" {
    nsxusername = "username"
    nsxpassword = "password"
    nsxserver = "apnsx020"
}
```

### Environment variables

You can provide your credentials via NSXUSERNAME, NSXPASSWORD and NSXSERVER
environment variables, representing your user name, password and NSX server
respectively.

```terra
provider "nsx" {}
```

Usage:

```bash
$ export NSXUSERNAME='username'
$ export NSXPASSWORD='password'
$ export NSXSERVER='apnsx020'
$ terraform plan
```

### <a name="provider-argument-reference"></a> Argument Reference

The following arguments are supported in the `provider` block:

* `debug` - (Optional) Show debug output of API operations to NSX
If omitted, default value is `false`.
* `insecure` - (Optional) Explicitly allow the provider to perform "insecure"
SSL requests. If omitted, default value is `false`.
* `nsxpassword` - (Optional) This is the password for connecting to the NSX
server.  It must be provided, but it can also be sourced from the `NSXPASSWORD`
environment variable.
* `nsxusername` - (Optional) This is the user name for connecting to the NSX
server.  It must be provided, but it can also be sourced from the `NSXUSERNAME`
environment variable.
* `nsxserver` - (Optional) This is the NSX server to connect to.  It must be
provided, but it can also be sourced from the `NSXSERVER` environment variable.

## NSX_LOGICAL_SWITCH Resource

The LOGICAL_SWITCH resource allows the creation and management of a logical
switch (sometimes virtual wire, port group or universal switch).

### Example Usage

```terra
resource "nsx_logical_switch" "virtual_wire" {
    desc = "Terraform managed Logical Switch"
    name = "tf_test"
    tenantid = "tf_testid"
    scopeid = "vdnscope-19"
}
```

### Argument Reference

The following arguments are supported:
 
* `desc` - (Required) A longer description for the logical switch.
* `name` - (Required) The name of the logical switch to be created.
* `tenantid` - (Required) The ID of the tenant that logical switch is to be
associated with.
* `scopeid` - (Required) The ID of transport zone that the logical switch is
to be associated with.


## NSX_EDGE_INTERFACE Resource

The EDGE_INTERFACE resource allows the creation and management of an edge
interface on the Distributed Logical Router (DLR).

### Example Usage

```terra
resource "nsx_edge_interface" "edge_interface" {
    edgeid = "edge-50"
    name = "app_virtualwire_one"
    virtualwireid = "virtualwire-271"
    gateway = "10.10.10.1"
    subnetmask = "255.255.255.0"
    interfacetype = "internal"
    mtu = "1500"
}
```

### Argument Reference

The following arguments are supported:
 
* `edgeid` - (Required) The NSX Edge ID for the Distributed Logical Router (DLR) we wish to use.
* `name` - (Required) The name of the edge interface we want to create on the DLR.
* `virtualwireid` - (Required) The ID of the virtual wire/logical switch (see [NSX Logical Switch](#nsx-logical-switch-resource).
* `gateway` - (Required) Gateway for network.
* `subnetmask` - (Required) Subnet mask for network.
* `interfacetype` - (Required) The interface type.
* `mtu` - (Required) Max transfer unit for the network.

## NSX_DHCP_RELAY Resource

The DHCP_RELAY resource allows the creation and management of a DHCP
relay for an edge interface on the Distributed Logical Router (DLR).

### Example Usage

```terra
resource "nsx_dhcp_relay" "dhcp_relay" {
    name = "tf_dhcp_relay"
    edgeid = "edge-50"
    vnicindex = "18"
    giaddress = "10.152.163.1"
    dhcpserverip = "10.152.160.10"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) A user supplied name for the DHCP relay.
* `edgeid` - (Required) The NSX Edge ID for the Distributed Logical Router (DLR) we wish to use.
* `vnicindex` - (Required) The VNIC Index.
* `giaddress` - (Required) The GIAddress is the IP address of the gateway on the edge interface (nsx_edge_interface gateway value above).
* `dhcpserverip` (Required) The IP address of the DHCP server to have requests
relayed to.

 

## NSX_SERVICE Resource

The SERVICE resource allows the creation of Services/Applications for use by
service groups and service policies.

### Example Usage

```terra
resource "nsx_service" "http" {
    name = "tf_service_http_80"
    scopeid = "globalroot-0"
    desc = "TCP port 80 - http"
    proto = "TCP"
    ports = "80"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name you want to call this service by.
* `scopeid` - (Required) The scopeid.
* `desc` - (Required) Description of the service.
* `proto` - (Required) The chosen protocol. E.g. TCP.
* `ports` - (Required) The ports assigned to this service. 



## NSX_SECURITY_GROUP Resource

The SECURITY_GROUP resource allows the creation of Security Groups for use by
service policies. Currently this will only have the one security tag to compare 
on within the group.

### Example Usage

```terra
resource "nsx_security_group" "web" {
    name = "tf_security_group_web"
    scopeid = "globalroot-0"
    setoperator = "OR"
    criteriaoperator = "OR"
    criteriakey = "VM.SECURITY_TAG"
    criteriavalue = "tag_name"
    criteria = "contains"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name you want to call this security group by.
* `scopeid` - (Required) The scopeid.
* `setoperator` - (Required) Set is used to combine the result of the dynamic set(s) evaluated previously with the result of this dynamic set. The possible values for this field are "AND" and "OR".
* `criteriaoperator` - (Required) The operator for the criteria. 
* `criteriakey` - (Required) The key in which the criteria should use to match.
* `criteriavalue` - (Required) The value in which the criteria should match.
* `criteria` - (Required) How the criteria should match.

## NSX_SECURITY_TAG Resource

The SECURITY_TAG resource allows the creation of Security Tags for use by
security groups and virtual machines. 

### Example Usage

```terra
resource "nsx_security_tag "web" {
    name = "tf_security_tag"
    desc = "TF Security Tag for web hosts"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name you want to call this security tag by.
* `desc` - (Required) A friendly description of the security tag.


## NSX_SECURITY_TAG_ATTACHMENT Resource

The SECURITY_TAG resource allows the attachment of Security Tags to 
virtual machines. 

### Example Usage

```terra
resource "nsx_security_tag_attachment "web" {
    tagid = "securitytag-1"
    moid = "vm-1"
}
```

### Argument Reference

The following arguments are supported:

* `tagid` - (Required) ID of security tag.
* `moid` - (Required) ID of vm.

## NSX_SECURITY_GROUPS Resource

The SECURITY_GROUPS resource allows the creation of Security Groups. 

### Example Usage

```terra
resource "nsx_security_group" "web" {
       name = "tf_web_security_group"
       scopeid = "globalroot-0"
       setoperator = "OR"
       criteriaoperator = "OR"
       criteriakey = "VM.SECURITY_TAG"
       criteriavalue = "tf_web_security_tag"
       criteria = "contains"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) Name of security group.
* `scopeid` - (Required) ID of scope.
* `setoperator` - (Required) "AND" or "OR" operator to match the 
criterias with the list of criterias.
* `criteriaoperator` - (Required) "AND" or "OR" operator to match within
 the criterias.
* `criteriakey` - (Required) The key to match the criterias on e.g. 
security tags.
* `criteriavalue` - (Required) Value to match criteria.
* `criteria` - (Required) What the criteria is.


## NSX_SECURITY_POLICY Resource

The SECURITY_POLICY resource allows the creation of Security Policies 
for use by security groups. 

### Example Usage

```terra
resource "nsx_security_policy" "web" {
       name = "tf_web_security_policy"
       description = "security policy for web role"
       precedence  = "55002"
       securitygroups = ["${nsx_security_group.web.id}"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) Name of security policy.
* `description` - (Required) Description of policy.
* `precedence` - (Required) Importance of the rule.
* `securitygroups` - (Required) List of security groups to attach policy
 to.


## NSX_SECURITY_POLICY_RULE resource

The SECURITY_POLICY_RULE creates rules on security policies.


### Example Usage

```terra
resource "nsx_security_policy_rule" "web" {
      name = "tf_web_security_policy_rule"
      securitypolicyname = "${nsx_security_policy.security_policy_web.name}"
      action = "allow"
      direction = "outbound"
      securitygroupids = ["${nsx_security_group.web.id}"]
      serviceids = ["${nsx_service.web.id}"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) Name of security policy rule.
* `securitypolicyname` - (Required) Name of policy to attach to.
* `action` - (Required) "ALLOW" or "BLOCK".
* `direction` - (Required) "OUTBOUND" or "INBOUND".
* `securitygroupids` - (Required) List of groups to add rule to.
* `serviceids` - (Required) List of services to apply.

 
 

### Limitations

This is currently a proof of concept and only has a very limited number of
supported resources.  These resources also have a very limited number
of attributes.

We have only implemented the ability to Create, Read and Delete resources.
Currently there is no implementation of Update.

