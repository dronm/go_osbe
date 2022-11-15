package fields

//***** Metadata text field:strings/texts ******************
type FieldBool struct {
	Field
}
func (f *FieldBool) GetDataType() FieldDataType {
	return FIELD_TYPE_BOOL
}

/*func (f *FieldBool) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *FieldBool) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}
*/
