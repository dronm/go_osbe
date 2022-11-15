package fields

import (
	"fmt"
	"errors"
)

type FieldAssocArray struct {
	Field
	MinCount ParamInt
	MaxCount ParamInt
}
func (f *FieldAssocArray) GetDataType() FieldDataType {
	return FIELD_TYPE_ARRAY
}
func (f *FieldAssocArray) GetMinCount() ParamInt {
	return f.MinCount
}
func (f *FieldAssocArray) GetMaxCount() ParamInt {
	return f.MaxCount
}
/*func (f *FieldAssocArray) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *FieldAssocArray) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}*/

type FielderAssocArray interface {
	Fielder
	GetMinCount() ParamInt
	GetMaxCount() ParamInt
}

//Array validaion
func ValidateAssocArray(f FielderAssocArray, val []interface{}) error {
	p := f.GetMaxCount()
	if p.IsSet && len(val) > p.Value {
		return errors.New(fmt.Sprintf(ER_VALID_AR_LEN, f.GetDescr()) )
	}

	return nil
}

