provider "nsx" {
    nsxusername = "username"
    nsxpassword = "password"
    nsxserver = "apnsx020"
}

resource "nsx_logical_switch" "foo" {
}
