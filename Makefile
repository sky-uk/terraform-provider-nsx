terraform-provider-nsx: config.go main.go provider.go resource_logical_switch.go resource_edge_interface.go
	go build -o terraform-provider-nsx
	strip terraform-provider-nsx
clean:
	rm -f terraform.tfstate  terraform.tfstate.backup terraform-provider-nsx crash.log terraform.log
test: terraform-provider-nsx
	TF_LOG=1 terraform apply 2>&1 | tee terraform.log
