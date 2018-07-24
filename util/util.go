package util

import (
	"errors"
	"fmt"
	"reflect"
)

func MakeStringSlice(from interface{}) ([]string, error) {
	empty := make([]string, 0)
	if from == nil {
		return empty, nil
	}

	slice := []string{}
	if reflect.TypeOf(from).Kind() == reflect.Slice {
		tmp := reflect.ValueOf(from)
		for i := 0; i < tmp.Len(); i++ {
			slice = append(slice, fmt.Sprintf("%v", tmp.Index(i)))
		}
		return slice, nil
	} else {
		return nil, errors.New("Unable to convert to string slice")
	}
}
