package fields

import (
	//"encoding/utf8"
	//"errors"
	//"fmt"
)

//***** Metadata text field:strings/texts ******************
type FieldTime struct {
	Field
	Time ParamBool
	Date ParamBool
	TZ ParamBool
}
func (f *FieldTime) GetDataType() FieldDataType {
	if f.Date.Value && f.Time.Value && f.TZ.Value {
		return FIELD_TYPE_DATETIMETZ
		
	}else if f.Date.Value && f.Time.Value {
		return FIELD_TYPE_DATETIME
		
	}else if f.Date.Value {
		return FIELD_TYPE_DATE
		
	}else if f.Time.Value {
		return FIELD_TYPE_TIME		
	}
	return FIELD_TYPE_DATETIMETZ
}
func (f *FieldTime) GetTime() ParamBool {
	return f.Time
}

func (f *FieldTime) SetTime(v ParamBool) {
	f.Time = v
}

func (f *FieldTime) GetDate() ParamBool {
	return f.Date
}

func (f *FieldTime) SetDate(v ParamBool) {
	f.Date = v
}

func (f *FieldTime) GetTZ() ParamBool {
	return f.TZ
}

func (f *FieldTime) SetTZ(v ParamBool) {
	f.TZ = v
}

type FielderTime interface {
	Fielder
	GetTime() ParamBool
	SetTime(ParamBool)
	GetDate() ParamBool
	SetDate(ParamBool)
	GetTZ() ParamBool
	SetTZ(ParamBool)
	
}

//String validaion
func ValidateTime(f FielderTime, val string) error {
	//ToDo

	return nil
}


