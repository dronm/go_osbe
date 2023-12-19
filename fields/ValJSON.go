package fields

import (
	"encoding/xml"
	"errors"
	"bytes"
	"database/sql/driver"	
//	"fmt"
)

const SINGLE_QUOTE_CHAR = 39

type ValJSON struct {
	Val
	TypedValue []byte
}

func (v ValJSON) String() string {
	return string(v.GetValue())
}

func (v ValJSON) GetValue() []byte{
	if v.IsNull {
		return []byte{}
	}else{
		return v.TypedValue
	}	
}

func (v *ValJSON) SetValue(vAr []byte){
	v.TypedValue = vAr
	v.IsSet = true
	v.IsNull = false
}

func (v *ValJSON) SetNull(){
	v.TypedValue = []byte{}
	v.IsSet = true
	v.IsNull = true
}

func (v ValJSON) GetIsNull() bool{
	return v.IsNull
}

func (v ValJSON) GetIsSet() bool{
	return v.IsSet
}


//Custom String unmarshal
func (v *ValJSON) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	
	if (data[0] == QUOTE_CHAR && byte(data[len(data)-1]) ==  QUOTE_CHAR) ||
	(data[0] == SINGLE_QUOTE_CHAR && byte(data[len(data)-1]) ==  SINGLE_QUOTE_CHAR) {
		//serialized string
		//data[0] = byte(39)
		//data[len(data)-1] = byte(39)
		data = bytes.ReplaceAll(data, []byte(`\"`),[]byte(`"`))
		
		//Added on 05/04/22 qotes inside JSON!!!
		data = bytes.ReplaceAll(data, []byte(`\\"`),[]byte(`\"`))
		
		data = data[1:len(data)-1]
	}
	v.TypedValue = data
//fmt.Println("ValJSON.UnmarshalJSON v.TypedValue=", v.TypedValue)
	return nil	
}

func (v *ValJSON) MarshalJSON() ([]byte, error) {
	if v.IsNull || len(v.TypedValue) == 0 {
		return []byte(JSON_NULL), nil
		
	}else{
		return v.TypedValue, nil
	}
}

func (v *ValJSON) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	t := xml.StartElement{Name: xml.Name{"", "json"}}
	value, err := v.MarshalJSON()
	if err != nil {
		return err
	}
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

//driver.Scanner, driver.Valuer interfaces
func (v *ValJSON) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
		var ok bool
		if v.TypedValue, ok = value.([]byte); ok {
			return nil

		}else if v_str, ok := value.(string); ok {
			v.TypedValue = []byte(v_str)
			return nil			
			
		}else{
			return errors.New(ER_UNMARSHAL_JSON + "unsupported value")
		}
	}
	return nil
}

func (v ValJSON) Value() (driver.Value, error) {	
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func NewValJSON(val []byte, isNull bool) ValJSON{
	return ValJSON{Val{true, isNull}, val}
}
