package fields

//***** Metadata text field:strings/texts ******************
type FieldJSON struct {
	Field
}
func (f *FieldJSON) GetDataType() FieldDataType {
	return FIELD_TYPE_JSON
}

type FielderJSON interface {
	Fielder
}

//String validaion
func ValidateJSON(f FielderJSON, val []byte) error {
	return nil
}


