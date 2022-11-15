package fields

//***** Metadata text field:strings/texts ******************
type FieldBytea struct {
	Field
}
func (f *FieldBytea) GetDataType() FieldDataType {
	return FIELD_TYPE_BYTEA
}

type FielderBytea interface {
	Fielder
}

//String validaion
func ValidateBytea(f FielderBytea, val []byte) error {
	return nil
}


