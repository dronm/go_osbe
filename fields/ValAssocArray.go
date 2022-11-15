package fields

import (
	"encoding/json"
	"errors"
)

type ValAssocArray struct {
	Val
	TypedValue map[string]interface{}
}
//Custom AssocArray unmarshal
func (v *ValAssocArray) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.New(ER_UNMARSHAL_ASSOC_ARRAY + err.Error())
	}
	v.TypedValue = temp
	
	return nil	
}

func (v ValAssocArray) GetValue() map[string]interface{}{
	if v.IsNull {
		return nil
	}else{
		return v.TypedValue
	}	
}

func (v ValAssocArray) GetIsNull() bool{
	return v.IsNull
}

func (v ValAssocArray) GetIsSet() bool{
	return v.IsSet
}

func (v *ValAssocArray) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		return json.Marshal(v.TypedValue)
	}
}

func (v ValAssocArray) String() string{
	return "<Not implemented>"
}

