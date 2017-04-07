package virtualwire

import (
	"fmt"
	. "github.com/lsegal/gucumber"
	"github.com/sky-uk/terraform-provider-nsx/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

// Client : NSXClient object
var Client = internal.GetNSXClient()

// StdOut : Command stdout
var StdOut []byte

// StdErr : Command stderr
var StdErr []byte

// ExitError : Command exitError
var ExitError *exec.ExitError

// TempTFDir : Terraform file temporary directory
var TempTFDir string

// TempTDFile : Temporary Terraform file
var TempTDFile string

// TerraformPlanChangesRegExp : Regular expression for terraform plan output
var TerraformPlanChangesRegExp = regexp.MustCompile(`\nPlan: (\d) to add, (\d) to change, (\d) to destroy\.\n`)

// TerraformApplyChangesRegExp : Regular expression for terraform apply output
var TerraformApplyChangesRegExp = regexp.MustCompile(`\nApply complete! Resources: (\d) added, (\d) changed, (\d) destroyed\.\n`)

func init() {

	Before("", func() {
		os.Remove("terraform.tfstate")
	})

	After("", func() {
		if TempTFDir != "" {
			os.Remove(TempTFDir)
		}
		os.Remove("terraform.tfstate")
	})

	Given(`^Virtual Wire with name "(.+?)" should not exist in Scope "(.+?)"$`, func(s1 string, s2 string) {
		assert.Empty(T, getVirtualWire(s1, s2))
	})

	And(`^I create a new Terraform manifest$`, func() {
		var err error
		TempTFDir, err = ioutil.TempDir(os.TempDir(), "gonsx-integration")
		internal.CheckError(err)
		TempTDFile = TempTFDir + "/test.tf"
	})

	And(`^I append a nsx_logical_switch resource with name "(.+?)", description "(.+?)", tenant "(.+?)" and scope "(.+?)"$`, func(s1 string, s2 string, s3 string, s4 string) {
		str := `provider nsx { insecure=true }
		        resource "nsx_logical_switch" "virtual_wire" {
				name = "` + s1 + `"
		    		desc = "` + s2 + `"
				tenantid = "` + s3 + `"
				scopeid = "` + s4 + `"
			}`
		var err error
		err = ioutil.WriteFile(TempTDFile, []byte(str), 0777)
		internal.CheckError(err)
	})

	When(`^I run terraform "(.+?)"`, func(s1 string) {
		StdOut, StdErr, ExitError = internal.ExecuteCommand("terraform", s1, "-no-color", TempTFDir)
	})

	Then(`^The command error code should be (\d+)$`, func(i1 int) {
		if i1 == 0 {
			assert.Nil(T, ExitError)
		} else {
			assert.Equal(T, fmt.Sprintf("exit status %d", i1), ExitError.Error())
		}
	})

	And(`^The Terraform plan output should have (\d+) to add, (\d+) to change and (\d+) to destroy`, func(i1 int, i2 int, i3 int) {
		planOutput := TerraformPlanChangesRegExp.FindStringSubmatch(string(StdOut))
		if len(planOutput) != 4 {
			printStdOutAndErr()
			log.Panic("Couldn't parse terraform plan output!")
		}
		toAdd, _ := strconv.Atoi(planOutput[1])
		toChange, _ := strconv.Atoi(planOutput[2])
		toDestroy, _ := strconv.Atoi(planOutput[3])
		assert.Exactly(T, i1, toAdd)
		assert.Exactly(T, i2, toChange)
		assert.Exactly(T, i3, toDestroy)
	})

	And(`^The Terraform apply output should have (\d+) added, (\d+) changed and (\d+) destroyed$`, func(i1 int, i2 int, i3 int) {
		applyOuput := TerraformApplyChangesRegExp.FindStringSubmatch(string(StdOut))
		if len(applyOuput) != 4 {
			printStdOutAndErr()
			log.Panic("Couldn't parse terraform apply output!")
		}
		added, _ := strconv.Atoi(applyOuput[1])
		changed, _ := strconv.Atoi(applyOuput[2])
		destroyed, _ := strconv.Atoi(applyOuput[3])
		assert.Equal(T, i1, added)
		assert.Equal(T, i2, changed)
		assert.Equal(T, i3, destroyed)
	})

	And(`^Virtual Wire with name "(.+?)" should exist in Scope "(.+?)"$`, func(s1 string, s2 string) {
		assert.NotEmpty(T, getVirtualWire(s1, s2).ObjectID)
	})

	Given(`^Virtual Wire with name "(.+?)", description: "(.+?)", tenant "(.+?)" and scope "(.+?)" exists$`, func(s1 string, s2 string, s3 string, s4 string) {
		createVirtualWire(s1, s2, s3, s4)
	})

	Given(`^All the Virtual Wires with name "(.+?)" do not exist in Scope "(.+?)"$`, func(s1 string, s2 string) {
		deleteAllVirtualWiresWithNameInScope(s1, s2)
	})

}

func printStdOutAndErr() {
	log.Println(string(StdOut))
	log.Println(string(StdErr))
}
