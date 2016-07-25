@virtualwire
Feature: VirtualWires

  Scenario: Terraform Plan with non existing Virtual Wire
    Given All the Virtual Wires with name "integration-test-01" do not exist in Scope "vdnscope-19"
    And I create a new Terraform manifest
    And I append a nsx_logical_switch resource with name "integration-test-01", description "ITest Virtual Wire", tenant "it-tenant" and scope "vdnscope-19"
    When I run terraform "plan"
    Then The command error code should be 0
    And The Terraform plan output should have 1 to add, 0 to change and 0 to destroy
    And Virtual Wire with name "integration-test-01" should not exist in Scope "vdnscope-19"

  Scenario: Terraform Apply with non existing Virtual Wire
    Given All the Virtual Wires with name "integration-test-01" do not exist in Scope "vdnscope-19"
    And I create a new Terraform manifest
    And I append a nsx_logical_switch resource with name "integration-test-01", description "IT Virtual Wire", tenant "it-tenant" and scope "vdnscope-19"
    When I run terraform "apply"
    Then The command error code should be 0
    And The Terraform apply output should have 1 added, 0 changed and 0 destroyed
    And Virtual Wire with name "integration-test-01" should exist in Scope "vdnscope-19"


#  Scenario: Terraform Plan with an already existing Virtual Wire
#    Given All the Virtual Wires with name "integration-test-01" do not exist in Scope "vdnscope-19"
#    And Virtual Wire with name "integration-test-01", description: "IT Virtual Wire", tenant "it-tenant" and scope "vdnscope-19" exists
#    And I create a new Terraform manifest
#    And I append a nsx_logical_switch resource with name "integration-test-01", description "IT Virtual Wire", tenant "it-tenant" and scope "vdnscope-19"
#    When I run terraform "plan"
#    Then The command error code should be 0
#    And The Terraform plan output should have 0 to add, 0 to change and 0 to destroy
#    And Virtual Wire with name "integration-test-01" should exist in Scope "vdnscope-19"
#
#  Scenario: Terraform Apply with an already existing Virtual Wire
#    Given All the Virtual Wires with name "integration-test-01" do not exist in Scope "vdnscope-19"
#    And Virtual Wire with name "integration-test-01", description: "IT Virtual Wire", tenant "it-tenant" and scope "vdnscope-19" exists
#    And I create a new Terraform manifest
#    And I append a nsx_logical_switch resource with name "integration-test-01", description "IT Virtual Wire", tenant "it-tenant" and scope "vdnscope-19"
#    When I run terraform "apply"
#    Then The command error code should be 0
#    And The Terraform apply output should have 0 added, 0 changed and 0 destroyed
#    And Virtual Wire with name "integration-test-01" should exist in Scope "vdnscope-19"

# TODO:
#  Scenario: Terraform create non existing Virtual Wire
#    Given There are 0 virtual wires attached to the DLR "test-dlr"
#    And I have a Terraform file with 1 virtual wires
#    When I run the command "terraform apply"
#    Then The Terraform output should have 1 changes
#    And There are 1 virtual wires attached to the DLR "test-dlr"