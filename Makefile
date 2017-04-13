terraform-provider-nsx: config.go main.go provider.go resource_logical_switch.go resource_edge_interface.go resource_dhcp_relay.go resource_service.go resource_security_group.go resource_security_tag.go resource_security_tag_attachment.go resource_security_policy_rules.go
	go build -o terraform-provider-nsx
	strip terraform-provider-nsx
install:
	cp terraform-provider-nsx /usr/local/terraform/
clean:
	#rm -f terraform.tfstate  terraform.tfstate.backup terraform-provider-nsx crash.log terraform.log
	rm -f terraform-provider-nsx terraform.log crash.log
apply: terraform-provider-nsx
	terraform apply 2>&1 | tee terraform.log
apply-debug: terraform-provider-nsx
	TF_LOG=1 terraform apply 2>&1 | tee terraform.log
destroy: terraform-provider-nsx
	terraform destroy -force 2>&1 | tee terraform.log
destroy-debug: terraform-provider-nsx
	TF_LOG=1 terraform destroy -force 2>&1 | tee terraform.log
