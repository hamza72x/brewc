package models

type OptionalArgs struct {
	Verbose bool
	Threads int
	// DeleteUnusedDependencies is a flag to delete unused dependencies after uninstalling a formula.
	DeleteUnusedDependencies bool
}
