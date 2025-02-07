package provider

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Returns true of the existing Set type field has the
//
//	same values as the list of values (as a slice) returned by
//	API call. This way, we prevent need for running plan
//	to think there is a need to update an attribute that is a
//	set if the only thing that changed is the order of all values.
func SetAndResponseMatch(setAttribute types.Set, responseData []int) bool {
	var existingSlices []int

	for val := range setAttribute.Elements() {
		existingSlices = append(existingSlices, val)
	}

	slices.Sort(responseData)
	slices.Sort(existingSlices)

	return slices.Equal(responseData, existingSlices)
}
