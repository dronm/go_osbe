package fields

import (
	"encoding/json"
	"errors"
	"database/sql/driver"
)

type ValText struct {
	Val
	TypedValue string
}

func (v ValText) String() string {
	return v.GetValue()
}

func (v ValText) GetValue() string{
	if v.IsNull {
		return ""
	}else{
		return v.TypedValue
	}	
}

func (v *ValText) SetValue(vStr string){
	v.TypedValue = vStr
	v.IsSet = true
	v.IsNull = false
}

func (v *ValText) SetNull(){
	v.TypedValue = ""
	v.IsSet = true
	v.IsNull = true
}

func (v ValText) GetIsNull() bool{
	return v.IsNull
}

func (v ValText) GetIsSet() bool{
	return v.IsSet
}


//Custom String unmarshal
func (v *ValText) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.New(ER_UNMARSHAL_STRING + err.Error())
	}
	
	v.TypedValue = temp
	
	return nil	
}

func (v *ValText) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(v.TypedValue)
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValText) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
		if val, err := driver.String.ConvertValue(value); err == nil {
			var ok bool
			if v.TypedValue, ok = val.(string); ok {
				return nil
				
			}else if v_bt, ok := val.([]byte); ok {
				v.TypedValue = string(v_bt)
				
			}else{
				return errors.New(ER_UNMARSHAL_STRING + "unsupported value")
			}
		}else{
			return errors.New(ER_UNMARSHAL_STRING + err.Error())
		}		
	}
	return nil
}

func (v ValText) Value() (driver.Value, error) {	
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func NewValText(val string, isNull bool) ValText{
	return ValText{Val{true, isNull}, val}
}
