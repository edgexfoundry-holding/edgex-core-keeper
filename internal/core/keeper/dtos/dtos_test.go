package dtos

import (
	"testing"
	"reflect"
	"unicode"
)

// CheckDtoJson verifies the JSON mapping of an EdgeX Data Transfer Object type.
// We check that the struct passed in has all its exported fields JSON-tagged
// with a name that starts with a lowercase letter.
// Embedded structs are recursively checked.
// Fails the test if criteria are not met.
func CheckDtoJson(t *testing.T, ourType reflect.Type) {
	if ourType.Kind() != reflect.Struct {
		t.Fatalf("Type %s in CheckDtoJson() is not a struct but a %s", ourType.String(), ourType.Kind())
	}
	for i := 0; i < ourType.NumField(); i++ {
		thisfield := ourType.Field(i)
		if ! thisfield.IsExported() {
			t.Logf("Skipping unexported field %s", thisfield.Name)
			continue
		}
		if thisfield.Anonymous {
			t.Logf("Checking embedded structure field %s", thisfield.Name)
			CheckDtoJson(t, thisfield.Type)
		} else {
			st := thisfield.Tag.Get("json")
			if st == "" {
				t.Fatalf("Type %s has an exported field %s with no JSON tag", ourType.String(), thisfield.Name)
			}
			first_rune := []rune(st)[0]
			if first_rune == ',' {
				t.Fatalf("Type %s has an exported field %s with JSON tag that has no name", ourType.String(), thisfield.Name)
			}
			if ! unicode.IsLower(first_rune) {
				t.Fatalf("Type %s has an exported field %s with a JSON tag %s that does not start with a lowercase letter",
					ourType.String(), thisfield.Name, st)
			}
			t.Logf("Type %s field %s has a lowercase-starting JSON tag %s", ourType.String(), thisfield.Name, st)
			if thisfield.Type.Kind() == reflect.Struct {
				t.Logf("Checking structure field %s", thisfield.Name)
				CheckDtoJson(t, thisfield.Type)
			}
		}
	}
}

func TestRegistrationDto(t *testing.T) {
	CheckDtoJson(t, reflect.TypeOf(Registration{}))
}
