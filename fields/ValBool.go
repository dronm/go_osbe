package fields

import (
	"errors"
	"database/sql/driver"
	"strings"
)

const (
	JSON_TRUE = "true"
	JSON_FALSE = "false"
)

type ValBool struct {
	Val
	TypedValue bool
}

func (v ValBool) GetValue() bool{
	if v.IsNull {
		return false
	}else{
		return v.TypedValue
	}	
}

func (v *ValBool) SetValue(vV bool){
	v.TypedValue = vV
	v.IsSet = true
	v.IsNull = false
}

func (v ValBool) GetIsNull() bool{
	return v.IsNull
}

func (v ValBool) GetIsSet() bool{
	return v.IsSet
}

//Custom Bool unmarshal
func (v *ValBool) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	v.TypedValue,_ = StrToBool(v_str)
	//v.TypedValue = (v_str == "true" || v_str == "yes" || v_str == "1" || v_str == "да")
	
	return nil	
}

func (v *ValBool) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else if v.TypedValue {
		return []byte(JSON_TRUE), nil
	}else{
		return []byte(JSON_FALSE), nil
	}
}

func (v ValBool) String() string {
	if v.IsNull {
		return ""
	}else if v.TypedValue {
		return JSON_TRUE
	}else{
		return JSON_FALSE
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValBool) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
	
		if val, err := driver.Bool.ConvertValue(value); err == nil {
			var ok bool
			if v.TypedValue, ok = val.(bool); ok {
				return nil
			}else{
				return errors.New(ER_UNMARSHAL_BOOL + "unsupported value")
			}
		}else{
			return errors.New(ER_UNMARSHAL_BOOL + err.Error())
		}		
	}
	return nil
}

func (v ValBool) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func (v ValBool) SetNull(){
	v.TypedValue = false
	v.IsSet = true
	v.IsNull = true
}

func StrToBool(vStr string) (bool, error) {
	return (vStr == "true"  || vStr == "yes" || vStr == "1" || vStr == "да" || strings.ToUpper(vStr) == "TRUE" || strings.ToUpper(vStr) == "YES" || strings.ToUpper(vStr) == "ДА"), nil
}

func NewValBool(val bool, isNull bool) ValBool{
	return ValBool{Val{true, isNull}, val}
}
