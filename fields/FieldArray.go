package fields

import (
	"fmt"
	"errors"
)

type FieldArray struct {
	Field
	MinCount ParamInt
	MaxCount ParamInt
}
func (f *FieldArray) GetDataType() FieldDataType {
	return FIELD_TYPE_ARRAY
}
func (f *FieldArray) GetMinCount() ParamInt {
	return f.MinCount
}
func (f *FieldArray) GetMaxCount() ParamInt {
	return f.MaxCount
}
/*func (f *FieldArray) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *FieldArray) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}*/

type FielderArray interface {
	Fielder
	GetMinCount() ParamInt
	GetMaxCount() ParamInt
}

//Array validaion
func ValidateArray(f FielderArray, val []interface{}) error {
	p := f.GetMaxCount()
	if p.IsSet && len(val) > p.Value {
		return errors.New(fmt.Sprintf(ER_VALID_AR_LEN, f.GetDescr()) )
	}

	return nil
}

