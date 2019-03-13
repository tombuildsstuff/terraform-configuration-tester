package runner

import (
	"fmt"
)

type TestRunInput struct {
	AvailableVariables []AvailableVariable
	ProviderName       string
	ProviderVersion    string
	TerraformVersion   string
}

func (input TestRunInput) Run(directory string) error {
	test := TerraformConfigurationTest{
		availableVars:         input.AvailableVariables,
		directory:             directory,
		providerName:          input.ProviderName,
		providerVersion:       input.ProviderVersion,
		terraformVersion:      input.TerraformVersion,
		terraformBinaryPath:   fmt.Sprintf("%s/terraform", directory),
		testFileName:          fmt.Sprintf("%s/test.tf", directory),
		testVariablesFileName: fmt.Sprintf("%s/test.tfvars", directory),
		variablesFileName:     fmt.Sprintf("%s/variables.tf", directory),
	}
	return test.run()
}
