package utslice

import "reflect"

func HasDifferentElements(arr interface{}) bool {
	// Get the value of the input
	v := reflect.ValueOf(arr)

	// If it's not a slice or array, return false
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}

	// If the array is empty, return false
	if v.Len() == 0 {
		return false
	}

	firstelem := v.Index(0).Interface()

	for i := 0; i < v.Len(); i++ {
		if !reflect.DeepEqual(firstelem, v.Index(i).Interface()) {
			return true
		}
	}

	return false
}
