package fields

import (
	"errors"
	"fmt"
)

const CONSTR_VALUE_SEP = ","

type FieldEnum struct {
	FieldText
	Values []string
}

func (f *FieldEnum) GetDataType() FieldDataType {
	return FIELD_TYPE_ENUM
}

func (f *FieldEnum) CheckValue(val string) bool{
	if f.Values == nil {
		return false
	}
	for _,v := range f.Values {
		if v == val {
			return true
		}
	}
	return false
}

//values - comma separated values
//func NewFieldEnum(values string) *FieldEnum {
//	return &FieldEnum{Values: strings.Split(values, CONSTR_VALUE_SEP)}
//}

type FielderEnum interface {
	FielderText
	CheckValue(string) bool
	//GetDescription(string, string) string	
}

//Enum validaion
func ValidateEnum(f FielderEnum, val string) error {
	if !f.CheckValue(val) {
		return errors.New(fmt.Sprintf(ER_VALID_ENUM, f.GetDescr()) )
	}
	return nil
}


