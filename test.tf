# Ensure that you have the environment variables NSXUSERNAME, NSXPASSWORD
# and NSXSERVER.
provider "nsx" {
        insecure = true 
}

resource "nsx_logical_switch" "virtual_wire" {
        desc = "Terraform managed Logical Switch"
        name = "tf_test"
        tenantid = "tf_testid"
        scopeid = "vdnscope-19"
}
