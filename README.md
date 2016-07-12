# Terraform-Provider-NSX

A Terraform provider for VMware NSX.  The NSX provider is used to interact
with resources supported by VMware NSX.  The provider needs to be configured
with the proper credentials before it can be used.

## Contents

* [Installation](#installation)
* [Getting Started](#getting-started)
* [Authentication](#authentication)
  * [Static Credentials](#static-credentials)
  * [Environment Variables](#environment-variables)
  * [Argument Reference](#provider-argument-reference)
* [NSX_LOGICAL_SWITCH Resource](#nsx-logical-switch-resource)
  * [Example Usage](#nsx-logical-switch-resource-example-usage)

## Installation

These instructions were tested against version 1.6.2 of Go and Terraform
version 0.6.16.

Assuming that you have already got GOPATH setup
(see https://golang.org/doc/code.html for details). Do the following:

Install the Terraform library from HashiCorp and workaround the new
Terraform API version introduced after 0.7.0:

```bash
$ go get -u github.com/hashicorp/terraform
$ cd $GOPATH/src/github.com/hashicorp/terraform
$ git checkout v0.6.16
Note: checking out 'v0.6.16'.

You are in 'detached HEAD' state. You can look around, make experimental
changes and commit them, and you can discard any commits you make in this
state without impacting any branches by performing another checkout.

If you want to create a new branch to retain commits you create, you may
do so (now or later) by using -b with the checkout command again. Example:

  git checkout -b <new-branch-name>

HEAD is now at 6e586c8... v0.6.16
```

As this project is hosted on GitLab and this does not yet support `go get ...`
manually checkout this project to the correct location:

```bash
mkdir -p $GOPATH/src/git.devops.int.ovp.bskyb.com/paas
cd $GOPATH/src/git.devops.int.ovp.bskyb.com/paas
git clone git@git.devops.int.ovp.bskyb.com:paas/terraform-provider-nsx.git
```

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

### <a name='nsx-logical-switch-resource-example-usage'></a>Example Usage

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
