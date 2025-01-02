// Package structs provides utilities to make working with struct resources easier
package structs

import (
	"encoding/json"
	"reflect"

	"goeasy.dev/errors"
)

var marshallers = map[reflect.Type]Marshaller{}
var unmarshallers = map[reflect.Type]Unmarshaller{}

// Marshaller is a function that will marshal an interface into a byte array
type Marshaller func(interface{}) ([]byte, error)

// Unmarshaller is a function that will populate a struct with the values of a byte array
type Unmarshaller func([]byte, interface{}) error

// RegisterMarshaller is used to set the function that should be used to marshal values of a given type
//
//	structs.RegisterMarshaller(MyStruct{}, json.Marshal)
func RegisterMarshaller(t interface{}, marshaller Marshaller) {
	marshallers[reflect.TypeOf(t)] = marshaller
}

// RegisterMarshaller is used to set the function that should be used to unmarshal bytes into a destination
//
//	structs.RegisterUnmarshaller(MyStruct{}, json.Unmarshal)
func RegisterUnmarshaller(t interface{}, unmarshaller Unmarshaller) {
	unmarshallers[reflect.TypeOf(t)] = unmarshaller
}

// Marshal converts the value into a byte array, returing the error of the marshal function
// If no marshaller is registered for the supplied type, `json.Marshal` is used
func Marshal(value interface{}) ([]byte, error) {
	marshaller, ok := marshallers[reflect.TypeOf(value)]
	if !ok {
		marshaller = json.Marshal
	}

	return marshaller(value)
}

// Unmarshal parses the data, storing the result in dest, returing the error of the marshal function
// If no unmarshaller is registered for the supplied type, `json.Unmarshal` is used
func Unmarshal(data []byte, dest interface{}) error {
	unmarshaller, ok := unmarshallers[reflect.TypeOf(dest)]
	if !ok {
		unmarshaller = json.Unmarshal
	}

	err := unmarshaller(data, dest)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshall data")
	}

	return nil
}
