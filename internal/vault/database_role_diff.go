package vault

// DatabaseRoleInfo is defined in database.go

// DatabaseRoleDiff represents the differences between two sets of database roles.
type DatabaseRoleDiff struct {
	// Added contains roles present in the new set but not the old set.
	Added []DatabaseRoleInfo
	// Removed contains roles present in the old set but not the new set.
	Removed []DatabaseRoleInfo
	// Changed contains roles present in both sets but with differing configurations.
	Changed []DatabaseRoleInfo
}

// HasChanges returns true if there are any added, removed, or changed roles.
func (d DatabaseRoleDiff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
