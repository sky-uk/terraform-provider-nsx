# Ensure that you have the environment variables NSXUSERNAME, NSXPASSWORD
# and NSXSERVER.
provider "nsx" {
        insecure = true 
}

resource "nsx_logical_switch" "virtual_wire" {
        desc = "Terraform managed Logical Switch"
        name = "tf_logical_switch"
        tenantid = "tf_testid"
        scopeid = "vdnscope-19"
}

resource "nsx_edge_interface" "edge_interface" {
        edgeid = "edge-50"
        name = "tf_edge_interface"
        virtualwireid = "${nsx_logical_switch.virtual_wire.id}"
        gateway = "10.10.10.1"
        subnetmask = "255.255.255.0"
        interfacetype = "internal"
        mtu = "1500"
}
