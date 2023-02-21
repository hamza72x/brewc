package models

type OptionalArgs struct {
	Verbose bool
	Threads int

	// DeleteUnusedDependencies is a flag to delete unused dependencies after uninstalling a formula.
	DeleteUnusedDependencies bool

	// DeleteAllNestedDependencies is a flag to delete all sub-dependencies after uninstalling a formula (it will delete all nested unused dependencies).
	DeleteAllNestedDependencies bool
}
