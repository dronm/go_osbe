package fields

import (
	"fmt"
	"errors"
)

type FieldInt struct {
	Field
	MinValue ParamInt64
	MaxValue ParamInt64
	NotZero ParamBool
}
func (f *FieldInt) GetDataType() FieldDataType {
	return FIELD_TYPE_INT
}
func (f *FieldInt) GetMinValue() ParamInt64 {
	return f.MinValue
}
func (f *FieldInt) SetMinValue(v ParamInt64) {
	f.MinValue = v
}

func (f *FieldInt) GetMaxValue() ParamInt64 {
	return f.MaxValue
}
func (f *FieldInt) SetMaxValue(v ParamInt64) {
	f.MaxValue = v
}

func (f *FieldInt) GetNotZero() ParamBool {
	return f.NotZero
}
func (f *FieldInt) SetNotZero(v ParamBool) {
	f.NotZero = v
}

/*func (f *FieldInt) GetPrimaryKey() bool {
	return f.PrimaryKey
}
func (f *FieldInt) SetPrimaryKey(v bool) {
	f.PrimaryKey = v
}*/

type FielderInt interface {
	Fielder
	GetMinValue() ParamInt64
	SetMinValue(ParamInt64)
	GetMaxValue() ParamInt64
	SetMaxValue(ParamInt64)
	GetNotZero() ParamBool
	SetNotZero(ParamBool)
}

//Int validaion
func ValidateInt(f FielderInt, val int64) error {
	if f.GetMinValue().IsSet && val < f.GetMinValue().Value {
		return errors.New(fmt.Sprintf(ER_VALID_MINVAL, f.GetDescr()) )
	}
	
	if f.GetMaxValue().IsSet && val > f.GetMaxValue().Value {
		return errors.New(fmt.Sprintf(ER_VALID_MAXVAL, f.GetDescr()) )
	}
	
	if f.GetNotZero().IsSet && val == 0 {
		return errors.New(fmt.Sprintf(ER_VALID_ZEROVAL, f.GetDescr()) )
	}

	return nil
}

