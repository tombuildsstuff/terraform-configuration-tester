package runner

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type terraformVariablesFile struct {
	Variables []terraformVariable `hcl:"variable,block"`

	// this ensures that other fields don't cause this to explode
	Ignore hcl.Body `hcl:",remain"`
}

type terraformVariable struct {
	Name        string  `hcl:"name,label"`
	Description *string `hcl:"description,attr"`
	Default     *string `hcl:"default,attr"`

	// this ensures that other fields don't cause this to explode
	Ignore hcl.Body `hcl:",remain"`
}

func parseVariables(fileName string) (*[]terraformVariable, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(fileName)
	if diags.HasErrors() {
		return nil, fmt.Errorf(diags.Error())
	}

	var variables terraformVariablesFile
	diags = gohcl.DecodeBody(file.Body, nil, &variables)
	if diags.HasErrors() {
		return nil, fmt.Errorf(diags.Error())
	}
	return &variables.Variables, nil
}
