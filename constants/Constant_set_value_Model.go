package constants

/**
 * Andrey Mikhalevich 16/12/21
 * This file is part of the OSBE framework
 */

//Controller method model
import (
	"osbe/fields"
	"encoding/json"
	"bytes"
	"errors"
//	"fmt"
)

const (
	CURLY_BR_OPEN_CHAR byte = 123
	CURLY_BR_CLOSE_CHAR byte = 125
)

type ValConst struct {
	fields.ValText
}
func (v *ValConst) UnmarshalJSON(data []byte) error {
	v.IsSet = true
	
	if fields.ExtValIsNull(data){
		v.IsNull = true
		return nil
	}
//json constant raises error, http always object!!!	
	var temp string
	if data[0] == CURLY_BR_OPEN_CHAR && byte(data[len(data)-1]) ==  CURLY_BR_CLOSE_CHAR {
		//object given
		//serialize all
		//!!!NO data = bytes.ReplaceAll(data, []byte(`\"`), []byte(`"`) )
		data = bytes.ReplaceAll(data, []byte(`"`), []byte(`\"`) )
		data = append(data[:1], data[0:]...)
		data[0] = fields.QUOTE_CHAR
		data = append(data, fields.QUOTE_CHAR)
		
//fmt.Println("data=", string(data))		
		//data[0] = QUOTE_CHAR
		//data[len(data)-1] = QUOTE_CHAR
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return errors.New(fields.ER_UNMARSHAL_STRING + err.Error())
	}
	
	v.TypedValue = temp
	
	return nil	
}

type Constant_set_value_argv struct {
	Argv *Constant_set_value `json:"argv"`	
}

//
type Constant_set_value struct {
	Id fields.ValText `json:"id" required:"true"`
	Val ValConst `json:"val"`
}

