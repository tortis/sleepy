package sleepy

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
const (
	sleepyRequired  = "required"
	sleepyReadOnly  = "readonly"
	sleepyWriteOnly = "writeonly"
	sleepyHidden    = "hidden"
)

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
type callDataModel struct {
	bodyIn         interface{}
	requiredFields [][]int
	roFields       [][]int
	woFields       [][]int
	hiddenFields   [][]int
	bodyOut        interface{}
	pathVarsDoc    []inputVar
	queryVarsDoc   []inputVar
}

func (model *callDataModel) identifyFieldTags(pos []int) {
	var curType reflect.Type
	if pos == nil {
		curType = reflect.TypeOf(model.bodyIn)
	} else {
		curType = reflect.TypeOf(model.bodyIn).FieldByIndex(pos).Type
	}
	// assert that kind at pos is struct
	if curType.Kind() != reflect.Struct {
		log.Fatal("The kind at pos is not a struct")
	}
	for curPos := 0; curPos < curType.NumField(); curPos++ {
		tags := strings.Split(curType.Field(curPos).Tag.Get("sleepy"), ",")
		for _, tag := range tags {
			switch tag {
			case sleepyRequired:
				model.requiredFields = append(model.requiredFields, append(pos, curPos))
			case sleepyReadOnly:
				model.roFields = append(model.roFields, append(pos, curPos))
			}
		}
		// recursive call
		if curType.Field(curPos).Type.Kind() == reflect.Struct {
			model.identifyFieldTags(append(pos, curPos))
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func (model *callDataModel) validateTagsIn(payload interface{}) Error {
	pValue := reflect.ValueOf(payload).Elem()
	mType := reflect.TypeOf(model.bodyIn)
	// Ensure required fields are not zeros
	// TODO add an exception for bool since default value is false.
	for _, reqField := range model.requiredFields {
		if isZero(pValue.FieldByIndex(reqField).Interface()) {
			return newRequestError("Required field: "+mType.FieldByIndex(reqField).Name+" is missing.", errors.New("Failed while validating tags for the payload."))
		}
	}

	// Ensure read only fields are not present
	for _, roField := range model.roFields {
		if isZero(pValue.FieldByIndex(roField).Interface()) {
			return newRequestError("Attempting to set read-only field: "+mType.FieldByIndex(roField).Name+".", errors.New("Failed while validating tags for the payload."))
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//
////////////////////////////////////////////////////////////////////////////////
func isZero(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

type inputVar struct {
	typ      string
	name     string
	desc     string
	required bool
}
