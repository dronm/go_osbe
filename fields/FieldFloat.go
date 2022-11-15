package fields

import (
	"fmt"
	"errors"
	"math"
	"strconv"
)

type FieldFloat struct {
	Field
	MinValue ParamFloat
	MaxValue ParamFloat
	NotZero ParamBool
	Precision ParamInt
	Length ParamInt
}
func (f *FieldFloat) GetDataType() FieldDataType {
	return FIELD_TYPE_FLOAT
}

func (f *FieldFloat) GetMinValue() ParamFloat {
	return f.MinValue
}
func (f *FieldFloat) SetMinValue(v ParamFloat) {
	f.MinValue = v
}

func (f *FieldFloat) GetMaxValue() ParamFloat {
	return f.MaxValue
}
func (f *FieldFloat) SetMaxValue(v ParamFloat) {
	f.MaxValue = v
}

func (f *FieldFloat) GetNotZero() ParamBool {
	return f.NotZero
}
func (f *FieldFloat) SetNotZero(v ParamBool) {
	f.NotZero = v
}

func (f *FieldFloat) GetPrecision() ParamInt {
	return f.Precision
}
func (f *FieldFloat) SetPrecision(v ParamInt) {
	f.Precision = v
}

func (f *FieldFloat) GetLength() ParamInt {
	return f.Length
}
func (f *FieldFloat) SetLength(v ParamInt) {
	f.Length = v
}


type FielderFloat interface {
	Fielder
	GetMinValue() ParamFloat
	SetMinValue(ParamFloat)
	GetMaxValue() ParamFloat
	SetMaxValue(ParamFloat)
	GetNotZero() ParamBool
	SetNotZero(ParamBool)
	GetPrecision() ParamInt
	SetPrecision(ParamInt)
	GetLength() ParamInt
	SetLength(ParamInt)
}

//Float validaion
func ValidateFloat(f FielderFloat, val float64) error {
	if f.GetMinValue().IsSet && val < f.GetMinValue().Value {
		return errors.New(fmt.Sprintf(ER_VALID_MINVAL,f.GetDescr()) )
	}
	
	if f.GetMaxValue().IsSet && val > f.GetMaxValue().Value {
		return errors.New(fmt.Sprintf(ER_VALID_MAXVAL,f.GetDescr()) )
	}
	
	if f.GetNotZero().IsSet && val == 0 {
		return errors.New(fmt.Sprintf(ER_VALID_ZEROVAL,f.GetDescr()) )
	}
	
	if f.GetPrecision().IsSet {		
		prec := f.GetPrecision().Value
		value_f := val * math.Pow(10.0, float64(prec))
		is_whole := int(value_f * 10.0 * float64(prec)) == int(value_f) * 10 * prec
		if !is_whole {
			return errors.New(fmt.Sprintf(ER_VALID_PRECISION,f.GetDescr()) )
		}		
	}
	
	//+Length
	if f.GetLength().IsSet {
		val_s := strconv.FormatFloat(val, 'f', -1, 64)
		if len(val_s) > f.GetLength().Value {
			return errors.New(fmt.Sprintf(ER_VALID_LEN,f.GetDescr()) )
		}
	}
	
	return nil
}

