# Terraform-Provider-NSX

##2016-07-18 - Release 0.1.1

### Summary

Minor improvement and bug fix.

### Bug Fixes
Fix for 'DHCP relays being updated in parallel causing unexpected results'. All DHCP relays are pushed on delete or create due to the design of the API. Fix permits one change at a time to avoid unexpected results.

### Improvements
Integrated with godep.

##2016-07-14 - Release 0.1.0

### Summary

The initial (proof of concept) release of the plugin.  Please refer to the
*Limitations* section of the README.

The following resources can be created and deleted:

* `nsx_logical_switch`
* `nsx_edge_interface`
* `nsx_dhcp_relay`
