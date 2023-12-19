package fields

import (
	"errors"
	"database/sql/driver"
	"time"
	"encoding/json"
	"strings"
//	"fmt"
)

const (
	FORMAT_DATE_TIME_TZ1 string = "2006-01-02T15:04:05.000-07"
	FORMAT_DATE_TIME_TZ2 string = "2006-01-02T15:04:05-07:00"
	FORMAT_DATE_TIME_TZ3 string = "2006-01-02T15:04:05Z07:00"
)

type ValDateTimeTZ struct {
	Val
	TypedValue time.Time
}

func (v ValDateTimeTZ) GetValue() time.Time{
	if v.IsNull {
		return time.Time{}
	}else{
		return v.TypedValue
	}	
}

func (v ValDateTimeTZ) GetIsNull() bool{
	return v.IsNull
}

func (v ValDateTimeTZ) GetIsSet() bool{
	return v.IsSet
}

func (v *ValDateTimeTZ) SetValue(vT time.Time){
	v.TypedValue = vT
	v.IsSet = true
	v.IsNull = false
}

func (v *ValDateTimeTZ) SetNull(){
	v.TypedValue = time.Time{}
	v.IsSet = true
	v.IsNull = true
}

//Custom Float unmarshal
func (v *ValDateTimeTZ) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	v.TypedValue = time.Time{} 
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	var dt_tmpl string
	if strings.Contains(v_str, "+") {
		dt_tmpl = FORMAT_DATE_TIME_TZ2
		
	}else if strings.Contains(v_str, "Z") {
		dt_tmpl = FORMAT_DATE_TIME_TZ3
		
	}else{
		dt_tmpl = FORMAT_DATE_TIME_TZ1
	}
	
	temp, err := StrToTime(v_str, dt_tmpl)
	if err != nil {
		return err
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValDateTimeTZ) String() string {
	if v.IsNull {
		return ""
	}
	return v.TypedValue.Format(FORMAT_DATE_TIME_TZ1)
}

func (v *ValDateTimeTZ) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(v.TypedValue)
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValDateTimeTZ) Scan(value interface{}) error {
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
				val_t, err := StrToTime(val, FORMAT_DATE_TIME_TZ1)	
				if err != nil {
					return err
				}
				v.TypedValue = val_t
		}	
		return errors.New(ER_UNMARSHAL_TIME + "unsupported value")
		
	}
	return nil
}

func (v ValDateTimeTZ) Value() (driver.Value, error) {
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func StrToTime(vStr string, tmpl string) (time.Time, error) {
	temp, err := time.Parse(tmpl, vStr)
	if err != nil {
		return time.Time{}, errors.New(ER_UNMARSHAL_TIME + err.Error())
	}
	return temp, nil
}

