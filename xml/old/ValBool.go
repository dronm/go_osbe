package xml

import (
	"encoding/xml"
	
	"osbe/fields"
)
type ValBool struct {
	fields.ValBool
}

func (v *ValBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	TextFieldToXML(v, e, start)
	return nil
}


