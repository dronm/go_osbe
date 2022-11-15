package fields

import (
	"strconv"
	"errors"
	"database/sql/driver"
)

//Ext type int
type ValUint struct {
	Val
	TypedValue uint64
}

func (v ValUint) GetValue() uint64{
	if v.IsNull {
		return 0
	}else{
		return v.TypedValue
	}	
}

func (v ValUint) GetIsNull() bool{
	return v.IsNull
}

func (v ValUint) GetIsSet() bool{
	return v.IsSet
}

func (v *ValUint) SetValue(vI uint64){
	v.TypedValue = vI
	v.IsSet = true
	v.IsNull = false
}

func (v ValUint) SetNull(){
	v.TypedValue = 0
	v.IsSet = true
	v.IsNull = true
}


//Custom Int unmarshal
func (v *ValUint) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToUint(v_str)
	if err != nil {
		return errors.New(ER_UNMARSHAL_INT + err.Error())
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValUint) String() string {
	return strconv.FormatUint(v.TypedValue, 10)
}

func (v *ValUint) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(v.String()), nil
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValUint) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
	
		if val, err := driver.Int32.ConvertValue(value); err == nil {
			var ok bool
			if v.TypedValue, ok = val.(uint64); ok {
				return nil
			}else{
				return errors.New(ER_UNMARSHAL_INT + "unsupported value")
			}
		}else{
			return errors.New(ER_UNMARSHAL_INT + err.Error())
		}		
	}
	return nil
}

func (v ValUint) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func StrToUint(vStr string) (uint64, error) {
	return strconv.ParseUint(vStr, 10, 64)
}

func NewValUint(val uint64, isNull bool) ValUint{
	return ValUint{Val{true, isNull}, val}
}
