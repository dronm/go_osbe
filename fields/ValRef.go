package fields

import (
	"encoding/json"
	"errors"
	"database/sql/driver"	
)

type ValRef struct {
	Val
	TypedValue Ref
}

func (v ValRef) String() string {
	return ""
}

func (v ValRef) GetValue() *Ref {
	if v.IsNull {
		return nil
	}else{
		return &v.TypedValue
	}	
}

func (v *ValRef) SetValue(vRef Ref){
	v.TypedValue = vRef
	v.IsSet = true
	v.IsNull = false
}

func (v *ValRef) SetNull(){
	v.TypedValue = Ref{}
	v.IsSet = true
	v.IsNull = true
}

func (v ValRef) GetIsNull() bool{
	return v.IsNull
}

func (v ValRef) GetIsSet() bool{
	return v.IsSet
}

func (v ValRef) GetKeyAsInt(id string) int64 {
	if v.IsNull {
		return 0
	}
	if v_if, ok := v.TypedValue.Keys[id]; ok {
		if v_int, ok := v_if.(int64); ok {
			return v_int
			
		}else if v_int, ok := v_if.(int); ok {
			return int64(v_int)
		}		
	}
	return 0
}
func (v ValRef) GetKeyAsString(id string) string {
	if v.IsNull {
		return ""
	}
	if v_if, ok := v.TypedValue.Keys[id]; ok {
		if v_str, ok := v_if.(string); ok {
			return v_str			
		}		
	}
	return ""
}


//Custom String unmarshal
func (v *ValRef) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	return json.Unmarshal(data, &v.TypedValue)
}

func (v *ValRef) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(&v.TypedValue)
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValRef) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
		if d, ok := value.([]byte); ok {
			return v.UnmarshalJSON(d)
			
		}else{
			return errors.New(ER_UNMARSHAL_STRING + "unsupported value")
		}
	}
	return nil
}

func (v ValRef) Value() (driver.Value, error) {	
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}


