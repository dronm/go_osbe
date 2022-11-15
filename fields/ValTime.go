package fields

import (
	"time"
	"errors"
	"fmt"
)


const FORMAT_TIME = "15:04:05"

type ValTime struct {
	ValDateTimeTZ
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
	temp, err := StrToTime(v_str, FORMAT_TIME)
	if err != nil {
		return err
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValTime) String() string {
	if v.IsNull {
		return ""
	}
	return v.TypedValue.Format(FORMAT_TIME)
}

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
				val_t, err := StrToTime(val, FORMAT_TIME)	
				if err != nil {
					return err
				}
				v.TypedValue = val_t
				return nil
		}	
		return errors.New(ER_UNMARSHAL_TIME + "unsupported value for time")
		
	}
	return nil
}

func (v *ValTime) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(fmt.Sprintf(`"%s"`, v.TypedValue.Format(FORMAT_DATE))), nil
	}
}

