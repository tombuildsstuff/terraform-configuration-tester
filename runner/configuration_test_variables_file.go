package runner

import (
	"fmt"
	"log"
	"os"
)

func (test *TerraformConfigurationTest) writeTerraformTestVariablesFile() error {
	output := ""

	for k, v := range *test.variables {
		output += fmt.Sprintf("%s=%q\n", k, v)
	}

	out, err := os.Create(test.testVariablesFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("[DEBUG] Writing Test Variables file to %q", test.testFileName)
	bytes := []byte(output)
	out.Write(bytes)
	return out.Sync()
}

func (test *TerraformConfigurationTest) deleteTerraformTestVariablesFile() {
	os.Remove(test.testVariablesFileName)
}
