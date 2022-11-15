package fields

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ValArray struct {
	Val
	TypedValue []interface{}
}

func (v ValArray) GetValue() []interface{}{
	if v.IsNull {
		return nil
	}else{
		return v.TypedValue
	}	
}

func (v ValArray) GetIsNull() bool{
	return v.IsNull
}

func (v ValArray) GetIsSet() bool{
	return v.IsSet
}

func (v ValArray) String() string {
	s := ""
	for _,vv := range v.TypedValue {
		if s != "" {
			s+= ","
		}
		s+= fmt.Sprintf("%s", vv)
	}
	return s
}

//Custom Array unmarshal
func (v *ValArray) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	var temp []interface{}
//fmt.Println("Unmarshal ",string(data))	
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.New(ER_UNMARSHAL_ARRAY + err.Error())
	}
	v.TypedValue = temp
	
	return nil	
}

func (v *ValArray) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(v.TypedValue)
	}
}

