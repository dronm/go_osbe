package fields

import (
	"strings"
	"database/sql/driver"
)

const (
	QUOTE_CHAR byte = 34
	JSON_NULL = "null"
)

type ValExt interface {
	GetIsSet() bool
	Value() (driver.Value, error)
}

//Base data type
type Val struct {
	IsSet bool
	IsNull bool
}

func ExtValIsNull(extVal []byte) bool {
	if len(extVal) == len(`"null"`) && (string(extVal) == `"null"` || strings.ToUpper(string(extVal)) == `"NULL"`) {
			return true
	}
	if len(extVal) == len(`null`) && (string(extVal) == `null` || strings.ToUpper(string(extVal)) == `NULL`) {
		return true
	}
	
	return false
}

func ExtRemoveQuotes(extVal []byte) string {
	var v_str string
	if extVal[0] == QUOTE_CHAR && byte(extVal[len(extVal)-1]) ==  QUOTE_CHAR {
		v_str = string(extVal[1:len(extVal)-1])
	}else {
		v_str = string(extVal)
	}
	return v_str
}
