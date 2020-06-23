package listflag

import (
	"fmt"
	"strings"
)

const (
	stringSeparator = ", "
)

// StringList can be used with std flag package to be used while setting multiple process arguments by repeating the named argument
type StringList struct {
	values        []string
	allowedValues map[string]bool
}

// NewStringList creates a new ListFlag and sets allowed values for a flag
// If list of allowed values contains no values then all possible string values are allowed
func NewStringList(allowedValues []string) *StringList {
	allowed := make(map[string]bool)

	for _, value := range allowedValues {
		allowed[value] = true
	}

	return &StringList{
		allowedValues: allowed,
	}
}

// Strings returns string representation of values that were set
func (f *StringList) String() string {
	if f == nil {
		return ""
	}

	return strings.Join(f.values, stringSeparator)
}

// Allowed returns whether the value is allowed to be set
// Always returns true if no allowed values were provided
func (f StringList) Allowed(value string) bool {
	if len(f.allowedValues) == 0 {
		return true
	}

	_, ok := f.allowedValues[value]

	return ok
}

// Values returns set values of the flag
func (f StringList) Values() []string {
	return f.values
}

// IsBoolFlag returns false as the list flag is not a bool flag; used by the flag package
func (f StringList) IsBoolFlag() bool {
	return false
}

// Set is used during setting of flags when used with flag package
// When the value is invalid (not allowed), error is being returned
func (f *StringList) Set(value string) error {
	if !f.Allowed(value) {
		return fmt.Errorf("value %s is not allowed", value)
	}

	f.values = append(f.values, value)

	return nil
}
