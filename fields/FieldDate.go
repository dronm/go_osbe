package fields

//***** Metadata text field:strings/texts ******************
type FieldDate struct {
	Field
}
func (f *FieldDate) GetDataType() FieldDataType {
	return FIELD_TYPE_DATE		
}

//String validaion
func ValidateDate(f Fielder, val string) error {
	//ToDo

	return nil
}


