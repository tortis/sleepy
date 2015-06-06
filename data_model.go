package sleepy

import (
	"net/http"
	"reflect"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// Sleepy field tags that can be used to signal if the fields of a struct     //
// that is used as the data in/out model are required, readonly, writeony,    //
// or hidden.                                                                 //
//                                                                            //
// - required:  Fields marked as required must be present (non-zero) in the   //
// 			    request body of a POST request. If a required field is not    //
//              present in the request, the handling will end before the API  //
//              call handler sees the request.                                //
//                                                                            //
// - readonly:  Fields marked as readonly must NOT be present in the request  //
//              body of a POST/PUT/PATCH/etc request. If a readonly field is  //
//              present, the handling will end before the API call handler    //
//              sees the request.                                             //
//                                                                            //
// - writeonly: Fields makred as write only will be removed before the        //
//              response is written to the client. This can be useful for     //
//              stopping the password field from being sent to the client.    //
//              This tag should normally accompany the json/xml omitempty     //
//              tag, so that the field will not be marshaled.                 //
////////////////////////////////////////////////////////////////////////////////
const (
	sleepyRequired  = "required"
	sleepyReadOnly  = "readonly"
	sleepyWriteOnly = "writeonly"
	sleepyHidden    = "hidden"
)

////////////////////////////////////////////////////////////////////////////////
// Every API call has a callDataModel that describes its inputs and outputs.  //
// Some of the specifications in callDataModel are puerly for documentation,  //
// but some specifications will be enforced by sleepy.                        //
//                                                                            //
// callDataModel is build when the user makes Call builder method calls.      //
////////////////////////////////////////////////////////////////////////////////
type callDataModel struct {
	bodyIn    modelIn
	bodyOut   modelOut
	pathVars  []inputVar
	queryVars []inputVar
}

////////////////////////////////////////////////////////////////////////////////
// The data model of the request body. This is set using the Call.Reads()     //
// method. The model must be a struct, and sleepy will identify all of fields //
// that are tagged with tags that are relevant to input (required, readonly)  //
// Sleepy tags on the model are enforced, so if a modelIn field is taged as   //
// required, but is not present in the request, the call will be terminateda  //
// with a semantic http error 422.                                            //
////////////////////////////////////////////////////////////////////////////////
type modelIn struct {
	model          interface{}
	requiredFields [][]int
	roFields       [][]int
}

////////////////////////////////////////////////////////////////////////////////
// The data model of the response body. This is set using the Call.Returns()  //
// method. The model must be a struct, and sleepy will identify all of fields //
// that are tagged with tags that are relevant to input (writeonly, hidden))  //
// Sleepy tags on the model are enforced, so if a modelOut field is taged as  //
// writeonly then that field will be zero'd in the response before it is sent //
// to the client.                                                             //
////////////////////////////////////////////////////////////////////////////////
type modelOut struct {
	model        interface{}
	woFields     [][]int
	hiddenFields [][]int
}

////////////////////////////////////////////////////////////////////////////////
// inputVar represents an input to the call that is not part of the request   //
// body. These may be URL path variables, URL query variables, or header      //
// variables. They are generated using Call builder methods, and their        //
// properties are currently not enforced, but may be in a later version.      //
////////////////////////////////////////////////////////////////////////////////
type inputVar struct {
	typ      string
	name     string
	desc     string
	required bool
}

////////////////////////////////////////////////////////////////////////////////
// Function responsible for identifying all sleepy field tags that are        //
// relevant to request body (required, readonly). It uses a recursive         //
// strategy to read all tags of fields in embeded structs. The discovered     //
// field indices are stored in their respective  slices in modelIn.           //
////////////////////////////////////////////////////////////////////////////////
func (cdm *callDataModel) identifyFieldTagsIn(pos []int) {
	var curType reflect.Type
	if pos == nil {
		curType = reflect.TypeOf(cdm.bodyIn.model)
	} else {
		curType = reflect.TypeOf(cdm.bodyIn.model).FieldByIndex(pos).Type
	}
	// assert that kind at pos is struct
	if curType.Kind() != reflect.Struct {
		log.Critical("The kind at pos is not a struct")
	}
	for curPos := 0; curPos < curType.NumField(); curPos++ {
		tags := strings.Split(curType.Field(curPos).Tag.Get("sleepy"), ",")
		for _, tag := range tags {
			switch tag {
			case sleepyRequired:
				cdm.bodyIn.requiredFields = append(cdm.bodyIn.requiredFields, append(pos, curPos))
			case sleepyReadOnly:
				cdm.bodyIn.roFields = append(cdm.bodyIn.roFields, append(pos, curPos))
			}
		}
		// recursive call
		if curType.Field(curPos).Type.Kind() == reflect.Struct {
			cdm.identifyFieldTagsIn(append(pos, curPos))
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Function responsible for identifying all sleepy field tags that are        //
// relevant to response body (writeonly, hidden). It uses a recusrive         //
// strategy to read all tags of fields in embeded structs. The discovered     //
// field indices are stored in their respective slices in modelOut.           //
////////////////////////////////////////////////////////////////////////////////
func (cdm *callDataModel) identifyFieldTagsOut(pos []int) {
	var curType reflect.Type
	if pos == nil {
		curType = reflect.TypeOf(cdm.bodyOut.model)
	} else {
		curType = reflect.TypeOf(cdm.bodyOut.model).FieldByIndex(pos).Type
	}
	// assert that kind at pos is struct
	if curType.Kind() != reflect.Struct {
		log.Info("The kind at pos is not a struct")
		return
	}
	for curPos := 0; curPos < curType.NumField(); curPos++ {
		tags := strings.Split(curType.Field(curPos).Tag.Get("sleepy"), ",")
		for _, tag := range tags {
			switch tag {
			case sleepyWriteOnly:
				cdm.bodyOut.woFields = append(cdm.bodyOut.woFields, append(pos, curPos))
			case sleepyHidden:
			}
		}
		// recursive call
		if curType.Field(curPos).Type.Kind() == reflect.Struct {
			cdm.identifyFieldTagsOut(append(pos, curPos))
		}
	}
}

func (cdm *callDataModel) validateQueryVars(r *http.Request, d CallData) *Error {
	for _, queryVar := range cdm.queryVars {
		// Check if required vars are missing
		if queryVar.required && r.FormValue(queryVar.name) == "" {
			return ErrBadRequest("Failed while validating query variables.", "Required query variable '"+queryVar.name+"' is missing.", ERR_FIELD_MISSING)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Function responsible for validating the fields of a payload against the    //
// sleepy tags of the dataIn data model. Loop over all required and readonly  //
// fields. Ensure that required fields are not zero values, ensure that       //
// readonly fields ARE zero values.                                           //
////////////////////////////////////////////////////////////////////////////////
func (cdm *callDataModel) validateTagsIn(payload interface{}) *Error {
	pValue := reflect.ValueOf(payload).Elem()
	mType := reflect.TypeOf(cdm.bodyIn.model)
	// Ensure required fields are not zeros
	// TODO add an exception for bool since default value is false.
	for _, reqField := range cdm.bodyIn.requiredFields {
		if isZero(pValue.FieldByIndex(reqField).Interface()) {
			return ErrBadRequest("Failed while validating tags for the payload.", "Required field: "+mType.FieldByIndex(reqField).Name+" is missing.", ERR_FIELD_MISSING)
		}
	}

	// Ensure read only fields are not present
	for _, roField := range cdm.bodyIn.roFields {
		if !isZero(pValue.FieldByIndex(roField).Interface()) {
			return ErrBadRequest("Failed while validating tags for the payload.", "Attempting to set read-only field: "+mType.FieldByIndex(roField).Name+".", ERR_MOD_RO_FIELD)
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// A helper function to check if the given variable is an instance of its     //
// type's zero value.                                                         //
////////////////////////////////////////////////////////////////////////////////
func isZero(x interface{}) bool {
	if reflect.TypeOf(x).Kind() == reflect.Slice {
		log.Info("The object is a slice.")
		if reflect.ValueOf(x).Len() == 0 {
			return true
		} else {
			return false
		}
	} else {
		return x == reflect.Zero(reflect.TypeOf(x)).Interface()
	}
}
