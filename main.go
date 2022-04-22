package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/oleiade/reflections"
)

type Attributes struct {
	Legal Legal
}

type Legal struct {
	Foo    *string
	Bar    string
	Energy Energy
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
			v := *value

			t := reflect.TypeOf(v)
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
		default:
			{
				output[prefix+field] = value
			}
		}
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func main() {

	output := make(map[string]interface{})

	value := &Attributes{
		Legal: Legal{
			Foo: strPtr("foo"),
			Bar: "bar",
			Energy: Energy{
				EPCLevel:         1,
				TotalConsumption: 2,
				Class:            "class",
			}}}

	err := deepStructToMap(value, output, "")
	if err != nil {
		panic(err)
	}

	for k, v := range output {
		fmt.Printf("%s: %v\n", k, v)
	}

}
