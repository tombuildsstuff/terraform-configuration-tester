package locator

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func DiscoverExamples(directory string) []string {
	directories := make(map[string]struct{})

	globs := []string{
		fmt.Sprintf("%s/*/main.tf", directory),
		fmt.Sprintf("%s/**/**/main.tf", directory),
		fmt.Sprintf("%s/**/**/**/main.tf", directory),
		fmt.Sprintf("%s/**/**/**/**/main.tf", directory),
	}
	for _, glob := range globs {
		files, err := filepath.Glob(glob)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range files {
			directoryName := filepath.Dir(v)
			directories[directoryName] = struct{}{}
		}
	}

	output := make([]string, 0)
	for v := range directories {
		// double-check if the path should be skipped via a `.skip-test` file
		directoryPath := path.Dir(v)
		_, err := os.Stat(fmt.Sprintf("%s/.skip-test", directoryPath))
		shouldSkipTest := !os.IsNotExist(err)
		if shouldSkipTest {
			log.Printf("[DEBUG] Skipping %q since a `.skip-test` is present", directoryPath)
		}

		output = append(output, v)
	}
	return output
}
