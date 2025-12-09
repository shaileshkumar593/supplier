package array

import (
	"reflect"
)

// InArray checks if the given value exists in an array
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// KeyExist checks if key exist in map string
func KeyExist(find string, source map[string]string) (exists bool) {
	_, exists = source[find]
	return
}

// Flip exchanges all keys with corresponding value
func Flip(source map[string]string) (output map[string]string) {
	output = make(map[string]string)
	for key, value := range source {
		output[value] = key
	}
	return
}

// Unique remove duplicates
func Unique(source []string) (output []string) {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	for v := range source {
		if encountered[source[v]] != true {
			encountered[source[v]] = true
			output = append(output, source[v])
		}
	}
	return
}
