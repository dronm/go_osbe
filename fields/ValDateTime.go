package fields

import (
	"time"
	"errors"
	"fmt"
)

const FORMAT_DATE_TIME = "2006-01-02T15:04:05"

type ValDateTime struct {
	ValDateTimeTZ
}

//Custom Float unmarshal
func (v *ValDateTime) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	v.TypedValue = time.Time{} 
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToTime(v_str, FORMAT_DATE_TIME)
	if err != nil {
		return err
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValDateTime) String() string {
	if v.IsNull {
		return ""
	}
	return v.TypedValue.Format(FORMAT_DATE_TIME)
}

func (v *ValDateTime) Scan(value interface{}) error {
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
				val_t, err := StrToTime(val, FORMAT_DATE_TIME)	
				if err != nil {
					return err
				}
				v.TypedValue = val_t
				return nil
		}	
		return errors.New(ER_UNMARSHAL_TIME + "unsupported value for date time")
		
	}
	return nil
}

func (v *ValDateTime) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(fmt.Sprintf(`"%s"`, v.TypedValue.Format(FORMAT_DATE))), nil
	}
}

