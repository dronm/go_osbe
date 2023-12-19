package fields

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"database/sql/driver"	
	"fmt"
	"strings"
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
	var s strings.Builder
	s.WriteString("[")
	for i, v := range v.TypedValue {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(`"`)
		s.WriteString(fmt.Sprintf("%v", v))
		s.WriteString(`"`)
	}
	s.WriteString("]")
	return s.String()
}

//Custom Array unmarshal
func (v *ValArray) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	var temp []interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.New(ER_UNMARSHAL_ARRAY + err.Error())
	}
	v.TypedValue = temp
	
	return nil	
}

func (v *ValArray) MarshalJSON() ([]byte, error) {
	if v.IsNull || v.TypedValue == nil || len(v.TypedValue) == 0 {
		return []byte(JSON_NULL), nil
		
	}else{
		//json.Marshal(v.TypedValue)
		return []byte(v.String()), nil
	}
}

func (v *ValArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	t := xml.StartElement{Name: xml.Name{"", "json"}}
	value := v.String()
	tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()	
}

func (v *ValArray) SetNull(){
	v.TypedValue = []interface{}{}
	v.IsSet = true
	v.IsNull = true
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValArray) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{

		val_s, ok := value.(string)
		if !ok {
			return errors.New(ER_UNMARSHAL_ARRAY + "unsupported value")
		}
		if val_s == "" {
			return nil
		}
		if val_s[0:1] != "{" || val_s[len(val_s)-1:] !=  "}" {
			return errors.New(ER_UNMARSHAL_ARRAY + "unsupported value")
		}		
		values := strings.Split(val_s[1:len(val_s)-1], ",")
		v.TypedValue = make([]interface{}, len(values))
		for i, val := range values {
			v.TypedValue[i] = val
		}
	}
	return nil
}

func (v ValArray) Value() (driver.Value, error) {	
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

