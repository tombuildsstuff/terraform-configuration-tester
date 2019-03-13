## Terraform Provider Example Tester

This library provides several helpers to allow for testing Examples of a Terraform Provider.

For each Directory found containing examples (requirements below) - this Test Helper does the following:

* Parses the `variables.tf` file within the directory and then either Generates or passes through a value from an Env. Variable
* Downloads Terraform Core from releases.hashicorp.com
* Writes out a `test.tf` which contains the following:

```hcl
terraform {
  required_version = "=0.11.13"
}

provider "azurerm" {
  version = "=1.23.0"
}
```

* Writes out a `test.tfvars` containing the values for each variable listed above
* Runs: `terraform init`, `terraform validate`, `terraform apply` and then `terraform destroy`

The intention is that this is used in an Acceptance Test to validate each example - see below for an example of this.

---

In order to be useful, this Tester makes several assumptions about the directory structure of each example:

* `main.tf` - which contains (at least some) Terraform Configuration.
* `variables.tf` - contains all of the variables required for the example.

In addition - tests can be skipped by placing a `.skip-test` file in the directory.

### Using this library

Given an examples folder which contains one or more Terraform Configurations - an Acceptance Test similar to below can be used:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/tombuildsstuff/terraform-configuration-tester/runner"
	"github.com/tombuildsstuff/terraform-configuration-tester/locator"
)

func TestRunExamples(t *testing.T) {
	if os.Getenv("TF_EXAMPLE_TEST") == "" {
		log.Printf("`TF_EXAMPLE_TEST` is not set - skipping")
		t.Skip()
	}

	examplesDirectory := "./examples"
	directories := locator.DiscoverExamples(examplesDirectory)

	input := runner.TestRunInput{
		ProviderVersion:  "1.23.0",
		ProviderName:     "azurerm",
		TerraformVersion: "0.11.13",
		AvailableVariables: []runner.AvailableVariable{
			{
				Name:     "prefix",
				Generate: true,
			},
			{
				Name:       "location",
				EnvKeyName: "ARM_LOCATION",
			},
			{
				Name:       "alt_location",
				EnvKeyName: "ARM_LOCATION_ALT",
			},
			{
				Name:       "kubernetes_client_id",
				EnvKeyName: "ARM_CLIENT_ID",
			},
			{
				Name:       "kubernetes_client_secret",
				EnvKeyName: "ARM_CLIENT_SECRET",
			},
		},
	}

	for _, directoryPath := range directories {
		shortDirName := strings.Replace(directoryPath, examplesDirectory, "", -1)
		testName := fmt.Sprintf("examples%s", shortDirName)
		t.Run(testName, func(t *testing.T) {
			if err := input.Run(directoryPath); err != nil {
				t.Fatalf("Error running %q: %s", shortDirName, err)
			}
		})
	}
}
```

### Requirements

- Go
- `unzip` available on your PATH
