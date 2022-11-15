package fields

//***** Metadata text field:strings/texts ******************
type FieldDateTime struct {
	Field
}
func (f *FieldDateTime) GetDataType() FieldDataType {
	return FIELD_TYPE_DATETIME		
}

//String validaion
func ValidateDateTime(f Fielder, val string) error {
	//ToDo

	return nil
}


