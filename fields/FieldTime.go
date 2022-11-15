package fields

import (
	//"encoding/utf8"
	//"errors"
	//"fmt"
)

//***** Metadata text field:strings/texts ******************
type FieldTime struct {
	Field
}
func (f *FieldTime) GetDataType() FieldDataType {
	return FIELD_TYPE_TIME		
}

//String validaion
func ValidateTime(f Fielder, val string) error {
	//ToDo

	return nil
}


