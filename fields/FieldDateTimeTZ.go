package fields

//***** Metadata text field:strings/texts ******************
type FieldDateTimeTZ struct {
	Field
}
func (f *FieldDateTimeTZ) GetDataType() FieldDataType {
	return FIELD_TYPE_DATETIMETZ		
}

//String validaion
func ValidateDateTimeTZ(f Fielder, val string) error {
	//ToDo

	return nil
}


