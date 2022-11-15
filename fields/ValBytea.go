package fields

import (
	"errors"
	"database/sql/driver"	
	"encoding/base64"	
	"encoding/hex"
)

type ValBytea struct {
	Val
	TypedValue []byte
}

func (v ValBytea) String() string {
	return hex.EncodeToString(v.TypedValue)
}

func (v ValBytea) GetValue() []byte{
	if v.IsNull {
		return []byte{}
	}else{
		return v.TypedValue
	}	
}

func (v *ValBytea) SetValue(vAr []byte){
	v.TypedValue = vAr
	v.IsSet = true
	v.IsNull = false
}

func (v ValBytea) SetNull(){
	v.TypedValue = []byte{}
	v.IsSet = true
	v.IsNull = true
}

func (v ValBytea) GetIsNull() bool{
	return v.IsNull
}

func (v ValBytea) GetIsSet() bool{
	return v.IsSet
}

func (v *ValBytea) Len() int {
	if !v.IsNull {
		return len(v.TypedValue)
	}
	return 0
}

//Custom String unmarshal
//incoming string in Base64
func (v *ValBytea) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
	var err error
	/*
	var temp []byte
	if err = json.Unmarshal(data, &temp); err != nil {
		return errors.New(ER_UNMARSHAL_BYTEA + err.Error())
	}
	*/
	//quotes
	if data[0] == QUOTE_CHAR && byte(data[len(data)-1]) ==  QUOTE_CHAR {
		data = data[1:len(data)-1]
	}		
	v.TypedValue, err = Base64Decode(data)
	if err != nil {
		return errors.New(ER_UNMARSHAL_BYTEA + err.Error())
	}
	
	return nil	
}

//to base 64
func (v *ValBytea) MarshalJSON() ([]byte, error) {
	if v.IsNull {
		return []byte(JSON_NULL), nil
		
	}else{
		vl := Base64Encode(v.TypedValue)
		return []byte(`"` + string(vl) + `"`), nil
	}
}

//driver.Scanner, driver.Valuer interfaces
func (v *ValBytea) Scan(value interface{}) error {
	v.IsSet = true
	v.IsNull = false
	if value == nil {
		v.IsNull = true
		return nil
	}else{
		var ok bool
		if v.TypedValue, ok = value.([]byte); ok {
			return nil
			
		}else{
			return errors.New(ER_UNMARSHAL_BYTEA + "unsupported value")
		}
	}
	return nil
}

func (v ValBytea) Value() (driver.Value, error) {	
	if v.IsNull {
		return driver.Value(nil),nil
	}
	return driver.Value(v.TypedValue), nil
}

func Base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func Base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}

func NewValBytea(val []byte, isNull bool) ValBytea{
	return ValBytea{Val{true, isNull}, val}
}
