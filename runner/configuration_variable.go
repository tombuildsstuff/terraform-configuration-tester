package runner

import "fmt"

func (test *TerraformConfigurationTest) populateVariables() error {
	// determine which variables are required from the `variables.tf` in the directory
	requiredVariables, err := parseVariables(test.variablesFileName)
	if err != nil {
		return fmt.Errorf("Error parsing variables from %q: %s", test.variablesFileName, err)
	}

	// then create some values for them
	variables, err := buildVariables(*requiredVariables, test.availableVars)
	if err != nil {
		return fmt.Errorf("Error building variables: %s", err)
	}

	test.variables = variables
	return nil
}
