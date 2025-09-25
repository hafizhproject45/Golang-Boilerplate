package utils

// -------------------------------------------------------------------
// FlagType
// -------------------------------------------------------------------
type FlagType string

const (
	FlagIsActive FlagType = "IS_ACTIVE"
)

// -------------------------------------------------------------------
// Validators
// -------------------------------------------------------------------

func IsValidFlagType(v string) bool {
	switch FlagType(v) {
	case FlagIsActive:
		return true
	}
	return false
}

// example use

/**
if !utils.IsValidFlagType(req.FlagName) {
    return fiber.NewError(fiber.StatusBadRequest, "invalid flag type")
}
*/
