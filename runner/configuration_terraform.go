package runner

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func (test TerraformConfigurationTest) downloadTerraform() error {
	url := fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s_%s.zip", test.terraformVersion, test.terraformVersion, runtime.GOOS, runtime.GOARCH)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	zipFileName := fmt.Sprintf("%s.zip", test.terraformBinaryPath)
	out, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	// then extract the zip file
	c := exec.Command("unzip", "-o", zipFileName)
	c.Env = os.Environ()
	c.Dir = test.directory
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	if err := c.Start(); err != nil {
		return fmt.Errorf("Error unzipping Terraform: %s", err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("Error unzipping Terraform: %s", err)
	}

	return nil
}

func (test TerraformConfigurationTest) runTerraform() error {
	log.Printf("[DEBUG] Running `terraform init`")
	if err := test.terraformInit(); err != nil {
		return fmt.Errorf("Error running `terraform init`: %s", err)
	}

	log.Printf("[DEBUG] Running `terraform validate`")
	if err := test.terraformValidate(); err != nil {
		return fmt.Errorf("Error running `terraform validate`: %s", err)
	}

	log.Printf("[DEBUG] Running `terraform apply`")
	if err := test.terraformApply(); err != nil {
		return fmt.Errorf("Error running `terraform apply`: %s", err)
	}

	log.Printf("[DEBUG] Running `terraform destroy`")
	if err := test.terraformDestroy(); err != nil {
		return fmt.Errorf("Error running `terraform destroy`: %s", err)
	}

	return nil
}

func (test TerraformConfigurationTest) terraformInit() error {
	c := exec.Command(test.terraformBinaryPath, "init")
	c.Env = os.Environ()
	c.Dir = test.directory
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	if err := c.Start(); err != nil {
		return fmt.Errorf("Error launching Init: %s", err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("Error running Init: %s", err)
	}

	return nil
}

func (test TerraformConfigurationTest) terraformApply() error {
	return test.runTerraformCommand("apply", "--auto-approve")
}

func (test TerraformConfigurationTest) terraformDestroy() error {
	return test.runTerraformCommand("destroy", "--auto-approve")
}

func (test TerraformConfigurationTest) terraformValidate() error {
	return test.runTerraformCommand("validate", "")
}

func (test TerraformConfigurationTest) runTerraformCommand(command string, extra string) error {
	c := exec.Command(test.terraformBinaryPath, command, fmt.Sprintf("-var-file=%s", test.testVariablesFileName), extra)
	c.Env = os.Environ()
	c.Dir = test.directory
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout

	if err := c.Start(); err != nil {
		return fmt.Errorf("Error launching %q: %s", command, err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("Error running %q: %s", command, err)
	}

	return nil
}
