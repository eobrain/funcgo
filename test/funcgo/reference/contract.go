package contract

// Set these booleans to control whether checking happens
var (
	CheckPreconditions = true
	CheckPostconditions = false
)

// precondition is a function returning bool
func Require(precondition) {
    if CheckPreconditions {
	    if !precondition() {
		    throw(new AssertionError("Precondition failed for "  str  precondition))
	    }
    }
}

// postcondition is function of single value, returning boolean
func Ensure(postcondition, result) {
	if CheckPostconditions {
		if !postcondition(result) {
			throw(new AssertionError("Postcondition failed for "  str  result))
		}
	}
	result
}
