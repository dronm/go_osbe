package fields

import (
	"errors"
	"database/sql/driver"
	"time"
	"encoding/json"
//	"fmt"
)

const (
	DATE_TIME_LAYOUT = "2006-01-02T15:04:05.000-07"
	DATE_LAYOUT = "2006-01-02"
)

type ValTime struct {
	Val
	TypedValue time.Time
}

func (v ValTime) GetValue() time.Time{
	if v.IsNull {
		return time.Time{}
	}else{
		return v.TypedValue
	}	
}

func (v ValTime) GetIsNull() bool{
	return v.IsNull
}

func (v ValTime) GetIsSet() bool{
	return v.IsSet
}

func (v *ValTime) SetValue(vT time.Time){
	v.TypedValue = vT
	v.IsSet = true
	v.IsNull = false
}

func (v ValTime) SetNull(){
	v.TypedValue = time.Time{}
	v.IsSet = true
	v.IsNull = true
}

//Custom Float unmarshal
func (v *ValTime) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	v.TypedValue = time.Time{} 
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToTime(v_str)
	if err != nil {
		return err
	}
	v.TypedValue = temp
	
	return nil	
}

func (v *ValTime) String() string {
	//return strconv.FormatFloat(v.TypedValue, 'f', -1, 64)
	return ""
}

func (v *ValTime) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(v.TypedValue)
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValTime) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
		switch val := value.(type) {
			case time.Time:
				v.TypedValue = val
				return nil
			case string:
				val_t, err := StrToTime(val)	
				if err != nil {
					return err
				}
				v.TypedValue = val_t
		}	
		return errors.New(ER_UNMARSHAL_TIME + "unsupported value")
		
	}
	return nil
}

func (v ValTime) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func StrToTime(vStr string) (time.Time, error) {
	temp, err := time.Parse(DATE_TIME_LAYOUT, vStr)
	if err != nil {
		return time.Time{}, errors.New(ER_UNMARSHAL_TIME + err.Error())
	}
	return temp, nil
}

