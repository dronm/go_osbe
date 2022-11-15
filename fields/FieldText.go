package fields

import (
	//"encoding/utf8"
	"errors"
	"fmt"
)

//***** Metadata text field:strings/texts ******************
type FieldText struct {
	Field
	Length ParamInt	
}
func (f *FieldText) GetDataType() FieldDataType {
	return FIELD_TYPE_TEXT
}
func (f *FieldText) GetLength() ParamInt {
	return f.Length
}

func (f *FieldText) SetLength(v ParamInt) {
	f.Length = v
}
/*func (f *FieldText) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *FieldText) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}*/

type FielderText interface {
	Fielder
	GetLength() ParamInt
	SetLength(ParamInt)
}

//String validaion
func ValidateText(f FielderText, val string) error {
	l := f.GetLength()
	if l.IsSet && len([]rune(val)) > l.Value {
		return errors.New(fmt.Sprintf(ER_VALID_LEN,f.GetDescr()) )
	}

	return nil
}


