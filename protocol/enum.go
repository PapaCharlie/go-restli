package protocol

import (
	"fmt"

	"github.com/PapaCharlie/go-restli/fnv1a"
)

// Enum is a class that all generated enums extend. It provides common logic and other utilities.
type Enum int32

// UnknownEnum is the 0-value enum, which is present for all generated enums.
const UnknownEnum = Enum(0)

func (e Enum) IsUnknown() bool {
	return e == UnknownEnum
}

// Equals checks the equality of the two enums. Returns false if either parameter is nil or if either enum in unknown.
func (e Enum) Equals(other Enum) bool {
	if e.IsUnknown() || other.IsUnknown() {
		return false
	}

	return e == other
}

func (e Enum) ComputeHash() fnv1a.Hash {
	if e.IsUnknown() {
		return fnv1a.ZeroHash()
	}
	hash := fnv1a.NewHash()
	hash.AddInt32(int32(e))
	return hash
}

type IllegalEnumConstant struct {
	Enum     string
	Constant int
}

func (i *IllegalEnumConstant) Error() string {
	return fmt.Sprintf("go-restli: Illegal constant for %q enum: %d", i.Enum, i.Constant)
}

type UnknownEnumValue struct {
	Enum  string
	Value string
}

func (u *UnknownEnumValue) Error() string {
	return fmt.Sprintf("go-restli: Unknown enum value for %q: %q", u.Enum, u.Value)
}
