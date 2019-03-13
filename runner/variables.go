package runner

type AvailableVariable struct {
	// The name of variable
	Name string

	// The name of the Environment Variable which houses the value for this
	EnvKeyName string

	// Should this variable be generated?
	Generate bool
}
