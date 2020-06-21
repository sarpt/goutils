package cmds

import (
	"fmt"
	"strings"
)

const (
	stringSeparator = ", "
)

// ListFlag can be used with std flag package to be used while setting multiple process arguments by repeating the named argument
type ListFlag struct {
	values        []string
	allowedValues map[string]bool
}

// NewListFlag creates a new ListFlag and sets allowed values for a flag
// If list of allowed values contains no values then all possible string values are allowed
func NewListFlag(allowedValues []string) ListFlag {
	allowed := make(map[string]bool)

	for _, value := range allowedValues {
		allowed[value] = true
	}

	return ListFlag{
		allowedValues: allowed,
	}
}

// Strings returns string representation of values that were set
func (f *ListFlag) String() string {
	if f == nil {
		return ""
	}

	return strings.Join(f.values, stringSeparator)
}

// Allowed returns whether the value is allowed to be set
// Always returns true if no allowed values were provided
func (f ListFlag) Allowed(value string) bool {
	if len(f.allowedValues) == 0 {
		return true
	}

	_, ok := f.allowedValues[value]

	return ok
}

// Values returns set values of the flag
func (f ListFlag) Values() []string {
	return f.values
}

// IsBoolFlag returns false as the list flag is not a bool flag; used by the flag package
func (f ListFlag) IsBoolFlag() bool {
	return false
}

// Set is used during setting of flags when used with flag package
// When the value is invalid (not allowed), error is being returned
func (f *ListFlag) Set(value string) error {
	if !f.Allowed(value) {
		return fmt.Errorf("value %s is not allowed", value)
	}

	f.values = append(f.values, value)

	return nil
}
