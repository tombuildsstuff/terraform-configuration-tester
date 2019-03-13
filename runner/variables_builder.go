package runner

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func buildVariables(requiredVariables []terraformVariable, availableVariables []AvailableVariable) (*map[string]string, error) {
	output := make(map[string]string)

	for _, variable := range requiredVariables {
		log.Printf("[DEBUG] Finding a variable for %q", variable.Name)

		found := false
		value := ""
		for _, v := range availableVariables {
			if variable.Name != v.Name {
				continue
			}

			found = true

			if v.EnvKeyName != "" {
				log.Printf("[DEBUG] Looking up the value for %q from Env.Variable %q", v.Name, v.EnvKeyName)
				if val := os.Getenv(v.EnvKeyName); val != "" {
					value = val
					break
				}
			}

			if v.Generate {
				log.Printf("[DEBUG] Generating a value for %q", v.Name)
				value = generateRandomValue(5)
				log.Printf("[DEBUG] Generated %q", value)
				break
			}
		}

		if value != "" {
			output[variable.Name] = value
			continue
		}

		if variable.Default != nil {
			log.Printf("[DEBUG] No override found for Variable %q but it has a default value", variable.Name)
			continue
		}

		if found {
			return nil, fmt.Errorf("Variable %q was found but no value was provided, and `generate` was not set", variable.Name)
		}

		return nil, fmt.Errorf("Unable to find a value for %q", variable.Name)
	}

	return &output, nil
}

func generateRandomValue(length int) string {
	availableChars := "abcdefghijklmnopqrstuvwxyz012346789"
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = availableChars[rand.Intn(len(availableChars))]
	}
	return string(result)
}
