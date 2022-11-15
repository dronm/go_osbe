package fields

import (
	"time"
	"errors"
	"fmt"
)

const FORMAT_DATE = "2006-01-02"

type ValDate struct {
	ValDateTimeTZ
}

//Custom Float unmarshal
func (v *ValDate) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	v.TypedValue = time.Time{} 
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	v_str := ExtRemoveQuotes(data)
	temp, err := StrToTime(v_str, FORMAT_DATE)
	if err != nil {
		return err
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValDate) String() string {
	if v.IsNull {
		return ""
	}
	return v.TypedValue.Format(FORMAT_DATE)
}

func (v *ValDate) Scan(value interface{}) error {
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
				val_t, err := StrToTime(val, FORMAT_DATE)	
				if err != nil {
					return err
				}
				v.TypedValue = val_t
		}	
		return errors.New(ER_UNMARSHAL_TIME + "unsupported value")
		
	}
	return nil
}

func (v *ValDate) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return []byte(fmt.Sprintf(`"%s"`, v.TypedValue.Format(FORMAT_DATE))), nil
	}
}

