package runner

import (
	"fmt"
	"log"
)

type TerraformConfigurationTest struct {
	availableVars    []AvailableVariable
	directory        string
	providerName     string
	providerVersion  string
	terraformVersion string

	testFileName          string
	testVariablesFileName string
	variables             *map[string]string
	variablesFileName     string
	terraformBinaryPath   string
}

func (test TerraformConfigurationTest) run() error {
	log.Printf("[DEBUG] Determining which Variables are required..")
	if err := test.populateVariables(); err != nil {
		return fmt.Errorf("Error populating variables: %s", err)
	}

	log.Printf("[DEBUG] Writing the Test Config..")
	if err := test.writeTerraformTestFile(); err != nil {
		return fmt.Errorf("Error writing test file: %s", err)
	}
	defer test.deleteTerraformTestFile()

	log.Printf("[DEBUG] Writing the Test Variables..")
	if err := test.writeTerraformTestVariablesFile(); err != nil {
		return fmt.Errorf("Error writing test file: %s", err)
	}
	defer test.deleteTerraformTestVariablesFile()

	log.Printf("[DEBUG] Downloading Terraform..")
	if err := test.downloadTerraform(); err != nil {
		return fmt.Errorf("Error downloading Terraform: %s", err)
	}

	log.Printf("[DEBUG] Running Terraform..")
	if err := test.runTerraform(); err != nil {
		return fmt.Errorf("Error running Terraform: %s", err)
	}

	log.Printf("[DEBUG] Test Completed")
	return nil
}
