package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/oleiade/reflections"
)

type Bathroom struct {
	Surface int64
	Name    string
}

type Attributes struct {
	Legal     *Legal
	Bathrooms []Bathroom
}

type Legal struct {
	Foo    *string
	Bar    string
	Energy *Energy
}

type Energy struct {
	EPCLevel         int64
	TotalConsumption int64
	Class            string
}

func deepStructToMap(input any, output map[string]interface{}, prefix string) error {
	fields, err := reflections.FieldsDeep(input)
	if err != nil {
		return err
	}

	for _, field := range fields {
		kind, err := reflections.GetFieldKind(input, field)
		if err != nil {
			return err
		}
		value, err := reflections.GetField(input, field)
		if err != nil {
			return err
		}

		if kind == reflect.Ptr {
			ptr := reflect.ValueOf(value)
			if ptr.IsNil() {
				output[prefix+field] = nil
				continue
			}
			value = ptr.Elem().Interface()
			t := reflect.TypeOf(value)
			kind = t.Kind()
		}

		switch kind {
		case reflect.Struct:
			{
				err = deepStructToMap(value, output, prefix+field+".")
				if err != nil {
					return err
				}
			}
		case reflect.Pointer:
			{
				return errors.New("shit")
			}
		case reflect.Slice:
			{
				slice := reflect.ValueOf(value)
				if slice.IsNil() {
					output[prefix+field] = nil
					continue
				}
				for i := 0; i < slice.Len(); i++ {
					err = deepStructToMap(slice.Index(i).Interface(), output, fmt.Sprintf("%s%s[%d].", prefix, field, i))
					if err != nil {
						return err
					}
				}
			}
		default:
			{
				output[prefix+field] = value
			}
		}
	}
	return nil
}

func getListOfZeroMapFields(input map[string]interface{}) []string {
	var fields []string
	for key, value := range input {
		v := reflect.ValueOf(value)

		if value == nil || v.IsZero() {
			fields = append(fields, key)
		}
	}
	return fields
}

func GetZeroFieldPaths(input any) (res []string, err error) {
	output := make(map[string]interface{})
	err = deepStructToMap(input, output, "")
	if err != nil {
		return
	}

	return getListOfZeroMapFields(output), nil
}

func strPtr(s string) *string {
	return &s
}

func main() {

	value := &Attributes{
		Bathrooms: []Bathroom{
			{
				Name:    "bathroom1",
				Surface: 100,
			},
			{
				Name:    "bathroom2",
				Surface: 200,
			},
		},
		Legal: &Legal{
			Foo: strPtr("foo"),
			Bar: "bar",
			Energy: &Energy{
				EPCLevel:         1,
				TotalConsumption: 2,
				Class:            "class",
			},
		},
	}

	empties, err := GetZeroFieldPaths(value)
	if err != nil {
		panic(err)
	}
	fmt.Printf("empties: %v\n", empties)

}
