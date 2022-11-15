package fields

import (
	"strconv"
	"errors"
	"database/sql/driver"
)

//Ext type int
type ValInt struct {
	Val
	TypedValue int64
}

func (v ValInt) GetValue() int64{
	if v.IsNull {
		return 0
	}else{
		return v.TypedValue
	}	
}

func (v ValInt) GetIsNull() bool{
	return v.IsNull
}

func (v ValInt) GetIsSet() bool{
	return v.IsSet
}

func (v *ValInt) SetValue(vI int64){
	v.TypedValue = vI
	v.IsSet = true
	v.IsNull = false
}

func (v ValInt) SetNull(){
	v.TypedValue = 0
	v.IsSet = true
	v.IsNull = true
}


//Custom Int unmarshal
func (v *ValInt) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToInt(v_str)
	if err != nil {
		return errors.New(ER_UNMARSHAL_INT + err.Error())
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValInt) String() string {
	return strconv.FormatInt(v.TypedValue, 10)
}

func (v *ValInt) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(v.String()), nil
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValInt) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
	
		if val, err := driver.Int32.ConvertValue(value); err == nil {
			var ok bool
			if v.TypedValue, ok = val.(int64); ok {
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

func (v ValInt) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func StrToInt(vStr string) (int64, error) {
	return strconv.ParseInt(vStr, 10, 64)
}

func NewValInt(val int64, isNull bool) ValInt {
	return ValInt{Val{true, isNull}, val}
}
