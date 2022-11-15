package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValDateTime struct {
	fields.ValDateTime
}

func (v *ValDateTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


