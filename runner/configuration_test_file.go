package runner

import (
	"fmt"
	"log"
	"os"
)

func (test *TerraformConfigurationTest) writeTerraformTestFile() error {
	config := fmt.Sprintf(`
terraform {
  required_version = "=%s"
}

provider "%s" {
  version = "=%s"
}
`, test.terraformVersion, test.providerName, test.providerVersion)

	out, err := os.Create(test.testFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("[DEBUG] Writing Test file to %q", test.testFileName)
	bytes := []byte(config)
	out.Write(bytes)
	return out.Sync()
}

func (test *TerraformConfigurationTest) deleteTerraformTestFile() {
	os.Remove(test.testFileName)
}
